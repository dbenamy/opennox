//go:build porttest

package opennox

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/opennox/opennox/v1/client"
	noxflags "github.com/opennox/opennox/v1/common/flags"
	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/legacy"
)

func TestEngineFatalBoundaries(t *testing.T) {
	const childEnv = "OPENNOX_PORT_FATAL_BOUNDARY_CHILD"
	if mode := os.Getenv(childEnv); mode != "" {
		defer func() {
			fmt.Fprintln(os.Stdout, "unexpected deferred cleanup")
			if v := recover(); v != nil {
				fmt.Fprintln(os.Stdout, "unexpected panic:", v)
				os.Exit(23)
			}
		}()
		switch mode {
		case "room", "painting":
			fmt.Fprintln(os.Stdout, "entering fatal boundary", mode)
			legacy.PortTestInvalidFixtureDispatch(mode == "painting")
		case "light-static", "light-dynamic":
			// Missing-drawable lookup needs only the real, empty drawable owner.
			c := &Client{Client: new(client.Client)}
			oldClient := legacy.GetClient
			legacy.GetClient = func() legacy.Client { return c }
			t.Cleanup(func() { legacy.GetClient = oldClient })
			core := newObjectXferOwner(t)
			u := newObjectXferTyped(t, core, "InvisibleLight")
			u.ObjFlags = 0
			u.Extent = 123
			u.NetCode = 124
			u.ScriptIDVal = 88
			u.ObjClass = 0
			if mode == "light-static" {
				u.ObjClass = 0x20000000
			}
			t.Cleanup(noxflags.PortTestGameFlags(0))
			if err := cryptfile.OpenGlobal(filepath.Join(t.TempDir(), "light.bin"), cryptfile.WriteOnly, -1); err != nil {
				t.Fatal(err)
			}
			fmt.Fprintln(os.Stdout, "entering fatal boundary", mode)
			if err := u.CallXfer(nil); err != nil {
				t.Fatal(err)
			}
		default:
			t.Fatal("unknown fatal contract mode")
		}
		fmt.Fprintln(os.Stdout, "unexpected fatal return")
		os.Exit(24)
	}
	for _, mode := range []string{"room", "painting", "light-static", "light-dynamic"} {
		t.Run(mode, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=^TestEngineFatalBoundaries$", "-test.count=1")
			cmd.Env = append(os.Environ(), childEnv+"="+mode, "GOTRACEBACK=single")
			output, err := cmd.CombinedOutput()
			exit, ok := err.(*exec.ExitError)
			if !ok || exit.ExitCode() != 2 || !bytes.Contains(output, []byte("entering fatal boundary "+mode)) || bytes.Contains(output, []byte("unexpected deferred cleanup")) || bytes.Contains(output, []byte("unexpected fatal return")) {
				t.Fatalf("fatal boundary %s: err=%v output=%.1800s", mode, err, output)
			}
		})
	}
}
