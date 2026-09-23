# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: 25 physical lines in three production `.c` files**, zero
reference C. C types, preambles, generated bridges, libc and third-party MP3 C
remain outside this metric. See [C_LOC.md](docs/porting/C_LOC.md).

## Current — safe-profile direct Go exports qualified

Actual-C baseline `9d5b4589` and conversion `e663a971` are committed/pushed.
The six safe memory/string
functions now export directly from Go using const-qualified pointer typedefs;
cgo_safe.c is removed. Allocator semantics, ASan and macro remapping are preserved.
All 142 frozen cases and the 129-case shop consumer pass under safe,porttest.
Safe build/static, fresh three production builds/ABI, exact known-suite comparison,
headless character creation and explicit save/load pass. Six old `_go` symbols
are absent; ten live callback identities remain distinct. Only the safe Go source
and deleted C file differ from the baseline fingerprint; no golden changed.
See [SAFE_BRIDGES.md](docs/porting/SAFE_BRIDGES.md) and qualification JSON.

Pipelines52607 and25811 JOINED PASS; scenario assets were safely deduplicated.
Finalizer finish-native.py is consumed after success. Check Git log/remote for
this checkpoint's commit/push. Artifacts are build/port-safe-direct/native-
{safe,preflight,production}; baseline c-source.json/c-probe captures are retained.

## Current work — empty callback C baseline qualified

The ten live empty callbacks are still C. Their 338-case capture is frozen at
848b76f173284c29edddd5d637863061d643428004ed8f51098f51863ea44b77. Final baseline
passes 69/69/69 roots in default/server/highres plus all 24 callback consumers
under safe, no skips. Existing goldens are unchanged. Production evidence is
explicitly reused after checking unchanged production source, four Go file
selections, binary hashes and ten callback addresses per binary. See
[EMPTY_CALLBACKS.md](docs/porting/EMPTY_CALLBACKS.md) and C qualification JSON.

A tagged spell-fixture repair was necessary: setup now saves/clears/restores only
fourteen type caches. Snapshot retains four historical zero columns for the
retired duration-state gap; live duration list/records remain captured separately.
The original safe run failed during setup, and the first repair failed in its
snapshot loop (primary missed that loop initially). Both failures are preserved.
A broadened safe run then failed in unrelated client-render fixture setup at
client_effects_environment_porttest.go:76 / registered offset1313532. That broader
safe suite is NOT qualified; its limitation is recorded for later renderer work.

Pipelines75016,98136,7877 JOINED PASS. Pipelines3034,42857,78446 JOINED FAILURE with
the dispositions above. All jobs joined. finish-c-reviewed.py is consumed after
success; original finish-c.py is obsolete. Final artifacts are c-reviewed-
{default,server,highres} and c-final-safe. Check Git log/remote for this baseline's
commit/push. Source must remain frozen during any Go/build job.

Next: after the baseline is committed/pushed, apply reviewed ignored
build/port-empty-callbacks/apply-native.py, source baseline env and gofmt the new
empty_callbacks.go. Run run-native-focused.sh, then run-native-rest.sh. Only then
run compare-dispatch.py with no other build/archive jobs; inspect timings and
write performance-review.json before finish-native.py can accept the conversion.
Preserve every callback name/address identity; all ten must be distinct Go-backed
C exports. Drafts/installers are not reusable after consumption. Do not change
frozen captures to accommodate the conversion.

Luna provided bounded identity/registry drafts and audits; primary reviewed and
wrote contracts/benchmarks. MP3 audit is under build/port-mp3-audit: shipped Dialog
assets include 1,246 MP3-in-WAV files, despite no .mp3 filenames. Primary caught
and corrected the extension-based fixture inference. Existing audio goldens are
known failures on this target, and retained decoder code uses x87 in inspected
paths; the cause of historical PCM mismatches is not established. No decoder
changes or new PCM baseline are implemented.

## Recent qualified milestones

- `f5970121` (pushed): remaining 44 globals/eight mapped buffers initialize from
  Go via the foreign allocator; 114→51 C lines. Full 2,291/2,280/2,291 roots pass.
  Fixture identity cleanup retains persistent aliases and removes stale transient
  addresses. Original intermittent stats mismatch cause remains unproven. See
  [RAW_STORAGE.md](docs/porting/RAW_STORAGE.md) and
  [FIXTURE_IDENTITIES.md](docs/porting/FIXTURE_IDENTITIES.md).
- `cd18015c` (pushed): orphan inline/empty-call cleanup, 51→45 C lines; eight
  additional header/preamble helpers retired. See [ORPHAN_INLINE.md](docs/porting/ORPHAN_INLINE.md).
- Prior `781901ef` (pushed): opaque audio addresses kept outside Go pointer
  scanning. Deterministic regression and full qualification passed. See
  [AUDIO_ADDRESS_GC.md](docs/porting/AUDIO_ADDRESS_GC.md).

## Disk and assets

Preserve `build/assets/extracted/drive_c/Nox` and untracked
`nox-iso-from-archive-org.7z`; never stage the archive. Recent cleanup plans and
consumed records are under `build/port-artifact-cleanup`:

- Old production executable rounds four/five/six removed 12/3/8 files,
  582,770,236 /145,417,452 /384,739,612 bytes after hash/evidence/host-use checks.
- Safe/preflight round one removed six obsolete outputs (301,307,320 bytes).
  Five hashes were observed, not inferred from qualification records.
- Old test ELF round one removed only five large outputs (288,789,800 bytes);
  the seven small diagnostic probes were kept. No production qualification claimed.
- Old capture archival preserved 78 captures in verified gzip files and reclaimed
  1,707,909,412 bytes. `old-capture-archive-record.json` maps original paths, hashes,
  modes/times and archives. Restore with `archive-old-captures.py --restore
  <original-relative-path>` before rerunning old scripts requiring plain JSON.

Capture archive round two also completed: 42 files, 751,577,141 bytes reclaimed
after metadata/hash/156 host-process checks and full decompression verification.
Restore via build/port-artifact-cleanup/archive-old-captures-round2.py --restore
<original-relative-path>; its round2 record maps every original/archive. It joined
before callback benchmarking began. Round three also completed: 121 captures,
1,510,131,029 bytes reclaimed after 158 host-process checks and full verified
round trips; use archive-old-captures-round3.py --restore and its round3 record.
It ran after the retained callback benchmark finished. Total free space after
these archives was about 2.9GiB; recheck before long builds.

Metadata, logs, source, assets and current storage/orphan/entry/audio evidence
were preserved. Earlier cleanup rounds are also consumed; inspect records.
Completed scenario assets are deduplicated with per-run restoration manifests.
Entry scenarios use `build/port-entry-direct/deduplicate-{c,native}-{preflight,
save}.py`; storage uses `build/port-remaining-storage/deduplicate-scoped.py` and
`deduplicate-save.py`; orphan uses `build/port-orphan-inline/deduplicate-preflight.py`
and `deduplicate-save.py`. Consult each run's deduplicated-assets.json before reuse.
