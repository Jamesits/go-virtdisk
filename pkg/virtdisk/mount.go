package virtdisk

import (
	"errors"
	"golang.org/x/sys/windows"
	"unsafe"
)

// functions related to opening and mounting the virtual disk

// Open an existing virtual disk.
// Implements:
// - Get-VHD -Path
func Open(path string, fileType VirtualStorageTypeDeviceType, openFlags OpenVirtualDiskFlag) (ret windows.Handle, err error) {
	storageType := VirtualStorageType{
		DeviceId: fileType,
		VendorId: VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}

	_, _, err = virtdisk.OpenVirtualDisk.Call(
		uintptr(unsafe.Pointer(&storageType)),
		uintptr(unsafe.Pointer(win32Path)),
		uintptr(openFlags),
		intPtrZero,
		uintptr(unsafe.Pointer(&ret)),
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return windows.InvalidHandle, err
	}

	return ret, nil
}

// Mount the virtual disk.
// Implements:
// - Mount-VHD
func Mount(handle windows.Handle, noDriveLetter bool, readOnly bool) (err error) {
	flags := AttachVirtualDiskFlagNone
	if noDriveLetter {
		flags |= AttachVirtualDiskFlagNoDriveLetter
	}
	if readOnly {
		flags |= AttachVirtualDiskFlagReadOnly
	}

	_, _, err = virtdisk.AttachVirtualDisk.Call(
		uintptr(handle),
		intPtrZero,
		uintptr(flags),
		intPtrZero,
		intPtrZero,
		intPtrZero,
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}
	return nil
}

// Dismount the virtual disk.
// Implements:
// - Dismount-VHD
func Dismount(handle windows.Handle) (err error) {
	_, _, err = virtdisk.DetachVirtualDisk.Call(
		uintptr(handle),
		uintptr(DetachVirtualDiskFlagNone),
		intPtrZero,
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}
	return nil
}

// Close the virtual disk handle.
func Close(handle windows.Handle) (err error) {
	err = windows.CloseHandle(handle)
	return
}
