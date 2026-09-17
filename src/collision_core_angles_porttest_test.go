//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"math"
	"testing"
	"unsafe"
)

func TestCollisionCoreAngleDrain(t *testing.T) {
	o := newCollisionCoreOwner(t)
	a, b := newObjectXferSimple(t, o.s), newObjectXferSimple(t, o.s)
	savedA, savedB := a.UpdateData, b.UpdateData
	a.UpdateData = collisionCoreGuarded(t, o, 64)
	b.UpdateData = collisionCoreGuarded(t, o, 64)
	t.Cleanup(func() { a.UpdateData = savedA; b.UpdateData = savedB })
	data := [2][]byte{unsafe.Slice((*byte)(a.UpdateData), 64), unsafe.Slice((*byte)(b.UpdateData), 64)}
	ids := collisionCoreIDs(a, b)
	ids[a.UpdateData] = 2001
	ids[b.UpdateData] = 2002
	type row struct {
		Base, Angle, Stop int
		Torque            uint32
		Data              [2][]byte
		Queues            [3]uint32
	}
	var rows []row
	for _, base := range []int{0, 1, 20, 31} {
		for _, angle := range []int{0, 1, 31, 127, 252, 255} {
			for _, torque := range []float32{-100, -1, -0.01, 0, 0.01, 1, 100, math.Nextafter32(0.01, 1)} {
				for stop := 0; stop < 3; stop++ {
					o.resetQueues()
					for i, d := range data {
						clear(d)
						binary.LittleEndian.PutUint32(d[4:], uint32(base))
						binary.LittleEndian.PutUint16(d[40:], uint16(angle))
						v := torque
						if i == 1 {
							v = -torque
						}
						binary.LittleEndian.PutUint32(d[32:], math.Float32bits(v))
						current := 0
						if stop == 1 {
							current = (base + 12) % 32
						}
						if stop == 2 {
							current = (base + 20) % 32
						}
						binary.LittleEndian.PutUint32(d[12:], uint32(current))
					}
					legacy.PortTestCollisionCore("angle-queue", a, nil, nil, 0)
					legacy.PortTestCollisionCore("angle-queue", b, nil, nil, 0)
					legacy.PortTestCollisionCore("angle-queue", a, nil, nil, 0)
					if o.queues(ids)[0] != 2002 {
						t.Fatal("angular queue insertion/dedup")
					}
					legacy.PortTestCollisionCore("angles", nil, nil, nil, 0)
					r := row{Base: base, Angle: angle, Stop: stop, Torque: math.Float32bits(torque), Queues: o.queues(ids)}
					if r.Queues[0] != 0 {
						t.Fatal("angular head not drained")
					}
					for i, d := range data {
						if binary.LittleEndian.Uint32(d[28:]) != 0 || binary.LittleEndian.Uint32(d[32:]) != 0 {
							t.Fatal("angular membership/torque not reset")
						}
						now := binary.LittleEndian.Uint16(d[40:])
						if now >= 256 || binary.LittleEndian.Uint32(d[12:]) != uint32(now)/8 {
							t.Fatal("angular direction/index", now, binary.LittleEndian.Uint32(d[12:]))
						}
						v := torque
						if i == 1 {
							v = -torque
						}
						if stop == 1 && v > 0 || stop == 2 && v < 0 {
							if now != uint16(angle) {
								t.Fatal("angular stop crossed", base, angle, stop, v, now)
							}
						}
						delta := (int(now) - angle + 256) % 256
						if delta > 4 && delta < 252 {
							t.Fatal("angular speed cap", delta)
						}
						r.Data[i] = append([]byte(nil), d...)
						raw := binary.LittleEndian.Uint32(d[36:])
						if raw != 0 {
							binary.LittleEndian.PutUint32(r.Data[i][36:], collisionCoreID(t, ids, raw))
						}
					}
					rows = append(rows, r)
				}
			}
		}
	}
	spellbookCapture(t, "collision-core-angle-drain", rows, "c424b4a2de0c097032b8bd49b41c14b6f88914a4071e58b317df1ce24e98b62f")
}
