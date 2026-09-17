# Team runtime, membership and map objectives

## Status and scope

Baseline investigation follows qualified match/roster conversion `12da9a9e`.
Last qualified production was **52,273 physical C lines / 82 files**, zero reference C.
The corrected C baseline is **52,275 lines** and fully qualified (details below).
No native team conversion is installed or qualified.

The original candidate scan had **39 functions / 1,011 C body lines**. Accepted
scope is now **37 functions / 992 current C body lines** after prerequisite fixes
and the UI owner boundary described below. It covers:
team naming/lookup and member lists, membership changes and notifications, score
messages, roster creation messages, team assignment/balancing and map flag/crown
setup. A preliminary range scan also included unrelated equipment removal and two
floating-point helpers; these are excluded. Compact one-line C wrappers were
recounted by their actual function extents, not distance to the next marker.

`nox_xxx_unused_418840` is a 62-line orphan: its only non-definition reference is
in `legacy/keep.go`, behind `//go:build none`, with no declaration, actual caller,
callback or registration. Remove the obsolete retention reference and C body in
the conversion; do not add test-only calls to retain an unused mode. The other
functions require behavioral coverage. GUI-adjacent team setup is included because
its actual behavior creates teams and associates map flag objects; GUI callbacks
remain dependencies with their own owners.

Source extents/references: build/port-team-runtime/scope-audit.json. This is the
working scope audit, not a frozen baseline. Prior 42-function range inventory is
superseded; no implementation has been removed.

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
wrappers with their client team-selection window owner. Their nonempty path opens
and positions that window (and initializes the ball); moving those orchestration
wrappers without the UI owner would add weakly covered code to this batch. The
shared flag scan/assignment is included and independently checked with real map
objects. Empty-wrapper regression cases remain useful checks on the C-to-Go
boundary after conversion. This reversible scope choice is recorded for the later
team-selection UI batch; no production behavior changes because of it.

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
