# Map polygon lifecycle, lookup, events and serialization

## Scope and result

C baseline **afdb26e7** was qualified, committed and pushed before translation.
The native conversion replaces **33 reachable functions / 1,000 original C body
lines**: GAME1_1 420C40..422140 and GAME1_2 serialization 428CD0. Five production
Go files implement storage, lookup, actor events, ambient color and serialization.
Seven thin exports retain live C callers; 26 function interfaces and three private
C globals are retired. Go callers and fixture adapters invoke Go directly.

Production C is **44,401 physical lines / 74 files / zero reference C**, a reduction
of **1,072 lines**. Root test source and frozen expectations are unchanged from the
C baseline. No C algorithms were retained solely for tests. Geometry helpers still
serve real callers and are candidates for the following batch.

## Baseline and fixture design

Own complete mapped vertex/polygon/selection arrays, control/default bytes, live
counters and remap head. Free only fixture-created allocations, then restore old
pointer-bearing state. Borrowed-reset contracts detach allocations and return them
to the real cleanup afterward. Seed nonzero names, RGB and levels. Use real C
allocation, cryptfile owners, player/monster factories, script VM callbacks,
notification/audio queues and renderer light-color state.

All 16 focused C roots pass (0.239s). Eleven polygon captures were repeated in
separate processes before freezing; independent contracts cover sparse allocation,
slot reuse, defaults, nearest-vertex thresholds/ties, selection/construction/remap,
ordinary containment and ray degeneracy, edge projection/input mutation, cached
lookup, actor event order, secret statistics and exact notification bytes, audio,
repeat suppression, ambient color, cache initialization and borrowed ownership.
Nil-player and nil-monster calls pass on the qualified compiler, with no C fix.
Serialization covers versions 1..4, editor metadata, vertex-ID remapping, bounds,
secret-count state, unsupported versions, bypass and exact modern wire round trips.

Broader C default/server/highres pass **269/268/269 roots**, **34,988/34,987/34,988
tests including subtests**, no skips. **87 captures / 51,433 records** are identical
across targets. Durations: **143.87/225.80/157.58s**; unchanged **2,175-file source
manifest**. Production source was verified identical to parent 5ecbdfc2, so its
qualified production results were reused. No C prerequisite correction was needed.

## Native qualification

All 16 focused native roots pass (0.223s), with all eleven polygon captures matching
C. Static memory-access checks pass. Broader default/server/highres pass
**269/268/269 roots**, **34,988/34,987/34,988 tests including subtests**, no skips.
All **87 captures / 51,433 records** match C and each other. Durations:
**108.70/215.00/152.76s**.

Fresh production qualification passes in **371.15s**: three 386/SSE2/CGO builds,
ABI/export/retirement checks, the exact known **1,553 asset failure entries** and
package outcomes (**15 passing / 3 failing / 32 no-test**), headless options gameplay,
save/load and forced flat-map regeneration against existing references. All four
native gates share an unchanged **2,180-file source manifest**. All sessions joined.
Client SHA: `0a9413f9dfb01e79c28603eecc3bed74967b537e222fc069a629de873fd399f2`.

The selected corpus covers polygon/map fixtures (excluding the isolated prerequisite
diagnostic), minimap, object reports, quest runtime, world collisions, AI main/
monster state, object rendering, player controls and numeric conversion. These cover
actual lookup/audio/minimap callers, map serialization and shared actor/script/
render owners. One existing client-only root is absent on the server target.

## Compatibility and review items

- Containment rotates a finite ray among world corners (0 or 5888). Its segment
  helper includes both endpoints, so a shared vertex can be counted twice and an
  interior point reported outside. Preserve this phase-sensitive behavior; separate
  independent contracts cover ordinary geometry and this existing limitation.
- A remote player leaving every polygon retains its previous region. The host's
  zero client-cache path clears server cache without resetting its level. These
  existing branches are explicit contracts, not corrected during conversion.
- Nearest-vertex arithmetic uses a double temporary and comparison, a float X
  square, and a float threshold updated after each accepted vertex. The first draft
  misread the temporary as an integer; the large-coordinate contract caught it.
  Corrected against the original declaration without changing expectations. The
  earlier C-baseline report's arithmetic description was also corrected.
- Preserve the constructor's initial decimal ID string before copying defaults,
  including trailing bytes after short default strings.
- Keep libc allocation/free for compatibility with existing C-heap owners. Reset
  detaches borrowed arrays without freeing them. Compile-time assertions check
  record/vertex/remap sizes and the vertex-pointer offset.
- If editor metadata allocation fails, return nil instead of attempting to free
  the fixed polygon record in backing storage. This narrow reversible correction
  is recorded for review; forced allocation failure is not covered by the fixtures.
  Normal allocation and borrowed/reset ownership are covered.
- Editor names use real 256-byte metadata and runtime events use real script
  callbacks. Serialized secret counts update both current/previous counters;
  fixtures own and restore both.

## Artifacts and recovery

Manifests: map-polygons-batch.json, map-polygons-tests.txt and
map-polygons-scope.json. Local evidence: build/port-polygons/c-{default,server,highres},
native-{default,server,highres,production}, c-audit.json, native-audit.json,
native-interface-audit.json, focused-11.log, native-focused-2.log and static logs.
Use the combined 33-function extraction; core-only scope.py/candidate.json exclude
the serializer. All fixture drafts, native drafts and installation scripts are
**consumed**. Never recopy drafts or regenerate frozen expectations during recovery.

## Disk cleanup

Before this batch, twelve completed panel/configuration run copies were verified
against original assets. Removing 19,648 byte-identical files reclaimed 6.184 GiB.
Per-run deduplicated-assets.json manifests retain hashes/restoration metadata.
Original assets/archive, modified maps, saves, captures and binaries are preserved.
Never repeat deletion mode on those runs; use the recorded restore command if a
full run data copy is needed. The polygon qualification uses fresh run copies.
