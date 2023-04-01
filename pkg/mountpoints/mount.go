package mountpoints

import (
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
)

func Mount(volume types.Volume, mountPoint types.MountPoint) error {
	// https://learn.microsoft.com/en-us/windows/win32/fileio/editing-drive-letter-assignments
	win32VolumePath, err := volume.AsFileName()
	if err != nil {
		return err
	}

	win32MountPoint, err := mountPoint.AsNormalizedDirectory()
	if err != nil {
		return err
	}

	err = windows.SetVolumeMountPoint(win32MountPoint, win32VolumePath)
	return err
}

func Dismount(mountPoint types.MountPoint) error {
	win32MountPoint, err := mountPoint.AsNormalizedDirectory()
	if err != nil {
		return err
	}

	err = windows.DeleteVolumeMountPoint(win32MountPoint)
	return err
}
