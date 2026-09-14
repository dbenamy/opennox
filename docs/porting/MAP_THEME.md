# Theme parser

All **33 routines** at 51E260..520D50 in GAME4_2.c are native Go. Only
mapGenReadTheme and sub_520D50 retain C ABI bridges; the other **31 entries** have
no remaining source references or production symbols. Physical C is **102,419
lines / 148 files / zero reference C**, **−1,881**. This conversion is qualified
and ready to commit/push; corrected C baseline **fc6f4252** was pushed first.

## Implementation and compatibility

The Go parser covers tokenization and nested conditions, player/percentage tests,
algorithm settings, spells/equipment/template modifiers, decoration properties,
wall/floor edging, weighted choices, foreach rules, copies, prefabs and cleanup.
It uses the existing file and engine owners. Records retain C allocator ownership
because remaining map-generator code consumes/frees them. Standard libc numeric
conversion and clock semantics remain production services. No old parser
algorithms are retained solely for testing.

Preserved behavior includes partial tokens at EOF, the shared token used by
conditional filtering even when reading into another destination, bytewise ASCII
case comparison, choice count clamping at 31, wildcard weight rounding, list order,
shallow decoration-copy children and existing cleanup scope. Edging-table matches
continue against the token changed by each nested read. File-open/key error logging
matches the previous wrappers. Ownership changes remain a separate review decision.

## C prerequisites

**ff4de53d** restored an actual 60-byte algorithm value buffer and a contiguous
four-element modifier-count array. Original normal numeric and four-slot probes
aborted; corrected tests and prior map checks passed in all variants. This removed
seven physical C lines. **4fcad1f5** corrected inherited-modifier removal: advance
the copy source after copying, preserving the next value instead of skipping it.
The original alpha/beta/gamma minus alpha probe failed; a 72-case weapon/armor,
slot and edit matrix verifies the correction. It changed no C line count.

Evidence lives in build/port-map-theme: algorithm-original.log,
modifier-original.log, template-removal-original.log, prerequisite-variants.json,
removal-variants.json and their referenced logs. These deliberate repairs were
committed/pushed before the complete corrected-C baseline.

## Coverage

**3,392 cases / 19 complete captures**, plus 74 focused prerequisite contracts.
Snapshots include all allocated record bytes and disposal, token/counter/file
state, random draws and unchanged player records. The fixture uses actual startup
string tables, sparse real player iteration, real file handles and scoped
allocation/clock observers. Every capture hash is mandatory; there is no bypass.

| Capture | Cases |
| --- | ---: |
| algorithms | 728 |
| choices | 226 |
| conditions | 688 |
| decor-copies | 192 |
| decor-properties | 60 |
| decor-sets | 219 |
| decorations | 78 |
| equipment-boundaries | 26 |
| equipment-sets | 256 |
| exits-prefabs | 68 |
| filtered-tokens | 87 |
| foreach | 24 |
| frequency-validation | 192 |
| full-files | 52 |
| operators | 12 |
| raw-tokens | 277 |
| skips | 90 |
| spell-sets | 45 |
| template-removals | 72 |

The 52 full-file cases use actual keyed theme files, including missing and partial
files, complete section combinations and final settings validation. The supplied
assets have no .thm files or AreaMap.lib: this establishes synthetic theme coverage,
not real theme/prefab asset integration. Ordinary asset-backed gameplay is tested.

## Qualification

Complete C captures match independent default/server/highres runs, alongside
existing map checks: **57.752 / 63.662 / 65.888s** wall. The locked C smoke passed
in **14.074s**. First native captures match all hashes in **13.353s**.

Accumulated **67,703 cases / 969 capture groups**, plus room/painting/hallway/theme
contracts, pass default/server/highres in **365.450 / 360.441 / 299.877s** wall.
After restoring original error logging, focused themes pass all variants in
**98.272 / 96.775 / 24.057s** wall, including compilation.

All three final production binaries build and verify ELF32/i386/SSE2/CGO; all 31
retired symbols and test adapters are absent, and both required bridges remain.
Build times: **69.317 / 9.749 / 69.820s** for normal/highres/server. Asset-backed
full-suite failures match the existing baseline exactly: **1,553 entries;
15 pass / 3 fail / 32 skip packages**, no additions/removals. Fresh unchanged
repeat-a gameplay passes in **36.461s** under Xvfb/null audio, override disabled.

Evidence: build/port-map-theme/qualification.json, complete-captures.json,
complete-variants.json, native-comparison.json, variants.json,
final-theme-variants.json, binary-verification.json and full-suite-comparison.json;
gameplay is build/baseline/runs/map-theme-port/result.json. No active build/test.
Do not reapply native-*.go.stage drafts or the one-time finalize-reporting.py.

## Recoverable local evidence

Capture JSON is losslessly gzipped with SHA-256 manifests:
complete-capture-archive.json, native-first-capture-archive.json,
removal-capture-archive.json and stream-capture-archive.json. Source tests and
mandatory hashes are committed; raw local snapshots are diagnostic evidence.

Older qualified binaries are losslessly compressed under their original bin
directories. Their metadata is in previous-binary-archive.json (population/hallways),
spell-binary-archive.json (effects/lifecycle/sustained), and
room-painting-binary-archive.json, all under build/port-map-theme. Restore one with
`python3 build/port-map-theme/restore-qualified-binary.py ABSOLUTE_BINARY_PATH MANIFEST_PATH`.
The script verifies SHA-256 and restores recorded permissions/timestamps.

Completed population/hallway/room/painting/sustained-spell/spell-effects/lifecycle
runs each had **1,654 identical asset files / 556,358,986 bytes** deduplicated
against canonical extracted assets. Changed/new files, logs and recordings remain.
Each run's deduplicated-assets.json has hashes and its restore command:
`python3 build/baseline/deduplicate-run-assets.py --restore RUN_NAME`.
The supplied archive is untouched. Fresh gameplay always uses fresh asset copies.

## Next batch

Seven map-growth/door routines, **1,082 C section lines**, 4D4790..4D5D20 in
GAME3_2.c; three external entries and four private helpers. Read-only inventory,
caller audit, actual table values and plan are in build/port-map-growth. No source
changes for that batch yet. Continue after this qualified conversion is pushed.
