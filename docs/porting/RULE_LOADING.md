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
