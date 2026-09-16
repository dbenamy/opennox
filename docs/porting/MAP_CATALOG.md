# Map catalog and rotation

C baseline **cd5337b6** qualified under `build/port-map-catalog/c-qualified`.
All eight roots pass in default, repeat, server and highres (240.234s total),
with unchanged source and no skips. Frozen expectations cover 1,789 records in
eight groups, including all 50 shipped map names through the real metadata reader.
The three-line prerequisite brings working production C to 69,345 lines in 88
files; the most recent completed conversion remains 69,342, zero reference C.
Native replacement is fully qualified.

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

Native focused checks match all eight frozen groups (`native-focused`, 105.200s
including compilation). Final review preserved safe null-list traversal and map
metadata error logging; the former has an additional independent ABI contract.
The actual list head and C-compatible allocations remain shared with live C users.
Go wrappers now call Go directly; five C-facing exports remain, fourteen internal
C symbols retire. The two removed blocks total 448 C lines, leaving 68,897 in 88
files, zero reference C (net −445 including the prerequisite).

Final gate: `build/port-map-catalog/native-qualified`, 662.181s with unchanged
source. All 47 affected roots pass without skips in default/server/highres
(158.790s / 158.547s / 76.836s including discovery/build). All three production
binaries are ELF32/i386/SSE2/CGO and pass retained/retired symbol checks. The full
asset suite exactly matches the known 1,553 failure entries and package outcomes
(15 pass / 3 known fail / 32 skip); this is not an all-green suite.

Fresh normal and flat gameplay runs (`map-catalog-native`,
`map-catalog-flat-native`) both exit zero, regenerate the warrior map exactly and
match their preserved frame references (50.366s / 50.137s). All readers joined.
No frozen expectation changed during conversion. The previous codec milestone
already ran the entire accumulated corpus; this bounded batch used its 47 known
affected roots. The next shared-list change will warrant the full corpus again.

Old derived compiler cache entries were pruned while no readers were active;
`cache-prune.json` records all 1,240 paths and 26,341,059,221 reclaimed bytes.
Original assets, archive and validation evidence were preserved. Ignored
apply/review scripts have been applied and are stale; do not rerun them.

Remaining-C caller audit identifies five interfaces to retain as Go-backed C
exports: maplist first/next, cycle enable get/set, and quest-map selection. Other
interfaces in this scope can become direct Go calls and retire their C symbols.
Keep the existing list-head storage while remaining C owners access it.
