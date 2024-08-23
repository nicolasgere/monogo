package runner

import (
	"fmt"
	"sync"
	"testing"
)

// func TestRunner(t *testing.T) {
// 	r := Runner{}

// 	task := Task{
// 		Cmd:  "cat hello.txt",
// 		Root: "./__playground__/a/",
// 	}
// 	fmt.Println("START")
// 	fmt.Println("running task")
// 	tf := r.RunTask(task)

// 	for {
// 		select {
// 		case stdout, ok := <-tf.Stdout:
// 			if !ok {
// 				tf.Stdout = nil
// 			}
// 			fmt.Println("stdout", string(stdout))
// 			break
// 		case stderr, ok := <-tf.Stderr:
// 			if !ok {
// 				tf.Stderr = nil
// 			}
// 			fmt.Println("stderr", string(stderr))
// 			break
// 		case result := <-tf.Done:
// 			fmt.Println("done")
// 			fmt.Println(result.Status)
// 			return
// 		}
// 	}
// }

func TestRunnerMultiple(t *testing.T) {
	r := Runner{}
	tasks := []Task{
		{
			Id:   "a",
			Cmd:  "cat hello.txt",
			Root: "./__playground__/a/",
		},
		{
			Id:   "b",
			Cmd:  "cat world.txt",
			Root: "./__playground__/b/",
		},
	}
	tfs := r.RunTasks(tasks)
	var wg sync.WaitGroup
	for _, tf := range tfs {
		wg.Add(1)
		go func(tf *TaskFuture) {
			for {
				select {
				case stdout, ok := <-tf.Stdout:
					if !ok {
						tf.Stdout = nil
					}
					if len(stdout) == 0 {
						break
					}
					fmt.Println("stdout", tf.Id, string(stdout))
					break
				case stderr, ok := <-tf.Stderr:
					if !ok {
						tf.Stderr = nil
					}
					if len(stderr) == 0 {
						break
					}
					fmt.Println("stderr", tf.Id, string(stderr))
					break
				case result := <-tf.Done:
					fmt.Println("done", tf.Id)
					fmt.Println(result.Status)
					wg.Done()
					break
				}
			}
		}(tf)
	}
	wg.Wait()
}
