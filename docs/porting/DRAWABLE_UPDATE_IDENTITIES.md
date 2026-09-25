# Drawable update callback identities

Scope: 28 C wrappers feeding the drawable update loop: the 23 trail/cloud/spell
update exports, orbit and monster-generator updates, color-light updates, and
both sprite-motion callbacks. Preserve the algorithms and 386 record layouts.
The original 25-wrapper proposal was expanded to include the remaining table and
secondary-slot siblings after inspecting every table entry: 25 table entries plus
three callbacks assigned directly to sprite fields. This is a reversible
batching decision, not an algorithm change.

## Behavior and original baseline

The root update loop saves the next list entry before invoking the primary
callback. It invokes the secondary callback only when the primary is absent or
returns nonzero; pause bypasses both. Preserve signed 32-bit results and the
secondary slot's discarded return. Distinct callback identities stay distinct.
Native identities must never enter a raw C call. Unknown callback addresses retain
the existing C calling conventions while other engine families are migrated.

The original magic-trail export has a void declaration, but the root loop observes
it through an integer-returning call. In the qualified Linux 386 binary, its normal
stack-canary comparison leaves the return register zero. The new
`TestClientUpdatesCallbackGate` confirms zero on all three original profiles and
checks that the actual magic owner ran and suppressed the secondary callback.
This is exact-target behavior, not a portable promise about void C functions.
The native adapter will call the same owner and return zero.

The contract also covers nil primary/secondary callbacks, raw primary results
0, 1, -1 and signed-32 limits, exact pointer arguments/call counts, cloud-height
wraparound and pause. It uses the real drawable pool and list owners, with existing
C observation fixtures. Original direct-result and loop-gate cases reset owner
state independently. Existing owner tests retain frozen mutations, RNG, deletion,
position, pixel and state expectations. Color-light fixtures invoke its native
owner directly; primary review of the simple pointer/int forwarding wrapper and
production qualification complement those tests, without claiming a separate
original-wrapper capture.

The first probe failed because the minimal test client does not initialize the
production callback table. A test-only Go getter now exposes the same C address,
as the existing cloud getter does. No new C observer code or broad global-table
initialization was needed. This fixture correction passed all three profiles;
failed logs remain under `original-probe/`, accepted runs under
`original-probe-fixed/` under `build/port-update-identities/`. Primary also corrected the helper's spawn count to account
for the recorded parent creation and added a successful-child assertion.

The conservative affected selection has 1,033 default/highres and 1,023 server
roots. It follows selected fixture APIs, stored callback fields and owner helpers;
ten explicit `!server` roots are excluded on server. An independent function-value
review found only local-name/field collisions outside the selected call closure.
The complete original selection passed against verified binaries on all profiles,
with the new gate contract repeated. Baseline accepted; no conversion installed. Production baseline is the
qualified drawing revision `b03fb880`, with only test changes for this baseline.

## Conversion and qualification

Use a separate typed update registry with static identities. Migrate the table,
assignments to both callback fields, sparse fixture getters, direct fixture calls,
and every raw field/alias consumer together. Preserve retained export comments
and remove obsolete prototypes. Do not change owners, assertions or captures.

After original baseline acceptance and commit, run the affected selection on all
three converted profiles, the complete default root corpus for the new shared
update dispatcher, safe/static checks, three production builds/ABI checks, exact
known-suite comparison, fresh preflight/final save-load scenarios and asset hashes.

## Delegation and recovery

Luna drafted the original gate probe and audits, and owns a fresh ignored legacy
migration overlay in `build/port-drawable-update-identities/luna-overlay/`. Primary owns
baseline acceptance, typed dispatch API, root/test callers, mapping review and
qualification. The initial audit overstated table membership and omitted the
retained target-motion table entry; primary inspected the table and requested a
corrected expanded audit. Stale older overlay artifacts in the parent directory
are not accepted input. No measured subscription savings are claimed.

The batch name was corrected during preparation: `update-identities` already
belongs to an earlier server batch. Its tracked documents were restored unchanged.
Original probe/baseline artifacts stay at their recorded paths; new conversion
artifacts and the replay manifest use `drawable-update-identities`. No engine source
or running test input changed during the baseline.
