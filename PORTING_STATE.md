# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: 45 physical lines in four production `.c` files**, zero
reference C. The latest change removes two cgo-preamble bodies, outside that
metric. C types, other preambles, generated bridges, libc and third-party C remain.
See [C_LOC.md](docs/porting/C_LOC.md).

## Current — entry classifiers call libc directly

The actual-C baseline is committed/pushed as `cbf62bc1`. The subsequent direct-libc
change is fully qualified and this checkpoint accompanies its commit. Check Git
status/log/remote for the final commit/push; no self-referential hash is inferred.

`uiEntryDigit`/`uiEntryAlnum` now call the same libc classifiers with explicit
`C.wint_t` widening. All 131,072 boolean classifications match three fresh C
captures: SHA256 `f5c39db5e866885891bbb1fe8d0b2244d9fbcea7d232cdd4b7e6f90dd03ff6e0`.
No locale substitution, reference C copy, or UI golden change was introduced.
Four legacy contracts and twelve widget roots pass in each profile, no skips.
Static, safe build, three production binaries/ABI, exact known-suite comparison,
headless gameplay and explicit save/load pass. Source fingerprints agree;
preflight matches production; all ten retained callback addresses stay distinct.
Safe runtime was not tested; known full-suite failures remain unchanged.

Native pipeline91437 JOINED PASS; its finalizer29700 JOINED PASS. All build/test/
scenario jobs are joined. Artifacts are `build/port-entry-direct/{c,native}-
{default,server,highres,safe,preflight,production}`. The migration/finalizers and
scenario cleanup scripts are CONSUMED. Tracked source supersedes ignored drafts.
See [ENTRY_CLASSIFIERS.md](docs/porting/ENTRY_CLASSIFIERS.md) and both qualification
JSON files. Do not regenerate captures or rerun consumed installers/finalizers.

## Next — optional safe-profile forwarding shims

After the current commit/push, inspect/run the isolated const-pointer export probe
under `build/port-safe-direct/probe`, sourcing `build/baseline/env.sh` first. It has
not been compiled. Its declarations require the generated Go export prototypes
to remain const-qualified. Only proceed with direct Go exports if that succeeds.
The production safe shims are still unchanged (20 C lines); ASan, macro remapping,
allocator behavior and compatibility symbol names must be preserved.

Read ignored `build/port-safe-direct/PLAN.md` and `contract-cases.md`. Luna drafted
the edge matrix; primary reviewed it and corrected report escaping/return-value
scope. Runtime support under safe,porttest still needs establishing. New fixtures
should use the existing root-package/Go-bridge pattern where practical to share
consumer builds. Captures should call the actual six C shims before conversion;
independent compare assertions use sign, while target-specific raw returns may be
recorded separately. Avoid undefined overlap, capacity and pointer inputs.

A read-only helper is tracing whether existing TestShopStockLoading actually
reaches FieldGuideXfer/strcpy; do not infer branch coverage from type29 alone.
Its prospective report is `build/port-safe-direct/shop-fieldguide-trace.md`.
Only one bounded GPT-6 Luna helper is used; implementation and acceptance remain
primary-owned. No safe-shim production/test source changes are applied yet.

Later candidates: ten distinct empty callback exports and third-party MP3 C.
Read `build/port-empty-callbacks/go-entrypoint-design.md` before callback work;
ABI, identity and callback-loop overhead require evidence. Other read-only audits
are under `build/port-inline-c-audit` and `build/port-safe-bridge-audit`.

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
