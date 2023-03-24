package disk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func testMkdirAllWrapper(t *testing.T, path string) {
	var err error
	err = Mkdir(path)
	assert.NoError(t, err)
	err = Rmdir(path)
	assert.NoError(t, err)
}

func TestMkdirAll(t *testing.T) {
	testMkdirAllWrapper(t, "\\\\.\\C:\\Windows\\Temp\\go-virtdisk\\114514")
	testMkdirAllWrapper(t, "C:\\Windows\\Temp\\go-virtdisk\\114514")
}
