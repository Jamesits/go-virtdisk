package types

// Device
// The returned string is an object under `\GLOBAL??`. Examples:
// - VHDX: `\\?\scsi#drives&ven_msft&prod_virtual_disk#2&1f4adffe&0&000001#{53f56307-b6bf-11d0-94f2-00a0c91efb8b}`
// - NVMe drives: `\\?\scsi#drives&ven_msft&prod_virtual_disk#2&1f4adffe&0&000001#{53f56307-b6bf-11d0-94f2-00a0c91efb8b}`
type Device Path

func (d Device) AsFileName() (*uint16, error) {
	return Path(d).asUTF16Ptr()
}

func (d Device) AsObjectPath() (*uint16, error) {
	return Path(d).asUTF16Ptr()
}
