package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"

	"github.com/Vulnware/gpm/memory"
)

func parseCommand(command string) []string {
	var parsedCommand []string

	// split command into array of strings by space
	parsedCommand = strings.Split(command, " ")
	return parsedCommand
}
func startProject(p project, gEnv *[]string, reply *Reply, channel chan bool) {
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
	signal.Notify(c, os.Interrupt, os.Kill)
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
	defer DeletePid(p.Name)
	defer LogFile.Close()
	defer LogErrFile.Close()
	defer delete(pids, p.Name)
	defer cmd.Process.Kill()

	*reply = Reply{Success: true, Message: "Project started"}
	channel <- true
	if err := cmd.Wait(); err != nil {
		log.Printf("%s command wait error: %s\n", p.Name, err)
	}

	// stdout, err := cmd.StdoutPipe()
	// if err != nil {
	// 	log.Printf("%s stdout pipe error: %s\n", p.Name, err)
	// 	return
	// }
	// // create a io.Reader to collect the output from the command (stderr)
	// var output bytes.Buffer

	// cmd.Stderr = &output
	// if CONFIGS["outputMode"] == "stdout" {
	// 	cmd.Stdout = os.Stdout
	// }
	// scanner := bufio.NewScanner(stdout)

	// if err := cmd.Start(); err != nil {
	// 	log.Printf("%s command start error: %s\n", p.Name, err)
	// 	return
	// }

	// defer func() {
	// 	if err := cmd.Wait(); err != nil {
	// 		log.Printf("%s command wait error: %s\n", p.Name, err)
	// 		// log the error from stderr if it exists
	// 		if output.Len() > 0 {
	// 			log.Printf("%s stderr: %s\n", p.Name, output.String())
	// 		}

	// 	}
	// }()

	// if CONFIGS["outputMode"] == "web" {
	// 	for scanner.Scan() {
	// 		output := scanner.Text()
	// 		// if err := decoder.Decode(&output); err != nil {
	// 		// 	log.Printf("%s stdout read error: %s\n", p.Name, err)
	// 		// 	break
	// 		// }

	// 		data := map[string]interface{}{
	// 			"method":  "service",
	// 			"service": p.Name,
	// 			"output":  output,
	// 		}
	// 		appendLogs(p.Name, output)
	// 		jsonData, err := json.Marshal(data)
	// 		if err != nil {
	// 			log.Printf("%s json marshal error: %s\n", p.Name, err)
	// 			break
	// 		}

	// 		// Broadcast message to all clients
	// 		for c := range clients {
	// 			mut.Lock()
	// 			if err := c.WriteMessage(websocket.TextMessage, jsonData); err != nil {
	// 				log.Printf("%s write error: %s\n", p.Name, err)
	// 				mut.Unlock()
	// 				continue
	// 			}
	// 			mut.Unlock()
	// 		}
	// 	}
	// 	if output.Len() > 0 {
	// 		for c := range clients {
	// 			mut.Lock()
	// 			if err := c.WriteMessage(websocket.TextMessage, output.Bytes()); err != nil {
	// 				log.Printf("%s write error: %s\n", p.Name, err)
	// 				mut.Unlock()
	// 				continue
	// 			}
	// 			mut.Unlock()
	// 		}
	// 	}
	// }
}

func stopProject(name string, reply *Reply) {
	pid := pids[name]
	if pid == 0 {
		log.Printf("Error: %s is not running\n", name)
		*reply = Reply{Success: false, Message: "Project is not running"}
		return
	}
	log.Printf("Stopping %s with pid %d\n", name, pid)
	p, err := os.FindProcess(pid)
	if err != nil {
		log.Printf("Error: %s\n", err)
		*reply = Reply{Success: false, Message: "Project is not running"}
		return
	}
	if err := p.Kill(); err != nil {
		log.Printf("Error: %s\n", err)
		*reply = Reply{Success: false, Message: "Project is not running"}
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
