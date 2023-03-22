package virtdisk

import (
	"fmt"
	"github.com/jamesits/go-virtdisk/pkg/disk"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDisk(t *testing.T) {
	disks, err := disk.GetDisks()
	assert.NoError(t, err)
	for _, d := range disks {
		fmt.Printf("%s\n", d)

		kp, err := disk.GetDiskKernelObjectPath(d)
		assert.NoError(t, err)
		fmt.Printf("\tPath: %s\n", kp)

		vhds, err := GetVirtualDiskBackingFiles(kp)
		assert.NoError(t, err)
		for _, vhd := range vhds {
			fmt.Printf("\tVHD parent: %s\n", vhd)
		}
	}
}
