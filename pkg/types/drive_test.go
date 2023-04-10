package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDriveId(t *testing.T) {
	for i := uint32(0); i < 100; i++ {
		d := DriveFromId(i)
		id, err := d.AsId()
		assert.NoError(t, err)
		assert.EqualValues(t, i, id)
	}
}
