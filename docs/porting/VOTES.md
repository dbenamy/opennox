# Client/server player voting

## Scope and status

The preceding prefab-script conversion **8e8db9a1 is committed and pushed**.
This connected candidate covers **34 C bodies / 999 original body lines** in
server vote allocation/lists/update, player/team thresholds, quest actions,
client vote-window state/selection and message construction. General message
decoding stays outside translation scope and is audited as a caller.
[Selection](votes-selection.json) records original source hashes. One C list-present
helper appears orphaned; its Go wrapper already reads the state directly. Confirm
reachability before selecting contracts. No algorithm is ported or oracle frozen.

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
No voting algorithm has been converted yet. The ignored native drafts await
integration after the qualified baseline is committed and pushed.
