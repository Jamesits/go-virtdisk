package virtdisk

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jamesits/go-virtdisk/pkg/disk"
	"golang.org/x/sys/windows"
	"sync"
	"unsafe"
)

// VHDX is a thin wrapper that implements the most common operations on a VHDX file. It tries to mimic [the PowerShell
// API](https://learn.microsoft.com/en-us/powershell/module/hyper-v/new-vhd?view=windowsserver2022-ps).
type VHDX struct {
	storageType VirtualStorageType
	path        *uint16
	handle      windows.Handle
	lock        sync.Mutex
}

// constructors

// New creates a new VHDX object from scratch.
// Implements:
// - New-VHD -Dynamic
// - New-VHD -Fixed
func New(path string, sizeBytes uint64, dynamic bool, blockSizeBytes uint32, logicalSectorSizeBytes uint32, physicalSectorSizeBytes uint32) (ret *VHDX, err error) {
	ret = &VHDX{
		storageType: VirtualStorageType{
			DeviceId: VirtualStorageTypeDeviceVhdx,
			VendorId: VirtualStorageTypeVendorMicrosoft,
		},
		handle: 0,
	}
	ret.path, err = windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
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
		uintptr(unsafe.Pointer(&ret.storageType)), // VirtualStorageType
		uintptr(unsafe.Pointer(ret.path)),         // Path
		uintptr(VirtualDiskAccessNone),            // VirtualDiskAccessMask (must be none if using struct v2)
		intPtrZero,                                // SecurityDescriptor
		uintptr(createVirtualDiskFlag),            // Flags
		intPtrZero,                                // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),           // Parameters
		intPtrZero,                                // Overlapped
		uintptr(unsafe.Pointer(&ret.handle)),      // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	return ret, nil
}

// Open creates a new VHDX object from an existing file.
// Implements:
// - Get-VHD -Path
func Open(path string, openFlags OpenVirtualDiskFlag) (ret *VHDX, err error) {
	ret = &VHDX{
		storageType: VirtualStorageType{
			DeviceId: VirtualStorageTypeDeviceVhdx,
			VendorId: VirtualStorageTypeVendorMicrosoft,
		},
		handle: 0,
	}
	ret.path, err = windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	ret = &VHDX{handle: 0}
	_, _, err = virtdisk.OpenVirtualDisk.Call(
		uintptr(unsafe.Pointer(&ret.storageType)),
		uintptr(unsafe.Pointer(ret.path)),
		uintptr(openFlags),
		intPtrZero,
		uintptr(unsafe.Pointer(&ret.handle)),
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	return ret, nil
}

func FromHandle(handle windows.Handle) (ret *VHDX) {
	ret = &VHDX{handle: handle}
	return
}

// FromParent creates a child VHDX from a parent one.
// Implements:
// - New-VHD -Differencing
func FromParent(path string, parentPath string, sizeBytes uint64, blockSizeBytes uint32, physicalSectorSizeBytes uint32) (ret *VHDX, err error) {
	ret = &VHDX{
		storageType: VirtualStorageType{
			DeviceId: VirtualStorageTypeDeviceVhdx,
			VendorId: VirtualStorageTypeVendorMicrosoft,
		},
		handle: 0,
	}
	ret.path, err = windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	win32ParentPath, err := windows.UTF16PtrFromString(parentPath)
	if err != nil {
		return nil, err
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
		uintptr(unsafe.Pointer(&ret.storageType)), // VirtualStorageType
		uintptr(unsafe.Pointer(ret.path)),         // Path
		uintptr(VirtualDiskAccessNone),            // VirtualDiskAccessMask (must be none if using struct v2)
		intPtrZero,                                // SecurityDescriptor
		uintptr(CreateVirtualDiskFlagNone),        // Flags
		intPtrZero,                                // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),           // Parameters
		intPtrZero,                                // Overlapped
		uintptr(unsafe.Pointer(&ret.handle)),      // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	return ret, nil
}

// FromSource copies a disk's content into a VHDX.
// Implements:
// - New-VHD -SourceDisk
func FromSource(path string, sourceDisk uint32, dynamic bool, blockSizeBytes uint32) (ret *VHDX, err error) {
	ret = &VHDX{
		storageType: VirtualStorageType{
			DeviceId: VirtualStorageTypeDeviceVhdx,
			VendorId: VirtualStorageTypeVendorMicrosoft,
		},
		handle: 0,
	}
	ret.path, err = windows.UTF16PtrFromString(path)
	if err != nil {
		return nil, err
	}
	win32SourcePath, err := windows.UTF16PtrFromString(fmt.Sprintf("\\\\.\\PhysicalDrive%d", sourceDisk))
	if err != nil {
		return nil, err
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
		uintptr(unsafe.Pointer(&ret.storageType)), // VirtualStorageType
		uintptr(unsafe.Pointer(ret.path)),         // Path
		uintptr(VirtualDiskAccessNone),            // VirtualDiskAccessMask (must be none if using struct v2)
		intPtrZero,                                // SecurityDescriptor
		uintptr(createVirtualDiskFlag),            // Flags
		intPtrZero,                                // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),           // Parameters
		intPtrZero,                                // Overlapped
		uintptr(unsafe.Pointer(&ret.handle)),      // handle
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}
	return ret, nil
}

func FromDiskId(sourceDisk uint32) (ret *VHDX, err error) {
	win32SourcePath, err := windows.UTF16PtrFromString(fmt.Sprintf("\\\\.\\PhysicalDrive%d", sourceDisk))
	if err != nil {
		return nil, err
	}

	// get a handle to the disk
	// TODO: refactor this using the disk package
	diskHandle, err := windows.CreateFile(win32SourcePath, windows.GENERIC_READ, windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING, 0, windows.Handle(0))
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(diskHandle)

	// query the dependencies
	// https://stackoverflow.com/a/14175877
	// https://github.com/microsoft/Windows-classic-samples/blob/7cbd99ac1d2b4a0beffbaba29ea63d024ceff700/Samples/Hyper-V/Storage/cpp/GetStorageDependencyInformation.cpp
	var bufSize uint64
	_, _, err = virtdisk.GetStorageDependencyInformation.Call(
		uintptr(diskHandle),
		uintptr(GetStorageDependencyFlagDiskHandle),
		intPtrZero,
		intPtrZero,
		uintptr(unsafe.Pointer(&bufSize)),
	)
	if err != windows.ERROR_SUCCESS {
		return nil, err
	}
	b := make([]byte, bufSize)
	_, _, err = virtdisk.GetStorageDependencyInformation.Call(
		uintptr(diskHandle),
		uintptr(GetStorageDependencyFlagDiskHandle),
		uintptr(bufSize),
		uintptr(unsafe.Pointer(&b[0])),
		uintptr(unsafe.Pointer(&bufSize)),
	)

	return nil, nil
}

// Mount the VHDX.
// Implements:
// - Mount-VHD
func (v *VHDX) Mount(noDriveLetter bool, readOnly bool) (err error) {
	v.lock.Lock()
	defer v.lock.Unlock()

	flags := AttachVirtualDiskFlagNone
	if noDriveLetter {
		flags |= AttachVirtualDiskFlagNoDriveLetter
	}
	if readOnly {
		flags |= AttachVirtualDiskFlagReadOnly
	}

	_, _, err = virtdisk.AttachVirtualDisk.Call(
		uintptr(v.handle),
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

// Dismount the VHDX.
// Implements:
// - Dismount-VHD
func (v *VHDX) Dismount() (err error) {
	_, _, err = virtdisk.DetachVirtualDisk.Call(
		uintptr(v.handle),
		uintptr(DetachVirtualDiskFlagNone),
		intPtrZero,
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return err
	}
	return nil
}

func (v *VHDX) getPhysicalPathUTF16() (path []uint16, err error) {
	virtualDiskPhysicalPathSize := uint32(0)
	_, _, err = virtdisk.GetVirtualDiskPhysicalPath.Call(
		uintptr(v.handle),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathSize)),
		intPtrZero,
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}

	virtualDiskPhysicalPathUtf16 := make([]uint16, virtualDiskPhysicalPathSize)
	_, _, err = virtdisk.GetVirtualDiskPhysicalPath.Call(
		uintptr(v.handle),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathSize)),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathUtf16[0])),
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}

	return virtualDiskPhysicalPathUtf16, nil
}

func (v *VHDX) GetPhysicalPath() (path string, err error) {
	p, err := v.getPhysicalPathUTF16()
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return "", err
	}

	virtualDiskPhysicalPath := windows.UTF16ToString(p)
	return virtualDiskPhysicalPath, nil
}

func (v *VHDX) GetDisk() (d *disk.Disk, err error) {
	p, err := v.getPhysicalPathUTF16()
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}

	diskHandle, err := windows.CreateFile(
		&p[0],
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		0,
		windows.Handle(0),
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}

	return disk.FromHandle(diskHandle), nil
}
