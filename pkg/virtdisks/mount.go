package virtdisks

import (
	"errors"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
)

// functions related to mounting the virtual drives

// Mount the virtual drives.
// Implements:
// - Mount-VHD
func Mount(handle types.VDiskHandle, noDriveLetter bool, readOnly bool) (err error) {
	flags := ffi.AttachVirtualDiskFlagNone
	if noDriveLetter {
		flags |= ffi.AttachVirtualDiskFlagNoDriveLetter
	}
	if readOnly {
		flags |= ffi.AttachVirtualDiskFlagReadOnly
	}
	flags |= ffi.AttachVirtualDiskFlagPermanentLifetime

	_, _, err = ffi.Virtdisk.AttachVirtualDisk.Call(
		uintptr(handle),
		types.IntPtrZero,
		uintptr(flags),
		types.IntPtrZero,
		types.IntPtrZero,
		types.IntPtrZero,
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}
	return nil
}

// Dismount the virtual drives.
// Implements:
// - Dismount-VHD
func Dismount(handle types.VDiskHandle) (err error) {
	_, _, err = ffi.Virtdisk.DetachVirtualDisk.Call(
		uintptr(handle),
		uintptr(ffi.DetachVirtualDiskFlagNone),
		types.IntPtrZero,
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}
	return nil
}

// Close the virtual drives handle.
func Close(handle types.VDiskHandle) (err error) {
	err = windows.CloseHandle(windows.Handle(handle))
	return
}
