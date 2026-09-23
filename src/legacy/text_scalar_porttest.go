//go:build porttest

package legacy

func PortTestTextCompareWide(a, b *uint16) int  { return textCompareWide(a, b) }
func PortTestTextCompareNarrow(a, b *byte) int  { return textCompareNarrow(a, b) }
func PortTestTextDecimal(p *uint16) int32       { return textDecimal(p) }
func PortTestTextCopy(dst, src *uint16) *uint16 { return textCopy(dst, src) }
