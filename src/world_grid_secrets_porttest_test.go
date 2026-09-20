//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"reflect"
	"testing"
)

func TestWorldGridSecretWalls(t *testing.T) {
	type capture struct {
		IDs     []uint16
		Ops     []legacy.PortTestSecretOperation
		Results []legacy.PortTestSecretResult
	}
	var rows []capture
	for _, ids := range [][]uint16{{0, 1, 2, 3}, {7, 7, 7, 7}, {32767, 32768, 65534, 65535}} {
		for a := 0; a < 4; a++ {
			for b := 0; b < 4; b++ {
				if b == a {
					continue
				}
				for c := 0; c < 4; c++ {
					if c == a || c == b {
						continue
					}
					d := 6 - a - b - c
					for mode := 0; mode < 5; mode++ {
						ops := []legacy.PortTestSecretOperation{{Op: 2, ID: 7}, {Op: 4}, {Op: 1, ID: a}}
						for _, i := range []int{a, b, c, d} {
							ops = append(ops, legacy.PortTestSecretOperation{Op: 0, ID: i})
						}
						for _, id := range []int{0, 1, 2, 3, 7, 32767, 32768, 65534, 65535, -1} {
							ops = append(ops, legacy.PortTestSecretOperation{Op: 2, ID: id})
						}
						if mode < 4 {
							ops = append(ops, legacy.PortTestSecretOperation{Op: 1, ID: mode}, legacy.PortTestSecretOperation{Op: 1, ID: mode})
						}
						ops = append(ops, legacy.PortTestSecretOperation{Op: 3}, legacy.PortTestSecretOperation{Op: 3}, legacy.PortTestSecretOperation{Op: 4}, legacy.PortTestSecretOperation{Op: 2, ID: 7})
						got := legacy.PortTestSecretWalls(ids, ops)
						var list []int
						for i, op := range ops {
							want := -1
							switch op.Op {
							case 0:
								want = op.ID
								list = append([]int{op.ID}, list...)
							case 1:
								for j, v := range list {
									if v == op.ID {
										want = v
										list = append(list[:j], list[j+1:]...)
										break
									}
								}
							case 2:
								for _, v := range list {
									if int(ids[v]) == int(int16(op.ID)) {
										want = v
										break
									}
								}
							case 3:
								list = nil
							}
							if len(list) == 0 {
								list = nil
							}
							if got[i].Return != want || !reflect.DeepEqual(got[i].List, list) || !got[i].PayloadOK {
								t.Fatalf("ids%v order%v mode%d op%d %+v: got%+v want return%d list%v", ids, []int{a, b, c, d}, mode, i, op, got[i], want, list)
							}
						}
						rows = append(rows, capture{ids, ops, got})
					}
				}
			}
		}
	}
	drawableStateCapture(t, "world-secrets", rows, "b1ce975c8e3d95e3e6214cd784341e3c0acbf1f2f2c9c1814423df6f6e8e69b8")
}
