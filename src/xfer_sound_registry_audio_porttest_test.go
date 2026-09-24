//go:build porttest

package opennox

import (
	"fmt"
	"testing"

	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Invoke the real sound algorithms through damageDefault and observe the normal
// deferred audio queue. This complements the replaceable-handler dispatch tests.
func TestXferSoundRegistryAudioOwner(t *testing.T) {
	s := newObjectXferOwner(t)
	u, source := newObjectXferSimple(t, s), newObjectXferSimple(t, s)
	u.ObjClass = 0
	source.ObjClass = 0
	hp, free := alloc.New(server.HealthData{})
	u.HealthData = hp
	t.Cleanup(func() { s.PortTestCombatAudioReset(); u.HealthData = nil; free() })
	for _, slot := range []string{"nil", "DefaultDamageSound", "PlayerDamageSound"} {
		for _, material := range []uint16{1, 8} {
			for _, state := range []uint32{0, 1} {
				for _, otherMaterial := range []int{-1, 0, 8, 0x4000} {
					t.Run(fmt.Sprintf("%s-m%d-state%d-other%d", slot, material, state, otherMaterial), func(t *testing.T) {
						*hp = server.HealthData{Cur: 80, Field2: 100, Max: 100}
						// The owner replaces prior damage kind before invoking sound.
						u.ObjFlags = 0
						u.Material = material
						u.Field131 = state
						u.DamageSound = nil
						if slot != "nil" {
							u.DamageSound = legacy.PortTestDamageSoundRegistryPointer(slot)
						}
						var other *server.Object
						if otherMaterial >= 0 {
							other = source
							source.Material = uint16(otherMaterial)
						}
						s.PortTestCombatAudioReset()
						if got := legacy.PortTestDamageSoundOwner(u, other, nil); got != 1 || hp.Cur != 68 || u.Field131 != 0 {
							t.Fatal("real sound damage result", got, hp.Cur)
						}
						var want sound.ID
						if slot == "PlayerDamageSound" {
							if otherMaterial == 8 {
								want = 857
							}
						} else if other != nil && otherMaterial != 0x4000 {
							if material == 8 {
								want = sound.SoundHitWoodBreakable
							} else if otherMaterial == 8 {
								want = 857
							}
						}
						events := s.PortTestCombatAudioSnapshot()
						if want == 0 {
							if len(events) != 0 {
								t.Fatal("unexpected sound", events)
							}
							return
						}
						if len(events) != 1 || events[0].ID != want || events[0].Obj != u || events[0].ByPos || events[0].Kind != 0 || events[0].Code != 0 {
							t.Fatal("damage sound event", events, "want", want)
						}
					})
				}
			}
		}
	}
}
