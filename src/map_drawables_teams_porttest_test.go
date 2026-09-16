//go:build porttest

package opennox

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
)

func TestMapDrawableTeamRegistration(t *testing.T) {
	type record struct {
		Modern, SpecialClass, Marker bool
		Flags                        uint32
		Team                         byte
		Members                      uint32
		Linked                       bool
		Cache                        uint32
	}
	var records []record
	for _, modern := range []bool{false, true} {
		for _, special := range []bool{false, true} {
			for _, flags := range []uint32{0, 1, 128, 129} {
				for _, marker := range []bool{false, true} {
					for _, team := range []byte{0, 1, 3} {
						t.Run(fmt.Sprintf("modern%v-special%v-flags%d-marker%v-team%d", modern, special, flags, marker, team), func(t *testing.T) {
							o := newObjectDrawingOwner(t, "FlagMarker")
							c := o.c
							t.Cleanup(c.srv.PortTestMapDrawableTeamMessages())
							t.Cleanup(noxflags.PortTestGameFlags(noxflags.GameFlag(flags)))
							gameBall := legacy.PortTestMapDrawableTeamWord()
							oldBall := *gameBall
							*gameBall = 1
							t.Cleanup(func() { *gameBall = oldBall })
							cache := memmap.PtrUint32(0x5D4594, 1309788)
							oldCache := *cache
							*cache = 0
							t.Cleanup(func() { *cache = oldCache })
							markerID := c.Things.IndByID("FlagMarker")
							id := 4
							if marker {
								id = markerID
							}
							typ := c.Things.TypeByInd(id)
							typ.ObjClass = 0
							if special {
								typ.ObjClass = object.Class(0x10000000)
							}
							var s mapDrawableStream
							if modern {
								s.modernBase(64, true, team, 0)
							} else {
								s.u16(5)
								s.oldBase(40, 5, team, 0)
							}
							path := filepath.Join(t.TempDir(), "team.bin")
							if err := os.WriteFile(path, s.Bytes(), 0600); err != nil {
								t.Fatal(err)
							}
							original := cryptfile.Global()
							cryptfile.SetGlobal(nil)
							defer func() { cryptfile.Close(); cryptfile.SetGlobal(original) }()
							if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
								t.Fatal(err)
							}
							var count uint32
							dr := legacy.PortTestMapDrawableBase(id, 40, 0, false, &count)
							if dr == nil {
								t.Fatal("actual drawable allocation failed")
							}
							if int(count) != s.Len() {
								t.Fatalf("consumed=%d want=%d", count, s.Len())
							}
							tm := c.srv.Teams.ByID(server.TeamID(team))
							var members uint32
							linked := false
							if tm != nil {
								members = *(*uint32)(unsafe.Add(tm.C(), 48))
								linked = *(*uint32)(unsafe.Add(tm.C(), 44)) == uint32(uintptr(dr.TeamPtr().C()))
							}
							allowed := flags&1 == 0 && team != 0
							if modern {
								allowed = allowed && !marker
							} else {
								allowed = allowed && (!special || flags&128 == 0)
							}
							if linked != allowed || members != uint32(bool2int(allowed)) {
								t.Fatalf("team registration: linked=%v count=%d want=%v", linked, members, allowed)
							}
							wantCache := uint32(0)
							if modern && flags&1 == 0 {
								wantCache = uint32(markerID)
							}
							if *cache != wantCache {
								t.Fatalf("flag marker cache=%d want=%d", *cache, wantCache)
							}
							records = append(records, record{modern, special, marker, flags, team, members, linked, *cache})
							// Only registered nodes belong to the actual team list. Other team bytes
							// are map attributes and must not be removed from an unrelated membership.
							if !linked {
								dr.TeamVal.ID = 0
							}
							c.Nox_xxx_spriteDeleteStatic_45A4E0_drawable(dr)
							if linked && (*(*uint32)(unsafe.Add(tm.C(), 44)) != 0 || *(*uint32)(unsafe.Add(tm.C(), 48)) != 0) {
								t.Fatal("actual team cleanup left the drawable linked")
							}
						})
					}
				}
			}
		}
	}
	spellbookCapture(t, "map-drawable-teams", records, "0a3b027d3095ebc8181e91ea904c0fe9d0f25851ec2c9e51b4e67b88656d23fe")
}
