//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"os"
	"testing"
)

var mapRoomHashes = map[string]string{
	"map-rooms-allocation":              "da4bbdde9b83b3644bfaae5223b3a682889d657120f25f2f3e2785b5addbf177",
	"map-rooms-connections":             "b11e5d4d1c62bce227534e5a492640806a5f0db7482330b13e962ad2fb938efd",
	"map-rooms-corridor-trimming":       "0900c53b945b33e10c28b046bfec1fe4589471f28cadba7581056ba660dbc6f8",
	"map-rooms-decoration-boundaries":   "27167d6b060a977467ce222ee4a8943482e0d80dc9db05b8f3edfdde18856a2e",
	"map-rooms-decoration-flags":        "235593102174a248ea87c6f724ebcab665efeaf5e14e4a7a1dab33d384d6ed6d",
	"map-rooms-exclusion-boundaries":    "385ac6559de1b2769a663e8ff2bc8be640f26ee33a1bd50e0a3a4bbaa80afbcd",
	"map-rooms-exclusion-sequence":      "c4f9b7ef7521566108d168ff335e689cbeec49e598351870bd247b0aabddd4b9",
	"map-rooms-free-list":               "eee7b302ef2adb60126815dea31fd910e505ad4f5b4a32a9ee68b51d366a975f",
	"map-rooms-geometry-15":             "e86f5091b1e8f5688ea06dac35e0900b1750415316151fd30b26e365a7ebd9b0",
	"map-rooms-geometry-16":             "e86f5091b1e8f5688ea06dac35e0900b1750415316151fd30b26e365a7ebd9b0",
	"map-rooms-geometry-17":             "79acbf48d82bd5a3802348d1c39cce1779f35a553002c4997568bd8b52904451",
	"map-rooms-geometry-18":             "79acbf48d82bd5a3802348d1c39cce1779f35a553002c4997568bd8b52904451",
	"map-rooms-geometry-26":             "6f9269d15508970a70dc0cff584a86c28188469f8988230ca05e48404aaea91d",
	"map-rooms-geometry-27":             "13b7eb4526f05cc49bca075e5c6b4ee8e45f63205709771d77172b068e26cc92",
	"map-rooms-geometry-28":             "13b7eb4526f05cc49bca075e5c6b4ee8e45f63205709771d77172b068e26cc92",
	"map-rooms-geometry-29":             "ad5d2333726085983d0814fa769191a3a345e3486bb49416ab43974cb7966d3a",
	"map-rooms-geometry-30":             "ad5d2333726085983d0814fa769191a3a345e3486bb49416ab43974cb7966d3a",
	"map-rooms-geometry-33":             "803dc520e5d310a5094d1caf36d1496fe22cb55fb9961506de81fbe61ab926c3",
	"map-rooms-geometry-34":             "9edbc559a4efc88f36fa53b1b1cfe61bf072fe1bf1bcbe3ca7f2be8df257dc07",
	"map-rooms-geometry-35":             "503e841ea205fffa39199487f6fdabdc8ed53f153ae66cd58d359b6f2e596c44",
	"map-rooms-geometry-36":             "9edbc559a4efc88f36fa53b1b1cfe61bf072fe1bf1bcbe3ca7f2be8df257dc07",
	"map-rooms-geometry-37":             "35f188cd28c32af41b052676dcbb6dd123263d42951174bafdb42b817895f165",
	"map-rooms-geometry-38":             "824d0df23614cc19552f55bacc1e5dbcdcb733d07a9a5b926b25bb3d300c4bbc",
	"map-rooms-geometry-43":             "dd0d5a54d24f8530273e081373f06224e7e1df684939878b3c9ac8a78f06d103",
	"map-rooms-geometry-44":             "d99b2b31d49f797a87b72b64e7417ca1f6ecea1440b1afa95818c044652bc027",
	"map-rooms-geometry-45":             "b678e780ff87c25b6ae61188b04798d48231df444678676edec38e9f04d5a291",
	"map-rooms-geometry-46":             "354c71b169a9413625d260a7114f861af69a981e22b38c20fddec1e386954fa8",
	"map-rooms-grid":                    "e7886b39252d51ffafbe1824c09187acfd9e13d1f7f3af171d8cd42f4993d8d1",
	"map-rooms-hall-boundaries":         "cc9c1a6ab8653e74b93cd7f4b21eeb54da733f304a4ff62d5b6d854596ee8f48",
	"map-rooms-limits-15":               "92e787e63367d333b830094f119c0cd7a06b92ab5f74471b5111970e53891658",
	"map-rooms-limits-16":               "92e787e63367d333b830094f119c0cd7a06b92ab5f74471b5111970e53891658",
	"map-rooms-limits-19":               "1ec19b214d640caa7cd9bf1d86075ca1555ba6123d82ed6550d8a72ec28fb412",
	"map-rooms-limits-26":               "5632edb32e398774c9b963ef8e344a2753f2d08b59e5b655c29d012e692b9ae9",
	"map-rooms-limits-33":               "051623366f09af07bda0d3edfd617a26cc4319c5ec2bb92c1b254e5f3dc50b2c",
	"map-rooms-limits-38":               "b5a352f82bc18e774d4a988aab99a8d7ff085cdded559a5e5963f922b8b77837",
	"map-rooms-limits-39":               "aed83ab53cda37fa1c252a9dd73faa76fd8b5312566a227892834590fff8f916",
	"map-rooms-limits-40":               "8904b0d179501b90550c72747ee2a40acceff6b37335e7e6b44c9829fe0dc8dd",
	"map-rooms-limits-41":               "ba7aa6f663cda4214e76db598eaca82ecf3dbcc8445e0cf4fd0df73c3fe3bf0c",
	"map-rooms-limits-42":               "c66f062710796b6334b3f9755179cd76e35c2dc9eae7034fb3d41468ba79cd1c",
	"map-rooms-limits-43":               "1d5b9d657e664c9aacffb98f81728d7b48015b252f7d8198ff330ec7152a9048",
	"map-rooms-limits-44":               "30a2a4fefa66d93f1b847e8573efd650a3dfbc07cf4a5b5dcdd8e76cdac79372",
	"map-rooms-limits-45":               "dd25fcde46352cff94c7fd30f565ef6f296a76ecbbe99890f2aa09df2a234573",
	"map-rooms-limits-46":               "949a045d146cd982daa6f46f3907316e3bc9ea6323358b0a3edab0a6c9fcadb0",
	"map-rooms-limits-52":               "668788b321cf4a6ea845e08169cd8ffc9daf8fdc491a9dc266ec7f4748c236b7",
	"map-rooms-limits-55":               "1ec19b214d640caa7cd9bf1d86075ca1555ba6123d82ed6550d8a72ec28fb412",
	"map-rooms-lists":                   "2c9824154cdda5a465529e4120a5cf4f0177c4bad34627726c5987c60f0e0ae2",
	"map-rooms-occupancy-sequence":      "25afb834900e332bcd5daba5f9ec11ccb6f8cf53ada7ca514ae4ca056380ae3d",
	"map-rooms-point-cache-capacity":    "52ac44b81a199c02fadd257c1f7fa342f4f2414957a51592bab6ba60f15f9497",
	"map-rooms-point-cache-epsilon":     "ff322873006e16a53049305fb08d0a85e52f66590d829b0a83b4e3c9f5daf391",
	"map-rooms-prepare-boundaries":      "8b8acadd1e55ef1b9b81eeaf7a62e3d5677ff72183ec721c1f366e107123fff5",
	"map-rooms-random":                  "8280f38eee80891fb889ce25ab1d39d0f3824f843c2e24b8785bdde77e92ddff",
	"map-rooms-random-variety":          "7a681469c0e984f4f3996414d454dd049595f64cd7e9598e4033fadb2f2a9502",
	"map-rooms-reciprocal":              "3c7cb5559b311f46a36d417542383b1dec990a7619f7eafc08150ae760c48508",
	"map-rooms-required-decorations":    "2707823a8c37df129ce2c4d89b1405e14b4d15a777a5617138b6c86f8d4e57e6",
	"map-rooms-resolve-iteration-limit": "2552931179aecf5fecb3bfae817fc1e554c4a3fa4bf6a5502c3498ea8d996b53",
	"map-rooms-resolve-overlap":         "f27ca162aa4c44d06842f48dfe507e0dfe5dd9e5db1a4b1b5e8f1c6b9ed12228",
	"map-rooms-round-contract-0":        "836287c0d8844ad67de28705473a1902aa0989cc95c53947366d0991a263f3d3",
	"map-rooms-round-contract-1":        "97379de10997fe04e0e8fbf6d32983174bae3ccb9be4d33a3434bc7986b384ff",
	"map-rooms-round-contract-2":        "f03dd1b46bbf001bbca1ffc226f791e633d1e392fd6bce647efe5de73a9f73a6",
	"map-rooms-round-contract-3":        "b604d374a3d4e4a4c731dac4573f647a3f1ae13a81b05826ced3b55f4ce4c624",
	"map-rooms-round-contract-4":        "8be38013a5bf0f07188077e9301a54b06c2a9fd64f86b94d2a5ba07312e14239",
	"map-rooms-rounding":                "8cd4ec34f260d3fd4c7ac85f520ac3bb1a0d1f2aa5b6c21a57e7a426f556b7d2",
	"map-rooms-smoke-00":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-01":                "656b3667d9255a14b1cbeb578f45ce724ca9942bd0d0c0a6051b6d4e64b39c13",
	"map-rooms-smoke-02":                "9a25fc6d778ab12de20661c9d5c8e6a3609417f3763d176a1d1daba0f8afe201",
	"map-rooms-smoke-03":                "ada10e98a605150bc416b111c9920cd660d74613b884f2481f5e02c3a59ca80a",
	"map-rooms-smoke-04":                "bf2dce02e4eb5b50ef3f7203786831660e704a3de05a6288b91dea7176c376ca",
	"map-rooms-smoke-05":                "593fcee719e4646a902621a7531c34aa3c6ab324a75e6ba89549bcde17028856",
	"map-rooms-smoke-06":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-07":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-08":                "fb43ebcd7cddbe93b40040bb4cd3d7ea3c642dd4ca1d908604a4f7d6ef92b3bf",
	"map-rooms-smoke-09":                "3c2ce554736ea403929b12be2e55fe7e63c885fd7019e25d6cf46cb23a576d78",
	"map-rooms-smoke-10":                "6339ed8afa4c1be02bb8c284d16dc80e1f9c78e4570638a639eae2593566a353",
	"map-rooms-smoke-11":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-12":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-13":                "9e195411f240dfb8d88a76deb0e33ba6c1a0e7468fa4f446c5a2580569fa4ab8",
	"map-rooms-smoke-14":                "593fcee719e4646a902621a7531c34aa3c6ab324a75e6ba89549bcde17028856",
	"map-rooms-smoke-15":                "9a25fc6d778ab12de20661c9d5c8e6a3609417f3763d176a1d1daba0f8afe201",
	"map-rooms-smoke-16":                "9a25fc6d778ab12de20661c9d5c8e6a3609417f3763d176a1d1daba0f8afe201",
	"map-rooms-smoke-17":                "301a7c2b624fc768f07e12245275f3a3bcd09516635cf70ab945973903fa0e91",
	"map-rooms-smoke-18":                "301a7c2b624fc768f07e12245275f3a3bcd09516635cf70ab945973903fa0e91",
	"map-rooms-smoke-19":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-20":                "91355ab8320ed6c4a05553a192b9b9a1b84bdee0c3b85b4a70a270fe907b9cfc",
	"map-rooms-smoke-21":                "ff9c12cf92730809ab040509de8bb9700eccde8c628c4e23f17c31d7f9c58b02",
	"map-rooms-smoke-22":                "f6ab117f10981ea7284c24ff1e013ecca05af4f101ae7548840d311d7be8b9e2",
	"map-rooms-smoke-23":                "40fadc846360c7dd45b1ab2c5c287c6df4ee7d5bf8670a4a47fa4b13d20b7789",
	"map-rooms-smoke-24":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-25":                "94472e210b3804f4db8db40e004566c04bcb8a9cffaed64ac391ec7dd5b6c192",
	"map-rooms-smoke-26":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-27":                "488525f99d80996de8a42f72b181e50644a27e4c8362fe3018f18e6aa822795b",
	"map-rooms-smoke-28":                "488525f99d80996de8a42f72b181e50644a27e4c8362fe3018f18e6aa822795b",
	"map-rooms-smoke-29":                "90884e44183364f8664b51fe843128bbce23683de31f50443866b06252d0ba65",
	"map-rooms-smoke-30":                "90884e44183364f8664b51fe843128bbce23683de31f50443866b06252d0ba65",
	"map-rooms-smoke-31":                "f3210140bafeef6f576ac58e290495d4ab9940376b227d11a8ee7af3a8da4754",
	"map-rooms-smoke-32":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-33":                "9a25fc6d778ab12de20661c9d5c8e6a3609417f3763d176a1d1daba0f8afe201",
	"map-rooms-smoke-34":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-35":                "a9c79a1c5afed5d1819ee702fe5a0f31dbde412a056e20893b2267f23c775a84",
	"map-rooms-smoke-36":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-37":                "1bf7e58e3a41abd9d815424a9179ecfb299d41850f72af54ed98eee10715dbf3",
	"map-rooms-smoke-38":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-39":                "e67c00f1c72b6a1618ff61f3e5d074507170f466894f783ff4465aed6b311c37",
	"map-rooms-smoke-40":                "9c47e710eea0fb13b021fd4a7bfa791217ca594767a5a6c925fa795d06e2392f",
	"map-rooms-smoke-41":                "9c47e710eea0fb13b021fd4a7bfa791217ca594767a5a6c925fa795d06e2392f",
	"map-rooms-smoke-42":                "d079cf09a83bceb3f09f9192e1d2957708c67edf23f9c9aaf41b7c1fff009cca",
	"map-rooms-smoke-43":                "9a25fc6d778ab12de20661c9d5c8e6a3609417f3763d176a1d1daba0f8afe201",
	"map-rooms-smoke-44":                "301a7c2b624fc768f07e12245275f3a3bcd09516635cf70ab945973903fa0e91",
	"map-rooms-smoke-45":                "301a7c2b624fc768f07e12245275f3a3bcd09516635cf70ab945973903fa0e91",
	"map-rooms-smoke-46":                "301a7c2b624fc768f07e12245275f3a3bcd09516635cf70ab945973903fa0e91",
	"map-rooms-smoke-47":                "9020e515dcb7d2a4e84330381ad2fb1e36c045f78d82b869e640078c508e5b08",
	"map-rooms-smoke-48":                "ba875cfa2616fc3551a3927f0f727a8421c0ab672b6b17f10dd5ed7b905fdfb7",
	"map-rooms-smoke-49":                "e899911db70baaffa23545bd99692ded02638779c45e2fcbfcbddb4b9967197b",
	"map-rooms-smoke-50":                "5ba78286be3a561378a29a662144c07a19558a1bda04fe69de5f34512bf7b8b9",
	"map-rooms-smoke-51":                "9a25fc6d778ab12de20661c9d5c8e6a3609417f3763d176a1d1daba0f8afe201",
	"map-rooms-smoke-52":                "9a25fc6d778ab12de20661c9d5c8e6a3609417f3763d176a1d1daba0f8afe201",
	"map-rooms-smoke-53":                "ba3f35c119748403a380d1bc0cef0ffaa6713053dcc071c82e0bbad82c64e3e7",
	"map-rooms-smoke-54":                "ccd7c65c6595e99bb2744aac689ed777451608f4d3cc28a9492874b7251f6fab",
	"map-rooms-smoke-55":                "9a25fc6d778ab12de20661c9d5c8e6a3609417f3763d176a1d1daba0f8afe201",
	"map-rooms-smoke-56":                "295a1c0f26fea2354f626e9c31ed5c8cde90bf2c922b7d3203cd89974797d1c9",
	"map-rooms-smoke-57":                "bedd31832ff49e1c709543c768fdd5c1734f512bb29030a8c641db062bf1bb1d",
	"map-rooms-smoke-58":                "c920834cfbc049fcbc12dac79123eb67ccccac1c81e0baa968c5d60dce5f85a2",
	"map-rooms-smoke-59":                "f85625211467ebe5ac1dedd7a050c0dce9109b9f39e4ba58668270563045e832",
}

func roomArg(slot int) legacy.PortTestMapRoomArg  { return legacy.PortTestMapRoomArg{Slot: slot} }
func roomValue(v int32) legacy.PortTestMapRoomArg { return legacy.PortTestMapRoomArg{Value: uint32(v)} }
func roomFloat(v float32) legacy.PortTestMapRoomArg {
	return legacy.PortTestMapRoomArg{Value: math.Float32bits(v)}
}
func roomAction(op int, args ...legacy.PortTestMapRoomArg) legacy.PortTestMapRoomAction {
	a := legacy.PortTestMapRoomAction{Op: op}
	copy(a.Args[:], args)
	return a
}
func roomRecord(size int) legacy.PortTestMapRoomRecord {
	return legacy.PortTestMapRoomRecord{Size: size, Words: map[int]uint32{}, Refs: map[int]legacy.PortTestMapRoomArg{}}
}
func roomSetGeometry(r *legacy.PortTestMapRoomRecord, kind, w, h int32, x, y float32) {
	r.Words[0] = uint32(kind)
	r.Words[4] = uint32(int32(math.Round(float64(x) / 32.526913)))
	r.Words[8] = uint32(int32(math.Round(float64(y) / 32.526913)))
	r.Words[12] = uint32(w)
	r.Words[16] = uint32(h)
	r.Words[20] = math.Float32bits(x)
	r.Words[24] = math.Float32bits(y)
	width, height := float32(float64(w)*32.526913), float32(float64(h)*32.526913)
	r.Words[28] = math.Float32bits(width)
	r.Words[32] = math.Float32bits(height)
	r.Words[36] = math.Float32bits(x)
	r.Words[40] = math.Float32bits(y)
	r.Words[44] = math.Float32bits(x + width)
	r.Words[48] = math.Float32bits(y + height)
	r.Words[364] = 1
}
func roomBase() legacy.PortTestMapRoomSpec {
	s := legacy.PortTestMapRoomSpec{Seed: 12345, Radius: 4}
	for _, n := range []int{1120, 376, 376, 32, 224, 224, 32, 32, 28} {
		s.Records = append(s.Records, roomRecord(n))
	}
	cfg := &s.Records[0]
	cfg.Words[4] = 5
	cfg.Words[8] = 1
	cfg.Words[32] = 8
	cfg.Words[36] = 2
	cfg.Words[64] = math.Float32bits(1000)
	cfg.Words[68] = 4
	for _, off := range []int{88, 120, 184} {
		cfg.Refs[off] = roomArg(5)
		for i := 2; i < 8; i++ {
			cfg.Words[off+i*4] = 10
		}
	}
	roomSetGeometry(&s.Records[1], 1, 4, 4, 0, 0)
	roomSetGeometry(&s.Records[2], 2, 2, 2, 32.526913, 32.526913)
	for _, i := range []int{4, 5} {
		r := &s.Records[i]
		r.Words[64] = 2 << 8
		r.Words[72] = 10
		r.Words[76] = 0
		r.Words[80] = 99
	}
	list := &s.Records[6]
	list.Refs[0] = roomArg(5)
	for i := 2; i < 8; i++ {
		list.Words[4*i] = 10
	}
	rect := &s.Records[8]
	for _, off := range []int{4, 8} {
		rect.Words[off] = math.Float32bits(1)
	}
	for _, off := range []int{12, 16} {
		rect.Words[off] = math.Float32bits(2)
	}
	return s
}
func roomSmoke(op int) legacy.PortTestMapRoomSpec {
	s := roomBase()
	a := roomAction(op, roomArg(2))
	switch op {
	case 0:
		a = roomAction(op, roomArg(4), roomArg(8))
	case 1, 7:
		a = roomAction(op, roomArg(4))
	case 2:
		s.NoGrid = true
		a = roomAction(op, roomArg(1))
	case 3, 9, 11, 24:
		a = roomAction(op)
	case 10:
		s.Buffer = true
		a = roomAction(op)
	case 8:
		a = roomAction(op, roomArg(2), roomArg(3))
	case 15, 16:
		a = roomAction(op, roomArg(1), roomArg(2))
	case 18:
		a = roomAction(op, roomArg(2), roomArg(4))
	case 19, 26:
		a = roomAction(op, roomArg(2), roomValue(0))
	case 20, 25:
		a = roomAction(op, roomArg(2), roomArg(3), roomValue(0))
	case 21:
		a = roomAction(op, roomValue(4), roomValue(3))
		a.Assign = 10
	case 22:
		a = roomAction(op, roomArg(1))
		a.Assign = 10
	case 27, 28, 29, 30:
		a = roomAction(op, roomArg(2), roomArg(3))
	case 31:
		a = roomAction(op, roomArg(2), roomArg(4), roomFloat(2), roomFloat(3))
		a.Assign = 10
	case 33, 34:
		a = roomAction(op, roomArg(2), roomArg(9))
	case 35:
		a = roomAction(op, roomArg(2), roomFloat(0.5), roomArg(8))
	case 36, 37:
		a = roomAction(op, roomArg(2), roomArg(4))
	case 39, 40, 41, 42:
		a = roomAction(op, roomValue(2))
	case 43:
		a = roomAction(op, roomArg(2), roomArg(3))
	case 44, 45, 46:
		a = roomAction(op, roomArg(2), roomArg(8))
	case 47:
		a = roomAction(op, roomValue(2), roomValue(3), roomValue(4))
		a.Assign = 10
	case 48:
		a = roomAction(op, roomArg(1), roomValue(2), roomValue(3))
		a.Assign = 10
	case 49:
		a = roomAction(op, roomArg(1), roomArg(2))
	case 50:
		a = roomAction(op, roomArg(2), roomArg(7))
	case 51, 52:
		a = roomAction(op, roomArg(5), roomArg(2))
	case 53:
		a = roomAction(op, roomArg(7), roomArg(5))
	case 54:
		a = roomAction(op, roomArg(1))
		s.Actions = append(s.Actions, roomAction(13, roomArg(2)))
	case 55:
		a = roomAction(op, roomFloat(1), roomFloat(1))
	case 56:
		a = roomAction(op, roomValue(17))
	case 57:
		a = roomAction(op, roomValue(-3), roomValue(7))
	case 58:
		a = roomAction(op, roomValue(10), roomValue(2))
	case 59:
		a = roomAction(op, roomFloat(-5), roomFloat(8))
	}
	s.Actions = append(s.Actions, a)
	return s
}
func mapRoomsHash(t *testing.T, label string, cases []legacy.PortTestMapRoomSpec) []legacy.PortTestMapRoomResult {
	t.Helper()
	r := legacy.PortTestMapRooms(cases)
	for i, v := range r {
		if !v.Intact || !v.ControlOK {
			t.Fatalf("case %d guards=%v control=%v", i, v.Intact, v.ControlOK)
		}
	}
	data, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_MAP_ROOMS_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d cases %s", label, len(cases), hash)
	if mapRoomHashes[label] == "" {
		t.Fatalf("unlocked map-room baseline: %s", label)
	}
	if want := mapRoomHashes[label]; want != "" && want != hash {
		t.Fatalf("hash %s want %s", hash, want)
	}
	return r
}
func TestMapRoomsSmoke(t *testing.T) {
	for op := 0; op < 60; op++ {
		t.Run(fmt.Sprintf("%02d", op), func(t *testing.T) {
			mapRoomsHash(t, fmt.Sprintf("map-rooms-smoke-%02d", op), []legacy.PortTestMapRoomSpec{roomSmoke(op)})
		})
	}
}
