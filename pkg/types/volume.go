package types

import (
	"golang.org/x/sys/windows"
	"strings"
)

// Volume GUID path in the format of `\\?\Volume{00000000-0000-0000-0000-00000000000}\`.
// Note that a volume's path always end with `\`.
// https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-volume
type Volume Path

func (v Volume) AsFileNameW() (*uint16, error) {
	// https://learn.microsoft.com/en-us/windows/win32/fileio/naming-a-volume
	return windows.UTF16PtrFromString(strings.TrimRight(string(v), PathSeparator))
}

func (v Volume) AsObjectPathW() (*uint16, error) {
	return Path(v).asUTF16Ptr()
}
