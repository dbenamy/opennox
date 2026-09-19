//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestClientInteractionActionMessages(t *testing.T) {
	o := newMeterOwner(t)
	words, restore := legacy.PortTestClientInteractionWords()
	defer restore()
	inv, restoreInv := legacy.PortTestInventoryWindowWords()
	defer restoreInv()
	sessions, restoreSession := legacy.PortTestSessionDialogWords()
	defer restoreSession()
	quit := o.c.GUI.NewWindowRaw(nil, 8, 0, 0, 10, 10, nil)
	defer quit.Destroy()
	*sessions["quit"] = uint32(uintptr(quit.C()))
	pl, freePlayer := alloc.New(server.Player{})
	defer freePlayer()
	dr, freeDrawable := alloc.New(client.Drawable{})
	defer freeDrawable()
	type row struct {
		Action                       int
		Flags, Active, Code          uint32
		Player, Object, Quit, Static bool
		Message                      []byte
	}
	var rows []row
	run := func(action int, flags, active, code uint32, player, object, shown, static bool) {
		t.Helper()
		o.c.srv.NetList.ResetAll()
		*words["dword_8531A0_2576"] = 0
		if player {
			*words["dword_8531A0_2576"] = uint32(uintptr(unsafe.Pointer(pl)))
		}
		pl.Field3680 = flags
		dr.NetCode32 = code
		dr.ObjClass = 0
		if static {
			dr.ObjClass = 0x20400000
		}
		*words["dword_5d4594_1123520"] = active
		*inv["dword_5d4594_1098624"] = active
		quit.SetHidden(!shown)
		var arg uintptr
		if object {
			arg = uintptr(unsafe.Pointer(dr))
		}
		ops := []string{"nox_xxx_clientTalk_42E7B0", "nox_xxx_clientCollideOrUse_42E810", "nox_xxx_clientTrade_42E850"}
		interactionCall(ops[action], arg)
		var got []byte
		packets := 0
		o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { got = append(got, b...); packets++; return false })
		allowed := object && (!player || flags&3 == 0) && (action == 1 || active != 1 && !shown)
		var want []byte
		if allowed {
			unit := uint16(code)
			if action != 0 {
				if code >= 0x8000 {
					unit = 0
				} else if static {
					unit |= 0x8000
				}
			}
			switch action {
			case 0:
				want = []byte{0xd0, 1, 0, 0}
				binary.LittleEndian.PutUint16(want[2:], unit)
			case 1:
				want = []byte{123, 0, 0}
				binary.LittleEndian.PutUint16(want[1:], unit)
			case 2:
				want = []byte{0xc9, 0x15, 0, 0}
				binary.LittleEndian.PutUint16(want[2:], unit)
			}
		}
		if !bytes.Equal(got, want) || packets > 1 {
			t.Fatalf("action%d flags%x active%d player%v object%v quit%v code%x static%v: %x want%x (%d packets)", action, flags, active, player, object, shown, code, static, got, want, packets)
		}
		rows = append(rows, row{action, flags, active, code, player, object, shown, static, got})
	}
	for action := 0; action < 3; action++ {
		for _, flags := range []uint32{0, 1, 2, 3, 4, 0x100, 0xffffffff} {
			for _, active := range []uint32{0, 1, 2, 0xffffffff} {
				for _, player := range []bool{false, true} {
					for _, object := range []bool{false, true} {
						for _, shown := range []bool{false, true} {
							run(action, flags, active, 0x1234, player, object, shown, false)
						}
					}
				}
			}
		}
		for _, code := range []uint32{0, 1, 0x7fff, 0x8000, 0xffff, 0x10001, 0xffffffff} {
			for _, static := range []bool{false, true} {
				run(action, 0, 0, code, true, true, false, static)
			}
		}
	}
	interactionCapture(t, "actions", rows)
}
