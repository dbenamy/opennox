//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/spell"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"io"
	"strings"
	"testing"
	"unsafe"
)

func TestResourceDefinitionsImages(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	c, _, _ := newEffectsFullOwner(t)
	imgs, restore := c.r.GetBag().PortTestSpriteImages([][]byte{spriteAnimationTestImage(0), spriteAnimationTestImage(1)})
	t.Cleanup(restore)
	typ, free := alloc.New(client.ObjectType{})
	defer free()
	raw := unsafe.Slice((*byte)(typ.C()), 128)
	arg, freeArg := alloc.CString("ignored")
	defer freeArg()
	var inputs [][]byte
	for _, id := range []uint32{0, 1} {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint32(b, id)
		copy(b[4:], []byte{1, 2, 3, 4})
		inputs = append(inputs, b)
	}
	for _, name := range []string{"", "icon", strings.Repeat("n", 126)} {
		for _, kind := range []byte{0, 1, 127, 128, 255} {
			b := []byte{255, 255, 255, 255, kind, byte(len(name) + 1)}
			b = append(b, []byte(name)...)
			b = append(b, 0, 0x7a, 0x7b)
			inputs = append(inputs, b)
		}
	}
	var rows []map[string]any
	for _, input := range inputs {
		for i := range raw {
			raw[i] = 0xa5
		}
		data, _ := alloc.CloneSlice(input)
		f := binfile.NewMemFile(unsafe.Pointer(&data[0]), len(data))
		ok := legacy.PortTestResourceClientParser("image", typ, f, unsafe.Pointer(arg))
		pos, _ := f.Seek(0, io.SeekCurrent)
		id := binary.LittleEndian.Uint32(input)
		ref := uint32(0)
		expected := uint32(0)
		consumed := 4
		if id < 2 {
			expected = uint32(uintptr(imgs[id].C()))
			ref = id + 1
		}
		if id == 0xffffffff {
			consumed = 6 + int(input[5])
		}
		want := bytes.Repeat([]byte{0xa5}, 128)
		resourcePut32(want, 112, expected)
		if !ok || pos != int64(consumed) || !bytes.Equal(raw, want) {
			t.Fatal("image", input, ok, pos, typ.PrettyImage)
		}
		snap := append([]byte(nil), raw...)
		resourcePut32(snap, 112, ref)
		rows = append(rows, map[string]any{"input": input, "ok": ok, "position": pos, "data": snap})
		f.Free()
	}
	spellbookCapture(t, "resource-definitions-images", rows, "9f5f274a1852436da768837d8f63c2b48733bac98aabd3df5be698abd3daefa6")
}
func TestResourceDefinitionsWandPartial(t *testing.T) {
	o := newReliableReportsOwner(t)
	o.s.SetTickRate(30)
	var rows []resourceParserRow
	for _, kind := range []string{"wand", "wandcast"} {
		for _, fill := range []byte{0, 0xa5} {
			inputs := []string{"", "5", "5 Missile", "5 Missile 0", "5 Missile 2 NONE", "5 Missile 2 MULTI_SHOT", "5 Missile 2 NONE missing"}
			if kind == "wandcast" {
				inputs = []string{"", "5", "5 2", "5 0 SPELL_MAGIC_MISSILE", "5 2 unknown"}
			}
			for _, input := range inputs {
				row := resourceParse(t, kind, input, fill)
				want := bytes.Repeat([]byte{fill}, 256)
				mode := uint32(0)
				if kind == "wandcast" {
					mode = 1
				}
				resourcePut32(want, 0, mode)
				fields := strings.Fields(input)
				ret := 0
				if len(fields) > 0 {
					want[108] = 5
					want[109] = 5
					resourcePut32(want, 112, 100)
				}
				if kind == "wand" && len(fields) > 1 {
					copy(want[4:], "Missile\x00")
					resourcePut32(want, 84, 0)
				}
				rateAt := 1
				if kind == "wand" {
					rateAt = 2
				}
				if len(fields) > rateAt {
					rate := uint32(15)
					if fields[rateAt] == "0" {
						rate = 0
					}
					resourcePut32(want, 100, rate)
				}
				if kind == "wandcast" && len(fields) > 2 {
					ret = 1
					id := uint32(0)
					if fields[2] == "SPELL_MAGIC_MISSILE" {
						id = uint32(spell.ParseID(fields[2]))
					}
					resourcePut32(want, 92, id)
				}
				if kind == "wand" && len(fields) > 2 {
					ret = 1
					if len(fields) > 3 && fields[3] == "MULTI_SHOT" {
						want[96] |= 1
					}
					if len(fields) > 4 {
						resourcePut32(want, 88, 0)
					}
				}
				if row.Return != ret || !bytes.Equal(row.Data, want) {
					t.Fatal("partial wand", kind, input, row.Data[:116], want[:116])
				}
				rows = append(rows, row)
			}
		}
	}
	spellbookCapture(t, "resource-definitions-wand-partial", rows, "edd04b87e361c4577d9df24ae64bd46f802802ed751f44e9ab3f5720d09146f4")
}
