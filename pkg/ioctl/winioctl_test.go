package ioctl

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func assertSameSize[T1 any, T2 any](t *testing.T) {
	assert.EqualValues(t, unsafe.Sizeof(new(T1)), unsafe.Sizeof(new(T2)))
}

func TestUnionPadding(t *testing.T) {
	assertSameSize[CreateDiskMbr, CreateDiskGpt](t)
	assertSameSize[DriveLayoutInformationMbr, DriveLayoutInformationGpt](t)
	assertSameSize[PartitionInformationMbr, PartitionInformationGpt](t)
}
