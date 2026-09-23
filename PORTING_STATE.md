# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: six physical lines in one production `.c` file**, zero
standalone reference C. That file includes the MP3 decoder implementation; C
preambles, headers, generated bridges and external libraries remain outside this
metric. See [C_LOC.md](docs/porting/C_LOC.md).

## Current — direct callback addresses qualified

Empty callbacks baseline `ae76f3c8` and conversion `dc9692f8` are committed/pushed.
Address-getter baseline `229bf34a` is committed/pushed. The working getter conversion
is fully qualified: five direct references replace five C preamble bodies across
four files. All 38 consumer roots pass per normal profile; safe build/static,
production/ABI, exact known-suite and headless creation/save/load pass. Generated
cgo references confirm the exact four target symbols. All four binaries retain
ten distinct Go-backed empty callback identities. No golden changed.
See [ADDRESS_ADAPTERS.md](docs/porting/ADDRESS_ADAPTERS.md) and native qualification
JSON. Check Git log/remote for the conversion's commit/push.

Pipelines3368,59256 (C baseline),20004 (native focused/safe),56013 (rest) and2544
(finalizer) JOINED PASS; all build/test jobs joined. finish-c.py, applied.json,
finish-native.py and both native scenario deduplication scripts are CONSUMED.
Artifacts: build/port-address-adapters/{c-*,native-*,cgo-addresses,inline-inventory}.
Do not rerun consumed phases. Standalone C stays six lines; heuristic production
preamble count falls 86→81 (76 shared dispatchers and five remaining adapters).

Luna's getter draft was correct but its initial coverage report overstated bot
allocation-failure, bot-identity and actor callback-slot coverage. Primary traced
setup/snapshots and corrected these claims. Exact source/cgo symbol mapping
complements the actual behavioral coverage; no forced allocation-failure test is
claimed. Optional-safe renderer fixture limitations remain documented separately.

Callback timing from dc9692f8 remains a review item: 84–140 ns/call additional
median cost (1.57–2.27x), accepted as reversible; actual frame impact is unmeasured.
Current callback binaries/benchmarks and frozen captures remain available.

## Next — MP3 baseline before numerical decisions

Luna is drafting a tagged ail MP3 asset-baseline test only under
build/port-mp3-audit/drafts/, not in src. Primary must review it before application,
collect repeated actual-C captures and freeze expectations before accepting.
The draft should observe input hashes, format, reported PCM/counts, Decode sequence,
full buffer writes and outside-buffer guards for MP3-in-WAV Dialog assets.
No old failing audio golden is a passing parity oracle; no decoder code changes
or new actual-C PCM baseline are implemented yet.

Read build/port-mp3-audit/{README,active-helper-plan,primary-scope-review}.md.
Dialog assets contain 1,246 MP3-in-WAV files despite no .mp3 extension. Inspected
retained C paths use x87, but the historical PCM mismatch cause is not established.
The bit-reader mutates parser state; inactive LayerI/II widths are not active
coverage. Avoid per-bit C→Go calls: a future Go decoder needs coherent internals
and numerical/state qualification before the production switch. Primary owns
that decision; the helper's source audit is planning, not implementation evidence.

Disk cleanup30220 JOINED PASS: ten verified obsolete entry executables removed,
489,077,624 bytes reclaimed after 159 host-process checks. Record:
build/port-artifact-cleanup/superseded-entry-removed.json; apply-superseded-entry.py
is CONSUMED. Historical entry and dependent safe-baseline finalizers need rebuilt
entry executables before reuse. All metadata/captures/current binaries remain.
Disk dropped near 139MiB during scenario copies and recovered after their verified
deduplication; recheck before further builds. Five safe-direct executables remain
separately inventoried but NOT deleted. Round7 general audit found no candidates.

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
