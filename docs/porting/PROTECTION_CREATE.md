# Protection record construction — 2026-09-10

Integer and float constructors now share Go record initialization. Allocation
still uses C calloc and the existing C insertion routine takes ownership. That
preserves failure returns, C heap ownership, list insertion and random-number
consumption. The pure Go initializer fills both encrypted words, clears links
and updates the checksum; a nil record returns false without changing the sum.

Pre-conversion tests are committed at `370b6889`. The original C passed 16,224
cases across 1,014 payload patterns, four keys and four C/Go call paths. Patterns
include signed zero, subnormal boundaries, infinities, signaling/quiet NaNs and
seeded random bits. Fixtures start an empty manager, verify payload/checksum/
count/endpoints and restore globals/free the allocation. Empty insertion consumes
no randomness; integration covers insertion into a populated manager.

## Float ABI finding and final boundary

A direct generated C-to-Go float export changed signaling NaN 0x7f800001 into
0x7fc00001 before Go received it. The pre-conversion test failed visibly; no
golden or expected value was relaxed. A raw-bit C shim passed those tests, but a
caller audit showed the only production caller of the old C float constructor
was its Go wrapper. Once that wrapper calls the shared Go initializer directly,
there is no reason to retain a C float entry point or shim. The Go wrapper keeps
its original public signature and uses math.Float32bits, avoiding the float ABI.
The integer constructor retains its C entry point for remaining C callers.

The old float C definition/declarations and temporary shim are removed. Current
constructor tests cover three live paths (C integer, Go integer and Go float),
12,168 cases with the same payload/key coverage. The old four-path harness is
recoverable from the pre-conversion commit. No C reference copy is retained.

Allocation-failure testing supplies nil to the initializer and checks unchanged
checksum state; actual OS allocation exhaustion was not induced. The wrapper
keeps calloc's null result and returns zero. The existing 386 record-size and
field-offset assertions apply. No float arithmetic is introduced in storage.

Commands from src with RECOVERY.md's target environment:

```bash
go test -count=1 ./internal/protection ./internal/e2etest
go test -tags porttest -count=1 -run '^TestProtection(CreateABI|BitsetABI|RecordsABI|ABI)$' .
go test -tags 'server porttest' -count=1 -run '^TestProtection(CreateABI|BitsetABI|RecordsABI|ABI)$' .
go test -tags 'highres porttest' -count=1 -run '^TestProtection(CreateABI|BitsetABI|RecordsABI|ABI)$' .
go run ./internal/noxbuild -o ../build/port-create/accepted-bin
```

Accepted final artifacts use the `accepted-` prefix under build/port-create.
Earlier build directories/logs capture diagnostic variants; float-abi-failure.log
records the rejected direct float export. Production C after the chunk is **142,458**
physical lines, down **45**, in 153 files. Test-reference C is **0**.

All accepted ABI checks pass for default/server/highres on 386; pure Go checks
pass on 386/amd64. All three accepted binaries build, contain the integer Go
export and omit the retired float exports. Full default suite: 15 passing,
3 known failing and 32 no-test packages, with no compilation/vet failures.
The three failures remain blobs, sprite references and audio references.
Fresh warrior scenario `create-port` exits 0 against both preserved screenshots
with override disabled, using the accepted-bin standard client.
