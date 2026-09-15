//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"image"
	"sort"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/client/gui"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

// Register the real linked drawable with the reused allocation/lifetime checker.
// Cleanup also runs on a failed contract, before the parent owner's pool check.
func (o *scoreboardOwner) rankSprite(t *testing.T, code uint32, pos image.Point) (*client.Drawable, func()) {
	t.Helper()
	dr := o.c.Nox_xxx_spriteLoadAdd_45A360_drawable(o.c.Things.TypeByID("Gold").Index(), pos)
	if dr == nil {
		t.Fatal("scoreboard drawable allocation")
	}
	dr.NetCode32, dr.ObjClass = code, 16
	index := len(o.objects)
	ref := uint32(0xec000001) + uint32(index)
	o.objects = append(o.objects, inventoryTransactionObject{Ref: ref, Live: true, Drawable: dr})
	o.identities[dr] = index
	o.c.refs[dr] = ref
	released := false
	release := func() {
		if released {
			return
		}
		released = true
		o.c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
		o.objects[index].Live = false
	}
	t.Cleanup(release)
	return dr, release
}

func TestScoreboardTeamsAndMembership(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows, renderRows []scoreboardResult
	for _, elimination := range []bool{false, true} {
		for _, count := range []int{0, 1, 2, 8} {
			o.resetRank(t)
			o.c.srv.Teams.Reset()
			o.c.srv.Teams.ActiveCnt = 0
			noxflags.SetGame(noxflags.GameModeCoopTeam)
			order := make([]int, count)
			for i := 0; i < count; i++ {
				tm := o.c.srv.Teams.Create(server.TeamID(i + 1))
				if tm == nil {
					t.Fatal("team allocation")
				}
				tm.SetNameAnd68(fmt.Sprintf("Team%d", i), 0)
				tm.Lessons = i%3 - 1
				tm.ColorInd = server.TeamColor(i + 1)
				order[i] = i
			}
			noxflags.ResetGame()
			if elimination {
				noxflags.SetGame(noxflags.GameModeElimination)
			}
			sort.SliceStable(order, func(a, b int) bool {
				a, b = order[a]%3-1, order[b]%3-1
				if elimination {
					return a < b
				}
				return a > b
			})
			var release []func()
			for i := 0; i < 6; i++ {
				p := &o.players[i]
				p.Active = 1
				p.PlayerInd = byte(i)
				p.Lessons = int32(6 - i)
				p.Field2140 = uint32(i)
				p.Field3680 = 0
				dr, free := o.rankSprite(t, p.NetCodeVal, image.Pt(100+i, 100))
				release = append(release, free)
				*(*byte)(unsafe.Add(unsafe.Pointer(p), 2251)) = byte(i % 3)
				dr.TeamVal.ID = server.TeamID(i % 4)
				if o.c.Objs.ByNetCodeDynamic(int(p.NetCodeVal)) != dr {
					t.Fatal("actual drawable registration")
				}

			}
			ret := legacy.PortTestScoreboard(11, 0, 0, 0)
			if int(*memmap.PtrUint8(0x5D4594, 1090116)) != count {
				t.Fatal("team count")
			}
			for row, ind := range order {
				off := uintptr(row * 56)
				id := *memmap.PtrUint32(0x5D4594, 1087248+off)
				score := int32(*memmap.PtrUint32(0x5D4594, 1087252+off))
				if id != uint32(ind+1) || score != int32(ind%3-1) {
					t.Fatalf("team row%d id%d/score%d want%d/%d", row, id, score, ind+1, ind%3-1)
				}
			}
			for row := 0; row < 6; row++ {
				off := uintptr(row * 80)
				id := *memmap.PtrUint32(0x5D4594, 1084192+off)
				ind := int(id) - 100
				want := uint32(ind % 4)
				if want == 0 || want > uint32(count) {
					want = 0xffffffff
				}
				if got := *memmap.PtrUint32(0x5D4594, 1084184+off); got != want {
					t.Fatalf("player%d team%x want%x", id, got, want)
				}
			}
			rows = append(rows, o.rankCapture(t, 11, ret))
			o.constructRank(t)
			*o.rankWords["dword_5d4594_1090120"] = 2
			o.rankWindow().Show()
			o.c.GUI.Draw()
			d := (*gui.ScrollListBoxData)(o.rankColumn(0, 0).WidgetData)
			colors := map[string]uint32{}
			for _, row := range unsafe.Slice(d.Items, int(d.Field_11_0)) {
				colors[alloc.GoString16S(row.Text[:])] = row.Field_129
			}
			for i := 0; i < 6; i++ {
				name := fmt.Sprintf("Player%d", i)
				got, ok := colors[name]
				if !ok {
					t.Fatalf("rendered player row missing: %s", name)
				}
				team := i % 4
				want := **(**uint32)(memmap.PtrOff(0x85B3FC, 132+3*4))
				if team > 0 && team <= count {
					var found bool
					want, found = colors[fmt.Sprintf("Team%d", team-1)]
					if !found {
						t.Fatal("rendered team row missing", team)
					}
				}
				if got != want {
					t.Fatalf("player%d rendered color%x want%x", i, got, want)
				}
			}
			renderRows = append(renderRows, o.rankCapture(t, 3, 1))

			for _, free := range release {
				free()
			}
		}
	}
	scoreboardCapture(t, "teams-membership", rows, "7be4586487abb2216d9224a13fb63884521f78680e7e826e793c1c48cedf140a")
	scoreboardCapture(t, "team-player-rendering", renderRows, "f4f702c69251b187038937efeed51efb3bf0be027dc8a046a809bbf142d6b6ca")
}
func TestScoreboardHeadlessHost(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	for _, mode := range []int{2, 5} {
		for _, elimination := range []bool{false, true} {
			for _, count := range []int{1, 2, 16, 32} {
				o.resetRank(t)
				noxflags.SetGame(noxflags.GameHost)
				noxflags.SetEngine(noxflags.EngineNoRendering)
				if elimination {
					noxflags.SetGame(noxflags.GameModeElimination)
				}
				*o.rankWords["dword_5d4594_1090120"] = uint32(mode)
				for i := 0; i < count-1; i++ {
					p := &o.players[i]
					p.Active = 1
					p.PlayerInd = byte(i)
					p.Field3680 = 0
					p.Lessons = int32(i)
					p.Field2140 = uint32(i)
					p.Field2108 = uint32(i + 1)
				}
				host := &o.players[31]
				host.Active = 1
				host.PlayerInd = 31
				host.Lessons = 123
				host.Field2140 = 123
				host.Field2108 = 123
				host.Field3680 = 1
				ret := legacy.PortTestScoreboard(11, 0, 0, 0)
				if got := int(*memmap.PtrUint8(0x5D4594, 1090117)); got != count-1 {
					t.Fatalf("headless count%d got%d", count, got)
				}
				seen := map[uint32]bool{}
				for row := 0; row < count-1; row++ {
					b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 1084132+uintptr(row*80))), 80)
					id := binary.LittleEndian.Uint32(b[60:])
					if id == host.NetCodeVal || id < 100 || id >= uint32(100+count-1) || seen[id] {
						t.Fatalf("headless row%d id%d", row, id)
					}
					seen[id] = true
				}
				if host.Lessons != 123 || host.Field2140 != 123 || host.Field2108 != 123 {
					t.Fatal("headless host temporary score not restored")
				}
				rows = append(rows, o.rankCapture(t, 11, ret))
				noxflags.UnsetEngine(noxflags.EngineNoRendering)
			}
		}
	}
	scoreboardCapture(t, "headless-host", rows, "d1d3f54c12cb9a858871480dd1f7f27738efd65059e31fed7d01d1bd99768f05")
}
func TestScoreboardObjectiveIcons(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	for _, mode := range []noxflags.GameFlag{0, noxflags.GameModeKOTR, noxflags.GameModeCTF, noxflags.GameModeFlagBall, noxflags.GameModeCTF | noxflags.GameModeFlagBall | noxflags.GameModeKOTR} {
		for _, buff := range []bool{false, true} {
			for _, owner := range []uint16{0, 100, 101} {
				o.resetRank(t)
				noxflags.SetGame(mode)
				p := &o.players[0]
				p.Active = 1
				dr, release := o.rankSprite(t, 100, image.Pt(100, 100))
				if buff {
					dr.Buffs = 1 << 30
				}
				*memmap.PtrUint16(0x5D4594, 1090128) = owner
				*memmap.PtrUint16(0x5D4594, 1090130) = 100
				*memmap.PtrUint16(0x5D4594, 1090132) = owner
				want := uint32(0)
				if mode.Has(noxflags.GameModeCTF) {
					want = 3
					if owner == 100 {
						want = 2
					}
				} else if mode.Has(noxflags.GameModeFlagBall) {
					if owner == 100 {
						want = 4
					}
				} else if mode.Has(noxflags.GameModeKOTR) && buff {
					want = 1
				}
				ret := legacy.PortTestScoreboard(12, uintptr(unsafe.Pointer(p)), 0, 0)
				if ret != want {
					t.Fatalf("mode%x buff%v owner%d got%d want%d", mode, buff, owner, ret, want)
				}
				rows = append(rows, o.rankCapture(t, 12, ret))
				release()
			}
		}
	}
	scoreboardCapture(t, "objective-icons", rows, "e091d8250582b64a76d3ef770c03b87ce30c24932c4675ab711a4a347cdea048")
}
