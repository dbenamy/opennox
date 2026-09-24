# Map-room and painting fixture owners

## Scope and original baseline

Replace the numeric C fixture dispatchers with existing native map-room and
painting owners, then retire 49 room and 15 painting exports and their prototypes.
Production algorithms, frozen expectations and pointer normalization remain
unchanged. Preserve live spell/door transfer addresses, allocation helpers and
x87 control-word fixtures that still require C.

All 89 selected roots passed without skips in fresh default/server/highres
processes on qualified `b5d7ee16`, using source-verified test binaries. The exact
selection covers rooms, painting, hallways, population and map sections. It excludes
only `TestMapPopulationPrerequisiteProbe`, an opt-in historical diagnostic that
normally skips. See [baseline](map-room-paint-owners-baseline.json),
[selection](map-room-paint-owners-tests.txt) and [manifest](map-room-paint-owners-batch.json).

Whole-source symbol scanning found no candidate address registrations. The root
package's same-named wall helper is a separate Go implementation and remains.
The fixture's live transfer callbacks are outside the retired group; no snapshot-ID
reservation change is needed for these exports.

GPT-6 Luna drafts the isolated conversion and per-operation mapping. Primary review
checks every sparse opcode, native owner, argument order, signed/unsigned narrowing,
32-bit address word, float32 input bits, float64 output bits and int64 rounding
result. Review caught a pointer result cast and an incorrect room-argument type
before installation. The ten-file conversion is installed and fully qualified.

## Qualification

All 89 baseline names pass without skips in each converted profile. Safe/static
checks, three fresh production builds/ABI checks, exact known-suite results and
headless creation/save/load/resume pass. The known suite retains 304 failure events,
with 17 passing, two failing and 32 skipped packages. Source fingerprints match
throughout; all ten changed/deleted files match the accepted hashes. Frozen
expectations and all 1,654 original asset hashes remain unchanged.

Selected production cgo files fall 227→225 (238/463 eliminated on net); exports
fall 1,053→989 (901/1,890 retired). Two test cgo imports remain for independent
control-word/allocation helpers and live transfer addresses. Production C callback
bodies remain 78; headers remain 157 files, now 3,712 physical lines. Standalone
production/test-reference C lines remain zero. External native bindings remain.
See [qualification](map-room-paint-owners-qualification.json) and
[inventory](map-room-paint-owners-inventory-after.json).

Primary reconstructed every edit and checked all 64 sparse opcode mappings against
original C cases and owner names. Existing native cases and remaining static C
helpers are unchanged. Review corrected the pointer encoding for room overlap and
the room-type argument. A follow-up draft mistakenly encoded the room-type uint32
result as a pointer; primary corrected it locally and passed the grid configuration
pointer directly. All corrections preceded compilation. The larger Luna draft was
useful for bulk migration, but its generated casts still required independent
signature review. No assertions, normalizers or algorithms changed.
Local artifacts: `build/port-map-room-paint-owners/`.
