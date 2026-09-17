//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"math"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"github.com/opennox/opennox/v1/server"
)

type visibilityEffectsRow struct {
	Name    string
	Return  uint32
	Packets [][]byte
}

func visibilityEffectsPlayers(t *testing.T) (*server.Server, []server.Object, func(int)) {
	t.Helper()
	t.Cleanup(handles.PortTestInit())
	s := newCreatureXferOwner(t)
	units, configure, _, free := s.PortTestEscortPlayers()
	t.Cleanup(free)
	old := s.NetList
	s.NetList = netlist.New()
	s.NetList.Init()
	t.Cleanup(func() { s.NetList.Free(); s.NetList = old })
	configure(3)
	return s, units, configure
}
func visibilityEffectsPackets(s *server.Server) [][]byte {
	out := make([][]byte, 32)
	for i := range out {
		out[i] = s.NetList.CopyPacketsA(ntype.PlayerInd(i), netlist.Kind1)
	}
	return out
}
func TestVisibilityEffectsViewport(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	payload := []byte{0x91, 0, 255, 3, 4}
	var rows []visibilityEffectsRow
	defer func() {
		spellbookCapture(t, "visibility-effects-viewport", rows, "89c3ed8a970f3ea32af3a626c7c0388fdafcf87ffb7087754f12f17c6dfd33fa")
	}()
	type pointCase struct {
		name   string
		p      types.Pointf
		inside bool
	}
	for _, count := range []int{0, 1, 2, 3} {
		configure(count)
		for _, mode := range []uint32{0, 1, 2, 3, 4, 255} {
			for _, camera := range []bool{false, true} {
				for i := range units {
					u := &units[i]
					u.PosVec = types.Pointf{X: 100, Y: 200}
					pl := u.UpdateDataPlayer().Player
					pl.Field10 = 20
					pl.Field12 = 30
					pl.Field3680 = mode
					pl.CameraFollowObj = nil
					if camera {
						pl.CameraFollowObj = &units[(i+1)%3]
					}
				}
				cases := []pointCase{{"center", types.Pointf{100, 200}, true}, {"left", types.Pointf{30, 200}, false}, {"right", types.Pointf{170, 200}, false}, {"top", types.Pointf{100, 120}, false}, {"bottom", types.Pointf{100, 280}, false}}
				for _, edge := range []struct {
					x, y   float32
					axis   int
					toward float32
				}{{30, 200, 0, 100}, {170, 200, 0, 100}, {100, 120, 1, 200}, {100, 280, 1, 200}} {
					p := types.Pointf{edge.x, edge.y}
					if edge.axis == 0 {
						p.X = math.Nextafter32(p.X, edge.toward)
					} else {
						p.Y = math.Nextafter32(p.Y, edge.toward)
					}
					cases = append(cases, pointCase{fmt.Sprintf("inner-%v-%v", edge.x, edge.y), p, true})
				}
				cases = append(cases, pointCase{"nan", types.Pointf{float32(math.NaN()), 200}, false}, pointCase{"inf", types.Pointf{100, float32(math.Inf(1))}, false})
				for _, tc := range cases {
					name := fmt.Sprintf("players%d/mode%d/camera%t/%s", count, mode, camera, tc.name)
					t.Run(name, func(t *testing.T) {
						s.NetList.ResetAll()
						before := make([][]byte, len(units))
						for i := range units {
							before[i] = bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(&units[i])), int(unsafe.Sizeof(units[i]))))
						}
						rv := legacy.PortTestVisibilityEffects(1, nil, nil, &tc.p, nil, [5]int32{}, payload, "")
						got := visibilityEffectsPackets(s)
						if rv != 0 {
							t.Fatalf("return %x", rv)
						}
						for i := 0; i < 32; i++ {
							want := []byte{}
							for j, slot := range []int{1, 7, 31} {
								if i == slot && j < count && tc.inside {
									want = payload
								}
							}
							if !bytes.Equal(got[i], want) {
								t.Fatalf("player %d: %x want %x", i, got[i], want)
							}
						}
						for i := range units {
							if !bytes.Equal(before[i], unsafe.Slice((*byte)(unsafe.Pointer(&units[i])), int(unsafe.Sizeof(units[i])))) {
								t.Fatal("unit mutated")
							}
						}
						s.NetList.ResetAll()
						s.Nox_xxx_netSendFxAllCli_523030(tc.p, payload)
						paired := visibilityEffectsPackets(s)
						for i := range got {
							if !bytes.Equal(got[i], paired[i]) {
								t.Fatalf("existing Go recipient %d differs", i)
							}
						}
						rows = append(rows, visibilityEffectsRow{name, rv, got})
					})
				}
			}
		}
	}
	// A representable coordinate below the exact X upper limit equals the rounded
	// Y upper limit. This independently distinguishes the two precision contracts.
	configure(1)
	units[0].PosVec = types.Pointf{16777216, 16777216}
	pl := units[0].UpdateDataPlayer().Player
	pl.Field3680 = 0
	pl.CameraFollowObj = nil
	pl.Field10 = 3
	pl.Field12 = 3
	for _, tc := range []pointCase{{"large-x", types.Pointf{16777268, 16777216}, true}, {"large-y", types.Pointf{16777216, 16777268}, false}} {
		t.Run(tc.name, func(t *testing.T) {
			s.NetList.ResetAll()
			rv := legacy.PortTestVisibilityEffects(1, nil, nil, &tc.p, nil, [5]int32{}, payload, "")
			got := visibilityEffectsPackets(s)
			if (len(got[1]) != 0) != tc.inside {
				t.Fatalf("packets %x", got[1])
			}
			rows = append(rows, visibilityEffectsRow{tc.name, rv, got})
		})
	}
}

func TestVisibilityEffectsObserverCamera(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	configure(1)
	units[0].PosVec = types.Pointf{100, 200}
	units[1].PosVec = types.Pointf{1000, 2000}
	p := units[0].UpdateDataPlayer().Player
	p.Field10 = 20
	p.Field12 = 30
	var rows []visibilityEffectsRow
	defer func() {
		spellbookCapture(t, "visibility-effects-camera", rows, "aaaa9ccf7961f10bdc6f395ea6d9803bec32c80c1084953d18fe96b764438ab6")
	}()
	for _, mode := range []uint32{0, 1, 2, 3, 4, 5, 128, 255, 256, 257} {
		for _, follow := range []bool{false, true} {
			p.Field3680 = mode
			p.CameraFollowObj = nil
			if follow {
				p.CameraFollowObj = &units[1]
			}
			for _, target := range []int{0, 1} {
				name := fmt.Sprintf("mode%d/follow%t/target%d", mode, follow, target)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					pos := units[target].PosVec
					rv := legacy.PortTestVisibilityEffects(1, nil, nil, &pos, nil, [5]int32{}, []byte{7, 8, 9}, "")
					got := visibilityEffectsPackets(s)
					wantTarget := 0
					if follow && mode&3 != 0 {
						wantTarget = 1
					}
					for i := range got {
						want := []byte{}
						if i == 1 && target == wantTarget {
							want = []byte{7, 8, 9}
						}
						if !bytes.Equal(got[i], want) {
							t.Fatalf("player %d: %x want %x", i, got[i], want)
						}
					}
					s.NetList.ResetAll()
					s.Nox_xxx_netSendFxAllCli_523030(pos, []byte{7, 8, 9})
					paired := visibilityEffectsPackets(s)
					for i := range got {
						if !bytes.Equal(got[i], paired[i]) {
							t.Fatalf("existing Go recipient %d differs", i)
						}
					}
					rows = append(rows, visibilityEffectsRow{name, rv, got})
				})
			}
		}
	}
}
