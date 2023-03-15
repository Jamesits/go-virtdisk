package virtdisk

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jamesits/go-bytebuilder"
	"github.com/jamesits/go-virtdisk/pkg/disk"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
	"log"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestCreateVhdx(t *testing.T) {
	var ret1 uintptr
	var err error

	// CreateVirtualDisk with parameters V1
	vsType := VirtualStorageType{
		DeviceId: VirtualStorageTypeDeviceVhdx,
		VendorId: VirtualStorageTypeVendorMicrosoft,
	}
	path, _ := windows.UTF16PtrFromString(filepath.Join(temporaryDirectory, "test.vhdx"))
	param := CreateVirtualDiskParametersV2{
		Version:     Version{Version: 2},
		UniqueId:    uuid.Nil,
		MaximumSize: 67108864,
	}
	handle := intPtrZero
	ret1, _, err = virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&vsType)),   // VirtualStorageType
		uintptr(unsafe.Pointer(path)),      // Path
		uintptr(VirtualDiskAccessNone),     // VirtualDiskAccessMask (must be none if using struct v2)
		intPtrZero,                         // SecurityDescriptor
		uintptr(CreateVirtualDiskFlagNone), // Flags
		intPtrZero,                         // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),    // Parameters
		intPtrZero,                         // Overlapped
		uintptr(unsafe.Pointer(&handle)),   // handle
	)
	fmt.Printf("handle = %v\n", handle)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret1) // ret should be the same as the error code
	assert.NotEqualValues(t, windows.InvalidHandle, handle)
	defer func() {
		// CloseHandle
		err = windows.CloseHandle(windows.Handle(handle))
		assert.NoError(t, err)
	}()

	// AttachVirtualDisk
	ret1, _, err = virtdisk.AttachVirtualDisk.Call(
		handle,     // VirtualDiskHandle
		intPtrZero, // SecurityDescriptor
		uintptr(AttachVirtualDiskFlagNoDriveLetter|AttachVirtualDiskFlagPermanentLifetime), // Flags
		intPtrZero, // ProviderSpecificFlags
		intPtrZero, // Parameters
		intPtrZero, // Overlapped
	)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret1)

	// GetVirtualDiskPhysicalPath: get the disk path of the mounted VHDX (usually "\\.\PhysicalDriveN")
	// https://social.msdn.microsoft.com/Forums/vstudio/en-US/1d429820-ce9b-4741-aa43-f0b8f85f8cb7/mounting-a-vhd?forum=vcgeneral
	virtualDiskPhysicalPathSize := uint32(0)
	_, _, err = virtdisk.GetVirtualDiskPhysicalPath.Call(
		handle,
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathSize)),
		intPtrZero,
	)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.NotZero(t, virtualDiskPhysicalPathSize)
	virtualDiskPhysicalPathUtf16 := make([]uint16, virtualDiskPhysicalPathSize)
	_, _, err = virtdisk.GetVirtualDiskPhysicalPath.Call(
		handle,
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathSize)),
		uintptr(unsafe.Pointer(&virtualDiskPhysicalPathUtf16[0])),
	)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	virtualDiskPhysicalPath := windows.UTF16ToString(virtualDiskPhysicalPathUtf16)
	log.Printf("physical path: %s\n", virtualDiskPhysicalPath)

	// get a handle to the disk
	diskHandle, err := windows.CreateFile(&virtualDiskPhysicalPathUtf16[0], windows.GENERIC_READ|windows.GENERIC_WRITE, windows.FILE_SHARE_READ|windows.FILE_SHARE_WRITE, nil, windows.OPEN_EXISTING, 0, windows.Handle(0))
	assert.NoError(t, err)
	assert.NotEqualValues(t, windows.InvalidHandle, diskHandle)
	defer func() {
		// close the handle to the disk
		err = windows.CloseHandle(diskHandle)
		assert.NoError(t, err)
	}()

	// initialize the disk
	// https://www.codeproject.com/script/Content/ViewAssociatedFile.aspx?rzp=%2FKB%2Fwinsdk%2FHard_drive_Information%2Fsmartsrc.zip&zep=SMART%2FDDKInclude%2Fntdddisk.h&obid=16671&obtid=2&ovid=1
	ioctlDiskCreateDisk := disk.CtlCode(disk.FileDeviceDisk, 0x0016, disk.MethodBuffered, disk.FileReadAccess|disk.FileWriteAccess)
	//createDisk := disk.CreateDiskMbr{
	//	PartitionStyle: disk.PartitionStyleMBR,
	//	Signature:      1,
	//}
	createDisk := disk.CreateDiskGpt{
		PartitionStyle:    disk.PartitionStyleGPT,
		DiskId:            uuid.Nil, // a random UUID will result in "The request is not supported."
		MaxPartitionCount: 128,
	}
	err = windows.DeviceIoControl(
		diskHandle,
		ioctlDiskCreateDisk,
		(*byte)(unsafe.Pointer(&createDisk)),
		uint32(unsafe.Sizeof(createDisk)),
		nil,
		0,
		nil,
		nil,
	)
	assert.NoError(t, err)

	// FIXME: wait for the MSR partition arrival??? HOW???
	// https://learn.microsoft.com/en-us/windows/win32/api/winioctl/ni-winioctl-ioctl_disk_create_disk#remarks
	// https://github.com/pbatard/rufus/blob/ca84a4f6c5f24891ffbe2834648bff0120bfc4e3/src/drive.c#L2542

	// Partition the disk
	partitionInformationEx := &bytebuilder.ByteBuilder{}
	_, _ = partitionInformationEx.WriteObject(disk.DriveLayoutInformationExGpt{
		DriveLayoutInformationEx: disk.DriveLayoutInformationEx{
			PartitionStyle: 1,
			PartitionCount: 1,
		},
		DriveLayoutInformationGpt: disk.DriveLayoutInformationGpt{
			DiskId:               uuid.Nil,
			StartingUsableOffset: 0,
			UsableLength:         3342336,
			MaxPartitionCount:    128,
		},
	})
	_, _ = partitionInformationEx.WriteObject(disk.PartitionInformationExGpt{
		PartitionInformationEx: disk.PartitionInformationEx{
			PartitionStyle:     1,
			StartingOffset:     1048576,
			PartitionLength:    3342336 - 1048576,
			PartitionNumber:    1,
			RewritePartition:   true,
			IsServicePartition: false,
		},
		PartitionInformationGpt: disk.PartitionInformationGpt{
			PartitionType: uuid.MustParse("EBD0A0A2-B9E5-4433-87C0-68B6B72699C7"),
			PartitionId:   uuid.Must(uuid.NewRandom()),
			Attributes:    0,
			Name:          [36]uint16{'t', 'e', 's', 't'},
		},
	})
	ioctlDiskSetDriveLayoutEx := disk.CtlCode(disk.FileDeviceDisk, 0x0015, disk.MethodBuffered, disk.FileReadAccess|disk.FileWriteAccess)
	b := partitionInformationEx.Bytes()
	err = windows.DeviceIoControl(
		diskHandle,
		ioctlDiskSetDriveLayoutEx,
		&b[0],
		uint32(len(b)),
		nil,
		0,
		nil,
		nil,
	)
	assert.NoError(t, err)

	// DetachVirtualDisk
	ret1, _, err = virtdisk.DetachVirtualDisk.Call(
		handle,
		uintptr(DetachVirtualDiskFlagNone),
		intPtrZero,
	)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret1)
}
