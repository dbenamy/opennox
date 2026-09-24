# Go types for internal primitive interfaces

Status: original baseline accepted from qualified source `e64ff24e`; conversion
not installed. See [baseline evidence](go-primitive-interfaces-baseline.json)
and [qualification manifest](go-primitive-interfaces-batch.json).

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
original results; there are no new fixtures or changed expectations. The repeated
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
with Go ASTs and target type probes, owns scope/baseline/acceptance, and is reviewing
caller boundaries. Luna is tracing the 20 local-scalar candidate files to their
native callees; unresolved mappings stay explicit rather than being guessed.
A syntax inventory does not establish C-header reachability or type correctness.

Local work: `build/port-go-primitive-interfaces/`. `baseline-source.json` freezes
all relevant source fingerprints; `ast-inventory.json`, `signature-names.json`
and `signature-references.json` record the initial syntax/caller boundary.
`primitive-types{,-repeat}.json` record the original target type probe. None of
these drafts authorizes source installation before review.
