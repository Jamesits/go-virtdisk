package virtdisk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
	"path/filepath"
	"testing"
	"unsafe"
)

func TestCreateVhd(t *testing.T) {
	var ret1 uintptr
	var err error

	// CreateVirtualDisk with parameters V1
	vsType := VirtualStorageType{
		DeviceId: VirtualStorageTypeDeviceVhd,
		VendorId: VirtualStorageTypeVendorMicrosoft,
	}
	path, _ := windows.UTF16PtrFromString(filepath.Join(temporaryDirectory, "test.vhd"))
	param := CreateVirtualDiskParametersV1{
		Version:     Version{Version: 1},
		MaximumSize: 67108864,
	}
	handle := intPtrZero
	ret1, _, err = virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&vsType)),   // VirtualStorageType
		uintptr(unsafe.Pointer(path)),      // Path
		uintptr(VirtualDiskAccessCreate),   // VirtualDiskAccessMask
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
	assert.NotEqualValues(t, 0, handle)

	// CloseHandle
	err = windows.CloseHandle(windows.Handle(handle))
	assert.NoError(t, err)
}
