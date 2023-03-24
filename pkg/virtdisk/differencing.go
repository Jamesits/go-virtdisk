package virtdisk

import (
	"bytes"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/utils"
	"golang.org/x/sys/windows"
	"unsafe"
)

// returns windows.ERROR_VHD_INVALID_TYPE if it is not a differencing disk
func GetParents(handle windows.Handle) (resolved bool, parents []string, err error) {
	// https://github.com/microsoft/Windows-classic-samples/blob/7cbd99ac1d2b4a0beffbaba29ea63d024ceff700/Samples/Hyper-V/Storage/cpp/GetVirtualDiskInformation.cpp
	info := &GetVirtualDiskInfo{
		Version: Version{Version: 3},
	}
	var ret uintptr

	//size := uint32(unsafe.Sizeof(*info))
	size := uint32(64) // just make up something large
	for i := 0; i < 2; i++ {
		buf := make([]byte, size)
		_, _ = bytebuilder.Copy(buf, info)

		ret, _, _ = virtdisk.GetVirtualDiskInformation.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&size)),
			uintptr(unsafe.Pointer(&buf[0])),
			intPtrZero,
		)

		if ret == uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
			continue
		}

		if ret != 0 {
			break
		}

		reader := bytes.NewReader(buf)
		head := &GetVirtualDiskInfoV3H{}
		_, _ = bytebuilder.ReadPartial(reader, head)
		resolved = head.ParentResolved
		parents = utils.UTF16ByteArrayToStringSlice(buf[unsafe.Sizeof(head)+4:])
		return
	}

	return false, nil, windows.Errno(ret)
}
