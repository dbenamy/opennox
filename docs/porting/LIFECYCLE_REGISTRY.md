# Object creation and initialization callbacks

## Scope and compatibility

Eight creation and twelve initialization registrations now dispatch through Go,
binding complete existing export wrappers and preserving names, callback addresses,
data sizes and parser order. Monster/Shopkeeper share one callback but keep distinct
sizes; both replaceable monster handlers remain resolved at invocation.

Creation runs after ordinary allocation and template copying. Existing CallInit
keeps its one-argument raw fallback and nil no-op. CallInitWithArg preserves the
two-pointer ABI of pending activation, player arrival and item respawning, keeping
the configured-slot precondition and existing caller guards. The respawn owner
loads Init by offset 688; named-field searches alone missed that caller.

## Qualified original baseline

New test-only fixtures verify all 20 registration identities and sizes. Registered
CallInit siblings reuse the nine reward initializer matrices and player initializer
matrix, discarding only the return that this void owner ignores. The existing
algorithm tests and frozen expectations remain unchanged. Separate tests exercise
the real pending-object owner, nil/raw behavior, exact argument forwarding and two
successive replacements of the shared Monster/Shopkeeper handler. Positive monster
initialization checks cover health scaling/history, position/direction, speed and
status flags. Allocator-based creation contracts now exercise all 8 registrations, nil/raw fallback and two successive MonsterCreate handler replacements. Positive checks distinguish weapon ammunition setup from armor durability.

The first capture 73701 passed five roots and failed the player sibling: primary
copied a byte-return assertion from the direct helper matrix although CallInit is
void. Only the new sibling now clears that expectation; the original return test
remains. No production behavior or frozen expectation was changed. Failed output
remains under build/port-lifecycle-registry/original*. Corrected init capture 19589 passed 8 roots and 344 observations, now frozen. Combined creation/init capture 25059 passed 10 roots. All 13 groups / 352 observations are frozen. First broad default 80742 completed 209 roots with no failures, but the strict runner rejected the opt-in TestMapPopulationPrerequisiteProbe skip. The explicit 208-root selection now excludes only that standalone diagnostic; no source or expected behavior changed. Rejected logs remain baseline-v1-skipped-probe. Corrected strict89145 passes all 208 roots in default/server/highres, with no skips or changed observations. Seven new and two changed porttest files are the only source changes; all other source/dependency fingerprints and four preceding binary hashes match. Production evidence is reused only for this test-only baseline. See [baseline qualification](lifecycle-registry-c-qualification.json).

The new two-argument API has no original method to call. A test-only bridge starts
with the original raw two-pointer invocation and switches to CallInitWithArg
with production conversion. Its nil, same-object and distinct-object argument
cases remain identical across the switch. Real caller tests additionally cover
pending activation, arrival and item respawning, so this bridge is not the sole
qualification of production routing.

## Delegation and recovery

One GPT-6 Luna helper drafts the production mapping and bounded creation fixture;
primary owns initialization contracts, integration and qualification. Primary
corrected the inventory's omitted mutable-handler variables, enclosing pending
owner name, and raw-offset respawn callsite. Static mapping checks alone are not
acceptance evidence. No measured subscription-cost savings are claimed.

After the first capture finished, primary removed 26 old reproducible repository
compiler archives totaling 1,646,054,562 bytes. Hash/stat/archive/module-marker and
host FD/executable/mapped-file checks passed with no Go build active. This is cache
eviction; old revisions may rebuild more slowly. Source, module downloads, assets,
binaries, original captures and qualification evidence remain. Exact audit/removal
records are under build/port-lifecycle-registry; the removal script is consumed.

Standalone production/test-reference C remains zero files/physical lines.
At the original-baseline checkpoint, production C preamble bodies remained 79;
production conversion had not started.

Creation-fixture review found a duplicate import, uninitialized modifier/health
values (`alloc.New` does not copy its argument), incomplete copied-buffer cleanup,
and incomplete global restoration. Primary corrected these before compilation,
restored 13 cache words, retained generator update data through FreeObject, and
freed saved non-owned buffers afterward. Known monster pointer slots must be nil
in this empty-definition fixture and are asserted as such; no pointer words are
silently erased from its capture. The helper changed its ignored draft during
handoff; primary caught an overlapping cleanup helper and removed a potential
double-free before running. Keep delivered artifacts stable during integration.
Allocation-heavy fixture ownership remains with primary pending better evidence.

First conversion 59981 failed during discovery because root object.go retained its
now-unused ccall import. Primary removed that import; original attempt remains
contracts-v1-unused-import. A selector-use audit now checks every changed file's
imports before retry 92843. No runtime test or frozen expectation changed. Both
the helper draft and primary review initially missed this import cleanup.

## Completed conversion qualification

All 208 affected roots pass without skips in default/server/highres. Thirteen
frozen groups retain all 352 observations, and the 300 explicit CallInit execution
counts repeat unchanged. Both raw ABIs, nil/cleared slots, real pending activation,
player arrival/respawn consumers and late-bound handlers pass. Six production
files change; only the new two-argument test bridge changes alongside them.
No original algorithm golden changes.

Safe build/static checks pass. Four fresh 386/SSE2/CGO binaries retain the required
Go-backed C exports and contain no test helpers. Full-suite comparison retains
exactly 304 known failure events and17 pass/2 fail/32 skip package outcomes. Headless
creation and explicit save/load/resume pass. See [conversion qualification](lifecycle-registry-qualification.json).

Original callback export identities and generic C fallback remain. Standalone C
remains zero files/physical lines; production C preamble bodies remain 79. No 64-bit
support or measured performance improvement is claimed. Completed scenario asset
copies were deduplicated against original hashes; per-run restore manifests and
ignored deduplicate-preflight/save.py --restore remain available.
