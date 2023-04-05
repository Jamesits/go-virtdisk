package types

type File interface {
	// AsFileNameW returns a pointer to a UTF-16 string which can be used as the first argument to CreateFileW.
	AsFileNameW() (*uint16, error)
}

type Directory interface {
	// AsDirectoryW returns a pointer to a UTF-16 string which contains the directory path representation of this object.
	AsDirectoryW() (*uint16, error)
}

type KernelObject interface {
	// AsObjectPathW returns a pointer to a UTF-16 string which is the object's kernel path.
	AsObjectPathW() (*uint16, error)
}
