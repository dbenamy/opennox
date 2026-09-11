# Line projection helpers — 2026-09-11

Scope: clamped projection with supplied length 57C790 and projection/bounds test
57C8A0. Each has one live C caller (GAME5.c and GAME4_1.c respectively). Existing
server.PointOnTheLine has different rounding and separate Go callers, so it is
outside this conversion. Keep both C entries and their output pointer behavior.

## Original-386 baseline

21,656 cases comprise 200 analytic axis-aligned projections, 1,456 IEEE edge and
output-overlap cases, and 20,000 deterministic raw-word cases across both helpers.
Inputs include signed zeros, subnormals, infinity, quiet/signaling NaNs, zero
length, degenerate segments and output overlap with either input (including
partial overlap). Guarded C-owned storage checks that only the two output words
change. The analytic cases independently check inclusive bounds and clamping.

The asset-free original-C fixture stores return/X/Y as three little-endian words
per case (259,872 bytes). SHA-256:
`a4d52cc0d1980c5779855e39a4a66715c08d828b4662fe5f10ae590969c38f42`.
After capture, a second original-C run verified the fixture and all analytic
checks. No C reference implementation is retained after conversion. Inputs remain
deterministic Go test code; local artifacts are under build/port-line-projection.

Disassembly, rather than decompiler local types, determines intermediate rounding:

- C790 rounds the X quotient to float32 before adding old line X; output X is
  stored as float32 but its wider sum is retained for clamp comparisons. Y adds
  line Y reloaded after the X store, then rounds before its store/comparisons.
- C8A0 rounds the squared denominator to float32; X uses the full dot product,
  then rounds before storing/comparing. Dot rounds to float32 before Y's multiply;
  Y adds line Y reloaded after the X store and retains the wider result for bounds.
- Both reload endpoint bounds after storing both output coordinates, so output
  aliasing can change those bounds. Their ordered comparisons also govern NaNs.

All cases pass against original C before replacement. Production C before this
chunk: **141,000 physical lines**, 153 files, zero reference C.

## Native implementation

The two C entries now bridge to Go. Finite arithmetic follows the verified PC53
stages above. Narrow private arithmetic/load/store helpers preserve x87 NaN sign
and payload: float32 loads quiet signaling NaNs; two NaNs select the larger
significand and positive sign on ties; invalid operations produce negative
canonical quiet NaN. This matters with this environment's GO386=softfloat,
whose ordinary conversions/arithmetic canonicalize NaNs to positive zero-payload
quiet NaN. An independent local x87 instruction probe verified propagation and
invalid-operation results under control word 0x027f.

Before this compatibility handling, all 6,681 differing output words were NaN
representations: no finite numeric value or return differed. Afterward all 21,656
cases match the unchanged original-C fixture exactly. No tolerance or fixture
regeneration was used. Original-C baseline is recoverable at `7947f8ca`.

## Validation and source size

Accumulated protection/network/waypoint/rules/spell-class/ping/glyph/collision/
projection tests pass in default, server, and highres 386 variants. All three
production ELF32 binaries build. Fresh `line-projection-port` headless gameplay
passes both preserved screenshot checks with overrides disabled. The full suite
has the exact prior 1,553 failure-entry multiset (15 passing, 3 known failing,
32 skipped/no-test packages). Logs/binaries are under build/port-line-projection.

Production C: **140,903 physical lines**, 153 files, **−97** from 141,000.
Reference C: **0**. Both remaining C callers keep their ABI bridges.
