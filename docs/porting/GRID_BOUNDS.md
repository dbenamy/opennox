# Grid-coordinate bounds defect — 2026-09-11

The next owner port is tile lookup 411160 in GAME1.c. Before porting, original-C
probes reproduce an overflow in its lower-bound check:

`i - 1 <= 0 || i >= 127 || j - 1 <= 0 || j >= 127`

Float-to-int conversion returns INT32_MIN for masked invalid conversions.
Subtracting 1 wraps to INT32_MAX under this build's -fno-strict-overflow, so both
lower and upper checks incorrectly allow INT32_MIN. Compiled 386 disassembly
confirms LEA -1 followed by signed tests, not a direct comparison to 1.

The baseline runs six invalid coordinate pairs in separate child processes:
NaN in either axis, positive infinity, negative infinity, and positive/negative
maximum finite float. Each reaches a null test grid and faults inside the C
lookup. The null grid intentionally detects any access: invalid coordinates
must return -1 without using map storage. Ordinary zero/negative/too-large
coordinates already reject safely. No crash was induced in a live game or
outside an isolated test process. The normal allocated-grid consequence depends
on coordinates and memory layout; it is not proven to crash for every input.

## Proposed behavior change

Use `i <= 1` and `j <= 1` directly, retaining both upper checks. This preserves
all previously accepted valid grid indices (2..126) and rejects INT32_MIN before
any table access. Return remains -1, the existing invalid-coordinate result.

User decision is needed before adopting this bounds repair and then porting
the owner. Prepare a tested patch; preserve the original-C baseline. Production
C remains **140,455 physical lines**, 153 files, zero reference C. Local artifacts:
build/port-grid-bounds. Float-to-int exports were deferred for measured C-caller
cost; see FLOAT_INT.md. Keep the native helper for future Go owners.
