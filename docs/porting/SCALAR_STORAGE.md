# Numeric global storage

Status: native conversion qualified. C is 166 physical lines /6 files (352 fewer),
zero reference C. Original-C baseline `bc90690c` is pushed. All 343 new owners
preserve their numeric types and initial bits. See
[scalar-storage-native-qualification.json](scalar-storage-native-qualification.json).

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
The applied draft passed its no-write preflight: 343 new Go owners, nine unused-definition
deletions, 3,576 selector replacements and 1,908 extern declarator removals across
270 source paths. It explicitly includes the new tagged fixture and preserves its
expected values. The script, patch and machine-readable inputs remain ignored
under `build/port-final-storage`; they are not qualification evidence.

The primary designed the baseline, integrated the change and ran a broad tagged
consumer sweep for this cross-cutting storage change, plus three-profile
production/ABI, exact known-suite comparison and headless gameplay/save-load.
The previous four orphan-removal tests alone do not cover this batch.

The primary removed an unnecessary `aaa_` filename prefix from the draft: no
initialization-order dependency justified it. The owner file is `scalar_globals.go`.
The migration installer is consumed; tracked source is authoritative. Native
storage contracts/static pass all three profiles. The full default consumer sweep
passes all 2,291 roots without skips; all remaining gates subsequently passed.

## Native qualification plan

A milestone sweep selects all 2,291 ordinary root-package tagged tests in 12
file-preserving groups of at most 200 roots. Run default/server/highres with shipped
asset inputs set, assert every discovered root finishes without skips, and keep
frozen hash assertions active. The standalone population prerequisite child is
excluded because its parent invokes all four prerequisite cases; other Probe-named
tests remain included. Grouping bounds the 386 process lifetime and memory growth.
Eleven client-only roots are excluded by build constraints under server, so its
expected count is 2,280. The legacy storage contract runs separately. Optional safe build, static checks,
three production binaries/ABI, exact known-suite comparison and fresh headless
character creation/gameplay plus explicit save/load complete the gates. Existing
full-suite failures are compared, not described as a green suite.

The first final-production launch failed before any build because its manifest
argument omitted `-native`. Corrected that documentation path and reran only the
production phase; existing source/test evidence was unchanged. The retry
production output is `native-production2`, not the failed launch directory.

## Completed gates and review

Default/server/highres pass all 2,291/2,280/2,291 ordinary tagged root tests, no
skips, plus the independent legacy storage contract and static checks. Every
storage capture matches the frozen 47,677-case C result. Safe build and symbol
checks pass (safe runtime was not tested). Three fresh production binaries and
ABI checks pass, as do exact known-suite comparison, headless gameplay and explicit
save/load. All final source fingerprints agree. No frozen expectation changed.

Primary independently compared the AST of all 268 modified existing Go files:
only intended owner selectors and eligible cgo imports changed. All 343 new
owner types/initializers match C; nine unused definitions have no active users.
Luna supplied the large mechanical draft and useful test/consumer inventories.
The primary removed an unjustified filename-order precaution. A later unapplied
fixture draft had a duplicate helper name; review caught it and Luna corrected it.
Continue one bounded Luna helper with primary integration and acceptance.

Next batching decision: combine the remaining 52 numeric owners used by C fixture
accessors (map, audio, sustained spells) because the ownership/accessor conversion
is the same. Preserve subsystem-focused tests and frozen original-C contracts;
avoid three separate full production cycles for identical migration mechanics.
The three pointer-typed numeric declarations and pointer/array owners stay deferred.
