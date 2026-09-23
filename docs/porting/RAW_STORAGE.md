# Remaining global and mapped-buffer storage

Status: actual-C storage baseline qualified together with the audio address
correction. C remains 114 physical lines /6 files, zero reference C. No remaining
storage migration is applied. Fresh production evidence comes from the corrected
source, superseding the initial production-reuse plan.

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
