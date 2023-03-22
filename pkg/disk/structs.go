package disk

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

type DriveLayoutInformationEx struct {
	PartitionStyle PartitionStyle
	_              uint16 // PartitionStyle in DRIVE_LAYOUT_INFORMATION_EX is defined as a DWORD, so we explicitly pad here
	PartitionCount uint32
}

type DriveLayoutInformationGpt struct {
	DiskId               uuid.UUID
	StartingUsableOffset uint64
	UsableLength         uint64
	MaxPartitionCount    uint32
}

type DriveLayoutInformationMbr struct {
	Signature uint32
	CheckSum  uint32
	_         [32]uint8
}

type DriveLayoutInformationExMbr struct {
	DriveLayoutInformationEx
	DriveLayoutInformationMbr
}

type DriveLayoutInformationExGpt struct {
	DriveLayoutInformationEx
	DriveLayoutInformationGpt
}

type PartitionInformationEx struct {
	PartitionStyle     PartitionStyle
	StartingOffset     uint64
	PartitionLength    uint64
	PartitionNumber    uint32
	RewritePartition   bool
	IsServicePartition bool
}

type PartitionInformationGpt struct {
	PartitionType uuid.UUID
	PartitionId   uuid.UUID
	Attributes    uint64
	Name          [36]uint16
}

type PartitionInformationMbr struct {
	PartitionType       byte
	BootIndicator       bool
	RecognizedPartition bool
	HiddenSectors       uint32
	PartitionId         uuid.UUID
	_                   [85]uint8
}

type PartitionInformationExGpt struct {
	PartitionInformationEx
	PartitionInformationGpt
}

type PartitionInformationExMbr struct {
	PartitionInformationEx
	PartitionInformationMbr
}

type StorageDeviceNumber struct {
	DeviceType      DeviceType
	DeviceNumber    uint32
	PartitionNumber uint32
}

// https://learn.microsoft.com/en-us/windows/win32/api/winioctl/ns-winioctl-storage_property_query
type StoragePropertyQueryH struct {
	PropertyId StoragePropertyId
	QueryType  StorageQueryType
}

type StoragePropertyQuery struct {
	StoragePropertyQueryH
	AdditionalParameters byte
}

// https://learn.microsoft.com/en-us/windows/win32/api/winioctl/ns-winioctl-storage_device_descriptor
type StorageDeviceDescriptorH struct {
	Version               uint32
	Size                  uint32
	DeviceType            byte
	DeviceTypeModifier    byte
	RemovableMedia        bool
	CommandQueueing       bool
	VendorIdOffset        uint32
	ProductIdOffset       uint32
	ProductRevisionOffset uint32
	SerialNumberOffset    uint32
	BusType               uint32
	RawPropertiesLength   uint32
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/winioctl/ns-winioctl-volume_disk_extents
type VolumeDiskExtentsH struct {
	NumberOfDiskExtents uint32
}

type VolumeDiskExtents struct {
	VolumeDiskExtentsH
	Extents DiskExtent
}

// https://learn.microsoft.com/zh-cn/windows/win32/api/winioctl/ns-winioctl-disk_extent
type DiskExtent struct {
	DiskNumber     uint32
	StartingOffset uint64
	ExtentLength   uint64
}

// https://learn.microsoft.com/en-us/windows/win32/api/winioctl/ns-winioctl-disk_partition_info
type DiskPartitionInfoH struct {
	SizeOfPartitionInfo uint32
	PartitionStyle      uint16
}

type DiskPartitionInfoMbr struct {
	Signature uint32
	CheckSum  uint32
}

type DiskPartitionInfoGpt struct {
	DiskID uuid.UUID
}

// setupapi.h

type SPDeviceInterfaceData struct {
	Size               uint32
	InterfaceClassGuid uuid.UUID
	Flags              uint32
	_                  uint64
}

type SPDeviceInterfaceDetailDataH struct {
	Size uint32
}

type SPDeviceInterfaceDetailDataA struct {
	SPDeviceInterfaceDetailDataH
	_ uint8
}

type SPDeviceInterfaceDetailDataW struct {
	SPDeviceInterfaceDetailDataH
	_ uint16
}
