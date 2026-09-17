# Spatial targeting

## Result and scope

**The Go conversion is fully qualified.** C baseline **76ba09f8** is pushed;
its parent is **6541718f**. This batch covers eleven live C
functions (554 body lines): projectile prediction and eligibility, indexed contact
candidates, cursor selection, ray forwarding, wall normals and local tile regions.
See [spatial-targeting-scope.json](spatial-targeting-scope.json).

The scope closes the remaining spatial C calls in motionTrace, permitting removal
of its temporary C allocation and the last collision-core C export. Keep this
connected batch smaller than the usual LOC target rather than adding unrelated
player-death or parsing code. Production C is **40,218 lines / 74 files / zero
reference C**, down **559 physical lines**. One C export remains for an actual
spell caller. Ten batch interfaces, the final collision-core export and one C
state global are retired. Cursor ownership is Go state; AI prediction and projectile
tracing no longer allocate temporary C point records.

## Qualified C baseline

The baseline changed only tests and documentation. Its tests called actual C helpers and reused real
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


## Native qualification and callback correction

The eleven new focused roots passed unchanged on the first native run (0.203s).
The broader default run then exposed an existing curve callback bridge defect:
`sub_4BEDE0` passed its integer token through `CallVoidPtr3`, so an otherwise valid
opaque token could happen to fall in the Go heap and trigger cgo pointer checking.
Server/highres runs were stopped and joined before source edits. This is unrelated
to spatial arithmetic but must be corrected for reliable qualification.

Keep both point arguments typed as pointers and the token typed as an integer at
the C callback boundary. Add a deterministic contract whose token bits equal the
address of a live Go object containing pointers, without interpreting the token as
a pointer. Preserve existing curve captures and all spatial goldens. This small,
reversible ABI correction is recorded for review; no gameplay algorithm changes.


The deterministic regression failed before the correction with the same cgo
pointer-check panic, then passed afterward. Native-focused-2 passes **13 roots**
in **0.299s**: eleven spatial groups plus existing curve captures and the new opaque
token contract. Frozen spatial and curve expectations were not regenerated.
Artifacts retain curve-token-before.log and the interrupted-run audit.


## Completed native qualification

All three affected-target sweeps pass, no skips, with **257 identical captures /
157,554 records** per target matching C. Root/subtest counts and timings are in
[spatial-targeting-native-qualification.json](spatial-targeting-native-qualification.json).
Fresh production passes three binaries and ABI/interface checks, the exact known
asset-suite failures, gameplay, save/load and forced flat-map regeneration. All four
gates share an unchanged source manifest. Every session is joined.

The only qualification-driven correction was the older curve callback argument
kind defect documented above. The spatial implementation matched the frozen C
captures on its first focused run. No goldens were changed. The private triangle
helper always returns a low byte in 0–3 even when upper bits remain, so the wall
switch's out-of-range fallback is not reachable through that helper.

Artifacts: build/port-spatial-targeting/native2-{default,server,highres},
native-production, native-focused-2, native-interface-audit.json and frozen-captures.json.
**install-native.py and the installed drafts are consumed; do not reinstall.**
The next read-only proposal covers monster control and definitions: 40 live
functions / 1,001 body lines plus three orphan candidates / 31 lines. Complete
reachability and numeric-table review before baseline work; no next source is installed.

Disk cleanup reclaimed **6.184 GiB** from twelve completed scenario copies after
verifying every removed file against the original asset hash. Per-run restoration
manifests preserve how to recreate them. Original assets/archive, changed files
and all reports remain. About **20 GiB** is free. The spatial-targeting
**deduplicate-completed-assets.py --apply is consumed; never repeat deletion mode.**
See build/port-spatial-targeting/completed-assets-{plan,audit}.json and applied log.
