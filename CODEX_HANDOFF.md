# OpenNox x86 Porting Handoff

## Current plan — revised 2026-09-10

This document is the working plan. Read PORTING_STATE.md for completed checks,
exact environment failures, and the next unfinished step. Update that checkpoint
at meaningful milestones so a lost session does not lose progress. For a lost
VM, follow [the recovery instructions](docs/porting/RECOVERY.md); the scenario is
tracked, while assets and generated outputs need separate backup or regeneration.

The user approved the revised approach: establish a 386 baseline and one
repeatable gameplay scenario, then port small cohesive leaves with tests chosen
for each conversion's risks. A comprehensive test framework or whole-engine
call graph is not a prerequisite to the first conversion.

The current VM has no local screen; use headless X. The UTM hardware/network
information below is historical and must not be assumed to describe this VM.
TigerVNC is installed; Xvfb was subsequently installed and used for baseline tests. The user confirmed
that the existing nox-iso-from-archive-org.7z contains ISO asset media for testing.
Extract the needed data into a separate local data directory, preserve the
archive, and keep assets and non-redistributable outputs out of Git.

## Purpose

Continue the investigation and implementation work needed to replace OpenNox's
remaining C code with Go. For now, build and test the actual 32-bit x86 Linux
client inside the UTM/QEMU guest. Establish a trustworthy baseline before
changing code, then use small, behavior-preserving porting steps with strong
feedback loops.

The broader questions motivating this work are:

- How much active C remains, rather than merely how much C text is present?
- What order should the remaining subsystems be ported in?
- Which existing and new tests can detect behavioral, ABI, memory-layout,
  network-protocol, graphics, audio, timing, and performance regressions?
- Once the C dependency is reduced, what work would remain for native macOS and
  browser/WebAssembly targets?

Do not begin with the browser or macOS port. First make the current Linux x86
build reproducible and measurable.

## Repository state recorded on the macOS host

- Upstream: `https://github.com/opennox/opennox.git`
- Branch: `dev`
- Commit: `b184030e76be2b681a7f6d2bcdef52b091d94b9b`
- Recorded: 2026-09-03
- The root `README.md` identifies `dev` as the release branch.
- `src/go.mod` requires Go 1.25.0.
- `src/internal/noxbuild/main.go` deliberately sets `GOARCH=386` and
  `CGO_ENABLED=1` for CGO builds.
- `docker/Dockerfile_client` is a useful reference environment: Go 1.25,
  `gcc-multilib`, and both amd64/i386 SDL2 and OpenAL development packages.

At handoff time the host working tree had one pre-existing untracked file:
`nox-iso-from-archive-org.7z`. It belongs to the user, contains game media/data,
and must not be changed, deleted, inspected unnecessarily, or committed.
This handoff file is also intentionally untracked initially.

Verify the guest checkout rather than assuming it exactly matches the host:

```bash
cd ~/src/opennox
git remote -v
git branch --show-current
git rev-parse HEAD
git status --short
```

Preserve all unrelated and pre-existing changes.

## VM environment

The guest is Ubuntu 26.04 Server amd64 under UTM/QEMU on an Apple Silicon Mac.
It is intentionally using full x86_64 emulation, not ARM virtualization, so it
can compile and execute the `linux/386` CGO client natively from the guest's
point of view.

Known configuration/state:

- UTM name: `OpenNox (Ubuntu 26.04 Server, x86_64, 1 core)`
- One emulated CPU, chosen because SMP had stability concerns
- 8 GiB RAM and a 2 GiB JIT cache
- VirtIO disk, network, and display devices
- Xfce/lightdm desktop installed and working
- SSH works from a normal macOS terminal, but the Codex desktop app was not
  granted macOS Local Network access. Do not depend on host-side Codex SSH.
- Codex CLI is installed in the VM and authenticated to the user's account.
- Apparent pauses and Linux watchdog messages occurred when the host/VM slept;
  keep the Mac awake during long builds and distinguish sleep starvation from a
  reproducible guest failure.

The host repository is exposed through UTM VirtFS with the 9p tag `share`.
A typical mount is:

```bash
sudo mkdir -p /mnt/utm
sudo mount -t 9p -o trans=virtio,version=9p2000.L share /mnt/utm
```

Do not build directly on the 9p mount; metadata-heavy Go builds are especially
slow there. Work from the guest's native virtual disk, for example:

```bash
mkdir -p ~/src
rsync -a --exclude='nox-iso-from-archive-org.7z' /mnt/utm/ ~/src/opennox/
cd ~/src/opennox
```

If edits must be returned to the Mac, prefer Git commits/patches or a carefully
scoped reverse `rsync`. Never use a broad destructive sync or overwrite unknown
host changes.

## Expected guest dependencies

These were the intended packages. Audit what is actually installed before
adding or changing anything:

```bash
sudo dpkg --add-architecture i386
sudo apt update
sudo apt install \
  git rsync curl ca-certificates \
  build-essential gcc-multilib pkg-config ccache \
  golang-go gdb-multiarch \
  libsdl2-dev libsdl2-dev:i386 \
  libopenal-dev libopenal-dev:i386 \
  libgl1-mesa-dev libgl1-mesa-dev:i386 \
  mesa-utils
```

Ubuntu 26.04 reportedly supplies Go 1.26, which should satisfy the module's Go
1.25 requirement. Confirm rather than assuming:

```bash
uname -a
go version
gcc --version
gcc -m32 -x c -o /tmp/opennox-c32-smoke - <<'EOF'
int main(void) { return 0; }
EOF
file /tmp/opennox-c32-smoke
```

The final `file` output should identify a 32-bit Intel/i386 ELF executable.

## First milestone: establish the baseline

Work in `~/src/opennox`, not `/mnt/utm`. Capture command lines, elapsed time,
peak memory if practical, and complete failure output.

1. Verify repository, current VM, and toolchain state. Establish writable,
   persistent Go module/build caches and download the pinned dependencies.
   Do not update dependency versions as part of baseline setup.
2. Run existing tests explicitly for 386 with CGO. Record passed, failed and
   skipped tests separately, including asset and external-service requirements.
3. Build server, standard client and highres client through the build driver.
   Preserve baseline binaries, hashes, revision, toolchain, flags and full logs
   in a local artifact directory before making engine changes.
4. Verify ELF architecture and help/startup behavior. The minimal 386 program
   ran outside the tool sandbox but received SIGSYS inside it; use the supported
   execution/approval mechanism for target tests rather than changing architecture.
5. Extract needed user-provided assets and establish one repeatable scenario:
   startup, menu/input, map load, a bounded stretch of gameplay, and clean exit.
   Use an isolated headless X display and isolated saves/configuration. Record
   display geometry, renderer and audio backend alongside the scenario.

Starting commands after cache and execution setup:

```bash
cd /root/src/opennox/src
go mod download
GOARCH=386 CGO_ENABLED=1 \
  CGO_CFLAGS_ALLOW='(-fshort-wchar)|(-fno-strict-aliasing)|(-fno-strict-overflow)' \
  go test -json ./...
go run ./internal/noxbuild
readelf -h opennox opennox-hd opennox-server
./opennox -h
./opennox-hd -h
./opennox-server -h
```

Capture complete outputs and exit statuses. Test relevant `server` and `highres`
build-tag variants as well; an untagged test run does not cover those branches.
The build driver sets 386/CGO itself and defaults to writing all three binaries
in the current directory. Use its `-o` option to preserve baseline outputs.
Separate environmental failures, source failures and skipped coverage; never
silently exclude a failing test to declare the baseline green.

Exit criterion: all three targets build, existing test outcomes are documented,
and the baseline scenario has repeatable observations. Resolve source failures
or explicitly document baseline defects before evaluating a conversion against
that baseline. Keep unverified behavior visible.

Reuse caches and avoid unnecessary clean builds. Capture elapsed time and peak
memory where practical; defer repeated cold builds until they answer a concrete
question. Measure performance against the same VM configuration, and do not
change VM CPU configuration as part of this work.

## Initial C inventory (not a completion estimate)

A simple source-tree count at the recorded commit found:

- 153 `.c` files, all under `src/legacy`
- approximately 142,665 physical lines across those `.c` files
- approximately 8,965 physical lines across headers under `src`

This raw count substantially overstates the useful porting denominator. The
large `GAME*.c` decompilation units, compatibility shims, CGO glue, vendored or
special-purpose code, already replaced/dead functions, and code compiled only
for particular targets must be classified separately. Do not report "142k
lines remain" as the answer.

Produce a reproducible, incremental active-C audit. First identify compiled
translation units per target and enough dependencies to choose the first leaf;
expand subsystem classification as porting proceeds. Track separately what is
compiled, retained by the linker, and observed executing in measured scenarios.
A retained symbol is not proof of semantic reachability; missing runtime coverage
is not proof of dead code. Include indirect callbacks and C/Go boundary edges.

Over time, determine:

- which C translation units are compiled into each of `opennox`, `opennox-hd`,
  and `opennox-server`;
- active C functions/symbols that are reachable from the final binaries;
- C-to-Go and Go-to-C call edges;
- active versus dead/replaced functions in the large `GAME*.c` files;
- code that should remain C because it is a codec, tiny portability shim, or
  external implementation, versus game-engine logic intended for Go;
- counts by subsystem and a dependency-aware migration order.

Prefer measurements from the actual build graph, object files, linker maps,
and symbols over filename or line-count guesses. Save scripts used for the
audit in a sensible repository location only if they are generally useful and
reviewable.

## Feedback loops for the port

Use the smallest meaningful feedback loop for each conversion, broadening it at
milestones or when new risks appear. The baseline has now run; see PORTING_STATE.md for outcomes and failures.
Reuse these existing tests and harnesses:

- `src/client/noxrender/*_test.go`: rendering primitives, particles and images;
  some paths explicitly skip. Establish which fixtures and branches actually run.
- `src/legacy/client/audio/ail/audio_test.go`: decoded PCM hash/size checks using
  external game data. Reuse these before adding another audio harness.
- `src/e2e.go`, `src/e2e_platform.go`, `src/e2e_seat.go`: YAML input, simulated
  time, seeded platform RNG, screenshots and save hashes. Two fresh standard-client
  warrior scenarios have matched decoded pixels at both gameplay checkpoints.
  The screenshot oracle has since been repaired: missing or mismatched goldens
  fail, and updates require NOX_E2E_OVERRIDE=true. Preserve the independent
  comparator as a cross-check.
- `src/replay.go`: additional scenario mechanism, not yet validated.

### Per-conversion loop

Source-rewriting blobs and noxfactor tests now use temporary source copies;
normal token and PNG diagnostics also use temporary paths. Verify the working
checkout remains unchanged after broader checks. When investigating unknown
source tools, use disposable copies, freshly created per variant. Keep the
checkout used to build baseline binaries clean.

1. Choose a small cohesive leaf; identify callers, callbacks, shared state and
   observable behavior. Add tests before replacing its C implementation. Retain
   C exports only where remaining C callers need them; do not create an ABI
   boundary solely to keep a test calling a retired internal entry point.
2. Keep a callable C reference while comparing C and Go on identical inputs.
   Cover boundary cases, return values, mutations, signedness/overflow,
   serialization, RNG consumption and timing effects where applicable. Reset
   global state or run each implementation in separate processes. Compare logical
   fields rather than raw pointer-containing memory or uninitialized padding.
3. Assert sizes, alignment and offsets for translated types crossing the C/Go
   boundary on the actual 386 target. Run formatting, affected-package tests and
   the relevant target builds. Use generated/fuzz inputs where they add coverage.
4. Run the scenario relevant to the change and compare against the preserved
   baseline. Only remove the C reference after equivalence checks pass; retain
   useful fixtures and a recoverable reference revision. Keep changes scoped and
   record exactly what was verified and what remains untested. After each completed
   conversion chunk, run `tools/porting/c_loc.py` and update
   `docs/porting/C_LOC.md` with the new production C line count and delta, keeping
   test-reference C separate. The user explicitly requested this ongoing checkpoint.

### Headless integration loop

Start with one bounded scenario, then extend it for the code being ported. Use
observable readiness and fixed simulation ticks rather than arbitrary UI sleeps.
Control seeds, input, configuration and assets where possible. Compare selected
logical state, serialized outputs and frames at defined checkpoints; first prove
that unchanged baseline runs agree. Normalize only known nondeterministic fields.

Prefer exact framebuffer comparisons for deterministic software rendering.
Introduce tolerances only for demonstrated backend variability, retaining useful
failure images. Headless X tests window/input/rendering integration; direct
renderer tests help isolate pixel-generation errors. Audio initialization and
PCM correctness are separate checks. Physical display and audible playback
quality remain manual release checks.

### Milestone checks

- Existing suite under explicit 386/CGO settings, relevant tag variants, all
  three builds, and startup/shutdown checks.
- Server map loading, scripted ticks, save/load and old/new client-server
  combinations; extend into multiplayer join/play/leave and protocol fixtures
  as network or stateful engine work requires.
- C warnings and supported sanitizers, plus applicable Go checkptr/race checks.
  Verify tool support first: do not assume race/ASan tooling works on 386 CGO.
  Tests on another supported architecture supplement, not replace, target checks.
- Repeatable tick/frame, allocation, load-time and memory measurements when a
  conversion could affect them. Use multiple warm runs in the same environment.
- Representative single-player and multiplayer manual release testing.

The testing plan is revisable: document new findings and adjust the next test to
address the actual failure modes of the next conversion. Avoid building a broad
harness ahead of evidence that it is needed.

## Constraints and cautions

- Original Nox assets are not part of this repository. Use only the user's
  legitimately owned data and never commit it or generated derivatives that
  cannot legally be redistributed.
- Treat all uncommitted files as user-owned.
- Avoid destructive Git operations and broad file synchronization.
- Do not silently update dependencies or reformat unrelated legacy files.
- `GOARCH=386` plus CGO is intentional for the current client; do not "fix" it
  to amd64 as a shortcut.
- Report environmental failures separately from source failures.
- The VM's speed is not representative of native x86 hardware. Use it for ABI
  correctness and functional compatibility; interpret performance results
  relative to a stable baseline in the same VM.

## Next session / next action

Read this plan and PORTING_STATE.md. The pristine 386 builds, help checks and
repeatable standard-client gameplay scenario are established. The full test suite
is red with documented baseline failures; do not treat these as port regressions
or regenerate goldens to hide them.

The E2E screenshot oracle, source-rewriting test isolation and stale test/tool
API compilation failures have been repaired in separate commits. The bounded rendering, audio and blob-tool diagnosis is complete; see
[the findings](docs/porting/FAILURE_DIAGNOSIS.md). Formatter and particle tests
have verified fixes; sprite color semantics, target-dependent PCM references and
blob storage modernization remain explicit follow-ups before touching those paths.
Establish save-load checks: map save bytes
matched across the two runs, but Player.plr bytes did not.

Continue the bounded C dependency audit and choose the next cohesive leaf,
using the checksum conversion as the first proven differential-test pattern. Update PORTING_STATE.md
with commands, outcomes, artifact locations and remaining work at each checkpoint.
The first conversion is the protection byte checksum; see PORTING_STATE.md for
validation and [C source-size checkpoints](docs/porting/C_LOC.md) for its reduction.

The checksum C reference was retired at the user’s request after validation.
Use its retained Go and tagged ABI tests; recover the original differential
harness from commit 66fa7bd4 if needed.

Protection record lookup/index/swap is also converted; see
[its validation](docs/porting/PROTECTION_RECORDS.md). Continue with the next
bounded section, committing/pushing and reporting each completed chunk. The user
authorized continuing onward without waiting for a new “go” unless a substantive
decision needs their input.

Protection spell/ability bitsets are converted too; see
[the checkpoint](docs/porting/PROTECTION_BITSET.md). Latest source counts are in
C_LOC.md. Next candidate is record construction, with exact float-bit tests.

Integer/float record construction is converted;
[its checkpoint](docs/porting/PROTECTION_CREATE.md) records the NaN ABI finding
and removed unused float entry point. Continue with deletion and cleanup,
retaining only the C entry points that have actual remaining callers.
