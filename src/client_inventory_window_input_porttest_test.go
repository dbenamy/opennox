//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/common/memmap"
	"testing"
)

func TestClientInventoryWindowOpeningGates(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	gates := []string{"none", "quit", "dialog", "dialog2", "cursor", "animation1", "animation2", "animation51", "no-player", "other", "hidden-quit", "hidden-dialog", "hidden-dialog2", "animation34"}
	for i, gate := range gates {
		o.reset(t)
		o.construct(t)
		*o.windowWords["dword_5d4594_1062516"] = 125
		*o.windowWords["dword_5d4594_1062512"] = 50
		blocked := i > 0 && i < 10
		switch gate {
		case "quit", "hidden-quit":
			*o.windowWords["nox_wnd_quitMenu_825760"] = uint32(txptr(o.mainWindow().C()))
		case "dialog", "hidden-dialog":
			*o.windowWords["dword_5d4594_1309820"] = uint32(txptr(o.mainWindow().C()))
		case "dialog2", "hidden-dialog2":
			*o.windowWords["dword_5d4594_1321228"] = uint32(txptr(o.mainWindow().C()))
		case "cursor":
			*memmap.PtrUint32(0x5D4594, 1096672) = 7
		case "animation1":
			*txword(o.items[0], 276) = 1
		case "animation2":
			*txword(o.items[0], 276) = 2
		case "animation51":
			*txword(o.items[0], 276) = 51
		case "animation34":
			*txword(o.items[0], 276) = 34
		case "no-player":
			*memmap.PtrPtr(0x852978, 8) = nil
		case "other":
			*o.windowWords["dword_5d4594_2523804"] = 3
		}
		if i >= 10 && i <= 12 {
			o.mainWindow().Hide()
		} else {
			o.mainWindow().Show()
		}
		r := o.capture(t, i, 0, 28, 0, 0, 0, 0)
		rows = append(rows, r)
		wantState, wantScroll := byte(1), uint32(125)
		if blocked {
			wantState, wantScroll = 0, 50
		}
		if got := *memmap.PtrUint8(0x5D4594, 1049868); got != wantState {
			t.Fatalf("%s state%d want%d", gate, got, wantState)
		}
		if got := *o.windowWords["dword_5d4594_1062512"]; got != wantScroll {
			t.Fatalf("%s scroll%d want%d", gate, got, wantScroll)
		}
	}
	inventoryWindowCapture(t, "opening-gates", rows, "2eeaf5823cd9e39b0aaaba36d586693390f10bea0727d863d0201ea20d7e68ef")
}

func TestClientInventoryWindowInputGates(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for _, event := range []uintptr{5, 6, 7, 8, 9, 19, 20, 99} {
		for gate := 0; gate < 8; gate++ {
			o.reset(t)
			o.construct(t)
			o.openInput()
			o.stack(t, 0, 0, 3, "RedApple", 100)
			*o.windowWords["dword_5d4594_1062512"] = 125
			switch gate {
			case 0:
				*o.windowWords["dword_5d4594_1047520"] = 1
			case 1:
				*memmap.PtrUint8(0x5D4594, 1049868) = 0
			case 2:
				*memmap.PtrUint8(0x5D4594, 1049868) = 1
			case 3:
				*memmap.PtrUint8(0x5D4594, 1049868) = 3
			case 4:
				*txword(o.items[0], 276) = 1
			case 5:
				*txword(o.items[0], 276) = 2
			case 6:
				*txword(o.items[0], 276) = 51
			case 7:
				*memmap.PtrPtr(0x852978, 8) = nil
			}
			r := o.capture(t, int(event)*10+gate, 0, 1, txptr(o.mainWindow().C()), event, inventoryWindowPoint(339, 38), 0)
			rows = append(rows, r)
			want := uint32(1)
			if gate >= 4 && (event == 9 || event == 99) {
				want = 0
			}
			if r.Return != want {
				t.Fatalf("event%d gate%d return%d want%d", event, gate, r.Return, want)
			}
			if *memmap.PtrUint32(0x5D4594, 1049848) != 0 || *o.windowWords["dword_5d4594_1062512"] != 125 {
				t.Fatalf("event%d gate%d changed drag or scroll", event, gate)
			}
		}
	}
	inventoryWindowCapture(t, "input-gates", rows, "9bcc687b976ce445ac2861ab2d842c528c98108b480fdb70d315dc44aea496a0")
}

func TestClientInventoryWindowButtonAndTooltip(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	for state := 0; state < 4; state++ {
		for _, event := range []uintptr{0, 4, 5, 6, 7, 8, 9, 19, 20, 0xffffffff} {
			o.reset(t)
			o.construct(t)
			*memmap.PtrUint8(0x5D4594, 1049868) = byte(state)
			rows = append(rows, o.capture(t, state*100+int(event), 0, 2, 0, 0, 0, 0))
			r := o.capture(t, state*100+int(event), 1, 18, txptr(o.mainWindow().C()), event, 0, 0)
			rows = append(rows, r)
			want := uint32(0)
			if event >= 5 && event <= 7 {
				want = 1
			}
			if r.Return != want {
				t.Fatalf("button event%d return%d want%d", event, r.Return, want)
			}
		}
	}
	inventoryWindowCapture(t, "button-tooltip", rows, "16c23252bcbd5cef77a982fb91a77e96fab6d59c0a9a414312ee0e26ca5c64d4")
}

// Click, secondary-click, release and hover at every edge of the actual embedded
// interaction rectangles, including one pixel outside each edge.
func TestClientInventoryWindowInputBoundaries(t *testing.T) {
	o := newInventoryWindowOwner(t)
	var rows []inventoryWindowResult
	points := [][2]int{{313, 13}, {314, 12}, {314, 13}, {513, 162}, {514, 162}, {513, 163}, {253, 13}, {254, 13}, {303, 62}, {303, 63}, {303, 112}, {303, 113}, {303, 162}, {304, 162}, {10, 15}, {11, 15}, {210, 214}, {211, 214}, {0, 0}, {65535, 65535}}
	for _, mode := range []uint32{0, 5, 6} {
		for _, event := range []uintptr{5, 6, 7, 9, 19, 20, 99} {
			for _, point := range points {
				o.reset(t)
				o.construct(t)
				o.openInput()
				*o.windowWords["dword_5d4594_1049864"] = mode
				o.stack(t, 0, 0, 1, "RedApple", 100)
				rows = append(rows, o.capture(t, len(rows), 0, 19, txptr(o.mainWindow().C()), 0, inventoryWindowPoint(point[0], point[1]), 0))
				rows = append(rows, o.capture(t, len(rows), 1, 1, txptr(o.mainWindow().C()), event, inventoryWindowPoint(point[0], point[1]), 0))
				// Each next case restores all ownership; an edge press may legitimately drag.
			}
		}
	}
	inventoryWindowCapture(t, "input-boundaries", rows, "6235def5ad55c3d79ceb968f6a553b065e45438d8983b063829390dbd35c49e6")
}
