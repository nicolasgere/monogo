package runner

type Task struct {
	Id   string
	Name string
	Root string
	Cmd  string
	Args []string
}

type TaskFuture struct {
	Stdout chan []byte
	Stderr chan []byte
	Done   chan TaskResult
	Id     string
}

type TaskResult struct {
	Err    error
	Status int
}
