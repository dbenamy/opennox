//go:build porttest

package opennox

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/libs/console"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/internal/netlist"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
)

type inventoryTransactionObject struct {
	Ref      uint32
	Live     bool
	Drawable *client.Drawable
}
type inventoryTransactionClient struct {
	*meterClient
	owner *inventoryTransactionOwner
}

func (c *inventoryTransactionClient) Nox_new_drawable_for_thing(typ int) *client.Drawable {
	dr := c.meterClient.Nox_new_drawable_for_thing(typ)
	if dr == nil {
		c.owner.events = append(c.owner.events, [3]uint32{0, uint32(typ), 0})
		return nil
	}
	o := c.owner
	ref := 0xec000001 + uint32(len(o.objects))
	o.objects = append(o.objects, inventoryTransactionObject{ref, true, dr})
	o.identities[dr] = len(o.objects) - 1
	o.c.refs[dr] = ref
	o.unlinked = append(o.unlinked, dr)
	o.events = append(o.events, [3]uint32{1, uint32(typ), ref})
	return dr
}
func (c *inventoryTransactionClient) Nox_xxx_spriteDelete_45A4B0(dr *client.Drawable) int {
	o := c.owner
	i, ok := o.identities[dr]
	if !ok || !o.objects[i].Live {
		panic("transaction deletion outside live case ownership")
	}
	o.events = append(o.events, [3]uint32{2, dr.TypeIDVal, o.objects[i].Ref})
	ret := c.meterClient.Nox_xxx_spriteDelete_45A4B0(dr)
	o.objects[i].Live = false
	return ret
}

type inventoryTransactionPrinter struct{ owner *inventoryTransactionOwner }

func (p inventoryTransactionPrinter) Print(c console.Color, s string) {
	p.owner.console = append(p.owner.console, fmt.Sprintf("%d:%s", c, s))
}
func (p inventoryTransactionPrinter) Printf(c console.Color, s string, a ...interface{}) {
	p.Print(c, fmt.Sprintf(s, a...))
}

type inventoryTransactionOwner struct {
	*uiInventoryOwner
	tx         *inventoryTransactionClient
	txwords    []*uint32
	txregions  [][]byte
	objects    []inventoryTransactionObject
	identities map[*client.Drawable]int
	events     [][3]uint32
	console    []string
	poolStart  [3]int
	countStart int
}

func newInventoryTransactionOwner(t *testing.T) *inventoryTransactionOwner {
	o := &inventoryTransactionOwner{uiInventoryOwner: newUIInventoryOwner(t)}
	o.tx = &inventoryTransactionClient{legacy.GetClient().(*meterClient), o}
	oldClient := legacy.GetClient
	legacy.GetClient = func() legacy.Client { return o.tx }
	t.Cleanup(func() { legacy.GetClient = oldClient })
	var restore func()
	o.txwords, restore = legacy.PortTestInventoryTransactionWords()
	t.Cleanup(restore)
	for _, r := range [][3]uintptr{{0x852978, 8, 4}, {0x5D4594, 1047764, 144}, {0x5D4594, 823804, 1932}, {0x5D4594, 1049724, 8}} {
		b := unsafe.Slice((*byte)(memmap.PtrOff(r[0], r[1])), r[2])
		old := append([]byte(nil), b...)
		o.txregions = append(o.txregions, b)
		t.Cleanup(func() { copy(b, old) })
	}
	entries := []strman.Entry{{ID: "guimsg.c:systemmsg", Vals: []strman.Variant{{Str: "System: %s"}}}}
	for _, name := range []string{"InventoryFull", "DrawablesExhausted", "EquippedNotFound", "TooManyEquipped", "DroppedNotFound"} {
		entries = append(entries, strman.Entry{ID: strman.ID("guiinv.c:" + name), Vals: []strman.Variant{{Str: name}}})
	}
	configure, restore := o.c.srv.Server.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
	oldStrings := strMan
	strMan = o.c.srv.Strings()
	t.Cleanup(func() { strMan = oldStrings })
	oldConsole := legacy.GetConsole
	con := console.New(inventoryTransactionPrinter{o})
	legacy.GetConsole = func() *console.Console { return con }
	t.Cleanup(func() { legacy.GetConsole = oldConsole })
	t.Cleanup(func() { o.clearObjects(t) })
	return o
}
func (o *inventoryTransactionOwner) clearObjects(t *testing.T) {
	for _, obj := range o.objects {
		if obj.Live {
			o.tx.Nox_xxx_spriteDelete_45A4B0(obj.Drawable)
		}
	}
	if o.identities != nil {
		if got := o.c.Objs.Alloc.PortTestCounts(); got != o.poolStart {
			t.Errorf("pool restoration %v want %v", got, o.poolStart)
		}
		if got := int(o.c.Objs.Count); got != o.countStart {
			t.Errorf("drawable count restoration %d want %d", got, o.countStart)
		}
	}
	for dr := range o.identities {
		delete(o.c.refs, dr)
	}
	o.objects = nil
	o.identities = nil
}
func (o *inventoryTransactionOwner) reset(t *testing.T) {
	t.Helper()
	o.clearObjects(t)
	o.uiInventoryOwner.reset(t)
	for _, p := range o.txwords {
		*p = 0
	}
	for _, b := range o.txregions {
		clear(b)
	}
	o.identities = make(map[*client.Drawable]int)
	o.events = nil
	o.console = nil
	o.poolStart = o.c.Objs.Alloc.PortTestCounts()
	o.countStart = int(o.c.Objs.Count)
	noxflags.ResetGame()
	noxflags.SetGame(1)
	o.tx.Nox_client_setCursorType(0)
	o.c.CursorPrev = 0
	o.c.dragndropItem = nil
}
func (o *inventoryTransactionOwner) item(t *testing.T, name string, code uint32) *client.Drawable {
	t.Helper()
	dr := o.tx.Nox_new_drawable_for_thing(o.c.Things.TypeByID(name).Index())
	if dr == nil {
		t.Fatal("case drawable allocation")
	}
	dr.NetCode32 = code
	return dr
}
func txptr(p unsafe.Pointer) uintptr                  { return uintptr(p) }
func txword(dr *client.Drawable, off uintptr) *uint32 { return (*uint32)(unsafe.Add(dr.C(), off)) }
func (o *inventoryTransactionOwner) norm(v uint32) uint32 {
	if i, ok := o.identities[(*client.Drawable)(unsafe.Pointer(uintptr(v)))]; ok {
		return o.objects[i].Ref
	}
	if n, ok := o.modRefs[v]; ok {
		return n
	}
	return o.uiInventoryOwner.norm(v)
}
func (o *inventoryTransactionOwner) call(op int, a, b, c, d uintptr) uint32 {
	return legacy.PortTestInventoryTransaction(op, a, b, c, d)
}
func (o *inventoryTransactionOwner) exhausted(t *testing.T, f func()) {
	t.Helper()
	oldPool, oldList := o.c.Objs.Alloc, o.c.Objs.List1
	pool := alloc.NewClassT("inventory exhaustion", client.Drawable{}, 1)
	if pool.NewObject() == nil {
		t.Fatal("exhaustion filler")
	}
	o.c.Objs.Alloc, o.c.Objs.List1 = pool, nil
	defer func() { o.c.Objs.Alloc, o.c.Objs.List1 = oldPool, oldList; pool.Free() }()
	f()
}

type inventoryTransactionResult struct {
	Sounds                 [][2]int
	Cursor                 [3]uint32
	Potions                [3][5]uint32
	MeterWindows           [2]uint32
	Case, Op               int
	Return                 uint32
	Named, Grid, Equipment []uint32
	Regions                [][]byte
	Objects                [][]uint32
	Events                 [][3]uint32
	Messages               [][]byte
	Console                []string
	Meters                 [2][2]uint32
	Pool                   [3]int
	Count                  int
}

func (o *inventoryTransactionOwner) snapshot(t *testing.T, id, op int, ret uint32) inventoryTransactionResult {
	t.Helper()
	r := inventoryTransactionResult{Case: id, Op: op, Return: o.norm(ret), Events: append([][3]uint32(nil), o.events...), Console: append([]string(nil), o.console...)}
	for _, ps := range [][]*uint32{o.words, o.txwords} {
		for _, p := range ps {
			r.Named = append(r.Named, o.norm(*p))
		}
	}
	r.Grid = append([]uint32(nil), unsafe.Slice((*uint32)(unsafe.Pointer(&o.grid[0])), len(o.grid)/4)...)
	for i := 0; i < 84; i++ {
		r.Grid[i*37] = o.norm(r.Grid[i*37])
	}
	for _, v := range o.equipment {
		r.Equipment = append(r.Equipment, o.norm(v))
	}
	for _, b := range append(append([][]byte(nil), o.regions...), o.txregions...) {
		cp := append([]byte(nil), b...)
		// Only declared pointer fields are normalized; UTF16 messages are bytes.
		if len(b) == 4 {
			binary.LittleEndian.PutUint32(cp, o.norm(binary.LittleEndian.Uint32(cp)))
		}
		if len(b) == 8 && unsafe.Pointer(&b[0]) == memmap.PtrOff(0x5D4594, 1049788) {
			binary.LittleEndian.PutUint32(cp, o.norm(binary.LittleEndian.Uint32(cp)))
		}
		r.Regions = append(r.Regions, cp)
	}
	for _, obj := range o.objects {
		row := []uint32{obj.Ref, 0}
		if obj.Live {
			row[1] = 1
			for _, off := range []uintptr{108, 112, 116, 120, 128, 292, 368, 372, 432, 436, 440, 444, 448, 452} {
				v := *txword(obj.Drawable, off)
				if off == 368 || off == 372 || off >= 432 && off < 448 {
					v = o.norm(v)
				}
				row = append(row, v)
			}
		}
		r.Objects = append(r.Objects, row)
	}
	o.c.srv.NetList.ByInd(31, netlist.Kind0).Each(func(b []byte) bool { r.Messages = append(r.Messages, append([]byte(nil), b...)); return false })
	for i := 0; i < 2; i++ {
		m := o.meters.Records[5+i]
		r.Meters[i] = [2]uint32{m.Current, m.Maximum}
	}
	r.Sounds = append([][2]int(nil), o.sounds...)
	r.Cursor = [3]uint32{uint32(o.c.Cursor), uint32(o.c.CursorPrev), o.norm(uint32(uintptr(unsafe.Pointer(o.c.dragndropItem))))}
	for i := 0; i < 3; i++ {
		off := uintptr(1090296 + i*536)
		for j := 0; j < 5; j++ {
			r.Potions[i][j] = memmap.Uint32(0x5D4594, off+uintptr(j*4))
		}
		r.Potions[i][0] = o.norm(r.Potions[i][0])
	}
	for i := 0; i < 2; i++ {
		r.MeterWindows[i] = uint32(o.meters.Records[5+i].Window.GetFlags())
	}
	r.Pool = o.c.Objs.Alloc.PortTestCounts()
	for i := range r.Pool {
		r.Pool[i] -= o.poolStart[i]
	}
	r.Count = int(o.c.Objs.Count) - o.countStart
	live := 0
	for _, obj := range o.objects {
		if obj.Live {
			live++
		}
	}
	if r.Pool[0] != live || r.Count != live {
		t.Fatalf("case %d live=%d pool=%v count=%d", id, live, r.Pool, r.Count)
	}
	return r
}
func inventoryTransactionCapture(t *testing.T, label string, rows []inventoryTransactionResult, want string) {
	t.Helper()
	data, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if p := os.Getenv("OPENNOX_INVENTORY_TRANSACTIONS_CAPTURE"); p != "" {
		if err := os.WriteFile(p+"-"+label+".json", data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(data))
	t.Logf("%s: %d results %s", label, len(rows), got)
	if got != want {
		t.Fatalf("%s hash %s want frozen C %s", label, got, want)
	}
}
