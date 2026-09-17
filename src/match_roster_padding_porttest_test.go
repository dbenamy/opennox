//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
)

// Both messages have fixed-width string fields. Bytes after the terminator are
// padding, and must not depend on earlier calls using the same C stack frame.
func TestMatchRosterMessagePadding(t *testing.T) {
	o := newQuestRuntimeOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	oldEngine := noxflags.GetEngine()
	t.Cleanup(func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) })
	noxflags.UnsetEngine(noxflags.EngineNoRendering)
	t.Cleanup(legacy.PortTestMatchRosterOwnCache())
	for _, mixed := range []bool{false, true} {
		label := "roster"
		if mixed {
			label = "mixed-roster"
		}
		t.Run(label, func(t *testing.T) {
			for _, name := range []string{"abcdefghijk", "abcdefghi", "a", "", "xyz", "b"} {
				expected := map[uint16][10]byte{}
				for pl := o.s.Players.First(); pl != nil; pl = o.s.Players.Next(pl) {
					value := name
					if mixed && pl.PlayerInd == 3 {
						value = "abcdefghijk"
					}
					b := unsafe.Slice((*byte)(unsafe.Add(pl.C(), 2096)), 12)
					clear(b)
					copy(b, value)
					code := uint16(1000 + int(pl.PlayerInd))
					objectXferSetWord(pl.C(), 2060, uint32(code))
					var field [10]byte
					copy(field[:], value)
					expected[code] = field
				}
				o.reset()
				legacy.PortTestMatchRosterSendPlayers(1)
				count := 0
				for _, node := range o.state().Nodes {
					if len(node.Data) != 129 || node.Data[0] != 45 {
						continue
					}
					count++
					code := binary.LittleEndian.Uint16(node.Data[1:3])
					field, ok := expected[code]
					if !ok {
						t.Fatalf("unknown roster code %d", code)
					}
					if !bytes.Equal(node.Data[119:129], field[:]) {
						t.Errorf("roster name %q, player %d has wrong fixed field/padding: %x want %x", name, code, node.Data[119:129], field)
					}
				}
				if count != 3 {
					t.Fatalf("roster message count %d, want 3", count)
				}
			}
		})
	}
	t.Run("settings", func(t *testing.T) {
		serverName := unsafe.Slice(memmap.PtrUint8(0x5D4594, 1324), 16)
		old := bytes.Clone(serverName)
		defer copy(serverName, old)
		for _, name := range []string{"abcdefghijklmno", "a", "", "xyz", "b"} {
			clear(serverName)
			copy(serverName, name)
			o.reset()
			legacy.PortTestMatchRosterSendSettings()
			count := 0
			for _, node := range o.state().Nodes {
				if len(node.Data) != 49 || node.Data[0] != 176 {
					continue
				}
				count++
				padding := node.Data[1+len(name)+1 : 17]
				if !bytes.Equal(padding, make([]byte, len(padding))) {
					t.Errorf("server name %q has nonzero fixed-field padding: %x", name, padding)
				}
			}
			if count != 1 {
				t.Fatalf("settings message count %d, want 1", count)
			}
		}
	})
}
