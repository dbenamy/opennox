//go:build porttest

package opennox

import (
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestClientInteractionClocksAndCursor(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	clock := serverConfigOwnBytes(t, 0x5D4594, 811908, 8)
	oldClock := legacy.PlatformTicks
	defer func() { legacy.PlatformTicks = oldClock }()
	type tickRow struct{ Input, Return, Stored uint64 }
	var ticks []tickRow
	for _, now := range []uint64{0, 1, 0x7fffffff, 0x80000000, 0xffffffff, 0x100000000, 0x123456789abcdef0} {
		legacy.PlatformTicks = func() uint64 { return now }
		for i := range clock {
			clock[i] = 0xff
		}
		got := interactionCall("nox_xxx_initTime_435570")
		stored := binary.LittleEndian.Uint64(clock)
		if got != uint64(uint32(now)) || stored != got {
			t.Fatal("tick widening", now, got, stored)
		}
		ticks = append(ticks, tickRow{now, got, stored})
	}
	o.c.Inp.SetMouseBounds(image.Rect(0, 0, 100, 100))
	for _, p := range []image.Point{{25, 30}, {40, 20}, {1, 1}, {-1, 101}} {
		interactionCall("nox_client_setMousePos_430B10", uintptr(uint32(p.X)), uintptr(uint32(p.Y)))
		want := p
		if want.X < 0 {
			want.X = 0
		}
		if want.Y > 100 {
			want.Y = 100
		}
		if got := o.c.Inp.GetMousePos(); got != want {
			t.Fatal("absolute mouse position", p, got, want)
		}
	}
	cursor := serverConfigOwnBytes(t, 0x5D4594, 1096672, 4)
	state := serverConfigOwnBytes(t, 0x5D4594, 1064848, 4)
	for _, value := range []uint32{0, 1, 0x80000000, 0xffffffff} {
		binary.LittleEndian.PutUint32(cursor, value)
		binary.LittleEndian.PutUint32(state, value)
		if uint32(interactionCall("nox_xxx_guiCursor_477600")) != value || uint32(interactionCall("sub_469FA0")) != value {
			t.Fatal("cursor/state getter", value)
		}
	}
	dr, free := alloc.New(client.Drawable{})
	defer free()
	for _, code := range []uint32{0, 1, 0x7fff, 0x8000, 0x10001, 0xffffffff} {
		for _, static := range []bool{false, true} {
			dr.NetCode32 = code
			dr.ObjClass = 0
			if static {
				dr.ObjClass = 0x400000
			}
			*words["dword_5d4594_1096640"] = uint32(uintptr(unsafe.Pointer(dr)))
			want := code
			if code >= 0x8000 {
				want = 0
			} else if static {
				want |= 0x8000
			}
			if uint32(interactionCall("nox_xxx_packetGetMarshall_476F40")) != want {
				t.Fatal("hover unit code", code, static)
			}
		}
	}
	*words["dword_5d4594_1096640"] = 0
	if interactionCall("nox_xxx_packetGetMarshall_476F40") != 0 {
		t.Fatal("no hover unit")
	}
	pl, freePlayer := alloc.New(server.Player{})
	defer freePlayer()
	w := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 10, 10, nil)
	defer w.Destroy()
	*words["dword_5d4594_1193712"] = uint32(uintptr(w.C()))
	for _, present := range []bool{false, true} {
		for _, flags := range []uint32{0, 1, 2, 3, 0x100, 0xffffffff} {
			*words["dword_8531A0_2576"] = 0
			if present {
				*words["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(pl)))
			}
			pl.Field3680 = flags
			interactionCall("nox_xxx_cliToggleObsWindow_4357A0")
			if w.GetFlags().IsHidden() == (present && flags&1 != 0) {
				t.Fatal("observer indicator gate", present, flags)
			}
		}
	}
	interactionCapture(t, "clock", ticks)
}

func TestClientInteractionModalFrameGate(t *testing.T) {
	o := sessionDialogResources(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	sessions, restoreSessions := legacy.PortTestSessionDialogWords()
	defer restoreSessions()
	options, restoreOptions := legacy.PortTestServerOptionsWords()
	defer restoreOptions()
	scoreboard, restoreScoreboard := legacy.PortTestScoreboardWords()
	defer restoreScoreboard()
	stamp := serverConfigOwnBytes(t, 0x5D4594, 811920, 4)
	oldFrame := o.c.srv.Frame()
	defer o.c.srv.SetFrame(oldFrame)
	quit := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 10, 10, nil)
	defer quit.Destroy()
	motd := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 10, 10, nil)
	defer motd.Destroy()
	*sessions["quit"] = uint32(uintptr(quit.C()))
	*sessions["motd"] = uint32(uintptr(motd.C()))
	oldConsole := legacy.Nox_gui_console_flagXxx_451410
	defer func() { legacy.Nox_gui_console_flagXxx_451410 = oldConsole }()
	var consoleCalls, consoleValue int
	legacy.Nox_gui_console_flagXxx_451410 = func() int { consoleCalls++; return consoleValue }
	type row struct {
		Gate                       int
		Frame, Old, Return, Stored uint32
		ConsoleCalls               int
	}
	var rows []row
	for gate := 0; gate < 7; gate++ {
		for _, frame := range []uint32{0, 1, 2, 3, 0x80000000, 0xffffffff} {
			for _, old := range []uint32{0, 1, 2, frame, frame - 1, frame - 2, 0xffffffff} {
				*options["root"] = 0
				quit.SetHidden(true)
				motd.SetHidden(true)
				*words["dword_5d4594_1305680"] = 0
				*scoreboard["dword_5d4594_1090048"] = 0
				*scoreboard["dword_5d4594_1090120"] = 0
				consoleValue = 0
				consoleCalls = 0
				switch gate {
				case 1:
					*options["root"] = uint32(uintptr(quit.C()))
				case 2:
					quit.Show()
				case 3:
					*words["dword_5d4594_1305680"] = uint32(uintptr(quit.C()))
				case 4:
					motd.Show()
				case 5:
					*scoreboard["dword_5d4594_1090048"] = uint32(uintptr(quit.C()))
					*scoreboard["dword_5d4594_1090120"] = 1
				case 6:
					consoleValue = 1
				}
				o.c.srv.SetFrame(frame)
				binary.LittleEndian.PutUint32(stamp, old)
				got := uint32(interactionCall("sub_436550"))
				stored := binary.LittleEndian.Uint32(stamp)
				want, keep := uint32(1), frame
				if gate == 0 && frame != 2 {
					want = 0
					if frame-old == 1 {
						want = 1
					}
					keep = old
				}
				calls := 0
				if gate == 0 || gate == 6 {
					calls = 1
				}
				if got != want || stored != keep || consoleCalls != calls {
					t.Fatal("modal frame gate", gate, frame, old, got, stored, consoleCalls, want, keep, calls)
				}
				rows = append(rows, row{gate, frame, old, got, stored, consoleCalls})
			}
		}
	}
	interactionCapture(t, "modal-frame", rows)
}
