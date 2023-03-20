package disk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDisk(t *testing.T) {
	disks, err := GetDisks()
	assert.NoError(t, err)
	for _, disk := range disks {
		fmt.Printf("%s\n", disk)

		n, err := GetDiskNumber(disk)
		assert.NoError(t, err)
		fmt.Printf("\tID: %d\n", n)

		kp, err := GetDiskKernelObjectPath(disk)
		assert.NoError(t, err)
		fmt.Printf("\tObject path: %s\n", kp)
	}
}
