//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"math"
	"testing"
)

func TestMatchRosterTimerCoordinateWidths(t *testing.T) {
	o := newMatchRosterOwner(t)
	type row struct {
		Name    string
		Packets [][]byte
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "match-roster-timer-coordinates", rows, "61a0f9c22f7cec236ced6db34b01b7e404081083f30274237e07f71547ab0ddc")
	}()
	bits := []uint32{0, 0x80000000, 0x3fc00000, 0xbfc00000, 0x477fffff, 0x47800000, 0x4f000001, 0xcf000001, 0x5effffff, 0x5f000000, 0xdf000000, 0xdf000001, 0x7f800000, 0xff800000, 0x7fc12345}
	expected := func(bits uint32) uint16 {
		x := float64(math.Float32frombits(bits))
		if math.IsNaN(x) || math.IsInf(x, 0) || x >= 0x1p63 || x < -0x1p63 {
			return 0
		}
		return uint16(int32(math.Mod(math.Trunc(x), 65536)))
	}
	for _, x := range bits {
		for _, y := range bits {
			name := fmt.Sprintf("x%x/y%x", x, y)
			t.Run(name, func(t *testing.T) {
				o.reset()
				o.s.PortTestCombatAudioReset()
				restore := noxflags.PortTestGameFlags(0)
				defer restore()
				noxServer.flag3592 = true
				*memmap.PtrUint64(0x5D4594, 3468) = 0
				*memmap.PtrUint32(0x587000, 4660) = 1
				for i := range o.units {
					p := o.units[i].UpdateDataPlayer().Player.C()
					objectXferSetWord(p, 3632, x)
					objectXferSetWord(p, 3636, y)
				}
				matchRosterCall("check-limit", nil, nil, nil)
				packets := visibilityEffectsPackets(o.s)
				for i := range o.units {
					p := o.units[i].UpdateDataPlayer().Player
					data := packets[p.PlayerInd]
					if len(data) != 5 || data[0] != 154 || binary.LittleEndian.Uint16(data[1:]) != expected(x) || binary.LittleEndian.Uint16(data[3:]) != expected(y) {
						t.Fatal("64-bit coordinate conversion", data)
					}
				}
				rows = append(rows, row{name, packets})
			})
		}
	}
}
