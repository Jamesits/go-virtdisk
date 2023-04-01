package virtdisks

import (
	"fmt"
	"github.com/jamesits/go-virtdisk/pkg/drives"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDisk(t *testing.T) {
	disks, err := drives.List()
	assert.NoError(t, err)
	for _, disk := range disks {
		vhds, err := GetVirtualDiskBackingFiles(disk)
		assert.NoError(t, err)
		for _, vhd := range vhds {
			fmt.Printf("\tVHD parent: %s\n", vhd)
		}
	}
}
