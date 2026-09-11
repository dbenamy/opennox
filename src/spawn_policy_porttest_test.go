//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
)

func spawnPolicyBase(n int) legacy.PortTestRoamSpec {
	s := generatorBase(9)
	p := &legacy.PortTestGeneratorSpawnPolicySpec{Point: [2]uint32{math.Float32bits(100), math.Float32bits(100)}}
	for i := 0; i < n; i++ {
		p.Records = append(p.Records, legacy.PortTestGeneratorSpawnPolicyRecord{Class: 2, Flags: 4, TypeInd: 25, Pos: [2]uint32{math.Float32bits(float32(100 + i*20)), math.Float32bits(100)}})
	}
	s.Callbacks.Generator.SpawnPolicy = p
	return s
}
func spawnPlayer(x float32, joined, busy uint32) legacy.PortTestSpawnPlayer {
	return legacy.PortTestSpawnPlayer{Pos: [2]uint32{math.Float32bits(x), math.Float32bits(100)}, View: [2]uint16{200, 200}, Joined: joined, Busy: busy, Flags: 4}
}
func spawnPolicyHash(t *testing.T, label string, specs []legacy.PortTestRoamSpec) []legacy.PortTestRoamResult {
	t.Helper()
	r := legacy.PortTestRoam(specs)
	for i, v := range r {
		for j, a := range v.Callbacks.Generator.Objects.SpawnPolicy {
			if !a.GuardsOK {
				t.Fatalf("case %d action %d guards", i, j)
			}
		}
	}
	callbackHash(t, label, r, spawnPolicyHashes[label])
	return r
}
func TestSpawnPolicyRegistration(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, n := range []int{1, 2, 3, 8} {
		for _, reverse := range []bool{false, true} {
			s := spawnPolicyBase(n)
			p := s.Callbacks.Generator.SpawnPolicy
			for i := 0; i < n; i++ {
				p.Actions = append(p.Actions, legacy.PortTestGeneratorSpawnPolicyAction{Op: legacy.PortTestGeneratorSpawnRegister, Record: i})
			}
			// Registering an associated child is idempotent even when the pool is full.
			p.Actions = append(p.Actions, legacy.PortTestGeneratorSpawnPolicyAction{Op: legacy.PortTestGeneratorSpawnRegister, Record: 0})
			for i := 0; i < n; i++ {
				j := i
				if reverse {
					j = n - 1 - i
				}
				p.Actions = append(p.Actions, legacy.PortTestGeneratorSpawnPolicyAction{Op: legacy.PortTestGeneratorSpawnRemove, Record: j})
			}
			specs = append(specs, s)
		}
	}
	r := spawnPolicyHash(t, "spawn-registration", specs)
	for i, v := range r {
		n := len(specs[i].Callbacks.Generator.SpawnPolicy.Records)
		actions := v.Callbacks.Generator.Objects.SpawnPolicy
		for j := 0; j < n; j++ {
			a := actions[j]
			want := min(j+1, 96)
			if int(a.GeneratorCount) != want || len(a.Spawn) != want || (a.Return == 1) != (j < 96) {
				t.Fatalf("case %d register%d count=%d nodes=%d return=%d", i, j, a.GeneratorCount, len(a.Spawn), a.Return)
			}
		}
		last := actions[len(actions)-1]
		if last.GeneratorCount != 0 || len(last.Spawn) != 0 {
			t.Fatalf("case %d final association not cleared", i)
		}
	}
}
func TestSpawnPolicyFarCull(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, distance := range []float32{0, 699, 700, 701, 1500} {
		for _, joined := range []uint32{0, 1, 2} {
			for _, busy := range []uint32{0, 1} {
				for _, players := range []int{0, 1, 2, 3} {
					s := spawnPolicyBase(3)
					p := s.Callbacks.Generator.SpawnPolicy
					p.Registers = []int{0, 1, 2}
					for i := 0; i < players; i++ {
						p.Players = append(p.Players, spawnPlayer(100+distance+float32(i*10), joined, busy))
					}
					p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnFarCull}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := spawnPolicyHash(t, "spawn-far-cull", specs)
	for i, v := range r {
		p := specs[i].Callbacks.Generator.SpawnPolicy
		a := v.Callbacks.Generator.Objects.SpawnPolicy[0]
		want := 0
		for _, m := range p.Records {
			near := false
			for _, pl := range p.Players {
				if pl.Joined == 1 && pl.Busy == 0 && math.Abs(float64(math.Float32frombits(pl.Pos[0])-math.Float32frombits(m.Pos[0]))) < 700 {
					near = true
				}
			}
			if !near {
				want++
			}
		}
		if len(a.Calls) != want*2 {
			t.Fatalf("case%d deletions %v want%d", i, a.Calls, want)
		}
	}
}
func TestSpawnPolicyVisibleCull(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, n := range []int{0, 1, 2, 4, 8} {
		for cap := 0; cap < 5; cap++ {
			for players := 0; players < 4; players++ {
				for layout := 0; layout < 4; layout++ {
					s := spawnPolicyBase(n)
					s.Callbacks.Generator.Balance["MaxOnscreenMonsterCount"] = float64(cap)
					p := s.Callbacks.Generator.SpawnPolicy
					for i := range p.Records {
						p.Registers = append(p.Registers, i)
						x := float32(110 + ((i*3+layout)%max(n, 1))*25)
						if layout == 3 {
							x += 1000
						}
						p.Records[i].Pos[0] = math.Float32bits(x)
					}
					for i := 0; i < players; i++ {
						p.Players = append(p.Players, spawnPlayer(float32(100+i*40), 1, 0))
					}
					p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnVisibleCull}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := spawnPolicyHash(t, "spawn-visible-cull", specs)
	for i, v := range r {
		a := v.Callbacks.Generator.Objects.SpawnPolicy[0]
		if len(a.MonsterList) != 0 || a.MonsterListCount != 0 || a.Return != 0 {
			t.Fatalf("case%d transient list not cleared", i)
		}
		cap := uint32(specs[i].Callbacks.Generator.Balance["MaxOnscreenMonsterCount"])
		p := specs[i].Callbacks.Generator.SpawnPolicy
		wantDeletes := 0
		if len(p.Players) > 0 && len(p.Records) > 0 && math.Float32frombits(p.Records[0].Pos[0]) < 1000 {
			wantDeletes = max(0, len(p.Records)-int(cap))
		}
		if len(a.Calls) != wantDeletes*2 {
			t.Fatalf("case%d visible deletes %v want%d", i, a.Calls, wantDeletes)
		}

		for _, n := range a.Counters {
			if n > cap {
				t.Fatalf("case%d visible count%d exceeds%d", i, n, cap)
			}
		}
	}
}
func TestSpawnPolicyTick(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, fps := range []uint32{1, 20, 30, 60} {
		for _, frame := range []uint32{0, 1, 14, 15, 16, 99, 100, 149, 150, 151, 299, 300, 0xffffffff} {
			s := spawnPolicyBase(3)
			s.Owner.FPS = fps
			s.MonsterState.FPS = fps
			s.Callbacks.Generator.Frame = frame
			p := s.Callbacks.Generator.SpawnPolicy
			p.Registers = []int{0, 1, 2}
			p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnTick}}
			specs = append(specs, s)
		}
	}
	r := spawnPolicyHash(t, "spawn-tick", specs)
	for i, v := range r {
		s := specs[i]
		a := v.Callbacks.Generator.Objects.SpawnPolicy[0]
		want := s.Callbacks.Generator.Frame / 15
		if s.Callbacks.Generator.Frame%15 == 0 {
			want = 0
		}
		if a.Return != want {
			t.Fatalf("case%d return%d want%d", i, a.Return, want)
		}
		deletes := 0
		if s.Callbacks.Generator.Frame%(5*s.Owner.FPS) == 0 {
			deletes = 3
		}
		if len(a.Calls) != 2*deletes {
			t.Fatalf("case%d tick deletion count", i)
		}
	}
}

func TestSpawnPolicyExhaustion(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, n := range []int{95, 96, 97, 100} {
		s := spawnPolicyBase(n)
		p := s.Callbacks.Generator.SpawnPolicy
		for i := 0; i < n-1; i++ {
			p.Registers = append(p.Registers, i)
		}
		p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{
			{Op: legacy.PortTestGeneratorSpawnRegister, Record: n - 1},
			{Op: legacy.PortTestGeneratorSpawnRegister, Record: 0},
			{Op: legacy.PortTestGeneratorSpawnRemove, Record: 0},
			{Op: legacy.PortTestGeneratorSpawnRegister, Record: n - 1},
			{Op: legacy.PortTestGeneratorSpawnRemove, Record: n - 1},
			{Op: legacy.PortTestGeneratorSpawnRemove, Record: n - 1},
		}
		specs = append(specs, s)
	}
	r := spawnPolicyHash(t, "spawn-exhaustion", specs)
	for i, v := range r {
		n := len(specs[i].Callbacks.Generator.SpawnPolicy.Records)
		a := v.Callbacks.Generator.Objects.SpawnPolicy
		if (a[0].Return == 1) != (n <= 96) || int(a[0].GeneratorCount) != min(n, 96) || a[1].GeneratorCount != a[0].GeneratorCount || a[3].Return != 1 || a[5].GeneratorCount != a[4].GeneratorCount {
			t.Fatalf("case%d pool capacity/reuse contract", i)
		}
	}
}

func TestSpawnPolicyZombieCleanup(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, typ := range []uint16{2, 3, 24, 25} {
		for _, op := range []legacy.PortTestGeneratorSpawnPolicyOp{legacy.PortTestGeneratorSpawnNonZombieCleanup, legacy.PortTestGeneratorSpawnGlyphCleanup} {
			s := spawnPolicyBase(1)
			p := s.Callbacks.Generator.SpawnPolicy
			p.Records[0].TypeInd = typ
			p.Registers = []int{0}
			p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: op}, {Op: legacy.PortTestGeneratorSpawnRemove}}
			specs = append(specs, s)
		}
	}
	r := spawnPolicyHash(t, "spawn-zombie-cleanup", specs)
	for i, v := range r {
		p := specs[i].Callbacks.Generator.SpawnPolicy
		a := v.Callbacks.Generator.Objects.SpawnPolicy
		want := byte(0)
		if p.Actions[0].Op == legacy.PortTestGeneratorSpawnNonZombieCleanup && p.Records[0].TypeInd <= 3 {
			want = 1
		}
		if a[0].GeneratorCount != want || a[1].GeneratorCount != 0 {
			t.Fatalf("case%d zombie cleanup count", i)
		}
	}
}
func TestSpawnPolicyGlyph(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, subclass := range []uint32{0, 0x2000, 0x2010} {
		for mask := 0; mask < 8; mask++ {
			s := spawnPolicyBase(1)
			p := s.Callbacks.Generator.SpawnPolicy
			p.Records[0].Subclass = subclass
			p.Records[0].UpdateWords = map[uintptr]uint32{}
			for i := 0; i < 3; i++ {
				if mask&(1<<i) != 0 {
					p.Records[0].UpdateWords[2044+4*uintptr(i)] = uint32(i + 1)
				}
			}
			p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnRegister}, {Op: legacy.PortTestGeneratorSpawnGlyphCleanup}, {Op: legacy.PortTestGeneratorSpawnGlyphCleanup}}
			specs = append(specs, s)
		}
	}
	r := spawnPolicyHash(t, "spawn-glyph", specs)
	for i, v := range r {
		p := specs[i].Callbacks.Generator.SpawnPolicy
		a := v.Callbacks.Generator.Objects.SpawnPolicy
		want := 0
		if p.Records[0].Subclass&0x2000 != 0 {
			want = 1
		}
		if len(a[0].Glyphs) != want || len(a[1].Calls) != want*2 || len(a[2].Calls) != 0 || a[1].GeneratorCount != 0 {
			t.Fatalf("case%d glyph allocation/cleanup", i)
		}
		if want == 1 {
			for j := 0; j < 3; j++ {
				if a[0].GlyphData[0][j] != p.Records[0].UpdateWords[2044+4*uintptr(j)] {
					t.Fatalf("case%d glyph spell%d", i, j)
				}
			}
		}
	}
}
func TestSpawnPolicyCandidate(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, class := range []uint32{0, 2, 4, math.Float32bits(1), math.Float32bits(2), math.Float32bits(3), math.Float32bits(258)} {
		for _, flags := range []uint32{4, 0x20, math.Float32bits(32), math.Float32bits(32768), math.Float32bits(32800)} {
			for _, typ := range []uint16{2, 3, 25} {
				for wall := 0; wall < 3; wall++ {
					s := spawnPolicyBase(1)
					s.Combat.Wall = wall
					p := s.Callbacks.Generator.SpawnPolicy
					p.Players = []legacy.PortTestSpawnPlayer{spawnPlayer(150, 1, 0)}
					p.Records[0].Class = class
					p.Records[0].Flags = flags
					p.Records[0].TypeInd = typ
					p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnCandidate}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := spawnPolicyHash(t, "spawn-candidate", specs)
	for i, v := range r {
		if specs[i].Combat.Wall != 0 {
			continue
		}
		rec := specs[i].Callbacks.Generator.SpawnPolicy.Records[0]
		class := uint8(math.Float32frombits(rec.Class))
		flags := uint32(math.Float32frombits(rec.Flags))
		want := uint32(0)
		if class&2 != 0 && flags&32 == 0 && (flags&0x8000 == 0 || rec.TypeInd == 2 || rec.TypeInd == 3) {
			want = 1
		}
		if got := v.Callbacks.Generator.Objects.SpawnPolicy[0].Occupied; got != want {
			t.Fatalf("case%d candidate count%d want%d", i, got, want)
		}
	}
}
func TestSpawnPolicyAdmission(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for players := 0; players < 4; players++ {
		for _, joined := range []uint32{0, 1, 2} {
			for _, cap := range []float64{0, 1, 2, 2.9} {
				for _, x := range []float32{-201, -200, 100, 400, 401, 1000} {
					s := spawnPolicyBase(0)
					s.Callbacks.Generator.Balance["MaxOnscreenMonsterCount"] = cap
					p := s.Callbacks.Generator.SpawnPolicy
					for i := 0; i < players; i++ {
						p.Players = append(p.Players, spawnPlayer(float32(100+i*40), joined, 0))
					}
					p.Point[0] = math.Float32bits(x)
					p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnAdmission}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := spawnPolicyHash(t, "spawn-admission", specs)
	for i, v := range r {
		p := specs[i].Callbacks.Generator.SpawnPolicy
		want := uint32(1)
		if specs[i].Callbacks.Generator.Balance["MaxOnscreenMonsterCount"] == 0 {
			x := math.Float32frombits(p.Point[0])
			for _, pl := range p.Players {
				cx := math.Float32frombits(pl.Pos[0])
				if pl.Joined != 0 && x >= cx-300 && x <= cx+300 {
					want = 0
				}
			}
		}
		if got := v.Callbacks.Generator.Objects.SpawnPolicy[0].Return; got != want {
			t.Fatalf("case%d admission%d want%d", i, got, want)
		}
	}
}

func TestSpawnPolicyIndexedAdmission(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for n := 0; n < 5; n++ {
		for cap := 0; cap < 6; cap++ {
			for wall := 0; wall < 3; wall++ {
				for _, numeric := range []bool{false, true} {
					s := spawnPolicyBase(n)
					s.Combat.Wall = wall
					s.Callbacks.Generator.Balance["MaxOnscreenMonsterCount"] = float64(cap)
					p := s.Callbacks.Generator.SpawnPolicy
					p.Players = []legacy.PortTestSpawnPlayer{spawnPlayer(150, 1, 0)}
					for i := range p.Records {
						p.Records[i].Indexed = true
						if numeric {
							p.Records[i].Class = math.Float32bits(2)
						}
					}
					p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnAdmission}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := spawnPolicyHash(t, "spawn-indexed-admission", specs)
	for i, v := range r {
		s := specs[i]
		if s.Combat.Wall != 0 {
			continue
		}
		p := s.Callbacks.Generator.SpawnPolicy
		a := v.Callbacks.Generator.Objects.SpawnPolicy[0]
		count := 0
		if len(p.Records) > 0 && p.Records[0].Class == math.Float32bits(2) {
			count = len(p.Records)
		}
		want := uint32(0)
		if count < int(s.Callbacks.Generator.Balance["MaxOnscreenMonsterCount"]) {
			want = 1
		}
		if a.Occupied != uint32(count) || a.Return != want {
			t.Fatalf("case%d occupied%d return%d want%d/%d", i, a.Occupied, a.Return, count, want)
		}
	}
}

func TestSpawnPolicyMonsterPool(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	for _, reserved := range []uint8{0, 94, 95, 96} {
		for cap := 0; cap < 4; cap++ {
			for players := 1; players < 4; players++ {
				s := spawnPolicyBase(4)
				p := s.Callbacks.Generator.SpawnPolicy
				p.ReserveMonster = reserved
				p.Registers = []int{0, 1, 2, 3}
				s.Callbacks.Generator.Balance["MaxOnscreenMonsterCount"] = float64(cap)
				for i := 0; i < players; i++ {
					p.Players = append(p.Players, spawnPlayer(float32(90+i*15), 1, 0))
				}
				p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnVisibleCull}, {Op: legacy.PortTestGeneratorSpawnVisibleCull}}
				specs = append(specs, s)
			}
		}
	}
	r := spawnPolicyHash(t, "spawn-monster-pool", specs)
	for i, v := range r {
		p := specs[i].Callbacks.Generator.SpawnPolicy
		a := v.Callbacks.Generator.Objects.SpawnPolicy
		for _, step := range a {
			if step.MonsterListCount != 0 || step.MonsterListHead != 0 {
				t.Fatalf("case%d monster-list cleanup", i)
			}
			if p.ReserveMonster == 96 && len(step.Calls) != 0 {
				t.Fatalf("case%d exhausted list deleted monsters", i)
			}
		}
	}
}

var spawnPolicyHashes = map[string]string{
	"spawn-visibility-groups": "849396b511571cc3c48ec1369cd5636e7aaf7674c3c55f6db61fe32bab2808bf",
	"spawn-admission":         "3a5309b60097403c0235f99ccbfb39022c2cba77a427d5504bbe7d138b5b07ef",
	"spawn-candidate":         "74ab41874a9e19e9093f93b5b192d992749e7014e6feef83ad623d4d02cbbd35",
	"spawn-exhaustion":        "de7e0f7d06a70204832a417bc2d11abfd45b143470a385384c17a76378486a9a",
	"spawn-far-cull":          "d4efecc7a6083a71943f9752252e6c9320378999ee11a8b861ba568c2b5436d5",
	"spawn-glyph":             "6b7d4b8e7b324fd69ee40fcf9b879c8485e9a032e91b7cd6169921844526ac7a",
	"spawn-indexed-admission": "39309f1c57c27d93de710e133b84912d049d433bb3b2eae3c0abe681fc7a10fb",
	"spawn-monster-pool":      "86f88754916c402762a15ec7de4d08366c323147b37e47d020143030c9a8d206",
	"spawn-registration":      "d98e0b1677ccecdec78f6d2497cd080e2d0345a9df3e9ab0d35397e5e692b1eb",
	"spawn-tick":              "204871ee4d807f63b7cbd09d2d96f99cc0484787097070005ec84b56142a153b",
	"spawn-visible-cull":      "52c40a9e5b8e204922e10dd7aa5ed941b747ea44f9d561bb84293414d6429a5f",
	"spawn-zombie-cleanup":    "62c44e9715634a58cdf112ab46058b66c408d5b968abe85026a5218ef0d53a3e",
}

func TestSpawnPolicyVisibilityGroups(t *testing.T) {
	var specs []legacy.PortTestRoamSpec
	positions := []float32{0, 100, 250, 400, 500, 750, 900, 1000, 1200}
	for _, n := range []int{3, 6, 9} {
		for cap := 0; cap < 5; cap++ {
			for players := 1; players < 4; players++ {
				for shift := 0; shift < 3; shift++ {
					s := spawnPolicyBase(n)
					p := s.Callbacks.Generator.SpawnPolicy
					s.Callbacks.Generator.Balance["MaxOnscreenMonsterCount"] = float64(cap)
					for i := range p.Records {
						p.Records[i].Pos[0] = math.Float32bits(positions[(i+shift)%len(positions)])
						p.Registers = append(p.Registers, i)
					}
					for i := 0; i < players; i++ {
						p.Players = append(p.Players, spawnPlayer(float32(100+400*i), 1, 0))
					}
					p.Actions = []legacy.PortTestGeneratorSpawnPolicyAction{{Op: legacy.PortTestGeneratorSpawnVisibleCull}}
					specs = append(specs, s)
				}
			}
		}
	}
	r := spawnPolicyHash(t, "spawn-visibility-groups", specs)
	for i, v := range r {
		limit := uint32(specs[i].Callbacks.Generator.Balance["MaxOnscreenMonsterCount"])
		a := v.Callbacks.Generator.Objects.SpawnPolicy[0]
		for _, n := range a.Counters {
			if n > limit {
				t.Fatalf("case%d shared visibility count%d exceeds%d", i, n, limit)
			}
		}
	}
}
