package ioctl

type IoctlMethod uint32

const (
	MethodBuffered           IoctlMethod = 0
	MethodInDirect           IoctlMethod = 1
	MethodDirectToHardware   IoctlMethod = 1
	MethodOutDirect          IoctlMethod = 2
	MethodDirectFromHardware IoctlMethod = 2
	MethodNeither            IoctlMethod = 3
)

type IoctlAccess uint32

const (
	FileAnyAccess     IoctlAccess = 0
	FileSpecialAccess IoctlAccess = FileAnyAccess
	FileReadAccess    IoctlAccess = 0x0001
	FileWriteAccess   IoctlAccess = 0x0002
)

func CtlCode(DeviceType DeviceType, Function uint32, Method IoctlMethod, Access IoctlAccess) uint32 {
	// devioctl.h
	// https://www.pinvoke.net/default.aspx/kernel32/CTL_CODE.html

	return (uint32(DeviceType) << 16) | (uint32(Access) << 14) | (Function << 2) | (uint32(Method))
}
