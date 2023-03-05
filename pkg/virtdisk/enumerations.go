package virtdisk

// Notes:
// - enum is by default uint16; v1_enum tagged enums are uint32.

type ApplySnapshotVhdSetFlag uint16

const (
	ApplySnapshotVhdSetFlagNone     ApplySnapshotVhdSetFlag = 0x00000000
	ApplySnapshotVhdSetFlagWritable ApplySnapshotVhdSetFlag = 0x00000001
)

type AttachVirtualDiskFlag uint16

const (
	AttachVirtualDiskFlagNone                          AttachVirtualDiskFlag = 0x00000000
	AttachVirtualDiskFlagReadOnly                      AttachVirtualDiskFlag = 0x00000001
	AttachVirtualDiskFlagNoDriveLetter                 AttachVirtualDiskFlag = 0x00000002
	AttachVirtualDiskFlagPermanentLifetime             AttachVirtualDiskFlag = 0x00000004
	AttachVirtualDiskFlagNoLocalHost                   AttachVirtualDiskFlag = 0x00000008
	AttachVirtualDiskFlagNoSecurityDescriptor          AttachVirtualDiskFlag = 0x00000010
	AttachVirtualDiskFlagBypassDefaultEncryptionPolicy AttachVirtualDiskFlag = 0x00000020
	AttachVirtualDiskFlagNonPnp                        AttachVirtualDiskFlag = 0x00000040
	AttachVirtualDiskFlagRestrictedRange               AttachVirtualDiskFlag = 0x0000080
	AttachVirtualDiskFlagSinglePartition               AttachVirtualDiskFlag = 0x00000100
	AttachVirtualDiskFlagRegisterVolume                AttachVirtualDiskFlag = 0x00000200
)

type CompactVirtualDiskFlag uint16

const (
	CompactVirtualDiskFlagNone         CompactVirtualDiskFlag = 0x00000000
	CompactVirtualDiskFlagNoZeroScan   CompactVirtualDiskFlag = 0x00000001
	CompactVirtualDiskFlagNoBlockMoves CompactVirtualDiskFlag = 0x00000002
)

type CreateVirtualDiskFlag uint16

const (
	CreateVirtualDiskFlagNone                              CreateVirtualDiskFlag = 0x0
	CreateVirtualDiskFlagFullPhysicalAllocation            CreateVirtualDiskFlag = 0x1
	CreateVirtualDiskFlagPreventWritesToSourceDisk         CreateVirtualDiskFlag = 0x2
	CreateVirtualDiskFlagDoNotCopyMetadataFromParent       CreateVirtualDiskFlag = 0x4
	CreateVirtualDiskFlagCreateBackingStorage              CreateVirtualDiskFlag = 0x8
	CreateVirtualDiskFlagUseChangeTrackingSourceLimit      CreateVirtualDiskFlag = 0x10
	CreateVirtualDiskFlagUseRctSourceLimit                 CreateVirtualDiskFlag = 0x10
	CreateVirtualDiskFlagPreserveParentChangeTrackingState CreateVirtualDiskFlag = 0x20
	CreateVirtualDiskFlagVhdSetUseOriginalBackingStorage   CreateVirtualDiskFlag = 0x40
	CreateVirtualDiskFlagSparseFile                        CreateVirtualDiskFlag = 0x80
	CreateVirtualDiskFlagPmemCompatible                                          = 0x100
)

type DeleteSnapshotVhdSetFlag uint16

const (
	DeleteSnapshotVhdSetFlagNone       DeleteSnapshotVhdSetFlag = 0x00000000
	DeleteSnapshotVhdSetFlagPersistRct DeleteSnapshotVhdSetFlag = 0x00000001
)

type DependentDiskFlag uint16

const (
	DependentDiskFlagNone                     DependentDiskFlag = 0x00000000
	DependentDiskFlagMultiBackingFiles        DependentDiskFlag = 0x00000001
	DependentDiskFlagFullyAllocated           DependentDiskFlag = 0x00000002
	DependentDiskFlagReadOnly                 DependentDiskFlag = 0x00000004
	DependentDiskFlagRemote                   DependentDiskFlag = 0x00000008
	DependentDiskFlagSystemVolume             DependentDiskFlag = 0x00000010
	DependentDiskFlagSystemVolumeParent       DependentDiskFlag = 0x00000020
	DependentDiskFlagRemovable                DependentDiskFlag = 0x00000040
	DependentDiskFlagNoDriveLetter            DependentDiskFlag = 0x00000080
	DependentDiskFlagParent                   DependentDiskFlag = 0x00000100
	DependentDiskFlagNoHostDisk               DependentDiskFlag = 0x00000200
	DependentDiskFlagPermanentLifetime        DependentDiskFlag = 0x00000400
	DependentDiskFlagSupportCompressedVolumes DependentDiskFlag = 0x00000800
)

type DetachVirtualDiskFlag uint16

const (
	DetachVirtualDiskFlagNone DetachVirtualDiskFlag = 0x00000000
)

type ExpandVirtualDiskFlag uint16

const (
	ExpandVirtualDiskFlagNone ExpandVirtualDiskFlag = 0x00000000
)

type GetStorageDependencyFlag uint16

const (
	GetStorageDependencyFlagNone        GetStorageDependencyFlag = 0x00000000
	GetStorageDependencyFlagHostVolumes GetStorageDependencyFlag = 0x00000001
	GetStorageDependencyFlagDiskHandle  GetStorageDependencyFlag = 0x00000002
)

type MergeVirtualDiskFlag uint16

const (
	MergeVirtualDiskFlagNone MergeVirtualDiskFlag = 0x00000000
)

type MirrorVirtualDiskFlag uint16

const (
	MirrorVirtualDiskFlagNone                 MirrorVirtualDiskFlag = 0x00000000
	MirrorVirtualDiskFlagExistingFile         MirrorVirtualDiskFlag = 0x00000001
	MirrorVirtualDiskFlagSkipMirrorActivation MirrorVirtualDiskFlag = 0x00000002
)

type ModifyVhdSetFlag uint16

const (
	ModifyVhdSetFlagNone              ModifyVhdSetFlag = 0x00000000
	ModifyVhdSetFlagWriteableSnapshot ModifyVhdSetFlag = 0x00000001
)

type OpenVirtualDiskFlag uint16

const (
	OpenVirtualDiskFlagNone                        OpenVirtualDiskFlag = 0x00000000
	OpenVirtualDiskFlagNoParents                   OpenVirtualDiskFlag = 0x00000001
	OpenVirtualDiskFlagBlankFile                   OpenVirtualDiskFlag = 0x00000002
	OpenVirtualDiskFlagBootDrive                   OpenVirtualDiskFlag = 0x00000004
	OpenVirtualDiskFlagCachedIo                    OpenVirtualDiskFlag = 0x00000008
	OpenVirtualDiskFlagCustomDiffChain             OpenVirtualDiskFlag = 0x00000010
	OpenVirtualDiskFlagParentCachedIo              OpenVirtualDiskFlag = 0x00000020
	OpenVirtualDiskFlagVhdSetFileOnly              OpenVirtualDiskFlag = 0x00000040
	OpenVirtualDiskFlagIgnoreRelativeParentLocator OpenVirtualDiskFlag = 0x00000080
	OpenVirtualDiskFlagNoWriteHardening            OpenVirtualDiskFlag = 0x00000100
	OpenVirtualDiskFlagSupportCompressedVolumes    OpenVirtualDiskFlag = 0x00000200
)

type RawScsiVirtualDiskFlag uint16

const (
	RawScsiVirtualDiskFlagNone RawScsiVirtualDiskFlag = 0x00000000
)

type ResizeVirtualDiskFlag uint16

const (
	ResizeVirtualDiskFlagNone                            ResizeVirtualDiskFlag = 0x0
	ResizeVirtualDiskFlagAllowUnsafeVirtualSize          ResizeVirtualDiskFlag = 0x1
	ResizeVirtualDiskFlagResizeToSmallestSafeVirtualSize ResizeVirtualDiskFlag = 0x2
)

type VirtualDiskAccessMask uint32

const (
	VirtualDiskAccessNone     VirtualDiskAccessMask = 0x00000000
	VirtualDiskAccessAttachRo VirtualDiskAccessMask = 0x00010000
	VirtualDiskAccessAttachRw VirtualDiskAccessMask = 0x00020000
	VirtualDiskAccessDetach   VirtualDiskAccessMask = 0x00040000
	VirtualDiskAccessGetInfo  VirtualDiskAccessMask = 0x00080000
	VirtualDiskAccessCreate   VirtualDiskAccessMask = 0x00100000
	VirtualDiskAccessMetaops  VirtualDiskAccessMask = 0x00200000
	VirtualDiskAccessRead     VirtualDiskAccessMask = 0x000d0000
	VirtualDiskAccessAll      VirtualDiskAccessMask = 0x003f0000
	VirtualDiskAccessWritable VirtualDiskAccessMask = 0x00320000
)
