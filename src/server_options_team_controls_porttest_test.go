//go:build porttest

package opennox

import (
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestServerOptionsTeamControls(t *testing.T) {
	type row struct {
		ID            uint
		Before, After uint
		Count         int
		Enabled       []bool
		Queue         []legacy.PortTestShopPacketResult
	}
	var rows []row
	for _, id := range []uint{10330, 10331, 10332, 10333} {
		for before := uint(0); before < 8; before++ {
			t.Run(fmt.Sprintf("%d-%d", id, before), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				t.Cleanup(o.c.srv.PortTestMapDrawableTeamMessages())
				defer noxflags.PortTestGameFlags(0)()
				old := noxflags.GetGamePlay()
				noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0))
				noxflags.SetGamePlay(noxflags.GameplayFlag(before))
				t.Cleanup(func() { noxflags.UnsetGamePlay(^noxflags.GameplayFlag(0)); noxflags.SetGamePlay(old) })
				if legacy.PortTestServerOptionsEvent(o.options, o.options.ChildByID(id), 16391, 0) != 1 {
					t.Fatal("team control return")
				}
				want := before
				count := 2
				switch id {
				case 10330:
					if before&4 != 0 {
						want &^= 4
						count = 0
					} else {
						want |= 4
						count = 4
					}
				case 10331:
					want ^= 2
				case 10333:
					want ^= 1
				}
				if uint(noxflags.GetGamePlay()) != want || o.c.srv.Teams.Count() != count {
					t.Fatalf("flags %d count %d want %d/%d", noxflags.GetGamePlay(), o.c.srv.Teams.Count(), want, count)
				}
				r := row{ID: id, Before: before, After: want, Count: count, Queue: o.packets()}
				for _, child := range []uint{10331, 10332, 10333} {
					enabled := o.options.ChildByID(child).Flags.IsEnabled()
					if id == 10330 && enabled != (before&4 == 0) {
						t.Fatal("team child enablement")
					}
					r.Enabled = append(r.Enabled, enabled)
				}
				rows = append(rows, r)
			})
		}
	}
	spellbookCapture(t, "server-options-team-controls", rows, "4d6b6061f5638350e7879cc7d5fcbbe453170cc1a56aa3bbb0216b0b08fdf3fd")
}
func TestServerOptionsTeamCount(t *testing.T) {
	type row struct {
		Flags, Count int
		Text         string
	}
	var rows []row
	for _, flags := range []int{0, 0x8000} {
		for _, count := range []int{0, 1, 2, 127, 128, 255, 256, 257} {
			t.Run(fmt.Sprintf("%x-%d", flags, count), func(t *testing.T) {
				o := newServerOptionsOwner(t)
				defer noxflags.PortTestGameFlags(noxflags.GameFlag(flags))()
				o.c.srv.Teams.ActiveCnt = count
				tm := o.c.srv.Teams.ByID(1)
				*(*uint32)(unsafe.Add(tm.C(), 60)) = 17
				o.call("team-count", 0, "")
				got := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(o.event(10110, 16386, 0, 0)))))
				want := int(byte(count))
				if flags != 0 {
					want = 1
				}
				if got != fmt.Sprintf("Teams: %d", want) {
					t.Fatalf("count label %q want %d", got, want)
				}
				rows = append(rows, row{flags, count, got})
			})
		}
	}
	spellbookCapture(t, "server-options-team-count", rows, "cc8fd8560879577aaedcdc9b8d0d65b8852ccadda387c1299c2081d89726ff47")
}
