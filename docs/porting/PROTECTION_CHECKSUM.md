# Protection checksum conversion — 2026-09-10

The first conversion removes the two checksum implementations from
`legacy/GAME5_2.c`. Existing C callers retain their function declarations and
call generated C-to-Go exports in `legacy/protection.go`. Both those exports and
the existing root `protectBytes` helper use `internal/protection.Checksum`.
It XORs complete little-endian words and ignores trailing bytes. The nullable
entry point preserves zero for null with any unsigned length.

The bridge processes bounded word-aligned chunks so the unsigned C byte count
cannot become a negative Go slice length on 386. It preserves all return bits
through C's signed int, does not retain the pointer, and does not mutate input.
Compile-time C assertions check both scalar argument/return widths. There is no
translated struct layout in this chunk. Invalid pointer/length combinations are
outside the old contract; no attempt is made to preserve undefined C reads.

## Evidence

The inventory and pre-conversion reference tests were committed as `00228a81`.
The reference functions in `internal/protectionref/reference.c` match revision
`0e9d2e1f` exactly apart from symbol names; provenance/build-tag comments were
added. They compile only with `porttest` and remain callable for future tests.

- Existing Go versus original C passed before replacing the C definitions.
- Fixed expected-value tests and deterministic differential tests pass on 386.
  Every length 0–1024 is tested at eight byte offsets, plus 1,000 larger generated
  buffers, chunk-boundary lengths through two MiB, null lengths through UINT_MAX,
  high-bit return values, trailing bytes, and input mutation checks. Both actual
  C ABI entry points are compared to the historical reference after conversion.
- Five-second fuzz runs passed: 587,107 pure checksum property cases and 524,264
  differential cases through the C ABI. These are bounded runs, not exhaustive
  proofs. Pure Go unit tests also pass on amd64.
- Differential tests pass under default, server and highres tags on 386.
- All three production targets build. Their symbol tables contain the generated
  Go export bridges and no `reference_checksum` test implementation.
- A fresh standard-client warrior scenario exits 0 with override disabled;
  both preserved gameplay screenshots match. This is integration coverage, not
  proof that every protection caller executed or that all player/item state is
  equivalent. Save/load byte differences remain the separately documented issue.
- Full default asset-backed suite: 15 passing packages, 3 known failing packages,
  32 without tests; no compilation/vet failures. The additional passing package
  is the new checksum package. Existing blobs, sprite and audio failures remain.

Commands below run from src with the environment in RECOVERY.md:

```bash
go test -count=1 ./internal/protection
go test -tags porttest -count=1 -run 'TestProtectionCReference|FuzzProtectionCReference' .
go test -tags 'server porttest' -count=1 -run 'TestProtectionCReference|FuzzProtectionCReference' .
go test -tags 'highres porttest' -count=1 -run 'TestProtectionCReference|FuzzProtectionCReference' .
go test -run '^$' -fuzz '^FuzzChecksum$' -fuzztime 5s -parallel 2 ./internal/protection
go test -tags porttest -run '^$' -fuzz '^FuzzProtectionCReference$' -fuzztime 5s -parallel 2 .
go run ./internal/noxbuild -o ../build/port-checksum/bin
```

Local logs are under `build/port-checksum`; the scenario is
`build/baseline/runs/checksum-port`. Recovery does not depend on those logs.
Use the tracked warrior scenario and the recovery instructions to regenerate it.

## Measured cost

Three 200 ms benchmark repetitions on this VM, after builds/scenario completed,
gave these median nanoseconds per operation (all zero Go allocations):

| Bytes | Direct Go | Go→C reference | Go→C→Go exported implementation |
| --- | ---: | ---: | ---: |
| 64 | 15.77 | 60.37 | 176.8 |
| 256 | 63.80 | 86.28 | 249.1 |
| 4096 | 930.2 | 537.2 | 2046.0 |

The remaining C callers incur a new C-to-Go transition, and the Go loop is slower
than the C reference for the measured larger buffer. These microbenchmarks do
not measure whole-game impact or direct C-to-C cost; both C columns include an
outer Go-to-C call from the benchmark. VM timings vary. This is a recorded cost
of the small migration boundary, not a performance improvement. Porting cohesive
callers can remove boundary crossings; profile relevant workloads before doing
that solely for speed. Reproduce with
`go test -tags porttest -run '^$' -bench '^BenchmarkProtectionChecksum$' -benchtime 200ms -count 3 .`.

## Source-size checkpoint and next work

Production C: **142,637 physical lines**, down **28**, across 153 `.c` files.
Test-reference C: 33 lines, separately counted. See [C_LOC.md](C_LOC.md).

Next select another cohesive leaf using caller/state evidence. Do not infer that
the surrounding protection manager is equally simple: it owns global keys,
records, mutation checks and object dependencies that this checksum does not.
