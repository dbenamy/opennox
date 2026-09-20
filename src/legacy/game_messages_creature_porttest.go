//go:build porttest

package legacy

import (
	"bytes"
	"encoding/binary"

	"github.com/opennox/opennox/v1/server"
)

type PortTestServerCreatureCommandSpec struct {
	Dispatch, Invoke bool
	Flags            uint32
	Code             uint16
}

func portTestServerCreatureCommand(proxy *portTestRoamOwnerServer, unit *server.Object, sp *PortTestMonsterStateSpec) {
	c := sp.ServerCommand
	owner := proxy.state.source
	pl := owner.UpdateDataPlayer().Player
	oldFlags := pl.Field3680
	pl.Field3680 = c.Flags
	oldList, oldPending := proxy.core.Objs.List, proxy.core.Objs.Pending
	oldNext, oldExtent := unit.ObjNext, unit.Extent
	proxy.core.Objs.List, proxy.core.Objs.Pending = unit, nil
	unit.ObjNext, unit.Extent = nil, 23
	restoreLookup := PortTestCreatureXferLookupOwner()
	defer func() {
		restoreLookup()
		proxy.core.Objs.List, proxy.core.Objs.Pending = oldList, oldPending
		unit.ObjNext, unit.Extent = oldNext, oldExtent
		if pl.Field3680 != c.Flags {
			panic("creature command changed player flags")
		}
		pl.Field3680 = oldFlags
	}()
	if c.Dispatch {
		data := []byte{120, 0, 0, byte(sp.Order)}
		binary.LittleEndian.PutUint16(data[1:], c.Code)
		before := bytes.Clone(data)
		if n := Nox_xxx_netOnPacketRecvServ_51BAD0_net_sdecode_switch(pl.PlayerIndex(), data, pl, owner, owner.UpdateData); n != 4 || !bytes.Equal(data, before) {
			panic("creature command length/input changed")
		}
	} else if c.Invoke {
		target := unit
		if sp.Broadcast {
			target = nil
		}
		monsterOrder(owner, target, sp.Order)
	}
}
