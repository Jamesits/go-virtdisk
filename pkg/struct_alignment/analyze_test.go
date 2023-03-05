package struct_alignment

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

type SelfPackageReflectionStub struct{}

func SelfPackagePath() string {
	// https://stackoverflow.com/a/25263604
	return reflect.TypeOf(SelfPackageReflectionStub{}).PkgPath()
}

func TestStructAlignmentAnalyzerSelf(t *testing.T) {
	// test ourselves (empty string means self package)
	assert.NotZero(t, Run([]string{}))
}

func TestStructAlignmentAnalyzerVictims(t *testing.T) {
	// test victims
	assert.Zero(t, Run([]string{strings.Join([]string{SelfPackagePath(), "internal", "testdata", "pass"}, "/")}))
	assert.NotZero(t, Run([]string{strings.Join([]string{SelfPackagePath(), "internal", "testdata", "fail"}, "/")}))
}
