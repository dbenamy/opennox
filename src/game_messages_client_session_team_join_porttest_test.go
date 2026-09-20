//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/object"
	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestGameMessageClientSessionTeamJoin(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	t.Cleanup(c.srv.PortTestObjectiveTypes(nil, nil, nil))
	t.Cleanup(c.srv.PortTestMapDrawableTeamMessages())
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	oldCode := legacy.ClientPlayerNetCode()
	legacy.ClientSetPlayerNetCode(999)
	t.Cleanup(func() { legacy.ClientSetPlayerNetCode(oldCode) })
	type row struct {
		On, Existing, Static, Present int
		Type                          uint16
		ID                            uint32
		Return, Objects               int
		Member                        byte
		Count                         uint32
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for existing := 0; existing < 2; existing++ {
			for static := 0; static < 2; static++ {
				for present := 0; present < 2; present++ {
					for _, typ := range []uint16{0, 4} {
						for _, id := range []uint32{1, 0x10001, 255} {
							// Original C forms a member pointer from a failed sprite allocation. Only
							// exercise defined paths here; native allocation-failure coverage is separate.
							if on != 0 && existing == 0 && typ == 0 && present != 0 && byte(id) == 1 {
								continue
							}
							c.resetCase(env, pix, 1, 100)
							c.srv.Teams.Reset()
							c.srv.Teams.ActiveCnt = 0
							var tm *server.Team
							if present != 0 {
								tm = c.srv.Teams.Create(1)
							}
							binary.LittleEndian.PutUint32(connected, uint32(on))
							var dr *client.Drawable
							if existing != 0 {
								dr = c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(30, 40))
								if dr == nil {
									t.Fatal("drawable")
								}
								dr.ObjClass = 0
								if static != 0 {
									dr.ObjClass = object.Class(0x20400000)
								}
								dr.NetCode32 = 17
							}
							code := uint16(17) | uint16(static<<15)
							data := make([]byte, 10)
							data[0], data[1] = 196, 1
							binary.LittleEndian.PutUint32(data[2:], id)
							binary.LittleEndian.PutUint16(data[6:], code)
							binary.LittleEndian.PutUint16(data[8:], typ)
							input := bytes.Clone(data)
							n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data)
							count := 0
							for it := c.Objs.List1; it != nil; it = it.NextPtr {
								count++
								dr = it
							}
							wantObjects := existing
							if on != 0 && existing == 0 && typ != 0 {
								wantObjects = 1
							}
							member := byte(0)
							teamCount := uint32(0)
							if on != 0 && present != 0 && byte(id) == 1 && wantObjects != 0 {
								member = 1
								teamCount = 1
							}
							if n != 10 || !bytes.Equal(input, data) || count != wantObjects {
								t.Fatal("team join dispatch", on, existing, static, present, typ, id, n, count)
							}
							if dr != nil && byte(dr.TeamPtr().ID) != member {
								t.Fatal("team membership", on, existing, static, present, typ, id)
							}
							if tm != nil {
								if objectXferGetWord(tm.C(), 48) != teamCount {
									t.Fatal("team count")
								}
								head := objectXferGetWord(tm.C(), 44)
								if teamCount == 0 && head != 0 || teamCount != 0 && head != uint32(uintptr(unsafe.Pointer(dr.TeamPtr()))) {
									t.Fatal("team list head")
								}
							}
							rows = append(rows, row{on, existing, static, present, typ, id, n, count, member, teamCount})
						}
					}
				}
			}
		}
	}
	interactionCapture(t, "game-client-session-team-join", rows)
}
