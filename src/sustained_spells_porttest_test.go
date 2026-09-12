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
	"sustained-duration-00":     "50db36983fbb6eb1eddb88a563ace072dda471df804259a56652f98ed3f0827a",
	"sustained-duration-05":     "349c508247e7edaa729c91967482ecce9be5edb177316165e614bd8cfea81470",
	"sustained-duration-06":     "65acec26a35133ca5fba547a6d6e755b37b7b4eeaeedce82a37cfe69505ef032",
	"sustained-duration-08":     "695f689c6d299459b38bed47c8d1e2a1589e6c3bfdca33d2a995d230291658a6",
	"sustained-duration-09":     "4bfa7ba648358f7dd01a1524b040fce7015adf7b9adb75be702c928ab717e120",
	"sustained-duration-10":     "76b35683900197027b2dacfc6b8bb58dd161e5af99d5eed3f4e48b7ac48ca1f5",
	"sustained-duration-11":     "38fb98418ae98aac16f17a9ae7152522ba2b2daf866b23637e6e638f12a9d497",
	"sustained-duration-12":     "ed8f4c756c6bad69f13e7ba73b3e4c164f521e90d89a57e0856fafe54e0021f9",
	"sustained-duration-13":     "bc98e8b2c47b25fe18903c8f19c0a481cd7b7f2a602abc5d40fd364e1bfcc922",
	"sustained-duration-14":     "c680a4aa91596c38fff7253c9e7361b67bacd494610c3429a8ae36d1d40d3a05",
	"sustained-duration-15":     "a0c8595b5cbaf89300b187d706e39e1144da5c42a70eb8ade9cd522327370e0a",
	"sustained-duration-16":     "24a993af87b07d76d6ed077d4b6e8e02c07e4af93e2489ea4d2d0a3b45301ff9",
	"sustained-duration-17":     "7aa57a2d1c5c7c986defc26de095a39f33174fbeb9ebef27ee32f3d6c03f439b",
	"sustained-duration-20":     "e698eb30bba41356e686ef97e29c3e34acb981fcd5af424b29a6edfedd28c8f2",
	"sustained-duration-21":     "3f96a64ff2802939ecbae6f796f370bd2f8b87c91c568104b3a7431a720cc5a8",
	"sustained-duration-24":     "38fb98418ae98aac16f17a9ae7152522ba2b2daf866b23637e6e638f12a9d497",
	"sustained-duration-25":     "80e5f4c5ac7459e295fd657361523dc3c2b70459b1500cd242703dd3ae0242db",
	"sustained-duration-26":     "24a993af87b07d76d6ed077d4b6e8e02c07e4af93e2489ea4d2d0a3b45301ff9",
	"sustained-duration-27":     "444741dbc64dd64c35a2ac3c0d7a69ce7b1b2c485624cc1040268ad2814c1ad9",
	"sustained-duration-28":     "a44feab83e681971b33de23191a0363f84a0c45b89d0d0aa80b62968c30f5c24",
	"sustained-duration-29":     "97e9ed63da43cf542e45d962d897aba09118ab47bd7b74e469bad55d836fa429",
	"sustained-duration-31":     "dec75a4052664034b11e0cb4ae139ca6d63e183d145f7a5d185f58b848333c93",
	"sustained-duration-32":     "6fa1b06fe723228a7e439900673c605e9a777b6bad8a55df67b35dc4cb2230c1",
	"sustained-duration-33":     "a44feab83e681971b33de23191a0363f84a0c45b89d0d0aa80b62968c30f5c24",
	"sustained-duration-34":     "e66dd493e8b86690b0178bab04b520e922dde136991e8db9c9adeaf0e2376b01",
	"sustained-duration-35":     "5a08a1599c85780496ede57b54801fadfa1d7c76c73830f4a2e522a45ca30a52",
	"sustained-duration-36":     "3a901eab1cd8792499de5a80fac7ba452b7e5f83f5cd3dea82b8c22efd6a7200",
	"sustained-duration-37":     "417bed2fdc9efca1e3de7633035ad45e7a9dfe817ad577df028119d60ab66ff9",
	"sustained-duration-38":     "50f7e17fc6fad60e499d01fd1cdad396ee46eab0d0aa73a7e2a3d59f9dd5e2ec",
	"sustained-duration-39":     "c44f4487a43f40c1ae1e8deb3d44822519bcf0020ac991f62f6b285b0151e44b",
	"sustained-duration-40":     "ca36ea3fcde65131380e13860a025073a5adfaeb27618455aa8fd7b1d26e156d",
	"sustained-duration-41":     "cf972da2f2b49a8e467880a438608eab83bced212d3217536a58c8a629623b7b",
	"sustained-duration-43":     "bafa7c9be8e07f71716cdcdc477daef4e5445cdd543747a9d985f34ed054128d",
	"sustained-duration-44":     "6eafe707e0b9f470e097fef3056a2ed934e1fcee274033ca34fddd333c1cb148",
	"sustained-duration-45":     "24a993af87b07d76d6ed077d4b6e8e02c07e4af93e2489ea4d2d0a3b45301ff9",
	"sustained-duration-46":     "08d9ada219e034b981872e573237b264cda13c416061082275f7d49f0c2a2c87",
	"sustained-duration-47":     "217219cc6fc424a2b61d5a24f9b86ad421c871d364d5e7902718e6aab09f18aa",
	"sustained-duration-48":     "e95514324f51b3730e9c8eb034aced356f3ebdb757bf83681aed0e55a106b532",
	"sustained-duration-50":     "38fb98418ae98aac16f17a9ae7152522ba2b2daf866b23637e6e638f12a9d497",
	"sustained-duration-51":     "c4940e8957a69855b492b596eb5541e1e60ccb376226d753870741240337445a",
	"sustained-duration-52":     "705f5cedcc6982f7dbc8cfc1677fe15ab8503fd71ad3121206a46d19f7e549ce",
	"sustained-factory-failure": "0d792b119db91f20552bfb96c62169c113c4177e3de5ba62ebafceb33f1803d0",
	"sustained-firewalk":        "2e77900bad20d15047b05caca57642b950d0440ba8f77e3729b7992d240f3dd9",
	"sustained-heal-channel-13": "d7dba6594850761c60ba0634d6b0aa1f3e8436a38dd33d574c2741c03e19061d",
	"sustained-heal-channel-14": "4c6a5429569397d6c23c52f9fb7f0139dfc2fcdd72a5ab5e4d0a1097f1e9adf3",
	"sustained-lightning-chain": "e1222a7e50ce3946b1f0c124e4f8a3aaaf4aa3f9805f85dc5c313fab18baec52",
	"sustained-mana-bomb":       "ee338a20c26d8907314b3c0f87c6f5fb9295f32a3dde6c53a69f2598868d418d",
	"sustained-mana-transfer":   "9a09c4ec815e69248c382a27e77c5243cfc5a092217925c28934dceb495593a5",
	"sustained-nil-casters":     "f2dc938108f1bf35f61e46f76184da22d33ce874bdea5b6ebfb84cdf363cb36b",
	"sustained-positive":        "16508459057e17f7ff110cebf039d6fdb40b178f766d1ce4d280ea2d6f764cf0",
	"sustained-shield-18":       "605087c9b7d7738b19722964bbca725da0aa8c1a8e4999c80d15dd3cc9579afd",
	"sustained-shield-19":       "7d584b235c093ef5143148d513839945a1c26d18e24d9bf4072967803ebf05c9",
	"sustained-smoke-00":        "d4f2602a157726170f0b24f6f9a39c4f9b1a23895f1a62fde5a8dc05e907fc3f",
	"sustained-smoke-01":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-02":        "1bcc8cdb4bdb5887c541c2b1e021963fef009fa8bf20dacb0431a65511105926",
	"sustained-smoke-03":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-04":        "d4f2602a157726170f0b24f6f9a39c4f9b1a23895f1a62fde5a8dc05e907fc3f",
	"sustained-smoke-05":        "0668da5550ead796612af311b337c1569ba567900c3e6bfccf5937da2f33f8dc",
	"sustained-smoke-06":        "104f3c09a658a9d42120d3179f3a0039aa1810c43e9a1340d097886cc501ad7c",
	"sustained-smoke-07":        "558f8debf703d70561c5bfe855e7cc3bab86553b5d88423622f9a3e5dad29434",
	"sustained-smoke-08":        "ad19731b3eff246c858f8339d1fee1566e0340821a7cecc5fde86e25c414383f",
	"sustained-smoke-09":        "29283a060ba2f28cd7d8bcea095c80db139a97d751d96ddd58ef5b6fb48dd441",
	"sustained-smoke-10":        "5545267db7be3e0c5c20bb1be2c74457b76d8927a70bf1710c366a42ae565a9b",
	"sustained-smoke-11":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-12":        "01d70741bb3a3bd7fc400a10a407ba1be140ff4925d776a90bb7148b70dd1e07",
	"sustained-smoke-13":        "121fa10d9c9a3f09914e2cfc587eb3194f0c734020fa2a21a525273dc21643fe",
	"sustained-smoke-14":        "d4f2602a157726170f0b24f6f9a39c4f9b1a23895f1a62fde5a8dc05e907fc3f",
	"sustained-smoke-15":        "29662413819959367df29b3c0fc94bd8ab3d82a40966d96afeb8a68cb6aec7d6",
	"sustained-smoke-16":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-17":        "5615a8d594a0890af6e62b2251b7e3844ea1a9782228823ef07b335c98740f87",
	"sustained-smoke-18":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-19":        "247b31f064b11155dacbf8b2a8235c95e1e026f2869d3ae0a6b671f2c522fb7a",
	"sustained-smoke-20":        "169b0737a89ee3f5ee27782f4e9edae0dcc0356ec6266c732626357585d8e2dc",
	"sustained-smoke-21":        "7497874f07bd3b0fba0d0cba2868822922f03ef9a925dd55b23cecee1a8fdb06",
	"sustained-smoke-22":        "558f8debf703d70561c5bfe855e7cc3bab86553b5d88423622f9a3e5dad29434",
	"sustained-smoke-23":        "562b9196127bda0de6a36de9fa29e741227b66da44bf463bbb65eff0b78f7e01",
	"sustained-smoke-24":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-25":        "46ed7dc1ce9d3d0220872093096a3fd9078ad6015e79f3820e0ce7fcbd74786c",
	"sustained-smoke-26":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-27":        "4364403c50f1b37e9f725eaba1f0a3bfc6e25372c3bc3f1bb0ff286c9ff76703",
	"sustained-smoke-28":        "8afa9addde9688761c5f86238f7edb2b4bc4f1baaa62184fd253445b196506cd",
	"sustained-smoke-29":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-30":        "571ff6d638b15bace1c79508f458118bf771dbd4b1f8d168cdc3f1d75f5bb9df",
	"sustained-smoke-31":        "b0bef013a81daf0da0b9d4eabf61f1dbff7bb1e1dd42d4528a02090b77002346",
	"sustained-smoke-32":        "b0bef013a81daf0da0b9d4eabf61f1dbff7bb1e1dd42d4528a02090b77002346",
	"sustained-smoke-33":        "8afa9addde9688761c5f86238f7edb2b4bc4f1baaa62184fd253445b196506cd",
	"sustained-smoke-34":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-35":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-36":        "8afa9addde9688761c5f86238f7edb2b4bc4f1baaa62184fd253445b196506cd",
	"sustained-smoke-37":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-38":        "00f479ff22bb13d7dc01e2598017087f454e9dc6248bd47f338ef132c7a8afe8",
	"sustained-smoke-39":        "0c3c211e66fbfbf5f5608cb559e4f0c1e3e4c7199fdb9352f3e84bf430b633b4",
	"sustained-smoke-40":        "58a7668a0725a901cdc9f8d1ed809e1bd518e54652cfcad0e08d60c0d2f73d65",
	"sustained-smoke-41":        "c34cb4f0c09a4a85ebdfa22fc258ee33e08891e097073142d649024c15d00f8f",
	"sustained-smoke-42":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-43":        "72e1011b003ddcc0c2373ad2a20c436d77efabd1576554fba56f5e60fba5d66c",
	"sustained-smoke-44":        "c7041956d33a1dce1aa3c67d0a43cfbdf99d55ba18c5b031b8314da9fbd5d911",
	"sustained-smoke-45":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-46":        "49425ebb6d9d8a1abbcd731369ca3b2468d08c560c8984d73035e0b9dda8bde9",
	"sustained-smoke-47":        "b994d72d5827b12f5216c465b5ed88075453c78647af27edd318811e0a2c488c",
	"sustained-smoke-48":        "40e855f01bb9a63e862dbc2de573f6a57af5093f74734d7dae156bf71511f72e",
	"sustained-smoke-49":        "558f8debf703d70561c5bfe855e7cc3bab86553b5d88423622f9a3e5dad29434",
	"sustained-smoke-50":        "5a1410267050ccccbe510b296001724db51d2db4846f9257238d7de66df7c651",
	"sustained-smoke-51":        "6eb0e91e34741b7a7b1c228ee13741b52b6e7422861955db873a7aeade2a0c05",
	"sustained-smoke-52":        "514029a4dfdfd7a81e325bc3d976787931c4712a9f5d7be0554beb2dc134c661",
	"sustained-spatial-00":      "e3744ebee447bbf4772aec8210cde89e259cb5dcddc8339ce5d752806ae07a5d",
	"sustained-spatial-02":      "94b34df91a5b9e955878564146e1b07c2f4f79fa0b2bf319ec720d8cce36091c",
	"sustained-spatial-03":      "d1268bd789f770f6b75b42f0ada43aa58ef4ce96c062913a1bff2431055fbf6b",
	"sustained-spatial-04":      "a0a0a948cb98c62c6ab8da1ad6364d5fcc3cc74fefbc9a727f0a3b3af2e5eccb",
	"sustained-spatial-06":      "ee60a09bee895d20e5f68f8143f804a32f15e0394d8fac35afe2578695d5de2e",
	"sustained-spatial-07":      "f87bd92cb6d5c6f343f9f99a34433fb76a0a4bf37106bde41aecfab261dae58c",
	"sustained-spatial-21":      "4f2a9600310ea2929932e08b1524a1e24bfd54447e4ba0a1c42c7c3b7da37a54",
	"sustained-spatial-22":      "edc173aaefb39c384bd74df84774f1134f7fc997b333b0f9f879e2f166a27014",
	"sustained-spatial-23":      "10a459572ad8c6f9121dd31b25cf1bc253e5fa119bf6ebabb5b350f9c4051026",
	"sustained-spatial-48":      "87c119d4d83275b60a9561a801b776edb71cd7d203a5a72adf00588baa237713",
	"sustained-spatial-49":      "9b09d1c70b4bc40c9704c1422a2d81560d4b99a71b5a9cac8e23440e5568e0f4",
	"sustained-tag-nil":         "5a973c8c26afa9cbdf3dc8b3c15b46101abefc1d88f918bce03fd06cdf4f9c57",
	"sustained-teleport-31":     "13afc5aee00211f6771ffa802763bb9db091c81eda141dbcb3e537c61f3e5b47",
	"sustained-teleport-32":     "3434e4c71f8ba76954bf0de77429774c0e278cc4a92d97532b4a19780b7b0ce9",
	"sustained-teleport-34":     "701230e72a27fd9bd4a36b67479c1930ccdd5aa143ed9129d70e0892f8d0af71",
	"sustained-teleport-35":     "8614c93ef7189f34c1e34550f518190b0299b9259c3d5494b3ae4ba877684b79",
	"sustained-teleport-37":     "0382e8c1dc0a8b8c29b8aa5de13e415c19af491599e0e7230cf6b97fcd63c279",
	"sustained-wand-09":         "3e1adae2c2faa4afef8a3f914c536105edf4e843cd3ba9e226c7527716ec47f2",
	"sustained-wand-11":         "4f6fca504449abcc3945f6a99bb1f1b8e5208cfe31f547174a64a96b39086945",
	"sustained-wand-20":         "4beb13cb7db02cc977d15f4bef3b1dfdeb13ad2e2a14c090792d49db45b3ebcc",
	"sustained-wand-24":         "818286a042124938e05f9ef19e8587004a0f0798418f240b523f365b4bc8ae3d",
	"sustained-wand-47":         "a55f335e18db7ae0c7aa8692eef8ce46559a06e3f9e11da6963e010371b9c4f6",
	"sustained-wand-48":         "7d5dd789d6134424acf914847fc048560a8656a6696ab29268430372b7318295",
	"sustained-wand-50":         "e2edcc7168f3b3006bfd00c59f3b2e3ad6655bb76898e8494b1be1050cafa3b3",
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
