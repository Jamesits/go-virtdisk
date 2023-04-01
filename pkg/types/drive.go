package types

import "fmt"

// Drive path in the format of `\\.\PhysicalDrive0`.
type Drive Path

func DriveFromId(id uint32) Drive {
	return Drive(fmt.Sprintf("\\\\.\\PhysicalDrive%d", id))
}
