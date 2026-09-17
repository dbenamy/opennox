# World geometry and wall/circle/box collision

## Scope and current state

The native conversion replaces **32 reachable functions / 1,385 corrected C body
lines**: segment/rectangle/shape primitives, direction/vector conversion,
map-coordinate transforms and wall/circle/box collision responses. C baseline
**9a945523** and supplemental player-wall contract **7d363edc** are pushed.

Production C is **42,982 physical lines / 74 files / zero reference C**, a reduction
of **1,414 lines** from the corrected baseline. All three native target sweeps and
fresh production qualification pass. Scope and frozen expectations are in the
tracked world-geometry manifests; detailed results follow below.

## Baseline plan

Reuse C-owned guarded word buffers for primitive input/output, preserving raw float
bits and aliasing. Add independent ordinary-geometry contracts plus boundaries,
endpoint and degenerate cases, float widths, NaN/infinity where defined, mutation
and no-hit output preservation. Reuse polygon captures for dependent geometry.

Initialize shipped direction/shape/wall tables explicitly and own extracted live
coefficients separately. Reuse actual server/object/wall owners for collision
response; isolate the real Hit allocation class, bucket array and event list.
Capture object mutations, velocity/acceleration, collision normals/order/dedup,
RNG consumption, wall/player state and callback/update effects. Include nonzero
force coefficients and masses. Snapshot and restore all ownership.

Independent contracts reproduced numeric conversions of float-backed flag words
and a lost quadrant condition. The corrected C is qualified below.
Do not silently fix collision rules or regenerate expectations during conversion.

Repeat/freeze/qualify and commit/push the C baseline before translating. Native
qualification must include direct/dependent three-target tests, ABI/interface
checks, fresh production/known-suite comparison and gameplay/save-load/flat maps.

## Qualified prerequisite corrections

Original-C focused-2/3 runs reproduce three box flag failures: first-object flag
8 and first-class 4/second-flag 0x2000 incorrectly allow force, and first-object
0x8000000 remains set. The second-object counterparts work. In 550F80 three
numeric conversions of float-backed class/flag fields are corrected to raw byte/
word reads, matching object layout and the adjacent circle collision checks.

The four-quadrant contract reproduces two failures: masks 1 and 8 are unreachable
because 550CB0 computes its X comparisons but branches on a constant-zero
decompiler condition byte. Restore the computed less-than/equal condition. The
16.263456 threshold is retained; the caller uses a rotated wall coordinate system.
These are reversible production corrections for later review, before freezing.
All eight primitive capture groups otherwise pass. Golden values were subsequently frozen after repeated corrected-C runs.
Projection draft is consumed; do not reinstall it.

## Coverage

Guarded primitive captures cover segment intersection/projection, inclusive integer
and float rectangles, unsigned wall bounds/clamps and aliased output, diamond
limits, all 256 directions and varying thresholds, normalize/atan2, shape dimensions
and map transforms, both diagonal wall projections, quadrant threshold neighbours,
ordered rectangle crossings and clipped midpoints. Special float bits are retained.

Real factory objects and the actual Hit allocator/list cover circle and box contacts,
force suppression and wake flags, overlap/tangent/separation, masses/friction,
blocked/transparent actor sight rays, deduplicated events, coincident RNG usage and
all four normals. Point/wall checks include cold/hot GameBall lookup, velocity
projection and independent force magnitude. Rotated wall axes cover crossings and
clipped overlap; wall-grid cases supply all eleven shipped shapes, neighbour and
window variants through actual wall indices. Gate cases cover angle, contact frame,
force, owner restriction, team lock, notification throttling, angular and collision
queues and the actual updatable list. Owned guards are checked and pointers are
normalized only when their fixture identities are known.

The broader selection has 543 candidate roots, including every new geometry root
and dependent polygon, AI/path, map generation/painting, object/creature transfer,
inventory, projectile/damage, spell lifecycle, client effects, quickbar/book/UI,
world mechanism and visibility tests. Actual starts/completions and target-specific
availability were checked by the runner. Static-2 passes.

Fixture corrections before freezing: inconsistent cgo spelling of an existing
unsigned-int global; actual ObjectType.Ind and uint16 TypeInd APIs; consuming message
queue semantics in the gate throttle assertion. These did not require production
changes. All projection/walls/gate installation drafts are consumed.

## Focused corrected-C evidence

Focused-7 and focused-8 pass all 16 roots in separate processes; every one of the
15 captures is byte-identical (15,782 records). Expectations are frozen in the root
tests and batch manifest. All baseline gates pass. Logs and repeated captures are under build/port-world-geometry.
Production must be fresh because the prerequisite corrections change C behavior.

The original proposal.json retains pre-correction C. Current frozen-source text is
recorded separately in build/port-world-geometry/corrected-c-functions.json; use
the committed corrected C for conversion, not the old proposal text. The tracked
scope records corrected body hashes. The prerequisite removes five unused local
lines (44,396 physical C lines); this is not conversion progress.

## Qualified C baseline

Default/server/highres pass **543/542/543 roots**, **42,759/42,758/42,759 tests
including subtests**, no skips. **204 captures / 81,403 records** are identical
across targets. Durations: **291.14/383.26/290.61s**. Fresh production passes in
**408.92s**: three builds with ABI/interface checks, the exact known 1,553 asset
failures (15 passing / 3 failing / 32 no-test packages), options gameplay, explicit
save/load and forced flat-map regeneration. References remain unchanged.

All four gates used the same unchanged **2,190-file source manifest**; all sessions
are joined. Artifacts: build/port-world-geometry/c-{default,server,highres,production}
and c-audit.json. Client SHA:
6a2c43da3f8ae1383b40031c440215aa57c443dd2b855e97c8e32f119de8b731.

Production is **44,396 lines / 74 C files / zero reference C**, five lines below
the parent because unused quadrant locals were removed. Corrected selected bodies
total **1,385 lines**. The conversion below keeps twelve interfaces for remaining
C callers and moves Go callers to direct calls. floatToInt32 is reused; review
found ShapeBox.Calc has different intermediate rounding. The two force
coefficients still have C readers outside this batch and remain shared.

### Supplemental player wall association contract

After baseline **9a945523** was pushed, translation review identified the player
update-data wall pointer at offset 296. A further 120 cases cover all eleven wall
shapes plus absent walls, nil/non-nil update data, exact wall identity, guards and
projection succeeding before the force/contact response. Initially the fixture
used ClassMonster (2); the C bit is ClassPlayer (4). Correcting the fixture class
made the intended branch observable. No production correction was needed.

Two default processes and server/highres all pass and produce identical captures
(player-wall-{a,b,server,highres}.log and associated JSON under the batch directory).
The supplementary hash is frozen; existing 15 captures remain unchanged. Production
source is identical to 9a945523, so its fresh production qualification applies.
The next native selection has 544 candidate roots / 205 expected captures; focused
coverage is now 17 roots / 16 captures / 15,902 records. No checks remain active.

## Native conversion and arithmetic review

The 32 selected functions are replaced by seven Go implementation files. Twelve
thin exports remain for real C callers; twenty function interfaces, one private
threshold and one wall-span table are retired. The two force coefficients remain
shared with C readers outside this batch. All Go callers and fixture operations
invoke Go implementations directly. Source-reference audit finds only the unrelated
existing client/sight.go function with the same sub_427C80 name.

The compiled baseline uses x87 intermediates. Merely translating each C float
local into a Go float32 rounds some expressions too early. Focused comparisons
exposed this in box corners, map/rotated coordinates, edge projection, wall forces,
wall overlap ratios and gate contact geometry. Inspection of the qualified C
binary identifies the actual stores/reloads. Go uses wide intermediates and
explicit float32 boundaries at those points. The existing ShapeBox.Calc is not
interchangeable; it rounds its products earlier. floatToInt32 remains reused.

Native-focused-7 matches all **17 focused groups / 16 captures / 15,902 records**
in 0.352s, after the final review of normalization, diagonal projection, quadrant
differences and circle/box response widths. Native static-2 also passes.
No frozen expectations were changed. Full compiled-C disassembly and native
comparison diagnostics are under build/port-world-geometry. These are local
implementation diagnostics, not a second retained C test implementation.

A local comparison script previously printed “all available captures match” even
if compilation produced no captures. The test process itself correctly failed;
the script now also requires all sixteen expected groups. One intermediate compile
failure was a stale renamed local and was corrected before rerunning.

## Qualified native results

| Target | Root tests | Tests including subtests | Duration |
| --- | ---: | ---: | ---: |
| Default client | 544 | 42,760 | 251.18s |
| Server | 543 | 42,759 | 375.39s |
| High-resolution client | 544 | 42,760 | 290.70s |

There are no skipped tests. All **205 captures / 81,523 records** match the frozen
C expectations and are identical across targets. Starts and completions are
checked, not merely discovery. Frozen root tests and hashes remain unchanged.

Fresh production passes in **398.08s**: three ELF32/i386/SSE2/CGO binaries,
ABI/retained/retired interface checks, the exact known 1,553 asset-suite failures
(15 passing / 3 failing / 32 no-test packages), options gameplay, explicit
save/load and forced flat-map regeneration. References remain unchanged.

All four gates share an unchanged **2,198-file source manifest**; all sessions
are joined. Artifacts: build/port-world-geometry/native-{default,server,highres,production},
native-audit.json and native-interface-source-audit.json. Client SHA:
3b0754cf3af539d8447c294af8a64a48429441760495861838d69429233a7737.

C remaining: **42,982 / 74 files / zero reference**, **−1,414** from 7d363edc.
All installation/width-correction drafts are consumed. Original assets and archive
are untouched. The four prerequisite behavior corrections remain the review
items; native rounding changes preserve the corrected C behavior.
