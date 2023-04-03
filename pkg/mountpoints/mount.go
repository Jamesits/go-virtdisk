package mountpoints

import (
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
)

func Mount(volume types.Volume, mountPoint types.MountPoint) error {
	// https://learn.microsoft.com/en-us/windows/win32/fileio/editing-drive-letter-assignments
	// The online documentation is wrong (likely some escaping problem); volume name require backslash in the end.
	win32VolumePath, err := volume.AsObjectPath()
	if err != nil {
		return err
	}

	win32MountPoint, err := mountPoint.AsDirectory()
	if err != nil {
		return err
	}

	err = windows.SetVolumeMountPoint(win32MountPoint, win32VolumePath)
	return err
}

func Dismount(mountPoint types.MountPoint) error {
	win32MountPoint, err := mountPoint.AsDirectory()
	if err != nil {
		return err
	}

	err = windows.DeleteVolumeMountPoint(win32MountPoint)
	return err
}
