package virtdisk

// Notes:
// - assuming uuid.UUID equals a [16]byte in the memory (which is not really guaranteed?)
//
// Types conversion:
// - UCHAR -> uint8
// - ULONG -> uint32
// - ULONGLONG -> uint64
// - PCWSTR (const wchar_t*) -> *uint16 (should be []uint16, but...)

import (
	"github.com/google/uuid"
)

// common

type Version struct {
	Version uint16
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

type CreateVirtualDiskParametersV3 struct {
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
	SourceLimitPath           *uint16
	BackingStorageType        VirtualStorageType
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
}

type GetVirtualDiskInfoV2 struct {
	Version
	Identifier uuid.UUID
}

type GetVirtualDiskInfoV3 struct {
	Version
	ParentResolved       bool
	ParentLocationBuffer *uint16
}

type GetVirtualDiskInfoV4 struct {
	Version
	ParentIdentifier uuid.UUID
}

type GetVirtualDiskInfoV5 struct {
	Version
	ParentTimestamp uint32
}

type GetVirtualDiskInfoV6 struct {
	Version
	VirtualStorageType VirtualStorageType
}

type GetVirtualDiskInfoV7 struct {
	Version
	ProviderSubtype VirtualDiskInfoProviderSubtype
}

type GetVirtualDiskInfoV8 struct {
	Version
	Is4kAligned bool
}

type GetVirtualDiskInfoV9 struct {
	Version
	IsLoaded bool
}

type GetVirtualDiskInfoV10 struct {
	Version
	LogicalSectorSize  uint32
	PhysicalSectorSize uint32
	IsRemote           bool
}

type GetVirtualDiskInfoV11 struct {
	Version
	VhyPhysicalSectorSize uint32
}

type GetVirtualDiskInfoV12 struct {
	Version
	SmallestSafeVirtualSize uint64
}

type GetVirtualDiskInfoV13 struct {
	Version
	FragmentationPercentage uint32
}

type GetVirtualDiskInfoV14 struct {
	Version
	VirtualDiskId uuid.UUID
}

type GetVirtualDiskInfoV15 struct {
	Version
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
	ScsiStatus         uint8
	SenseInfoLength    uint8
	DataTransferLength uint32
}

type ResizeVirtualDiskParametersV1 struct {
	Version
	NewSize uint64
}

type SetVirtualDiskInfoV1 struct {
	Version
	ParentFilePath *uint16
}

type SetVirtualDiskInfoV2 struct {
	Version
	UniqueIdentifier uuid.UUID
}

type SetVirtualDiskInfoV3 struct {
	Version
	ChildDepth     uint32
	ParentFilePath *uint16
}

type SetVirtualDiskInfoV4 struct {
	Version
	VhdPhysicalSectorSize uint32
}

type SetVirtualDiskInfoV5 struct {
	Version
	VirtualDiskId uuid.UUID
}

type SetVirtualDiskInfoV6 struct {
	Version
	ChangeTrackingEnabled bool
}

type SetVirtualDiskInfoV7 struct {
	Version
	LinkageId      uuid.UUID
	ParentFilePath *uint16
}

type StorageDependencyInfoV1 struct {
	Version
	NumberEntries   uint32
	Version1Entries []*StorageDependencyInfoType1
}

type StorageDependencyInfoType1 struct {
	DependencyTypeFlags   DependentDiskFlag
	ProviderSpecificFlags uint32
	VirtualStorageType    VirtualStorageType
}

type StorageDependencyInfoV2 struct {
	Version
	NumberEntries   uint32
	Version2Entries []*StorageDependencyInfoType2
}

type StorageDependencyInfoType2 struct {
	StorageDependencyInfoType1
	AncestorLevel               uint32
	DependencyDeviceName        *uint16
	HostVolumeName              *uint16
	DependentVolumeName         *uint16
	DependentVolumeRelativePath *uint16
}

type TakeSnapshotVhdSetParametersV1 struct {
	Version
	SnapshotId uuid.UUID
}

type VirtualDiskProgress struct {
	OperationStatus uint32 // windows.Errno
	CurrentValue    uint64
	CompletionValue uint64
}

type VirtualStorageType struct {
	DeviceId VirtualStorageTypeDeviceType
	VendorId uuid.UUID
}
