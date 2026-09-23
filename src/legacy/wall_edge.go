package legacy

/*
#include "defs.h"

*/
import "C"
import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/client/noxrender"
	"image"
	"unsafe"
)

// wallEdgeDraw shades opaque RLE runs into the existing shared pixel-row owner.
// Its former C return was scratch state ignored by its sole production caller;
// the former flags argument was unused. Keep both out of the private Go API.
func wallEdgeDraw(handle noxrender.ImageHandle, pos image.Point, first, second *[3]uint32, bottom, left, right, crop int) {
	if handle == nil {
		return
	}
	img := GetClient().R2().GetBag().AsImage(handle)
	if img.Type()&63 != 3 {
		return
	}
	data := img.Pixdata()
	if len(data) == 0 {
		return
	}
	width := int(int32(binary.LittleEndian.Uint32(data)))
	height := int(int32(binary.LittleEndian.Uint32(data[4:]))) - crop
	if height <= 0 {
		return
	}
	pos = pos.Add(image.Pt(int(int32(binary.LittleEndian.Uint32(data[8:]))), int(int32(binary.LittleEndian.Uint32(data[12:])))))
	data = data[17:]
	minX, minY := int(int32(dword_5d4594_3807140)), int(int32(dword_5d4594_3807136))
	maxX, maxY := int(int32(dword_5d4594_3807116)), int(int32(dword_5d4594_3807152))
	if pos.X > maxX || pos.Y > maxY {
		return
	}
	skipX, skipY, drawWidth := 0, 0, width
	if pos.X < minX {
		if pos.X+width < minX {
			return
		}
		skipX = minX - pos.X
		drawWidth -= skipX
		pos.X = minX
	}
	if pos.X+drawWidth > maxX {
		drawWidth = maxX - pos.X
	}
	if pos.Y < minY {
		if pos.Y+height < minY {
			return
		}
		skipY = minY - pos.Y
		height -= skipY
		pos.Y = minY
	}
	if pos.Y+height > maxY {
		height = maxY - pos.Y
	}
	if drawWidth <= 0 || height <= 0 {
		return
	}
	height = min(height, max(bottom, pos.Y)-pos.Y)
	extra := 0
	if left > pos.X+skipX {
		extra = left - skipX - pos.X
		skipX = left - pos.X
		drawWidth -= extra
	}
	if right < pos.X+drawWidth+skipX {
		drawWidth = right - skipX - pos.X
	}
	if drawWidth <= 0 {
		return
	}
	var start [3]uint32
	var shading [6]uint32
	color, step := shading[:3], shading[3:]
	for i := range start {
		start[i] = first[i] << 8
		step[i] = uint32(int32((second[i]<<8)-start[i]) / int32(width))
	}
	// Decode every row to retain its stream position even when no pixels are drawn.
	nextRun := func() (op byte, count int, pixels []uint16) {
		op, count = data[0], int(data[1])
		data = data[2:]
		if op == 2 {
			pixels = unsafe.Slice((*uint16)(unsafe.Pointer(&data[0])), count)
			data = data[2*count:]
		}
		return
	}
	for y := 0; y < skipY; y++ {
		for x := 0; x < width; {
			_, n, _ := nextRun()
			x += n
		}
	}
	if height == 0 {
		return
	}
	rows := unsafe.Slice((*unsafe.Pointer)(unsafe.Pointer(legacyGlobals.nox_pixbuffer_rows_3798784)), pos.Y+1)
	dst := unsafe.Add(rows[pos.Y], 2*(pos.X+extra))
	pitch := Nox_getBackbufferPitch()
	parity := Sub_473970(image.Pt(0, pos.Y)).Y
	high := nox_client_highResFrontWalls_80820 != 0
	var previous unsafe.Pointer
	copyPixels := 0
	for y := 0; y < height; y++ {
		draw := high || parity&1 == 0
		countCopied := 0
		if !draw && previous != nil && copyPixels != 0 {
			copy(unsafe.Slice((*uint16)(dst), copyPixels), unsafe.Slice((*uint16)(previous), copyPixels))
		}
		for x := 0; x < width; {
			op, n, pixels := nextRun()
			end := x + n
			if draw && end > skipX && x < skipX+drawWidth {
				lo, hi := max(x, skipX), min(end, skipX+drawWidth)
				if op == 2 {
					for i := range color {
						color[i] = start[i] + uint32(lo)*step[i]
					}
					target := unsafe.Slice((*uint16)(unsafe.Add(dst, 2*(lo-skipX))), hi-lo)
					Sub_480860(target, pixels[lo-x:hi-x], hi-lo, color, step)
				}
				// Preserve legacy low-resolution copy length: an initial partial transparent
				// run contributes nothing, while later transparent runs count in full even
				// when their end extends past the draw interval.
				if x < skipX {
					if op == 2 {
						countCopied += hi - lo
					}
				} else if op == 2 {
					countCopied += hi - lo
				} else {
					countCopied += n
				}
			}
			x = end
		}
		previous = dst
		dst = unsafe.Add(dst, pitch)
		parity++
		copyPixels = countCopied
	}
}
