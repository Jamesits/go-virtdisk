package virtdisks

import (
	"fmt"
	"github.com/jamesits/go-virtdisk/pkg/ffi"
	"github.com/jamesits/go-virtdisk/pkg/types"
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
	vsType := ffi.VirtualStorageType{
		DeviceId: VirtualStorageTypeDeviceVhd,
		VendorId: VirtualStorageTypeVendorMicrosoft,
	}
	path, _ := windows.UTF16PtrFromString(filepath.Join(temporaryDirectory, "test.vhd"))
	param := ffi.CreateVirtualDiskParametersV1{
		Version:     ffi.Version{Version: 1},
		MaximumSize: 67108864,
	}
	handle := types.IntPtrZero
	ret1, _, err = ffi.Virtdisk.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&vsType)),   // VirtualStorageType
		uintptr(unsafe.Pointer(path)),      // Path
		uintptr(VirtualDiskAccessCreate),   // VirtualDiskAccessMask
		types.IntPtrZero,                   // SecurityDescriptor
		uintptr(CreateVirtualDiskFlagNone), // Flags
		types.IntPtrZero,                   // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),    // Parameters
		types.IntPtrZero,                   // Overlapped
		uintptr(unsafe.Pointer(&handle)),   // handle
	)
	fmt.Printf("handle = %v\n", handle)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret1) // ret should be the same as the error code
	assert.NotEqualValues(t, 0, handle)

	// CloseHandle
	err = windows.CloseHandle(windows.Handle(handle))
	assert.NoError(t, err)
}
