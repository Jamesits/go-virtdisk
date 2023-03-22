package disk

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/utils"
	"golang.org/x/sys/windows"
	"strings"
	"unsafe"
)

// GetVolumes returns GUID paths of all the volumes exist on the system.
// Example: `\\?\Volume{7cb86808-7fee-44d2-ae19-137066372203}\`
func GetVolumes() (ret []string, err error) {
	s := uint32(65536)
	b := make([]uint16, s)

	handle, err := windows.FindFirstVolume(&b[0], s)
	if handle == windows.InvalidHandle {
		return ret, err
	}
	defer windows.FindVolumeClose(handle)
	for {
		ret = append(ret, windows.UTF16ToString(b))
		err := windows.FindNextVolume(handle, &b[0], s)
		if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
			break
		}
	}

	return ret, nil
}

// GetVolumeMountPoints returns a list of mount points (drive name: `C:\` or directory) for a volume.
func GetVolumeMountPoints(VolumeGUIDPath string) (ret []string, err error) {
	v, err := windows.UTF16PtrFromString(VolumeGUIDPath)
	if err != nil {
		return nil, err
	}

	// test buffer length
	var bufLength uint32
	err = windows.GetVolumePathNamesForVolumeName(v, nil, 0, &bufLength)
	if !errors.Is(err, windows.ERROR_MORE_DATA) {
		return nil, err
	}

	b := make([]uint16, bufLength)
	err = windows.GetVolumePathNamesForVolumeName(v, &b[0], bufLength, &bufLength)
	if err != nil {
		return nil, err
	}

	ret = utils.UTF16ToStringSlice(b)

	return ret, nil
}

func GetVolumeBackingDrives(VolumeGUIDPath string) (drives []string, err error) {
	// https://stackoverflow.com/questions/29212597/how-to-enumerate-disk-volume-names
	// https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-volume
	vp, err := windows.UTF16PtrFromString(strings.TrimRight(VolumeGUIDPath, "\\"))
	if err != nil {
		return nil, err
	}

	vHandle, err := windows.CreateFile(
		vp,
		// Even if GENERIC_READ is not specified, metadata can still be read; GENERIC_READ requires administrator privileges.
		// https://stackoverflow.com/questions/327718/how-to-list-physical-disks#comment89593842_11683906
		0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		windows.Handle(0),
	)
	if err != nil || vHandle == windows.InvalidHandle {
		return nil, err
	}
	defer windows.CloseHandle(vHandle)

	// https://stackoverflow.com/a/5664841
	extents := VolumeDiskExtents{}
	bytesReturned := uint32(unsafe.Sizeof(extents))

	for i := 0; i < 2; i++ {
		b := make([]byte, bytesReturned)
		err = windows.DeviceIoControl(
			vHandle,
			ioctlVolumeGetVolumeDiskExtents,
			nil,
			0,
			&b[0],
			bytesReturned,
			&bytesReturned,
			nil,
		)
		if err != nil {
			return nil, err
		}
		_ = bytebuilder.Unmarshal(b, &extents)
		if bytesReturned > uint32(unsafe.Sizeof(extents)) {
			continue
		}

		// parse extents
		reader := bytes.NewReader(b)
		header := VolumeDiskExtentsH{}
		extent := DiskExtent{}
		_, _ = bytebuilder.ReadPartial(reader, &header)
		_, _ = bytebuilder.Skip(reader, 4)
		for j := uint32(0); j < header.NumberOfDiskExtents; j++ {
			_, _ = bytebuilder.ReadPartial(reader, &extent)
			drives = append(drives, DiskKernelObjectPathById(extent.DiskNumber))
		}

		return
	}

	return nil, ErrorRetryLimitExceeded
}

func GetVolumeSerial(VolumeGUIDPath string) (ret string, err error) {
	v, err := windows.UTF16PtrFromString(VolumeGUIDPath)
	if err != nil {
		return "", err
	}

	var volSerial uint32
	err = windows.GetVolumeInformation(
		v,
		(*uint16)(nil),
		0,
		&volSerial,
		(*uint32)(nil),
		(*uint32)(nil),
		(*uint16)(nil),
		0,
	)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", volSerial), nil
}

func GetVolumeLabel(VolumeGUIDPath string) (ret string, err error) {
	v, err := windows.UTF16PtrFromString(VolumeGUIDPath)
	if err != nil {
		return "", err
	}

	b := make([]uint16, MaxPath+1)
	err = windows.GetVolumeInformation(
		v,
		&b[0],
		MaxPath+1,
		(*uint32)(nil),
		(*uint32)(nil),
		(*uint32)(nil),
		(*uint16)(nil),
		0,
	)
	if err != nil {
		return "", err
	}

	return windows.UTF16ToString(b), nil
}
