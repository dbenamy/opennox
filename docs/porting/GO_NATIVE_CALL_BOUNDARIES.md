# Native Go calls at internal adapter boundaries

Status: original baseline accepted from fully qualified `78ff20f9`; implementation
is being drafted and is not installed. See the
[baseline](go-native-call-boundaries-baseline.json) and
[qualification manifest](go-native-call-boundaries-batch.json).

## Scope and baseline

The preceding primitive-type batch left 31 candidate files whose 99 remaining
C selectors serve Go adapter arguments (93 occurrences) or Go-defined scalar
fields (six writes). No selector calls a C function directly. This connected batch
moves those callers to native Go helpers and normalizes the necessary private
signatures/fields. Include connected definitions and callers outside those files.
Keep actual C exports, their signatures and external native-library bindings.

All four preceding qualification phases have exactly the current source hashes:
reuse their frozen owner captures, complete root suites, safe build and production
results as the original baseline. No expected values are regenerated. Existing
C export contracts continue to exercise wrappers where native cores are extracted.
The recently proven inventory identity regression remains in every root profile.

## Boundaries to preserve

- Integer widths, signedness, explicit narrowing and float32-to-float64 widening.
  The previous target probe establishes primitive representation on 386/SSE2.
- Raw object/session address words retain 32-bit conversion before pointer
  reinterpretation. No allocation ownership or lifetime changes.
- Nil `server.Obj` interfaces remain nil, including object/owner creation arguments
  and projectile exclusions. Temporary damage already passes literal nil.
- String adapters retain interned-pointer lifetime and truncation at the first NUL.
- Window hit testing retains global coordinates and inclusive edges; the existing
  nil-window behavior accepts only the origin. Do not substitute a different API.
- Inventory display and failure-message helpers retain their entire bodies,
  rendering effects, integer returns and exact three-byte encoding.
- Checksum/validation cores retain nil handling, unsigned byte counts, complete-word
  truncation, bounded chunking, XOR results and signed validation-ID rejection.
- Four browser fields and one existing raw-storage field use identical Go primitive
  widths/alignment. Deadline addition wraps in uint32 before widening to uint64.
  Keep unmanaged storage and all raw/typed aliases; migrate every field writer.

Five proposed private signature changes have only identified Go callers, including
test bridges. Primary review checks the complete source references, not just the
initial candidate set. The original C exports are retained around exact extracted
native bodies, so export removal is not a goal of this batch.

## Work and qualification

Luna owns an uninstalled 17-caller draft; primary owns the remaining callers,
helper/field changes, integration and acceptance. Exact original hashes and old/new
edit manifests permit independent reconstruction. Keep source frozen during tests.

After integration, repeat all seven storage captures, the three complete root
suites with exact names/skips, safe build/static checks, three fresh production/ABI
checks, the exact known full-suite comparison and headless save/load/resume. Verify
original assets and measure selected production cgo files rather than assuming all
31 imports disappear. Standalone C remains zero.

Local drafts, source fingerprints and review evidence live under
`build/port-go-native-call-boundaries/`. Candidate planning and five-helper caller
references are retained under `build/port-go-primitive-interfaces/next-*`.
