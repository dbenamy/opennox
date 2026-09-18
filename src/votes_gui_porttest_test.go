//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"image"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

// Match the shipped window's IDs, widget types, capacities and selection modes;
// controlled labels and existing fixture artwork keep this contract asset-free.
const voteGUIResource = `FONT = small;
WINDOW
4300 0 50 250 300 USER;
STATUS = ENABLED+ABOVE+HIDDEN+NOFOCUS;
CHILD
 WINDOW
 4301 10 2 230 40 STATICTEXT;
 STATUS = ENABLED+NOFOCUS;
 DATA = 1 0 GUIVote.c:VoteKickPlayer;
 END
 WINDOW
 4311 15 270 55 20 PUSHBUTTON;
 STATUS = ENABLED+NOFOCUS;
 END
 WINDOW
 4312 180 270 55 20 PUSHBUTTON;
 STATUS = ENABLED+NOFOCUS;
 END
 WINDOW
 4320 15 42 220 210 SCROLLLISTBOX;
 STATUS = ENABLED+IMAGE+NOFOCUS+ONE_LINE;
 DATA = 32 10 0 0 1 1 1;
 END
 WINDOW
 4321 15 42 220 210 SCROLLLISTBOX;
 STATUS = ENABLED+IMAGE+NOFOCUS+ONE_LINE;
 DATA = 32 10 0 0 1 0 1;
 END
END
END
`

func newVoteGUIOwner(t *testing.T) (*teamUIOwner, map[string]*uint32, *gui.Window) {
	o := newTeamUIOwner(t)
	o.c.srv.Teams.Reset()
	o.c.srv.Teams.ActiveCnt = 0
	words, restore := legacy.PortTestVoteGUIOwner()
	t.Cleanup(restore)
	var entries []strman.Entry
	for _, name := range []string{"VoteKickPlayer", "SelectVoteTopic", "VoteTopicLabel", "VoteResetServer", "Vote:ResetQuest", "WindowDir:Yes", "WindowDir:No"} {
		id := name
		if !strings.Contains(id, ":") {
			id = "GUIVote.c:" + id
		}
		entries = append(entries, strman.Entry{ID: strman.ID(id), Vals: []strman.Variant{{Str: name}}})
	}
	configure, restoreStrings := o.c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restoreStrings)
	configure(0)
	oldParser := legacy.Nox_new_window_from_file
	legacy.Nox_new_window_from_file = func(name string, fn gui.WindowFunc) *gui.Window {
		if name != "GuiKick.wnd" {
			t.Fatalf("unexpected resource %q", name)
		}
		return newWindowFromString(o.c.GUI, voteGUIResource, fn)
	}
	t.Cleanup(func() { legacy.Nox_new_window_from_file = oldParser })
	for i := 0; i < 3; i++ {
		o.players[i].Active = 1
		o.players[i].Field4792 = 1
		o.players[i].PlayerInd = byte(i)
		o.players[i].SetName(fmt.Sprintf("voter%d", i))
	}
	if legacy.PortTestVoteGUI("init", 0) != 1 {
		t.Fatal("window initialization")
	}
	w := legacy.PortTestVoteGUIWindow()
	if w == nil {
		t.Fatal("nil window")
	}
	return o, words, w
}

func TestPortVotesGUISelection(t *testing.T) {
	o, words, w := newVoteGUIOwner(t)
	list := w.ChildByID(4320)
	legacy.PortTestVoteGUI("show", 0)
	// Excluding the local player leaves voter1 and voter2. Selection changes
	// travel through the actual multi-list widget and window callback.
	var capture [][][]byte
	for step, selected := range [][]int{{0}, {0}, {0, 1}, {1}, {}} {
		o.c.srv.NetList.ResetAll()
		list.Func94(&gui.RawEvent{Event: 16403, Arg1: ^uintptr(0)})
		for _, i := range selected {
			list.Func94(&gui.RawEvent{Event: 16405, Arg1: uintptr(i)})
		}
		w.Func94(&gui.RawEvent{Event: 16391, Arg1: uintptr(unsafe.Pointer(w.ChildByID(4311)))})
		var packets [][]byte
		o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { packets = append(packets, bytes.Clone(b)); return false })
		capture = append(capture, packets)
		expected := []struct {
			action byte
			name   string
		}{{0, "voter1"}, {0, ""}, {0, "voter2"}, {2, "voter1"}, {2, "voter2"}}[step]
		n := 1
		if expected.name == "" {
			n = 0
		}
		if len(packets) != n {
			t.Fatalf("step %d: %d messages, want %d", step, len(packets), n)
		}
		if n != 0 {
			want := make([]byte, 52)
			want[0] = 238
			want[1] = expected.action
			for i, c := range []byte(expected.name) {
				want[2+2*i] = c
			}
			if !bytes.Equal(packets[0], want) {
				t.Fatalf("step %d: message %x, want %x", step, packets[0], want)
			}
		}
		if *words["count"] != uint32(len(selected)) {
			t.Fatalf("selection %v produced count %d", selected, *words["count"])
		}
		legacy.PortTestVoteGUI("show", 0)
	}
	spellbookCapture(t, "votes-gui-selection", capture, "4ae8b85d28d86ce5b0229b302e7f315001051479594b48b34780904445eb44b5")
	for _, topic := range []int{1, 3, 4, 2} {
		legacy.PortTestVoteGUI("show", topic)
		if *words["topic"] != uint32(topic) {
			t.Fatal("topic")
		}
	}
	if legacy.PortTestVoteGUI("hide", 0) != 1 || legacy.PortTestVoteGUI("hide", 0) != 0 {
		t.Fatal("hide transition")
	}
}

func TestPortVotesGUITopicsAndReset(t *testing.T) {
	_, words, w := newVoteGUIOwner(t)
	reset, packets, free := legacy.PortTestVisibilityReliableOwner()
	t.Cleanup(free)
	reset()
	var capture []struct {
		Topic    uint32
		Names    []string
		Choice   uint32
		Messages [][]byte
	}
	topics := w.ChildByID(4321)
	players := w.ChildByID(4320)
	legacy.PortTestVoteGUI("show", 4)
	if !players.GetFlags().IsHidden() || topics.GetFlags().IsHidden() {
		t.Fatal("topic list visibility")
	}
	if names := teamUIRowNames(topics); len(names) != 2 || names[0] != "VoteTopicLabel VoteResetServer" || names[1] != "VoteTopicLabel VoteKickPlayer" {
		t.Fatalf("topics %v", names)
	}
	for _, selected := range []int{0, 0, 1, 1, 0} {
		legacy.PortTestVoteGUI("show", 2)
		topics.Func94(&gui.RawEvent{Event: 16403, Arg1: uintptr(selected)})
		w.Func94(&gui.RawEvent{Event: 16391, Arg1: uintptr(unsafe.Pointer(w.ChildByID(4311)))})
		want := uint32(1 - selected)
		if *words["choice"] != want || *words["previousChoice"] != want {
			t.Fatal("reset choice state")
		}
		row := struct {
			Topic    uint32
			Names    []string
			Choice   uint32
			Messages [][]byte
		}{*words["topic"], teamUIRowNames(topics), *words["choice"], nil}
		for _, p := range packets() {
			row.Messages = append(row.Messages, p.Data)
		}
		capture = append(capture, row)
	}
	all := packets()
	if len(all) != 3 || !bytes.Equal(all[0].Data, []byte{238, 4}) || !bytes.Equal(all[1].Data, []byte{238, 5}) || !bytes.Equal(all[2].Data, []byte{238, 4}) {
		t.Fatalf("reset message transitions %+v", all)
	}
	// Topic confirmation opens the corresponding panel without closing it.
	for _, selected := range []int{0, 1} {
		legacy.PortTestVoteGUI("show", 4)
		topics.Func94(&gui.RawEvent{Event: 16403, Arg1: uintptr(selected)})
		w.Func94(&gui.RawEvent{Event: 16391, Arg1: uintptr(unsafe.Pointer(w.ChildByID(4311)))})
		if *words["topic"] != uint32(2+selected) || w.GetFlags().IsHidden() {
			t.Fatal("topic confirmation")
		}
	}
	w.Func94(&gui.RawEvent{Event: 16391, Arg1: uintptr(unsafe.Pointer(w.ChildByID(4312)))})
	if !w.GetFlags().IsHidden() || len(packets()) != 3 {
		t.Fatal("cancel")
	}
	spellbookCapture(t, "votes-gui-topics", capture, "76d58750ee8c0425d48d1216ff399a82f9702934553cbd8a5e7e8378a88b3c82")
}

func TestPortVotesGUITeamFilter(t *testing.T) {
	o, words, w := newVoteGUIOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	t.Cleanup(o.c.srv.PortTestMapDrawableTeamMessages())
	team := o.c.srv.Teams.Create(1)
	o.c.srv.Teams.Create(2)
	// Team palette rows are controlled nonzero values, not zero-filled blob data.
	palette := memmap.PtrUint32(0x587000, 156400+8*uintptr(team.ID()%10))
	oldPalette := *palette
	*palette = 7
	t.Cleanup(func() { *palette = oldPalette })
	for i := 0; i < 3; i++ {
		code := 7 + i
		o.players[i].NetCodeVal = uint32(code)
		dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(10+i*10, 10))
		dr.NetCode32 = uint32(code)
		dr.ObjClass = 4
		id := 1
		if i == 2 {
			id = 2
		}
		legacy.Nox_xxx_createAtImpl_4191D0(server.TeamID(id), dr.TeamPtr(), 0, code, 0)
	}
	legacy.PortTestVoteGUI("show", 0)
	list := w.ChildByID(4320)
	if names := teamUIRowNames(list); !reflect.DeepEqual(names, []string{"voter1"}) {
		t.Fatalf("team player list %v", names)
	}
	list.Func94(&gui.RawEvent{Event: 16405, Arg1: 0})
	legacy.PortTestVoteGUI("selection", 0)
	if *words["count"] != 1 {
		t.Fatal("team selection")
	}
	legacy.PortTestVoteGUI("show", 0)
	d := (*gui.ScrollListBoxData)(list.WidgetData)
	if d.Items.Field_129 != legacy.Get_nox_color_red_2589776() {
		t.Fatalf("team row color %x", d.Items.Field_129)
	}
	spellbookCapture(t, "votes-gui-team", struct {
		Names []string
		Color uint32
	}{teamUIRowNames(list), d.Items.Field_129}, "84c8ac471f2d0f14366578bd0a65aa5eabccae1ea5043ad64d3b2e3f09ba11ed")
	if sel := unsafe.Slice((*int32)(unsafe.Pointer(uintptr(d.Field_12))), 2); sel[0] != 0 || sel[1] != -1 {
		t.Fatalf("restored team selection %v", sel)
	}
}

func TestPortVotesGUIFullRoster(t *testing.T) {
	o, words, w := newVoteGUIOwner(t)
	for i := range o.players {
		o.players[i].Active = 1
		o.players[i].PlayerInd = byte(i)
		o.players[i].SetName(fmt.Sprintf("player%02d", i))
	}
	legacy.PortTestVoteGUI("show", 0)
	list := w.ChildByID(4320)
	if len(teamUIRowNames(list)) != 31 {
		t.Fatal("full roster")
	}
	for i := 0; i < 31; i++ {
		list.Func94(&gui.RawEvent{Event: 16405, Arg1: uintptr(i)})
	}
	var capture [][][]byte
	for step := 0; step < 2; step++ {
		o.c.srv.NetList.ResetAll()
		if step == 1 {
			list.Func94(&gui.RawEvent{Event: 16403, Arg1: ^uintptr(0)})
		}
		legacy.PortTestVoteGUI("selection", 0)
		var packets [][]byte
		o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { packets = append(packets, bytes.Clone(b)); return false })
		if len(packets) != 31 || *words["count"] != uint32(31*(1-step)) {
			t.Fatalf("full roster step%d: %d messages %d selected", step, len(packets), *words["count"])
		}
		for i, p := range packets {
			want := make([]byte, 52)
			want[0] = 238
			want[1] = byte(step * 2)
			for j, c := range []byte(fmt.Sprintf("player%02d", i+1)) {
				want[2+2*j] = c
			}
			if !bytes.Equal(p, want) {
				t.Fatalf("step%d player%d message", step, i+1)
			}
		}
		capture = append(capture, packets)
	}
	spellbookCapture(t, "votes-gui-full-roster", capture, "c6e34a15de40137ed5ce0735a64e415f4ce851b8b2ae1473dd6a2d53c92ce2ba")
}
