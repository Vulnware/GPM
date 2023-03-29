package main

type StartArgs struct {
	Name    string
	Path    string
	Command string
}

type StartReply struct {
	Success bool
	Message string
}
type Reply struct {
	Success bool
	Message string
}

type Service struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Command string `json:"command"`
	// Env is a map of environment iables to set for the command. which can be nil.
	Env map[string]interface{} `json:"env"` // optional environment variables to set
}
