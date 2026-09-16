//go:build porttest

package blobdata

// PortTestBookGuideFamily returns the actual startup guide family and pointer
// table. Fixtures must rebase its one pointer after installing these bytes.
func PortTestBookGuideFamily() []byte {
	return append([]byte(nil), data587000[132100:132132]...)
}

func PortTestBookQuickbarSpacing() []byte { return append([]byte(nil), data587000[133488:133508]...) }
