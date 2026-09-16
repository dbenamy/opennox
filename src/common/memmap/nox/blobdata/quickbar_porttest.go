//go:build porttest

package blobdata

// PortTestQuickbarConstants copies the production startup ability records and
// the image-name suffix read by the constructor. These contain no pointers.
func PortTestQuickbarConstants() []struct {
	Offset uintptr
	Data   []byte
} {
	return []struct {
		Offset uintptr
		Data   []byte
	}{
		{133536, append([]byte(nil), data587000[133536:133656]...)},
		{134996, append([]byte(nil), data587000[134996:135000]...)},
	}
}
