# Server configuration, admission persistence and rule picker

## Scope and current state

Parent server-panel conversion **80c9ca2c** is qualified, committed and pushed.
The corrected C expectations are frozen after two identical process runs; no
native conversion is installed. Candidate scope is **66 functions / 836 original C body
lines** in GAME1 409E40..40A740 and 4161E0..4169F0, GAME3_2 4CEBA0..4CF060,
and GAME3_3 4E41B0..4E43F0. Corrected scope is 838 body lines (64 live functions / 816 lines plus two dead
functions / 22 lines). Working production C is 46,393 lines in 74 files (+2).

`sub_416690` has live callers but only a literal `if (0)` body. Its private setter
`sub_4164F0` is used only there. Audit and remove the no-op calls and dead helper
when converting; do not restore the disabled report code.

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
- Rule entry completion interprets its UTF-16 text as narrow text, and the draw
  callback updates button availability separately. Capture these existing paths
  explicitly rather than silently changing them during translation.

- Rule filenames retain the existing byte-widening behavior: UTF-8 filename bytes
  display as separate wide characters. Broader Unicode support is deferred.
- Refresh can report changes repeatedly even when bytes remain unchanged: signed
  weapon-mask bytes compare against unsigned storage; headless empty counts compare
  stored 255 against -1; a preserved video flag can differ from the masked input.
  Independent contracts cover these notification semantics without changing them.

## Fixtures and qualification work

Reuse real GUI/list/entry/font/render owners, reliable queues, players, spell and
item definitions, actual list allocation/removal, and isolated temporary files.
Seed filename suffixes, section headers and nonzero lookup tables explicitly.
Own the mapped scratch buffer independently of extracted live globals. The new
rule fixture uses the shipped entry data (`9 130 0 2`), leaving earlier fixtures
unchanged.

Focused contracts cover scalar setters/getters, exact name/password storage,
settings slots and selection, report reset dispatch, timer widths and notifications,
rate notification ordering, lookup boundaries, admission lists/expiry/persistence,
and rule-picker lifecycle/files/actions, pixels/control state and full settings
refresh. Tests include independent assertions before frozen-C comparisons.

All 26 focused roots pass in 1.706s (`focused-11.log`). Two separate process runs
produce 25 byte-identical captures (`capture-ready-*`, `capture-repeat-*` and
`repeat-audit.json`). These are now frozen in tests and the batch manifest.
The broader selection extends the parent's affected corpus with all new roots;
its manifest asserts 240 captures (215 parent plus 25 new).

Local evidence and consumed drafts are under `build/port-server-config/`.
Never recopy consumed drafts over the corrected fixtures or regenerate goldens
while translating. Next: qualify/commit/push this C baseline, then convert and
qualify native Go.

Corrected-C broader default/server/highres qualification passes **535/534/535
roots**, **58,651/58,650/58,651 leaves**, no skips. All **240 captures / 76,002
records** match across targets and all inherited expectations remain unchanged.
Durations: **215.15/296.09/229.65s**; static preflight **0.592s**. All three
source manifests match (2,161 files). Sessions are joined. Fresh production passes in **371.49s**: all three builds and ABI/interface checks,
the exact known 1,553 asset failures (15 passing / 3 failing / 32 no-test packages),
options gameplay, save/load and forced flat-map regeneration. All four gates share
one unchanged 2,161-file source manifest; all sessions joined. Client SHA:
`d3d78b8b24f265044e2289e4bf7a15d86e8a1b6d79f7c6b143b237296f5ae219`.
Artifacts: `c-{default,server,highres,production}`, `c-audit.json`. No native source
has been installed; ignored draft files require integration and validation.
