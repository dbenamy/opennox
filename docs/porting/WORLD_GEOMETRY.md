# World geometry and wall/circle/box collision

## Scope and current state

Parent polygon conversion **10d3294d** is qualified and pushed. This batch proposes
**32 reachable functions / 1,390 original C body lines**: segment/rectangle/shape
primitives, direction/vector conversion, map-coordinate transform, and the
54FFC0..551A90 wall/circle/box collision family. Read-only extraction and closure:
build/port-world-geometry/proposal.json and combined-reachability.json. All 32 are
live in the current whole-source reference/internal-call closure.

Production remains **44,401 C lines / 74 files / zero reference**. Fifteen repeated baseline
captures are frozen (15,782 records); no native conversion is installed yet. Four
prerequisite corrections are installed; corrected C qualification is complete.

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

## Coverage under preparation

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
availability will be checked by the runner. Static-2 passes.

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
total **1,385 lines**. Next: convert these 32 functions, keep the 12 interfaces
with remaining C callers, move Go callers to direct calls, reuse the existing
ShapeBox.Calc and floatToInt32 where their behavior matches, then repeat all gates
without changing frozen expectations. The two force coefficients still have C
readers outside this batch and must remain shared.
