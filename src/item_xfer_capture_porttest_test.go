//go:build porttest

package opennox

import (
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/internal/cryptfile"
	"github.com/opennox/opennox/v1/server"
)

type itemXferCaptureRow struct {
	Case            string
	State           objectXferState
	Health          []byte
	Modifiers       [4]string
	Children        [12]*objectXferState
	StreamChecksum  uint32
	WrittenChecksum uint32
}

func itemXferCaptureCase(t *testing.T, rows *[]itemXferCaptureRow, u *server.Object, checksum ...uint32) {
	r := itemXferCaptureRow{Case: t.Name(), State: objectXferSnapshot(u)}
	if u.HealthData != nil {
		r.Health = append([]byte(nil), unsafe.Slice((*byte)(unsafe.Pointer(u.HealthData)), int(unsafe.Sizeof(*u.HealthData)))...)
	}
	name := u.Server().Types.ByInd(int(u.TypeInd)).ID()
	switch name {
	case "portweapon", "portarmor", "portammo", "portteam":
		mods := (*[4]*server.ModifierEff)(u.InitData)
		for i, m := range mods {
			v := uint32(0)
			if m != nil {
				v = 1
				r.Modifiers[i] = m.Name()
			}
			binary.LittleEndian.PutUint32(r.State.Data["init"][4*i:], v)
		}
	case "portmonstergenerator":
		children := (*[12]*server.Object)(u.UpdateData)
		for i, ch := range children {
			v := uint32(0)
			if ch != nil {
				v = 1
				state := objectXferSnapshot(ch)
				r.Children[i] = &state
			}
			binary.LittleEndian.PutUint32(r.State.Data["update"][4*i:], v)
		}
	}
	if len(checksum) > 0 {
		r.StreamChecksum = checksum[0]
	} else if f := cryptfile.Global(); f != nil {
		r.StreamChecksum = f.PortTestChecksum()
	}
	if len(checksum) > 1 {
		r.WrittenChecksum = checksum[1]
	}
	*rows = append(*rows, r)
}
func itemXferReadRecord(t *testing.T, u *server.Object, path string) uint32 {
	t.Helper()
	if err := cryptfile.OpenGlobal(path, cryptfile.ReadOnly, -1); err != nil {
		t.Fatal(err)
	}
	defer cryptfile.Close()
	if err := u.CallXfer(nil); err != nil {
		t.Fatal(err)
	}
	return cryptfile.Global().PortTestChecksum()
}
