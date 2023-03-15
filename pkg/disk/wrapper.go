package disk

import "golang.org/x/sys/windows"

type Disk struct {
	handle windows.Handle
}

func FromHandle(handle windows.Handle) (ret *Disk) {
	ret = &Disk{handle: handle}
	return
}
