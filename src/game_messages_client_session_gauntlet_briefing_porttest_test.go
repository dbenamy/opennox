//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"reflect"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionGauntletBriefing(t *testing.T) {
	o := newBriefingOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	words, restore := legacy.PortTestMinimapWords()
	t.Cleanup(restore)
	packet, free := alloc.New([90]byte{})
	t.Cleanup(free)
	type row struct {
		On, Kind, Flag, Named, Mask int
		State                       briefingResult
		Return                      int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{12, 13, 14} {
			for _, flag := range []byte{0, 1, 2, 3, 128, 255} {
				for named := 0; named < 2; named++ {
					for _, mask := range []int{0, 1, 21, 63} {
						clear(packet[:])
						packet[0], packet[1], packet[4] = 240, byte(kind), flag
						binary.LittleEndian.PutUint16(packet[2:], 65535)
						length := 69
						if kind == 12 {
							length = 90
							binary.LittleEndian.PutUint16(packet[4:], 32768)
							j := 0
							for i := 0; i < 6; i++ {
								if mask&(1<<i) == 0 {
									continue
								}
								r := packet[6+j*14:]
								binary.LittleEndian.PutUint16(r, uint16(100+i))
								for f := 1; f <= 4; f++ {
									binary.LittleEndian.PutUint16(r[f*2:], uint16(1000*i+f))
								}
								binary.LittleEndian.PutUint32(r[10:], uint32(i*31+7))
								j++
							}
						} else {
							image := "MissingBriefingImage"
							if named != 0 {
								image = "CustomBriefingImage"
								copy(packet[37:], "Briefing:Custom")
							}
							copy(packet[5:37], image)
						}
						setup := func() {
							o.resetBriefing(t)
							o.constructBriefing(t)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							*words["gui"] = 0xcafe1234
						}
						setup()
						if on != 0 {
							show := uintptr(1)
							if kind == 13 {
								show = uintptr(flag & 1)
							}
							legacy.PortTestBriefing(kind-8, txptr(unsafe.Pointer(packet)), show, 0)
						}
						*words["gui"] = 0
						want := o.snapshotBriefing(t, kind, 0)
						setup()
						input := bytes.Clone(packet[:])
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(240), packet[:length])
						got := o.snapshotBriefing(t, kind, 0)
						if n != length || !bytes.Equal(input, packet[:]) || *words["gui"] != 0 || !reflect.DeepEqual(got, want) {
							t.Fatal("quest briefing route", on, kind, flag, named, mask, n)
						}
						rows = append(rows, row{on, kind, int(flag), named, mask, got, n})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-gauntlet-briefing", rows)
}
