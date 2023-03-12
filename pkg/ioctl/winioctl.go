package ioctl

import "github.com/google/uuid"

// Winioctl.h

type PartitionStyle uint16

const (
	PartitionStyleMBR PartitionStyle = 0
	PartitionStyleGPT PartitionStyle = 1
	PartitionStyleRAW PartitionStyle = 2
)

type CreateDiskGpt struct {
	PartitionStyle    PartitionStyle
	DiskId            uuid.UUID
	MaxPartitionCount uint32
}

type CreateDiskMbr struct {
	PartitionStyle PartitionStyle
	Signature      uint32
	_              [16]uint8 // padded to the same length as CreateDiskGpt
}
