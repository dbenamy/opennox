//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
)

var sustainedHashes = map[string]string{
	"sustained-duration-00":     "30e2f42703d8ddc5c2e2874239bd42fd2042f7383d570ab8384a5b664753f337",
	"sustained-duration-05":     "d831d98b422a72b62b05e8954dfd6b5c41357d392734cefd4ba0e967c375c890",
	"sustained-duration-06":     "c8f0f7ae25f1dbe7b034fb58d278d0e441f22a062dd1a247e7bad7940b0c72e6",
	"sustained-duration-08":     "d5edcdb2c34be6808c3a2c7c9f4e04fb6e65dbae5ad44638b17ba20ca822f646",
	"sustained-duration-09":     "04e4c8720ce26bcba2cbf8562dd20d5d82e53a892231379e6fc4aaca197cc83b",
	"sustained-duration-10":     "768fdd7072a49972c2c8ad88811e471d72faa736136196f534b6baf97683b595",
	"sustained-duration-11":     "f392d74ba311c144ed37f02d910f798cbae923a23592554d80dcb25238574baf",
	"sustained-duration-12":     "7648d7a334054815a4ee278a81cc15c3c9ce4f4fbacd5d170e92bd92aa37b972",
	"sustained-duration-13":     "3694ab5faa33cdf298a19240b658e9f8da0ad5f7db72e94082a58ba68709fb27",
	"sustained-duration-14":     "0d7a706d6ecd2a2c25791faa265af84af789c451d9c7e4adbae3d3e410319afb",
	"sustained-duration-15":     "7ac3ba26554ad2951be5bc8e4003b94ea0f473cc16148594be3b1dd5d560f12e",
	"sustained-duration-16":     "0585d81b636eaa8863fc90dce7d2711860dc49cff34a1c57c51e1212df3dd416",
	"sustained-duration-17":     "68d0cc3269be274fc85b9e807e0bd992221f15951efcecdddcefc35c82999ed5",
	"sustained-duration-20":     "37b63d539deda34d94e61a692ee327093e160d0131f9fd828f26e2922b1dbc19",
	"sustained-duration-21":     "c3ab9baed04fa38aa653eb9fa83f9f62d0d8a47d52aa2ddf24651cf282664997",
	"sustained-duration-24":     "f392d74ba311c144ed37f02d910f798cbae923a23592554d80dcb25238574baf",
	"sustained-duration-25":     "f5a78743baf48f840e0a7bcc1f2255bcbbd209c4e959032ffbe2956fac6585b5",
	"sustained-duration-26":     "0585d81b636eaa8863fc90dce7d2711860dc49cff34a1c57c51e1212df3dd416",
	"sustained-duration-27":     "2a8ef0d450f15bea5f693951797c14ab1d6d4cf721660be678a5809b133f6642",
	"sustained-duration-28":     "98197a5327e9c0a1736fa6bee0294c96a04ba1b53b73319bd74bebae542e07ce",
	"sustained-duration-29":     "9a27ed520d5d313f02b03d0874b5b34c86d62a7202d7d47ed4afc30af4910a42",
	"sustained-duration-31":     "77de2fffe9372e45231bf614e5bcb46087b8277c5bcdb19ea3f873f49aee4b3f",
	"sustained-duration-32":     "79bd4b31ac19ccc787334f988a25b8c22ab57be3e60e33aaca6b2afe44d3d038",
	"sustained-duration-33":     "98197a5327e9c0a1736fa6bee0294c96a04ba1b53b73319bd74bebae542e07ce",
	"sustained-duration-34":     "4cce34cb9c3f809013f89b83f854b31ea77df2a97a59637be03c03b1ed7aad84",
	"sustained-duration-35":     "3bcb8507644e417bed2e55d50023bb3de9b4dabf90ed2d613408c56c72620907",
	"sustained-duration-36":     "39930f3f8728d4c58e8218b7b30402a243864650ae119144148ee69924c8e5c3",
	"sustained-duration-37":     "b0e35adad22c8bb1282ef386eed66e22e998b3cac732ca39efff085f9324bd39",
	"sustained-duration-38":     "3496ebf3d89d2b888842d99f67285b559a99d09febdd979422b3d448d8f50766",
	"sustained-duration-39":     "668053fa4ff346b97c86fa285ba13e69dad4b9fb7144e970568c46ea8114032e",
	"sustained-duration-40":     "6e61faf90851a5247e9cfeaddd49fe26fa7d8ce80894ede9a315fa9262aeeaec",
	"sustained-duration-41":     "246100c0f8aedc927e0b2be3dcc945d27dddf7d28aa61878bd938c76c0b95a00",
	"sustained-duration-43":     "274633acafa9c7669c160135f3ed013574a7137125ea8e73aa7a862cf18d9d24",
	"sustained-duration-44":     "cdce7a0edb613e5e5c1919a6a31258422ddf971146091355a5cf7d140b36c073",
	"sustained-duration-45":     "0585d81b636eaa8863fc90dce7d2711860dc49cff34a1c57c51e1212df3dd416",
	"sustained-duration-46":     "30fc4c14c8cc6fea07fd3d77f4c45f88f975a265ee46821a75747544e9d1d617",
	"sustained-duration-47":     "84c506edfafeeec0b99a38906720aef622deac037ce4b94f40b6fb80c978ba71",
	"sustained-duration-48":     "f142c739fab60ae2611bdaf3867b0f21b97076c2ec84bb1b421432be6cfdd2ea",
	"sustained-duration-50":     "f392d74ba311c144ed37f02d910f798cbae923a23592554d80dcb25238574baf",
	"sustained-duration-51":     "ab56e9c96f8aacaf8fd031133439aec82bca1cc70afb254eb21b0aee0c5b7be8",
	"sustained-duration-52":     "d07dcbc89e9243adaae339368d02a4d328d17f1e5445fba2dc43c824a982fd3a",
	"sustained-factory-failure": "ac5c65fae982bc7d33bf4990fe5c0a9a452ff950508f7afa01c39e19bc80b18a",
	"sustained-firewalk":        "e5d3c799e86713f43b2bf7194424874cdcd947099ee16b8be38a963ad51cc02f",
	"sustained-heal-channel-13": "17a68de74cc2cb096aa062b23504b5533b490e3f0c97ac03bbc264253b92945d",
	"sustained-heal-channel-14": "64710910ba6a77413aae68377039d6bf13ef399b102eb5478d3be5962140ea0d",
	"sustained-lightning-chain": "c564fb856853dfe9c4202e9e3b24c86ce5907cf0b7d2db0d2712cdeb377e8736",
	"sustained-mana-bomb":       "94241c8bce29d1ba094da4e68e66bcf5a974f3a2ad53bd9ff0df78e634198f3a",
	"sustained-mana-transfer":   "5b85e06f21b95a5d5735717e3d3cc839b36c39f22aea082004aa6a3537b52189",
	"sustained-nil-casters":     "3bb7720971aa3154d72c8493fcfacb7aa2a7b7ce17c02c8df18b833e3db18e95",
	"sustained-positive":        "9c1154f299b1c7be04b53adb377f30172be8dda8e7ae78a9e27e327a56def820",
	"sustained-shield-18":       "a2974629ab560c50de41b532fb57f39c5c1e04d1964a9cd2812c7c6ba998c48c",
	"sustained-shield-19":       "d7f7aa36d3df41eb9721035ce5da320fd6c95236617b7027decc0e6371791988",
	"sustained-smoke-00":        "8fcdfae28b138d1398707c0cfd3aa1c2a48981047011d448de8f3ed7eafa51ca",
	"sustained-smoke-01":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-02":        "cfbd70413b73e93241513aed3ad319ccc1512ff8510722e08e0bb31334c84466",
	"sustained-smoke-03":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-04":        "8fcdfae28b138d1398707c0cfd3aa1c2a48981047011d448de8f3ed7eafa51ca",
	"sustained-smoke-05":        "18c98627d986f40a383f35031339a467edf5a2616c339a3cd224ec753c9b5634",
	"sustained-smoke-06":        "dd8269fcf837b4476d183c3244fd56ea4df77ab8a7b39690e9846b2daa6f2852",
	"sustained-smoke-07":        "6ce9c961ef0b9260679778a7c56e4aed637b6f06eb07fdfc89515ae6f03842fa",
	"sustained-smoke-08":        "5e4447d1c0cafb0331df2c348144eaa83fc52c6518ea5d45be63c9f700218724",
	"sustained-smoke-09":        "9a458e6c47b221dd05019daa649d7153ea6d6d127c652a1db25b2a232347911d",
	"sustained-smoke-10":        "e7fce86cd609aa6e95475f4d6484ae3f5415d896390978f936de9cffe8fdd404",
	"sustained-smoke-11":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-12":        "3b229a262657e853f02dc352071ecd569b18d0444117d8cc7fe365b0e6ff347e",
	"sustained-smoke-13":        "b9f42b54fcf36b15786b66bae66db28bceddd7bb272b0399fe8e16845a6c955d",
	"sustained-smoke-14":        "8fcdfae28b138d1398707c0cfd3aa1c2a48981047011d448de8f3ed7eafa51ca",
	"sustained-smoke-15":        "13f1cdf9639c022515d3b1dbf7ab20d1ff4644a7cfeb79313d034adeb2ff1ce5",
	"sustained-smoke-16":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-17":        "e281e132401ab021b3afcb181faede4b7203bfbd73fec1766a6b980f610927a2",
	"sustained-smoke-18":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-19":        "1cba29dd71c5764de34614c88fb4f944314302c81a71368d05e29ecaf20c1aa5",
	"sustained-smoke-20":        "48ef2c09a372c43da1d38eaa20eb4afc6bbecc592b1e54865188e115b386d309",
	"sustained-smoke-21":        "2110f8480bb4c024023ae371b32e3d81239cf31b94fda4456e730713e66766f1",
	"sustained-smoke-22":        "6ce9c961ef0b9260679778a7c56e4aed637b6f06eb07fdfc89515ae6f03842fa",
	"sustained-smoke-23":        "721d11abd24ba20cfdb5f45b1a9819202f7cf4ea6212dc55747f8e16fea387b5",
	"sustained-smoke-24":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-25":        "59e9ed849ca3b61d0ae714ff71d04511bc5be8a4304a9d1ce35f049c700e3135",
	"sustained-smoke-26":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-27":        "8138d5d6b4e8f843814a3356bb51d9889d839d182616c1e1069b448adbf79cf4",
	"sustained-smoke-28":        "7576a73474433f3d3b41aebde30a6055d93057f25badc71bc9474f4fa2c2f659",
	"sustained-smoke-29":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-30":        "e515b22919b2261f4bfa8cbc0acd9e995c70eeb18778c125fb1b533c4eaee4d5",
	"sustained-smoke-31":        "aeb43c312a288142f2c30bc3ce11c122147ace987608bcc7b1bcb60aed8f2971",
	"sustained-smoke-32":        "aeb43c312a288142f2c30bc3ce11c122147ace987608bcc7b1bcb60aed8f2971",
	"sustained-smoke-33":        "7576a73474433f3d3b41aebde30a6055d93057f25badc71bc9474f4fa2c2f659",
	"sustained-smoke-34":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-35":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-36":        "7576a73474433f3d3b41aebde30a6055d93057f25badc71bc9474f4fa2c2f659",
	"sustained-smoke-37":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-38":        "03d75690fd94ef692ea43f2953260ecb2d6bc126f9e6824770acce88ce6aaa22",
	"sustained-smoke-39":        "09970a30f0c6f22dd6239ee75b4d18b7ab299225e2482f6dadacd1722c187e69",
	"sustained-smoke-40":        "9324559d81b486c24d957161a307f4d3d1ac2059a8368dcb2db852f401e6ecde",
	"sustained-smoke-41":        "4b41dd24fe4b9288b59331faa94eefc0f98cb678f598c758bdbc1327635f3f4d",
	"sustained-smoke-42":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-43":        "06d062632ff552191888ca691cf0ec7d19f2357c7f6c0e429efd7cbdbf29f801",
	"sustained-smoke-44":        "12771f94f72f43af4b9dd478a626cd7118c6372612207afcf3e8532567413840",
	"sustained-smoke-45":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-46":        "d728e66a293a44b31904af510f2924d7ddbf98cea8ca39e70ff4d15ec38eed47",
	"sustained-smoke-47":        "901b653e3a6c3a0c61189a6a2faa2c42bf1336baecd50f2969f7d3f4c6a44c0f",
	"sustained-smoke-48":        "bf63f6e72860ca2df8215236afdea6ca1d5ba35ef4ff08b575e7b289fceaec1c",
	"sustained-smoke-49":        "6ce9c961ef0b9260679778a7c56e4aed637b6f06eb07fdfc89515ae6f03842fa",
	"sustained-smoke-50":        "97b6b9b0902f688a9baef2b15d131349edd28575980f105661d5340e5ef4b183",
	"sustained-smoke-51":        "9fd040e329a2339ff4dc57d8abfb09b92be9504e5eaf2765a2c2c230b72a784a",
	"sustained-smoke-52":        "05fbe5d5ba89558786e693021ae2df8bcc6873be26f460155bc7f153c8c9433b",
	"sustained-spatial-00":      "dc7800e2a2c0103f59f01312bdf0417cd598131677a24c69d60ae8e0ea9c8881",
	"sustained-spatial-02":      "f10d0efbc0a35b0bed5161d3a244e8b52510fe2d190185e9d6deddd36ff5187d",
	"sustained-spatial-03":      "411a06f30b251b6298767e96fb085f51b093e00168fd0e2f14e89397fd008f75",
	"sustained-spatial-04":      "2737e09893f611f0e479f77fef9ef8756eac481cf96adf25e9403b9ce6e78ab1",
	"sustained-spatial-06":      "d161dc0079a9fe4e7996af88ad0bd87a435c98801e90e3933126934800daf840",
	"sustained-spatial-07":      "665d37a2b1946c889977115eff2c10d156c42355f96e50939c68f50d7db32fe9",
	"sustained-spatial-21":      "67164dd020ef5293ab1758a88839ac003f7b72c18122a38aa478ffec1d2bb586",
	"sustained-spatial-22":      "ac4c24364c9e4f95914893820604b1af9d8c0c4d78c63c9b248d02186ca37cc2",
	"sustained-spatial-23":      "080930914f24baceccd22a07064e2ea414a408e012dfc4373704b2d5368039b2",
	"sustained-spatial-48":      "594165872a8c98ce52be8a1ff01a2fa3d3eb955661ad5cbdb5a8270e90ff9f9a",
	"sustained-spatial-49":      "e36f867b16a2113c485575b169b83a813c87d7d71289a1a07b64722080bf22c4",
	"sustained-tag-nil":         "ec1a85abab72abba075434d5a4bb23d438a088c273b68b5fe1b71d5f4b201e89",
	"sustained-teleport-31":     "23e07501268bb0fbf39788fb8c75bcce1b36ce3fe549a5affc7aac4cd4384cf4",
	"sustained-teleport-32":     "d91c309df16dc76e2e09af3424354299c5361f95074fbbb927c663879b799d1a",
	"sustained-teleport-34":     "a75770fb1ae072c2fe7f057463ae20aa810b841b5666910348c275c510f029e2",
	"sustained-teleport-35":     "7bc1cee177d284a3c19131ca11545d9c4d742905508cf7fde6ef843b396c8fc0",
	"sustained-teleport-37":     "d726f32d0713496b05278c7548d06fac0ef63ee1f8818717f6389581cd600fc1",
	"sustained-wand-09":         "2a86b99e5a2fe63768cccdf247b9a1095b5ce7817164cb677b1be0574e389ed4",
	"sustained-wand-11":         "8a28ad294c8ecbb262dafb11a32e7f4ba5f49118aa921082297dab3ee511824f",
	"sustained-wand-20":         "4b1ef33be73a83f96428e3836cd80aeecc64cc8ea76dcc9a01ef95a8b82d8768",
	"sustained-wand-24":         "a482ecd6732cdd6bb7854b8707e84f796cd080e8b742ddfa1dac2eafd47135be",
	"sustained-wand-47":         "5f554e9c5866c2f64e3b3373b3ff5641ca421a5039843d4dad86acb6b44a3faa",
	"sustained-wand-48":         "9f623d4feec7dc79ccba2c55569443f4754121030060701cbf35b42b00eb50e5",
	"sustained-wand-50":         "c71b836ba00bcd43e0bc8b778a2a58d999d33236e111f59491bbe276992aab2b",
}

func sustainedBase(op int) legacy.PortTestRoamSpec {
	s := spellEffectsBase(0)
	s.Owner.Frame = 100
	s.Spells.TargetCur = 100
	s.Spells.TargetMax = 100
	p := s.Callbacks.Shop
	a := p.TemporaryUpdates.World.Objectives.Attack
	sp := a.Controls.SpellLifecycle
	sp.Effects.Sustained = &legacy.PortTestSustainedSpellsSpec{Caches: map[uintptr]uint32{}, HealTable: [5]float32{0.25, 0.5, 1, 2, 4}}
	sp.Effects.Args[0] = 2
	sp.Effects.Ints[0] = 7
	sp.Effects.Output.Words = map[int]uint32{0: 20, 4: math.Float32bits(520)}
	sp.Record = legacy.PortTestSpellLifeRecord{Words: map[int]uint32{4: 24, 8: 3, 28: math.Float32bits(512), 32: math.Float32bits(512), 52: math.Float32bits(520), 56: math.Float32bits(512), 60: 100, 64: 100, 68: 200}, Refs: map[int]int{12: 1, 16: 1, 24: 1, 48: 2}}
	p.Sequence = []legacy.PortTestShopAction{{Op: 1700 + op}}
	for _, name := range []string{"ManaDrainRange", "LightningRange", "LightningSearchTime", "PlasmaSearchTime", "TagDurationPerLevel", "MoonglowEnchantmentDuration", "ManaBombGlyphDuration", "ManaBombInRadius", "ManaBombOutRadius", "ManaBombShakeMag", "PlasmaDamage", "PlasmaDamageHecubah"} {
		p.EffectsUse.Balance[name] = []float64{30}
	}
	for _, name := range []string{"ManaDrainCoeff", "EnergyBoltDamage", "EnergyBoltGlyphDamage", "LightningDamage", "LightningGlyphDamage", "ChannelLifeCoeff", "ShieldDuration", "ShieldHealth", "TeleportDelay", "ManaBombInitPower", "ManaBombDeltaPower", "TurnUndeadKillPoints"} {
		p.EffectsUse.Balance[name] = []float64{1, 2.5, 5, 10, 20}
	}
	return s
}
func TestSustainedSpellsSmoke(t *testing.T) {
	for op := 0; op < 53; op++ {
		t.Run(fmt.Sprintf("%02d", op), func(t *testing.T) {
			s := sustainedBase(op)
			sp := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
			if op == 2 || op == 30 {
				sp.Record.Refs = nil
				sp.Record.Words = map[int]uint32{0: math.Float32bits(512), 4: math.Float32bits(512)}
				sp.Effects.Output.Words = map[int]uint32{0: math.Float32bits(520), 4: math.Float32bits(520)}
			}
			if op == 31 || op == 32 {
				sp.Record.Words[4] = 122
			}
			sustainedHash(t, fmt.Sprintf("sustained-smoke-%02d", op), []legacy.PortTestRoamSpec{s})
		})
	}
}
func sustainedHash(t *testing.T, name string, cases []legacy.PortTestRoamSpec) []legacy.PortTestRoamResult {
	// Actual startup data from blob_587000.dat: one through five lightning targets.
	var priorCounts [5]uint32
	for i := range priorCounts {
		ptr := memmap.PtrUint32(0x587000, uintptr(260384+4*i))
		priorCounts[i] = *ptr
		*ptr = uint32(i + 1)
	}
	priorWalls := doDamageWalls
	defer func() {
		for i, v := range priorCounts {
			*memmap.PtrUint32(0x587000, uintptr(260384+4*i)) = v
		}
		doDamageWalls = priorWalls
	}()

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
	callbackHash(t, name, r, sustainedHashes[name])
	return r
}

func TestSustainedSpellsTagNilCaster(t *testing.T) {
	s := sustainedBase(25)
	s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle.Record.Refs[16] = 0
	r := sustainedHash(t, "sustained-tag-nil", []legacy.PortTestRoamSpec{s})
	if r[0].Callbacks.Shop.Sequence[0].Return != 1 {
		t.Fatal("nil caster must reject tagging")
	}
}

func TestSustainedSpellsDurationBoundaries(t *testing.T) {
	for _, op := range []int{0, 5, 6, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 20, 21, 24, 25, 26, 27, 28, 29, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 43, 44, 45, 46, 47, 48, 50, 51, 52} {
		var cases []legacy.PortTestRoamSpec
		for _, level := range []uint32{1, 3, 5} {
			for mode := 0; mode < 6; mode++ {
				s := sustainedBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				sp.Record.Words[8] = level
				if op == 31 || op == 32 {
					sp.Record.Words[4] = 122
				}
				switch mode {
				case 1:
					sp.Record.Refs[48] = 0
				case 2:
					s.Spells.TargetFlags = 0x8024
				case 3:
					a.ActorWords = map[int]uint32{340: 1 << 8}
				case 4:
					sp.Record.Words[20] = 1
				case 5:
					sp.Record.Words[68] = 101
				}
				cases = append(cases, s)
			}
		}
		sustainedHash(t, fmt.Sprintf("sustained-duration-%02d", op), cases)
	}
}
func TestSustainedSpellsPositiveContracts(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{9, 15, 30, 38, 41, 44, 51} {
		s := sustainedBase(op)
		sp := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
		sp.Record.Refs[48] = 1
		if op == 30 {
			sp.Record.Refs = nil
			sp.Record.Words = map[int]uint32{0: math.Float32bits(512), 4: math.Float32bits(512)}
			sp.Effects.Output.Words = map[int]uint32{0: math.Float32bits(520), 4: math.Float32bits(520)}
		}
		cases = append(cases, s)
	}
	r := sustainedHash(t, "sustained-positive", cases)
	for _, v := range []struct{ i, n int }{{0, 1}, {2, 1}, {3, 1}, {4, 43}, {6, 1}} {
		if got := len(r[v.i].Lifecycle.Created); got != v.n {
			t.Errorf("positive %d created %d want %d", v.i, got, v.n)
		}
	}
	for _, v := range []struct{ i, bit int }{{1, 26}, {5, 27}, {6, 1}} {
		if u := r[v.i].Callbacks.Shop.Sequence[0].ResourceData[1]; u[85]&(1<<uint(v.bit)) == 0 {
			t.Errorf("positive %d missing enchant %d", v.i, v.bit)
		}
	}
}

func TestSustainedSpellsManaTransfer(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, kind := range []int{0, 2} {
		for _, amount := range []int32{-1, 0, 1, 7, 40, 41, 65535} {
			for _, flags := range []uint32{0, 4096} {
				s := sustainedBase(1)
				p := s.Callbacks.Shop
				o := p.TemporaryUpdates.World.Objectives
				a := o.Attack
				sp := a.Controls.SpellLifecycle
				p.Resources.Mana = 100
				p.Resources.OldMana = 100
				sp.Effects.Ints[0] = amount
				s.Lifecycle.GameFlags |= flags
				o.PlayerUpdateWords[1] = map[int]uint32{4: 40 | 40<<16, 8: 200}
				sp.Effects.Args[0] = 101
				if kind == 2 {
					sp.Effects.Args[0] = 3
					p.Items[0].Class = 1 << 22
					p.Items[0].Subclass = 0x18
					p.TemporaryUpdates.UpdateWords[0] = map[int]uint32{0: 40}
				}
				cases = append(cases, s)
			}
		}
	}
	r := sustainedHash(t, "sustained-mana-transfer", cases)
	// A positive player transfer preserves total mana and applies the requested seven.
	for i, s := range cases {
		sp := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
		if sp.Effects.Ints[0] == 7 && s.Lifecycle.GameFlags&4096 == 0 && sp.Effects.Args[0] == 101 {
			_, source, _ := objectivePlayerData(r[i].Callbacks.Shop.Sequence[0], 0)
			_, donor, _ := objectivePlayerData(r[i].Callbacks.Shop.Sequence[0], 1)
			if source[1]&65535 != 107 || donor[1]&65535 != 33 {
				t.Fatalf("mana transfer source=%d donor=%d", source[1]&65535, donor[1]&65535)
			}
		}
	}
}
func TestSustainedSpellsShieldDamage(t *testing.T) {
	for _, op := range []int{18, 19} {
		var cases []legacy.PortTestRoamSpec
		for _, hp := range []uint16{1, 2, 20, 100, 65535} {
			for _, damage := range []int32{-1, 0, 1, 2, 19, 20, 100} {
				for mode := 0; mode < 3; mode++ {
					s := sustainedBase(op)
					p := s.Callbacks.Shop
					a := p.TemporaryUpdates.World.Objectives.Attack
					sp := a.Controls.SpellLifecycle
					p.Resources.HP = hp
					p.Resources.MaxHP = 65535
					sp.Effects.Ints[0] = damage
					sp.Effects.Output.Words[0] = uint32(damage)
					if mode != 0 {
						sp.Durations = []legacy.PortTestSpellLifeRecord{{Words: map[int]uint32{4: 51, 8: 3, 72: 30}, Refs: map[int]int{16: 1, 48: 1}}}
					}
					if mode == 2 {
						a.ActorWords = map[int]uint32{16: 0x8024}
					}
					cases = append(cases, s)
				}
			}
		}
		sustainedHash(t, fmt.Sprintf("sustained-shield-%02d", op), cases)
	}
}

func TestSustainedSpellsSpatial(t *testing.T) {
	for _, op := range []int{0, 2, 3, 4, 6, 7, 21, 22, 23, 48, 49} {
		var cases []legacy.PortTestRoamSpec
		for _, dist := range []float32{0, 1, 29.99, 30, 30.01, 100, 400, 401} {
			for _, level := range []uint32{1, 3, 5} {
				for mode := 0; mode < 3; mode++ {
					s := sustainedBase(op)
					p := s.Callbacks.Shop
					o := p.TemporaryUpdates.World.Objectives
					a := o.Attack
					sp := a.Controls.SpellLifecycle
					p.Resources.Mana = 100
					p.Resources.OldMana = 100
					o.PlayerWords[0] = map[int]uint32{16: 4, 56: math.Float32bits(512), 60: math.Float32bits(512), 124: 0}
					s.Combat.Target = [2]uint32{math.Float32bits(512 + dist), math.Float32bits(512)}
					s.Spells.TargetFlags = 4
					s.Spells.TargetCur = 500
					s.Spells.TargetMax = 500
					if op != 2 {
						a.Actor = 2
						sp.ActorType = ""
						a.UpdateWords = map[int]uint32{4: 50 | 50<<16, 8: 100}
					}
					p.TemporaryUpdates.Indexed = []int{2}
					sp.Record.Words[8] = level
					sp.Effects.Sustained.Caches = map[uintptr]uint32{2487836: math.Float32bits(512), 2487840: math.Float32bits(512), 2487876: math.Float32bits(1e9), 2487820: math.Float32bits(512), 2487824: math.Float32bits(512)}
					sp.Effects.Sustained.Globals[0] = math.Float32bits(1e9)
					sp.Effects.Sustained.Globals[5] = math.Float32bits(1e9)
					sp.Effects.Sustained.GlobalRefs = map[int]int{2: 1}
					if op == 3 || op == 4 || op == 7 || op == 22 || op == 23 || op == 49 {
						a.Actor = 2
						sp.Effects.Args[0] = 1
						sp.ActorType = ""
					}
					if op == 2 {
						sp.Record.Refs = nil
						sp.Record.Words = map[int]uint32{0: math.Float32bits(512), 4: math.Float32bits(512)}
					}
					if mode == 1 {
						s.Spells.TargetFlags = 0x8024
					}
					if mode == 2 {
						s.Combat.Wall = 1
					}
					if op == 0 || op == 6 || op == 21 || op == 48 {
						p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: 1700 + op})
					}
					cases = append(cases, s)
				}
			}
		}
		r := sustainedHash(t, fmt.Sprintf("sustained-spatial-%02d", op), cases)
		if op == 0 {
			if r[0].Callbacks.Shop.Sequence[0].Sustained.Record[12] == 0 {
				t.Fatal("mana drain selected no donor")
			}
			_, ud, _ := objectivePlayerData(r[0].Callbacks.Shop.Sequence[1], 0)
			if ud[1]&65535 != 102 {
				t.Fatalf("two mana drain ticks: mana%d want102", ud[1]&65535)
			}
		}
		if op == 7 && r[9].Callbacks.Shop.Sequence[0].Sustained.Globals[6] != 100 {
			t.Fatal("energy candidate did not select target")
		}
		if op == 6 || op == 21 || op == 48 {
			if len(r[9].Callbacks.Damage) == 0 {
				t.Fatalf("spell %d produced no damage callback", op)
			}
		}
	}
}
func TestSustainedSpellsTeleport(t *testing.T) {
	for _, op := range []int{31, 32, 34, 35, 37} {
		var cases []legacy.PortTestRoamSpec
		for slot := 0; slot < 4; slot++ {
			for mode := 0; mode < 6; mode++ {
				s := sustainedBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				sp.Record.Words[4] = uint32(122 + slot)
				sp.Record.Words[68] = 101
				sp.Record.Refs[48] = 1
				a.UpdateRefs = map[int]int{116 + 4*slot: 3}
				a.UpdateWords = map[int]uint32{156: 0x02020202}
				p.TemporaryUpdates.ItemWords[0][56] = math.Float32bits(560)
				p.TemporaryUpdates.ItemWords[0][60] = math.Float32bits(580)
				if op == 37 {
					sp.Record.Refs[48] = 2
					sp.Record.Words[20] = 1
				}
				switch mode {
				case 1:
					sp.Record.Words[68] = 102
				case 2:
					a.ActorWords = map[int]uint32{340: 1 << 14}
				case 3:
					a.ActorWords = map[int]uint32{340: 1}
				case 4:
					a.UpdateWords[156] = 0x01010101
				case 5:
					a.UpdateRefs = nil
				}
				cases = append(cases, s)
			}
		}
		r := sustainedHash(t, fmt.Sprintf("sustained-teleport-%02d", op), cases)
		if op == 32 || op == 34 {
			u := r[0].Callbacks.Shop.Sequence[0].ResourceData[1]
			if u[14] != math.Float32bits(560) || u[15] != math.Float32bits(580) {
				t.Fatalf("teleport %d did not move player to glyph", op)
			}
		}
	}
}

func TestSustainedSpellsWandTransitions(t *testing.T) {
	for _, op := range []int{9, 11, 20, 24, 47, 48, 50} {
		var cases []legacy.PortTestRoamSpec
		for _, charges := range []byte{0, 1, 2, 127, 128, 255} {
			for mode := 0; mode < 3; mode++ {
				s := sustainedBase(op)
				p := s.Callbacks.Shop
				a := p.TemporaryUpdates.World.Objectives.Attack
				sp := a.Controls.SpellLifecycle
				p.Items[0].Class = 0x1000
				p.Items[0].Subclass = 0x4000000
				if op == 20 || op == 24 {
					p.Items[0].Subclass = 0x40000
				}
				if op == 9 || op == 11 {
					p.Items[0].Subclass = 0x200000
				}
				a.UpdateRefs = map[int]int{104: 3}
				p.TemporaryUpdates.World.Objectives.UseWords = []map[int]uint32{{96: 4, 108: uint32(charges) | 255<<8, 112: 100}}
				if op == 11 {
					sp.Record.Refs[76] = 3
				}
				if op == 24 || op == 48 || op == 50 {
					sp.Record.Refs[72] = 3
				}
				if mode == 1 {
					p.TemporaryUpdates.World.Objectives.UseWords[0][96] = 0
				}
				if mode == 2 {
					a.UpdateRefs = nil
					sp.Record.Refs[72] = 0
					sp.Record.Refs[76] = 0
				}
				cases = append(cases, s)
			}
		}
		sustainedHash(t, fmt.Sprintf("sustained-wand-%02d", op), cases)
	}
}
func TestSustainedSpellsFirewalk(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, level := range []uint32{1, 2, 3, 4, 5} {
		for _, dist := range []float32{0, 14.99, 15, 15.01, 30, 100} {
			for mode := 0; mode < 3; mode++ {
				s := sustainedBase(8)
				p := s.Callbacks.Shop
				sp := p.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
				sp.Record.Words[8] = level
				sp.Record.Words[64] = 99
				sp.Record.Words[72] = math.Float32bits(512)
				sp.Record.Words[76] = math.Float32bits(512)
				sp.Record.Words[80] = math.Float32bits(500)
				sp.Record.Words[84] = math.Float32bits(501)
				s.Combat.Target = [2]uint32{math.Float32bits(512 + dist), math.Float32bits(512 + dist)}
				if mode == 1 {
					sp.Record.Words[64] = 100
				}
				if mode == 2 {
					s.Spells.TargetFlags = 0x8024
				}
				cases = append(cases, s)
			}
		}
	}
	r := sustainedHash(t, "sustained-firewalk", cases)
	if got := len(r[15].Lifecycle.Created); got != 2 {
		t.Fatalf("firewalk distance100 creates %d flames want2", got)
	}
}
func TestSustainedSpellsLightningChain(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, level := range []uint32{1, 2, 3, 4, 5} {
		for count := 1; count <= 5; count++ {
			for mode := 0; mode < 3; mode++ {
				s := sustainedBase(21)
				p := s.Callbacks.Shop
				for len(p.Items) < count {
					p.Items = append(p.Items, p.Items[0])
					p.TemporaryUpdates.ItemWords = append(p.TemporaryUpdates.ItemWords, map[int]uint32{})
					p.TemporaryUpdates.UpdateWords = append(p.TemporaryUpdates.UpdateWords, map[int]uint32{})
				}
				o := p.TemporaryUpdates.World.Objectives
				a := o.Attack
				sp := a.Controls.SpellLifecycle
				sp.Record.Words[8] = level
				sp.Record.Refs[48] = 0
				o.PlayerWords[0] = map[int]uint32{56: math.Float32bits(512), 60: math.Float32bits(512), 16: 4}
				p.EffectsUse.Balance["LightningRange"] = []float64{200}
				p.EffectsUse.Balance["LightningDamage"] = []float64{2.5}
				for j := 0; j < count; j++ {
					p.Items[j].Class = 0x20000
					p.Items[j].Subclass = 0
					p.TemporaryUpdates.ItemWords[j][16] = 4
					p.TemporaryUpdates.ItemWords[j][172] = 2
					p.TemporaryUpdates.ItemWords[j][176] = math.Float32bits(1)
					p.TemporaryUpdates.ItemWords[j][56] = math.Float32bits(522 + float32(j)*10)
					p.TemporaryUpdates.ItemWords[j][60] = math.Float32bits(512)
					p.TemporaryUpdates.Indexed = append(p.TemporaryUpdates.Indexed, 3+j)
					sp.Effects.Sustained.RecordDamageRefs = append(sp.Effects.Sustained.RecordDamageRefs, 3+j)
				}
				if mode == 1 {
					sp.Record.Words[20] = 1
				}
				if mode == 2 {
					p.Resources.Mana = 0
					p.Resources.OldMana = 0
				}
				p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: 1721}, legacy.PortTestShopAction{Op: 1724})
				cases = append(cases, s)
			}
		}
	}
	r := sustainedHash(t, "sustained-lightning-chain", cases)
	for i, s := range cases {
		sp := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
		if sp.Record.Words[20] == 0 && s.Callbacks.Shop.Resources.Mana > 0 {
			step := r[i].Callbacks.Shop.Sequence
			want := int(sp.Record.Words[8])
			count := (i%15)/3 + 1
			if want > count {
				want = count
			}
			if got := len(step[0].Sustained.Children); got != want {
				t.Fatalf("chain case %d created %d segments want%d", i, got, want)
			}
			if len(step[2].Sustained.Children) != 0 {
				t.Fatalf("chain case%d cancellation left segments", i)
			}
		}
	}
}

func TestSustainedSpellsHealChannel(t *testing.T) {
	for _, op := range []int{13, 14} {
		var cases []legacy.PortTestRoamSpec
		for _, hp := range []uint16{0, 1, 2, 99, 100} {
			for _, mana := range []uint16{0, 1, 199, 200} {
				for mode := 0; mode < 3; mode++ {
					s := sustainedBase(op)
					p := s.Callbacks.Shop
					sp := p.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
					p.Resources.HP = hp
					p.Resources.MaxHP = 100
					p.Resources.Mana = mana
					p.Resources.OldMana = mana
					sp.Record.Refs[48] = 1
					sp.Record.Words[72] = math.Float32bits(0.75)
					if mode == 1 {
						sp.Record.Words[20] = 1
					}
					if mode == 2 {
						p.Resources.PlayerClass = 2
						s.Lifecycle.GameFlags |= 4096
					}
					p.Sequence = append(p.Sequence, legacy.PortTestShopAction{Op: 1700 + op})
					cases = append(cases, s)
				}
			}
		}
		sustainedHash(t, fmt.Sprintf("sustained-heal-channel-%02d", op), cases)
	}
}
func TestSustainedSpellsNilCasters(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, op := range []int{0, 6, 12, 21, 25, 27, 32, 34, 35, 36, 37, 39, 40, 47, 48} {
		s := sustainedBase(op)
		sp := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.SpellLifecycle
		sp.Record.Refs[16] = 0
		if op == 32 {
			sp.Record.Words[4] = 122
		}
		cases = append(cases, s)
	}
	sustainedHash(t, "sustained-nil-casters", cases)
}
func TestSustainedSpellsFactoryFailure(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, v := range []struct {
		op   int
		name string
	}{{9, "ForceOfNatureCharge"}, {10, "DeathBall"}, {51, "Moonglow"}} {
		s := sustainedBase(v.op)
		a := s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack
		a.MissingTypes = []string{v.name}
		a.Controls.SpellLifecycle.Record.Words[68] = 101
		a.Controls.SpellLifecycle.Record.Refs[48] = 1
		cases = append(cases, s)
	}
	r := sustainedHash(t, "sustained-factory-failure", cases)
	for i, v := range r {
		if len(v.Lifecycle.Created) != 0 {
			t.Fatalf("unavailable named factory case%d created object", i)
		}
	}
}

func TestSustainedSpellsManaBombSequence(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for _, mana := range []uint16{0, 1, 14, 15, 200} {
		for mode := 0; mode < 3; mode++ {
			s := sustainedBase(38)
			p := s.Callbacks.Shop
			a := p.TemporaryUpdates.World.Objectives.Attack
			sp := a.Controls.SpellLifecycle
			p.Resources.Mana = mana
			p.Resources.OldMana = mana
			a.ActorWords = map[int]uint32{120: math.Float32bits(37)}
			sp.Record.Words[68] = 101
			if mode == 1 {
				sp.Record.Words[68] = 200
			}
			if mode == 2 {
				sp.Record.Words[20] = 1
			}
			p.Sequence = []legacy.PortTestShopAction{{Op: 1738}, {Op: 1739}, {Op: 1740}}
			cases = append(cases, s)
		}
	}
	r := sustainedHash(t, "sustained-mana-bomb", cases)
	for i, v := range r {
		steps := v.Callbacks.Shop.Sequence
		if len(v.Lifecycle.Created) != 1 {
			t.Fatalf("mana bomb case%d did not create charge", i)
		}
		if i%3 != 2 {
			if steps[0].ResourceData[1][85]&((1<<5)|(1<<14)|(1<<29)) != ((1 << 5) | (1 << 14) | (1 << 29)) {
				t.Fatalf("mana bomb case%d missing charge buffs", i)
			}
			if steps[2].ResourceData[1][30] != math.Float32bits(37) || steps[2].ResourceData[1][85]&((1<<5)|(1<<14)|(1<<29)) != 0 {
				t.Fatalf("mana bomb case%d did not restore caster", i)
			}
		}
	}
}
