//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxcolor "github.com/opennox/libs/color"
	"github.com/opennox/libs/console"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/input"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"golang.org/x/image/font/basicfont"
	"testing"
	"unsafe"
)

type combatOverlayOwner struct {
	*chatBubbleRenderOwner
	words   map[string]*uint32
	pools   map[string]*unsafe.Pointer
	font    unsafe.Pointer
	console []string
}

type combatOverlayPrinter struct{ owner *combatOverlayOwner }

func (p combatOverlayPrinter) Print(_ console.Color, s string) {
	p.owner.console = append(p.owner.console, s)
}
func (p combatOverlayPrinter) Printf(c console.Color, s string, a ...interface{}) {
	p.Print(c, fmt.Sprintf(s, a...))
}

func newCombatOverlayOwner(t *testing.T) *combatOverlayOwner {
	t.Helper()
	o := &combatOverlayOwner{chatBubbleRenderOwner: newChatBubbleRenderOwner(t, "ArcherBolt", "ArcherArrow", "Bow", "CrossBow")}
	oldConsole := legacy.GetConsole
	con := console.New(combatOverlayPrinter{o})
	legacy.GetConsole = func() *console.Console { return con }
	t.Cleanup(func() { legacy.GetConsole = oldConsole })
	configure, restoreStrings := o.c.srv.PortTestMeterStrings(
		strman.Entry{ID: "die.c:LocalizeAttacker", Vals: []strman.Variant{{Str: "A:%s"}}},
		strman.Entry{ID: "die.c:LocalizeVictim", Vals: []strman.Variant{{Str: "V:%s"}}},
		strman.Entry{ID: "die.c:AttackerNasty", Vals: []strman.Variant{{Str: "nature"}}})
	t.Cleanup(restoreStrings)
	configure(0)

	var restore func()
	o.words, o.pools, restore = legacy.PortTestCombatOverlayGlobals()
	t.Cleanup(restore)
	o.font, restore = o.c.r.GetFonts().PortTestWindowFont(basicfont.Face7x13, "numbers", "large")
	t.Cleanup(restore)
	for _, region := range [][2]uintptr{{1200916, 512}, {1201428, 2400}, {1203828, 4}, {1203844, 16}, {1203872, 4}} {
		clear(serverConfigOwnBytes(t, 0x5D4594, region[0], int(region[1])))
	}
	binary.LittleEndian.PutUint32(serverConfigOwnBytes(t, 0x85B3FC, 940, 4), noxcolor.RGB5551Color(255, 255, 255).Color32())
	binary.LittleEndian.PutUint32(serverConfigOwnBytes(t, 0x5D4594, 2597996, 4), noxcolor.RGB5551Color(255, 255, 255).Color32())
	*o.words["yellow"] = noxcolor.RGB5551Color(255, 255, 0).Color32()
	for _, table := range blobdata.PortTestCombatOverlayTables() {
		copy(serverConfigOwnBytes(t, table.Base, table.Offset, len(table.Data)), table.Data)
	}
	o.c.Inp = input.New(o.c.Log, &entrySeat{}, false, 0)
	if !legacy.PortTestCombatFriendInit() || !legacy.PortTestCombatHealthInit() {
		t.Fatal("combat allocation class initialization")
	}
	friend, health := alloc.AsClass(*o.pools["friend"]), alloc.AsClass(*o.pools["health"])
	t.Cleanup(func() { friend.Free(); health.Free() })
	if *o.words["font"] != uint32(uintptr(o.font)) {
		t.Fatal("numbers font lookup")
	}
	return o
}
func (o *combatOverlayOwner) friendCodes(t *testing.T) []uint32 {
	t.Helper()
	var out []uint32
	for p := unsafe.Pointer(uintptr(*o.words["friends"])); p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 4)) {
		if len(out) >= 128 {
			t.Fatal("friend list cycle/capacity")
		}
		out = append(out, *(*uint32)(p))
	}
	return out
}

type combatHealthRecord struct {
	Code   uint32
	Amount int16
	Stamp  uint32
}

func (o *combatOverlayOwner) healthRecords(t *testing.T) []combatHealthRecord {
	t.Helper()
	var out []combatHealthRecord
	var prev unsafe.Pointer
	for p := unsafe.Pointer(uintptr(*o.words["health"])); p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 12)) {
		if len(out) >= 32 {
			t.Fatal("health list cycle/capacity")
		}
		if *(*unsafe.Pointer)(unsafe.Add(p, 16)) != prev {
			t.Fatal("health previous link")
		}
		out = append(out, combatHealthRecord{objectXferGetWord(p, 0), *(*int16)(unsafe.Add(p, 4)), objectXferGetWord(p, 8)})
		prev = p
	}
	return out
}
