package virtdisk

import "github.com/jamesits/libiferr/exception"

var virtdisk Virtdisk

const intPtrZero = uintptr(0)

func init() {
	// load DLLs
	virtdisk = Virtdisk{}
	err := virtdisk.Unmarshal("virtdisk.dll")
	exception.HardFailWithReason("unable to load virtdisk.dll", err)
}
