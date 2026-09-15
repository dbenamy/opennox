//go:build porttest

package noxrender

// PortTestSpriteImages installs indexed images in the actual handle owner.
// Pixel data is interned by Image.Pixdata exactly as for video.bag images.
// The returned cleanup restores the previous owner without closing its bag.
func (b *RenderSprites) PortTestSpriteImages(data [][]byte) ([]*Image, func()) {
	oldIndex, oldHandles := b.byIndex, b.byHandle
	b.byIndex = make([]*Image, len(data))
	b.byHandle = make(map[ImageHandle]*Image)
	for i, raw := range data {
		b.byIndex[i] = &Image{c: b, typ: 6, raw: append([]byte(nil), raw...)}
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
