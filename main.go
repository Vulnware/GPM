package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
)

func main() {

	// check program is run by go run or go build
	isGoRun := false
	for _, arg := range os.Args {
		if arg == "run" {
			isGoRun = true
		}
	}
	fmt.Println("isGoRun", isGoRun)

	if len(os.Args) == 2 {
		if os.Args[1] == "daemon" {
			SaveDaemonPid()
			// set logs to file
			// go
			logFile, err := os.OpenFile(GetGpmLogFile(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				log.Fatalf("Failed to open log file: %s", err)
			}
			defer logFile.Close()

			// Redirect stdout and stderr to the log file
			log.SetOutput(logFile)
			log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)
			defer DeleteDaemonPid()
			defer logFile.Close()
			StartIpcServer()

			return
		}
	}
	/*
		This part will be moved to cli project
	*/
	// create new parser object
	parser := argparse.NewParser("GPM", "A simple process manager made in Go")
	// start command which get command
	startCmd := parser.NewCommand("start", "Start a service")
	startFileArg := startCmd.StringPositionalWithName("file", &argparse.Options{Required: true, Help: "The file to run"})

	startName := startCmd.String("", "name", &argparse.Options{Required: false, Help: "The name of the service"})
	startWatch := startCmd.Flag("", "watch", &argparse.Options{Required: false, Help: "Watch the service for changes"})
	startForce := startCmd.Flag("", "force", &argparse.Options{Required: false, Help: "Force the service to start"})
	// stop command which get command
	stopCmd := parser.NewCommand("stop", "Stop a service")
	stopName := stopCmd.StringPositionalWithName("name", &argparse.Options{Required: true, Help: "The name of the service"})
	// restart command which get command
	restartCmd := parser.NewCommand("restart", "Restart a service")
	// status command which get command
	statusCmd := parser.NewCommand("status", "Get the status of a service")
	// list command which get command
	listCmd := parser.NewCommand("list", "List all services")
	// logs command which get command
	logsCmd := parser.NewCommand("logs", "Get the logs of a service")
	// config command which get command
	configCmd := parser.NewCommand("config", "Get the config of a service")
	// test command which get command
	testCmd := parser.NewCommand("test", "not for production")
	_ = testCmd // remove unused warning for testCmd variable
	// x := testCmd.String("m", "memory", &argparse.Options{Required: false, Help: "memory dump"})
	// parse command arguments

	// start command arguments

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if !DoesDaemonWork() {
		StartDaemon()
	}
	Client := ConnectToIpcServer()
	if startCmd.Happened() {

		fmt.Println("start command")
		fmt.Println("command: ", *startFileArg)
		fmt.Println("name: ", *startName)
		fmt.Println("watch: ", *startWatch)
		fmt.Println("force: ", *startForce)
		// get working directory
		// if *path == "" {
		// 	*path = "."
		// }
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			return
		}
		if *startName == "" {
			*startName = *startFileArg
		}
		var reply Reply
		// send the command to the server
		Client.Call("Server.Start", &Service{Name: *startName, Path: cwd, Command: *startFileArg}, &reply)

		log.Println(reply.Message)

	}

	if stopCmd.Happened() {

		fmt.Println("name: ", *stopName)
		var reply Reply
		// send the command to the server
		Client.Call("Server.Stop", &Service{Name: *stopName}, &reply)
		log.Println(reply.Message)

	}
	if restartCmd.Happened() {
		fmt.Println("restart command")
	}
	if statusCmd.Happened() {
		fmt.Println("status command")
	}
	if listCmd.Happened() {
		fmt.Println("list command")
	}
	if logsCmd.Happened() {
		fmt.Println("logs command")
	}
	if configCmd.Happened() {
		fmt.Println("config command")
	}
}
