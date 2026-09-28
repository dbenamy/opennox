# Retire safe-only C memory bridges

## Scope and review

Remove eleven thin safe-only C export wrappers and their unused libc preprocessor
redirects from `legacy/cgo_safe.go`. Six calls in the safe memory fixture instead
invoke the exact existing `alloc` Go routines. The other five wrappers have no
remaining callers after redirects are removed. Preserve sanitizer compile/link
flags, `-O0`, the safe-mode constant, runtime memory checks, tracker ownership,
allocation failure behavior and all frozen assertions/captures.

A whole-source symbol audit finds the exports only in their definitions, the six
fixture callers and historical callgraph-tool exclusions. Project headers have no
standard memory/allocation calls. Remaining legacy C selectors do not call those
standard helpers; allocation observer bodies use `__real_calloc`/`__real_free`.
The low-level allocator is a separate cgo package, unaffected by legacy's defines.
No observer, allocator implementation or external library changes in this chunk.

On the qualified 386 target, the fixture's bounded sizes preserve the former
unsigned-int conversions; returned lengths/comparison values keep their int32
conversion. Destination pointers are compared directly. The AST review finds one
changed fixture function, eleven removed wrappers and unchanged safe initialization.
The proposed result is zero safe-only exports and six fixture C-import files,
with normal-profile production cgo still 4 client/highres and 5 server.
The safe flags file still imports C intentionally. Standalone C stays zero;
headers remain 157 files / 2,731 physical lines.

## Original qualification

The original path passed nine safe roots, including `TestSafeMemoryBridges` and
its frozen capture. Eight allocation/resource/thread-scope/aligned-input owners
passed in each normal profile, followed by three focused repeats per profile.
Five direct allocator contracts and five private string/clock contracts passed
in both default and safe. Original root binaries are reused from the numeric
chunk only after exact source, supplemental-source, environment and binary hash
checks; these test selections execute freshly. No original source correction or
new golden is required.

Normal production source is unchanged. Reuse the preceding numeric production
baseline before conversion, then rerun native safe contracts, library contracts,
all selected profiles, static checks, safe build, production builds/ABI, exact
known-suite outcomes, fresh save/load and original asset hashes.

The selector has the distinct `safe-bridge-retirement` prefix; historical
`safe-bridges` manifests/selectors remain intact. Existing assertions cover byte
mutations/guards, comparison signs, exact string terminators, destination identity,
allocation/reallocation ownership, failure cleanup and observer thread scope.
Do not expand this chunk into allocator redesign.

Evidence: [original baseline](safe-bridge-retirement-baseline.json).
Working evidence and the uninstalled draft: `build/port-safe-bridges/`.
Primary handles this batch; Luna quota is unavailable and no substitute is used.
