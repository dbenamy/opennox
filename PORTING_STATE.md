# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: six physical lines in one production `.c` file**, zero
standalone reference C. That file includes the MP3 decoder implementation; C
preambles, headers, generated bridges and external libraries remain outside this
metric. See [C_LOC.md](docs/porting/C_LOC.md).

## Current — distinct empty callbacks qualified

Actual-C baseline `ae76f3c8` is committed/pushed. The working conversion replaces
ten distinct empty C bodies with Go exports, deleting GAME5_2.c and
common__object__modifier.c. All 338 frozen cases match; 69/69/69 normal-profile
consumer roots and 24 safe roots pass with no skips. Safe build/static, three
production builds/ABI, exact known-suite comparison, headless creation and explicit
save/load pass. All four binaries retain ten distinct Go-backed callback addresses.
No golden changed. See [EMPTY_CALLBACKS.md](docs/porting/EMPTY_CALLBACKS.md) and
empty-callbacks-native-qualification.json. Check Git log/remote for commit/push.

Pipelines24759,27583 and finalizer75180 JOINED PASS; all jobs joined. apply-native.py,
compare-dispatch.py, finish-native.py and both native scenario deduplication scripts
are CONSUMED. Artifacts are build/port-empty-callbacks/native-{default,server,
highres,safe,preflight,production} and paired-timing. Do not rerun completed phases.
Paired timing shows 84–140 ns/call added median cost (1.57–2.27x). Primary accepted
this reversible cost in performance-review.json; real game-frame impact remains
unmeasured. Review profiling before any further no-op dispatch optimization.

Baseline fixture limitations remain documented: the spell fixture's retired-state
range was repaired without golden changes; an unrelated renderer fixture still
blocks the expanded optional-safe suite. Failed runs remain recorded and are not
qualification evidence. Original C benchmark/captures remain available.

## Next — five pure address getters

Luna drafted build/port-address-adapters/{draft.patch.txt,README.md}. Primary
reviewed the exact five substitutions across four files and prepared ignored
replacements.json, tests.txt.draft, c-batch.json.draft and primary-plan.md. None is
applied. Run the 38 existing player-control/orchestration roots in three profiles
with the original getters first, verify unchanged production evidence, then apply
and qualify the conversion. No duplicate C oracle or changed golden is needed.

Luna's first coverage report overstated rare-branch and snapshot coverage. Primary
traced setup and snapshots: bot allocation-failure fallback and observer bot-identity
early return are not proven by current fixtures, and actor callback slot744 is not
explicitly captured there. Preserve these limits in the accepted test plan; source
mapping and generated cgo references must also prove the exact target symbols.

Luna's round7 obsolete-production-ELF audit found no eligible candidates under its
exclusions. No files were deleted; no host-use check was needed. Helper is idle.
Free disk was around 1GiB during final qualification; recheck before more builds.

MP3 audit is under build/port-mp3-audit. Dialog assets contain 1,246 MP3-in-WAV
files despite no .mp3 filenames; primary corrected that initial audit inference.
Historical audio goldens are known failures on this target. Inspected retained
C decoder paths use x87, but the cause of historical PCM mismatches is not proven.
No decoder changes or new actual-C PCM baseline are implemented.

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
