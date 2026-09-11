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
