//go:build porttest

package legacy

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/common/ntype"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

type PortTestUnitReadSpec struct {
	Warp, NonPlayer, Blocked                         bool
	Frame, FPS, Timestamp, Stage, Allowed, Threshold uint32
	Text                                             string
}

func (p *portTestShopPools) unitReadContract() []uint32 {
	sp := p.proxy.callbacks.shop.spec.TemporaryUpdates.World.Objectives.Attack.Controls.UnitRead
	u, it := p.resources.unit, p.items[0].u
	oldClass := u.ObjClass
	defer func() { u.ObjClass = oldClass }()
	u.PosVec = types.Ptf(100, 100)
	it.PosVec = types.Ptf(200, 100)
	if p.proxy.core.MapTraceVision(u, it) == sp.Blocked {
		panic("read fixture visibility")
	}
	name := "ReadUse"
	if sp.Warp {
		name = "WarpReadUse"
	}
	fn, size := server.PortTestUnitGameplayRegistration(name, true)
	if size != 260 {
		panic("read use registration data size")
	}
	data := p.objectiveRegion(int(size))
	b := unsafe.Slice((*byte)(data), 260)
	copy(b, sp.Text)
	binary.LittleEndian.PutUint32(b[256:], sp.Timestamp)
	p.identify(fn, 93066+uint32(bool2int(sp.Warp)))
	it.Use = server.UseFuncPtr{Ptr: fn}
	oldUseData := it.UseData.Ptr
	defer func() { it.UseData.Ptr = oldUseData }()
	it.UseData.Ptr = data
	allowed, stage := memmap.PtrUint32(0x5D4594, 1556120), memmap.PtrUint32(0x587000, 202028)
	oldAllowed, oldStage := *allowed, *stage
	defer func() { *allowed, *stage = oldAllowed, oldStage }()
	*allowed, *stage = sp.Allowed, sp.Stage
	p.proxy.core.SetFrame(sp.Frame)
	p.proxy.core.SetTickRate(sp.FPS)
	if sp.NonPlayer {
		u.ObjClass = 0
	}
	before := p.proxy.core.NetList.CopyPacketsA(ntype.PlayerInd(1), 1)
	if !it.Use.Get()(u, it) {
		panic("read use return")
	}
	active := !sp.NonPlayer && (sp.Timestamp == 0 || sp.Frame-sp.Timestamp > 3*sp.FPS) && !sp.Blocked
	wantTime := sp.Timestamp
	var want []byte
	if active {
		wantTime = sp.Frame
		if sp.Warp && sp.Allowed != 0 {
			want = []byte{169, 21, 0, 0, 0, 0}
			binary.LittleEndian.PutUint32(want[2:], sp.Threshold)
		} else {
			text := sp.Text
			if sp.Warp {
				text = "GeneralPrint:WarpClosed"
			}
			if len(text) > 0 && len(text) <= 48 {
				want = append([]byte{169, 15, 1}, []byte(text)...)
				want = append(want, 0)
			}
		}
	}
	if got := binary.LittleEndian.Uint32(b[256:]); got != wantTime {
		panic(fmt.Sprintf("read timestamp %d want %d", got, wantTime))
	}
	after := p.proxy.core.NetList.CopyPacketsA(ntype.PlayerInd(1), 1)
	expected := append(append([]byte(nil), before...), want...)
	if !bytes.Equal(after, expected) {
		panic(fmt.Sprintf("read message %x want %x", after, expected))
	}
	return []uint32{wantTime, uint32(len(want))}
}
