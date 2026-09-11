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

Baseline log: build/port-roam-history/c-before.log. Production C is unchanged at
140,260 physical lines, 153 files, zero test-reference C. Fixture uses the existing
isolated server binding and minimal type-table hooks. No copied C reference body.
