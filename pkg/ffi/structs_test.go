package ffi

import (
	"github.com/jamesits/libiferr/testing_utils"
	"testing"
)

func TestUnionPadding(t *testing.T) {
	testing_utils.AssertSameSize[CreateDiskGpt, CreateDiskMbr](t)
	testing_utils.AssertSameSize[DriveLayoutInformationGpt, DriveLayoutInformationMbr](t)
	testing_utils.AssertSameSize[PartitionInformationGpt, PartitionInformationMbr](t)
}
