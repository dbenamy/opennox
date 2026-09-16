//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"io"
	"math"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/cryptfile"
)

func TestObjectXferLightReadContracts(t *testing.T) {
	core := newObjectXferOwner(t)
	mapDrawableTables(t)
	path := filepath.Join(t.TempDir(), "light.bin")
	for _, version := range []int16{-1, 0, 1, 2, 39, 40, 41, 42, 60} {
		for variant := 0; variant < 2; variant++ {
			t.Run(fmt.Sprintf("v%d-variant%d", version, variant), func(t *testing.T) {
				u := newObjectXferTyped(t, core, "InvisibleLight")
				u.ObjFlags = 0
				u.Field34 = 999
				var stream, payload mapDrawableStream
				stream.u16(uint16(version))
				objectXferEmptyBase(&stream, version)
				payload.light(version, variant)
				stream.Write(payload.Bytes())
				if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if err := u.CallXfer(nil); err != nil {
					t.Fatal(err)
				}
				pos, err := cryptfile.Global().File.Seek(0, io.SeekCurrent)
				if err != nil || pos != int64(stream.Len()) || u.Field34 != 999 {
					t.Fatalf("stream position=%d want=%d lifetime=%d err=%v", pos, stream.Len(), u.Field34, err)
				}
				want := make([]byte, 140)
				wire := payload.Bytes()
				copy(want[:36], wire[:36])
				if version < 2 {
					want[138] = 128
					if variant == 1 {
						binary.LittleEndian.PutUint32(want[4:], math.Float32bits(63))
						binary.LittleEndian.PutUint32(want[8:], uint32(client.LightRadius(63)))
					}
				} else {
					copy(want[40:132], wire[36:128])
					copy(want[134:139], wire[128:133])
					if version == 41 {
						binary.LittleEndian.PutUint32(want[36:], uint32(wire[133]))
					} else if version >= 42 {
						copy(want[36:40], wire[133:137])
					}
				}
				got := unsafe.Slice((*byte)(unsafe.Add(u.Field189, 2432)), 140)
				if !bytes.Equal(got, want) {
					t.Fatalf("light data=%x want=%x", got, want)
				}
			})
		}
	}
}

func TestObjectXferLightWriteContracts(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "light.bin")
	for _, matched := range []bool{false, true} {
		t.Run(fmt.Sprint(matched), func(t *testing.T) {
			u := newObjectXferTyped(t, core, "InvisibleLight")
			u.ObjFlags = 0
			u.Extent = 123
			u.ScriptIDVal = 88
			u.PosVec = types.Pointf{X: 64.25, Y: 128.75}
			u.Field34 = 777
			dr := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(64, 128))
			if dr == nil {
				t.Fatal("light drawable allocation")
			}
			defer c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
			dr.NetCode32 = 124
			if matched {
				dr.NetCode32 = 123
			}
			raw := unsafe.Slice((*byte)(unsafe.Add(dr.C(), 136)), 140)
			for i := range raw {
				raw[i] = byte(3*i + 7)
			}
			light := make([]byte, 140)
			if matched {
				copy(light, raw)
			}
			var want mapDrawableStream
			want.u16(60)
			want.u16(64)
			want.u32(123)
			want.u32(88)
			want.f32(64.25)
			want.f32(128.75)
			want.u8(0)
			want.Write(light[:36])
			want.Write(light[40:132])
			want.Write(light[134:139])
			want.Write(light[36:40])
			if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			if err := u.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
			if err := cryptfile.Close(); err != nil {
				t.Fatal(err)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, want.Bytes()) || u.Field34 != 777 {
				t.Fatalf("light bytes=%x want=%x lifetime=%d", got, want.Bytes(), u.Field34)
			}
		})
	}
}
