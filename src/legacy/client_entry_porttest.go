//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
)

type PortTestEntryEnvironment struct {
	context bool
	active  *gui.Window
	blink   byte
}

func PortTestNewEntryEnvironment() *PortTestEntryEnvironment {
	e := &PortTestEntryEnvironment{uiEntryContext, uiEntryActive, *memmap.PtrUint8(0x5D4594, 1193344)}
	uiEntryContext = false
	uiEntryActive = nil
	*memmap.PtrUint8(0x5D4594, 1193344) = 0
	return e
}
func (e *PortTestEntryEnvironment) Context(on bool) { uiEntryContext = on }
func (e *PortTestEntryEnvironment) ClearActive()    { uiEntryActive = nil }
func (e *PortTestEntryEnvironment) State() [3]uint32 {
	context := uint32(0)
	if uiEntryContext {
		context = 1
	}
	return [3]uint32{context, uint32(uintptr(uiEntryActive.C())), uint32(*memmap.PtrUint8(0x5D4594, 1193344))}
}
func (e *PortTestEntryEnvironment) Blink(v byte) { *memmap.PtrUint8(0x5D4594, 1193344) = v }
func (e *PortTestEntryEnvironment) Restore() {
	uiEntryContext = e.context
	uiEntryActive = e.active
	*memmap.PtrUint8(0x5D4594, 1193344) = e.blink
}
