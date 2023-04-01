package ffi

type Method uint32

const (
	MethodBuffered           Method = 0
	MethodInDirect           Method = 1
	MethodDirectToHardware   Method = 1
	MethodOutDirect          Method = 2
	MethodDirectFromHardware Method = 2
	MethodNeither            Method = 3
)

type Access uint32

const (
	FileAnyAccess     Access = 0
	FileSpecialAccess Access = FileAnyAccess
	FileReadAccess    Access = 0x0001
	FileWriteAccess   Access = 0x0002
)

func CtlCode(DeviceType DeviceType, Function uint32, Method Method, Access Access) uint32 {
	// devioctl.h
	// https://www.pinvoke.net/default.aspx/kernel32/CTL_CODE.html

	return (uint32(DeviceType) << 16) | (uint32(Access) << 14) | (Function << 2) | (uint32(Method))
}

var (
	IoctlDiskCreateDisk             = CtlCode(FileDeviceDisk, 0x0016, MethodBuffered, FileReadAccess|FileWriteAccess)
	IoctlStorageGetDeviceNumber     = uint32(0x2d1080)
	IoctlStorageQueryProperty       = uint32(0x2D1400)
	IoctlVolumeGetVolumeDiskExtents = CtlCode(0x00000056, 0, MethodBuffered, FileAnyAccess)
)
