//go:build porttest

package noxrender

import "golang.org/x/image/font"

// PortTestDefaultFont owns a default face without touching the shared font files.
func (r *RenderFonts) PortTestDefaultFont(face font.Face) func() {
	old := r.def
	r.def = face
	return func() { r.def = old }
}
