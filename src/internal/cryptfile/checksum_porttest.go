//go:build porttest

package cryptfile

// PortTestChecksum exposes stream state for serialization compatibility captures.
// Keeping call boundaries is observable even when payload bytes are unchanged.
func (f *CryptFile) PortTestChecksum() uint32 { return f.crcSum }
