# Client/server player voting

## Scope and status

The repaired C baseline **5572b500 is committed and pushed**. Native Go now
implements 35 live routines covering server vote allocation/lists/update,
player/team thresholds, quest actions, client window lifecycle/selection and
message construction. One orphan C helper is removed. General message decoding
remains outside translation scope and is tested as a caller.
[Selection](votes-selection.json) records the 36 original bodies and hashes.

The working conversion removes **1,178 physical C lines**, leaving **34,343 in
71 production files**, zero reference C. Five C entrypoints remain for live C
callers; 31 interfaces and 12 C globals are retired. Both name snapshot arrays
are Go-owned. Native qualification passes; see PORTING_STATE.md.

## Baseline strategy

Use real player/team/object owners, the allocation class, outbound queues and
actual GUI widgets/string tables. Capture list linkage, 52-byte vote records,
counts/bitsets, frame values, thresholds, effects and serialized messages. Cover
empty/single/multiple records, duplicate/missing participants, bit31 and byte
count boundaries, repeated removal, deleted/nonparticipating targets, team/global
thresholds and live update dispatch. Do not invent expiration behavior: no timeout
check was found in the selected updater; it dispatches on vote type.

For GUI behavior, cover all topics, admission, label/visibility/list selection,
selection deltas, repeated confirmations and actual outgoing messages. Use real
window parsing and callbacks; supply shipped data or deliberately controlled
nonzero values for consumed tables. Audit caller-observed results before freezing
decompiler returns that sometimes contain addresses or unused low pointer bytes.

Two 52-byte name-message temporaries are not initialized. Reproduce any padding
instability with independent contracts before deciding a correction; do not hide
it in a normalized golden. Shipped thresholds and the shared quest minimums must
be installed deliberately. Reuse the preceding qualified production baseline only
if production source remains identical; otherwise qualify fresh production.

Repeat/freeze C contracts, qualify affected default/server/highres owners, commit
and push the baseline, then convert/retire, compare without regenerating goldens,
qualify fresh production, document C LOC and commit/push.

## Caller and ownership audit

Server initialization/shutdown and tick enter through legacy/server.go; object
removal enters through legacy/object.go. The general server message decoder calls
cast/withdraw dispatch for actions 0–5. Those dispatch exports must remain while
that decoder is C. The updater saves the next link before actions may remove a
record. Allocation owns 52-byte records; the team field borrows object+48 and the
target field borrows the player unit. Empty/nonempty transitions broadcast the
three-byte vote-status message. Deleting a reset vote additionally notifies each
remaining voter. Sparse player slots, including slot31, use a uint32 mask.

Client initialization enters through legacy/game.go. The quit window invokes the
vote window, whose real callback dispatches confirmation/cancellation. The hide
helper also has an external C caller in GAME2_3.c around line 2941, outside the
selected bodies: retain that interface. The list-present C helper sub_5071C0 has
no live caller; its existing Go wrapper reads the list head directly. Header
prototypes and its own definition are its only C references.

The normal threshold reads a full uint32 for admission but stores its low byte
in a record. Controlled baseline values are 5 for normal votes and 9 for unknown
kinds. Quest kinds initialize their record minimum to literal 6; shared settings
named quest minimums gate admission but do not supply that stored value. Preserve
this observed distinction unless a separate independent contract justifies a fix.
Normal withdrawal currently searches kind 0 even when dispatch receives kind 1;
this needs an explicit regression/behavior decision before freezing.

Initial independent contracts are installed for allocation/zeroing, frame/team
ownership, prepend/delete linkage, transition messages, repeated player removal,
bit31 and complete52-byte client name messages. They are not yet a complete or
frozen baseline. Initial compilation/results are under build/port-votes/initial.*.

## Reproduced prerequisite correction

Initial independent message contracts failed for short ASCII, longer ASCII and
Unicode names, in both cast and withdrawal messages: bytes after the UTF-16
terminator contained unrelated stack contents. Both52-byte C temporaries now
use zero initialization. Message length, action, name and terminator are unchanged;
unused bytes are deterministic zero. This is a deliberate reversible correction,
not capture normalization. Four initial contract roots pass after the repair
(second.json). The first allocation assertion was corrected to reflect the
actual newest-first reliable queue. Fresh production qualification is required
because these two production C declarations changed.

The GUI fixture uses the real parser, widgets and callback with controlled labels
and artwork. Its IDs, dimensions, list capacity and single/multiple selection
modes follow shipped GuiKick.wnd. A test-only string-ID type error in the first
GUI compile was corrected; that failed run is not qualification evidence.

The decoder-level alternate-withdrawal regression also failed: after casting
kind 0 and kind 1 for the same target, action 3 removed kind 0 instead of kind 1.
The private normal-withdraw helper now accepts the dispatched kind and searches
that kind. Both callers are inside the selected dispatcher; the external decoder
interface is unchanged. The independent contract requires action 3 to leave kind 0,
then action 2 to remove it. This correction is made before freezing expectations.

## Current contract coverage and fixture limits

Thirteen roots passed together in eighth.json, including pool exhaustion/reuse,
team/global thresholds, quest withdrawal/target guards, real decoder routing and
successful normal/quest actions. Successful-action tests retain the real vote
updater and blocked-player list; existing disconnect/removal service hooks record
the requested player/reason, since those services are outside this translation.
A fixed clock makes the actual temporary blocked-player expiry deterministic.

GUI reset/topic contracts use real single-selection widgets and window callbacks,
checking labels, visibility, repeated choices, cancellation and two-byte messages.
The quest-reset contract uses actual participant counting, stage mutation, catalog
choice, list disposal and notification queues. It observes map-switch requests at
the existing server interface; players are in the actual respawn entry's guarded
quest-transition state. Full respawn equipment/positioning remains covered by its
accumulated contracts rather than recreated here. This limitation must remain
explicit in qualification. Fixture bring-up corrected catalog name offset 4 and
used the runtime quest-stage owner (587000+202028), distinct from the prefab
stage field. No production change was needed for these fixture corrections.

## Qualified repaired C baseline

All sessions are joined. Eighteen focused roots pass; eight captures / 381 records
repeat byte-for-byte in separate processes and are frozen in source. The final
team-color contract supplies palette index 7 and checks the actual red row color,
so default-white rendering cannot hide a missing table read. Full 31-player GUI
selection/withdrawal checks all 62 messages. Bringing up that contract corrected a
fixture iterator callback that had stopped after its first message.

Default/server/highres each pass 220 roots / 43,576 including subtests, with no skips.
All 168 captures / 48,969 records match across targets. The four gates share
unchanged 2,348-file source. Static mapped-memory checks pass. Fresh production
passes three 32-bit/SSE2/CGO builds and ABI/interfaces; the full asset suite matches
the known 1,553 failure entries and 15 pass / 3 fail / 32 skip package results exactly.
Gameplay, save/load and flat regeneration pass against existing references.
Production qualification took 369.7s. No original assets were changed.

[C qualification](votes-c-qualification.json), [capture hashes](votes-captures.json)
and [reviewed test selection](votes-tests.txt) record the evidence. Raw artifacts
are build/port-votes/final-{a,b}, c-{default,server,highres,production} and
 target-capture-audit.json. Production C remains 35,521 / 72 files /zero reference.
This baseline remains the oracle for the native conversion. Ignored installation
scripts are consumed and draft files are stale; installed source is authoritative.

## Lifecycle audit extension

The first native compile exposed two neighboring lifecycle routines omitted from
the original selection: `sub_48D450` (window disposal) and `sub_48D4A0` (choice
reset). Their reads of moved state required conversion together. Before translating
them, an isolated checkout of 5572b500 exercised the original C with real windows,
capture/stack ownership, repeat disposal and reinitialization. Six records matched
across default/server/highres and a separate repeat. The added fixture and C adapter
are recoverable as [a test-only patch](votes-lifecycle-c.patch); its hash, original
source identity and results are in [the extension report](votes-lifecycle-c-qualification.json).
The original baseline's 220-root result does not include this extension.

The window indicator helper uses separate state and stays outside scope. Disposal
now has a direct Go caller; choice reset retains its C entrypoint for the client
decoder. This raises the complete frozen set to nine captures /387 records.

## Native implementation and review decisions

The fixed-capacity allocator preserves the 52-byte record layout and borrowed
object/team ownership. Closing clears the allocator handle explicitly because
`ClassT.Free` has a value receiver. Native contracts close twice and restart three
times. Update/removal traversal preserves next-link capture and the C callback
order, including reading the player through its borrowed update data after the
service callback. Quest settings still gate admission separately from the literal
record minimum. Threshold subtraction keeps its unsigned zero-player behavior.

GUI player rows and name comparison preserve raw UTF-16 units. Native copies bound
outgoing names to 24 units plus terminator in the existing 52-byte message, and
selection snapshots to 32 names of 27 units plus terminator. Valid baseline cases
are unchanged. A missing local team object now returns without showing the vote
window; this replaces a nil dereference. These reversible boundary decisions have
independent native contracts and are recorded for review.

All actual callers ignore the old cast dispatch and show-window return values.
Their retained C declarations now return void instead of exposing incidental
pointer-derived decompiler results. Three other C entrypoints retain their types.
Production qualification checks the retained and retired symbols.

The first missing-team native test panicked during fixture team creation because
its string/message service was not installed. Reusing the existing team-fixture
setup fixed it without a production change. The final focused run passes 22 roots,
all nine captures match C byte-for-byte, and static mapped-memory checks pass.

## Qualified native conversion

Default/server/highres each pass **224 roots /43,580 cases**, with no skips.
All **169 captures /48,975 records** match across targets, including every original
C capture. The four gates share unchanged 2,352-file source. Fresh production
passes all three 386/SSE2/CGO builds, ABI and retained/retired interface checks.
The asset suite exactly matches the known 1,553 failure entries and package
outcomes. Gameplay/options panels, save/load and flat-map regeneration match
existing references. All test/build sessions are joined.

[Native qualification](votes-native-qualification.json) records the results.
Local evidence: build/port-votes/native-{default,server,highres,production},
native-third.*, native-static.log and native-capture-audit.json. No C algorithms
remain solely for tests. Completed C scenario asset copies were SHA-256 verified
against originals before deduplication, reclaiming 1,660,044,319 bytes; each run
retains its restoration manifest, changed maps/saves and complete reports.
