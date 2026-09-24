# Native Go calls at internal adapter boundaries

Status: qualified on Linux 386/SSE2. **31 production C imports removed;
selected project cgo files fall from 302 to 271.** All 50 changed source files
match the independently reconstructed draft. C exports remain 1,179; embedded
callback bodies remain 79; standalone production/test C remains zero.
External native-library bindings are unchanged.

Baseline: fully qualified `78ff20f9`, frozen at `3e892775`. See the
[baseline](go-native-call-boundaries-baseline.json),
[qualification](go-native-call-boundaries-qualification.json),
[manifest](go-native-call-boundaries-batch.json), and
[dependency inventory](go-native-call-boundaries-inventory-after.json).

## Final qualification

- All seven frozen storage captures match: raw/scalar in default, server and
  highres; scalar separately in safe. Exact test sets pass without skips.
- Complete root suites: default/highres each 2,426 pass; server 2,415 pass.
  Each has only the expected prerequisite-probe skip; discovered, started and
  completed name sets match the qualified baseline, including the identity regression.
- Safe build/static checks and three fresh production builds/ABI checks pass.
- Headless character creation, save/load and resume pass against the frozen reference.
- Full-suite outcomes exactly match the known baseline: 304 failure events,
  17 passing, two failing and 32 skipped packages.
- Identical source fingerprints throughout; all 1,654 original asset hashes match.
  No root assertion, frozen expectation or production behavior correction was needed.

## Scope and baseline

The preceding primitive-type batch left 31 candidate files whose 99 remaining
C selectors serve Go adapter arguments (93 occurrences) or Go-defined scalar
fields (six writes). No selector calls a C function directly. This connected batch
moves those callers to native Go helpers and normalizes the necessary private
signatures/fields. Include connected definitions and callers outside those files.
Keep actual C exports, their signatures and external native-library bindings.

Before conversion, all four preceding qualification phases matched the accepted
baseline source hashes. Their frozen owner captures, complete root suites, safe
build and production results supplied the original baseline. No expected values
were regenerated. Existing
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

The five private signature changes have only identified Go callers, including
test bridges. Primary review checked the complete source references, not just the
initial candidate set. The original C exports are retained around exact extracted
native bodies, so export removal is not a goal of this batch.

## Work and qualification

Luna supplied a 17-caller draft; primary handled the remaining callers,
helper/field changes, integration and acceptance. Exact original hashes and old/new
edit manifests permit independent reconstruction. Source stayed frozen during tests.

Qualification repeated all seven storage captures, the three complete root
suites with exact names/skips, safe build/static checks, three fresh production/ABI
checks, the exact known full-suite comparison and headless save/load/resume.
Original assets and selected production cgo counts were independently verified.
Standalone C remains zero.

Local drafts, source fingerprints and review evidence live under
`build/port-go-native-call-boundaries/`. Candidate planning and five-helper caller
references are retained under `build/port-go-primitive-interfaces/next-*`.

## Integration review

Primary reconstructed all 50 changed files byte-for-byte from the two reviewed
manifests. The 31 candidate C imports disappear, with no C export signature
changes. Five private signatures and all identified callers change together;
four browser fields and one raw-storage field retain width, alignment and aliases.

Review preserved evaluation order around the replaceable `GetServer` hook:
area-damage callers use an exact native helper so argument reads precede that
hook; string interning occurs before the receiver lookup and string reading
remains after it. Both nil creation arguments and projectile exclusions keep nil
interfaces. Two shortened helper names in Luna's caller draft were corrected to
the existing names during integration; definitions were never renamed.

The existing checksum contracts cover every short length/alignment, seeded data,
chunk boundaries, huge null lengths and rejected IDs with unreadable storage.
The window matrix retains nil windows and boundary coordinates. No new expected
values or root fixtures were needed for these exact extractions. The first 386
compiler preflight passed without fixes. All qualification gates
passed on the installed source. Draft reconstruction supplements compilation and
owner contracts.

## Delegation and follow-on review

The bounded Luna caller draft was useful after primary review of every conversion
and exact reconstruction. Primary corrected two helper names and preserved the
replaceable server-hook evaluation order before installation. No subscription
savings estimate is claimed. A follow-on field audit found 26 C scalar fields in
two Go-owned aggregates; field-only normalization would remove no C imports, so
combine it with a connected dependency-removal batch if selected.

Seven obsolete scalar-storage executables were removed after independent source,
replacement, hash and host-use checks, reclaiming 395,452,416 allocated bytes.
Recovery uses `e64ff24e`; replacement source matches `78ff20f9`. The helper initially
compared replacement maps with the old revision; primary caught the mismatch and
required the corrected proof before cleanup. Retain the final proposal and
`cleanup-scalar-{approved.json,deleted.jsonl}` under this batch directory.
