//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"unsafe"

	"github.com/opennox/libs/object"
	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

type journalOwner struct {
	*inventoryWindowOwner
	units      []server.Object
	data       []server.PlayerUpdateData
	refs       map[*server.PlayerJournal]uint32
	next       uint32
	queueReset func()
	unchanged  [][]byte
}
type journalNode struct {
	Ref            uint32
	Name           [64]byte
	Next, Prev     uint32
	Flags, Padding uint16
}
type journalResult struct {
	Saved          []byte `json:",omitempty"`
	Op             int
	Return, Height uint32
	Lists          [][]journalNode
	Packets        []legacy.PortTestShopPacketResult
	Render         *objectRenderResult `json:",omitempty"`
}

func newJournalOwner(t *testing.T) *journalOwner {
	o := &journalOwner{inventoryWindowOwner: newInventoryWindowOwner(t)}
	var free func()
	o.units, free = alloc.Make([]server.Object{}, 3)
	t.Cleanup(free)
	o.data, free = alloc.Make([]server.PlayerUpdateData{}, 3)
	t.Cleanup(free)
	var restore func()
	o.queueReset, restore = legacy.PortTestJournalQueue()
	t.Cleanup(restore)
	path := filepath.Join(t.TempDir(), "journal-strings.json")
	if err := o.c.srv.Strings().WriteJSON(path, false); err != nil {
		t.Fatal(err)
	}
	var lang struct {
		Lang    int            `json:"lang"`
		Entries []strman.Entry `json:"entries"`
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(raw, &lang); err != nil {
		t.Fatal(err)
	}
	for _, v := range [][2]string{{"Short", "A short entry."}, {"Wrapped", "This journal entry has enough words to wrap over several lines in the real renderer."}, {"Multiline", "First line\nSecond line\nThird line"}, {"Empty", ""}, {"Unicode", "A tale of café and Ω."}} {
		lang.Entries = append(lang.Entries, strman.Entry{ID: strman.ID("Journal:" + v[0]), Vals: []strman.Variant{{Str: v[1]}}})
	}
	raw, err = json.Marshal(lang)
	if err != nil {
		t.Fatal(err)
	}
	if err = os.WriteFile(path, raw, 0600); err != nil {
		t.Fatal(err)
	}
	if err = o.c.srv.Strings().ReadJSON(path); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { o.release() })
	o.resetJournal(t)
	return o
}
func (o *journalOwner) release() {
	for i := range o.players {
		for n := o.players[i].Journal; n != nil; {
			next := n.Next
			legacy.StrFree(n)
			n = next
		}
		o.players[i].Journal = nil
	}
}
func (o *journalOwner) resetJournal(t *testing.T) {
	o.release()
	o.inventoryWindowOwner.reset(t)
	clear(o.players)
	o.queueReset()
	o.refs = map[*server.PlayerJournal]uint32{}
	o.next = 0
	for i, slot := range []int{1, 7, 31} {
		p := &o.players[slot]
		p.Active = 1
		p.PlayerInd = byte(slot)
		p.PlayerUnit = &o.units[i]
		o.units[i].ObjClass = object.ClassPlayer
		o.units[i].UpdateData = unsafe.Pointer(&o.data[i])
		o.data[i].Player = p
	}
	o.players[3].Active = 1
	o.players[3].PlayerInd = 3 // active slot with no unit
	legacy.Set_dword_8531A0_2576(&o.players[31])
	*memmap.PtrUint32(0x5D4594, 1064848) = 0x12345678
	o.unchanged = nil
	for i := range o.players {
		b := bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(&o.players[i])), int(unsafe.Sizeof(server.Player{}))))
		clear(b[3644:3648])
		o.unchanged = append(o.unchanged, b)
	}
}
func (o *journalOwner) remember(n *server.PlayerJournal) {
	if n == nil {
		return
	}
	o.next++
	o.refs[n] = o.next
}
func (o *journalOwner) ref(t *testing.T, n *server.PlayerJournal) uint32 {
	if n == nil {
		return 0
	}
	v, ok := o.refs[n]
	if !ok {
		t.Fatal("unowned journal node")
	}
	return v
}
func (o *journalOwner) call(t *testing.T, op, slot int, name string, flags uint16) journalResult {
	text, free := alloc.CString(name)
	defer free()
	a := uintptr(unsafe.Pointer(&o.players[slot]))
	if op == 1 || op == 3 || op == 6 || op == 8 {
		a = uintptr(o.players[slot].PlayerUnit.CObj())
	}
	ret := legacy.PortTestJournal(op, a, uintptr(unsafe.Pointer(text)), uintptr(flags))
	if op == 0 || op == 1 {
		o.remember(o.players[slot].Journal)
	}
	if op == 0 || op == 5 || op == 6 && slot == 31 {
		if ret != 0 {
			ret = o.ref(t, (*server.PlayerJournal)(unsafe.Pointer(uintptr(ret))))
		}
	}
	return o.snapshot(t, op, ret, false)
}
func (o *journalOwner) snapshot(t *testing.T, op int, ret uint32, render bool) journalResult {
	r := journalResult{Op: op, Return: ret, Height: memmap.Uint32(0x5D4594, 1064848), Packets: legacy.PortTestJournalPackets()}
	for i := range o.players {
		p := &o.players[i]
		b := bytes.Clone(unsafe.Slice((*byte)(unsafe.Pointer(p)), int(unsafe.Sizeof(server.Player{}))))
		clear(b[3644:3648])
		if !bytes.Equal(b, o.unchanged[i]) {
			t.Fatalf("unrelated player fields changed at slot %d", i)
		}
	}
	for _, slot := range []int{1, 7, 31} {
		var row []journalNode
		var prev *server.PlayerJournal
		seen := map[*server.PlayerJournal]bool{}
		for n := o.players[slot].Journal; n != nil; n = n.Next {
			if seen[n] {
				t.Fatal("journal list cycle")
			}
			seen[n] = true
			if n.Prev != prev {
				t.Fatal("broken journal previous link")
			}
			row = append(row, journalNode{o.ref(t, n), n.EntryBuf, o.ref(t, n.Next), o.ref(t, n.Prev), n.Field3, n.Field4})
			prev = n
		}
		r.Lists = append(r.Lists, row)
	}
	if render {
		v := o.renderResult(t, 0, 0, 0)
		r.Render = &v
	}
	return r
}
func journalCapture(t *testing.T, label string, rows []journalResult, want string) {
	t.Helper()
	b, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	if p := os.Getenv("OPENNOX_JOURNAL_CAPTURE"); p != "" {
		if err = os.WriteFile(p+"-"+label+".json", b, 0600); err != nil {
			t.Fatal(err)
		}
	}
	got := fmt.Sprintf("%x", sha256.Sum256(b))
	t.Logf("%s: %d results %s", label, len(rows), got)
	if want != got {
		t.Fatalf("%s hash %s want %s", label, got, want)
	}
}
