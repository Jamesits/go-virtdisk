package virtdisks

import (
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"syscall"
	"unsafe"
)

// Open an existing virtual drives.
// Implements:
// - Get-VHD -Path
func Open(path types.Path, fileType ffi.VirtualStorageTypeDeviceType, accessMask ffi.VirtualDiskAccessMask, openFlags ffi.OpenVirtualDiskFlag) (handle types.VDiskHandle, err error) {
	storageType := ffi.VirtualStorageType{
		DeviceId: fileType,
		VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := path.AsFileName()
	if err != nil {
		return types.InvalidVDiskHandle, err
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
		return types.InvalidVDiskHandle, syscall.Errno(ret)
	}

	return handle, nil
}
