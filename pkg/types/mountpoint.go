package types

import (
	"strings"
)

// MountPoint is a reparse point.
// Examples:
// - `\\.\C:`
// - `\\.\C:\xxx\yyy`
type MountPoint Path

func (mp MountPoint) AsFileNameW() (*uint16, error) {
	return Path(mp).asUTF16Ptr()
}

func (mp MountPoint) AsDirectoryW() (*uint16, error) {
	if !strings.HasSuffix(string(mp), PathSeparator) {
		mp = mp + MountPoint(PathSeparator)
	}

	return Path(mp).asUTF16Ptr()
}
