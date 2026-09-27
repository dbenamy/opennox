# Unused generated dispatch signatures

Status: original baseline accepted; draft reviewed but not installed.
Original source is `62873a54`, the qualified remaining-fixture conversion.

The generator emits 76 generic C callback signatures, but only 19 have Go callers
in the repository (one is fixture-only). An AST import/selector audit and a literal
search through every tracked Go, assembly, C and header file found the other 57
names only in their own definitions. Selected dependency metadata lists only
project packages as importers. This is unused source cleanup, not a claim that
57 live callback routes have been translated.

The draft limits the existing generator to those 19 signatures and regenerates
its output. Every retained C definition is byte-identical; every retained Go
wrapper is AST-identical. No caller, owner, identity, fallback behavior, ABI
layout, scalar conversion or external binding changes. All ten live legacy C
exports remain. Expected embedded production bodies: 77→20 (19 generic + 1 specialized).
Cgo package/file counts, header counts and standalone C LOC do not change.

Reuse the preceding exact-source full-default and focused profile qualification.
The existing `TestLegacyCallbackAdapters` contract also passed twice per original
profile, checking callback choice, pointer arguments, signed return boundaries
and object callback invocation. After conversion, run it in all three profiles,
the full default root corpus, focused server/high-resolution roots, safe/static,
all production/ABI builds, known-suite comparison, final fresh save/load and
asset hashes. A separate GUI preview is unnecessary for deleting unused generic
signatures with identical retained bodies; the final gameplay gate remains.

Primary owns this batch while Luna is quota-limited. See
[baseline](unused-dispatch-signatures-baseline.json),
[manifest](unused-dispatch-signatures-batch.json) and
[selection](unused-dispatch-signatures-tests.txt).
Local audit/draft: `build/port-unused-dispatch-signatures/`.
