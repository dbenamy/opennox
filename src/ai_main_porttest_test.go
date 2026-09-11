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

func mainBase(op int) legacy.PortTestRoamSpec {
	b := math.Float32bits
	s := stateBase(0)
	s.Op = 13
	s.Main = &legacy.PortTestMainSpec{Op: op, Class: 2, Danger: 0x12345678, Nearest: b(1e9), TargetType: 1, TargetNewPos: [2]uint32{b(140), b(100)}}
	s.MonsterState.Frame = 128
	s.MonsterState.FPS = 30
	s.MonsterState.Deadline = 128
	s.MonsterState.HealthCur = 100
	s.MonsterState.HealthMax = 100
	s.MonsterState.Aggression = b(.5)
	s.Combat.TargetClass = 4
	s.Combat.TargetFlags = 0
	s.Combat.Radius = b(5)
	s.Combat.Target = [2]uint32{b(140), b(100)}
	s.Combat.Velocity = [2]uint32{b(-2), 0}
	s.Main.Words = [][2]uint32{{496, 128}, {500, b(100)}, {504, b(100)}}
	return s
}
func mainCorpus() []legacy.PortTestRoamSpec {
	b := math.Float32bits
	var out []legacy.PortTestRoamSpec
	for n := 0; n < 256; n++ {
		for op := 0; op < 7; op++ {
			s := mainBase(op)
			v := s.Main
			m := s.MonsterState
			s.Seed = n + 1
			m.Action = []uint32{0, 1, 4, 7, 8, 9, 10, 13, 16, 17, 18, 20, 21, 23, 24, 25, 26, 27, 28, 30, 31, 36, 37}[n%23]
			m.Stack = int8(n%20 + 1)
			m.Frame = []uint32{0, 15, 16, 127, 128, 0x7fffffff, 0x80000000, 0xffffffff}[n%8]
			m.FPS = []uint32{1, 30, 60}[n%3]
			m.NetCode = uint32(n % 32)
			m.MimicAge = uint32(n % 16)
			m.Deadline = m.Frame - []uint32{0, 1, 89, 90, 91}[n%5]
			m.Aggression = []uint32{0, b(.08), b(.33), b(.5), b(.66), b(1), 0x7fc12345}[n%7]
			m.Speed = []uint32{0, b(.01), b(1), 0x7fc12345}[n%4]
			m.Status = []uint32{0, 4, 0x20, 0x40, 0x80, 0x200}[n%6]
			m.ObjFlags = []uint32{0, 4, 0x8000}[n%3]
			m.HealthCur = uint16(n % 101)
			if op == 3 && m.Action == 31 {
				m.ObjFlags |= 0x8000
			}
			s.Owner.Buffs = []uint32{0, 1 << 3, 1 << 11, 1 << 29}[n%4]
			s.Owner.Enemy = n%2 == 0
			m.Subclass = []uint32{0, 0x10, 0x80, 0x200, 0x400}[n%5]
			m.PlayerWeapon = []uint32{0, 0x400, 0x1000}[n%3]
			m.PlayerShield = []uint32{0, 0x1000000}[n%2]
			m.Direction = int16(n)
			v.TargetType = []uint32{1, 10, 11, 0xffff}[n%4]
			v.TargetSubclass = uint32(n%4) << 9
			s.Combat.TargetClass = []uint32{1, 4, 0x2000, 0x10000}[n%4]
			v.Words = append(v.Words, [2]uint32{524, []uint32{0, b(300), b(301), 0x7fc12345}[n%4]}, [2]uint32{528, m.Frame - uint32(n%3)}, [2]uint32{1336, []uint32{0, b(.3), b(.5), b(1)}[n%4]}, [2]uint32{1356, []uint32{0, b(20), b(100)}[n%3]})
			v.CacheCloud = []uint32{0, 10, 1}[n%3]
			v.CacheSmall = []uint32{0, 11, 4}[n%3]
			if op == 2 {
				v.Actions = []uint32{1, 25, 26, 27, 28, 16, 17, 18, 19, 20}
				if n%2 == 0 {
					v.Actions[5] = 4
				}
			}
			if op >= 5 {
				s.Combat.Shield = true
				v.Missile = true
				m.ObjFlags = 0x80
				m.Direction = 0
				v.SecondMissile = n%3 == 0
				v.SecondPos = [2]uint32{b(110 + float32(n%50)), b(100)}
				v.SecondVelocity = [2]uint32{b(-2), 0}
			}
			out = append(out, s)
		}
	}
	return out
}
func mainHash(t *testing.T, label string, got []legacy.PortTestRoamResult, want string) {
	t.Helper()
	for i, r := range got {
		if !r.Intact || !r.Main.Intact || !r.MonsterState.Intact || !r.Combat.Intact {
			t.Fatalf("%s case %d guards: %s", label, i, mustMainJSON(r))
		}
	}
	data, err := json.Marshal(got)
	if err != nil {
		t.Fatal(err)
	}
	if p := os.Getenv("OPENNOX_MAIN_CAPTURE"); p != "" {
		if err := os.WriteFile(p+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	sum := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s cases=%d complete-state=%s", label, len(got), sum)
	if want != "" && sum != want {
		t.Fatalf("want %s got %s", want, sum)
	}
}
func mustMainJSON(v any) string { b, _ := json.Marshal(v); return string(b) }
func TestAIMainCorpus(t *testing.T) {
	specs := mainCorpus()
	if s := os.Getenv("OPENNOX_MAIN_CASE"); s != "" {
		var i int
		if _, e := fmt.Sscanf(s, "%d", &i); e != nil {
			t.Fatal(e)
		}
		specs = specs[i : i+1]
	}
	want := "9a775ff5fa2e30b043c7d7f076e9d0d0eaf0d7fde28d47914de2a97a471adc6f"
	if os.Getenv("OPENNOX_MAIN_CASE") != "" {
		want = ""
	}
	mainHash(t, "corpus", legacy.PortTestRoam(specs), want)
}
func TestAIMainContracts(t *testing.T) {
	b := math.Float32bits
	type contract struct {
		name  string
		spec  legacy.PortTestRoamSpec
		check func(*testing.T, legacy.PortTestRoamResult)
	}
	var all []contract
	add := func(name string, s legacy.PortTestRoamSpec, check func(*testing.T, legacy.PortTestRoamResult)) {
		all = append(all, contract{name, s, check})
	}
	for _, action := range []uint32{0, 4, 1, 7} {
		for _, due := range []bool{false, true} {
			s := mainBase(0)
			s.MonsterState.Action = action
			deadline := uint32(129)
			if due {
				deadline = 128
			}
			s.Main.Words = append(s.Main.Words, [2]uint32{528, deadline})
			plays := due && (action == 0 || action == 4)
			add(fmt.Sprintf("idle-%d-%t", action, due), s, func(t *testing.T, r legacy.PortTestRoamResult) {
				if (len(r.Combat.Sounds) > 0) != plays || r.Logic != map[bool]int{false: 1, true: 2}[plays] {
					t.Fatalf("idle sound/RNG: %s", mustMainJSON(r))
				}
			})
		}
	}
	for _, actions := range [][]uint32{{1, 16, 17, 18, 19, 20, 25, 26, 27, 28}, {1, 16, 4, 25, 26}, {4, 1}, {0}} {
		s := mainBase(2)
		s.Main.Actions = actions
		want := len(actions) - 1
		for want > 0 {
			a := actions[want]
			if !(a >= 16 && a <= 20 || a >= 25 && a <= 28) {
				break
			}
			want--
		}
		add(fmt.Sprint("unwind-", actions), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if int(r.Stack) != want {
				t.Fatalf("stack %d want %d", r.Stack, want)
			}
		})
	}
	for _, tile := range []uint32{0, 6} {
		s := mainBase(4)
		s.Main.Tile = tile
		add(fmt.Sprint("dodge-tile-", tile), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := uint32(1)
			other := 3
			if tile == 6 {
				want = 0
				other = 11
			}
			if r.Main.Return != want || r.Logic != other {
				t.Fatalf("dodge result/RNG: %s", mustMainJSON(r))
			}
			if want == 1 {
				a := r.Combat.Actions
				if a[len(a)-1].Action != 9 || a[len(a)-2].Action != 41 {
					t.Fatal("dodge action order")
				}
			}
		})
	}
	for _, delta := range []uint32{0, 1, 15, 16} {
		s := mainBase(3)
		s.MonsterState.Action = 0
		s.MonsterState.Frame = 128 + delta
		s.Owner.Buffs = 1 << 11
		add(fmt.Sprint("cadence-", delta), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := uint32(0)
			if delta&15 == 0 {
				want = 24
			}
			a := r.Combat.Actions
			if a[len(a)-1].Action != want {
				t.Fatalf("cadence head %d want %d", a[len(a)-1].Action, want)
			}
		})
	}
	for _, buff := range []uint32{1 << 3, 1 << 11, (1 << 3) | (1 << 11)} {
		s := mainBase(3)
		s.Owner.Buffs = buff
		add(fmt.Sprint("buff-", buff), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			a := r.Combat.Actions
			want := uint32(36)
			if buff&(1<<11) != 0 {
				want = 24
			}
			if a[len(a)-1].Action != want {
				t.Fatalf("buff head %d want %d", a[len(a)-1].Action, want)
			}
		})
	}
	for _, hp := range []uint16{49, 50, 51} {
		s := mainBase(3)
		s.MonsterState.HealthCur = hp
		s.MonsterState.Deadline = 0
		s.Main.Words = append(s.Main.Words, [2]uint32{1336, b(.5)})
		add(fmt.Sprint("health-retreat-", hp), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			a := r.Combat.Actions
			want := uint32(1)
			if hp <= 50 {
				want = 6
			}
			if a[len(a)-1].Action != want {
				t.Fatalf("health head %d want %d", a[len(a)-1].Action, want)
			}
		})
	}

	for _, class := range []uint32{4, 0x2000, 0x10000, 0x12000} {
		for _, typ := range []uint32{1, 10, 11} {
			for _, sub := range []uint32{0, 0x200, 0x400, 0x12345678} {
				s := mainBase(1)
				s.Combat.TargetClass = class
				s.Main.TargetType = typ
				s.MonsterState.Subclass = sub
				clearDanger := false
				var value uint32
				if class&0x2000 != 0 {
					value = (sub >> 10) & 1
					clearDanger = value == 0
				} else if typ == 10 || typ == 11 {
					value = sub
					clearDanger = sub&0x200 == 0
				} else {
					value = typ
					clearDanger = class&0x10000 != 0
				}
				wantRet := uint32(int32(int16(value)))
				wantDanger := uint32(0x12345678)
				if clearDanger {
					wantDanger = 0
				}
				add(fmt.Sprintf("danger-%x-%d-%x", class, typ, sub), s, func(t *testing.T, r legacy.PortTestRoamResult) {
					if r.Main.Return != wantRet || r.Main.Globals[6] != wantDanger || r.Main.Globals[0] != 10 || r.Main.Globals[1] != 11 {
						t.Fatalf("danger state %s", mustMainJSON(r.Main))
					}
				})
			}
		}
	}
	for _, x := range []float32{119, 120, 121} {
		s := mainBase(5)
		s.Combat.Shield = true
		s.Main.Missile = true
		s.MonsterState.ObjFlags = 0x80
		s.Combat.Target = [2]uint32{b(140), b(x)}
		s.Main.TargetNewPos = s.Combat.Target
		add(fmt.Sprint("shield-lateral-", x), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := uint32(0)
			if x < 120 {
				want = 100
			}
			if r.Main.Return != want {
				t.Fatalf("shield target %d want %d", r.Main.Return, want)
			}
		})
	}
	for _, shield := range []uint32{0, 4} {
		for _, weapon := range []uint32{0, 0x400} {
			s := mainBase(3)
			s.Combat.Shield = true
			s.Main.Missile = true
			s.MonsterState.ObjFlags = 0x80
			s.MonsterState.Status = shield
			s.MonsterState.Subclass = 0x10
			s.MonsterState.PlayerWeapon = weapon
			if shield != 0 {
				s.MonsterState.PlayerShield = 0x1000000
			}
			add(fmt.Sprintf("defense-%x-%x", shield, weapon), s, func(t *testing.T, r legacy.PortTestRoamResult) {
				a := r.Combat.Actions
				want := uint32(1)
				if weapon == 0 && shield != 0 {
					want = 21
				}
				if a[len(a)-1].Action != want {
					t.Fatalf("defense head %d want %d", a[len(a)-1].Action, want)
				}
			})
		}
	}
	for _, host := range []bool{false, true} {
		for _, cursor := range []uint32{0, 1} {
			for _, dx := range []int32{0, 9, 10} {
				s := mainBase(3)
				s.Lifecycle.GameFlags = 2048
				s.Main.Flags2 = 16
				s.Main.Host = host
				s.Main.Cursor = cursor
				s.Main.CursorObject = true
				s.Main.Mouse = [2]int32{100 + dx, 100}
				add(fmt.Sprintf("cursor-%t-%d-%d", host, cursor, dx), s, func(t *testing.T, r legacy.PortTestRoamResult) {
					a := r.Combat.Actions
					want := uint32(1)
					if host && cursor == 0 && dx == 0 {
						want = 26
					}
					if a[len(a)-1].Action != want {
						t.Fatalf("cursor head %d want %d", a[len(a)-1].Action, want)
					}
				})
			}
		}
	}
	for _, dist := range []float32{49, 50, 51} {
		s := mainBase(3)
		s.Owner.Enemy = true
		s.MonsterState.Status = 0x20
		s.MonsterState.Deadline = 0
		s.Combat.Target = [2]uint32{b(100 + dist), b(100)}
		s.Main.Words = append(s.Main.Words, [2]uint32{1356, b(100)}, [2]uint32{1504, 1}, [2]uint32{1480, 2 | (5 << 16)}, [2]uint32{1484, 128})
		add(fmt.Sprint("blink-half-distance-", dist), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			a := r.Combat.Actions
			want := uint32(24)
			if dist < 50 {
				want = 1
			}
			if a[len(a)-1].Action != want {
				t.Fatalf("blink head %d want %d", a[len(a)-1].Action, want)
			}
		})
	}

	for _, started := range []bool{false, true} {
		for _, act := range []uint32{16, 17, 21, 23} {
			s := mainBase(3)
			s.Main.Started = started
			s.MonsterState.Action = act
			s.Owner.Buffs = 1 << 3
			add(fmt.Sprintf("cancel-push-%t-%d", started, act), s, func(t *testing.T, r legacy.PortTestRoamResult) {
				a := r.Combat.Actions
				if a[len(a)-1].Action != 36 {
					t.Fatal("missing confusion")
				}
			})
		}
	}
	for _, dx := range []float32{0, 15, 16} {
		s := mainBase(3)
		s.MonsterState.Action = 7
		s.Main.Words = append(s.Main.Words, [2]uint32{496, 0}, [2]uint32{500, b(100 + dx)})
		add(fmt.Sprint("movement-progress-", dx), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if (r.Combat.Status&0x200000 != 0) != (dx <= 15) {
				t.Fatal("movement frustration threshold")
			}
		})
	}
	{
		s := mainBase(3)
		s.MonsterState.Action = 7
		s.MonsterState.Pos[0] = b(-1e-7)
		s.Main.Words = append(s.Main.Words, [2]uint32{496, 0}, [2]uint32{500, b(15)})
		add("movement-double-delta", s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Combat.Status&0x200000 != 0 {
				t.Fatal("delta rounded before square")
			}
			if float32(float64(float32(15))-float64(float32(-1e-7))) != 15 {
				t.Fatal("fixture no longer distinguishes rounding")
			}
		})
	}
	for _, idx := range []int{21, 22, 23} {
		s := mainBase(3)
		s.MonsterState.Stack = int8(idx)
		s.Owner.Buffs = 1 << 11
		add(fmt.Sprint("fear-capacity-", idx), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Stack != 23 || len(r.Combat.Sounds) != 1 || r.Combat.Sounds[0] != 312 {
				t.Fatal("capacity changed prior effects")
			}
		})
	}
	for _, near := range []float32{10, 40, 60} {
		s := mainBase(5)
		s.Combat.Shield = true
		s.Main.Missile = true
		s.Main.SecondMissile = true
		s.MonsterState.ObjFlags = 0x80
		s.Main.SecondPos = [2]uint32{b(100 + near), b(100)}
		s.Main.SecondVelocity = [2]uint32{b(-2), 0}
		add(fmt.Sprint("shield-nearest-", near), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			want := uint32(100)
			if near <= 40 {
				want = 102
			}
			if r.Main.Return != want {
				t.Fatalf("nearest %d want %d", r.Main.Return, want)
			}
		})
	}

	for i, vel := range [][2]uint32{{3224937172, 3238389944}, {3249429723, 3262966429}} {
		s := mainBase(6)
		s.Combat.Shield = true
		s.Main.Missile = true
		s.MonsterState.ObjFlags = 0x80
		s.Combat.Target = [2]uint32{b(140), b(119)}
		s.Main.TargetNewPos = s.Combat.Target
		s.Combat.Velocity = vel
		add(fmt.Sprint("shield-distance-spill-", i), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Main.Globals[8] != 100 {
				t.Fatal("shield distance rounded before angular test")
			}
			vx, vy := float64(math.Float32frombits(vel[0])), float64(math.Float32frombits(vel[1]))
			speed := math.Sqrt(vx*vx + vy*vy)
			dist := math.Sqrt(40*40 + 19*19)
			nx := float64(float32(vx / speed))
			full := vy/speed*(-19/dist) + nx*(-40/dist)
			early := vy/speed*(-19/float64(float32(dist))) + nx*(-40/float64(float32(dist)))
			if !(full > .69999999) || early > .69999999 {
				t.Fatal("fixture must discriminate early float32 rounding")
			}
		})
	}

	{
		s := mainBase(3)
		s.Main.DefFlags = 8
		s.Lifecycle.GameFlags = 2048
		s.Combat.Shield = true
		s.Main.Missile = true
		s.MonsterState.ObjFlags = 0x84
		add("quest-defensive-dodge", s, func(t *testing.T, r legacy.PortTestRoamResult) {
			a := r.Combat.Actions
			if a[len(a)-1].Action != 9 || r.Logic != 3 {
				t.Fatal("quest dodge path")
			}
		})
	}
	for _, morphed := range []bool{false, true} {
		s := mainBase(3)
		s.Combat.Scan = true
		s.Combat.TargetFlags = 4
		s.Lifecycle.Eligible = true
		s.Lifecycle.Use = true
		s.Main.TargetSubclass = 0x10
		if morphed {
			s.Combat.TargetClass = 0x1000000
			s.Combat.PlayerUpdate = true
			s.MonsterState.Status = 0x20000
			s.MonsterState.Subclass = 0x10
		} else {
			s.Combat.TargetClass = 0x10
			s.MonsterState.HealthCur = 99
		}
		add(fmt.Sprint("inventory-", morphed), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			found := false
			for i, v := range r.Trace {
				if v == 30 && i+4 < len(r.Trace) && r.Trace[i+1] == 1 && r.Trace[i+2] == 1 {
					found = true
				}
			}
			if !found {
				t.Fatal("missing inventory placement")
			}
			if !morphed && (len(r.Lifecycle.Calls) == 0 || r.Lifecycle.Calls[0] != 3) {
				t.Fatal("missing use after placement")
			}
		})
	}
	for _, enemy := range []bool{false, true} {
		s := mainBase(3)
		s.Main.Actions = []uint32{4, 1}
		s.Main.TargetPlayer = true
		s.Combat.Scan = true
		s.Combat.TargetFlags = 4
		s.Combat.TargetClass = 4
		s.Owner.Enemy = enemy
		add(fmt.Sprint("guard-aggro-", enemy), s, func(t *testing.T, r legacy.PortTestRoamResult) {
			if r.Combat.Status&0x200 == 0 {
				t.Fatal("missing guard aggro status")
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
	mainHash(t, "contracts", got, "8a46dbb2e63ae4253721a2358dd5a5f879f146fd1c9d86922d8614c9698c831d")
}
