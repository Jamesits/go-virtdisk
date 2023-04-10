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
func New(path types.Path, diskType ffi.VirtualStorageTypeDeviceType, sizeBytes uint64, dynamic bool, blockSizeBytes uint32, logicalSectorSizeBytes uint32, physicalSectorSizeBytes uint32) (handle types.VDiskHandle, err error) {
	storageType := ffi.VirtualStorageType{
		DeviceId: diskType,
		VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := path.AsFileNameW()
	if err != nil {
		return types.InvalidVDiskHandle, err
	}
	param := ffi.CreateVirtualDiskParametersV2{
		Version:                   ffi.Version{Version: 2},
		UniqueId:                  uuid.Nil,
		MaximumSize:               sizeBytes,
		BlockSizeInBytes:          blockSizeBytes,
		SectorSizeInBytes:         logicalSectorSizeBytes,
		PhysicalSectorSizeInBytes: physicalSectorSizeBytes,
		OpenFlags:                 ffi.OpenVirtualDiskFlagNone,
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
		return types.InvalidVDiskHandle, err
	}

	return handle, nil
}

// Mirror a drives into a new virtual drives.
// Implements:
// - New-VHD -SourceDisk
func Mirror(path types.Path, sourceDiskPath types.Path, diskType ffi.VirtualStorageTypeDeviceType, dynamic bool, blockSizeBytes uint32) (handle types.VDiskHandle, err error) {
	storageType := ffi.VirtualStorageType{
		DeviceId: diskType,
		VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := path.AsFileNameW()
	if err != nil {
		return types.InvalidVDiskHandle, err
	}
	win32SourcePath, err := sourceDiskPath.AsFileNameW()
	if err != nil {
		return types.InvalidVDiskHandle, err
	}
	param := ffi.CreateVirtualDiskParametersV2{
		Version:                   ffi.Version{Version: 2},
		UniqueId:                  uuid.Nil,
		MaximumSize:               0, // FIXME: ???
		BlockSizeInBytes:          blockSizeBytes,
		SectorSizeInBytes:         0,
		PhysicalSectorSizeInBytes: 0,
		SourcePath:                win32SourcePath,
		OpenFlags:                 ffi.OpenVirtualDiskFlagNone,
		SourceVirtualStorageType: ffi.VirtualStorageType{
			DeviceId: diskType,
			VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
		},
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
		return types.InvalidVDiskHandle, err
	}
	return handle, nil
}
