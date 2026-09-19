//go:build porttest

package legacy

import "github.com/opennox/opennox/v1/legacy/common/alloc"

func PortTestUnitActionName(id int32, restricted bool) (string, bool) {
	p := unitActionName(id, restricted)
	return alloc.GoString(p), p != nil
}
func PortTestUnitActionIndex(name string, restricted bool) int32 {
	return unitActionIndex(name, restricted)
}
