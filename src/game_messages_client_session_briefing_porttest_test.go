//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/player"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionBriefing(t *testing.T) {
	o := newBriefingWindowOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Class, Chapter, Begin int
		Seen, After               uint32
		State                     briefingWindowResult
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for class := 0; class < 3; class++ {
			for _, chapter := range []int{0, 1, 10, 254, 255} {
				for _, begin := range []int{0, 1, 127, 128, 255} {
					for _, seen := range []uint32{0, 1} {
						var after *uint32
						setup := func() {
							o.resetWindow(t)
							o.constructBriefing(t)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							p := &o.players[0]
							legacy.Set_dword_8531A0_2576(p)
							p.Info().SetPlayerClass(player.Class(class))
							after = (*uint32)(unsafe.Add(p.C(), 4408+4*(chapter%11)))
							*after = seen
						}
						setup()
						if on != 0 {
							legacy.PortTestBriefingWindow(7, uintptr(chapter), uintptr(begin), 0, 0)
						}
						want := o.windowCapture(t, 0, 0)
						wantAfter := *after
						setup()
						data := []byte{214, byte(chapter), byte(begin)}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(214), data)
						got := o.windowCapture(t, 0, 0)
						if n != 3 || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) || *after != wantAfter {
							t.Fatal("chapter-end route", on, class, chapter, begin, seen)
						}
						if chapter < 11 {
							expected := seen
							if on != 0 && begin == 0 {
								expected = 1
							}
							if *after != expected {
								t.Fatal("chapter loss seen state")
							}
						}
						rows = append(rows, row{on, class, chapter, begin, seen, *after, got})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-briefing", rows)
}
