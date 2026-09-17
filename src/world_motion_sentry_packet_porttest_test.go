//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func TestWorldMotionSentryPacket(t *testing.T) {
	o := newCollisionCoreOwner(t)
	u := newObjectXferSimple(t, o.s)
	saved := *u
	t.Cleanup(func() { *u = saved })
	scratch, restore := legacy.PortTestWorldMotionRoundScratch()
	t.Cleanup(restore)
	type row struct {
		Input   [4]uint32
		Full    bool
		Return  int32
		Packet  []byte
		Scratch [2]uint32
	}
	var rows []row
	for _, v := range []float32{-100, -0.5, 0, 0.49, 0.5, 1.5, 2.5, 65534.5, 65535.5, 65536, 8388607, 8388608, 16777216} {
		for _, full := range []bool{false, true} {
			o.s.NetList.ResetAll()
			*scratch = [4]uint32{}
			if full && !o.s.NetList.AddToMsgListCli(1, netlist.Kind1, bytes.Repeat([]byte{0x44}, 2048)) {
				t.Fatal("fill sentry recipient buffer")
			}
			u.PosVec = types.Pointf{v, v + 1}
			u.Pos39 = types.Pointf{v + 2, v + 3}
			rv := legacy.PortTestWorldMotionSentryPacket(1, u)
			packet := o.s.NetList.CopyPacketsA(1, netlist.Kind1)
			if full {
				if rv != 0 || !bytes.Equal(packet, bytes.Repeat([]byte{0x44}, 2048)) {
					t.Fatal("sentry report capacity contract", rv, len(packet))
				}
			} else {
				if rv != 1 || len(packet) != 9 || packet[0] != 0x95 {
					t.Fatal("sentry report acceptance", rv, packet)
				}
				if v < 65534 {
					for i, x := range []float32{v, v + 1, v + 2, v + 3} {
						want := uint16(0)
						if x >= 0 {
							want = uint16(math.RoundToEven(float64(x)))
						}
						if binary.LittleEndian.Uint16(packet[1+i*2:]) != want {
							t.Fatal("sentry coordinate ties-even/clamp", v, i, packet, want)
						}
					}
				}
			}
			rows = append(rows, row{[4]uint32{math.Float32bits(v), math.Float32bits(v + 1), math.Float32bits(v + 2), math.Float32bits(v + 3)}, full, rv, append([]byte(nil), packet...), [2]uint32{scratch[2], scratch[3]}})
		}
	}
	spellbookCapture(t, "world-motion-sentry-packet", rows, "")
}
