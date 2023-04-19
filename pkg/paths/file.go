package paths

import (
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
)

func Remove(path types.Path) error {
	win32Path, err := path.AsFileNameW()
	if err != nil {
		return err
	}

	return windows.DeleteFile(win32Path)
}
