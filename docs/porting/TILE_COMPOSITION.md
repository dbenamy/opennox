# Tile composition, overlays and scrolling

Status: **qualified Go conversion**. Frozen C baseline: **6819d3ff**.
Removed **579 C lines**; **73,450 C lines /91 files /zero reference C** remain.

Seven routines in six marked blocks total **558 C lines**: counter reset, edge
overlay callback and linked traversal, horizontal/vertical incremental redraw,
redraw predicate and full tile composition. The private no-op edge callback and
both callback slots become unnecessary when their final callers move to Go.

## Ownership and contracts

Fixtures reuse the actual client tile buffer, shape tables, image bag, configuration
and callback selection from tile-raster tests. Map cells come from the production C
grid allocator. Cleanup frees rows through the actual row-free owner, then frees its
pointer table and restores the previous grid. Definition/image arrays and linked
edge nodes use C-compatible allocations, detached before freeing. Fixtures save and
restore actual usage tables, shared words, counter, definitions and callbacks.

- Full redraw: both callback modes, empty and asymmetric populated grids, variation
  and definition offsets, rounding/clamps near map origin and the upper grid edge;
  captures entire ring pixels, usage flags and state. Independent contracts check
  scroll/counter reset, empty-grid flags and unchanged grid/viewport inputs.
- Incremental redraw: horizontal/vertical boundaries, large jumps, no-op paths and
  map limits, plus repeated scrolling around the ring in all four directions.
  Captures pixels and state at every step; independently checks origin transitions.
- Redraw predicate: either tile half, enable and definition flags, each inclusion/
  exclusion edge, signed Y and unsigned X near-origin behavior. Independent expected
  selection and complete state immutability checks supplement frozen results.
- Overlays: actual linked edge dispatch, transparent/literal/tile-source runs,
  partial/full first and last rows, nil image, texture/flat modes and ring splits.
  Independent per-byte application checks the entire buffer; distinct tile sources
  make linked overlay order observable. Usage flags are checked independently.

## Compatibility decisions

The redraw predicate computes X from uint32 viewport words, so values below 11
wrap before division and clamp at the upper map edge. Full redraw first stores
the delta in signed int. Preserve this difference rather than silently sharing a
signed bounds calculation. Both use signed Y. The near-origin cases must establish
this behavior against original C before accepting the translation.

Counter reset has a Go server caller; retain its Go wrapper and actual shared
counter storage. All tile/edge callback consumers belong to this composition batch.
Move selection/state to Go with the callers, retiring private exports and C slots.
Actual-owner fixtures follow the new owner after translation, with unchanged
frozen behavioral expectations. No C algorithm remains solely for testing.

The root draw orchestration and final screen lighting/copy are already Go. Use a
complete accumulated corpus for the native floor-composition milestone, with an
explicit 900-second per-package timeout: the prior highres milestone took 598.823s
under the600-second default. The runner's new optional timeout leaves defaults
unchanged and records its value; nonpositive inputs fail before invoking Go.
Also qualify all production builds/ABI, exact full-assets known failures, normal 12
and GUI flat-floor 14 gameplay frames. Solo replay does not establish multiplayer coverage.
Malformed image streams remain outside the valid asset contract.

Evidence/drafts live under build/port-tile-composition. Applied stages are stale;
follow actual source and PORTING_STATE.md. Original assets/archive stay unchanged.

## Development findings

Development-a joined with failure in 187.296s before any composition behavior ran:
the fixture attempted to assign the grid capacity, which is actually a const 128
in production C despite mutable extern declarations elsewhere. The corrected fixture
only reads/asserts that capacity and restores the mutable grid pointer. Production
is unchanged; no algorithm or expectation correction is involved.

Development-b completed all five roots in 188.089s: four passed, while the overlay
pixel contract caught an owner mistake. Creating a second temporary image bag hid
the tile-source handles, so the real loader correctly returned nil and drew nothing.
The corrected fixture installs tile and overlay images in the same actual bag.
It also adds full-composition cases with edges attached to either cell half.
Production remains unchanged. Earlier successful group hashes are retained.

Development-c stopped in 22.765s at overlay origin(275,229). That arbitrary X came
from raster tests, but the legacy overlay path requires the row's diamond start
not to lie beyond the ring end; production composition's maximum relative X is
width−23 (253 here). The fixture now uses(253,229), testing an exact-end start and
split without violating that caller precondition. The retained crash identifies
the invalid-input limit; it is not a qualified behavior or a porting correction.

Development-d stopped in 22.381s at(31,450) with a late first overlay row: the legacy
setup performs only one initial ring wrap after adding that row, so this combination
starts beyond a second complete ring. Use(31,400) for the one-wrap matrix, retaining
coverage of all first/last rows. Keep the invalid-input evidence and document the
one-wrap precondition rather than treating arbitrary raster-helper inputs as valid
for every overlay caller. No production behavior or frozen result was changed.

Initial C qualification passed 265 default roots in 138.988s, then server discovery
failed in 83.730s: the public Go wrappers in legacy/client_draw.go are client-only.
The test adapter now calls the same production C routines directly under porttest
in all targets. Production is unchanged; expectations stay frozen. Repeat all three
affected targets with the corrected adapter, preserving the initial artifacts.

Pixel-row ownership audit correction: no remaining C function reads
nox_pixbuffer_rows_3798784. Its storage is still C, but the live reader is the Go
wall-edge renderer and the setter/allocation path is already Go. Earlier notes
about waiting for C tile consumers were inaccurate. Consolidating that storage
is a separate ownership cleanup; this composition batch leaves it unchanged.

## Frozen C qualification

Development-e completed all six roots in **24.492s**. All **3,096 records** are
frozen; the four groups that already passed in development-b remain byte-for-byte
identical. No production algorithm or frozen expectation was corrected.

| Group | Records | SHA-256 |
| --- | ---: | --- |
| attached-edges | 32 | `e4e2ef4f01efc47182b48d8fa0aea2017a80aaade63599b5f63bceaf2f40006f` |
| full | 88 | `504daeb5c996ed08ddb9f5df460b449c0589b3645baf67beef0d9b235fe62d82` |
| overlays | 1440 | `6fd9c980bfebc961818c9916ec2ca6e14bbb15b87985b45d23375a18b3c232ba` |
| predicate | 1152 | `0d23c8d5d142fff19f2ab5df3844d0a30fd7663b18530c0464d1c27006836064` |
| scroll-sequence | 64 | `b6014d8bee88c4473779b42fb6cb38d9de4d04c2ac1e462e2ccc6f1bb7824952` |
| scroll | 320 | `f084ce81b75b785fa1cf2340951eadaf101c926d2b2f1f6ffeaad189f1f3e8aa` |

Affected default/server/highres checks completed **265 / 263 / 265 roots**
in **226.828 / 220.753 / 146.716s**. Every selected root started and finished, with the expected
optional prerequisite skip. All frozen hashes match in every target. Source
fingerprints remain unchanged. Production is identical to **fb882b61**, allowing
reuse of its three builds/ABI, exact known 1,553 full-suite failure entries and
fresh normal 12+GUI flat-floor 14 gameplay. The native conversion requires new evidence.
Evidence: build/port-tile-composition/c-qualification.json. No composition algorithm
has been replaced in this baseline.

Review added one independent counter/no-op contract: actual scrolling resets the
shared counter from 17 to −1, while a no-op preserves every renderer word and pixel.
The complete focused family then passed 7 roots in all three targets, with all frozen
expectations unchanged. All earlier source fingerprints are identical, validating
reuse of the broader affected runs; the added test is the only source change.

## Go conversion and complete rendering milestone

The first Go focused run passed all seven roots and 3,096 frozen records in
**188.594s**, without an algorithm or expectation correction. Removing 558
function-block lines, one private no-op line, two callback-slot definitions and
18 unused extern declarations reduces C by 579 lines.

Go now owns tile callback selection, full and incremental composition, redraw
selection, edge-list traversal and overlay copies. Actual grid/definition/buffer
and counter owners remain. Private callbacks are called directly in Go; twelve
C interfaces/storage slots retire, including the prior raster exports and private
no-op. A later caller audit also found the width/height C getter adapters unused;
retire those in the following asset batch, retaining their live Go owners. Frozen expectations remain unchanged; no test-only C algorithm remains.
The unused C scratch returns and upper-word callback argument residue are absent
from private Go APIs. Actual callers never consumed them.

The **complete accumulated port corpus** completed **1023 / 1019 / 1023 roots**
in **583.533 / 668.732 / 600.907s** for default/server/highres, with an explicit 900-second
per-package timeout. Every selected root started and finished; only established
optional prerequisite skips remain. All 3,096 composition records match frozen C in
each target, alongside the earlier raster/edge/wall/UI/gameplay port contracts.
All three production builds pass ELF32/i386/SSE2/CGO and retired-interface checks,
with no test helpers. Full-assets failures/package outcomes remain exactly 1,553
entries and 15 pass/3 fail/32 skip. Fresh normal chapter/minimap 12 and GUI flat-floor 14
frames match, with updates disabled. Source fingerprints stay unchanged throughout.
Evidence: build/port-tile-composition/native-qualification.json. Invalid image
streams and out-of-contract overlay coordinates remain outside this qualification;
solo replay does not establish remote multiplayer coverage.
