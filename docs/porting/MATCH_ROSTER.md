# Match results and player-roster synchronization

## Status and scope

The native conversion is qualified against corrected C baseline `bacee82b`.
It replaces **27 live functions** and removes one proven unreachable 19-line C
implementation. Scope covers roster/settings messages, match winner/limit
selection, remembered player identities, object report masks, flag-state records
and team assignment. Full player-arrival orchestration, server-ready UI transitions
and the GUI-widget-backed settings reader remain with separate owners.

The conversion removes **773 physical C lines**: 767 function lines, four private
global lines and two EOF separators. Production C is **52,273 / 82 files**, zero
reference C. Ten real C interfaces remain and eighteen symbols retire. The sections
below record prerequisite discoveries, frozen baseline and final qualification.

## Message-padding prerequisite

The original roster sender uses an uninitialized 132-byte scratch buffer and
sends 129 bytes. Short player identifiers leave trailing bytes unspecified. The
settings sender likewise leaves unused bytes of its 16-byte server-name field
uninitialized. An independent regression sends long, short and empty strings in
sequence and checks the complete fixed-width field.

The first probe exposed settings padding but also corrected a fixture expectation:
the shared player owner includes an active nil-unit slot, so roster iteration sends
three records to the selected recipient, not two. The corrected original-C probe
then reproduced both padding failures (37.48s). Raw bytes and logs stay under
build/port-match-roster/padding-before-2; no stack data is frozen as expected output.

Before freezing the C baseline, clear the roster scratch buffer **inside each
player iteration** and initialize the settings message buffer. The regression now
also mixes long and short identifiers within one roster send. Preserve the actual
12-byte player-identifier storage and the existing 129-byte wire length: long
identifiers can occupy all ten transmitted field bytes. This is a reversible,
confirmed correction within the authorized workflow. The expanded corrected padding probe passes (131.42s). Fresh corrected-C production qualification passes; details below.

## Flagball timeout draw prerequisite

The C timeout selector calls the Flagball winner reporter with a null team on a
tie or with no teams. Both the original C reporter (before its earlier port) and
the current Go reporter require an actual team. The independent empty-match probe
reproduced the nil-team crash (132.67s); this is an existing caller bug, not a new
translation difference.

The client decoder's generic flag-result message (87, team 65535, time-limit marker
1) already handles a draw. Use that existing message for the no-winner branch;
keep message 86 for actual Flagball winners. This removes four C lines. The low-level
winner reporter retains its real-team precondition. No client protocol addition is
needed. Empty, tied and actual-winner cases now exercise real team allocation and
the reliable queue. Combined corrected-C prerequisites pass (135.54s).

## Ownership and qualification plan

Reuse real sparse player/unit owners, the qualified reliable queue, object/type
owners, actual team arrays and C list allocation. Preserve the active nil-unit
roster slot as distinct from unit iteration. Own and restore mapped settings,
score/time limits, cached type IDs, team state, remembered-list state and clocks.
Use shipped or controlled nonzero constants/tables, with independent format and
boundary assertions, before freezing any captures.

Cover string lengths and caller buffer extents, integer widths/overflow, signed
team/player scores and ties, spectators and unitless slots, recipient/host filtering,
object report masks, per-team state offsets and adjacent bytes, timer boundaries,
team selection/creation, list identity matching, cleanup and repeated use. Exercise
the actual message owner and callbacks, not copied C algorithms.

The settings message embeds the C compile-time protocol version. Source audit and
actual default capture show that legacy/video_highres.go has no build constraint:
its NOX_HIGH_RES C flag is enabled for **all** targets. C therefore emits 0x000F039A,
while the root Go version varies by target. Preserve the C value in this selected
message and independently assert it. Review the unconditional C flag separately;
changing it also changes rendering constants and is outside this conversion.

The lowercase C nox_xxx_wallSendDestroyed_4DF0A0 has only its definition and header
declaration. Whole-repository search shows callers use the existing Go method
Server.Nox_xxx_wallSendDestroyed_4DF0A0. Remove the dead C body/declaration in the
conversion, without building an oracle for it. Keep coverage on the live Go path.

After repeated C captures and independent contracts: qualify default/server/highres
and corrected-C production; commit/push the baseline. Then translate the batch,
retain only actual C interfaces, qualify affected callers and fresh production,
record C LOC/decisions, commit/push and continue.

## Player identifier layout prerequisite

The broader encoder contract compares C and the existing Go encoder, including
nonempty identifiers. It found that `Player.Field2096Buf` actually started at byte
2093: the preceding `Active` byte had no explicit trailing padding. C uses 2096.
Add three explicit padding bytes before the identifier and a compile-time offset
assertion. The following uint32 already supplied three implicit padding bytes,
so total player size and subsequent offsets stay unchanged. This repairs the Go
identifier accessor/setter and lets the shared encoder preserve the C layout.

The initial broader run also exposed a fixture lifetime mistake: `Teams.Reset`
clears records but leaves `ActiveCnt`; repeated cases must explicitly reset that
count. Neither golden output nor production team-reset behavior was changed to
accommodate this fixture correction. The c-first-3 run passes all 12 initial roots (135.97s). Subsequent expansion
covers roster/inventory/minimap, settings/IDs and victory limits. Queue ordering
and the player-holder high bit required fixture-expectation corrections. Settings
and roster/minimap/IDs pass in c-roster-2; inventory's holder-marker correction and
new victory contracts are under test in c-victory.

## Coverage expansion checkpoint

The 19-root sweep through timer modes passes (40.80s). Team assignment is exercised
with real C linked membership, actual team creation and the reliable queue. Roster
membership must use that list: merely setting an object's team ID is insufficient.
The quest map request contract checks both the full filename and the extensionless
basename; the latter is the existing storage convention, not a port change.

The timer position report uses a C `long long` conversion before narrowing to a
16-bit coordinate. Its separate contract includes values beyond signed 32-bit
range, signed 64-bit boundaries, infinities and NaNs. Ordinary object positions use
the already-qualified 32-bit conversion; preserve this distinction.

## Qualified corrected C baseline

**24 roots / 10,027 leaf cases / 21 frozen groups / 10,034 records**, no skips.

| Gate | Result | Seconds |
| --- | --- | ---: |
| Complete C capture | Pass | 44.95 |
| Frozen default repeat | Pass, all hashes | 56.54 |
| Frozen server | Pass, all hashes | 149.58 |
| Frozen highres | Pass, all hashes | 67.17 |
| Fresh production and headless gates | Pass | 400.55 |

The three frozen target runs and production have identical fingerprints for all
2,054 source files. All three production binaries and ABI checks pass. The full
asset suite matches exactly 1,553 known failure entries and its 15 pass / 3 fail /
32 no-test package results. Gameplay and actual save/load match the preceding quest
native scenes. Flat rendering also matches, with 51 loose maps removed and exact
regeneration of the selected compressed map. Client SHA-256:
`1be0263beff5541876b9d4f3fef383df8ee0458bbc5d6d54aa32e92e4786fd97`.

Artifacts: build/port-match-roster/c-capture, c-frozen-repeat, c-frozen-server,
c-frozen-highres and c-production. Headless scenes are match-roster-c,
match-roster-save-c and match-roster-flat-c. Expectations are committed in fixture
hashes and match-roster-c-batch.json; local capture payloads remain ignored.

Working production C after prerequisites: **53,046 lines / 82 files**, zero
reference C. This is the committed pre-conversion baseline.

## Native implementation

Baseline `bacee82b` is committed/pushed. The Go implementation is qualified. It consolidates the existing root roster encoder in server, moves
match/roster state and algorithms to Go, and replaces the private remembered-name
C allocation/list with Go records. The frozen list payload/order and identity
contracts remain unchanged. Ten C entry points remain for real C callers; eighteen
retired symbols have zero remaining source references.

Removal totals **773 physical C lines**: 767 function lines, four private global
lines and two trailing EOF separators. Working production count is **52,273 / 82
files**, zero reference C.

Native bring-up: first discovery found two nonexistent controlWord references
(85.25s), replaced with the existing equipmentWord accessor. The next full focused
run (137.74s) matched 20/21 groups; wall open/close used the ordered sender instead
of the original unordered reliable sender. Payloads matched; the queue metadata
identified the difference. Correct that flag without changing frozen expectations.
native-focused-3 passes (136.42s), with all frozen groups unchanged. Affected selection is 283 roots / 111 test source files,
covering direct callers and real report/queue/team/list/world/quest owners.

Completed C scene deduplication reclaimed **1,660,044,319 bytes**; per-scene verified
restoration manifests preserve the removed duplicate assets. Original assets and
the user's archive remain untouched.

### Translation review

Preserve fixed 129-byte roster messages and the larger 132-byte caller scratch
buffer, including identifiers that fill all ten transmitted bytes. The shared
Go encoder accepts both caller extents; caller initialization controls padding.
The remembered-name check preserves C's signed-char comparison against an unsigned
stored class byte, including the unusual high-byte cases.

Object mask shifts retain x86 count wrapping. Flag records retain their mapped
six-byte layout, untouched padding and C-visible pointers. Winner selection keeps
the separate team and solo-player tie passes and signed score bounds. The timeout
coordinate converter preserves the low 16 bits of the C 64-bit conversion using
float bits; large finite floats and masked invalid conversions both yield zero
low bits. Existing C team membership/notification owners remain shared dependencies.

All private Go callers now invoke Go directly. The retained C interfaces are GUI
settings, complete settings, roster fanout, player IDs, objective minimap, simple
object reports, team roster, both flag-record pointers and minimum-score winner
selection. The unreachable C wall-destroy implementation is removed; the existing
live Go method remains covered by the same wall fixtures.

## Final native qualification

| Gate | Result | Seconds |
| --- | --- | ---: |
| Frozen focused contracts | Pass, all 21 groups | 136.42 |
| Affected default | Pass | 115.52 |
| Affected server | Pass | 225.84 |
| Affected highres | Pass | 153.84 |
| Fresh production and headless gates | Pass | 388.13 |

Each broader target runs **283 roots / 52,837 leaves / 127 capture groups /
55,862 records**, without skips. All 127 groups match across targets; new frozen
contracts and the 106 preceding groups remain unchanged. Selection covers 111 test
source files and is reproduced by match-roster-tests.txt / match-roster-batch.json.
The three sweeps and production share identical fingerprints for 2,060 source files.

All three production binaries and retained/retired ABI checks pass. Full asset
qualification preserves the exact 1,553 known failure entries and 15 pass / 3 fail /
32 no-test package results. Gameplay, actual save/load and decoded flat-rendered
frames match corrected C. The flat scene removes 51 loose maps and regenerates one
selected compressed map exactly. This is bounded integration coverage, not a claim
that the whole game or full suite is green. Native client SHA-256:
`1676f4ba3ae34a7bb2abeced8d08021765721817841b4828dd5218355af45486`.

Artifacts: build/port-match-roster/native-focused-3, native-default, native-server,
native-highres, native-production and native-coverage-audit.json. Headless scenes:
match-roster-native, match-roster-save-native and match-roster-flat-native.
No C test algorithms remain. The committed C baseline and unchanged expectations
provide recovery without depending on local captures.

Completed native scene deduplication also reclaimed **1,660,044,319 bytes**, with
verified per-scene restoration manifests. Combined C/native cleanup freed 3.32 GB;
original assets and the archive are preserved. See dedup-native.log locally.
