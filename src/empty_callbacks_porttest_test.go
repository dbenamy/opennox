//go:build porttest

package opennox

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"testing"
	"unsafe"

	"github.com/opennox/opennox/v1/legacy"
	"github.com/opennox/opennox/v1/legacy/common/alloc"
	"github.com/opennox/opennox/v1/server"
)

func emptyCallbackCall(c legacy.PortTestEmptyCallback, args [6]unsafe.Pointer) func() {
	if c.Kind == "duration" {
		return func() { server.CallDurSpellDiscard(c.Pointer, (*server.DurSpell)(args[0])) }
	}
	switch c.Arity {
	case 3:
		return func() {
			server.CallModifierEffect3Discard(c.Pointer, (*server.ModifierEff)(args[0]), (*server.Object)(args[1]), (*server.Object)(args[2]))
		}
	case 5:
		return func() {
			server.CallModifierEffect5(c.Pointer, (*server.ModifierEff)(args[0]), (*server.Object)(args[1]), (*server.Object)(args[2]), (*server.Object)(args[3]), args[4])
		}
	case 6:
		return func() {
			server.CallModifierEffect6(c.Pointer, (*server.ModifierEff)(args[0]), (*server.Object)(args[1]), (*server.Object)(args[2]), (*server.Object)(args[3]), (*server.Object)(args[4]), args[5])
		}
	default:
		panic("unsupported empty callback arity")
	}
}

func TestEmptyCallbackContracts(t *testing.T) {
	const want = "848b76f173284c29edddd5d637863061d643428004ed8f51098f51863ea44b77" // Frozen from three actual-C runs.
	callbacks := legacy.PortTestEmptyCallbacks()
	if len(callbacks) != 10 {
		t.Fatalf("callbacks: got %d, want 10", len(callbacks))
	}
	expected := []struct {
		symbol, effect, kind string
		arity                int
	}{
		{"nullsub_22", "ReadinessEffect", "damage", 5},
		{"nullsub_36", "ReplenishmentEffect", "damage", 5},
		{"nullsub_38", "FrostEffect", "damage", 5},
		{"nullsub_39", "PanicEffect", "damage", 5},
		{"nullsub_40", "ResilienceEffect", "defend", 6},
		{"nullsub_41", "BreakingEffect", "defend", 6},
		{"nullsub_42", "PunctureProneEffect", "defend", 6},
		{"nullsub_43", "ParasiteUpdate", "update", 3},
		{"nullsub_44", "AttractionUpdate", "update", 3},
		{"nullsub_29", "EnergyBoltDestroy", "duration", 1},
	}
	seen := make(map[unsafe.Pointer]string)
	type row struct {
		Symbol      string
		Arity, Mask int
		State       string
	}
	var rows []row
	for index, c := range callbacks {
		e := expected[index]
		if c.Symbol != e.symbol || c.Effect != e.effect || c.Kind != e.kind || c.Arity != e.arity {
			t.Fatalf("callback %d metadata mismatch: %+v", index, c)
		}
		if c.Pointer == nil {
			t.Fatalf("%s: nil callback", c.Symbol)
		}
		if old, ok := seen[c.Pointer]; ok {
			t.Fatalf("%s shares identity with %s", c.Symbol, old)
		}
		seen[c.Pointer] = c.Symbol
		if c.Kind == "duration" {
			if c.Pointer != legacy.Get_nullsub_29() {
				t.Fatalf("duration getter identity mismatch")
			}
		} else if p := server.PortTestModifierCallback(c.Kind, c.Effect); p != c.Pointer {
			t.Fatalf("%s registration identity mismatch", c.Symbol)
		}
		buffers, free := alloc.Make([]byte(nil), 6*96)
		for mask := 0; mask < 1<<c.Arity; mask++ {
			for i := range buffers {
				buffers[i] = byte(i*37 + mask*13 + index*19)
			}
			before := bytes.Clone(buffers)
			var args [6]unsafe.Pointer
			for i := 0; i < c.Arity; i++ {
				if mask&(1<<i) != 0 {
					args[i] = unsafe.Pointer(&buffers[i*96+16])
				}
			}
			emptyCallbackCall(c, args)()
			if !bytes.Equal(buffers, before) {
				free()
				t.Fatalf("%s mask %d changed guarded argument records", c.Symbol, mask)
			}
			sum := sha256.Sum256(buffers)
			rows = append(rows, row{c.Symbol, c.Arity, mask, hex.EncodeToString(sum[:])})
		}
		free()
	}
	if len(rows) != 338 {
		t.Fatalf("cases: got %d, want 338", len(rows))
	}
	data, err := json.Marshal(rows)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(data)
	hash := hex.EncodeToString(sum[:])
	t.Logf("empty callback cases=%d sha256=%s", len(rows), hash)
	if hash != want {
		t.Fatalf("callback capture got %s want %s", hash, want)
	}
	if p := os.Getenv("OPENNOX_EMPTY_CALLBACKS_CAPTURE"); p != "" {
		if err := os.WriteFile(p, data, 0600); err != nil {
			t.Fatal(err)
		}
	}
}

// BenchmarkEmptyCallbackDispatch includes each selected native callback registry
// boundary. It measures a callback boundary, not a game-frame time or call frequency.
func BenchmarkEmptyCallbackDispatch(b *testing.B) {
	for _, c := range legacy.PortTestEmptyCallbacks() {
		b.Run(c.Symbol, func(b *testing.B) {
			buffers, free := alloc.Make([]byte(nil), 6*96)
			defer free()
			var args [6]unsafe.Pointer
			for i := range args {
				args[i] = unsafe.Pointer(&buffers[i*96+16])
			}
			call := emptyCallbackCall(c, args)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				call()
			}
		})
	}
}
