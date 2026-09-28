# Fixture raw allocation consolidation

## Scope and ownership

Twenty-five porttest files route direct C.calloc/C.free calls through the existing
legacyCalloc/legacyFree adapters. Normal builds retain raw libc allocation; safe
builds retain alloc.Calloc/FreePtr tracking, zero-size behavior and invalid-free
checks. Allocation sizes, release order, failure checks and ownership remain
profile-specific. The uintptr casts retain the target size_t width. Fourteen
fixtures lose their C import; remaining scalar types, aligned_alloc, CString/CBytes
and real C observers stay in place. Production allocator implementation is unchanged.

The linker calloc/free wrappers still apply to the centralized allocator's C
stubs. Preserve thread pinning, observer activation/stop order and cross-fixture
releases. Resource observer storage is allocated before pinning and freed after
observation stops. Existing C strings, aligned storage and native-owner buffers
retain their existing per-profile release behavior. AST review: 32 changed
functions, 134 unchanged, none added or removed.

## Coverage

A first function-only graph missed method owners. Receiver-type references cover
list close, map-catalog add and network-alias allocation/release. The shared pool's
controlsPrepare method returns immediately when Controls is nil; select consumers
of PortTestPlayerControlsSpec rather than every unrelated shared-pool fixture.
The reviewed selection contains 622 compiled roots in each profile, plus 24 focused
repeats and six safe contracts. Include allocation failures, release observations,
thread scope, list ordering, map/catalog owners, protection records and bot cleanup.

The previously excluded rejected-list release contract passed freshly on the
original path and is included. Edge thresholds remain included. Omit only the
standalone population prerequisite diagnostic, which intentionally skips without
an environment selector: its four cases run through the selected
TestMapPopulationPrerequisiteRegressions root. No captures are regenerated.

Run exact owner selections in default/server/highres, safe contracts/build, static
mapped-state checks, production builds/ABI, exact known-suite outcomes and fresh
save/load. Commit the accepted original baseline before installing the draft.
All target binaries compile every root, even though runtime selection is scoped.

## Qualified result

All gates passed. Exact original and native selections passed 622 roots in each
profile, with 24 focused repeats per profile and six safe contracts, without skips.
Every target test binary compiled all roots. Static checks, safe build, production
builds/ABI, exact known-suite outcomes and fresh save/load passed. Allocation and
release observers retained their expected results; original assets and all test
assertions/captures are unchanged.

Fixture C imports fell 34→20. Production remains 4 client/highres and 5 server,
with zero embedded production C bodies, legacy exports and standalone C. Headers
remain 157 files /2,731 physical lines. The centralized allocator still uses libc;
this consolidation preserves normal and safe implementations and their distinct
bookkeeping.

Evidence: [baseline](fixture-raw-allocation-baseline.json),
[qualification](fixture-raw-allocation-qualification.json),
[inventory](fixture-raw-allocation-inventory-after.json).
Luna quota remained unavailable; primary handled review and qualification.
Artifacts: `build/port-fixture-raw-allocation/`.

## Safe-routing correction

The initial draft called RawCalloc/RawFree directly. Normal owner sweeps passed,
but safe TestProtectionHandlesABI caught an incorrect-free panic: cgo_safe.go
redirects legacy-package C allocation calls to tracked nox_calloc/nox_free via
compiler `-D` flags. A header-only macro audit missed those package-wide settings.
The existing legacyCalloc/legacyFree adapters preserve both normal and safe paths.
All 25 files were corrected without changing original assertions or captures.
Qualification restarted with safe contracts, then passed every gate on the
corrected source under `corrected/`. The initial results are retained separately.

## Scenario artifact recovery

After qualification, 1,654 byte-verified original-asset copies were removed from
the inactive scenario, reclaiming 559,869,952 allocated bytes. Original assets,
run outputs and restoration records remain. Restore with:

```sh
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/fixture-raw-allocation-save/deduplicated-assets.json
```

Keep the initial raw-routing attempt separate from accepted binaries under
`build/port-fixture-raw-allocation/corrected/` when reusing evidence.

The qualified root/safe/production binaries for this batch were later removed
after source and newer replacements were verified. Rebuild d2df2eb6; see
[string-copy cleanup](FIXTURE_STRING_COPIES.md#artifact-cleanup-and-recovery).
Recorded results and reports remain available.
