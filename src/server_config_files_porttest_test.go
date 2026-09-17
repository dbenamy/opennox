//go:build porttest

package opennox

import (
	"github.com/opennox/opennox/v1/legacy"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func serverConfigFileOwner(t *testing.T) string {
	t.Helper()
	o := serverConfigListOwner(t)
	dir := serverOptionsRulesOwner(t, o)
	serverConfigOwnBytes(t, 0x5D4594, 1563936, 1024)
	for _, v := range []struct {
		off  uintptr
		size int
		text string
	}{{202212, 16, "[Banned Users]"}, {202228, 16, "[Allowed Users]"}, {54280, 8, "ban.txt"}} {
		b := serverConfigOwnBytes(t, 0x587000, v.off, v.size)
		clear(b)
		copy(b, v.text)
	}
	return dir
}
func TestServerConfigAdmissionFileRead(t *testing.T) {
	type row struct {
		Name, Op         string
		Result           int
		Allowed, Blocked []serverConfigListRow
	}
	var rows []row
	for _, tc := range []struct {
		name, op, text   string
		missing          bool
		result           int
		allowed, blocked []serverConfigListRow
	}{
		{name: "missing", op: "read", missing: true, result: 0},
		{name: "empty", op: "read", result: 1},
		{name: "unrecognized", op: "read", text: "unknown\nignored\n", result: 1},
		{name: "allowed", op: "read", text: "[Allowed Users]\nAlice\nBob Smith\n\n", result: 1, allowed: []serverConfigListRow{{Name: "Alice"}, {Name: "Bob Smith"}}},
		{name: "blocked", op: "read", text: "[Banned Users]\nAlice\n0\nBob\n127.0.0.1\n\n", result: 1, blocked: []serverConfigListRow{{Name: "Alice"}, {Name: "Bob", Address: "127.0.0.1"}}},
		{name: "both-crlf", op: "read", text: "[Banned Users]\r\nAlice\r\n0\r\n\r\n[Allowed Users]\r\nBob\r\n", result: 1, blocked: []serverConfigListRow{{Name: "Alice"}}, allowed: []serverConfigListRow{{Name: "Bob"}}},
		{name: "reverse", op: "read", text: "[Allowed Users]\nBob\n\n[Banned Users]\nAlice\n0\n", result: 1, blocked: []serverConfigListRow{{Name: "Alice"}}, allowed: []serverConfigListRow{{Name: "Bob"}}},
		{name: "truncated-pair", op: "read", text: "[Banned Users]\nAlice\n", result: 1},
		{name: "blank-address", op: "read", text: "[Banned Users]\nAlice\n\n", result: 0},
		{name: "direct-allowed", op: "allowed", text: "One\nTwo\n\nThree\n", result: 1, allowed: []serverConfigListRow{{Name: "One"}, {Name: "Two"}}},
		{name: "direct-blocked", op: "blocked", text: "One\n0\nTwo\n10.0.0.2\n", result: 1, blocked: []serverConfigListRow{{Name: "One"}, {Name: "Two", Address: "10.0.0.2"}}},
		{name: "tabs", op: "allowed", text: "\tOne\tTwo\n\t\n", result: 1, allowed: []serverConfigListRow{{Name: "One"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := serverConfigFileOwner(t)
			if !tc.missing {
				if err := os.WriteFile(filepath.Join(dir, "input.txt"), []byte(tc.text), 0600); err != nil {
					t.Fatal(err)
				}
			}
			result := legacy.PortTestServerConfigFile(tc.op, "input.txt")
			allowed, _ := serverConfigListSnapshot(t, "allowed")
			blocked, _ := serverConfigListSnapshot(t, "blocked")
			if result != tc.result || !reflect.DeepEqual(allowed, tc.allowed) || !reflect.DeepEqual(blocked, tc.blocked) {
				t.Fatal("read contract", result, allowed, blocked, tc)
			}
			rows = append(rows, row{tc.name, tc.op, result, allowed, blocked})
		})
	}
	spellbookCapture(t, "server-config-admission-file-read", rows, "db0c9135af9a761b94e5880cf57fbe0e6e2733996531c2d8ba31869d69bf7164")
}
func TestServerConfigAdmissionFileWrite(t *testing.T) {
	dir := serverConfigFileOwner(t)
	old := legacy.PlatformTicks
	legacy.PlatformTicks = func() uint64 { return 100 }
	t.Cleanup(func() { legacy.PlatformTicks = old })
	serverConfigListAdd("blocked", "Permanent", "", 0)
	serverConfigListAdd("blocked", "Address", "192.0.2.1", 0)
	serverConfigListAdd("blocked", "Temporary", "", 1)
	serverConfigListAdd("allowed", "Alice", "", 0)
	serverConfigListAdd("allowed", "Bob Smith", "", 0)
	if legacy.PortTestServerConfigFile("write", "output.txt") != 1 {
		t.Fatal("write result")
	}
	raw, err := os.ReadFile(filepath.Join(dir, "output.txt"))
	if err != nil {
		t.Fatal(err)
	}
	want := "[Banned Users]\nPermanent\n0\nAddress\n192.0.2.1\n\n[Allowed Users]\nAlice\nBob Smith\n"
	if string(raw) != want {
		t.Fatalf("file bytes %q want %q", raw, want)
	}
	if err := os.Mkdir(filepath.Join(dir, "directory"), 0700); err != nil {
		t.Fatal(err)
	}
	if legacy.PortTestServerConfigFile("write", "directory") != 0 {
		t.Fatal("directory destination must fail")
	}
	before, nodes := serverConfigListSnapshot(t, "blocked")
	if len(before) != 3 {
		t.Fatal("writer changed live entries")
	}
	// Shutdown returns the original first blocked-node address, after freeing it.
	if legacy.PortTestServerConfigAdmission("close", 0, nil, nil) != nodes[0] {
		t.Fatal("close result")
	}
	for _, kind := range []string{"allowed", "blocked"} {
		s, _ := serverConfigListSnapshot(t, kind)
		if len(s) != 0 {
			t.Fatal("close retained nodes", kind, s)
		}
	}
	saved, err := os.ReadFile(filepath.Join(dir, "ban.txt"))
	if err != nil || string(saved) != want {
		t.Fatal("close persistence", err, string(saved))
	}
	if got := legacy.PortTestServerConfigAdmissionInit(); got != 1 {
		t.Fatal("initialize result", got)
	}
	allowed, _ := serverConfigListSnapshot(t, "allowed")
	blocked, _ := serverConfigListSnapshot(t, "blocked")
	if len(allowed) != 2 || len(blocked) != 2 || blocked[0].Name != "Permanent" || blocked[1].Address != "192.0.2.1" {
		t.Fatal("roundtrip", allowed, blocked)
	}
	spellbookCapture(t, "server-config-admission-file-write", []any{string(raw), before, allowed, blocked}, "a483e0a382033d0cd41f4070c97ba70b35ac5878b60ab648f0c67f116c935345")
}

func TestServerConfigAdmissionFileEncoding(t *testing.T) {
	type row struct {
		Input, File string
		Restored    []serverConfigListRow
	}
	var rows []row
	for _, tc := range []struct{ name, narrow, restored string }{{"Alice", "Alice", "Alice"}, {"Åke", "\xc5ke", "Åke"}, {"Àßÿ", "\xc0\xdf\xff", "Àßÿ"}, {"abcdefghijklmnopqrstuvwxy", "abcdefghijklmnopqrstuvwxy", "abcdefghijklmnopqrstuvwxy"}, {"AxĀBy", "Ax", "Ax"}} {
		t.Run(tc.name, func(t *testing.T) {
			dir := serverConfigFileOwner(t)
			serverConfigListAdd("allowed", tc.name, "", 0)
			if legacy.PortTestServerConfigFile("write", "roundtrip.txt") != 1 {
				t.Fatal("encoding write")
			}
			raw, err := os.ReadFile(filepath.Join(dir, "roundtrip.txt"))
			if err != nil {
				t.Fatal(err)
			}
			want := "[Banned Users]\n\n[Allowed Users]\n" + tc.narrow + "\n"
			if string(raw) != want {
				t.Fatalf("legacy byte encoding %q want %q", raw, want)
			}
			legacy.PortTestServerConfigAdmission("allowed-remove", 0, nil, nil)
			if legacy.PortTestServerConfigFile("read", "roundtrip.txt") != 1 {
				t.Fatal("encoding read")
			}
			restored, _ := serverConfigListSnapshot(t, "allowed")
			if len(restored) != 1 || restored[0].Name != tc.restored {
				t.Fatal("legacy decoding", restored, tc)
			}
			rows = append(rows, row{tc.name, string(raw), restored})
		})
	}
	spellbookCapture(t, "server-config-admission-file-encoding", rows, "c079ebdafc623b0d66ec3bb342151b1e443d743b7d9ff3a2d0726d8658d3c0c3")
}
