# Go types for internal primitive interfaces

Status: qualified on Linux 386/SSE2. **52 production files and three test files
no longer require cgo; selected project cgo files fall from 354 to 302.**
The engine change normalizes 103 private interfaces (102 original candidates plus
the locked-door string parameter) and their callers, with three small native
helper extractions. C exports remain 1,179; callback bodies remain 79; standalone
C remains zero. External native-library bindings are unchanged.

Baseline: qualified `e64ff24e`, recorded at `ff4e8440`. See the
[baseline](go-primitive-interfaces-baseline.json),
[qualification](go-primitive-interfaces-qualification.json), and
[dependency inventory](go-primitive-interfaces-inventory-after.json).

## Final qualification

- All seven frozen storage captures match: raw/scalar owners in normal profiles,
  scalar owners separately in safe. Exact expected test sets pass without skips.
- Default/highres each pass 2,426 root tests; server passes 2,415. Each has only
  the expected map-population diagnostic skip. Exact root names match the prior
  inventory plus the collision regression; no omitted tests or failure events.
- Safe build/static check and three fresh production binaries/ABI checks pass.
- The full suite matches the known baseline exactly: 304 failure events, with
  17 passing, two failing and 32 skipped packages. This is not an all-green claim.
- Headless character creation and explicit save/load/resume pass.
- All phases use identical source fingerprints. Exact reconstruction covers
  129 source files; the only root fixture change is the independently proven
  normalization fix and its regression. Frozen goldens remain unchanged.
- All 1,654 original asset sizes and hashes match after production qualification.

Final artifacts are under `build/port-go-primitive-interfaces/`:
`storage-final/`, `contracts-final/`, `safe/`, and `production/`.
Initial failed qualification is retained under `contracts/`, with logs archived
losslessly. The fixture proof below explains that failure and its correction.

## Boundary and scope

Normalize primitive types in 102 unexported Go function signatures and their Go
callers. The discovery set is 83 selected production files containing 465 numeric
C selectors. All are selected in default/server/highres; none has an active C
export directive or cgo build flag. The primary AST audit finds 255 direct calls
across `src/legacy/*.go`, including porttest helpers, and no function-value uses.
The candidate functions keep their names and behavior. Actual C exports and
external native-library interfaces keep their existing signatures.

Use fixed-width Go types with the same underlying kind, size and alignment on
Linux 386/SSE2. An isolated cgo probe repeated in fresh processes confirms all 12
mappings: signed `int`/`char`/`short`, unsigned integer variants and fixed-width
aliases, 32-bit `size_t`, and `float`/`double`. `char` is signed. This is a type
representation change in existing Go code, not new floating-point arithmetic.
Keep every explicit narrowing step, pointer alias, null value, allocation domain,
return convention and initializer. Do not turn raw address words into managed
pointers or change pointed-to storage. Existing C-structure aliases are outside
this scalar-type boundary.

Review Go callers with C-typed wrappers separately. Prefer the existing native
implementation where its argument/result conversions can be traced exactly;
otherwise retain explicit conversion at that surviving boundary. Likewise retain
casts for actual C fields and calls. Removing all 83 C imports is not promised:
measure the genuinely empty imports after integration. Never hide the dependency
behind a new central alias to a C type.

## Baseline and qualification

All production/test source fingerprints equal the just-qualified scalar batch
for storage, complete root sweeps, safe/static and production. Reuse those exact
original results as the engine baseline. The separately proven fixture correction
and new collision regression below leave all frozen expectations unchanged. The repeated
primitive-type probe supplements the existing behavior/layout contracts. Keep
its original results, not regenerated expectations chosen to fit the conversion.

Before installation, independently check whole-selector edit spans, function
bindings, all changed argument positions and unchanged signature parts. Go-only
return changes can require further conversions in callers; compiler diagnostics
are an integration aid, not proof of behavior. Review pointer conversions,
signedness and preserved truncation explicitly. Check full preambles/headers,
initializers and exports before removing a C import. Preserve the `ai_combat.go`
registration initializer; its registry rejects duplicate actions. Preserve all
compile-time layout assertions.

After conversion, repeat scalar/raw owner captures in three normal profiles and
scalar-only safe; complete all three root sweeps with exact name sets; run safe
build/static and fresh production/ABI/known-suite/headless save/load/resume.
Compare the same frozen captures and original assets. Use the existing sequential
build/two-prebuilt-test process limit. No source edits while a build/test reads
this checkout.

## Delegation and local recovery

Luna supplied the bounded candidate/group inventories. The primary verified them
with Go ASTs and target type probes, reviewed caller boundaries, and owns
scope/baseline/acceptance. Luna produced a 99-file mechanical draft and a bounded native-call
audit. Primary reconstruction checks all bytes against the explicit edit spans;
the independent AST inventory agrees on all 465 selectors and all 102 signatures.
Primary review found redundant retained C casts in callers outside the original
83 files and requested a separate correction manifest, including pointer casts.
The 29 string-pointer TODOs resolve through the simultaneous `internCStr` return
conversion; existing string-reader adapters preserve the original pointer and
NUL/null behavior. Two animation output pointers become `*int32` with their
locals; the remaining chat pointer uses the reviewed native call. No allocation
or fixture expectations change. Unresolved boundaries remain explicit.
A syntax inventory does not establish C-header reachability or type correctness.

Local work: `build/port-go-primitive-interfaces/`. `baseline-source.json` freezes
all relevant source fingerprints; `ast-inventory.json`, `signature-names.json`
and `signature-references.json` record the initial syntax/caller boundary.
`primitive-types{,-repeat}.json` record the original target type probe. None of
these drafts authorizes source installation before review.

## Additional reviewed boundaries

Extract the existing background/selected image setters and drawable unit-code
predicate into three small native Go helpers. Keep their C entry points and
signatures, forwarding with the original 32-bit conversions. This allows the
reviewed Go callers to avoid those C-typed adapters. Preserve nil returns,
image-handle words, code-range rejection, class bits and unchanged object bytes.
The existing 5,040-snapshot window matrix includes nil/live windows and nil/two
live image handles; the existing client-code ABI contract includes nulls, boundary
codes, walking/complement class bits and 2,048 seeded cases. No fixture or golden
change is needed. These extractions add two files to the initial candidate scope;
they do not change allocation ownership or add a new C implementation.

Primary also inspected generated cgo declarations: these 12 primitives are
ordinary numeric Go named types/aliases, with no not-in-heap annotations. The
45-header candidate closure has warnings, balanced packing and type assertions;
full preambles contain includes plus three prototypes. Keep all structure layout
and C boundary conversions. The reviewed conversion is installed for compiler preflight.

## Integration checkpoint

The separate caller correction pass accounts for all 255 references and fixes
65 redundant scalar C casts plus 34 matching pointer casts outside the initial
selector pass. Primary verifies source/draft hashes, exact spans and unchanged
cast operands against the original argument expressions before integration.
The first preflight includes 102 changed files, 83 production C-import removals
and three test C-import removals. These are draft counts until compilation and
qualification establish the final boundary. An export-bearing caller retains
its C import even though its last C selector disappeared. No ABI retirement or frozen-expectation changes are part of this batch.
The later normalization regression addresses a proven fixture defect.

The first compile reported 160 fixed-width mismatches at surviving C boundaries.
Primary mapped each diagnostic to an exact Go AST expression, reviewed the target
primitive against the repeated type probe, and inserted 154 explicit conversions
in 56 files. This restores 31 original C imports/preambles. Pointer conversions
are between pointers to identical underlying numeric types; no pointer arithmetic,
allocation, storage or ownership changes. Keep the explicit narrowing before the
boundary even when it becomes a redundant same-width conversion.

Six diagnostics resolve through reviewed native boundaries: two score/death calls
use the exact existing gameplay helpers, and two drawable-code calls use the new
helper. The two locked-door string calls already receive `*int8` from `internCStr`;
normalize that one additional private helper's parameter and cast only at its
existing C export/string reader. Whole-source search finds those two callers and
one unchanged C export adapter. Its null checks, length checks, zero padding and
message bytes are unchanged. This avoids introducing cgo into previously pure-Go
callers. The second preflight has 128 changed files, 52 production C-import removals
and three test removals; no new C imports relative to the original source.

The second default porttest compile passes. Primary exact reconstruction of all
128 changed source files matches the reviewed draft, 154 C-boundary conversions
and the recorded native fixes. Root assertions and golden files are unchanged.

Storage qualification passes on unchanged source: default 67.53 s, server 66.90 s,
highres 4.88 s and scalar-only safe 71.77 s. The phase verifies all seven frozen
capture hashes and the complete expected storage test sets, with no skips.
Initial server/highres root profiles pass (2,414/2,425 passes plus the expected
skip); default finishes with only `TestClientInventoryDisplayIdentification`
failing. All 60 differing leaves are nested cached-return identities:
`0xec000001` becomes `0xe8000000`; pixels, text and all other state match.

The fixture normalized the return twice. A canonical drawable token can equal a
real image handle in the 32-bit handle arena and become an image token on the
second normalization. The deterministic regression forces that map collision and
reproduces the exact mismatch on unchanged pre-conversion engine source
`e64ff24e`. Passing the original return to each snapshot fixes it; the regression
and all existing inventory-display tests then pass on that original engine, with
unchanged frozen captures. It also checks null/non-null text returns. This changes
only the fixture; no engine behavior or golden is adjusted.

The corrected fixture and regression are installed. Exact reconstruction
covers 129 source files, including the explicit two-line normalization correction
and added test. Final storage, complete root sweeps, safe and production gates
all pass on this source. The new root is explicitly added to each historical name
set; historical test inventories remain unchanged. Proof logs and results live in
`build/port-go-primitive-interfaces/identify-proof/{unfixed,fixed}/`.
Final acceptance is recorded above.

## Local artifact maintenance

Luna identified three superseded 378-export root test executables. Primary review
caught that the qualification's recorded Git HEAD was the pre-conversion revision;
Luna's corrected recovery audit and primary independent `git cat-file` verification
match all 3,090 source/supplemental fingerprints to `f6f5ee4c`. The later successor
is `12bc387d`. After host process/open-file/mapping checks plus exact inode, size,
mtime and binary hash checks, primary removed those three rebuildable binaries,
reclaiming 204,169,216 allocated bytes. Assets, original behavior captures and current
binaries remain. The ignored approved plan and deletion journal preserve recovery.
The draft report's total typo was also corrected before cleanup. This was useful
bounded delegation, with primary checks still necessary; no cost-savings claim.

The completed initial root logs are losslessly gzip-compressed after host
open-file checks, reclaiming 306,057,005 bytes. Uncompressed and restored hashes
match; `contract-log-archive.json` records hashes, sizes and restore commands.

A second bounded Luna audit identified seven superseded Go-only-export
test/safe/production executables. Primary independently matched all five source
maps (3,062 unique files) to `12bc387d`, verified retained qualified replacements
and recorded hashes, and checked host processes/open files/mappings. Removing
those binaries reclaimed 395,526,144 allocated bytes. Commands, captures, logs
and recovery manifests remain; the approved plan/deletion journal is retained.
