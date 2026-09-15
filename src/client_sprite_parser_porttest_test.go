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

func TestClientSpriteAnimationParsers(t *testing.T) {
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
	imageRef := func(p uint32) uint32 {
		n, ok := refs[p]
		if !ok {
			t.Fatal("parser returned unowned image handle")
		}
		return n
	}
	type result struct {
		Op, Count, Delay, Kind, States, Return, Remaining int
		Metadata                                          []uint32
		Data                                              [][]uint32
		Scratch                                           []byte
	}
	var out []result
	kinds := []string{"OneShot", "OneShotRemove", "Loop", "LoopAndFade", "Random", "Slave", "loop", "", "unknown"}
	for _, op := range []int{0, 1, 2, 3, 9} {
		for _, count := range []int{0, 1, 2, 31, 255} {
			for _, delay := range []int{0, 1, 255} {
				for ki, kind := range kinds {
					states := 1
					if op == 1 {
						states = ki % 6
					}
					input := []byte{}
					header := func() {
						input = append(input, byte(count), byte(delay), byte(len(kind)))
						input = append(input, kind...)
					}
					frames := func(n int) {
						for i := 0; i < n; i++ {
							if ki == 8 && i%7 == 0 {
								input = binary.LittleEndian.AppendUint32(input, 0xffffffff)
								input = append(input, 6, 7, 'm', 'i', 's', 's', 'i', 'n', 'g')
							} else {
								input = binary.LittleEndian.AppendUint32(input, uint32((i+ki)%255))
							}
						}
					}
					switch op {
					case 0:
						header()
						frames(count)
					case 1:
						input = append(input, byte(states))
						for i := 0; i < states; i++ {
							header()
							frames(count)
						}
					case 2:
						frames(1)
					case 3, 9:
						input = append(input, byte(count))
						frames(count)
					}
					payload := append(append([]byte(nil), input...), 0xa6, 0x5b)
					bytes, _ := alloc.CloneSlice(payload)
					mf := binfile.NewMemFile(unsafe.Pointer(&bytes[0]), len(bytes))
					obj, freeObj := alloc.New(client.ObjectType{})
					obj.Field_54, obj.Field_60 = 0xabc, 0xdef
					scratch, freeScratch := alloc.Make([]byte{}, 512)
					for i := range scratch {
						scratch[i] = 0x5a
					}
					got := legacy.PortTestSpriteParse(op, obj, mf, unsafe.Pointer(&scratch[0]), nil)
					if got != 1 || obj.DrawData == nil {
						t.Fatalf("parser %d failed valid input", op)
					}
					if len(mf.Data()) != 2 {
						t.Fatalf("parser %d consumed %d bytes, want %d", op, len(payload)-len(mf.Data()), len(input))
					}
					callbackOp := op
					if op == 9 {
						callbackOp = 5
					}
					if obj.DrawFunc != legacy.PortTestSpriteAnimationCallback(callbackOp) {
						t.Fatalf("parser %d callback mismatch", op)
					}
					size := map[int]int{0: 4, 1: 14, 2: 2, 3: 3, 9: 3}[op]
					words := append([]uint32(nil), unsafe.Slice((*uint32)(obj.DrawData), size)...)
					groups := [][]uint32{words}
					if op == 2 {
						words[1] = imageRef(words[1])
					} else {
						slots := 1
						if op == 1 {
							slots = 5
						}
						for slot := 1; slot <= slots; slot++ {
							p := words[slot]
							if p == 0 {
								continue
							}
							n := count
							if op == 1 && slot > states {
								t.Fatal("unparsed conditional state allocated frames")
							}
							a := append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(uintptr(p))), n)...)
							for i := range a {
								a[i] = imageRef(a[i])
							}
							words[slot] = 0xe6000000 + uint32(slot)
							groups = append(groups, a)
						}
					}
					out = append(out, result{op, count, delay, ki, states, got, len(mf.Data()), []uint32{obj.Field_54, obj.Field_60}, groups, append([]byte(nil), scratch[:256]...)})
					legacy.PortTestSpriteFree(obj.DrawData, map[int]int{0: 3, 1: 4, 2: 1, 3: 2, 9: 2}[op])
					freeScratch()
					freeObj()
					mf.Free()
				}
			}
		}
	}
	effectsCapture(t, "sprite-animation-parsers", out, len(out), "fe24ac26e809b9b3544791850d802ac5011f588d3e78bbb8b674e82972eaf91c")
}
