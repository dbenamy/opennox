//go:build porttest

package blobdata

// PortTestScoreboardTables returns the original literal formats and color tables.
// The fixture separately owns the three relocated class-name input pointers.
func PortTestScoreboardTables() []PortTestClientEffectsTable {
	return []PortTestClientEffectsTable{{Base: 0x587000, Offset: 4704, Data: append([]byte(nil), data587000[4704:4728]...)}, {Base: 0x587000, Offset: 145584, Data: append([]byte(nil), data587000[145584:147904]...)}}
}
