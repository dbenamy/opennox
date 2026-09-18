//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/client/noxrender"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapMetadataAmbient(t *testing.T) {
	o := newObjectRenderOwner(t)
	owned := serverConfigOwnBytes(t, 0x587000, 142292, 20)
	oldFlags := noxflags.GetGame()
	t.Cleanup(func() { noxflags.UnsetGame(^noxflags.GameFlag(0)); noxflags.SetGame(oldFlags) })
	dir := t.TempDir()
	type record struct {
		Name  string
		IO    legacy.PortTestMapSectionWire
		State []byte
		Light noxrender.RGB
	}
	var records []record
	for _, read := range []bool{true, false} {
		for _, ver := range []uint16{0, 1, 2, 32767, 32768, 65535} {
			if !read && ver != 1 {
				continue
			}
			for _, flags := range []uint32{0, 1, 2, 0x200000, 0x200002, 0x400000} {
				for vi, values := range [][3]uint32{{0, 0, 0}, {31, 127, 255}, {256, 65535, 0x80000000}, {0xffffffff, 0x7fffffff, 0x12345678}} {
					name := fmt.Sprintf("read=%t/v=%d/flags=%x/values=%d", read, ver, flags, vi)
					t.Run(name, func(t *testing.T) {
						noxflags.UnsetGame(^noxflags.GameFlag(0))
						noxflags.SetGame(noxflags.GameFlag(flags))
						for i := range owned {
							owned[i] = 0x91
						}
						initial := noxrender.RGB{R: 301, G: -702, B: 913}
						o.c.r.Data().SetLightColor(initial)
						if !read {
							for i, v := range values {
								binary.LittleEndian.PutUint32(owned[4+4*i:], v)
							}
						}
						want := bytes.Clone(owned)
						light := initial
						wire := binary.LittleEndian.AppendUint16(nil, ver)
						for _, v := range values {
							wire = binary.LittleEndian.AppendUint32(wire, v)
						}
						ret := uint32(1)
						pos := int64(14)
						if int16(ver) < 1 {
							ret = 0
							pos = 2
						} else if read {
							copy(want[4:16], wire[2:])
							if flags&0x200002 != 0 {
								light = noxrender.RGB{R: int(int32(values[0])), G: int(int32(values[1])), B: int(int32(values[2]))}
							}
						}
						result := legacy.PortTestMapMetadata(legacy.PortTestMapSectionIO{Function: "ambient", Read: read, Data: wire}, dir)
						gotLight := o.c.r.Data().GetLightColor()
						if result.Return != ret || result.Position != pos || !bytes.Equal(result.Data, wire) {
							t.Fatalf("IO got %+v want return %d position %d bytes %x", result, ret, pos, wire)
						}
						if !bytes.Equal(owned, want) || gotLight != light {
							t.Fatalf("state %x light %+v want %x %+v", owned, gotLight, want, light)
						}
						records = append(records, record{name, result, bytes.Clone(owned), gotLight})
					})
				}
			}
		}
	}
	spellbookCapture(t, "map-metadata-ambient", records, "9c5ae3da0286943222f6f9a7abfeaed1528f04945c9c9ef10556654dda43e4f3")
}
