//go:build porttest

package noxrender

// Fixtures use the actual bag handle and lazily allocated pixel-data owners.
func (b *RenderSprites) PortTestWallEdgeImages(data [][]byte, types []int) ([]*Image, func()) {
	oldIndex, oldHandles := b.byIndex, b.byHandle
	b.byIndex = make([]*Image, len(data))
	b.byHandle = make(map[ImageHandle]*Image)
	for i, raw := range data {
		b.byIndex[i] = &Image{c: b, typ: types[i], raw: append([]byte(nil), raw...)}
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
