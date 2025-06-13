package services

type CmdConfig struct {
	Cmd     string
	Args    []string
	Dir     string
	Stdin   string
	Env     map[string]string
	Stdout  string
	Stderr  string
}

type CmdResult struct {
	isErr       bool
	errMsg      string
	ExitCode    int
	Stdout      string
	Stderr      string
}
