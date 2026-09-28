# Shared fixture types and callers

## Scope and review

Twenty porttest files migrate Go-only adapter signatures and their callers from
C scalar/vector types to explicit-width Go types. Production behavior, casts,
evaluation order, pointer words, allocations, assertions and captures are unchanged.
Seventeen files lose their C import. Edge-mapping, prefab and worklist fixtures
retain their real C allocation/grid observer calls and matching argument types.

Target GCC checks confirm signed eight-bit char, 32-bit int/uint, 16-bit text words,
and float2/float4 size, alignment and field offsets. Pointer-only vector adapters
use two-/four-float Go arrays with the same layout. Removed int-width C assertions
are replaced by explicit int32/uint32 types; this does not qualify other targets.
The 18-header removed-include closure has no startup hooks. AST review finds
73 changed functions, 45 unchanged, none added or removed.

The shared file has 52 helpers with references in 43 other files. A scoped audit
of fixture entrypoints and root-test helper references selects 203 roots across
30 entrypoints. The server build omits two client hover tests, leaving 201 compiled
owners. A broad unqualified identifier graph over-selected common names and was
not used as the coverage claim. Existing entrypoint switches conservatively select
all their consumers, including cases whose branch is unchanged.

## Qualification plan

Run all selected original owners freshly: 203 default, 201 server, 203 highres;
repeat 22 focused roots per profile and six safe contracts. Reuse only the preceding
verified same-source binaries and unchanged production baseline. Commit accepted
original evidence before installing the draft.

After conversion, compile every root in all three profiles and run those exact
owner sets, focused repeats and safe contracts. Run static checks, safe build,
production builds/ABI, exact known-suite outcomes and fresh save/load against
unchanged references. The preceding cleanup retains broad default coverage; the
callback batch retains the latest broad three-profile sweep. This batch scopes
runtime tests to affected fixture interfaces and their consumers.

The previously excluded TestEdgeNormalizationThresholds passed freshly on the
original path and is included here. An initial launcher used the wrong pattern
argument and failed before executing a test; the corrected launch passed. Initial
selection also assumed equal profile counts; discovery identified the two
server-excluded hover tests before any owner sweep began.

## Status

Original baseline accepted: 203/201/203 owners, 22 focused repeats per profile
and six safe contracts all passed; draft not installed. Expected fixture C imports fall
51→34; production remains 4 client/highres and 5 server. Standalone C, embedded
production bodies and legacy exports remain zero. Header counts remain 157/2,731.
Luna quota remains unavailable; primary performed implementation and review.
Artifacts: `build/port-fixture-native-types/`.
