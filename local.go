package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"
)

func GetGpmPath() string {
	Homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path := Homedir + "/.gpm"
	// check if path exists and create if not
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	return path
}
func GetGpmLogFile() string {
	return GetGpmPath() + "/gpm.log"
}

func GetPidFolder() {
	// check if pid folder exists and create if not
	if _, err := os.Stat(GetGpmPath() + "/pids"); os.IsNotExist(err) {
		err := os.Mkdir(GetGpmPath()+"/pids", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetLogFolder() {
	// check if log folder exists and create if not
	if _, err := os.Stat(GetGpmPath() + "/logs"); os.IsNotExist(err) {
		err := os.Mkdir(GetGpmPath()+"/logs", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func GetProcessFolder() {
	// check if process folder exists and create if not
	if _, err := os.Stat(GetGpmPath() + "/processes"); os.IsNotExist(err) {
		err := os.Mkdir(GetGpmPath()+"/processes", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func SavePid(pid int, name string) {
	GetPidFolder()
	// save pid to file
	file, err := os.Create(GetGpmPath() + "/pids/" + name + ".pid")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("%d", pid))

	pids[name] = pid
}

func SaveProcess(process string, name string, t TOML) {
	GetProcessFolder()
	// save process to file .toml
	// check if file exists and create if not
	if _, err := os.Stat(GetGpmPath() + "/processes/" + name + ".toml"); os.IsNotExist(err) {
		file, err := os.Create(GetGpmPath() + "/processes/" + name + ".toml")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		// if file exists, delete it
		err := os.Remove(GetGpmPath() + "/processes/" + name + ".toml")
		if err != nil {
			log.Fatal(err)
		}
		// create new file
		file, err := os.Create(GetGpmPath() + "/processes/" + name + ".toml")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
	// save process to file
	file, err := os.OpenFile(GetGpmPath()+"/processes/"+name+".toml", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	toml.NewEncoder(file).Encode(t)

}

func GetLogFile(name string) *os.File {
	GetLogFolder()
	// check if log file exists and create if not
	if _, err := os.Stat(GetGpmPath() + "/logs/" + name + ".log"); os.IsNotExist(err) {
		file, err := os.Create(GetGpmPath() + "/logs/" + name + ".log")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		return file
	}

	// if file exists, open it with read and write permissions

	file, err := os.OpenFile(GetGpmPath()+"/logs/"+name+".log", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file

}

func GetLogFileError(name string) *os.File {
	GetLogFolder()
	// check if log file exists and create if not
	if _, err := os.Stat(GetGpmPath() + "/logs/" + name + ".err"); os.IsNotExist(err) {
		file, err := os.Create(GetGpmPath() + "/logs/" + name + ".err")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		return file
	}

	// if file exists, open it with read and write permissions
	file, err := os.OpenFile(GetGpmPath()+"/logs/"+name+".err", os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	return file
}

func GetPid(name string) int {
	// check if pid file exists
	if _, err := os.Stat(GetGpmPath() + "/pids/" + name + ".pid"); os.IsNotExist(err) {
		log.Fatal("Process does not exist")
	}
	// read pid from file
	file, err := os.Open(GetGpmPath() + "/pids/" + name + ".pid")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var pid int
	_, err = fmt.Fscanf(file, "%d", &pid)
	if err != nil {
		log.Fatal(err)
	}
	return pid
}

func GetProcess(name string) TOML {
	// check if process file exists
	if _, err := os.Stat(GetGpmPath() + "/processes/" + name + ".toml"); os.IsNotExist(err) {
		log.Fatal("Process does not exist")
	}
	// read process from file
	file, err := os.Open(GetGpmPath() + "/processes/" + name + ".toml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var t TOML
	_, err = toml.NewDecoder(file).Decode(&t)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func DeletePid(name string) {
	// check if pid file exists
	if _, err := os.Stat(GetGpmPath() + "/pids/" + name + ".pid"); os.IsNotExist(err) {
		log.Fatal("Process does not exist")
	}
	// delete pid file
	err := os.Remove(GetGpmPath() + "/pids/" + name + ".pid")
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteProcess(name string) {
	// check if process file exists
	if _, err := os.Stat(GetGpmPath() + "/processes/" + name + ".toml"); os.IsNotExist(err) {
		log.Fatal("Process does not exist")
	}
	// delete process file
	err := os.Remove(GetGpmPath() + "/processes/" + name + ".toml")
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteLog(name string) {
	// check if log file exists
	if _, err := os.Stat(GetGpmPath() + "/logs/" + name + ".log"); os.IsNotExist(err) {
		log.Fatal("Process does not exist")
	}
	// delete log file
	err := os.Remove(GetGpmPath() + "/logs/" + name + ".log")
	if err != nil {
		log.Fatal(err)
	}
}

func GetDaemonPid() int {
	// check if pid file exists
	if _, err := os.Stat(GetGpmPath() + "/daemon.pid"); os.IsNotExist(err) {
		log.Fatal("Daemon does not exist")
	}
	// read pid from file
	file, err := os.Open(GetGpmPath() + "/daemon.pid")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var pid int
	_, err = fmt.Fscanf(file, "%d", &pid)
	if err != nil {
		log.Fatal(err)
	}
	return pid

}

func DoesDaemonWork() bool {
	// check if daemon is running
	if _, err := os.Stat(GetGpmPath() + "/daemon.pid"); os.IsNotExist(err) {
		return false
	}
	// find process by pid and check if it is daemon
	pid := GetDaemonPid()
	_, err := os.FindProcess(pid)

	return err == nil

}
func StartDaemon() {
	// check if daemon is running
	if DoesDaemonWork() {
		log.Fatal("Daemon is already running")
	}
	// start daemon
	cmd := exec.Command("gpm", "daemon")
	err := cmd.Start()

	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Process.Release() // release process from parent
	if err != nil {
		log.Fatal(err)
	}

	// save pid to file
}

func SaveDaemonPid() {
	// get own pid
	pid := os.Getpid()
	// check if file exists

	if _, err := os.Stat(GetGpmPath() + "/daemon.pid"); os.IsNotExist(err) {
		file, err := os.Create(GetGpmPath() + "/daemon.pid")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	} else {
		// if file exists, delete it
		err := os.Remove(GetGpmPath() + "/daemon.pid")
		if err != nil {
			log.Fatal(err)
		}
		// create new file
		file, err := os.Create(GetGpmPath() + "/daemon.pid")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
	}
	// save pid to file
	file, err := os.OpenFile(GetGpmPath()+"/daemon.pid", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = fmt.Fprintln(file, pid)
	if err != nil {
		log.Fatal(err)
	}
}

func DeleteDaemonPid() {
	// check if pid file exists
	if _, err := os.Stat(GetGpmPath() + "/daemon.pid"); os.IsNotExist(err) {
		log.Fatal("Daemon does not exist")
	}
	// delete pid file
	err := os.Remove(GetGpmPath() + "/daemon.pid")
	if err != nil {
		log.Fatal(err)
	}
}
