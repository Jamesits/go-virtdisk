package virtdisk

// Notes:
// - assuming uuid.UUID equals a [16]byte in the memory (which is not really guaranteed?)
//
// Types conversion:
// - ULONG -> uint32
// - ULONGLONG -> uint64
// - PCWSTR (const wchar_t*) -> uintptr (should be []uint16, but...)

import "github.com/google/uuid"

// common

type Version struct {
	Version uint16
}

type VirtualStorageType struct {
	DeviceId uint32
	VendorId uuid.UUID
}

// specific

type ApplySnapshotVhdSetParametersV1 struct {
	Version
	SnapshotId     uuid.UUID
	LeafSnapshotId uuid.UUID
}

type AttachVirtualDiskParametersV1 struct {
	Version
	Reserved uint32
}

type AttachVirtualDiskParametersV2 struct {
	Version
	RestrictedOffset uint64
	RestrictedLength uint64
}

type CompactVirtualDiskParametersV1 struct {
	Version
	Reserved uint32
}

type CreateVirtualDiskParametersV1 struct {
	Version
	UniqueId          uuid.UUID
	MaximumSize       uint64
	BlockSizeInBytes  uint32
	SectorSizeInBytes uint32
	ParentPath        uintptr
	SourcePath        uintptr
}

type CreateVirtualDiskParametersV2 struct {
	Version
	UniqueId                  uuid.UUID
	MaximumSize               uint64
	BlockSizeInBytes          uint32
	SectorSizeInBytes         uint32
	PhysicalSectorSizeInBytes uint32
	ParentPath                uintptr
	SourcePath                uintptr
	OpenFlags                 OpenVirtualDiskFlag
	ParentVirtualStorageType  VirtualStorageType
	SourceVirtualStorageType  VirtualStorageType
	ResiliencyGuid            uuid.UUID
}
