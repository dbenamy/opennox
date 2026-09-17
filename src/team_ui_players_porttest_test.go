//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"image/color"
	"reflect"
	"strings"
	"testing"
	"unsafe"
)

func (o *teamUIOwner) openPlayers(t *testing.T) *gui.Window {
	t.Helper()
	clear(o.players) // Start with the two real teams and an empty player roster.
	result := teamUICall("players-construct", int(uintptr(o.parent.C())), 0)
	w := o.window("players")
	if w == nil || result != int(uintptr(w.C())) {
		t.Fatalf("constructor result %d, window %p", result, w)
	}
	if result := teamUICall("players-construct", 0, 0); result != 0 {
		t.Fatalf("repeat construction: %d", result)
	}
	return w
}
func teamUIRowNames(w *gui.Window) []string {
	d := (*gui.ScrollListBoxData)(w.WidgetData)
	var names []string
	for _, row := range unsafe.Slice(d.Items, int(d.Field_11_0)) {
		names = append(names, alloc.GoString16S(row.Text[:]))
	}
	return names
}
func TestTeamUIPlayerListRefresh(t *testing.T) {
	o := newTeamUIOwner(t)
	w := o.openPlayers(t)
	for pass := 0; pass < 3; pass++ {
		if pass > 0 {
			teamUICall("players-refresh", 0, 0)
		}
		names := teamUIRowNames(w.ChildByID(10502))
		rows := legacy.PortTestTeamUIRows(true)
		if len(names) != 2 || names[0] != "Red team" || names[1] != "Blue team" {
			t.Fatalf("pass %d displayed teams %v", pass, names)
		}
		if len(rows) != len(names) {
			t.Fatalf("pass %d: %d backing rows for %d displayed teams", pass, len(rows), len(names))
		}
		for i, row := range rows {
			if row.Name != names[i] {
				t.Errorf("row %d metadata %q display %q", i, row.Name, names[i])
			}
		}
	}
}
func TestTeamUIPlayerListTeamNames(t *testing.T) {
	o := newTeamUIOwner(t)
	o.openPlayers(t)
	for i, name := range []string{"Red team", "Blue team"} {
		got := legacy.PortTestTeamUI("team-find", nil, nil, 0, 0, name)
		if got != i {
			t.Errorf("lookup %q=%d, want %d", name, got, i)
		}
		if got := legacy.PortTestTeamUISelectedName(i); got != name {
			t.Errorf("selected row %d name %q, want %q", i, got, name)
		}
	}
}
func TestTeamUIPlayerListRename(t *testing.T) {
	for _, selected := range []int{-1, 0, 1} {
		t.Run(fmt.Sprint(selected), func(t *testing.T) {
			o := newTeamUIOwner(t)
			w := o.openPlayers(t)
			list := w.ChildByID(10502)
			list.Func94(&gui.RawEvent{Event: 16403, Arg1: uintptr(selected)})
			tm := o.c.srv.Teams.ByID(2)
			legacy.PortTestTeamUI("team-rename", nil, tm, 0, 0, "Azure team")
			names := teamUIRowNames(list)
			if !reflect.DeepEqual(names, []string{"Red team", "Azure team"}) {
				t.Errorf("rename team 2 with selected row %d: %v", selected, names)
			}
			rows := legacy.PortTestTeamUIRows(true)
			if len(rows) != 2 || rows[1].Name != "Azure team" {
				t.Errorf("rename left stale metadata: %+v", rows)
			}
			if got := int(int32((*gui.ScrollListBoxData)(list.WidgetData).Field_12)); got != selected {
				t.Errorf("selection changed %d -> %d", selected, got)
			}
		})
	}
}
func TestTeamUIPlayerListMutations(t *testing.T) {
	o := newTeamUIOwner(t)
	w := o.openPlayers(t)
	list := w.ChildByID(10501)
	names := []string{"Ada", "Ben Player", "Zoë"}
	for i, name := range names {
		legacy.PortTestTeamUI("player-add", nil, nil, 70+i, 0, name)
	}
	if got := teamUIRowNames(list); !reflect.DeepEqual(got, names) {
		t.Fatalf("added player rows %v", got)
	}
	for i, name := range names {
		if got := teamUICall("player-find", 70+i, 0); got != i {
			t.Errorf("find %s=%d", name, got)
		}
	}
	var records []any
	for _, code := range []int{71, 999, 70, 72, 72} {
		teamUICall("player-remove", code, 0)
		rows := legacy.PortTestTeamUIRows(false)
		shown := teamUIRowNames(list)
		if len(rows) != len(shown) {
			t.Fatalf("remove %d: rows=%v, displayed=%v", code, rows, shown)
		}
		for i, row := range rows {
			if row.Name != shown[i] {
				t.Fatalf("row %d mismatch %q/%q", i, row.Name, shown[i])
			}
		}
		records = append(records, struct {
			Code  int
			Rows  []legacy.PortTestTeamUIRow
			Shown []string
		}{code, rows, shown})
	}
	if got := teamUIRowNames(list); len(got) != 0 {
		t.Errorf("players left %v", got)
	}
	spellbookCapture(t, "team-ui-player-mutations", records, "4db1b83effbbf65432e3362577b32f8ee88c7921c4a670460481225f81d049a1")
}

func TestTeamUIPlayerListRosterRefresh(t *testing.T) {
	o := newTeamUIOwner(t)
	w := o.openPlayers(t)
	for i, name := range []string{"First player", "Last player"} {
		p := &o.players[i]
		p.Active = 1
		p.PlayerInd = byte(i)
		p.NetCodeVal = uint32(7 + i)
		p.SetName(name)
	}
	for pass := 0; pass < 3; pass++ {
		teamUICall("players-refresh", 0, 0)
		shown := teamUIRowNames(w.ChildByID(10501))
		rows := legacy.PortTestTeamUIRows(false)
		if !reflect.DeepEqual(shown, []string{"First player", "Last player"}) || len(rows) != 2 {
			t.Fatalf("pass %d rows=%v display=%v", pass, rows, shown)
		}
		for i, row := range rows {
			if row.Name != shown[i] || row.Code != uint32(7+i) {
				t.Errorf("pass %d row %d: %+v", pass, i, row)
			}
		}
	}
}
func TestTeamUIPlayerListSelectionEvents(t *testing.T) {
	type row struct {
		Mode, Selected, Pending int
		Locked                  bool
		Join, Rename            bool
		Result                  int
	}
	var records []row
	for _, mode := range []int{0, 1, 128, 129, 0x8000, 0x8001} {
		for _, selected := range []int{-1, 0, 1} {
			for _, pending := range []int{0, 1, 2} {
				for _, locked := range []bool{false, true} {
					t.Run(fmt.Sprintf("%d/%d/%d/%v", mode, selected, pending, locked), func(t *testing.T) {
						o := newTeamUIOwner(t)
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(mode))()
						w := o.openPlayers(t)
						*memmap.PtrUint32(0x5D4594, 1045696) = uint32(pending)
						if locked {
							*memmap.PtrUint8(0x5D4594, 371380+53) = 128
						}
						list := w.ChildByID(10502)
						list.Func94(&gui.RawEvent{Event: 16403, Arg1: uintptr(selected)})
						result := legacy.PortTestTeamUIEvent(w, list, 16400, 0)
						join := w.ChildByID(10507).Flags.IsEnabled()
						rename := w.ChildByID(10509).Flags.IsEnabled()
						allowed := selected >= 0 && mode&0x8000 == 0 && !locked
						if join != (allowed && (mode&128 != 0 || pending == 0)) || rename != (allowed && mode&1 != 0) {
							t.Errorf("join=%v rename=%v", join, rename)
						}
						records = append(records, row{mode, selected, pending, locked, join, rename, result})
					})
				}
			}
		}
	}
	spellbookCapture(t, "team-ui-list-selection-events", records, "e229e8600ceb5397f6c41b814db9df22c7b201760832746ed87e173beae234b7")
}
func TestTeamUIPlayerListFormatting(t *testing.T) {
	type record struct {
		Mode, Settings int
		Names          []string
		Rows           []legacy.PortTestTeamUIRow
		Colors         []uint32
	}
	var records []record
	for _, mode := range []int{0, 32, 64, 96, 128} {
		for _, settings := range []int{0, 32, 64, 96, 128} {
			t.Run(fmt.Sprintf("%d/%d", mode, settings), func(t *testing.T) {
				o := newTeamUIOwner(t)
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(mode))()
				*memmap.PtrUint8(0x5D4594, 371380+52) = byte(settings)
				w := o.openPlayers(t)
				list := w.ChildByID(10502)
				names := teamUIRowNames(list)
				want := []string{"Red team", "Blue team"}
				if mode&96 != 0 || settings&96 != 0 {
					want[0] += "\tred flag"
					want[1] += "\tblue flag"
				}
				if !reflect.DeepEqual(names, want) {
					t.Fatalf("formatted rows %v want %v", names, want)
				}
				data := (*gui.ScrollListBoxData)(list.WidgetData)
				var colors []uint32
				for i, item := range unsafe.Slice(data.Items, int(data.Field_11_0)) {
					color := **(**uint32)(memmap.PtrOff(0x85B3FC, 132+uintptr(4*(i+2))))
					if item.Field_129 != color {
						t.Errorf("row %d color %#x want %#x", i, item.Field_129, color)
					}
					colors = append(colors, item.Field_129)
					if got := legacy.PortTestTeamUISelectedName(i); got != []string{"Red team", "Blue team"}[i] {
						t.Errorf("bare selection %q", got)
					}
				}
				records = append(records, record{mode, settings, names, legacy.PortTestTeamUIRows(true), colors})
			})
		}
	}
	spellbookCapture(t, "team-ui-list-format", records, "528cb058213d41ba25479f5c293a7100f4009d6e97a7144dc6fc59ef2c408d3e")
}
func TestTeamUIPlayerListColors(t *testing.T) {
	o := newTeamUIOwner(t)
	tm := o.c.srv.Teams.ByID(1)
	var records []int
	for color := 0; color < 256; color++ {
		tm.ColorInd = server.TeamColor(color)
		got := legacy.PortTestTeamUI("team-color", nil, tm, 0, 0, "")
		if got != color%10+1 {
			t.Fatalf("color %d mapped to %d", color, got)
		}
		records = append(records, got)
	}
	spellbookCapture(t, "team-ui-list-colors", records, "9347f27317796d33bf7d252c232f425b1ceb517bd128fc533fec9214e2c3f3fb")
}
func TestTeamUIPlayerListDrawing(t *testing.T) {
	type record struct {
		Mode, Selected, Player, Background int
		Assign, Join, Rename               bool
		Pixels                             string
	}
	var records []record
	for _, mode := range []int{0, 1, 129, 0x8001} {
		for _, selected := range []int{-1, 0} {
			for _, player := range []int{-1, 0} {
				for background := 0; background < 3; background++ {
					t.Run(fmt.Sprintf("%d/%d/%d/%d", mode, selected, player, background), func(t *testing.T) {
						o := newTeamUIOwner(t)
						defer noxflags.PortTestGameFlags(noxflags.GameFlag(mode))()
						w := o.openPlayers(t)
						legacy.PortTestTeamUI("player-add", nil, nil, 777, 0, "Player")
						w.ChildByID(10502).Func94(&gui.RawEvent{Event: 16403, Arg1: uintptr(selected)})
						w.ChildByID(10501).Func94(&gui.RawEvent{Event: 16403, Arg1: uintptr(player)})
						for _, id := range []uint{10503, 10507, 10509} {
							w.ChildByID(id).Flags |= gui.StatusEnabled
						}
						w.SetPos(image.Pt(5, 5))
						w.DrawData().SetBackgroundColor(color.Transparent)
						if background == 1 {
							w.DrawData().SetBackgroundColor(color.White)
						}
						if background == 2 {
							w.Flags |= gui.StatusImage
							w.DrawData().SetBackgroundImage(o.images[3])
						}
						for i := range o.pix.Pix {
							o.pix.Pix[i] = 0xffff
						}
						before := effectsPixelHash(o.pix)
						if ret := legacy.PortTestTeamUIDraw("players", w); ret != 1 {
							t.Errorf("draw return %d", ret)
						}
						pix := effectsPixelHash(o.pix)
						if (pix != before) != (background != 0) {
							t.Errorf("background %d draw effect", background)
						}
						assign := w.ChildByID(10503).Flags.IsEnabled()
						wantAssign := true
						if mode&1 != 0 && mode&0x8000 == 0 {
							wantAssign = selected >= 0 && player >= 0
						}
						if assign != wantAssign {
							t.Errorf("assign enabled %v want %v", assign, wantAssign)
						}
						join, rename := w.ChildByID(10507).Flags.IsEnabled(), w.ChildByID(10509).Flags.IsEnabled()
						if join != (selected >= 0) || rename != (selected >= 0) {
							t.Errorf("empty selection button clearing: %v %v", join, rename)
						}
						records = append(records, record{mode, selected, player, background, assign, join, rename, pix})
					})
				}
			}
		}
	}
	spellbookCapture(t, "team-ui-list-drawing", records, "b61f4cfa9959dfc87caaab5ed78e2d848452d2553a6d018a1109e88df44acb10")
}
func TestTeamUITeamRemovalAndRefresh(t *testing.T) {
	type record struct {
		Step    string
		Teams   []string
		Rows    []legacy.PortTestTeamUIRow
		Players []string
	}
	var records []record
	for _, order := range [][2]int{{1, 2}, {2, 1}} {
		t.Run(fmt.Sprint(order), func(t *testing.T) {
			o := newTeamUIOwner(t)
			defer noxflags.PortTestGameFlags(1)()
			w := o.openPlayers(t)
			p := &o.players[0]
			p.Active = 1
			p.PlayerInd = 0
			p.NetCodeVal = 7
			p.SetName("Active player")
			capture := func(step string) {
				names := teamUIRowNames(w.ChildByID(10502))
				rows := legacy.PortTestTeamUIRows(true)
				if len(names) != len(rows) {
					t.Fatalf("%s metadata/display mismatch", step)
				}
				players := teamUIRowNames(w.ChildByID(10501))
				if !reflect.DeepEqual(players, []string{"Active player"}) || len(legacy.PortTestTeamUIRows(false)) != 1 {
					t.Fatalf("%s player refresh mismatch: %v", step, players)
				}
				records = append(records, record{step, names, rows, players})
			}
			for i := 0; i < 2; i++ {
				teamUICall("team-clear", 0, 0)
				capture(fmt.Sprintf("clear%d", i))
			}
			for _, id := range order {
				o.c.srv.TeamRemove(o.c.srv.Teams.ByID(server.TeamID(id)), false)
				if got := len(teamUIRowNames(w.ChildByID(10502))); got != o.c.srv.Teams.Count() {
					t.Fatalf("remaining teams %d vs %d", got, o.c.srv.Teams.Count())
				}
				capture(fmt.Sprintf("remove%d", id))
			}
			legacy.PortTestTeamUI("team-remove", nil, nil, 0, 0, "missing")
			capture("remove missing")
		})
	}
	spellbookCapture(t, "team-ui-team-removal-refresh", records, "e687b863002fe83e5a27c562c2dc4e1abe94ac0a6d8639720fb71691b1d9ea4a")
}
func TestTeamUINameBoundaries(t *testing.T) {
	type record struct {
		Kind, Name string
		Rows       []string
		Metadata   []legacy.PortTestTeamUIRow
		Found      int
	}
	var records []record
	for _, name := range []string{"A", "abcdefghijklmnopqrstu", "Two words Ω", "Короткое имя"} {
		t.Run(name, func(t *testing.T) {
			o := newTeamUIOwner(t)
			w := o.openPlayers(t)
			tm := o.c.srv.Teams.ByID(1)
			legacy.PortTestTeamUI("team-rename", nil, tm, 0, 0, name)
			rows := teamUIRowNames(w.ChildByID(10502))
			metadata := legacy.PortTestTeamUIRows(true)
			found := legacy.PortTestTeamUI("team-find", nil, nil, 0, 0, strings.Map(func(r rune) rune {
				if r >= 'a' && r <= 'z' {
					return r - 'a' + 'A'
				}
				return r
			}, name))
			if found != 0 || rows[0] != name || metadata[0].Name != name {
				t.Fatalf("full name round trip: %d %v %v", found, rows, metadata)
			}
			records = append(records, record{"team", name, rows, metadata, found})
		})
	}
	for _, name := range []string{"", "A", "abcdefghijklmnopqrstuvw", "Player Ω", "Игрок"} {
		t.Run("player/"+name, func(t *testing.T) {
			o := newTeamUIOwner(t)
			w := o.openPlayers(t)
			legacy.PortTestTeamUI("player-add", nil, nil, 77, 0, name)
			rows := teamUIRowNames(w.ChildByID(10501))
			metadata := legacy.PortTestTeamUIRows(false)
			found := teamUICall("player-find", 77, 0)
			if found != 0 || len(rows) != 1 || rows[0] != name || len(metadata) != 1 || metadata[0].Name != name {
				t.Fatalf("player name round trip: %d %v %v", found, rows, metadata)
			}
			records = append(records, record{"player", name, rows, metadata, found})
		})
	}
	spellbookCapture(t, "team-ui-name-boundaries", records, "cda5433ef0dcb489e173e76e3888c388993f258ac5ce307d6b5a4d859f37832e")
}
