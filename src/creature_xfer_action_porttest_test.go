//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/sound"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"os"
	"path/filepath"
	"testing"
	"unsafe"
)

var creatureXferReadTimers = []int{280, 404, 496}
var creatureXferAlwaysTimers = []int{536, 540, 1204, 508, 512, 516, 520, 528, 532, 548, 1208, 1212}

func creatureXferSetTimers(u *server.Object, value uint32) {
	for _, off := range creatureXferReadTimers {
		objectXferSetWord(u.UpdateData, off, value)
	}
	for _, off := range creatureXferAlwaysTimers {
		objectXferSetWord(u.UpdateData, off, value)
	}
	*(*byte)(unsafe.Add(u.UpdateData, 544)) = 255 // Empty signed action stack.
}
func creatureXferEmptyActionRecord(u *server.Object, version uint16, frame uint32) *mapDrawableStream {
	p := new(mapDrawableStream)
	p.u16(version)
	if version >= 2 {
		p.u8(1)
	}
	p.u32(frame)
	p.u32(0)
	at := func(off, n int) { p.Write(unsafe.Slice((*byte)(unsafe.Add(u.UpdateData, off)), n)) }
	at(8, 4)
	at(12, 0)
	at(268, 4)
	at(272, 8)
	at(280, 4)
	at(284, 1)
	at(296, 4)
	at(364, 4)
	at(368, 8)
	at(376, 4)
	at(380, 8)
	itemXferRewardName(p, sound.ID(0).String())
	at(396, 8)
	at(404, 4)
	at(481, 1)
	at(482, 1)
	at(483, 1)
	if version < 3 {
		p.u32(0)
	}
	at(496, 4)
	at(500, 8)
	// The legacy reserved word carries the completed waypoint-loop count.
	p.u32(objectXferGetWord(u.UpdateData, 296))
	at(536, 4)
	at(540, 4)
	at(544, 1)
	at(1129, 1)
	p.u32(0)
	at(1204, 4)
	at(2096, 4)
	at(2100, 4)
	at(2104, 1)
	at(2105, 1)
	p.u8(0)
	if version >= 4 {
		at(4, 4)
		at(288, 4)
		at(292, 4)
		p.u32(0)
		at(492, 4)
		for _, off := range []int{508, 512, 516, 520, 528, 532, 524, 548} {
			at(off, 4)
		}
		at(1128, 1)
		p.u32(0)
		at(1208, 4)
		at(1212, 4)
		p.u32(0)
		at(2172, 1)
	}
	return p
}
func creatureXferClampTimer(value, delta uint32) uint32 {
	out := value + delta
	if int32(out) < 1 {
		return 1
	}
	return out
}
func TestCreatureXferActionCurrent(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "action.bin")
	var rows []struct {
		Writer, Reader itemXferCaptureRow
		Wire           []byte
	}
	defer func() { spellbookCapture(t, "creature-xfer-action-current", rows, "") }()
	for _, frames := range [][2]uint32{{123, 123}, {123, 0}, {0, 123}, {0xffffffff, 0}, {0, 0xffffffff}} {
		for _, value := range []uint32{0, 1, 0x7fffffff, 0xffffffff} {
			t.Run(fmt.Sprintf("frame%08x-%08x-value%08x", frames[0], frames[1], value), func(t *testing.T) {
				u := newCreatureXferObject(t, s, "Monster")
				creatureXferSetTimers(u, value)
				want := creatureXferEmptyActionRecord(u, 4, frames[0])
				t.Cleanup(noxflags.PortTestGameFlags(1))
				s.SetFrame(frames[0])
				if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
					t.Fatal(err)
				}
				if ret := legacy.PortTestCreatureXferHelper(0, u, nil, 0); ret != 1 {
					t.Fatalf("writer result=%d", ret)
				}
				written := cryptfile.Global().PortTestChecksum()
				cryptfile.Close()
				got, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				if !bytes.Equal(got, want.Bytes()) {
					t.Fatalf("action wire length=%d want=%d", len(got), want.Len())
				}
				for _, off := range creatureXferReadTimers {
					if objectXferGetWord(u.UpdateData, off) != value {
						t.Fatal("writer changed read-only timer")
					}
				}
				for _, off := range creatureXferAlwaysTimers {
					if objectXferGetWord(u.UpdateData, off) != creatureXferClampTimer(value, 0) {
						t.Fatalf("writer timer %d", off)
					}
				}
				var wr []itemXferCaptureRow
				itemXferCaptureCase(t, &wr, u, 0, written)
				v := newCreatureXferObject(t, s, "Monster")
				s.SetFrame(frames[1])
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if ret := legacy.PortTestCreatureXferHelper(0, v, nil, 0); ret != 1 {
					t.Fatalf("reader result=%d", ret)
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != int64(len(got)) {
					t.Fatal("action position")
				}
				for _, off := range append(append([]int(nil), creatureXferReadTimers...), creatureXferAlwaysTimers...) {
					if objectXferGetWord(v.UpdateData, off) != creatureXferClampTimer(value, frames[1]-frames[0]) {
						t.Fatalf("reader timer %d", off)
					}
				}
				var rd []itemXferCaptureRow
				itemXferCaptureCase(t, &rd, v)
				rows = append(rows, struct {
					Writer, Reader itemXferCaptureRow
					Wire           []byte
				}{wr[0], rd[0], got})
			})
		}
	}
}
func TestCreatureXferActionGate(t *testing.T) {
	s := newCreatureXferOwner(t)
	path := filepath.Join(t.TempDir(), "gate.bin")
	for _, flags := range []uint32{0, 0x400000, 0x400001} {
		t.Run(fmt.Sprintf("flags%x", flags), func(t *testing.T) {
			u := newCreatureXferObject(t, s, "Monster")
			creatureXferSetTimers(u, 0)
			t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			ret := legacy.PortTestCreatureXferHelper(0, u, nil, 0)
			cryptfile.Close()
			data, err := os.ReadFile(path)
			if err != nil || ret != 1 || !bytes.Equal(data, []byte{4, 0, 0}) {
				t.Fatal("writer action gate")
			}
			for _, off := range creatureXferAlwaysTimers {
				if objectXferGetWord(u.UpdateData, off) != 0 {
					t.Fatal("gated writer changed timer")
				}
			}
		})
	}
	for _, version := range []uint16{2, 3, 4, 5, 32767} {
		t.Run(fmt.Sprintf("v%d", version), func(t *testing.T) {
			u := newCreatureXferObject(t, s, "Monster")
			creatureXferSetTimers(u, 77)
			p := new(mapDrawableStream)
			p.u16(version)
			p.u8(0)
			p.u32(0x12345678)
			if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
				t.Fatal(err)
			}
			defer cryptfile.Close()
			ret := legacy.PortTestCreatureXferHelper(0, u, nil, 0)
			want, stop := uint32(1), int64(3)
			if version > 4 {
				want, stop = 0, 2
			}
			pos, err := cryptfile.Global().File.Seek(0, 1)
			if err != nil || ret != want || pos != stop {
				t.Fatal("reader action gate")
			}
			for _, off := range creatureXferAlwaysTimers {
				if objectXferGetWord(u.UpdateData, off) != 77 {
					t.Fatal("gated reader changed timer")
				}
			}
		})
	}
}
