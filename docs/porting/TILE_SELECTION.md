# Tile selection state — 2026-09-11

Scope: tile name/image/variation selection and the boolean tile flag setter
51D4D0/51D540/51D570/51D5C0. All four retain live C callers. Larger map-placement
consumers remain with their owner.

Name lookup scans all 176 physical C entries, including entries beyond
nox_tile_def_cnt. It uses the shared nox_strcmpi byte/locale comparator, lets the
last matching entry win, then gives NONE precedence with selected index 255.
A miss resets selected index to zero. Preserve the existing string helper as a
production dependency during this conversion; Unicode case folding is different.

Image selection accepts signed values 0..175 and resets selection to zero for
other values. Variation requires a valid selected index and accepts any signed
value <= width*height−1, including negatives. Failed variation resets variation
to zero. Existing C callers select an image immediately before setting variation;
an invalid image selection resets to tile zero. Calling variation after NONE
would index outside the C array and has no defined compatibility requirement.
The flag setter accepts exactly 0 or 1 and preserves state on other inputs.

Shared state is the full C-backed tile array, selected index at 0x973F18+35912,
flag at +35916, and C variable dword_5d4594_3835348 for variation. Tests must
restore the raw table, active count, all state, and adjacent guard words.

## Original-C baseline

The installed fixture passes 459,124 checks against original C: 12 name cases,
180 image indices, 180 flag inputs, and 458,752 variation cases (all 256×256
width/height pairs with seven signed boundary values). Name cases cover last
case-insensitive duplicate, a table entry named NONE overridden by the sentinel,
empty/high-byte/mixed-case names, positive 31-byte matches, longer-prefix misses,
and embedded NULs. The active count is deliberately smaller than the physical
lookup table.

The fixture checks exact full-table bytes, read-only input names, active count,
unrelated state and adjacent guard words. Independent before/after snapshots
verify complete restoration after each route. Layout assertions cover the
60-byte TileDef and its name/dimension offsets. No copied C implementation or
assets are used. Local artifacts: build/port-tile-selection.

Production C before this chunk:
**140,845 physical lines**, 153 files, zero reference C.

## Native implementation

All four C entries now bridge to private Go state helpers. The same 459,124
checks pass after conversion. Name comparison continues to call the shared
production C nox_strcmpi; no new C reference implementation is retained.
Original-C baseline: `375fc1b3`. Production C: **140,785 physical lines**
(−60), 153 files, zero reference C. Accumulated default/server/highres 386 tests
pass, all three ELF32 binaries build, and fresh `tile-selection-port` headless
gameplay passes both preserved screenshots with overrides disabled. The
immediately preceding waypoint chunk ran the full suite with the exact known
1,553 failure entries unchanged; this small chunk does not repeat that run.
