# Hallway routing

Scope: six connected routines in GAME5.c, 54B2D0, 54B810, 54BA60, 54BB20,
54BD90 and 54BF20 (531 C section lines). Population conversion **07229da1**
is committed and pushed. Its native prefab connector is the only external
caller; all six C entry points can retire together.

## Corrected-C prerequisite

A positive northbound route from a 4×4 prefab at (0,0) to a 4×4 room at
(6,-12), with a one-cell entrance, failed in the original C path. Construction
writes the second corridor to blob offset 2491616, but ten connection reads use
the separate named dword_5d4594_2491616, which remains zero. The real room
connection helper receives a null second room and fails. Evidence:
build/port-map-hallways/bent-original.log; regression TestMapHallwaysBentConnection.

Correct those ten reads to the same blob array used for construction, positioning
and admission. This is a deliberate repair of an existing failure, not a claim
of identical behavior for the failing path. The regression covers four directions,
both lateral orientations and widths one/two, requiring admission and exactly
three live corridor records. Existing population captures must remain unchanged.
Validation passes with unchanged population/painting/room captures in all three
variants: default 23.869s test time; server 102.602s and highres 34.029s wall
time including their builds. Evidence: prerequisite-{default,server,highres}.log
and prerequisite-variants.json. Physical C remains 104,839
lines / 148 files / zero reference C; no conversion has occurred in this batch.

## Conversion plan

Extend complete captures through actual routing with straight, bent and
three-segment routes, rectangular target rooms, entrance widths, close/large gaps,
obstructions, fallback candidates and exact connected topology. Record discarded
corridors before release without reading freed storage. Use direct helper tests
where entry-point coverage leaves a real branch gap. Keep all 36 population
hashes unchanged. Repeat and lock corrected C captures in default/server/highres;
commit/push that baseline before replacing the algorithms with Go.

Then compare every captured byte, run accumulated three-variant tests, build and
audit all production binaries, compare the asset-backed known-failure suite and
run fresh unchanged headless repeat-a gameplay. Update C_LOC and recovery docs,
commit/push, summarize and continue to the next connected batch.

## Expanded C baseline — in progress

Prerequisite **9ec45f48** is committed/pushed. The corpus now contains **2,512
captured cases / three groups**: 2,376 routes, 72 obstacles and 64 candidate
fallback cases, plus 16 explicit bent-route contracts. Routes span gaps 1,2,3,4,8,12,
three widths, eleven lateral offsets, four directions and three rectangular
room shapes. Route counts are 230 rejections, 568 single-segment, 241 two-segment
and 1,337 three-segment admissions. Topology checks require live reciprocal edges
and reachability of the intended target. Barriers require rejection and no live
temporary corridor; fallback skips an undersized candidate in either list order.

The stronger contracts exposed fixture issues, corrected before baseline lock:
initialize/restore the real startup opposite-direction table (1,0,3,2), avoid
shared mutable room maps, and place the undersized candidate off the route to
its fallback. Existing 36 population hashes retain their original inputs and
remain mandatory. The release observer now records discarded corridor identities
before free, without reading released memory. Scope this to the connector call.

The close-gap matrix captures **491 zero/negative-length live corridors** produced
by existing routing arithmetic. Preserve this defined behavior during conversion;
changing admission and route selection is a separate gameplay correction for
review. Evidence: degenerate-corridors.json. No claim is made that these shapes
are desirable. The existing direct room allocator already permits signed lengths.

Current evidence: c-boundaries.log (routes/obstacles pass), fallback-corrected.log
(64 cases pass), c-boundaries-{routes,obstructions,fallback}.json. Independent
repeat and server/highres checks run via repeat-baseline.py; lock mandatory hashes
and commit/push only after all complete captures repeat exactly. Early exploratory
captures are losslessly gzipped with checksums in early-capture-archive.json.
Do not reapply staged corpus/topology/fallback files; the working tests include
subsequent corrections.

## Repeated corrected-C captures qualified

All **2,512 cases / three groups** match byte-for-byte in independent default,
server and highres, with existing population/painting/room tests passing:
**42.316s / 119.848s / 52.549s** wall time including builds. See
baseline-qualification.json, baseline-variants.json and c-{repeat,server,highres}.log.
Every capture hash is mandatory in src/map_hallways_baseline_porttest_test.go.
The final mandatory-hash smoke passes in **19.890s** (locked-baseline.log).
The Go draft is staged separately and production still uses the six C routines.

Reproduce after sourcing build/baseline/env.sh, from src:

```sh
GOMAXPROCS=2 go test -p 2 -tags porttest -count=1 -run '^TestMap(Hallways|Population|Painting|Room)' .
GOMAXPROCS=2 go test -p 2 -tags porttest,server -count=1 -run '^TestMap(Hallways|Population|Painting|Room)' .
GOMAXPROCS=2 go test -p 2 -tags porttest,highres -count=1 -run '^TestMap(Hallways|Population|Painting|Room)' .
```

Optional OPENNOX_MAP_HALLWAYS_CAPTURE is an absolute output filename prefix;
mandatory comparisons stay enabled. Duplicate baseline snapshots are gzipped
with checksums in baseline-capture-archive.json. Older population completed
captures and the early hallway topology run were also losslessly compressed to
manage disk use; see completed-capture-archive.json under build/port-map-hallways.
Decompress named evidence before rerunning an older script requiring raw JSON.

## Native implementation — qualification pending

Corrected-C baseline **45961868** is committed/pushed; mandatory hashes pass in
19.890s. All six algorithms are now native in map_hallways.go. Their C entry
points retire completely; native population calls directly into Go. The C bodies,
six prototypes and one unused extern are removed. Current count: **104,307 lines /
148 files / zero reference C**, −532 physical lines. Zero retired-name references
remain in source. Native-first.log is the focused comparison; qualify-all.py
then runs broad validation. Do not reapply the staged native draft.

All **2,512 cases / three captures** match on the first native run, with unchanged
mandatory hashes. The 16 bent contracts plus existing population/painting/room
checks pass in **41.914s**. Complete snapshots were losslessly archived after
comparison; see native-comparison.json and qualified-capture-archive.json.
Broad qualification is running under qualify-all.py. Next scope is the 33-routine
theme parser, documented locally under build/port-map-theme; no theme source
changes are included in the hallway batch.

Accumulated **64,311 cases / 950 capture groups**, plus the room, painting and
16 bent-hallway contracts, pass in default/server/highres: **266.526s / 342.367s /
281.183s** wall time. Production builds and final full-suite/gameplay checks
remain in progress; see variants.json and qualify-all.log.

## Completed native qualification — 2026-09-14

All six routines are native and all six old entry points retire. **104,307 physical
production C lines / 148 files / zero reference C**, a reduction of **532 lines**.
All 2,512 cases / three complete captures match baseline **45961868** unchanged
on the first native run. Accumulated 64,311 cases / 950 groups and the additional
contracts pass in all three variants (266.526s / 342.367s / 281.183s). All three
production binaries pass architecture and retired-symbol verification.
The asset-backed full suite matches the known failure multiset exactly: 1,553
entries, 15 pass / 3 fail / 32 skip packages. Fresh unchanged repeat-a gameplay
passes in **36.527s** with Xvfb/null audio. See qualification.json and the
map-hallways-port run artifacts. This supersedes intermediate pending status.
Next connected batch: 33 theme-parser routines, with source and plan under
build/port-map-theme. No theme source changes are part of this conversion.
