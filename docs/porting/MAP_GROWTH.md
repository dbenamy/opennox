# Procedural map growth

Seven connected routines from GAME3_2.c now use Go: initial normal/ring layouts,
frontier traversal, recursive dispatch, room expansion, hallway expansion, branch
selection, and doors with exclusion/waypoint connections.

The corrected C baseline is `eca05817`. Native source is in
`src/legacy/map_growth_{layout,expand,doors,exports}.go`. Three C callers retain
bridges (`nox_xxx_mapgen_Doors_4D4790`, `nox_xxx_mapGenMkSmallRoom_4D4F40`,
`sub_4D52F0`). Four private helpers and their header declarations are retired;
source and binary audits find no remaining references. No C algorithms remain
solely for testing. Shared records retain their existing allocator ownership.

## C corrections established before conversion

`2fc007db` fixes the four direction choices being indexed through separate scalar
locals: the stronger original probe observed an invalid hallway kind before its
release. The same prerequisite switches four merge-rate reads from an unrelated
named global to the configuration word actually written by the theme reader.
A positive merge case failed with configured rate 100 and named word 0; 192 additional
direction/rate/seed cases prove independence from that unrelated word.

`eca05817` fixes door scanning reading integer GridX/GridY words as floats. The
original translated-room case 6 created no door or waypoints instead of one/two.
Sixteen translated cases now assert counts, full state, and projected door-position
shifts. All eight previous capture hashes are unchanged by this correction.

These corrections were qualified in C before native conversion. Original failure
logs and source snapshots are under `build/port-map-growth`. The fixture also
corrected its headless flag from `1<<22` to the existing `GameFlag22` constant;
that was a fixture correction, with no production rendering change.

## Locked captures and contracts

| Group | Cases |
| --- | ---: |
| Branch choices, including unsigned dimension comparisons | 3,528 |
| Initial normal and ring layouts | 168 |
| Recursion limits | 100 |
| Room expansion | 128 |
| Hallway expansion | 144 |
| Doors and adjacent-room geometry | 1,216 |
| Translated door positions | 16 |
| Empty, generated and repeated frontiers | 97 |
| Obstructions and configured merges | 384 |
| **Complete captured cases** | **5,781** |

The 335 separate contracts cover bootstrap, direction selection, rejected hallway
contents, configured merges and positive door creation. The 16 translated-door
contracts are already included in captured cases above. Captures include complete
registered room/exclusion records, discarded candidate contents before release,
objects, waypoints and links, tile/wall state, globals, and RNG tails. Guard regions
and floating-point control state are checked. Every hash is mandatory; the temporary
baseline-extension bypass is removed.

Corrected C default-repeat/server/highres runs reproduced all nine hashes with
prior map checks in 82.945/161.009/93.334s wall. Locked smoke passed in 25.113s.
The first complete native run matched every hash in 25.082s, plus all contracts.
An initial compile-only fixture array-size mismatch was corrected beforehand.

## Final qualification

Accumulated **73,484 captured cases / 978 groups**, plus contracts, pass in all
three variants: default **393.463s**, server **388.575s**, highres **320.343s** wall.
All production builds pass (68.687 / 9.612 / 67.813s for normal / HD / server).
The binaries are ELF32/i386 with SSE2 and CGO; all four retired C symbols and
fixture helpers are absent, and all three required C bridges are present.

The asset-backed full-suite failure multiset is unchanged: **1,553 entries**,
**15 pass / 3 fail / 32 skip packages**, with no added or removed failures.
Fresh unchanged repeat-a gameplay under Xvfb/null audio passes in **37.042s**.
No expected gameplay artifacts or C capture hashes were updated for the port.

The conversion removes **1,081 physical C lines**: 1,079 section lines and two
obsolete declarations. Remaining production C is **101,335 lines / 148 files**;
reference C is **zero**. See [C_LOC.md](C_LOC.md) for reproducible checkpoints.

## Evidence and recovery

Under `build/port-map-growth`:

- `baseline-captures.json`, `baseline-variants.json`, `baseline-qualification.json`:
  repeated corrected-C hashes and locked smoke.
- `native-comparison.json`, `variants.json`, `builds.json`,
  `binary-verification.json`, `full-suite-comparison.json`, `qualification.json`:
  native comparison and final checks.
- `retired.json`, `abi-audit.json`, `retired-source-audit.json`: boundary audit.
- `baseline-capture-archive.json`, `early-capture-archive.json`,
  `native-first-archive.json`: complete gzip captures, byte counts and hashes.
- `fill-directions-original.log`, `merge-original.log`,
  `door-translations-original.log`: focused original-C failures.

The mandatory hashes and fixture source are committed; full captures and raw logs
are local evidence. The corrected C implementation is recoverable from Git at
`eca05817`. Build scripts run with `source build/baseline/env.sh`, `GOMAXPROCS=2`
and Go `-p 2`; asset-backed tests require the extracted Nox data directory.
The unchanged gameplay result is in `build/baseline/runs/map-growth-port`.

The next candidate is the five remaining generator orchestration routines.
Their baseline should include real generation through a small synthetic theme,
plus isolated retry/filesystem/ownership checks. Existing warrior gameplay does
not by itself validate procedural map generation. Read-only inventory and
unapplied fixture/probe drafts are in `build/port-map-orchestration`.
