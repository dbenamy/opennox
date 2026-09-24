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


## Qualified conversion

All 222 focused roots pass in each default/server/high-resolution profile without
skips. Safe/static checks, three production/ABI builds and headless character
creation with explicit save/load/resume pass. The full asset suite exactly matches
the known baseline: 304 failure events, 17 passing/two failing/32 skipped packages.
All 1,654 original asset hashes and existing assertions/captures remain unchanged.
The preceding item batch's full default corpus passed 2,441 roots plus one known
diagnostic skip; it was not repeated for this unchanged registry implementation.

All accepted phases have identical source fingerprints. All seventeen changed/new/
deleted files match review, all 30 exports are absent, and retained exported
function signatures and bodies are unchanged. Production cgo files fall 208→202
(261/463 eliminated); selected legacy exports fall 825→795 (1,095/1,890 retired).
Headers remain 157 files with 3,524 physical lines; embedded production callback
bodies remain 77. Standalone production and test-reference C remain zero.
External native-library bindings are unchanged.

Luna produced the bounded overlay; primary reviewed every owner mapping, the
default hook/logging body, all key/getter consumers and fixture normalization.
Case-sensitive export/API audit corrections were made before installation. Primary
removed the final redundant int32 C cast in the map-painting consumer.
The first compile found a missing stdint.h include in the raw C observer after
old engine headers were removed; primary added the explicit standard header.
Accepted contracts are under contracts-fixed. Qualification found no behavior
differences and required no assertion or capture changes.

[Qualification](xfer-identities-qualification.json),
[inventory](xfer-identities-inventory-after.json).
