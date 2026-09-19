//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/common/memmap/nox/blobdata"
	"github.com/opennox/opennox/v1/legacy"
	"testing"
)

func TestPlayerDeathCorpseCache(t *testing.T) {
	s := newObjectXferOwner(t)
	suffixData := blobdata.PortTestPlayerCorpseSuffixes()
	suffixStorage := serverConfigOwnBytes(t, 0x587000, 281216, 32)
	copy(suffixStorage, suffixData)
	parts := []string{"Skull", "RibCage", "Pelvis", "LeftLowerLeg", "LeftUpperLeg", "LeftLowerArm", "LeftUpperArm", "RightLowerLeg", "RightUpperLeg", "RightLowerArm", "RightUpperArm"}
	dirs := []int{0, 1, 2, 3, 5, 6, 7, 8}
	names := make([]string, 0, 88)
	for i := range dirs {
		raw := suffixData[4*i : 4*i+4]
		end := bytes.IndexByte(raw, 0)
		if end < 0 {
			t.Fatal("unterminated shipped suffix")
		}
		for _, part := range parts {
			names = append(names, "Corpse"+part+string(raw[:end]))
		}
	}
	owned := serverConfigOwnBytes(t, 0x5D4594, 2488728, 416)
	type record struct {
		Name  string
		State []byte
	}
	var rows []record
	for _, missing := range []int{-1, 0, 10, 44, 87} {
		for _, seed := range []uint32{0, 1, 0xffffffff} {
			name := fmt.Sprintf("missing=%d/seed=%x", missing, seed)
			t.Run(name, func(t *testing.T) {
				var absent []string
				if missing >= 0 {
					absent = []string{names[missing]}
				}
				t.Cleanup(s.PortTestRewardTypes(names, absent, true, 0, 0))
				for i := range owned {
					owned[i] = 0xa5
				}
				binary.LittleEndian.PutUint32(owned[8:], seed)
				want := bytes.Clone(owned)
				binary.LittleEndian.PutUint32(want[8:], 1)
				for i, dir := range dirs {
					for j := range parts {
						binary.LittleEndian.PutUint32(want[12+4*(11*dir+j):], uint32(s.Types.IndByID(names[11*i+j])))
					}
				}
				if s.Types.IndByID(names[1]) == 0 {
					t.Fatal("fixture must install a real corpse definition")
				}
				legacy.PortTestPlayerCorpseCache()
				if !bytes.Equal(owned, want) {
					for i := range want {
						if owned[i] != want[i] {
							t.Fatalf("cache byte %d got %x want %x", i, owned[i], want[i])
						}
					}
				}
				if !bytes.Equal(suffixStorage, suffixData) {
					t.Fatal("suffix table changed")
				}
				rows = append(rows, record{name, bytes.Clone(owned)})
			})
		}
	}
	spellbookCapture(t, "player-death-corpse-cache", rows, "0bd912a573c12f4630d349fa5a93e0211ece9909b57529c5bc54dd4d197eba47")
}
