//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

// Reuse the presentation fixture's actual GUI, font/string, renderer and player
// ownership. Scoreboard routines use their actual production implementations.
type scoreboardOwner struct {
	*briefingOwner
	rankWords   map[string]*uint32
	rankRegions [][]byte
}
type scoreboardResult struct {
	Op       int
	Return   uint32
	Named    map[string]uint32
	Regions  [][]byte
	Windows  []inventoryWindowWidget
	Text     []inventoryDisplayText
	Render   objectRenderResult
	Captured uint32
}

func newScoreboardOwner(t *testing.T) *scoreboardOwner {
	o := &scoreboardOwner{briefingOwner: newBriefingOwner(t)}
	installListboxPalette(t)
	var restore func()
	o.rankWords, restore = legacy.PortTestScoreboardWords()
	t.Cleanup(restore)
	for _, r := range [][3]uintptr{{0x5D4594, 1083972, 64}, {0x5D4594, 1084036, 32}, {0x5D4594, 1084068, 64}, {0x5D4594, 1084132, 2560}, {0x5D4594, 1087204, 504}, {0x5D4594, 1088996, 1044}, {0x5D4594, 1090052, 48}, {0x5D4594, 1090104, 4}, {0x5D4594, 1090116, 4}, {0x5D4594, 1090124, 16}, {0x5D4594, 1086692, 512}, {0x5D4594, 3468, 4}, {0x5D4594, 3488, 18}, {0x5D4594, 371380, 58}, {0x587000, 4660, 4}, {0x587000, 4704, 24}, {0x587000, 145584, 2320}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), r[2])
		old := append([]byte(nil), b...)
		o.rankRegions = append(o.rankRegions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	for i, p := range legacy.PortTestScoreboardCallbacks() {
		o.c.callbackRefs[p] = 0xf0500000 + uint32(i)
	}
	for _, id := range []string{"Ball", "Conjurer", "Flag", "HealthHeading", "InternalError", "King", "LivesHeading", "Noxworld.c:Quest", "Noxworld.c:Stage", "TeamPlayerRank", "Teams", "Top3", "Warrior", "Wizard", "WolRank", "class", "ping", "player", "rank", "score", "team", "yourrank", "yourteamrank"} {
		key := id
		if !strings.Contains(key, ":") {
			key = "guirank.c:" + key
		}
		o.strings[key] = strman.Variant{Str: id}
	}
	o.strings["WindowDir:Empty"] = strman.Variant{Str: ""}
	o.strings["guirank.c:TimeRemaining"] = strman.Variant{Str: "Time %d:%02d"}
	o.strings["guirank.c:LessonLimit"] = strman.Variant{Str: "Limit %d"}
	o.installStrings(t)
	t.Cleanup(o.releaseRank)
	o.resetRank(t)
	return o
}
func (o *scoreboardOwner) rankWindow() *gui.Window {
	return (*gui.Window)(unsafe.Pointer(uintptr(*o.rankWords["dword_5d4594_1090048"])))
}
func (o *scoreboardOwner) releaseRank() {
	if w := o.rankWindow(); w != nil {
		w.Destroy()
		*o.rankWords["dword_5d4594_1090048"] = 0
		o.c.GUI.FreeDestroyed()
	}
}
func (o *scoreboardOwner) resetRank(t *testing.T) {
	o.releaseRank()
	o.resetBriefing(t)
	for _, p := range o.rankWords {
		*p = 0
	}
	for _, b := range o.rankRegions {
		clear(b)
	}
	for _, table := range blobdata.PortTestScoreboardTables() {
		copy(unsafe.Slice((*byte)(memmap.PtrOff(table.Base, table.Offset)), len(table.Data)), table.Data)
	}
	for i, s := range []string{"Warrior", "Wizard", "Conjurer"} {
		*memmap.PtrPtr(0x587000, 145676+uintptr(i*4)) = unsafe.Pointer(alloc.InternCString(s))
	}
	for i := range o.players {
		p := &o.players[i]
		p.Active = 0
		p.NetCodeVal = uint32(100 + i)
		p.SetName(fmt.Sprintf("Player%d", i))
	}
	*o.rankWords["nox_player_netCode_85319C"] = 100
	noxflags.ResetGame()
	o.loads = nil
	o.displayText = nil
}
func (o *scoreboardOwner) constructRank(t *testing.T) {
	t.Helper()
	if legacy.PortTestScoreboard(1, 0, 0, 0) == 0 {
		t.Fatal("scoreboard construction")
	}
	o.collect()
}
func (o *scoreboardOwner) rankCapture(t *testing.T, op int, ret uint32) scoreboardResult {
	o.collect()
	ws := o.inventoryWindowOwner.capture(t, 0, 0, 7, 0, 0, 0, 0)
	r := scoreboardResult{Op: op, Return: o.norm(ret), Named: map[string]uint32{}, Windows: ws.Windows, Text: append([]inventoryDisplayText(nil), o.displayText...), Render: ws.Display.Render, Captured: ws.Captured}
	for n, p := range o.rankWords {
		v := *p
		switch n {
		case "dword_5d4594_1090048", "dword_5d4594_1090100", "dword_5d4594_1090108", "dword_5d4594_1090112":
			v = o.norm(v)
		}
		r.Named[n] = v
	}
	for i, b := range o.rankRegions {
		if i == len(o.rankRegions)-1 {
			continue
		}
		cp := append([]byte(nil), b...)
		norm := func(off int) { binary.LittleEndian.PutUint32(cp[off:], o.norm(binary.LittleEndian.Uint32(cp[off:]))) }
		switch i {
		case 1:
			for off := 20; off < 32; off += 4 {
				norm(off)
			}
		case 6, 7:
			for off := 0; off < len(cp); off += 4 {
				norm(off)
			}
		}
		r.Regions = append(r.Regions, cp)
	}
	return r
}
func scoreboardCapture(t *testing.T, label string, rows []scoreboardResult, want string) {
	t.Helper()
	raw, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if prefix := os.Getenv("OPENNOX_SCOREBOARD_CAPTURE"); prefix != "" {
		if err = os.WriteFile(prefix+"-"+label+".json", raw, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(raw))
	t.Logf("%s: %d results %s", label, len(rows), got)
	// Expectations were captured from the original C owners.
	if got != want {
		t.Fatalf("%s hash %s want %s", label, got, want)
	}
}
