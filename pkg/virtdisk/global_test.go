package virtdisk

import (
	"github.com/jamesits/libiferr/exception"
	"log"
	"os"
	"testing"
)

const intPtrZero = uintptr(0)

var temporaryDirectory string
var d Virtdisk

func TestMain(m *testing.M) {
	var err error

	// create temporary directory
	temporaryDirectory, err = os.MkdirTemp(os.Getenv("TEMP"), "go-virtdisk.*")
	exception.HardFailWithReason("unable to create temporary directory", err)
	log.Printf("temporary directory: %s\n", temporaryDirectory)

	// load DLL
	d = Virtdisk{}
	err = d.Unmarshal("virtdisk.dll")
	exception.HardFailWithReason("unable to load DLL", err)

	ret := m.Run()

	err = os.RemoveAll(temporaryDirectory)
	if err != nil {
		log.Printf("unable to clean up temporary directory: %s %v", temporaryDirectory, err)
	}

	os.Exit(ret)
}
