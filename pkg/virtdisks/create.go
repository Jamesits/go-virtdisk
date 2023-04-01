package virtdisks

import (
	"errors"
	"github.com/google/uuid"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
	"unsafe"
)

// New virtual drives from scratch.
// Implements:
// - New-VHD -Dynamic
// - New-VHD -Fixed
func New(path string, diskType ffi.VirtualStorageTypeDeviceType, sizeBytes uint64, dynamic bool, blockSizeBytes uint32, logicalSectorSizeBytes uint32, physicalSectorSizeBytes uint32) (handle windows.Handle, err error) {
	storageType := ffi.VirtualStorageType{
		DeviceId: diskType,
		VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}
	param := ffi.CreateVirtualDiskParametersV2{
		Version:                   ffi.Version{Version: 2},
		UniqueId:                  uuid.Nil,
		MaximumSize:               sizeBytes,
		BlockSizeInBytes:          blockSizeBytes,
		SectorSizeInBytes:         logicalSectorSizeBytes,
		PhysicalSectorSizeInBytes: physicalSectorSizeBytes,
		ParentPath:                nil,
		SourcePath:                nil,
		OpenFlags:                 ffi.OpenVirtualDiskFlagNone,
		ParentVirtualStorageType:  ffi.VirtualStorageType{},
		SourceVirtualStorageType:  ffi.VirtualStorageType{},
		ResiliencyGuid:            uuid.Nil,
	}
	createVirtualDiskFlag := ffi.CreateVirtualDiskFlagNone
	if dynamic {
		createVirtualDiskFlag |= ffi.CreateVirtualDiskFlagFullPhysicalAllocation
	} else {
		createVirtualDiskFlag |= ffi.CreateVirtualDiskFlagSparseFile
	}

	_, _, err = ffi.Virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&storageType)), // VirtualStorageType
		uintptr(unsafe.Pointer(win32Path)),    // Path
		uintptr(ffi.VirtualDiskAccessNone),    // VirtualDiskAccessMask (must be none if using struct v2)
		types.IntPtrZero,                      // SecurityDescriptor
		uintptr(createVirtualDiskFlag),        // Flags
		types.IntPtrZero,                      // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),       // Parameters
		types.IntPtrZero,                      // Overlapped
		uintptr(unsafe.Pointer(&handle)),      // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return windows.InvalidHandle, err
	}

	return handle, nil
}

// Fork a differencing virtual drives from a parent.
// Implements:
// - New-VHD -Differencing
func Fork(path string, parentPath string, diskType ffi.VirtualStorageTypeDeviceType, sizeBytes uint64, blockSizeBytes uint32, physicalSectorSizeBytes uint32) (handle windows.Handle, err error) {
	storageType := ffi.VirtualStorageType{
		DeviceId: diskType,
		VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}
	win32ParentPath, err := windows.UTF16PtrFromString(parentPath)
	if err != nil {
		return windows.InvalidHandle, err
	}
	param := ffi.CreateVirtualDiskParametersV2{
		Version:                   ffi.Version{Version: 2},
		UniqueId:                  uuid.Nil,
		MaximumSize:               sizeBytes,
		BlockSizeInBytes:          blockSizeBytes,
		SectorSizeInBytes:         0,
		PhysicalSectorSizeInBytes: physicalSectorSizeBytes,
		ParentPath:                win32ParentPath,
		SourcePath:                nil,
		OpenFlags:                 ffi.OpenVirtualDiskFlagNone,
		ParentVirtualStorageType:  ffi.VirtualStorageType{}, // TODO: check if this works
		SourceVirtualStorageType:  ffi.VirtualStorageType{},
		ResiliencyGuid:            uuid.Nil,
	}

	_, _, err = ffi.Virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&storageType)),  // VirtualStorageType
		uintptr(unsafe.Pointer(win32Path)),     // Path
		uintptr(ffi.VirtualDiskAccessNone),     // VirtualDiskAccessMask (must be none if using struct v2)
		types.IntPtrZero,                       // SecurityDescriptor
		uintptr(ffi.CreateVirtualDiskFlagNone), // Flags
		types.IntPtrZero,                       // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),        // Parameters
		types.IntPtrZero,                       // Overlapped
		uintptr(unsafe.Pointer(&handle)),       // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return windows.InvalidHandle, err
	}
	return handle, nil
}

// Mirror a drives into a new virtual drives.
// Implements:
// - New-VHD -SourceDisk
func Mirror(path string, sourceDiskPath string, diskType ffi.VirtualStorageTypeDeviceType, dynamic bool, blockSizeBytes uint32) (handle windows.Handle, err error) {
	storageType := ffi.VirtualStorageType{
		DeviceId: diskType,
		VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}
	win32SourcePath, err := windows.UTF16PtrFromString(sourceDiskPath)
	if err != nil {
		return windows.InvalidHandle, err
	}
	param := ffi.CreateVirtualDiskParametersV2{
		Version:                   ffi.Version{Version: 2},
		UniqueId:                  uuid.Nil,
		MaximumSize:               0, // FIXME: ???
		BlockSizeInBytes:          blockSizeBytes,
		SectorSizeInBytes:         0,
		PhysicalSectorSizeInBytes: 0,
		ParentPath:                nil,
		SourcePath:                win32SourcePath,
		OpenFlags:                 ffi.OpenVirtualDiskFlagNone,
		ParentVirtualStorageType:  ffi.VirtualStorageType{},
		SourceVirtualStorageType:  ffi.VirtualStorageType{},
		ResiliencyGuid:            uuid.Nil,
	}
	createVirtualDiskFlag := ffi.CreateVirtualDiskFlagNone
	if dynamic {
		createVirtualDiskFlag |= ffi.CreateVirtualDiskFlagFullPhysicalAllocation
	} else {
		createVirtualDiskFlag |= ffi.CreateVirtualDiskFlagSparseFile
	}

	_, _, err = ffi.Virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&storageType)), // VirtualStorageType
		uintptr(unsafe.Pointer(win32Path)),    // Path
		uintptr(ffi.VirtualDiskAccessNone),    // VirtualDiskAccessMask (must be none if using struct v2)
		types.IntPtrZero,                      // SecurityDescriptor
		uintptr(createVirtualDiskFlag),        // Flags
		types.IntPtrZero,                      // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),       // Parameters
		types.IntPtrZero,                      // Overlapped
		uintptr(unsafe.Pointer(&handle)),      // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return windows.InvalidHandle, err
	}
	return handle, nil
}
