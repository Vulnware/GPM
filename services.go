package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/Vulnware/gpm/memory"
)

func parseCommand(command string) []string {
	var parsedCommand []string

	// split command into array of strings by space
	parsedCommand = strings.Split(command, " ")
	return parsedCommand
}
func startProject(p Service, gEnv *[]string, reply *Reply, channel chan bool) {
	if pids[p.Name] != 0 {
		log.Printf("Project %s is already running\n", p.Name)
		*reply = Reply{Success: false, Message: "Project is already running"}
		channel <- true

		return
	}
	log.Printf("Starting %s in %s with command '%s'\n", p.Name, p.Path, p.Command)
	var cmd *exec.Cmd
	// check if OS is windows
	// parsedCommand := parseCommand(p.Command)

	cmd = exec.Command("cmd", "/C", p.Command)
	// set group id of process to be the same as the parent process

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// cleanup
		cmd.Process.Kill()
		os.Exit(1)
	}()
	cmd.Dir = p.Path
	cmd.Env = os.Environ()
	if p.Env != nil {
		for k, v := range p.Env {
			switch v.(type) {
			case bool:
				globalEnv = append(globalEnv, k+"="+strconv.FormatBool(v.(bool)))
			case int:
				globalEnv = append(globalEnv, k+"="+strconv.Itoa(v.(int)))

			case float64:
				globalEnv = append(globalEnv, k+"="+strconv.FormatFloat(v.(float64), 'f', -1, 64))
			case string:
				globalEnv = append(globalEnv, k+"="+v.(string))
			case []interface{}:

				panic("Error: environment variables cannot be an array")
			case map[string]interface{}:
				panic("Error: environment variables cannot be a map")
			case int64:
				globalEnv = append(globalEnv, k+"="+strconv.FormatInt(v.(int64), 10))

			default:
				k = fmt.Sprintf("%v", k)
				vType := fmt.Sprintf("%T", v)
				panic(fmt.Sprintf("Error: environment variable %s has an unknown type %s", k, vType))

			}

		}
	}
	if gEnv != nil {
		cmd.Env = append(cmd.Env, *gEnv...)
	}

	LogFile := GetLogFile(p.Name)
	LogErrFile := GetLogFileError(p.Name)
	cmd.Stdout = LogFile
	cmd.Stderr = LogErrFile
	// limit memory usage to 100MB
	// cmd.SysProcAttr.CreationFlags = 0x08000000 // CREATE_NO_WINDOW flag

	if err := cmd.Start(); err != nil {
		log.Printf("%s command start error: %s\n", p.Name, err)
		*reply = Reply{Success: false, Message: "Project failed to start"}
		channel <- true
	}
	SavePid(cmd.Process.Pid, p.Name)
	defer func() {
		DeletePid(p.Name)
		LogFile.Close()
		LogErrFile.Close()
		delete(pids, p.Name)
		delete(cmds, p.Name)
		cmd.Process.Kill()
	}() // delete pid from map when function exits (deferred)

	*reply = Reply{Success: true, Message: "Project started"}
	channel <- true
	if err := cmd.Wait(); err != nil {
		log.Printf("%s command wait error: %s\n", p.Name, err)
	}

}

func stopProject(name string, reply *Reply) {
	cmd := cmds[name]
	if cmd == nil {
		log.Printf("Error: %s is not running\n", name)
		*reply = Reply{Success: false, Message: "Project is not running"}
		return
	}
	// kill the process
	if err := cmd.Process.Kill(); err != nil {
		log.Printf("Error: %s failed to stop\n", name)
		*reply = Reply{Success: false, Message: "Project failed to stop"}
		return
	}

	*reply = Reply{Success: true, Message: "Project stopped"}
}

func dumpProject(name string, reply *Reply) {
	// dump the memory of the process
	pid := pids[name]
	if pid == 0 {
		log.Printf("Error: %s is not running\n", name)
		*reply = Reply{Success: false, Message: "Project is not running"}
		return
	}
	memory.DumpMemory(pid, "dumpmemory")

	*reply = Reply{Success: true, Message: "Project memory dumped"}

}

func restartProject(name string, reply *Reply, channel chan bool) {
	stopProject(name, reply)
	if !reply.Success {
		return
	}
	// find service  in services
	for _, p := range Services {
		if p.Name == name {
			go startProject(p, nil, reply, channel)
			if <-channel {
				*reply = Reply{Success: true, Message: "Project restarted"}
			}
			return
		}
	}
}
