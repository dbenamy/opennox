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
)

func stateBase(op int) legacy.PortTestRoamSpec {
	b := math.Float32bits
	return legacy.PortTestRoamSpec{Op: 12, Seed: 1, Stack: 1,
		Owner:        &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, X: b(100), Y: b(100)},
		Combat:       &legacy.PortTestCombatSpec{Sound: true, Target: [2]uint32{b(150), b(100)}},
		Lifecycle:    &legacy.PortTestLifecycleSpec{Cur: 3, Max: 100},
		MonsterState: &legacy.PortTestMonsterStateSpec{Op: op, Stack: 1, Action: 1, Frame: 123, FPS: 30, Speed: b(1), Pos: [2]uint32{b(100), b(100)}, Arg: [2]uint32{b(110), b(100)}, HealthCur: 3, HealthMax: 100, MeleeRange: b(10), FrameA: 2, FrameB: 4}}
}
func stateCorpus() []legacy.PortTestRoamSpec {
	b := math.Float32bits
	var out []legacy.PortTestRoamSpec
	for n := 0; n < 512; n++ {
		for op := 0; op < 28; op++ {
			s := stateBase(op)
			v := s.MonsterState
			s.Seed = n + 1
			v.Type = n % 5
			v.Stack = int8(n % 24)
			v.Action = uint32(n % 73)
			if op == 0 && n%7 == 0 {
				v.Stack = -1
			}
			v.Status = []uint32{0, 0x20, 0x40, 0x80, 0x200, 0x4000, 0x8000, 0x10000, 0x40000, 0x80000, 0xffffffff}[n%11]
			v.Subclass = []uint32{0, 0x10, 0x20, 0x30}[n%4]
			v.ObjFlags = []uint32{0, 4, 0x8000, 0x8004}[n%4]
			v.PlayerWeapon = uint32(1) << uint(n%32)
			v.PlayerShield = uint32(1) << uint((n/2)%32)
			v.MissileName = byte(n % 2)
			v.Poison = byte(n % 3)
			v.NilHealth = n%4 == 0
			v.Speed = []uint32{0, b(.01), b(1), b(-1), 0x7fc12345}[n%5]
			v.MeleeRange = []uint32{0, b(1), b(-1), 0x7fc12345, 0x7f800000}[n%5]
			v.Aggression = []uint32{0, b(.08), b(.33), b(.66), b(1), 0x7fc12345}[n%6]
			v.Frame = []uint32{0, 1, 123, 0xffffffff}[n%4]
			v.FPS = []uint32{1, 30, 60}[n%3]
			v.Deadline = v.Frame - []uint32{0, 2, 3, 89, 90, 91, 0xffffffff}[n%7]
			v.MimicAge = v.Frame - []uint32{0, 1, 29, 30, 31, 0xffffffff}[n%6]
			v.Direction = int16(n % 256)
			v.Arg = [2]uint32{b(float32(100 + n%17)), b(float32(100 + n%13))}
			v.NetCode = uint32(n)
			v.FrameA = uint32(n % 7)
			v.FrameB = uint32(n % 9)
			s.Combat.Anim = byte(n % 9)
			s.Combat.Progress = byte(n % 2)
			v.AnimData = n%2 == 0
			v.FallbackAnim = uint32(n % 64)
			v.Source = n % 3
			v.Order = n % 8
			v.Broadcast = n%2 == 0
			v.Own = n%3 != 0
			v.Second = n%3 == 0
			v.SecondEnabled = n%5 != 0
			v.NilUnit = n%17 == 0
			s.Owner.Enemy = n%2 == 0
			if n%5 == 0 {
				v.MimicCache = 1
				v.PlantCache = 3
				v.ZombieCache = 7
				v.VileZombieCache = 9
			}
			out = append(out, s)
		}
	}
	return out
}
func stateHash(t *testing.T, label string, got []legacy.PortTestRoamResult, want string) {
	t.Helper()
	for i, r := range got {
		if !r.Intact {
			t.Fatalf("%s case %d guarded state changed: %+v", label, i, r)
		}
	}
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	h := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s cases=%d complete-state-sha256=%s", label, len(got), h)
	if p := os.Getenv("OPENNOX_STATE_CAPTURE"); p != "" {
		if err := os.WriteFile(p+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if want != "" && os.Getenv("OPENNOX_STATE_CASE") == "" && h != want {
		t.Fatal("original C state differs", h)
	}
}
func TestAIMonsterStateOriginal(t *testing.T) {
	specs := stateCorpus()
	if v := os.Getenv("OPENNOX_STATE_CASE"); v != "" {
		var n int
		if _, err := fmt.Sscan(v, &n); err != nil {
			t.Fatal(err)
		}
		specs = specs[n : n+1]
	}
	stateHash(t, "corpus", legacy.PortTestRoam(specs), "cce74b363d6c5c4d2210781f9b252f503ee34763880ddf8e14037d7551bab1f5")
}

func TestAIMonsterStateContracts(t *testing.T) {
	b := math.Float32bits
	var specs []legacy.PortTestRoamSpec
	var checks []func(*testing.T, legacy.PortTestRoamResult)
	add := func(s legacy.PortTestRoamSpec, f func(*testing.T, legacy.PortTestRoamResult)) {
		specs = append(specs, s)
		checks = append(checks, f)
	}
	result := func(want uint64) func(*testing.T, legacy.PortTestRoamResult) {
		return func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.MonsterState.Return != want {
				t.Fatal("return", r.MonsterState.Return, want)
			}
		}
	}
	for _, v := range []struct {
		action, status uint32
		kind           int
		enemy          bool
		want           uint64
	}{
		{1, 0, 0, false, 8}, {7, 0, 0, false, 12}, {7, 0x4000, 0, false, 13}, {16, 0, 0, false, 1}, {17, 0, 0, false, 3}, {30, 0, 3, false, 9}, {30, 0x80000, 3, false, 15}, {31, 0x80000, 3, false, 10}, {1, 0x40000, 1, false, 0}, {1, 0, 1, false, 8}, {1, 0, 2, false, 14}, {1, 0, 2, true, 8}, {34, 0, 0, false, 15}, {35, 0, 0, false, 14},
	} {
		s := stateBase(0)
		s.MonsterState.Action = v.action
		s.MonsterState.Status = v.status
		s.MonsterState.Type = v.kind
		s.Owner.Enemy = v.enemy
		add(s, result(v.want))
	}
	for _, v := range []struct {
		raw  uint32
		want uint64
	}{{b(.01), 0}, {b(math.Nextafter32(.01, 1)), 1}, {0x7fc12345, 0}, {0x7f800000, 1}, {b(-1), 0}} {
		s := stateBase(5)
		s.MonsterState.Speed = v.raw
		add(s, result(v.want))
	}
	for _, v := range []struct {
		raw  uint32
		want uint64
	}{{0, 0}, {b(-1), 0}, {b(1), 1}, {0x7fc12345, 1}, {0x7f800000, 1}} {
		s := stateBase(1)
		s.MonsterState.MeleeRange = v.raw
		add(s, result(v.want))
	}
	for _, weapon := range []uint32{0, 1, 2, 4, 0x100, 0x10000, 0x4000000, 0x8000000} {
		for _, op := range []int{1, 2} {
			s := stateBase(op)
			s.MonsterState.Subclass = 0x10
			s.MonsterState.PlayerWeapon = weapon
			want := uint64(0)
			if weapon&0x47f00fe != 0 {
				want = 1
			}
			if op == 1 {
				want = 1 - want
			}
			add(s, result(want))
		}
	}
	for _, op := range []int{13, 14} {
		for _, status := range []uint32{0, 0x4000, 0x8000, 0xc000, 0x10000, 0x14000, 0xffffffff} {
			s := stateBase(op)
			s.MonsterState.Status = status
			add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
				want := status
				if op == 13 && status&0x10000 == 0 {
					want |= 0x4000
				}
				if op == 14 && status&0x8000 == 0 {
					want &^= 0x4000
				}
				if r.MonsterState.Return != uint64(want) || r.Combat.Status != want {
					t.Fatal("running flags/return", r.MonsterState.Return, r.Combat.Status, want)
				}
			})
		}
	}
	for _, fps := range []uint32{1, 30, 60, 0x80000000} {
		for _, frame := range []uint32{0, 123, 0xffffffff} {
			for _, delta := range []uint32{0, 3*fps - 1, 3 * fps, 3*fps + 1} {
				s := stateBase(18)
				s.MonsterState.FPS = fps
				s.MonsterState.Frame = frame
				s.MonsterState.Deadline = frame - delta
				want := uint64(0)
				if delta < 3*fps {
					want = 1
				}
				add(s, result(want))
			}
		}
	}
	for _, order := range []int{2, 3, 4, 5} {
		for _, source := range []int{1, 2} {
			for _, broadcast := range []bool{false, true} {
				for _, enabled := range []bool{false, true} {
					s := stateBase(23)
					v := s.MonsterState
					v.Order = order
					v.Source = source
					v.Broadcast = broadcast
					v.Own = true
					v.Second = true
					v.SecondEnabled = true
					v.Status = 0x40
					if enabled {
						v.Status |= 0x80
					}
					v.Direction = -1
					add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
						active := !broadcast || enabled
						if active {
							actions := map[int]uint32{2: 0, 3: 4, 4: 3, 5: 5}
							a := r.Combat.Actions[len(r.Combat.Actions)-1]
							if a.Action != actions[order] {
								t.Fatal("command action", a, order)
							}
							aggression := b(.5)
							if order >= 4 {
								aggression = b(.83)
							}
							if lifeWord(r, 4096+1304, 0) != aggression {
								t.Fatal("command aggression")
							}
							if order == 3 && (a.Args[0] != uintptr(b(100)) || a.Args[1] != uintptr(b(100)) || uint32(a.Args[2]) != 0xffffffff) {
								t.Fatal("guard arguments", a)
							}
							if order == 4 {
								owner := uint32(102)
								pos := uint32(0)
								if source == 2 {
									owner = 200
									pos = b(100)
								}
								if uint32(a.Args[0]) != pos || uint32(a.Args[1]) != pos || uint32(a.Args[2]) != owner {
									t.Fatal("escort arguments", a)
								}
							}
						} else if r.Stack != v.Stack || r.Combat.Actions[len(r.Combat.Actions)-1].Action != v.Action {
							t.Fatal("broadcast ordered disabled unit")
						}
						cmd := uint32(0)
						if source == 2 && broadcast && order >= 3 {
							cmd = uint32(order)
						}
						if r.MonsterState.Command != cmd {
							t.Fatal("player broadcast command", r.MonsterState.Command, cmd)
						}
						sounds := 0
						if active {
							sounds++
						}
						if broadcast {
							sounds++
						}
						if len(r.Combat.Sounds) != sounds {
							t.Fatal("command sounds", r.Combat.Sounds, sounds)
						}
						for _, sound := range r.Combat.Sounds {
							if sound != 317 {
								t.Fatal("command sound id", sound)
							}
						}
					})
				}
			}
		}
	}
	for _, op := range []int{0, 19, 21, 22} {
		s := stateBase(op)
		s.MonsterState.Type = 0
		s.MonsterState.MimicCache = 1
		s.MonsterState.PlantCache = 1
		s.MonsterState.ZombieCache = 1
		s.MonsterState.VileZombieCache = 42
		want := uint64(1)
		if op == 0 {
			s.MonsterState.Status = 0x40000
			want = 0
		}
		add(s, result(want))
	}
	for _, action := range []uint32{1, 16, 18, 21, 23, 30, 31} {
		s := stateBase(24)
		s.MonsterState.Subclass = 0x10
		s.MonsterState.Action = action
		s.MonsterState.FallbackAnim = 6
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			idx := map[uint32]int{16: 6, 18: 21, 21: 40, 23: 47, 30: 1, 31: 2}
			count, delay := 0, 0
			if i, ok := idx[action]; ok {
				count = 3 + i%7
				delay = i % 3
			}
			want := uint32(count<<8 | delay<<16)
			if r.MonsterState.Return != 800 || r.MonsterState.Animation[0] != 0 || r.MonsterState.Animation[1] != 0 || r.MonsterState.Animation[2] != want || r.MonsterState.Animation[3] != 0 {
				t.Fatal("NPC player range", r.MonsterState)
			}
		})
	}
	for _, anim := range []bool{false, true} {
		s := stateBase(24)
		s.MonsterState.AnimData = anim
		s.MonsterState.Action = 16
		want := uint64(0)
		if anim {
			want = 701
		}
		add(s, result(want))
	}
	for _, point := range []uint32{b(.9), b(math.Nextafter32(.9, 1)), b(1), b(-1)} {
		s := stateBase(26)
		s.MonsterState.Direction = 0
		s.MonsterState.Arg = [2]uint32{point, 0}
		want := uint64(0)
		if float64(math.Float32frombits(point)) > .89999998 {
			want = 1
		}
		add(s, result(want))
	}
	for _, d := range []struct {
		x, y float32
		want uint32
	}{{1, 0, 0}, {0, 1, 64}, {-1, 0, 128}, {0, -1, 192}} {
		s := stateBase(27)
		s.MonsterState.Direction = 17
		s.MonsterState.Arg = [2]uint32{b(100 + d.x), b(100 + d.y)}
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if lifeWord(r, 124, 17|17<<16) != 17|d.want<<16 {
				t.Fatal("facing direction", r.Changes)
			}
		})
	}
	for _, seeds := range [][2]uint32{{0, 42}, {1, 0}, {7, 0}, {0, 3}} {
		s := stateBase(22)
		s.MonsterState.Type = 4
		s.MonsterState.ZombieCache = seeds[0]
		s.MonsterState.VileZombieCache = seeds[1]
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			z, v := seeds[0], seeds[1]
			if z == 0 {
				z, v = 2, 3
			}
			if r.MonsterState.Caches[2] != z || r.MonsterState.Caches[3] != v {
				t.Fatal("joint zombie cache", r.MonsterState.Caches)
			}
			want := uint64(0)
			if z == 3 || v == 3 {
				want = 1
			}
			if r.MonsterState.Return != want {
				t.Fatal("zombie cache reuse", r.MonsterState.Return, want)
			}
		})
	}
	for _, fallback := range []uint32{0, 6, 0x100, 0x106, 255} {
		s := stateBase(24)
		s.MonsterState.Subclass = 0x10
		s.MonsterState.Action = 16
		s.MonsterState.FallbackAnim = fallback
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			idx := int(byte(fallback))
			count, delay := 0, 0
			if idx < 64 {
				count, delay = 3+idx%7, idx%3
			}
			if r.MonsterState.Animation[2] != uint32(count<<8|delay<<16) {
				t.Fatal("fallback animation byte/zero", fallback, r.MonsterState.Animation)
			}
		})
	}
	for _, source := range []int{0, 1, 2} {
		for _, own := range []bool{false, true} {
			for _, kind := range []int{0, 3} {
				for _, dead := range []bool{false, true} {
					for _, order := range []int{0, 1} {
						s := stateBase(23)
						v := s.MonsterState
						v.Source = source
						v.Own = own
						v.Type = kind
						v.Order = order
						if dead {
							v.ObjFlags = 0x8000
						}
						add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
							allowed := source != 0 && (!dead || kind == 3 && order == 0)
							banish := allowed && order == 0 && own
							observe := allowed && order == 1 && source == 2
							gotB := len(r.Trace) > 0 && r.Trace[0] == 10
							gotO := len(r.Trace) > 0 && r.Trace[0] == 34
							if gotB != banish || gotO != observe {
								t.Fatal("banish/observe gates", gotB, gotO, banish, observe, r.Trace)
							}
						})
					}
				}
			}
		}
	}

	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		t.Run(fmt.Sprint(i), func(t *testing.T) { checks[i](t, r) })
	}
	stateHash(t, "contracts", got, "44bab745a7d8d2e58461eafc5821f61a083f1367898d3dd03e0a5d006a978c03")
}
