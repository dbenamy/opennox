//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
)

func TestGameMessageClientCreatures(t *testing.T) {
	o := newSummonOwner(t)
	originalList := o.c.Objs.List1
	release := func() {
		for o.c.Objs.List1 != originalList {
			if o.c.Objs.List1 == nil {
				t.Fatal("creature list owner lost")
			}
			o.c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(o.c.Objs.List1)
		}
	}
	t.Cleanup(release)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		AcquireOn, LoseOn, High, Quiet, Present int
		Name                                    string
		Acquire, Lose                           summonResult
		AcquireInput, LoseInput                 []byte
		Allies                                  [2]bool
		Marks                                   [2][2]byte
		Counts                                  [2]int
	}
	var rows []row
	for acquireOn := 0; acquireOn < 2; acquireOn++ {
		for loseOn := 0; loseOn < 2; loseOn++ {
			for high := 0; high < 2; high++ {
				for quiet := 0; quiet < 2; quiet++ {
					for present := 0; present < 3; present++ {
						for _, name := range []string{"PortSmallCreature", "PortMediumCreature", "PortLargeCreature"} {
							release()
							o.reset(t)
							o.init(t)
							typ := o.c.Things.TypeByID(name)
							typ.ObjClass = 0
							var dynamic, static *client.Drawable
							for i := 0; i < present; i++ {
								dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(typ.Index(), image.Pt(100, 200))
								if dr == nil {
									t.Fatal("creature fixture allocation")
								}
								dr.NetCode32 = 7
								if i == 0 {
									dynamic = dr
								} else {
									static = dr
									dr.ObjClass = object.Class(0x20400000)
								}
								o.c.Objs.MinimapAdd(dr, 2)
							}
							code := uint16(7 | high<<15)
							typeCode := uint16(typ.Index() | quiet<<15)
							binary.LittleEndian.PutUint32(connected, uint32(acquireOn))
							data := []byte{108, byte(code), byte(code >> 8), byte(typeCode), byte(typeCode >> 8)}
							want := bytes.Clone(data)
							if acquireOn != 0 {
								want[4] &= 0x7f
							}
							if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(108), data); n != 5 || !bytes.Equal(data, want) {
								t.Fatal("creature acquire length/type mutation")
							}
							selected := dynamic
							if high != 0 {
								selected = static
							}
							if acquireOn != 0 {
								if selected == nil {
									selected = o.c.Objs.ByNetCode(7)
									if selected == nil {
										t.Fatal("creature fallback creation")
									}
									if dynamic == nil {
										dynamic = selected
									}
								}
								if selected.Field_71_0&1 == 0 {
									t.Fatalf("creature acquire minimap membership on%d lose%d high%d quiet%d present%d name%s mark%d code%d class%x count%d", acquireOn, loseOn, high, quiet, present, name, selected.Field_71_0, selected.NetCode32, selected.ObjClass, o.c.Objs.Count)
								}
							}
							ally := legacy.PortTestCombatIsAlly(uint32(code))
							if ally != (acquireOn != 0) {
								t.Fatal("creature acquire full-code ally")
							}
							marks := func() [2]byte {
								var v [2]byte
								if dynamic != nil {
									v[0] = dynamic.Field_71_0
								}
								if static != nil {
									v[1] = static.Field_71_0
								}
								return v
							}
							beforeMarks := marks()
							count := o.c.Objs.Count
							r := row{AcquireOn: acquireOn, LoseOn: loseOn, High: high, Quiet: quiet, Present: present, Name: name, Acquire: o.snapshot("acquire", 5), AcquireInput: bytes.Clone(data)}
							binary.LittleEndian.PutUint32(connected, uint32(loseOn))
							lost := []byte{109, byte(code), byte(code >> 8)}
							want = bytes.Clone(lost)
							if loseOn != 0 {
								want[2] &= 0x7f
							}
							if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(109), lost); n != 3 || !bytes.Equal(lost, want) {
								t.Fatal("creature loss length/code mutation")
							}
							afterMarks := marks()
							expected := beforeMarks
							if loseOn != 0 {
								expected[0] &^= 1
							}
							if afterMarks != expected || o.c.Objs.Count != count {
								t.Fatal("creature loss dynamic-only minimap update")
							}
							afterAlly := legacy.PortTestCombatIsAlly(uint32(code))
							wantAlly := ally && (loseOn == 0 || high != 0)
							if afterAlly != wantAlly {
								t.Fatal("creature loss masked ally lookup")
							}
							r.Lose = o.snapshot("lose", 3)
							r.LoseInput = bytes.Clone(lost)
							r.Allies = [2]bool{ally, afterAlly}
							r.Marks = [2][2]byte{beforeMarks, afterMarks}
							r.Counts = [2]int{count, o.c.Objs.Count}
							rows = append(rows, r)
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-progress-creatures", rows)
}
