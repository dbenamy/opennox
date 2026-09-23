# Remaining global and mapped-buffer storage

Status: native conversion qualified against pushed baseline `781901ef`. C is now
51 physical lines in 4 files (−63), zero reference C. The underlying foreign
allocator and libc/CGO dependencies remain.

## Scope and memory ownership

Move the 44 remaining vardefs owners and eight cgo_blobs arrays into Go initialization.
They include 37 pointer slots, four integer words, two arrays and a list head;
the eight buffers total 4,014,405 bytes. Two definitions rely on implicit static
zero initialization. No remaining owner was proven unused.

Use the existing legacy/common/alloc foreign allocator from Go. One typed record
holds the 44 fields; the eight buffers retain their individual extents and logical
addresses. Allocate once and retain for the process, like the original C globals.
This preserves unmanaged raw pointer/integer storage, C address retention and
address stability without introducing GC scanning of pointer-shaped integer bits.
No finalizer or cleanup hook frees these process-lifetime regions.

This decision removes C global definitions; it does **not** replace the underlying
libc allocator or make the program CGO-free. Foreign allocation remains explicit
follow-up work. A Go pointer-bearing global would change GC behavior; raw Go
buffers would require a separate review of pinning and retained C addresses.
The additional allocation/initialization step is a deliberate, reversible choice.
The physical `.c` count excludes headers, inline cgo preambles and third-party
dependencies. Removing these storage definitions does not remove those remaining
C types, callback boundaries or build requirements.

Keep current C record/field types while preserving the actual definition widths.
Three uint32 definitions have pointer-like extern declarations in some callers.
The briefing word already has integer-valued active selectors; options and
scoreboard have active pointer-valued selectors. Their pointer reads require
explicit conversions when moved into canonical uint32 fields. Address-only casts
keep the same four-byte storage. Never infer adjacency from the order of owners.

## Baseline contracts

The 52-region contract checks frozen widths, pre-init zero state, target alignment,
non-overlap (including the 395 numeric owners), seventeen full-region byte patterns
per owner /884 cases, neighboring owners, alias visibility, GC stability and
restoration. Odd-sized blob tails are included. Typed aliases separately check
pointer/integer slots, all nine equipment slots, all 84 inventory cells and the
three list links. Layout expectations are independent constants.

Blob checks cover registry identity, initial snapshots, reverse pointer lookup,
exact extents and offsets of all registered variables that fit the backing region.
Production inventory/list/scoreboard aliases and the three mixed-declaration window
getters are checked directly. The game-level pointer setter identified as a gap in
Luna's consumer audit also has a direct contract. Synthetic pointer-bit patterns
remain raw bytes; typed pointer stores use real foreign addresses.

The first compile identified cgo's inconsistent-definition check for the two
pointer-valued active declarations. The fixture now follows those declarations,
while independent size/layout expectations still come from the actual C definitions.
The next compile caught a missing invocation of the new pre-init snapshot closure;
primary corrected it. These are fixture repairs before capture freezing, not
production changes.

The corrected baseline passes all three target contracts and the full consumer
sweep. The audio correction required fresh production evidence instead of reusing
`05182782`. After migration, require the full 12-shard consumer sweep in
default/server/highres, static checks, safe build/symbol audit, fresh client
preflight, three production builds/ABI, exact known full-suite comparison and
explicit save/load. Source fingerprints must agree across final gates.

## Delegation and review

Luna supplied read-only owner/reference and focused consumer inventories, including
479 candidate consumer roots and a specific direct-setter coverage gap. Including
the eight shared buffers expands acceptance to the full consumer sweep. Luna supplied
the mechanical migration draft; primary owns the memory model, contracts,
independent review, integration and qualification. Ignored drafts are not recovery
artifacts; tracked source and committed reports are authoritative.

## Baseline interruption

The isolated 52-region probe passes and both raw/numeric formal default contracts
match their captures. The first broad consumer shard then failed in audio stream
buffer tests: an opaque data-address word was queued by a Go pointer write barrier
and rejected during GC. No remaining-storage migration had been applied. A focused
subprocess regression and bounded audio representation repair precede baseline
acceptance; production qualification must be refreshed after that repair.

The ignored mechanical draft required three source-review corrections before
compilation: convert blob sizeof expressions before replacing C selectors, strip
C array-declarator brackets before emitting Go brackets, and use AsWindowP for the
scoreboard getter's unsafe pointer. The generator's Python check did not establish
Go syntax/type correctness. After correction, primary independently compared all
44 fields and 8 buffer lengths/registrations directly with committed C source;
that review passes. AST comparison and compilation remain required after application.

## Accepted baseline

The audio representation correction passes the complete three-profile consumer
sweep and fresh production qualification. Original C definitions in vardefs.c and
cgo_blobs.c are byte-identical to05182782. The 52-region capture is unchanged from
the original probe and matches across target profiles:
`d83791277fd3c4d8cbc1744a3058abe27732e796e0cb165f7daac9e25dac0fd3`.
The existing 395-owner capture also remains unchanged. The successful acceptance
manifest is [audio-address-native-batch.json](audio-address-native-batch.json),
which includes all three independent contracts and the broad consumer sweep.
See [raw-storage-c-qualification.json](raw-storage-c-qualification.json).

## Native implementation and review

One process-lifetime foreign record owns the 44 original field types and array
dimensions. Eight foreign byte slices retain the original extents and registry
addresses. The conversion removes vardefs.c and cgo_blobs.c, 52 C storage symbols,
and 217 extern declarations. The canonical integer words used by pointer getters
are converted explicitly at the getter boundary.

Independent source review matches all field types, dimensions, zero initialization,
buffer sizes and logical registrations to committed C. AST comparison verifies 59
files differ only by owner selectors (ignoring C declarations); five Go files with
larger changes receive separate review. No retired names remain in C preambles.
Direct extent checks find the expected 84-cell and nine-word array views, with no
direct owner adjacency dependency in the expressions inspected. Compilation and
full native qualification pass; frozen captures are unchanged.

## Qualified native result

The first compiled probe passes all three legacy contracts. The final default,
server and high-resolution profiles pass 2,291/2,280/2,291 consumer roots, no skips,
and all 884 raw-region /54,905 numeric patterns with their original capture hashes.
The audio GC regression passes in each profile. Static checks, safe build/symbol
checks, three production binaries/ABI, exact known-suite comparison, fresh headless
gameplay and explicit save/load also pass. Safe runtime was not exercised; the
known-suite comparison preserves recorded failures rather than claiming a green
full suite. Final source fingerprints match across every phase, and the gameplay
preflight binary matches the final default production binary.

The production storage patch needed no source correction after application.
Qualification exposed a fixture identity-lifetime leak: per-case freed addresses
remained in the shared normalization map. A deterministic boundary assertion failed
on the second stats case before repair; restoring the original address set after each complete
case, with updated persistent aliases, preserves all frozen captures. Failure-only callback captures now retain JSON
if a mismatch recurs. See [FIXTURE_IDENTITIES.md](FIXTURE_IDENTITIES.md). Luna's ignored
draft received the three review corrections described above before its first
compile. The primary reviewed, integrated and qualified it; no measured subscription
savings are claimed. See [raw-storage-native-qualification.json](raw-storage-native-qualification.json).

Remaining standalone C files are the optional safe adapters (20 lines), modifier
callback identities (15), other empty callbacks (10), and the MP3 implementation
include (6). Headers, inline cgo bodies and third-party implementation size are
outside that count and require their own inventory before any claim of C removal.
