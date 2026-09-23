# Registered object updates

## Scope and design

Move all53 registered object-update callbacks through a typed Go registry while
retaining their C addresses, data extents, parser ordering and raw fallback.
Object.CallUpdate retains its nil-slot no-op and ignores callback return values.
The proposed bindings call the existing Go-export wrapper directly, preserving
its complete body and all argument conversions without crossing through C.
Six replaceable handlers remain resolved when invoked. Damage sounds are deferred
so this larger batch has one shared dispatch boundary.

The production patch remains unapplied in build/port-update-registry. Independent
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
[baseline qualification](update-registry-c-qualification.json). Conversion is pending.

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
