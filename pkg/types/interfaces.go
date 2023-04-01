package types

type File interface {
	// AsFileName returns a pointer to a UTF-16 string which can be used as the first argument to CreateFileW.
	AsFileName() (*uint16, error)
}

type Directory interface {
	// AsDirectory returns a pointer to a UTF-16 string which contains the directory path representation of this object.
	AsDirectory() (*uint16, error)
}

type KernelObject interface {
	// AsObjectPath returns a pointer to a UTF-16 string which is the object's kernel path.
	AsObjectPath() (*uint16, error)
}
