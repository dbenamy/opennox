# Dynamic unit codes and extent lookup — 2026-09-10

The dynamic-code resolver (578B40) and extent lookup (4ED020) now use the existing
server.Objs.GetObjectByInd implementation. Unmarked dynamic codes return their
full raw32 input without accessing the server. Marked codes clear only bit 15,
lookup the remaining full extent, and return the matching object's full NetCode
or zero. Lookup skips destroyed objects and returns the first live match.

Both C signatures remain unchanged. Extent lookup returns the original raw32
object-address ABI for its two remaining decompiled C consumers. The dynamic
resolver directly uses typed Go objects, avoiding a round trip through that
address bridge. There is no truncation of high input bits to uint16.

Original-C baseline: `35bb2ce4`. The fixture uses C-allocated server.Object arrays
and installs their list on an isolated server proxy, avoiding Go pointers stored
in a C-visible list. It checks exact returned addresses by matching allocations,
full network-code results and the complete object bytes remaining unchanged.
It restores the server callback before freeing fixture memory.

Tests cover 163,840 unmarked low-word/upper-word combinations with a callback
that panics on server access, empty lists, marked hits/misses, destroyed-first
and duplicate live entries, zero/max network codes and full-width extents.
Another 400 deterministic randomized list traces check both the direct lookup
and dynamic resolver. Independent tests use raw uint32 equality, per-bit
arithmetic and a first-live-match model.

Production C: **141,786 physical lines (−35)** in 153 files; test-reference C: **0**.
Local artifacts: `build/port-extent/`.

All accumulated protection/network tests pass on 386 default/server/highres.
All three production binaries build and expose both C entries through Go-backed
bridges. The `extent-port` warrior scenario accepts both preserved screenshots
with overrides disabled. The latest full-suite milestone remains the RNG chunk,
which matched the known failure set exactly.
