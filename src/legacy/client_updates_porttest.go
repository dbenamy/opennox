//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"unsafe"
)

// The owner supplies a real client, drawable pool/list/index and seeded RNG.
// Results retain the original 32-bit callback ABI.
func PortTestClientUpdate(op int, vp *noxrender.Viewport, dr *client.Drawable, a [4]int32) uint32 {
	switch op {
	case 0:
		updateTransfer(int(a[0]), vp, dr, a[1] != 0, true)
		return 0
	case 1:
		return uint32(updateCloud(dr, int(a[0]), int(a[1])))
	case 2:
		return updateDeathBallSparks(dr, int(a[0]))
	case 3:
		return uint32(int32(updateDeathBallCharge(dr)))
	case 4:
		return uint32(int32(updateHealDrain(vp, dr, false)))
	case 5:
		updateFireball(dr, int(a[0]))
		return 0
	case 6:
		return uint32(int32(updateHealDrain(vp, dr, true)))
	case 7:
		return uint32(int32(updateManaBomb(dr)))
	case 8:
		return uint32(int32(updateMagicMissile(dr)))
	case 9:
		updateMagicTrail(dr)
		return 0
	case 10:
		return uint32(int32(updateTrailSparks(dr, false)))
	case 11:
		return uint32(int32(updateTeleportWake(dr)))
	case 12:
		return uint32(int32(updateVortex(dr)))
	case 13:
		return 1
	case 14:
		return uint32(int32(updateHeight(dr, false)))
	case 15:
		return uint32(int32(updateHeight(dr, true)))
	case 16:
		return uint32(int32(updateFireballFrame(dr, 5)))
	case 17:
		return uint32(int32(updateFireballFrame(dr, 4)))
	case 18:
		return uint32(int32(updateFireballFrame(dr, 3)))
	case 19:
		return uint32(int32(updateFireballFrame(dr, 2)))
	case 20:
		return uint32(int32(updateFireballFrame(dr, 1)))
	case 21:
		return uint32(int32(updateCharm(vp, dr)))
	case 22:
		updateDeathBallSparks(dr, 3)
		return 1
	case 23:
		updateDeathBallSparks(dr, 1)
		return 1
	case 24:
		return uint32(int32(updateCloudFrame(dr, 75)))
	case 25:
		return uint32(int32(updateCloudRise(dr)))
	case 26:
		return uint32(int32(updateCloudFrame(dr, 35)))
	default:
		panic("unknown drawable update")
	}
}

func PortTestClientUpdateCloudCallback() unsafe.Pointer {
	return drawableUpdateIdentity(updateID_sub_4CE340)
}
func PortTestClientUpdateMagicCallback() unsafe.Pointer {
	return drawableUpdateIdentity(updateID_magic)
}

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
	dword_5d4594_1522956 = uint32(p.named[0])
	dword_5d4594_1522968 = uint32(p.named[1])
	*memmap.PtrUint32(0x587000, 190108) = p.density
}
func (p *PortTestClientUpdateEnvironment) Snapshot() []uint32 {
	out := []uint32{uint32(dword_5d4594_1522956), uint32(dword_5d4594_1522968), *memmap.PtrUint32(0x587000, 190108)}
	return append(out, unsafe.Slice(memmap.PtrUint32(0x5D4594, 1522944), 20)...)
}
