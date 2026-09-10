# Protection record helpers — 2026-09-10

This chunk moves three GAME5_2.c helpers into Go:

- `sub_56F590`: return the first record whose encrypted ID matches the requested
  ID XOR the current key; return null for an empty list or missing ID.
- `sub_56F6F0`: return the zero-based list element, or null if absent.
- `sub_56F720`: swap the two payload words, leave links untouched, and increment
  the existing 32-bit swap counter. A self-swap still increments it. A missing
  operand changes nothing. Counter overflow wraps.

All three symbols were retained in the standard baseline binary. Lookup is used
by creation/deletion/value-validation paths; index lookup and swap are used by
protection rekeying. The only diagnostic callback, `nullsub_31`, is an empty
function. Its calls are omitted. The manager's allocation, deletion, random
selection and rekeying remain C. No random calls were added or removed.

`internal/protection.Record` names the four-word layout. On 386, compile-time
assertions in `legacy/protection_records.go` verify size 16, alignment 4, payload
offsets 0/4 and pointer offsets 8/12. Records remain C-allocated and C-owned.
Exports return original record pointers; no record is copied into the Go heap.
Lists must remain valid and acyclic, as required by the original code. Negative
index lookup returns null immediately rather than walking the list to its end;
there were no observable mutations during that walk.

## Validation

Pre-conversion behavior tests and pure Go helpers are committed at `ab782be4`.
The same 2,000 deterministic ABI scenarios run before and after replacing the C:
empty/single/multiple lists, four keys including zero/high bits, duplicate IDs,
misses, extreme positive/negative indices, null/self/adjacent/non-adjacent swaps,
all payload words, unchanged links/head/key and wrapping counter values.
The fixture owns C-heap storage and restores globals after each serial call.
No C reference is added; the test-before-port commit preserves the old executable
path and the same expected-value checks for reproduction.

Commands from src, with RECOVERY.md's target environment:

```bash
go test -count=1 ./internal/protection
go test -tags porttest -count=1 -run '^TestProtection(RecordsABI|ABI)$' .
go test -tags 'server porttest' -count=1 -run '^TestProtection(RecordsABI|ABI)$' .
go test -tags 'highres porttest' -count=1 -run '^TestProtection(RecordsABI|ABI)$' .
go run ./internal/noxbuild -o ../build/port-records/bin
```

Logs and generated binaries are local under build/port-records. The scalar/layout
checks and ABI tests target 386; pure Go tests also run on amd64. The new exports
introduce C-to-Go crossings at these remaining C callers, as in the checksum
chunk. No performance improvement or full manager equivalence is claimed.

Production C after this chunk: **142,570 lines**, down **67**, with 153 `.c`
files and **0** test-reference C lines. See [C_LOC.md](C_LOC.md).

Post-conversion ABI checks pass under default/server/highres tags on 386. All
three production binaries build and contain the corresponding Go export bridges.
Pure Go tests pass on 386 and amd64; root compilation and E2E oracle tests pass.
Fresh standard-client warrior scenario `records-port` exits 0 against both
preserved screenshots with overrides disabled. This covers integration, not
every protection-manager path. Existing unrelated baseline failures remain.
