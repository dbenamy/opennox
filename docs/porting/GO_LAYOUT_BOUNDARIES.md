# Native geometry and shared-state boundaries

Status: original baseline accepted from fully qualified `4214ea8f`; no conversion
is installed. See the [baseline](go-layout-boundaries-baseline.json) and
[qualification manifest](go-layout-boundaries-batch.json).

## Connected scope

The initial inventory identifies 23 files with 224 C type selectors and no C
export or build directives. Twenty contain active owners; three contain only
unused private wrappers. Include the remaining 26 C scalar fields in `browserUI`
and `legacyGlobalStorage`, their readers/writers, and callers of changed private
interfaces. The initial file count is not the integration scope or a promised
C-import reduction: shared adapters have callers outside those files.

Retire 132 private orphan helpers where complete reference and body review prove
they are unused. The primary source scan and Luna's independent scan of 4,216
tracked text files find only their definitions in runtime/build source; additional
references are historical documentation. None has a C export or linkname
directive. Preserve all other declarations and native owners. Inspect affected
preambles and imports before deleting candidate-only files.

The nine helpers Luna flagged for pointer/layout operations are ordinary unused
function bodies, not initialization or registration. Primary review found no
reason to preserve those dead wrappers solely because their bodies use pointers.
Recheck every deletion against the frozen source before integration.

Keep actual C exports, callback identities, shared raw fallbacks, external native
bindings and existing memory ownership. Do not introduce C-name aliases merely
to reduce the import count. Prefer exact native owners and meaningful private
Go signatures, updating all callers at the remaining C boundary.

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

Use one Luna helper for a bounded caller draft; primary owns shared signatures,
state fields, orphan retirement, integration and acceptance. Require original
hashes and exact edit manifests, then independently reconstruct the installation.
Keep source frozen during builds/tests and preserve every frozen expectation.

Run focused owner checks during implementation. The completed connected batch
changes widely shared interfaces and state, so repeat all seven storage captures,
all three complete root suites, safe build/static checks, three production/ABI
checks, the exact known-suite comparison and headless save/load/resume. Check
exact test-name sets and skips, source identity, assets and dependency counts.
Standalone C remains zero; measure other counts after qualification.

Luna's planning report needed corrections to point-value arguments, raw waypoint
returns and the distinction between initial candidates and complete caller scope.
Named tests exist, but operation-selector coverage still needs primary review;
an adjacent test name does not establish coverage of a changed path. These are
planning corrections, not changes to engine behavior.

Local planning evidence is under `build/port-go-native-call-boundaries/next-*`;
the new batch's baseline, drafts and runs belong under
`build/port-go-layout-boundaries/`.
