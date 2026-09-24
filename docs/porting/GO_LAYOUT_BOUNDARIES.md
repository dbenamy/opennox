# Native geometry and shared-state boundaries

Status: qualified on Linux 386/SSE2. **24 production C imports removed;
selected project cgo files fall from 271 to 247.** The 74-file batch converts
26 scalar fields, retires 145 private wrappers and deletes ten empty files.
Actual C exports remain 1,179; embedded callback bodies remain 79; standalone
production/test C remains zero. External native-library bindings are unchanged.

Baseline: fully qualified `4214ea8f`, recorded at `80af9879`. See the
[baseline](go-layout-boundaries-baseline.json),
[qualification](go-layout-boundaries-qualification.json),
[manifest](go-layout-boundaries-batch.json) and
[dependency inventory](go-layout-boundaries-inventory-after.json).

## Final qualification

- All seven frozen storage captures match: raw/scalar in default/server/highres,
  scalar separately in safe. Exact test sets pass without skips.
- Complete roots: default/highres each 2,426 pass; server 2,415 pass. Each has
  only the expected prerequisite-probe skip; exact root-name sets match baseline.
- Safe build/static checks and three fresh production/ABI checks pass.
- Headless character creation, save/load and resume match the frozen reference.
- Known-suite outcomes match exactly: 304 failure events, 17 passing, two failing
  and 32 skipped packages.
- Identical 3,052 source fingerprints throughout; all 1,654 original assets match.
  No root assertion or frozen expectation changed. Two existing fixture bridges
  received native parameter types; raw-storage writes retain the same word bits.

## Connected scope

The initial inventory contained 23 files with 224 C type selectors and no C
export or build directives. Integration included their connected callers, all
26 remaining scalar fields in `browserUI` and `legacyGlobalStorage`, and private
wrapper retirement. The final scope is 74 files.

Whole-source reference and body review proved 137 wrappers unused. The primary
scan and Luna's independent scan covered 4,216 tracked text files, distinguishing
historical documentation and translation-rule strings from runtime references.
Eight further wrappers retired after their active callers moved together.
None of the retired functions was an initializer, C export or linkname target.

Actual C exports, callback identities, shared raw fallbacks, external native
bindings and memory ownership remain unchanged. Native owners and private Go
signatures replace the adapters without introducing C-name aliases.

## Compatibility and baseline

All four preceding qualification phases match the baseline's 3,062 source
fingerprints. Reuse their frozen captures and qualified production behavior.
The target probe compares the actual headers with Go layouts: integer/float
two-word points are 8 bytes aligned to 4, four-word records are 16 bytes aligned
to 4, and bool is 1 byte aligned to 1. The second point field is at offset 4;
cgo's bool has Go bool kind. Probe results are embedded in the baseline record.

Preserve signed narrowing, float32 rounding/widening, point-copy timing, null
handling, raw address words, output-buffer aliases and return conventions.
Preserve evaluation order around replaceable hooks and function variables.
`ApplyForce` takes a point value, read after receiver lookup in its existing
adapter. Waypoint owners already return raw `uint32` address words; they need no
extra pointer-to-word conversion. String paths retain NUL termination, nil
distinctions, interned lifetime and case-insensitive `NONE` handling where present.

State fields retain size, alignment, wraparound and existing pointer-like word
aliases. Updating direct writes alone is insufficient: include readers, indirect
column writes, typed aliases and fixture bridges. Keep unmanaged backing storage.

## Implementation and qualification

One Luna helper supplied a bounded caller draft; primary owned shared signatures,
state fields, orphan retirement, integration and acceptance. Original hashes and
exact edit manifests allowed independent reconstruction of every installed file.
Source remained frozen during all builds and tests.

Luna's planning report needed corrections to point-value arguments, raw waypoint
returns and the distinction between initial candidates and complete caller scope.
Primary traced fixture operations to the changed owners before acceptance.
These were planning corrections, not changes to engine behavior.

Local baseline, drafts and qualification artifacts are under
`build/port-go-layout-boundaries/`.

Integration review reconstructed all 74 changed files exactly, including ten
files left empty by wrapper retirement. The connected batch removes 24 C imports
and 145 private helpers: the initial 132, five additional orphans, and eight
wrappers whose active callers were migrated together. All actual C export
signatures remain unchanged. The 24 removed preambles contain only declarations
and includes; their 32-header local closure has no initialization/registration
body. Existing scalar storage and root expectations remain unchanged.

Primary review caught and corrected two stale struct-field accesses after Luna
converted a direction scratch value to an array. This was a draft correction
before installation, not a fixture or behavior change.

Storage and full root/production qualification all pass on the same frozen source.

## Compatibility decisions and review

The five additional orphan functions have only historical comment or translation
rule references; those strings are not runtime calls. Eight additional live
wrappers retired only after their complete legacy caller set moved to the native
owners. Same-named root-package owners and public function variables remain.
No remaining legacy source identifier references any of the 145 retired functions.

The target layout probe, full reader/writer audit and storage captures establish
same-width fields and aliases. Pointer-like words remain words; unmanaged storage
and existing memory ownership are unchanged. Mana deltas retain signed int16
extension; event codes retain signed int32 while raw arguments retain uint32 bits.
Strings keep the original NUL reader and interned lifetime. Point values retain
their copy timing relative to replaceable server hooks.

One Luna helper supplied 13 caller files and an independent review of the primary
mappings. Primary supplied the shared state/interfaces and wrapper retirements,
reviewed each owner mapping, reconstructed both drafts and performed acceptance.
The array-field draft error was corrected before installation. Reusable test
routes were checked explicitly; no claim of universal branch coverage is made.

Local reconstruction, preamble closure and owner coverage records are under
`build/port-go-layout-boundaries/`. Current source and qualification metadata are
committed; local binaries/drafts/logs are rebuildable and not backed up by Git.
