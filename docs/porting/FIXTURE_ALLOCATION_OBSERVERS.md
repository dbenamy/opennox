# Native fixture allocation observers and header retirement

## Scope

Replace four connected C observer fixtures with a Linux porttest-only registry
at the centralized engine allocation boundary. Retire their two C exports and
linker wrapping of calloc/free. Normal production dispatch forwards directly to
the existing libc primitives; raw malloc/realloc and allocator implementation are
unchanged. Aligned room/painting buffers remain for the allocator phase.

The abort adapter's only C call is standard-library abort. Remove its obsolete
defs.h/GAME4.h includes alongside the observers' last project-header include.
A whole-source include audit then finds only standard-library headers in repository
C imports, including Windows-only socket bindings. Retire all 157 project headers
(2,731 physical lines). They contain declarations/types, no startup constructors,
static function bodies or abort override. Lexer tests use literal include text;
historical capture tooling already recovers its retired source from Git.

The proposed result is two fixture C-import files, zero fixture C exports and
zero project C headers. Normal production imports stay 4 client/highres and 5
server; the safe flags file remains. Standalone C and production C exports/bodies
stay zero. External SDL2/OpenGL/OpenAL dependencies are preserved.

## Ownership and thread review

Register only while pinned to an OS thread, and restore before unpinning. A registry
keyed by Linux thread ID gives each thread its own observer. Theme, grid and
resource observation share that thread's state but retain independent activation
and cleanup. Repeated stop is inert. Private grid statistics are read within their
observation interval; inactive theme state reads zero.

Preserve the shape-dependent failure countdown, 130-slot grid limit, nil suppression,
32-bit size multiplication and release order: resource record, grid record, theme
callback, then raw free. Allocation callbacks precede safe tracker registration;
release callbacks follow tracker removal. Safe calloc's existing nil marker is
preserved. Public allocation API signatures stay unchanged.

Observation now covers engine calls through RawCalloc/RawFree. Linker wrapping also
saw incidental runtime/external-library calls; those are outside engine ownership
contracts. The audited engine consumers all cross the centralized boundary, and
frozen real-owner outputs must remain unchanged. Do not claim a whole-process
libc allocation monitor. Production builds contain no registry or observer dispatch.

## Independent original contract

The existing thread test checked only an activation flag. Add
`TestAllocationObserverThreadIsolation` before conversion: while failure injection
is active on the pinned owner thread, a C pthread allocates/releases through the
actual wrappers. That allocation must succeed without consuming the owner's failure
or changing its counters. The owner's next allocation fails, the following one
succeeds and releases, and tracking returns to its prior count. Preserve then remove
only the fixture-created nil marker in safe mode.

This new contract passed twice in fresh processes in safe/default/server/highres,
with three independent thread lifecycles per invocation. The native equivalent uses
another pinned Go thread, so it exercises the same isolation at the engine boundary.

## Baseline and planned qualification

The active-observer graph reaches 76 roots across theme/growth/orchestration,
map sections, prefab disposal, resource teardown and grid ownership. Add allocation,
CString and aligned-buffer boundary contracts: 80 roots per normal profile.
All 80 passed, plus eight focused repeats per profile and eleven safe roots.
Five allocator and five string/clock library contracts passed in normal and safe.
The baseline adds only the original contract; production remains byte-for-byte
source-identical to the preceding safe-bridge qualification.

Native qualification starts safe-first, then direct library contracts, focused
preflight, matching owners/repeats and the complete compiled root corpus in every
profile. Only the opt-in population diagnostic may skip; its four cases already
have regression coverage. This broad sweep is warranted because test allocation
dispatch is shared and the last project headers retire. Finish with static checks,
safe build, production builds/ABI, exact known-suite outcomes, fresh save/load and
original asset hashes. Verify C observers/export symbols are absent and sanitizer
initialization remains present. Preserve frozen root assertions/captures.

Evidence: [baseline](fixture-allocation-observers-baseline.json).
Working evidence and the uninstalled native draft:
`build/port-fixture-allocation-observers/`. The original contract is installed;
no native observer/header deletion is installed yet. Completed scripts are single-use.
Primary handles this batch; Luna quota remains unavailable.
