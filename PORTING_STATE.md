# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: 25 physical lines in three production `.c` files**, zero
reference C. C types, preambles, generated bridges, libc and third-party MP3 C
remain outside this metric. See [C_LOC.md](docs/porting/C_LOC.md).

## Current — safe-profile direct Go exports qualified

Actual-C baseline `9d5b4589` is committed/pushed. The six safe memory/string
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

## Next — empty callback behavior and performance baseline

Ten distinct empty C callbacks remain in common__object__modifier.c/GAME5_2.c.
Read build/port-empty-callbacks/acceptance-plan.md. Readiness/Replenishment are
both identity-sensitive and actually invoked; modifier paths pass 3/5/6 arguments,
and Energy Bolt destruction passes one to a void(void) definition on this target.
Preserve names/identities and measure the production Go→ccall dispatch path before
accepting extra Go callback transitions. Do not port solely to lower the C count.

Luna drafted bounded declaration-only identity/registry helpers in ignored
build/port-empty-callbacks/drafts; primary drafted guarded-state contracts and a
representative benchmark there. None are installed or qualified yet. Root test
must independently validate expected metadata, not blindly trust the helper table.
Only one Luna helper; primary reviews source and acceptance. No source changes
or concurrent Go/build jobs during active qualification.

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

Metadata, logs, source, assets and current storage/orphan/entry/audio evidence
were preserved. Earlier cleanup rounds are also consumed; inspect records.
Completed scenario assets are deduplicated with per-run restoration manifests.
Entry scenarios use `build/port-entry-direct/deduplicate-{c,native}-{preflight,
save}.py`; storage uses `build/port-remaining-storage/deduplicate-scoped.py` and
`deduplicate-save.py`; orphan uses `build/port-orphan-inline/deduplicate-preflight.py`
and `deduplicate-save.py`. Consult each run's deduplicated-assets.json before reuse.
