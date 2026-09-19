//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"
)

type playerFileServerOwner struct {
	legacy.Server
	s *server.Server
}

func (o *playerFileServerOwner) S() *server.Server { return o.s }
func playerFileInfoOwner(t *testing.T) []byte {
	t.Helper()
	s := new(server.Server)
	set, restore := s.PortTestWorldCollisionBalance()
	t.Cleanup(restore)
	set(map[string]float64{"MaxExtraLives": 5})
	old := legacy.GetServer
	legacy.GetServer = func() legacy.Server { return &playerFileServerOwner{s: s} }
	t.Cleanup(func() { legacy.GetServer = old })
	p, free := alloc.New(server.PlayerInfo{})
	t.Cleanup(free)
	return unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(*p)))
}
func playerFileName(s string) []byte {
	var out []byte
	for _, v := range utf16.Encode([]rune(s)) {
		out = binary.LittleEndian.AppendUint16(out, v)
	}
	return out
}
func playerFileInfoReset(raw []byte, name string) {
	for i := range raw {
		raw[i] = 0xa5
	}
	for i := 50; i < 89; i++ {
		raw[i] = byte(i*7 + 3)
	}
	n := playerFileName(name)
	copy(raw, n)
	binary.LittleEndian.PutUint16(raw[len(n):], 0)
}
func playerFileAttributes(version uint16, mode uint32, name string, fields []byte) []byte {
	out := binary.LittleEndian.AppendUint16(nil, version)
	if int16(version) >= 5 {
		out = binary.LittleEndian.AppendUint32(out, mode)
	}
	n := playerFileName(name)
	out = append(out, byte(len(n)/2))
	out = append(out, n...)
	out = append(out, fields[50:83]...)
	if int16(version) >= 2 {
		out = append(out, fields[83:88]...)
	}
	out = append(out, fields[88])
	if int16(version) >= 3 {
		out = binary.LittleEndian.AppendUint32(out, 0)
		if version == 3 {
			out = append(out, bytes.Repeat([]byte{0x7d}, 36)...)
		}
	}
	if int16(version) >= 4 {
		out = binary.LittleEndian.AppendUint32(out, 0)
	}
	return out
}
func TestPlayerFilesAttributesWrite(t *testing.T) {
	raw := playerFileInfoOwner(t)
	ptr := uint32(uintptr(unsafe.Pointer(&raw[0])))
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, gameFlags := range []flags.GameFlag{0, 2048, 4096, 8192} {
		for _, name := range []string{"", "A", "Café Ω", "😀", strings.Repeat("X", 24), strings.Repeat("Y", 25)} {
			flags.ResetGame()
			flags.SetGame(gameFlags)
			playerFileInfoReset(raw, name)
			before := append([]byte(nil), raw...)
			mode := uint32(2)
			if gameFlags&4096 != 0 {
				mode = 4
			} else if gameFlags&2048 != 0 {
				mode = 1
			}
			want := playerFileAttributes(5, mode, name, raw)
			wantRet := uint32(1)
			if len(playerFileName(name))/2 >= 25 {
				want = want[:7]
				wantRet = 0
			}
			ret, got, pos := playerFileSection(t, "sub_41A590", nil, 0, ptr)
			if ret != wantRet || !bytes.Equal(got, want) || pos != int64(len(want)) || !bytes.Equal(raw, before) {
				t.Fatal("attributes write", gameFlags, name, ret, wantRet, pos, len(want))
			}
			rows = append(rows, map[string]any{"flags": uint32(gameFlags), "name": name, "return": ret, "bytes": got, "position": pos})
		}
	}
	ret, got, pos := playerFileSection(t, "sub_41A590", nil, 0, 0)
	if ret != 0 || len(got) != 0 || pos != 0 {
		t.Fatal("nil player info")
	}
	rows = append(rows, map[string]any{"nil_info": true, "return": ret, "position": pos})
	spellbookCapture(t, "player-files-attributes-write", rows, "5e075027ec6a87113bb3b8c513c3fdbbd153ea1652e34a00d863307e4689fe1a")
}
func TestPlayerFilesAttributesRead(t *testing.T) {
	raw := playerFileInfoOwner(t)
	ptr := uint32(uintptr(unsafe.Pointer(&raw[0])))
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, version := range []uint16{0, 1, 2, 3, 4, 5, 6, 0x8000, 0xffff} {
		for _, name := range []string{"", "A", "Café Ω", "😀", strings.Repeat("X", 24), strings.Repeat("Y", 25)} {
			playerFileInfoReset(raw, "Before")
			before := append([]byte(nil), raw...)
			fields := append([]byte(nil), raw...)
			for i := 50; i < 89; i++ {
				fields[i] = byte(i*11 + 9)
			}
			input := playerFileAttributes(version, 2, name, fields)
			length := len(input)
			input = append(input, 0xde, 0xad, 0xbe, 0xef)
			expected := append([]byte(nil), before...)
			wantRet, wantPos := uint32(1), int64(length)
			nameBytes := playerFileName(name)
			if int16(version) > 5 {
				wantRet = 0
				wantPos = 2
			} else if len(nameBytes)/2 >= 25 {
				wantRet = 0
				wantPos = 3
				if int16(version) >= 5 {
					wantPos += 4
				}
			} else {
				copy(expected, nameBytes)
				binary.LittleEndian.PutUint16(expected[len(nameBytes):], 0)
				copy(expected[50:83], fields[50:83])
				if int16(version) >= 2 {
					copy(expected[83:88], fields[83:88])
				}
				expected[88] = fields[88]
			}
			ret, got, pos := playerFileSection(t, "sub_41A590", input, 0, ptr)
			if ret != wantRet || pos != wantPos || !bytes.Equal(got, input) || !bytes.Equal(raw, expected) {
				t.Fatal("attributes read", version, name, ret, wantRet, pos, wantPos)
			}
			rows = append(rows, map[string]any{"version": version, "name": name, "return": ret, "position": pos, "info": append([]byte(nil), raw...)})
		}
	}
	spellbookCapture(t, "player-files-attributes-read", rows, "1af6495bd98c46fda79ea158a8904d98cefa456204300ad97a66bc33df6b00d8")
}
