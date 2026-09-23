# Registered object updates

## Scope and design

Move all53 registered object-update callbacks through a typed Go registry while
retaining their C addresses, data extents, parser ordering and raw fallback.
Object.CallUpdate retains its nil-slot no-op and ignores callback return values.
The bindings call the existing Go-export wrapper directly, preserving
its complete body and all argument conversions without crossing through C.
Six replaceable handlers remain resolved when invoked. Damage sounds are deferred
so this larger batch has one shared dispatch boundary.

The three-file production patch is applied after baseline90f2a227. Independent
review verifies all53 names, addresses, sizes, ordering, wrapper identifiers and
argument types against the actual source. C-int object words remain32-bit; no
new64-bit support is claimed. Existing wrappers/exports/layout remain unchanged.

## Original owners and contracts

The new fixture selects actual resource-name registrations, independently checks
address/data-size identity, installs the callback in Object.Update and invokes
Object.CallUpdate. Existing direct-helper tests and goldens remain unchanged.
Forty capture roots pass:36 frozen groups containing8251 owner cases,53 registration
identity checks, raw nil/forward/cleared-slot contracts and two successive handler
replacements for each of six mutable callbacks. Execution records show12188 calls
across52 names; existing TestUnitGameplayUndeadUpdate exercises the remaining name
through its already-registered owner. All281 baseline roots pass in default/server/highres with no skips;36 hashes and
exact execution counts repeat unchanged. Nine new porttest files and four existing
porttest adapters are the only source changes relative to0a94c36f. All other source
and dependency fingerprints and four preceding binary hashes match; all53 exports
remain. Production qualification is reused for this test-only baseline. See
[baseline qualification](update-registry-c-qualification.json). Conversion14362 passes all281 roots/profile with unchanged captures, execution counts and dynamic/raw contracts. Fresh safe/static, production/ABI, exact known-suite comparison and headless creation/save-load also pass. See [conversion qualification](update-registry-qualification.json).

Shared temporary/world/objective adapters cover42 callback names. Generator gate
and spawning matrices include12 independently asserted positive spawns. Dedicated
mover, sentry and shooting-trap siblings preserve movement, list, projectile,
audio, cooldown and guard assertions. Their previous helper return assertions
are deliberately absent from the new void-owner route; old tests retain them.
The invisible-teleport sibling checks activation clearing and powered destination
coordinates, alongside its frozen observations.

The first capture60903 ran42 roots and failed three execution-count assertions:
remember-owner, spark-creation and float-filter matrices never call a registered
update. Removed only these redundant new siblings and added the missing invisible
pentagram route. Original gameplay assertions did not fail. Corrected77044 passed
before any hash was frozen. Failed logs/captures are retained separately.
A primary draft-generator replacement briefly renamed PortTest selectors too;
text review corrected it before compilation. No production changes were involved.

## Delegation and recovery

Luna drafted the53 bindings and three motion siblings; primary owns shared-owner
routing, generator/dynamic/raw contracts, acceptance and qualification. The first
inventory search missed handler assignments outside legacy; whole-src search finds
all six in legacy_exports.go and a Monster override in player_controls_porttest_test.go.
No callback-variable absence is inferred from a subtree-only search.

After the first capture finished,40 older reproducible repository compiler
archives were removed following stat/hash/archive/package-marker and host-use
checks, reclaiming2168958680bytes. Source, module cache, assets, binaries and
qualification/capture evidence remain. This is ordinary cache eviction; historical
rebuilds may take longer. Exact plan/removal records are under build/port-update-registry.

Standalone C remains zero physical files/lines; production C preamble bodies79.

## Completed qualification and limits

Only three production files changed. All four fresh binaries retain the selected
Go-backed C exports and contain no porttest helpers. The known full-suite failure
multiset remains exactly304 events, with17 passing,2 failing and32 skipped packages.
Headless creation and explicit save/load/resume succeed. No conversion correction
or golden changes were needed after the qualified baseline.

This removes runtime C round trips for registered updates; C export identities
and the generic fallback dispatchers remain. No measured speedup or64-bit support
is claimed. Current configured target remains386/SSE2 with CGO.
Luna's production and motion-fixture drafts passed primary mapping/body/contract
review and qualification; the earlier inventory scope correction is recorded above.
No measured subscription savings. Keep one bounded helper with primary acceptance.

Completed preflight/save scenario assets were deduplicated against original hashes.
Per-run restore manifests and ignored deduplicate-preflight/save.py --restore
remain available; original assets, results, logs and binaries are retained.
