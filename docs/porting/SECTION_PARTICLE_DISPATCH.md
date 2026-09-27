# Player-section and screen-particle dispatch

## Scope and baseline

Convert the player-file section/metadata and screen-particle callback dispatchers
and their connected observers to native Go. Retire the single-uintptr integer C
adapter and three fixture C imports. Preserve the current 386 target and external
backends. Baseline source is the qualified modifier revision `5ba6fce6`.

The selected owners comprise 139 existing root contracts in each of default,
server and highres: player files, session entry, client screen effects, final
callback boundaries, client-session messages and native layouts. Six focused
contracts repeat independently in each profile; two asset-free section/map
contracts also run with `safe`. Existing assertions and frozen captures remain
unchanged. This is an affected-owner selection, not the full root corpus.

Baseline accepted: 139/139/139 owner roots, six separately repeated focused
contracts per profile, and two safe contracts passed without skips. Original
production evidence is reused from the identical source at `5ba6fce6`; all
production gates will run again after conversion. See
[baseline](section-particle-dispatch-baseline.json). Native eight-file conversion qualified against baseline `18d66990`.

## Reviewed behavior

The client player-section table starts at blob offset55936, with 12-byte rows.
Blob initialization supplies the GUI/metadata/music identities through
`GAME_data_init.go`; framing reads callback slots at row+8, while session metadata
reads55956 directly. Both fixture table owners restore their saved bytes.

Particle creation supplies its native draw identity; the traversal fixture
replaces it with an observer and restores it. Preserve viewport and particle
pointer identity, saved-next traversal, callback mutation/new-head insertion,
ignored return values and the 52-byte layout. Existing particle contracts cover
pool reuse, byte fields, expiry and viewport edges.

The pointer observer is shared with a map-section adapter. Convert that adapter
with the observer, retaining nil/unknown-section guards, the non-nil pointer
argument and signed zero/nonzero results. Observer functions remain replaceable
and restore the preceding hook. Player-section dispatch still supplies nil.

The fixture-only client-write wrapper preserves signed 32-bit mode conversion.
Particle fixture parameters preserve signed positions/velocities and low-byte
truncation. Distance arguments/results preserve their 32-bit representations.
Retire only the three private C-typed wrappers whose callers move in this batch.

For review later: unknown section/particle callback keys become explicit panics
after the producer audit. This reversible correction is scoped to these families.
Tooltip/flame registrations and other dispatch families remain unchanged.

## Delegation and recovery

Primary implementation/review; Luna remains unavailable due to quota.
Ignored draft, source snapshots, caller searches, function comparison and run
records: `build/port-section-particle-dispatch/`. Generated dispatch is compared
byte-for-byte with its edited generator. No source changes occur during tests.

Before this batch, verified nine superseded binaries against committed source,
replacement binaries and host use, then removed 482,979,840 allocated bytes.
Rebuild `7e698a49` or `f5970121` for those historical outputs. Exact paths/hashes:
`old-binaries-{approved,result}.json` in the batch directory.

Twelve historical JSONL logs were gzip-compressed and verified by decompressed
SHA256 before removing originals, reclaiming 1,097,605,120 bytes. Restore with
`gzip -dk FILE.jsonl.gz`; exact paths, hashes and commands: `log-archive.json`.
Current modifier binaries/raw logs, original assets and archive remain intact.

## Native qualification and counts

All matching 139/139/139 owner roots and six separately repeated focused contracts
per profile passed without skips. Both original/native safe contracts passed.
Safe build/static, three production builds/ABI, exact known-suite outcomes and
fresh headless save/load passed. All 1,654 original asset hashes are unchanged.
Retired observer/dispatcher symbols are absent; retained embedded C bodies,
engine helpers, assertions and layouts are unchanged.

Embedded production C bodies: **12→11** (68/79 retired). Fixture C imports:
**94→91**. Production cgo files remain **5 client/highres, 6 server**; project cgo
packages remain three. Selected legacy exports remain zero. Headers remain
157 files /2,731 physical lines. Standalone C remains **0 production /0 test**.

Evidence: [qualification](section-particle-dispatch-qualification.json),
[inventory](section-particle-dispatch-inventory-after.json).
The initial ignored generator invocation was corrected from `go run` to a built
host generator after argument handling failed; generated output then matched
byte-for-byte before installation. No source/test correction was needed during
native qualification.
