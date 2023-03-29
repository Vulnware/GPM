package main

/*
This project will be moved to a new repo soon and will be renamed to gpm-cli
*/
import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"os"
	"runtime"

	"github.com/akamensky/argparse"
	"github.com/shirou/gopsutil/cpu"
	"gopkg.in/natefinch/npipe.v2"
)

func GeneratePipeName() string {
	// Get cpu info and mac address
	cpuInfo, _ := cpu.Info()
	macAddress, _ := net.Interfaces()
	// Generate a hash from the cpu info and mac address
	hash := sha256.New()
	hash.Write([]byte(cpuInfo[0].ModelName + macAddress[0].HardwareAddr.String()))
	// Return the hash as a string
	if runtime.GOOS == "windows" {
		return fmt.Sprintf(`\\.\pipe\%s.gpm.rpc`, hex.EncodeToString(hash.Sum(nil)))
	} else {
		return fmt.Sprintf(`/tmp/%s.gpm.rpc`, hex.EncodeToString(hash.Sum(nil)))
	}

}

var pipeName = GeneratePipeName()

func ConnectToIpcServer() *rpc.Client {
	var client *rpc.Client
	var err error

	// Determine the operating system
	if runtime.GOOS == "windows" {
		// Use npipe package to connect to named pipe on Windows
		c, err := npipe.Dial(GeneratePipeName())
		// send a message to the client

		if err != nil {
			fmt.Println("Error connecting to named pipe:", err)
			os.Exit(1)
		}
		client = rpc.NewClient(c)

		fmt.Println("Connected to named pipe:", GeneratePipeName())
	} else {
		// Use net package to connect to Unix domain socket on Unix
		client, err = rpc.Dial("unix", GeneratePipeName())
		if err != nil {
			fmt.Println("Error connecting to Unix domain socket:", err)
			os.Exit(1)
		}
		fmt.Println("Connected to Unix domain socket:", GeneratePipeName())
	}

	// Call the RPC method

	return client
}

func main() {
	// check program is run by go run or go build
	isGoRun := false
	for _, arg := range os.Args {
		if arg == "run" {
			isGoRun = true
		}
	}
	fmt.Println("isGoRun", isGoRun)

	// create new parser object
	parser := argparse.NewParser("GPM-CLI", "GPM-CLI is a command line interface for GPM")
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
	// parse command arguments

	// start command arguments

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
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
