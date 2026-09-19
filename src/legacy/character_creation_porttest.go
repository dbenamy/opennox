//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"unsafe"
)

func PortTestCharacterName(p unsafe.Pointer) int { return characterName((*uint16)(p)) }
func PortTestCharacterPaletteOwner(root, menu unsafe.Pointer) func() {
	a, b, c := characterUI.colorRoot, characterUI.palette, characterUI.target
	characterUI.colorRoot = (*gui.Window)(root)
	characterUI.palette = (*gui.Window)(menu)
	return func() { characterUI.colorRoot = a; characterUI.palette = b; characterUI.target = c }
}
func PortTestCharacterPaletteMatch(w unsafe.Pointer, bank int, rgb unsafe.Pointer) uintptr {
	return characterPaletteMatch((*gui.Window)(w), bank, (*byte)(rgb))
}
func PortTestCharacterPaletteFill(bank uint16) uintptr {
	return uintptr(characterPaletteFill(bank).C())
}
func PortTestCharacterPaletteClose(id int, index uint16) uintptr {
	characterUI.target = uint32(id)
	return characterPaletteClose(index)
}
func PortTestCharacterPaletteDraw(w unsafe.Pointer) int {
	return characterPaletteDraw((*gui.Window)(w))
}
func PortTestCharacterAppearanceOwner(player unsafe.Pointer, windows []unsafe.Pointer) func() {
	host, controls := characterUI.colorHost, characterUI.controls
	characterUI.colorHost = player
	for i, p := range windows {
		characterUI.controls[i] = (*gui.Window)(p)
	}
	return func() { characterUI.colorHost = host; characterUI.controls = controls }
}
func PortTestCharacterAppearance() unsafe.Pointer { return characterAppearance() }
func PortTestCharacterCreateFile() int            { return characterCreateFile() }
func PortTestCharacterClassWindow(root unsafe.Pointer) func() {
	old := characterUI.classRoot
	characterUI.classRoot = (*gui.Window)(root)
	return func() { characterUI.classRoot = old }
}
func PortTestCharacterClassEvent(root, child unsafe.Pointer, event int) int {
	return gui.EventRespInt(characterClassEvent((*gui.Window)(root), gui.AsWindowEvent(event, uintptr(child), 0)))
}
func PortTestCharacterClassDraw(w unsafe.Pointer) int { return characterClassDraw((*gui.Window)(w)) }
func PortTestCharacterConstruct(color bool) int {
	if color {
		return characterShowColor()
	}
	return characterShowClass()
}
func PortTestCharacterDefaultSet(v int32) int32   { Sub_4A7A60(int(v)); return v }
func PortTestCharacterAdmissionSet(v int32) int32 { Sub_4A7A70(int(v)); return v }
func PortTestCharacterDefaultGet() uint32         { return characterUI.defaults }
func PortTestCharacterDefaultOwner(value uint32) func() {
	old := characterUI.defaults
	characterUI.defaults = value
	return func() { characterUI.defaults = old }
}
func PortTestCharacterSetup() uint32 { return characterSetup() }
func PortTestCharacterColorEvent(root, child unsafe.Pointer, event int, point uint32) int {
	return gui.EventRespInt(characterColorEvent((*gui.Window)(root), gui.AsWindowEvent(event, uintptr(child), uintptr(point))))
}
func PortTestCharacterPaletteOutside(menu unsafe.Pointer, event int, point uint32) int {
	return gui.EventRespInt(characterPaletteOutside((*gui.Window)(menu), gui.AsWindowEvent(event, uintptr(point), 0)))
}
func PortTestCharacterAnimationOwner(color bool, anim unsafe.Pointer) func() {
	p := &characterUI.classAnim
	if color {
		p = &characterUI.colorAnim
	}
	old := *p
	*p = (*gui.Anim)(anim)
	return func() { *p = old }
}
func PortTestCharacterAnimationDone(color bool) int {
	if color {
		return characterColorDone()
	}
	return characterClassDone()
}
func PortTestCharacterClassStart() int { return characterClassStart() }
func PortTestCharacterPreviewOwner(pants, shirt, shoes unsafe.Pointer) func() {
	old := characterUI.modifiers
	characterUI.modifiers = [3]unsafe.Pointer{pants, shirt, shoes}
	return func() { characterUI.modifiers = old }
}
func PortTestCharacterPreview(w unsafe.Pointer) int { return characterPreview((*gui.Window)(w)) }
func PortTestCharacterClassOwner(p unsafe.Pointer) func() {
	old := characterUI.classHost
	characterUI.classHost = p
	return func() { characterUI.classHost = old }
}
func PortTestCharacterQuickbar(mode int) uintptr { return characterQuickbar(mode) }

func PortTestCharacterClassNewChild(id uint) int {
	return gui.EventRespInt(characterClassEvent(nil, gui.WindowNewChild{ID: id}))
}
