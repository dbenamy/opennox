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
