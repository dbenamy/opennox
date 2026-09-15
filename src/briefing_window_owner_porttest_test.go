//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/client/audio/ail"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/dialog"
	"github.com/opennox/opennox/v1/legacy/timer"
)

type briefingWindowOwner struct {
	*briefingOwner
	dialogState [6]uint32
	ticks       uint64
}
type briefingWindowResult struct {
	Base        briefingResult
	Value       uint64
	Dialogue    string
	MapDisabled int
	Sounds      [][2]int
}

func newBriefingWindowOwner(t *testing.T) *briefingWindowOwner {
	o := &briefingWindowOwner{briefingOwner: newBriefingOwner(t)}
	var driver ail.Driver
	var timers [4]timer.TimerGroup
	old := legacy.Dialogs
	legacy.Dialogs = dialog.NewDialog("dialog", &o.dialogState[0], &o.dialogState[1], &o.dialogState[2], &o.dialogState[3], &o.dialogState[4], &driver, &o.dialogState[5], o.c.srv.Strings, &timers[0], &timers[1], &timers[2], &timers[3], nil, nil, nil, nil, nil, nil, nil)
	t.Cleanup(func() { legacy.Dialogs = old })
	gate := memmap.PtrUint32(0x5D4594, 527720)
	saved := *gate
	t.Cleanup(func() { *gate = saved })
	mapDisabled := legacy.Get_nox_gameDisableMapDraw_5d4594_2650672()
	t.Cleanup(func() { legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(mapDisabled) })
	oldClass, oldEnding := dword_587000_311372, dword_5d4594_2516476
	t.Cleanup(func() { dword_587000_311372, dword_5d4594_2516476 = oldClass, oldEnding })
	for i, p := range legacy.PortTestBriefingWindowCallbacks() {
		o.c.callbackRefs[p] = 0xf0110000 + uint32(i)
	}
	for i, s := range []string{"gmcap14e.wav", "gmcap13e.wav", "gmcap17e.wav", "gmcap18e.wav", "gmcap15e.wav", "gmcap16e.wav"} {
		o.c.dataRefs[uint32(txptr(unsafe.Pointer(alloc.InternCString(s))))] = 0xf0120000 + uint32(i)
	}
	legacy.PlatformTicks = func() uint64 { return o.ticks }
	o.resetWindow(t)
	return o
}
func (o *briefingWindowOwner) resetWindow(t *testing.T) {
	o.resetBriefing(t)
	clear(o.dialogState[:])
	o.dialogState[0] = 1
	legacy.Dialogs.Sub_44D8F0()
	*memmap.PtrUint32(0x5D4594, 527720) = 0
	legacy.Set_nox_gameDisableMapDraw_5d4594_2650672(0)
	dword_587000_311372, dword_5d4594_2516476 = 0, 0
	o.ticks = 123456
	// Exercise the real save-menu boundary in quest mode; its early return avoids
	// introducing unrelated save-selector resources into transition unit fixtures.
	noxflags.SetGame(noxflags.GameModeQuest)
}
func (o *briefingWindowOwner) windowCapture(t *testing.T, op int, value uint64) briefingWindowResult {
	return briefingWindowResult{Base: o.snapshotBriefing(t, op, uint32(value)), Value: value, Dialogue: legacy.Dialogs.FileToRead(), MapDisabled: legacy.Get_nox_gameDisableMapDraw_5d4594_2650672(), Sounds: append([][2]int(nil), o.sounds...)}
}
func briefingWindowCapture(t *testing.T, label string, rows []briefingWindowResult, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_BRIEFING_WINDOW_CAPTURE"); prefix != "" {
		if err = os.WriteFile(prefix+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Logf("%s: %d results %s", label, len(rows), got)
	if got != want {
		t.Fatalf("%s hash %s want %s", label, got, want)
	}
}
