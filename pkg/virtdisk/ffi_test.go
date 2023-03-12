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

var TemporaryVhdxPath = filepath.Join(".", "test.vhdx")

func TestCreateVirtualDisk(t *testing.T) {
	d := Virtdisk{}
	assert.NoError(t, d.Unmarshal("virtdisk.dll"))

	// CreateVirtualDisk
	vsType := VirtualStorageType{
		DeviceId: VirtualStorageTypeDeviceVhdx,
		VendorId: VirtualStorageTypeVendorMicrosoft,
	}
	path, _ := windows.UTF16PtrFromString(TemporaryVhdxPath)
	param := CreateVirtualDiskParametersV1{
		Version:           Version{Version: 1},
		UniqueId:          uuid.Nil,
		MaximumSize:       67108864,
		BlockSizeInBytes:  0,
		SectorSizeInBytes: 0,
		ParentPath:        nil,
		SourcePath:        nil,
	}
	handle := uintptr(0)
	ret, _, err := d.CreateVirtualDisk.Call(
		uintptr(unsafe.Pointer(&vsType)),   // VirtualStorageType
		uintptr(unsafe.Pointer(path)),      // Path
		uintptr(VirtualDiskAccessCreate),   // VirtualDiskAccessMask
		uintptr(0),                         // SecurityDescriptor
		uintptr(CreateVirtualDiskFlagNone), // Flags
		uintptr(0),                         // ProviderSpecificFlags
		uintptr(unsafe.Pointer(&param)),    // Parameters
		uintptr(0),                         // Overlapped
		uintptr(unsafe.Pointer(&handle)),   // Handle
	)
	fmt.Printf("handle = %d\n", handle)
	assert.ErrorIs(t, err, windows.ERROR_SUCCESS)
	assert.Zero(t, ret) // ret should be the same as the error code
	assert.NotEqualValues(t, 0, handle)
}
