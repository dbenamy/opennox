//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/legacy"
	"image"
	"testing"
	"unsafe"
)

func TestColorLightViewportUpdate(t *testing.T) {
	c, pix, effects := newEffectsFullOwner(t)
	env := legacy.PortTestNewClientParticleEnvironment()
	t.Cleanup(env.Restore)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	vp := c.Viewport()
	old := *vp
	t.Cleanup(func() { *vp = old })
	vp.World = image.Rect(200, 300, 264, 364)
	vp.Size = image.Pt(64, 64)
	type record struct {
		X, Y, Count int
		Suppressed  bool
		Return      int
		Linked      bool
		Light       []uint32
	}
	var records []record
	for _, x := range []int{99, 100, 200, 364, 365} {
		for _, y := range []int{199, 200, 300, 464, 465} {
			for _, count := range []int{0, 1, 2} {
				for _, suppressed := range []bool{false, true} {
					t.Run(fmt.Sprintf("x%d-y%d-count%d-suppressed%v", x, y, count, suppressed), func(t *testing.T) {
						c.resetCase(effects, pix, 17, 0)
						env.Reset()
						flags := uint32(0)
						if suppressed {
							flags = 0x200000
						}
						restore := noxflags.PortTestGameFlags(noxflags.GameFlag(flags))
						defer restore()
						d := c.Nox_xxx_spriteLoadAdd_45A360_drawable(4, image.Pt(x, y))
						raw := unsafe.Slice((*byte)(d.C()), int(unsafe.Sizeof(*d)))
						clear(raw[136:276])
						clear(raw[432:435])
						raw[432] = byte(count)
						raw[178] = 7
						raw[179] = 13
						raw[180] = 19
						binary.LittleEndian.PutUint16(raw[258:], 1)
						c.Objs.List5Add(d)
						got := legacy.PortTestColorLight(5, vp, d)
						if got != 1 || (d.InClientUpdateList != 0) != (count != 0) {
							t.Fatal("update return/list membership")
						}
						active := count > 1 && !suppressed && x >= 100 && x <= 364 && y >= 200 && y <= 464
						light := make([]uint32, 10)
						for i := range light {
							light[i] = binary.LittleEndian.Uint32(raw[136+4*i:])
						}
						if active {
							if light[0] != 2 || light[4] != 7 || light[5] != 13 || light[6] != 19 {
								t.Fatal("visible light did not update")
							}
						} else {
							for _, v := range light {
								if v != 0 {
									t.Fatal("culled or inactive light changed")
								}
							}
						}
						records = append(records, record{x, y, count, suppressed, got, d.InClientUpdateList != 0, light})
					})
				}
			}
		}
	}
	spellbookCapture(t, "color-light-update", records, "973af366685512326ddd2f3699d1d8964e4f90767a4329da0d1c7bf965707fb0")
}
