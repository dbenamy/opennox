//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/object"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"math"
	"testing"
	"unsafe"
)

type objectReportEncodingRow struct {
	Name   string
	Return int
	Packet []byte
}

func TestObjectReportsSimpleEncoding(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	u := &units[1]
	var rows []objectReportEncodingRow
	defer func() {
		spellbookCapture(t, "object-reports-simple", rows, "7837a08b6508021c13b225a0acf1a8c072e114e1e4253e86c64900e94fc6d1e9")
	}()
	coords := []float32{0, -0.75, 0.75, 1.99, -1.99, 32767, 32768, 65535, 65536, -65537, 2147483520, 2147483648, float32(math.Inf(1)), float32(math.NaN())}
	for _, slot := range []int{1, 7, 31} {
		for _, class := range []object.Class{0, object.ClassSimple, object.ClassClientPersist, object.ClassImmobile} {
			for _, code := range []uint32{0, 1, 32767, 32768, 65535} {
				for i, x := range coords {
					name := fmt.Sprintf("slot%d/class%x/code%d/coord%d", slot, uint32(class), code, i)
					t.Run(name, func(t *testing.T) {
						s.NetList.ResetAll()
						u.ObjClass = class
						u.TypeInd = 0xabcd
						u.NetCode = code
						u.Extent = code
						u.PosVec = types.Pointf{x, coords[(i+3)%len(coords)]}
						rv := legacy.PortTestObjectReports(2, nil, u, slot, 0, 0, nil)
						got := s.NetList.CopyPacketsA(ntype.PlayerInd(slot), netlist.Kind2)
						want := make([]byte, 9)
						want[0] = 47
						netcode := uint16(code)
						if code >= 32768 {
							netcode = 0
						} else if class.HasAny(object.ClassClientPersist | object.ClassImmobile) {
							netcode |= 0x8000
						}
						binary.LittleEndian.PutUint16(want[1:], netcode)
						binary.LittleEndian.PutUint16(want[3:], 0xabcd)
						word := func(v float32) uint16 {
							if math.IsNaN(float64(v)) || v >= 2147483648 || v < -2147483648 {
								return 0
							}
							return uint16(int32(v))
						}
						binary.LittleEndian.PutUint16(want[5:], word(u.PosVec.X))
						binary.LittleEndian.PutUint16(want[7:], word(u.PosVec.Y))
						if rv != 1 || !bytes.Equal(got, want) {
							t.Fatalf("return=%d packet=%x want=%x", rv, got, want)
						}
						for _, other := range []int{1, 7, 31} {
							if other != slot && s.NetList.ByInd(ntype.PlayerInd(other), netlist.Kind2).Count() != 0 {
								t.Fatal("wrong recipient")
							}
						}
						if len(visibilityEffectsPackets(s)[slot]) != 0 {
							t.Fatal("wrong queue")
						}
						rows = append(rows, objectReportEncodingRow{name, rv, got})
					})
				}
			}
		}
	}
}

func TestObjectReportsPhantomDirection(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	u := &units[1]
	threshold := legacy.PortTestObjectReportsDirectionThreshold()
	old := *threshold
	*threshold = 0
	t.Cleanup(func() { *threshold = old })
	var rows []objectReportEncodingRow
	defer func() {
		spellbookCapture(t, "object-reports-phantom", rows, "1b317d3a2391bb76627a1b2d69c3fc7a6269e018fae665e8ac7e3414fb7dcd49")
	}()
	for _, class := range []object.Class{0, object.ClassMissile} {
		for _, sub := range []uint32{0, 0x10, 0x20, 0x30, 0x40} {
			for d := 0; d < 256; d++ {
				name := fmt.Sprintf("class%d/sub%d/dir%d", class, sub, d)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					u.ObjClass = class
					u.ObjSubClass = object.SubClass(sub)
					u.TypeInd = 0x3456
					u.NetCode = 123
					u.PosVec = types.Pointf{10.9, -20.9}
					u.Direction1 = server.Dir16(d)
					rv := legacy.PortTestObjectReports(1, nil, u, 7, 0, 0, nil)
					got := s.NetList.CopyPacketsA(7, netlist.Kind2)
					if rv != 1 || len(got) != 11 || !bytes.Equal(got[:9], []byte{48, 123, 0, 0x56, 0x34, 10, 0, 236, 255}) {
						t.Fatalf("return=%d packet=%x", rv, got)
					}
					extra := byte(255)
					if class.Has(object.ClassMissile) && sub&0x30 != 0 {
						extra = byte(d >> 3)
					}
					if got[10] != extra || got[9]&15 != 0 || got[9] > 112 {
						t.Fatalf("direction/extra %x", got[9:])
					}
					// Cardinal directions independently pin down the packed orientation ordering.
					cardinal := map[int]byte{0: 64, 64: 96, 128: 48, 192: 16}
					if want, ok := cardinal[d]; ok && got[9] != want {
						t.Fatalf("cardinal %d got %d want %d", d, got[9], want)
					}
					rows = append(rows, objectReportEncodingRow{name, rv, got})
				})
			}
		}
	}
}

// Own lazy monster lookup caches reached by report animation helpers.
func objectReportsPlayers(t *testing.T) (*server.Server, []server.Object, func(int)) {
	s, units, configure := visibilityEffectsPlayers(t)
	caches := unsafe.Slice(memmap.PtrUint32(0x5D4594, 2488524), 4)
	old := append([]uint32(nil), caches...)
	clear(caches)
	t.Cleanup(func() { copy(caches, old) })
	threshold := legacy.PortTestObjectReportsDirectionThreshold()
	oldThreshold := *threshold
	*threshold = 6
	t.Cleanup(func() { *threshold = oldThreshold })
	return s, units, configure
}
