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
before installation. No source is installed yet.

## Qualification plan

Reconstruct exact edits, format and review before installation. Require all 89
baseline names in each converted profile, safe/static checks, three fresh production
builds/ABI checks, exact known-suite results, headless creation/save/load/resume,
identical source fingerprints throughout and unchanged original asset hashes.
Measure selected cgo/export/header counts after qualification. Local artifacts:
`build/port-map-room-paint-owners/`.
