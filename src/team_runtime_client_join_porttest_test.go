//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"image"
	"testing"
	"unicode/utf16"
	"unsafe"

	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
)

func TestTeamRuntimeClientJoinNotifications(t *testing.T) {
	o := newObjectDrawingOwner(t)
	t.Cleanup(o.c.srv.PortTestMapDrawableTeamMessages())
	oldClient, oldServer := noxClient, noxServer
	noxClient, noxServer = o.c.Client, o.c.srv
	t.Cleanup(func() { noxClient, noxServer = oldClient, oldServer })
	printer, index, region := teamRuntimeJoinTextOwner(t, o.c.srv.Server)
	players, free := o.c.srv.PortTestObjectRenderPlayers()
	t.Cleanup(free)
	pl := &players[0]
	pl.Active = 1
	pl.NetCodeVal = 777
	name := unsafe.Slice((*uint16)(unsafe.Add(pl.C(), 4704)), 25)
	clear(name)
	copy(name, utf16.Encode([]rune("Ada")))
	dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(10, 10))
	dr.ObjClass = 4
	dr.NetCode32 = 777
	cache := legacy.PortTestMapDrawableTeamWord()
	oldCache := *cache
	*cache = 1
	t.Cleanup(func() { *cache = oldCache })
	type row struct {
		Mode     int
		Text     string
		Deadline uint32
		Count    uint32
	}
	var rows []row
	for _, mode := range []int{0, 4096} {
		t.Run(fmt.Sprint(mode), func(t *testing.T) {
			defer noxflags.PortTestGameFlags(noxflags.GameFlag(mode))()
			o.c.srv.Teams.Reset()
			o.c.srv.Teams.ActiveCnt = 0
			tm := o.c.srv.Teams.Create(1)
			tm.SetNameAnd68("Red", 0)
			dr.TeamPtr().ID = 0
			dr.TeamPtr().Field0 = 0
			clear(region)
			*index = 0
			printer.lines = nil
			legacy.Nox_xxx_createAtImpl_4191D0(1, dr.TeamPtr(), 1, 777, 0)
			want := "Ada joined Red"
			if mode == 4096 {
				want = "Ada joined quest"
			}
			if len(printer.lines) != 1 || printer.lines[0] != want || *index != 1 {
				t.Fatal("client join text", printer.lines, *index)
			}
			if objectXferGetWord(tm.C(), 44) != uint32(uintptr(dr.TeamPtr().C())) || objectXferGetWord(tm.C(), 48) != 1 || dr.TeamPtr().ID != 1 {
				t.Fatal("client linked membership")
			}
			deadline := binary.LittleEndian.Uint32(region[1280:])
			if deadline != o.c.srv.Frame()+5*uint32(o.c.srv.TickRate()) {
				t.Fatal("client join deadline")
			}
			rows = append(rows, row{mode, printer.lines[0], deadline, objectXferGetWord(tm.C(), 48)})
		})
	}
	spellbookCapture(t, "team-runtime-client-join", rows, "a48b2279dc5ae1b25dcdb38c0dfa0c4e39094b3ea7d847ea53e1cc0163acff82")
}
