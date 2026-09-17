//go:build porttest

package opennox

import (
	"fmt"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/client/noxrender"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type serverOptionsOwner struct {
	*teamUIOwner
	options     *gui.Window
	optionWords map[string]*uint32
	settings    []byte
	packets     func() []legacy.PortTestShopPacketResult
}

func serverOptionsResource() string {
	var b strings.Builder
	b.WriteString("FONT = small; WINDOW 10100 20 30 520 420 USER; STATUS = ENABLED; CHILD\n")
	for _, id := range []int{10101, 10134, 10135} {
		fmt.Fprintf(&b, "WINDOW %d 10 10 180 20 ENTRYFIELD; STATUS = ENABLED; DATA = 64 -1 0 0; END\n", id)
	}
	for _, id := range []int{10114, 10120} {
		fmt.Fprintf(&b, "WINDOW %d 10 40 190 160 SCROLLLISTBOX; STATUS = ENABLED; DATA = 128 1 0 0 1 0 0; END\n", id)
	}
	for _, id := range []int{10110, 10117, 10118, 10121} {
		fmt.Fprintf(&b, "WINDOW %d 10 10 180 20 STATICTEXT; STATUS = ENABLED; DATA = 0 0 WindowDir:Blank; END\n", id)
	}
	for _, id := range []int{10119, 10122, 10141, 10145, 10146, 10149, 10152, 10159, 10160, 10161, 10162, 10163, 10180, 10181, 10330, 10331, 10332, 10333} {
		kind := "PUSHBUTTON"
		if id == 10122 || id == 10330 || id == 10331 || id == 10333 {
			kind = "CHECKBOX"
		}
		fmt.Fprintf(&b, "WINDOW %d 10 240 40 20 %s; STATUS = ENABLED; END\n", id, kind)
	}
	for _, id := range []int{10150, 10153, 10183, 10196, 10197, 10199} {
		fmt.Fprintf(&b, "WINDOW %d 10 10 200 200 USER; STATUS = ENABLED; END\n", id)
	}
	b.WriteString("WINDOW 10182 210 40 16 160 VERTSLIDER; STATUS = ENABLED; DATA = 0 100; END\nEND END")
	return b.String()
}

func newServerOptionsOwner(t *testing.T) *serverOptionsOwner {
	t.Helper()
	o := &serverOptionsOwner{teamUIOwner: newTeamUIOwner(t)}
	reset, packets, freePackets := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(freePackets)
	reset()
	o.packets = packets
	oldSend := o.c.srv.NetSendPacketXxx
	o.c.srv.NetSendPacketXxx = legacy.Nox_xxx_netSendPacket_4E5030
	t.Cleanup(func() { o.c.srv.NetSendPacketXxx = oldSend })

	var entries []strman.Entry
	for _, x := range [][2]string{
		{"guiserv.c:CTF", "Capture flag"}, {"guiserv.c:Arena", "Arena battle"},
		{"guiserv.c:Highlander", "Last survivor"}, {"guiserv.c:KotR", "King of the realm"},
		{"guiserv.c:Flagball", "Flag ball"}, {"guiserv.c:Quest", "Quest adventure"},
		{"Noxworld.c:Chat", "Chat lobby"}, {"WindowDir:Blank", ""},
		{"guiserv.c:SettingsMsg", "Settings for %S"}, {"guiserv.c:GameType", "Choose mode"},
		{"guiserv.c:GameTypeIs", "Mode: %s"}, {"guiserv.c:GoMessage", "Start game"},
		{"guiserv.c:OptsMessage", "Options"}, {"guiserv.c:NumTeamsMsg", "Teams: %d"},
		{"Servopts.wnd:CaptureLimit", "Capture limit"}, {"Servopts.wnd:DeathLimit", "Death limit"},
		{"Servopts.wnd:KillLimit", "Kill limit"},
		{"guiserv.c:AutoAssignOnTT", "assign on"}, {"guiserv.c:AutoAssignOffTT", "assign off"}, {"guiserv.c:TeamDamageOnTT", "damage on"}, {"guiserv.c:TeamDamageOffTT", "damage off"},
		{"guiserv.c:RecPlayers", "%S\t%d-%d"},
		{"servopts.wnd:teams", "Teams"}, {"servopts.wnd:TeamTT", "Choose teams"},
	} {
		entries = append(entries, strman.Entry{ID: strman.ID(x[0]), Vals: []strman.Variant{{Str: x[1]}}})
	}
	language, restore := o.c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	language(0)
	o.configureLanguage = language
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	t.Cleanup(legacy.PortTestServerOptionsModes())
	var done func()
	o.optionWords, done = legacy.PortTestServerOptionsWords()
	t.Cleanup(done)
	for _, region := range [][2]int{{371380, 320}, {371688, 4}, {1324, 16}, {3452, 64}, {2598188, 80}, {1045452, 8}, {1045488, 20}, {1045700, 956}} {
		buf := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, uintptr(region[0]))), region[1])
		old := append([]byte(nil), buf...)
		clear(buf)
		t.Cleanup(func() { copy(buf, old) })
	}
	// This localization lookup uses a filename stored in the original data blob.
	fileName := unsafe.Slice(memmap.PtrUint8(0x587000, 131072), 64)
	oldFileName := append([]byte(nil), fileName...)
	t.Cleanup(func() { copy(fileName, oldFileName) })
	clear(fileName)
	copy(fileName, "guiserv.c")
	// Restore the shipped mode-to-limit index table, and seed distinct live limits.
	for _, table := range blobdata.PortTestScoreboardTables() {
		if table.Base == 0x587000 && table.Offset == 4704 {
			buf := unsafe.Slice((*byte)(memmap.PtrOff(table.Base, table.Offset)), len(table.Data))
			old := append([]byte(nil), buf...)
			t.Cleanup(func() { copy(buf, old) })
			copy(buf, table.Data)
		}
	}
	for i := uintptr(0); i < 6; i++ {
		*memmap.PtrUint16(0x5D4594, 3488+2*i) = uint16(101 + 11*i)
		*memmap.PtrUint8(0x5D4594, 3500+i) = uint8(13 + 7*i)
	}
	head := memmap.PtrOff(0x5D4594, 1045956)
	addr := uint32(uintptr(head))
	*(*[3]uint32)(head) = [3]uint32{addr, addr, addr}
	o.settings = unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 371380)), 58)
	o.options = newWindowFromString(o.c.GUI, serverOptionsResource(), nil)
	if o.options == nil {
		t.Fatal("server options fixture resource")
	}
	serverOptionsArrange(o.options)
	*o.optionWords["root"] = uint32(uintptr(o.options.C()))
	for name, id := range map[string]uint{"maps": 10114, "map-controls": 10183, "limits-panel": 10197, "teams-panel": 10199, "name": 10101, "score": 10134, "time": 10135, "main-panel": 10150, "advanced": 10153} {
		*o.optionWords[name] = uint32(uintptr(o.options.ChildByID(id).C()))
	}
	t.Cleanup(func() {
		if *o.optionWords["root"] != 0 {
			o.call("close", 0, "")
		} else if o.options != nil {
			o.options.Destroy()
		}
		o.c.GUI.FreeDestroyed()
	})
	return o
}
func (o *serverOptionsOwner) call(op string, value int, text string) int {
	return legacy.PortTestServerOptions(op, unsafe.Pointer(&o.settings[0]), value, text)
}
func (o *serverOptionsOwner) text(id uint, text string) {
	o.options.ChildByID(id).Func94(&gui.RawEvent{Event: 16414, Arg1: uintptr(unsafe.Pointer(alloc.InternCString16(text)))})
}
func (o *serverOptionsOwner) event(id uint, code int, a, b uintptr) int {
	return gui.EventRespInt(o.options.ChildByID(id).Func94(gui.AsWindowEvent(code, a, b)))
}
func (o *serverOptionsOwner) entry(id uint) string {
	return alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(o.event(id, 16413, 0, 0)))))
}

func serverOptionsArrange(w *gui.Window) {
	for parent, ids := range map[uint][]uint{10197: {10134, 10135, 10152, 10141}, 10183: {10114, 10180, 10181, 10182}, 10199: {10330, 10331, 10332, 10333}} {
		for _, id := range ids {
			w.ChildByID(id).SetParent(w.ChildByID(parent))
		}
	}
}

func (o *serverOptionsOwner) installConstructor(t *testing.T) *[]string {
	t.Helper()
	table := unsafe.Slice((*byte)(memmap.PtrOff(0x587000, 129760)), 36)
	oldTable := append([]byte(nil), table...)
	t.Cleanup(func() { copy(table, oldTable) })
	for i := 0; i < 9; i++ {
		*memmap.PtrPtr(0x587000, 129760+4*uintptr(i)) = unsafe.Pointer(alloc.InternCString(fmt.Sprintf("server-options-%d.wnd", i)))
	}
	oldLoad := legacy.Nox_xxx_gLoadImg
	legacy.Nox_xxx_gLoadImg = func(name string) *noxrender.Image {
		switch name {
		case "UITabs1":
			return o.images[0]
		case "UITabs2":
			return o.images[1]
		case "UITabs3":
			return o.images[2]
		}
		return oldLoad(name)
	}
	t.Cleanup(func() { legacy.Nox_xxx_gLoadImg = oldLoad })
	var loads []string
	oldParser := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if !strings.HasPrefix(name, "server-options-") {
			return oldParser(name, fn)
		}
		loads = append(loads, name)
		w := newWindowFromString(o.c.GUI, serverOptionsResource(), fn)
		serverOptionsArrange(w)
		o.options = w
		return w
	}
	t.Cleanup(func() { legacy.Nox_new_window_from_file = oldParser })
	o.options.Destroy()
	o.c.GUI.FreeDestroyed()
	o.options = nil
	*o.optionWords["root"] = 0
	// These children belonged to the discarded manual root.
	for _, key := range []string{"maps", "map-controls", "limits-panel", "teams-panel", "name", "score", "time", "main-panel", "advanced"} {
		*o.optionWords[key] = 0
	}
	return &loads
}
