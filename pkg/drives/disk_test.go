package drives

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDisk(t *testing.T) {
	disks, err := List()
	assert.NoError(t, err)
	for _, disk := range disks {
		fmt.Printf("%s\n", disk)

		n, err := GetStorageDeviceNumber(disk)
		assert.NoError(t, err)
		fmt.Printf("\tID: %d\n", n)

		kp, err := FromDevice(disk)
		assert.NoError(t, err)
		fmt.Printf("\tObject path: %s\n", kp)

		serial, err := GetSerial(disk)
		assert.NoError(t, err)
		fmt.Printf("\tSerial: %s\n", serial)
	}
}
