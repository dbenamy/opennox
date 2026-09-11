//go:build porttest

package blobdata

// PortTestEffectsInventoryTable returns the original six descriptor lookup rows.
// The caller relocates the callback addresses as Init does for the live table.
func PortTestEffectsInventoryTable() []byte {
	return append([]byte(nil), data587000[200160:200280]...)
}
