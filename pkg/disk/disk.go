package disk

import (
	"errors"
	"fmt"
	"github.com/jamesits/go-bytebuilder"
	"golang.org/x/sys/windows"
	"unsafe"
)

var GuidDevInterfaceDisk = &windows.GUID{
	Data1: 0x53F56307,
	Data2: 0xB6BF,
	Data3: 0x11D0,
	Data4: [8]byte{0x94, 0xF2, 0x00, 0xA0, 0xC9, 0x1E, 0xFB, 0x8B},
}

var setupapi Setupapi

func init() {
	err := setupapi.Unmarshal("setupapi.dll")
	if err != nil {
		panic(err)
	}
}

// GetDisks returns the device path of the disk.
// The returned string is a object under `\GLOBAL??`. Examples:
// - VHDX: `\\?\scsi#disk&ven_msft&prod_virtual_disk#2&1f4adffe&0&000001#{53f56307-b6bf-11d0-94f2-00a0c91efb8b}`
// - NVMe disk: `\\?\scsi#disk&ven_msft&prod_virtual_disk#2&1f4adffe&0&000001#{53f56307-b6bf-11d0-94f2-00a0c91efb8b}`
func GetDisks() (ret []string, err error) {
	// https://stackoverflow.com/a/18183115
	handle, err := windows.SetupDiGetClassDevsEx(GuidDevInterfaceDisk, "", uintptr(0), windows.DIGCF_PRESENT|windows.DIGCF_DEVICEINTERFACE, windows.DevInfo(0), "")
	if windows.Handle(handle) == windows.InvalidHandle || err != nil {
		return nil, err
	}

	deviceIndex := uintptr(0)
	devInterfaceData := SPDeviceInterfaceData{}
	devInterfaceData.Size = uint32(unsafe.Sizeof(devInterfaceData))
	for {
		successful, _, _ := setupapi.SetupDiEnumDeviceInterfaces.Call(
			uintptr(handle),
			uintptr(0),
			uintptr(unsafe.Pointer(GuidDevInterfaceDisk)),
			deviceIndex,
			uintptr(unsafe.Pointer(&devInterfaceData)),
		)
		if successful == 0 {
			break
		}

		// query size
		var s uintptr
		_, _, err := setupapi.SetupDiGetDeviceInterfaceDetailW.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&devInterfaceData)),
			uintptr(0),
			uintptr(0),
			uintptr(unsafe.Pointer(&s)),
			uintptr(0),
		)

		if !errors.Is(err, windows.ERROR_INSUFFICIENT_BUFFER) {
			return nil, err
		}

		var interfaceDetailData SPDeviceInterfaceDetailData
		interfaceDetailData.Size = uint32(unsafe.Sizeof(interfaceDetailData)) + 4 // must add a WCHAR
		b := make([]byte, s)
		_, _ = bytebuilder.Copy(b, &interfaceDetailData)

		successful, _, err = setupapi.SetupDiGetDeviceInterfaceDetailW.Call(
			uintptr(handle),
			uintptr(unsafe.Pointer(&devInterfaceData)),
			uintptr(unsafe.Pointer(&b[0])),
			s,
			uintptr(0),
			uintptr(0),
		)
		if successful != 1 {
			return nil, err
		}

		_, pathBuffer := bytebuilder.CarCdr[SPDeviceInterfaceDetailData](b)
		path := windows.UTF16ToString(bytebuilder.SliceCast[uint8, uint16](pathBuffer)) // string is NUL terminated
		ret = append(ret, path)

		deviceIndex++
	}

	return
}

// GetDiskNumber returns the PhysicalDrive number of the hard disk device.
func GetDiskNumber(diskDevicePath string) (uint32, error) {
	dp, err := windows.UTF16PtrFromString(diskDevicePath)
	if err != nil {
		return 0, err
	}

	dHandle, err := windows.CreateFile(
		dp,
		// Even if GENERIC_READ is not specified, metadata can still be read; GENERIC_READ requires administrator privileges.
		// https://stackoverflow.com/questions/327718/how-to-list-physical-disks#comment89593842_11683906
		0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		windows.Handle(0),
	)
	if err != nil || dHandle == windows.InvalidHandle {
		return 0, err
	}

	var n StorageDeviceNumber
	var bytesReturned uint32
	err = windows.DeviceIoControl(
		dHandle,
		ioctlStorageGetDeviceNumber,
		nil,
		0,
		(*byte)(unsafe.Pointer(&n)),
		uint32(unsafe.Sizeof(n)),
		&bytesReturned,
		nil,
	)
	if err != nil {
		return 0, err
	}

	return n.DeviceNumber, nil
}

// GetDiskKernelObjectPath returns the (symlink'd) kernel object name for the DR.
// e.g. `\\.\PhysicalDrive0`
func GetDiskKernelObjectPath(diskDevicePath string) (string, error) {
	n, err := GetDiskNumber(diskDevicePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("\\\\.\\PhysicalDrive%d", n), nil
}
