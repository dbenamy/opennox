# Native transfer identities

## Original baseline and scope

All 222 selected roots pass in default, server and high-resolution profiles,
without skips, against qualified item checkpoint `00044175` plus one new test file.
Production fingerprints match that checkpoint exactly. Reuse its production/ABI,
known-suite and headless save/load evidence as the production baseline.

Replace all 28 transfer callback addresses with distinct static native identities,
keeping the existing Go registry and unknown-address C fallback unchanged. Retire
30 C exports: 28 callbacks plus two unused lowercase common/inventory wrappers.
Their distinct capitalized Go helpers remain, as do the original diagnostic strings.
Primary whole-source scanning found no callers for the two extra wrappers.

Preserve all 27 nondefault native owners and the default mutable hook, current
cryptfile argument, logging and 0/1 result. Migrate four getters, type defaults,
shop reward comparisons, map exit/door/reward classifiers and prefab script event
and reference selection. Three fixture identity tables retain their existing
normalization IDs; raw C observers and floating-point control-word probes remain.

New original-path contracts check all 28 identities in foreign ObjectType/Object
storage across GC/stack growth, getters/defaults, and actual invisible-light
script-storage allocation across four game-flag combinations. Existing focused
roots cover registry/raw dispatch, object/item/creature serialization, map
painting/population, prefab scripts, shop, player files and orchestration. The
established MapPopulationPrerequisiteProbe diagnostic is intentionally excluded;
the preceding full item corpus confirms its expected skip.

The map-painting consumer's remaining C.int cast only wraps an int32 result.
Remove that redundant roundtrip and its cgo import within this tested boundary.
No algorithm, layout, registry API or ownership changes are intended. This
reversible cleanup removes another cgo dependency without changing result bits.

One GPT-6 Luna helper drafts the bounded production/fixture overlay; primary
owns contracts, exact owner/identity mapping review and acceptance. The helper's
initial audit needed correction for case-sensitive export names and existing API
capabilities. Primary verified every native-owner mapping directly against the
original wrappers; prose alone is not evidence of ABI retirement.

After conversion run all focused profiles, safe/static, three production/ABI
builds, exact known-suite comparison, headless save/load and asset hashes. The
full default port corpus is the preceding item checkpoint; repeat it if a failure
or review widens the affected scope. Existing assertions and captures stay fixed.

[Original evidence](xfer-identities-baseline.json),
[commands](xfer-identities-batch.json), [selection](xfer-identities-tests.txt).
Local artifacts: `build/port-xfer-identities/`.
