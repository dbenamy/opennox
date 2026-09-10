# Map-rule loading — 2026-09-10

Scope: mode-header lookup 57A1B0, top-level loader 57A1E0, file reader 57A3F0,
UTF-16 tokenizer/rejected-line collection 57A4D0 and directive dispatcher 57A620.
Rule saving/removal and generic C-owned list lifetime are separate chunks.

## Original-C baseline

The `porttest` fixtures execute the original five C functions. The settings
buffer, header strings, token pointers and rejected-line nodes are C allocated;
all global context, flags, header table, working directory and handle arena are
restored. A bare typed server supplies controlled spell definitions and the
shipped weapon/armor tables through the actual existing Go semantic lookups.
No asset loading or global blob initialization is required for these tests.

- All 65,536 mode-header inputs check exact original string-pointer identity.
- Directive traces cover every valid shipped weapon/armor entry, spell IDs and
  titles, invalid/non-rule spells, case rules, four operations, ten requested
  mode combinations, three section headers and host/non-host execution.
- Explicit line/token pairs cover whitespace, comments, malformed and quoted
  commands, UTF-16 titles, embedded NUL, 32 tokens and quoting quirks. Each runs
  across four modes and both host states, with/without trailing whitespace.
- 640 selection/presence/online combinations cover user-file precedence,
  fallback to map rules, online rules applied last, mode resets even on failed
  opens, reset of settings bytes24..51, and the final flagball spell132 clear.
  Eight additional cases cover custom names, the Go wrapper, an empty user
  file blocking fallback, case-insensitive paths, eight-byte map names,
  short names, a nil rejected list and unused selection bits.
- Forty direct-file cases cover CRLF, final lines without LF, bare CR, empty
  and missing files, NUL, physical-line truncation at 255 bytes, and Latin-1
  byte widening versus UTF-8 bytes. Direct reads preserve pre-existing settings
  and rejected entries; top-level loading clears the rejected list first.

The independent model consumes explicit expected tokens rather than sharing
production tokenization. It calculates settings bit positions independently.
Snapshots check all 60 settings bytes (including padding), surrounding guards,
list contents and bidirectional links, unchanged header table and file contents,
and no leaked registered C file handles.

Historical behavior matters: COMMON commands can modify settings while still
being retained in the rejected list; spell `on` never re-enables a disabled
spell, while weapon/armor operations other than `off` enable their bits. File
bytes are widened to U+00xx, not decoded as UTF-8. The C filesystem bridge reads
an entire physical line before copying at most 255 bytes, discarding the rest.
Quotes are recognized only at the legacy tokenizer's specific cursor positions.

Baseline inputs stay within defined C stack-buffer/token limits. Directory
opens and other non-EOF read errors are excluded because the original loop can
fail to terminate. The port must handle ordinary files faithfully without
reproducing memory corruption or an endless read-error loop.

Local artifacts: `build/port-rules/`; no C reference implementation is retained
solely for testing after conversion. Production C before conversion: **141,745
physical lines**, 153 files, zero test-reference C lines.

All accumulated protection/network/waypoint/rule ABI tests passed against the
original loader on 386 default, server and highres before replacement.

Six additional directive cases passed against C before integration: Unicode
lookalikes must not become ASCII command keywords, while byte narrowing of
header/spell names can introduce a terminating NUL. These distinguish wide
keyword comparisons from the separate legacy narrow-name conversion.

## Go conversion

Original-C checkpoints: `2bf05750` (main baseline), `9c86046c` (wide/narrow
character distinctions). The loader now lives in `legacy/rules.go`. Only the
header-pointer and top-level loader C bridges remain because other C code calls
them. File reading, parsing and dispatch are private Go functions, and the old
C context global is replaced by private Go state. The Go caller bypasses C while
preserving filename NUL truncation and the distinction between an empty filename
and the C nil pointer that selects `user.rul`.

Typed Settings2 fields replace offsets for resets, spell bits and equipment.
Header strings still come from the existing table to preserve pointer identity.
Rejected nodes remain calloc-owned for the C writer and generic list cleanup.
The pure Go token slices are not exposed to C. Existing typed server methods
continue to handle spell and item lookup. Review corrected the draft's blank-line
handling and wide-keyword comparison before final validation.

Read errors now stop the loop, including an explicit empty filename opening a
map directory. Two Go-only regression cases cover the live C entry and Go caller;
this is a termination improvement, not a claim that the nonterminating C case
produced equivalent output. Ordinary file behavior remains baseline-compatible.

Production C: **141,455 physical lines (−290)** in 153 files; test-reference C:
**0**. No test-only C copy is retained. The 290-line reduction includes the five
functions and obsolete C context declaration/definition; generated live bridges
remain outside this physical C-file count.

All accumulated protection/network/waypoint/rule tests, including the new read
error termination case, pass on 386 default, server and highres after conversion.

All three production binaries build. Symbol inspection confirms the two live
bridges plus native Go loader/parser and no retired internal symbols/context.
`rules-port` passes the preserved headless gameplay scenario with both screenshot
checks and overrides disabled. The full suite exactly matches the waypoint
milestone: 15 passing, 3 known failing, 32 skipped/no-test packages, with the same
1,553 failure entries and no additions or removals.
