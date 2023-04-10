package types

import (
	"fmt"
	"strconv"
)

// Drive path in the format of `\\.\PhysicalDrive0`.
type Drive Path

func (d Drive) AsFileNameW() (*uint16, error) {
	return Path(d).asUTF16Ptr()
}

func (d Drive) AsObjectPathW() (*uint16, error) {
	return Path(d).asUTF16Ptr()
}

func (d Drive) AsId() (uint32, error) {
	if len(d) <= 17 {
		return 0, ErrorTyping
	}
	u, err := strconv.ParseUint(string(d[17:]), 10, 32)
	return uint32(u), err
}

func DriveFromId(id uint32) Drive {
	return Drive(fmt.Sprintf("\\\\.\\PhysicalDrive%d", id))
}
