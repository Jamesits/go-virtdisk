package devices

import (
	"fmt"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TODO: make tests reproducible

func TestListDrives(t *testing.T) {
	devices, err := ListDrives()
	assert.NoError(t, err)

	for _, device := range devices {
		fmt.Printf("%s\n", device)

		id, err := GetStorageDeviceNumber(device)
		assert.NoError(t, err)

		drive := types.DriveFromId(id)
		fmt.Printf("\tDrive: %s\n", drive)
	}
}
