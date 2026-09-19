//go:build porttest

package opennox

import (
	"encoding/binary"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
	"unsafe"
)

type chatBubbleStorage struct{ head, tail *unsafe.Pointer }
type chatBubbleOwner struct {
	*matchRosterOwner
	*chatBubbleStorage
}

func newChatBubbleStorage(t *testing.T) *chatBubbleStorage {
	t.Helper()
	o := new(chatBubbleStorage)
	old := legacy.Get_nox_alloc_chat_1197364()
	pool := alloc.NewClass("Chat", 692, 64)
	legacy.Set_nox_alloc_chat_1197364(pool.UPtr())
	t.Cleanup(func() { pool.Free(); legacy.Set_nox_alloc_chat_1197364(old) })
	var restore func()
	o.head, o.tail, restore = legacy.PortTestChatBubbleGlobals()
	t.Cleanup(restore)
	return o
}
func newChatBubbleOwner(t *testing.T) *chatBubbleOwner {
	t.Helper()
	o := &chatBubbleOwner{matchRosterOwner: newMatchRosterOwner(t), chatBubbleStorage: newChatBubbleStorage(t)}
	o.reset()
	return o
}
func (o *chatBubbleStorage) create(t *testing.T, code uint16, text string, length byte, pos uint32, duration uint16) unsafe.Pointer {
	t.Helper()
	var b [11]byte
	binary.LittleEndian.PutUint16(b[1:], code)
	binary.LittleEndian.PutUint32(b[4:], pos)
	b[8] = length
	binary.LittleEndian.PutUint16(b[9:], duration)
	str, free := alloc.CString16(text)
	defer free()
	legacy.PortTestChatBubbleCreate(unsafe.Pointer(&b[0]), str)
	p := legacy.PortTestChatBubbleLookup(uint32(code))
	if p == nil {
		t.Fatal("created bubble missing", code)
	}
	return p
}

type chatBubbleRecord struct {
	Code                uint32
	Text                string
	Length, Expiry, Pos uint32
	Previous, Next      uint32
}

func (o *chatBubbleStorage) snapshot(t *testing.T) []chatBubbleRecord {
	t.Helper()
	var out []chatBubbleRecord
	var prev unsafe.Pointer
	for p := *o.head; p != nil; p = *(*unsafe.Pointer)(unsafe.Add(p, 684)) {
		if len(out) >= 64 {
			t.Fatal("bubble list cycle")
		}
		if *(*unsafe.Pointer)(unsafe.Add(p, 688)) != prev {
			t.Fatal("bubble previous link")
		}
		r := chatBubbleRecord{Code: objectXferGetWord(p, 656), Text: alloc.GoString16((*uint16)(p)), Length: objectXferGetWord(p, 636), Expiry: objectXferGetWord(p, 640), Pos: objectXferGetWord(p, 644)}
		if prev != nil {
			r.Previous = objectXferGetWord(prev, 656)
		}
		if next := *(*unsafe.Pointer)(unsafe.Add(p, 684)); next != nil {
			r.Next = objectXferGetWord(next, 656)
		}
		out = append(out, r)
		prev = p
	}
	if *o.head != nil && *o.tail != prev {
		t.Fatal("bubble tail does not identify last live entry")
	}
	return out
}
