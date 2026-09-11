//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/libs/prand"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"os"
	"strconv"
	"testing"
)

func callbackBase(op int) legacy.PortTestRoamSpec {
	s := aiSpellBase(0)
	s.Op = 15
	s.Callbacks = &legacy.PortTestAICallbackSpec{Op: op, Range: math.Float32bits(50), Force: math.Float32bits(2), Damage: 7, DamageType: 1, PoisonPower: 1, PoisonMax: 5, TargetRadius: math.Float32bits(5), ActorRadius: math.Float32bits(5), Nearest: math.Float32bits(50), CloudLifetime: 2.5, LootName: "FanChakram"}
	s.Spells.TargetClass = 4
	s.Main.TargetPlayer = true
	s.Spells.TargetFlags = 4
	s.Combat.Friendly = false
	s.Combat.MeleeRange = math.Float32bits(50)
	s.MonsterState.Direction = 0
	return s
}
func callbackHash(t *testing.T, label string, r []legacy.PortTestRoamResult, want string) {
	t.Helper()
	for i, v := range r {
		if !v.Intact || !v.Callbacks.Intact || !v.Spells.Intact || !v.Combat.Intact || !v.MonsterState.Intact {
			t.Fatalf("case %d guards: outer=%v callback=%v spells=%v combat=%v state=%v", i, v.Intact, v.Callbacks.Intact, v.Spells.Intact, v.Combat.Intact, v.MonsterState.Intact)
		}
	}
	b, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if path := os.Getenv("OPENNOX_CALLBACK_CAPTURE"); path != "" {
		if err := os.WriteFile(path+"-"+label+".json", b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s: %d cases %s", label, len(r), got)
	if want == "" && os.Getenv("OPENNOX_CALLBACK_CASE") == "" && os.Getenv("OPENNOX_CALLBACK_CORPUS_CASE") == "" {
		want = callbackHashes[label]
	}
	if want != "" && got != want {
		t.Fatalf("hash %s want %s", got, want)
	}
}
func TestAICallbackSmoke(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for op := 0; op < 33; op++ {
		specs = append(specs, callbackBase(op))
	}
	if n := os.Getenv("OPENNOX_CALLBACK_CASE"); n != "" {
		i, err := strconv.Atoi(n)
		if err != nil {
			t.Fatal(err)
		}
		specs = specs[i : i+1]
	}
	callbackHash(t, "smoke", legacy.PortTestRoam(specs), "")
}

func TestAICallbackLoaders(t *testing.T) {
	got := legacy.PortTestCallbackLoaders()
	if len(got) != 88 {
		t.Fatalf("loader cases %d", len(got))
	}
	for _, r := range got {
		if !r.Intact || !r.FieldCorrect || r.Success != r.WantSuccess || r.Slot != r.WantSlot {
			t.Fatalf("%+v", r)
		}
	}
	b, _ := json.Marshal(got)
	if h := fmt.Sprintf("%x", sha256.Sum256(b)); h != "4c630d08dca8d187a0634cfdb80e6f25b3fae8c022aec015fe61ebe017380c67" {
		t.Fatalf("loader hash %s", h)
	}
}

func init() {
	legacy.PortTestCallbackServer = func(core *server.Server) (legacy.Server, func()) {
		old := doDamageWalls
		return &Server{Server: core}, func() { doDamageWalls = old }
	}
}

func callbackCorpus() []legacy.PortTestRoamSpec {
	b := math.Float32bits
	var all []legacy.PortTestRoamSpec
	for n := 0; n < 128; n++ {
		for op := 0; op < 33; op++ {
			s := callbackBase(op)
			v := s.Callbacks
			s.Seed = n + 1
			v.Range = []uint32{b(0), b(10), b(30), b(50), b(100)}[n%5]
			v.Force = []uint32{b(-1), 0, b(.01), b(2), b(10)}[n%5]
			v.PoisonChance = []uint32{0, 1, 50, 99, 100, 101, 0xffffffff}[n%7]
			v.PoisonPower = uint32(n % 4)
			v.PoisonMax = uint32(n % 7)
			v.TargetPoison = uint32(n % 6)
			v.TargetBuffs = []uint32{0, 1 << 23, 1 << 5}[n%3]
			v.Damage = uint32(n % 20)
			v.DamageType = uint32(n % 8)
			v.ActorRadius = b(float32(n % 10))
			v.TargetRadius = b(float32(n % 8))
			v.AllTargets = uint32(n % 2)
			v.Nearest = b(float32(n % 80))
			v.DebrisIndex = uint32(n % 2)
			v.BoneIndex = uint32(n % 2)
			v.CloudLifetime = []float64{0, .5, 2.5, 10}[n%4]
			v.SelfTarget = n%7 == 0
			s.Lifecycle.GameFlags = []uint32{0, 2048, 4096}[n%3]
			s.MonsterState.FPS = []uint32{1, 30, 60}[n%3]
			s.MonsterState.Frame = []uint32{0, 128, 0x80000000, 0xffffffff}[n%4]
			s.MonsterState.Direction = int16(n * 2)
			s.MonsterState.Stack = int8(n%21 + 1)
			s.Combat.Wall = n % 3
			s.Combat.Target = [2]uint32{b(float32(90 + n%90)), b(float32(70 + n%60))}
			s.Spells.TargetClass = []uint32{2, 4, 8}[n%3]
			s.Main.TargetPlayer = s.Spells.TargetClass == 4
			s.Spells.TargetFlags = []uint32{4, 5, 0x8004, 0x14}[n%4]
			s.Spells.TargetSubclass = []uint32{0, 0x200}[n%2]
			s.Spells.TargetMax = uint16(100 * (n % 2))
			s.Spells.NilTargetHealth = n%7 == 0
			s.Spells.Second = n%3 == 0
			s.Spells.SecondPos = [2]uint32{b(float32(110 + n%25)), b(100)}
			if n%5 == 0 {
				v.Enabled = map[string]bool{"skull": false}
			}
			if n%5 == 1 {
				v.Enabled = map[string]bool{"porttestdebrisa": false}
			}
			if n%5 == 2 {
				v.Enabled = map[string]bool{"porttestdebrisb": false}
			}
			all = append(all, s)
		}
	}
	return all
}
func TestAICallbackCorpus(t *testing.T) {
	specs := callbackCorpus()
	if n := os.Getenv("OPENNOX_CALLBACK_CORPUS_CASE"); n != "" {
		i, err := strconv.Atoi(n)
		if err != nil {
			t.Fatal(err)
		}
		specs = specs[i : i+1]
	}
	callbackHash(t, "corpus", legacy.PortTestRoam(specs), "")
}
func callbackSeed(roll int) int {
	for seed := 0; seed < 4096; seed++ {
		if prand.New(seed).IntClamp(0, 100) == roll {
			return seed
		}
	}
	panic("missing RNG roll")
}
func TestAICallbackLootContracts(t *testing.T) {
	type entry struct {
		s     legacy.PortTestRoamSpec
		want  uint32
		draws int
		name  string
	}
	var cases []entry
	for _, op := range []int{17, 18, 20, 21, 22, 23, 24} {
		for _, roll := range []int{20, 21, 25, 26, 50, 51} {
			for _, online := range []bool{false, true} {
				for _, missing := range []bool{false, true} {
					s := callbackBase(op)
					s.Seed = callbackSeed(roll)
					s.Callbacks.Enabled = map[string]bool{"skull": false}
					var want uint32
					name := ""
					switch op {
					case 17, 18, 22:
						if roll > 20 {
							if roll <= 50 {
								want = 15
								name = "sword"
							} else if op == 18 {
								want = 17
								name = "steelshield"
							} else {
								want = 16
								name = "woodenshield"
							}
						}
					case 20:
						if roll > 20 {
							if roll <= 50 {
								want = 19
								name = "bow"
							} else {
								want = 20
								name = "quiver"
							}
						}
					case 21:
						if roll > 25 {
							want = 21
							name = "ogreaxe"
						}
					case 23:
						if roll > 25 {
							want = 18
							name = "staffwooden"
						}
					case 24:
						if roll > 25 {
							want = 22
							name = "fanchakram"
						}
					}
					if online {
						s.Lifecycle.GameFlags = 2048
					}
					if missing && name != "" {
						s.Callbacks.Enabled[name] = false
					}
					if !online || missing {
						want = 0
					}
					draws := 1
					if want != 0 {
						draws++
					} // one random placement angle after successful allocation
					cases = append(cases, entry{s, want, draws, fmt.Sprintf("op%d-roll%d-online%v-missing%v", op, roll, online, missing)})
				}
			}
		}
	}
	for _, online := range []bool{false, true} {
		for _, missing := range []bool{false, true} {
			s := callbackBase(31)
			if online {
				s.Lifecycle.GameFlags = 2048
			}
			if missing {
				s.Callbacks.Enabled = map[string]bool{"fanchakram": false}
			}
			var want uint32
			draws := 0
			if online && !missing {
				want = 22
				draws = 1
			}
			cases = append(cases, entry{s, want, draws, fmt.Sprintf("direct-online%v-missing%v", online, missing)})
		}
	}
	specs := make([]legacy.PortTestRoamSpec, len(cases))
	for i, c := range cases {
		specs[i] = c.s
	}
	got := legacy.PortTestRoam(specs)
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			r := got[i]
			created := r.Lifecycle.Created
			if c.want == 0 {
				if len(created) != 0 {
					t.Fatal("unexpected creation")
				}
			} else if len(created) != 1 || created[0][1]&65535 != c.want {
				t.Fatalf("wrong loot for %s", c.name)
			}
			if r.Logic != c.s.Seed+c.draws {
				t.Fatalf("RNG index %d want %d", r.Logic, c.s.Seed+c.draws)
			}
			if c.want != 0 {
				mods := r.Callbacks.Modifiers
				if len(mods) != 6 {
					t.Fatal("modifier payload missing")
				}
				want0, want1, ammo := uint32(0), uint32(0), uint32(0)
				switch c.s.Callbacks.Op {
				case 17, 18:
					want0, want1 = 970, 972
				case 21:
					want0, want1 = 970, 972
				case 22:
					want0, want1 = 970, 971
				case 23:
					want0 = 970
				case 24, 31:
					ammo = 5 << 8
				}
				if c.want == 16 || c.want == 17 {
					want0 = 0
				}
				if mods[1] != want0 || mods[2] != want1 || mods[3] != 0 || mods[4] != 0 || mods[5] != ammo {
					t.Fatalf("modifiers/ammo %v want %d %d 0 0 %d", mods, want0, want1, ammo)
				}
			}

			if c.s.Callbacks.Op != 31 && r.Callbacks.Return != 1 {
				t.Fatal("callback success return")
			}
		})
	}
	callbackHash(t, "loot", got, "")
}

func TestAICallbackStrikeContracts(t *testing.T) {
	type entry struct {
		name  string
		s     legacy.PortTestRoamSpec
		check func(*testing.T, legacy.PortTestRoamResult)
	}
	var cases []entry
	add := func(name string, s legacy.PortTestRoamSpec, f func(*testing.T, legacy.PortTestRoamResult)) {
		cases = append(cases, entry{name, s, f})
	}
	b := math.Float32bits
	for op := 0; op < 11; op++ {
		for _, kind := range []int{0, 1, 2, 3, 4} {
			s := callbackBase(op)
			switch kind {
			case 1:
				s.Combat.Target = [2]uint32{b(300), b(100)}
			case 2:
				s.Combat.Target = [2]uint32{b(60), b(100)}
			case 3:
				s.Combat.Target = [2]uint32{b(200), b(100)}
				s.Callbacks.Range = b(150)
				s.Combat.Wall = 1
			case 4:
				s.Combat.Target = [2]uint32{b(200), b(100)}
				s.Callbacks.Range = b(150)
				s.Combat.Wall = 2
			}
			hit := (kind == 0 || kind == 4) && op != 9
			want := uint32(1)
			if !hit && (op == 0 || op == 3 || op == 4 || op == 5 || op == 6 || op == 7) {
				want = 0
			}
			if kind == 3 && op != 9 {
				want = 0
			} // selected target with failed ray returns zero even for default/ghost
			add(fmt.Sprintf("strike-%d-geometry-%d", op, kind), s, func(t *testing.T, r legacy.PortTestRoamResult) {
				if (len(r.Callbacks.Damage) == 5) != hit || r.Callbacks.Return != want {
					t.Fatalf("hit/return: %+v want hit=%v return=%d", r.Callbacks, hit, want)
				}
			})
		}
	}
	for _, op := range []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 10} {
		for _, force := range []float32{-1, 0, 3} {
			s := callbackBase(op)
			s.Callbacks.MutateOnDamage = true
			s.Callbacks.ForceAfterDamage = b(force)
			add(fmt.Sprintf("force-after-damage-%d-%g", op, force), s, func(t *testing.T, r legacy.PortTestRoamResult) {
				calls := 0
				for i, v := range r.Trace {
					if v == 47 && i+5 < len(r.Trace) {
						calls++
						if r.Trace[i+4] != b(force) {
							t.Fatal("force was read before damage callback")
						}
					}
				}
				want := 0
				if op == 0 || force > 0 {
					want = 1
				}
				if calls != want {
					t.Fatalf("force calls %d want %d", calls, want)
				}
				if r.Callbacks.Definition[120/4] != b(force) {
					t.Fatal("damage callback mutation lost")
				}
			})
		}
	}
	for _, op := range []int{1, 2, 5, 6, 7} {
		s := callbackBase(op)
		s.Callbacks.PoisonChance = 100
		add(fmt.Sprintf("poison-force-order-%d", op), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := uint32(0)
			if op == 5 {
				want = 1
			}
			found := false
			for i, v := range r.Trace {
				if v == 47 && i+5 < len(r.Trace) {
					found = true
					if r.Trace[i+5] != want {
						t.Fatal("poison/force order")
					}
				}
			}
			if !found {
				t.Fatal("force missing")
			}
		})
	}
	for _, kind := range []int{0, 1, 2, 3} {
		s := callbackBase(27)
		s.Callbacks.AllTargets = 1
		want := uint32(100)
		switch kind {
		case 0:
			s.Spells.NilTargetHealth = true
			want = 0 // ordinary class bits numerically convert from float to zero
		case 1:
			s.Spells.TargetFlags = 1 // flags are likewise numerically converted, not reinterpreted
		case 2:
			s.Spells.TargetFlags = b(17)
			want = 0
		case 3:
			s.Spells.TargetFlags = b(17)
			s.Spells.TargetClass = b(6)
			s.Spells.NilTargetHealth = true
		}
		add(fmt.Sprintf("numeric-float-byte-gates-%d", kind), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Callbacks.Globals[2] != want {
				t.Fatalf("selected %d want %d", r.Callbacks.Globals[2], want)
			}
		})
	}
	gap := float32(float64(float32(40+.001)) - 10)
	for _, nearest := range []float32{math.Nextafter32(gap, float32(math.Inf(-1))), gap, math.Nextafter32(gap, float32(math.Inf(1)))} {
		s := callbackBase(27)
		s.Callbacks.Nearest = b(nearest)
		want := uint32(0)
		if nearest > gap {
			want = 100
		}
		add(fmt.Sprintf("strict-nearest-%08x", b(nearest)), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Callbacks.Globals[2] != want {
				t.Fatal("strict nearest tie/spill")
			}
		})
	}
	for _, secondX := range []float32{130, 140, 150} {
		s := callbackBase(26)
		s.Spells.Second = true
		s.Spells.SecondPos = [2]uint32{b(secondX), b(100)}
		want := uint32(102)
		if secondX > 140 {
			want = 100
		}
		add(fmt.Sprintf("nearest-target-%g", secondX), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Callbacks.Return != want {
				t.Fatalf("target %d want %d", r.Callbacks.Return, want)
			}
		})
	}
	specs := make([]legacy.PortTestRoamSpec, len(cases))
	for i, c := range cases {
		specs[i] = c.s
	}
	got := legacy.PortTestRoam(specs)
	for i, c := range cases {
		t.Run(c.name, func(t *testing.T) { c.check(t, got[i]) })
	}
	callbackHash(t, "strike", got, "")
}

func TestAICallbackPoisonContracts(t *testing.T) {
	type entry struct {
		s      legacy.PortTestRoamSpec
		ret    uint32
		draws  int
		poison byte
	}
	var all []entry
	for _, chance := range []uint32{0, 1, 50, 100, 101, 0xffffffff} {
		for _, roll := range []int{1, 50, 100} {
			for _, immune := range []bool{false, true} {
				s := callbackBase(28)
				for seed := 0; seed < 4096; seed++ {
					if prand.New(seed).IntClamp(1, 100) == roll {
						s.Seed = seed
						break
					}
				}
				s.Callbacks.PoisonChance = chance
				if immune {
					s.Callbacks.TargetBuffs = 1 << 23
				}
				e := entry{s: s}
				if chance != 0 {
					e.draws = 1
					if int32(chance) >= int32(roll) && !immune {
						e.draws = 2
						e.ret = 1
						e.poison = 1
					}
				}
				all = append(all, e)
			}
		}
	}
	specs := make([]legacy.PortTestRoamSpec, len(all))
	for i, c := range all {
		specs[i] = c.s
	}
	got := legacy.PortTestRoam(specs)
	for i, c := range all {
		r := got[i]
		if r.Callbacks.Return != c.ret || r.Logic != c.s.Seed+c.draws || byte(r.MonsterState.Objects[0][540/4]) != c.poison {
			t.Fatalf("case %d poison return/RNG/value", i)
		}
		want := uint32(0)
		if c.poison != 0 {
			want = 1024
		}
		if r.Callbacks.PlayerStatus != want {
			t.Fatalf("case %d player poison status", i)
		}
	}
	callbackHash(t, "poison", got, "")
}

func TestAICallbackDebrisContracts(t *testing.T) {
	type entry struct {
		s     legacy.PortTestRoamSpec
		types []uint32
		draws int
		rot   uint32
	}
	var all []entry
	for _, op := range []int{14, 15, 32} {
		for _, online := range []bool{false, true} {
			for _, start := range []uint32{0, 1} {
				for _, missing := range []int{0, 1, 2, 3} {
					s := callbackBase(op)
					s.Seed = 7
					s.Callbacks.DebrisIndex = start
					s.Callbacks.BoneIndex = start
					if online {
						s.Lifecycle.GameFlags = 2048
					}
					s.Callbacks.Enabled = map[string]bool{}
					switch missing {
					case 1:
						s.Callbacks.Enabled["porttestdebrisa"] = false
					case 2:
						s.Callbacks.Enabled["porttestdebrisb"] = false
					case 3:
						s.Callbacks.Enabled["skull"] = false
					}
					e := entry{s: s, rot: start}
					count, each := 2, 4
					if op == 15 {
						count, each = 6, 5
						if online {
							count = prand.New(s.Seed).IntClamp(20, 30)
							e.draws = 1
						}
					}
					if op == 32 {
						if missing == 3 {
							all = append(all, e)
							continue
						}
						e.types = append(e.types, 12)
						e.draws = 5
						each = 5
						if online {
							count = prand.New(s.Seed+4).IntClamp(10, 20)
						} else {
							count = prand.New(s.Seed+4).IntClamp(5, 10)
						}
					}
					for i := 0; i < count; i++ {
						index := (int(start) + i) % 2
						if op == 14 {
							index = i
						}
						if missing == index+1 {
							break
						}
						e.types = append(e.types, uint32(13+index))
						e.draws += each
						if op != 14 {
							e.rot = uint32((index + 1) % 2)
						}
					}
					all = append(all, e)
				}
			}
		}
	}
	specs := make([]legacy.PortTestRoamSpec, len(all))
	for i, e := range all {
		specs[i] = e.s
	}
	got := legacy.PortTestRoam(specs)
	for i, e := range all {
		r := got[i]
		if r.Logic != e.s.Seed+e.draws {
			t.Fatalf("case%d RNG %d want%d", i, r.Logic, e.s.Seed+e.draws)
		}
		if len(r.Lifecycle.Created) != len(e.types) {
			t.Fatalf("case%d count %d want%d", i, len(r.Lifecycle.Created), len(e.types))
		}
		for j, v := range r.Lifecycle.Created {
			if v[1]&65535 != e.types[j] {
				t.Fatalf("case%d creation order", i)
			}
		}
		idx := 7
		if e.s.Callbacks.Op == 32 {
			idx = 8
		}
		if r.Callbacks.Globals[idx] != e.rot {
			t.Fatalf("case%d rotating index", i)
		}
	}
	callbackHash(t, "debris", got, "")
}

func TestAICallbackPrecisionContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, xy := range [][2]uint32{{0x4c80000d, 0x4cddb3e4}, {0x4c7fffec, 0x4cddb3bc}} {
		s := callbackBase(27)
		s.Combat.Target = xy
		s.Spells.TargetFlags = 0
		s.Callbacks.Nearest = math.Float32bits(2e8)
		specs = append(specs, s)
	}
	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		want := uint32(0)
		if i == 0 {
			want = 100
		}
		if r.Callbacks.Globals[2] != want {
			t.Fatalf("precision case%d selected%d want%d", i, r.Callbacks.Globals[2], want)
		}
	}
	callbackHash(t, "precision", got, "")
}
func TestAICallbackCloudContracts(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, fps := range []uint32{1, 30, 65537, 0x80000000, 0xffffffff} {
		for _, duration := range []float64{0, .25, 2.5, -.75} {
			s := callbackBase(19)
			s.MonsterState.FPS = fps
			s.Callbacks.CloudLifetime = duration
			specs = append(specs, s)
		}
	}
	got := legacy.PortTestRoam(specs)
	for i, r := range got {
		s := specs[i]
		v := float64(float32(s.Callbacks.CloudLifetime * float64(int32(s.MonsterState.FPS))))
		want := uint32(0x80000000)
		if v >= -2147483648 && v < 2147483648 {
			want = uint32(int32(math.Trunc(v)))
		}
		if len(r.Callbacks.CreatedData) != 2 || r.Callbacks.CreatedData[1] != want || r.Callbacks.Globals[6] != 11 {
			t.Fatalf("cloud case%d: %+v want%d", i, r.Callbacks, want)
		}
	}
	callbackHash(t, "cloud", got, "")
}

var callbackHashes = map[string]string{
	"smoke":     "9c728bf2e7aca4122a4fe999d773a25a95f11a5a8d56dd52ef622865d8e819df",
	"corpus":    "dfb3eabe2a2740dce4e5b72ce5f58b646450aee06716bd1702a8703920263376",
	"loot":      "656c517c1c229997a5451ed1e965f0552c538935c8c90e495b37ca141fa9af70",
	"strike":    "643d99fbd75072d3379992f93d7252c457440574afc170c1f532429d2cc2dffb",
	"poison":    "0935568717e4b8b976cb6ea4ffc2b1acf30e1fb5619f5da0b04d75b6678c2e2e",
	"debris":    "56b8186e6e254e6ce2cd46ce1fb6af07cf271d9aa40d63c85562b4ba9637d196",
	"precision": "313607d3b331bd5ba6e4e8c551c65edc8bdf2aab37871762e14a08445d6f933d",
	"cloud":     "a8ab0325f416409b1ffcc70e7cd9537fdf83a1e00f7a6c6b233c8dbb5785f987",
}
