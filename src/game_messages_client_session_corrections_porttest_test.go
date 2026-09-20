//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"github.com/opennox/libs/noxnet/netmsg"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"strings"
	"testing"
)

// Original C's 20-unit temporary overflows at these sizes. These are native
// contracts, deliberately separate from the original-C capture corpus.
func TestGameMessageClientSessionLongTeamNames(t *testing.T) {
	o := newMatchRosterOwner(t)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	binary.LittleEndian.PutUint32(connected, 1)
	for _, count := range []int{20, 21, 254, 255} {
		o.s.Teams.Reset()
		o.s.Teams.ActiveCnt = 0
		data := make([]byte, 18+count*2)
		data[0] = 196
		data[2] = 1
		data[15] = byte(count)
		data[16] = 2
		for i := 0; i < count; i++ {
			binary.LittleEndian.PutUint16(data[18+2*i:], 'x')
		}
		before := bytes.Clone(data)
		n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, netmsg.Op(196), data)
		tm := o.s.Teams.ByID(1)
		if n != len(data) || !bytes.Equal(data, before) || tm == nil || tm.Name() != strings.Repeat("x", 20) {
			t.Fatalf("team name length=%d returned=%d", count, n)
		}
	}
}

func TestGameMessageClientSessionFailedDrawableCreation(t *testing.T) {
	c, pix, env := newEffectsFullOwner(t, "GreenZap")
	t.Cleanup(noxflags.PortTestGameFlags(0))
	t.Cleanup(c.srv.PortTestObjectiveTypes(nil, nil, nil))
	t.Cleanup(c.srv.PortTestMapDrawableTeamMessages())
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	cache := serverConfigOwnBytes(t, 0x5D4594, 1200904, 4)
	for _, code := range []uint16{17, 0x8011, 65535} {
		c.resetCase(env, pix, 1, 100)
		binary.LittleEndian.PutUint32(connected, 1)
		c.srv.Teams.Reset()
		c.srv.Teams.ActiveCnt = 0
		tm := c.srv.Teams.Create(1)
		data := []byte{196, 1, 1, 0, 0, 0, 0, 0, 0, 0}
		binary.LittleEndian.PutUint16(data[6:], code)
		before := bytes.Clone(data)
		if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, 196, data); n != 10 || !bytes.Equal(data, before) || c.Objs.List1 != nil || objectXferGetWord(tm.C(), 44) != 0 || objectXferGetWord(tm.C(), 48) != 0 {
			t.Fatalf("failed team drawable code=%d", code)
		}
	}
	for _, warm := range []bool{false, true} {
		c.resetCase(env, pix, 1, 100)
		c.FailEvery = 1
		binary.LittleEndian.PutUint32(connected, 1)
		clear(cache)
		if warm {
			binary.LittleEndian.PutUint32(cache, uint32(c.Things.IndByID("GreenZap")))
		}
		data := []byte{240, 16, 1, 0, 2, 0, 3, 0, 4, 0, 5, 0}
		before := bytes.Clone(data)
		if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, 240, data); n != 12 || !bytes.Equal(data, before) || c.Objs.List1 != nil {
			t.Fatal("failed green bolt allocation", warm, n)
		}
	}
}

func TestGameMessageClientSessionQuestNoticeSelectorBounds(t *testing.T) {
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	for on := 0; on < 2; on++ {
		binary.LittleEndian.PutUint32(connected, uint32(on))
		want := 52
		if on != 0 {
			want = -1
		}
		for selector := 5; selector <= 255; selector++ {
			data := make([]byte, 52)
			data[0], data[1], data[51] = 240, 33, byte(selector)
			before := bytes.Clone(data)
			if n := legacy.Nox_xxx_netOnPacketRecvCli_48EA70_switch(0, 240, data); n != want || !bytes.Equal(data, before) {
				t.Fatal("quest class selector", on, selector, n)
			}
		}
	}
}
