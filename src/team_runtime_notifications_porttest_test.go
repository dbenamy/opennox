//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unicode/utf16"
	"unsafe"

	"github.com/opennox/libs/console"
	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

type teamRuntimePrinter struct{ lines []string }

func (p *teamRuntimePrinter) Print(_ console.Color, s string) { p.lines = append(p.lines, s) }
func (p *teamRuntimePrinter) Printf(c console.Color, s string, args ...any) {
	p.Print(c, fmt.Sprintf(s, args...))
}

func teamRuntimeJoinTextOwner(t *testing.T, s *server.Server) (*teamRuntimePrinter, *uint32, []byte) {
	configure, restore := s.PortTestMeterStrings(
		strman.Entry{ID: "team.c:NewMember", Vals: []strman.Variant{{Str: "%s joined %s"}}},
		strman.Entry{ID: "GeneralPrint:PlayerJoinQuest", Vals: []strman.Variant{{Str: "%s joined quest"}}},
		strman.Entry{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "%s"}}},
	)
	t.Cleanup(restore)
	configure(0)
	printer := new(teamRuntimePrinter)
	con := console.New(printer)
	oldConsole := legacy.GetConsole
	legacy.GetConsole = func() *console.Console { return con }
	t.Cleanup(func() { legacy.GetConsole = oldConsole })
	index := legacy.PortTestTeamRuntimeTextIndex()
	oldIndex := *index
	t.Cleanup(func() { *index = oldIndex })
	region := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 823804)), 1932)
	oldRegion := bytes.Clone(region)
	t.Cleanup(func() { copy(region, oldRegion) })
	return printer, index, region
}

func TestTeamRuntimeJoinNotifications(t *testing.T) {
	o := newMatchRosterOwner(t)
	printer, index, region := teamRuntimeJoinTextOwner(t, o.s)
	u := &o.units[0]
	u.TypeInd = uint16(o.s.Types.ByID("PortCreatureMonster").Ind())
	pl := u.UpdateDataPlayer().Player
	name := unsafe.Slice((*uint16)(unsafe.Add(pl.C(), 4704)), 25)
	clear(name)
	copy(name, utf16.Encode([]rune("Ada")))
	pl.NetCodeVal = u.NetCode
	type row struct {
		Mode     int
		Notify   bool
		Text     []string
		Centered []byte
		Queue    legacy.PortTestReliableReportState
	}
	var rows []row
	for _, mode := range []int{1, 0x2001, 0x3001, 0x2081} {
		for _, notify := range []bool{false, true} {
			t.Run(fmt.Sprintf("mode%x/notify%t", mode, notify), func(t *testing.T) {
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(mode))()
				o.s.Teams.Reset()
				o.s.Teams.ActiveCnt = 0
				tm := o.s.Teams.Create(1)
				tm.SetNameAnd68("Red", 0)
				for i := range o.units {
					o.units[i].TeamVal = server.ObjectTeam{}
				}
				clear(region)
				*index = 0
				printer.lines = nil
				o.reset()
				n := 0
				if notify {
					n = 1
				}
				legacy.Nox_xxx_createAtImpl_4191D0(1, u.TeamPtr(), n, int(u.NetCode), 0)
				wantText := ""
				if mode&0x2000 != 0 {
					wantText = "Ada joined Red"
					if mode&4096 != 0 {
						wantText = "Ada joined quest"
					}
				}
				if wantText == "" {
					if len(printer.lines) != 0 || *index != 0 {
						t.Fatal("unexpected join text")
					}
				} else {
					if len(printer.lines) != 1 || printer.lines[0] != wantText || *index != 1 {
						t.Fatalf("join text %q index%d", printer.lines, *index)
					}
					want := make([]byte, 636)
					for i, c := range utf16.Encode([]rune(wantText)) {
						binary.LittleEndian.PutUint16(want[2*i:], c)
					}
					if !bytes.Equal(region[644:1280], want) {
						t.Fatal("centered text storage")
					}
					if binary.LittleEndian.Uint32(region[1280:]) != o.s.Frame()+5*uint32(o.s.TickRate()) {
						t.Fatal("centered text deadline")
					}
				}
				state := o.state()
				teamMessages := 0
				for _, node := range state.Nodes {
					if len(node.Data) > 1 && node.Data[0] == 196 {
						teamMessages++
						want := []byte{196, 1, 1, 0, 0, 0, 0, 0, 0, 0}
						binary.LittleEndian.PutUint16(want[6:], uint16(u.NetCode))
						binary.LittleEndian.PutUint16(want[8:], u.TypeInd)
						if !bytes.Equal(node.Data, want) {
							if !bytes.Equal(node.Data[:8], want[:8]) {
								t.Fatal("join membership header fields differ")
							}
							t.Fatalf("join membership object type %d, want %d", binary.LittleEndian.Uint16(node.Data[8:]), u.TypeInd)
						}
					}
				}
				wantMessages := 0
				if notify && mode&0x2000 != 0 {
					wantMessages = 1
				}
				if teamMessages != wantMessages {
					t.Fatal("join message filter")
				}
				rows = append(rows, row{mode, notify, append([]string(nil), printer.lines...), bytes.Clone(region), state})
			})
		}
	}
	spellbookCapture(t, "team-runtime-join-notifications", rows, "9bb9242c9b7eb2e1260b84648edeb00c0e73fd726c080df27f180c5804a3cbfc")
}
