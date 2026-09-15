//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"sort"
	"testing"
	"unsafe"
)

func TestBriefingWinReportPlayersAndFields(t *testing.T) {
	o := newBriefingOwner(t)
	packet, free := alloc.New([90]byte{})
	defer free()
	var rows []briefingResult
	for mask := 0; mask < 64; mask++ {
		o.resetBriefing(t)
		o.constructBriefing(t)
		clear(packet[:])
		packet[0], packet[1] = 240, 12
		binary.LittleEndian.PutUint16(packet[2:], 37)
		binary.LittleEndian.PutUint16(packet[4:], 65535)
		var slots []int
		for slot := 0; slot < 6; slot++ {
			if mask&(1<<slot) != 0 {
				slots = append(slots, slot)
			}
		}
		for i, slot := range slots {
			r := packet[6+i*14:]
			binary.LittleEndian.PutUint16(r, uint16(100+slot))
			for f := 1; f <= 4; f++ {
				binary.LittleEndian.PutUint16(r[f*2:], uint16(1000*slot+f))
			}
			binary.LittleEndian.PutUint32(r[10:], uint32(slot*31+7))
		}
		ret := legacy.PortTestBriefing(4, txptr(unsafe.Pointer(packet)), 0, 0)
		sort.Slice(slots, func(i, j int) bool { return slots[i] > slots[j] })
		for i := 0; i < 6; i++ {
			r := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 832364+uintptr(i*16))), 16)
			if i >= len(slots) {
				for _, v := range r {
					if v != 0 {
						t.Fatal("unused score slot changed")
					}
				}
				continue
			}
			slot := slots[i]
			if binary.LittleEndian.Uint32(r) != uint32(txptr(unsafe.Pointer(&o.players[slot]))) || binary.LittleEndian.Uint32(r[12:]) != uint32(slot*31+7) {
				t.Fatal("score order/player mapping")
			}
			for j, src := range []int{4, 1, 2, 3} {
				if binary.LittleEndian.Uint16(r[4+j*2:]) != uint16(1000*slot+src) {
					t.Fatalf("score field row%d field%d", i, j)
				}
			}
		}
		if memmap.Uint32(0x5D4594, 832356) != 37 || memmap.Uint32(0x5D4594, 831228) != 65535 || memmap.Uint32(0x5D4594, 832472) != 1 {
			t.Fatal("score heading/presentation")
		}
		if o.c.GUI.Captured() == nil {
			t.Fatal("score report did not open actual modal owner")
		}
		rows = append(rows, o.snapshotBriefing(t, 4, ret))
	}
	briefingCapture(t, "win-reports", rows, "0bd03ef1c9b7e8870a3270a18eadfce560273bb9c0f9b6b14b26123c63267f77")
}
func TestBriefingSelectionReports(t *testing.T) {
	o := newBriefingOwner(t)
	packet, free := alloc.New([69]byte{})
	defer free()
	var rows []briefingResult
	for _, op := range []int{5, 6} {
		for _, show := range []uintptr{0, 1} {
			for _, flag := range []byte{0, 2, 3} {
				for _, named := range []bool{false, true} {
					o.resetBriefing(t)
					o.constructBriefing(t)
					clear(packet[:])
					packet[0], packet[1] = 240, byte(op+8)
					packet[4] = flag
					binary.LittleEndian.PutUint16(packet[2:], 65535)
					imageName := "MissingBriefingImage"
					if named {
						imageName = "CustomBriefingImage"
						copy(packet[37:], "Briefing:Custom")
					}
					copy(packet[5:37], imageName)
					ret := legacy.PortTestBriefing(op, txptr(unsafe.Pointer(packet)), show, 0)
					if memmap.Uint32(0x5D4594, 832468) != 65535 || memmap.Uint32(0x5D4594, 832460) == 0 {
						t.Fatal("briefing stage/image selection")
					}
					text := alloc.GoString16((*uint16)(*memmap.PtrPtr(0x5D4594, 832464)))
					want := ""
					if named {
						want = "A custom briefing."
					}
					if text != want {
						t.Fatal("briefing caption selection")
					}
					if show != 0 {
						mode := uint32(2)
						if op == 6 {
							mode = 4
						}
						if memmap.Uint32(0x5D4594, 832472) != mode || o.c.GUI.Captured() == nil {
							t.Fatal("briefing did not enter expected presentation mode")
						}
					}
					rows = append(rows, o.snapshotBriefing(t, op, ret))
				}
			}
		}
	}
	briefingCapture(t, "selection-reports", rows, "ef25a95ac5036e2cd617b51cb034d4ea9a24571103644c33112b07d50da1daec")
}
func TestBriefingWinReportMissingAndDuplicateIDs(t *testing.T) {
	o := newBriefingOwner(t)
	packet, free := alloc.New([90]byte{})
	defer free()
	var rows []briefingResult
	for _, ids := range [][6]uint16{{999, 100}, {100, 100, 101}, {0, 100, 0, 101}, {105, 104, 103, 102, 101, 100}} {
		o.resetBriefing(t)
		o.constructBriefing(t)
		clear(packet[:])
		for i, id := range ids {
			r := packet[6+i*14:]
			binary.LittleEndian.PutUint16(r, id)
			binary.LittleEndian.PutUint32(r[10:], uint32(i%2))
		}
		ret := legacy.PortTestBriefing(4, txptr(unsafe.Pointer(packet)), 0, 0)
		rows = append(rows, o.snapshotBriefing(t, 4, ret))
	}
	briefingCapture(t, "missing-duplicate-ids", rows, "620310c8ca7616fc5bb15c033c05b6345ff85af22013ef7bc582332bb120866e")
}
func TestBriefingWindowConstructionLifecycle(t *testing.T) {
	o := newBriefingOwner(t)
	for cycle := 0; cycle < 100; cycle++ {
		o.resetBriefing(t)
		o.constructBriefing(t)
		w := o.briefWindow()
		if w.Size().X != 640 || w.Size().Y != 480 {
			t.Fatal(fmt.Sprintf("briefing dimensions %v", w.Size()))
		}
		o.releaseBriefing()
		if *o.briefWords["nox_wnd_briefing_831232"] != 0 || *o.briefWords["dword_5d4594_831236"] != 0 {
			t.Fatal("briefing owner not released")
		}
	}
}
