//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func aiSpellBase(op int) legacy.PortTestRoamSpec {
	b := math.Float32bits
	s := mainBase(0)
	s.Op = 14
	s.MonsterState.ObjFlags = 0x1000004
	s.MonsterState.Status = 0x20
	s.MonsterState.Frame = 128
	s.Spells = &legacy.PortTestAISpellSpec{Op: op, Spell: 41, InversionRange: 200, TargetClass: 2, TargetFlags: 4, TargetCur: 20, TargetMax: 100, SecondCur: 10, SecondMax: 100, Accuracy: b(.5), ArgPos: [2]uint32{b(140), b(100)}, TargetArg: true, SummonAllowed: true}
	s.Main.Words = append(s.Main.Words, [2]uint32{1448, 2 | (5 << 16)}, [2]uint32{1456, 3 | (7 << 16)}, [2]uint32{1464, 4 | (9 << 16)}, [2]uint32{1472, 5 | (11 << 16)}, [2]uint32{1480, 6 | (13 << 16)})
	s.Spells.Permissions = [][2]uint32{{41, 0xf8000000}}
	s.Spells.Definitions = []server.PortTestSpellClassDef{{Index: 41, Flags: 16, Valid: true}}
	return s
}
func aiSpellCorpus() []legacy.PortTestRoamSpec {
	b := math.Float32bits
	var all []legacy.PortTestRoamSpec
	for n := 0; n < 256; n++ {
		for op := 0; op < 15; op++ {
			s := aiSpellBase(op)
			v, m := s.Spells, s.MonsterState
			s.Seed = n + 1
			m.Stack = int8(n%21 + 1)
			m.Action = []uint32{1, 4, 7, 18, 19, 20, 21, 24, 26}[n%9]
			m.Frame = []uint32{0, 1, 31, 32, 127, 128, 0x7fffffff, 0x80000000, 0xffffffff}[n%9]
			m.FPS = []uint32{1, 30, 60}[n%3]
			m.ObjFlags = []uint32{4, 0x1000004}[n%2]
			m.Status = []uint32{0, 0x20, 0x820, 0x1020, 0x1820}[n%5]
			m.HealthCur = uint16(n % 100)
			m.HealthMax = uint16(100 + n%2)
			v.Spell = []int{-1, 0, 1, 41, 74, 75, 90, 114, 115, 136, 137}[n%11]
			v.Permissions = nil
			v.Definitions = nil
			for j := 0; j < 4; j++ {
				id := uint32(1 + (n+j*41)%136)
				flags := []uint32{0, 16, 18, 20, 22}[((n/3)+j)%5]
				v.Permissions = append(v.Permissions, [2]uint32{id, []uint32{0, 0x08000000, 0x10000000, 0x20000000, 0x40000000, 0x80000000, 0xf8000000}[((n/2)+j)%7]})
				v.Definitions = append(v.Definitions, server.PortTestSpellClassDef{Index: id, Flags: flags, Valid: n%3 != 0})
			}
			for _, off := range []uint32{1452, 1460, 1468, 1476, 1484} {
				s.Main.Words = append(s.Main.Words, [2]uint32{off, m.Frame - uint32(n%3) + 1})
			}
			v.DurationSpell = [2]uint32{uint32(n % 137), uint32((n + 50) % 137)}
			v.DurationOwner = [2]bool{n%2 == 0, n%3 == 0}
			v.DurationEmpty = n%7 == 0
			v.SummonAllowed = n%3 != 0
			v.Missile = n%2 == 0
			v.MissileTarget = n%3 != 0
			v.TargetSubclass = uint32(n % 4)
			v.TargetFlags = []uint32{0, 4, 0x8004}[n%3]
			v.TargetCur = uint16(n % 100)
			v.TargetMax = uint16(100 + n%2)
			v.NilTargetHealth = n%7 == 0
			v.Second = n%5 == 0
			v.SecondPos = [2]uint32{b(float32(90 + n%80)), b(100)}
			v.Mode = n % 4
			v.Accuracy = []uint32{0, b(.5), b(1), b(2), b(-.5)}[n%5]
			s.Combat.AttackFrame = uint32(n % 4)
			s.Combat.Anim = byte(n % 5)
			s.Combat.Progress = byte(n % 2)
			s.Combat.Target = [2]uint32{b(110 + float32(n%40)), b(100 + float32(n%20))}
			s.Combat.Velocity = [2]uint32{b(float32(n%9 - 4)), b(float32(n%7 - 3))}
			v.SelfTarget = n%5 == 0
			if op == 12 || op == 13 {
				v.Spell = 41
				if n%4 == 0 {
					m.Status |= 0x20000
					s.Combat.PlayerUpdate = true
					m.Subclass = 0x10
				}
			}
			// Recoil requires a real target object, even when the action is blocked.
			if op == 13 && v.Mode == 0 {
				v.TargetArg = true
			}
			if v.Missile && v.SelfTarget {
				v.MissileTarget = true
			}
			all = append(all, s)
		}
	}
	return all
}
func aiSpellHash(t *testing.T, label string, got []legacy.PortTestRoamResult, want string) {
	t.Helper()
	for i, r := range got {
		if !r.Intact || !r.Spells.Intact || !r.MonsterState.Intact || !r.Combat.Intact {
			t.Fatalf("%s case %d guarded state: %s", label, i, mustMainJSON(r))
		}
	}
	data, e := json.Marshal(got)
	if e != nil {
		t.Fatal(e)
	}
	if p := os.Getenv("OPENNOX_SPELL_CAPTURE"); p != "" {
		if e := os.WriteFile(p+"-"+label+".json", data, 0600); e != nil {
			t.Fatal(e)
		}
	}
	sum := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s cases=%d complete-state=%s", label, len(got), sum)
	if want != "" && want != sum {
		t.Fatalf("want %s got %s", want, sum)
	}
}
func TestAISpellCorpus(t *testing.T) {
	specs := aiSpellCorpus()
	if s := os.Getenv("OPENNOX_SPELL_CASE"); s != "" {
		var i int
		if _, e := fmt.Sscanf(s, "%d", &i); e != nil {
			t.Fatal(e)
		}
		specs = specs[i : i+1]
	}
	aiSpellHash(t, "corpus", legacy.PortTestRoam(specs), "984b1523ba137d057e07c021202a8dc167c7bd067e27c6257a9f91c17993dc0f")
}
func aiSpellCastID(r legacy.PortTestRoamResult) uint32 {
	for _, a := range r.Combat.Actions {
		if a.Action >= 18 && a.Action <= 20 {
			return uint32(a.Args[0])
		}
	}
	return 0
}
func TestAISpellContracts(t *testing.T) {
	type contract struct {
		name  string
		spec  legacy.PortTestRoamSpec
		check func(*testing.T, legacy.PortTestRoamResult)
	}
	var all []contract
	add := func(name string, s legacy.PortTestRoamSpec, check func(*testing.T, legacy.PortTestRoamResult)) {
		all = append(all, contract{name, s, check})
	}
	for _, op := range []int{1, 3, 7, 8, 9} {
		for _, id := range []uint32{1, 41, 136} {
			for _, flags := range []uint32{0, 16, 18} {
				s := aiSpellBase(op)
				s.Spells.SelfTarget = true
				s.Spells.Permissions = [][2]uint32{{id, 0xf8000000}}
				s.Spells.Definitions = []server.PortTestSpellClassDef{{Index: id, Flags: flags, Valid: true}}
				if op == 1 {
					s.Spells.Missile = true
					s.Spells.MissileTarget = true
					s.Spells.TargetSubclass = 2
				}
				add(fmt.Sprintf("select-%d-%d-%d", op, id, flags), s, func(t *testing.T, r legacy.PortTestRoamResult) {
					want := uint32(0)
					if flags&16 != 0 {
						want = id
					}
					if aiSpellCastID(r) != want || r.Spells.Return != uint32(boolToInt(want != 0)) {
						t.Fatalf("selected %d return %d want %d", aiSpellCastID(r), r.Spells.Return, want)
					}
				})
			}
		}
	}
	for _, op := range []int{3, 9} {
		s := aiSpellBase(op)
		sp := uint32(server.EnchantID(3).Spell())
		s.Owner.Buffs = 1 << 3
		s.Spells.Permissions = [][2]uint32{{41, 0xf8000000}, {sp, 0xf8000000}}
		s.Spells.Definitions = append(s.Spells.Definitions, server.PortTestSpellClassDef{Index: sp, Flags: 16, Valid: true})
		add(fmt.Sprint("any-existing-enchant-", op), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Spells.Return != 0 || r.Logic != 1 || r.Changed {
				t.Fatal("candidate-wide enchant rejection")
			}
		})
	}
	for _, allowed := range []bool{false, true} {
		for _, active := range []bool{false, true} {
			for _, mixed := range []bool{false, true} {
				s := aiSpellBase(7)
				s.Spells.SelfTarget = true
				s.Spells.SummonAllowed = allowed
				s.Spells.Permissions = [][2]uint32{{75, 0x40000000}}
				s.Spells.Definitions = []server.PortTestSpellClassDef{{Index: 75, Flags: 16, Valid: true}, {Index: 41, Flags: 16, Valid: true}}
				if active {
					s.Spells.DurationSpell[0] = 75
					s.Spells.DurationOwner[0] = true
				}
				if mixed {
					s.Spells.Permissions = append(s.Spells.Permissions, [2]uint32{41, 0x40000000})
				}
				add(fmt.Sprintf("summon-%t-%t-%t", allowed, active, mixed), s, func(t *testing.T, r legacy.PortTestRoamResult) {
					got := aiSpellCastID(r)
					if !mixed && (!allowed || active) && got != 0 {
						t.Fatal("summon must reject")
					}
					if mixed && (!allowed || active) && got != 41 {
						t.Fatal("must retry to non-summon")
					}
					if !mixed && allowed && !active && got != 75 {
						t.Fatal("summon allowed")
					}
				})
			}
		}
	}
	for _, max := range []uint16{100, 101} {
		for _, cur := range []uint16{49, 50, 51} {
			s := aiSpellBase(10)
			s.MonsterState.Status = 0x820
			s.MonsterState.HealthCur = cur
			s.MonsterState.HealthMax = max
			add(fmt.Sprintf("self-heal-%d-%d", cur, max), s, func(t *testing.T, r legacy.PortTestRoamResult) {
				want := uint32(0)
				if cur < max>>1 {
					want = 41
				}
				if aiSpellCastID(r) != want || r.Spells.Return != 0 {
					t.Fatal("self heal integer-half/return")
				}
			})
		}
	}
	for _, other := range []bool{false, true} {
		s := aiSpellBase(10)
		s.MonsterState.Status = 0x1020
		s.Spells.Second = other
		s.Spells.SecondPos = s.Combat.Target
		add(fmt.Sprint("heal-last-target-", other), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Spells.Return != 1 || r.Spells.Globals[1] != 100 || aiSpellCastID(r) != 41 {
				t.Fatalf("heal result %s", mustMainJSON(r.Spells))
			}
		})
	}

	for _, op := range []int{1, 3, 7, 8, 9} {
		for _, frame := range []uint32{0, 128, 0x80000000, 0xffffffff} {
			for _, delta := range []int32{-1, 0, 1} {
				s := aiSpellBase(op)
				s.MonsterState.Frame = frame
				s.Spells.SelfTarget = true
				if op == 1 {
					s.Spells.Missile = true
					s.Spells.MissileTarget = true
					s.Spells.TargetSubclass = 2
				}
				off := map[int]uint32{1: 1452, 3: 1460, 7: 1476, 8: 1468, 9: 1484}[op]
				deadline := frame + uint32(delta)
				s.Main.Words = append(s.Main.Words, [2]uint32{off, deadline})
				add(fmt.Sprintf("cooldown-%d-%x-%d", op, frame, delta), s, func(t *testing.T, r legacy.PortTestRoamResult) {
					eligible := frame >= deadline
					want := uint32(0)
					if eligible {
						want = 1
					}
					if r.Spells.Return != want {
						t.Fatal("cooldown unsigned conversion")
					}
					if eligible && lifeWord(r, 4096+off, deadline) == deadline {
						t.Fatal("cooldown was not rescheduled")
					}
					if !eligible && r.Logic != 1 {
						t.Fatal("blocked cooldown consumed RNG")
					}
				})
			}
		}
	}
	for _, mode := range []int{0, 1, 2} {
		for _, frame := range []byte{0, 1, 2} {
			for _, progress := range []byte{0, 1} {
				for _, muted := range []bool{false, true} {
					s := aiSpellBase(13)
					s.Spells.Mode = mode
					s.Combat.AttackFrame = 2
					s.Combat.Anim = frame
					s.Combat.Progress = progress
					if muted {
						s.Owner.Buffs = 1 << 29
					}
					add(fmt.Sprintf("cast-frame-%d-%d-%d-%t", mode, frame, progress, muted), s, func(t *testing.T, r legacy.PortTestRoamResult) {
						cast := false
						for i, v := range r.Trace {
							if v == 40 && i+2 < len(r.Trace) && r.Trace[i+1] == 41 && r.Trace[i+2] == 101 {
								cast = true
							}
						}
						if cast != (progress == 0 && frame == 2 && mode < 2 && !muted) {
							t.Fatal("cast frame/mode/suppression")
						}
						wantSound := progress == 0 && frame == 1
						if (len(r.Combat.Sounds) > 0) != wantSound {
							t.Fatal("cast audio frame")
						}
						draws := 1
						if progress == 0 && frame == 2 && mode == 0 {
							draws = 4
						}
						if r.Logic != draws {
							t.Fatal("recoil RNG before suppression", r.Logic, draws)
						}
					})
				}
			}
		}
	}
	for _, morph := range []bool{false, true} {
		s := aiSpellBase(12)
		if morph {
			s.MonsterState.Status |= 0x20000
			s.Combat.PlayerUpdate = true
		}
		add(fmt.Sprint("direct-cast-morph-", morph), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := uint32(2)
			if morph {
				want = 4
			}
			found := false
			for i, v := range r.Trace {
				if v == 46 && i+3 < len(r.Trace) {
					found = true
					if r.Trace[i+1] != want || r.Trace[i+2] != 0 {
						t.Fatal("cast callback class/subclass")
					}
				}
			}
			if !found {
				t.Fatal("cast callback missing")
			}
		})
	}
	for _, quest := range []bool{false, true} {
		s := aiSpellBase(10)
		s.MonsterState.Status = 0x1020
		s.Combat.Target = [2]uint32{math.Float32bits(400), math.Float32bits(100)}
		if quest {
			s.Lifecycle.GameFlags = 4096
		}
		add(fmt.Sprint("heal-rectangle-", quest), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := uint32(0)
			if quest {
				want = 1
			}
			if r.Spells.Return != want {
				t.Fatal("heal rectangle range")
			}
		})
	}
	for _, kind := range []int{0, 1, 2, 3} {
		s := aiSpellBase(10)
		s.MonsterState.Status = 0x1020
		switch kind {
		case 0:
			s.Spells.NilTargetHealth = true
		case 1:
			s.Spells.TargetFlags |= 0x8000
		case 2:
			s.Spells.TargetCur = 50
			s.Spells.TargetMax = 101
		case 3:
			s.Spells.TargetClass = 4
			s.Main.TargetPlayer = true
		}
		add(fmt.Sprint("heal-reject-", kind), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Spells.Return != 0 || aiSpellCastID(r) != 0 {
				t.Fatal("heal rejected target")
			}
		})
	}
	recoilWant := [][2]uint32{{1126244515, 1119999876}, {1121905324, 1122758101}, {1124630886, 1119804006}, {1118958158, 1122549300}, {3442263653, 1309958949}}
	for i, accuracy := range []uint32{0, math.Float32bits(.5), math.Float32bits(1), 0x3e800001, math.Float32bits(-16777216)} {
		s := aiSpellBase(14)
		s.Spells.Accuracy = accuracy
		s.Seed = 3 + i
		s.Combat.Velocity = [2]uint32{math.Float32bits(2.3), math.Float32bits(-1.7)}
		add(fmt.Sprint("recoil-", i), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Spells.Output != recoilWant[i] {
				t.Fatalf("recoil bits: got %v want %v", r.Spells.Output, recoilWant[i])
			}
			if r.Logic != s.Seed+3 {
				t.Fatal("three recoil draws")
			}
		})
	}
	specs := make([]legacy.PortTestRoamSpec, len(all))
	for i, v := range all {
		specs[i] = v.spec
	}
	got := legacy.PortTestRoam(specs)
	for i, c := range all {
		t.Run(c.name, func(t *testing.T) { c.check(t, got[i]) })
	}
	aiSpellHash(t, "contracts", got, "3a1f3cbadc33dd52919d8c064995513c5e4526e0710a8cbb1b70483cef5b51ee")
}
func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
