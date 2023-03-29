//go:build !windows
// +build !windows

package memory

import (
	"fmt"
	"os"
)

func DumpMemory(pid int, file string) {
	// dump memory of the process using gdb

	// open the /proc/<pid>/mem file for reading
	memfile, err := os.Open(fmt.Sprintf("/proc/%d/mem", pid))
	if err != nil {
		fmt.Printf("Error opening mem file: %v\n", err)
		return
	}
	defer memfile.Close()

	// seek to the start of the memory region to dump
	offset := uintptr(0x10000) // example offset
	_, err = memfile.Seek(int64(offset), 0)
	if err != nil {
		fmt.Printf("Error seeking to offset %d: %v\n", offset, err)
		return
	}

	// read the memory region
	buf := make([]byte, 4096) // example buffer size
	n, err := memfile.Read(buf)
	if err != nil {
		fmt.Printf("Error reading memory: %v\n", err)
		return
	}
	// write the memory region to a file
	dumpfile, err := os.Create(file)
	if err != nil {
		fmt.Printf("Error creating dump file: %v\n", err)
		return
	}
	defer dumpfile.Close()
	_, err = dumpfile.Write(buf[:n])
	if err != nil {
		fmt.Printf("Error writing dump file: %v\n", err)
		return
	}

}
