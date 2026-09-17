//go:build porttest

package opennox

import (
	"fmt"
	"github.com/opennox/libs/types"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestObjectReportsInitialPolygon(t *testing.T) {
	s, units, _ := objectReportsPlayers(t)
	t.Cleanup(s.PortTestObjectReportAudio())
	noxServer.ai.Init(noxServer)
	oldSize := videoGetWindowSize()
	videoSetWindowSize(image.Pt(640, 480))
	t.Cleanup(func() { videoSetWindowSize(oldSize) })
	o := objectReportsPolygonOwner(t)
	o.polygon(t, 0, image.Rect(10, 10, 20, 20), 7)
	viewer, emitter := &units[0], &units[1]
	viewer.PosVec = types.Pointf{15, 16}
	emitter.PosVec = viewer.PosVec
	pl := viewer.UpdateDataPlayer().Player
	*(*types.Pointf)(unsafe.Add(unsafe.Pointer(pl), 3632)) = viewer.PosVec
	*(*byte)(unsafe.Add(unsafe.Pointer(emitter.UpdateDataPlayer().Player), 3668)) = 200
	var rows []struct {
		Name   string
		Index  uint32
		Level  byte
		Packet []byte
	}
	defer func() {
		spellbookCapture(t, "object-reports-initial-polygon", rows, "f9688c4d426531a5700933cdde466a8dfea7c6b3f8518a4f9560363e1a202da3")
	}()
	for _, mode := range []uint32{0, 1, 4} {
		for _, follow := range []bool{false, true} {
			name := fmt.Sprintf("mode%d/follow%t", mode, follow)
			t.Run(name, func(t *testing.T) {
				s.NetList.ResetAll()
				s.Audio.Reset()
				pl.Field3680 = mode
				pl.CameraFollowObj = nil
				if follow {
					pl.CameraFollowObj = emitter
				}
				*(*uint32)(unsafe.Add(unsafe.Pointer(pl), 3664)) = 0xdeadface
				*(*byte)(unsafe.Add(unsafe.Pointer(pl), 3668)) = 99
				s.Audio.EventObj(1, emitter, 0, 0)
				legacy.PortTestObjectReports(9, viewer, nil, 0, 0, 0, nil)
				index := *(*uint32)(unsafe.Add(unsafe.Pointer(pl), 3664))
				level := *(*byte)(unsafe.Add(unsafe.Pointer(pl), 3668))
				packet := s.NetList.CopyPacketsA(1, netlist.Kind1)
				if mode&3 != 0 && follow {
					if index != 0xdeadface || level != 99 || len(packet) != 4 {
						t.Fatal("observer must retain viewer polygon")
					}
				} else if index != 1 || level != 7 || len(packet) != 0 {
					t.Fatalf("initial polygon index%x level%d packet%x", index, level, packet)
				}
				rows = append(rows, struct {
					Name   string
					Index  uint32
					Level  byte
					Packet []byte
				}{name, index, level, packet})
			})
		}
	}
}
