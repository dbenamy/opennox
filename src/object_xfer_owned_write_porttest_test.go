//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestObjectXferOwnedReferencesWrite(t *testing.T) {
	core := newObjectXferOwner(t)
	path := filepath.Join(t.TempDir(), "owned.bin")
	var captures []objectXferCaptureRow
	defer func() { spellbookCapture(t, "object-xfer-owned-write", captures, "44c7ed44a24050b30e90cda145d94733de7ebafbe3df1dc825019b94cc948925") }()
	for _, count := range []int{0, 1, 5} {
		for _, inventory := range []int{0, 1, 3} {
			for _, hidden := range []bool{false, true} {
				for _, destroyed := range []bool{false, true} {
					t.Run(fmt.Sprintf("refs%d-inv%d-hidden%v-destroyed%v", count, inventory, hidden, destroyed), func(t *testing.T) {
						u := newObjectXferSimple(t, core)
						objectXferSetCommon(u)
						u.Field34 = 222
						defer objectXferCaptureCase(t, &captures, u)
						core.Types.ByInd(2).Icon = 0
						if hidden {
							core.Types.ByInd(2).Icon = -1
						}
						var included []uint32
						var previous *server.Object
						for i := 0; i < count; i++ {
							typ := 1 + i%2
							child := core.NewObjectByTypeInd(typ)
							trackObjectXferTyped(t, core, child)
							child.ScriptIDVal = 101 + i
							if destroyed && i%3 == 0 {
								child.ObjFlags |= object.FlagDestroyed
							}
							if previous == nil {
								u.Field129 = child
							} else {
								previous.Field128 = child
							}
							previous = child
							if !(destroyed && i%3 == 0) && !(hidden && typ == 2) {
								included = append(included, uint32(child.ScriptIDVal))
							}
						}
						for i := 0; i < inventory; i++ {
							child := newObjectXferSimple(t, core)
							child.InvNextItem = u.InvFirstItem
							child.InvHolder = u
							if u.InvFirstItem != nil {
								u.InvFirstItem.Field125 = child
							}
							u.InvFirstItem = child
						}
						var want mapDrawableStream
						want.u16(64)
						want.u32(u.Extent)
						want.u32(88)
						want.f32(64.25)
						want.f32(128.75)
						if count == 0 && inventory == 0 {
							want.u8(0)
						} else {
							want.u8(255)
							want.u32(0)
							want.u8(0)
							want.u8(0)
							want.u8(byte(inventory))
							want.u16(uint16(len(included)))
							for _, id := range included {
								want.u32(id)
							}
							want.u32(0)
							objectXferScript(&want, 0)
							want.u32(99)
						}
						if err := cryptfile.OpenGlobal(path, cryptfile.WriteOnly, -1); err != nil {
							t.Fatal(err)
						}
						if legacy.Nox_xxx_mapReadWriteObjData_4F4530(u, 60) != 1 {
							t.Fatal("owned writer failed")
						}
						if err := cryptfile.Close(); err != nil {
							t.Fatal(err)
						}
						got, err := os.ReadFile(path)
						if err != nil {
							t.Fatal(err)
						}
						if !bytes.Equal(got, want.Bytes()) {
							t.Fatalf("owned bytes=%x want=%x", got, want.Bytes())
						}
					})
				}
			}
		}
	}
}
