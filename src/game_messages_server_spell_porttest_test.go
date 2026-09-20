//go:build porttest

package opennox

import (
	"encoding/binary"
	"reflect"
	"testing"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestGameMessageServerSpellRequest(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, ids := range [][5]int32{{1}, {1, 1}, {1, 0, 1}, {1, 1, 1, 1, 1}, {137}, {-1}} {
		for _, flags := range []uint32{0, 1, 2, 3, 0x100} {
			for _, game := range []uint32{0, 128, 2048, 2176} {
				for _, dead := range []bool{false, true} {
					for _, self := range []byte{0, 1, 255} {
						makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
							s := spellLifeBase(13)
							p := s.Callbacks.Shop
							o := p.TemporaryUpdates.World.Objectives
							a := o.Attack
							p.Sequence = []legacy.PortTestShopAction{{Op: 1477}}
							s.Lifecycle.GameFlags |= game
							o.PlayerDataWords[0][3680] = flags
							a.ActorWords = map[int]uint32{16: 0}
							if dead {
								a.ActorWords[16] = 0x4000
							}
							count := int32(0)
							for _, id := range ids {
								if id != 0 {
									count++
								}
							}
							insert := flags&3 == 0 && game&128 == 0 && (!dead || game&2048 != 0)
							a.Controls.GameMessageSpell = &legacy.PortTestServerSpellSpec{IDs: ids, Warnings: byte(flags & 3), Insert: insert, Count: count, Self: int32(self), WantBook: insert && ids[0] == 1}
							if dispatch {
								data := make([]byte, 22)
								data[0], data[21] = 121, self
								for i, id := range ids {
									binary.LittleEndian.PutUint32(data[1+4*i:], uint32(id))
								}
								a.Controls.GameMessage, a.Controls.GameMessageLength = data, 22
							}
							return s
						}
						direct = append(direct, makeCase(false))
						messages = append(messages, makeCase(true))
					}
				}
			}
		}
	}
	want, got := controlsRun(t, direct), controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("spell request differs from qualified queue operation: case %d", i)
		}
	}
	interactionCapture(t, "game-server-spell-request", got)
}

func TestGameMessageServerSpellFriendlyTarget(t *testing.T) {
	var direct, messages []legacy.PortTestRoamSpec
	for _, special := range []bool{false, true} {
		for _, target := range []bool{false, true} {
			for _, enemy := range []bool{false, true} {
				for _, quest := range []bool{false, true} {
					for _, count := range []int32{1, 2} {
						makeCase := func(dispatch bool) legacy.PortTestRoamSpec {
							s := spellLifeBase(13)
							p := s.Callbacks.Shop
							a := p.TemporaryUpdates.World.Objectives.Attack
							p.Sequence = []legacy.PortTestShopAction{{Op: 1477}}
							if quest {
								s.Lifecycle.GameFlags |= 4096
							}
							if target {
								ref := 1 // The caster itself is unambiguously friendly.
								if enemy {
									ref = 4
									p.Items[1].Class = uint32(object.ClassMonsterGenerator)
								}
								a.UpdateRefs = map[int]int{288: ref}
							}
							flags := uint32(0x1000000)
							if special {
								flags |= 32
							}
							a.Controls.SpellLifecycle.Definitions = []server.PortTestSpellLifecycleDef{{Index: 1, Flags: flags, Valid: true, Enabled: true, ManaCost: 10, Phonemes: []int{0}, Sounds: [3]int{101, 301, 501}}}
							ids := [5]int32{1}
							if count == 2 {
								ids[1] = 1
							}
							insert := !special || !target || enemy || quest || count == 2
							a.Controls.GameMessageSpell = &legacy.PortTestServerSpellSpec{IDs: ids, Insert: insert, Count: count, WantBook: insert}
							if dispatch {
								data := make([]byte, 22)
								data[0] = 121
								for i, id := range ids {
									binary.LittleEndian.PutUint32(data[1+4*i:], uint32(id))
								}
								a.Controls.GameMessage, a.Controls.GameMessageLength = data, 22
							}
							return s
						}
						direct = append(direct, makeCase(false))
						messages = append(messages, makeCase(true))
					}
				}
			}
		}
	}
	want, got := controlsRun(t, direct), controlsRun(t, messages)
	gameMessageStateGuards(t, want)
	gameMessageStateGuards(t, got)
	for i := range got {
		if !reflect.DeepEqual(got[i], want[i]) {
			t.Fatalf("spell friendly-target admission: case %d", i)
		}
	}
	interactionCapture(t, "game-server-spell-friendly-target", got)
}
