package virtdisk

import (
	"github.com/jamesits/go-virtdisk/pkg/struct_alignment"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStructAlignment(t *testing.T) {
	assert.Zero(t, struct_alignment.Run([]string{}))
}
