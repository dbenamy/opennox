//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
	"unicode/utf16"
	"unsafe"
)

func TestTeamRuntimeNameStorage(t *testing.T) {
	o := newMatchRosterOwner(t)
	tm := o.s.Teams.Create(1)
	type row struct {
		Name string
		Flag uint32
		Data []byte
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "team-runtime-name-storage", rows, "c318c6a4a51929473c272c41b421932a46de5f3f248855dd621a06e9fdcb7194")
	}()
	for _, name := range []string{"", "A", "Team Ω", "équipe", "😀", strings.Repeat("x", 19), strings.Repeat("y", 20), strings.Repeat("z", 21), strings.Repeat("😀", 11), "A\x00B"} {
		for _, flag := range []uint32{0, 1, 255, 0x80000001, 0xffffffff} {
			t.Run(fmt.Sprintf("%q/%x", name, flag), func(t *testing.T) {
				raw := unsafe.Slice((*byte)(tm.C()), 80)
				for i := 0; i < 44; i++ {
					raw[i] = 0xa5
				}
				old := bytes.Clone(raw)
				legacy.PortTestTeamRuntimeName(tm, name, flag)
				text := utf16.Encode([]rune(strings.SplitN(name, "\x00", 2)[0]))
				if len(text) > 20 {
					text = text[:20]
				}
				want := bytes.Clone(old)
				clear(want[:42])
				for i, c := range text {
					binary.LittleEndian.PutUint16(want[2*i:], c)
				}
				binary.LittleEndian.PutUint32(want[68:], flag)
				if !bytes.Equal(raw, want) {
					t.Fatal("bounded name copy or adjacent team fields")
				}
				rows = append(rows, row{name, flag, bytes.Clone(raw[:44])})
			})
		}
	}
	legacy.PortTestTeamRuntimeName(nil, "nil", 9)
}
