//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"strings"
	"testing"
	"unsafe"
)

func TestResourceDefinitionsNames(t *testing.T) {
	var rows []resourceParserRow
	for _, fill := range []byte{0, 0xa5} {
		for _, name := range []string{"Monster", "n", strings.Repeat("x", 63)} {
			row := resourceParse(t, "skull", " \t"+name+" extra", fill)
			want := bytes.Repeat([]byte{fill}, 256)
			resourcePut32(want, 12, 0)
			copy(want[16:], name)
			want[16+len(name)] = 0
			if row.Return != 1 || !bytes.Equal(row.Data, want) {
				t.Fatal("skull", name)
			}
			rows = append(rows, row)
		}
		for _, name := range []string{sound.ID(226).String(), sound.ID(925).String(), "unknown_sound"} {
			for _, kind := range []string{"audio", "spawn"} {
				input := name
				off := 0
				if kind == "spawn" {
					input = "Monster " + name
					off = 128
				}
				row := resourceParse(t, kind, input, fill)
				want := bytes.Repeat([]byte{fill}, 256)
				resourcePut32(want, off, uint32(sound.ByName(name)))
				ret := 1
				if kind == "audio" && sound.ByName(name) == 0 {
					ret = 0
				}
				if kind == "spawn" {
					copy(want, "Monster\x00")
				}
				if row.Return != ret || !bytes.Equal(row.Data, want) {
					t.Fatal(kind, name)
				}
				rows = append(rows, row)
			}
		}
		for _, input := range []string{"", sound.ID(226).String(), sound.ID(226).String() + " " + sound.ID(925).String(), "unknown_sound unknown_sound", sound.ID(226).String() + "\t" + sound.ID(925).String()} {
			row := resourceParse(t, "trigger", input, fill)
			want := bytes.Repeat([]byte{fill}, 256)
			fields := strings.FieldsFunc(input, func(r rune) bool { return r == ' ' })
			for i := 0; i < len(fields) && i < 2; i++ {
				resourcePut32(want, 36+4*i, uint32(sound.ByName(fields[i])))
			}
			if row.Return != 1 || !bytes.Equal(row.Data, want) {
				t.Fatal("trigger", input)
			}
			rows = append(rows, row)
		}
	}
	spellbookCapture(t, "resource-definitions-names", rows, "4c05d91e19a8ee0d2e8cf662db9f091ba7a8b715928a3d998a19174288176b3f")
}
func TestResourceDefinitionsWands(t *testing.T) {
	o := newReliableReportsOwner(t)
	var rows []map[string]any
	for _, fps := range []uint32{15, 30, 60} {
		o.s.SetTickRate(fps)
		for _, fill := range []byte{0, 0xa5} {
			for _, charges := range []int{-1, 0, 1, 255, 256} {
				for _, rate := range []int{-3, 1, 2, 7} {
					for _, kind := range []string{"wand", "wandcast"} {
						name := "SPELL_MAGIC_MISSILE"
						if spell.ParseID(name) == 0 {
							t.Fatal("invalid fixture spell", name)
						}
						input := fmt.Sprintf("%d %d %s", charges, rate, name)
						if kind == "wand" {
							input = fmt.Sprintf("%d Missile %d MULTI_SHOT %s", charges, rate, sound.ID(226).String())
						}
						row := resourceParse(t, kind, input, fill)
						want := bytes.Repeat([]byte{fill}, 256)
						want[108] = byte(charges)
						want[109] = byte(charges)
						resourcePut32(want, 112, 100)
						resourcePut32(want, 100, uint32(int32(fps)/int32(rate)))
						if kind == "wandcast" {
							resourcePut32(want, 0, 1)
							resourcePut32(want, 92, uint32(spell.ParseID(name)))
						} else {
							resourcePut32(want, 0, 0)
							copy(want[4:], "Missile\x00")
							resourcePut32(want, 84, 0)
							want[96] |= 1
							resourcePut32(want, 88, uint32(sound.ByName(sound.ID(226).String())))
						}
						if row.Return != 1 || !bytes.Equal(row.Data, want) {
							t.Fatal(kind, fps, input, row.Data[:116], want[:116])
						}
						rows = append(rows, map[string]any{"fps": fps, "row": row})
					}
				}
			}
		}
	}
	spellbookCapture(t, "resource-definitions-wands", rows, "3444ef187eef4184171d60689e31d99be7c4f839c8b981cb39413918a7f121b2")
}
func TestResourceDefinitionsRegistrations(t *testing.T) {
	resourceDamageNames(t)
	newReliableReportsOwner(t)
	routes := []struct{ group, name, kind, input string }{
		{"update", "PushUpdate", "push", "1.5 3"}, {"update", "TriggerUpdate", "trigger", sound.ID(226).String()},
		{"update", "ToggleUpdate", "trigger", sound.ID(925).String()}, {"update", "LoopAndDamageUpdate", "triple", "1 2 3"},
		{"update", "LifetimeUpdate", "lifetime", "-1"}, {"update", "SkullUpdate", "skull", "Monster"},
		{"use", "WandUse", "wand", "5 Missile 2 MULTI_SHOT"}, {"use", "WandCastUse", "wandcast", "5 2 SPELL_MAGIC_MISSILE"},
		{"death", "CreateObjectDie", "spawn", "Monster " + sound.ID(226).String()},
		{"death", "SpawnObjectDie", "spawn", "Monster " + sound.ID(925).String()},
		{"collide", "ProjectileCollide", "projectile", "7"}, {"collide", "ProjectileSparkCollide", "projectile", "8"},
		{"collide", "DamageCollide", "damage", "12 BLADE"}, {"collide", "ManaDrainCollide", "mana", "257"},
		{"collide", "SparkExplosionCollide", "spark", "255"}, {"collide", "WallReflectCollide", "projectile", "3"},
		{"collide", "WallReflectSparkCollide", "projectile", "4"}, {"collide", "PixieCollide", "projectile", "5"},
		{"collide", "AudioEventCollide", "audio", sound.ID(226).String()},
		{"collide", "MonsterArrowCollide", "arrow", "1 99"}, {"collide", "YellowStarShotCollide", "projectile", "9"},
	}
	var rows []map[string]any
	for _, r := range routes {
		for _, fill := range []byte{0, 0xa5} {
			baseline := resourceParse(t, r.kind, r.input, fill)
			buf, free := alloc.Make([]byte{}, 256)
			for i := range buf {
				buf[i] = fill
			}
			typ := &server.ObjectType{UpdateData: unsafe.Pointer(&buf[0]), CollideData: unsafe.Pointer(&buf[0]), DeathData: unsafe.Pointer(&buf[0])}
			typ.UseData.Ptr = unsafe.Pointer(&buf[0])
			err := server.PortTestResourceRegisteredParser(r.group, r.name, typ, strings.Split(r.input, " "))
			if (err == nil) != (baseline.Return != 0) || !bytes.Equal(buf, baseline.Data) {
				t.Fatal("registration", r.name, err)
			}
			rows = append(rows, map[string]any{"group": r.group, "name": r.name, "fill": fill, "ok": err == nil, "data": append([]byte(nil), buf...), "word": binary.LittleEndian.Uint32(buf)})
			free()
		}
	}
	spellbookCapture(t, "resource-definitions-registrations", rows, "1140e7333f923de5511bd0db84f4aabebacaadb4b199f9188cff788bfc8216df")
}

func resourceDamageNames(t *testing.T) []string {
	t.Helper()
	table := unsafe.Slice(memmap.PtrUint32(0x587000, 200728), 18)
	saved := append([]uint32(nil), table...)
	t.Cleanup(func() { copy(table, saved) })
	names := make([]string, 18)
	for i := range names {
		names[i] = fmt.Sprintf("TYPE%d", i)
		if i == 0 {
			names[i] = "BLADE"
		}
		p, free := alloc.CString(names[i])
		t.Cleanup(free)
		table[i] = uint32(uintptr(unsafe.Pointer(p)))
	}
	return names
}
func TestResourceDefinitionsDamage(t *testing.T) {
	names := resourceDamageNames(t)
	var rows []resourceParserRow
	for i, name := range append(names, "unknown", "blade") {
		typ := i
		if typ > 18 {
			typ = 18
		}
		for _, value := range []int{-1, 0, 255, 256} {
			for _, fill := range []byte{0, 0xa5} {
				row := resourceParse(t, "damage", fmt.Sprintf("%d %s", value, name), fill)
				want := bytes.Repeat([]byte{fill}, 256)
				want[0] = byte(value)
				resourcePut32(want, 4, uint32(typ))
				ret := 1
				if typ == 18 {
					ret = 0
				}
				if row.Return != ret || !bytes.Equal(row.Data, want) {
					t.Fatal("damage", name, value)
				}
				rows = append(rows, row)
			}
		}
	}
	spellbookCapture(t, "resource-definitions-damage", rows, "cd08a99e8925b080d9cba84ad3894dd2de4ed3868958782fa108c16c869d8193")
}
