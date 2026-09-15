//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

func (o *briefingOwner) drawBriefing(t *testing.T, op int, frame, state uint32) briefingResult {
	o.c.srv.SetFrame(frame)
	var flag *uint32
	switch op {
	case 1:
		flag = o.briefWords["dword_587000_122956"]
	case 2:
		flag = o.briefWords["nox_xxx_aSpellphoneme_3_587000_123008"]
	case 3:
		flag = memmap.PtrUint32(0x587000, 123012)
	}
	*flag = state
	o.displayText = nil
	o.loads = nil
	w := o.briefWindow().ChildByID(1010)
	ret := legacy.PortTestBriefing(op, txptr(w.C()), txptr(w.DrawData().C()), 0)
	want := state
	if frame%30 == 0 {
		if state == 1 {
			want = 0
		} else {
			want = 1
		}
	}
	if *flag != want {
		t.Fatalf("blink state op%d frame%d got%d want%d", op, frame, *flag, want)
	}
	prompt := 0
	for _, v := range o.displayText {
		if v.Text == "Press a key to continue" {
			prompt++
		}
	}
	if (prompt == 1) != (want == 1) || prompt > 1 {
		t.Fatalf("blink prompt op%d frame%d state%d count%d", op, frame, state, prompt)
	}
	if want != 1 {
		expected := frame / 30
		if frame%30 == 0 {
			expected = 1
		}
		if ret != expected {
			t.Fatal("hidden prompt return")
		}
	}
	return o.snapshotBriefing(t, op, ret)
}
func TestBriefingTitleFramesAndCaptions(t *testing.T) {
	o := newBriefingOwner(t)
	var rows []briefingResult
	for _, stage := range []uint32{0, 1, 65535, 0x80000000, 0xffffffff} {
		for _, caption := range []bool{false, true} {
			for _, frame := range []uint32{0, 29, 30, 31, 59, 60, 0xffffffff} {
				for _, state := range []uint32{0, 1, 2} {
					o.resetBriefing(t)
					o.constructBriefing(t)
					*memmap.PtrUint32(0x5D4594, 832468) = stage
					if caption {
						*memmap.PtrPtr(0x5D4594, 832464) = unsafe.Pointer(alloc.InternCString16("A custom briefing."))
					}
					rows = append(rows, o.drawBriefing(t, 2, frame, state))
				}
			}
		}
	}
	briefingCapture(t, "title-frames", rows, "95f69e75ca5ccc15a6d5cbcc62dc45426dc2d9047dcd1906d7b09c61a7db2691")
}
func TestBriefingStatisticsLayouts(t *testing.T) {
	o := newBriefingOwner(t)
	var rows []briefingResult
	for _, mask := range []uint{0, 1, 2, 5, 63} {
		for _, narrow := range []bool{false, true} {
			for _, frame := range []uint32{29, 30} {
				for _, state := range []uint32{0, 1} {
					o.resetBriefing(t)
					o.constructBriefing(t)
					legacy.Set_dword_8531A0_2576(&o.players[0])
					*o.briefWords["dword_5d4594_832476"] = 85
					*memmap.PtrUint32(0x5D4594, 831228) = 65535
					*memmap.PtrUint32(0x5D4594, 832356) = 27
					if narrow {
						*memmap.PtrUint32(0x587000, 122968) = 130
					}
					for i := 0; i < 6; i++ {
						if mask&(1<<i) == 0 {
							continue
						}
						p := &o.players[i]
						if narrow {
							p.SetName(strings.Repeat("😀", 10) + "WWWWΩ")
						}
						b := unsafe.Slice((*byte)(memmap.PtrOff(0x5D4594, 832364+uintptr(16*i))), 16)
						binary.LittleEndian.PutUint32(b, uint32(txptr(unsafe.Pointer(p))))
						for j := 0; j < 4; j++ {
							binary.LittleEndian.PutUint16(b[4+j*2:], uint16(i*11+j))
						}
						binary.LittleEndian.PutUint32(b[12:], []uint32{0, 1, 65535, 0x7fffffff, 0x80000000, 0xffffffff}[i])
					}
					rows = append(rows, o.drawBriefing(t, 1, frame, state))
				}
			}
		}
	}
	briefingCapture(t, "statistics-layouts", rows, "16380893b8863b712b14ff8c301da8ac16c7610fa8753bba3589bad71bed8fc6")
}
func TestBriefingInstructionsLayoutAndFrames(t *testing.T) {
	o := newBriefingOwner(t)
	var rows []briefingResult
	for _, long := range []bool{false, true} {
		value := "Instruction 3a"
		if long {
			value = strings.Repeat("Wide instruction ", 12)
		}
		o.strings["GeneralPrint:QuestSplash3a"] = strman.Variant{Str: value}
		o.installStrings(t)
		for _, width := range []int{640, 800} {
			for _, frame := range []uint32{0, 29, 30, 31, 60} {
				for _, state := range []uint32{0, 1, 2} {
					o.resetBriefing(t)
					o.resize(width, 480)
					o.constructBriefing(t)
					rows = append(rows, o.drawBriefing(t, 3, frame, state))
					if len(o.objects) != 12 {
						t.Fatal("instruction sprites not allocated")
					}
				}
			}
		}
	}
	briefingCapture(t, "instructions-layouts", rows, "4e62403c41d948dd6d6f0dae238b9592af299496b83b805dac7e6d4178cd9c35")
}
