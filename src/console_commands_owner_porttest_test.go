//go:build porttest

package opennox

import (
	"context"
	"fmt"
	"github.com/opennox/libs/console"
	"github.com/opennox/libs/strman"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/common/memmap"
	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"strings"
	"testing"
	"unsafe"
)

type consoleCommandLine struct {
	Color console.Color
	Text  string
}
type consoleCommandPrinter struct{ lines []consoleCommandLine }

func (p *consoleCommandPrinter) Print(c console.Color, s string) {
	p.lines = append(p.lines, consoleCommandLine{c, s})
}
func (p *consoleCommandPrinter) Printf(c console.Color, s string, args ...interface{}) {
	p.Print(c, fmt.Sprintf(s, args...))
}

type consoleCommandOwner struct {
	*serverOptionsOwner
	console *console.Console
	printer *consoleCommandPrinter
}

func newConsoleCommandOwner(t *testing.T) *consoleCommandOwner {
	t.Helper()
	o := &consoleCommandOwner{serverOptionsOwner: newServerOptionsOwner(t), printer: &consoleCommandPrinter{}}
	t.Cleanup(legacy.PortTestConsoleContext())
	serverConfigOwnBytes(t, 0x5D4594, 823804, 1932)
	serverConfigOwnBytes(t, 0x5D4594, 822660, 1024)
	connected := serverConfigOwnBytes(t, 0x5D4594, 815764, 4)
	clear(connected)
	t.Cleanup(noxflags.PortTestGameFlags(0))
	flags := legacy.Nox_xxx_getServerSubFlags_409E60()
	t.Cleanup(func() { legacy.Sub_409EC0(-1); legacy.Sub_409E70(int(flags)) })
	var entries []strman.Entry
	for _, pair := range [][2]string{
		{"guimsg.c:systemmsg", "system %s"}, {"cmd_token:on", "ON"}, {"cmd_token:off", "OFF"}, {"weapons", "weapons %s"}, {"staffs", "staffs %s"},
		{"monsters", "monsters %s"}, {"monsterrespawn", "respawn %s"}, {"MapCycleOn", "cycle ON"}, {"MapCycleOff", "cycle OFF"},
		{"setgamename", "name %S"}, {"sysoppasswordset", "password set"}, {"playersset", "players %d"}, {"lessonsset", "lessons %d"},
		{"cantbanyourself", "cannot ban self"}, {"cantkickyourself", "cannot kick self"}, {"banned", "banned %s"}, {"banDisallow", "blocked %s"}, {"kicked", "kicked %s"}, {"notyetimplemented", "not implemented"}, {"userslist", "users"}, {"SysMuted", "server muted"}, {"ClientMuted", "client muted"},
		{"set", "set"}, {"observermode", "observer %s"}, {"ExecutingFunction", "execute %s"}, {"InvalidFunction", "invalid function %s"}, {"RemoteSysop", "remote %s %s"}, {"notinobserver", "not observer"}, {"invalidattempt", "invalid attempt %s %s"}, {"NotInChat", "not in chat %s"}, {"spellenabled", "enabled %s"}, {"spelldisabled", "disabled %s"}, {"invalidspell", "invalid %s"}, {"Muted", "muted %s"}, {"UnMuted", "unmuted %s"}, {"UserNotFound", "missing %s"},
	} {
		id := pair[0]
		if !strings.Contains(id, ":") {
			id = "parsecmd.c:" + id
		}
		entries = append(entries, strman.Entry{ID: strman.ID(id), Vals: []strman.Variant{{Str: pair[1]}}})
	}
	configure, restore := o.c.srv.PortTestMeterStrings(entries...)
	t.Cleanup(restore)
	configure(0)
	o.console = console.New(o.printer)
	o.console.SetCheats(true)
	for _, cmd := range noxCommands {
		o.console.Register(cmd)
	}
	old := legacy.GetConsole
	legacy.GetConsole = func() *console.Console { return o.console }
	t.Cleanup(func() { legacy.GetConsole = old })
	// Original dispatch compares these localized token-table entries directly.
	for i, text := range map[uintptr]string{6: "respawn", 7: "server"} {
		ptr := memmap.PtrPtr(0x587000, 94468+4*i)
		before := *ptr
		wide, free := alloc.CString16(text)
		*ptr = unsafe.Pointer(wide)
		t.Cleanup(func() { *ptr = before; free() })
	}
	return o
}
func (o *consoleCommandOwner) call(t *testing.T, path string, client bool, args ...string) bool {
	t.Helper()
	parts := strings.Fields(path)
	cmds := o.console.Commands()
	var cmd *console.Command
	for _, part := range parts {
		cmd = nil
		for _, c := range cmds {
			if c.Token == part {
				cmd = c
				break
			}
		}
		if cmd == nil {
			t.Fatal("unregistered command", path)
		}
		cmds = cmd.Sub
	}
	if cmd.LegacyFunc == nil {
		t.Fatal("missing legacy command", path)
	}
	ctx := console.AsServer(context.Background())
	if client {
		ctx = console.AsClient(ctx)
	}
	o.printer.lines = nil
	return cmd.LegacyFunc(ctx, o.console, len(parts), append(parts, args...))
}
