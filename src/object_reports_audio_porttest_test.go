//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
	"image"
	"testing"
	"unsafe"
)

func TestObjectReportsObserverAudio(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	t.Cleanup(s.PortTestObjectReportAudio())
	noxServer.ai.Init(noxServer)
	oldSize := videoGetWindowSize()
	videoSetWindowSize(image.Pt(640, 480))
	t.Cleanup(func() { videoSetWindowSize(oldSize) })
	o := objectReportsPolygonOwner(t)
	o.polygon(t, 0, image.Rect(10, 10, 20, 20), 7)
	camera, free := alloc.New(server.Object{})
	t.Cleanup(free)
	viewer, emitter, follow := &units[0], &units[1], &units[2]
	emitter.PosVec = types.Pointf{15, 16}
	*(*byte)(unsafe.Add(unsafe.Pointer(emitter.UpdateDataPlayer().Player), 3668)) = 7
	*(*byte)(unsafe.Add(unsafe.Pointer(follow.UpdateDataPlayer().Player), 3668)) = 200
	pl := viewer.UpdateDataPlayer().Player
	*(*types.Pointf)(unsafe.Add(unsafe.Pointer(pl), 3632)) = types.Pointf{15, 16}
	var rows []struct {
		Name    string
		Packets []byte
	}
	defer func() {
		spellbookCapture(t, "object-reports-observer-audio", rows, "c24c4dabf2e9263649f481637dda56c05f2a971ea4e43bc519a377e61cbe13ed")
	}()
	for _, mode := range []uint32{0, 1, 2, 3, 4, 255} {
		for _, cached := range []byte{0, 7, 200, 255} {
			for _, target := range []string{"nil", "player", "polygon", "outside"} {
				name := fmt.Sprintf("mode%d/cache%d/%s", mode, cached, target)
				t.Run(name, func(t *testing.T) {
					s.NetList.ResetAll()
					s.Audio.Reset()
					pl.Field3680 = mode
					pl.CameraFollowObj = nil
					*(*byte)(unsafe.Add(unsafe.Pointer(pl), 3668)) = cached
					level := cached
					switch target {
					case "player":
						pl.CameraFollowObj = follow
						if mode&3 != 0 {
							level = 200
						}
					case "polygon":
						camera.PosVec = types.Pointf{15, 16}
						pl.CameraFollowObj = camera
						if mode&3 != 0 {
							level = 7
						}
					case "outside":
						camera.PosVec = types.Pointf{60, 60}
						pl.CameraFollowObj = camera
						if mode&3 != 0 {
							level = 0
						}
					}
					s.Audio.EventObj(1, emitter, 0, 0)
					s.Audio.EventPos(2, types.Pointf{60, 60}, 0, 0)
					legacy.PortTestObjectReports(9, viewer, nil, 0, 0, 0, nil)
					got := s.NetList.CopyPacketsA(1, netlist.Kind1)
					wantLen := 4
					if level == 7 {
						wantLen = 8
					}
					if len(got) != wantLen {
						t.Fatalf("level%d packets%x want%d bytes", level, got, wantLen)
					}
					if s.NetList.ByInd(7, netlist.Kind1).Count() != 0 || s.NetList.ByInd(31, netlist.Kind1).Count() != 0 {
						t.Fatal("wrong audio recipient")
					}
					rows = append(rows, struct {
						Name    string
						Packets []byte
					}{name, got})
				})
			}
		}
	}
}
