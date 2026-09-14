//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
	"unicode/utf16"

	"github.com/opennox/opennox/v1/legacy"
)

var gameplayTextHashes = map[string]string{
	"long-chat":        "312f4ea48263662c46b38e6aa446da0869a859b45cef521c278d57b96b0e110c",
	"chat-codes":       "c342c7567901bdf0bf9069dbc05cabd6e4ddb56c87bd07edf51818ad27c95fcf",
	"chat":             "7d44523383372157f60554cf5ba4cfcbf9eca4b3a9ab11e06de586e294f50c88",
	"classifier":       "0a59bb8b9edaa212bb57cca5d8dccc70da922205019cd07ab06ca8a8b761c7d6",
	"formats":          "b023fdee087a63049314846fd10ead68b8e548093aede926ea7bd09160648796",
	"formatted":        "fab891bbb238161465ae9ea87dc26cd6b9bd15997c45fa4c4ec69fb993fce4e9",
	"information-all":  "c788021e18a891c40e4ec7aa103a445e03147b4913bd06cd43bab07a360ee56c",
	"information":      "a7e9cae89154e116e17f359bdec9eb89c6e18da3bde87bf03fb7186ed0dbd719",
	"iteration-guards": "3b6ca7bf1cb2b2f8315e4069e5ef93f97eb10df730c04441beb29ee441ca341a",
	"private":          "2606b1dd48a9c8886f29e9b454cfbceac815182170cff8c47c54662d5608e615",
	"sparse-players":   "1cc57a744b71917814248819a4e6a7df73ae0203c52516ab63960a92bb61eecf",
}

func textArg(ref int) legacy.PortTestGameplayReportArg {
	return legacy.PortTestGameplayReportArg{Kind: "text", Ref: ref}
}
func textBase(op int, args ...legacy.PortTestGameplayReportArg) legacy.PortTestRoamSpec {
	return gameplayReportsBase(100+op, args...)
}
func textSpec(s *legacy.PortTestRoamSpec) *legacy.PortTestGameplayReportsSpec {
	return s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.Controls.Reports
}
func textUnits(s string) []uint16 { return utf16.Encode([]rune(s)) }
func textPrefix(s []uint16) []uint16 {
	for i, v := range s {
		if v == 0 {
			return s[:i]
		}
	}
	return s
}
func textWire(s []uint16, flags byte, code, x, y, extra uint16) []byte {
	s = textPrefix(s)
	wide := false
	for _, v := range s {
		if v > 255 {
			wide = true
		}
	}
	if wide {
		flags |= 4
	} else {
		flags |= 2
	}
	b := make([]byte, 11)
	b[0] = 168
	b[3] = flags
	b[8] = byte(len(s) + 1)
	binary.LittleEndian.PutUint16(b[1:], code)
	binary.LittleEndian.PutUint16(b[4:], x)
	binary.LittleEndian.PutUint16(b[6:], y)
	binary.LittleEndian.PutUint16(b[9:], extra)
	if flags&4 != 0 {
		for _, v := range s {
			b = append(b, byte(v), byte(v>>8))
		}
		b = append(b, 0, 0)
		return b[:11+2*int(b[8])]
	}
	for _, v := range s {
		b = append(b, byte(v))
	}
	b = append(b, 0)
	return b[:11+int(b[8])]
}
func textCapture(t *testing.T, label string, out []legacy.PortTestRoamResult) {
	t.Helper()
	for i, r := range out {
		if !r.Intact || !r.Callbacks.Intact || !r.Spells.Intact || !r.Combat.Intact || !r.MonsterState.Intact || !r.Callbacks.Shop.Intact {
			t.Fatalf("text %s case%d guards", label, i)
		}
	}
	data, err := json.Marshal(out)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_GAMEPLAY_TEXT_CAPTURE"); prefix != "" {
		if err := os.WriteFile(prefix+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d cases %s", label, len(out), hash)
	want, ok := gameplayTextHashes[label]
	if !ok {
		t.Fatalf("missing audited C capture: %s", label)
	}
	if hash != want {
		t.Fatalf("%s hash %s want %s", label, hash, want)
	}
}
func textCheck(t *testing.T, i int, r legacy.PortTestRoamResult, want [3][]byte) {
	t.Helper()
	step := r.Callbacks.Shop.Sequence[0]
	for j := range want {
		if !bytes.Equal(step.ResourceMessages[j], want[j]) {
			t.Fatalf("text case%d recipient%d bytes got%x want%x", i, j, step.ResourceMessages[j], want[j])
		}
	}
	if len(step.Packets) != 0 {
		t.Fatalf("text case%d unexpected reliable queue", i)
	}
}
func textInputs() [][]uint16 {
	out := [][]uint16{{}, {'A'}, {0x7f}, {0x80}, {0xff}, {0x100}, {0xd800}, {0xdc00}, {0xd83d, 0xde00}, {'a', 0, 'z'}, {'A', 0xff, 0x100}}
	for _, n := range []int{253, 254, 255} {
		out = append(out, textUnits(strings.Repeat("x", n)))
	}
	out = append(out, func() []uint16 {
		a := make([]uint16, 253)
		for i := range a {
			a[i] = 0x100
		}
		return a
	}())
	return out
}
func TestGameplayTextClassifier(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var expected []uint32
	for _, v := range textInputs() {
		s := textBase(8, textArg(0))
		textSpec(&s).Texts = [][]uint16{v}
		cases = append(cases, s)
		want := uint32(1)
		for _, v := range textPrefix(v) {
			if v > 255 {
				want = 0
			}
		}
		expected = append(expected, want)
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		if r.Callbacks.Shop.Sequence[0].Return != expected[i] {
			t.Fatalf("classifier case%d", i)
		}
	}
	textCapture(t, "classifier", out)
}
func TestGameplayTextFormatted(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	inputs := textInputs()
	for _, n := range []int{254, 255} {
		v := make([]uint16, n)
		for i := range v {
			v[i] = 0x100
		}
		inputs = append(inputs, v)
	}
	for _, op := range []int{0, 1} {
		for _, v := range inputs {
			for _, flag := range []uint32{0, 1, 4, 16, 255} {
				for players := 0; players <= 3; players++ {
					s := textBase(op, reportObject(100), textArg(0))
					s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Players = players
					textSpec(&s).Texts = [][]uint16{v}
					if op == 1 {
						textSpec(&s).Args[0] = reportValue(flag)
					}
					var want [3][]byte
					if op == 0 {
						want[0] = textWire(v, 0, 0, 0, 0, 0)
					} else {
						for j := 0; j < players; j++ {
							want[j] = textWire(v, byte(flag), 0, 0, 0, 0)
						}
					}
					cases = append(cases, s)
					checks = append(checks, want)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "formatted", out)
}
func TestGameplayTextFormats(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for _, op := range []int{0, 1} {
		for _, v := range []uint32{0, 1, 255, 65535, 0x7fffffff, 0x80000000, 0xffffffff} {
			s := textBase(op, reportObject(100), textArg(0), reportValue(v))
			if op == 1 {
				textSpec(&s).Args[0] = reportValue(0)
			}
			textSpec(&s).Texts = [][]uint16{textUnits("Value:%u %%")}
			var want [3][]byte
			want[0] = textWire(textUnits(fmt.Sprintf("Value:%d %%", int32(v))), 0, 0, 0, 0, 0)
			cases = append(cases, s)
			checks = append(checks, want)
		}
	}
	for _, op := range []int{0, 1} {
		s := textBase(op, reportObject(100), textArg(0), textArg(1))
		if op == 1 {
			textSpec(&s).Args[0] = reportValue(0)
		}
		textSpec(&s).Texts = [][]uint16{textUnits("[%s]"), {'A', 0x100, 0xd800}}
		var want [3][]byte
		want[0] = textWire([]uint16{'[', 'A', 0x100, 0xd800, ']'}, 0, 0, 0, 0, 0)
		cases = append(cases, s)
		checks = append(checks, want)
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "formats", out)
}
func TestGameplayTextInformation(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for _, kind := range append(func() []uint32 {
		var out []uint32
		for i := uint32(0); i <= 25; i++ {
			out = append(out, i)
		}
		return out
	}(), 255, 256, 65536, 0x80000000, 0xffffffff) {
		for j, to := range []uint32{1, 7, 31} {
			for _, value := range []uint32{0, 1, 0x80000000, 0xffffffff} {
				s := textBase(2, reportValue(to), reportValue(kind), legacy.PortTestGameplayReportArg{Kind: "record"})
				s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Attack.RecordWords = map[int]uint32{0: value}
				var b []byte
				switch kind {
				case 0, 1, 2, 12, 13, 16, 20, 21:
					b = []byte{169, byte(kind), byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24)}
				}
				if kind == 17 {
					b = []byte{169, 17}
				}
				var want [3][]byte
				want[j] = b
				cases = append(cases, s)
				checks = append(checks, want)
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "information", out)
}
func TestGameplayTextInformationAll(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for _, kind := range append(func() []uint32 {
		var out []uint32
		for i := uint32(0); i <= 25; i++ {
			out = append(out, i)
		}
		return out
	}(), 255, 256, 65536, 0x80000000, 0xffffffff) {
		for players := 0; players <= 3; players++ {
			s := textBase(3, reportValue(kind), legacy.PortTestGameplayReportArg{Kind: "record"})
			o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
			o.Players = players
			o.Attack.RecordWords = map[int]uint32{0: 0x44332211, 4: 0x88776655, 8: 0xccbbaa99}
			var b []byte
			switch kind {
			case 3, 4, 8, 18, 19, 21:
				b = []byte{169, byte(kind), 0x11, 0x22, 0x33, 0x44}
			case 5, 6, 7, 9, 10, 11:
				b = []byte{169, byte(kind), 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa}
			case 14:
				b = []byte{169, 14, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88, 0x99, 0xaa, 0xbb}
			}
			var want [3][]byte
			for j := 0; j < players; j++ {
				want[j] = b
			}
			cases = append(cases, s)
			checks = append(checks, want)
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "information-all", out)
}
func TestGameplayTextPrivate(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for _, op := range []int{4, 5} {
		for _, length := range []int{0, 1, 47, 48, 49, 63} {
			for _, mask := range []uint32{0, 2, 128, 0x80000082} {
				for players := 0; players <= 3; players++ {
					name := strings.Repeat("m", length)
					s := textBase(op, reportObject(100), legacy.PortTestGameplayReportArg{Kind: "record", Ref: 1}, reportValue(255))
					sp := textSpec(&s)
					sp.Records = []legacy.PortTestGameplayReportRecord{{Text: name}}
					sp.SuppressedPlayers = &mask
					s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Players = players
					if op == 5 {
						sp.Args[0] = sp.Args[1]
					}
					var want [3][]byte
					if length > 0 && length <= 48 {
						for j, ind := range []uint{1, 7, 31} {
							if (op == 4 && j == 0) || (op == 5 && j < players) {
								if mask&(1<<ind) == 0 {
									flag := byte(0)
									if op == 4 {
										flag = 255
									}
									want[j] = append(append([]byte{169, 15, flag}, []byte(name)...), 0)
								}
							}
						}
					}
					cases = append(cases, s)
					checks = append(checks, want)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "private", out)
}
func TestGameplayTextChat(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for _, v := range textInputs() {
		for _, pos := range []float32{-65537.75, -1.75, 0, 1.75, 65535.75, 65536.75} {
			for players := 0; players <= 3; players++ {
				for _, extra := range []uint32{0, 65535, 65536} {
					s := textBase(9, reportObject(100), textArg(0), reportValue(extra))
					textSpec(&s).Texts = [][]uint16{v}
					o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
					o.Players = players
					o.PlayerWords[0] = map[int]uint32{36: 123, 56: math.Float32bits(pos), 60: math.Float32bits(-pos)}
					var want [3][]byte
					for j := 0; j < players; j++ {
						want[j] = textWire(v, 0, 123, uint16(int64(pos)), uint16(int64(-pos)), uint16(extra))
					}
					cases = append(cases, s)
					checks = append(checks, want)
				}
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "chat", out)
}
func TestGameplayTextIterationAndGuards(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	for players := 0; players <= 3; players++ {
		for _, op := range []int{6, 7} {
			for _, ref := range []int{0, 4, 100, 101, 102} {
				s := textBase(op, reportObject(ref))
				s.Callbacks.Shop.TemporaryUpdates.World.Objectives.Players = players
				cases = append(cases, s)
			}
		}
	}
	for _, ref := range []int{0, 4} {
		for _, op := range []int{0, 4} {
			s := textBase(op, reportObject(ref), reportValue(0))
			cases = append(cases, s)
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, [3][]byte{})
	}
	textCapture(t, "iteration-guards", out)
}

func TestGameplayTextSparsePlayers(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for active := 0; active < 8; active++ {
		for units := 0; units < 8; units++ {
			for _, op := range []int{1, 6, 7} {
				s := textBase(op, reportValue(0), textArg(0))
				if op == 7 {
					textSpec(&s).Args[0] = reportObject(100)
				}
				textSpec(&s).Texts = [][]uint16{{'s'}}
				o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
				o.Players = 3
				var want [3][]byte
				for j := 0; j < 3; j++ {
					o.PlayerDataWords[j] = map[int]uint32{2092: uint32((active >> j) & 1)}
					ref := 0
					if units&(1<<j) != 0 {
						ref = 100 + j
					}
					o.PlayerDataRefs[j] = map[int]int{2056: ref}
					if op == 1 && active&units&(1<<j) != 0 {
						want[j] = textWire([]uint16{'s'}, 0, 0, 0, 0, 0)
					}
				}
				cases = append(cases, s)
				checks = append(checks, want)
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "sparse-players", out)
}

func TestGameplayTextChatCodes(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for _, code := range []uint32{0, 1, 32767, 32768, 65535} {
		for _, extent := range []uint32{0, 123, 32767, 32768} {
			for _, class := range []uint32{4, 4 | 0x400000, 4 | 0x20000000} {
				s := textBase(9, reportObject(100), textArg(0), reportValue(0x1234))
				textSpec(&s).Texts = [][]uint16{{0xff, 0x100}}
				o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
				o.PlayerWords[0] = map[int]uint32{8: class, 36: code, 40: extent, 56: 0, 60: 0}
				id := uint16(code)
				if code >= 32768 || extent >= 32768 {
					id = 0
				} else if class&(0x400000|0x20000000) != 0 {
					id = uint16(extent) | 32768
				}
				cases = append(cases, s)
				checks = append(checks, [3][]byte{textWire([]uint16{0xff, 0x100}, 0, id, 0, 0, 0x1234)})
			}
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "chat-codes", out)
}

func TestGameplayTextLongChat(t *testing.T) {
	var cases []legacy.PortTestRoamSpec
	var checks [][3][]byte
	for _, length := range []int{256, 257, 507, 508} {
		for players := 0; players <= 3; players++ {
			text := make([]uint16, length)
			for i := range text {
				text[i] = uint16('a' + i%26)
			}
			s := textBase(9, reportObject(100), textArg(0), reportValue(0x4321))
			textSpec(&s).Texts = [][]uint16{text}
			o := s.Callbacks.Shop.TemporaryUpdates.World.Objectives
			o.Players = players
			o.PlayerWords[0] = map[int]uint32{36: 123, 56: 0, 60: 0}
			b := []byte{168, 123, 0, 2, 0, 0, 0, 0, byte(length + 1), 0x21, 0x43}
			for i := 0; i < int(byte(length+1)); i++ {
				b = append(b, byte(text[i]))
			}
			var want [3][]byte
			for j := 0; j < players; j++ {
				want[j] = b
			}
			cases = append(cases, s)
			checks = append(checks, want)
		}
	}
	out := controlsRun(t, cases)
	for i, r := range out {
		textCheck(t, i, r, checks[i])
	}
	textCapture(t, "long-chat", out)
}
