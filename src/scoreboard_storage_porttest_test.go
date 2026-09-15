//go:build porttest

package opennox

import (
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client/gui"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestScoreboardConstructionLifetime(t *testing.T) {
	o := newScoreboardOwner(t)
	for _, op := range []int{17, 18} {
		if got := legacy.PortTestScoreboard(op, 0, 0, 0); got != 0 {
			t.Fatalf("empty callback%d returned%d", op, got)
		}
	}
	for i := 0; i < 100; i++ {
		o.resetRank(t)
		o.constructRank(t)

		if resp := o.rankWindow().Func93(&gui.RawEvent{Event: 23}); resp != nil {
			t.Fatal("rank input callback must return nil for zero")
		}
		if resp := (*gui.Window)(*memmap.PtrPtr(0x5D4594, 1090052)).Func94(&gui.RawEvent{Event: 23}); resp != nil {
			t.Fatal("rank column owner callback must return nil for zero")
		}
		for off := uintptr(1090052); off < 1090100; off += 4 {
			if *memmap.PtrUint32(0x5D4594, off) == 0 {
				t.Fatalf("cycle%d missing listbox %d", i, off)
			}
		}
		if *o.rankWords["dword_587000_145668"] != 6 || *o.rankWords["dword_587000_145664"] != 1 {
			t.Fatal("constructor mode count/refresh")
		}
		o.releaseRank()
		if o.rankWindow() != nil {
			t.Fatal("released scoreboard retained")
		}
	}
}
func TestScoreboardObjectiveFlags(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	flags := unsafe.Slice(memmap.PtrUint16(0x5D4594, 1090128), 3)
	for _, action := range []byte{0, 1, 2, 3, 4, 127, 128, 255} {
		for _, team := range []byte{0, 1, 2, 3, 127, 128, 255} {
			for _, code := range []uint16{0, 1, 127, 128, 255, 256, 32767, 32768, 65535} {
				flags[0], flags[1], flags[2] = 111, 222, 333
				got := legacy.PortTestScoreboard(27, uintptr(action), uintptr(team), uintptr(code))
				want := [3]uint16{111, 222, 333}
				if team == 1 || team == 2 {
					if action == 0 || action == 2 {
						want[team-1] = 0
					}
					if action == 1 {
						want[team-1] = code
					}
				}
				if [3]uint16(flags) != want {
					t.Fatalf("flag action%d/team%d/code%d got%v want%v", action, team, code, flags, want)
				}
				rows = append(rows, o.rankCapture(t, 27, got))
			}
		}
	}
	scoreboardCapture(t, "team-objective-flags", rows, "8234d0d689fc98d15348ab9196f1393afa2a185dc2248c5299783484be3ee1bb")
	rows = nil
	for _, action := range []byte{0, 1, 2, 3, 4, 127, 128, 255} {
		for _, code := range []uint16{0, 1, 127, 128, 255, 256, 32767, 32768, 65535} {
			flags[0], flags[1], flags[2] = 111, 222, 333
			got := legacy.PortTestScoreboard(28, uintptr(action), uintptr(code), 0)
			want := uint16(333)
			if action == 0 || action == 1 {
				want = 0
			}
			if action == 2 || action == 4 {
				want = code
			}
			if flags[0] != 111 || flags[1] != 222 || flags[2] != want {
				t.Fatalf("ball action%d/code%d got%v", action, code, flags)
			}
			rows = append(rows, o.rankCapture(t, 28, got))
		}
	}
	if got := legacy.PortTestScoreboard(29, 0, 0, 0); got != 0 || [3]uint16(flags) != [3]uint16{} {
		t.Fatal("clear objective ownership", got, flags)
	}
	rows = append(rows, o.rankCapture(t, 29, 0))
	scoreboardCapture(t, "ball-objective-flags", rows, "43378a273a7980a6df836a1c8e368b7586cda88d03f4eecd6dedfac67bd45abe")
}
func TestScoreboardMembership(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	for _, spec := range []struct{ op, count, base, stride, max int }{{13, 1090116, 1087248, 56, 9}, {15, 1090117, 1084192, 80, 32}} {
		for count := 0; count <= spec.max; count++ {
			*memmap.PtrUint8(0x5D4594, uintptr(spec.count)) = byte(count)
			for i := 0; i < spec.max; i++ {
				*memmap.PtrUint32(0x5D4594, uintptr(spec.base+i*spec.stride)) = uint32(100 + i*7)
			}
			for _, id := range []int{0, 99, 100, 107, 114, 156, 317, 324, -1} {
				got := legacy.PortTestScoreboard(spec.op, uintptr(id), 0, 0)
				want := uint32(0)
				for i := 0; i < count; i++ {
					if id == 100+i*7 {
						want = 1
					}
				}
				if got != want {
					t.Fatalf("op%d/count%d/id%d got%d want%d", spec.op, count, id, got, want)
				}
				rows = append(rows, o.rankCapture(t, spec.op, got))
			}
		}
	}
	scoreboardCapture(t, "membership", rows, "d4bc613f85624e67c582d0f48162514cc7621d2595493af9cf3a99d1567685b6")
}
func TestScoreboardNamesAndClasses(t *testing.T) {
	o := newScoreboardOwner(t)
	var rows []scoreboardResult
	ret := legacy.PortTestScoreboard(2, 0, 0, 0)
	for i, want := range []string{"Warrior", "Wizard", "Conjurer"} {
		got := alloc.GoString16(*(**uint16)(memmap.PtrOff(0x5D4594, 1084056+uintptr(4*i))))
		if got != want {
			t.Fatalf("class%d %q want%q", i, got, want)
		}
	}
	rows = append(rows, o.rankCapture(t, 2, ret))
	for _, s := range []string{"", "A", "Player", "A long name that will clip", "é界😀XY"} {
		for _, width := range []int{0, 5, 12, 19, 40, 80, 200} {
			p, free := alloc.Make([]uint16{}, 128)
			alloc.StrCopy16(p, s)
			*memmap.PtrUint32(0x5D4594, 1084036) = uint32(width)
			got := legacy.PortTestScoreboard(14, uintptr(unsafe.Pointer(&p[0])), 0, 0)
			out := alloc.GoString16(&p[0])
			t.Logf("name%q width%d ->%q return%v", s, width, out, got != 0)
			if width == 0 && out != s {
				t.Fatal("zero width changed name")
			}
			if width == 200 && out != s {
				t.Fatalf("wide name unexpectedly clipped %q", out)
			}
			// Store content in owned scratch for the frozen capture; addresses are local.
			copy(unsafe.Slice(memmap.PtrUint16(0x5D4594, 1084068), 32), p[:32])
			if got != 0 {
				got = 1
			}
			rows = append(rows, o.rankCapture(t, 14, got))
			free()
		}
	}
	for _, id := range []int{-1, 0, 1, 2, 3, 4, 5, 127, 255} {
		color, free := alloc.New(byte(0))
		*color = 99
		ret = legacy.PortTestScoreboard(5, uintptr(id), uintptr(unsafe.Pointer(color)), 0)
		got := alloc.GoString16((*uint16)(unsafe.Pointer(uintptr(ret))))
		want, shade := "\n", byte(4)
		switch id {
		case 1:
			want = "<King"
		case 2:
			want = "<Flag"
			shade = 7
		case 3:
			want = "<Flag"
			shade = 13
		case 4:
			want = "<Ball"
		}
		if got != want || *color != shade {
			t.Fatal(fmt.Sprintf("status%d %q/%d want%q/%d", id, got, *color, want, shade))
		}
		rows = append(rows, o.rankCapture(t, 5, uint32(*color)))
		free()
	}
	scoreboardCapture(t, "names-classes-status", rows, "37056a6b41f29c8a2582fd1d26117c0d3425f6e9a4c633d53a9960d55eba8790")
}
