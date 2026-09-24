# Centralized libc allocation

Status: original allocation-domain and owner baseline qualified; implementation
draft is not installed yet.

Move 49 allocation/free calls across 21 legacy owner files behind profile-specific
Go helpers. Move the tracked allocator's four libc calls behind the same raw
backend, leaving its tracking, zero-size and failure behavior unchanged. Normal
legacy calls remain untracked; the `safe` build tag selects tracked calls exactly
as the old C macros did. Runtime `NOX_SAFE=true` does not select those C macros.
The raw backend still uses libc; replacing storage ownership is later work.

## Contracts and qualification

The original bridge calls libc through the legacy package, including its safe
macro route. Thirty-six bounded allocation cases each exercise four grow/shrink
reallocations, calloc zeroing, preserved prefix bytes, tracker liveness/counts,
old-entry retirement when the pointer moves, and complete release. Captures omit
addresses and whether realloc moved; they include inputs and observed data.
Normal and safe captures each match across three separate processes and are
frozen before conversion. A fourth normal run with `NOX_SAFE=true` confirms the
compile-time domain distinction. Six CString cases independently verify copied
bytes (including embedded NUL/high bytes), termination, and actual CString/StrFree
tracking and cleanup in both profiles.

The targeted owner selection passes 111 roots without skips in each of default,
server and highres. Six safe-profile roots pass without skips: the two new domain
contracts, resource-graph and rejected-list release owners, existing memory
bridges and shop loading. Existing owner tests provide independent contracts;
this report does not claim every optional historical capture was regenerated or
compared. The preceding production qualification at `4e1e86a6` is reused only for
the test-only baseline; production source remains identical. Direct allocator
and memory-helper package tests also pass in all three profiles. After conversion,
repeat these gates, safe build/static checks, fresh production builds/ABI, exact
known-suite comparison and headless save/load.

The selector audit found 13 existing roots absent from the accumulated selection,
plus the two new contracts. All 15 are added to the accumulated pattern. Two of
those existing roots also lacked any focused selector: resource free graph and
server rejected-list free. The broader 111-root union catches the other omissions
through existing owner families. Baseline/native selected-name sets must match.

## Scope and review

Keep allocation sizes, zero-fill, pointer layouts and failure handling unchanged.
On 386, size_t and uintptr are both 32 bits; safe calloc's conversion to int matches
the previous C.uint-to-int route. Raw free/realloc must not acquire tracker
semantics. The known failed tracked-realloc bookkeeping issue remains outside
this mechanical batch. Existing string/input lifetime issues are also unchanged.
C.CString/C.CBytes and test-only direct libc fixture calls remain separate work.

Luna drafted the owner migration and reviewed the primary's bounded contracts.
Luna corrected an initial 48-site count to 49 while drafting; primary review caught
three unnecessary imports before integration and verified the 21 owner bodies
differ only in the reviewed call/type substitutions. The four tracked-backend
changes preserve all map/lock/failure statements. The import audit checks
selectors, export directives and C build directives before removing ten legacy
preambles. Moving alloc.go's C dependency into raw.go avoids adding a second cgo
file to that package. No delegation time/cost saving is claimed.

Artifacts: `build/port-raw-allocation/`; reviewed draft and audit history:
`build/port-go-memory/raw-*`. See [original qualification](raw-allocation-c-qualification.json).
