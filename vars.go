package main

import "os/exec"

var (
	pids      map[string]int       = make(map[string]int)
	cmds      map[string]*exec.Cmd = make(map[string]*exec.Cmd) // map of project names to commands (for killing processes)
	globalEnv []string

	tomlVersion     = "0.1.0"
	tomlVersionNums = parseVersion(tomlVersion)
	Services        []Service
)
