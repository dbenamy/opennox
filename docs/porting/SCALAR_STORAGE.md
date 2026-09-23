# Numeric global storage

Status: original-C baseline qualified; native storage draft not installed.
Production remains the qualified unused-bridge checkpoint `48b7f80b`: 518 physical
C lines in six files, no reference C. Only two tagged fixture/test files changed.
Production/ABI, known-suite and headless evidence is reused from that checkpoint.

## Scope and ownership

Move 343 numeric owners from `vardefs.c` into process-lifetime Go variables while
preserving their numeric C typedef, initializer bits and raw-word representation.
Go callers use those variables directly. This stage changes storage ownership;
it deliberately limits type/API changes. Go's local 1.26 runtime treats
linker-allocated globals as pinned; numeric fields do not acquire GC pointer slots.
Existing aliases must keep observing the same owner for the process lifetime.

The inventory found 407 scalar definitions: 346 with Go-only uses, 52 also used by
C fixture bodies, and nine unused definitions. Three of the 346 have pointer-typed
extern declarations and are deferred: `dword_5d4594_1090100`,
`dword_5d4594_1309720`, and `dword_5d4594_831236`. Pointer/array/struct owners and
C fixture-body users remain for later batches. A numeric definition alone does not
prove all declarations give it a numeric type.

The nine unused C definitions can be deleted; same-name native constants and
struct members are independent owners and remain. Whole-source literal review
found declarations, metadata, or separate native owners only. No selected scalar
references were found in assembly or C++ files. Review found no direct wider
access or address arithmetic for the selected scalar addresses. Three same-name
locals in `initMusic` refer to the package variable in their own initializer, before
the local enters scope; preserve those expressions.

## C baseline

One isolated legacy-package contract checks 343 actual C owners with 139 patterns
each: 47,677 cases. It checks declared width, initial definition bits, signed reads,
distinct/nonoverlapping storage, typed/raw alias visibility, unchanged neighboring
owners and restoration. Patterns include every walking bit and its complement,
32/64-bit boundaries and alternating bytes. GC runs while patterned raw numeric
values are installed. Expected widths and initializer bits are frozen from the
original C definitions; observations are also captured and hashed.

All cases pass in two independent default processes, server and highres, with no
skips; all captures match and static checks pass. Source snapshots agree. See
[scalar-storage-c-qualification.json](scalar-storage-c-qualification.json).
No C algorithm is added or retained for the test: it accesses real storage directly.

The first probe failed cgo's declaration-identity check for eleven new fixture
externs: existing Go users spelled them `unsigned int`, whereas the fixture used
`uint32_t`. Match the existing spelling; their width and bit semantics remain the
same. The corrected probe and every formal baseline gate pass. This is a fixture
correction, not a change to production storage or behavior.

## Integration and delegation

GPT-6 Luna produced the reference inventory and guarded migration draft. Primary
review requested complete preamble declaration types and independently checked
same-name locals, unused definitions, assembly/C++ references and storage lifetime.
Luna's draft passes a no-write preflight: 343 new Go owners, nine unused-definition
deletions, 3,576 selector replacements and 1,908 extern declarator removals across
270 source paths. It explicitly includes the new tagged fixture and preserves its
expected values. The script, patch and machine-readable inputs remain ignored
under `build/port-final-storage`; they are not qualification evidence.

The primary designed the baseline, owns integration and will run a broad tagged
consumer sweep for this cross-cutting storage change, plus three-profile
production/ABI, exact known-suite comparison and headless gameplay/save-load.
The previous four orphan-removal tests alone do not cover this batch.
