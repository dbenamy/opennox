//go:build porttest

package blobdata

// PortTestObjectStateLoot returns the shipped barrel/crate tables and strings.
// The fixture relocates their name pointers as runtime Init does.
func PortTestObjectStateLoot() []byte { return append([]byte(nil), data587000[203080:203724]...) }
