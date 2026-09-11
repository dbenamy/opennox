//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"

	"github.com/opennox/opennox/v1/legacy"
)

func TestNetworkAliasCallers(t *testing.T) {
	for _, stream := range []int{1, 2} {
		for _, frame := range []uint32{0, 100, 0xffffffef} {
			for slot := 0; slot < 255; slot++ {
				t.Run(fmt.Sprintf("stream%d/frame%x/slot%d", stream, frame, slot), func(t *testing.T) {
					const key1, key2 = uint16(0x3401), uint16(9)
					initial := make([]byte, 255*8)
					for i := 0; i < 255; i++ {
						binary.LittleEndian.PutUint16(initial[i*8:], 0xaabb)
						binary.LittleEndian.PutUint16(initial[i*8+2:], 0xccdd)
						binary.LittleEndian.PutUint32(initial[i*8+4:], 0xffffffff)
					}
					if slot != 0 {
						binary.LittleEndian.PutUint16(initial[slot*8:], key1)
						binary.LittleEndian.PutUint16(initial[slot*8+2:], key2)
					}
					got := legacy.PortTestAliasCallers(initial, []legacy.PortTestAliasCallerCall{{Stream: stream, Frame: frame, Key1: key1, Key2: key2}})[0]
					wantTable := append([]byte(nil), initial...)
					var wantMsg []byte
					if slot != 0 {
						expiry := uint32(0xffffffff)
						if stream == 2 {
							expiry = frame + 60
						}
						binary.LittleEndian.PutUint32(wantTable[slot*8+4:], expiry)
						wantMsg = []byte{0xa5, byte(slot), 0x01, 0x34, 9, 0, 0, 0, 0, 0}
						binary.LittleEndian.PutUint32(wantMsg[6:], expiry)
					}
					if !bytes.Equal(got.Table, wantTable) || !bytes.Equal(got.Guard, bytes.Repeat([]byte{0xd7}, 8)) || !bytes.Equal(got.Announcements, wantMsg) || got.AnnouncementCnt != len(wantMsg)/10 {
						t.Errorf("alias storage/message mismatch: guard=%x announcements=%x count=%d", got.Guard, got.Announcements, got.AnnouncementCnt)
					}
					pos := [2]int32{321, 654}
					consumed := uint32(11)
					camera := [2]int{321, 654}
					cameraCalls := 1
					if stream == 2 {
						pos = [2]int32{305, 393}
						consumed = 8
						camera = [2]int{-999, -999}
						cameraCalls = 0
					}
					calls := [][4]int{{int(key2), int(key1), int(pos[0]), int(pos[1])}}
					if got.ReturnRaw != consumed || got.Position != pos || got.Camera != camera || got.CameraCalls != cameraCalls || !reflect.DeepEqual(got.SpriteCalls, calls) || got.SpriteFrame != frame || got.SpriteAnim != 7 || got.SpriteDirection != 0 {
						t.Errorf("packet processing: ret=%x pos=%v camera=%v/%d sprite=%v frame=%x anim=%d dir=%d", got.ReturnRaw, got.Position, got.Camera, got.CameraCalls, got.SpriteCalls, got.SpriteFrame, got.SpriteAnim, got.SpriteDirection)
					}
				})
			}
		}
	}
}
