# GPM - Go Processes Manager

GPM is a command-line tool written in Go that helps to manage processes on your server or local machine. It consists of two modules - `gpm.exe` which runs as a daemon in the background, and `gpm-cli` which provides a command-line interface for managing processes. With GPM, you can start, stop, restart, and monitor processes on your system.

## Installation

### Prerequisites

- Go 1.16 or higher installed on your system.

### Installation Steps

1. Clone the repository from GitHub:

     ```bash
     git clone https://github.com/Vulnware/GPM.git
     ```

2. Build the binary for `gpm.exe` using the following command:
    For windows:

     ```bash
     make.bat
     ```

    For linux:

     ```bash
        make
    ```

3. Build the binary for `gpm-cli` using the following command:
    For windows:

     ```bash
     cd cli
     make.bat
     ```

    For linux:

     ```bash
        cd cli
        make
    ```

4. Move both binaries to the `/usr/local/bin` directory so that they can be executed from anywhere in the terminal for linux or `C:\Windows` for windows (or any other directory that is in your `PATH` environment variable). You can also add the directory where the binaries are located to your `PATH` environment variable.
    For linux:

    ```bash
    mv gpm.exe /usr/local/bin
    mv gpm-cli /usr/local/bin
    ```

    For windows:

    ```bash
    mv gpm.exe C:\Windows
    mv gpm-cli C:\Windows
    ```

    Or add the directory where the binaries are located to your `PATH` environment variable.

5. Verify that GPM is installed correctly by running the following command:

    ```bash
    gpm-cli --version
    ```

    You should see the version number of GPM printed to the console.

## Usage

### Starting a process

To start a process with GPM, run the following command:

```gpm-cli start <command> <args>```

For example, to start a Node.js server with GPM, you could use the following command:

````gpm-cli start "node server.js" --name "node-server"````

### Stopping a process

To stop a process with GPM, run the following command:

````gpm-cli stop <process-name>````

For example, to stop the Node.js server that we started earlier, you could use the following command:

````gpm-cli stop node-server````

### Restarting a process (Not implemented yet)

To restart a process with GPM, run the following command:

````gpm-cli restart <process-name>````

For example, to restart the Node.js server that we started earlier, you could use the following command:

````gpm-cli restart node-server````

### Listing processes (Not implemented yet)

To list all the processes that are currently running with GPM, run the following command:

````gpm-cli ls````

### Viewing process logs (Not implemented yet)

To view the logs for a process that is running with GPM, run the following command:

````gpm-cli logs <process-name>````

For example, to view the logs for the Node.js server that we started earlier, you could use the following command:

````gpm-cli logs node-server````

## Contributing

If you would like to contribute to GPM, please open a pull request on GitHub. We welcome any contributions, including bug fixes, new features, and documentation improvements.

## License

GPM is licensed under the MIT license. See the `LICENSE` file for more information.
