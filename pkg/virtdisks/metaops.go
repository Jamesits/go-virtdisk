package virtdisks

import (
	"errors"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
	"unsafe"
)

func Compact(vDiskHandle types.VDiskHandle, flags ffi.CompactVirtualDiskFlag) error {
	parameters := &ffi.CompactVirtualDiskParametersV1{
		Version:  ffi.Version{Version: 1},
		Reserved: 0,
	}

	_, _, err := ffi.Virtdisk.CompactVirtualDisk.Call(
		uintptr(vDiskHandle),
		uintptr(flags),
		uintptr(unsafe.Pointer(parameters)),
		types.IntPtrZero,
	)

	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}

	return nil
}
