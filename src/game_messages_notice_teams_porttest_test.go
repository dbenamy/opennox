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
	"github.com/opennox/opennox/v1/server"
)

func TestGameMessageNoticeTeams(t *testing.T) {
	o := newInventoryTransactionOwner(t)
	o.reset(t)
	printer, index, region := teamRuntimeJoinTextOwner(t, o.c.srv.Server)
	formats := map[string]string{
		"guimsg.c:systemmsg": "System: %s", "team.c:NoTeam": "No team",
		"objcoll.c:FlagCaptureNotice": "Capture %s / %s", "objcoll.c:FlagPickupNotice": "Pickup %s / %s",
		"drop.c:FlagDropNotice": "Drop %s / %s", "update.c:FlagRespawnNotice": "Respawn %s",
		"objcoll.c:FlagBallUnknownNotice": "Unknown / %s", "objcoll.c:FlagBallNotice": "Ball %s / %s",
		"pickup.c:PickUpTeamCrown": "Crown %s / %s", "pickup.c:PickUpCrown": "Crown %s",
		"drop.c:DropTeamCrown": "Lost %s / %s", "drop.c:DropCrown": "Lost %s",
		"pickup.c:WrongTeam": "Wrong %s",
	}
	var entries []strman.Entry
	for k, v := range formats {
		entries = append(entries, strman.Entry{ID: strman.ID(k), Vals: []strman.Variant{{Str: v}}})
	}
	set, restore := o.c.srv.PortTestMeterStrings(entries...)
	defer restore()
	set(0)
	defer o.c.srv.PortTestGameMessageTeamTitles()()
	for i := range o.players {
		o.players[i].Active = 0
	}
	pl := &o.players[0]
	pl.Active = 1
	pl.NetCodeVal = 7
	clear(pl.NameFinal[:])
	copy(pl.NameFinal[:], utf16.Encode([]rune("Ada Ω")))
	o.c.srv.Teams.ByID(1).SetNameAnd68("Ruby Ω", 0)
	o.c.srv.Teams.ByID(2).SetNameAnd68("Azure é", 0)
	o.c.srv.SetFrame(123)
	o.c.srv.SetTickRate(30)
	type row struct {
		Kind         int
		Player, Team uint32
		Return       int
		Text         []string
		Sounds       [][2]int
		Ring         []byte
		Head         uint32
	}
	var rows []row
	for _, kind := range []int{5, 6, 7, 8, 9, 10, 11, 16} {
		for _, id := range []uint32{7, 999} {
			for _, team := range []uint32{0, 1, 2, 255, 257, 0xffffffff} {
				data := bytes.Repeat([]byte{0xa5}, 24)
				data[0], data[1] = 169, byte(kind)
				binary.LittleEndian.PutUint32(data[2:], id)
				binary.LittleEndian.PutUint32(data[6:], team)
				if kind == 8 || kind == 16 {
					binary.LittleEndian.PutUint32(data[2:], team)
				}
				before := bytes.Clone(data)
				clear(region)
				*index = 2
				printer.lines = nil
				o.sounds = nil
				title := "No team"
				name := ""
				switch byte(team) {
				case 1:
					title = "Red Ω"
					name = "Ruby Ω"
				case 2:
					title = "Blue é"
					name = "Azure é"
				}
				text := ""
				sound := 0
				wantN := 10
				switch kind {
				case 5, 6, 7:
					if id == 7 {
						text = fmt.Sprintf("%s Ada Ω / %s", map[int]string{5: "Capture", 6: "Pickup", 7: "Drop"}[kind], title)
						sound = map[int]int{5: 306, 6: 303, 7: 304}[kind]
					}
				case 8:
					text = "Respawn " + title
					sound = 305
					wantN = 6
				case 9:
					if name != "" {
						if id == 7 {
							text = "Ball Ada Ω / " + name
						} else {
							text = "Unknown / " + name
						}
					}
				case 10, 11:
					if id == 7 {
						text = map[int]string{10: "Crown", 11: "Lost"}[kind] + " Ada Ω"
						if name != "" {
							text += " / " + name
						}
					}
				case 16:
					text = "Wrong " + title
					wantN = 6
				}
				got := legacy.PortTestGameNotice(data)
				var wantText []string
				wantRing := make([]byte, len(region))
				wantHead := uint32(2)
				if text != "" {
					wantText = []string{"System: " + text}
					wantHead = 0
					for j, c := range utf16.Encode([]rune(text)) {
						binary.LittleEndian.PutUint16(wantRing[2*j:], c)
					}
					binary.LittleEndian.PutUint32(wantRing[636:], 273)
				}
				var wantSounds [][2]int
				if sound != 0 {
					wantSounds = [][2]int{{sound, 100}}
				}
				if got != wantN || !bytes.Equal(data, before) || !slices.Equal(printer.lines, wantText) || !bytes.Equal(region, wantRing) || *index != wantHead || !slices.Equal(o.sounds, wantSounds) {
					t.Fatalf("team notice kind=%d id=%d team=%d ret=%d text=%q want=%q sounds=%v", kind, id, team, got, printer.lines, wantText, o.sounds)
				}
				if (o.c.srv.Teams.ByID(server.TeamID(team)) != nil) != (name != "") {
					t.Fatal("fixture team ownership")
				}
				rows = append(rows, row{kind, id, team, got, slices.Clone(printer.lines), slices.Clone(o.sounds), bytes.Clone(region), *index})
			}
		}
	}
	gameMessageCapture(t, "game-notice-teams", rows)
}
