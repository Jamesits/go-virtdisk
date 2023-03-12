package virtdisk

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jamesits/go-virtdisk/pkg/ioctl"
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
		uintptr(unsafe.Pointer(&handle)),   // Handle
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
	ioctlDiskCreateDisk := ioctl.CtlCode(ioctl.FileDeviceDisk, 0x0016, ioctl.MethodBuffered, ioctl.FileReadAccess|ioctl.FileWriteAccess)
	gpt := ioctl.CreateDiskGpt{
		PartitionStyle:    ioctl.PartitionStyleGPT,
		DiskId:            uuid.New(),
		MaxPartitionCount: 128,
	}
	err = windows.DeviceIoControl(
		diskHandle,
		ioctlDiskCreateDisk,
		(*byte)(unsafe.Pointer(&gpt)),
		uint32(unsafe.Sizeof(gpt)),
		nil,
		0,
		nil,
		nil,
	)
	assert.NoError(t, err)
	//mbr := ioctl.CreateDiskMbr{
	//	PartitionStyle: ioctl.PartitionStyleMBR,
	//	Signature:      1,
	//}
	//err = windows.DeviceIoControl(
	//	diskHandle,
	//	ioctlDiskCreateDisk,
	//	(*byte)(unsafe.Pointer(&mbr)),
	//	uint32(unsafe.Sizeof(mbr)),
	//	nil,
	//	0,
	//	nil,
	//	nil,
	//)
	//assert.NoError(t, err)

	// DetachVirtualDisk
	ret1, _, err = virtdisk.DetachVirtualDisk.Call(
		handle,
		uintptr(DetachVirtualDiskFlagNone),
		intPtrZero,
	)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret1)
}
