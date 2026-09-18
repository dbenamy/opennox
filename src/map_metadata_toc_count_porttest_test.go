//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestMapMetadataTOCCount(t *testing.T) {
	newObjectXferOwner(t)
	oldCount := objectTypeCode16ByInd_len
	t.Cleanup(func() { objectTypeCode16ByInd_len = oldCount })
	dir := t.TempDir()
	type record struct {
		Count int
		IO    legacy.PortTestMapSectionWire
		Table []uint16
	}
	var records []record
	for _, count := range []int{0, 1, 32767, 32768, 65535} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			objectTypeCode16ByInd = make([]uint16, 64)
			objectTypeCode16ByInd[0] = 42
			objectTypeCode16ByInd_len = 99
			wire := binary.LittleEndian.AppendUint16([]byte{1, 0}, uint16(count))
			for i := 0; i < count; i++ {
				wire = binary.LittleEndian.AppendUint16(wire, uint16(i+1))
				wire = append(wire, 0)
			}
			result := legacy.PortTestMapMetadata(legacy.PortTestMapSectionIO{Function: "toc", Read: true, Data: wire}, dir)
			if result.Return != 1 || result.Position != int64(len(wire)) || !bytes.Equal(result.Data, wire) {
				t.Fatalf("return %d position %d want %d", result.Return, result.Position, len(wire))
			}
			if objectTypeCode16ByInd[0] != uint16(count) || objectTypeCode16ByInd_len != 0 {
				t.Fatalf("table %v count %d", objectTypeCode16ByInd, objectTypeCode16ByInd_len)
			}
			for _, v := range objectTypeCode16ByInd[1:] {
				if v != 0 {
					t.Fatal("uncleared table")
				}
			}
			records = append(records, record{count, result, append([]uint16(nil), objectTypeCode16ByInd...)})
		})
	}
	spellbookCapture(t, "map-metadata-toc-count", records, "a5476323127608808b41982d4718fcd6084cd8c5f149eee0b2f8ef1b500be715")
}
