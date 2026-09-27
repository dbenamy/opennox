# Libc helper boundaries

## Scope and baseline

This batch replaces libc use in entry text classification, theme numeric parsing
and clock reads, browser address parsing, and shop field-guide string copying.
Production baseline: `6e9681f2`; original fixture baseline: `e1926909`.
The conversion is fully qualified. Local artifacts are under
`build/port-libc-boundaries/`.

Four production and two fixture cgo imports are retired. Selected production cgo
files are **9 client / 10 server**, down from 13/14. Standalone production/test C
remains **0/0**, selected legacy exports remain **0**, embedded production C bodies
remain **20**, and headers remain **157 files / 2,731 lines**. Existing
C-string allocation, trigonometry, allocator implementation and callback fallback
work remain separate. External native libraries remain in scope only as retained
dependencies, as specified in PORT.md.

## Behavioral evidence

- Independent 386 libc probes captured 58 numeric inputs twice, preserving exact
  int32 results and float64 bits. Cases include saturation, signed zero, malformed
  prefixes/exponents, hexadecimal numbers, infinities, NaN payloads, subnormals,
  overflow/underflow, whitespace and embedded NUL.
- Two original address captures agree on 5,615 explicit, boundary and deterministic
  generated cases. Contracts exercise the actual selected-browser-record entrypoint
  and check guards and record immutability. Legacy short/octal/hex forms and suffix
  acceptance after ASCII whitespace are intentional compatibility behavior.
- Existing entry contracts cover every uint16 input for both predicates: 131,072
  classification results. Native ASCII predicates independently reproduce the
  original C-locale hash. No locale-changing initialization was found in engine
  source or the inspected OpenNox libraries.
- New clock contracts cover five 32-bit epochs, repeated activation and idempotent
  restoration. Existing theme file captures exercise clock consumption by the
  real owner. Allocation observation retains thread-local exclusion and OS-thread
  pinning; the native fixture clock hook is global and used serially.
- Existing shop-stock loading captures include field-guide creation and copied
  payloads. Native `alloc.Strcpy` includes the terminal NUL and is already the
  safe-build implementation behind C strcpy. Object/type ownership stays intact.

Frozen numeric/address expectations live in `src/legacy/testdata/`; original
contracts are in `src/legacy/libc_boundaries_porttest_test.go`. C probes and their
raw repetitions stay in ignored build artifacts, rather than retaining C
implementations solely for testing.

## Review and decisions

The existing resource float parser differs from atof on incomplete negative
hexadecimal prefixes: libc preserves negative zero. The correction is
local to the theme boundary, leaving resource-parser behavior unchanged.

Browser private C scalar parameters/returns become exact-width native integers.
The whole-source caller scan found no remaining C callers or preamble/header uses;
fixture callers migrate with them. The centralized wchar alias and formatter
CString allocation/free ownership remain unchanged. Invalid addresses and valid
broadcast addresses retain their shared all-ones return value.

The AST review counts 32 changed, 46 unchanged and one added function (one counted
change is an attached comment). Regression selection follows helper references,
expands complete assertion files, and explicitly includes theme, browser, entry,
shop and shared allocation-observer owners: 813 candidate root names before
profile filtering. Private legacy-package contracts run separately.

Luna remains unavailable due to its reported usage limit. Primary completed the
bounded draft, capture work and caller review locally; no substitute model ran.

## Qualification

Original-path four-root preflight passed. Default/server/high-resolution root
selections pass 813/811/813 tests twice each, with exact names and no skips. The
four private legacy contracts pass twice each in default, server, high-resolution
and safe profiles. Source fingerprints match throughout. Production source is
identical to the preceding qualified revision; its production evidence is reused
only for this original baseline. Production gates were rerun after conversion.
See [baseline](libc-boundaries-baseline.json) and [captures](libc-boundaries-captures.json).
Converted four-root preflight and static check pass. The first compile caught a
missing explicit unsafe.Pointer cast for the shop name helper; no tests ran in
that failed compile. The corrected build passes, and a fresh headless character
creation/save/load/resume preview matches the reference. All converted root
selections pass 813/811/813 with exact original names and no skips. Private contracts pass twice in all four profiles. The safe production
build/static check passes. All three production builds/ABI checks and the exact
known-suite failure/package comparison pass. The final fresh save/load/resume run also passes, and all 1,654
original asset hashes are unchanged. Accepted gates share the reviewed source
fingerprints. See [qualification](libc-boundaries-qualification.json) and
[inventory](libc-boundaries-inventory-after.json). The last complete default
port-test corpus remains the preceding `6e9681f2` qualification; this bounded
batch uses audited affected selections.

## Artifact headroom

Seven superseded scalar-batch binaries were deleted after verifying their source
fingerprints against `1b80dbe1`, replacements against `6e9681f2`, artifact hashes,
file identities and host process use. Reclaimed 387,121,152 allocated bytes.
Current replacements, original captures, source and assets remain. Local plan and
journal: `build/port-libc-boundaries/scalar-cleanup-{approved.json,deleted.jsonl}`.

The completed preview also had 1,654 hash-identical asset copies removed
(559,874,048 allocated bytes). Saves and results remain; restore with
`build/baseline/runs/libc-boundaries-preview-save/deduplicated-assets.json`.

After all test jobs exited, 22 inactive Linux386 compiler-cache archives were
removed following hash/stat and host compiler/open-file checks (589,103,104
allocated bytes). Their cache misses rebuild normally. Current binaries, source
and caches used within six hours were retained. Records are under
`build/port-libc-boundaries/cache-headroom-*`.

The final scenario had another 1,654 verified asset copies removed (559,910,912
allocated bytes), preserving originals and save/comparison evidence. Restore via
`build/baseline/runs/libc-boundaries-save/deduplicated-assets.json`. Completed
installation, qualification and cleanup scripts are consumed.
