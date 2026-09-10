# Recovering the porting workspace

The plan, checkpoint, source changes and `warrior-smoke.yaml` are tracked in Git.
Read `CODEX_HANDOFF.md` and `PORTING_STATE.md` after cloning. Work continues on
`dev` in `https://github.com/dbenamy/opennox.git`; no C conversion has begun yet.
Infrastructure changes through `e694e1ac` were pushed on 2026-09-10.

## What needs a separate backup

Git does not contain the user's original Nox assets, extracted data, screenshots,
saves, binaries, caches or raw test logs. Keep the source asset archive in your
own backup. Rebuilding can regenerate binaries and test output, but cannot recover
lost original assets or prove old results without rerunning the checks.

Do not commit SSH private keys, access tokens, game serials, assets, screenshots
or raw logs. In particular, `internal/version/TestLatestGithub` logs the token
returned by its registry request; raw suite logs need review before sharing.
Git author names/emails, repository URLs, revisions and ordinary local paths are
normal project metadata. SSH credentials are provisioned separately on a new VM.

## Recreate the 386 build environment

Use an x86_64 Linux machine with 32-bit execution support. The measured environment
was Ubuntu 26.04.1 with Go 1.26.0, multilib GCC, i386 SDL2 and OpenAL development
packages. Preserve the toolchain version when comparing against recorded goldens.
The old UTM settings in the handoff are historical, not requirements for a new VM.

Relevant Ubuntu packages include `gcc-multilib`, `g++-multilib`, `pkg-config`,
`libsdl2-dev:i386`, `libopenal-dev:i386`, `libgl1-mesa-dev:i386`, `xvfb`, `xauth`,
`libarchive-tools` and `squashfs-tools`. Enable the i386 package architecture
before installing its packages. Install a Go toolchain satisfying `src/go.mod`.

From the repository root, configure caches and target settings in your shell:

```bash
export NOX_PORT_ROOT="$PWD"
export GOMODCACHE="$NOX_PORT_ROOT/build/cache/gomod"
export GOCACHE="$NOX_PORT_ROOT/build/cache/go-build"
export GOARCH=386 CGO_ENABLED=1 CC=gcc CXX=g++
export CGO_CFLAGS_ALLOW='(-fshort-wchar)|(-fno-strict-aliasing)|(-fno-strict-overflow)'
export PKG_CONFIG_LIBDIR=/usr/lib/i386-linux-gnu/pkgconfig:/usr/share/pkgconfig
mkdir -p build/cache/gomod build/cache/go-build build/recovery/bin build/recovery/logs
cd src
go mod download
go run ./internal/noxbuild -o ../build/recovery/bin
readelf -h ../build/recovery/bin/opennox
# Expected known failures are documented in PORTING_STATE.md.
go test -json -count=1 -timeout 3m ./... > ../build/recovery/logs/tests.jsonl 2>&1
cd ..
```

Record exit codes, including failures. Repeat relevant tests with `-tags server`
and `-tags highres`. Use `NOX_DATA` pointing to supplied assets to exercise the
asset-dependent tests; a run without them does not establish their correctness.

## Recover the gameplay scenario

The supplied `.7z` contains `Nox.wsquashfs` with installed files at `drive_c/Nox`.
Extract only this game subtree into a new ignored local directory:

```bash
mkdir -p build/recovery/media
bsdtar -xf nox-iso-from-archive-org.7z -C build/recovery/media
unsquashfs -no-progress -d build/recovery/extracted \
  build/recovery/media/Nox.wsquashfs drive_c/Nox
```

Use a fresh data copy for each scenario run, excluding the archive's saved games.
Do not run bundled executables. Keep the original data and archive unchanged.
For example, from the repository root (the destination must not already exist):

```bash
python3 - <<'PY'
from pathlib import Path
import shutil
src = Path('build/recovery/extracted/drive_c/Nox')
dst = Path('build/recovery/run-a/data')
shutil.copytree(src, dst, ignore=lambda path, names: [n for n in names if n.lower() == 'save'])
(dst / 'save').mkdir()
PY
mkdir -p build/recovery/run-a/testdata
cp docs/porting/warrior-smoke.yaml build/recovery/run-a/scenario.yaml
```

The script navigates the menus, creates a warrior, loads `war01a`, closes the
captain dialogue, walks, captures two frames and quits. Raw mouse coordinates
assume an X display of **1280x960**, menus of 640x480 and gameplay of 1024x768.
OpenAL's null backend exercises initialization without a physical audio device.

Normal validation requires existing `gameplay.png` and `after_walk.png` reference
images in the run's `testdata` directory:

```bash
NOX_E2E="$PWD/build/recovery/run-a/scenario.yaml" \
NOX_E2E_OVERRIDE=false ALSOFT_DRIVERS=null \
xvfb-run -a -s '-screen 0 1280x960x24 -nolisten tcp' \
  timeout 120s "$PWD/build/recovery/bin/opennox" \
  -data "$PWD/build/recovery/run-a/data" -window -pprof 127.0.0.1:0 \
  > build/recovery/logs/scenario.log 2>&1
```

If the goldens were lost, rebuild the original baseline revision
`b184030e76be2b681a7f6d2bcdef52b091d94b9b` in a separate detached worktree using the
same environment. Run that binary with `NOX_E2E_OVERRIDE=true` to generate local
references, then repeat on fresh data and confirm the decoded pixels agree.
The original baseline's comparator is defective; use the current
`internal/e2etest.CheckScreen` helper to check its two outputs independently.
Only then use those goldens with the current binary and override disabled. Never
bless current output simply because a regression check fails.

For the recorded standard-client scenario, SHA-256 over decoded NRGBA pixel
bytes (1024x768) was:

- `gameplay.png`: `0a113a2a37440e95c4e457e2771d36aa14de1d50becc806273b5db4e91a827ee`
- `after_walk.png`: `32a18672dd49064bb512016d622b53083b0bd1c7f2d5ae8cf8c6d31b227f457e`

These are pixel hashes, not PNG-file hashes. A mismatch after recovery requires
investigation of assets/toolchain/configuration as well as engine behavior.
Player save bytes were not deterministic in the recorded baseline; map save bytes
matched. Save/load compatibility still needs a dedicated check.

## Resume development

Choose a small C leaf using the build dependency graph and current callers.
Establish relevant passing checks and C/Go differential coverage before replacing
it. Unrelated documented failures need not block that conversion. Commit each
coherent verified step, update the checkpoint and push it. Local commits alone
will not survive loss of this VM.
