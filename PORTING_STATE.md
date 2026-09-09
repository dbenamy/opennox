# Porting checkpoint — 2026-09-09

Read CODEX_HANDOFF.md for the working plan. This is the resume checkpoint.

## Repository and environment

- dev at b184030e76be2b681a7f6d2bcdef52b091d94b9b; origin is the user's fork,
  https://github.com/dbenamy/opennox.git.
- No pre-existing tracked changes. Only the handoff and media archive were
  untracked initially. No engine edits have been retained or porting begun.
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
- Why independent? src/e2e.go Screen passes nil to e2eError on a mismatch, and
  auto-creates missing goldens. Harness exit status alone is insufficient.
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

1. Repair E2E screenshot failure handling and isolate source-rewriting tests.
2. Fix stale API/test compile errors in scoped changes; retain baseline logs.
3. Diagnose PNG/PCM goldens without masking real behavior changes. Investigate
   Player.plr nondeterminism and add a save-load scenario.
4. Use build graphs for a bounded dependency audit; choose a cohesive C leaf and
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
