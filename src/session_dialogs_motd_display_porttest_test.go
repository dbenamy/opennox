//go:build porttest

package opennox

import (
	"encoding/binary"
	"slices"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestSessionMOTDDisplayLinesAndGates(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestSessionDialogWords()
	defer restore()
	restoreFlags := noxflags.PortTestGameFlags(0)
	defer restoreFlags()
	oldEngine := noxflags.GetEngine()
	defer func() { noxflags.ResetEngine(); noxflags.SetEngine(oldEngine) }()
	oldDialog := legacy.Sub_44A4A0
	defer func() { legacy.Sub_44A4A0 = oldDialog }()
	dialog := 0
	legacy.Sub_44A4A0 = func() int { return dialog }
	quest := serverConfigOwnBytes(t, 0x5D4594, 1556160, 4)
	message := serverConfigOwnBytes(t, 0x5D4594, 826060, 4)
	pending := serverConfigOwnBytes(t, 0x5D4594, 826068, 4)
	w := (*gui.Window)(unsafe.Pointer(sessionCall("motdOpen", nil, 0, 0, 0)))
	if w == nil {
		t.Fatal("MOTD constructor")
	}
	defer w.Destroy()
	q := sessionQuitWindow(t, o, words)
	defer q.Destroy()
	type row struct {
		Length, Gate int
		Shown        bool
		Pending      uint32
		Text         []string
	}
	var rows []row
	// Byte widening is intentional, including high bytes; each GUI row owns 255 units.
	for _, length := range []int{0, 1, 254, 255, 256, 257, 1024, 4096} {
		line := strings.Repeat("x", length)
		text := line + "\r\n\n\xff\xc3\xa9\rlast"
		ptr, free := alloc.CString(text)
		binary.LittleEndian.PutUint32(message, uint32(uintptr(unsafe.Pointer(ptr))))
		for gate := 0; gate < 10; gate++ {
			noxflags.ResetEngine()
			dialog = 0
			*words["otherDialogA"] = 0
			*words["otherDialogB"] = 0
			q.SetHidden(gate != 1)
			w.SetHidden(gate != 2)
			if gate == 3 {
				noxflags.SetEngine(noxflags.EngineNoRendering)
			}
			if gate == 4 {
				dialog = 1
			}
			if gate == 5 {
				*words["otherDialogA"] = 1
			}
			if gate == 6 {
				*words["otherDialogB"] = 1
			}
			clear(quest)
			if gate == 7 || gate == 8 {
				binary.LittleEndian.PutUint32(quest, 1)
			}
			flags := noxflags.GameFlag(0)
			if gate == 7 || gate == 9 {
				flags = noxflags.GameFlag(4096 | 128)
			}
			restoreCase := noxflags.PortTestGameFlags(flags)
			binary.LittleEndian.PutUint32(pending, 1)
			sessionCall("motdShow", nil, 0, 0, 0)
			names := serverPanelsListNames(w.ChildByID(4203))
			want := []string{}
			if gate == 0 || gate == 8 || gate == 9 {
				first := line
				if length == 0 {
					first = " "
				}
				if len(first) > 255 {
					first = first[:255]
				}
				want = []string{first, " ", "ÿÃ©", "last"}
			}
			if !slices.Equal(names, want) {
				t.Fatalf("length %d gate %d: rows %q want %q", length, gate, names, want)
			}
			shown := sessionCall("motdShown", nil, 0, 0, 0) != 0
			if shown != (gate == 0 || gate == 2 || gate == 8 || gate == 9) {
				t.Fatal("show gate", length, gate, shown)
			}
			p := binary.LittleEndian.Uint32(pending)
			expected := uint32(0)
			if gate == 1 || gate == 7 {
				expected = 1
			}
			if p != expected {
				t.Fatal("pending flag", gate, p, expected)
			}
			rows = append(rows, row{length, gate, shown, p, names})
			sessionCall("motdClose", nil, 0, 0, 0)
			restoreCase()
		}
		if alloc.GoString(ptr) != text {
			t.Fatal("message mutated")
		}
		free()
	}
	sessionCapture(t, "motd-display", rows)
}
