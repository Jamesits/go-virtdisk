package paths

import (
	"fmt"
	"github.com/jamesits/go-virtdisk/pkg/types"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func testMkdirAllWrapper(t *testing.T, path types.Path) {
	var err error

	fmt.Printf("Path: %s\n", path)
	err = Mkdir(path)
	assert.NoError(t, err)
	err = Rmdir(path)
	assert.NoError(t, err)
}

func TestMkdirAll(t *testing.T) {
	tmp := os.Getenv("TEMP")
	testMkdirAllWrapper(t, types.Path("\\\\.\\"+filepath.Join(tmp, "go-virtdisk", "114514")))
	testMkdirAllWrapper(t, types.Path(filepath.Join(tmp, "go-virtdisk", "114514")))
}
