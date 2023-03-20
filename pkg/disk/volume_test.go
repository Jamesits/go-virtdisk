package disk

import (
	"fmt"
	"github.com/stretchr/testify/assert"
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
	}
}
