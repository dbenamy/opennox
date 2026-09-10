# Waypoint allocation, links and predicates — 2026-09-10

The two next-link helpers (579870/5798A0), raw waypoint allocator (579E70), and
composite enabled/mask predicate (547EE0) now execute Go. They retain existing
32-bit pointer ABIs, nil-safe next-link results, and the typed Waypoint layout.
Allocation still calls raw calloc for 516 bytes, returns nil on failure, and
sets only Flags bit 0x01000000. Using alloc.New would change failure behavior to
a panic, so it is intentionally not used here.

The composite predicate short-circuits nil and disabled waypoints before calling
the existing typed HasFlag2Mask method. Its sole-use C mask helper (579EE0) and
declaration are removed. Remaining C consumers keep their four bridge symbols;
no C implementation is retained solely for tests.

Original-C baseline: `30d5ce15`. Exhaustive Flags2/mask byte pairs across seven
main-flag patterns plus nil cases produce **459,008 predicate cases**. Independent
per-bit arithmetic provides the mask oracle. All complete waypoint bytes remain
unchanged. Both next entries cover the 16 combinations of nil/three source nodes
and nil/three destinations, including self-links; returned addresses must match
the exact C allocations and previous links must not be mistaken for next links.
256 raw allocation snapshots assert all 516 bytes are zero except byte483=1.
Fixture memory is poisoned before freeing to exercise calloc zeroing on reuse.
OOM injection is not added; the explicit nil-return branch remains in production.

Production C: **141,745 physical lines (−41)** in 153 files; test-reference C: **0**.
Local artifacts: `build/port-waypoint/`.

All accumulated protection/network/waypoint tests pass on 386 default/server/
highres. All three binaries build; symbol checks show the four live entries
are Go-backed and the retired mask symbol is absent. `waypoint-port` accepts
both preserved gameplay screenshots with overrides disabled. The full suite
exactly matches the RNG milestone: 15 passing, 3 known failing, 32 skipped/no-test
packages and no changed failure entries.
