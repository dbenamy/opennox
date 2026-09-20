//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageClientSessionAudio(t *testing.T) {
	type row struct {
		On, Kind, ID, EncodedVolume, Pan, Return       int
		Event                                          bool
		State, Volume, CurrentPan, TargetPan, Priority uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []int{166, 167} {
			for _, id := range []int{1, 1022} {
				for _, vol := range []int{0, 25, 50, 63} {
					for _, pan := range []int{-128, -51, -50, 0, 50, 51, 127} {
						t.Run(fmt.Sprintf("on%d/kind%d/id%d/volume%d/pan%d", on, kind, id, vol, pan), func(t *testing.T) {
							o := newAudioEventsOwner(t, 1)
							metadata := o.metadata(id, 0, 1)
							connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
							binary.LittleEndian.PutUint32(connected, uint32(on))
							packed := uint16(id | vol<<10)
							data := []byte{byte(kind), byte(pan), byte(packed), byte(packed >> 8)}
							before := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
							if n != 4 || !bytes.Equal(data, before) {
								t.Fatal("audio message length/input", n)
							}
							root := uint32(uintptr(memmap.PtrOff(0x5D4594, 840612)))
							p := audioStreamWords(root, 3)[0]
							r := row{On: on, Kind: kind, ID: id, EncodedVolume: vol, Pan: pan, Return: n, Event: p != root}
							if r.Event != (on != 0) {
								t.Fatal("audio connection gate", r)
							}
							if r.Event {
								w := audioStreamWords(p, 144)
								r.State, r.Volume, r.CurrentPan, r.TargetPan, r.Priority = w[7], w[47], w[63], w[64], w[75]
								volume := 2 * vol
								if volume > 100 {
									volume = 100
								}
								clamped := pan
								if clamped < -50 {
									clamped = -50
								}
								if clamped > 50 {
									clamped = 50
								}
								priority := uint32(0)
								if kind == 167 {
									priority = 2
								}
								if w[9] != metadata || r.State != 1 || r.Volume != uint32(volume*163)<<16 || r.CurrentPan != 8192<<16 || r.TargetPan != uint32(8192+clamped*8192/50)<<16 || r.Priority != priority {
									t.Fatal("audio packed fields/signed pan/priority", r)
								}
							}
							rows = append(rows, r)
						})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-audio", rows)
}
