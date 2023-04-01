package drives

import (
	"bytes"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/devices"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"golang.org/x/sys/windows"
	"unsafe"
)

func List() ([]types.Drive, error) {
	devices, err := devices.ListDrives()
	if err != nil {
		return nil, err
	}

	var ret []types.Drive
	for _, dev := range devices {
		drive, err := FromDevice(dev)
		if err != nil {
			return nil, err
		}
		ret = append(ret, drive)
	}

	return ret, nil
}

// FromDevice returns the (symlink'd) kernel object name for the DR.
// e.g. `\\.\PhysicalDrive0`
func FromDevice(device types.Device) (types.Drive, error) {
	n, err := devices.GetStorageDeviceNumber(device)
	if err != nil {
		return "", err
	}

	return types.DriveFromId(n), nil
}

func FromVolume(volume types.Volume) (drives []types.Drive, err error) {
	// https://stackoverflow.com/questions/29212597/how-to-enumerate-disk-volume-names
	vp, err := volume.AsFileName()
	if err != nil {
		return nil, err
	}

	vHandle, err := windows.CreateFile(
		vp,
		// Even if GENERIC_READ is not specified, metadata can still be read; GENERIC_READ requires administrator privileges.
		// https://stackoverflow.com/questions/327718/how-to-list-physical-disks#comment89593842_11683906
		0,
		windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		windows.Handle(0),
	)
	if err != nil || vHandle == windows.InvalidHandle {
		return nil, err
	}
	defer windows.CloseHandle(vHandle)

	// https://stackoverflow.com/a/5664841
	extents := ffi.VolumeDiskExtents{}
	bytesReturned := uint32(unsafe.Sizeof(extents))

	for i := 0; i < 2; i++ {
		b := make([]byte, bytesReturned)
		err = windows.DeviceIoControl(
			vHandle,
			ffi.IoctlVolumeGetVolumeDiskExtents,
			nil,
			0,
			&b[0],
			bytesReturned,
			&bytesReturned,
			nil,
		)
		if err != nil {
			return nil, err
		}
		_ = bytebuilder.Unmarshal(b, &extents)
		if bytesReturned > uint32(unsafe.Sizeof(extents)) {
			continue
		}

		// parse extents
		reader := bytes.NewReader(b)
		header := ffi.VolumeDiskExtentsH{}
		extent := ffi.DiskExtent{}
		_, _ = bytebuilder.ReadPartial(reader, &header)
		_, _ = bytebuilder.Skip(reader, 4)
		for j := uint32(0); j < header.NumberOfDiskExtents; j++ {
			_, _ = bytebuilder.ReadPartial(reader, &extent)
			drives = append(drives, types.DriveFromId(extent.DiskNumber))
		}

		return
	}

	return nil, types.ErrorRetryLimitExceeded
}

func GetSerial(disk types.Drive) (string, error) {
	dp, err := disk.AsFileName()
	if err != nil {
		return "", err
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
		return "", err
	}
	defer windows.CloseHandle(dHandle)

	// https://forums.codeguru.com/showthread.php?545649-WinAPI-How-to-Get-Hard-Disk-Serial-Number
	var query ffi.StoragePropertyQuery
	query.PropertyId = ffi.StorageDeviceProperty
	query.QueryType = ffi.PropertyStandardQuery
	var storageDeviceDescriptor ffi.StorageDeviceDescriptorH
	bytesReturned := uint32(unsafe.Sizeof(storageDeviceDescriptor))

	for i := 0; i < 2; i++ {
		ret := make([]byte, bytesReturned)
		err = windows.DeviceIoControl(
			dHandle,
			ffi.IoctlStorageQueryProperty,
			(*byte)(unsafe.Pointer(&query)),
			uint32(unsafe.Sizeof(query)),
			&ret[0],
			bytesReturned,
			nil,
			nil,
		)
		if err != nil {
			return "", err
		}

		_ = bytebuilder.Unmarshal(ret, &storageDeviceDescriptor)
		if storageDeviceDescriptor.Size != bytesReturned {
			bytesReturned = storageDeviceDescriptor.Size
			continue
		}

		if storageDeviceDescriptor.SerialNumberOffset == 0 {
			return "", nil
		}

		serial := windows.BytePtrToString(&ret[storageDeviceDescriptor.SerialNumberOffset])
		return serial, nil
	}

	return "", types.ErrorRetryLimitExceeded
}
