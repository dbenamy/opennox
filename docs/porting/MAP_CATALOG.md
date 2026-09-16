# Map catalog and rotation

Status: C baseline qualified under `build/port-map-catalog/c-qualified`.
All eight roots pass in default, repeat, server and highres (240.234s total),
with unchanged source and no skips. Frozen expectations cover 1,789 records in
eight groups, including all 50 shipped map names through the real metadata reader.
The three-line prerequisite brings working production C to 69,345 lines in 88
files; the most recent completed conversion remains 69,342, zero reference C.
Native integration is next.

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
and recent-history exclusions. The original IntClamp(0,-1) returns -1 without advancing
RNG, so C returns a string pointer before its row array. Implemented prerequisite before freezing C:
fall back to a uniform choice from the full nonempty catalog when exclusions leave
no candidate. This is a bounded, reversible bug correction under the standing
user authorization. Independent contracts require a valid selected name and exactly one RNG draw for
each choice from a catalog larger than one. Successful candidate-selection
behavior is preserved. The C baseline commit is the recoverable oracle.

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

Coverage includes bytewise sorting/duplicate order through 180 nodes; quest
capacity 127/128/129; six-character families; 1,050 seeded selection/history steps
including counter wrap; all cycle mode combinations and signed cursor boundaries;
first-line, NUL, CR/LF, EOF and 124–1,024-byte line cases. Fixtures restore list,
quest/cycle state, metadata, startup pointers and path, and check file-handle
cleanup. Frozen hashes are in `map-catalog-captures.json`.

Ignored `native.go.draft` under build/port-map-catalog still needs integration and
review. Other fixture drafts are stale: use the committed source fixtures.

Remaining-C caller audit identifies five interfaces to retain as Go-backed C
exports: maplist first/next, cycle enable get/set, and quest-map selection. Other
interfaces in this scope can become direct Go calls and retire their C symbols.
Keep the existing list-head storage while remaining C owners access it.
