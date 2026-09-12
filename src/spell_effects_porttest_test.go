//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

var spellEffectsHashes = map[string]string{
	"spell-effects-area-09":             "7bd6d587db5ccafa997ada2cd400da41dfa244643a4b28c88829e62faca96ae4",
	"spell-effects-area-12":             "ac7e0e448c44c70cf744904ad6d9463134a8ae320d6ad74ad6913d75218dc2e4",
	"spell-effects-area-13":             "0067da4ac6c68b142902b0dcbab66ddd5f225cb07dabe8789f0962c90e719481",
	"spell-effects-area-37":             "4e30752f85e8a9083aabf64e1e275f911c1e0d081209e2247c471bb262f11812",
	"spell-effects-area-38":             "b646b9e7554462f293d789de7884251dfdcd65290a3ebd3c22fc8643fccf57e8",
	"spell-effects-charm-05":            "35267d37477e298bf46bb5d619b2b656ecfcceed8f8c6c92ba25255e45949d3c",
	"spell-effects-charm-06":            "0059dde19e68b1e1c235df6869ea86010bca6b85f5dfb62dbc633118a67fe2ca",
	"spell-effects-charm-07":            "d3b8095dd64ea7c121c4cfd0fab2842cb81fca4a8145b64a7468a87c84d4c87b",
	"spell-effects-charm-08":            "3e334e8ca92c658866aa48b9cc1c0ad4905a35c7219bbd717113c589b49e1e42",
	"spell-effects-doors-25":            "e886e05e4ff08e69f0b1d3235aedec45a472da629ea92b10762e64035fe75a63",
	"spell-effects-doors-26":            "a65968c073827f743b795c52da8a563aaf3a80843bd703408e98d355dd6addfe",
	"spell-effects-doors-27":            "3cb692b6dd9ae7aae086d98da9bc0cfcc9e2b9fafd46fc4433b91e82736a1237",
	"spell-effects-doors-28":            "3e35c4b2949eba3785586431f058e59aa4a315e10c7367e343da465479b63836",
	"spell-effects-factory-failure-20":  "1c724cc3135d4a09292e5b6a03fb52b1a7cc815de3209ea5fa7d6d4a56498e86",
	"spell-effects-factory-failure-29":  "1c724cc3135d4a09292e5b6a03fb52b1a7cc815de3209ea5fa7d6d4a56498e86",
	"spell-effects-force-39":            "8d5392fda3c5c8bed127d1f76625340a4f2a7cc63f3ea1e22a11ecb1ba50f56d",
	"spell-effects-force-40":            "ed7f60ee7bf0b66f83e9b61ea8eb7dbc8a3ab9e6841663dbca7d031d21e6a28c",
	"spell-effects-force-41":            "c67b66ff486a21decc9539390b08c9e6edd8ffc8c873358a25858286c6e6e8c1",
	"spell-effects-glyph-selection":     "3a03ac07cdcf52633d35aa7a5e3b118ac317a8443bf628d984d7d39f063cfbe9",
	"spell-effects-heal-fumble-14":      "6431969437c0d203c3e43670ef221dc6e08c4e784a108d5c9058e5f061c53a60",
	"spell-effects-heal-fumble-36":      "36e67e006bddd6a1285c1ff5f273b93c6c5da0c424f2e9d11791d17f125e9ffe",
	"spell-effects-inversion":           "83f9aa5f60ce693970af23ee76b172576739b4682303e101591efa7d7e98970a",
	"spell-effects-portals-21":          "b78bb254e25b7dba8a4fdd49d909aa0e83bbc1908873f365af3095567c5e71ba",
	"spell-effects-portals-22":          "009df8d8f5c211fb8a4614c92d8220cea3a5fb2d5a5407c143a59fb0c9d42fec",
	"spell-effects-positive":            "fd19d851079ac3870d2bce2bd53b0d1838f0fa3d1f4b7254502fb7268b5c1833",
	"spell-effects-projectiles-17":      "25ccad54ba76dd3d747572bf17cae65ad590b5ddbebb76f80e81373c2e1c7abc",
	"spell-effects-projectiles-20":      "3c54f58003e4af2095a7592a223156f663ef8e82b0fd4e458695965b9aa433b0",
	"spell-effects-projectiles-29":      "67427699c277a608abce35be614e68b7cd56f7fb0051a49fb9a9c1b8f84749c3",
	"spell-effects-projectiles-30":      "d4bd0e7cc2b68a03ae2a353603face8798436515208d97b6047f5aaf00f0053e",
	"spell-effects-projectiles-31":      "b9c36dafca0a4586f4b8b63c6882c1d0435df130dd9db8730317d58f33cfb0aa",
	"spell-effects-projectiles-32":      "670d3753b7fa5c0f5bbc64ade260773bce1b224507dd1e179d6d4257da9a9ad8",
	"spell-effects-projectiles-33":      "c65202d1f9129e7a5a00ec56637aff67431b7ca1b76c200c94f1892dde358c3f",
	"spell-effects-projectiles-34":      "e8357cab9ff163504e538f1751c264b5b702095c855f870b12d184b6f8b11ab5",
	"spell-effects-projectiles-35":      "3b2f6495947e0dff11064845fa74e21ed33bff99cc38f8c81b354cb8799ccdea",
	"spell-effects-resource-buffs-10":   "5ccdad82dc8d7d3072658bc0f8761b6362f8b0a1b6cc5d62b67fede77f9734e2",
	"spell-effects-resource-buffs-11":   "b848d577402d1347f98ec0b905b9596b56567f38e71e315bc88bdf61510667c3",
	"spell-effects-resource-buffs-15":   "6c80987cc12b7c6c232304b3532e79435870f3037792d5014faa05fefda3ced8",
	"spell-effects-resource-buffs-16":   "36eadf63d343c5b7c96026bce16569995d93846be335699d05a638a89818177c",
	"spell-effects-resource-buffs-18":   "5c13f96fe05b0276052929df46e2499e04d7281ba4b25a7108960371480d7946",
	"spell-effects-resource-buffs-19":   "54c85059c026a01e0c8fac634f1b0f96fee6b4501f757c9a11baf46a05e088c7",
	"spell-effects-resource-buffs-24":   "4c26cf6025a8977e07dc77e6ea3f7d88bc283206d701c0f08838c29586038879",
	"spell-effects-smoke-00":            "6f30dd00cefde84e0aea1f2b184a479bbc7e860d07f8e75e25193d64d17efd31",
	"spell-effects-smoke-01":            "35a11f5d4783085b9936974840763a397f57ebabf2de6294807b51275c235d31",
	"spell-effects-smoke-02":            "8fe904643f433a7b3ae1a0a2e0ff1c707df0896f7d185ba3bdb73fb406a1e0c9",
	"spell-effects-smoke-03":            "4e2e460857e3aef0a030374af3efdd8697b8acd23c0872d86269056944950053",
	"spell-effects-smoke-04":            "5ac829cc63e5762700fe76656a7f78b71d9190c243c2989e4aec663a3a7a2abc",
	"spell-effects-smoke-05":            "01b0cb59172f91816ff75ae28d4dbb64eb4258d1a394188cb81872830f42affe",
	"spell-effects-smoke-06":            "61c0c968dba0a4dcb7ac527e8fed9861dedc1cae378fb4bbe618c261d7fb50f1",
	"spell-effects-smoke-07":            "be45ed96550aca0f32cef17ceb2197ad020321f91cb7f67d3c5bddb90a06411a",
	"spell-effects-smoke-08":            "c64077fee9db40d0ef75a5b9f934fbc688c4b9e906a8fc08e04acbde39a29c05",
	"spell-effects-smoke-09":            "b11d7b3331260a7e17613a6947fb6ada0a76695c5145eabd278dcf1986d7a21a",
	"spell-effects-smoke-10":            "b2c185ff90fa8bf563842d32a7c7d14fd662e2a078af0cb2ed6db136daeec0bc",
	"spell-effects-smoke-11":            "7cfe35d1cbf4fa10f7a4e6eb3ea7fb92932246a9c5713de0b1147f0ad533ac9d",
	"spell-effects-smoke-12":            "b6990adfff61cfcb67db6b4182db51e7971bd95c48583f298943032c99a2d78b",
	"spell-effects-smoke-13":            "099e03551af6f6180187871ecf0af8b4ab781ea199ac2ab5b295203e40af1c97",
	"spell-effects-smoke-14":            "65c716cb318e1902368ac76452abc1f89ef69563ae980437abd7ee828cb33711",
	"spell-effects-smoke-15":            "59a20681a9b5594358f77f60645552a106d9e91fcfcf40624b92dd14288152e2",
	"spell-effects-smoke-16":            "d63698c9f193dd0f48160f78efbee9594f23ccf7795279bb7bcf3e7c0cf1ba44",
	"spell-effects-smoke-17":            "41728ff452386aa3c1fc12f98cc59a4c7270d31f381681f4bb818807d4509822",
	"spell-effects-smoke-18":            "bfe8be3ce7bede41a659d5c2b3852608012601c6178d90b65741af33d3f32579",
	"spell-effects-smoke-19":            "0ac7c39c1692ca93f74cfae13c4d527576e8dfaa3100ec537302731e9e7acd65",
	"spell-effects-smoke-20":            "741b799426e2983428f9abc5291b73609b5336a59558964edc7a2ed33c89bfce",
	"spell-effects-smoke-21":            "8078bba7548d11113f6267dc9dff3013739913537aaaf2c29af2bf75bff49630",
	"spell-effects-smoke-22":            "7b5d854ba3a5eb324705c420a3e17b657f9b075fc77a74e0b885cdca9ca1f237",
	"spell-effects-smoke-23":            "6f30dd00cefde84e0aea1f2b184a479bbc7e860d07f8e75e25193d64d17efd31",
	"spell-effects-smoke-24":            "802f27f407f8e62477fd3755fbabec37ff6ffe095e7f8f524fd4ae700e08458b",
	"spell-effects-smoke-25":            "6f30dd00cefde84e0aea1f2b184a479bbc7e860d07f8e75e25193d64d17efd31",
	"spell-effects-smoke-26":            "41f5bbba51b6f50d7c5db421bea722bd49f566d2c1f491904e9187d3c1fc9450",
	"spell-effects-smoke-27":            "6f30dd00cefde84e0aea1f2b184a479bbc7e860d07f8e75e25193d64d17efd31",
	"spell-effects-smoke-28":            "6f30dd00cefde84e0aea1f2b184a479bbc7e860d07f8e75e25193d64d17efd31",
	"spell-effects-smoke-29":            "e70dd74379ea9f4f453ac38e10ed73b685f84667d38350e3ea2909605b5f98ff",
	"spell-effects-smoke-30":            "a5a6ed6c4427d8cdf2319d9789afd50a856ef3ddc043609b24e4505e80e0beb7",
	"spell-effects-smoke-31":            "173583e2dd9eef5dc838cb7b6b957c3764d24fe198bda2680f52ceadfe249118",
	"spell-effects-smoke-32":            "907a6093632eb1727baf45772ca3a302d7406c69454d79396abe4aa868af1cfa",
	"spell-effects-smoke-33":            "3d1c3c248122f51233ccc4a540efee51c3b9caf026c70294b8a1e673ef900466",
	"spell-effects-smoke-34":            "ef896a805b650b7c20b0a675fd91f0b300c10f012f221be40f237b4a604d4c66",
	"spell-effects-smoke-35":            "ca2ed5c3d97a6b03ddca3bd7d689f1bda288f4787a610b0bb79320d201b93306",
	"spell-effects-smoke-36":            "4fea9b33fc999b1ba21d543aff3a68a22cc9141aa3db0ff766a10c8fc9cf27d2",
	"spell-effects-smoke-37":            "b11d7b3331260a7e17613a6947fb6ada0a76695c5145eabd278dcf1986d7a21a",
	"spell-effects-smoke-38":            "9fca5ef33a5a8c851dd7dce33bb6919141248d770dd4c02f7561a87b06db0501",
	"spell-effects-smoke-39":            "82701a07c5fb8bd8ffbbb6b51c66c8f5ffd4ec55f3c805aff9fda75f8c228a71",
	"spell-effects-smoke-40":            "17f3416211c0fab56f05d353d0ca9f67b9e5be8b6022e0ab24ebbb8ce2f44232",
	"spell-effects-smoke-41":            "f038eac38de7478e67546ec55beec9083d588df4583da0f3cc958ef2d760858a",
	"spell-effects-summon-admission-01": "a8ff24d7fc7f96311392effea23395540c1c0c5b58638401b8f092719a297a46",
	"spell-effects-summon-admission-02": "56fec07e85bd62e9342fe37424682a2479b1cde133b1b2e2429aba9cbb8f2aed",
	"spell-effects-summon-admission-03": "61ac32cbd4308a86a31af77dbb12040c41b7cb00ae38f24a09b7ca080efffe0c",
	"spell-effects-summon-admission-04": "dcfa52fd43d2d8d0cf40af1a610a27cf322deb1b4f671276c332b9cd9e2b6a6a",
	"spell-effects-summon-capacity":     "4bdf09d1ad4bba9b423b6302df5bec8fb94b28f24844f52acfb908af68928b80",
	"spell-effects-summon-cost":         "a8732bffd3bd1f3528d3f96dfcc0dc79a49e6ef13033e850e7d78979f7a97cbb",
	"spell-effects-summon-sequence":     "40aa6c0f80153dcf5e24ae26582caeafd4dbbfef936fe106713fa2eea82a76ec",
}

func spellEffectsBase(op int) legacy.PortTestRoamSpec {
	s := spellLifeBase(0)
	p := s.Callbacks.Shop
	a := p.TemporaryUpdates.World.Objectives.Attack
	sp := a.Controls.SpellLifecycle
	sp.ActorType = "NPC"
	sp.Effects = &legacy.PortTestSpellEffectsSpec{Args: [3]int{1, 1, 1}, Caches: map[uintptr]uint32{2487700: 3}}
	a.Controls.X = 1
	sp.Z = 3
	sp.Record = legacy.PortTestSpellLifeRecord{Words: map[int]uint32{4: math.Float32bits(512), 8: math.Float32bits(512)}, Refs: map[int]int{0: 2}}
	p.Sequence = []legacy.PortTestShopAction{{Op: 1600 + op}}
	for _, name := range []string{"BurnDuration", "ConfuseEnchantDuration", "LesserHealAmount", "PullPowerCoeff", "PushPowerCoeff", "ShockEnchantDuration", "EarthquakeRange", "EarthquakeJiggle", "InversionRange", "StunEnchantDuration"} {
		p.EffectsUse.Balance[name] = []float64{30}
	}
	for _, name := range []string{"CharmLargeDuration", "CharmMediumDuration", "CharmSmallDuration", "EarthquakeDamage", "FireballSpeedCoeff", "FistOfVengeanceDamage", "FistSpeed", "MeteorDamage", "MeteorSpeed", "ShockTrapDamage", "SummonDuration", "ToxicCloudLifetime"} {
		p.EffectsUse.Balance[name] = []float64{1, 2.5, 5, 10, 20}
	}
	return s
}
func TestSpellEffectsSmoke(t *testing.T) {
	for op := 0; op < 42; op++ {
		t.Run(fmt.Sprintf("%02d", op), func(t *testing.T) {
			s := spellEffectsBase(op)
			p := s.Callbacks.Shop
			a := p.TemporaryUpdates.World.Objectives.Attack
			sp := a.Controls.SpellLifecycle
			if op == 0 {
				a.Controls.X = 75
			}
			if op >= 1 && op <= 7 {
				sp.Record = legacy.PortTestSpellLifeRecord{Words: map[int]uint32{4: 75, 8: 3, 52: math.Float32bits(512), 56: math.Float32bits(512)}, Refs: map[int]int{16: 1, 48: 2}}
			}
			if op == 5 || op == 6 || op == 7 {
				sp.Record.Words[4] = 7
			}
			if op == 22 {
				a.Controls.X = 46
			}
			if op == 38 {
				a.Actor = 2
			}
			if op == 40 {
				sp.Record.Refs = nil
				sp.Record.Words = map[int]uint32{0: math.Float32bits(512), 4: math.Float32bits(512)}
				sp.Effects.Floats = [3]uint32{math.Float32bits(10), math.Float32bits(5), math.Float32bits(2)}
			}
			if op == 41 {
				sp.Effects.RecordOutput = true
				sp.Effects.Output.Words = map[int]uint32{0: math.Float32bits(512), 4: math.Float32bits(512)}
				sp.Record.Words = map[int]uint32{4: math.Float32bits(10), 8: math.Float32bits(5), 12: math.Float32bits(20)}
			}
			spellEffectsHash(t, fmt.Sprintf("spell-effects-smoke-%02d", op), []legacy.PortTestRoamSpec{s})
		})
	}
}

func TestSpellEffectsResourcesAndBuffs(t *testing.T) {
	for _, op := range []int{10, 11, 15, 16, 18, 19, 24} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 3} {
			for _, level := range []int32{1, 5, 128} {
				for _, duration := range []float64{-1, 0, 1.5, 65536} {
					for mode := 0; mode < 2; mode++ {
						s := spellEffectsBase(op)
						p := s.Callbacks.Shop
						a := p.TemporaryUpdates.World.Objectives.Attack
						sp := a.Controls.SpellLifecycle
						p.Resources.Subject = subject
						p.Resources.HP = 7
						p.Resources.MaxHP = 100
						p.Resources.Mana = 7
						p.Resources.MaxMana = 200
						p.Resources.Poison = 5
						sp.Z = level
						sp.Record.Refs[0] = 1
						if mode == 1 {
							sp.Record.Refs[0] = 0
						}
						for _, key := range []string{"ConfuseEnchantDuration", "StunEnchantDuration", "ShockEnchantDuration"} {
							p.EffectsUse.Balance[key] = []float64{duration}
						}
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: 1600 + op})
						cases = append(cases, s)
					}
				}
			}
		}
		spellEffectsHash(t, fmt.Sprintf("spell-effects-resource-buffs-%02d", op), cases)
	}
}
func TestSpellEffectsProjectiles(t *testing.T) {
	for _, op := range []int{17, 20, 29, 30, 31, 32, 33, 34, 35} {
		var cases []legacy.PortTestRoamSpec
		for _, dir := range []uint32{0, 1, 31, 64, 127, 192, 255} {
			for mode := 0; mode < 4; mode++ {
				s := spellEffectsBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				a.ActorWords = map[int]uint32{124: dir | dir<<16, 80: math.Float32bits(0.125), 84: math.Float32bits(-0.375), 176: math.Float32bits(8.25)}
				sp.Z = int32(1 + mode)
				sp.Record.Refs[0] = 1
				if mode&1 != 0 {
					s.Combat.Wall = 1
				}
				if mode&2 != 0 {
					s.Lifecycle.GameFlags |= 2048
				}
				cases = append(cases, s)
			}
		}
		r := spellEffectsHash(t, fmt.Sprintf("spell-effects-projectiles-%02d", op), cases)
		if op == 20 {
			for i, x := range r {
				if len(x.Lifecycle.Created) != 1 {
					t.Errorf("fireball case %d did not create a projectile", i)
				}
			}
		}
	}
}
func TestSpellEffectsPositiveContracts(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{10, 15, 17, 21, 22, 29, 30, 31} {
		s := spellEffectsBase(op)
		p := s.Callbacks.Shop
		sp := p.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
		p.Resources.HP = 7
		p.Resources.MaxHP = 100
		sp.Record.Refs[0] = 1
		if op == 22 {
			p.TemporaryUpdates.World.Objectives.Attack.Controls.X = 46
		}
		cases = append(cases, s)
	}
	r := spellEffectsHash(t, "spell-effects-positive", cases)
	if hp := r[0].Callbacks.Shop.Sequence[0].ResourceData[2][0] & 65535; hp != 100 {
		t.Errorf("full heal HP=%d want 100", hp)
	}
	u, _, _ := objectivePlayerData(r[1].Callbacks.Shop.Sequence[0], 0)
	if u[85]&(1<<3) == 0 {
		t.Error("confusion enchantment was not applied")
	}
	for i, want := range []int{1, 1, 1, 1, 1, 48} {
		if got := len(r[i+2].Lifecycle.Created); got != want {
			t.Errorf("positive case %d created %d want %d", i+2, got, want)
		}
	}
}

func spellEffectsHash(t *testing.T, name string, cases []legacy.PortTestRoamSpec) []legacy.PortTestRoamResult {
	t.Helper()
	oldLimit := legacy.PortTestSpellEffectsSummonLimit
	legacy.PortTestSpellEffectsSummonLimit = nox_xxx_checkSummonedCreaturesLimit_500D70
	defer func() { legacy.PortTestSpellEffectsSummonLimit = oldLimit }()
	old := legacy.Nox_xxx_unitDoSummonAt_5016C0
	legacy.Nox_xxx_unitDoSummonAt_5016C0 = func(id int, pos types.Pointf, owner *server.Object, dir server.Dir16) *server.Object {
		prior := noxServer
		noxServer = &Server{Server: legacy.GetServer().S()}
		defer func() { noxServer = prior }()
		return old(id, pos, owner, dir)
	}
	defer func() { legacy.Nox_xxx_unitDoSummonAt_5016C0 = old }()
	oldPlayer := legacy.PortTestSpellLifePlayerSpell
	legacy.PortTestSpellLifePlayerSpell = func(u *server.Object) { owner := &Server{Server: legacy.GetServer().S()}; owner.PlayerSpell(u) }
	defer func() { legacy.PortTestSpellLifePlayerSpell = oldPlayer }()
	r := controlsRun(t, cases)
	callbackHash(t, name, r, spellEffectsHashes[name])
	return r
}

func TestSpellEffectsSummonSequence(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, subject := range []int{1, 3} {
		for _, size := range []byte{1, 2, 4} {
			for _, seq := range []uint32{0, 64999, 65000, 65535} {
				s := spellEffectsBase(1)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				p.TemporaryUpdates.World.Objectives.PlayerDataWords[0][3648] = 4
				p.Resources.Subject = subject
				tile := 0
				sp.Effects.Tile = &tile
				sp.Effects.Guides = []legacy.PortTestSpellEffectGuide{{Index: 1, Size: size}}
				sp.Effects.Caches[1570276] = seq
				sp.Record = legacy.PortTestSpellLifeRecord{Words: map[int]uint32{4: 75, 8: 3, 52: math.Float32bits(520), 56: math.Float32bits(515)}, Refs: map[int]int{16: 1}}
				// A one-frame duration lets the actual start and finish run at the same fixture frame.
				p.EffectsUse.Balance["SummonDuration"] = []float64{1, 1, 1}
				p.Sequence = []legacy.PortTestShopAction{{Op: 1601}, {Op: 1603}, {Op: 1604}}
				cases = append(cases, s)
			}
		}
	}
	r := spellEffectsHash(t, "spell-effects-summon-sequence", cases)
	for i, x := range r {
		if x.Callbacks.Shop.Sequence[0].Return != 0 || len(x.Lifecycle.Created) != 1 {
			t.Errorf("summon %d: start=%d created=%d", i, x.Callbacks.Shop.Sequence[0].Return, len(x.Lifecycle.Created))
		}
	}
}

func TestSpellEffectsSummonAdmission(t *testing.T) {
	for _, op := range []int{1, 2, 3, 4} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 3} {
			for mode := 0; mode < 9; mode++ {
				s := spellEffectsBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				p.Resources.Subject = subject
				tile := 0
				sp.Effects.Tile = &tile
				sp.Record = legacy.PortTestSpellLifeRecord{Words: map[int]uint32{4: 75, 8: 3, 52: math.Float32bits(612.25), 56: math.Float32bits(522.75), 68: s.Owner.Frame + 2}, Refs: map[int]int{16: 1}}
				switch mode {
				case 1:
					sp.Record.Refs[16] = 0
				case 2:
					a.ActorWords = map[int]uint32{16: 0x8020}
				case 3:
					sp.Record.Words[20] = 1
				case 4:
					s.Combat.Wall = 1
				case 5:
					sp.Effects.NullOutput = true
				case 6:
					s.Lifecycle.GameFlags |= 512
				case 7:
					sp.Record.Words[84] = 0x100
				case 8:
					tile = 1
				}
				cases = append(cases, s)
			}
		}
		spellEffectsHash(t, fmt.Sprintf("spell-effects-summon-admission-%02d", op), cases)
	}
}

func TestSpellEffectsForce(t *testing.T) {
	for _, op := range []int{39, 40, 41} {
		var cases []legacy.PortTestRoamSpec
		for _, distance := range []float32{0, 1, 4.9, 5, 9.9, 10, 10.1, 20} {
			for _, mass := range []float32{0.5, 1, 10} {
				for mode := 0; mode < 4; mode++ {
					s := spellEffectsBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					sp := a.Controls.SpellLifecycle
					a.ActorWords = map[int]uint32{56: math.Float32bits(512 + distance), 60: math.Float32bits(512), 88: 0, 92: 0, 120: math.Float32bits(mass), 172: 1, 16: 4}
					if mode == 1 {
						a.ActorWords[16] = 0x24
					}
					if mode == 2 {
						a.ActorWords[8] = 4 | 1<<22
					}
					if mode == 3 {
						s.Combat.Wall = 1
					}
					sp.Record = legacy.PortTestSpellLifeRecord{Words: map[int]uint32{0: math.Float32bits(512), 4: math.Float32bits(512)}}
					sp.Effects.RecordForce = true
					sp.Effects.Ints[1] = 17
					if op == 40 {
						p.TemporaryUpdates.Indexed = []int{1}
						sp.Effects.Floats = [3]uint32{math.Float32bits(10), math.Float32bits(5), math.Float32bits(2)}
					}
					if op == 41 {
						sp.Effects.RecordOutput = true
						sp.Effects.Output.Words = map[int]uint32{0: math.Float32bits(512), 4: math.Float32bits(512)}
						sp.Record.Words = map[int]uint32{4: math.Float32bits(10), 8: math.Float32bits(5), 12: math.Float32bits(20), 24: 17}
					}
					cases = append(cases, s)
				}
			}
		}
		r := spellEffectsHash(t, fmt.Sprintf("spell-effects-force-%02d", op), cases)
		if op != 39 {
			d := r[12].Callbacks.Shop.Sequence[0].TemporaryUpdatesData
			if d[len(d)-321] != 3 || d[len(d)-320] != 54000 || d[len(d)-318] != 17 {
				t.Errorf("op %d force callback arguments missing or incorrect", op)
			}
			u, _, _ := objectivePlayerData(r[12].Callbacks.Shop.Sequence[0], 0)
			if math.Float32frombits(u[22]) <= 0 {
				t.Errorf("op %d did not push a movable unit away from the source: %08x", op, u[22])
			}
		}
	}
}

func TestSpellEffectsCharm(t *testing.T) {
	for _, op := range []int{5, 6, 7, 8} {
		var cases []legacy.PortTestRoamSpec
		for _, size := range []byte{1, 2, 4} {
			for mode := 0; mode < 8; mode++ {
				s := spellEffectsBase(op)
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				a := o.Attack
				sp := a.Controls.SpellLifecycle
				p.Resources.Subject = 3
				a.ActorWords = map[int]uint32{56: math.Float32bits(520), 60: math.Float32bits(512), 12: 0, 16: 4, 340: 1 << 28, 172: 1}
				a.ActorRefs = map[int]int{508: 0}
				sp.Effects.ObjectTypes = map[int]string{1: "Bat"}
				sp.Effects.Guides = []legacy.PortTestSpellEffectGuide{{Index: 1, Size: size}}
				sp.Record = legacy.PortTestSpellLifeRecord{Words: map[int]uint32{4: 7, 8: 1, 52: math.Float32bits(520), 56: math.Float32bits(512), 68: s.Owner.Frame + 1}, Refs: map[int]int{16: 100, 48: 1}}
				o.PlayerDataWords[0][4248] = 1
				o.PlayerDataWords[0][3648] = 4
				p.TemporaryUpdates.Indexed = []int{1}
				switch mode {
				case 1:
					sp.Record.Refs[48] = 0
				case 2:
					a.ActorWords[16] = 0x8020
				case 3:
					a.ActorWords[56] = math.Float32bits(900)
				case 4:
					sp.Record.Words[68] = s.Owner.Frame + 2
				case 5:
					a.ActorWords[12] = 0x2000
				case 6:
					a.ActorRefs[508] = 100
				case 7:
					sp.Effects.CharmAll = true
					sp.Record.Words[20] = 1
				}
				cases = append(cases, s)
			}
		}
		r := spellEffectsHash(t, fmt.Sprintf("spell-effects-charm-%02d", op), cases)
		if op == 5 && r[0].Callbacks.Shop.Sequence[0].Return != 0 {
			t.Error("charm did not begin on an eligible creature")
		}
		if op == 6 {
			for _, i := range []int{0, 8, 16} {
				if r[i].Callbacks.Shop.Sequence[0].Return != 1 || r[i].Callbacks.Shop.Sequence[0].ResourceData[1][127] != 200 {
					t.Errorf("charm finish %d did not complete", i)
				}
			}
		}
	}
}

func TestSpellEffectsDoors(t *testing.T) {
	for _, op := range []int{25, 26, 27, 28} {
		var cases []legacy.PortTestRoamSpec
		for _, dir := range []uint32{0, 8, 16, 24} {
			for mode := 0; mode < 6; mode++ {
				s := spellEffectsBase(op)
				p := s.Callbacks.Shop
				tmp := p.TemporaryUpdates
				a := tmp.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				a.Actor = 3
				p.Items[0].Class = 0x80
				p.Items[0].Flags = 4
				tmp.ItemWords[0][56] = math.Float32bits(529)
				tmp.ItemWords[0][60] = math.Float32bits(506)
				tmp.ItemWords[0][172] = 1
				tmp.UpdateWords[0] = map[int]uint32{12: dir, 16: 23, 20: 22}
				tmp.ItemRefs[0] = map[int]int{508: 0}
				tmp.Indexed = []int{3, 4}
				p.Items[1].Class = 0x80
				p.Items[1].Flags = 4
				tmp.ItemWords[1][56] = math.Float32bits(552)
				tmp.ItemWords[1][60] = math.Float32bits(506)
				tmp.ItemWords[1][172] = 1
				tmp.UpdateWords[1] = map[int]uint32{12: dir, 16: 24, 20: 22}
				sp.Effects.Caches[2487704] = math.Float32bits(1e8)
				switch mode {
				case 1:
					tmp.ItemRefs[0][508] = 1
				case 2:
					tmp.ItemRefs[0][508] = 101
				case 3:
					s.Combat.Wall = 1
				case 4:
					p.Items[0].Class = 1
				case 5:
					tmp.Indexed = nil
				}
				cases = append(cases, s)
			}
		}
		r := spellEffectsHash(t, fmt.Sprintf("spell-effects-doors-%02d", op), cases)
		if op == 26 || op == 28 {
			for _, u := range r[0].Callbacks.Shop.Sequence[0].Objects[:2] {
				if u[128] != 54000 {
					t.Errorf("linked door did not acquire caster owner: %d", u[128])
				}
			}
		}
	}
}

func TestSpellEffectsPortals(t *testing.T) {
	for _, op := range []int{21, 22} {
		var cases []legacy.PortTestRoamSpec
		for slot := 0; slot < 4; slot++ {
			for mode := 0; mode < 4; mode++ {
				s := spellEffectsBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				a.Controls.X = int32(46 + slot)
				a.UpdateRefs = map[int]int{}
				for j := 0; j < slot; j++ {
					a.UpdateRefs[116+4*j] = 3 + j
				}
				if mode == 1 {
					a.UpdateRefs[116+4*slot] = 4
				}
				if mode == 2 {
					sp.Effects.Args[2] = 3
					sp.Effects.ObjectTypes = map[int]string{3: "Glyph"}
				}
				if mode == 3 {
					a.MissingTypes = []string{fmt.Sprintf("TeleportGlyph%d", slot+1)}
				}
				for j := 0; j < 3; j++ {
					p.TemporaryUpdates.ItemWords[j][136] = uint32(j + 1)
				}
				cases = append(cases, s)
			}
		}
		spellEffectsHash(t, fmt.Sprintf("spell-effects-portals-%02d", op), cases)
	}
}

func TestSpellEffectsAreaCasts(t *testing.T) {
	for _, op := range []int{9, 12, 13, 37, 38} {
		var cases []legacy.PortTestRoamSpec
		for _, distance := range []float32{0, 1, 29, 30, 31, 599, 600, 601} {
			for level := int32(1); level <= 5; level++ {
				for mode := 0; mode < 2; mode++ {
					s := spellEffectsBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					sp := a.Controls.SpellLifecycle
					a.Actor = 2
					a.ActorWords = map[int]uint32{56: math.Float32bits(512 + distance), 60: math.Float32bits(512), 88: 0, 92: 0, 120: math.Float32bits(10), 172: 1, 16: 4}
					if mode == 1 {
						a.ActorWords[16] = 0x4004
					}
					p.TemporaryUpdates.Indexed = []int{2}
					sp.Z = level
					sp.Effects.Caches[2487700] = uint32(level)
					p.EffectsUse.Balance["EarthquakeJiggle"] = []float64{1, 2, 3, 4, 5}
					cases = append(cases, s)
				}
			}
		}
		spellEffectsHash(t, fmt.Sprintf("spell-effects-area-%02d", op), cases)
	}
}

func TestSpellEffectsLesserHealAndFumble(t *testing.T) {
	for _, op := range []int{14, 36} {
		var cases []legacy.PortTestRoamSpec
		for _, subject := range []int{1, 2, 3} {
			for _, hp := range []uint16{0, 1, 99, 100} {
				for mode := 0; mode < 4; mode++ {
					s := spellEffectsBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					sp := a.Controls.SpellLifecycle
					p.Resources.Subject = subject
					p.Resources.HP = hp
					p.Resources.MaxHP = 100
					p.Resources.PlayerClass = byte(mode % 3)
					sp.Record.Refs[0] = 1
					p.Inventory.Linked = []int{0, 1, 2}
					p.Inventory.Owned = []int{0, 1, 2}
					if mode == 1 {
						p.Resources.Subclass = 0x10
					}
					if mode == 2 {
						p.Resources.Subclass = 0x2000
					}
					if mode == 3 {
						sp.Record.Refs[0] = 0
					}
					p.EffectsUse.Balance["LesserHealAmount"] = []float64{0.125}
					cases = append(cases, s)
				}
			}
		}
		spellEffectsHash(t, fmt.Sprintf("spell-effects-heal-fumble-%02d", op), cases)
	}
}

func TestSpellEffectsInversion(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, distance := range []float32{0, 1, 29.999, 30, 30.001, 60} {
		for mode := 0; mode < 4; mode++ {
			s := spellEffectsBase(9)
			p := s.Callbacks.Shop
			tmp := p.TemporaryUpdates
			p.Items[0].Class = uint32(object.ClassMissile)
			p.Items[0].Subclass = uint32(object.MissileMagic)
			p.Items[0].Flags = 4
			tmp.ItemWords[0][56] = math.Float32bits(512 + distance)
			tmp.ItemWords[0][60] = math.Float32bits(512)
			tmp.ItemWords[0][172] = 1
			tmp.ItemRefs[0] = map[int]int{508: 101}
			tmp.UpdateRefs[0] = map[int]int{0: 101, 4: 1}
			tmp.UpdateWords[0] = map[int]uint32{12: 1}
			tmp.Indexed = []int{3}
			switch mode {
			case 1:
				p.Items[0].Subclass = 0
			case 2:
				tmp.UpdateRefs[0][4] = 101
			case 3:
				tmp.Indexed = nil
			}
			cases = append(cases, s)
		}
	}
	r := spellEffectsHash(t, "spell-effects-inversion", cases)
	if r[0].Callbacks.Shop.Sequence[0].Objects[0][128] != 54000 {
		t.Error("inversion did not transfer missile ownership")
	}
}

func TestSpellEffectsGlyphSelection(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for nearest := 0; nearest < 3; nearest++ {
		for mode := 0; mode < 4; mode++ {
			s := spellEffectsBase(23)
			p := s.Callbacks.Shop
			tmp := p.TemporaryUpdates
			sp := tmp.World.Objectives.Attack.Controls.SpellLifecycle
			sp.Effects.ObjectTypes = map[int]string{3: "Glyph", 4: "Glyph", 5: "Glyph"}
			for i := 0; i < 3; i++ {
				tmp.ItemWords[i][56] = math.Float32bits(550 + float32(i))
				tmp.ItemWords[i][60] = math.Float32bits(512)
				tmp.ItemRefs[i] = map[int]int{508: 1}
			}
			tmp.ItemWords[nearest][56] = math.Float32bits(513)
			switch mode {
			case 1:
				tmp.ItemRefs[nearest][508] = 101
			case 2:
				sp.Effects.ObjectTypes[3+nearest] = "NPC"
			case 3:
				tmp.World.Objectives.ObjectList = nil
			}
			cases = append(cases, s)
		}
	}
	r := spellEffectsHash(t, "spell-effects-glyph-selection", cases)
	for _, i := range []int{0, 4, 8} {
		d := r[i].Callbacks.Shop.Sequence[0].TemporaryUpdatesData
		if d[len(d)-319] != uint32(70000+i/4) {
			t.Errorf("nearest owned glyph %d was not selected", i/4)
		}
	}
}

func TestSpellEffectsSummonCapacity(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, size := range []byte{1, 2, 4} {
		for _, cl := range []object.MonsterClass{object.MonsterSmall, object.MonsterMedium, object.MonsterLarge} {
			s := spellEffectsBase(1)
			p := s.Callbacks.Shop
			o := p.TemporaryUpdates.World.Objectives
			a := o.Attack
			sp := a.Controls.SpellLifecycle
			p.Resources.Subject = 3
			a.ActorWords = map[int]uint32{12: uint32(cl)}
			a.UpdateWords = map[int]uint32{1440: 0x80}
			a.ActorRefs = map[int]int{508: 100}
			o.PlayerRefs[0] = map[int]int{516: 1}
			tile := 0
			sp.Effects.Tile = &tile
			sp.Effects.Guides = []legacy.PortTestSpellEffectGuide{{Index: 1, Size: size}}
			sp.Record = legacy.PortTestSpellLifeRecord{Words: map[int]uint32{4: 75, 8: 1, 52: math.Float32bits(520), 56: math.Float32bits(515)}, Refs: map[int]int{16: 100}}
			cases = append(cases, s)
		}
	}
	r := spellEffectsHash(t, "spell-effects-summon-capacity", cases)
	for i, x := range r {
		want := uint32(0)
		if []int{1, 2, 4}[i/3]+[]int{1, 2, 4}[i%3] > 4 {
			want = 1
		}
		if x.Callbacks.Shop.Sequence[0].Return != want {
			t.Errorf("capacity case %d return %d want %d", i, x.Callbacks.Shop.Sequence[0].Return, want)
		}
	}
}

func TestSpellEffectsFactoryFailures(t *testing.T) {
	names := []string{"MediumFlame", "Fireball", "StrongFireball", "TitanFireball", "TelekinesisHand", "SmallFist", "MediumFist", "LargeFist", "FlameCleanse", "BlueFlameCleanse", "SmallFlameCleanse", "SmallBlueFlameCleanse", "MediumFlameCleanse", "MediumBlueFlameCleanse", "LargeFlameCleanse", "LargeBlueFlameCleanse", "MeteorShower", "Meteor", "ToxicCloud", "ArachnaphobiaFocus"}
	for _, op := range []int{20, 29} {
		var cases []legacy.PortTestRoamSpec
		for _, level := range []int32{1, 3, 5} {
			s := spellEffectsBase(op)
			a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
			a.MissingTypes = names
			a.Controls.SpellLifecycle.Z = level
			cases = append(cases, s)
		}
		r := spellEffectsHash(t, fmt.Sprintf("spell-effects-factory-failure-%02d", op), cases)
		for i, x := range r {
			if len(x.Lifecycle.Created) != 0 {
				t.Errorf("op %d level %d created an unavailable type", op, i)
			}
		}
	}
}

func TestSpellEffectsSummonCost(t *testing.T) {
	ids := []int{75, 76, 88, 100, 114}
	values := []uint32{0, 1, 200, 0x80000000, 0xffffffff}
	saved := make([]uint32, len(ids))
	for i, id := range ids {
		ptr := memmap.PtrUint32(0x587000, 217668+uintptr(4*id))
		saved[i] = *ptr
		*ptr = values[i]
	}
	defer func() {
		for i, id := range ids {
			*memmap.PtrUint32(0x587000, 217668+uintptr(4*id)) = saved[i]
		}
	}()
	var cases []legacy.PortTestRoamSpec
	for _, id := range ids {
		for _, subject := range []int{1, 3} {
			s := spellEffectsBase(0)
			s.Callbacks.Shop.Resources.Subject = subject
			s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.X = int32(id)
			cases = append(cases, s)
		}
	}
	r := spellEffectsHash(t, "spell-effects-summon-cost", cases)
	for i, x := range r {
		want := uint32(0)
		if i%2 == 0 {
			want = values[i/2]
		}
		if x.Callbacks.Shop.Sequence[0].Return != want {
			t.Errorf("cost case %d got %08x want %08x", i, x.Callbacks.Shop.Sequence[0].Return, want)
		}
	}
}
