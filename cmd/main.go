package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	analyzer "mono/lib/analyser"
	"mono/lib/git"
	"mono/lib/runner"
	"mono/lib/utils"

	"github.com/urfave/cli/v2"
)

var defaultDir = "."

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandling(cancel)

	r := runner.NewRunner(ctx, 3)
	app := createCliApp(&r)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func setupSignalHandling(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()
}

func createCliApp(r *runner.Runner) *cli.App {
	return &cli.App{
		Commands: []*cli.Command{
			createCommand("install", "Install dependency for every modules", "go mod download -x", r),
			createCommand("fmt", "Format every modules", "go fmt ./...", r),
			createCommand("test", "Test every modules", "go test ./...", r),
		},
	}
}

func createCommand(name, usage, cmd string, r *runner.Runner) *cli.Command {
	var target string
	var dependency bool
	var branch string

	return &cli.Command{
		Name:  name,
		Usage: usage,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "Path",
				Usage:       "Path to the root directory of the project",
				Aliases:     []string{"p"},
				Destination: &defaultDir,
			},
			&cli.StringFlag{
				Name:        "target",
				Usage:       "Targuetted module",
				Aliases:     []string{"t"},
				Destination: &target,
			},
			&cli.BoolFlag{
				Name:        "dependency",
				Usage:       "Run with all dependency of the targuet. Descendent and ascendent",
				Aliases:     []string{"d"},
				Destination: &dependency,
				Value:       false,
			},
			&cli.StringFlag{
				Name:        "branch",
				Usage:       "Compare the current branch with the master branch, and found affected modules",
				Aliases:     []string{"b"},
				Destination: &branch,
			},
		},
		Action: func(*cli.Context) error {
			modules, err := analyzer.ListModule(defaultDir)
			if err != nil {
				return err
			}
			allModules := make(map[string]analyzer.Module)
			for _, module := range modules {
				allModules[module.Path] = module
			}
			modulesToRun := modules
			if branch != "" {
				gitModules, err := git.GetAffectedRootDirectories("master", defaultDir)
				if err != nil {
					return err
				}
				filteredModule := make([]analyzer.Module, 0)
				// Filter modules from git
				for _, m := range modules {
					for _, gm := range gitModules {
						if m.Path == gm {
							filteredModule = append(filteredModule, m)
						}
					}

				}
				modulesToRun = filteredModule
			}

			if target != "" {
				filteredModule := make([]analyzer.Module, 0)
				for _, m := range modules {
					if m.Path == target {
						filteredModule = append(filteredModule, m)
					}
				}
				modulesToRun = filteredModule
				// Add dependency
				if dependency {
					dependencyModules := map[string]string{}
					g, err := analyzer.BuildDependencyGraph(modules)
					if err != nil {
						return err
					}
					for _, m := range modulesToRun {
						paths, err := analyzer.GetDependencyPaths(g, m.Path)
						if err != nil {
							return err
						}
						for _, p := range paths {
							dependencyModules[p] = p
						}
					}
					existingMap := make(map[string]bool)
					for _, module := range modulesToRun {
						existingMap[module.Path] = true
					}

					// Add new unique modules
					for path := range dependencyModules {
						if !existingMap[path] {
							modulesToRun = append(modulesToRun, allModules[path])
							existingMap[path] = true
						}
					}
				}
			}
			runOnModules(defaultDir, cmd, r, modulesToRun)
			return nil
		},
	}
}

func runOnModules(dir, cmd string, r *runner.Runner, modules []analyzer.Module) {

	tasks := createTasks(modules, cmd)
	tfs := r.RunTasks(tasks)

	var wg sync.WaitGroup
	wg.Add(len(tfs))

	for _, tf := range tfs {
		go handleTaskFuture(tf, &wg)
	}

	wg.Wait()
	return
}

func createTasks(modules []analyzer.Module, cmd string) []runner.Task {
	tasks := make([]runner.Task, len(modules))
	for i, module := range modules {
		tasks[i] = runner.Task{
			Id:   module.Path,
			Cmd:  cmd,
			Root: module.Dir,
		}
	}
	return tasks
}

func handleTaskFuture(tf *runner.TaskFuture, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case stdout, ok := <-tf.Stdout:
			handleOutput(tf.Id, stdout, ok, &tf.Stdout)
		case stderr, ok := <-tf.Stderr:
			handleOutput(tf.Id, stderr, ok, &tf.Stderr)
		case result := <-tf.Done:
			utils.LogWithTaskId(tf.Id, fmt.Sprintf("Done with status %d", result.Status), utils.INFO)
			return
		}
	}
}

func handleOutput(id string, output []byte, ok bool, channel *chan []byte) {
	if !ok {
		*channel = nil
		return
	}
	if len(output) > 0 {
		utils.LogWithTaskId(id, string(output), utils.INFO)
	}
}
