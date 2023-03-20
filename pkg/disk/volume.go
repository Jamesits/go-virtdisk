package disk

import (
	"errors"
	"github.com/jamesits/go-virtdisk/pkg/utils"
	"golang.org/x/sys/windows"
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
