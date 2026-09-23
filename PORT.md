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

Five raw damage callers now use exact integer Go dispatch for eleven canonical
callbacks, preserving boolean APIs, raw fallback and full-word/low-byte results.
All132 roots/profile,2071 frozen cases and independent raw/compatibility contracts
pass, plus safe/static, production/ABI, exact known-suite and headless save/load.
See [DAMAGE_VALUES.md](docs/porting/DAMAGE_VALUES.md).
**Standalone C remains zero; production C preamble bodies remain79.**
The next53-callback update baseline is qualified:281 roots/profile,8251 frozen
cases and dynamic-handler contracts. Production conversion is pending. See
[UPDATE_REGISTRY.md](docs/porting/UPDATE_REGISTRY.md).

### Earlier checkpoints

Go MP3 synthesis now avoids overwritten mono outputs and accumulates each lane
locally while preserving its exact float32 operation order. Controlled medians
improved 34% on one shipped mono asset and 8% on a synthetic stereo stream;
both remain allocation-free. Frozen PCM/state, all 1,246 asset observations,
production/ABI, known-suite and headless creation/save-load checks pass.
**Standalone C remains zero; production C preamble bodies remain 79.**
Next: review a bounded group of redundant Go→C→Go book callbacks before broader
callback architecture work. See
[MP3_SYNTHESIS_PERFORMANCE.md](docs/porting/MP3_SYNTHESIS_PERFORMANCE.md).

Two typed callback shims now reuse existing shared dispatchers, preserving both
Go APIs. All 24 forwarding/return cases pass per profile; client consumer tests,
fresh production/ABI, known-suite and headless creation/save-load checks pass.
**Standalone C stays zero; production C preamble bodies fall 81→79.** The remaining
bodies are callback invocation glue. Next: profile the documented MP3 performance
gap before considering a broader callback representation redesign. See
[TYPED_CALLBACK_ADAPTERS.md](docs/porting/TYPED_CALLBACK_ADAPTERS.md).

Production audio now uses the fully qualified Go MP3 decoder. All 1,246 shipped
asset observations and historical PCM goldens match unchanged; all decoder
profiles, production/safe builds, ABI checks, exact known-suite comparison and
headless creation/save-load pass. **Standalone production C: zero files/lines.**
The 1,890-line decoder header is retired too; 81 production C preamble bodies
and external cgo dependencies remain. Go decode is allocation-free per frame,
but measured 2.79× slower than C on one mono asset; retain this performance
review item. Next: deduplicate two typed callback adapters through existing
shared dispatchers while preserving their Go APIs. See
[MP3_GO.md](docs/porting/MP3_GO.md).

The complete Go MP3 decoder matches2,682 generated C frame calls in544 sequences,
including PCM, metadata and observable state. All ten test roots pass in four
profiles and with cgo disabled; vet passes. Original unused-QMF and private-bit
initialization findings are documented for review. Next: wire audio and qualify
unchanged shipped-asset goldens and production builds. C remains six standalone
lines plus the active header. See [MP3_FRAME.md](docs/porting/MP3_FRAME.md).

Go MP3 Huffman/dequantization matches34,353 exact C cases, including every
codebook leaf and directed escape extremes/signs. All nine helper roots pass
in four profiles and with cgo disabled; vet passes. Next: complete-frame decoder
assembly and PCM/state qualification before switching production audio. C remains
six standalone lines plus the active header. See
[MP3_HUFFMAN.md](docs/porting/MP3_HUFFMAN.md).

Go MP3 PCM synthesis/filterbank state matches197,774 C records, including
196,616 rounding cases and576 stateful granule steps. All eight helper roots
pass in four profiles and with cgo disabled; vet passes. C remains six standalone
lines plus the active header. Next: Huffman/dequantization and decoder wrappers.
Production remains on C. See [MP3_SYNTHESIS.md](docs/porting/MP3_SYNTHESIS.md).

Go MP3 synthesis DCT matches2,172 exact original-C cases in four profiles and
with cgo disabled; all seven helper roots and vet pass. C remains six standalone
lines plus the active header. Next: PCM synthesis and persistent filterbank state.
Production evidence is reused for this unimported package. See
[MP3_DCT.md](docs/porting/MP3_DCT.md).

Go MP3 inverse transforms and overlap state match2,380 exact original-C records,
including1,258 granule steps. All six helper roots pass in four profiles and with
cgo disabled; vet passes. Production evidence is explicitly reused for the
unimported package. C remains six standalone lines plus the active header.
Next: synthesis transforms. See [MP3_IMDCT.md](docs/porting/MP3_IMDCT.md).

Go MP3 stereo, spectral reordering and antialiasing match4,364 exact C state
records in four profiles and with cgo disabled; all previous helper roots pass.
A documented legacy8kHz mixed-block extent is preserved in bounded workspace
slices for later integration review. The decoder remains unwired; existing
production evidence is explicitly reused. C remains six standalone lines plus
the active header. Next: inverse transforms/overlap state. See
[MP3_SPECTRUM.md](docs/porting/MP3_SPECTRUM.md).

Go MP3 scalefactor byte parsing and quarter-power scaling now match57,002 frozen
scalar-SSE2 C cases, with exact float bits and every reuse mask. All four helper
roots pass in four profiles and with cgo disabled. The complete decoder remains
unwired; production evidence is explicitly reused after source/binary checks.
C remains six standalone lines plus the active header. Next: stereo/reordering/
antialiasing. See [MP3_SCALEFACTORS.md](docs/porting/MP3_SCALEFACTORS.md).

Go MP3 frame scanning, initialization and reservoir helpers now match2,774 frozen
C cases, including the real48KiB input-buffer boundary and retained state/tails.
All previous MP3 helper tests still pass in four profiles and with cgo disabled.
The decoder remains unwired; existing production evidence is explicitly reused.
C remains six standalone lines plus the active decoder header. Next: scalefactors
and exact floating-point results. See [MP3_STREAM.md](docs/porting/MP3_STREAM.md).

Go Layer III side-information parsing now matches 36,864 frozen C state records,
including partial writes on errors. Table-row alias and CRC-offset regressions
pass alongside the integer helpers in four profiles and with cgo disabled.
The Go decoder remains unwired; production evidence is reused with source and
binary identity checks. C remains six standalone lines plus the active header.
Next: frame scanning and reservoir state. See [MP3_SIDEINFO.md](docs/porting/MP3_SIDEINFO.md).

The first Go MP3 internals (bit reader/header arithmetic) match 819,207 frozen
actual-C cases in four profiles and with cgo disabled. They remain unimported by
production while the complete decoder is assembled; existing production evidence
is explicitly reused after source/dependency/binary checks. C stays six lines,
with the third-party header still active. Next: Layer III side-information parsing.
See [MP3_INTEGER.md](docs/porting/MP3_INTEGER.md).

The audio package now uses SSE2 scalar arithmetic on 386. All 1,246 historical
MP3 dialog PCM goldens pass unchanged; the full suite loses exactly 1,249 audio
failure events and retains the 304 unrelated events. Client audio/stream tests,
safe build/static, fresh production/ABI and headless creation/save-load pass.
The old and corrected guarded decoder captures are preserved. Standalone C stays
six lines; next is coherent Go decoder internals with additional format/state
contracts. See [MP3_DECODER.md](docs/porting/MP3_DECODER.md).

Five callback-address getters now reference the same C functions directly from
Go. All 38 consumer roots pass per normal profile, with safe build/static, fresh
production/ABI, exact known-suite and headless creation/save-load qualification.
Standalone C stays **six lines/one file**; production C preamble bodies fall from
86 to 81. See [ADDRESS_ADAPTERS.md](docs/porting/ADDRESS_ADAPTERS.md).

The ten live empty callback identities now export from Go. All 338 frozen cases,
69 consumer roots per normal profile and 24 under safe pass, alongside fresh
production/ABI, exact known-suite comparison and headless creation/save/load.
Standalone C is **six physical lines in one file**, down 19; that file includes
the third-party MP3 decoder. Preambles, headers and generated bridges remain.
The measured extra callback cost is 84–140 ns/call in this VM; accepted as a
reversible compatibility cost, with game-frame impact still unmeasured. See
[EMPTY_CALLBACKS.md](docs/porting/EMPTY_CALLBACKS.md).

The six optional safe-profile forwarding shims now export directly from Go.
All 142 frozen cases, the shop-loading consumer, safe build, production/ABI,
known-suite and headless/save-load gates pass. Standalone C is **25 physical lines
in three files**, down 20; C headers, generated bridges and third-party decoder
implementation remain. See [SAFE_BRIDGES.md](docs/porting/SAFE_BRIDGES.md).

Entry-character predicates now call the same libc classifiers directly, removing
two custom cgo-preamble bodies. All 131,072 classifications and twelve widget roots
per profile match the C baseline; production/ABI/gameplay/save-load gates pass.
Standalone C stays 45 lines/four files because preamble bodies are outside that
metric. See [ENTRY_CLASSIFIERS.md](docs/porting/ENTRY_CLASSIFIERS.md).

The two entry-character predicates now have a qualified actual-C baseline covering
all 131,072 boolean classifications plus twelve widget test roots per profile.
Full native/ABI/gameplay/save-load gates pass; that baseline supports the
subsequent direct-libc conversion above. C remains 45 lines/four files. See [ENTRY_CLASSIFIERS.md](docs/porting/ENTRY_CLASSIFIERS.md).

Unused header/preamble helpers and unregistered empty callbacks are retired,
along with two empty Obelisk calls. Synchronization and all ten live callback
identities are preserved. Focused tests in all three profiles, storage/GC
contracts, production/ABI and gameplay/save-load qualification pass. Standalone
production C is now **45 lines in four files** (−6); headers, preambles and
third-party C remain outside that count. See [ORPHAN_INLINE.md](docs/porting/ORPHAN_INLINE.md).

The remaining 44 globals and eight mapped buffers now initialize from Go, using
the existing foreign allocator for stable process-lifetime storage. Both frozen
storage captures match C; all 2,291/2,280/2,291 consumer roots pass across the three
profiles, alongside fresh production/ABI/gameplay/save-load qualification. This
removes 63 C lines and two translation units. C is now **51 physical lines in four
files**. C types, callbacks, inline preambles and libc/CGO dependencies remain
explicit follow-up work. See [RAW_STORAGE.md](docs/porting/RAW_STORAGE.md).

An audio stream GC regression is corrected: opaque sample addresses remain raw
32-bit words through buffer/chunk/voice bookkeeping and convert to pointers at
sample access. The regression failed before the fix and passes after it. The full
2,291/2,280/2,291 consumer sweep and fresh production/gameplay/save-load qualification
pass. The remaining 44 globals and 8 shared buffers now have a qualified actual-C
storage baseline, used by the subsequent conversion above. That repair left C at
114 lines. See
[AUDIO_ADDRESS_GC.md](docs/porting/AUDIO_ADDRESS_GC.md) and
[RAW_STORAGE.md](docs/porting/RAW_STORAGE.md).

Another 52 numeric globals (map generation, audio and sustained spells) now have
Go owners; 11 C fixture accessor/table bodies are retired. All 54,905 storage
cases match the frozen C capture. Each profile passes 183 focused consumer roots
without skips. Safe build, production/ABI, exact known-suite comparison and fresh
headless gameplay/save-load pass. This removes 52 C lines. See
[FIXTURE_STORAGE.md](docs/porting/FIXTURE_STORAGE.md).

343 numeric globals now have process-lifetime Go owners; nine unused C definitions
are removed. The 47,677-case storage capture matches C in all three profiles.
The full tagged consumer sweep passes 2,291/2,280/2,291 roots in default/server/highres,
with no skips. Safe build, production/ABI, known-suite comparison and fresh headless
gameplay/save-load pass. This removes 352 C lines. See
[SCALAR_STORAGE.md](docs/porting/SCALAR_STORAGE.md).

Unused C memory accessors and five GUI adapters are retired. The existing Go
registry now supplies the durability fixture's threshold address. Four focused
tests and static checks pass in default/server/highres; the optional safe build,
production/ABI, known-suite comparison and headless gameplay/save-load also pass.
This removes 117 C lines and 24 interfaces. See [ORPHAN_BRIDGES.md](docs/porting/ORPHAN_BRIDGES.md).

Custom text formatting, scalar strings and the audio catalog directory check are
Go. The 375/373/375 affected roots pass in default/server/highres, along with 668,876 frozen
cases and native consumer contracts. Production/ABI, known-suite comparison,
headless gameplay and save/load pass. This removes 584 C lines and 28 interfaces.
See [TEXT_FORMAT.md](docs/porting/TEXT_FORMAT.md).

Extension player lookup, weapon cycling, trap drop, flag index and server listing
are Go. All 173 affected tests pass in default/server/highres, and 21,135 cases
match the corrected C baseline. Production/ABI, known-suite comparison, headless
gameplay and save/load pass. This removes 541 C lines from the corrected baseline
(521 net since the preceding conversion). See [SERVER_TEXT.md](docs/porting/SERVER_TEXT.md).

The remaining runtime numeric and lookup helpers, pause lifecycle and saved-creature
ownership are Go. This removes 965 C lines and 34 old interfaces. Default/server/highres pass 304/303/304 affected
tests; 55,986 captured cases match C. Production/ABI, the
known full-suite comparison and headless gameplay/save-load pass. See
[SERVER_RUNTIME.md](docs/porting/SERVER_RUNTIME.md).

Client resource teardown, timing, colors, caches and menu/modal lifecycle are Go,
removing 838 C lines and 18 old interfaces. Default/server/highres pass 135/135/135
affected roots and all 6,415 C cases; fresh production/ABI, known-suite comparison
and headless gameplay/save-load pass. See [CLIENT_RESOURCES.md](docs/porting/CLIENT_RESOURCES.md).

The remaining client render helpers and audio lifecycle are Go, removing 915 C lines
and 15 unused exports. Default/server/highres pass 96/91/96 affected roots, no skips,
and 6,129 captured cases match C. Fresh production/ABI, known-suite comparison and
headless gameplay/save-load pass. See [CLIENT_RENDER_HELPERS.md](docs/porting/CLIENT_RENDER_HELPERS.md).

The ordered client message queue is now Go, removing another 145 C lines and
four unused interfaces. Its 32 affected tests pass across all three profiles;
four captures /1,699 cases match the C baseline, and fresh production/ABI,
known-suite comparison and headless gameplay/save-load pass. See
[CLIENT_SEQUENCE.md](docs/porting/CLIENT_SEQUENCE.md).

The client settings, team, trade and quest dispatcher and its private ball HUD
helper are now Go and fully qualified. The conversion removes 1,536 C lines and 146
unused C interfaces. Default/server/highres pass 632/628/632 test roots, no skips,
and all 134 frozen captures (202,586 cases). Fresh production/ABI, exact full-suite
comparison, headless gameplay and save/load checks pass.

C remaining is **six physical lines in one file**, zero reference C; the active
MP3 implementation header and C preambles are counted separately. See
[C_LOC.md](docs/porting/C_LOC.md).
See [PORTING_STATE.md](PORTING_STATE.md) for recovery details.

The preceding world-grid conversion is recorded in
[WORLD_GRID.md](docs/porting/WORLD_GRID.md).

Previous completed GUI batches include
[client interaction](docs/porting/CLIENT_INTERACTION.md), the
[server browser](docs/porting/SERVER_BROWSER.md) and
[session dialogs](docs/porting/SESSION_DIALOGS.md).

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

The [two-round process trial](docs/porting/PROCESS_TRIAL.md) established the default
process: use focused package checks inside each coherent batch and full
qualification at meaningful boundaries. Recovery commits
may precede full qualification when their evidence and remaining gates are explicit.

1. Select a connected behavior batch, aiming for roughly
   1,000–3,000 C lines where dependencies permit. Identify callers, callbacks, shared state,
   ownership and observable effects. Move callers with private helpers when useful.
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
2. Build a recoverable C baseline using real owners and reusable fixtures. Cover
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
   Before the first compile, format new files, check the whitespace diff, and compare
   new export signatures with every existing header declaration. When removing a
   cgo import, check for `//export` directives too: those still need cgo even when
   no `C.` calls remain. Imports with `#cgo` directives also carry build settings
   without direct calls; preserve them. Limit import cleanup to the files changed
   by the batch. A small late source fix can invalidate the whole cgo
   package build and repeat the remaining C compile.
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
Batch size is a guide, not a LOC quota or reason to weaken coverage.

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

The balance-getter audit exposed this failure mode: multiple repeated Python
scans ran for minutes and competed with qualification. Primary verified the
parent scripts in the host namespace, stopped the remaining scan and confirmed
all had exited. A replacement single-pass inventory took about 0.2 seconds. This
was an audit-efficiency failure, not a source/test mismatch; no build timing
comparison from the contended run should be treated as a performance benchmark.
Keep broad inventory algorithm design with the primary until bounded execution
is demonstrated; Luna remains useful for exact-list drafts, fixture drafts,
coverage review and storage inventories with primary verification.

Record delegation outcomes briefly in the batch report: task/model, acceptance
checks, meaningful corrections, missed issues, and whether handoff/review/rework
appeared worthwhile. Record actual time or usage only when available; do not infer
subscription savings from a successful test. Reflect after the next two completed
batches, then at normal batch boundaries without pausing for user approval. If a
task needs repeated steering, substantial rewrites or duplicate qualification,
finish it locally and narrow future delegation of that task type. Expand the
helper's scope gradually when results support it. Keep this section and the
checkpoint current when changing the process.

Evidence so far: the [Luna scalar trial](docs/porting/LUNA_TRIAL.md) passed all
655,391 captured cases without behavior corrections; the primary requested one
readability cleanup. Earlier Terra trials also succeeded for
[randomized insertion](docs/porting/PROTECTION_INSERT.md) and
[integer/byte/word setters](docs/porting/PROTECTION_SET.md). These support bounded
delegation, not blanket trust in every subsystem or measured cost savings.

After the first two batches under these guidelines, retain one bounded Luna helper.
The scalar implementation and orphan-removal draft passed qualification after
primary review. Mechanical edits and exact test selection were useful. Import
cleanup, reachability and disk audits required corrections; require executable
preflight checks and machine-generated evidence for those tasks. Continue to own
baseline acceptance, storage/ABI choices and final qualification in the primary.
The scalar storage audit caught pointer-typed redeclarations worth separating
from ordinary numeric owners. No subscription-cost reduction has been measured.

Recent storage work reinforces that boundary. Luna's numeric-owner draft lost a
sparse fixture index; independent comparison caught it before compilation. The
remaining-storage draft needed fixes to blob-size rewrite ordering, array brackets
and a pointer getter. Its generator checks did not prove Go syntax or types.
Compare emitted output with the original definitions, then compile and qualify it.
For reachability audits, require a concrete caller per symbol and distinguish
production, test-only and macro-remapped calls; the safe-adapter audit initially
generalized one live caller to all six wrappers. Cleanup inventories have been
useful when each artifact includes its hash and retained qualification evidence;
the primary still verifies host process references and performs deletion.
Test-selection reports must trace the enclosing function and fixture dispatch
chain, including numeric operation selectors. Searching only a guessed helper
name missed existing obelisk coverage in a later cleanup draft.
A bounded read-only lifetime review also identified stale per-case address IDs in
the shared test fixture. The primary confirmed that defect with a deterministic
red; broader captures then caught the need to retain aliases for persistent
objects. Such reviews are useful hypothesis generation, while full qualification
still decides whether the repair preserves the established behavior.
Require reachability reports to include the exact whole-source search command
and scope. The entry-classifier audit missed two production constructor callers
outside `src/legacy`; the primary found them with a literal search across `src`.
Do not accept “no in-repository callers” from a directory-local search.
Generate path/line inventories directly from search output instead of manually
transcribing them; verify each reported path and line before accepting the report.
Recent callback/decoder audits add two checks: identify media by container/codec
headers rather than filename extensions, and distinguish a captured callback slot
from proof that a particular branch executed. Existing bot-update cases do not
establish allocation-failure coverage just because their snapshots include the
fallback assignment's destination. Primary review caught both overstatements.

Recent decoder work suggests keeping numerical/ownership contract design and C
capture with the primary, while Luna drafts a bounded connected implementation
and independently reviews capture inputs. Primary-owned fixture writing can overlap
that work. Review representation and caller bounds explicitly: a table family is
not a retained table-row alias, and a capture's size limit is not automatically a
production limit. Recheck even small cleanup edits for remaining symbol uses.
This is an observed workflow adjustment, not measured model cost/speed savings.

Damage-registry work reinforces the need to test primary review assumptions too:
Luna's initial late callback selection matched the original compiled behavior;
primary's requested early selection did not. An original-path mutation test exposed
that before production changes or frozen expectations. Keep the original failure
and record who corrected what. Publish helper audit artifacts atomically so a
partially written candidate list cannot be mistaken for a finished result.

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
separately before relocating established tests: the entry-classifier batch spent
about 87 seconds on a legacy gate whose tests took 1.9 seconds, and 149 seconds on
a root gate whose tests took about four seconds after a cgo source edit.

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

The selected-test runner gives discovery/build a separate `1536MiB` GOMEMLIMIT,
configurable with `--build-memory-limit`, while preserving the execution budget
and explicit runtime overrides above. Discovery compiles and lists tests; the
same-source execution reuses those build outputs. Both environments are recorded
in the result JSON. This is a reversible attempt to reduce compiler GC pressure,
not a measured speed improvement. Nine Python accounting/environment checks pass;
actual baseline runs must still verify runtime settings and complete execution.

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
