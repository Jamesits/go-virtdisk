package devices

import (
	"errors"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
	"unsafe"
)

// ListDrives returns the devices path of the drives.
func ListDrives() (ret []types.Device, err error) {
	// https://stackoverflow.com/a/18183115
	handle, err := windows.SetupDiGetClassDevsEx(ffi.GuidDevInterfaceDisk, "", uintptr(0), windows.DIGCF_PRESENT|windows.DIGCF_DEVICEINTERFACE, windows.DevInfo(0), "")
	if windows.Handle(handle) == windows.InvalidHandle || err != nil {
		return nil, err
	}
	defer ffi.Setupapi.SetupDiDestroyDeviceInfoList.Call(uintptr(handle))

	deviceIndex := uintptr(0)
	devInterfaceData := ffi.SPDeviceInterfaceData{}
	devInterfaceData.Size = uint32(unsafe.Sizeof(devInterfaceData))
	for {
		successful, _, _ := ffi.Setupapi.SetupDiEnumDeviceInterfaces.Call(
			uintptr(handle),
			uintptr(0),
			uintptr(unsafe.Pointer(ffi.GuidDevInterfaceDisk)),
			deviceIndex,
			uintptr(unsafe.Pointer(&devInterfaceData)),
		)
		if successful == 0 {
			break
		}

		var s uintptr
		_, _, err := ffi.Setupapi.SetupDiGetDeviceInterfaceDetailW.Call(
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

		b := make([]byte, s)
		// an additional WCHAR[1] must be calculated in the initial size so that the function does not complain about
		// windows.ERROR_INVALID_USER_BUFFER (The supplied user buffer is not valid for the requested operation.)
		var interfaceDetailData ffi.SPDeviceInterfaceDetailDataW
		interfaceDetailData.Size = uint32(unsafe.Sizeof(interfaceDetailData))
		_, _ = bytebuilder.Copy(b, &interfaceDetailData)

		successful, _, err = ffi.Setupapi.SetupDiGetDeviceInterfaceDetailW.Call(
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

		_, pathBuffer := bytebuilder.CarCdr[ffi.SPDeviceInterfaceDetailDataH](b)
		path := windows.UTF16ToString(bytebuilder.SliceCast[uint8, uint16](pathBuffer)) // string is NUL terminated
		ret = append(ret, types.Device(path))

		deviceIndex++
	}

	return
}

// GetStorageDeviceNumber returns the PhysicalDrive number of the hard drives devices.
func GetStorageDeviceNumber(device types.Device) (uint32, error) {
	dp, err := device.AsFileNameW()
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
	defer windows.CloseHandle(dHandle)

	var n ffi.StorageDeviceNumber
	var bytesReturned uint32
	err = windows.DeviceIoControl(
		dHandle,
		ffi.IoctlStorageGetDeviceNumber,
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
