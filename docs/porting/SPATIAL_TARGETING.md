# Spatial targeting

## Scope and checkpoint

Parent **6541718f** is committed and pushed. This batch covers eleven live C
functions (554 body lines): projectile prediction and eligibility, indexed contact
candidates, cursor selection, ray forwarding, wall normals and local tile regions.
See [spatial-targeting-scope.json](spatial-targeting-scope.json).

The scope closes the remaining spatial C calls in motionTrace, permitting removal
of its temporary C allocation and the last collision-core C export. Keep this
connected batch smaller than the usual LOC target rather than adding unrelated
player-death or parsing code. Production C remains **40,777 lines / 74 files**.

## Qualified C baseline

No production behavior has changed. Tests call actual C helpers and reuse real
wall, map-index, object, player and lookup-table owners. Initial geometry coverage
includes every integer point in a tile, fractional boundaries, negative and large
coordinates, integer grid overflow, all eleven wall shapes, neighbours/windows,
and unchanged inputs and output-on-miss. Independent region and normal contracts
supplement byte-exact captures.

All eleven focused roots pass: **11 captures / 26,142 records**, repeated in
separate processes and frozen. Coverage includes projectile owner/team/subclass
branches, actual index and shipped door contacts, prediction raw float bits and
return identity, cursor filters/shapes/radii/height/ties/headless state, and ray flags.
Static mapped-state checks pass. All three broader targets pass with no skips;
see spatial-targeting-c-qualification.json for counts, durations and identical
captures. Every original source fingerprint and three parent production binary
hashes match; the six new files are porttest-only. Production reuse is recorded
in spatial-targeting-c-production-reuse.json. Fresh production is required after
conversion.

## Compatibility review

The tile-triangle C helper writes only the low byte of an integer temporary.
Out-of-tile results can therefore retain upper bits; capture the complete return,
not only the region values 0–3 valid inside a tile. Negative quadrant remainders
are narrowed to unsigned bytes. Preserve these behaviors unless evidence supports
a separately qualified correction.

Wall direction 2 can return success without writing a normal for an out-of-range
starting region. Cursor scores start at zero and ties do not replace the current
candidate. Numeric table review identified the shipped door deltas at
0x587000:196184 (256 bytes); nontrivial door contacts must exercise this table.

Local artifacts: build/port-spatial-targeting. No previous install, freeze or
cleanup script should be replayed. The original asset archive is preserved.

Fixture development caught missing Active flags on indexed objects and a cgo
unsigned-int declaration mismatch; corrected before freezing. The existing
server.LineTraceXxx canonicalizes diagonal endpoints, so a general mathematical
segment oracle is not interchangeable. Axis midpoint crossings provide an
independent door contract; diagonal captures retain current dependency behavior.
This helper is already Go and changing its geometry is outside this conversion.

Compiled-C review: tile offsets spill to float32 before truncation. Cursor
projected Y and ranking score spill to float32; equal rounded scores retain the
first candidate. Circle interval endpoints stay wide during center-line tests.
The compiler removes the box fallback's impossible equal-bound comparisons.
Disassembly artifacts are under build/port-spatial-targeting/c-disassembly.
