//go:build porttest

package noxrender

// Install opaque RLE wall images in the actual image/handle owner. Type 3 is the
// format consumed by the production edge renderer as well as ordinary drawing.
func (b *RenderSprites) PortTestWorldWallImages(data [][]byte) ([]*Image, func()) {
	oldIndex, oldHandles := b.byIndex, b.byHandle
	b.byIndex = make([]*Image, len(data))
	b.byHandle = make(map[ImageHandle]*Image)
	for i, raw := range data {
		b.byIndex[i] = &Image{c: b, typ: 3, raw: append([]byte(nil), raw...)}
		b.byIndex[i].C()
	}
	owned := b.byIndex
	return owned, func() {
		for _, img := range owned {
			img.Free()
		}
		b.byIndex, b.byHandle = oldIndex, oldHandles
	}
}
