# Object creation and initialization callbacks

## Scope and compatibility

Qualify eight creation and twelve initialization registrations together. Production
conversion is not applied yet. The candidate binds each complete existing Go export
wrapper while retaining its resource name, C address, data size and parser order.
MonsterInit and ShopkeeperInit share one address but retain distinct data sizes.
Monster creation and initialization resolve their replaceable Go handlers when called.

Creation dispatch occurs inside the real object allocator, after template copying.
Initialization has three live two-argument owners: pending-object activation,
player arrival and respawn-item creation. They pass the object and nil. The separate
CallInit method uses one argument and has a nil-slot no-op. The candidate preserves
both raw calling conventions explicitly; its new CallInitWithArg method requires a
configured slot, keeping the existing caller guards. No object layout changes.
The third owner accesses the Init field by offset688; a named-field-only search
missed it. Primary found it while tracing the player initializer's item creation.

## Qualified original baseline

New test-only fixtures verify all20 registration identities and sizes. Registered
CallInit siblings reuse the nine reward initializer matrices and player initializer
matrix, discarding only the return that this void owner ignores. The existing
algorithm tests and frozen expectations remain unchanged. Separate tests exercise
the real pending-object owner, nil/raw behavior, exact argument forwarding and two
successive replacements of the shared Monster/Shopkeeper handler. Positive monster
initialization checks cover health scaling/history, position/direction, speed and
status flags. Allocator-based creation contracts now exercise all8 registrations, nil/raw fallback and two successive MonsterCreate handler replacements. Positive checks distinguish weapon ammunition setup from armor durability.

The first capture73701 passed five roots and failed the player sibling: primary
copied a byte-return assertion from the direct helper matrix although CallInit is
void. Only the new sibling now clears that expectation; the original return test
remains. No production behavior or frozen expectation was changed. Failed output
remains under build/port-lifecycle-registry/original*. Corrected init capture19589 passed8 roots and344 observations, now frozen. Combined creation/init capture25059 passed10 roots. All13 groups352 observations are frozen. First broad default80742 completed209 roots with no failures, but the strict runner rejected the opt-in TestMapPopulationPrerequisiteProbe skip. The explicit208-root selection now excludes only that standalone diagnostic; no source or expected behavior changed. Rejected logs remain baseline-v1-skipped-probe. Corrected strict89145 passes all208 roots in default/server/highres, with no skips or changed observations. Seven new and two changed porttest files are the only source changes; all other source/dependency fingerprints and four preceding binary hashes match. Production evidence is reused only for this test-only baseline. See [baseline qualification](lifecycle-registry-c-qualification.json).

The new two-argument API has no original method to call. A test-only bridge starts
with the original raw two-pointer invocation and will switch to CallInitWithArg
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

After the first capture finished, primary removed26 old reproducible repository
compiler archives totaling1,646,054,562 bytes. Hash/stat/archive/module-marker and
host FD/executable/mapped-file checks passed with no Go build active. This is cache
eviction; old revisions may rebuild more slowly. Source, module downloads, assets,
binaries, original captures and qualification evidence remain. Exact audit/removal
records are under build/port-lifecycle-registry; the removal script is consumed.

Standalone production/test-reference C remains zero files/physical lines.
Production C preamble bodies remain79; no conversion progress is claimed yet.

Creation-fixture review found a duplicate import, uninitialized modifier/health
values (`alloc.New` does not copy its argument), incomplete copied-buffer cleanup,
and incomplete global restoration. Primary corrected these before compilation,
restored13 cache words, retained generator update data through FreeObject, and
freed saved non-owned buffers afterward. Known monster pointer slots must be nil
in this empty-definition fixture and are asserted as such; no pointer words are
silently erased from its capture. The helper changed its ignored draft during
handoff; primary caught an overlapping cleanup helper and removed a potential
double-free before running. Keep delivered artifacts stable during integration.
Allocation-heavy fixture ownership remains with primary pending better evidence.
