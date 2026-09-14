//go:build porttest

package blobdata

func PortTestMapThemeTableData() []byte {
	return append([]byte(nil), data587000[253144:253512]...)
}
