package types

import "strings"

// MountPoint is a reparse point.
// Examples:
// - `\\.\C:`
// - `\\.\C:\xxx\yyy`
type MountPoint Path

func (mp MountPoint) AsFileName() (*uint16, error) {
	return Path(mp).asUTF16Ptr()
}

func (mp MountPoint) AsNormalizedDirectory() (*uint16, error) {
	//if strings.HasPrefix(mountPoint, "\\\\.\\") {
	//	mountPoint = mountPoint[4:]
	//}

	if !strings.HasSuffix(string(mp), "\\") {
		mp = mp + "\\"
	}

	return Path(mp).asUTF16Ptr()
}
