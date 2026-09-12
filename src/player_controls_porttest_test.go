//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

var controlsHashes = map[string]string{
	"controls-admission-36":    "056d7dd23051984fc76109e2ccb56c3e5126e04cca89b3e32293dfa45f4511bb",
	"controls-admission-37":    "f06248b0fa857c700a5382ac058cff02a1f0a9ecf8c51aed49995badc6df0cb1",
	"controls-stamina-28":      "634fbbcf2c81773fc6cfa5c8211ae9fe5200509e0734645bfbd1d30a5ba186c8",
	"controls-stamina-29":      "5e680e4f6d904801d42d0918683082a7bcb16f1e162576281befe608ac3ba531",
	"controls-weapon-map-31":   "c3b71aa404aae40b4a34720458be31843f108fb17c785f899c6c1ec3332d23cf",
	"controls-weapon-map-41":   "1b2ef838cfbf86be0176b1ffc079d045390e31d2f29f93926b7ea1eefec750f5",
	"controls-action-map":      "75ccf98684dab6682f54c0abfda0ec2beb7c72d8899c61d0dc29f0bf3b140d5f",
	"controls-confusion":       "7956441ff6ffad9a22615b9d63a8df60e172cf178d6696f36bffc51cf305bb57",
	"controls-bolt-damage":     "4990efebe87a9af5681ad0c63d700ab729edcd88f09e0340db18ae7600c51a92",
	"controls-messages-50":     "4f54e650f7d57fbf8476f14222a7d7d920226426b074b0223f4c606e34143235",
	"controls-messages-51":     "95b332f86e0c5e83f1a82faa5188584ad733f832f712850375343e9e39b1969b",
	"controls-waypoints-23":    "17467201e4e80aed74c91b75059d19cd1f091c179768c8c4374507cb7f3390ba",
	"controls-waypoints-24":    "38f254ed374213ad2eb02a0afcda9245da2b220e3d1057651e13c649547e719f",
	"controls-waypoints-34":    "04adb6d2a956b40295d67c7abe49b2fd50a440eceb4f4689150f19b82b1669d3",
	"controls-waypoints-35":    "152e7f9733c577703dc758e2c67e4eec404ac6db6db1b643288ae4f828d98e56",
	"controls-ownership-00":    "e27c57a19f9c943869ace2f934b37c5c2af87745853cde9e26ad6779edf7beef",
	"controls-ownership-01":    "688e40f04fe44a22015cfa97d23b08e8b1dfece3a5a586aef429b4414b20423f",
	"controls-ownership-04":    "8c2aa5e83bd1be7477098a00febb953b2b8d187deaa1f5c88dc460196ae9fc55",
	"controls-ownership-06":    "3f5211efa2d662a220d100e4b3def341aeb429a8cca26d580454f4a3c478396e",
	"controls-ownership-07":    "21bf45136a0abcc9aa6e11ba27502ec3b676965062c695bcd216ee4693cda4d6",
	"controls-ownership-08":    "0c625374221bc8fe6bb9bc5df966cdd04065595099b900d9583f9f2553be1a7f",
	"controls-ownership-09":    "136017f6604369d7773404c644139c51172bb6c3aefd020f0ec24d79a17db3cc",
	"controls-ownership-16":    "a912fce879fe2a3bab9ccf943e162caa398fcb85d6912d1baa517473ae38d585",
	"controls-ownership-22":    "3d5866ffd53e197701defe40f8031853c374c45a4282b98c0526557dc4477121",
	"controls-start-selection": "02f9935909787c29740647f2270cf4d5ef6378c8db569c47b9ffe28b2f535b89",
	"controls-stats-10":        "cb60bf8440c5cbf7019e789df0ff5ea2707c8369bc7567c45082d6e6f1887833",
	"controls-stats-11":        "6139c3e22dbd5b9fa6b452e24e3ecd2567dfb623bc91dd2d5e035c8b26775c53",
	"controls-stats-12":        "1eb30e0b28bfd067925f8232ea8b451b87157747abe0cec2d972f14bb46901f9",
	"controls-stats-14":        "da1c87627057955ea2177455dc4b8540b5da4c4fbdea42741155b61da513f303",
	"controls-bots-44":         "0ec444618afa125d1b488a6c7431c5089ef9b252fdcd1f1105c2b3b19f9a30d1",
	"controls-bots-48":         "27b2cb5526d538f0a7154a88241e9321216a05379d8091fb47bd9e976e4b1c76",
	"controls-observers-02":    "5dd6d03cd000ba161870f5cafd73d394f2352854d377b96299f20cbc20279a0f",
	"controls-observers-03":    "151f39a79dd0a31cd0ac21e4e204310dfd3e9e6cb27e4f10d2e21bdc1656b164",
	"controls-observers-05":    "c3aeacf10f30590fb69136e5df3eed05c2c8e7942a0ec026fe4a7a3a15225c6b",
	"controls-observers-39":    "3ff1375a509bcfbb4990f38fd5fa368d0e127ec8281c20dd805bd67a253fe353",
	"controls-observers-40":    "ab321bac629e96d2c72bea63ccd801b68e2f5e72d9f654ea933b76302ab3dfa3",
	"controls-scheduled-54":    "00393fce258d84dc75d52dd0a8fff862d57e064f472fefc0c97c242a3173a51e",
	"controls-scheduled-55":    "592a5932aa057491e34e87eabf84357699bfeb3f15ef567e9db9af3c10c440fd",
	"controls-respawn-items":   "efe2c5b21cf716c77d3d2ff86eb09fa5f089b3b64ba1116c5be191bb4dd9fff4",
	"controls-respawn-15":      "77e817af3997900cc7e58b873081ff5d4acd9c8473a767475ec9121cd39bb17d",
	"controls-respawn-18":      "f86148ae8a364efebb94999403c554f00ec7e0beb74fba9da4479441fff8c25b",
	"controls-respawn-19":      "29dd22cae7a7f91545267bbab4728dd3ed12c015a2352d9178404585bbe16083",
	"controls-respawn-20":      "3e699769bfd10cc7066a19cffc6d9caf954967012a1efee84d944687b44e206f",
	"controls-respawn-21":      "8e70bcf9d354b9be28c6b2f9b1bbd1f03570940fdbc21b5629b640f84472a759",
	"controls-respawn-32":      "0638faddbb182fc8017d926cae3f2129ea6f168ac9565e93de9896d0c3d3c26a",
	"controls-respawn-49":      "1a09e487f82cc158e9318b985f8245c5886f28f50d3af1470f40f1b196dc7433",
	"controls-remaining-27":    "8fd8c0c57adeff60fa501c33490b344ef0e148d73f53c98134efed7f3fd67b6c",
	"controls-remaining-30":    "dd034e53c80d2c3393c02a2d058aa2a22f3ea1e8f7564d0f0d44676a6504b2a2",
	"controls-remaining-33":    "7da5d13fe937b14559eb6b0c5dcf4b805353ad4306e29cfe59e0b83bdabf893e",
	"controls-remaining-38":    "f14ead6e02c3bb2cf1796285e5d58a381ff5bc6a75c0503e918fd1de80b583d6",
	"controls-remaining-43":    "47509415432b82aa1729cd8a69ecc82184e5052b963601f8ed36a377d6dace7b",
	"controls-remaining-52":    "b8e0a83ce62f1e13b433789b9fe80c022ceb3f53f4824ec2b5aecdf29d177c1b",
	"controls-remaining-53":    "3f7c08c56d489ac259030879d73e880516aea2a64ec2d317f314ede7d182c45a",
	"controls-respawn-flags":   "ca9bf86a6b9987a291a5e6be73e3e0da2551be9b639553673fb7a8081fef514c",
	"controls-morphs":          "8ba223073a6bce3276f72c05a652509ae4f652fa5bd3066cfb1eb15d5bb59cf1",
	"controls-bot-update":      "75dd94117ef8a0608a647168802b23f0f407f0bb3fca737451a3adb33534754f",
	"controls-inversion":       "6f1da123c2be2a53803ec5c5d03a3702f31ab7df2d7d873ce2777a4cbfc33e3c",
	"controls-positive":        "746747c5dc299b3f5d12a22e4c9de767d6aa38f2d6aac199d542655210885aeb",
}

func controlsBase(op int) legacy.PortTestRoamSpec {
	s := attackBase()
	p := s.Callbacks.Shop
	a := p.TemporaryUpdates.World.Objectives.Attack
	a.Actor = 1
	a.Controls = &legacy.PortTestPlayerControlsSpec{Target: 4}
	a.RecordWords = map[int]uint32{}
	a.RecordRefs = nil
	p.Inventory.Linked = nil
	p.Inventory.Owned = nil
	p.Equipment.ActiveWeapon = 0
	p.Equipment.WeaponFlags = 0
	p.Sequence = []legacy.PortTestShopAction{{Op: 1400 + op}}
	return s
}
func controlsHash(t *testing.T, name string, cases []legacy.PortTestRoamSpec) []legacy.PortTestRoamResult {
	t.Helper()
	results := controlsRun(t, cases)
	callbackHash(t, name, results, controlsHashes[name])
	return results
}
func controlsRun(t *testing.T, cases []legacy.PortTestRoamSpec) []legacy.PortTestRoamResult {
	t.Helper()
	oldUpdate := legacy.Nox_xxx_unitUpdateMonster_50A5C0
	legacy.Nox_xxx_unitUpdateMonster_50A5C0 = func(u *server.Object) {
		core := u.Server()
		oldOwner := core.ExtServer
		owner := &Server{Server: core}
		owner.ai.Init(owner)
		core.ExtServer = unsafe.Pointer(owner)
		defer func() { core.ExtServer = oldOwner }()
		oldUpdate(u)
	}
	defer func() { legacy.Nox_xxx_unitUpdateMonster_50A5C0 = oldUpdate }()
	old := legacy.Nox_xxx_playerCancelAbils_4FC180
	legacy.Nox_xxx_playerCancelAbils_4FC180 = func(u *server.Object) {
		a := serverAbilities{s: &Server{Server: legacy.GetServer().S()}}
		a.CancelAbilities(u)
	}
	defer func() { legacy.Nox_xxx_playerCancelAbils_4FC180 = old }()
	return effectsTimedRun(t, cases)
}
func TestPlayerControlsAdmission(t *testing.T) {
	for _, op := range []int{36, 37} {
		var cases []legacy.PortTestRoamSpec
		for _, buffs := range []uint32{0, 1 << 25, 1 << 5, 1<<25 | 1<<5} {
			for _, state := range []uint32{0, 1, 2, 23, 255} {
				for mode := 0; mode < 4; mode++ {
					s := controlsBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.ActorWords = map[int]uint32{340: buffs}
					a.UpdateWords = map[int]uint32{88: state, 280: uint32(mode & 1)}
					if mode&2 != 0 {
						s.Lifecycle.GameFlags |= 4096
					}
					if state == 1 {
						a.UpdateRefs = map[int]int{104: 3}
						p.Items[0].Class = 0x1000000
						p.Items[0].Subclass = uint32(mode * 4)
					}
					cases = append(cases, s)
				}
			}
		}
		controlsHash(t, fmt.Sprintf("controls-admission-%02d", op), cases)
	}
}
func TestPlayerControlsStamina(t *testing.T) {
	for _, op := range []int{28, 29} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 2, 3} {
			for _, stamina := range []uint32{0, 1, 44, 45, 100, 255} {
				for _, amount := range []int32{-257, -1, 0, 1, 44, 45, 100, 255, 256, 0x7fffffff} {
					s := controlsBase(op)
					p := s.Callbacks.Shop
					p.Resources.Subject = subject
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.Controls.X = amount
					if subject == 1 {
						a.UpdateWords = map[int]uint32{88: stamina << 24}
					} else if subject == 3 {
						a.UpdateWords = map[int]uint32{1128: stamina}
					}
					p.Sequence = append(p.Sequence, p.Sequence[0])
					cases = append(cases, s)
				}
			}
		}
		controlsHash(t, fmt.Sprintf("controls-stamina-%02d", op), cases)
	}
}
func TestPlayerControlsWeaponMappings(t *testing.T) {
	for _, op := range []int{31, 41} {
		var cases []legacy.PortTestRoamSpec
		masks := []uint32{0, 1, 2, 3, 0xffffffff, 0x80000000, 0x47f0000, 0x7ff8000}
		for i := 0; i < 32; i++ {
			masks = append(masks, 1<<i, 1<<i|0x400)
		}
		for _, mask := range masks {
			s := controlsBase(op)
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.X = int32(mask)
			cases = append(cases, s)
		}
		controlsHash(t, fmt.Sprintf("controls-weapon-map-%02d", op), cases)
	}
}
func TestPlayerControlsActionMapping(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for state := uint32(0); state < 36; state++ {
		for _, weapon := range []uint32{0, 1, 4, 0x100, 0x400, 0x800, 0x4000} {
			s := controlsBase(42)
			p := s.Callbacks.Shop
			o := p.TemporaryUpdates.World.Objectives
			p.Equipment.WeaponFlags = weapon
			o.Attack.UpdateWords = map[int]uint32{88: state}
			o.PlayerDataWords[0] = map[int]uint32{8: 7}
			cases = append(cases, s)
		}
	}
	controlsHash(t, "controls-action-map", cases)
}
func TestPlayerControlsConfusion(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, dir := range []uint32{0, 1, 32, 127, 255, 32767, 32768, 65535} {
		for _, power := range []uint32{0, 1, 5, 255} {
			for _, frame := range []uint32{0, 1, 9, 10, 19, 20, 29, 30, 39, 40, 0xffffffff} {
				s := controlsBase(25)
				s.Owner.Frame = frame
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
				a.ActorWords = map[int]uint32{124: dir << 16, 408: power << 24}
				cases = append(cases, s)
			}
		}
	}
	controlsHash(t, "controls-confusion", cases)
}
func TestPlayerControlsBoltDamage(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, level := range []int32{-1, 0, 1, 10, 32767} {
		for _, base := range []uint32{0, 1, 65535} {
			for _, rate := range []float32{0, 0.1, 1.25, -0.5} {
				s := controlsBase(13)
				a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
				a.Controls.X = level
				a.RecordWords = map[int]uint32{60: base, 64: math.Float32bits(rate)}
				cases = append(cases, s)
			}
		}
	}
	controlsHash(t, "controls-bolt-damage", cases)
}

func TestPlayerControlsMessages(t *testing.T) {
	for _, op := range []int{50, 51} {
		var cases []legacy.PortTestRoamSpec
		for _, actor := range []int{0, 1, 3} {
			for _, kind := range []int32{-1, 0, 1, 2, 3} {
				for _, selector := range []int32{0, 1, 4} {
					s := controlsBase(op)
					a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
					a.Actor = actor
					a.Controls.X = kind
					a.Controls.Y = selector
					if op == 51 {
						a.Controls.X = selector
						name := []string{"", "objcoll.c:GateLockedKey", "objcoll.c:DoorLockedKey", "a", "unknown"}[kind+1]
						a.Controls.Name = &name
					}
					cases = append(cases, s)
				}
			}
		}
		controlsHash(t, fmt.Sprintf("controls-messages-%02d", op), cases)
	}
}
func TestPlayerControlsWaypoints(t *testing.T) {
	for _, op := range []int{23, 24, 34, 35} {
		var cases []legacy.PortTestRoamSpec
		for index := uint32(0); index < 3; index++ {
			for _, present := range []bool{false, true} {
				for _, distance := range []float32{0, 9.99, 10, 10.01, 30} {
					s := controlsBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.UpdateWords = map[int]uint32{180: index | index<<8}
					a.Controls.X = int32(math.Float32bits(520))
					a.Controls.Y = int32(math.Float32bits(530))
					if present {
						a.UpdateRefs = map[int]int{168 + 4*int(index): 4}
					}
					p.TemporaryUpdates.ItemWords[1][56] = math.Float32bits(512 + distance)
					p.TemporaryUpdates.ItemWords[1][60] = math.Float32bits(512)
					cases = append(cases, s)
				}
			}
		}
		controlsHash(t, fmt.Sprintf("controls-waypoints-%02d", op), cases)
	}
}
func TestPlayerControlsOwnership(t *testing.T) {
	for _, op := range []int{0, 1, 4, 6, 7, 8, 9, 16, 22} {
		var cases []legacy.PortTestRoamSpec
		for mode := 0; mode < 5; mode++ {
			for _, flags := range []uint32{0, 0x20, 0x8000, 0x8020} {
				s := controlsBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				if mode != 0 {
					p.Inventory.Owned = []int{0, 1, 2}
					p.Inventory.Linked = []int{0, 1, 2}
				}
				a.Controls.MonsterRefs = []int{3, 4, 5}
				a.Controls.UpdateByRef = map[int]map[int]uint32{3: {1440: 0x80}, 4: {1440: 0x80}, 5: {1440: 0x80}}
				for i := 0; i < 3; i++ {
					p.Items[i].Class = 2
					p.Items[i].Flags = flags
					p.Items[i].Subclass = 0x80
				}
				p.TemporaryUpdates.World.Objectives.ObjectList = []int{3, 4, 5}
				if mode == 2 {
					p.Items[1].Flags = 0
					a.Controls.UpdateByRef[3][1440] = 0
				}
				if mode == 3 {
					p.TemporaryUpdates.World.ItemNames = []string{"Glyph", "Glyph", "Glyph"}
				}
				if op == 7 {
					a.Actor = 3
				}
				if op == 9 {
					a.ActorRefs = map[int]int{508: 101}
				}
				if op == 4 && mode == 4 {
					p.TemporaryUpdates.World.Objectives.PlayerDataRefs[0] = map[int]int{3628: 3}
				}
				a.Controls.X = 77
				if op == 22 {
					p.TemporaryUpdates.ItemWords[1][36] = 77
				}
				cases = append(cases, s)
			}
		}
		controlsHash(t, fmt.Sprintf("controls-ownership-%02d", op), cases)
	}
}
func TestPlayerControlsStartSelection(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for players := 0; players < 4; players++ {
		for mask := 0; mask < 8; mask++ {
			for seed := 1; seed <= 4; seed++ {
				s := controlsBase(26)
				s.Seed = seed
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				o.Players = players
				o.ObjectList = []int{3, 4, 5}
				p.TemporaryUpdates.World.ItemNames = []string{"PlayerStart", "PlayerStart", "PlayerStart"}
				for i := 0; i < 3; i++ {
					p.Items[i].Flags = uint32((mask>>i)&1) << 24
					p.TemporaryUpdates.ItemWords[i][56] = math.Float32bits(float32(480 + 40*i))
					p.TemporaryUpdates.ItemWords[i][60] = math.Float32bits(float32(490 + 15*i))
				}
				o.PlayerWords[1][56] = math.Float32bits(470)
				o.PlayerWords[2][56] = math.Float32bits(580)
				cases = append(cases, s)
			}
		}
	}
	controlsHash(t, "controls-start-selection", cases)
}

func controlsStats(s *legacy.PortTestRoamSpec) {
	p := s.Callbacks.Shop
	p.Resources.ExtraProtection = true
	p.Resources.Protected = true
	p.TemporaryUpdates.World.Objectives.Attack.Controls.Stats = []server.ClassStats{
		{Health: 25.25, Mana: 15.125, Speed: 1500.25, Strength: 8.75},
		{Health: 150.5, Mana: 60.25, Speed: 4000.75, Strength: 40.25},
		{Health: 80.25, Mana: 150.5, Speed: 3500.125, Strength: 20.75},
		{Health: 110.75, Mana: 110.25, Speed: 3750.5, Strength: 30.125},
	}
	p.EffectsUse.Balance["XPTable"] = []float64{0, 10, 30, 60, 100, 150, 210, 280, 360, 450, 550, 660}
}
func TestPlayerControlsStats(t *testing.T) {
	for _, op := range []int{10, 11, 12, 14} {
		var cases []legacy.PortTestRoamSpec
		for class := byte(0); class < 3; class++ {
			for _, level := range []uint32{0, 1, 2, 5, 9, 10} {
				for _, mode := range []uint32{0, 4096, 8192} {
					s := controlsBase(op)
					controlsStats(&s)
					p := s.Callbacks.Shop
					p.Resources.PlayerClass = class
					s.Lifecycle.GameFlags |= mode
					o := p.TemporaryUpdates.World.Objectives
					o.PlayerDataWords[0] = map[int]uint32{3684: level}
					o.Attack.Controls.X = int32(level)
					o.Attack.ActorWords = map[int]uint32{28: math.Float32bits(float32(level) * 50.1)}
					if op == 10 {
						o.Attack.Controls.X = int32(level)
						o.PlayerDataWords[0][3700] = 3
						o.PlayerDataWords[0][3704] = 4
					}
					cases = append(cases, s)
				}
			}
		}
		controlsHash(t, fmt.Sprintf("controls-stats-%02d", op), cases)
	}
}
func TestPlayerControlsBots(t *testing.T) {
	for _, op := range []int{44, 48} {
		var cases []legacy.PortTestRoamSpec
		for class := byte(0); class < 3; class++ {
			for _, fps := range []uint32{1, 30, 60, 65535} {
				for i := 0; i < 16; i++ {
					s := controlsBase(op)
					p := s.Callbacks.Shop
					p.Resources.PlayerClass = class
					s.Owner.FPS = fps
					s.Seed = i + 1
					a := p.TemporaryUpdates.World.Objectives.Attack
					a.Controls.Bot = op == 48 || i%2 == 0
					a.Controls.BotWords = map[int]uint32{544: 0, 552: uint32(i * 3), 1440: uint32(i%2) << 14}
					if op == 44 {
						a.Controls.BotWords[0] = 0x87654321
					}
					cases = append(cases, s)
				}
			}
		}
		controlsHash(t, fmt.Sprintf("controls-bots-%02d", op), cases)
	}
}
func TestPlayerControlsObservers(t *testing.T) {
	for _, op := range []int{2, 3, 5, 39, 40} {
		var cases []legacy.PortTestRoamSpec
		for players := 1; players <= 3; players++ {
			for mode := 0; mode < 8; mode++ {
				s := controlsBase(op)
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				o.Players = players
				if mode&1 != 0 {
					o.PlayerWords[1][16] = 0x20
					o.PlayerDataWords[0] = map[int]uint32{3680: 1}
				}
				if mode&2 != 0 {
					o.PlayerDataRefs[0] = map[int]int{3628: 101}
				}
				if mode&4 != 0 {
					s.Lifecycle.GameFlags |= 64
					p.TemporaryUpdates.World.ItemNames = []string{"GameBall"}
					o.ObjectList = []int{3}
				}
				if op == 39 {
					o.Attack.UpdateRefs = map[int]int{288: 101}
					s.Combat.Friendly = mode&1 != 0
				}
				if op == 40 {
					o.Attack.ActorRefs = map[int]int{520: 101}
				}
				cases = append(cases, s)
			}
		}
		controlsHash(t, fmt.Sprintf("controls-observers-%02d", op), cases)
	}
}
func TestPlayerControlsScheduledSpells(t *testing.T) {
	for _, op := range []int{54, 55} {
		var cases []legacy.PortTestRoamSpec
		for count := uint32(0); count <= 5; count++ {
			for _, buffs := range []uint32{0, 1 << 29} {
				for _, coord := range []int32{-16777217, -1, 0, 16777217} {
					s := controlsBase(op)
					a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
					a.ActorWords = map[int]uint32{340: buffs}
					a.UpdateWords = map[int]uint32{212: count, 220: uint32(coord), 224: uint32(-coord)}
					for i := 0; i < 5; i++ {
						a.UpdateWords[192+4*i] = uint32(1 + i)
					}
					s.Callbacks.Shop.Sequence = append(s.Callbacks.Shop.Sequence, s.Callbacks.Shop.Sequence[0])
					cases = append(cases, s)
				}
			}
		}
		controlsHash(t, fmt.Sprintf("controls-scheduled-%02d", op), cases)
	}
}

func TestPlayerControlsRespawnItems(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, name := range []string{"StreetShirt", "Longsword", "missing"} {
		for _, attrs := range []bool{false, true} {
			for x := int32(0); x < 2; x++ {
				for y := int32(0); y < 2; y++ {
					s := controlsBase(17)
					a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
					a.Controls.Equipment = true
					a.Controls.Name = &name
					a.Controls.X = x
					a.Controls.Y = y
					a.Controls.NullRecord = !attrs
					if attrs {
						a.RecordRefs = map[int]int{0: 3, 4: 4}
					} // modifier identity fixture inputs, no active effect callbacks.
					cases = append(cases, s)
				}
			}
		}
	}
	controlsHash(t, "controls-respawn-items", cases)
}
func TestPlayerControlsRespawn(t *testing.T) {
	for _, op := range []int{15, 18, 19, 20, 21, 32, 49} {
		var cases []legacy.PortTestRoamSpec
		for class := byte(0); class < 3; class++ {
			for mode := 0; mode < 4; mode++ {
				s := controlsBase(op)
				controlsStats(&s)
				p := s.Callbacks.Shop
				p.Resources.PlayerClass = class
				a := p.TemporaryUpdates.World.Objectives.Attack
				a.Controls.Equipment = true
				a.Controls.X = int32(mode & 1)
				a.Controls.Y = int32(mode >> 1)
				if op == 32 || op == 49 {
					a.Controls.Corpse = true
				}
				if op == 18 || op == 20 {
					if mode&1 != 0 {
						a.Controls.ByteReturn = 1
					} else if class < 2 && (op == 20 || mode>>1 == 0) {
						a.Controls.ByteReturn = 2
					}
				}
				if op == 49 {
					a.Controls.Bot = true
					a.Controls.BotWords = map[int]uint32{548: 0}
					if mode&1 != 0 {
						p.Resources.HP = 0
					}
				}
				p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{3684: 5, 4700: uint32(mode & 1)}
				// The isolated setup supplies local respawn positioning and item callbacks.
				cases = append(cases, s)
			}
		}
		controlsHash(t, fmt.Sprintf("controls-respawn-%02d", op), cases)
	}
}
func TestPlayerControlsRemaining(t *testing.T) {
	for _, op := range []int{27, 30, 33, 38, 43, 52, 53} {
		var cases []legacy.PortTestRoamSpec
		for mode := 0; mode < 8; mode++ {
			s := controlsBase(op)
			p := s.Callbacks.Shop
			a := p.TemporaryUpdates.World.Objectives.Attack
			a.Controls.X = int32(mode % 3)
			if op == 27 {
				a.Actor = 3
				p.Items[0].Flags = uint32(mode&1) << 24
			}
			if op == 30 && mode&1 != 0 {
				p.Inventory.Owned = []int{0}
				p.TemporaryUpdates.World.ItemNames = []string{"GameBall"}
			}
			if op == 38 {
				p.Equipment.WeaponFlags = uint32(mode&1) * 0x100
				a.UpdateWords = map[int]uint32{88: uint32(mode%3) | 100<<24}
			}
			if op == 43 {
				p.Inventory.Linked = []int{0}
				p.Items[0].Flags = uint32(mode&1) << 8
			}
			if op == 52 || op == 53 {
				a.Controls.Guide = true
				p.TemporaryUpdates.World.Objectives.PlayerDataWords[0] = map[int]uint32{4248: uint32(mode & 1)}
				p.Items[1].Class = 2
				p.TemporaryUpdates.World.ItemNames = []string{"", "Bat"}
				p.EffectsUse.Balance["FieldGuideDamageBonus"] = []float64{1.25}
				a.RecordWords = map[int]uint32{0: uint32(mode * 10)}
			}
			cases = append(cases, s)
		}
		controlsHash(t, fmt.Sprintf("controls-remaining-%02d", op), cases)
	}
}

func TestPlayerControlsRespawnFlags(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for mask := 0; mask < 256; mask++ {
		s := controlsBase(15)
		a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
		a.Controls.Equipment = true
		a.Controls.Disallowed = uint8(mask)
		cases = append(cases, s)
	}
	results := controlsHash(t, "controls-respawn-flags", cases)
	for i, r := range results {
		if got, want := r.Callbacks.Shop.Sequence[0].Return, uint32(int32(int8(^uint8(i)))); got != want {
			t.Fatalf("mask %d flags %x want %x", i, got, want)
		}
	}
}
func TestPlayerControlsMorphs(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for class := byte(0); class < 3; class++ {
		for _, flags := range []uint32{4, 5, 0x80000004} {
			s := controlsBase(56)
			p := s.Callbacks.Shop
			p.Resources.PlayerClass = class
			a := p.TemporaryUpdates.World.Objectives.Attack
			a.Controls.Bot = true
			a.ActorWords = map[int]uint32{8: flags, 12: 123}
			p.Sequence = append(p.Sequence, p.Sequence[0])
			cases = append(cases, s)
		}
	}
	controlsHash(t, "controls-morphs", cases)
}
func TestPlayerControlsBotUpdate(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for class := byte(0); class < 3; class++ {
		for _, frame := range []uint32{1, 30, 60, 120} {
			s := controlsBase(47)
			s.Owner.Frame = frame
			p := s.Callbacks.Shop
			p.Resources.PlayerClass = class
			a := p.TemporaryUpdates.World.Objectives.Attack
			a.Controls.Bot = true
			p.Sequence = []legacy.PortTestShopAction{{Op: 1444}, {Op: 1447}}
			cases = append(cases, s)
		}
	}
	controlsHash(t, "controls-bot-update", cases)
}

func TestPlayerControlsInversion(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for mask := 0; mask < 16; mask++ {
		for _, value := range []uint32{0, 1, 0xffffffff} {
			for _, flags := range []uint32{0, 0x100} {
				s := controlsBase(43)
				p := s.Callbacks.Shop
				p.Inventory.Linked = []int{0}
				p.Items[0].Flags = flags
				for i := 0; i < 4; i++ {
					p.Items[0].Mods[i] = mask&(1<<i) != 0
					p.EffectsUse.Modifiers[i] = legacy.PortTestEffectsModifier{Collide: effectsIdentity(legacy.PortTestEffects4E03D0), Words: map[int]uint32{96: value}}
				}
				cases = append(cases, s)
			}
		}
	}
	controlsHash(t, "controls-inversion", cases)
}

func TestPlayerControlsPositive(t *testing.T) {
	stats := controlsBase(11)
	controlsStats(&stats)
	stats.Callbacks.Shop.Resources.PlayerClass = 1
	stats.Lifecycle.GameFlags |= 8192
	respawn := controlsBase(18)
	controlsStats(&respawn)
	respawn.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.Equipment = true
	respawn.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.ByteReturn = 2
	glyph := controlsBase(16)
	gp := glyph.Callbacks.Shop
	gp.Inventory.Linked = []int{0, 1, 2}
	gp.TemporaryUpdates.World.ItemNames = []string{"Glyph", "Glyph", "Glyph"}
	inversion := controlsBase(43)
	ip := inversion.Callbacks.Shop
	ip.Inventory.Linked = []int{0}
	ip.Items[0].Flags = 0x100
	ip.Items[0].Mods[3] = true
	ip.EffectsUse.Modifiers[3] = legacy.PortTestEffectsModifier{Collide: effectsIdentity(legacy.PortTestEffects4E03D0), Words: map[int]uint32{96: 1}}
	front, back := controlsBase(54), controlsBase(55)
	for _, s := range []*legacy.PortTestRoamSpec{&front, &back} {
		s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.UpdateWords = map[int]uint32{192: 1, 196: 2, 200: 3, 212: 3}
	}
	r := controlsHash(t, "controls-positive", []legacy.PortTestRoamSpec{stats, respawn, glyph, inversion, front, back})
	hp := r[0].Callbacks.Shop.Sequence[0].ResourceData[2][0] & 65535
	_, ud, _ := objectivePlayerData(r[0].Callbacks.Shop.Sequence[0], 0)
	if hp != 80 || ud[2]&65535 != 150 {
		t.Fatalf("max-stat restoration: HP=%d mana=%d", hp, ud[2]&65535)
	}
	if len(r[1].Lifecycle.Created) != 4 {
		t.Fatalf("default gear: %d allocations", len(r[1].Lifecycle.Created))
	}
	if got := r[2].Callbacks.Shop.Sequence[0].Return; got != 3 {
		t.Fatalf("glyph count %d", got)
	}
	if got := r[3].Callbacks.Shop.Sequence[0].Return; got != 1 {
		t.Fatalf("inversion result %d", got)
	}
	_, f, _ := objectivePlayerData(r[4].Callbacks.Shop.Sequence[0], 0)
	_, b, _ := objectivePlayerData(r[5].Callbacks.Shop.Sequence[0], 0)
	if f[48] != 2 || f[49] != 3 || f[50] != 0 || f[53]&255 != 2 {
		t.Fatal("front spell removal did not shift and clear")
	}
	if b[48] != 1 || b[49] != 2 || b[50] != 3 || b[53]&255 != 2 {
		t.Fatal("back spell removal did not retain queue words")
	}
}
