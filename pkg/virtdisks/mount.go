package virtdisks

import (
	"errors"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

// functions related to opening and mounting the virtual drives

// Open an existing virtual drives.
// Implements:
// - Get-VHD -Path
func Open(path string, fileType ffi.VirtualStorageTypeDeviceType, accessMask ffi.VirtualDiskAccessMask, openFlags ffi.OpenVirtualDiskFlag) (handle windows.Handle, err error) {
	storageType := ffi.VirtualStorageType{
		DeviceId: fileType,
		VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}

	ret, _, _ := ffi.Virtdisk.OpenVirtualDisk.Call(
		uintptr(unsafe.Pointer(&storageType)),
		uintptr(unsafe.Pointer(win32Path)),
		uintptr(accessMask),
		uintptr(openFlags),
		types.IntPtrZero,
		uintptr(unsafe.Pointer(&handle)),
	)
	if ret != 0 {
		return windows.InvalidHandle, syscall.Errno(ret)
	}

	return handle, nil
}

// Mount the virtual drives.
// Implements:
// - Mount-VHD
func Mount(handle windows.Handle, noDriveLetter bool, readOnly bool) (err error) {
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
func Dismount(handle windows.Handle) (err error) {
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
func Close(handle windows.Handle) (err error) {
	err = windows.CloseHandle(handle)
	return
}
