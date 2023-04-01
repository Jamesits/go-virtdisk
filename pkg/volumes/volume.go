package volumes

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
	"unsafe"
)

// List returns GUID paths of all the volumes exist on the system.
func List() (ret []types.Volume, err error) {
	s := uint32(65536)
	b := make([]uint16, s)

	handle, err := windows.FindFirstVolume(&b[0], s)
	if handle == windows.InvalidHandle {
		return ret, err
	}
	defer windows.FindVolumeClose(handle)
	for {
		ret = append(ret, types.Volume(types.PathFromUTF16(b)))
		err := windows.FindNextVolume(handle, &b[0], s)
		if errors.Is(err, windows.ERROR_NO_MORE_FILES) {
			break
		}
	}

	return ret, nil
}

func FromMountPoint(path types.MountPoint) (ret types.Volume, err error) {
	v, err := path.AsFileName()
	if err != nil {
		return "", err
	}

	b := make([]uint16, ffi.MaxPath+1)
	err = windows.GetVolumeNameForVolumeMountPoint(v, &b[0], ffi.MaxPath+1)
	if err != nil {
		return "", err
	}

	return types.Volume(windows.UTF16ToString(b)), nil
}

func GetBackingDrives(volume types.Volume) (drives []types.Drive, err error) {
	// https://stackoverflow.com/questions/29212597/how-to-enumerate-disk-volume-names
	vp, err := volume.AsFileName()
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
	extents := ffi.VolumeDiskExtents{}
	bytesReturned := uint32(unsafe.Sizeof(extents))

	for i := 0; i < 2; i++ {
		b := make([]byte, bytesReturned)
		err = windows.DeviceIoControl(
			vHandle,
			ffi.IoctlVolumeGetVolumeDiskExtents,
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
		header := ffi.VolumeDiskExtentsH{}
		extent := ffi.DiskExtent{}
		_, _ = bytebuilder.ReadPartial(reader, &header)
		_, _ = bytebuilder.Skip(reader, 4)
		for j := uint32(0); j < header.NumberOfDiskExtents; j++ {
			_, _ = bytebuilder.ReadPartial(reader, &extent)
			drives = append(drives, types.DriveFromId(extent.DiskNumber))
		}

		return
	}

	return nil, types.ErrorRetryLimitExceeded
}

func GetSerial(volume types.Volume) (ret string, err error) {
	v, err := volume.AsObjectPath()
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

func GetLabel(volume types.Volume) (ret string, err error) {
	v, err := volume.AsObjectPath()
	if err != nil {
		return "", err
	}

	b := make([]uint16, ffi.MaxPath+1)
	err = windows.GetVolumeInformation(
		v,
		&b[0],
		ffi.MaxPath+1,
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
