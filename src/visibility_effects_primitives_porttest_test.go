//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
	"unsafe"
)

func TestVisibilityEffectsKillable(t *testing.T) {
	s := newCreatureXferOwner(t)
	u := newCreatureXferObject(t, s, "Monster")
	health := u.HealthData
	var rows [][4]uint32
	defer func() {
		spellbookCapture(t, "visibility-effects-killable", rows, "d0f92b2ea6d5f7b87a991ab9452a8c6cdf89527e2d096d023178656d9db04ca5")
	}()
	for _, present := range []bool{false, true} {
		for _, cur := range []uint16{0, 1, 32767, 32768, 65535} {
			for _, max := range []uint16{0, 1, 32767, 32768, 65535} {
				name := fmt.Sprintf("present%t/cur%d/max%d", present, cur, max)
				t.Run(name, func(t *testing.T) {
					u.HealthData = health
					health.Cur = cur
					health.Max = max
					if !present {
						u.HealthData = nil
					}
					rv := legacy.PortTestVisibilityEffects(14, u, nil, nil, nil, [5]int32{}, nil, "")
					want := uint32(0)
					if present && (cur != 0 || max == 0) {
						want = 1
					}
					if rv != want {
						t.Fatalf("return%d want%d", rv, want)
					}
					if health.Cur != cur || health.Max != max {
						t.Fatal("health mutated")
					}
					p := uint32(0)
					if present {
						p = 1
					}
					rows = append(rows, [4]uint32{p, uint32(cur), uint32(max), rv})
					u.HealthData = health
				})
			}
		}
	}
}
func TestVisibilityEffectsFrameCopy(t *testing.T) {
	s := newCreatureXferOwner(t)
	dst := memmap.PtrUint32(0x5D4594, 2487684)
	old := *dst
	t.Cleanup(func() { *dst = old })
	before, after := *(*uint32)(unsafe.Add(unsafe.Pointer(dst), -4)), *(*uint32)(unsafe.Add(unsafe.Pointer(dst), 4))
	var rows [][3]uint32
	defer func() {
		spellbookCapture(t, "visibility-effects-frame-copy", rows, "f2eb7ab62acd056881288d9dc8c95e6298da2abe5236f7b3c6dc5d15ffd2d36c")
	}()
	for _, frame := range []uint32{0, 1, 65535, 65536, 0x7fffffff, 0x80000000, 0xfffffffe, 0xffffffff} {
		for _, op := range []int{15, 16} {
			t.Run(fmt.Sprintf("frame%x/op%d", frame, op), func(t *testing.T) {
				s.SetFrame(frame)
				*dst = 0xAABBCCDD
				rv := legacy.PortTestVisibilityEffects(op, nil, nil, nil, nil, [5]int32{}, nil, "")
				want := frame
				if op == 15 {
					want++
				}
				if rv != want || *dst != want {
					t.Fatalf("return%x stored%x want%x", rv, *dst, want)
				}
				if before != *(*uint32)(unsafe.Add(unsafe.Pointer(dst), -4)) || after != *(*uint32)(unsafe.Add(unsafe.Pointer(dst), 4)) {
					t.Fatal("neighbor modified")
				}
				rows = append(rows, [3]uint32{uint32(op), frame, rv})
			})
		}
	}
}
