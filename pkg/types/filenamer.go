package types

type FileNamer interface {
	// AsFileName returns a pointer to a UTF-16 string which can be used as the first argument to CreateFileW.
	AsFileName() (*uint16, error)
}
