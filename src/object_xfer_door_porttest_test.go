//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/common/memmap"
)

func TestObjectXferDoorGeometry(t *testing.T) {
	core := newObjectXferOwner(t)
	mapDrawableTables(t)
	path := filepath.Join(t.TempDir(), "door.bin")
	var captures []objectXferCaptureRow
	defer func() { spellbookCapture(t, "object-xfer-door-geometry", captures, "20c9ae17e8ca3d17e152622396e7dab9776b3321b76c4763b12866577b9329ee") }()
	for _, direction := range []uint32{0, 7, 31} {
		for _, orientation := range []uint32{0, 1, 7, 16, 31} {
			for _, x := range []float32{-23, -0.5, 0, 11.5, 23, 46, 64.25, 128.75, 1000} {
				for _, y := range []float32{-0.5, 128.75} {
					t.Run(fmt.Sprintf("dir%d-orientation%d-x%g-y%g", direction, orientation, x, y), func(t *testing.T) {
						u := newObjectXferTyped(t, core, "Door")
						defer objectXferCaptureCase(t, &captures, u)
						stream := objectXferCurrentRecord(60)
						binary.LittleEndian.PutUint32(stream.Bytes()[12:], math.Float32bits(x))
						binary.LittleEndian.PutUint32(stream.Bytes()[16:], math.Float32bits(y))
						stream.u32(direction)
						stream.u32(71)
						stream.u32(orientation)
						if err := os.WriteFile(path, stream.Bytes(), 0600); err != nil {
							t.Fatal(err)
						}
						objectXferReadRecord(t, u, path)
						halfX := memmap.Int32(0x587000, 196184+8*uintptr(orientation)) / 2
						halfY := memmap.Int32(0x587000, 196188+8*uintptr(orientation)) / 2
						wantX := int32((float64(halfX) + float64(x)) * 0.043478262)
						wantY := int32((float64(halfY) + float64(y)) * 0.043478262)
						if objectXferGetWord(u.UpdateData, 12) != direction || objectXferGetWord(u.UpdateData, 8) != direction || objectXferGetWord(u.UpdateData, 4) != orientation || *(*uint16)(unsafe.Add(u.UpdateData, 40)) != uint16(direction*8) || *(*byte)(unsafe.Add(u.UpdateData, 1)) != 71 {
							t.Fatal("door direction/state")
						}
						if int32(objectXferGetWord(u.UpdateData, 16)) != wantX || int32(objectXferGetWord(u.UpdateData, 20)) != wantY {
							t.Fatalf("wall coordinates=%d/%d want=%d/%d", objectXferGetWord(u.UpdateData, 16), objectXferGetWord(u.UpdateData, 20), wantX, wantY)
						}
					})
				}
			}
		}
	}
}
