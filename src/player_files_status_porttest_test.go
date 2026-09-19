//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	flags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"testing"
	"unsafe"
)

func TestPlayerFilesStatusWrite(t *testing.T) {
	u, playerRaw := playerFileBookOwner(t)
	hp, free := alloc.New(server.HealthData{})
	t.Cleanup(free)
	u.HealthData = hp
	ud := u.UpdateDataPlayer()
	raw := unsafe.Slice((*byte)(unsafe.Pointer(u)), int(unsafe.Sizeof(*u)))
	data := unsafe.Slice((*byte)(unsafe.Pointer(ud)), int(unsafe.Sizeof(*ud)))
	saved := unsafe.Slice(memmap.PtrUint32(0x5D4594, 527696), 2)
	old := append([]uint32(nil), saved...)
	t.Cleanup(func() { copy(saved, old) })
	defer flags.PortTestGameFlags(0)()
	var rows []map[string]any
	for _, gf := range []flags.GameFlag{0, 2048, 4096, 8192, 2048 | 8192} {
		for _, v := range []uint16{0, 1, 32767, 32768, 65535} {
			flags.ResetGame()
			flags.SetGame(gf)
			hp.Cur = v ^ 0x1357
			hp.Max = v
			ud.ManaCur = v ^ 0x2468
			ud.ManaMax = v ^ 0x4321
			u.Poison540 = byte(v)
			u.Field541 = byte(v >> 8)
			u.Field542 = v ^ 0xabcd
			binary.LittleEndian.PutUint32(raw[28:], 0x3f800000+uint32(v))
			binary.LittleEndian.PutUint16(raw[124:], v^0x9876)
			saved[0] = 0xdeadbeef
			saved[1] = 0x12345678
			beforeP, beforeU, beforeD := bytes.Clone(playerRaw), bytes.Clone(raw), bytes.Clone(data)
			beforeHP := *hp
			want := []byte{2, 0, 0}
			wantSaved := []uint32{saved[0], saved[1]}
			if gf&2048 != 0 {
				want[2] = 1
				for _, word := range []uint16{hp.Max, ud.ManaMax, hp.Cur, ud.ManaCur} {
					want = binary.LittleEndian.AppendUint16(want, word)
				}
				want = append(want, u.Poison540, u.Field541)
				want = binary.LittleEndian.AppendUint16(want, u.Field542)
				want = append(want, raw[28:32]...)
				want = append(want, raw[124:126]...)
				wantSaved = []uint32{uint32(hp.Cur), uint32(ud.ManaCur)}
			}
			ret, got, pos := playerFileSection(t, "sub_41AA30", nil, uint32(uintptr(unsafe.Pointer(u))), 0)
			if ret != 1 || pos != int64(len(want)) || !bytes.Equal(got, want) || !bytes.Equal(playerRaw, beforeP) || !bytes.Equal(raw, beforeU) || !bytes.Equal(data, beforeD) || *hp != beforeHP || saved[0] != wantSaved[0] || saved[1] != wantSaved[1] {
				t.Fatal("status write", gf, v, ret, pos, len(want), saved, wantSaved)
			}
			rows = append(rows, map[string]any{"flags": uint32(gf), "value": v, "return": ret, "bytes": got, "position": pos, "saved": append([]uint32(nil), saved...)})
		}
	}
	spellbookCapture(t, "player-files-status-write", rows, "")
}
