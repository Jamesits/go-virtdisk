package ffi

var Setupapi setupapi
var Kernel32 kernel32
var Virtdisk virtdisk

func init() {
	// Errors are ignored if some functions are not found in the dll.
	// I don't know it is a good way towards better compatibility (since Kernel.dll sometimes delete functions)
	// but let's just accept there will be errors during any syscall.
	_ = Setupapi.Unmarshal("setupapi.dll")
	_ = Kernel32.Unmarshal("kernel32.dll")
	_ = Virtdisk.Unmarshal("virtdisk.dll")
}
