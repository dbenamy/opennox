//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"image"
	"testing"
	"unsafe"

	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/libs/wall"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

func TestGameMessageClientSecretWalls(t *testing.T) {
	c, _, _ := newEffectsFullOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	t.Cleanup(c.srv.PortTestMinimapWalls())
	w := c.srv.Walls.CreateAtGrid(image.Pt(0, 0))
	if w == nil {
		t.Fatal("wall allocation")
	}
	t.Cleanup(legacy.PortTestGameSecretWall(w.C()))
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	payload, free := alloc.Make([]byte{}, 32)
	t.Cleanup(free)
	type row struct {
		On, Host, Kind, ID, Match, Secret, Data, Return int
		Payload                                         []byte
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for host := 0; host < 2; host++ {
			restoreFlags := noxflags.PortTestGameFlags(noxflags.GameFlag(host))
			for _, kind := range []int{59, 60} {
				for _, id := range []uint16{7, 32767, 32768, 65535} {
					for match := 0; match < 2; match++ {
						for secret := 0; secret < 2; secret++ {
							modes := []int{1}
							if kind == 60 {
								modes = []int{0, 1}
							}
							for _, hasData := range modes {
								binary.LittleEndian.PutUint32(connected, uint32(on))
								w.Field10 = id
								w.Flags4 = wall.Flags(secret * 4)
								for i := range payload {
									payload[i] = byte(i*7 + 3)
								}
								w.Data = nil
								if hasData != 0 {
									w.Data = unsafe.Pointer(&payload[0])
								}
								code := id
								if match == 0 {
									code ^= 1
								}
								data := []byte{byte(kind), byte(code), byte(code >> 8)}
								before := bytes.Clone(data)
								want := bytes.Clone(payload)
								wallBefore := bytes.Clone(unsafe.Slice((*byte)(w.C()), 36))
								if on != 0 && host == 0 && match != 0 && id <= 32767 && secret != 0 && hasData != 0 {
									if kind == 59 {
										want[21], want[22] = 3, 23
									} else {
										want[21], want[22] = 1, 0
									}
								}
								n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(kind), data)
								if n != 3 || !bytes.Equal(data, before) || !bytes.Equal(payload, want) || !bytes.Equal(unsafe.Slice((*byte)(w.C()), 36), wallBefore) {
									t.Fatalf("secret on%d host%d kind%d id%d match%d flag%d data%d ret%d payload%x want%x", on, host, kind, id, match, secret, hasData, n, payload, want)
								}
								rows = append(rows, row{on, host, kind, int(id), match, secret, hasData, n, bytes.Clone(payload)})
							}
						}
					}
				}
			}
			restoreFlags()
		}
	}
	interactionCapture(t, "game-client-secret-walls", rows)
}
