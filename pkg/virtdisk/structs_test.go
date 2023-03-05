package virtdisk

import (
	"github.com/jamesits/nopadding/analyzer_wrapper"
	"github.com/jamesits/nopadding/padding"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
	"testing"
)

func TestStructAlignment(t *testing.T) {
	assert.Zero(t, analyzer_wrapper.Run(
		[]string{},
		[]*analysis.Analyzer{padding.Analyzer},
	))
}
