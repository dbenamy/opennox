//go:build porttest

package legacy

import "unsafe"
import "github.com/opennox/opennox/v1/legacy/common/alloc"
import "github.com/opennox/opennox/v1/client/gui"

// Exercise native production helpers with caller-owned guarded storage.
func PortTestServerPanelsByte(p *byte, mask byte, on int) bool {
	return unsafe.Pointer(serverPanelsByte(p, mask, on)) == unsafe.Pointer(p)
}
func PortTestServerPanelsWord(p *uint32, mask uint32, on int) bool {
	return unsafe.Pointer(serverPanelsWordMask(p, mask, on)) == unsafe.Pointer(p)
}
func PortTestServerPanelsBit(p *uint32, index, on int) (uint32, bool) {
	r := uint32(serverPanelsSpellBit(p, index, on))
	return r, serverPanelsSpellQuery(p, index)
}
func PortTestServerPanelsMaskQuery(mask uint32, armor bool) bool {
	if armor {
		return serverPanelsArmorQuery(mask)
	}
	return serverPanelsWeaponQuery(mask)
}
func PortTestServerPanelsClass(index byte, armor bool) bool {
	old := *serverPanelsWord(1045460)
	*serverPanelsWord(1045460) = uint32(bool2int(armor))
	defer func() { *serverPanelsWord(1045460) = old }()
	return serverPanelsClass(index)
}

func PortTestServerPanelsWeaponStore(p *uint32) bool {
	return unsafe.Pointer(serverPanelsWeaponStore(p)) == unsafe.Pointer(p)
}
func PortTestServerPanelsWeaponPointer() *uint32 {
	return (*uint32)(unsafe.Pointer(serverPanelsWeaponPointer()))
}
func PortTestServerPanelsArmorStore(v uint32) uint32 { return uint32(serverPanelsArmorStore(v)) }
func PortTestServerPanelsArmorLoad() uint32          { return uint32(serverPanelsArmorLoad()) }
func PortTestServerPanelsSpellStore(p *uint32)       { serverPanelsSpellStore(p) }
func PortTestServerPanelsSpellPointer() *uint32 {
	return (*uint32)(unsafe.Pointer(serverPanelsSpellPointer()))
}

// The snapshot helper's last return is either a scalar or a pointer into its
// output. Identify owned pointers before recording; never capture an address.
func PortTestServerPanelsSpellSnapshot(p *uint32) (uint32, int) {
	r := uint32(serverPanelsSpellSnapshot(p))
	for i := 0; i < 5; i++ {
		if uintptr(r) == uintptr(unsafe.Pointer(p))+uintptr(4*i) {
			return 0, i
		}
	}
	return r, -1
}
func PortTestServerPanelsSpellApply(p *uint32) int {
	return int(serverPanelsSpellApply(p))
}

func PortTestServerPanelsObjectWords() (map[string]*uint32, func()) {
	words := map[string]*uint32{
		"kind": serverPanelsWord(1045460),
		"list": serverPanelsWord(1045464),
		"root": serverPanelsWord(1045468),
	}
	old := make(map[string]uint32)
	for k, p := range words {
		old[k] = *p
		*p = 0
	}
	return words, func() {
		for k, p := range words {
			*p = old[k]
		}
	}
}
func PortTestServerPanelsConstruct(kind string, parent *gui.Window, settings unsafe.Pointer) int {
	p := parent
	switch kind {
	case "weapon":
		return int(serverPanelsObjectOpen(p, 0x1000000))
	case "armor":
		return int(serverPanelsObjectOpen(p, 0x2000000))
	case "spell":
		return int(serverPanelsSpellOpen(p))
	case "access":
		return int(serverPanelsAccessOpen(p))
	case "general":
		return int(serverPanelsGeneralOpen(p))
	case "advanced":
		return int(serverPanelsAdvancedOpen(settings))
	case "advserv":
		return int(serverPanelsAdvancedServerOpen())
	default:
		panic(kind)
	}
}
func PortTestServerPanelsRefresh(kind string) {
	switch kind {
	case "object":
		serverPanelsObjectRefresh()
	case "spell":
		serverPanelsSpellRefresh()
	case "access":
		serverPanelsAccessPopulate()
	case "general":
		serverPanelsGeneralRefresh()
	default:
		panic(kind)
	}
}
func PortTestServerPanelsWeaponSnapshot(p *uint32) int {
	return int(serverPanelsWeaponSnapshot(p))
}
func PortTestServerPanelsArmorSnapshot() int { return int(serverPanelsArmorSnapshot()) }

func PortTestServerPanelsAccessLookup(name string) int {
	return int(serverPanelsPlayerLookup(alloc.InternCString16(name)))
}
func PortTestServerPanelsAccessSelected() bool     { return serverPanelsAccessSelected() }
func PortTestServerPanelsAccessRefresh()           { serverPanelsAccessRefresh() }
func PortTestServerPanelsAccessClose(destroy bool) { serverPanelsAccessClose(destroy) }
func PortTestServerPanelsAdvancedUpdate(settings unsafe.Pointer) {
	serverPanelsAdvancedUpdate(settings)
}
func PortTestServerPanelsCallbackTable() func() {
	old := serverPanelsRefreshCallbacks
	serverPanelsRefreshCallbacks = [4]func() int{nil, serverPanelsSpellRefresh, serverPanelsObjectRefresh, serverPanelsObjectRefresh}
	return func() { serverPanelsRefreshCallbacks = old }
}
