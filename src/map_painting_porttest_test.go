//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

var mapPaintingHashes = map[string]string{
	"map-painting-border-neighbors":       "5727ef035a610eb33e2ed8af6cca67a216553babd2d7588e1d10e6e15e2a2248",
	"map-painting-border-propagation-43":  "3ba6b31f1004697bb1596d3e3d45f37fc18eab55277eb25324c70638c704a5d9",
	"map-painting-border-propagation-44":  "5a70e70a69e148625041dd24ea10a17555187079e84c54dc81e29c66ac069313",
	"map-painting-border-propagation-47":  "d06e47abde3fd3d13e9fac7bb0b39e3a393c66fff7c55cc3e6f9aff307a20ab4",
	"map-painting-connected-patterns":     "345bb3e55b177d3ab75dca55404dc14c7594a11fcaadfbef636084846cccfda3",
	"map-painting-connected-rooms":        "b1eae2ffff8975337b1f23627931c902f80c38d8705987608d401abe0f899ad5",
	"map-painting-coordinate-origin":      "cf40deefcaf13d7691ae19ba371aacebeec086d5b25bf7408ae6c953f22c336f",
	"map-painting-coordinates":            "c5a901406560c2ea0d409e31e9e1aa09288f5f3126c5c6f097003fed1b4d05f3",
	"map-painting-direction-setters":      "170c4dc385ad0ad3fea6decb126d18a93d3dad686b4014294e7db75205f31fad",
	"map-painting-direction-union":        "bb32729f55e0d98f47a20fe4eb756f9390b37730ecbcf6036f8b6c5f3dd0f499",
	"map-painting-floor-lines":            "a62cdd5f9b1cbbdc2cf70741d90842e44589e05645cfda88ec707565d478ff2c",
	"map-painting-floor-rectangles":       "c43d660282d480eaaef911e510e6ec5e8ce8a5aa83537c75d85de6d5fc56fb1f",
	"map-painting-layout-selection":       "1dca7444b9a05d78b688da359d8e47352e31672b8f51201078f84787d7370047",
	"map-painting-null-objects":           "a6b0ec24eb7ac2b4283a4b75a7ebb42b663811a406bcfdb049955491dce49b3e",
	"map-painting-object-name-gates":      "761ab29fd48ad6043631fe87224c3ceb406d0302558d5d9c5edbdb009f224c4a",
	"map-painting-object-orientation":     "3de394cf185fabe3ddb90fad8b4b15e05d618ad96f58c7465c5caa6ca5203337",
	"map-painting-object-placement":       "59c78029810aeedc54f7dcd1908a298cad8499941319f46525b3991878b945f8",
	"map-painting-patterns-16":            "70a726c0f9809ae02fb9326fa27bf5fb5588688e352f4c278c8f892537764a03",
	"map-painting-patterns-17":            "e6a31345e80573b2ee1369e8c47dd838c3280a7bd1a6231e57c98e2379b22d0a",
	"map-painting-patterns-18":            "65212a9b477a5ed320a9a7db5c070db0b9e93d23de2c915308ad61261cb8653c",
	"map-painting-repeat-object-move":     "a6ce83c5b99c663b94feadbbaf01579d021163c46b41194f27a954b6868a2d07",
	"map-painting-rooms-15":               "ee3ffac5d6e7dbdd8909e09291095175a5531f7a5751c870f8e7e54694da884b",
	"map-painting-rooms-19":               "6ffa70b4acd2634a3419d3c391280cca0712a5e5899002fb445b3814803cfb49",
	"map-painting-rooms-20":               "8af4786fb7fe6c96814f141be6c5bef178c5d9dfc470d0c1b8aa7095ba39892a",
	"map-painting-rooms-24":               "86fab20a07bcba00c54e07a89edbe3e3c685042bbd2c4fa9701042879225a194",
	"map-painting-rooms-25":               "99f3cccda4a6c06f2055e723e98032543badc0d33a66634f9fbb69ad3a10ae30",
	"map-painting-secret-wall-removal":    "81b63ef47ace7e4aea1704d743ba4d18414d8135658f499b388ab0a2cf776d94",
	"map-painting-smoke-00":               "7b24d3843f89b44fbfa6daabecd343f4fa04c89d6903d2eadff0ff3d71775b7f",
	"map-painting-smoke-01":               "5ca697739fd0ac13b610f32e0d635bee8702fa4eab20839e3775d7d6257afac1",
	"map-painting-smoke-02":               "606b657c49197ccf780a9b32ac15aab06cd3c6809e4fcb93337197773ad93c8d",
	"map-painting-smoke-03":               "076952768edec46a72e3727d42c6891802c0adb3612802b2d91cdba02d4cf96d",
	"map-painting-smoke-04":               "cf40deefcaf13d7691ae19ba371aacebeec086d5b25bf7408ae6c953f22c336f",
	"map-painting-smoke-05":               "606b657c49197ccf780a9b32ac15aab06cd3c6809e4fcb93337197773ad93c8d",
	"map-painting-smoke-06":               "d4530c866583937596e67da79ab421635a47ac76c9737605454c02e725e6deee",
	"map-painting-smoke-07":               "d4530c866583937596e67da79ab421635a47ac76c9737605454c02e725e6deee",
	"map-painting-smoke-08":               "ee9fdf52c2d31fc162b5508b223494cc827852e996a139c65d37e3798da1bde9",
	"map-painting-smoke-09":               "0e49cedb5b7140fd6796759152c2b4732e25a6100622c35e40cdc0792095091b",
	"map-painting-smoke-10":               "21ef2602d414fc0481b129eeaebbce3bd6a87dc9529f5e47c2de05706fc2e7d3",
	"map-painting-smoke-11":               "745a99af7de90c7ced3ea362ca6647612e9e320de37883f958731b7f807d5fb9",
	"map-painting-smoke-12":               "29fcab2f42822e0471c19b629a8f6fccb8adef3cb8b68c93c934cc9c61c93b29",
	"map-painting-smoke-13":               "f8ef4fce503c3fa0392743562db009f9a829736a270113c995856351eaee8b1d",
	"map-painting-smoke-14":               "ba7d65a98ed06fd32cf94ca38b2669e2e5c21b780f2dc11dc6b2e8d5c3e8dd0c",
	"map-painting-smoke-15":               "b9b8de62923c7b692358b47ef66bf258cc516369525d6be85484b3d35a6c371c",
	"map-painting-smoke-16":               "c2576239aeefdd8525d676124b9e763b92ce8ec6c8a6033a337d1781649e60e9",
	"map-painting-smoke-17":               "0f886d73a6ef853dfd5f53dba9afe1c2ccc1b1fe0e7535e2906769accd6b9885",
	"map-painting-smoke-18":               "c545b7cb543d51818b7ca18c7e5b33ded8b243f039b3779eb0b0a53e611ab2a0",
	"map-painting-smoke-19":               "4f083f21787653eeb54f041c93cb1e22283755320181720d7079b074a468b8ed",
	"map-painting-smoke-20":               "a16a05257e0ebdc632e2946d77ab303b2c623f446e56ccd6b1c98b6d0b222235",
	"map-painting-smoke-21":               "606b657c49197ccf780a9b32ac15aab06cd3c6809e4fcb93337197773ad93c8d",
	"map-painting-smoke-22":               "606b657c49197ccf780a9b32ac15aab06cd3c6809e4fcb93337197773ad93c8d",
	"map-painting-smoke-23":               "ea2def5dafc8ac1839255c10a44c63e6ee52d4e33f32d51924322f47c2098e24",
	"map-painting-smoke-24":               "606b657c49197ccf780a9b32ac15aab06cd3c6809e4fcb93337197773ad93c8d",
	"map-painting-smoke-25":               "0f2f204bcefbfd47349e26b081f0e8b4cd4b4dc00e869b6de63abf2634b7bea6",
	"map-painting-smoke-26":               "2d143e8329c4ed949348d11dd7a66753fe934f30b3b9ab055dc87879b066e754",
	"map-painting-smoke-27":               "3e1a0a038db9fe10836fb26e2c93fc06d769e1143603af08ebd0267fa7f20936",
	"map-painting-smoke-28":               "14cea7c2164a7b360b4146f3d06886ceb0375e8156335eca59c8409b65dd0ce4",
	"map-painting-smoke-29":               "a9ba1cc74117b509e681517db0a95d64218f567dce7d76ebdd05db24cb2daa74",
	"map-painting-smoke-30":               "dfa202788be0e540f1a7beb9f84e37f80eaa20c0c329b7b18c2821e07a67addb",
	"map-painting-smoke-31":               "da8be833bb2b924b7a0f6415a292644e360acfa311e256054dc9794774ee8875",
	"map-painting-smoke-32":               "bd52af87e32db0514af202fb6dfc712cbee75e1260d1a710e715f6d41886497b",
	"map-painting-smoke-33":               "606b657c49197ccf780a9b32ac15aab06cd3c6809e4fcb93337197773ad93c8d",
	"map-painting-smoke-34":               "e37cd3e2e2492c7db488c937ef57842cc5d641171f593318b7aa0b2da0966691",
	"map-painting-smoke-35":               "6030a43a2a7f331955e1e3e7825bf899be1fe381e94de0b446d4a0e5d6e9fbce",
	"map-painting-smoke-36":               "f50e2c0d734dcd700703ff343243dafa32692818e07f5f24aa4c3d1b1818ff74",
	"map-painting-smoke-37":               "75dd0a3e2031a3548e66118d3ab66838c7a331db217434b7816df3a8de6e0de1",
	"map-painting-smoke-38":               "e99aeaec45105ae02162c951f04b95cf36da71fcc651987cb8b971d13aa9a0fd",
	"map-painting-smoke-39":               "053653c0a805fc58db3b257fc4e938ece94a341598e1ccc3a157f0fb647d4964",
	"map-painting-smoke-40":               "648ae0cc943ddf0270540b0d9c6f4d122e2e5fcf6fe3ba886a44c3c5fee7d533",
	"map-painting-smoke-41":               "0f6a297af866581ee2e00367794644cbeac86fb8aa04e130e1a2917a1210adde",
	"map-painting-smoke-42":               "ff181f9b913da8da4ad83e9662f7b1532800447cf34891811dcdd3d41140bacb",
	"map-painting-smoke-43":               "f50e2c0d734dcd700703ff343243dafa32692818e07f5f24aa4c3d1b1818ff74",
	"map-painting-smoke-44":               "f50e2c0d734dcd700703ff343243dafa32692818e07f5f24aa4c3d1b1818ff74",
	"map-painting-smoke-45":               "a7657396935ffa39e78d35fb20006ac5dbd3a952dabaa2a45d39c9b2595333b4",
	"map-painting-smoke-46":               "984c1eceff450042cd563d17f65ded380591edbb73983d67dd58733d85c8d2cd",
	"map-painting-smoke-47":               "f50e2c0d734dcd700703ff343243dafa32692818e07f5f24aa4c3d1b1818ff74",
	"map-painting-spellbook-type-check":   "d9e748cb4e29aec4ee320762a05add0726edae0eeedf4544fced7bfab94ea845",
	"map-painting-subtile-chain-free":     "aa26e8278e5209775cc542a1ecbd7bea89886a630d4a7051d20c928cfb661f52",
	"map-painting-subtile-none":           "683b15475ee8b0d58b5da12796b2a1cae48ba68a94ebc2e88bb4bd6a7ad4796a",
	"map-painting-subtile-pool":           "5afae77700dc58b1c8851bed78b4627cd5209e6e0936f812b97eaff56b6414a9",
	"map-painting-subtile-sequences":      "9f17945dd408948de9026eff4e21069ee2f012fec5de17faa502aef9d6a936e5",
	"map-painting-tile-anchors":           "d428575cf822026da85f57ed68c1169410f2d7bca215e113e0d9b58e9a96e3db",
	"map-painting-tile-routing":           "092c754dee84075aa245c27aa802bc1e8a851a012754f8f0fdf125c0603cdb02",
	"map-painting-wall-coordinates":       "e8d8a4c4737360cb25bb5c3f8a4a30e28eea72362e1000ff7f3c0186c76bc506",
	"map-painting-wall-floor-masks":       "c75e2295eed65fd64dfa84bdfce79d0a7012727e7d2c3a1881b0388c5e09d46b",
	"map-painting-wall-lifecycle":         "fb3aa4604b1af5aa64c49f65c11b2e62d446a5fc058c2e5d02b5be5203051048",
	"map-painting-wall-neighbor-masks":    "a9ddc59900c380d6c810d11a15634d69a63abd1f4caee999a46caed04a2829ec",
	"map-painting-wall-updates":           "0d97cb7e11138a0f5775edabbf54e6206152c6e9330dcaf3e853338acfd83ff1",
	"map-painting-world-tile-coordinates": "aed2cec23045012c3e85117eea088bfafd5143011c754274ca10853fcd6bad43",
}

func paintAction(op int, args ...legacy.PortTestMapRoomArg) legacy.PortTestPaintAction {
	a := legacy.PortTestPaintAction{Op: op}
	copy(a.Args[:], args)
	return a
}
func paintString(r *legacy.PortTestMapRoomRecord, off int, s string) {
	b := make([]byte, (len(s)+1+3)&^3)
	copy(b, s)
	for i := 0; i < len(b); i += 4 {
		r.Words[off+i] = binary.LittleEndian.Uint32(b[i:])
	}
}
func paintBase() legacy.PortTestPaintSpec {
	s := legacy.PortTestPaintSpec{Seed: 12345, Globals: map[string]legacy.PortTestMapRoomArg{}}
	for _, n := range []int{1120, 376, 376, 32, 224, 128, 128, 32, 20, 64, 32, 32, 32, 20} {
		s.Records = append(s.Records, roomRecord(n))
	}
	roomSetGeometry(&s.Records[1], 1, 4, 4, 0, 0)
	roomSetGeometry(&s.Records[2], 1, 4, 4, 32.526913, 0)
	s.Records[1].Refs[372] = roomArg(5)
	s.Records[2].Refs[372] = roomArg(5)
	s.Records[4].Refs[84] = roomArg(6)
	s.Records[4].Words[88] = 1
	paintString(&s.Records[4], 100, "PaintDoor")
	paintString(&s.Records[4], 160, "PaintDoor")
	paintString(&s.Records[5], 0, "PaintWall")
	paintString(&s.Records[5], 60, "PaintTile")
	s.Records[6].Words[0] = 1
	paintString(&s.Records[6], 4, "PaintBorder")
	paintString(&s.Records[6], 64, "PaintTileTwo")
	s.Records[7].Words[0] = 4
	s.Records[7].Words[4] = 4
	s.Records[8].Words[0] = 1
	paintString(&s.Records[9], 0, "PaintDoor")
	s.Records[10].Words[0] = math.Float32bits(2957)
	s.Records[10].Words[4] = math.Float32bits(2956)
	s.Records[11].Words[0] = 128
	s.Records[11].Words[4] = 128
	s.Records[12].Words[0] = 1
	return s
}
func paintSmoke(op int) legacy.PortTestPaintSpec {
	s := paintBase()
	v := roomValue
	p := roomArg
	var a legacy.PortTestPaintAction
	switch op {
	case 0:
		a = paintAction(op, v(1), v(2), v(0), v(0))
		a.Assign = 50
	case 1, 2:
		a = paintAction(op, p(9))
	case 3:
		a = paintAction(op, v(0), v(1))
	case 4:
		a = paintAction(op, p(4), p(8))
	case 5:
		a = paintAction(op, v(1))
	case 6:
		a = paintAction(op, p(4))
	case 7:
		a = paintAction(op, p(11))
	case 8:
		a = paintAction(op, v(64), v(64), v(23), v(23), p(13))
	case 9:
		a = paintAction(op, v(64), v(64), p(13), v(1), v(0))
	case 10:
		a = paintAction(op, p(5))
	case 11, 12:
		a = paintAction(op, p(4), v(3))
	case 13:
		a = paintAction(op, p(1), p(4), v(3), v(3))
	case 14:
		a = paintAction(op, p(1), p(4), v(3))
	case 15:
		a = paintAction(op, p(1), p(2), p(6))
	case 16, 17, 18:
		a = paintAction(op, p(1), p(7), p(4), p(8))
	case 19:
		a = paintAction(op, p(1), p(2))
	case 20:
		a = paintAction(op, p(2), p(3), v(0))
	case 21, 22:
		a = paintAction(op, p(4), v(3))
	case 23:
		a = paintAction(op, p(4))
	case 24:
		a = paintAction(op, p(1), p(2))
	case 25:
		a = paintAction(op, p(1), p(2), p(3), v(0))
	case 26, 27, 28, 29:
		a = paintAction(op, p(2), p(4), v(4))
	case 30, 31:
		a = paintAction(op, v(1))
	case 32:
		a = paintAction(op, v(8))
	case 33:
		a = paintAction(op, p(4), p(8))
	case 34, 35:
		a = paintAction(op, p(4))
		if op == 35 {
			s.Walls = []legacy.PortTestPaintWall{{X: 128, Y: 128}}
		}
	case 36:
		a = paintAction(op, p(11))
	case 37:
		a = paintAction(op, p(12))
	case 38:
		a = paintAction(op, p(10))
	case 39:
		a = paintAction(op, p(4))
		a.Assign = 50
		s.Globals["dword_5d4594_3835388"] = v(2)
	case 40:
		a = paintAction(op, p(100), p(4))
		s.Objects = []legacy.PortTestPaintObject{{Slot: 100, Type: 2}}
	case 41:
		a = paintAction(op, p(100), v(5))
		s.Objects = []legacy.PortTestPaintObject{{Slot: 100, Type: 2}}
	case 42:
		a = paintAction(op, p(100), v(37))
		s.Objects = []legacy.PortTestPaintObject{{Slot: 100, Type: 3}}
	case 43:
		a = paintAction(op, p(4))
	case 44:
		a = paintAction(op, p(14), p(13), v(46))
	case 45:
		a = paintAction(op, v(64), v(64), v(1), v(0), p(13), v(0))
	case 46:
		a = paintAction(op, p(9), v(1), v(0), v(0), v(0), v(0))
	case 47:
		a = paintAction(op, p(11))
	default:
		panic(op)
	}
	s.Actions = []legacy.PortTestPaintAction{a}
	return s
}
func paintRun(cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	return legacy.PortTestMapPainting(cases, func(core *server.Server) (legacy.Server, func()) {
		old := noxServer
		wrapped := &Server{Server: core}
		core.ExtServer = unsafe.Pointer(wrapped)
		noxServer = wrapped
		return wrapped, func() { noxServer = old; core.ExtServer = nil }
	})
}
func paintingHash(t *testing.T, label string, cases []legacy.PortTestPaintSpec) []legacy.PortTestPaintResult {
	t.Helper()
	out := paintRun(cases)
	for i, r := range out {
		if !r.Intact || !r.ControlOK {
			t.Fatalf("case %d fixture guard/control state", i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_MAP_PAINTING_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d cases %s", label, len(cases), hash)
	if want := mapPaintingHashes[label]; want != hash {
		t.Fatalf("hash %s want %s", hash, want)
	}
	return out
}
func TestMapPaintingSmoke(t *testing.T) {
	for i := 0; i < 48; i++ {
		t.Run(fmt.Sprintf("%02d", i), func(t *testing.T) {
			paintingHash(t, fmt.Sprintf("map-painting-smoke-%02d", i), []legacy.PortTestPaintSpec{paintSmoke(i)})
		})
	}
}
