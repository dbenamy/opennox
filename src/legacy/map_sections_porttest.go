//go:build porttest

package legacy

import (
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/server"
	"os"
	"path/filepath"
	"unsafe"
)

type PortTestMapSectionIO struct {
	Function       string
	Read, Previous bool
	Data           []byte
}
type PortTestMapSectionSpec struct {
	WallSprites  bool
	ScratchWalls []PortTestPaintWall
	Name         string
	Paint        PortTestPaintSpec
	Flags        uint32
	ChainSlot    int
	Chain        [][4]uint32
	IO           []PortTestMapSectionIO
}
type PortTestMapSectionWire struct {
	Breakable []uint32
	Scratch   [][9]uint32
	Tile      [][4]uint32
	Data      []byte
	Position  int64
	Checksum  uint32
	Return    uint32
}
type PortTestMapSectionResult struct {
	Name  string
	Paint PortTestPaintResult
	IO    []PortTestMapSectionWire
}

// Reuse the real grid, subtile pool and wall owners. The only new adapter work is
// opening an owned file and invoking the selected C section entrypoint.
func PortTestMapSections(cases []PortTestMapSectionSpec, dir string, owner func(*server.Server) (Server, func())) []PortTestMapSectionResult {
	ext := &paintTestExtension{globals: map[string]*uint32{
		"section-max-x": &mapSectionMaxX,
		"section-max-y": &mapSectionMaxY,
	}}
	for name, off := range map[string]uintptr{"map-min-x": 739980, "map-min-y": 739984, "map-width": 739988, "wall-load-flags": 739992, "window-count": 741336, "breakable-count": 741340, "breakable-index": 741344, "secret-count": 741348, "secret-index": 741352, "min-x": 741360, "min-y": 741368, "magic-wall": 741372} {
		ext.globals["section-"+name] = memmap.PtrUint32(0x5D4594, off)
	}
	prefabWords, restorePrefab := PortTestPrefabRuntimeGlobals()
	defer restorePrefab()
	for name, p := range prefabWords {
		ext.globals["section-prefab-"+name] = p
	}
	ext.globals["section-prefab-tile-count"] = memmap.PtrUint32(0x5D4594, 1599560)
	var active *paintTestFixture
	ext.setup = func(*server.PortTestPaintOwners) func() {
		oldAlloc, oldFree := themeTestOnAllocate, themeTestOnRelease
		themeTestOnAllocate = func(p unsafe.Pointer, n int) {
			if active == nil || n == 200 {
				return
			} // The existing subtile-pool discovery owns these blocks.
			switch n {
			case 12, 20, 24, 32, 36:
				kind := "section-allocation"
				if n == 20 {
					kind = "input"
				}
				r := active.register(p, n, kind, false)
				active.owned[r] = true
			}
		}
		themeTestOnRelease = func(p unsafe.Pointer) {
			if active != nil {
				if r := active.known(p); r != nil {
					r.alive = false
					delete(active.owned, r)
				}
			}
		}
		return func() { themeObserve(false, 0); themeTestOnAllocate, themeTestOnRelease = oldAlloc, oldFree }
	}
	ext.finish = func(f *paintTestFixture) { themeObserve(false, 0); active = nil; f.owners.S.Walls.ClearBreakable() }
	savedFile := cryptfile.Global()
	cryptfile.SetGlobal(nil)
	defer func() { cryptfile.Close(); cryptfile.SetGlobal(savedFile) }()
	current := -1
	out := make([]PortTestMapSectionResult, len(cases))
	specs := make([]PortTestPaintSpec, len(cases))
	for i, c := range cases {
		specs[i] = c.Paint
		out[i].Name = c.Name
	}
	ext.reset = func(PortTestPaintSpec) { active = nil; current++ }
	ext.before = func(f *paintTestFixture, _ PortTestPaintSpec) {
		active = f
		c := cases[current]
		noxflags.ResetGame()
		noxflags.SetGame(noxflags.GameFlag(c.Flags))
		themeObserve(true, 0)
		for _, spec := range c.ScratchWalls {
			node := prefabWallNew(byte(spec.X), byte(spec.Y))
			wall := mapRoomPointer(*prefabWord(node, 0))
			for off, v := range spec.Words {
				*(*uint32)(unsafe.Add(wall, off)) = v
			}
			if spec.Data.Slot != 0 || spec.Data.Value != 0 {
				*(*uint32)(unsafe.Add(wall, 28)) = f.resolve(spec.Data)
			}
		}
		themeObserve(false, 0)
		if len(c.Chain) != 0 {
			root := (*[5]uint32)(f.slots[c.ChainSlot].ptr)
			for _, value := range c.Chain {
				node := mapPaintSubtileNew(int32(value[0]), int32(value[1]), int32(value[2]), int32(value[3]))
				root[4] = mapRoomRaw(unsafe.Pointer(node))
				root = node
			}
			f.discover(0, -1)
		}
	}
	ext.invoke = func(f *paintTestFixture, op int, args [6]uint32) uint32 {
		spec := cases[current].IO[op]
		path := filepath.Join(dir, "section.bin")
		if spec.Read {
			data := spec.Data
			if spec.Previous {
				data = out[current].IO[len(out[current].IO)-1].Data
			}
			if err := os.WriteFile(path, data, 0600); err != nil {
				panic(err)
			}
		}
		mode := cryptfile.WriteOnly
		if spec.Read {
			mode = cryptfile.ReadOnly
		}
		if err := cryptfile.OpenGlobal(path, mode, -1); err != nil {
			panic(err)
		}
		// Version 2 selects a variation only if its real wall definition has a
		// sprite handle. The section checks presence without dereferencing it.
		// Use owned backing and restore before the definition snapshot.
		restoreSprites := func() {}
		if cases[current].WallSprites {
			var old [3][4][15][16]unsafe.Pointer
			for i := range old {
				d := f.owners.S.Walls.DefByInd(i)
				old[i] = d.Sprite8432
				for a := range d.Sprite8432 {
					for b := range d.Sprite8432[a] {
						for c := range d.Sprite8432[a][b] {
							d.Sprite8432[a][b][c] = f.rows[0]
						}
					}
				}
			}
			restoreSprites = func() {
				for i := range old {
					f.owners.S.Walls.DefByInd(i).Sprite8432 = old[i]
				}
			}
		}
		defer restoreSprites()
		var ret uint32
		themeObserve(true, 0)
		defer themeObserve(false, 0)
		switch spec.Function {
		case "tile":
			ret = mapSectionTile((*[5]uint32)(mapRoomPointer(args[1])))
		case "floor":
			ret = mapSectionFloor((*[8]uint32)(mapRoomPointer(args[0])))
		case "walls":
			ret = mapSectionWalls((*[8]uint32)(mapRoomPointer(args[0])))
		case "windows":
			ret = mapSectionMetadata(0, (*[8]uint32)(mapRoomPointer(args[0])))
		case "breakable":
			ret = mapSectionMetadata(1, (*[8]uint32)(mapRoomPointer(args[0])))
		case "secret":
			ret = mapSectionMetadata(2, (*[8]uint32)(mapRoomPointer(args[0])))
		default:
			panic(spec.Function)
		}
		themeObserve(false, 0)
		restoreSprites()
		pos, err := cryptfile.Global().File.Seek(0, 1)
		if err != nil {
			panic(err)
		}
		sum := cryptfile.Global().PortTestChecksum()
		if err := cryptfile.Close(); err != nil {
			panic(err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			panic(err)
		}
		row := PortTestMapSectionWire{Data: data, Position: pos, Checksum: sum, Return: ret}
		for p := f.owners.S.Walls.FirstBreakable(); p != nil; p = p.Next() {
			row.Breakable = append(row.Breakable, f.norm(mapRoomRaw(unsafe.Pointer(p.Wall))))
		}
		if spec.Function == "tile" {
			for p := mapRoomPointer(args[1]); p != nil; p = *controlPtr(p, 16) {
				if len(row.Tile) > 512 {
					panic("unexpected tile chain cycle")
				}
				row.Tile = append(row.Tile, *(*[4]uint32)(p))
			}
		}
		for node := *prefabGlobal(prefabWalls); node != 0; node = *prefabWord(node, 4) {
			data := *(*[9]uint32)(mapRoomPointer(*prefabWord(node, 0)))
			for i, v := range data {
				data[i] = f.norm(v)
			}
			row.Scratch = append(row.Scratch, data)
		}
		out[current].IO = append(out[current].IO, row)
		return ret
	}
	states := portTestMapPainting(specs, owner, ext)
	for i, state := range states {
		out[i].Paint = state
	}
	return out
}
