# Map catalog and rotation — next batch audit

Status: audit and ignored drafts only. The compressor's accumulated qualification
must finish before these source edits begin. No catalog C baseline is frozen yet.

Complete the connected map-list owner and rotation APIs: sorted insertion/free,
six-mode mapcycle parsing and indices, quest grouping, selection and played-map
history. The scope spans the maplist/rotation section of GAME3_2.c and
getSomeMapName/sub_4D10F0 in server__system__server.c, roughly 450 C lines. This is
a complete owner scope; do not add unrelated functions to meet a size quota.

Planned contracts use actual list ownership, RNG, startup tables and file readers:
bytewise sort and duplicate order, 0/1/128/129 catalog boundaries, six-character
quest families, missing/empty cycle files, 25-entry limits, first-line behavior,
mode masks, cursor wrap, RNG consumption and repeated selection/history sequences.
Read actual shipped maps through the metadata owner using a disposable mapcycle
file. Preserve original assets, global metadata, data path and handle ownership.

A valid small catalog can have zero eligible quest candidates after the family
and recent-history exclusions. Current IntClamp(0,-1) returns -1 without advancing
RNG, so C returns a string pointer before its row array. Proposed prerequisite:
fall back to a uniform choice from the full nonempty catalog when exclusions leave
no candidate. This is a bounded, reversible bug correction under the standing
user authorization. Implement and independently test it before freezing the C
baseline; document the deliberate fallback/RNG change and preserve all successful
candidate-selection behavior. No implementation has changed yet.

Compatibility details from the audit:

- The first non-section mapcycle line bypasses normal file/type/empty checks.
- The actual fgets facade consumes a complete line, copies at most 126 bytes and
  rejects an EOF-terminated tail. CR replacement precedes LF replacement, so an
  LF-before-CR input can mutate both positions.
- Next-cycle selection tests index > count, not >=. Preserve this defined behavior
  on addressable slots and record it for later review; do not infer a wrap fix.
- Quest selection uses raw map flag 2 mapped to game flag 0x1000. Case-insensitive
  grouping is bytewise and stops at NUL, not Unicode case folding.
- Preserve string-buffer tails and unsigned history/counter arithmetic.

Ignored drafts are under build/port-map-catalog, with AUDIT.md describing remaining
review and integration work. They have not been compiled or tested and must not
be treated as completed source after a session loss.

Remaining-C caller audit identifies five interfaces to retain as Go-backed C
exports: maplist first/next, cycle enable get/set, and quest-map selection. Other
interfaces in this scope can become direct Go calls and retire their C symbols.
Keep the existing list-head storage while remaining C owners access it.
