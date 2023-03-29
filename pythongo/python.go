package pythongo

// #cgo CFLAGS: -I"C:/msys64/mingw64/include/python3.10"
// #cgo LDFLAGS: -L"C:/msys64/mingw64/lib" -lpython3.10
/* #include <Python.h> */
import "C"
import (
	"unsafe"
)

func makeNameSpace() *C.PyObject {
	C.Py_InitializeEx(C.int(1))

	// Create a new module
	mainModule := C.PyImport_AddModule(C.CString("__main__"))
	// Get the module's dictionary
	mainNamespace := C.PyModule_GetDict(mainModule)
	return mainNamespace

}

func runFile(filename string, mainNamespace *C.PyObject) *C.PyObject {
	fileName := C.CString(filename)
	defer C.free(unsafe.Pointer(fileName))

	// Open the file
	file := C.fopen(fileName, C.CString("r"))

	// capture the result of PyRun_* so it is cleaned up later
	var result *C.PyObject
	if file == nil {
		result = C.PyRun_StringFlags(C.CString("print('Could not open file')"), C.Py_file_input, mainNamespace, nil, nil)
	} else {
		result = C.PyRun_FileEx(file, fileName, C.Py_file_input, mainNamespace, nil, C.int(1))
	}
	return result
}

func RunPythonFile(filename string) {
	mainNamespace := makeNameSpace()
	result := runFile(filename, mainNamespace)
	// print mainNamespace to stdout (for debugging)

	// Clean up

	defer C.Py_Finalize()
	defer C.free(unsafe.Pointer(result))
}

func Test() {
	RunPythonFile("./main.py")
}
