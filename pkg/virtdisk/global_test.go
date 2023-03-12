package virtdisk

import (
	"github.com/jamesits/go-virtdisk/pkg/privilege"
	"github.com/jamesits/libiferr/exception"
	"log"
	"os"
	"testing"
)

const intPtrZero = uintptr(0)

var temporaryDirectory string
var virtdisk Virtdisk

func TestMain(m *testing.M) {
	var err error

	// elevate
	// https://stackoverflow.com/questions/24396644/programmatically-mount-a-microsoft-virtual-hard-drive-vhd
	// https://social.msdn.microsoft.com/Forums/en-US/415436ce-4f5e-4808-9d46-f9614d0c2e44/a-privilege-problem-on-win7-about-semanagevolumename?forum=vcgeneral
	err = privilege.EnablePrivilege("SeManageVolumePrivilege")
	exception.HardFailWithReason("unable to enable SeManageVolumePrivilege", err)

	// create temporary directory
	temporaryDirectory, err = os.MkdirTemp(os.Getenv("TEMP"), "go-virtdisk.*")
	exception.HardFailWithReason("unable to create temporary directory", err)
	log.Printf("temporary directory: %s\n", temporaryDirectory)

	// load DLLs
	virtdisk = Virtdisk{}
	err = virtdisk.Unmarshal("virtdisk.dll")
	exception.HardFailWithReason("unable to load virtdisk.dll", err)

	ret := m.Run()

	err = os.RemoveAll(temporaryDirectory)
	if err != nil {
		log.Printf("unable to clean up temporary directory: %s %v", temporaryDirectory, err)
	}

	os.Exit(ret)
}
