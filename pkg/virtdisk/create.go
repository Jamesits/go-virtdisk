package virtdisk

import (
	"errors"
	"github.com/google/uuid"
	"golang.org/x/sys/windows"
	"unsafe"
)

// New virtual disk from scratch.
// Implements:
// - New-VHD -Dynamic
// - New-VHD -Fixed
func New(path string, diskType VirtualStorageTypeDeviceType, sizeBytes uint64, dynamic bool, blockSizeBytes uint32, logicalSectorSizeBytes uint32, physicalSectorSizeBytes uint32) (handle windows.Handle, err error) {
	storageType := VirtualStorageType{
		DeviceId: diskType,
		VendorId: VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}
	param := CreateVirtualDiskParametersV2{
		Version:                   Version{Version: 2},
		UniqueId:                  uuid.Nil,
		MaximumSize:               sizeBytes,
		BlockSizeInBytes:          blockSizeBytes,
		SectorSizeInBytes:         logicalSectorSizeBytes,
		PhysicalSectorSizeInBytes: physicalSectorSizeBytes,
		ParentPath:                nil,
		SourcePath:                nil,
		OpenFlags:                 OpenVirtualDiskFlagNone,
		ParentVirtualStorageType:  VirtualStorageType{},
		SourceVirtualStorageType:  VirtualStorageType{},
		ResiliencyGuid:            uuid.Nil,
	}
	createVirtualDiskFlag := CreateVirtualDiskFlagNone
	if dynamic {
		createVirtualDiskFlag |= CreateVirtualDiskFlagFullPhysicalAllocation
	} else {
		createVirtualDiskFlag |= CreateVirtualDiskFlagSparseFile
	}

	_, _, err = virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&storageType)), // VirtualStorageType
		uintptr(unsafe.Pointer(win32Path)),    // Path
		uintptr(VirtualDiskAccessNone),        // VirtualDiskAccessMask (must be none if using struct v2)
		intPtrZero,                            // SecurityDescriptor
		uintptr(createVirtualDiskFlag),        // Flags
		intPtrZero,                            // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),       // Parameters
		intPtrZero,                            // Overlapped
		uintptr(unsafe.Pointer(&handle)),      // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return windows.InvalidHandle, err
	}

	return handle, nil
}

// Fork a differencing virtual disk from a parent.
// Implements:
// - New-VHD -Differencing
func Fork(path string, parentPath string, diskType VirtualStorageTypeDeviceType, sizeBytes uint64, blockSizeBytes uint32, physicalSectorSizeBytes uint32) (handle windows.Handle, err error) {
	storageType := VirtualStorageType{
		DeviceId: diskType,
		VendorId: VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}
	win32ParentPath, err := windows.UTF16PtrFromString(parentPath)
	if err != nil {
		return windows.InvalidHandle, err
	}
	param := CreateVirtualDiskParametersV2{
		Version:                   Version{Version: 2},
		UniqueId:                  uuid.Nil,
		MaximumSize:               sizeBytes,
		BlockSizeInBytes:          blockSizeBytes,
		SectorSizeInBytes:         0,
		PhysicalSectorSizeInBytes: physicalSectorSizeBytes,
		ParentPath:                win32ParentPath,
		SourcePath:                nil,
		OpenFlags:                 OpenVirtualDiskFlagNone,
		ParentVirtualStorageType:  VirtualStorageType{}, // TODO: check if this works
		SourceVirtualStorageType:  VirtualStorageType{},
		ResiliencyGuid:            uuid.Nil,
	}

	_, _, err = virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&storageType)), // VirtualStorageType
		uintptr(unsafe.Pointer(win32Path)),    // Path
		uintptr(VirtualDiskAccessNone),        // VirtualDiskAccessMask (must be none if using struct v2)
		intPtrZero,                            // SecurityDescriptor
		uintptr(CreateVirtualDiskFlagNone),    // Flags
		intPtrZero,                            // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),       // Parameters
		intPtrZero,                            // Overlapped
		uintptr(unsafe.Pointer(&handle)),      // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return windows.InvalidHandle, err
	}
	return handle, nil
}

// Mirror a disk into a new virtual disk.
// Implements:
// - New-VHD -SourceDisk
func Mirror(path string, sourceDiskPath string, diskType VirtualStorageTypeDeviceType, dynamic bool, blockSizeBytes uint32) (handle windows.Handle, err error) {
	storageType := VirtualStorageType{
		DeviceId: diskType,
		VendorId: VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return windows.InvalidHandle, err
	}
	win32SourcePath, err := windows.UTF16PtrFromString(sourceDiskPath)
	if err != nil {
		return windows.InvalidHandle, err
	}
	param := CreateVirtualDiskParametersV2{
		Version:                   Version{Version: 2},
		UniqueId:                  uuid.Nil,
		MaximumSize:               0, // FIXME: ???
		BlockSizeInBytes:          blockSizeBytes,
		SectorSizeInBytes:         0,
		PhysicalSectorSizeInBytes: 0,
		ParentPath:                nil,
		SourcePath:                win32SourcePath,
		OpenFlags:                 OpenVirtualDiskFlagNone,
		ParentVirtualStorageType:  VirtualStorageType{},
		SourceVirtualStorageType:  VirtualStorageType{},
		ResiliencyGuid:            uuid.Nil,
	}
	createVirtualDiskFlag := CreateVirtualDiskFlagNone
	if dynamic {
		createVirtualDiskFlag |= CreateVirtualDiskFlagFullPhysicalAllocation
	} else {
		createVirtualDiskFlag |= CreateVirtualDiskFlagSparseFile
	}

	_, _, err = virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&storageType)), // VirtualStorageType
		uintptr(unsafe.Pointer(win32Path)),    // Path
		uintptr(VirtualDiskAccessNone),        // VirtualDiskAccessMask (must be none if using struct v2)
		intPtrZero,                            // SecurityDescriptor
		uintptr(createVirtualDiskFlag),        // Flags
		intPtrZero,                            // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),       // Parameters
		intPtrZero,                            // Overlapped
		uintptr(unsafe.Pointer(&handle)),      // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return windows.InvalidHandle, err
	}
	return handle, nil
}
