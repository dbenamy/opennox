//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func lifecycleBase(op, phase int) legacy.PortTestRoamSpec {
	b := math.Float32bits
	return legacy.PortTestRoamSpec{Op: 11, Seed: 1, Stack: 1,
		Owner:     &legacy.PortTestRoamOwnerSpec{Frame: 123, FPS: 30, X: b(100), Y: b(100), Buffs: 1 << 29},
		Combat:    &legacy.PortTestCombatSpec{Sound: true, Target: [2]uint32{b(110), b(100)}, TargetClass: 0x10, TargetFlags: 4, Direction: 73},
		Lifecycle: &legacy.PortTestLifecycleSpec{Op: op, Phase: phase, Cur: 3, Max: 100, DelayMin: 5, DelayMax: 5, Range: b(100), Duration: 5, Frame137: 117, HasTarget: true, Place: true, Use: true, Eligible: true, ItemSubclass: 0x10}}
}
func lifecycleCorpus() []legacy.PortTestRoamSpec {
	b := math.Float32bits
	var out []legacy.PortTestRoamSpec
	for n := 0; n < 512; n++ {
		for op := 0; op < 11; op++ {
			s := lifecycleBase(op, 0)
			l, c, o := s.Lifecycle, s.Combat, s.Owner
			s.Seed = n + 1
			s.Stack = int8(n % 24)
			l.Phase = n % 4
			l.Kind = n % 3
			l.Callback = n%2 == 0
			l.Soul = n%3 != 0
			l.Subclass = []uint32{0, 0x10, 0x80, 0x10000, 0x10090}[n%5]
			l.ObjectFlags = []uint32{0, 0x8058, 0x10000, 0x200}[n%4]
			l.DeleteDef = n%2 != 0
			l.Updatable = n%3 != 0
			l.PlayerOwner = n%5 == 0
			l.HistoryCount = uint32(n % 17)
			l.SeenCount = byte(n % 17)
			l.Motion = [6]uint32{b(1), b(-2), b(3), b(-4), 0x80000000, b(6)}
			l.Frame137 = []uint32{0, 117, 118, 119, 0xfffffff0}[n%5]
			l.Duration = []uint32{0, 5, 6, 0xffffffff}[n%4]
			l.DelayMin = []float64{0, 2.49999999, 2.50000001, 7.5}[n%4]
			l.DelayMax = l.DelayMin + 3
			l.GameFlags = []uint32{0, 0x2000, 0x8000000, 0x200000}[n%4]
			o.Status = []uint32{0, 0x80000, 0x100000, 0x180000}[n%4]
			o.Enemy = n%2 == 0
			o.Frame = []uint32{0, 123, 0xffffffff, 2}[n%4]
			c.Done = byte(n % 2)
			c.Wall = n % 3
			l.HasTarget = n%7 != 0
			l.Place = n%2 == 0
			l.Swap = n%3 == 0
			l.Use = n%5 != 0
			l.ItemSubclass = []uint32{0, 4, 8, 0x10, 0x18, 0x80, 0x90, 0x98}[n%8]
			l.Poison = byte(n % 3)
			l.Eligible = n%2 == 0
			c.TargetClass = []uint32{0, 0x10, 0x1000000, 0x1000010}[n%4]
			c.Target[0] = []uint32{b(110), b(175), math.Float32bits(math.Nextafter32(175, 0)), b(176), b(200)}[n%5]
			c.Scan = op == 9 || op == 10
			out = append(out, s)
		}
	}
	return out
}
func lifecycleHash(t *testing.T, label string, got []legacy.PortTestRoamResult, baseline string) {
	t.Helper()
	for i, r := range got {
		if !r.Intact {
			t.Fatalf("%s case %d changed guarded/read-only state: %+v", label, i, r)
		}
	}
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	h := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s cases=%d complete-state-sha256=%s", label, len(got), h)
	if path := os.Getenv("OPENNOX_LIFECYCLE_CAPTURE"); path != "" {
		if err := os.WriteFile(path+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	if baseline != "" && os.Getenv("OPENNOX_LIFECYCLE_CASE") == "" && h != baseline {
		t.Fatal("original C state differs", h)
	}
}
func TestAILifecycleOriginalState(t *testing.T) {
	specs := lifecycleCorpus()
	if v := os.Getenv("OPENNOX_LIFECYCLE_CASE"); v != "" {
		var n int
		if _, err := fmt.Sscan(v, &n); err != nil {
			t.Fatal(err)
		}
		specs = specs[n : n+1]
	}
	lifecycleHash(t, "corpus", legacy.PortTestRoam(specs), "15d5986a08af2c44eb120f178e45366cb014817a42c4788a21db0d892ce4c485")
}
func lifeWord(r legacy.PortTestRoamResult, off, initial uint32) uint32 {
	for i := 0; i < len(r.Changes); i += 2 {
		if r.Changes[i] == off {
			return r.Changes[i+1]
		}
	}
	return initial
}
func TestAILifecycleContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	var checks []func(*testing.T, legacy.PortTestRoamResult)
	add := func(s legacy.PortTestRoamSpec, f func(*testing.T, legacy.PortTestRoamResult)) {
		specs = append(specs, s)
		checks = append(checks, f)
	}
	for _, op := range []int{2, 4} {
		for _, done := range []byte{0, 1} {
			s := lifecycleBase(op, 0)
			s.Combat.Done = done
			add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
				if r.Stack != 1-int8(done) {
					t.Fatal("animation completion", r.Stack)
				}
			})
		}
	}
	for kind := 0; kind < 3; kind++ {
		for _, soul := range []bool{false, true} {
			for _, sub := range []uint32{0, 0x10000} {
				s := lifecycleBase(3, 1)
				s.Lifecycle.Kind = kind
				s.Lifecycle.Callback = true
				s.Lifecycle.Soul = soul
				s.Lifecycle.Subclass = sub
				s.Lifecycle.Motion = [6]uint32{1, 2, 3, 4, 5, 6}
				add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
					if !reflect.DeepEqual(r.Lifecycle.Calls, []uint32{2, 0, 0, 0, 0, 0, 0}) {
						t.Fatal("dead callback preceded motion clearing", r.Lifecycle.Calls)
					}
					flags := uint32(0x18)
					rng := 1
					if kind != 0 {
						flags = 0x10
						rng = 2
						if lifeWord(r, 4096+492, 5) != 5 {
							t.Fatal("duration")
						}
					}
					if lifeWord(r, 16, 0) != flags || r.Logic != rng {
						t.Fatal("death flags/RNG", r)
					}
					want := 0
					if soul && sub != 0 {
						want = 1
					}
					if len(r.Lifecycle.Created) != want {
						t.Fatal("soul gate", r.Lifecycle.Created)
					}
					if want != 0 && r.Lifecycle.Created[0][31] != 73|73<<16 {
						t.Fatal("soul direction")
					}
				})
			}
		}
	}
	for _, kind := range []int{0, 1, 2} {
		for _, enemy := range []bool{false, true} {
			for _, age := range []uint32{4, 5, 6, 0xffffffff} {
				for _, status := range []uint32{0, 0x100000} {
					s := lifecycleBase(3, 0)
					s.Lifecycle.Kind = kind
					s.Lifecycle.Frame137 = s.Owner.Frame - age
					s.Owner.Enemy = enemy
					s.Owner.Status = status
					s.Lifecycle.ObjectFlags = 0x8058
					s.Lifecycle.Updatable = true
					s.Lifecycle.HistoryCount = 16
					s.Lifecycle.SeenCount = 16
					add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
						if kind == 0 {
							if r.Lifecycle.Updatable != 0 || lifeWord(r, 744, 1) != 0 || lifeWord(r, 4096+1196, func() uint32 {
								if enemy {
									return 100
								}
								return 0
							}()) != 0 {
								t.Fatal("dead cleanup", r)
							}
							return
						}
						raised := enemy && age > 5 && status == 0
						if raised {
							if r.Stack != 2 || r.Combat.Actions[1].Action != 61 || r.Combat.Actions[2].Action != 35 || r.Lifecycle.Health[0]&0xffff != 100 || lifeWord(r, 16, 0x8058) != 0 {
								t.Fatal("revival", r)
							}
						} else if r.Stack != 1 || len(r.Combat.Sounds) != 0 {
							t.Fatal("unexpected revival", r)
						}
					})
				}
			}
		}
	}
	for _, place := range []bool{false, true} {
		for _, swap := range []bool{false, true} {
			for _, x := range []float32{110, 175, 176} {
				s := lifecycleBase(1, 0)
				s.Combat.Target[0] = math.Float32bits(x)
				s.Lifecycle.Place = place
				s.Lifecycle.Swap = swap
				add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
					if r.Stack != 0 {
						t.Fatal("pickup must pop")
					}
					var want []uint32
					if x < 175 {
						target := uint32(100)
						if swap {
							target = 102
						}
						want = []uint32{3, 101, target}
					}
					if !reflect.DeepEqual(r.Lifecycle.Calls, want) {
						t.Fatal("pickup/use after placement", r.Lifecycle.Calls, want)
					}
				})
			}
		}
	}
	for _, op := range []int{9, 10} {
		for _, eligible := range []bool{false, true} {
			for _, poison := range []byte{0, 1} {
				for _, sub := range []uint32{0, 4, 8, 0x10, 0x18, 0x80, 0x98} {
					s := lifecycleBase(op, 0)
					s.Combat.Scan = true
					s.Lifecycle.ItemSubclass = sub
					s.Lifecycle.Poison = poison
					s.Lifecycle.Eligible = eligible
					if op == 10 {
						s.Combat.TargetClass = 0x1000000
					}
					add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
						yes := eligible
						if op == 9 {
							yes = sub&4 == 0 && sub&8 == 0 && (sub&0x80 == 0 || poison != 0)
						}
						want := 0
						if yes {
							want = 100
						}
						if r.Return != want {
							t.Fatal("item filter", r.Return, want)
						}
					})
				}
			}
		}
	}

	for _, stack := range []int8{0, 1, 23} {
		s := lifecycleBase(0, 0)
		s.Stack = stack
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if lifeWord(r, 4096+1304, 0) != 0x3f547ae1 {
				t.Fatal("hunt aggression")
			}
			want := stack
			if stack < 23 {
				want++
			}
			if r.Stack != want {
				t.Fatal("hunt stack", r.Stack, want)
			}
			if stack < 23 {
				a := r.Combat.Actions[len(r.Combat.Actions)-1]
				if a.Action != 10 || a.Args[0] != 0 || a.Args[2] != 128 {
					t.Fatal("hunt roam", a)
				}
			}
		})
	}
	for _, v := range []float64{0, 2.49999999, 2.50000001, 3.5, 7.5, 16777217} {
		s := lifecycleBase(3, 1)
		s.Lifecycle.Kind = 1
		s.Lifecycle.DelayMin = v
		s.Lifecycle.DelayMax = v
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := uint32(math.Trunc(float64(float32(v))))
			if lifeWord(r, 4096+492, 5) != want || r.Logic != 2 {
				t.Fatal("float32 duration conversion", v, want, r)
			}
		})
	}
	for _, owner := range []bool{false, true} {
		for _, op := range []int{2, 3, 8} {
			s := lifecycleBase(op, 0)
			if op == 2 {
				s.Lifecycle.Phase = 2
			}
			s.Lifecycle.Kind = 1
			s.Owner.Status = 0x80000
			s.Lifecycle.PlayerOwner = owner
			s.Lifecycle.Subclass = 0x10080
			add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
				if len(r.Lifecycle.Created) != 1 || r.Lifecycle.Created[0][1]&0xffff != 6 || r.Lifecycle.Decay != 1001 {
					t.Fatal("burn scorch", r.Lifecycle)
				}
				if owner && lifeWord(r, 12, 0x10080) != 0x10000 {
					t.Fatal("owner marker not cleared")
				}
				shield, spark := false, false
				for _, p := range r.Lifecycle.Packets {
					if len(p) == 3 && p[0] == 109 {
						shield = true
					}
					if len(p) >= 6 && p[0] == 0x93 && p[5] == 100 {
						spark = true
					}
				}
				if shield != owner || spark != (op == 3) {
					t.Fatal("burn packets", r.Lifecycle.Packets)
				}
				if r.Logic != 3 {
					t.Fatal("scorch RNG draws", r.Logic)
				}
			})
		}
	}
	for _, op := range []int{9, 10} {
		for _, x := range []float32{105, 110, 115} {
			s := lifecycleBase(op, 0)
			s.Combat.Scan = true
			s.Lifecycle.Second = true
			s.Lifecycle.SecondPos = [2]uint32{math.Float32bits(x), math.Float32bits(100)}
			if op == 10 {
				s.Combat.TargetClass = 0x1000000
			}
			add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
				if x <= 110 && r.Return != 102 || x > 110 && r.Return != 100 {
					t.Fatal("nearest candidate", r.Return, x)
				}
			})
		}
	}

	for _, op := range []int{9, 10} {
		s := lifecycleBase(op, 0)
		s.Owner.X = math.Float32bits(100)
		s.Owner.Y = 0xb3800000
		s.Combat.Target = [2]uint32{math.Float32bits(100), 0x3f800001}
		s.Combat.Scan = true
		if op == 10 {
			s.Combat.TargetClass = 0x1000000
		}
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			off := uint32(2489444)
			if op == 10 {
				off = 2489448
			}
			minimum := uint32(0)
			for i := 0; i+1 < len(r.Trace); i++ {
				if r.Trace[i] == off {
					minimum = r.Trace[i+1]
					break
				}
			}
			if r.Return != 100 || minimum != 0x3f800003 {
				t.Fatal("search full-delta precision", r.Return, minimum)
			}
		})
	}

	for _, kind := range []int{0, 1, 2} {
		for _, action := range []uint32{10, 31, 35} {
			s := lifecycleBase(5, 0)
			s.Lifecycle.Kind = kind
			s.Lifecycle.HeadAction = action
			s.Lifecycle.ObjectFlags = 0x10058
			add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
				want := 0
				if kind != 0 {
					want = int(action)
					if action == 31 {
						want = 0x10000
					}
				}
				if r.Return != want {
					t.Fatal("raise ABI return", r.Return, want)
				}
			})
		}
	}
	for _, callback := range []bool{false, true} {
		s := lifecycleBase(2, 1)
		s.Lifecycle.Callback = callback
		s.Lifecycle.Motion[0] = 123
		add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if !reflect.DeepEqual(r.Combat.Sounds, []uint32{315}) || len(r.Trace) < 5 || !reflect.DeepEqual(r.Trace[:5], []uint32{10, 7, 1264, 0, 1}) {
				t.Fatal("dying sound/script", r.Trace)
			}
			var want []uint32
			if callback {
				want = []uint32{1, 123}
			}
			if !reflect.DeepEqual(r.Lifecycle.Calls, want) {
				t.Fatal("die callback", r.Lifecycle.Calls)
			}
		})
	}
	for _, game := range []uint32{0, 0x2000} {
		for _, actor := range []uint32{0, 0x10} {
			for _, poison := range []byte{0, 1} {
				s := lifecycleBase(9, 0)
				s.Combat.Scan = true
				s.Lifecycle.ItemSubclass = 0x98
				s.Lifecycle.GameFlags = game
				s.Lifecycle.Subclass = actor
				s.Lifecycle.Poison = poison
				add(s, func(t *testing.T, r legacy.PortTestRoamResult) {
					want := 0
					if game != 0 && actor != 0 && poison != 0 {
						want = 100
					}
					if r.Return != want {
						t.Fatal("special edible gate", r.Return, want)
					}
				})
			}
		}
	}

	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		t.Run(fmt.Sprint(i), func(t *testing.T) { checks[i](t, r) })
	}
	lifecycleHash(t, "contracts", got, "b6509bbe7747b05c7b8f48f58bc8505f520ab2510743f731b33292cf8b42a6e8")
}
