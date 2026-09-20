//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"slices"
	"testing"
	"unsafe"
)

func TestGameMessageClientMapProgress(t *testing.T) {
	o := newInventoryDisplayOwner(t)
	o.reset(t)
	screen := legacy.PortTestNewScreenEnvironment()
	t.Cleanup(screen.Restore)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	storage := serverConfigOwnBytes(t, 0x5D4594, 1309516, 160)
	dim := legacy.PortTestBindingDimensions()
	w, h := *dim[0], *dim[1]
	t.Cleanup(func() { *dim[0], *dim[1] = w, h })
	*dim[0], *dim[1] = 320, 240
	set, restore := o.c.srv.Server.PortTestMeterStrings(strman.Entry{ID: "guigen.c:Generating", Vals: []strman.Variant{{Str: "Generating"}}}, strman.Entry{ID: "guigen.c:Assembling", Vals: []strman.Variant{{Str: "Assembling"}}}, strman.Entry{ID: "guigen.c:Populating", Vals: []strman.Variant{{Str: "Populating"}}})
	t.Cleanup(restore)
	set(0)
	var loads []string
	oldLoad := legacy.Nox_xxx_gLoadImg
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		loads = append(loads, name)
		return o.images[(len(loads)-1)%len(o.images)]
	}
	binary.LittleEndian.PutUint32(storage[148:], 0x1234)
	if legacy.Nox_xxx_compassGenStrings_4A9C80() != 1 || binary.LittleEndian.Uint32(storage[148:]) != 0 {
		t.Fatal("compass initializer return/reset")
	}
	var names []string
	for i := 1; i <= 4; i++ {
		names = append(names, fmt.Sprintf("Compass%d", i))
	}
	for i := 1; i <= 32; i++ {
		names = append(names, fmt.Sprintf("CompassMainArrow%d", i))
	}
	if !slices.Equal(loads, names) {
		t.Fatal("compass image load order")
	}
	for i := 0; i < 36; i++ {
		off := 128 + 4*i
		if i >= 4 {
			off = 4 * (i - 4)
		}
		if binary.LittleEndian.Uint32(storage[off:]) != uint32(uintptr(o.images[i%len(o.images)].C())) {
			t.Fatal("compass image placement")
		}
	}
	oldCopy := legacy.Nox_video_callCopyBackBuffer_4AD170
	t.Cleanup(func() { legacy.Nox_video_callCopyBackBuffer_4AD170 = oldCopy })
	copies := 0
	legacy.Nox_video_callCopyBackBuffer_4AD170 = func() { copies++ }
	type row struct {
		On, Kind, Frame, Repeat, Value int
		Stored, Next                   uint32
		Text, Pixels                   string
		Draws                          []objectImageDraw
		Sounds                         [][2]int
		Copies                         int
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, kind := range []byte{155, 156, 157} {
			for frame := 0; frame < 32; frame++ {
				for _, value := range []uint16{0, 1, 32768, 65535} {
					for repeat := 0; repeat < 2; repeat++ {
						binary.LittleEndian.PutUint32(connected, uint32(on))
						prior := uint32(value) ^ 0xffff
						if repeat != 0 {
							prior = uint32(value)
						}
						*memmap.PtrUint32(0x5D4594, 1309668) = prior
						*memmap.PtrUint32(0x5D4594, 1309672) = uint32(frame)
						*memmap.PtrPtr(0x5D4594, 1309660) = unsafe.Pointer(alloc.InternCString16("Previous"))
						o.drawTrace = nil
						o.sounds = nil
						copies = 0
						clear(o.pix.Pix)
						data := []byte{kind, byte(value), byte(value >> 8)}
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
						changed := on != 0 && repeat == 0
						wantStored, wantFrame, wantText := prior, uint32(frame), "Previous"
						wantCopies := 0
						if changed {
							wantStored = uint32(value)
							wantFrame = uint32((frame + 1) % 32)
							wantText = []string{"Generating", "Assembling", "Populating"}[kind-155]
							wantCopies = 1
						}
						text := alloc.GoString16(*(**uint16)(memmap.PtrOff(0x5D4594, 1309660)))
						if n != 3 || !bytes.Equal(input, data) || memmap.Uint32(0x5D4594, 1309668) != wantStored || memmap.Uint32(0x5D4594, 1309672) != wantFrame || text != wantText || copies != wantCopies || len(o.drawTrace) != 2*wantCopies {
							t.Fatalf("map progress on%d kind%d frame%d repeat%d value%d: length%d stored%d next%d text%q copies%d draws%d", on, kind, frame, repeat, value, n, memmap.Uint32(0x5D4594, 1309668), memmap.Uint32(0x5D4594, 1309672), text, copies, len(o.drawTrace))
						}
						if changed {
							ids := []int{frame % 4, 4 + frame}
							for i, id := range ids {
								wantImage := o.c.imageRefs[uint32(uintptr(o.images[id%len(o.images)].C()))]
								if o.drawTrace[i].Image != wantImage {
									t.Fatal("map progress compass/arrow selection")
								}
							}
							if !slices.Equal(screen.State()[:4], []uint32{0, 0, 319, 239}) {
								t.Fatal("map progress screen bounds")
							}
						}
						var sound [][2]int
						if changed {
							sound = append(sound, [2]int{897, 50})
						}
						if !slices.Equal(o.sounds, sound) {
							t.Fatal("map progress sound")
						}
						rows = append(rows, row{on, int(kind), frame, repeat, int(value), wantStored, wantFrame, text, effectsPixelHash(o.pix), slices.Clone(o.drawTrace), slices.Clone(o.sounds), copies})
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-map-progress", rows)
}
