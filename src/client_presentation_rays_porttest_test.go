//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestClientPresentationRayEvents(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "CharmRay", "DrainManaRay", "HealRay", "HarpoonRope")
	state := serverConfigOwnBytes(t, 0x5D4594, 1303540, 860)
	names := map[byte]string{1: "PlasmaRay", 2: "CharmRay", 3: "DynamicChainLightning", 4: "DynamicEnergyBolt", 5: "DrainManaRay", 6: "HealRay", 7: "HarpoonRope", 140: "DynamicLightning"}
	type record struct {
		Kind                      byte
		Classes, Missing, Failure int
		Position                  image.Point
		Parameter                 byte
		Calls                     []effectsSpawnCall
		Stored                    [13]byte
		Deleted                   []uint32
	}
	var rows []record
	for _, kind := range []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 140, 255} {
		for classes := 0; classes < 4; classes++ {
			for missing := 0; missing < 3; missing++ {
				for _, fail := range []int{0, 1} {
					for _, param := range []byte{0, 255} {
						c.resetCase(env, pix, 41, 120)
						clear(state)
						a, b := image.Pt(123, -456), image.Pt(-124, 457)
						endpoints := [2]*client.Drawable{}
						for i, pos := range []image.Point{a, b} {
							if missing == i+1 {
								continue
							}
							dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, pos)
							dr.NetCode32 = uint32(7 + i)
							dr.ObjClass = 0
							if classes&(1<<i) != 0 {
								dr.ObjClass = 0x20000000
							}
							endpoints[i] = dr
						}
						c.Calls = nil
						c.FailEvery = fail
						event := [7]byte{0, kind, param}
						codeA, codeB := uint16(7), uint16(8)
						if classes&1 != 0 {
							codeA |= 0x8000
						}
						if classes&2 != 0 {
							codeB |= 0x8000
						}
						binary.LittleEndian.PutUint16(event[3:], codeA)
						binary.LittleEndian.PutUint16(event[5:], codeB)
						legacy.PortTestPresentationRayAdd(&event)
						name, valid := names[kind]
						attempt := valid && missing == 0
						wantCalls := 0
						if attempt {
							wantCalls = 1
						}
						if len(c.Calls) != wantCalls {
							t.Fatal("ray kind/endpoint allocation boundary")
						}
						dr := (*client.Drawable)(*memmap.PtrPtr(0x5D4594, 1303924))
						if (dr != nil) != (attempt && fail == 0) {
							t.Fatal("ray ownership after allocation")
						}
						var stored [13]byte
						mid := image.Pt(a.X+(b.X-a.X)/2, a.Y+(b.Y-a.Y)/2)
						if attempt && (c.Calls[0].Type != c.Things.IndByID(name) || c.Calls[0].Position != mid) {
							t.Fatal("ray type/midpoint")
						}
						if dr != nil {
							copy(stored[:], unsafe.Slice((*byte)(unsafe.Pointer(&dr.Union)), 13))
							var want [13]byte
							want[0] = 1
							binary.LittleEndian.PutUint32(want[1:], uint32(param))
							binary.LittleEndian.PutUint32(want[5:], uint32(codeA))
							binary.LittleEndian.PutUint32(want[9:], uint32(codeB))
							if stored != want || dr.PosVec != mid || !legacy.PortTestPresentationRayContains(dr) {
								t.Fatal("ray payload/membership")
							}
							for _, end := range endpoints {
								if legacy.PortTestPresentationRayContains(end) {
									t.Fatal("endpoint counted as ray")
								}
							}
						}
						legacy.PortTestPresentationRayRemove(&event)
						wantDeleted := 0
						if dr != nil {
							wantDeleted = 1
						}
						if len(c.Deleted) != wantDeleted || *memmap.PtrPtr(0x5D4594, 1303924) != nil {
							t.Fatal("ray removal")
						}
						legacy.PortTestPresentationRayRemove(&event)
						if len(c.Deleted) != wantDeleted {
							t.Fatal("ray removed twice")
						}
						rows = append(rows, record{kind, classes, missing, fail, mid, param, append([]effectsSpawnCall(nil), c.Calls...), stored, append([]uint32(nil), c.Deleted...)})
					}
				}
			}
		}
	}
	spellbookCapture(t, "client-presentation-ray-events", rows, "fb17d15996b07ed5de6bd7c7b3a7101565866bfa4e5bb47df3b1144670de5dbc")
}

func TestClientPresentationRayCapacity(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t, "CharmRay", "DrainManaRay", "HealRay", "HarpoonRope")
	clear(serverConfigOwnBytes(t, 0x5D4594, 1303540, 860))
	for i := 0; i < 2; i++ {
		dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(i*3, i*5))
		dr.NetCode32 = uint32(7 + i)
		dr.ObjClass = 0
	}
	event := [7]byte{0, 1, 255, 7, 0, 8, 0}
	for i := 0; i < 96; i++ {
		legacy.PortTestPresentationRayAdd(&event)
	}
	refs := make([]uint32, 96)
	for i := range refs {
		dr := (*client.Drawable)(*memmap.PtrPtr(0x5D4594, 1303924+4*uintptr(i)))
		if dr == nil || !legacy.PortTestPresentationRayContains(dr) {
			t.Fatal("ray capacity ownership")
		}
		refs[i] = c.refs[dr]
	}
	before := c.Objs.Count
	legacy.PortTestPresentationRayAdd(&event)
	if c.Objs.Count != before {
		t.Fatalf("full ray table allocated an untracked sprite: count %d -> %d", before, c.Objs.Count)
	}
	// Duplicate endpoint pairs remove one entry at a time, leaving the others live.
	legacy.PortTestPresentationRayRemove(&event)
	if len(c.Deleted) != 1 || c.Deleted[0] != refs[0] {
		t.Fatal("duplicate ray removal order")
	}
	legacy.PortTestPresentationRayAdd(&event)
	first := (*client.Drawable)(*memmap.PtrPtr(0x5D4594, 1303924))
	if first == nil {
		t.Fatal("ray failed to reuse first free slot")
	}
	legacy.PortTestPresentationRayClear()
	if c.Objs.Count != 2 || len(c.Deleted) != 97 {
		t.Fatal("ray bulk cleanup leaked a sprite")
	}
	for i := 0; i < 96; i++ {
		if *memmap.PtrPtr(0x5D4594, 1303924+4*uintptr(i)) != nil {
			t.Fatal("ray cleanup left a pointer")
		}
	}
	spellbookCapture(t, "client-presentation-ray-capacity", append(refs, c.Deleted...), "ec224161b48942a1f1d7eb48f629c27f1b073ea56c10cb7b7d2cfc891e5b8e16")
}
