//go:build porttest

package opennox

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"slices"
	"testing"
	"unicode/utf16"

	"github.com/opennox/libs/strman"
	"github.com/opennox/opennox/v1/legacy"
)

func TestGameMessageNoticeText(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	printer, index, region := teamRuntimeJoinTextOwner(t, o.c.srv.Server)
	formats := map[string]string{
		"guimsg.c:systemmsg":            "System: %s",
		"netserv.c:PlayerTimeout":       "Timeout %s",
		"objcoll.c:FlagRetrieveNotice":  "Retrieve %s",
		"Netserv.c:InObservationMode":   "Observing",
		"Netserv.c:PressJump":           "Press jump",
		"atckexec.c:PlayerStunned":      "Stunned",
		"atckexec.c:PlayerConfused":     "Confused",
		"atckexec.c:PlayerPoisoned":     "Poisoned",
		"player.c:TooHeavy":             "Too heavy",
		"objcoll.c:PlayerExited":        "Exited %s",
		"objcoll.c:PlayerExitedWarp":    "Warp %s",
		"GeneralPrint:SecretFoundOther": "Secret %s",
		"GeneralPrint:AdvanceToStage1":  "Stage %d",
		"use.c:SignSays":                "Sign [%s]",
		"inform.c:NoticeFixture":        "Café Ω",
	}
	var entries []strman.Entry
	for k, v := range formats {
		entries = append(entries, strman.Entry{ID: strman.ID(k), Vals: []strman.Variant{{Str: v}}})
	}
	set, restore := o.c.srv.PortTestMeterStrings(entries...)
	defer restore()
	set(0)
	oldPlay := legacy.GameGetPlayState
	defer func() { legacy.GameGetPlayState = oldPlay }()
	for i := range o.players {
		o.players[i].Active = 0
	}
	pl := &o.players[0]
	pl.Active = 1
	pl.NetCodeVal = 0x8007
	copy(pl.NameFinal[:], utf16.Encode([]rune("Ada Ω")))
	o.c.srv.SetFrame(123)
	o.c.srv.SetTickRate(30)
	type row struct {
		Kind   int
		Value  uint32
		State  int
		Return int
		Text   []string
		Sounds [][2]int
		Ring   []byte
		Head   uint32
	}
	var rows []row
	run := func(kind int, value uint32, state int, text []string, sound int, wantReturn int, sign string) {
		t.Helper()
		data := bytes.Repeat([]byte{0xa5}, 80)
		data[0] = 0xa9
		data[1] = byte(kind)
		binary.LittleEndian.PutUint32(data[2:], value)
		if kind == 15 {
			copy(data[3:], sign)
			data[3+len(sign)] = 0
		}
		before := bytes.Clone(data)
		clear(region)
		*index = 2
		printer.lines = nil
		o.sounds = nil
		legacy.GameGetPlayState = func() int { return state }
		got := legacy.PortTestGameNotice(data)
		if got != wantReturn || !bytes.Equal(before, data) {
			t.Fatalf("notice kind=%d value=%x return/input %d want %d", kind, value, got, wantReturn)
		}
		wantConsole := make([]string, len(text))
		wantRing := make([]byte, len(region))
		for i, s := range text {
			wantConsole[i] = "System: " + s
			for j, c := range utf16.Encode([]rune(s)) {
				binary.LittleEndian.PutUint16(wantRing[i*644+j*2:], c)
			}
			binary.LittleEndian.PutUint32(wantRing[i*644+636:], 273)
		}
		wantHead := uint32((2 + len(text)) % 3)
		if !slices.Equal(printer.lines, wantConsole) || !bytes.Equal(region, wantRing) || *index != wantHead {
			t.Fatalf("notice kind=%d value=%x text=%q want=%q head=%d/%d", kind, value, printer.lines, wantConsole, *index, wantHead)
		}
		var sounds [][2]int
		if sound != 0 {
			sounds = [][2]int{{sound, 100}}
		}
		if !slices.Equal(o.sounds, sounds) {
			t.Fatalf("notice %d sounds %v want %v", kind, o.sounds, sounds)
		}
		rows = append(rows, row{kind, value, state, got, slices.Clone(printer.lines), slices.Clone(o.sounds), bytes.Clone(region), *index})
	}
	for _, kind := range []int{3, 4, 18, 19, 20} {
		prefix := map[int]string{3: "Timeout ", 4: "Retrieve ", 18: "Exited ", 19: "Warp ", 20: "Secret "}[kind]
		for _, id := range []uint32{0, 0x8007, 0x10008007, 0xffffffff} {
			var text []string
			sound := 0
			if id == 0x8007 {
				text = []string{prefix + "Ada Ω"}
				if kind == 4 {
					sound = 305
				}
			}
			run(kind, id, 3, text, sound, 6, "")
		}
	}
	for _, value := range []uint32{0, 1, 2, 3, 4, 0x80000000, 0xffffffff} {
		text := []string{"Observing"}
		if value != 0 {
			text = append(text, "Press jump")
		}
		run(12, value, 3, text, 0, 6, "")
		var status []string
		if value < 4 {
			status = []string{[]string{"Stunned", "Confused", "Poisoned", "Too heavy"}[value]}
		}
		run(13, value, 3, status, 0, 6, "")
		run(21, value, 3, []string{fmt.Sprintf("Sign [Stage %d]", int32(value))}, 0, 6, "")
	}
	for _, state := range []int{0, 2, 3, 4} {
		for _, flag := range []uint32{0, 1, 255} {
			var text []string
			if state == 3 {
				s := "Café Ω"
				if flag != 0 {
					s = "Sign [" + s + "]"
				}
				text = []string{s}
			}
			run(15, flag, state, text, 0, 3+len("NoticeFixture")+1, "NoticeFixture")
		}
	}
	for kind := 22; kind < 256; kind++ {
		run(kind, 0xffffffff, 3, nil, 0, 0, "")
	}
	gameMessageCapture(t, "game-notice-text", rows)
}
