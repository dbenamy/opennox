//go:build porttest

package blobdata

// PortTestMeterTables returns the actual charge-row geometry and durability
// threshold without initializing unrelated live game state.
func PortTestMeterTables() []PortTestClientEffectsTable {
	return []PortTestClientEffectsTable{
		{Base: 0x587000, Offset: 147904, Data: append([]byte(nil), data587000[147904:148392]...)},
		{Base: 0x581450, Offset: 9608, Data: append([]byte(nil), data581450[9608:9616]...)},
		{Base: 0x581450, Offset: 9760, Data: append([]byte(nil), data581450[9760:9776]...)},
	}
}
