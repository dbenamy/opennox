//go:build porttest

package opennox

import (
	"encoding/binary"
	"fmt"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"testing"
)

func TestConsoleCommandsNumbers(t *testing.T) {
	type row struct {
		Command, Input string
		Result         bool
		Value          uint64
		Updated        uint32
		Output         []consoleCommandLine
	}
	var rows []row
	inputs := []struct {
		text  string
		value int32
	}{{"", 0}, {"junk", 0}, {"１２", 0}, {"\u00a05", 0}, {"12Ω", 12}, {" 12suffix", 12}, {"+13", 13}, {"-1", -1}, {"0x10", 0}, {"255", 255}, {"256", 256}, {"65535", 65535}, {"65536", 65536}, {"999", 999}, {"1000", 1000}, {"2147483647", 2147483647}, {"2147483648", 2147483647}, {"-2147483649", -2147483648}}
	for _, command := range []string{"players", "time", "lessons"} {
		for _, in := range inputs {
			t.Run(command+"/"+in.text, func(t *testing.T) {
				o := newConsoleCommandOwner(t)
				binary.LittleEndian.PutUint16(o.settings[52:], 0x100)
				*o.optionWords["settings-updated"] = 0
				got := o.call(t, "set "+command, false, in.text)
				var value, want uint64
				switch command {
				case "players":
					value = legacy.PortTestServerConfigScalar("limit-get", 0, 0)
					v := in.value
					if v < 0 {
						v = 0
					}
					if v > 999 {
						v = 999
					}
					want = uint64(v)
				case "time":
					value = legacy.PortTestServerConfigScalar("minutes", 0x100, 0)
					want = uint64(uint8(in.value))
				case "lessons":
					value = legacy.PortTestServerConfigScalar("score", 0x100, 0)
					v := uint16(in.value)
					if v > 999 {
						v = 999
					}
					want = uint64(v)
				}
				if !got || value != want {
					t.Fatalf("result %t value %d want %d", got, value, want)
				}
				rows = append(rows, row{command, in.text, got, value, *o.optionWords["settings-updated"], o.printer.lines})
			})
		}
	}
	spellbookCapture(t, "console-commands-numbers", rows, "b00aa6cb69f5532331972152ec99224081762eead747c19046d219f335ce55fd")
}
func TestConsoleCommandsNames(t *testing.T) {
	type row struct {
		Command string
		Args    []string
		Result  bool
		Stored  string
		Bytes   []byte
		Output  []consoleCommandLine
	}
	var rows []row
	for _, tc := range []struct {
		command string
		args    []string
		want    string
	}{
		{"name", []string{"Alpha"}, "Alpha"}, {"name", []string{"Ω界"}, "\xa9L"}, {"name", []string{"Alpha\x00tail"}, "Alpha"}, {"name", []string{"Alpha", "Beta"}, "Alpha Beta"},
		{"name", []string{"12345678901234567890"}, "123456789012345"},
		{"sysop", []string{"pass"}, "pass"}, {"sysop", []string{"Ω界"}, "Ω界"}, {"sysop", []string{""}, ""},
	} {
		t.Run(fmt.Sprint(tc.command, tc.args), func(t *testing.T) {
			o := newConsoleCommandOwner(t)
			// Seed the original byte separator used when joining server-name tokens.
			sep := serverConfigOwnBytes(t, 0x587000, 104484, 2)
			copy(sep, []byte{' ', 0})
			old := legacy.Nox_xxx_sysopGetPass_40A630()
			t.Cleanup(func() { legacy.Nox_xxx_sysopSetPass_40A610(old) })
			got := o.call(t, "set "+tc.command, false, tc.args...)
			stored := legacy.Nox_xxx_sysopGetPass_40A630()
			if tc.command == "name" {
				stored = alloc.GoString((*byte)(legacy.PortTestServerConfigPointer("name-get", 0, nil)))
			}
			if !got || stored != tc.want {
				t.Fatalf("result=%t stored=%q want=%q", got, stored, tc.want)
			}
			rows = append(rows, row{tc.command, tc.args, got, stored, []byte(stored), o.printer.lines})
		})
	}
	spellbookCapture(t, "console-commands-names", rows, "c649194609d121c8ece1be13e3edbaabb3ac82dc27beafa81d45dcd349b667fe")
}
func TestConsoleCommandsQuality(t *testing.T) {
	type row struct {
		Command      string
		Result       bool
		Type, Online int
		Rate         uint64
	}
	var rows []row
	o := newConsoleCommandOwner(t)
	oldType, oldOnline := legacy.Get_nox_server_connectionType_3596(), legacy.Get_dword_5d4594_2650652()
	t.Cleanup(func() { legacy.Set_nox_server_connectionType_3596(oldType); legacy.Set_dword_5d4594_2650652(oldOnline) })
	table := serverConfigOwnBytes(t, 0x587000, 4664, 40)
	serverConfigOwnBytes(t, 0x587000, 4728, 4)
	for i := 0; i < 5; i++ {
		binary.LittleEndian.PutUint32(table[8*i:], uint32(i))
		binary.LittleEndian.PutUint32(table[8*i+4:], uint32(17+11*i))
	}
	for _, tc := range []struct {
		name         string
		kind, online int
	}{{"modem", 4, 1}, {"isdn", 3, 1}, {"cable", 2, 1}, {"T1", 1, 1}, {"LAN", 1, 0}} {
		got := o.call(t, "set quality "+tc.name, false)
		rate := legacy.PortTestServerConfigScalar("rate-get", 0, 0)
		if !got || rate != uint64(17+11*tc.kind) || legacy.Get_nox_server_connectionType_3596() != tc.kind || legacy.Get_dword_5d4594_2650652() != tc.online {
			t.Fatal(tc, got, rate)
		}
		rows = append(rows, row{tc.name, got, tc.kind, tc.online, rate})
	}
	spellbookCapture(t, "console-commands-quality", rows, "234363794ddcf6f5bfeca6644e419ee290bbd64cefee287945347676c1075aad")
}
