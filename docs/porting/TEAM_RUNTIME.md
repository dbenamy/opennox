# Team runtime, membership and map objectives

## Status and scope

The native conversion is fully qualified. It replaces **36 live C functions**
and removes one **62-line orphan**, retaining **25 C interfaces** and retiring
**12**. Two private globals move to Go. Production C is **51,203 physical lines /
82 files**, zero reference C: **1,072 lines removed** from corrected baseline
`0b4b853c` (52,275 lines). All frozen expectations remain unchanged.

Default/server/highres pass **390 / 389 / 390 roots**, about **56,100 leaf cases**,
and **151 identical capture groups / 58,911 records**. Fresh production passes
all builds/ABI, the exact known asset failures, gameplay, save/load and flat rendering.
See the final qualification below; earlier investigation notes are historical.

The original candidate scan had **39 functions / 1,011 C body lines**. Accepted
scope is **37 functions / 992 baseline C body lines** after prerequisite fixes
and the UI owner boundary described below. It covers:
team naming/lookup and member lists, membership changes and notifications, score
messages, roster creation messages, team assignment/balancing and map flag/crown
setup. A preliminary range scan also included unrelated equipment removal and two
floating-point helpers; these are excluded. Compact one-line C wrappers were
recounted by their actual function extents, not distance to the next marker.

`nox_xxx_unused_418840` is a 62-line orphan: its only non-definition reference is
in `legacy/keep.go`, behind `//go:build none`, with no declaration, actual caller,
callback or registration. The obsolete retention reference and C body are removed;
no test-only calls retain an unused mode. The other
functions require behavioral coverage. GUI-adjacent team setup is included because
its actual behavior creates teams and associates map flag objects; GUI callbacks
remain dependencies with their own owners.

Source extents/references: build/port-team-runtime/scope-audit.json, captured before
conversion. The initial range inventory is superseded. The investigation sections
below record intermediate states; the final qualification sections take precedence.

## Baseline investigation

The first independent contract exercises real C rename, membership-change and
join/change-request messages through actual team allocation and existing queues.
It checks complete fixed-width records, including short/empty/Unicode names and
the reserved final word. C currently leaves rename padding and the final word of
three team messages uninitialized. Reproduce the behavior before changing C; never
freeze incidental stack bytes or normalize arbitrary payload differences.

The name field contains 20 transmitted UTF-16 code units in the rename message;
use at most 19 plus a terminator for this format. Boundary behavior of longer names
must be audited against callers before choosing a correction. Membership-change
and request decoders consume their existing 10-byte format; preserve that length.

## Ownership and coverage plan

Reuse the match/roster real sparse player owner, reliable message queue, team array,
world objects/type registry, client drawable owner and map-object creation/deletion
owners. Exercise host and client lookup separately. Tests must model actual linked
membership: assigning only an ID does not satisfy the C membership predicate.
Capture list ordering, node IDs, counts, report-mask and friend-list effects, and
controlled unit updates using the actual shared implementations.

Cover nil/missing teams and members, head/middle/tail removal, repeated join/remove,
team switching, signed scores and byte-sized counts, identifier/name widths,
case-insensitive lookup, recipient/mode filters and both client/host requests.
A direct unlink deliberately does not decrement the team count; its callers own
that operation. Preserve and test this division of responsibility.

For map setup use nonzero type IDs and real flag/crown update storage. Check object
retention/deletion, minimap reports, team flag/crown pointers and on/off state.
Balancing contracts must capture RNG consumption (50 pair swaps for multiple
eligible players), spectator/host filters, clan mapping, capacity and ties.
Nearest-flag assignment requires explicit 386 floating-point boundary review.

Audit existing Go duplicates before reuse: ObjectTeam.Has/SameAs, team lessons,
object-team lookup and SetNameAnd68. In particular C wcsncpy zero-fills the remaining
name field whereas the existing Go setter uses a different copy helper; compare
raw field padding and the untouched adjacent word, not just the displayed name.

Freeze only after independent contracts and repeated C captures. Reuse production
`12da9a9e` only if production source remains identical; any C prerequisite fix
requires fresh production qualification before its baseline commit. Then port,
qualify all three targets and affected callers, run fresh production/headless
gates, record C LOC, commit/push and continue.

### Initial probe results

The original-C padding probe reproduces short/empty/Unicode rename padding
failures (125.30s). Its player switch then reached the real join-notification text
path without a string manager: a fixture ownership error, not a production bug.
The second probe uses an allocated non-player ObjectTeam linked through real C
create/move calls, isolating complete message formatting from player notifications.
Player notifications still require a separate real string/UI owner in broader
coverage. No original-C expectation is frozen yet.

### Confirmed prerequisite: initialize fixed-width message buffers

The corrected original-C probe (padding-before-3, 33.26s) reaches all four paths
and confirms uninitialized trailing bytes in rename, membership change, join
request and change request. Initialize those four local buffers to zero without
changing field values or message lengths. The receiving implementations consume
46-byte rename and 10-byte change/request messages; their reserved tail has no
semantic input. This is a reversible prerequisite correction authorized by the
batch process. Physical C LOC is unchanged. Corrected-C qualification is pending.

The first retry failed discovery on a missing test-helper argument (1.85s), fixed
before the complete original-C probe. No capture expectation was generated or
changed. Initial corrected-C contracts add raw UTF-16 field/adjacent-byte boundaries
and actual player-member list removal/duplicate joins. They are running together;
no native source has been installed.

Initial corrected-C checks pass: three roots covering message padding, raw name
storage and linked player membership (124.12s), with no skips. Expanded original-C
coverage now adds membership predicates (including detached matching IDs), name
lookup, least-populated/available-team selection and empty flag-map returns.
No captures are frozen until broader ownership and boundary review is complete.

Expanded seven-root C checks pass (124.00s), including zero-return/zero-count
contracts for empty flag maps. The next extension uses actual allocated objects,
controlled nonzero flag-name/color table entries, server list linkage and real
on/off functions. It checks 0–3 flags against 0–3 teams, including missing team
matches and repeated toggles. This avoids accepting zero-filled color tables as
an oracle for map setup.

### Remaining review before freezing

Broader coverage still needs player join notifications with a real string manager,
client drawable lookup, group-mode assignment, reset/rebalance, nearest flag
floating-point ties, and crown deletion/minimap ownership. Host-balancing contracts
are drafted with independent eligibility/count checks and an exact 100-draw RNG
consumption check when more than one player is eligible. Do not treat the current
eight-root discovery suite as full qualification of the 39-function scope.

The C flag-map wrappers write a boolean into only the low byte of an int local.
Empty-map contracts currently return zero, but nonempty and target variants must
be checked; do not silently copy an uninitialized-high-byte assumption into Go.

The eight-root map-flag extension passes (122.90s). Ten-root assignment checks
now include 96 host-balancing combinations and 108 nearest-flag cases. The latter
uses an independently known x=5 bisector, exact/adjacent float32 values, missing
flags, existing membership, the initial distance cutoff, infinities and NaN.
Balancing checks eligibility, team-size bounds and random consumption; frozen
captures will additionally preserve actual member order and messages.

Assignment contracts pass (35.14s). The new player-join probe installs a real string
manager and console with an observing printer, owns the centered-message ring,
and checks text content/deadline, membership-message payload and mode/notify
filters. It exercises the original path that exposed the earlier incomplete
fixture. All remaining C team entry points now have thin test accessors available,
so further Go-only fixture additions can reuse the compiled C package.

Crown setup also calls the retained C respawn registry when it removes a crown.
The registry must be allocated and populated through its real allocator/add
functions; an empty fake global would make unregister dereference a missing head.
The next fixture owns/restores that shared registry rather than stubbing out its
behavior. Delayed deletion and minimap effects must likewise run through real
owners and be observed before teardown.

### Confirmed prerequisite: join-message object type

The full join probe passes text/ring/deadline checks but finds three message type
failures (123.20s). A narrower diagnostic confirms the first eight bytes match;
the final object-type word differs from the actual registered type (33.04s).
C createAt writes that word into a separate `short v22`, while sending the first
10 bytes of `int v21[3]`. Store it explicitly at byte 8 of that array and remove
the unused local. This preserves the intended existing record format and fixes
client object reconstruction when it needs the type. No packet goldens are changed.

Working C is **52,272 lines / 82 files**, zero reference C (−1 prerequisite line).
The C padding corrections and this type-field fix await fresh qualification before
the baseline commit. Crown setup now uses real object allocation, delayed deletion,
minimap lists and the real populated C respawn registry; its probe is running with
all preceding team contracts. Nothing has been translated yet in this batch.

Twelve-root corrected-C checks pass (123.15s), including the repaired object-type
field and 32 crown/flag presence/team-mode cases with actual respawn removal,
minimap membership and delayed-deletion state. Lifecycle expansion now adds
client dynamic-vs-static drawable lookup, C/existing-Go signed score messages,
and linked/detached member departures with team lifetime and group reset behavior.

### Accepted owner boundary

Accept **37 functions / 989 current C body lines**, including the 62-line orphan,
after the join-type prerequisite. Keep the 21 lines in the CTF/Flagball map-entry
wrappers with their client team HUD/window owner. Their nonempty path opens
and positions that window (and initializes the ball); moving those orchestration
wrappers without the UI owner would add weakly covered code to this batch. The
shared flag scan/assignment is included and independently checked with real map
objects. Empty-wrapper regression cases remain useful checks on the C-to-Go
boundary after conversion. This reversible scope choice is recorded for the later
team HUD and player-list UI batch; no production behavior changes because of it.

The original 39-function candidate list is superseded by scope-audit.json plus
deferred-ui-wrappers.json. The two wrappers' partial-width boolean locals still
need review with the nonempty UI path. Do not claim their positive branch is
qualified by the empty-map contracts.

Fifteen-root lifecycle checks pass (32.85s). The reset/capacity run (33.50s) passes
those new contracts but corrects a map-team fixture expectation: C truncates a
localized title to 20 UTF-16 code units, including the fixture's long fallback
string ID. Apply that independently specified bound; do not change production
name behavior. The next group/rebalance probe also checks that each team's member
counter agrees with assigned members after the clear-and-reassign path. The low-level
unlink contract remains caller-owned; the higher-level clear operation may be
missing its counter update. No baseline is frozen while this is unresolved.

### Confirmed prerequisite: clear/rebalance member counts

The 19-root group/rebalance probe (34.35s) passes all other roots, including the
corrected map-team title expectation. Every reset-and-reassign group case exposes
a stale counter: clearing three players leaves three counted members, and each
new assignment adds to that stale value. The low-level unlink helper still leaves
the counter to its caller, as required by remove/switch owners. Correct the
higher-level clear-player caller to decrement after a successful unlink. Check
that the member ID was cleared so a detached matching-ID entry cannot underflow
the count. The independent group-membership invariant and a detached-member case
cover this correction; update the earlier clear orchestration expectation to zero.
The direct unlink expectation remains unchanged. No frozen hash is changed.

The first combined corrected run (122.76s) stopped at the new client-join fixture:
its team-specific message manager was absent even though the general string manager
was installed. Reuse PortTestMapDrawableTeamMessages for that owner too. This is
fixture setup, not a production failure. The retry runs all 21 team roots.
Working C after all prerequisite edits is **52,275 / 82 files**, zero reference C;
net **+2 lines** versus the preceding qualified conversion. The accepted function
body scope is now **992 lines** after adding the guarded counter update.

All 21 corrected-C contract roots pass (corrected-c-2, 35.56s), with no skips.
The combined capture now selects 46 roots: those 21, all 24 preceding match/roster
roots, and real client map-team registration. The manifest includes fresh
production qualification because this baseline changes production C. Expectations
are not frozen yet; successful discovery alone is not baseline acceptance.

## Frozen corrected C capture

Combined C capture passes in **17.89s**, selecting 46 roots with no skips.
It produces **40 groups / 10,791 captured records** (array entries, or one record
for the clear/reset object). New expected hashes are now installed in source and
the manifest; the 22 preceding match/roster/map groups remain unchanged. Frozen
default/server/highres repeats are running from the same source. No source edits
until all gates finish. No production qualification or baseline commit yet.

## Corrected C baseline qualification

All frozen repeats pass with identical 40 capture files: default **48.70s**,
server **134.14s**, highres **57.89s**. The corpus executes **46 roots / 10,835
leaves / 10,791 records**, no skips. The 22 preceding capture groups are unchanged.
Fresh production passes in **367.10s**: all three builds and ABI checks; exact
1,553 known failure entries and package outcomes (15 pass / 3 fail / 32 no-test);
gameplay, save/load and flat rendering match the preceding native references.
Flat validation removes 51 maps and regenerates one. All four gates use identical
2,078-file source manifests and report unchanged source. Every process is joined.
Client SHA256: `7113b83605f11d38a258ee9e4954daeb95c9f2cbb3d5b2a06f6cee007c8ce1b9`.

Artifacts are under build/port-team-runtime/c-frozen-{default,server,highres} and
c-production. The baseline manifest and expectations are committed. Earlier
investigation paragraphs describe intermediate states, superseded by this result.
No C implementation is retained solely as a test oracle after conversion.

Completed C scenario copies were deduplicated against the unchanged original
assets, reclaiming **1,660,044,319 bytes**. Per-run restoration manifests preserve
paths, hashes and metadata; captures and unique saves remain. This cleanup is
consumed; do not rerun it unless the copies have first been explicitly restored.

## Native conversion in progress

The corrected C baseline is committed/pushed as `0b4b853c`. All 36 live functions
are translated; the 62-line unreachable retention function is removed. The caller
audit retains 25 C interfaces and retires 12, including the orphan. Two private
C globals become Go state; fixtures now own that actual state. Existing Go callers
invoke the new owners directly, and existing ObjectTeam predicates are reused.
Score updates preserve the original direct reliable-queue dispatch (see below). Deferred map-entry wrappers
still call the retained boolean scan interface.

Review preserves raw bounded UTF-16 copying (including zero-fill and untouched
adjacent storage), signed byte returns, 32-bit identifiers, queue ordering and
message operation codes. Membership mutation keeps direct unlink count-neutral;
clear, leave and switch own their increments/decrements. Assignment retains the
50 pair swaps / 100 random draws, eligibility filters and clan mapping. Nearest
flag selection retains float32 subtraction, double distance arithmetic, float32
best-distance storage and strict first-match ties. Crown setup keeps real delayed
deletion and retained C respawn-registry removal.

Join text uses the actual string manager and centered-message owner. The two
shipped format strings were inspected and use ordinary string substitutions.
Relocation continues through the previously ported random-placement and unit-move
owners; the team notification corpus does not itself exercise the relocation
branch with both chat mode and movement enabled. Full production and accumulated
placement/player-control coverage remain required; do not describe the smoke
scenario as complete team-mode integration.

The first native discovery stopped on a missing GAME1.h include for the retained
clan-mode helper (45.45s); fixed before retry. No expectations changed. Broader
selection is **389 roots**, taking the preceding 283-root match/roster selection,
new team contracts and affected damage, spell lifecycle/effects, map-drawable and
client-object rendering owners. The exact selection and native qualification
manifest are committed with the completed batch. Native qualification is pending.

Focused native qualification passes in **133.39s** (native-focused-3): all 46 roots
complete without skips, all 40 frozen captures unchanged. The preceding retry
stopped on an explicit uint32-to-int caller conversion (78.54s), corrected before
this pass. No C or Go oracle changed. Final source cleanup removes obsolete
function markers/blank lines; working physical C is **51,203 / 82 files**, zero
reference C, a **1,072-line reduction** from the committed corrected baseline.
Broader default/server/highres checks are now running; production is pending.

### Broader qualification caught a dispatch-owner mismatch

The first broader default/server sweeps finish in 285.68 / 280.54s with exactly
two failures: TestObjectivesHomeScore and TestObjectivesCTFScore. All other roots
pass. Their independent scoring/position assertions pass, but complete state
captures differ. The existing Go TeamChangeLessons setter uses the replaceable
Server.NetSendPacketXxx hook; the original C sender's already-ported bridge calls
the actual reliable queue directly. The objective fixture intentionally observes
that hook, exposing a distinction hidden by comparing normal queue output alone.

Preserve the original dispatch owner with a direct score-update helper and add
an independent hook-versus-queue contract. Keep the existing public Go setter's
behavior unchanged. Do not change any frozen captures. Highres first sweep is
still running; no source edits until joined. The earlier claim of full score
setter equivalence is limited to default-hook output and is superseded here.

All first broader processes are joined. Highres completes in **344.79s** with the
same two failures and no others. Targets execute 389 / 388 / 389 roots with no
skips; server excludes the client-only TestClientObjectRenderOcclusion by build tag.
The direct score-update helper is installed, with a new independent
TestTeamRuntimeLessonsDispatch contract. Focused retry now selects 49 roots: the
46 frozen C roots, that dispatch contract, and both affected objective-score roots.
Broader selection is now **390 default/highres roots / 389 server roots**.
No expected hashes were edited. Source is frozen while focused retry runs.

Focused retry **native-focused-4 passes in 134.10s**: all 49 roots complete without
skips, both original objective-score hashes restored and all 40 C capture files
unchanged. The new contract proves the legacy entry point bypasses the replaceable
server send hook and still queues the complete score record. Broader retries are
running from this fixed source. Earlier C baseline captures and production remain
the references; no re-recording was necessary.

## Native broader qualification

All corrected sweeps pass: default **158.90s**, server **261.03s**, highres
**197.54s**. They execute **390 / 389 / 390 roots** and **56,105 / 56,104 / 56,105
leaf cases**, with no skips or missing completions. The server's one omitted root
is client-only render occlusion, excluded at build time. Every target produces
identical **151 capture groups / 58,911 records**, including all 40 frozen team
baseline groups. All three share an identical **2,083-file source manifest** and
report unchanged source. Processes are joined. Fresh production is running.

Dispatch review is now part of PORT.md's existing-Go reuse guidance: check the
actual hook/queue owner as well as byte output and state changes. This adds a
specific boundary check rather than rerunning unrelated suites for every helper.

## Native production and final audit

Fresh production passes in **369.47s**. Standard/highres/server binaries build and
pass ABI/export checks (70.44 / 9.72 / 75.65s). The full asset suite preserves the
exact **1,553 failure entries**, package results **15 pass / 3 fail / 32 no-test**.
Gameplay and save/load pass against corrected C references. Flat rendering also
matches, with 51 maps removed and one regenerated. All four final gates share the
same 2,083-file source manifest, report unchanged source, and are joined.
Client SHA256: `5edffa51d483b8868f7511fa93a755e7ed2fb9acbcf9b4d194c1557a949a0929`.

Final artifacts: build/port-team-runtime/native-{default,server,highres}-2 and
native-production; native-qualified-audit.json summarizes counts. The 12 retired
functions have no source references. The two retired private-global names/offsets
remain only as historical blob-map labels, with no live reads. All 25 remaining
C entry points are thin Go bridges; no algorithm is retained only for testing.

The next connected candidate is the team HUD and server player-list UI, including
the two deferred map-entry wrappers and the CTF construction/tooltip owner. The
preliminary audit is 35 functions / 921 C body lines; ownership and reachability
still require review before freezing its baseline. Earlier references to a
“team-selection window” were imprecise: these map-entry helpers open team/flag HUDs.

Completed native scenario copies were also deduplicated against original assets,
reclaiming **1,660,044,319 bytes** (3.32 GB across both C/native trios). Restoration
manifests preserve every removed duplicate; unique captures and saves remain.
Both cleanup invocations are consumed. Original assets and archive are untouched.
