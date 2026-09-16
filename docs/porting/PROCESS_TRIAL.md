# Two-round process trial: results and adopted process

Both rounds are qualified. The user reviewed the trial and **adopted the revised
process for continued successive batches**, lifting the trial's scheduled pause. The trial removed **1,290 C lines**, leaving **71,252 physical C
lines in 90 files**, with no test-reference C. These counts measure source removed,
not remaining effort or a percentage of the port completed.

## Assessment

Keep the revised process. Explicit package selection, coherent batches, separate
recovery commits and qualification gates, and boundary-oriented contracts all
proved useful. The larger map-decoder scope was manageable as a single batch;
both Go implementations matched their frozen C behavior on their first actual
behavioral run. This is encouraging evidence, not a general guarantee.

The main remaining fixed cost is the accumulated corpus. Its default run alone
took 765.871s, while the map decoder's first focused Go run took 4.613s.
The combined three-target sweep took 2,108.020s, about 35 minutes. Continue focused/affected checks during implementation and
full qualification at coherent round/subsystem boundaries. Do not pay for a full
accumulated run at every recovery commit. Keep a broader milestone after related
rounds, shared infrastructure changes or uncertain failure scope.

Prefer complete owners or APIs, aiming for roughly 1–3k C lines when dependencies
permit. The audio round deliberately completed an existing small scope to establish
the runner. The 969-line decoder then exercised larger scope and a non-root package.
Do not turn line count into a quota or combine unrelated code just to hit it.

Keep the common runner, but its JSON manifests are still verbose: commands, tags
and hash lists repeat. A later small cleanup can reference shared capture manifests
and target lists. Avoid expanding the harness before another concrete need appears.
No subagents were used in this trial; it provides no measured model-cost or
subscription-limit comparison.

## Measured gates

| Round | C removed | Frozen records | Affected roots per target | Final C qualification | Native qualification attempts |
| --- | ---: | ---: | ---: | ---: | ---: |
| Client audio assets | 321 | 1,433 | 51 | 165.227s | 499.084s |
| Complete map decoder | 969 | 174 | 9 | 14.792s | 439.645s |

These are qualification command times, excluding fixture authoring and focused
development runs. The map native column includes both failed harness-verification
attempts and the successful correction. Its final attempt reused unchanged
build/ABI/full-suite evidence only after source, binary-hash and gate comparisons.
Map qualification overlapped the accumulated job; times are not additive.

A separate same-source warm-cache driver check ran the initial two audio tests
in 6.041s for the root package versus 12.707s across the tree. The earlier 187.797s
run included compilation and is not a comparable speedup baseline. Different
batch complexity, cache state, I/O and concurrent jobs prevent claiming an overall
speedup factor from these two rounds.

## What the checks found

- The audio real-catalog fixture needed the existing file-handle initializer and
  cleanup. One obsolete Go import also needed removal after deleting C calls.
- Empty compression input hits a pre-existing zero-allocation panic. The decoder's
  empty-output behavior and malformed-body errors are explicit, reversible changes
  recorded in DECISIONS.md; the compressor's empty-input limitation remains.
- Enumerating real filenames added three map pairs missed by the old directory-
  name convention. All 50 shipped compressed maps are now covered by the new tests.
- Screenshot success did not prove decoder coverage: the warrior map had no
  compressed counterpart. The new integration gate rejected that run. It now
  creates a compressed counterpart only in the disposable asset copy using the
  unchanged production compressor, then requires the game to regenerate the map.
- The game creates lowercase war01a.map. The first verifier expected the original
  War01A.map spelling; its corrected lookup verifies the actual created filename
  and bytes. Both normal and flat-floor replays then passed all 26 reference frames.

These were fixture, coverage and compile issues, not reasons to change frozen
valid-file expectations. Harness development cost time; its fail-closed checks
also prevented recording a replay as decoder coverage when it wasn't.

## Final regression evidence

Both rounds pass three ELF32/i386/SSE2/CGO production builds, their C-interface
audits and the exact known full asset suite: **1,553 failure entries**, with
**15 pass /3 fail /32 skip packages**. Existing failures remain visible.

The complete accumulated corpus ran across packages in all three variants:

| Target | Selected tests | Root-package tests | Driver seconds |
| --- | ---: | ---: | ---: |
| default | 1058 | 1047 | 765.871 |
| server | 1054 | 1043 | 717.866 |
| highres | 1058 | 1047 | 623.616 |

Each target skips only TestMapPopulationPrerequisiteProbe, an existing opt-in
isolated diagnostic requiring OPENNOX_POPULATION_PROBE. Neither new batch skips
any selected test. Source fingerprints remained unchanged during validation.
Fifteen tooling tests cover execution accounting, artifact hashes, source changes,
failed commands and invalid evidence reuse. See CLIENT_AUDIO_ASSETS.md and
MAP_DECOMPRESSION.md for batch-specific contracts, limitations and local artifacts.

## Reusable commands and recovery

Source `build/baseline/env.sh`, then use a reviewed manifest with a fresh directory:

```
python3 tools/porting/run_batch.py docs/porting/map-decode-batch.json \
  --phase native --out build/map-decode-review
```

Give scenario runs fresh names too; never overwrite earlier evidence.
A manifest declares `batch`, optional `env`, and ordered steps under `phases`.
Each step names an argument-array command, optional cwd and expected artifact
hashes. Strings support `{root}`, `{out}` and `{phase}`. The runner records
commands, timing, source identity and C LOC; it never edits source or expectations.
`run_tests.py --package .` is the default; repeat `--package` for affected packages.
Use `--package ./...` deliberately for the complete corpus and
`--require-no-skips` for required asset fixtures.

`qualify_production.py` reads build/ABI and scenario declarations from that same
manifest. Evidence reuse requires exact source fingerprints, binary hashes,
matching requested gates and a qualified known-failure suite. Failed attempts
remain separate; no implicit retry or golden regeneration is performed.

`run_scenario.py` is tracked and supports forced map expansion in fresh copies.
Build its small compressor helper from the src module when needed:

```
source build/baseline/env.sh
(cd src && go build -p 2 -o ../build/port-map-decompression/map-compress \
  ../tools/porting/map_compress.go)
```

The helper calls the existing production CompressFile. Original assets stay
unchanged. The local verified asset-deduplication helper under build/baseline
is optional; without it, successful run copies stay intact. See RECOVERY.md for
restoring the environment and references. Do not rerun stale ignored apply/freeze/finalize
scripts. Consult PORTING_STATE.md and the committed manifests after session loss.
