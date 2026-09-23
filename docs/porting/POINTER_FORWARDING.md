# Direct Go pointer forwarding

Planned scope: 74 calls across 36 existing Go wrappers and 36 caller files.
Remove only the C selector prefix; arguments, casts, result handling and evaluation
order stay unchanged. No new export retirements are selected in this batch.
Standalone production/reference C stays zero lines/files; production C preamble
bodies stay 79. Headers, generated bridges and external libraries are excluded.

## Review

Primary compared every wrapper signature with freshly generated cgo metadata from
qualified revision `06bd85b6`, normalizing explicit Go/C aliases only. See the
[signature inventory](pointer-forwarding-signatures.json). No selected function
name has a macro definition. Independent reconstruction of Luna's patch matches
all 74 intended prefix removals byte for byte; every touched file still uses C.
The draft is not applied until the baseline passes.

Primary read all 36 wrapper bodies. Object arguments and returned references keep
the existing server owners. State-sync return pointers refer to object storage.
Point inputs are copied by movement/force helpers or consumed synchronously;
random placement writes the same caller output and returns its original pointer.
Object factories copy input C strings to Go strings before use; unit-name results
are interned. Tile names are compared synchronously. Information-message helpers
retain their caller-buffer mutations, and NetList.allocBufferRaw copies bytes to
owned message storage before return. No RNG, arithmetic, object lifecycle or
message-content change is intended. Target remains 386/SSE2.

## Baseline and qualification

Reuse the preceding 490/489/490 affected roots only after exact production/test
fingerprint, runtime environment and four-binary identity checks. Add 196 owner
roots covering AI main/callbacks/spells, shop/trade, gameplay text, world
mechanisms, map population, tile selection, object creation, effects, generators
and resources. All 196 source declarations use the porttest build constraint;
actual discovery and execution must agree in all three profiles.

The separately invoked TestMapPopulationPrerequisiteProbe intentionally skips
without a diagnostic selector; exclude that root and retain the prerequisite
regression test, which exercises all four probes. No fixture expectations change.
After conversion, the full combined 686/685/686 roots must pass without skips,
then safe/static, fresh production/ABI, exact known full-suite comparison and
headless creation/save-load checks. This is affected-owner coverage, not proof of
every caller branch. Verify all 36 redundant C-call bridges are absent from all
four new binaries while retaining required exports.

## Delegation and recovery

Luna drafted the exact-list patch, which passed primary byte-for-byte review.
Initial ownership notes included names from the previous batch; corrected before
acceptance. Its owner-test suggestions needed a broader second pass and primary
additions for AI main, generators/resources/effects. Keep mechanical drafts with
Luna and final scope/coverage/lifetime decisions with the primary. No measured
model cost or build-speed saving is claimed.

Artifacts: build/port-pointer-forwarding. direct-calls.patch is unapplied;
primary-review.md and patch-review.json record source review. All 196 additional contracts pass in each profile with no skips. Exact discovered
names, source/runtime-environment identity and retained binary hashes pass. See
[baseline evidence](pointer-forwarding-c-qualification.json). Native conversion
and production qualification remain pending.

The selected-test runner now separates discovery/build memory (1536MiB default,
configurable) from the unchanged execution budget (768MiB here). Nine Python
accounting/environment tests pass. Actual profile results record both settings;
existing runtime overrides remain intact. Luna drafted the minimal runner/test
change; primary added an explanatory comment and cleaned formatting. This is a
reversible compiler-GC improvement attempt, with no benchmark claim.

Six completed historical capture groups were losslessly archived after primary
provenance/link checks and host open-file verification. 400,569,409 raw bytes became
8,579,873 gzip bytes. All 107 hardlink paths and six symlink aliases are recorded
in capture-archive-record.json. Use archive-captures.py --restore before historical
finalizers; archival is CONSUMED. Assets, fixtures and binaries are preserved.
Cleanup briefly overlapped the additional baseline; timings are not benchmarks.
