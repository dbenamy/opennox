# Complete port-test corpus

Status: qualified. The complete corpus passes in all three production profiles.
Engine source and frozen expectations are unchanged; the corrected fixture is qualified.

The old accumulated pattern selected 1,527 of 2,426 default/highres roots and
1,523 of 2,415 server roots. Independent compiled-name and source audits establish
that all omitted roots are porttest-tagged. Earlier focused qualification remains
valid for its recorded scope; the former accumulated run was not the entire corpus.

Use `^Test` to select every compiled root automatically. Keep affected-owner
patterns for focused batch tests. Run profiles sequentially with original assets,
null audio, a 768 MiB runtime heap limit and a 3,600-second package timeout.
Review exact discovered/executed/completed name sets and every skip. Only the
standalone `TestMapPopulationPrerequisiteProbe` diagnostic may skip; its regression
owner invokes its cases. `TestSafeMemoryBridges` requires the separate safe profile.
Existing frozen expectations remain unchanged.

Luna reviewed missing-root prerequisites. Primary independently computed selection
using the complete regex, correcting the initial helper literal-alternative count.
Existing manifest asset variables cover the newly included tests; no ambient
capture-output variables are set. Broader inventory algorithm design stays with
the primary. No additional production rebuild is needed for a selector-only change;
reuse `9adcf3d3` production evidence after verifying that the corrected test fixture
is the sole source difference.

Inventory: [complete-port-corpus-inventory.json](complete-port-corpus-inventory.json).
Manifest: [complete-port-corpus-batch.json](complete-port-corpus-batch.json).
Original run: `build/port-complete-corpus/contracts/`. Corrected full run:
`build/port-complete-corpus/contracts-valid-pointers/` (server/highres/default order).

## Player-reset fixture correction

The first complete default run passed all 2,426 roots (2,425 pass, one diagnostic
skip). Server passed 1,865 roots and skipped that diagnostic, then crashed during
`TestServerOrchestrationPlayerReset`; highres had not started. The failing root was
one of the omitted tests. Runtime reported `0x76543210` as an invalid Go pointer
while flushing a write barrier. The fixture put that integer in `Object.Obj130`,
which production clears through a typed pointer assignment. Repository consumers
use that field as an object pointer, not an integer/pointer union. An isolated
20-repeat run with `GOGC=1` passed on unchanged source, so that focused run alone
does not reproduce or rule out the heap-state-dependent fault.

Seed live, non-null object pointers in Obj130, PlayerUpdateData.Field29 and Field77.
Keep numeric seeds in numeric fields, all 40 input combinations, full byte-write
assertions and the original frozen hash. Do not bypass Go's write barrier in the
engine to accommodate invalid fixture input. Primary traced the crash; Luna's
bounded field-layout review independently confirmed pointer versus numeric fields.
Run 20 repeats per profile under frequent GC, static checks, then a fresh complete
sweep. Original failure logs remain under `build/port-complete-corpus/contracts/`;
corrected runs use distinct output paths. The earlier default result does not
qualify the changed fixture. Production evidence can be reused only after verifying
that this test file is the sole source difference from `9adcf3d3`.

The first corrected stress attempt mistakenly applied `GOGC=1` to the compiler
as well. Default completed all 20 repeats/800 subcases with the original hash;
primary stopped the subsequent server compiler and retained the failed-attempt
record. The replacement manifest compiles with normal GC and applies `GOGC=1`
only to the test binary. Reuse the completed default stress result by exact source
identity; run server/highres under `fixture-runtime-gc/`. No engine failure is
inferred from the explicitly interrupted compiler attempt.

Focused qualification is complete: 20 root repetitions and 800 subcases pass in
each profile at `GOGC=1`, with no skips/failures and the unchanged frozen hash.
Static checks pass. Source fingerprints differ from qualified production only in
`server_orchestration_player_porttest_test.go`. Evidence:
[player-reset-fixture-qualification.json](player-reset-fixture-qualification.json).
The fresh complete sweep also passes: 2,425 pass plus one diagnostic skip in each
of default/highres; 2,414 pass plus that skip in server. All discovered, executed
and completed root-name sets exactly match the inventory. There are no failed
events or other skips. Only the fixture differs from qualified production source,
so production evidence from `9adcf3d3` is reused; no fresh production build is claimed.
See [complete-port-corpus-qualification.json](complete-port-corpus-qualification.json).

Initial default/server raw logs were losslessly gzip-archived after host-use and
hash checks. Restore using `build/port-complete-corpus/initial-log-archive.json`
before tools that require their original `.jsonl` paths. The recorded failure
evidence and SHA256 values are unchanged.

The completed corrected full-sweep logs are also losslessly archived; restore
commands and hashes are in `build/port-complete-corpus/qualified-log-archive.json`.
