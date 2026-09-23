# Porting checkpoint

Read [PORT.md](PORT.md) for the working plan. This is the resume checkpoint.

**Qualified C remaining: six physical lines in one production `.c` file**, zero
standalone reference C. That file includes the MP3 decoder implementation; C
preambles, headers, generated bridges and external libraries remain outside this
metric. See [C_LOC.md](docs/porting/C_LOC.md).

## Current — direct callback addresses qualified

Empty callbacks baseline `ae76f3c8` and conversion `dc9692f8` are committed/pushed.
Address-getter baseline `229bf34a` is committed/pushed. Conversion `78a5a21c` is committed/pushed and
fully qualified: five direct references replace five C preamble bodies across
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

## Current — MP3 scalar SSE2 correction qualified

Original-C asset baseline `22627e7f` and SSE correction `08774225` are committed/pushed. The correction
adds package-local `#cgo 386 CFLAGS: -msse2 -mfpmath=sse` in ail/audio_mp3.go.
All 1,246 historical TestAudioDecode PCM goldens pass UNCHANGED in default/highres.
The newer guarded observation expectation deliberately changes from the recorded
x87 hash to the SSE hash, restoring independently established historical behavior:
old 7dc3362576eb00660fd68a836d17af6012c4d1c0595d137958109d62228add0b;
new e0688114198dc50b4acf9b5695ea4cd8bf6dbb9f91d53d1e9d550c02ff1e2763.
Inputs, formats, sample counts and Decode sequences are unchanged. Old captures
remain. This is a separately qualified correctness fix, not a Go decoder port or
controlled speed claim. See [MP3_DECODER.md](docs/porting/MP3_DECODER.md) and
mp3-sse-qualification.json. Check Git log/remote for correction commit/push.

Pipelines45585 (audio clients),56838 (reviewed GC consumers/safe/production/scenarios)
and3009 (finalizer) JOINED PASS. All jobs joined; finish.py and scenario deduplication
scripts are CONSUMED. Pipeline61175 JOINED FAILURE only because primary omitted
--package ./legacy and selected zero tests; preserve default/ failure. Reviewed
runs supply qualification. Artifacts: build/port-mp3-sse/{audio-default,
audio-highres,reviewed-default,reviewed-highres,safe,preflight,production}.
All four binaries retain distinct empty callbacks and expected ABI identities;
fresh decoder disassembly confirms XMM instructions. Source changes are exactly
the compiler directive and new observation hash/comment.

Full suite now has exactly 304 non-audio failure events, 16 pass/2 fail/32 skip
packages. The 1,249 audio failure events are gone. No actual result is filtered:
mp3-sse-expected-suite.jsonl requires audio pass and all old non-audio failures
unchanged. expected-suite-review.json records the preceding log/hash and rule.
Standalone C remains six lines/one file; preamble bodies remain 81.

## Next — coherent Go MP3 integer helpers

Luna drafted build/port-mp3-integer/drafts/{bits_headers.go,README.md}; nothing is
installed. Primary must review against minimp3.h and build independent frozen
actual-C vectors before acceptance. Scope: bs_init/get_bits plus seven header
helpers. Preserve table entries, hdr_compare asymmetry, position advance on overrun,
uint32 reads, int32 free-format fallback; widths0..32 within defined padded input
bounds. LayerI/II headers still participate in frame detection even though their
decode routines are excluded. No per-bit C→Go callbacks: build coherent Go decoder
internals while existing C remains production until the full path qualifies.

Potential baseline approach (not implemented): a temporary ignored probe includes
the actual production minimp3 header with matching macros/SSE2 flags, calls its
static integer helpers, and emits frozen vectors. Do not copy algorithm bodies or
retain an alternate C decoder solely as a test oracle. Keep source/input hashes
and provenance; test the Go draft against those immutable outputs plus independent
bit/order/state invariants. C LOC will not fall until production decoder switches.
Primary owns test design, numerical behavior and final integration; helper is idle.

The shipped MP3 corpus only covers mono/22,050Hz. Stereo/other rates, seek, short
streams and frame/state boundaries still need contracts before a Go replacement.
Current guards detect adjacent overwrites only. See build/port-mp3-audit planning
notes and MP3_DECODER.md. Preserve public-domain attribution.

Superseded-safe/callback cleanup32904 JOINED PASS: ten verified executables removed,
489,093,068 bytes reclaimed after 163 host-process checks. Consumed record:
build/port-artifact-cleanup/superseded-safe-callback-removed.json. Current address
binaries and both callback benchmark binaries remain, alongside source/captures/
metadata/assets/cache. Historical safe/callback/address-C finalizers require
rebuilt old binaries before replay. Free disk after qualification about 522MiB;
check before more broad builds. Small integer-package probes should be bounded.

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

## Current — MP3 integer baseline

819,207 actual-C helper records repeat identically across three native processes;
a UBSan build also matches without diagnostics. Fixture and capture tool are
committed before Go integration. See [MP3_INTEGER.md](docs/porting/MP3_INTEGER.md).
Luna supplied ignored implementation/test drafts; primary reviewed the arithmetic,
all 90 table entries and captured C independently. Production remains unchanged.
Capture sessions71135 and92923 joined PASS; initial sandbox11583 failed SIGSYS.
Standalone C remains six lines/one file, reference C zero. Next: install reviewed
Go helpers and frozen runner, qualify the pure-Go package, prove no production
dependency before reusing the SSE production qualification.

## Current — Go MP3 integer helpers qualified, unwired

C baseline `91f520ae` is committed/pushed. Go helpers and runner pass all819207
vectors in default/server/highres/safe and CGO_ENABLED=0; vet passes (74167 joined).
Four production dependency selections exclude this package. Finalizer verifies
all previous source hashes and four retained SSE binary hashes, so production/
ABI/gameplay evidence is explicitly reused. See mp3-integer-go-qualification.json.
No source/build jobs remain active. C stays six lines/one file, reference zero;
81 production preamble bodies unchanged. Next full suite needs one additional
passing package; the304 unrelated failure events remain the existing baseline.
Luna is drafting side-information parser in ignored build/port-mp3-sideinfo/drafts;
primary must independently capture C behavior/partial writes before integration.
Disk around480MiB free; recheck before broad builds. No raw assets changed.

## Current — MP3 side-information C baseline

Integer helpers `e55c8385` are committed/pushed. The next actual-C baseline contains
36864 side-info records, identical in three native captures and a UBSan capture
(11698/5915 joined PASS). See MP3_SIDEINFO.md and mp3-sideinfo-c-capture.json.
Production remains unchanged. Luna's initial draft lost the selected band-table
row; primary caught this and requested a borrowed-slice representation before
integration. Draft runner/implementation are still ignored and unqualified.
Baseline is committed before installing Go. Next: review corrected representation,
compare all static tables, qualify frozen output/independent partial-state tests.
C remains six lines/one file, reference zero; free disk roughly400MiB.

## Current — Go side-information parser qualified, unwired

Both test roots pass in default/server/highres/safe and CGO_ENABLED=0:819207 integer
and36864 side-info vectors each. Vet passes;82406 joined PASS. All824 table entries
match C. Primary corrected Luna's row representation and constructed-test CRC
prefix before qualification; frozen outputs stayed unchanged. See
mp3-sideinfo-go-qualification.json and MP3_SIDEINFO.md. All prior sources and four
SSE binaries remain identical; production/ABI/gameplay evidence explicitly reused.
C stays six lines/one file, reference zero;81 production preamble bodies unchanged.

Next draft: build/port-mp3-stream/{primary-plan.md,capture_stream.py,generate.py,run.py}.
Preliminary probe12577 joined PASS:1637find/185match/5init/425restore/504save cases,
three native processes plus UBSan identical. Still ignored/unqualified; make capture
reproducible in tools, review domains and commit baseline before Go installation.
Luna is drafting Go stream primitives in build/port-mp3-stream/drafts; no source
edits/builds authorized for helper. All primary build/test jobs joined. Disk about
390MiB free; preserve source assets/archive and current binaries.

## Current — MP3 stream C baseline

Side-info conversion `f4fe693d` committed/pushed. Stream baseline now has a
self-contained tools/porting/capture_mp3_stream.py and2756 frozen records. Three
native runs plus UBSan match the preliminary probe exactly;28198 joined PASS.
See MP3_STREAM.md and mp3-stream-c-capture.json. Baseline committed before Go
integration. Luna drafts stream primitives; primary owns fixture runner this time
so capture/test work can overlap. No active builds. C remains six lines/one file,
reference zero. Disk around380MiB. Next: review helper, install Go/runner, qualify
whole package and unchanged-production proof, commit/push then continue.
