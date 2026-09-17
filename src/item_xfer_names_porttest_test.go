//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestItemXferNameBoundaries(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "names.bin")
	var rows []itemXferCaptureRow
	defer func() {
		spellbookCapture(t, "item-xfer-name-boundaries", rows, "eef3c2673164b967d3bcb50983f239c98fcc7145f1b673a530e0070fb0f4156e")
	}()
	for _, name := range []string{"SpellReward", "AbilityReward", "FieldGuide"} {
		for _, n := range []int{0, 1, 31, 63, 64, 127, 128, 255} {
			t.Run(fmt.Sprintf("%s-n%d", name, n), func(t *testing.T) {
				u := newItemXferObject(t, s, name)
				u.Field34 = 999
				if name != "FieldGuide" {
					*(*byte)(u.UseData.Ptr) = 1
				}
				p := objectXferCurrentRecord(60)
				prefix := p.Len()
				p.u8(byte(n))
				p.WriteString(strings.Repeat("x", n))
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				err := u.CallXfer(nil)
				limit := 128
				if name == "FieldGuide" {
					limit = 64
				}
				reject := n >= limit
				if (err != nil) != reject {
					t.Fatalf("rejected=%v want=%v", err != nil, reject)
				}
				pos, seekErr := cryptfile.Global().File.Seek(0, 1)
				wantPos := p.Len()
				if reject {
					wantPos = prefix + 1
				}
				if seekErr != nil || pos != int64(wantPos) {
					t.Fatalf("position=%d want=%d err=%v", pos, wantPos, seekErr)
				}
				if reject {
					if u.Field34 != 0 {
						t.Fatal("failed callback restored lifetime unexpectedly")
					}
					if name != "FieldGuide" && *(*byte)(u.UseData.Ptr) != 1 {
						t.Fatal("rejected name mutated identifier")
					}
				} else {
					if u.Field34 != 999 {
						t.Fatal("successful callback lost lifetime")
					}
					if name == "FieldGuide" {
						if alloc.GoString((*byte)(u.UseData.Ptr)) != strings.Repeat("x", n) {
							t.Fatal("guide bytes")
						}
					} else if *(*byte)(u.UseData.Ptr) != 0 {
						t.Fatal("unknown reward name should resolve to zero")
					}
				}
				itemXferCaptureCase(t, &rows, u)
			})
		}
	}
}
func TestItemXferOldSpellNameRejection(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "old-names.bin")
	for _, v := range []int16{31, 40} {
		for rejectAt := 0; rejectAt < 3; rejectAt++ {
			t.Run(fmt.Sprintf("v%d-slot%d", v, rejectAt), func(t *testing.T) {
				u := newItemXferObject(t, s, "SpellReward")
				*(*byte)(u.UseData.Ptr) = 7
				var p mapDrawableStream
				p.u16(uint16(v))
				objectXferEmptyBase(&p, v)
				for i := 0; i < rejectAt; i++ {
					p.u8(1)
					p.u8('x')
				}
				p.u8(128)
				wantPos := p.Len()
				p.Write(make([]byte, 128))
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if u.CallXfer(nil) == nil {
					t.Fatal("oversized old spell name accepted")
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != int64(wantPos) {
					t.Fatalf("position=%d want=%d err=%v", pos, wantPos, err)
				}
				if *(*byte)(u.UseData.Ptr) != 7 {
					t.Fatal("partial rejected record changed selected spell")
				}
			})
		}
	}
}
func TestItemXferNonpositiveVersions(t *testing.T) {
	s := newItemXferOwner(t)
	path := filepath.Join(t.TempDir(), "nonpositive.bin")
	for _, name := range []string{"ToxicCloud", "MonsterGenerator", "RewardMarker"} {
		for _, v := range []int16{-32768, -1, 0} {
			t.Run(fmt.Sprintf("%s-v%d", name, v), func(t *testing.T) {
				u := newItemXferObject(t, s, name)
				before := *(*[193]uint32)(unsafe.Pointer(u))
				var p mapDrawableStream
				p.u16(uint16(v))
				p.u32(0x12345678)
				if err := os.WriteFile(path, p.Bytes(), 0600); err != nil {
					t.Fatal(err)
				}
				if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
					t.Fatal(err)
				}
				defer cryptfile.Close()
				if u.CallXfer(nil) == nil {
					t.Fatal("nonpositive version accepted")
				}
				pos, err := cryptfile.Global().File.Seek(0, 1)
				if err != nil || pos != 2 {
					t.Fatal("rejection read payload")
				}
				if before != *(*[193]uint32)(unsafe.Pointer(u)) {
					t.Fatal("rejected object changed")
				}
			})
		}
	}
}
