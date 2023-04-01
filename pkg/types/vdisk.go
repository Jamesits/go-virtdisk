package types

import "golang.org/x/sys/windows"

type VDiskHandle windows.Handle

const InvalidVDiskHandle = VDiskHandle(windows.InvalidHandle)
