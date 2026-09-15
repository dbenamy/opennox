//go:build porttest

package noxrender

import (
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"golang.org/x/image/font"
	"unsafe"
)

// PortTestWindowFont registers a real renderer font reference, with the same
// opaque handle ownership as loaded fonts, without requiring files on disk.
func (r *RenderFonts) PortTestWindowFont(face font.Face, aliases ...string) (unsafe.Pointer, func()) {
	oldNames, oldPointers := r.byName, r.byPtr
	f := &fontFile{Name: "port-window", Font: face, Ptr: handles.NewPtr()}
	r.byName = make(map[string]*fontFile, len(oldNames)+1)
	for k, v := range oldNames {
		r.byName[k] = v
	}
	r.byName[f.Name] = f
	for _, name := range aliases {
		r.byName[name] = f
	}
	r.byPtr = make(map[unsafe.Pointer]*fontFile, len(oldPointers)+1)
	for k, v := range oldPointers {
		r.byPtr[k] = v
	}
	r.byPtr[f.Ptr] = f
	return f.Ptr, func() { r.byName = oldNames; r.byPtr = oldPointers }
}
