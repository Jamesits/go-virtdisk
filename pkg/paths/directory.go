package paths

import (
	"errors"
	"github.com/jamesits/go-virtdisk/pkg/mountpoints"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
	"strings"
)

// Mkdir creates a directory.
// Notes:
// - The operation is always idempotent.
// - All the directories on the way will be created.
// - All types of paths including "\\.\C:\dir" and "\\?\Volume{GUID}\path" is supported
func Mkdir(path types.Path) error {
	vol, rel, err := mountpoints.FromPath(path)
	if err != nil {
		return err
	}
	if rel == "" {
		return nil
	}

	segments := rel.Segments()
	cPath := vol
	for _, seg := range segments {
		cPath = types.MountPoint(strings.Join([]string{string(cPath), string(seg)}, types.PathSeparator))
		win32cPath, err := cPath.AsFileNameW()
		if err != nil {
			return err
		}

		err = windows.CreateDirectory(win32cPath, nil)
		if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
			continue
		}
		if err != nil {
			return err
		}
	}

	return nil
}

// Rmdir ensures a directory does not exist.
func Rmdir(path types.Path) error {
	vol, rel, err := mountpoints.FromPath(path)
	if err != nil {
		return err
	}
	if rel == "" {
		return nil
	}

	if rel == "." {
		// is a volumes mount point
		return mountpoints.Dismount(vol)
	}

	win32Path, err := path.AsFileNameW()
	if err != nil {
		return err
	}

	err = windows.RemoveDirectory(win32Path)
	if err != nil {
		return err
	}

	return nil
}
