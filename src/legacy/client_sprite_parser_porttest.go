//go:build porttest

package legacy

import (
	"github.com/opennox/opennox/v1/client"
	"github.com/opennox/opennox/v1/internal/binfile"
	"unsafe"
)

func PortTestSpriteParse(op int, obj *client.ObjectType, mf *binfile.MemFile, attr unsafe.Pointer, vector unsafe.Pointer) int {
	switch op {
	case 0:
		scratch := unsafe.Slice((*byte)(attr), 256)
		if spriteParseAnimate(obj, mf, scratch) {
			return 1
		}
		return 0
	case 1:
		scratch := unsafe.Slice((*byte)(attr), 256)
		if spriteParseConditional(obj, mf, scratch) {
			return 1
		}
		return 0
	case 2:
		scratch := unsafe.Slice((*byte)(attr), 256)
		if spriteParseStatic(obj, mf, scratch) {
			return 1
		}
		return 0
	case 3:
		scratch := unsafe.Slice((*byte)(attr), 256)
		if spriteParseRandom(obj, mf, scratch, false) {
			return 1
		}
		return 0
	case 4:
		return spriteVectorHeader((*client.AnimationVector)(vector), mf)
	case 6:
		return spriteVectorFrames((*client.AnimationVector)(vector), mf)
	case 7:
		if spriteParseState(obj, mf) {
			return 1
		}
		return 0
	case 9:
		scratch := unsafe.Slice((*byte)(attr), 256)
		if spriteParseRandom(obj, mf, scratch, true) {
			return 1
		}
		return 0
	case 8:
		return int(client.ParseAnimKind(GoStringP(attr)))
	}
	panic("unknown sprite parser")
}
func PortTestSpriteFree(data unsafe.Pointer, kind int) {
	spriteDataFreeKind(data, kind)
}

func PortTestSpriteFreeVectorFrames(data unsafe.Pointer) {
	spriteVectorFree(unsafe.Add(data, 4))
}
