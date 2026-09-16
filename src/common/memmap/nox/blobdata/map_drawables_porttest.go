//go:build porttest

package blobdata

func PortTestMapDrawableTables() []PortTestClientEffectsTable {
	return []PortTestClientEffectsTable{
		{Base: 0x581450, Offset: 9736, Data: append([]byte(nil), data581450[9736:9760]...)},
		{Base: 0x587000, Offset: 196184, Data: append([]byte(nil), data587000[196184:196440]...)},
	}
}
