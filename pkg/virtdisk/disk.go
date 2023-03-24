package virtdisk

import (
	"bytes"
	"errors"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/utils"
	"golang.org/x/sys/windows"
	"unsafe"
)

// GetVirtualDiskBackingFiles returns filesystem paths to all the virtual disk backing files, sorted from the child to the parent.
func GetVirtualDiskBackingFiles(diskDevicePath string) ([]string, error) {
	var err error
	win32SourcePath, err := windows.UTF16PtrFromString(diskDevicePath)
	if err != nil {
		return nil, err
	}

	// get a handle to the disk
	// we don't need any access permission here
	diskHandle, err := windows.CreateFile(win32SourcePath, 0, windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING, 0, windows.Handle(0))
	if err != nil {
		return nil, err
	}
	defer windows.CloseHandle(diskHandle)

	// query the dependencies
	// https://stackoverflow.com/a/14175877
	// https://github.com/microsoft/Windows-classic-samples/blob/7cbd99ac1d2b4a0beffbaba29ea63d024ceff700/Samples/Hyper-V/Storage/cpp/GetStorageDependencyInformation.cpp
	sdInfo := StorageDependencyInfo{
		StorageDependencyInfoH: StorageDependencyInfoH{Version: Version{Version: 2}},
	}
	bufSize := uint64(unsafe.Sizeof(sdInfo))
	for i := 0; i < 2; i++ {
		b := make([]byte, bufSize)
		_, _ = bytebuilder.Copy(b, &sdInfo)
		errcode, _, _ := virtdisk.GetStorageDependencyInformation.Call(
			uintptr(diskHandle),
			uintptr(GetStorageDependencyFlagHostVolumes|GetStorageDependencyFlagDiskHandle),
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
			return nil, err
		}

		// parse data
		var sdh StorageDependencyInfoH
		var sde StorageDependencyInfoType2
		var ret []string
		reader := bytes.NewReader(b)
		_, _ = bytebuilder.ReadPartial(reader, &sdh)
		for j := uint32(0); j < sdh.NumberEntries; j++ {
			_, _ = bytebuilder.ReadPartial(reader, &sde)
			depPath := windows.UTF16PtrToString(sde.HostVolumeName) + windows.UTF16PtrToString(sde.DependentVolumeRelativePath)
			ret = append(ret, depPath)
		}

		return ret, nil
	}

	return nil, utils.ErrorRetryLimitExceeded
}

func getPhysicalPathUTF16(handle windows.Handle) (path []uint16, err error) {
	virtualDiskPhysicalPathSize := uint32(0)
	_, _, err = virtdisk.GetVirtualDiskPhysicalPath.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathSize)),
		intPtrZero,
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}

	virtualDiskPhysicalPathUtf16 := make([]uint16, virtualDiskPhysicalPathSize)
	_, _, err = virtdisk.GetVirtualDiskPhysicalPath.Call(
		uintptr(handle),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathSize)),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathUtf16[0])),
	)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return nil, err
	}

	return virtualDiskPhysicalPathUtf16, nil
}

// GetPhysicalPath returns normalized disk path of a opened virtual disk.
func GetPhysicalPath(handle windows.Handle) (path string, err error) {
	p, err := getPhysicalPathUTF16(handle)
	if !errors.Is(err, windows.ERROR_SUCCESS) {
		return "", err
	}

	virtualDiskPhysicalPath := windows.UTF16ToString(p)
	return virtualDiskPhysicalPath, nil
}
