package virtdisks

// converts between virtual disks, disks, handles and paths paths.

import (
	"bytes"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
	"unsafe"
)

// GetVirtualDiskBackingFiles returns filesystem paths to all the virtual drives backing paths, sorted from the child to the parent.
func GetVirtualDiskBackingFiles(drive types.Drive) ([]types.Path, error) {
	var err error
	win32SourcePath, err := drive.AsFileNameW()
	if err != nil {
		return nil, err
	}

	// get a handle to the drives
	// we don't need any access permission here
	diskHandle, err := windows.CreateFile(win32SourcePath, 0, windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING, 0, windows.Handle(0))
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(diskHandle)

	// query the dependencies
	// https://stackoverflow.com/a/14175877
	// https://github.com/microsoft/Windows-classic-samples/blob/7cbd99ac1d2b4a0beffbaba29ea63d024ceff700/Samples/Hyper-V/Storage/cpp/GetStorageDependencyInformation.cpp
	sdInfo := ffi.StorageDependencyInfo{
		StorageDependencyInfoH: ffi.StorageDependencyInfoH{Version: ffi.Version{Version: 2}},
	}
	bufSize := uint64(unsafe.Sizeof(sdInfo))
	for i := 0; i < 2; i++ {
		b := make([]byte, bufSize)
		_, _ = bytebuilder.Copy(b, &sdInfo)
		errcode, _, _ := ffi.Virtdisk.GetStorageDependencyInformation.Call(
			uintptr(diskHandle),
			uintptr(ffi.GetStorageDependencyFlagHostVolumes|ffi.GetStorageDependencyFlagDiskHandle),
			uintptr(bufSize),
			uintptr(unsafe.Pointer(&b[0])),
			uintptr(unsafe.Pointer(&bufSize)),
		)
		if errcode == uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
			continue
		}
		if errcode == uintptr(windows.SEVERITY_ERROR) {
			// not a VHD
			return nil, nil
		}
		if errcode != uintptr(windows.ERROR_SUCCESS) {
			return nil, windows.Errno(errcode)
		}

		// parse data
		var sdh ffi.StorageDependencyInfoH
		var sde ffi.StorageDependencyInfoType2
		var ret []types.Path
		reader := bytes.NewReader(b)
		_, _ = bytebuilder.ReadPartial(reader, &sdh)
		for j := uint32(0); j < sdh.NumberEntries; j++ {
			_, _ = bytebuilder.ReadPartial(reader, &sde)
			depPath := windows.UTF16PtrToString(sde.HostVolumeName) + windows.UTF16PtrToString(sde.DependentVolumeRelativePath)
			ret = append(ret, types.Path(depPath))
		}

		return ret, nil
	}

	return nil, types.ErrorRetryLimitExceeded
}

func getPhysicalPathUTF16(handle types.VDiskHandle) (path []uint16, err error) {
	virtualDiskPhysicalPathSize := uint32(0)
	errcode, _, _ := ffi.Virtdisk.GetVirtualDiskPhysicalPath.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathSize)),
		types.IntPtrZero,
	)
	if errcode != uintptr(windows.ERROR_INSUFFICIENT_BUFFER) {
		return nil, windows.Errno(errcode)
	}

	virtualDiskPhysicalPathUtf16 := make([]uint16, virtualDiskPhysicalPathSize)
	errcode, _, _ = ffi.Virtdisk.GetVirtualDiskPhysicalPath.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathSize)),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathUtf16[0])),
	)
	if errcode != 0 {
		return nil, windows.Errno(errcode)
	}

	return virtualDiskPhysicalPathUtf16, nil
}

// GetPhysicalPath returns normalized drives path of an opened virtual drive.
// Required permission: virtdisks.VirtualDiskAccessGetInfo
func GetPhysicalPath(handle types.VDiskHandle) (path types.Drive, err error) {
	p, err := getPhysicalPathUTF16(handle)
	if err != nil {
		return "", err
	}

	virtualDiskPhysicalPath := types.Drive(types.PathFromUTF16(p))
	return virtualDiskPhysicalPath, nil
}
