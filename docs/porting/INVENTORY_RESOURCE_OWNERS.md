# Equipment, inventory and resource owners

## Scope and original baseline

Starting from qualified `f2088976`, remove 61 engine C exports: 32 equipment,
seven inventory and 22 resource wrappers. Move fixture dispatch and the respawn
armor-mask caller to existing Go owners. Preserve the 22 inventory callbacks and
gold-pickup callback still used by C. Move the equipment pointer conversion helper
to native Go and update its 12 callers; remove the unused resource return helper.

All 263 selected root tests passed without skips in fresh default, server and
highres processes before conversion. The reused test binaries were checked against
identical source fingerprints from the preceding qualified batch. Frozen assertions
and expected captures are unchanged. See the [baseline](inventory-resource-owners-baseline.json),
[test selection](inventory-resource-owners-tests.txt) and [manifest](inventory-resource-owners-batch.json).

The 18-file draft was independently reconstructed from exact edit manifests and
formatted before installation. Whole-source scanning leaves only four historical
symbol comments for the retired names. No retired address is registered in fixture
identity maps, so this batch needs no new identity reservation or normalization.

Review covered signed 16-bit getter results, unsigned gold arithmetic, C-int input
narrowing, raw 32-bit pointer-shaped results, nil position propagation and exact
binary64 defense results. Independent modifier/drop/death callback fixtures remain.
Retained export signatures remain identical.

## Delegation and qualification

GPT-6 Luna drafted the equipment/inventory fixture changes and independently
reviewed the primary resource/helper changes. Primary reconstruction verified every
edit against the original bytes. Review requested restoration of untouched C-case
whitespace and removal of an unnecessary C object conversion in the position
calculation; both were corrected before acceptance. This bounded delegation was
useful; no cost savings are inferred from it.

Qualification passes: all 263 roots in each profile without skips, safe/static
checks, three fresh production binaries and ABI checks, exact known full-suite
comparison, headless character creation/save/load/resume and all 1,654 original
asset hashes. The known suite retains its existing 304 failure events (17 passing,
two failing and 32 skipped packages); this is not a wholly green project suite.
All phases have identical source fingerprints; all 18 converted files match the
reviewed bytes, and frozen expectations remain unchanged.

Selected cgo files fall from 230 to 229 (234/463 eliminated on net); selected
legacy C exports fall from 1,142 to 1,081 (809/1,890 retired). The 78 production
C callback bodies remain. Headers remain 157 files, now 3,804 physical lines.
Standalone production/test-reference C lines remain zero. External native-library
bindings are unchanged. See [qualification](inventory-resource-owners-qualification.json)
and [inventory](inventory-resource-owners-inventory-after.json).
Local artifacts: `build/port-inventory-resource-owners/`.
