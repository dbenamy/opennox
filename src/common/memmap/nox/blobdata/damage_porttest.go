//go:build porttest

package blobdata

// PortTestDamageTable returns the shipped damage-name table, strings, and
// conductivity constant. The fixture relocates the eighteen name pointers.
func PortTestDamageTable() []byte { return append([]byte(nil), data587000[200728:201112]...) }
