//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/binfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/legacy/common/alloc/handles"
	"testing"
	"unsafe"
)

func TestClientSpriteAnimationVectors(t *testing.T) {
	t.Cleanup(handles.PortTestInit())
	c, _, _ := newEffectsFullOwner(t)
	raw := make([][]byte, 255)
	for i := range raw {
		raw[i] = spriteAnimationTestImage(i)
	}
	imgs, free := c.r.GetBag().PortTestSpriteImages(raw)
	t.Cleanup(free)
	refs := map[uint32]uint32{0: 0}
	for i, img := range imgs {
		refs[uint32(uintptr(img.C()))] = uint32(i + 1)
	}
	captureVector := func(p unsafe.Pointer) [][]uint32 {
		words := append([]uint32(nil), unsafe.Slice((*uint32)(p), 12)...)
		count := int(uint16(words[10]))
		groups := [][]uint32{words}
		for _, slot := range []int{1, 2, 3, 4, 6, 7, 8, 9} {
			if words[slot] == 0 {
				continue
			}
			a := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(uintptr(words[slot]))), count)...)
			for i, v := range a {
				n, ok := refs[v]
				if !ok {
					t.Fatal("vector parser returned unowned image")
				}
				a[i] = n
			}
			words[slot] = 0xe7000000 + uint32(slot)
			groups = append(groups, a)
		}
		return groups
	}
	type result struct {
		Op, Count, Delay, Kind, Return, Remaining int
		Metadata                                  []uint32
		Data                                      [][][]uint32
	}
	var out []result
	kinds := []string{"OneShot", "OneShotRemove", "Loop", "LoopAndFade", "Random", "Slave", "unknown"}
	for op := 4; op <= 7; op++ {
		if op == 5 {
			continue
		} // Unused C vector-array loader has an existing Go replacement.
		for _, count := range []int{0, 1, 2, 31, 255} {
			for _, delay := range []int{0, 1, 255} {
				for ki, kind := range kinds {
					input := []byte{}
					header := func() {
						input = append(input, byte(count), byte(delay), byte(len(kind)))
						input = append(input, kind...)
					}
					frames := func(n int) {
						for i := 0; i < n; i++ {
							input = binary.LittleEndian.AppendUint32(input, uint32((i+ki)%255))
						}
					}
					switch op {
					case 4:
						header()
					case 5:
						frames(8 * count)
					case 6:
						frames(count)
					case 7:
						for _, flag := range []uint32{8, 2, 4} {
							input = binary.LittleEndian.AppendUint32(input, 0x4652414d)
							input = binary.LittleEndian.AppendUint32(input, flag)
							input = append(input, 4, 'N', 'U', 'L', 'L', 1, 0)
							header()
							frames(count)
						}
						input = binary.LittleEndian.AppendUint32(input, 0x454e4420)
					}
					payload := append(input, 0xa6, 0x5b)
					bytes, _ := alloc.CloneSlice(payload)
					mf := binfile.NewMemFile(unsafe.Pointer(&bytes[0]), len(bytes))
					obj, freeObj := alloc.New(client.ObjectType{})
					obj.Field_54, obj.Field_60 = 0xabc, 0xdef
					vector, freeVector := alloc.New(client.AnimationVector{})
					vector.Cnt40 = uint16(count)
					vector.Val42 = uint16(delay)
					vector.Kind = client.AnimKind(ki % 6)
					scratch, freeScratch := alloc.Make([]byte{}, 512)
					got := legacy.PortTestSpriteParse(op, obj, mf, unsafe.Pointer(&scratch[0]), unsafe.Pointer(vector))
					if got != 1 || len(mf.Data()) != 2 {
						t.Fatalf("vector parser %d return %d remaining %d", op, got, len(mf.Data()))
					}
					groups := [][][]uint32{}
					if op == 7 {
						if obj.DrawData == nil || obj.Field_54 != 2 || *(*uint32)(obj.DrawData) != 148 || obj.DrawFunc != legacy.PortTestSpriteAnimationCallback(6) {
							t.Fatal("state parser did not install data/kind")
						}
						for state := 0; state < 3; state++ {
							groups = append(groups, captureVector(unsafe.Add(obj.DrawData, 4+state*48)))
						}
					} else {
						groups = append(groups, captureVector(unsafe.Pointer(vector)))
					}
					out = append(out, result{op, count, delay, ki, got, len(mf.Data()), []uint32{obj.Field_54, obj.Field_60}, groups})
					if op == 7 {
						legacy.PortTestSpriteFree(obj.DrawData, 8)
					} else if op != 4 {
						legacy.PortTestSpriteFreeVectorFrames(unsafe.Pointer(vector))
					}
					freeScratch()
					freeVector()
					freeObj()
					mf.Free()
				}
			}
		}
	}
	effectsCapture(t, "sprite-animation-vectors", out, len(out), "cad133b687a342d021a0316636a25ffc6082f6151bc11ce3a6008b90c831d661")
}
