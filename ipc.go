package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net"
	"net/rpc"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"gopkg.in/natefinch/npipe.v2"
)

type Args struct {
	A, B int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

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

type Server struct{}

func (t *Server) Start(args *StartArgs, reply *Reply) error {
	// create a channel to wait for the process to start
	c := make(chan bool, 1)

	// start the process
	go startProject(Service{Name: args.Name, Path: args.Path, Command: args.Command}, &globalEnv, reply, c)

	<-c
	// return success
	return nil
}

func (t *Server) Stop(args *StartArgs, reply *Reply) error {
	// stop the process

	stopProject(args.Name, reply)
	// return success
	return nil

}
func (t *Server) Dump(args *StartArgs, reply *Reply) error {
	// dump memory of the process

	dumpProject(args.Name, reply)

	// return success
	return nil
}

func (t *Server) Restart(args *StartArgs, reply *Reply) error {
	// restart the process
	c := make(chan bool, 1)
	restartProject(args.Name, reply, c)
	<-c // wait for the process to restart before returning
	return nil
}
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

func StartIpcServer() {
	var listener net.Listener
	var err error

	// Determine the operating system
	if runtime.GOOS == "windows" {
		// Use npipe package to create named pipe on Windows
		listener, err = npipe.Listen(pipeName)
		if err != nil {
			fmt.Println("Error creating named pipe:", err)
			os.Exit(1)
		}
		fmt.Println("Listening on named pipe:", GeneratePipeName())
		// send a message to the client

	} else {
		// Use net package to create Unix domain socket on Unix
		os.Remove(GeneratePipeName())
		listener, err = net.Listen("unix", GeneratePipeName())
		if err != nil {
			fmt.Println("Error listening:", err.Error())
			os.Exit(1)
		}
		fmt.Println("Listening on Unix domain socket:", GeneratePipeName())
	}
	defer listener.Close()

	// Register the RPC server
	rpc.Register(new(Server))

	// Accept incoming connections
	for {
		conn, err := listener.Accept()
		// send a message to the client
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		// read the message

		go rpc.ServeConn(conn)
	}
}

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
