package files

import (
	"github.com/jamesits/go-virtdisk/pkg/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testMkdirAllWrapper(t *testing.T, path types.Path) {
	var err error
	err = Mkdir(path)
	assert.NoError(t, err)
	err = Rmdir(path)
	assert.NoError(t, err)
}

func TestMkdirAll(t *testing.T) {
	testMkdirAllWrapper(t, "\\\\.\\C:\\Windows\\Temp\\go-virtdisks\\114514")
	testMkdirAllWrapper(t, "C:\\Windows\\Temp\\go-virtdisks\\114514")
}
