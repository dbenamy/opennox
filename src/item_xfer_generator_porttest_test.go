//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	crypt "github.com/opennox/noxcrypt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/server"
)

func itemXferGeneratorPrefix(p *mapDrawableStream) {
	p.u8(3)
	p.Write([]byte{1, 2, 3})
	p.u8(7)
	p.u8(8)
	p.u32(1234)
	for i := 0; i < 4; i++ {
		objectXferScript(p, 0)
	}
}
func itemXferGeneratorTail(p *mapDrawableStream) { p.u8(3); p.Write([]byte{4, 5, 6}); p.u32(5678) }
func itemXferGoldRecord(u *server.Object) *mapDrawableStream {
	p := new(mapDrawableStream)
	p.u16(60)
	p.u16(64)
	p.u32(u.Extent)
	p.u32(uint32(u.ScriptIDVal))
	p.f32(u.PosVec.X)
	p.f32(u.PosVec.Y)
	p.u8(0)
	p.u32(objectXferGetWord(u.InitData, 0))
	return p
}
func itemXferTrackGeneratorChildren(t *testing.T, s *server.Server, u *server.Object) {
	children := (*[12]*server.Object)(u.UpdateData)
	for _, ch := range children {
		if ch != nil {
			trackObjectXferTyped(t, s, ch)
		}
	}
	t.Cleanup(func() { clear(children[:]) })
}
func TestItemXferGeneratorRecords(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "generator.bin")
	var rows []struct {
		State   itemXferCaptureRow
		WireSHA string
	}
	defer func() {
		spellbookCapture(t, "item-xfer-generator-records", rows, "4d0309d232cb1ce62afa7c862e480e1fe7afef1280e87e38da73bb76d5d395c0")
	}()
	for _, mask := range []uint16{0, 1, 2, 8, 0x10, 0x80, 0x100, 0x800, 0x111, 0x842, 0xfff} {
		t.Run(fmt.Sprintf("mask%03x", mask), func(t *testing.T) {
			u := newItemXferObject(t, s, "MonsterGenerator")
			objectXferSetCommon(u)
			u.Field34 = 777
			children := (*[12]*server.Object)(u.UpdateData)
			for i := 0; i < 12; i++ {
				if mask&(1<<i) != 0 {
					ch := newItemXferObject(t, s, "Gold")
					objectXferSetCommon(ch)
					ch.Extent = uint32(100 + i)
					ch.ScriptIDVal = 200 + i
					objectXferSetWord(ch.InitData, 0, uint32(300+i))
					children[i] = ch
				}
			}
			t.Cleanup(func() { clear(children[:]) })
			copy(unsafe.Slice((*byte)(unsafe.Add(u.UpdateData, 80)), 8), []byte{1, 2, 3, 4, 5, 6, 7, 8})
			objectXferSetWord(u.UpdateData, 88, 1234)
			objectXferSetWord(u.UpdateData, 92, 5678)
			want := objectXferCurrentRecord(63)
			var padding []int
			align := func() {
				for want.Len()%crypt.Block != 0 {
					padding = append(padding, want.Len())
					want.u8(0)
				}
			}
			itemXferGeneratorPrefix(want)
			want.u8(3)
			var ordered [3][]*server.Object
			for row := 0; row < 3; row++ {
				for col := 0; col < 4; col++ {
					if ch := children[4*row+col]; ch != nil {
						ordered[row] = append(ordered[row], ch)
					}
				}
				want.u8(byte(len(ordered[row])))
				for _, ch := range ordered[row] {
					name := "portgold"
					want.u8(byte(len(name)))
					want.WriteString(name)
					want.u16(uint16(1000 + ch.TypeInd))
					rec := itemXferGoldRecord(ch)
					align()
					want.u32(uint32(rec.Len()))
					want.u32(0)
					want.Write(rec.Bytes())
				}
			}
			itemXferGeneratorTail(want)
			logicalEnd := want.Len()
			align()
			// Real saves use cipher blocks: each child length occupies an aligned
			// eight-byte block. Compatibility padding is excluded from the field
			// contract and retained in the separate frozen wire hash.
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, crypt.MapKey); err != nil {
				t.Fatal(err)
			}
			if err := u.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			written := cryptfile.Global().PortTestChecksum()
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			encoded, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if len(encoded) != want.Len() {
				t.Fatalf("wire length=%d want=%d", len(encoded), want.Len())
			}
			cipher, err := crypt.NewCipher(crypt.MapKey)
			if err != nil {
				t.Fatal(err)
			}
			got := append([]byte(nil), encoded...)
			for off := 0; off < len(got); off += crypt.Block {
				cipher.Decrypt(got[off:off+crypt.Block], got[off:off+crypt.Block])
			}
			for _, off := range padding {
				got[off] = 0
			}
			if !bytes.Equal(got, want.Bytes()) {
				for off := range got {
					if got[off] != want.Bytes()[off] {
						t.Fatalf("plaintext offset %d=%02x want=%02x", off, got[off], want.Bytes()[off])
					}
				}
			}
			v := newItemXferObject(t, s, "MonsterGenerator")
			v.Field34 = 999
			before := s.Objs.Alive
			if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, crypt.MapKey); err != nil {
				t.Fatal(err)
			}
			defer cryptfile.Close()
			err = v.CallXfer(nil)
			itemXferTrackGeneratorChildren(t, s, v)
			if err != nil {
				t.Fatal(err)
			}
			loaded := (*[12]*server.Object)(v.UpdateData)
			count := 0
			for row := 0; row < 3; row++ {
				for col := 0; col < 4; col++ {
					ch := loaded[4*row+col]
					if col >= len(ordered[row]) {
						if ch != nil {
							t.Fatal("unexpected child slot")
						}
						continue
					}
					original := ordered[row][col]
					count++
					if ch == nil || ch == original || ch.TypeInd != original.TypeInd || ch.Extent != original.Extent || ch.ScriptIDVal != original.ScriptIDVal || objectXferGetWord(ch.InitData, 0) != objectXferGetWord(original.InitData, 0) {
						t.Fatal("nested child identity/state/order")
					}
				}
			}
			if s.Objs.Alive != before+count || v.Field34 != 999 || u.Field34 != 777 {
				t.Fatal("generator ownership/lifetime")
			}
			pos, err := cryptfile.Global().File.Seek(0, 1)
			if err != nil || pos != int64(logicalEnd) {
				t.Fatalf("position=%d want=%d err=%v", pos, logicalEnd, err)
			}
			var captured []itemXferCaptureRow
			itemXferCaptureCase(t, &captured, v, cryptfile.Global().PortTestChecksum(), written)
			rows = append(rows, struct {
				State   itemXferCaptureRow
				WireSHA string
			}{captured[0], fmt.Sprintf("%x", sha256.Sum256(encoded))})
		})
	}
}
func TestItemXferGeneratorPartialFailure(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "partial-generator.bin")
	s.Types.ByID("PortInvisibleLight").Weight = 255
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-generator-partial", rows, "d2dc20501579610e180c7d21ac00a89dab1b1cfeae99a50d7422926583a83716")
	}()
	for _, first := range []bool{false, true} {
		for _, bad := range []string{"", "MissingGeneratorType", "PortInvisibleLight"} {
			t.Run(fmt.Sprintf("first%v-bad%s", first, bad), func(t *testing.T) {
				u := newItemXferObject(t, s, "MonsterGenerator")
				u.Field34 = 999
				before := s.Objs.Alive
				p := objectXferCurrentRecord(63)
				itemXferGeneratorPrefix(p)
				p.u8(1)
				n := byte(1)
				if first {
					n = 2
				}
				p.u8(n)
				if first {
					name := "PortGold"
					p.u8(byte(len(name)))
					p.WriteString(name)
					p.u16(0)
					p.u32(0xffffffff)
					rec := objectXferCurrentRecord(60)
					rec.u32(123)
					p.Write(rec.Bytes())
				}
				p.u8(byte(len(bad)))
				p.WriteString(bad)
				if bad == "PortInvisibleLight" {
					p.u16(0)
					p.u32(2)
					p.u16(61)
				}
				wantPos := p.Len()
				p.u32(0xdeadbeef)
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				err := u.CallXfer(nil)
				itemXferTrackGeneratorChildren(t, s, u)
				if err == nil {
					t.Fatal("bad child accepted")
				}
				expected := before
				if first {
					expected++
				}
				if s.Objs.Alive != expected || u.Field34 != 0 {
					t.Fatal("partial failure ownership/lifetime")
				}
				children := (*[12]*server.Object)(u.UpdateData)
				if first {
					if children[0] == nil || objectXferGetWord(children[0].InitData, 0) != 123 {
						t.Fatal("earlier child lost")
					}
				} else if children[0] != nil {
					t.Fatal("rejected child linked")
				}
				for _, ch := range children[1:] {
					if ch != nil {
						t.Fatal("unexpected linked child")
					}
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != int64(wantPos) {
					t.Fatalf("position=%d want=%d err=%v", pos, wantPos, err)
				}
				itemXferCaptureCase(t, &rows, u)
			})
		}
	}
}
