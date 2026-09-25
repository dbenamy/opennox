# OpenNox C-to-Go port

This file holds durable goals and working rules. Read
[PORTING_STATE.md](PORTING_STATE.md) for the current checkpoint, remaining work,
pause status and local artifact recovery. Keep batch history in `docs/porting/`
and Git, not in this guide.

## Contents

- [Goal and target](#goal-and-target)
- [Batch workflow](#batch-workflow)
- [Subagent use](#subagent-use)
- [Explaining the work and reporting diagnostics](#explaining-the-work-and-reporting-diagnostics)
- [Testing strategy](#testing-strategy)
- [Build and test environment](#build-and-test-environment)
- [Commits and recovery](#commits-and-recovery)
- [Progress and decisions](#progress-and-decisions)

## Goal and target

The immediate goal is to **remove the engine's internal C glue** while preserving
observable behavior. This includes engine-owned C calls, callback dispatchers,
export bridges, C type/header dependencies and libc allocation/string/memory
helpers. Qualify the actual Linux x86 client, high-resolution client and server.

Keep SDL2, OpenGL, OpenAL and similar external native libraries and their bindings
for this phase. Their use of cgo is allowed; `CGO_ENABLED=0` for the complete build
is not an acceptance requirement for this milestone. After the internal glue is
removed, discuss whether and how to eliminate cgo from external bindings. Do not
replace native backends merely to meet the earlier whole-build cgo-free goal.

The target remains **386/SSE2**; support for older CPUs is unnecessary. Keeping
x86/32-bit-specific behavior, layouts and assumptions is acceptable. Change memory
ownership/layout only where needed to remove internal C dependencies correctly;
general layout modernization, 64-bit, native macOS and browser/WebAssembly support
are outside this milestone. Preserve behavior throughout; a successful build alone
is insufficient.

Use headless X for window/input integration and deterministic software rendering.
A local screen is unnecessary. OpenAL's null backend exercises audio initialization;
PCM tests separately check decoding. Physical display and audible playback quality
remain manual release checks. Compare performance against a stable baseline in
the same VM; do not extrapolate its speed to native hardware.

## Batch workflow

The [two-round process trial](docs/porting/PROCESS_TRIAL.md) established the default
process: use focused package checks inside each coherent batch and full
qualification at meaningful boundaries. Recovery commits
may precede full qualification when their evidence and remaining gates are explicit.

1. Select a connected behavior or dependency-removal batch. Size it by ownership
   boundaries and qualification cost, not a standalone C-line quota. Identify
   callers, callbacks, shared state, ownership and observable effects. Move callers with private helpers when useful.
   Check whole-repository reachability before building fixtures: a function with no
   external callers may be a live private helper or completely orphaned. Audit
   callbacks, registrations and C preambles too. A prototype filter must never
   discard a `return function(...)` call. After removing C bodies, also run a
   literal symbol search across every remaining C file; the player-file audit
   caught a live character-creation caller that its declaration filter missed.
   Inspect the enclosing caller
   conditions: a textual reference inside a constant-false branch is not a live
   entrypoint. Follow the reachable private-helper graph from actual roots. Remove proven unreachable code
   with documented evidence instead of translating it solely to keep tests alive.
2. Build a recoverable original-behavior baseline using real owners and reusable
   fixtures (C captures where the original path still uses C). Cover
   boundaries, return values, mutations, signedness/overflow, layout, serialization,
   RNG consumption, timing and pixels as relevant. Repeat original-C captures in
   separate processes; normalize only identified nondeterministic fields.
   Reuse the preceding qualified production baseline when the new baseline changes
   only tests/docs and production source is identical. Record that identity and the
   reused artifacts; always rerun production qualification after the conversion.
   Reuse a just-qualified affected test selection only when all production/test
   source fingerprints and relevant environment settings match exactly. Run new
   owner selections separately, verify discovered-name sets, and require the full
   combined selection after conversion. Record reused versus newly run evidence;
   do not infer coverage from a passing package or a matching test count alone.
   Before freezing, audit numeric constants and lookup tables read by the selected
   C functions. Supply shipped data or explicitly controlled values, and check a
   nontrivial result so zero-filled fixture state cannot hide behavior. The quest
   scoring review found a real width bug that zero-exponent fixtures had masked.
3. Add independent contracts so matching a baseline is not the only correctness
   check. If these uncover an existing bug, make a justified, reversible correction
   before freezing the baseline and record it for later review. Commit the baseline.
4. Translate the batch, keeping C exports only for remaining C callers/callbacks.
   Go callers should invoke Go directly. Compare against frozen expectations while
   implementing; diagnose differences without regenerating goldens to hide them.
   Before launching long milestone gates, review arithmetic widths, signedness,
   pointer construction and callback behavior against C. Matching captured cases
   does not replace that review; add C contracts for newly identified boundaries.
   For floating-point code, inspect the qualified C binary when source types do not
   explain a mismatch. The current C compiler can keep float expressions wide in
   x87 registers; preserve its observable store/reload boundaries. Review related
   calculations together before rebuilding instead of rounding every C float local
   to float32 in Go. The world-geometry conversion demonstrated this distinction.
   For mapped ABI words that can contain integers, write addresses as raw integer
   words. A Go pointer assignment can make its write barrier scan the previous
   integer bits as a managed pointer; the world-grid broad sweep caught this under
   active GC. Keep foreign list links in their original raw representation too.
   Distinguish real pointer fields from integer/pointer unions: seed pointer-only
   fixture fields with live non-null pointers, not arbitrary integer patterns.
   Even native-backed records can pass their old values through Go's write barrier.
   The player-reset full sweep exposed such an invalid fixture seed; keep typed
   production assignments and correct the fixture without changing its goldens.
   Keep raw addresses separate from canonical snapshot identities. Nested fixtures
   must normalize the original value independently, not normalize an identity
   token again: a 32-bit image handle can happen to equal that token. Prove such
   fixture corrections against pre-conversion source with a forced collision;
   preserve frozen captures and rerun qualification on the corrected source.
   Preserve evaluation order around replaceable function hooks such as GetServer.
   Replacing a wrapper with a method call can move argument reads after receiver
   lookup; retain a small native helper or explicit temporary when order matters.
   Before the first compile, format new files, check the whitespace diff, and compare
   new export signatures with every existing header declaration. When removing a
   cgo import, check for `//export` directives too: those still need cgo even when
   no `C.` calls remain. Preserve each retained `//export` as an adjacent function
   doc comment; deleting preceding declarations must not join it to an import or
   closing brace. The remaining-draw overlay review caught this before compilation.
   Imports with `#cgo` directives also carry build settings
   without direct calls; preserve them. When removing engine headers from a
   retained C observer, add its own required standard headers (for example,
   stdint.h for uintptr_t); generated cgo glue is not a substitute. The transfer
   identity compile caught this missing dependency after header cleanup.
   Limit import cleanup to the files changed
   by the batch. Resolve actual declared package names before pruning imports;
   directory basenames can differ (`common/flags` declares `noxflags`). When the
   name is unknown, keep the import for compiler review. A small late source fix
   can invalidate the whole cgo package build and repeat the remaining C compile.
   For GUI dispatch, check the event kind before decoding its arguments. The
   character-creation scenario caught a numeric WindowNewChild ID interpreted as
   a window pointer; button-only fixtures had valid pointers and missed it.
   Include resource-parser notifications in independent event contracts.
   For GUI batches, run a fresh default-client scenario before the full production
   sweep. The browser scenario caught static-label pointer lifetime and a thumb
   child/parent mix-up that focused contracts missed. Keep the original reference
   screens unchanged and repeat final qualification on the corrected source.
   Trace the existing C adapter when choosing a Go API: similar names can hide
   differences in coordinate space, return conventions or ownership.
   Preserve failed lookups separately from valid empty strings. The notice review
   caught a discarded spell-title lookup result: the C formatter prints NULL as
   "(null)", while a valid empty title remains empty. Cover both before milestone
   gates when replacing pointer-returning string adapters.
   For libc parsers, establish saturation, direct float32 rounding, incomplete
   tokens, ASCII keyword matching and NaN payloads before final target sweeps.
   Small ignored library probes can settle these cheaply; retain independent Go
   contracts for the results. Review file adapters' diagnostics and closure at the
   same time. The resource-definition batch found these details late and repeated
   qualification unnecessarily; complete this review before starting long gates.
   When replacing indexed C accessors, extract actual numeric keys and returned
   owners from the C source, including holes and special/null slots. Independently
   compare every mapping before compiling; an owner list alone loses sparse keys.
   The fixture-storage review caught index28 incorrectly compacted to14 in a draft.
   When retiring callback addresses, follow each key into every callback field,
   including test setup. Trace pointer-returning getters and saved aliases into
   raw dispatch calls too; five player-death tests bypassed typed dispatch through
   a helper getter that did not mention the original C symbol. Collision fixtures
   reused Pentagram as a Death callback;
   that cross-family route needed an explicit typed fixture adapter after the key
   became non-executable data. Preserve the original assertions and captures.
   Check dispatch ownership when reusing an existing Go implementation: equal
   output under the default configuration can hide different hooks or queues.
   The team score port caught this through accumulated objective-scoring captures.
5. Run the completed-batch qualification below, review the diff and measure C LOC.
   Update the batch report, decision log where needed, size table and checkpoint.
   Commit and push the conversion before starting another batch.

The user authorized confident, reasonably reversible implementation decisions:
make the decision and record it for later review. Ask when a meaningful product
choice, major compatibility change or costly irreversible action needs their
judgment. Continue one chunk at a time, summarizing each pushed qualification in the
conversation and immediately proceeding to the next. Keep notable issues and
reversible decisions in the batch report and decision log for later review. Stop
only for a substantial blocker or a decision whose answer changes the result in a
way that is hard to undo; honor explicit user pauses.

Do not retain C algorithms solely for tests. A committed C baseline and frozen
expectations provide recovery after conversion. Reuse fixtures across related
functions; avoid building a broad framework before there is a demonstrated need.
Batch size must not weaken coverage.

## Subagent use

The user authorizes proactive delegation to **GPT-6 Luna (`gpt-6-luna`)** for
suitable work throughout the remaining port. Treat this as the default process to
try, and adjust it when evidence shows it is not helping. Use the primary agent
plus **at most one helper** at a time; do not create a parallel agent fleet or
silently substitute another model if Luna is unavailable.

At each batch, identify a substantial, bounded task Luna can own. Delegate it when
specifying and reviewing the result is likely cheaper than doing it locally.
Prefer handing over a complete small task rather than dictating every edit or
having both agents implement the same thing. Keep useful independent work for the
primary while the helper runs. Do tiny edits locally; there is no delegation quota.

Good default assignments:

- **Implementation and caller migration:** translate a bounded helper or move an
  identified set of callers once the behavior, replacement API and acceptance
  tests are established. Use frozen C expectations for integration.
- **Reachability and ABI audits:** enumerate callers, callbacks, registrations,
  preamble references and dependencies; propose removals with inspectable evidence.
  The primary checks the evidence before retiring code or interfaces.
- **Independent test-gap review:** inspect C and a proposed conversion for missing
  boundary cases, ownership issues and signedness/rounding differences. Include
  focused test drafts where useful. Keep final baseline design with the primary.
- **Disk-cleanup audits:** inventory obsolete builds, caches and duplicate assets;
  provide exact paths, sizes, retention reasons and proposed verification steps.
  The primary reviews and executes cleanup after checking active jobs, open files,
  symlink targets and required recovery artifacts. Process/open-file checks must
  see the host process namespace, not only an isolated sandbox view. Do not let the helper delete
  files during an audit. Preserve original assets, current evidence and source;
  use verified deduplication or clearly reproducible obsolete outputs where possible.
- **Documentation and mechanical checks:** draft batch reports, caller inventories,
  LOC counts and qualification summaries from completed artifacts. The primary
  verifies claims against the logs before committing.

Keep ambiguous behavior, architecture/API choices, shared-state ownership,
baseline acceptance, integration and final qualification with the primary agent.
A helper's report alone is not acceptance evidence. The primary reviews the code
against C and independent contracts, runs appropriate checks, and commits/pushes.
Normal reversible choices remain authorized; delegation does not add a new user
approval step.

Give each assignment a compact handoff: objective, exact files and ownership,
relevant baseline/context, required behavior, acceptance checks, prohibited
mutations, and expected deliverables. Use disjoint files or an ignored draft path.
Tell the helper which commands it may run; the primary schedules expensive builds
and tests. Honor the no-source-edits-during-builds rule across both agents. Neither
agent may alter source consumed by an active build/test; a separate uninstalled
draft or read-only audit is suitable overlapping work.

Keep source inventories bounded too. Collect selector names once per file and
look them up in a dictionary; do not scan every source line once per exported
symbol. Use an explicit short timeout (normally 20 seconds) for an inventory
script, and stop/report if it exceeds that bound. Preserve the complete tool
result, including any running session ID. Join or terminate that session before
launching a replacement scan; blank output is not completion. Primary owns host
process checks. A sandbox `ps` cannot establish that host jobs have exited.

Record delegation outcomes briefly in the batch report: task/model, acceptance
checks, meaningful corrections, missed issues, and whether handoff/review/rework
appeared worthwhile. Record actual time or usage only when available; do not infer
subscription savings from a successful test. Reflect at normal batch boundaries
without pausing for user approval. If a task needs repeated steering, substantial
rewrites or duplicate qualification,
finish it locally and narrow future delegation of that task type. Expand the
helper's scope gradually when results support it. Keep this section and the
checkpoint current when changing the process.

Apply these review rules learned from earlier batches:

- Keep numerical/ownership contract design, original-path capture and
  allocation-heavy fixtures with the primary. Review initialization, cleanup and
  restoration of shared state explicitly.
- For generated edits, compare every emitted mapping with the original, including
  sparse keys, array sizes and pointer types. Generator checks do not establish
  valid Go syntax or types; format, compile and qualify the output. When caller
  migration extends beyond the initial selector-edit files, verify casts in those
  callers explicitly too; a matched-reference count does not prove they were
  rewritten.
- Require reachability reports to show whole-source search commands and a concrete
  caller per symbol, distinguishing production, test-only and macro-remapped uses.
  Generate path/line references from search output and verify them. Treat the comment
  immediately before `import "C"` as executable C input, not an ordinary Go
  comment; scan its full contents even across blank lines. A follow-on export
  audit missed test-preamble calls by treating them as comments. Keep broad
  reachability algorithm design with the primary; prefer explicit edit manifests
  for helper coding tasks until those audits demonstrate reliable coverage.
- When moving wrapper bodies into test fixtures, carry over the package imports
  those bodies use and remove imports made obsolete in their old files. Matching
  function bodies alone does not establish that the destination compiles. Prefer
  exact-width Go types where the fixture no longer needs a C interface, and create
  identity keys only for addresses with actual consumers. The server-fixture batch
  needed these corrections before its first build.
  When removing a fixture address table, retain native identities already mixed
  into its getter/registration loop. The player-action draft dropped two live
  normalization keys while retiring C entries; primary restored them before
  compilation. Reserve only removed entries when map size determines capture IDs.
  Preserve side effects before constant returns in every route, including test
  adapters. The drawable-update draft kept the production spark call but omitted
  it in two fixture operations; primary review restored it before compilation.
  Derive the adapter list from actual calls too: a symbol referenced only by a
  prototype or snapshot identity needs no callable replacement. The UI fixture
  draft copied unused wrappers that primary removed before qualification.
- Trace test selection through enclosing functions and fixture operation selectors.
  Match calls when following helper methods: common names such as `state`,
  `result` and `call` also appear as unrelated variables. Audit function-valued
  references separately, and record conservative over-selection explicitly.
  A captured callback slot does not prove a branch ran. Identify media by actual
  container/codec headers rather than filename extensions.
- Test primary review assumptions against the original path too. Do not overrule
  a helper's behavior choice solely on intuition or regenerate frozen expectations
  to accommodate a review mistake.
- Publish helper artifacts atomically, hand off a stable draft and stop editing it
  while the primary integrates. Keep broad inventory algorithm design with the
  primary until bounded execution is demonstrated.

Historical evidence includes the [Luna trial](docs/porting/LUNA_TRIAL.md);
subsequent outcomes belong in their batch reports. Successful bounded drafts do
not establish blanket trust or measured cost savings.

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

For new focused bridge fixtures, prefer the existing pattern of a `porttest` Go
bridge in `legacy` with assertions in the root test package when it exposes the
actual function cleanly. Consumer checks can then share that root test build.
Keep private invariant tests where their access is needed. Measure build time
separately before relocating established tests; compilation can dominate execution
after a cgo source edit.

When retiring C callbacks, preserve names used to assign stable capture IDs:
replace their addresses with nil instead of deleting sorted-table entries. Keep
frozen expectations unchanged; do not mistake an ID shift for a behavior change.

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
  Do not schedule the full corpus merely because another helper, recovery commit
  or arbitrary number of batches has completed.
  Select affected packages explicitly (`--package`, root by default).
  Check that selected tests actually start and finish; discovery success or a
  process exit alone does not establish coverage. The test driver defaults to
  `GOMEMLIMIT=768MiB` to limit Go heap growth within the 386 address space; explicit
  environment settings override it and the effective settings are recorded.
  Independent target sweeps may run concurrently when their output directories
  and fixtures are isolated and total CPU/memory fit this VM.
  The optional `run_profiles.py --jobs 2` controller builds root test binaries
  sequentially, verifies their source/binary/profile records, then runs at most
  two test processes. Use it only for the audited root porttest corpus; keep
  capture/diagnostic output variables unset and dependencies/source frozen.
  Production builds and headless scenarios stay sequential. See the bounded
  [prebuilt-profile trial](docs/porting/PREBUILT_PROFILES.md); fall back to
  `--jobs 1` if concurrent resource or isolation problems appear.
  The accumulated selector is `^Test`: run every root test compiled with `porttest`
  instead of maintaining a historical list of names. Keep focused owner patterns
  separate. Compare discovered and completed root-name sets for each profile;
  document any intentional exclusions or diagnostic skips explicitly. The selector
  repair found 899 existing roots omitted by the former hand-maintained pattern.
  Safe-only roots require their own applicable safe selection.
- **Static memory accesses:** run `go test ./common/memmap/nox -run
  '^TestCodeStatic$' -count=1` from `src` when adding or changing fixtures that
  touch mapped state. This inexpensive preflight also scans porttest files.
  Own raw backing-blob snapshots explicitly, separately from extracted live globals.
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
are recorded in the expectation file linked from PORTING_STATE.md and remain visible.
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
python3 tools/porting/run_batch.py docs/porting/map-encode-batch.json \
  --phase milestone --out build/port-check
python3 tools/porting/c_loc.py
```

Use a fresh output directory. The batch manifest supplies the asset environment
and explicit package selections for all three targets.

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

Do not retain every historical log or intermediate run indefinitely. Keep committed
expectations/reports, current qualification artifacts and useful original-behavior
references. Old raw captures can be compressed with byte-for-byte verification;
superseded text logs and reproducible outputs can be discarded after checking
references and active use. For recovery, match recorded source fingerprints to
the final commit: a build's Git HEAD may name the earlier baseline while its
working tree already contains the conversion. Record removals and any restoration
steps in the checkpoint. Retained capture bytes help diagnose hash mismatches even when tests
use committed hashes rather than reading the historical capture files.

The selected-test runner gives discovery/build a separate `1536MiB` GOMEMLIMIT,
configurable with `--build-memory-limit`, while preserving the execution budget
and explicit runtime overrides above. Discovery compiles and lists tests; the
same-source execution reuses those build outputs. Both environments are recorded
in the result JSON. This is a reversible attempt to reduce compiler GC pressure,
not a measured speed improvement. Baseline runs must still verify runtime settings
and complete execution.

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

See [PORTING_STATE.md](PORTING_STATE.md) for current counts and remaining work.
Use [C_INVENTORY.md](docs/porting/C_INVENTORY.md) for historical build/linker evidence and [DECISIONS.md](docs/porting/DECISIONS.md)
for corrections and tradeoffs to review. Keep the current checkpoint concise;
put detailed qualification and limitations in the corresponding batch report.
