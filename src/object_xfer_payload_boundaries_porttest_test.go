//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/spell"
	"github.com/opennox/libs/types"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func objectXferWriteRecord(t *testing.T, u *server.Object, path string) []byte {
	t.Helper()
	if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
		t.Fatal(err)
	}
	if err := u.CallXfer(nil); err != nil {
		t.Fatal(err)
	}
	if err := cryptfile.Close(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func objectXferReadRecord(t *testing.T, u *server.Object, path string) {
	t.Helper()
	if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
		t.Fatal(err)
	}
	defer cryptfile.Close()
	if err := u.CallXfer(nil); err != nil {
		t.Fatal(err)
	}
}
func objectXferCurrentRecord(version uint16) *mapDrawableStream {
	s := new(mapDrawableStream)
	s.u16(version)
	objectXferEmptyBase(s, 60)
	return s
}
func objectXferSetCommon(u *server.Object) {
	u.ObjFlags = 0
	u.Extent = 0x12345678
	u.ScriptIDVal = 88
	u.PosVec = types.Pointf{X: 64.25, Y: 128.75}
}
func objectXferSetWord(p unsafe.Pointer, off int, v uint32) {
	binary.LittleEndian.PutUint32(unsafe.Slice((*byte)(unsafe.Add(p, off)), 4), v)
}
func objectXferGetWord(p unsafe.Pointer, off int) uint32 {
	return binary.LittleEndian.Uint32(unsafe.Slice((*byte)(unsafe.Add(p, off)), 4))
}

func TestObjectXferTriggerDimensions(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "trigger.bin")
	var captures []objectXferCaptureRow
	defer func() { spellbookCapture(t, "object-xfer-trigger-dimensions", captures, "ccaa201395324e84466376f5406aa82be0c098c3b35c7c9b6026fed8e21463ba") }()
	for _, w := range []int32{-2147483648, -1, 0, 1, 59, 60, 61, 2147483647} {
		for _, h := range []int32{-2147483648, -1, 0, 1, 59, 60, 61, 2147483647} {
			t.Run(fmt.Sprintf("read-%d-%d", w, h), func(t *testing.T) {
				u := newObjectXferTyped(t, core, "Trigger")
				defer objectXferCaptureCase(t, &captures, u)
				s := objectXferCurrentRecord(61)
				objectXferHistoricalPayload(s, "Trigger", 61)
				binary.LittleEndian.PutUint32(s.Bytes()[21:], uint32(w))
				binary.LittleEndian.PutUint32(s.Bytes()[25:], uint32(h))
				if err := os.WriteFile(path, s.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				objectXferReadRecord(t, u, path)
				wantW, wantH := float32(w), float32(h)
				if wantW > 60 {
					wantW = 60
				}
				if wantH > 60 {
					wantH = 60
				}
				if u.Shape.Box.W != wantW || u.Shape.Box.H != wantH {
					t.Fatalf("shape=%v/%v want=%v/%v", u.Shape.Box.W, u.Shape.Box.H, wantW, wantH)
				}
			})
		}
	}
	for _, w := range []float32{-65535.75, -0.75, 0, 0.75, 59.75, 60, 61.75} {
		for _, h := range []float32{-65535.75, -0.75, 0, 0.75, 59.75, 60, 61.75} {
			t.Run(fmt.Sprintf("write-%g-%g", w, h), func(t *testing.T) {
				u := newObjectXferTyped(t, core, "Trigger")
				objectXferSetCommon(u)
				defer objectXferCaptureCase(t, &captures, u)
				u.Shape.Box.W = w
				u.Shape.Box.H = h
				data := objectXferWriteRecord(t, u, path)
				if len(data) != 85 || int32(binary.LittleEndian.Uint32(data[21:])) != int32(w) || int32(binary.LittleEndian.Uint32(data[25:])) != int32(h) {
					t.Fatal("trigger dimensions not serialized as truncating integers")
				}
				if u.Shape.Box.W != w || u.Shape.Box.H != h {
					t.Fatal("writer clamped live shape")
				}
			})
		}
	}
}

func TestObjectXferLinkedPayloads(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "linked.bin")
	a := core.WPs.Nox_xxx_waypointNewNotMap_579970(17, types.Pointf{})
	b := core.WPs.Nox_xxx_waypointNewNotMap_579970(29, types.Pointf{})
	t.Cleanup(func() { core.WPs.Pending = nil; alloc.Free(a); alloc.Free(b) })
	var captures []objectXferCaptureRow
	defer func() { spellbookCapture(t, "object-xfer-linked-payloads", captures, "6946cc7db7dc6e75cadd9a1b6258c4d403fd18a9db8b74526109d68f20fb40e1") }()
	for _, name := range []string{"Transporter", "Mover", "Sentry"} {
		for _, flags := range []uint32{0, 0x200000, 0x400000} {
			for _, linked := range []bool{false, true} {
				for _, value := range []uint32{0, 1, 255, 0x7fffffff, 0x80000000, 0xffffffff} {
					t.Run(fmt.Sprintf("%s-flags%x-linked%v-value%x", name, flags, linked, value), func(t *testing.T) {
						t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
						u := newObjectXferTyped(t, core, name)
						objectXferSetCommon(u)
						version := uint16(60)
						if name == "Sentry" {
							version = 61
						}
						want := objectXferCurrentRecord(version)
						switch name {
						case "Transporter":
							objectXferSetWord(u.UpdateData, 16, value)
							if linked {
								*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 12)) = u.CObj()
								want.u32(value)
							} else {
								want.u32(0)
							}
						case "Mover":
							objectXferSetWord(u.UpdateData, 4, value)
							objectXferSetWord(u.UpdateData, 8, value^0xabcdef01)
							objectXferSetWord(u.UpdateData, 32, value+2)
							*(*byte)(u.UpdateData) = byte(value)
							want.u32(value)
							want.u32(value ^ 0xabcdef01)
							want.u32(value + 2)
							want.u8(byte(value))
							if linked {
								*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 12)) = a.C()
								*(*unsafe.Pointer)(unsafe.Add(u.UpdateData, 20)) = b.C()
								want.u32(17)
								want.u32(29)
							} else {
								want.u32(0)
								want.u32(0)
							}
							u.SpeedBase = 3.5
							u.SpeedCur = -2.25
							want.f32(3.5)
							want.f32(-2.25)
						case "Sentry":
							objectXferSetWord(u.UpdateData, 0, value^0xabcdef01)
							objectXferSetWord(u.UpdateData, 4, value)
							objectXferSetWord(u.UpdateData, 8, value+2)
							want.u32(value)
							want.u32(value + 2)
							if flags == 0x200000 {
								want.u32(value)
							} else {
								want.u32(value ^ 0xabcdef01)
							}
						}
						got := objectXferWriteRecord(t, u, path)
						if !bytes.Equal(got, want.Bytes()) {
							t.Fatalf("payload=%x want=%x", got, want.Bytes())
						}
						v := newObjectXferTyped(t, core, name)
						defer objectXferCaptureCase(t, &captures, v)
						objectXferReadRecord(t, v, path)
						switch name {
						case "Transporter":
							expected := uint32(0)
							if linked {
								expected = value
							}
							if objectXferGetWord(v.UpdateData, 16) != expected {
								t.Fatal("transporter reference")
							}
						case "Mover":
							p, q := uint32(0), uint32(0)
							if linked {
								p, q = 17, 29
							}
							if objectXferGetWord(v.UpdateData, 16) != p || objectXferGetWord(v.UpdateData, 24) != q || v.SpeedBase != 3.5 || v.SpeedCur != -2.25 {
								t.Fatal("mover references/speeds")
							}
						case "Sentry":
							expected := value ^ 0xabcdef01
							if flags == 0x200000 {
								expected = value
							}
							if objectXferGetWord(v.UpdateData, 0) != expected {
								t.Fatal("sentry direction")
							}
						}
					})
				}
			}
		}
	}
}

func TestObjectXferGlyphSpellNames(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "glyph.bin")
	var captures []objectXferCaptureRow
	defer func() { spellbookCapture(t, "object-xfer-glyph-spells", captures, "0964c4ccf793cb6cca791706d267b4aac0f560b99d50d96250d88bdfa5b76629") }()
	for count := 0; count <= 5; count++ {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			u := newObjectXferTyped(t, core, "Glyph")
			objectXferSetCommon(u)
			u.Direction1 = 0xab35
			*(*byte)(unsafe.Add(u.InitData, 20)) = byte(count)
			objectXferSetWord(u.InitData, 28, 0x12345678)
			objectXferSetWord(u.InitData, 32, 0x87654321)
			want := objectXferCurrentRecord(60)
			want.u8(0x35)
			want.u32(0x12345678)
			want.u32(0x87654321)
			want.u8(byte(count))
			for i := 0; i < count; i++ {
				id := spell.ID(i + 1)
				objectXferSetWord(u.InitData, 4*i, uint32(id))
				name := id.String()
				want.u8(byte(len(name)))
				want.WriteString(name)
			}
			got := objectXferWriteRecord(t, u, path)
			if !bytes.Equal(got, want.Bytes()) {
				t.Fatalf("glyph=%x want=%x", got, want.Bytes())
			}
			v := newObjectXferTyped(t, core, "Glyph")
			v.Direction1 = 0xcd00
			defer objectXferCaptureCase(t, &captures, v)
			objectXferReadRecord(t, v, path)
			if v.Direction1 != 0xcd35 || v.Direction2 != 0xcd35 {
				t.Fatal("glyph direction byte and copy")
			}
			for i := 0; i < count; i++ {
				if objectXferGetWord(v.InitData, 4*i) != uint32(i+1) {
					t.Fatal("glyph spell-name mapping")
				}
			}
		})
	}
}
