package ffi

var Setupapi setupapi
var Kernel32 kernel32
var Virtdisk virtdisk

func init() {
	var err error

	err = Setupapi.Unmarshal("setupapi.dll")
	if err != nil {
		panic(err)
	}

	err = Kernel32.Unmarshal("kernel32.dll")
	if err != nil {
		panic(err)
	}

	err = Virtdisk.Unmarshal("virtdisk.dll")
	if err != nil {
		panic(err)
	}
}
