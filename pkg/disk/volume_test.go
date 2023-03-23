package disk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
	"testing"
)

func TestGetVolumes(t *testing.T) {
	v, err := GetVolumes()
	assert.NoError(t, err)
	assert.Less(t, 0, len(v))

	fmt.Printf("Volumes:\n")
	for _, volPath := range v {
		fmt.Printf("\t%s\n", volPath)

		mps, err := GetVolumeMountPoints(volPath)
		assert.NoError(t, err)
		for _, mp := range mps {
			fmt.Printf("\t\t%s\n", mp)
		}

		disks, err := GetVolumeBackingDrives(volPath)
		assert.NoError(t, err)
		for _, disk := range disks {
			fmt.Printf("\t\tdep=%s\n", disk)
		}

		serial, err := GetVolumeSerial(volPath)
		assert.NoError(t, err)
		fmt.Printf("\t\tserial=%s\n", serial)

		label, err := GetVolumeLabel(volPath)
		assert.NoError(t, err)
		fmt.Printf("\t\tlabel=%s\n", label)
	}
}

func TestFilePathConversion(t *testing.T) {
	path := "\\\\.\\C:\\Windows\\System32\\notepad.exe"
	mp, rel, err := GetMountPointByFileName(path)
	assert.NoError(t, err)
	fmt.Printf("mp=%s\n", mp)
	fmt.Printf("rel=%s\n", rel)

	vol, err := GetVolumeGUIDPathByMountPoint(mp)
	assert.NoError(t, err)
	fmt.Printf("vol=%s\n", vol)

	vol, err = GetVolumeGUIDPathByMountPoint(path)
	assert.ErrorIs(t, err, windows.ERROR_INVALID_NAME)
}
