# Porting checkpoint — 2026-09-10

Read CODEX_HANDOFF.md for the working plan. This is the resume checkpoint.

## GitHub backup and recovery

On 2026-09-10, SSH authentication to GitHub succeeded as dbenamytravis and all
seven commits through e694e1ac were pushed to dbenamy/opennox, branch dev. The
remote branch was verified at e694e1ac40ab12da869cccb50248881fd537b3f4 before this
checkpoint update. The plan and this log were already included in that push.

See [recovery instructions](docs/porting/RECOVERY.md) and the tracked
[warrior smoke scenario](docs/porting/warrior-smoke.yaml). They preserve the
asset-free setup needed to resume without the ignored build/ directory. Original
assets still require a separate user backup; raw logs, binaries and screenshots
are not on GitHub. Provision SSH credentials separately; no private key or token
is recorded here. The SSH push URL is git@github.com:dbenamy/opennox.git; origin
still uses HTTPS, so an explicit SSH URL was used for the push.

## Current status after infrastructure repairs

Plan and infrastructure changes are committed. Screenshot checks now fail
reliably, automatic test writes are isolated, and stale API/vet failures are
repaired. Full default suite with assets is down from seven failing packages to
three: blobs tooling, renderer goldens and audio goldens. None are suppressed.
No C-to-Go conversion has begun. Historical baseline results below remain useful;
see the final infrastructure section for the latest checks.

## Repository and environment

- Original baseline: b184030e76be2b681a7f6d2bcdef52b091d94b9b on dev; origin is the user's fork,
  https://github.com/dbenamy/opennox.git.
- No pre-existing tracked changes. Only the handoff and media archive were
  untracked initially. No C-to-Go porting has begun; subsequent infrastructure edits are committed.
- Ubuntu 26.04.1 x86_64, Go 1.26.0, multilib GCC, i386 SDL2/OpenAL dev packages.
  Current Codex sandbox is disabled at the user's request. Stay within this host.
- Added archive/headless tools; bsdtar, 7z, unsquashfs, Xvfb and xdotool available.
  dpkg --audit is clean. No reboot performed. Old UTM hardware details are historical.
- Initial sandbox SIGSYS on 386 execution, network denial and unwritable module
  cache are resolved. Explicit CC=gcc/CXX=g++ avoids Go selecting absent
  i686-linux-gnu-gcc. Target remains 386 with CGO.

## Local artifacts and reproducibility

Everything under build/ is ignored by Git. These artifacts persist in this VM,
not in commits. See build/baseline/README.md for details and commands.

- build/baseline/env.sh: target flags, i386 pkg-config and build/cache Go caches.
  Pinned modules downloaded successfully; no go.mod/go.sum changes.
- build/baseline/bin/{opennox,opennox-hd,opennox-server}: pristine baseline builds.
  binaries.json records SHA-256, ELF32/Intel 80386 headers and help exit codes (0).
  environment.json records toolchain, revision, flags and source archive hash.
- logs/build-pristine.log: all targets built successfully in 7.66 seconds warm,
  peak RSS 455332 KiB. No cold-build or gameplay performance claim.
- The user-confirmed asset archive actually contains Nox.wsquashfs, with an
  installed drive_c/Nox tree. Extracted only that subtree to
  build/assets/extracted/drive_c/Nox. Original archive untouched. Test runs use
  separate copies and empty save directories; no bundled executables were run.
- Full JSON test outputs, exit files and build graphs are in build/baseline/logs.
  test-summary.json gives package/test event counts.

## Test outcomes: baseline is red

All three full suite variants were attempted with assets, GOARCH=386,
CGO_ENABLED=1 and allowed C flags. Each exits 1:

| Variant | Passed packages | Failed packages | Packages without tests |
| --- | ---: | ---: | ---: |
| default | 11 | 7 | 30 |
| server | 10 | 7 | 30 |
| highres | 11 | 7 | 30 |

One explicit skipped test with assets: internal/blobs/TestSplitBlob. Failed
package counts include compilation/vet failures. Test/subtest counts are not
independent scenarios.

Failures:
- cmd/noxmovie uses old sdl.New signature (missing logger).
- internal/netstr tests have obsolete callback types and missing netmsg import.
- internal/offalign tests pass uint where uintptr is required.
- Root package vet rejects log.Printf's %w in maps.go.
- internal/blobs has obsolete memmap.go path and a formatter failure.
- noxrender image/particle goldens fail. Particle tests fail even without assets.
  They hash encoded PNG, not just pixels; root cause is not yet established.
- Audio PCM hashes mismatch supplied data. Asset/decoder/toolchain differences
  have not been isolated; do not simply regenerate the expected hashes.
- server-tag legacy/dialog expects "empty" and receives an empty string.

Important: some existing tests MODIFY SOURCE. The first full run changed ten
tracked files and emitted token files. Saved its diff to
logs/test-generated-changes.patch, restored only those generated edits, moved
emitted files to logs, and rebuilt all binaries from clean source afterward.
Subsequent suites ran in build/baseline/test-checkout. That disposable copy was
reused across variants; future runs must recreate it per variant to avoid any
cross-run contamination. The working checkout is clean apart from planning docs
and the original archive.

## Gameplay scenario and oracle

- Existing src/e2e*.go harness supports YAML input, simulated time, platform RNG,
  screen comparisons and save hashing. Use it before inventing another harness.
- build/baseline/e2e/warrior-smoke.yaml: title → Solo → new warrior → war01a →
  captain dialogue → short walk → screenshots → clean quit.
- run-scenario.py RUN_NAME [BINARY] creates a fresh data/save copy, uses Xvfb
  1280x960 and OpenAL null backend (audio enabled), enforces 120s timeout, and
  records command/exit/elapsed/logs in runs/RUN_NAME. Run names must be new.
- runs/repeat-a and runs/repeat-b both exited 0 in about 40 seconds. Visually
  inspected gameplay and moved-player frames. Decoded NRGBA pixels match exactly
  at both checkpoints (frame-comparison.json). This covers this short scenario,
  not broad gameplay correctness or physical audio/display quality.
- compare-frames.go and its binary perform independent decoded-pixel comparison.
  It correctly rejects a deliberately replaced frame (comparator-negative/).
- Why independent at baseline? The original Screen passed nil to e2eError on a
  mismatch and auto-created missing goldens. This was repaired in 3575f443; current
  checks require explicit golden updates and reject mismatches.
- HD client also completed the same fresh scenario with exit 0 in 39.3 seconds.
  Both captured frames match the standard-client pixels at this 1024x768 game
  resolution (hd-comparison.json); higher resolutions are not covered.
- Save comparison: war01a.map bytes match across runs; Player.plr differs in both
  WORKING and AUTOSAVE. See save-comparison.json. Cause and save-load compatibility
  are unverified. Missing-script-object warnings also remain in baseline logs.

## C inventory started

Saved go list -deps -json graphs for default/server/highres. Local compiled C:
152 for clients (148 legacy + 3 cnxz + 1 ail), 151 for server (no ail C file).
This is compiled translation-unit inventory, not linked/reachable/active logic.
Indirect callbacks, retained symbols and subsystem classification remain pending.

## Next work

1. Diagnose PNG/PCM goldens without masking real behavior changes. Investigate
   Player.plr nondeterminism and add a save-load scenario.
2. Use build graphs for a bounded dependency audit; choose a cohesive C leaf and
   establish differential coverage before porting it.

Do not call the suite green. Dedicated-server map/tick scenarios, multiplayer,
replay validation, sanitizer compatibility and performance work remain pending.

## Infrastructure changes after baseline

- Plan/checkpoint committed as af08739d; repository-local Git author configured
  from the user's supplied identity.
- Screenshot oracle: moved comparison into internal/e2etest. Normal checks fail
  for missing goldens, pixel differences and dimension differences; only explicit
  NOX_E2E_OVERRIDE updates goldens. Decoder formats/origins are normalized, actual
  and diff frames are retained, input buffers are not mutated, and file errors
  propagate. Focused 386 tests pass; all three targets build into build/infra-bin.
  A real headless client with a deliberately wrong-size golden exits 2 with a
  screen mismatch (build/screen-negative), confirming integration fails visibly.
- Source-rewriting blobs tests now copy src into t.TempDir and restore the global
  tool path afterward. Token diagnostics also go into t.TempDir. Package execution
  left tracked engine files untouched. Token tests pass; blobs retains the known
  obsolete memmap.go-path and formatter parse failures (tests-isolated.log).
  Those failures are not skipped or relabeled as passing.
- Full-suite follow-up identified internal/noxfactor/TestNoxFactor as another
  source rewriter. Its generated enum substitutions were saved to
  logs/noxfactor-generated.patch and restored; the test now uses a temporary
  source copy. Default rendering/heatmap PNGs also use temporary output paths.
  Focused noxfactor, memmap and primitive-render tests pass on 386 and leave
  engine source untouched (tests-isolation-followup.log).
- Updated obsolete movie-player SDL call, offalign fixture integer types, and
  maps.go's invalid logging format. Reworked the netstr test for current typed
  callbacks, synchronous server binding, loopback payload verification, a bounded
  wait and goroutine cleanup. Target tests pass, including five repeated netstr
  runs; supplementary amd64 netstr race test passes. Root package passes default
  vet/compilation. No production network behavior changed.
- Full default suite with assets now has 14 passing packages, 3 failing packages,
  32 packages without tests, and no compilation/vet failures. Remaining failures
  are internal/blobs, client/noxrender, and legacy/client/audio/ail, already known
  from baseline. See logs/tests-infra-default.*. Final test-isolation follow-up
  above was validated separately after this full run.
- Tagged follow-up: the movie command lacked !server even though its movie library
  is client-only. Added that matching constraint; it is omitted by server go list
  ./... and still compiles for default/highres. Focused server root/netstr/offalign/
  e2etest checks pass; highres focused checks pass. Initial variant logs preserve
  this discovered setup failure (their shell's final exit reflected highres only).
- Final production validation: all three targets build after infrastructure/API
  changes (logs/build-infra-final.log). The patched standard client completed
  the fresh warrior scenario against both pre-existing baseline goldens with
  NOX_E2E_OVERRIDE=false and exit 0 in 38.9s (runs/infra-positive). This exercises
  the repaired in-engine screenshot oracle's success path; screen-negative
  exercises its failure path. Only the original media archive remains untracked.

Next session should start with the bounded active-C dependency audit and select
an independent leaf whose relevant tests pass. The remaining blob-tool and
render/audio failures need diagnosis before touching those areas, but are not a
blanket blocker for unrelated conversions. Baseline Player.plr differences also
remain unexplained. No C implementation has been replaced yet.

## Bounded failure diagnosis — 2026-09-10

- Blob formatter: combining two dynamic Go offset terms dropped the joining +,
  e.g. uintptr(x)*13+71276+uintptr(y) became invalid Go. Added regression cases
  covering operand order, nested sums, subtraction and zero offsets, checking
  parsing, idempotence and independently evaluated arithmetic. Tests fail before
  the fix and pass afterward; TestFormatAccesses now passes on the source copy.
- Remaining ReadBlobs failure is a separate obsolete storage-format assumption:
  it expects root memmap.go/memmap.c/GAME_data.c, while current code uses legacy
  shims plus embedded .dat files and generated pointer initialization. Simply
  changing paths to cgo_blobs.c would incorrectly treat initialized data as zero.
  Updating the split/write tool to that format is separate work, not needed for
  the initial dependency inventory. Do not use it to rewrite current blobs yet.
- Particle rendering diagnosis: 386 and amd64 produce identical pixels for all
  six particle cases. Current PNGs re-encode identically with Go 1.19.13, 1.23.12
  and 1.26.0. The old color library used RGB max 248; the pinned current library
  expands that value to 255. Applying the old expansion reproduces all six
  original PNG goldens exactly. Migrated ONLY those independently verified cases
  to dimensioned, little-endian framebuffer-word SHA-256 hashes; pixel/size
  mutation and subimage-stride checks pass, as do particle tests on 386/amd64.
  No production rendering behavior or sprite goldens were changed.
- Sprite diagnosis: the exact golden-era revision c62202f5 passes the sampled
  APA00001/default case with the same assets. In a disposable current dependency
  copy, restoring only historical color/rgba5551.go makes the entire sprite test
  matrix pass. Export-only normalization was insufficient because the conversion
  also affects intermediate inputs. Keep sprite goldens unchanged pending an
  explicit color-behavior decision before porting that path.
- Audio diagnosis: three dialogue files have equal sample counts on 386/amd64;
  only 206/205/52 samples differ respectively, each by at most one int16 unit.
  Diagnostic 386 SSE floating-point flags produce byte-identical amd64 PCM and
  original hashes. Production flags and goldens remain unchanged. This finding
  covers these samples, not all audio. See the tracked detailed report below.
- Final full default 386/CGO suite with assets: 14 passing packages, 3 failing
  packages, 32 without tests, with no compilation/vet failures. Remaining failing
  packages are blobs (only TestReadBlobs), noxrender (sprite references), and ail
  (audio references). Log: build/diagnosis/final-suite.jsonl. Tracked engine files
  remained unchanged by execution. No new game behavior was introduced, so the
  previously verified three builds and gameplay baseline were not rerun for
  these formatter/test-only changes.

The bounded diagnosis is complete. Durable findings, measurements, reproduction
commands and follow-up decisions are in
[docs/porting/FAILURE_DIAGNOSIS.md](docs/porting/FAILURE_DIAGNOSIS.md).
Next: use the saved build graphs for a bounded compiled/linked C inventory,
choose an independent leaf, and establish its C-reference differential tests
before conversion. Do not require the entire baseline suite to be green, and do
not use obsolete blob writers or unresolved render/audio goldens as port oracles.

## First conversion: protection checksum — 2026-09-10

- Bounded build inventory: standard/highres select 152 repository C translation
  units, server 151. Both checksum symbols are retained in all baseline binaries.
  Read-only reproduction tool and dependency scope: docs/porting/C_INVENTORY.md.
- Committed pre-conversion reference tests and inventory as 00228a81. Existing
  Go checksum agrees with the untouched historical C functions on 386.
- Replaced the two production C checksum definitions with Go exports calling a
  shared internal/protection implementation. Retained the original C only behind
  porttest. C ABI width, return bits, null handling, word/tail boundaries,
  unaligned buffers, chunk boundaries and non-mutation checks pass. Differential
  tests pass for default/server/highres; both bounded fuzz runs pass.
- All three production builds succeed. Test reference symbols are absent from
  their binaries. Fresh warrior scenario checksum-port exits 0 against both
  preserved screenshots with overrides disabled. Full suite has 15 passing,
  3 known failing, 32 no-test packages, with no compile/vet failures.
- Remaining C-to-Go calls have measurable overhead; the local benchmark and
  interpretation are recorded in docs/porting/PROTECTION_CHECKSUM.md. Do not claim
  this conversion improves performance or covers every surrounding caller.
- Production .c physical LOC is now 142,637 (−28), 153 files. Test-only reference
  is 33 lines separately. The user requested counts after EVERY conversion chunk;
  tools/porting/c_loc.py and docs/porting/C_LOC.md define and track this measure.

Next session: read docs/porting/PROTECTION_CHECKSUM.md, retain its differential
reference tests, and select the next cohesive leaf with caller/state evidence.
One checksum implementation plus its nullable wrapper has now moved out of C;
no broad protection-manager or render/audio conversion has been attempted.

Conversion commit: 66fa7bd4; pushed to dbenamy/opennox dev with the preceding
reference-test commit 00228a81. Only the original media archive is untracked.

## Retire checksum C test reference — 2026-09-10

At the user's request, removed internal/protectionref after the completed
conversion comparisons. The historical C implementation and differential harness
remain recoverable from 66fa7bd4. Keep the permanent Go fixed-value/property tests
and tagged ABI tests; the latter now calculate expected results independently by
byte lane, preserving alignment, chunk, mutation and nullable-length coverage.
The benchmark retains direct Go and C-to-Go paths only. Production code is unchanged.

Validation: 386 Go unit tests and TestProtectionABI pass; a five-second ABI fuzz
run passes 551,441 cases. Logs: build/port-checksum/retire-{unit,abi,fuzz}.log.
Production C count stays 142,637 lines across 153 files; test-reference C drops
from 33 to 0. C_LOC.md and the handoff reflect this retirement. No engine rebuild
or gameplay rerun was needed for this test-only removal.

## Protection record helpers — in progress, 2026-09-10

Selected sub_56F590 (decoded-ID lookup), sub_56F6F0 (index lookup), and
sub_56F720 (payload swap). All are retained in the standard baseline binary.
Their only diagnostic callback, nullsub_31, is an empty C function. Tests execute
2,000 deterministic scenarios against current C, with temporary C-heap records
and restored globals: empty/single/multiple lists, duplicates, high-bit keys/IDs,
missing and extreme indices, null/self/adjacent/non-adjacent swaps, preserved
links and modulo-32-bit counter increments. Current C and new pure Go helpers
pass separately before rewiring the ABI. Logs: build/port-records/c-before.log
and unit.log. Production C count is still 142,637; no extra C reference is needed.

Protection record helper conversion completed: the same ABI scenarios pass on
386 for default/server/highres; pure Go tests pass on 386/amd64, all three targets
build with Go export bridges, and records-port exits 0 against both preserved
screenshots. The global layout assertions pass. No C reference was added.
Production C: 142,570 lines (−67 this chunk), 153 files; test-reference C: 0.
See docs/porting/PROTECTION_RECORDS.md for scope, commands and limitations.
Next chunk: inspect and test the protection spell/ability bitset operations;
keep allocation, rekeying and floating-point state outside that scope.

## Protection bitset operations — in progress, 2026-09-10

Current C passes 5,000 deterministic state scenarios plus 12,291 direct bit
checks before conversion. Tests cover signed handle thresholds, empty/missing
records, key zero/high bits, enabled/disabled awards, checksum deltas, successful
and failed validation, ignored entry zero, modulo-32 collisions, signed truthy
values and count <= 1 with null data. Inputs and non-payload state stay unchanged.
Logs: build/port-bitset/c-before.log and unit.log. No C reference is copied.

Protection bitset conversion completed: the same C-before/Go-after ABI checks
pass for default/server/highres; pure Go tests pass on 386/amd64. All three
binaries build with the Go exports, and bitset-port exits 0 against both preserved
screenshots. Existing Go wrappers avoid a C round trip. Production C is 142,503
lines (−67 this chunk), 153 files, with zero C reference lines. Details:
docs/porting/PROTECTION_BITSET.md. Next inspect integer/float record construction,
including exact float bit patterns and allocation-failure state handling.

## Protection constructors — in progress, 2026-09-10

Current C constructors pass 16,224 cases (1,014 bit patterns × four keys × four
C/Go integer/float call paths) before replacement. Patterns include signed zero,
subnormal boundaries, infinities and NaN payloads plus seeded random values.
Tests use an empty C-owned manager, verify exact words/checksum/list endpoints,
and restore globals/free records. Empty insertion draws no randomness. The new
pure Go initializer separately checks nil-allocation state and reset links.
Logs: build/port-create/c-before.log and unit.log. No C reference is copied.

Protection construction completed. A generated C-to-Go float export failed the
signaling-NaN test (7f800001 became 7fc00001); no expected values were relaxed.
A caller audit showed no remaining C caller for the float constructor once its
Go wrapper calls the shared initializer directly. Removed that unused C entry
point/declarations rather than retaining a float shim. Integer C entry remains.
All 12,168 live-path constructor cases and earlier ABI tests pass under all three
tags; pure Go tests pass on 386/amd64. All accepted binaries build and have the
expected symbols. The create-port scenario exits 0 against both screenshots.
Full suite: 15 passing, 3 known failing, 32 no-test packages, no compile/vet errors.
Use build/port-create/accepted-* artifacts; earlier outputs are diagnostics.
Production C is 142,458 lines (−45), 153 files; test-reference C is zero.
See docs/porting/PROTECTION_CREATE.md. Next: record deletion and manager cleanup;
only the delete-and-clear operation has remaining C callers, so preserve that
ABI while routing existing Go cleanup directly to Go.

## Protection deletion/cleanup — in progress, 2026-09-10

Current C passes 1,000 deterministic removal/cleanup sequences, including exact
surviving payloads/links/endpoints, head/middle/tail and duplicate-ID removal,
misses/repeated deletion, zero and UINT_MAX IDs, checksum updates, uint16 count
wrap and cleanup resets. The handle sequence remains unchanged. Fixtures own
individual C allocations, call the public cleanup wrapper and restore globals.
Pure Go unlink tests pass separately. Final pre-port log:
build/port-remove/c-before-final.log. Only delete-and-clear has remaining C
callers; cleanup and the internal delete-by-ID entry can become direct Go.

Protection deletion/cleanup completed: original-C and Go-after state sequences
pass, all accumulated ABI checks pass for default/server/highres, and pure Go
checks pass on 386/amd64. All three builds have the expected retained/removed
symbols. remove-port exits 0 against both preserved screenshots. Production C
is 142,393 lines (−65 this chunk), 153 files; C references remain zero.
See docs/porting/PROTECTION_REMOVE.md. Next: randomized record insertion, with
explicit comparison of list order and RNG index/consumption under fixed seeds.

## Protection randomized insertion — in progress, 2026-09-10

Original C passes 500 deterministic insertion sequences plus prepopulated
32,768/65,535-record boundaries. Tests compare exact list order, back links,
endpoints, checksum, count wrapping, both constructor paths and Logic/Other RNG
indices. Pure Go InsertBefore tests pass separately. Logs:
build/port-insert/c-before-final.log and unit.log. Production C remains 142,393
lines. The sole production caller is the Go constructor; remove the obsolete C
entry point after equivalence validation. Trial delegation: Terra implements
this bounded conversion; the primary agent reviews and runs integration checks.

Randomized insertion completed: primary review accepted Terra's implementation
without corrections. All accumulated protection ABI tests pass under three tags,
pure helpers pass on 386/amd64, all production binaries build with expected
symbols, and insertion-port exits 0 against both screenshots. Evidence:
build/port-insert. Production C: 142,351 lines (−42), 153 files; C references: 0.
See docs/porting/PROTECTION_INSERT.md. Next: reserved-record initialization and
handle allocation, including uint32 sequence wrap and return-value semantics.

## Protection reserved records/handles — in progress, 2026-09-10

Original C passes 400 deterministic mixed operation sequences against full list,
checksum, handle sequence, return-value and RNG snapshots. Cases include empty
and prepopulated lists, zero/high-bit IDs, the threshold and uint32 wraparound.
Final pre-port evidence: build/port-handles/c-before.log. Both C entries remain
needed by sub_56F1C0. Reserved initialization increments the sequence after each
attempt; ordinary allocation increments only on success. Allocator exhaustion is
not injected. Production C remains 142,351 lines; reference C remains zero.

Reserved records/handles completed: original-C and Go-after sequences pass,
including all accumulated protection tests for default/server/highres. All three
production builds have the required Go-backed C exports. handles-port exits 0
against both preserved screenshots. Production C: 142,327 lines (−24), 153
files; reference C: 0. Evidence: build/port-handles; see
docs/porting/PROTECTION_HANDLES.md. Next: record rekey/shuffle, testing exact
payload order, unchanged links, checksum resets and RNG/counter consumption.
