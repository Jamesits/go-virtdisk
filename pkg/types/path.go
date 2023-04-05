package types

import (
	"fmt"
	"golang.org/x/sys/windows"
	"os"
	"strings"
)

var PathSeparator = fmt.Sprintf("%c", os.PathSeparator)

// Path represents a generic Windows path in the paths system, without assumption of its kind or existence.
type Path string

func (p Path) AsDirectory() string {
	if !strings.HasSuffix(string(p), PathSeparator) {
		return string(p) + PathSeparator
	}

	return string(p)
}

func (p Path) AsFileName() string {
	for strings.HasSuffix(string(p), PathSeparator) {
		p = Path(strings.TrimSuffix(string(p), PathSeparator))
	}
	return string(p)
}

func (p Path) AsFileNameW() (*uint16, error) {
	return p.asUTF16Ptr()
}

func (p Path) asUTF16Ptr() (*uint16, error) {
	return windows.UTF16PtrFromString(string(p))
}

func (p Path) Segments() (ret []Path) {
	for _, seg := range strings.Split(string(p), PathSeparator) {
		ret = append(ret, Path(seg))
	}
	return
}

func PathFromUTF16(s []uint16) Path {
	return Path(windows.UTF16ToString(s))
}
