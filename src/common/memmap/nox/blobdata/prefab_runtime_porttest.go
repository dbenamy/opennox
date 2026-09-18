//go:build porttest

package blobdata

// PortTestPrefabRuntimeTables returns shipped descriptor bytes and path punctuation.
// Script scanning reads the count/kind fields, never descriptor name pointers.
func PortTestPrefabRuntimeTables() map[uintptr][]byte {
	return map[uintptr][]byte{
		218632: append([]byte(nil), data587000[218632:228548]...),
		229828: append([]byte(nil), data587000[229828:229845]...),
		229976: append([]byte(nil), data587000[229976:229978]...),
	}
}
