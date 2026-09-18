# Console command handling

The qualified conversion replaces **47 live C bodies**: all 45 routines in
`client__system__parsecmd.c`, map-column formatter `sub_4417E0`, and score-limit
setter `sub_409FB0_settings`. It removes two C files, 46 private function
interfaces, two C globals and three private scratch buffers. One C-to-Go dispatcher
entrypoint remains for the quit dialog's observer action. Production C is
**33,267 physical lines /69 files /zero reference C**, down **1,077** from the
repaired C baseline.

## Ownership and behavior

The actual console registry now invokes Go handlers without C token-array dispatch.
Handlers use existing settings, player, team, map catalog, spell, GUI, script VM,
message-queue and audio owners. The score setter's Go callers also invoke Go
directly. The borrowed remote-command sender is cleared when dispatch completes.

Preserve ASCII case folding, UTF16 low-byte conversion for server names and script
identifiers, embedded-NUL termination, signed32 decimal prefix parsing/saturation,
and the original byte/word narrowing order. Player limits clamp signed values;
score limits narrow to uint16 before clamping. Distinct requested scores above 999
still trigger the original update comparison before clamping. T1/LAN share their
rate/type but differ in online state. Local/server mute routes retain the reserved
slot31 distinction. Go strings replace fixed temporary buffers while final stored
server names remain limited to 15 bytes.

The C formatter ignores string width/precision and recognizes `%%` and `%!`.
The Go command formatter preserves these forms, including broadcasts with no
arguments and escaped conversion text. Unsupported argument-taking directives
without arguments have no defined C result; the Go no-argument path treats them
as text. See [DECISIONS.md](DECISIONS.md) for decisions to review.

## Baseline corrections and review findings

Repaired C baseline **15df0163** is committed and pushed. Independent contracts
reproduced two defects before freezing:

- “Respawn off” printed OFF while adding bit8. Change it to the existing removal
  helper; eight client/server, initial-state and case variants reproduced the bug.
- Mode-restricted remote dispatch returned with its borrowed sender still set.
  Clear it on the early return; fifteen mode/action combinations reproduced it.

An initial map-list expectation incorrectly assumed `%S` width padding. Inspection
of the actual formatter established the unpadded name+suffix+tabs contract. Fixture
setup also needed all three sequence-list sentinel words and the actual direct
Kind1 broadcast queue rather than the reliable-report queue.

The first native compile caught a missed caller in `client__gui__guiquit.c`.
Retain its dispatcher signature and a thin Go export. The remote fixture exercises
that C entrypoint, including the quit dialog's nil-text observer call. A fresh
whole-source reference audit found no other selected-body C callers.

Review after the first passing native target sweep found missing literal-format
semantics. Nine added contracts passed C and failed the native draft. Ten final
cases, including string width and escaped conversion text, passed twice from the
committed C baseline with identical captures. Preserve the reproducible
[fixture patch](console-commands-format-c-fixture.patch) and
[C evidence](console-commands-format-c-qualification.json). No original frozen
expectation was changed to accept the native implementation.

## Qualification

- Original baseline: 23 focused roots /527 passing tests, no skips; 17 identical
  captures /571 records from two independent processes. Selected-body hashes and
  frozen capture hashes are tracked in the selection/capture JSON files.
- Extended native focus: 25 roots /539 tests; all 18 captures /581 records match C.
- Final default/server/highres: **300 roots /20,221 tests each**, no skips;
  **169 captures /32,019 records** match across targets and qualified C evidence.
- Static mapped-memory validation passes; all four final gates have identical
  source fingerprints. Three fresh ELF32/SSE2/CGO binaries pass ABI checks.
- Full asset suite exactly matches 1,553 known failure entries and package outcomes
  (15 pass /3 fail /32 skip). Options/gameplay, save/load and flat-render scenarios
  match references; save/load explicitly saves, reloads and resumes.

Reports: [C baseline](console-commands-c-qualification.json) and
[native qualification](console-commands-native-qualification.json). The selected
regression scope includes rules, settings, options, teams, player controls, roster
and gameplay reports, reliable reports, map catalog, scripts, voting and audio.
The shared player-control fixture adds only an isolated level-command operation;
existing allocation, capture and operation behavior remains unchanged.

Contracts observe actual owners wherever practical. Existing service hooks observe
console execution, observer admission and ability cancellation outside this batch.
The headless MOTD check covers dispatch/reset, and GUI owners cover options actions.
Physical display and audible playback remain release checks. Existing unimplemented
commands retain their current return behavior.

Local evidence is under `build/port-console-commands`: `final-a/b.*`,
`native-format-fixed.*`, `native-final-{default,server,highres}`, `native-production`,
`native-static-qualified.log`, and the isolated `c-format` worktree. The earlier
native sweeps precede the formatter correction and are not the final gates.

## Disk maintenance

Completed voting-native and console-C asset copies were deduplicated against
SHA256-verified originals, retaining changed files and per-run restoration manifests.
Original assets and the archive were untouched. Another 274 large JSON logs/captures
from four completed older batches were losslessly compressed, reclaiming
2,869,209,914 bytes. Every compressed file was decompressed and SHA256-verified before
removing its original. The local `archived-old-evidence.json` records paths, hashes,
metadata and `gzip -dk` restoration instructions. Current-run evidence and the
known-suite oracle were excluded.
