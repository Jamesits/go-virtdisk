package utils

import "strings"

func FixMountPointPath(s string) string {
	//if strings.HasPrefix(mountPoint, "\\\\.\\") {
	//	mountPoint = mountPoint[4:]
	//}

	if !strings.HasSuffix(s, "\\") {
		s = s + "\\"
	}

	return s
}
