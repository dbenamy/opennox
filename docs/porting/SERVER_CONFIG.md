# Server configuration, admission persistence and rule picker

## Scope and current state

The native conversion replaces **64 live functions / 816 corrected C body lines**
in GAME1 409E40..40A740 and 4161E0..4169F0, GAME3_2 4CEBA0..4CF060, and
GAME3_3 4E41B0..4E43F0. The disabled `sub_416690`, its private setter `sub_4164F0`
and their no-op callers are removed (two functions / 22 body lines).
Eight private C globals and 41 old function interfaces are retired; 32 thin
exports remain for actual C callers. Six panel mask bridges and the settings-label
bridge became unnecessary. Go callers and GUI callbacks invoke Go directly.
Shared configuration storage, list layouts and allocation ownership remain compatible.
No C algorithm is retained solely for testing.

Current production C is **45,473 physical lines / 74 files / zero reference C**,
a reduction of **920 lines** from corrected C. The original scope was 66 functions /
836 body lines; prerequisite fixes made it 838. Corrected C baseline **ce1399b2**
was committed/pushed before translation; its parent is panel conversion **80c9ca2c**.

Native focused, broader three-target and fresh production qualification all pass
with unchanged frozen expectations and gameplay references.

## Prerequisites found by independent C contracts

The standing authorization for confident reversible fixes applies to these:

- Missing `rulelist.wnd` crashes during child access. Return zero immediately
  when the root load fails; retry succeeds. Original failure: focused-1.log.
- Expiring consecutive admission entries leaves the second entry behind because
  the loop advances its index after deletion. Advance only when retaining the
  current node. Original failure: focused-3.log. The fix passes all 32 mixtures
  of five permanent/timed entries, strict deadline equality, and repeated cleanup.
- Admission persistence passes UTF-16 game names to libc's `%S`, which expects a
  different character width on this target and omits names. Convert through the
  existing game formatter locally, then write `%s`; leave the shared file formatter
  unchanged. Original failure: focused-4.log. ASCII and legacy byte-encoding
  round trips pass, including 25-character names and low-byte truncation.
- The rule picker reads the alternate filename field, left empty by the Linux
  file adapter, and shows blank rows. Use actual filenames, strip `.rul`, filter
  the current map and `user` case-insensitively, and skip directories. Increase the
  local title and action filename buffers to accommodate listed filenames.
  Original failure: focused-5.log. Real-file enumeration, a 100-character name,
  create, duplicate, overwrite, load and delete contracts pass; an unrelated file
  remains unchanged. The actual entry resource still limits newly typed names.

- Settings refresh copies 12 map-name bytes into the established 9-byte field,
  overwriting the first three server-name bytes. Normalize to eight map-name bytes
  plus a terminator and copy/compare exactly that field. Original failure:
  focused-9.log. This matches the already-qualified settings editor's field boundary;
  contracts cover short, full-length and overlong names, both host/client paths,
  repeated refresh, actual players/spells/items, and unchanged neighboring fields.

These production changes require fresh corrected-C production qualification.

## Compatibility details to preserve

- The C clock adapter truncates platform ticks to uint32. Deadline reset/remaining
  and timed admission expiry have distinct arithmetic widths; boundary tests cover
  each. The first deadline test incorrectly assumed full-width ticks and was fixed
  after tracing the adapter; this was a fixture correction, not a C change.
- Shutdown of admission lists returns the original first blocked-node address
  after freeing it. Compare identity with the owned node before recording; never
  dereference or capture a raw freed address. The initial round-trip test's null
  return assumption was corrected without changing C.
- Admission files use the existing low-byte string encoding, not UTF-8. Preserve
  truncation at a narrowed NUL. Broader Unicode file compatibility is separate work.
  Preserve the filesystem adapter's final-line/EOF behavior and its 1,023-byte
  scratch copy; parsing still mutates the shared scratch buffer.
- Rule entry completion interprets its UTF-16 text as narrow text, and the draw
  callback updates button availability separately. Capture these existing paths
  explicitly rather than silently changing them during translation.

- Rule filenames retain the existing byte-widening behavior: UTF-8 filename bytes
  display as separate wide characters. Broader Unicode support is deferred.
- Refresh can report changes repeatedly even when bytes remain unchanged: signed
  weapon-mask bytes compare against unsigned storage; headless empty counts compare
  stored 255 against -1; a preserved video flag can differ from the masked input.
  Independent contracts cover these notification semantics without changing them.

## Fixtures and qualification

Reuse real GUI/list/entry/font/render owners, reliable queues, players, spell and
item definitions, actual list allocation/removal, and isolated temporary files.
Seed filename suffixes, section headers and nonzero lookup tables explicitly.
Own the mapped scratch buffer independently of extracted live globals. The rule
fixture uses the shipped entry data (`9 130 0 2`), leaving earlier fixtures unchanged.

The 26 focused roots cover scalar state, exact name/password storage, settings
slots/selection, report reset dispatch, timer widths/notifications, rate notification
ordering, lookup boundaries, admission lists/expiry/persistence, rule-picker
lifecycle/files/actions, pixels/control state and full settings refresh. Independent
contracts run before the frozen-C comparisons. Twenty-five new C captures repeat
byte-for-byte in separate processes before freezing. The broader selection extends
the parent's affected corpus and asserts 240 captures (215 parent plus 25 new).

| Gate | Corrected C | Native Go |
| --- | --- | --- |
| Focused 26 roots | 1.706s | 2.127s |
| Static memory preflight | 0.592s | 0.398s |
| Default: 535 roots / 58,651 leaves | 215.15s | 178.11s |
| Server: 534 roots / 58,650 leaves | 296.09s | 283.13s |
| Highres: 535 roots / 58,651 leaves | 229.65s | 222.88s |
| Fresh production | 371.49s | 377.75s |

All broader runs have no skips and **240 identical captures / 76,002 records**.
Native results match corrected C and every inherited capture remains unchanged.
Corrected-C production passes three builds/ABI/interface checks, the exact known
1,553 asset failures (15 passing / 3 failing / 32 no-test packages), options
interaction, save/load and forced flat-map regeneration. Its four gates share an
unchanged 2,161-file source manifest. All four native gates share 2,167 unchanged
source files; all sessions have joined. Native production passes the same exact
asset failure set, all three builds/ABI/interface checks, options gameplay,
save/load and forced flat-map regeneration against corrected C.

C client SHA: `d3d78b8b24f265044e2289e4bf7a15d86e8a1b6d79f7c6b143b237296f5ae219`.

Native client SHA: `7c3a16ffb2531aac2c9b6a4368380b41982b89de2af36483c5c5b9d30017ef95`.

## Native review and recovery

The pixel capture caught reuse of a blended panel background for an opaque picker
background. The native code now calls the original opaque primitive using the
renderer color; all nine drawing cases match without changing goldens. Review
also preserves byte-by-byte string conversion, signed comparisons and clock/read
ordering. The Cgo admission declaration drops an unrepresentable const qualifier
without changing its binary interface. All 149 required exports are checked present,
including the live spell-apply bridge; no Go caller routes through the selected C
functions. Existing root Go helper names/comments are not retired C interfaces.

Manifests: `server-config-batch.json`, `server-config-tests.txt`,
`server-config-scope.json`. Local evidence is under `build/port-server-config/`:
`c-{default,server,highres,production}`, `native-{default,server,highres,production}`,
`c-audit.json`, `native-audit.json`, `repeat-audit.json`, `native-interface-audit.json`,
focused and static logs. The original fixture drafts, native drafts and installation
scripts are **consumed**: never recopy them over the corrected source or regenerate
frozen expectations during translation. Original assets/archive remain untouched.
