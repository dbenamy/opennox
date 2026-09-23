//go:build porttest

package legacy

/*
#include "GAME3_1.h"
#include "GAME3_2.h"
#include "client__drawable__update__charmup.h"
#include "client__drawable__update__cloud.h"
#include "client__drawable__update__dball.h"
#include "client__drawable__update__drainup.h"
#include "client__drawable__update__fireball.h"
#include "client__drawable__update__healup.h"
#include "client__drawable__update__manabomb.h"
#include "client__drawable__update__mmislup.h"
#include "client__drawable__update__mtailup.h"
#include "client__drawable__update__sparklup.h"
#include "client__drawable__update__telwake.h"
#include "client__drawable__update__vortexup.h"
*/
import "C"

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

// The owner supplies a real client, drawable pool/list/index and seeded RNG.
// Integer pointer parameters are the production 386 callback ABI.
func PortTestClientUpdate(op int, vp *noxrender.Viewport, dr *client.Drawable, a [4]int32) uint32 {
	viewport := C.int(uintptr(vp.C()))
	drawable := C.int(uintptr(dr.C()))
	words := (*C.uint32_t)(dr.C())
	viewWords := (*C.uint32_t)(vp.C())
	switch op {
	case 0:
		updateTransfer(int(a[0]), vp, dr, a[1] != 0, true)
		return 0
	case 1:
		return uint32(updateCloud(dr, int(a[0]), int(a[1])))
	case 2:
		return updateDeathBallSparks(dr, int(a[0]))
	case 3:
		return uint32(C.nox_xxx_updDrawDBallCharge_4CE0C0(viewport, drawable))
	case 4:
		return uint32(C.sub_4CD690(viewWords, drawable))
	case 5:
		updateFireball(dr, int(a[0]))
		return 0
	case 6:
		return uint32(C.sub_4CD450(viewWords, drawable))
	case 7:
		return uint32(C.nox_xxx_updDrawManabombCharge_4CCAC0(viewport, words))
	case 8:
		return uint32(C.nox_xxx_updDrawMagicMissile_4CD9E0(viewport, words))
	case 9:
		C.nox_xxx_updDrawMagic_4CDD80(viewport, words)
		return 0
	case 10:
		return uint32(C.nox_xxx_updDrawSparkleTrail_4CDBF0(viewport, words))
	case 11:
		return uint32(C.nox_xxx_updDrawTeleportWake_4CD8D0(viewport, drawable))
	case 12:
		return uint32(C.nox_xxx_updDrawVortexSource_4CC950(viewport, drawable))
	case 13:
		return uint32(C.nox_xxx_updDrawUndeadKiller_4CCCF0())
	case 14:
		return uint32(C.sub_4CCD00(viewport, drawable))
	case 15:
		return uint32(C.nox_xxx_updDrawFist_4CCDB0(viewport, drawable))
	case 16:
		return uint32(C.sub_4CCE70(viewport, words))
	case 17:
		return uint32(C.sub_4CD090(viewport, words))
	case 18:
		return uint32(C.sub_4CD0C0(viewport, words))
	case 19:
		return uint32(C.sub_4CD0F0(viewport, words))
	case 20:
		return uint32(C.sub_4CD120(viewport, words))
	case 21:
		return uint32(C.sub_4CD400(viewWords, drawable))
	case 22:
		return uint32(C.nox_xxx_updDrawDBall_4CDF80(viewport, drawable))
	case 23:
		return uint32(C.sub_4CE0A0(viewport, drawable))
	case 24:
		return uint32(C.nox_xxx_updDrawCloud_4CE1D0(viewport, drawable))
	case 25:
		return uint32(C.sub_4CE340(viewport, drawable))
	case 26:
		return uint32(C.sub_4CE360(viewport, drawable))
	default:
		panic("unknown drawable update")
	}
}

func PortTestClientUpdateCloudCallback() unsafe.Pointer { return unsafe.Pointer(C.sub_4CE340) }

type PortTestClientUpdateEnvironment struct {
	mapped  []uint32
	named   [2]uint32
	density uint32
}

func PortTestNewClientUpdateEnvironment() *PortTestClientUpdateEnvironment {
	p := &PortTestClientUpdateEnvironment{mapped: append([]uint32(nil), unsafe.Slice(memmap.PtrUint32(0x5D4594, 1522944), 20)...), named: [2]uint32{uint32(dword_5d4594_1522956), uint32(dword_5d4594_1522968)}, density: *memmap.PtrUint32(0x587000, 190108)}
	p.Reset(0)
	return p
}
func (p *PortTestClientUpdateEnvironment) Reset(density uint32) {
	clear(unsafe.Slice(memmap.PtrUint32(0x5D4594, 1522944), 20))
	dword_5d4594_1522956 = 0
	dword_5d4594_1522968 = 0
	*memmap.PtrUint32(0x587000, 190108) = density
}
func (p *PortTestClientUpdateEnvironment) Restore() {
	copy(unsafe.Slice(memmap.PtrUint32(0x5D4594, 1522944), 20), p.mapped)
	dword_5d4594_1522956 = C.uint32_t(p.named[0])
	dword_5d4594_1522968 = C.uint32_t(p.named[1])
	*memmap.PtrUint32(0x587000, 190108) = p.density
}
func (p *PortTestClientUpdateEnvironment) Snapshot() []uint32 {
	out := []uint32{uint32(dword_5d4594_1522956), uint32(dword_5d4594_1522968), *memmap.PtrUint32(0x587000, 190108)}
	return append(out, unsafe.Slice(memmap.PtrUint32(0x5D4594, 1522944), 20)...)
}
