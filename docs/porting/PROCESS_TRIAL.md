# Two-round process trial

The user authorized trying all five process recommendations for two completed
conversion rounds, then pausing to review the results. Round 1 starts from the
existing client audio asset fixtures; round 2 should select a larger connected
scope after its owner/caller audit. Do not inflate scope to meet a line quota.

## Changes being tried

- Focused tests select explicit packages (root by default), with discovery and
  execution accounting for every selected package. Whole-tree checks remain at
  qualified round boundaries. Required asset tests must not silently skip.
- Use coherent scopes, aiming for 1–3k C lines where shared ownership permits.
  Keep baseline and implementation recovery commits separate from expensive
  completed-round qualification. Mark intermediate evidence accurately.
- A small tracked manifest runner records commands, timings, source hashes,
  frozen-capture verification and C LOC. Use a fresh output directory per run;
  never overwrite failed evidence or regenerate expectations to match Go.
- Audit length/capacity and empty inputs; signedness/narrowing; owner lifecycle;
  partial failures and cursor position; surviving C callers; actual assets.
  Each report explains applicable boundaries and exclusions.
- Report lines removed separately from time, fixture/debugging effort, coverage
  gaps and remaining dependency complexity. No claim that LOC equals effort.

## Runner

Source `build/baseline/env.sh`, then use:

```
python3 tools/porting/run_batch.py docs/porting/BATCH.json \
  --phase c --out build/BATCH/c-repeat-a
```

A manifest has `batch`, optional `env`, and `phases` mapping names to ordered
steps. Each step has `name`, an argument-array `command`, optional `cwd`, and
optional `hashes` entries with `path` and `sha256`. Strings support `{root}`,
`{out}` and `{phase}` substitutions. Commands must exit zero, including separate
comparison commands that explicitly accept the exact known legacy failure set.
The runner checks source integrity after each step and writes a failure report
when a command or capture check fails. It performs no source edits or commits.
Batch-specific ownership and ABI checks remain explicit; consolidate proven
repetition rather than inventing a generalized fixture framework.

`run_tests.py --package .` is the default. Repeat `--package` for additional
actual dependencies with matching tests; use `--package ./...` only when intended.
`--require-no-skips` makes missing optional prerequisites fail qualification.

## Measurements and checkpoint

Baseline: initial two-root audio C-fixture run took 187.797 seconds using the old
whole-tree driver. This included compilation and is not a controlled benchmark.
Compare warm runs separately; record cache/source differences and avoid attributing
all elapsed-time variation to the runner.

Tooling setup and the first round are in progress. No trial round is complete.
The initial audio fixtures passed; audio production code is unchanged.
Pause after the second fully qualified conversion, even if more work is available.

Initial driver validation: six accounting tests pass (empty discovery, missing
execution despite exit zero, package-aware accounting, failed process, required
prerequisite skip and default package selection). Four manifest-runner tests
cover successful artifact verification, hash mismatch, source mutation and failed
commands. Actual warm-cache audio runs: root 6.041s; whole tree 12.707s, both two
roots completed without skips. Source was unchanged between these two runs.
These are single observations, not a cold-build benchmark.

Qualification cadence for this trial: each round has three-configuration affected
contracts, three production builds/ABI, the exact full asset suite, and relevant
reference gameplay. Run the entire accumulated port corpus at the end of round 2
as the combined milestone. Driver accounting itself is covered by dedicated
failure-path tests plus actual package selections in both rounds. This avoids
repeating the accumulated corpus merely for each tooling/fixture commit.

The common production qualifier now takes ABI and scenario declarations from the
same manifest, recording builds, exact failure/package comparison and integration
outcomes. It reuses the proven local gameplay runner; that remaining recovery
requirement is explicit rather than silently embedding it in another batch script.

Round 1 completed: 321 C lines removed, eight frozen groups /1,433 records and
51 affected roots in all targets pass. One fixture ownership correction (file
handles) and one compile-only unused-import correction; no Go behavioral
mismatch. Full native qualification took 499.084s. See CLIENT_AUDIO_ASSETS.md.
The small pre-existing scope established the runner; round 2 will test a larger
969-line self-contained scope and a non-root test package.
