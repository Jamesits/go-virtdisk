package mountpoints

import (
	"errors"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"github.com/jamesits/go-virtdisk/pkg/utils"
	"golang.org/x/sys/windows"
	"path/filepath"
)

// FromPath returns the mount point, and its relative path from the mount point, from a paths.
func FromPath(path types.Path) (mp types.MountPoint, rel types.Path, err error) {
	v, err := path.AsFileNameW()
	if err != nil {
		return "", "", err
	}

	b := make([]uint16, ffi.MaxPath+1)
	err = windows.GetVolumePathName(v, &b[0], ffi.MaxPath+1)
	if err != nil {
		return "", "", err
	}

	mp = types.MountPoint(types.PathFromUTF16(b))
	relTmp, err := filepath.Rel(string(mp), string(path))
	rel = types.Path(relTmp)
	return mp, rel, err
}

// FromVolume returns a list of mount points (drives name: `C:\` or directory) for a volumes.
func FromVolume(volume types.Volume) (ret []types.MountPoint, err error) {
	v, err := volume.AsObjectPathW()
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

	mps := utils.UTF16ToStringSlice(b)
	for _, mp := range mps {
		ret = append(ret, types.MountPoint(mp))
	}

	return ret, nil
}
