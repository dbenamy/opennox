//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func TestCreatureXferDefinitionDefaults(t *testing.T) {
	s := newCreatureXferOwner(t)
	type row struct {
		Case   string
		Health server.HealthData
		Data   []byte
	}
	var rows []row
	defer func() {
		spellbookCapture(t, "creature-xfer-definition-defaults", rows, "2694cdd92e3157113a7851b6462f234742057c5ddeaa4bf3166b3cf4587910e3")
	}()
	for _, present := range []bool{false, true} {
		for _, mask := range []uint32{0, 0xffffffff, 0x19c40, 0x20, 0x1800, 0xffc00000} {
			for _, defaults := range []byte{0, 1, 2} {
				for _, healthSaved := range []byte{0, 1, 255} {
					t.Run(fmt.Sprintf("present%v-mask%08x-default%d-health%d", present, mask, defaults, healthSaved), func(t *testing.T) {
						u := newCreatureXferObject(t, s, "Monster")
						def, free := alloc.New(server.MonsterDef{})
						defer free()
						// Include a nonmatching head to exercise actual linked-list lookup.
						head, freeHead := alloc.New(server.MonsterDef{})
						defer freeHead()
						head.TypeInd240 = 0xffffffff
						head.Next244 = def
						def.TypeInd240 = uint32(u.TypeInd)
						if !present {
							def.TypeInd240 = 0xfffffffe
						}
						def.Health68 = 0x12345678
						def.StatusFlags92 = object.MonsterStatus(mask)
						def.RetreatRatio80 = 0.25
						def.ResumeRatio84 = 0.75
						restore := legacy.PortTestCreatureXferDefinitions(unsafe.Pointer(head))
						defer restore()
						objectXferSetWord(u.UpdateData, 1440, 0xa5a5a5a5)
						objectXferSetWord(u.UpdateData, 1336, 0x3f000000)
						objectXferSetWord(u.UpdateData, 1344, 0x3f800000)
						*(*byte)(unsafe.Add(u.UpdateData, 1444)) = defaults
						*(*byte)(unsafe.Add(u.UpdateData, 1445)) = healthSaved
						*(*byte)(unsafe.Add(u.UpdateData, 1340)) = defaults
						*(*byte)(unsafe.Add(u.UpdateData, 1348)) = defaults
						// Auto-spell assignment needs separate owner coverage.
						*(*byte)(unsafe.Add(u.UpdateData, 2036)) = 0
						data := unsafe.Slice((*byte)(u.UpdateData), int(unsafe.Sizeof(server.MonsterUpdateData{})))
						want := append([]byte(nil), data...)
						hp := *u.HealthData
						if present {
							if healthSaved == 0 {
								hp.Cur = 0x5678
								hp.Max = 0x5678
								hp.Field2 = 0x5678
							}
							flags := uint32(0xa5a5a5a5)
							if defaults == 1 {
								flags = mask
							} else {
								const replace = uint32(0x3fffff &^ 0x19c40)
								flags = (flags &^ replace) | (mask & replace)
							}
							if flags&0x20 == 0 {
								flags &^= 0x1800
							}
							binary.LittleEndian.PutUint32(want[1440:], flags)
							if defaults == 1 {
								binary.LittleEndian.PutUint32(want[1336:], 0x3e800000)
								binary.LittleEndian.PutUint32(want[1344:], 0x3f400000)
							}
						}
						legacy.PortTestCreatureXferHelper(4, u, nil, 0)
						if *u.HealthData != hp || !bytes.Equal(data, want) {
							t.Fatalf("definition policy differs: health=%+v want=%+v", *u.HealthData, hp)
						}
						rows = append(rows, row{t.Name(), *u.HealthData, append([]byte(nil), data...)})
					})
				}
			}
		}
	}
	legacy.PortTestCreatureXferHelper(4, nil, nil, 0)
}
