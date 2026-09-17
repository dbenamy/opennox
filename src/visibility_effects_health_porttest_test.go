//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestVisibilityEffectsHealthThrottle(t *testing.T) {
	s, units, configure := visibilityEffectsPlayers(t)
	configure(1)
	u := &units[0]
	u.NetCode = 0x1234
	hp, freeHP := alloc.New(server.HealthData{})
	u.HealthData = hp
	t.Cleanup(func() { u.HealthData = nil; freeHP() })
	pl := u.UpdateDataPlayer().Player
	*(*byte)(unsafe.Add(unsafe.Pointer(pl), 2251)) = 1
	reset, snapshot, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	half := func(off int) *uint16 { return (*uint16)(unsafe.Add(u.UpdateData, off)) }
	stamp := func(off int) *uint32 { return (*uint32)(unsafe.Add(unsafe.Pointer(pl), off)) }
	type row struct {
		Name    string
		Return  uint32
		Last    [2]uint16
		Times   [2]uint32
		Packets []legacy.PortTestShopPacketResult
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "visibility-effects-health-throttle", rows, "c698c3302555dc05f4419f73e4870d7bc824ec17a715d5c135ab776e62990512")
	}()
	for _, fps := range []uint32{0, 3, 30, 0x80000000, 0xffffffff} {
		for _, frame := range []uint32{0, 1, 65535, 65536, 0x80000000, 0xffffffff} {
			for _, elapsed := range []uint32{0, 7, 8, 0xffffffff} {
				for _, values := range [][4]uint16{{100, 100, 100, 100}, {101, 100, 101, 100}, {110, 100, 110, 100}, {99, 100, 99, 100}, {0, 65535, 65535, 0}, {100, 100, 32768, 32768}, {100, 100, 65535, 65535}} {
					for _, max := range []uint16{0, 9, 100, 65535} {
						name := fmt.Sprintf("fps%x/frame%x/elapsed%x/values%v/max%d", fps, frame, elapsed, values, max)
						t.Run(name, func(t *testing.T) {
							reset()
							s.SetFrame(frame)
							s.SetTickRate(fps)
							hp.Cur = values[0]
							hp.Max = max
							*half(10) = values[1]
							*half(4) = values[2]
							*half(6) = values[3]
							*half(8) = max
							*stamp(2176) = frame - elapsed
							*stamp(2180) = frame - elapsed
							delta := func(a, b uint16) int {
								d := int(a) - int(b)
								if d < 0 {
									return -d
								}
								return d
							}
							health := values[0] != values[1] && (delta(values[0], values[1]) >= int(max)/10 || elapsed > uint32(int32(fps)>>2))
							manaChanged := int(values[3]) != int(int16(values[2]))
							mana := manaChanged && (delta(values[2], values[3]) >= int(max)/10 || elapsed > fps>>2)
							wantReturn := uint32(int32(int16(values[2])))
							if manaChanged {
								wantReturn = 0
								if mana {
									wantReturn = uint32(int32(int16(frame)))
								}
							}
							rv := legacy.PortTestVisibilityEffects(13, u, nil, nil, nil, [5]int32{}, nil, "")
							if rv != wantReturn {
								t.Fatalf("return%x want%x", rv, wantReturn)
							}
							wantLast := [2]uint16{values[1], values[3]}
							wantTimes := [2]uint32{frame - elapsed, frame - elapsed}
							var wantPackets [][]byte
							if mana {
								wantLast[1] = values[2]
								wantTimes[1] = frame
								wantPackets = append(wantPackets, []byte{69, 0x34, 0x12, byte(values[2]), byte(values[2] >> 8)})
							}
							if health {
								wantLast[0] = values[0]
								wantTimes[0] = frame
								wantPackets = append(wantPackets, []byte{67, byte(values[0]), byte(values[0] >> 8)})
							}
							got := row{name, rv, [2]uint16{*half(10), *half(6)}, [2]uint32{*stamp(2176), *stamp(2180)}, snapshot()}
							if got.Last != wantLast || got.Times != wantTimes {
								t.Fatalf("state%v/%v want%v/%v", got.Last, got.Times, wantLast, wantTimes)
							}
							if len(got.Packets) != len(wantPackets) {
								t.Fatalf("packets%d want%d", len(got.Packets), len(wantPackets))
							}
							for i, p := range got.Packets {
								if p.Recipient != 1 || !bytes.Equal(p.Data, wantPackets[i]) {
									t.Fatalf("packet%+v want%x", p, wantPackets[i])
								}
							}
							rows = append(rows, got)
						})
					}
				}
			}
		}
	}
}
