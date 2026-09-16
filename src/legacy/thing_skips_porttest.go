//go:build porttest

package legacy

import "github.com/opennox/opennox/v1/internal/binfile"

func PortTestThingSkip(op int, f *binfile.MemFile, scratch []byte) int {
	switch op {
	case 0:
		return thingSkipAUD(f)
	case 1:
		return thingSkipSpells(f)
	case 2:
		return thingSkipAbilities(f)
	case 3:
		return thingSkipImages(f)
	case 4:
		return thingSkipAVNT(f)
	case 5:
		return thingSkipAVNTInner(f)
	case 6:
		return thingSkipWall(f, scratch)
	default:
		panic(op)
	}
}
