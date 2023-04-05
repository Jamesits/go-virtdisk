package paths

import (
	"github.com/jamesits/go-virtdisk/pkg/mountpoints"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"github.com/jamesits/go-virtdisk/pkg/volumes"
)

func NormalizeAsVolumeAbsolute(path types.Path) (types.Path, error) {
	// TODO: wrap error?
	vol, rel, err := mountpoints.FromPath(path)
	if err != nil {
		return "", err
	}

	volAbs, err := volumes.FromMountPoint(vol)
	if err != nil {
		return "", err
	}

	return types.Path(string(volAbs) + string(rel)), nil
}
