# OpenNox C-to-Go port

## Contents

- [Current status](#current-status)
- [Goal and target](#goal-and-target)
- [Batch workflow](#batch-workflow)
- [Subagent use](#subagent-use)
- [Explaining the work and reporting diagnostics](#explaining-the-work-and-reporting-diagnostics)
- [Testing strategy](#testing-strategy)
- [Build and test environment](#build-and-test-environment)
- [Commits and recovery](#commits-and-recovery)
- [Progress and decisions](#progress-and-decisions)

## Current status

Client inventory, quantity-dialog, player-to-player trade, shop UI and quest
journal storage/rendering are converted. The journal batch removed **339 C lines /
11 routines**, with all **717 captured results unchanged**, independent storage,
report and save/load contracts, three-target qualification and fresh gameplay.
There are about **80k C lines left** (**79,924** physical lines in **93** production
files). Next: connected quest briefing presentation and its helpers, a candidate
**13 routines / 947 C lines**. See [the journal report](docs/porting/JOURNAL.md),
[the shop report](docs/porting/CLIENT_SHOP_UI.md) and [PORTING_STATE.md](PORTING_STATE.md).
An isolated hallway mismatch and a later identification-display mismatch remain
unexplained; future failures automatically preserve full captures.

## Goal and target

Replace OpenNox's remaining C implementation with Go while preserving observable
behavior. Qualify the actual Linux x86 client, high-resolution client and server.
The current target is **386/SSE2 with CGO**; support for older CPUs is unnecessary.
Native macOS and browser/WebAssembly work comes after reducing the C dependency
and understanding the remaining platform requirements.

Use headless X for window/input integration and deterministic software rendering.
A local screen is unnecessary. OpenAL's null backend exercises audio initialization;
PCM tests separately check decoding. Physical display and audible playback quality
remain manual release checks. Compare performance against a stable baseline in
the same VM; do not extrapolate its speed to native hardware.

## Batch workflow

1. Select a connected behavior batch, usually several hundred C lines (roughly
   300–1,000 where dependencies permit). Identify callers, callbacks, shared state,
   ownership and observable effects. Move callers with private helpers when useful.
2. Build a recoverable C baseline using real owners and reusable fixtures. Cover
   boundaries, return values, mutations, signedness/overflow, layout, serialization,
   RNG consumption, timing and pixels as relevant. Repeat original-C captures in
   separate processes; normalize only identified nondeterministic fields.
3. Add independent contracts so matching a baseline is not the only correctness
   check. If these uncover an existing bug, make a justified, reversible correction
   before freezing the baseline and record it for later review. Commit the baseline.
4. Translate the batch, keeping C exports only for remaining C callers/callbacks.
   Go callers should invoke Go directly. Compare against frozen expectations while
   implementing; diagnose differences without regenerating goldens to hide them.
5. Run the completed-batch qualification below, review the diff and measure C LOC.
   Update the batch report, decision log where needed, size table and checkpoint.
   Commit and push the conversion before starting another batch.

The user authorized confident, reasonably reversible implementation decisions:
make the decision and record it for later review. Ask when a meaningful product
choice, major compatibility change or costly irreversible action needs their
judgment. Normally continue one chunk at a time; honor any explicit pause requested by the user.

Do not retain C algorithms solely for tests. A committed C baseline and frozen
expectations provide recovery after conversion. Reuse fixtures across related
functions; avoid building a broad framework before there is a demonstrated need.
Batch size is a guide, not a LOC quota or reason to weaken coverage.

## Subagent use

The user supports bounded delegation to a cheaper model such as Terra when the
primary agent is confident it will save total work without weakening the result.
Use the primary agent plus at most one helper by default; avoid an agent fleet.
Follow the active session's delegation rules and available model choices. This
plan does not override restrictions on spawning agents.

A useful split has been a small, well-specified Go implementation draft, with the
primary agent owning the C baseline, tests, review and integration. The recorded
Terra randomized-insertion trial was accepted without corrections, and the
integer/byte/word setter draft also passed review. See
[the insertion report](docs/porting/PROTECTION_INSERT.md) and
[the setter report](docs/porting/PROTECTION_SET.md). These are successful bounded
trials, not evidence that every subsystem is equally easy to delegate or that
subscription savings have been measured.

- Delegate an independent, concrete task: a caller/ABI audit, a bounded helper
  translation against frozen expectations, or a focused review. Give the helper
  exact scope, ownership rules, relevant files, expected behavior and acceptance
  checks. Keep useful independent work for the primary agent while it runs.
- Keep baseline design, ambiguous behavior, shared-state ownership and final
  qualification with the primary agent. Review the draft against C and the
  independent contracts; a helper's report alone is not acceptance evidence.
- Give helpers disjoint files or an ignored draft path. Coordinate all source
  edits with the no-edits-during-builds rule; no concurrent source mutation while
  another agent's tests are reading it. The primary agent integrates and commits.
- Count context transfer, review, corrections and duplicate builds as delegation
  costs. Stop delegating a task if these outweigh the saved work. Larger connected
  batches and reused fixtures remain the main way to reduce qualification overhead.

## Explaining the work and reporting diagnostics

The user has seen repeated UI cybersecurity-classifier interruptions during this
port. The exact triggers are unknown; we have no diagnostic evidence identifying
particular words or logs as the cause. Clear context and bounded output are useful
reporting practices, not a guaranteed remedy.

Explain the concrete game-engine behavior being preserved, the local fixture or
headless scenario exercising it, and the observed result. For example: “Port the
sprite opacity calculation to Go and compare pixels and renderer state against
the committed C baseline.” For legacy protection/checksum or packet-processing
code, name its actual role in the game and the specific compatibility checks.
Keep necessary technical terms, source identifiers and failure details accurate.

Keep full compiler, crash and suite logs in ignored local artifacts. Inspect them
as needed, then report the relevant error, affected function, expected/actual
result and the artifact path. Prefer a short diagnostic excerpt or exact failure
comparison to repeatedly dumping entire logs into the conversation. Preserve the
complete evidence locally and record material failures and fixes in the batch
report. This also makes reviews easier and avoids publishing unrelated log data.

Do not disguise the task, use euphemisms to conceal its purpose, change algorithms
or omit tests to influence a classifier. If an interruption recurs, checkpoint the
actual source and qualification state so work can resume without repeating it.
Do not promise wording that prevents interruptions or weaken port quality to try
to avoid them.

## Testing strategy

Use focused tests during implementation, then broaden at the completed batch
boundary. Do not repeat the entire qualification for each small internal helper.
Reconsider the tests as the behavior and failure modes become clearer.

- **Baseline:** repeated C captures and independent contracts; qualify affected
  targets and relevant integration before replacement. Record exactly which source
  state and cases supplied the oracle.
- **Native batch:** focused captures/contracts plus accumulated tests for affected
  callers, owners and dependencies in default/server/highres. Audit the selection
  against real callers and shared state, and record its pattern and coverage.
  Run the complete accumulated port corpus at subsystem milestones, when shared
  infrastructure changes, or when a failure leaves the affected scope uncertain.
  Check that selected tests actually start and finish; discovery success or a
  process exit alone does not establish coverage.
- **Builds and ABI:** all three production binaries on 386/SSE2/CGO. Check expected
  C callbacks/exports, retired symbols and absence of test helpers. Assert sizes,
  alignment and offsets for types crossing the C/Go boundary.
- **Integration:** fresh asset copies and the relevant headless gameplay scenario,
  with reference comparison enabled. Extend beyond the warrior smoke scenario
  when the changed subsystem needs it; a smoke test cannot cover every branch.
- **Broader checks:** full asset suite at subsystem milestones and shared changes,
  comparing the exact known failure set and package outcomes. Add save/load,
  multiplayer/protocol or meaningful performance checks when the batch warrants
  them. Check tool support before choosing race, sanitizer or checkptr runs on 386.

A completely green legacy suite is not a prerequisite for porting. Known failures
in blob tooling, renderer goldens and audio goldens are recorded and remain visible.
New failures or changes to the established failure set require investigation.
Keep environmental failures separate from source or fixture failures.

Use actual allocation/list/index, rendering, player/team, lighting and asset owners
for stateful fixtures. Reset global state between cases and restore ownership on
cleanup. Avoid reading uninitialized padding or comparing raw pointer addresses.
Case counts document coverage size; discriminating cases and independent contracts
are what make that coverage useful.

Keep source-rewriting tests on temporary copies, and verify broader checks leave
the checkout unchanged. **Do not edit Go/C source while tests or builds are reading
it.** Drafts and documentation can be prepared separately while checks run.

## Build and test environment

Run from the repository root. The current VM has an ignored environment helper;
source it in **every shell that invokes Go**:

```bash
source build/baseline/env.sh
python3 tools/porting/run_tests.py \
  --pattern-file docs/porting/accumulated-test-pattern.txt \
  --tags porttest \
  --log build/port-check.jsonl \
  --result build/port-check-result.json
python3 tools/porting/c_loc.py
```

The helper configures `GOARCH=386 GO386=sse2 CGO_ENABLED=1`, GCC/G++, i386
pkg-config paths and ignored Go caches. Preserve the qualified toolchain when
comparing recorded outputs. Do not silently change dependencies or switch to
amd64 to get a failing check to run.

For a new VM or missing helper, follow [RECOVERY.md](docs/porting/RECOVERY.md).
It contains environment setup, three-target builds, asset extraction, the tracked
[warrior scenario](docs/porting/warrior-smoke.yaml) and independent regeneration
of lost gameplay references. Headless gameplay uses Xvfb at 1280×960, a fresh
save directory, null audio and `NOX_E2E_OVERRIDE=false` during validation.

Original assets and the archive stay unchanged and outside Git. Use separate data
copies for runs. Raw logs, captures, screenshots, binaries and caches live under
ignored `build/`; retain relevant manifests for local reproducibility. Completed
run copies may have verified deduplication/restoration manifests. Do not remove
active run data or the original assets to reclaim space.

## Commits and recovery

The working branch is `dev` in [dbenamy/opennox](https://github.com/dbenamy/opennox).
Commit as Daniel Benamy `<daniel@benamy.info>`. Commit and push qualified batch
checkpoints so a lost session or VM does not lose source changes and expectations.
The existing explicit push command is:

```bash
git -c core.sshCommand='ssh -o BatchMode=yes' \
  push git@github.com:dbenamy/opennox.git dev:dev
```

Review local changes before resuming and preserve unrelated user files. Stage an
explicit file list; avoid destructive Git operations and unrelated reformatting.
Provision credentials separately on a replacement VM. Review logs before sharing;
[recovery guidance](docs/porting/RECOVERY.md#what-needs-a-separate-backup) identifies
outputs that should stay outside Git.

After a lost session, read this plan and the **current checkpoint** at the top of
PORTING_STATE.md, inspect `git status`, then read the relevant batch report.
Historical notes and ignored draft scripts may be stale. Actual source, committed
expectations and recorded qualification take precedence; never rerun a completed
integration script merely because it still exists under `build/`.

## Progress and decisions

Record the new C count after every completed conversion in
[C_LOC.md](docs/porting/C_LOC.md). The counter measures physical lines, including
blanks/comments, in tracked `src/**/*.c`, with test-reference C reported separately.
It excludes headers, C in Go preambles, dependencies and generated build outputs.
It measures source size, not active-code coverage or remaining effort.

The initial count was 142,665 production C lines. See PORTING_STATE.md for the
current rough and exact counts. Use [C_INVENTORY.md](docs/porting/C_INVENTORY.md)
for historical build/linker evidence and [DECISIONS.md](docs/porting/DECISIONS.md)
for corrections and tradeoffs to review. Keep the current checkpoint concise;
put detailed qualification and limitations in the corresponding batch report.
