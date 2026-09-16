# Map decompression

Round 2 of the process trial selects the complete 969-line C map decoder: initial
format tables, adaptive symbol ordering, 64 KiB dictionary, copy expansion,
byte-aligned block boundaries and lifecycle. The C compressor remains live
production code. Its allocator helper is shared and remains necessary.

The replacement uses Go-owned arrays and checked bit reads. The existing Go
reader in the pinned libs dependency was useful for comparison, but its handling
of the end-block symbol does not implement the C decoder's byte alignment. This
batch owns its decoder locally rather than changing dependencies or silently
substituting a reader with different chunk semantics.

## Contracts

Freeze original C output for deterministic synthetic data around length-code,
frequency-rebuild, dictionary and 500000-byte chunk boundaries; five data patterns
include incompressible and highly repeated input. Assert expansion equals the
independent original input, and freeze compressed bytes too. An unchanged C
compressor provides production-format inputs after conversion, not a test-only
reference implementation. Fifty shipped .nxz maps have independent .map files.

Hand-encoded files independently exercise block alignment, cross-block dictionary
history and overlapping copies. File contracts check path errors, short headers,
destination preservation on pre-open failure, overwrite and same-path behavior.
Native internal checks cover signed16-bit frequency ordering, ties, arithmetic
halving and increment wrap, plus malformed-stream errors.

Malformed bodies: the old C decoder contains unchecked reads and can emit a
zero-filled partial map without reporting failure. New code returns a descriptive
error before opening the output. This is a deliberate reversible correction under
standing authorization, recorded for review; it is not a claim to reproduce
undefined C accesses. Valid original captures remain immutable.

## Qualification

The explicit affected package is ./legacy/cnxz. Existing compression/decompression
asset tests and the new contracts must all run without skips in three variants.
The complete accumulated corpus across packages is the combined two-round
milestone. Production builds/ABI and the exact full-suite comparison remain.

Gameplay will use fresh copied assets with expanded maps removed only where a
compressed counterpart exists. Require a real regenerated map, compare its bytes
to the source asset, and retain all normal/flat-floor pixel references. Thus the
integration check exercises the actual changed loading path.

Status: final original-C baseline qualified; native decoder remains unapplied.

Baseline development found two useful edges. Empty source data panics in the
existing compressor's zero-sized allocation before it reaches C; the old decoder
also allocates zero-sized buffers for an empty declared output. The unchanged
compressor is therefore tested on positive sizes. Native decoding will accept a
four-byte zero-size header as an empty map, with an explicit independent test.
The compressor's empty-input limitation remains outside this decoder conversion.

The initial asset loop inherited the older directory-name convention and covered
47 pairs. Enumerating actual .nxz files adds BankShot/Bankshot, Fortress/fortress
and TreeHaus/treehaus without changing assets. The final oracle contains **174
records**: 115 synthetic, 50 real maps, nine hand-encoded block/history cases.
File API checks are independent assertions, not additional frozen records.

Final C qualification: default 4.061s, repeat 3.309s, server 3.538s, highres 3.418s.
All six selected package roots (including existing codec asset tests) completed
without skips; 174 frozen records matched on every run, with source unchanged.
