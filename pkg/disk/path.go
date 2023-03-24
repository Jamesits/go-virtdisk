package disk

import (
	"errors"
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"strings"
)

var osPathSeparatorString = fmt.Sprintf("%c", os.PathSeparator)

// Mkdir creates a directory.
// Notes:
// - The operation is always idempotent.
// - All the directories on the way will be created.
// - All types of paths including "\\.\C:\dir" and "\\?\Volume{GUID}\path" is supported
func Mkdir(path string) error {
	vol, rel, err := GetMountPointByFileName(path)
	if err != nil {
		return err
	}
	if rel == "" {
		return nil
	}

	segments := strings.Split(rel, osPathSeparatorString)
	cPath := vol
	for _, seg := range segments {
		cPath = strings.Join([]string{cPath, seg}, osPathSeparatorString)
		win32cPath, err := windows.UTF16PtrFromString(cPath)
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
func Rmdir(path string) error {
	vol, rel, err := GetMountPointByFileName(path)
	if err != nil {
		return err
	}
	if rel == "" {
		return nil
	}

	if rel == "." {
		// is a volume mount point
		return Dismount(vol)
	}

	win32Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	err = windows.RemoveDirectory(win32Path)
	if err != nil {
		return err
	}

	return nil
}
