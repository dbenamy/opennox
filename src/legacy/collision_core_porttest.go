//go:build porttest

package legacy

import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/server"
	"runtime"
	"unsafe"
)

// The original observer was aligned to 256 bytes for callback low-byte contracts.
var portTestCoreContactKeyStorage [256]byte

func portTestCoreContactKey() unsafe.Pointer {
	p := unsafe.Pointer(&portTestCoreContactKeyStorage[0])
	return unsafe.Add(p, (-uintptr(p))&255)
}

var portTestCoreContacts [1024]uint32
var portTestCoreContactCount uint32

func init() {
	server.PortTestRegisterCollideCallback(portTestCoreContactKey(), func(u *server.Object, a, b uintptr) uint32 {
		if portTestCoreContactCount >= 256 {
			return 0
		}
		i := 4 * portTestCoreContactCount
		portTestCoreContactCount++
		portTestCoreContacts[i] = uint32(uintptr(unsafe.Pointer(u)))
		portTestCoreContacts[i+1] = uint32(a)
		p := (*[2]uint32)(unsafe.Pointer(b))
		portTestCoreContacts[i+2], portTestCoreContacts[i+3] = p[0], p[1]
		return 0
	})
}

func PortTestCollisionCore(op string, a, b *server.Object, p *types.Pointf, mode int32) int32 {

	switch op {
	case "angle-queue":
		worldAngleQueue(a.UpdateData)
	case "shape":
		geometryShapeBox(&a.Shape)
	case "contains":
		return collisionObjectContains(a, p)
	case "eligible":
		return collisionEligible(a, b)
	case "retained":
		return collisionRetained(a, b)
	case "pair":
		collisionPair(a, b)
	case "scan":
		collisionScan(a)
	case "circle-walls":
		collisionCircleWalls(a)
	case "circle-box":
		collisionCircleBox(a, b, mode)
	case "gate-circle":
		collisionGateCircle(a, b, mode)
	case "dispatch":
		collisionDispatch()
	case "angles":
		collisionDrainAngles()
	case "elevator":
		collisionElevator(a, b, mode)
	case "shaft":
		collisionShaft(a, b)
	case "types":
		collisionInitTypes()
	case "active":
		return int32(a.Field116 & 1)
	case "activate":
		return int32(collisionActivate(a))
	case "remove":
		collisionRemoveActive(a)
	case "head":
		return int32(collisionActiveHead)
	case "next":
		return int32(collisionNextActive(a))
	case "pop":
		return int32(collisionObjectAddress(collisionPopActive()))
	default:
		panic(op)
	}
	return 0
}
func PortTestCollisionCoreDistance(p *types.Pointf, radius float32, box *server.Object, out *types.Pointf) float64 {
	return collisionBoxDistance(p, radius, box, out)
}
func PortTestCollisionCoreWallOpen(grid *[2]int32, u *server.Object) {
	collisionWallOpen(grid, u)
}
func PortTestCollisionCoreAddHit(a, b *server.Object, sentinel uint32, normal *types.Pointf) {
	target := uint32(sentinel)
	if b != nil {
		target = uint32(uintptr(b.CObj()))
	}
	collisionAddHit(a, uint32(target), normal)
}
func PortTestCollisionCoreObserver() (unsafe.Pointer, func(), func(map[unsafe.Pointer]uint32) [][4]uint32, func()) {
	data := portTestCoreContacts[:]
	count := &portTestCoreContactCount
	old := append([]uint32(nil), data...)
	oldN := *count
	reset := func() { clear(data); *count = 0 }
	reset()
	snapshot := func(ids map[unsafe.Pointer]uint32) [][4]uint32 {
		out := make([][4]uint32, *count)
		for i := range out {
			copy(out[i][:], data[4*i:4*i+4])
			for j := 0; j < 2; j++ {
				if out[i][j] != 0 {
					id, ok := ids[unsafe.Pointer(uintptr(out[i][j]))]
					if !ok {
						panic("contact outside fixture")
					}
					out[i][j] = id
				}
			}
		}
		return out
	}
	return portTestCoreContactKey(), reset, snapshot, func() { copy(data, old); *count = oldN }
}
func PortTestCollisionCorePentagram() unsafe.Pointer { return collisionKey(collisionIdentityPentagram) }
func PortTestCollisionCoreGlobals() (map[string]*uint32, func()) {
	words := map[string]*uint32{"trigger": &collisionTrigger, "powder": &collisionPowder, "hand": &collisionHand, "small": &collisionSmallFist, "medium": &collisionMediumFist, "large": &collisionLargeFist, "meteor": &collisionMeteor, "ready": &collisionTypesReady}
	old := map[string]uint32{}
	for n, p := range words {
		old[n] = *p
		*p = 0
	}
	return words, func() {
		for n, p := range words {
			*p = old[n]
		}
	}
}

// Inspect every real bucket link and validate its bucket and object identities.
func PortTestCollisionCoreBuckets(ids map[unsafe.Pointer]uint32) map[uint32][][2]uint32 {
	out := map[uint32][][2]uint32{}
	for i := uint32(0); i < 256; i++ {
		p := collisionBuckets[i]
		for n := 0; p != 0; n++ {
			if n >= 1024 {
				panic("collision bucket cycle")
			}
			r := (*[7]uint32)(unsafe.Pointer(uintptr(p)))
			if r[6] != i {
				panic("wrong collision bucket")
			}
			pair := [2]uint32{r[2], r[3]}
			for j := range pair {
				if pair[j] > 6 {
					id, ok := ids[unsafe.Pointer(uintptr(pair[j]))]
					if !ok {
						panic("bucket outside fixture")
					}
					pair[j] = id
				}
			}
			out[i] = append(out[i], pair)
			p = r[0]
		}
	}
	return out
}

// Observe the production radial predicate through its native callback.
func PortTestCollisionCoreRadial(u *server.Object, p *types.Pointf, radius float32, code uint32) [3]uint32 {
	var out [3]uint32
	motionRadialCandidate(u, p, radius, func(u *server.Object) {
		out[0]++
		out[1] = uint32(uintptr(unsafe.Pointer(u)))
		out[2] = code
		runtime.KeepAlive(u)
	})
	return out
}
