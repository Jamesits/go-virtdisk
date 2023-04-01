package types

import "fmt"

// Drive path in the format of `\\.\PhysicalDrive0`.
type Drive Path

func (d Drive) AsFileName() (*uint16, error) {
	return Path(d).asUTF16Ptr()
}

func (d Drive) AsObjectPath() (*uint16, error) {
	return Path(d).asUTF16Ptr()
}

func DriveFromId(id uint32) Drive {
	return Drive(fmt.Sprintf("\\\\.\\PhysicalDrive%d", id))
}
