# Fixture raw allocation consolidation

## Scope and ownership

Twenty-five porttest files route direct C.calloc/C.free calls through the existing
alloc.RawCalloc/RawFree functions. Those functions still call the same libc APIs,
without allocs registration or safe-mode bookkeeping. Preserve zeroing, nil and
overflow handling, allocation sizes, release order, allocator domain and failure
checks. The uintptr casts retain the target size_t width. Fourteen fixtures lose
their C import; remaining scalar types, aligned_alloc, CString/CBytes and real C
observers stay in place. Production allocator implementation is unchanged.

The linker calloc/free wrappers still apply to the centralized allocator's C
stubs. Preserve thread pinning, observer activation/stop order and cross-fixture
releases. Resource observer storage is allocated before pinning and freed after
observation stops. Existing C strings, aligned storage and native-owner buffers
continue to be freed by libc. AST review: 32 changed functions, 134 unchanged,
none added or removed.

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

## Status

Original baseline accepted: 622 roots and 24 focused repeats per profile, plus
six safe contracts, passed without skips; draft not installed. Expected fixture C imports: 34→20.
Production remains 4 client/highres and 5 server, with zero embedded production C
bodies, legacy exports and standalone C. Headers remain 157 /2,731 physical lines.
Luna quota remains unavailable; primary handles review and qualification.
Artifacts: `build/port-fixture-raw-allocation/`.
