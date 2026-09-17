//go:build porttest

package legacy

/*
extern unsigned int nox_player_netCode_85319C;
*/
import "C"
import (
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/server"
	"unsafe"
)

func PortTestSpatialQuadrant(p *types.Pointf) int8                 { return spatialQuadrant(p) }
func PortTestSpatialTriangle(p *types.Pointf, grid [2]int32) int32 { return spatialTriangle(p, grid) }
func PortTestSpatialNormal(grid *[2]int32, ray *[4]float32, normal *types.Pointf) int32 {
	return spatialNormal(grid, ray, normal)
}
func PortTestSpatialEligible(a, b *server.Object, teamFilter bool) int32 {
	if teamFilter {
		return int32(bool2int(spatialTeamEligible(a, b)))
	}
	return int32(bool2int(spatialEligible(a, b)))
}
func PortTestSpatialPredict(a, b *server.Object, speed float32, out *types.Pointf) uint32 {
	return motionAddress(spatialPredict(a, b, speed, out))
}
func PortTestSpatialProbe(a *server.Object, next, prev *types.Pointf) uint32 {
	return motionAddress(spatialProbe(a, next, prev))
}
func PortTestSpatialCandidate(a, b *server.Object, next, prev *types.Pointf) uint32 {
	var hit *server.Object
	spatialCandidate(a, b, next, prev, &hit)
	return motionAddress(hit)
}
func PortTestSpatialRay(ray *[4]float32) int32                         { return int32(bool2int(spatialRay(ray))) }
func PortTestSpatialCursorCandidate(u *server.Object, p *types.Pointf) { spatialCursorCandidate(u, p) }
func PortTestSpatialCursor(u *server.Object) uint32                    { return motionAddress(spatialCursor(u)) }
func PortTestSpatialGlobals() (map[string]*uint32, func()) {
	out := map[string]*uint32{"owner": &spatialCursorOwner, "local": (*uint32)(unsafe.Pointer(&C.nox_player_netCode_85319C))}
	for name, off := range map[string]uintptr{"chosen": 2491596, "score": 2491600, "polyp": 2491604} {
		out[name] = memmap.PtrUint32(0x5D4594, off)
	}
	old := map[string]uint32{}
	for n, p := range out {
		old[n] = *p
		*p = 0
	}
	return out, func() {
		for n, p := range out {
			*p = old[n]
		}
	}
}
