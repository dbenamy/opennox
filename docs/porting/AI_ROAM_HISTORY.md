# Roaming history and successor selection — 2026-09-11

Scope: 545790 start, 5457C0 cancel, 545B00 insert, 545B60 reverse search,
545BB0 dead-end handling and 545C60 successor selection. The larger roam update
5457E0 remains C for this batch. History insert/search/dead-end entries retain
live C callers; successor selection is private to dead-end handling.

Original-C baseline: 98,304 cases, six operations over 16,384 deterministic
scenarios. One guarded C-owned object, full monster data and 34 waypoint records
are reused. Tests cover every ring index, all byte masks, neighbor counts 0–32,
all stack indices 0–23, nil and duplicate entries, mixed eligibility, unvisited
selection and history-only fallback. An independent ID-based model checks
history, selected pointer, stack, fields and both RNG indices. A SHA-256 baseline
also captures every changed object/monster word after normalizing pointer IDs;
all waypoint storage and object/monster guards must remain unchanged.

Preserve these details: start clears slot zero only; cancel clears Args[0] only;
reverse search excludes the current slot; forward fallback includes all 16 slots;
duplicate outgoing edges retain their random-selection weight. A one-candidate
selection uses the real random helper. A zero-neighbor dead end preserves Args[0]
before the normal pop/reset, while a nonempty unsuccessful search writes nil.
Dead-end failure still performs type-name lookup and real action pop. Debug logging
is disabled in this fixture; the final implementation must retain its message.

Baseline commit: 9c70753d; log: build/port-roam-history/c-before.log. At baseline C was
140,260 physical lines, 153 files, zero test-reference C. Fixture uses the existing
isolated server binding and minimal type-table hooks. No copied C reference body.

Native conversion: all six bodies removed; start/cancel now use the real native
registration, three C entries remain for the C update owner, and successor
selection is private Go. All 98,304 baseline hashes/model checks pass after
conversion. The dead-end message uses the existing Go AI logger, with explicit
signed 32-bit formatting for the original C percent-d fields. The C logging
scratch buffer is no longer needed by this native message.

Production C: **140,082 physical lines (minus 178)**, 153 files, zero reference C.
Accumulated default/server/highres tests pass and all three production binaries
build as ELF32/80386 with GO386=sse2. The previous
AI batch already qualified the exact full-suite baseline under SSE2; this batch
changes no shared toolchain or test infrastructure. Fresh roam-history-port
gameplay exits 0 against both preserved screenshots with overrides disabled.
No isolated speedup is claimed: the remaining C update still pays exports for
history helpers. Benchmark the full roaming update when its caller is moved,
rather than extrapolating from fixture timings dominated by state snapshots.

Next owner audit: 5457E0 uses an existing native retaliation method (545E60),
nearest-waypoint lookup (518460), detailed-path callback and fallback-waypoint
callback. Actual movement (50D3B0), sound investigation (5466F0) and move audio
remain C. Cover attack interruption, timed idle pushes, invalid/absent waypoint,
arrival distance, path failure and completion before retiring the owner. Existing
Go aggression predicates compare float32 against rounded constants whereas the C
predicates compare against double literals; check threshold neighbors against
actual compiled C before reusing those helpers. No behavior change is approved
merely by this audit. Do not let the next batch become a pathfinder rewrite.

The follow-up main-update port now retires all three history exports; see
[AI_ROAM_UPDATE.md](AI_ROAM_UPDATE.md) for its separate baseline and qualification.
