package virtdisks

import (
	"bytes"
	"errors"
	"github.com/google/uuid"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"github.com/jamesits/go-virtdisk/pkg/utils"
	"golang.org/x/sys/windows"
	"unsafe"
)

// Fork a differencing virtual drives from a parent.
// Implements:
// - New-VHD -Differencing
func Fork(path types.Path, parentPath types.Path, diskType ffi.VirtualStorageTypeDeviceType, sizeBytes uint64, blockSizeBytes uint32, physicalSectorSizeBytes uint32) (handle types.VDiskHandle, err error) {
	storageType := ffi.VirtualStorageType{
		DeviceId: diskType,
		VendorId: ffi.VirtualStorageTypeVendorMicrosoft,
	}
	win32Path, err := path.AsFileName()
	if err != nil {
		return types.InvalidVDiskHandle, err
	}
	win32ParentPath, err := parentPath.AsFileName()
	if err != nil {
		return types.InvalidVDiskHandle, err
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
		return types.InvalidVDiskHandle, err
	}
	return handle, nil
}

// GetParents get the parents of the child disk.
// returns windows.ERROR_VHD_INVALID_TYPE if it is not a differencing drives
func GetParents(handle types.VDiskHandle) (resolved bool, parents []types.Path, err error) {
	// https://github.com/microsoft/Windows-classic-samples/blob/7cbd99ac1d2b4a0beffbaba29ea63d024ceff700/Samples/Hyper-V/Storage/cpp/GetVirtualDiskInformation.cpp
	info := &ffi.GetVirtualDiskInfo{
		Version: ffi.Version{Version: 3},
	}
	var ret uintptr

	//size := uint32(unsafe.Sizeof(*info))
	size := uint32(64) // just make up something large
	for i := 0; i < 2; i++ {
		buf := make([]byte, size)
		_, _ = bytebuilder.Copy(buf, info)

		ret, _, _ = ffi.Virtdisk.GetVirtualDiskInformation.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&size)),
			uintptr(unsafe.Pointer(&buf[0])),
			types.IntPtrZero,
		)

		if ret == uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
			continue
		}

		if ret != 0 {
			break
		}

		reader := bytes.NewReader(buf)
		head := &ffi.GetVirtualDiskInfoV3H{}
		_, _ = bytebuilder.ReadPartial(reader, head)
		resolved = head.ParentResolved
		ps := utils.UTF16ByteArrayToStringSlice(buf[unsafe.Sizeof(head)+4:])
		for _, p := range ps {
			parents = append(parents, types.Path(p))
		}
		return
	}

	return false, nil, windows.Errno(ret)
}
