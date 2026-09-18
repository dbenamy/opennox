# Session lifecycle and map entry

Repaired C baseline in progress after qualified player-state conversion **50c06879**.
Selection: **33 live routines and one orphan**, originally749 body lines across
GAME3_2.c, map-name/mode helpers in GAME1.c, and saved-character metadata loading in
GAME1_1.c. Current production C is **32,608 physical lines /69 files /zero reference C**.

The whole-source caller audit includes C conditionals, Go wrappers, headers and
embedded C preambles. Eight entrypoints must remain for outside C callers;26
interfaces can retire. The C `nox_xxx_mapGetTypeMB_4CFFA0` is orphaned; the similarly
named root Go function is already native. Two existing variadic text wrappers in
GAME3_2.c remain compatibility glue. See [selection](session-entry-selection.json).

## Contracts and owners

Twenty new test roots produce21 captures /1,070 records. The completed
focused corpus also includes the existing console-script fixture:21 roots /224 tests
including subtests. Additional loops check every16-bit input against three mode
tables,512 map-flag combinations and all16,384 tile cells. Case counts do not
replace the independent assertions below.

- Save slots: AUTOSAVE and SAVE0001..SAVE0013, sparse/gapped/full sets, each slot,
  and ignored SAVE0000/SAVE0014 directories; actual isolated filesystem.
- Character files and metadata: actual encrypted-file writer/reader, independent
  version10..13 layout, unknown-section skipping, missing/empty/rejected files,
  global restoration, output padding and working-directory behavior.
- Map validation: actual key19 files, both known magics, unknown magic, matching
  and differing CRC, missing/nil names, directory paths and case normalization.
- Map names/modes: case-insensitive no-op, filename truncation, basename extraction,
  short names, selected-name terminator/tail, first matching mode row, signedness,
  narrowing and separate file-to-game flag encoding.
- Game operator roster: actual list allocation, invalid/inactive/unitless players,
  duplicate add, init/clear/reuse, and head/middle/tail removal. The existing console
  fixture now uses the real owner; its original frozen hash remains unchanged.
- Shadow list: skip flags, insertion, removal, reuse and retained removed-node links.
- Latency: actual player iteration and report queues, controlled timing boundary,
  zero-timing search limit, local-player selection, payload narrowing and recipients.
- Tile/crown cleanup: actual grid allocation/free list, idempotence, retained tile
  bytes, object type lookup/cache, delayed deletion and real respawn ownership.
- State/save availability: full-width state words, resets,30-frame boundary,
  frame wraparound, absent local player/unit and each activity/UI blocker.
- Readiness: status, remembered position and arrival buff across game modes and
  sparse/unitless players. Local cooperative GUI construction is covered by the
  existing server-options contracts and headless integration rather than this owner.
- Arrival: actual alias reset, player/quest state, report owners and minimap tracking;
  reuses the existing initializer callback log to verify invocation and arguments.
- Reset: active/inactive/unitless players, score/stage/health restoration, connection
  mask and game flags. Established-loadout state avoids duplicating the separately
  qualified default-item generation corpus.
- Departure: script callback ordering, roster removal, object deletion, peer reference
  clearing and reports. Raw active state is independently checked against the
  intentional client/server Game2 difference; the shared capture records whether
  that target-specific contract matched.
- Settings broadcasts: actual GUI owner and report queue, defined blank fields,
  numeric narrowing, changed/unchanged snapshots and exclusion of local slot31.

## Corrections and review items

Independent contracts identified these prerequisites before freezing:

1. Save-slot counting used the count as the next slot index, missing files after a
   gap. It now iterates1..13 independently, matching the actual save-selection UI.
   First run reproduced37 failing subcases.
2. Path basenames shorter than four bytes produced a negative copy length. The
   directory branch now has the same zero clamp as the filename branch. This was
   established by source review; the original undefined copy was not executed.
3. Character counting reused a prior record's flag after an empty or rejected
   metadata file. Tenth run reproduced both overcounts. Clear each output and
   require successful loading before inspecting its flag.
4. Settings broadcast compared/sent uninitialized local bytes when GUI fields were
   blank. Tenth run reproduced changing snapshots and repeated broadcasts without
   an edit. Zero-initialize the temporary snapshot.

These are reversible corrections for review and require fresh production
qualification. No expectations have been regenerated to conceal a port difference.

Fixture development corrected a working-directory assumption, the metadata parser's
untouched string tail, two adapter/type compile errors, and missing ability-manager
initialization in the reset owner. Twelfth and thirteenth focused runs pass.
Raw evidence stays under build/port-session-entry. Thirteenth run and its independent
repeat match all20 new captures byte-for-byte; expectations are frozen. Static
mapped-memory checks pass. Default/server/highres each pass336 roots /38,970 tests
without skips. All206 accumulated captures /51,868 records match C and each other;
all three final target gates have identical2,407-file source. The production gate
has identical production source; its sole difference is the added porttest file. Three fresh ELF32/SSE2/CGO binaries
and ABI assertions pass; the full asset suite has the exact1,553 known failure
entries and15 passing /3 failing /32 skipped packages. Options/gameplay, save/load
and flat regeneration scenarios pass. See [qualification](session-entry-c-qualification.json).

The initial production gate stopped at a manifest assertion that incorrectly
expected Go-backed exports for C baseline routines. Corrected retained-C assertions
pass in c-production-final; no source change was needed. All gate sessions are joined.
Baseline **adb06d4e is committed and pushed**. A six-case departure extension now
covers match reset with zero/one remaining player and elimination winner selection.
Its first run and independent repeat match; the new capture is frozen. Final
accumulated target sweeps pass. Production-source identity is verified, so the
qualified production builds/scenarios are reused. All sessions are joined.

## Local disk recovery

Verified deduplication of the three completed player-state native scenarios
reclaimed **1,660,044,319 bytes**. Original assets/archive and changed scenario files
remain. Each run has a restoration manifest. The consumed cleanup script is
`build/port-session-entry/deduplicate-player-state-native-assets.py`; never replay
its audit/apply modes. Its --restore mode remains available for those three runs.
Original assets and the untracked archive are unchanged.
