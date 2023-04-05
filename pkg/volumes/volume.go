package volumes

import (
	"errors"
	"fmt"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
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
	v, err := path.AsFileNameW()
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

func GetSerial(volume types.Volume) (ret string, err error) {
	v, err := volume.AsObjectPathW()
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
	v, err := volume.AsObjectPathW()
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
