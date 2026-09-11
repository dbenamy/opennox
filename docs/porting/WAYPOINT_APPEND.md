# Waypoint link insertion — 2026-09-11

Scope: map-generator link insertion `sub_51D300` and the blob-kind wrapper
`sub_51D2C0`, both with live C callers in GAME4_1.c. Existing server.Waypoint
matches the C layout: 516 bytes, 32 eight-byte link slots at offset 92, and
one-byte active count at 476. Keep both ABI entries.

Insertion rejects counts >=31 (leaving physical slot 31 unused) and source==target.
It scans only active entries, rejecting the first same-pointer/same-kind match;
otherwise it writes the next pointer and kind and increments count. Nil target
is allowed. Pointer identity matters even for waypoints with equal Index fields.
Only the pointer, kind byte, and count may change; slot padding stays intact.

The compiled 386 duplicate comparison zero-extends the stored uint8 kind and
sign-extends the incoming C char. Kinds 128..255 therefore never match an existing
stored byte. Preserve this behavior rather than silently changing the classifier
to Go byte equality. Disassembly confirms `movzbl` versus `movsbl`. The wrapper
reads its kind from 0x973F18+35972; the direct call uses its explicit argument.

Local artifacts and original disassembly: build/port-waypoint-append.

## Baseline fixture

532,480 cases exercise both live C routes. Every incoming kind and count byte is
covered, with additional every-slot cases at counts 0/1/2/15/30/31/32/255. These
include exact and apparent high-bit duplicates, equal-ID/different-pointer
waypoints, nil targets, self-links, and matches outside the active count. The
direct and wrapper routes deliberately receive different explicit/blob kinds.

C-owned source/target storage has poisoned padding and distinct guards. Tests
check all 32 normalized pointer/kind pairs, exact count/return values, permitted
write offsets, unchanged targets/guards, and unchanged/restored blob state.
Layout assertions cover waypoint size and slot/count offsets plus link size and
pointer/kind offsets. All cases pass against original C before replacement.


Production C before this chunk: **140,879 physical lines**, 153 files, zero
reference C. No copied C reference implementation is needed.

## Native implementation

Both C entries now bridge to a private Go helper using the existing typed
Waypoint layout. It assigns pointer and kind fields separately, preserving
padding, and keeps the original signed comparison. The same 532,480 cases pass
after conversion. Original-C baseline: `23f32960`.

Production C: **140,845 physical lines** (−34), 153 files, zero reference
C. All accumulated default/server/highres 386 tests pass and all three ELF32
binaries build. Fresh `waypoint-append-port` gameplay passes both preserved
screenshots with overrides disabled. The full suite exactly matches the known
1,553 failure-entry multiset: 15 passing, 3 known failing, 32 skipped/no-test
packages. Artifacts: build/port-waypoint-append.
