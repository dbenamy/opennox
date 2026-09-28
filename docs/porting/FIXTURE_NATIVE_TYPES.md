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

## Qualified result

All gates passed without source corrections. Exact native owner sets passed:
203 default, 201 server, 203 highres; 22 focused repeats per profile and six safe
contracts passed without skips. All three test binaries compiled every root.
Safe build/static, three production builds/ABI, exact known-suite outcomes and
fresh save/load passed. Original assets and test assertions/captures are unchanged.
The previously excluded edge-threshold contract passed before and after conversion.

Fixture C imports fell 51→34. Production remains 4 client/highres and 5 server
across two project packages. Standalone C, embedded production bodies and legacy
exports remain zero. Headers remain 157 files /2,731 physical lines. External
native backends and the 32-bit target are unchanged.

Evidence: [baseline](fixture-native-types-baseline.json),
[qualification](fixture-native-types-qualification.json),
[inventory](fixture-native-types-inventory-after.json).
Luna quota remained unavailable; primary performed implementation and review.
Artifacts: `build/port-fixture-native-types/`.

## Qualification interruption and recovery

The first production gate built all three binaries and matched the exact known
suite, then exhausted disk while staging the fresh scenario assets. Its scenario
returned 120 and the outer result file was truncated. Source and expectations were
unchanged. After verified cleanup, production qualification reran successfully in
`production-retry/`, with a fresh `fixture-native-types-save-retry` scenario.
All 1,654 original assets retained their recorded hashes.

The initial evidence collector had a generated expected-count typo (62 instead of
22 focused preflight roots). Correcting that assertion accepted the already-passed
22 exact owner names; no source, test, selection or expectation was changed.

Obsolete binary cleanup removed 31 reproducible outputs (1,557,561,344 allocated
bytes), after matching committed source, retained qualified replacements and host
usage. Rebuild revisions `5cc464f4`, `1bb4b78a`, `64376c9f`, `89c91e08` and the
other revisions recorded in [cleanup evidence](fixture-native-types-cleanup.json) if needed.
Current baseline and native binaries remain.

The full disk initially prevented writing cleanup scripts. One verified asset copy reclaiming 122,011,648 bytes from the failed run was removed after printing its source/hash record;
that record was persisted as soon as space was available. Restore its `data/video.bag`
from the unchanged `build/assets/extracted/drive_c/Nox/video.bag`. Another 1,507
verified failed-run copies reclaimed 414,855,168 bytes; partial nonmatching files
and diagnostics remain. Their restoration command is:

```sh
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/fixture-native-types-save/deduplicated-assets.json
```

The accepted retry has its own restoration manifest:

```sh
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/fixture-native-types-save-retry/deduplicated-assets.json
```

Original assets remain untouched. Cleanup logs and precise binary revisions are
under `build/port-fixture-native-types/`.

After every build/test joined, 87 unchanged Linux386 Go cache archives older than
six hours were removed, reclaiming 3,806,949,376 bytes. Accepted retry deduplication
reclaimed 559,841,280 bytes from 1,654 verified copies. Current binaries, module
sources and original assets were preserved; physical free space is roughly 5 GB.

The qualified root/safe/production binaries for this batch were later removed
after source and newer replacements were verified. Rebuild 6048add7; see
[string-copy cleanup](FIXTURE_STRING_COPIES.md#artifact-cleanup-and-recovery).
Recorded results and reports remain available.
