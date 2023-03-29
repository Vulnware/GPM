// build +windows
package memory

import (
	"fmt"
	"os"
	"syscall"
)

var (
	kernel32          = syscall.NewLazyDLL("kernel32.dll")
	openProcess       = kernel32.NewProc("OpenProcess")
	readProcessMemory = kernel32.NewProc("ReadProcessMemory")
	closeHandle       = kernel32.NewProc("CloseHandle")
	dbghelp           = syscall.NewLazyDLL("dbghelp.dll")
	miniDumpWriteDump = dbghelp.NewProc("MiniDumpWriteDump")
)

const PROCESS_ALL_ACCESS = 0x1F0FFF
const MINIDUMP_TYPE_FULL = 0x00000003

func DumpMemory(pid int, file string) {
	processHandle, err := syscall.OpenProcess(PROCESS_ALL_ACCESS, false, uint32(pid))
	if err != nil {
		fmt.Println("Error opening process:", err)
		os.Exit(1)
	}
	defer syscall.CloseHandle(processHandle)

	dumpFile, err := os.Create(file)
	if err != nil {
		fmt.Println("Error creating dump file:", err)
		os.Exit(1)
	}
	defer dumpFile.Close()

	_, _, _ = miniDumpWriteDump.Call(
		uintptr(processHandle),
		uintptr(pid),
		uintptr(dumpFile.Fd()),
		MINIDUMP_TYPE_FULL,
		0,
		0,
		0,
	)

}
