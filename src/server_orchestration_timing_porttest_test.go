//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestServerOrchestrationTimeout(t *testing.T) {
	serverConfigOwnBytes(t, 0x5D4594, 2523788, 16)
	start, enabled := memmap.PtrUint64(0x5D4594, 2523788), memmap.PtrUint64(0x5D4594, 2523796)
	clock := legacy.PlatformTicks
	t.Cleanup(func() { legacy.PlatformTicks = clock })
	type row struct {
		Enabled, Start, Now uint64
		Result              uint32
		Reads               int
	}
	var rows []row
	for _, on := range []uint64{0, 1, 1 << 63, math.MaxUint64} {
		for _, base := range []uint64{0, 123, math.MaxUint64 - 5000, math.MaxUint64 - 4999, math.MaxUint64} {
			for _, now := range []uint64{0, base + 4999, base + 5000, base + 5001, math.MaxUint64} {
				*start, *enabled = base, on
				reads := 0
				legacy.PlatformTicks = func() uint64 { reads++; return now }
				got := legacy.PortTestServerOrchestration("timeout", nil, 0)
				want := uint32(0)
				if on != 0 && base+5000 < now {
					want = 1
				}
				wantReads := 0
				if on != 0 {
					wantReads = 1
				}
				if got != want || reads != wantReads || *start != base || *enabled != on {
					t.Fatalf("timeout on%x/start%x/now%x got%d reads%d", on, base, now, got, reads)
				}
				rows = append(rows, row{on, base, now, got, reads})
			}
		}
	}
	spellbookCapture(t, "server-orchestration-timeout", rows, "b65850caac85b5853b8bde9049aa58adbe5182b0424ab8f7a49e5d69b2301ba4")
}

func TestServerOrchestrationLoadState(t *testing.T) {
	serverConfigOwnBytes(t, 0x5D4594, 1563072, 4)
	p := memmap.PtrUint32(0x5D4594, 1563072)
	var rows [][2]uint32
	for _, v := range []uint32{0, 1, 2, 0x7fffffff, 0x80000000, 0xffffffff} {
		*p = v
		got := legacy.PortTestServerOrchestration("load-state", nil, 0)
		if got != v || *p != v {
			t.Fatal("state bits", v, got)
		}
		rows = append(rows, [2]uint32{v, got})
	}
	spellbookCapture(t, "server-orchestration-load-state", rows, "06bc7ceb21e4856acda9505ef77ee3b0f98be4968aefcff276f21d005f92f82a")
}

func TestServerOrchestrationDifficultyTiming(t *testing.T) {
	o := newMatchRosterOwner(t)
	o.balance(map[string]float64{"PlayerDifficultyDelta": 0.5, "GeneratorMaxHealth": 999, "PlayerDamageDiffInit": 2, "SystemHealthDiffInit": 3, "PlayerDamageDiffCoeff": 0.25, "SystemHealthDiffCoeff": 0.5, "PlayerDamageCap": 100, "SystemHealthCap": 100})
	oldObjects := o.s.Objs.First()
	o.s.Objs.SetObjects(nil)
	t.Cleanup(func() { o.s.Objs.SetObjects(oldObjects) })
	type row struct{ FPS, Frame, Before, After, Damage, Health uint32 }
	var rows []row
	// Positive FPS is the runtime precondition; zero would divide by zero in C.
	for mask := 0; mask < 8; mask++ {
		count := 0
		for i := range o.units {
			v := uint32(0)
			if mask&(1<<i) != 0 {
				v = 1
				count++
			}
			objectXferSetWord(o.units[i].UpdateDataPlayer().Player.C(), 4792, v)
		}
		for _, fps := range []uint32{1, 30, 60, 0x80000001} {
			for _, frame := range []uint32{0, 1, 149, 150, 151, 5 * fps, 5*fps + 1, 0xffffffff} {
				for _, same := range []bool{false, true} {
					name := fmt.Sprintf("mask%d/fps%x/frame%x/same%v", mask, fps, frame, same)
					t.Run(name, func(t *testing.T) {
						o.s.SetTickRate(fps)
						o.s.SetFrame(frame)
						*memmap.PtrUint32(0x587000, 202028) = 2
						*memmap.PtrUint32(0x5D4594, 1563928) = 0
						*memmap.PtrUint32(0x5D4594, 1563932) = 0
						// Count comes from the real active-player owner; derive the expected stage
						// formula from that owner, then separately assert the tick boundary.
						// Only the explicitly participating units count; the unitless player does not.
						expected := float32(2 * (1 + float64(count-1)*0.5))
						before := float32(99)
						if same {
							before = expected
						}
						*memmap.PtrFloat32(0x587000, 202024) = before
						*memmap.PtrFloat32(0x587000, 202032) = 17
						*memmap.PtrFloat32(0x587000, 202036) = 19
						legacy.PortTestServerOrchestration("difficulty", nil, 0)
						after, damage, health := *memmap.PtrFloat32(0x587000, 202024), *memmap.PtrFloat32(0x587000, 202032), *memmap.PtrFloat32(0x587000, 202036)
						want, wd, wh := before, float32(17), float32(19)
						if frame%(5*fps) == 0 {
							want = expected
							if expected != before {
								wd = 2 + (expected-1)*0.25
								wh = 3 + (expected-1)*0.5
							}
						}
						if after != want || damage != wd || health != wh {
							t.Fatalf("timing got%v,%v,%v want%v,%v,%v count%d", after, damage, health, want, wd, wh, count)
						}
						rows = append(rows, row{fps, frame, math.Float32bits(before), math.Float32bits(after), math.Float32bits(damage), math.Float32bits(health)})
					})
				}
			}
		}
	}
	spellbookCapture(t, "server-orchestration-difficulty-timing", rows, "6359cabe21d8bbc4cf1a2777e04decc953103b293ba6b4ba12b26b886625492f")
}
