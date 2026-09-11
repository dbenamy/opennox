//go:build porttest

package client

import "github.com/opennox/opennox/v1/legacy/common/alloc"

// PortTestGlyphClient creates only the lookup state used by the legacy Glyph
// cache. It deliberately avoids client/server initialization and assets.
func PortTestGlyphClient(index uint32) (*Client, func()) {
	c := new(Client)
	if index == 0 {
		return c, func() {}
	}
	typ, free := alloc.New(ObjectType{})
	typ.Field_1c = int32(index)
	c.Things.byID = map[string]*ObjectType{"glyph": typ}
	return c, free
}
