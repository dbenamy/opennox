# Native window dispatch and tooltip observation

## Scope and baseline

Retire the unproduced raw window event/draw callback routes, translate the tooltip
observer, and remove the unused raw-pointer quickbar-constructor fixture adapter.
Keep all native window handlers, layout fields, snapshot identities, assertions
and frozen captures unchanged. Baseline source: qualified duration/audio revision
`d3d759ba`. Seven-file conversion remains an ignored draft.

Original baseline accepted without skips. The broad audited root selection passed:
2,482 default /2,471 server /2,482 highres, plus six separately repeated focused
contracts and the asset-free safe deferred-cleanup contract. Native qualification also runs
a fresh default-client save/load scenario immediately after focused preflight,
before the broad sweep and complete production qualification. GUI event/draw and
destruction code underpins many shared UI owners, making the broad selection
appropriate here. The existing eight excluded compiled roots remain explicit in
the selection record; they are outside this inherited audited selection, not
all declared known failures. Full production gates retain the exact known-suite
comparison. Do not describe this as a complete root-corpus pass.

## Producer and behavior review

No writers of the raw event/draw slots were found by field names, byte offsets
372/376/380, word indices 93/94/95 or hex offsets. Window pool allocation zeros the
record and construction installs all three native extension handlers. The retained
C headers declare layout fields only. Keep nil/dead-window and nil/dead-word
slot guards, native extension priority and the native deferred-destroy snapshot.
Non-nil, non-sentinel unsupported raw callbacks now panic explicitly instead of
entering C. This is a scoped invalid-input correction after the producer audit.

Tooltip producers register native identities; migrate its one remaining foreign
observer through the existing registration API. Preserve its four uint32 capture
words (count, window, draw data, argument), exact state restore, argument boundaries
and replacement behavior. Native tooltip KeepAlive remains unchanged.

The quickbar raw-pointer constructor selector is uncalled; the native constructor
remains used by production and its typed fixture. Remove only that selector case
and its unused import. **Keep its nil entry in the callback snapshot map** because
sorted names determine canonical IDs of other callbacks. This avoids changing
existing capture identities simply by removing an unused adapter.

For later review: popup construction currently wraps its parent's raw slot 376,
which has no producer, and consequently passes nil to the resource parser. Keep
that behavior here, with explicit rejection of a non-nil raw key. Substituting
`parent.Func94` would change behavior and belongs in a separate owner-level review.
Do not confuse this with deferred destruction, which already captures the native
extension callback and must continue doing so.

## Tooling, delegation and recovery

Function comparison now keys methods by receiver type and distinguishes init
functions. Its first receiver-formatting attempt failed loudly and was corrected
before accepting the draft review. Generated dispatcher output matches the
edited generator byte-for-byte. Ignored evidence: `build/port-window-dispatch/`.
Primary implementation/review; Luna quota remains unavailable.

Removed seven superseded modifier binaries after committed-source, replacement,
hash and host-use checks, reclaiming 387,018,752 bytes. Rebuild `5ba6fce6` for those
historical outputs; logs/metadata remain. Exact paths and hashes:
`modifier-cleanup-approved.json` and deleted journal in the batch directory.
Current duration/audio binaries and original assets remain intact.

After all baseline jobs finished, removed 39 inactive Linux 386 cache archives
untouched for six hours (1,941,274,624 allocated bytes) after metadata/hash and
host compiler/fd/maps checks. These are rebuildable caches; module sources, current
binaries, logs and assets remain. Records: `cache-headroom-{approved,result}.json`
and deletion journal in the batch directory. Physical free space afterward:
3,030,851,584 bytes.
