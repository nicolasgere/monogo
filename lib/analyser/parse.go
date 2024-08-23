package analyzer

import (
	"encoding/json"
	"fmt"
	"github.com/dominikbraun/graph"
	"os/exec"
	"strings"
)

func ListModule(dir string) (modules []Module, err error) {
	output, err := runCommand(dir, "go list -m -json")
	if err != nil {
		return
	}

	d := json.NewDecoder(strings.NewReader(output))
	for d.More() {
		var m Module
		if err = d.Decode(&m); err != nil {
			return
		}
		modules = append(modules, m)
	}
	return
}

func GetDependencyForModule(dir string) (details []ModuleDetail, err error) {
	output, err := runCommand(dir, "go list -json ./...")
	if err != nil {
		return
	}

	d := json.NewDecoder(strings.NewReader(output))
	for d.More() {
		var detail ModuleDetail
		if err = d.Decode(&detail); err != nil {
			return
		}
		details = append(details, detail)
	}
	return
}

func GetImportFromModules(details []ModuleDetail) (imports map[string]bool) {
	imports = make(map[string]bool)
	for _, m := range details {
		for _, i := range m.Imports {
			imports[ strings.Split(i, "/")[0]] = true
			
		}
	}
	return imports
}

func BuildDependencyGraph(modules []Module) (gr *graph.Graph[string, string], err error) {
	g := graph.New(graph.StringHash, graph.Directed(), graph.Acyclic())
	moduleMap := make(map[string]*Module)

	for i := range modules {
		moduleMap[modules[i].Path] = &modules[i]
	}

	for _, m := range modules {
		var details []ModuleDetail
		details, err = GetDependencyForModule(m.Dir)
		if err != nil {
			return nil, err
		}
		moduleMap[m.Path].Imports = GetImportFromModules(details)

		if err = g.AddVertex(m.Path); err != nil {
			// return nil, fmt.Errorf("failed to add vertex %s: %w", m.Path, err)
		}
	}
	
	for _, m := range moduleMap {
		for	i := range m.Imports{
			_, ok := moduleMap[i]
			if ok {
				fmt.Println("ADD EDGE", i, m.Path)
				if err = g.AddEdge(m.Path,i ); err != nil {
					//				// return nil, fmt.Errorf("failed to add edge %s -> %s: %w", m.Path, importPath, err)
					//			}
				}
			}
	
		}
	}
	gr = &g
	return
}

func GetDependencyPaths(g *graph.Graph[string, string], vertex string) ([]string, error) {
	var dependencyPaths []string

	// Define a visitor function for DFS
	visitor := func(v string) bool {
		if v != vertex { // Don't include the starting vertex itself
			dependencyPaths = append(dependencyPaths, v)
		}
		return false // Continue traversal
	}

	// Perform DFS
	err := graph.DFS(*g, vertex, visitor)
	if err != nil {
		return nil, fmt.Errorf("failed to perform DFS for vertex %s: %w", vertex, err)
	}

	return dependencyPaths, nil
}

func runCommand(dir, command string) (output string, err error) {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = dir
	var outputBytes []byte
	outputBytes, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w\nOutput: %s", err, outputBytes)
	}
	return string(outputBytes), nil
}
