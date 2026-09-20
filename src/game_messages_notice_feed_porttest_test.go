//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestGameMessageNoticeDeathFeed(t *testing.T) {
	type record struct {
		Name    string
		Packet  string
		Console []string
		Ring    [6]uint32
	}
	var rows []record
	for _, killer := range []uint16{0, 7, 999} {
		for _, assist := range []uint16{0, 8, 999} {
			for _, victim := range []uint16{0, 9, 999} {
				name := fmt.Sprintf("killer=%d/assist=%d/victim=%d", killer, assist, victim)
				t.Run(name, func(t *testing.T) {
					o := newCombatOverlayOwner(t)
					for i, n := range []string{"Killer", "Helper", "Victim"} {
						o.players[i].Active = 1
						o.players[i].PlayerInd = byte(i)
						o.players[i].NetCodeVal = uint32(7 + i)
						o.players[i].SetName(n)
					}
					var packet [11]byte
					packet[0], packet[1] = 169, 14
					binary.LittleEndian.PutUint16(packet[2:], killer)
					binary.LittleEndian.PutUint16(packet[4:], assist)
					binary.LittleEndian.PutUint16(packet[6:], victim)
					binary.LittleEndian.PutUint16(packet[8:], 17)
					packet[10] = 2
					// A missing participant must not reuse the previous notification's name.
					prime := packet
					binary.LittleEndian.PutUint16(prime[6:], 9)
					if n := legacy.PortTestGameNotice(prime[:]); n != 11 {
						t.Fatal("prime consumed length", n)
					}
					o.console = nil
					*o.words["feedWrite"], *o.words["feedRead"] = 0, 0
					before := bytes.Clone(packet[:])
					if n := legacy.PortTestGameNotice(packet[:]); n != 11 || !bytes.Equal(packet[:], before) {
						t.Fatal("notice consumed length/input", n)
					}
					attacker := "A:nature"
					if killer == 7 {
						attacker = "A:Killer"
						if assist == 8 {
							attacker += " + Helper"
						}
					}
					who := ""
					if victim == 9 {
						who = "V:Victim"
					}
					want := []string{who + " " + attacker}
					if !reflect.DeepEqual(o.console, want) {
						t.Fatalf("notification %q want %q", o.console, want)
					}
					ring := *(*[6]uint32)(memmap.PtrOff(0x5D4594, 1201428))
					if ring != [6]uint32{uint32(killer), uint32(assist), uint32(victim), 17, 2, 123} || *o.words["feedWrite"] != 1 || *o.words["feedRead"] != 0 {
						t.Fatal("feed storage")
					}
					rows = append(rows, record{name, fmt.Sprintf("%x", packet), append([]string(nil), o.console...), ring})
				})
			}
		}
	}
	gameMessageCapture(t, "game-notice-death-feed", rows)
}
