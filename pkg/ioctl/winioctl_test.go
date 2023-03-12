package ioctl

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func assertSameSize[Expected any, Actual any](t *testing.T) {
	assert.EqualValues(t, unsafe.Sizeof(*new(Expected)), unsafe.Sizeof(*new(Actual)))
}

func TestUnionPadding(t *testing.T) {
	assertSameSize[CreateDiskGpt, CreateDiskMbr](t)
	assertSameSize[DriveLayoutInformationGpt, DriveLayoutInformationMbr](t)
	assertSameSize[PartitionInformationGpt, PartitionInformationMbr](t)
}
