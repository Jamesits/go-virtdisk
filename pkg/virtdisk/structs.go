package virtdisk

// Notes:
// - assuming uuid.UUID equals a [16]byte in the memory (which is not really guaranteed?)
//
// Types conversion:
// - ULONG -> uint32
// - ULONGLONG -> uint64
// - PCWSTR (const wchar_t*) -> *uint16 (should be []uint16, but...)

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
	ParentPath        *uint16
	SourcePath        *uint16
}

type CreateVirtualDiskParametersV2 struct {
	Version
	UniqueId                  uuid.UUID
	MaximumSize               uint64
	BlockSizeInBytes          uint32
	SectorSizeInBytes         uint32
	PhysicalSectorSizeInBytes uint32
	ParentPath                *uint16
	SourcePath                *uint16
	OpenFlags                 OpenVirtualDiskFlag
	ParentVirtualStorageType  VirtualStorageType
	SourceVirtualStorageType  VirtualStorageType
	ResiliencyGuid            uuid.UUID
}

type DeleteSnapshotVhdSetParametersV1 struct {
	Version
	SnapshotId uuid.UUID
}

type ExpandVirtualDiskParametersV1 struct {
	Version
	NewSize uint64
}

type GetVirtualDiskInfoV1 struct {
	Version

	VirtualSize  uint64
	PhysicalSize uint64
	BlockSize    uint32
	SectorSize   uint32

	Identifier uuid.UUID

	ParentResolved       bool
	ParentLocationBuffer *uint16

	ParentIdentifier   uuid.UUID
	ParentTimestamp    uint32
	VirtualStorageType VirtualStorageType
	ProviderSubtype    uint32
	Is4kAligned        bool
	IsLoaded           bool

	LogicalSectorSize  uint32
	PhysicalSectorSize uint32
	IsRemote           bool

	VhyPhysicalSectorSize   uint32
	SmallestSafeVirtualSize uint64
	FragmentationPercentage uint32
	VirtualDiskId           uuid.UUID

	Enabled      bool
	NewerChanges bool
	MostRecentId *uint16
}

type MergeVirtualDiskParametersV1 struct {
	Version
	MergeDepth uint32
}

type MergeVirtualDiskParametersV2 struct {
	Version
	MergeSourceDepth uint32
	MergeTargetDepth uint32
}

type MirrorVirtualDiskParametersV1 struct {
	Version
	MirrorVirtualDiskPath *uint16
}

type ModifyVhdSetParametersV1 struct {
	Version

	SnapshotId       uuid.UUID
	SnapshotFilePath *uint16

	DefaultSnapshotId uuid.UUID
	DefaultFilePath   *uint16
}

type OpenVirtualDiskParametersV1 struct {
	Version
	RWDepth uint32
}

type OpenVirtualDiskParametersV2 struct {
	Version
	GetInfoOnly    bool
	ReadOnly       bool
	ResiliencyGuid uuid.UUID
}

type QueryChangesVirtualDiskRange struct {
	ByteOffset uint64
	ByteLength uint64
	Reserved   uint64
}

type RawScsiVirtualDiskParametersV1 struct {
	Version
	RSVDHandle         bool
	DataIn             uint8
	CdbLength          uint8
	SenseInfoLength    uint8
	SrbFlags           uint32
	DataTransferLength uint32
	DataBuffer         uintptr
	SenseInfo          *uint8
	Cdb                *uint8
}

type RawScsiVirtualDiskResponseV1 struct {
	Version
}
