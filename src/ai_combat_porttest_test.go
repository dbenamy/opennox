//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"os"
	"testing"
)

func combatCorpus() []legacy.PortTestRoamSpec {
	bits := math.Float32bits
	var out []legacy.PortTestRoamSpec
	for n := 0; n < 1024; n++ {
		for op := 0; op < 11; op++ {
			o := &legacy.PortTestRoamOwnerSpec{Frame: []uint32{0, 1, 123, 0x7fffffff, 0x80000000, 0xfffffff0, 0xffffffff}[n%7], FPS: []uint32{1, 30, 60}[n%3], X: bits(100), Y: bits(100), Buffs: 1 << 29, Status: []uint32{0, 0x20, 0x40, 0x100, 0x4000, 0x8000, 0x10000}[n%7], Enemy: n%3 == 0}
			c := &legacy.PortTestCombatSpec{Op: op, Phase: 0, Strike: n%3 - 1, Cooldown: o.Frame + []uint32{0, 1, 0xffffffff}[n%3], ArgFrame: o.Frame - []uint32{0, 1, 299, 300, 301, 600, 601, 0xffffffff}[n%8], MeleeRange: []uint32{0, bits(10), bits(50), 0x7fc12345}[n%4], MissileRange: bits(100), AttackFrame: uint32(n % 4), DelayMin: 2, DelayMax: 7, Target: [2]uint32{bits(float32(100 + n%60)), bits(100)}, Direction: uint16(n % 256), Anim: byte(n % 3), Progress: byte(n / 3 % 2), Done: byte(n / 6 % 2), Sound: n%2 == 0, Shoot: n%3 == 0, Killable: n%2 == 0, TargetFlags: []uint32{0, 4, 0x8004}[n%3], TargetClass: 2, ScanCode: 123, EnemyCode: 123, HeardCode: 123, Heard: 777}
			// Player weapon engines are retained; their owner paths get dedicated valid fixtures.
			if op == 0 {
				c.Phase = n % 3
			}
			if op == 4 || op == 5 {
				c.Phase = n % 2
			}
			if op == 1 || op == 2 || op == 3 {
				c.Playerlike = n%2 == 0
			}
			if op == 6 || op == 7 || op == 8 || op == 9 {
				c.Playerlike = n%5 == 0
				c.Weapon = []uint32{0, 2, 0x100}[n%3]
			}
			if op == 5 {
				c.TargetArg = n%2 == 0
				c.Wall = n % 3
			}
			if op == 0 && c.Phase == 0 && !o.Enemy {
				c.Scan = n%2 == 0
			}
			out = append(out, legacy.PortTestRoamSpec{Op: 10, Seed: n + 1, Stack: int8(n % 24), Owner: o, Combat: c})
		}
	}
	for n := 0; n < 4096; n++ {
		o := &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, X: bits(100), Y: bits(100), Buffs: 1 << 29}
		c := &legacy.PortTestCombatSpec{Op: 10, Strike: 0, Direction: uint16(n % 256), TargetClass: 2, Target: [2]uint32{bits(100 + float32(n%137-68)/3), bits(100 + float32(n%97-48)/5)}}
		out = append(out, legacy.PortTestRoamSpec{Op: 10, Seed: 1, Stack: 1, Owner: o, Combat: c})
	}
	for n := 0; n < 2048; n++ {
		o := &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, X: bits(100), Y: bits(100), Buffs: 1 << 29}
		c := &legacy.PortTestCombatSpec{Op: 5, Strike: 0, Shoot: true, Sound: true, Anim: 2, AttackFrame: 2, Direction: uint16(n % 256), Radius: bits(float32(n%37) / 19), Target: [2]uint32{bits(float32(n%131) / 17), bits(float32(n%137) / 13)}}
		out = append(out, legacy.PortTestRoamSpec{Op: 10, Seed: 1, Stack: 1, Owner: o, Combat: c})
	}

	return out
}
func TestAICombatOriginalState(t *testing.T) {
	specs := combatCorpus()
	if v := os.Getenv("OPENNOX_COMBAT_CASE"); v != "" {
		var n int
		if _, err := fmt.Sscan(v, &n); err != nil {
			t.Fatal(err)
		}
		specs = specs[n : n+1]
	}
	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		if !r.Intact {
			t.Fatalf("case %d changed guarded/read-only state: %+v", i, r)
		}
	}
	data, _ := json.Marshal(got)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("cases=%d complete-state-sha256=%s", len(got), hash)
	if os.Getenv("OPENNOX_COMBAT_CAPTURE") != "" {
		if err := os.WriteFile("../build/port-ai-combat/c-cases.json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	const baseline = "68c25cbb61de2e5f1603f76f7349e9595ebaae813a266f61b7f8099ead1221fe"
	if baseline != "" && os.Getenv("OPENNOX_COMBAT_CASE") == "" && hash != baseline {
		t.Fatal("original C state differs", hash)
	}
}

func TestAICombatContracts(t *testing.T) {
	bits := math.Float32bits
	base := func(op, phase int) legacy.PortTestRoamSpec {
		return legacy.PortTestRoamSpec{Op: 10, Seed: 1, Stack: 1, Owner: &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, X: bits(100), Y: bits(100), Buffs: 1 << 29}, Combat: &legacy.PortTestCombatSpec{Op: op, Phase: phase, Strike: 1, Sound: true, Cooldown: 123, ArgFrame: 123, MeleeRange: bits(10), MissileRange: bits(100), DelayMin: 5, DelayMax: 5, AttackFrame: 2, Anim: 2, Target: [2]uint32{bits(200), bits(100)}, TargetClass: 2}}
	}
	var specs []legacy.PortTestRoamSpec
	var checks []func(*testing.T, legacy.PortTestRoamResult)
	add := func(s legacy.PortTestRoamSpec, f func(*testing.T, legacy.PortTestRoamResult)) {
		specs = append(specs, s)
		checks = append(checks, f)
	}
	sound := func(t *testing.T, r legacy.PortTestRoamResult, want uint32) {
		t.Helper()
		if len(r.Combat.Sounds) != 1 || r.Combat.Sounds[0] != want {
			t.Fatalf("sound: %v want %d", r.Combat.Sounds, want)
		}
	}
	for _, op := range []int{4, 5} {
		for _, delta := range []uint32{0, 1, 0xffffffff} {
			s := base(op, 1)
			s.Combat.Cooldown += delta
			add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
				if delta == 1 {
					if r.Stack != 3 || r.Combat.Actions[2].Action != 50 || r.Combat.Actions[3].Action != 1 || r.Combat.Actions[3].Args[0] != 124 {
						t.Fatalf("cooldown wait chain: %+v", r.Combat.Actions)
					}
					if len(r.Combat.Sounds) != 0 {
						t.Fatal("waiting played sound")
					}
				} else {
					if r.Combat.Cooldown != 128 || r.Logic != 2 {
						t.Fatalf("cooldown/RNG: %+v", r)
					}
					want := uint32(306)
					if op == 5 {
						want = 310
					}
					sound(t, r, want)
				}
			})
		}
	}
	for _, hit := range []int{-1, 0, 1} {
		for _, anim := range []byte{1, 2, 3} {
			for _, progress := range []byte{0, 1} {
				for _, done := range []byte{0, 1} {
					s := base(4, 0)
					s.Combat.Strike = hit
					s.Combat.Anim = anim
					s.Combat.Progress = progress
					s.Combat.Done = done
					add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
						strikes := 0
						if hit >= 0 && anim == 2 && progress == 0 {
							strikes = 1
						}
						if r.Combat.Strikes != strikes {
							t.Fatal("strike count", r.Combat.Strikes, strikes)
						}
						pop := hit < 0 || done != 0
						want := int8(1)
						if pop {
							want = 0
						}
						if r.Stack != want {
							t.Fatal("completion", r.Stack, want)
						}
						if strikes == 1 {
							id := uint32(308)
							if hit == 0 {
								id = 309
							}
							sound(t, r, id)
						} else if len(r.Combat.Sounds) != 0 {
							t.Fatal("unexpected strike sound")
						}
					})
				}
			}
		}
	}
	for _, shoot := range []bool{false, true} {
		for _, wall := range []int{0, 1, 2} {
			for _, targetArg := range []bool{false, true} {
				for _, dir := range []uint16{0, 32, 64, 128, 192, 255} {
					s := base(5, 0)
					s.Combat.Shoot = shoot
					s.Combat.Wall = wall
					if wall != 0 {
						s.Combat.Radius = bits(100)
					}
					s.Combat.TargetArg = targetArg
					s.Combat.Direction = dir
					s.Combat.Velocity = [2]uint32{bits(1), bits(2)}
					add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
						sound(t, r, 311)
						if !shoot {
							if len(r.Combat.Projectile) != 0 {
								t.Fatal("missing type allocated")
							}
							return
						}
						if len(r.Combat.Projectile) != 193 {
							t.Fatal("projectile not captured")
						}
						p := r.Combat.Projectile
						if wall == 1 && dir == 0 {
							if p[20] != 0 || p[21] != 0 {
								t.Fatal("blocked projectile moved")
							}
						}
						if wall == 0 {
							if p[31] != uint32(dir)|uint32(dir)<<16 {
								t.Fatal("projectile direction", p[31], dir)
							}
							if p[20] == 0 && p[21] == 0 {
								t.Fatal("projectile has no velocity")
							}
						}
					})
				}
			}
		}
	}
	for _, op := range []int{1, 2, 3} {
		for _, player := range []bool{false, true} {
			for _, done := range []byte{0, 1} {
				for _, deadline := range []uint32{122, 123, 124} {
					s := base(op, 0)
					s.Combat.Playerlike = player
					s.Combat.Done = done
					s.Combat.ArgFrame = deadline
					add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
						want := int8(1)
						if op == 1 {
							if deadline < 123 && player {
								want = 0
							}
							if deadline < 123 && !player && r.Combat.Actions[len(r.Combat.Actions)-1].Action != 22 {
								t.Fatal("missing block finish")
							}
						} else if done != 0 {
							want = 0
						}
						if r.Stack != want {
							t.Fatal("block completion", r.Stack, want)
						}
					})
				}
			}
		}
	}
	s := base(0, 1)
	add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
		sound(t, r, 305)
		if r.Combat.Status&0x4100 != 0x4100 {
			t.Fatal("fight flags")
		}
		if len(r.Trace) < 5 || r.Trace[0] != 10 || r.Trace[1] != 13 || r.Trace[2] != 1240 || r.Trace[4] != 1 {
			t.Fatal("fight script", r.Trace)
		}
	})
	for _, morph := range []bool{false, true} {
		for _, op := range []int{4, 5} {
			s := base(op, 0)
			s.Combat.Playerlike = true
			s.Combat.PlayerUpdate = true
			if morph {
				s.Owner.Status = 0x20000
			}
			add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
				if r.Stack != 0 || r.Combat.Strikes != 0 || len(r.Combat.Projectile) != 0 {
					t.Fatal("player attack false path", r)
				}
			})
		}
	}
	for _, stamina := range []byte{0, 9, 10, 11, 69, 70, 71, 255} {
		s := base(4, 1)
		s.Combat.Playerlike = true
		s.Combat.Stamina = stamina
		s.Combat.Weapon = 0x200
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := byte(0)
			if stamina < 70 {
				want = stamina - 70
			}
			if r.Combat.Stamina != want {
				t.Fatal("stamina compatibility", r.Combat.Stamina, want)
			}
		})
	}
	for _, friend := range []bool{false, true} {
		s := base(4, 1)
		s.Combat.Scan = true
		s.Combat.TargetFlags = 4
		s.Combat.Target[0] = bits(110)
		s.Combat.Friendly = true
		s.Combat.Playerlike = friend
		s.Combat.WeaponSubclass = 0x4000
		s.Combat.Direction = 0
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Combat.Selected != 100 {
				t.Fatal("front unit not selected")
			}
			if friend {
				a := r.Combat.Actions
				if len(a) != 4 || a[1].Action != 27 || a[2].Action != 41 || a[3].Action != 24 || a[3].Args[0] != uintptr(bits(110)) {
					t.Fatal("friend avoidance chain", a)
				}
			} else {
				sound(t, r, 306)
			}
		})
	}
	for _, dead := range []bool{false, true} {
		for _, matching := range []bool{false, true} {
			s := base(0, 0)
			s.Combat.Scan = true
			s.Combat.TargetFlags = 4
			s.Combat.ScanCode = 12
			s.Combat.EnemyCode = 12
			s.Combat.HeardCode = 12
			s.Combat.Heard = 987
			if dead {
				s.Combat.TargetFlags |= 0x8000
			}
			if !matching {
				s.Combat.ScanCode++
			}
			add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
				want := int8(3)
				if dead && matching {
					want = 0
				}
				if r.Stack != want {
					t.Fatal("fight scan", r.Stack, want)
				}
			})
		}
	}
	s = base(1, 0)
	s.Combat.Scan = true
	s.Combat.Shield = true
	s.Combat.TargetClass = 1
	s.Combat.TargetFlags = 4
	s.Combat.Target[0] = bits(110)
	s.Combat.Velocity = [2]uint32{bits(-3), 0}
	s.Combat.ArgFrame = 0
	add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
		if r.Combat.Actions[len(r.Combat.Actions)-1].Args[0] != 138 {
			t.Fatal("shield deadline not extended", r.Combat.Actions, r.Trace)
		}
	})
	results := legacy.PortTestRoam(specs)
	for i, r := range results {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			if !r.Intact {
				t.Fatal("guard/read-only mutation")
			}
			checks[i](t, r)
		})
	}
	data, _ := json.Marshal(results)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("contracts=%d complete-state-sha256=%s", len(results), hash)
	const baseline = "3375bac303a78fc6473ce0dfe489428def83a47b48146ac3d34224bc55add2b3"
	if baseline != "" && hash != baseline {
		t.Fatal("original C contracts differ", hash)
	}
}
