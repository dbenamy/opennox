# Player admission, status and displayed equipment

Repaired C baseline **81146f36 is committed and pushed**. Native conversion
replaces **20 live routines**, removes **two orphan bodies**, retires **11 C
interfaces** and retains **11 exports** for decoder/server lifecycle callers.
Production C is **32,605 physical lines /69 files /zero reference C**, **−663**.

## Baseline correction for review

The original active-competitor loop counted empty teams whenever any eligible
player existed elsewhere. Independent contracts using actual memberships failed
35 team-mode cases and no non-team cases. Add the same membership predicate used
by the neighboring per-team counter before freezing. Non-team counting still
includes active player records without units. The separate multiple-participants
query keeps its status bit 0x20 behavior and does not exclude the headless host slot.
This reversible correction adds one C line to the baseline (33,268).

## Contracts and native behavior

- Actual sparse player records, units, teams, report queues and modifier owners.
  Quest limit of 6, nil admission callers, re-entry, signed score/flag-capacity values,
  remembered team identities, team limits and selected special-mode bytes.
- Sudden-death data conversion, status reset, frame wrap and strict 20-second
  boundary; re-entry setter preserves all 32 bits.
- Exact seven-byte status messages, low 16-bit IDs, mask 0x423, host gating and recipients.
  Timer startup distinguishes host/chat/participants/configuration/previous state;
  clearing observer status again after time advances does not restart it.
- Whole-roster status includes the active unitless player. Lesson reset clears all
  active records and emits messages only for players with units. Queue contracts
  account for head insertion and the connection mask versus destination fields.
- Name lookup preserves first-match ordering, nil inputs and raw UTF-16 with
  ASCII case folding, matching the existing C locale; non-ASCII units stay intact.
- Minimap tests use real allocated circular lists: fanout, duplicate mark/merged
  flags, partial removal, middle/head/final unlink, invalid indices and nil objects.
- Equipment tests cover full/empty slot arrays, unknown players/modifiers, every
  respawn mask bit, classes and game modes. Unused modifier fields and each slot's
  sixth word are preserved. Missing base-color data returns after clearing slots.

Go callers invoke native helpers directly. The two orphan scalar helpers have no
callers/registrations across source and are removed without replacement tests.
No test-reference C remains. Scope stays with this connected behavior batch.

## Qualification

Focused native checks pass **16 roots /18,542 tests including subtests**. All 16
captures /18,779 records match the independently repeated frozen C baseline.
Default/server/highres each pass **316 roots /38,763 tests**, without skips.
All 185 captures /50,798 records match each other and C byte for byte. Static
mapped-memory validation and interface audits pass. Fresh native production passes three ELF32/SSE2/CGO builds and ABI checks, the
exact known full-suite results (1,553 failure entries; 15 pass /3 fail /32 skip
packages), and options/gameplay, save/load and flat-map scenarios. All four gates
share unchanged 2,386-file source. All sessions are joined. See
[C qualification](player-state-c-qualification.json) and
[native qualification](player-state-native-qualification.json).

## Process and storage notes

The installer's first import cleanup touched 14 unrelated files, including directive-only
cgo files. The check was interrupted and all 14 files restored exactly from the
baseline before restarting. An unused import in a changed file was then removed.
PORT.md now explicitly preserves `#cgo` imports and confines cleanup to batch files.
The final change does not alter those unrelated files; the installer is consumed.

Hash-verified deduplication of completed console-native and player-state-C scenario
assets recovered 1,660,044,319 bytes each. Per-run manifests preserve restoration
paths, hashes and metadata. Twelve inactive compiler temporary directories older
than 24 hours were removed after checking process references (579,880,242 bytes).
Original assets and the archive are unchanged. Lossless compression of 775 old evidence files recovered another 2,877,720,135
bytes. Decompressed hashes were verified before deleting originals; restoration
manifests remain under build/port-player-state. Cleanup scripts are consumed.
