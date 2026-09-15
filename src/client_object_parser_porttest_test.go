//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

func TestClientObjectDrawingDoorParser(t *testing.T) {
	o := newObjectDrawingOwner(t)
	type result struct {
		Count, Variant, Remaining int
		Metadata, Frames          []uint32
		Scratch                   []byte
	}
	var out []result
	for _, count := range []int{0, 1, 2, 31, 255} {
		for variant := 0; variant < 3; variant++ {
			payload := []byte{byte(count)}
			for i := 0; i < count; i++ {
				payload = binary.LittleEndian.AppendUint32(payload, uint32((i+variant*7)%32))
			}
			payload = append(payload, 0xab, 0xcd)
			raw, freeRaw := alloc.CloneSlice(payload)
			mf := binfile.NewMemFile(unsafe.Pointer(&raw[0]), len(raw))
			obj, freeObj := alloc.New(client.ObjectType{})
			obj.Field_54 = 0xabc
			obj.Field_60 = 0xdef
			scratch, freeScratch := alloc.Make([]byte{}, 512)
			for i := range scratch {
				scratch[i] = byte(variant * 57)
			}
			ok := legacy.PortTestObjectDoorParse(obj, mf, unsafe.Pointer(&scratch[0]))
			if !ok || obj.DrawData == nil || obj.DrawFunc != legacy.PortTestObjectDrawCallback(0) || len(mf.Data()) != 2 {
				t.Fatal("door parser contract")
			}
			data := unsafe.Slice((*uint32)(obj.DrawData), 3)
			frames := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(uintptr(data[1]))), count)...)
			for i, p := range frames {
				ref, ok := o.c.imageRefs[p]
				if !ok {
					t.Fatal("unowned parsed door frame")
				}
				frames[i] = ref
			}
			out = append(out, result{count, variant, len(mf.Data()), []uint32{obj.Field_54, obj.Field_60, data[0], data[2]}, frames, append([]byte(nil), scratch...)})
			legacy.PortTestSpriteFree(obj.DrawData, 2)
			freeScratch()
			freeObj()
			freeRaw()
		}
	}
	effectsCapture(t, "object-drawing-door-parser", out, len(out), "363fdd6579402ea968fd231b79e37180f6b0093a0522cd8a57c50914dc1a3bfe")
}
