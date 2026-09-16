//go:build porttest

package noxrender

// PortTestAssetBag loads the real archive at an explicit fixture path and
// restores the existing indexed-image owner after its consumers are destroyed.
func (b *RenderSprites) PortTestAssetBag(path string) (func(), error) {
	oldBag, oldIndex, oldHandles := b.bag, b.byIndex, b.byHandle
	b.byHandle = make(map[ImageHandle]*Image)
	if err := b.readVideobag(path); err != nil {
		b.bag, b.byIndex, b.byHandle = oldBag, oldIndex, oldHandles
		return nil, err
	}
	return func() { b.Free(); b.bag, b.byIndex, b.byHandle = oldBag, oldIndex, oldHandles }, nil
}
