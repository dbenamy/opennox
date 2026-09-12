//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

func spellLifeBase(op int) legacy.PortTestRoamSpec {
	s := controlsBase(0)
	p := s.Callbacks.Shop
	o := p.TemporaryUpdates.World.Objectives
	a := o.Attack
	a.Controls.SpellLifecycle = &legacy.PortTestSpellLifecycleSpec{Record: legacy.PortTestSpellLifeRecord{Words: map[int]uint32{}, Refs: map[int]int{}}}
	p.Sequence = []legacy.PortTestShopAction{{Op: 1500 + op}}
	s.Lifecycle.GameFlags |= 1
	p.Resources.PlayerClass = 1
	p.Resources.Mana = 200
	p.Resources.OldMana = 200
	p.Resources.MaxMana = 200
	a.Controls.X = 1
	o.ObjectList = []int{100, 101, 102, 3, 4, 5}
	o.PlayerDataWords[0] = map[int]uint32{3700: 3}
	p.EffectsUse.Balance["ShockDamage"] = []float64{1, 2.5, 10, 20, 30}
	p.EffectsUse.Balance["MagicMissileCount"] = []float64{1, 2, 3, 4, 5}
	p.EffectsUse.Balance["PixieCount"] = []float64{1, 2, 3, 4, 5}
	p.EffectsUse.Balance["MaxTrapCount"] = []float64{10}
	p.EffectsUse.Balance["MaxBomberCount"] = []float64{10}
	p.EffectsUse.Balance["SpellCastDelay"] = []float64{1}
	return s
}
func TestSpellLifecycleSmoke(t *testing.T) {
	for op := 0; op < 29; op++ {
		t.Run(fmt.Sprintf("%02d", op), func(t *testing.T) {
			s := spellLifeBase(op)
			a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
			sp := a.Controls.SpellLifecycle
			switch op {
			case 0, 12:
				a.Controls.X = 0
			case 1, 3:
				a.Controls.X = 0
				a.Controls.Y = 0
			case 4:
				sp.Record.Words[0] = 1
			case 5:
				a.Controls.Y = 1
			case 13:
				sp.Record.Words[0] = 1
			case 14:
				a.Controls.X = int32(math.Float32bits(0.5))
			case 20:
				sp.Record.Words[4] = 7
				sp.Record.Words[8] = 3
				sp.Record.Refs[16] = 1
				sp.Record.Refs[48] = 4
			case 22, 24, 25, 27:
				a.Controls.X = 3
				a.ActorWords = map[int]uint32{340: 1 << 3, 408: 0x04030201}
			case 23:
				a.Controls.X = 3
				a.Controls.Y = 60
				sp.Z = 4
			case 28:
				a.ActorWords = map[int]uint32{340: 1 << 3, 348: 1 << 16}
			}
			spellLifeHash(t, fmt.Sprintf("spell-life-smoke-%02d", op), []legacy.PortTestRoamSpec{s})
		})
	}
}

func TestSpellLifecycleBuffAccessors(t *testing.T) {
	for _, op := range []int{22, 24, 25} {
		var cases []legacy.PortTestRoamSpec
		for _, slot := range []int32{-128, -1, 0, 1, 3, 16, 21, 22, 29, 30, 31, 32, 127} {
			for mode := 0; mode < 4; mode++ {
				s := spellLifeBase(op)
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
				a.Controls.X = slot
				a.ActorWords = map[int]uint32{340: 0xaaaaaaaa}
				for i := 0; i < 16; i++ {
					a.ActorWords[344+4*i] = uint32(65535-i) | uint32(i*100+1)<<16
				}
				for i := 0; i < 8; i++ {
					a.ActorWords[408+4*i] = 0xff804001 + uint32(i)
				}
				if mode&1 != 0 {
					a.ActorWords[340] = 0xffffffff
				}
				if mode&2 != 0 {
					a.Actor = 0
				}
				cases = append(cases, s)
			}
		}
		spellLifeHash(t, fmt.Sprintf("spell-life-accessors-%02d", op), cases)
	}
}
func TestSpellLifecycleBuffMutation(t *testing.T) {
	for _, op := range []int{23, 26, 27, 28} {
		var cases []legacy.PortTestRoamSpec
		for slot := int32(0); slot < 32; slot++ {
			for _, timer := range []uint16{0, 1, 65535} {
				s := spellLifeBase(op)
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
				a.Controls.X = slot
				a.Controls.Y = int32(timer)
				a.Controls.SpellLifecycle.Z = 128 + slot
				a.ActorWords = map[int]uint32{340: 1 << uint(slot)}
				off := 344 + 4*int(slot/2)
				a.ActorWords[off] = uint32(timer) << uint(16*(slot%2))
				off = 408 + 4*int(slot/4)
				a.ActorWords[off] = uint32(slot+17) << uint(8*(slot%4))
				s.Callbacks.Shop.Sequence = append(s.Callbacks.Shop.Sequence, legacy.PortTestShopAction{Op: 1500 + op})
				cases = append(cases, s)
			}
		}
		spellLifeHash(t, fmt.Sprintf("spell-life-buff-mutation-%02d", op), cases)
	}
}

func spellLifeHash(t *testing.T, name string, cases []legacy.PortTestRoamSpec) []legacy.PortTestRoamResult {
	t.Helper()
	old := legacy.PortTestSpellLifePlayerSpell
	legacy.PortTestSpellLifePlayerSpell = func(u *server.Object) { owner := &Server{Server: legacy.GetServer().S()}; owner.PlayerSpell(u) }
	defer func() { legacy.PortTestSpellLifePlayerSpell = old }()
	r := controlsRun(t, cases)
	callbackHash(t, name, r, spellLifeHashes[name])
	return r
}
func TestSpellLifecycleBooks(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for count := 1; count <= 3; count++ {
		for mode := 0; mode < 8; mode++ {
			s := spellLifeBase(2)
			p := s.Callbacks.Shop
			o := p.TemporaryUpdates.World.Objectives
			a := o.Attack
			sp := a.Controls.SpellLifecycle
			for i := 0; i < count; i++ {
				r := legacy.PortTestSpellLifeRecord{Words: map[int]uint32{8: 1, 40: 0, 44: 0}, Refs: map[int]int{4: 100 + i}}
				switch mode {
				case 1:
					r.Words[40] = s.Owner.Frame + 1
				case 2:
					o.PlayerWords[i][16] = 0x20
				case 3:
					r.Words[28] = 1 << 8
				case 4:
					r.Words[12] = 1
				case 5:
					r.Words[48] = 1
				case 6:
					r.Words[28] = 1 << 8
					o.PlayerUpdateWords[i] = map[int]uint32{192: 1, 212: 1}
				case 7:
					r.Words[28] = 1 << 8
					p.Resources.Mana = 0
				}
				sp.Books = append(sp.Books, r)
			}
			p.Sequence = []legacy.PortTestShopAction{{Op: 1502}, {Op: 1502}, {Op: 1502}, {Op: 1502}}
			cases = append(cases, s)
		}
	}
	spellLifeHash(t, "spell-life-books", cases)
}
func TestSpellLifecycleDurationLists(t *testing.T) {
	for _, op := range []int{1, 3, 17, 18, 19, 21} {
		var cases []legacy.PortTestRoamSpec
		for n := 0; n < 4; n++ {
			for mode := 0; mode < 4; mode++ {
				s := spellLifeBase(op)
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				a.Controls.X = int32(mode & 1)
				if op == 21 {
					a.Controls.X = 24
				}
				if op == 18 {
					a.Controls.Target = 3
					s.Callbacks.Shop.Items[0].Class = 0x1000
					s.Callbacks.Shop.Items[0].Subclass = 0x4040000
				}
				for i := 0; i < n; i++ {
					id := uint32(24)
					if i == 1 {
						id = 43
					}
					if i == 2 {
						id = 59
					}
					refs := map[int]int{16: 1, 48: 1}
					if mode&2 != 0 {
						refs[16] = 101
						if i%2 == 0 {
							refs[48] = 2
						}
					}
					sp.Durations = append(sp.Durations, legacy.PortTestSpellLifeRecord{Words: map[int]uint32{4: id, 8: 3, 88: uint32(i & 1)}, Refs: refs})
				}
				s.Callbacks.Shop.Sequence = append(s.Callbacks.Shop.Sequence, legacy.PortTestShopAction{Op: 1500 + op})
				cases = append(cases, s)
			}
		}
		spellLifeHash(t, fmt.Sprintf("spell-life-durations-%02d", op), cases)
	}
}
func TestSpellLifecycleRayMessages(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, id := range []uint32{0, 7, 9, 22, 24, 35, 43, 59, 136} {
		for mode := 0; mode < 4; mode++ {
			s := spellLifeBase(20)
			sp := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
			sp.Record.Words = map[int]uint32{4: id, 8: uint32(127 + mode)}
			sp.Record.Refs = map[int]int{16: 1, 48: 4}
			if mode&1 != 0 {
				sp.Record.Refs[48] = 1
			}
			if mode&2 != 0 {
				sp.Record.Refs[48] = 0
			}
			if id == 43 {
				sp.Durations = []legacy.PortTestSpellLifeRecord{{Words: map[int]uint32{4: 7, 8: 3}, Refs: map[int]int{16: 1, 48: 4}}, {Words: map[int]uint32{4: 9, 8: 5}, Refs: map[int]int{16: 100, 48: 101}}}
				sp.RecordChildren = map[int]int{108: 0}
			}
			cases = append(cases, s)
		}
	}
	spellLifeHash(t, "spell-life-ray-messages", cases)
}

func TestSpellLifecycleMana(t *testing.T) {
	for _, op := range []int{4, 5, 6} {
		var cases []legacy.PortTestRoamSpec
		for _, mana := range []uint16{0, 1, 199, 200, 65535} {
			for _, cost := range []int32{-32768, -1, 0, 1, 200, 32767, 65536} {
				for mode := 0; mode < 6; mode++ {
					s := spellLifeBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					sp := a.Controls.SpellLifecycle
					p.Resources.Mana = mana
					p.Resources.OldMana = mana
					p.Resources.God = mode == 1
					if mode == 2 {
						p.Resources.Subject = 3
					}
					if mode == 3 {
						p.Resources.Subject = 2
					}
					sp.Definitions = []server.PortTestSpellLifecycleDef{{Index: 1, Valid: true, Enabled: true, Flags: 0x1000000, ManaCost: int(cost), Phonemes: []int{0}}}
					a.Controls.X = 1
					a.Controls.Y = 1
					if op == 4 {
						sp.Record.Words = map[int]uint32{0: 1, 4: 1, 8: 1}
						a.Controls.X = 3
						if mode == 4 {
							a.Controls.X = -1
						}
						if mode == 5 {
							sp.NullRecord = true
						}
					}
					if op == 5 {
						if mode == 4 {
							a.Controls.X = 0
						}
						if mode == 5 {
							a.Controls.Y = 2
						}
					}
					if op == 6 {
						a.Controls.X = cost
						sp.HalfPointerReturn = mode == 2 || mode == 3
					}
					cases = append(cases, s)
				}
			}
		}
		spellLifeHash(t, fmt.Sprintf("spell-life-mana-%02d", op), cases)
	}
}
func TestSpellLifecycleClassEligibility(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for class := byte(0); class < 3; class++ {
		for _, flags := range []uint32{0, 0x1000000, 0x2000000, 0x4000000, 0x6000000, 0x7000000} {
			for mode := 0; mode < 4; mode++ {
				s := spellLifeBase(7)
				p := s.Callbacks.Shop
				p.Resources.PlayerClass = class
				p.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle.Definitions = []server.PortTestSpellLifecycleDef{{Index: 1, Flags: flags, Valid: mode&1 == 0, Enabled: mode&2 == 0}}
				cases = append(cases, s)
			}
		}
	}
	spellLifeHash(t, "spell-life-class-eligibility", cases)
}
func TestSpellLifecyclePositiveContracts(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{22, 24, 25, 5, 4} {
		s := spellLifeBase(op)
		a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
		a.ActorWords = map[int]uint32{340: 1 << 3, 348: 1234 << 16, 408: 7 << 24}
		a.Controls.X = 3
		if op == 5 {
			a.Controls.X = 1
			a.Controls.Y = 1
		}
		if op == 4 {
			a.Controls.X = 1
			a.Controls.SpellLifecycle.Record.Words[0] = 1
		}
		cases = append(cases, s)
	}
	results := spellLifeHash(t, "spell-life-positive-contracts", cases)
	for i, want := range []uint32{1, 1234, 7, 10, 1} {
		if got := results[i].Callbacks.Shop.Sequence[0].Return; got != want {
			t.Errorf("case %d return=%d want %d", i, got, want)
		}
	}
}

var spellLifeHashes = map[string]string{
	"spell-life-accessors-22":       "47ec93279cbce22bd051c8191863f267304abdd9f29d822d89abb80a6854f283",
	"spell-life-accessors-24":       "85e38b94faa84e85df2886303f5f20c1aeb48cd4384ee4fc1742480438de2445",
	"spell-life-accessors-25":       "7b32672050ab95e8f4f905d1ca3dd9f8ea82df2a08216939bfaa53292d92f076",
	"spell-life-book-counter":       "dafa3b62fe20f5661dbafa56ef2310147ea304f077d98c9f3d8f369949c1af37",
	"spell-life-book-insertion":     "194b4c7372dfbd54d8ff26f504bc34a6469d26da7c1f9f9bb5d8f06f9e9c65ce",
	"spell-life-books":              "780ecc09a9d2b084d9924197646bce4a1a7108cf10a86331cf81f718ce44e64e",
	"spell-life-buff-exclusions":    "644d3f5737431ecfa7a19686455dec0f183f812c3a59738e0891bce7bc06f8b5",
	"spell-life-buff-mutation-23":   "5286a3d5a3d1ea0c9eaa808021201c1decfa42123ae79e14ed75c1d79325bd06",
	"spell-life-buff-mutation-26":   "65b73ad027056e8be05b8ac0507e59123a58cb82b51f71f9bd8e4da67cafeebb",
	"spell-life-buff-mutation-27":   "5744ec004f1b66f30cbec332d0e4e2cb0efb5ef501266503c458bc5cf1050638",
	"spell-life-buff-mutation-28":   "fa4457c6cc67d81378b2e5089d9db0cf8ae912bed1144831dd434738912302f3",
	"spell-life-class-eligibility":  "721d361f62758fa575d92e787a0aedd4cac62a1782edf1f792feb2ce92b69332",
	"spell-life-durations-01":       "faca2a0228ab89395a10d157533c91a11de3263505e73d1e0f3fdae786c7c998",
	"spell-life-durations-03":       "3768496e004b1bac5ed80bb69c00800cd086b47ce56795c3082bbf3462dfd2f9",
	"spell-life-durations-17":       "9a81ca073ec32b0b723f6a494572978e506b8d734df9c1bbbfd5b11a4b57e2be",
	"spell-life-durations-18":       "991f21736c5bdffff447ee4fd0ec062dbc8ae2aedae79774ff34e128baa07329",
	"spell-life-durations-19":       "9a81ca073ec32b0b723f6a494572978e506b8d734df9c1bbbfd5b11a4b57e2be",
	"spell-life-durations-21":       "1a0b719df4b5d46f3f679dde9a71cd9b804d8d68b862073b4bed31f14168b66e",
	"spell-life-mana-04":            "ff4f54ff15643c8c212a44099eb6cbd7a7a84f5a304a50b36a0ba565def4c73d",
	"spell-life-mana-05":            "1bc166d5cb581acec2f91c45ee6912399624b39b6024b34b25ac4d768c332293",
	"spell-life-mana-06":            "beefa28e7f651f4d77a1c0a68c1aa7330fbb17c337f448ffba52bb2625aaf81f",
	"spell-life-phonemes-00":        "8fdcc577166a3745577340258f72f4679bee7b6728d57ee6e0e5156b7f1a1660",
	"spell-life-phonemes-12":        "0571925cf84906958f807666678a1944b0da77cb2f6fc43adccf5e47ae604c5e",
	"spell-life-position":           "48f2721f344c5a21fcf027db7561274ec6f29d578bd71f176172e8d8265066aa",
	"spell-life-positive-contracts": "dd2f391efb0c3c7b7e8d402901a1dc33017f2807c471c8ba3a720c21788ea4d4",
	"spell-life-power":              "34d16c03cb69cbd9d11a9475dc387853d5bfafdceddbcae165150ae7af3f968d",
	"spell-life-projectiles":        "f172786b23da86066e129df7263fe6f8b99b92905599b36a4aedbc1d41f7a638",
	"spell-life-ray-messages":       "529006657f2224726113f313e6864a8813ef0ca2fc7cdb7d7994a04215e63b13",
	"spell-life-restrictions-08":    "6fd76f28232dc411319c4285e04943a04397916f501b0e14ea8c0070e9088823",
	"spell-life-restrictions-09":    "3439cba2cdeb33fffc0e8a24cf0aa46f643c329b2295dc856505e4be6743eea6",
	"spell-life-shock":              "88d3b1da52626625ff40581be07c6ff68cf6f20ebf46ce526b30a511858e82aa",
	"spell-life-smoke-00":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-01":           "2c36a2ec0d8a268b2d518e9ec0a7b0159a9960cf378fe5867fe31b50edfaae1b",
	"spell-life-smoke-02":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-03":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-04":           "1ee6873797094079e4ece96a2e6b67f569659aa0c4526e3b5e90ff74273b8796",
	"spell-life-smoke-05":           "67b993d5933ab30dfb097e98e306e2d49e65a4f53fcf7b343b2bc9825be0b88f",
	"spell-life-smoke-06":           "1a7cf28e6385c53d0af181a0b20aef440e47d003944eeda86a45da64c5ed1267",
	"spell-life-smoke-07":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-08":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-09":           "24e8bfdd0d89e8908a7930a1a275cebd727566e599ba90478af81c18765013e1",
	"spell-life-smoke-10":           "638276358f75e2f10545cb090f7e5ac1f6602d1c106ed397f95562c6661c9416",
	"spell-life-smoke-11":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-12":           "425844ec19cb241f5813cae0125e80e942df26d766ac09d53bb94bce5870915c",
	"spell-life-smoke-13":           "a487672bc7f86be9bc5858f077c7ff831eabec7d48b38865c01d462c624be05c",
	"spell-life-smoke-14":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-15":           "d13fdb0260df87af3145936b9896bf4575d6b00d5cf23eba0acac7ec70b5380b",
	"spell-life-smoke-16":           "2c36a2ec0d8a268b2d518e9ec0a7b0159a9960cf378fe5867fe31b50edfaae1b",
	"spell-life-smoke-17":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-18":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-19":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-20":           "dfeb264d4e57d144d254a8d76d9f14245d0ceabbade1719eaa150dd0c714a26b",
	"spell-life-smoke-21":           "9de1cd18f60d3e7f7d613054513efb9fda81e0923bd8671247c9f1f19ee63ec5",
	"spell-life-smoke-22":           "5222191367b84d130976091baa8355c798360c22f0fda43e0622974e40eba6d3",
	"spell-life-smoke-23":           "72b54b915bdaaed7c66506242a6378214a668d865b95689b1e2f06867f1f6dea",
	"spell-life-smoke-24":           "bb83e3db6c4260ac20bc11af956087bdaa0d4f035d77d5671c47c242e1af592c",
	"spell-life-smoke-25":           "9a54afe40289a6939927ec4876a9a5d3d1f54e6edae4e7c62b427e77eb160efa",
	"spell-life-smoke-26":           "8774f58f4900fbdb09ea8cbd30ac262cfc370f506b34a6b3aa98e69bda517d15",
	"spell-life-smoke-27":           "390085615c5481f804671e63ffeebe5580501134a84d6814b2ad5dcec0d96507",
	"spell-life-smoke-28":           "df99bff23a1ef656ef79728007effaf9513da44edba0b671e602a16e85a5b1f7",
}

func TestSpellLifecyclePhonemes(t *testing.T) {
	for _, op := range []int{0, 12} {
		var cases []legacy.PortTestRoamSpec
		for _, phon := range []int32{-128, -1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 127} {
			for mode := 0; mode < 8; mode++ {
				s := spellLifeBase(op)
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				a := o.Attack
				a.Controls.X = phon
				o.PlayerDataWords[0][2252] = uint32(mode & 1)
				if mode&2 != 0 {
					p.Resources.Subject = 3
					o.ObjectList = append(o.ObjectList, 1)
				}
				if mode&4 != 0 {
					s.Lifecycle.GameFlags &^= 1
					a.Controls.SpellLifecycle.ClientSprite = true
				}
				cases = append(cases, s)
			}
		}
		r := spellLifeHash(t, fmt.Sprintf("spell-life-phonemes-%02d", op), cases)
		if op == 12 {
			for mode, want := range []uint32{193, 201, 193, 193, 193, 201, 193, 193} {
				if got := r[16+mode].Callbacks.Shop.Sequence[0].Return; got != want {
					t.Errorf("phoneme zero mode %d got %d want %d", mode, got, want)
				}
			}
		}
	}
}
func TestSpellLifecycleBookInsertion(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for n := 1; n <= 5; n++ {
		for mode := 0; mode < 12; mode++ {
			s := spellLifeBase(13)
			p := s.Callbacks.Shop
			o := p.TemporaryUpdates.World.Objectives
			a := o.Attack
			sp := a.Controls.SpellLifecycle
			a.Controls.X = int32(n)
			a.Controls.Y = int32(mode)
			sp.Z = int32(mode & 1)
			for i := 0; i < n; i++ {
				sp.Record.Words[i*4] = 1
			}
			switch mode {
			case 1:
				a.ActorWords = map[int]uint32{16: 0x20}
			case 2:
				a.UpdateWords = map[int]uint32{280: 1}
			case 3:
				a.UpdateWords = map[int]uint32{216: 1}
			case 4:
				sp.Record.Words[16] = 137
			case 5:
				o.PlayerDataWords[0][3700] = 0
			case 6:
				p.Resources.PlayerClass = 0
			case 7:
				p.Resources.PlayerClass = 2
			case 8:
				sp.Record.Words[0] = 34
				o.PlayerDataWords[0][3696+4*34] = 3
			case 9:
				sp.Record.Words[0] = 34
				o.PlayerDataWords[0][3696+4*34] = 3
				p.Resources.OldMana = 0
			case 10:
				sp.Books = []legacy.PortTestSpellLifeRecord{{Words: map[int]uint32{8: 1}, Refs: map[int]int{4: 101}}}
			case 11:
				sp.Record.Words[0] = 0xffffffff
			}
			p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: 1502}, legacy.PortTestShopAction{Op: 1517})
			cases = append(cases, s)
		}
	}
	r := spellLifeHash(t, "spell-life-book-insertion", cases)
	if got := r[0].Callbacks.Shop.Sequence[0].Return; got != 1 {
		t.Fatalf("positive book insertion returned %d", got)
	}
}

func TestSpellLifecyclePowerAndPosition(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, subject := range []int{1, 2, 3} {
		for _, flags := range []uint32{0, 16, 32, 64, 256, 1024} {
			for _, power := range []uint32{0, 1, 3, 255, 65535} {
				s := spellLifeBase(15)
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				a := o.Attack
				p.Resources.Subject = subject
				s.Lifecycle.GameFlags = 1 | flags
				o.PlayerDataWords[0][3700] = power
				if subject == 3 {
					a.UpdateWords = map[int]uint32{2040: power}
				}
				if power == 65535 {
					a.Controls.SpellLifecycle.ActorType = "ImaginaryCaster"
				}
				cases = append(cases, s)
			}
		}
	}
	spellLifeHash(t, "spell-life-power", cases)
	cases = nil
	for _, origin := range []float32{0, 100, -100, 16777216} {
		for _, dx := range []float32{-5.000001, -5, -4.999999, 0, 4.999999, 5, 5.000001} {
			for _, dy := range []float32{-5, 0, 5} {
				s := spellLifeBase(16)
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
				a.ActorWords = map[int]uint32{56: math.Float32bits(origin), 60: math.Float32bits(origin)}
				a.Controls.SpellLifecycle.Record.Words = map[int]uint32{0: math.Float32bits(origin + dx), 4: math.Float32bits(origin + dy)}
				cases = append(cases, s)
			}
		}
	}
	spellLifeHash(t, "spell-life-position", cases)
}
func TestSpellLifecycleBuffExclusions(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, typ := range []string{"", "Hecubah", "Necromancer"} {
		for _, flags := range []uint32{0, 2048, 4096, 6144} {
			for _, buff := range []int32{0, 3, 11, 29} {
				for mode := 0; mode < 4; mode++ {
					s := spellLifeBase(23)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					s.Lifecycle.GameFlags = 1 | flags
					p.Resources.Subject = 3
					a.Controls.X = buff
					a.Controls.Y = 30
					a.Controls.SpellLifecycle.Z = 5
					a.Controls.SpellLifecycle.ActorType = typ
					a.ActorWords = map[int]uint32{12: 0x1000, 340: 1}
					if mode == 1 {
						a.ActorWords[12] = 0
					}
					if mode == 2 {
						a.ActorWords[16] = 0x20
					}
					if mode == 3 {
						a.Actor = 0
						a.Controls.SpellLifecycle.ActorType = ""
					}
					cases = append(cases, s)
				}
			}
		}
	}
	spellLifeHash(t, "spell-life-buff-exclusions", cases)
}
func TestSpellLifecycleProjectilesAndShock(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, dir := range []uint32{0, 1, 31, 64, 127, 192, 255} {
		for mode := 0; mode < 8; mode++ {
			s := spellLifeBase(10)
			p := s.Callbacks.Shop
			a := p.TemporaryUpdates.World.Objectives.Attack
			a.ActorWords = map[int]uint32{124: dir | dir<<16, 80: math.Float32bits(0.125), 84: math.Float32bits(-0.375), 176: math.Float32bits(8.25)}
			a.ProjectileSpeed = 3.25
			if mode&1 != 0 {
				a.Controls.Target = 0
			}
			if mode&2 != 0 {
				s.Combat.Wall = 1
			}
			if mode&4 != 0 {
				a.ActorWords[340] = 1 << 21
				a.ActorWords[384] = 123 << 16
				a.ActorWords[428] = 7 << 8
			}
			cases = append(cases, s)
		}
	}
	r := spellLifeHash(t, "spell-life-projectiles", cases)
	if got := r[0].Callbacks.Shop.Sequence[0].Return; got == 0 {
		t.Fatal("clear ray did not create Magic projectile")
	}
	cases = nil
	for _, power := range []uint32{1, 2, 3, 4, 5} {
		for mode := 0; mode < 8; mode++ {
			s := spellLifeBase(11)
			a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
			a.Controls.Target = 2
			a.ActorWords = map[int]uint32{340: 1 | (1 << 22), 428: power << 16}
			s.Combat.Friendly = mode&1 != 0
			if mode&2 != 0 {
				s.Combat.TargetFlags = 8
			}
			if mode&4 != 0 {
				a.ActorWords[340] = 1
			}
			cases = append(cases, s)
		}
	}
	spellLifeHash(t, "spell-life-shock", cases)
}

func TestSpellLifecycleCastingRestrictions(t *testing.T) {
	for _, op := range []int{8, 9} {
		var cases []legacy.PortTestRoamSpec
		for _, flags := range []uint32{0, 16, 32, 64} {
			for mode := 0; mode < 8; mode++ {
				s := spellLifeBase(op)
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				a := o.Attack
				s.Lifecycle.GameFlags = 1 | flags
				a.Controls.Y = int32(mode & 1)
				a.Controls.SpellLifecycle.Definitions = []server.PortTestSpellLifecycleDef{{Index: 1, Valid: true, Enabled: true, Flags: 0x1080000}}
				if mode&2 != 0 {
					a.ActorWords = map[int]uint32{340: 1 << 29}
					o.PlayerDataWords[0][4] = 1
				}
				if mode&4 != 0 {
					p.Inventory.Linked = []int{0}
					p.Items[0].Class = 0x10000000
				}
				cases = append(cases, s)
			}
		}
		spellLifeHash(t, fmt.Sprintf("spell-life-restrictions-%02d", op), cases)
	}
}
func TestSpellLifecycleBookCounter(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, radius := range []float32{0, 0.1, 0.10001, 10, 1000} {
		for mode := 0; mode < 4; mode++ {
			s := spellLifeBase(14)
			p := s.Callbacks.Shop
			o := p.TemporaryUpdates.World.Objectives
			a := o.Attack
			a.Controls.X = int32(math.Float32bits(radius))
			a.ActorWords = map[int]uint32{56: math.Float32bits(100), 60: math.Float32bits(100)}
			for i := 0; i < 3; i++ {
				o.PlayerWords[i][56] = math.Float32bits(float32(100 + 10*i))
				o.PlayerWords[i][60] = math.Float32bits(100)
			}
			a.Controls.SpellLifecycle.Books = []legacy.PortTestSpellLifeRecord{{Words: map[int]uint32{8: 1}, Refs: map[int]int{4: 101}}, {Words: map[int]uint32{8: 1}, Refs: map[int]int{4: 100}}, {Words: map[int]uint32{8: 1}, Refs: map[int]int{4: 102}}}
			if mode&1 != 0 {
				s.Combat.Wall = 1
			}
			if mode&2 != 0 {
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: 1501})
			}
			cases = append(cases, s)
		}
	}
	spellLifeHash(t, "spell-life-book-counter", cases)
}
