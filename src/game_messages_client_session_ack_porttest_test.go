//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestGameMessageClientSessionAcknowledge(t *testing.T) {
	o := newReliableReportsOwner(t)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	type row struct {
		On, Recipient, Target, Return int
		Ordered                       bool
		Frame                         uint32
		State                         legacy.PortTestReliableReportState
	}
	var rows []row
	for on := 0; on < 2; on++ {
		for _, recipient := range []int{1, 7, 31} {
			for _, target := range []int{1, 7, 31, 255} {
				for _, ordered := range []bool{false, true} {
					for _, frame := range []uint32{0, 122, 123, 124, 0xffffffff} {
						setup := func() {
							o.reset()
							binary.LittleEndian.PutUint32(connected, uint32(on))
							legacy.PortTestReliableReports(0, 0, 0, 0, nil, nil, 0, false)
							for i := 0; i < 3; i++ {
								legacy.PortTestReliableReports(8, target, 0, 0, []byte{77, byte(i)}, nil, 0, ordered)
							}
							for _, to := range []int{1, 7, 31} {
								legacy.PortTestReliableReports(19, to, 0, 1, nil, nil, 0, false)
							}
						}
						setup()
						if on != 0 {
							legacy.PortTestReliableReports(16, recipient, 0, frame, nil, nil, 0, false)
						}
						want := o.state()
						setup()
						data := binary.LittleEndian.AppendUint32([]byte{171}, frame)
						input := bytes.Clone(data)
						n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(ntype.PlayerInd(recipient), netmsg.Op(171), data)
						got := o.state()
						if n != 5 || !bytes.Equal(input, data) || !reflect.DeepEqual(got, want) {
							t.Fatal("ack route/frame/connection", on, recipient, target, ordered, frame)
						}
						rows = append(rows, row{on, recipient, target, n, ordered, frame, got})
					}
				}
			}
		}
	}
	// This reserved message consumes its two payload bytes without changing state.
	for value := 0; value < 256; value++ {
		before := o.state()
		data := []byte{174, byte(value), byte(255 - value)}
		input := bytes.Clone(data)
		n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(31, netmsg.Op(174), data)
		if n != 3 || !bytes.Equal(data, input) || !reflect.DeepEqual(o.state(), before) {
			t.Fatal("outgoing-client no-op", value)
		}
	}
	interactionCapture(t, "game-client-session-ack", rows)
}
