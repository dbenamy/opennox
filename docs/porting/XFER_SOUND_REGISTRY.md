# Object transfer and damage-sound dispatch

## Scope

Candidate: bind 28 object-transfer registrations and two damage-sound registrations
to their complete existing Go export wrappers. Preserve C callback identities,
object layout, unknown-callback calling conventions and late handler replacement.
Production conversion is not installed. The preceding qualified collision
checkpoint is `d44f33d7`.

`server.Object.CallXfer` treats every nonzero integer as success and converts zero
to its existing generic error. Its callback must be configured. The outer
`opennox.Object.CallXfer` has a separate DefaultXfer shortcut: it calls the local
implementation and returns its detailed error. Keep that shortcut unchanged; it
does not consult the replaceable legacy handler.

The damage owner chooses weapon over source and ignores the callback return.
A nil sound slot calls the fixed default wrapper, even if the registry's default
pointer changes. Keep that branch and the owner's admission, ordering and sound
suppression conditions unchanged. The two wrappers resolve their Go handlers at
invocation, as does the server DefaultXfer wrapper.

## Baseline and qualification plan

Five new test-only files add eight contract roots. They cover all registration
identities, raw zero/nonzero return boundaries and arguments, successive handler
replacement, generic versus detailed errors, exact default serialized bytes and
read state, sound selection and invocation timing, ignored sound returns, and
real deferred audio events through the unchanged damage owner.

Existing object/item/creature transfer fixtures already call the registered
server API for all 27 nondefault types. Reuse their serialized-record and ownership
goldens. A new identity test verifies every fixture's actual callback slot.
Select 228 exact roots including transfer, damage, player files, prefab runtime,
resource definitions, map sections and orchestration. Require original and
converted runs in default/server/highres without skips; compare discovered names,
source/environment fingerprints and capture hashes, not just test counts.

Reuse preceding production evidence only for the test-only original baseline,
after checking source and binary hashes. After conversion run safe/static,
fresh production/ABI checks, exact known-suite comparison and headless
character creation plus explicit save/load/resume. Standalone C remains zero;
79 production C preamble bodies remain before this batch.

## Delegation and recovery

One GPT-6 Luna helper produced the 30 registration bindings and reviewed existing
fixture coverage. Primary independently checked every name, C address and Go
wrapper signature/call. Primary owns the API, fixtures, integration and acceptance.
The first inventory blurred the outer default shortcut with the mutable legacy
handler; that was corrected before implementation. The final review found all
27 nondefault callbacks covered; its version-gate wording was clarified because
rejected records enter the callback and return at its version check.

Current artifacts are in `build/port-xfer-sound-registry`; initial ignored drafts
and audits remain in `build/port-collision-registry/next-xfer-sound-*`. Production
and tests must remain fixed while a Go job runs. New production conversion must
wait for a qualified, committed original baseline.

## Original-path progress

All eight new roots passed in original-v2 before production conversion. The first
attempt stopped during compilation because primary's new bridge omitted the
GAME4_3 declaration header; adding the test-only include resolved it. No runtime
expectations were changed. The combined 228-root original selection subsequently passed
in default/server/highres with all 220 frozen capture groups unchanged.

With no Go jobs active, primary verified Luna's second bounded cache plan and
removed 30 old single-link repository compiler archives, reclaiming 1,442,158,216
bytes. Exact stat/hash/archive/module and host process/file-use checks passed.
The strict cutoff was 2026-09-23 22:00 UTC. Original assets, module cache, source,
binaries and qualification evidence remain; older rebuilds may take longer.

## Qualified original baseline

All 228 exact roots pass in each profile without skips. All 220 frozen capture
groups match across original default/server/highres runs. The new contracts add
36 raw transfer forwarding/return cases, 12 successive default-handler cases,
128 damage-owner slot/return/source/weapon cases, 48 real audio-event cases,
two exact default record round trips and two detailed outer-error cases.
Existing record matrices cover all 27 nondefault transfer registrations; the
new identity test confirms their real callback slots.

Only five new porttest files differ from the preceding qualified collision
checkpoint. All other source/dependency fingerprints and four binary hashes are
identical; all 30 callback exports remain in those binaries. Production evidence
is reused only at this test-only checkpoint. See
[original qualification](xfer-sound-registry-c-qualification.json).
The production patch is still unapplied; conversion qualification follows.
