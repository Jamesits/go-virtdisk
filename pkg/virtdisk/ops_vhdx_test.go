package virtdisk

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
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
	handle := uintptr(0)
	ret1, _, err = d.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&vsType)),   // VirtualStorageType
		uintptr(unsafe.Pointer(path)),      // Path
		uintptr(VirtualDiskAccessNone),     // VirtualDiskAccessMask (must be none if using struct v2)
		uintptr(0),                         // SecurityDescriptor
		uintptr(CreateVirtualDiskFlagNone), // Flags
		uintptr(0),                         // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),    // Parameters
		uintptr(0),                         // Overlapped
		uintptr(unsafe.Pointer(&handle)),   // Handle
	)
	fmt.Printf("handle = %d\n", handle)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret1) // ret should be the same as the error code
	assert.NotEqualValues(t, 0, handle)

	// AttachVirtualDisk
	ret1, _, err = d.AttachVirtualDisk.Call(
		uintptr(unsafe.Pointer(handle)), // VirtualDiskHandle
		uintptr(0),                      // SecurityDescriptor
		uintptr(AttachVirtualDiskFlagNoDriveLetter|AttachVirtualDiskFlagPermanentLifetime), // Flags
		uintptr(0), // ProviderSpecificFlags
		uintptr(0), // Parameters
		uintptr(0), // Overlapped
	)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret1)

	// DetachVirtualDisk
	ret1, _, err = d.DetachVirtualDisk.Call(
		uintptr(unsafe.Pointer(handle)),
		uintptr(DetachVirtualDiskFlagNone),
		uintptr(0),
	)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret1)

	// CloseHandle
	err = windows.CloseHandle(windows.Handle(handle))
	assert.NoError(t, err)
}
