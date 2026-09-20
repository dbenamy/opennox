# Client drawable state and update streams

Qualified parent `ad78ea93` is pushed: 9,963 C lines /39 files /zero reference C.
Original-C baseline `8266e740` is pushed; the native conversion is fully qualified.
Current C: **9,267 lines /39 files /zero reference C**, down 696 lines.

The selection contains 16 C definitions /532 body lines: the two compact drawable
update readers, drawable state/list accessors, two motion callbacks, summon/shield
creation and lookup, and effect-type initialization. Whole-repository caller
inventory is recorded alongside the selection. The minimap next helper currently
has no callers; registration/literal auditing confirms it can retire without translation. The live C message decoder, map traversal and registered motion
callbacks still require bridges after conversion. Private Go callers should move
together with their helpers.

The existing NetworkAlias tests cover alias selection and announcements but use a
boundary recorder for sprite creation. New stream contracts use the real drawable
pool, lookup, position index and outgoing queue. They cover explicit/alias records,
absolute/signed-relative coordinates, all status bytes, static/dynamic and
player/monster animation behavior, repeated updates, camera calls and frame wrap.
Preserve the different operation order: stream 1 updates animation before its
frame helper, whereas stream 2 calls that helper before changing animation.

Additional contracts cover coordinate boundaries/early returns, missing
and newly created objects, state accessors, real motion deletion/list membership,
shield lookup/duplicate suppression and summon ownership. Use original-C captures
plus independent assertions. Inspect compiled motion arithmetic before freezing:
its x87 intermediates and explicit float stores need separate treatment. Never
feed impossible frame ages that would iterate billions of times.

Production qualification may reuse the parent only after source identity proves
all baseline changes are test-only. The Go conversion needs fresh three-target
production/ABI and headless gameplay/save-load checks. Affected selection will
include alias, effects, motion, presentation and relevant drawable owners.

Completed spell-start scenario asset copies were verified against original assets
and deduplicated, reclaiming 1,112,747,701 bytes. Their results, saves, screenshots
and restoration manifests remain; the ignored apply script is consumed.

## Frozen focused C evidence

Nine focused roots pass final and separate repeat processes. Eight frozen JSON
captures contain 18,414 result rows: 12,800 stream updates, 74 rejected coordinate
records, 144 shield operations, 13 type-initialization sequences, 875 state cases,
2,304 summon cases, 2,160 motion cases and 44 creation/failure cases. Additional
contracts cover 432 shield-scan combinations, two terminator records and link
accessors. All 256 summon directions use independently audited shipped-table
ranges. No selected production source has changed.

An initial shield-scan fixture allocated its key using alloc.New but did not
assign it; that API initializes storage to zero, ignoring the exemplar value.
The fixture now assigns the key explicitly. All other focused cases passed;
this was a test setup correction before freezing, not a production change.

The minimap C accessor is a proven unused duplicate of the directly called Go
Drawable method. Symbol/address/registration searches identify no C caller.
Retire its C body without translating it. The summon cache is also used by the
Go object-animation renderer; move that consumer with its owner and include
ClientObjectDrawing tests. Blue/Violet spark caches still have C decoder callers
and retain their existing C ownership.

Static memory checks pass. All ten source changes are porttest-only; all three
parent production binaries rehash to their recorded values and retain the selected
raw-C symbols. Fresh C production/scenario runs are unnecessary; the Go conversion
will run fresh qualification. All three targets pass 88 selected roots without skips. Full capture inventories
are identical; all jobs are joined. Eight new hashes were enforced during the
sweeps; inherited tests enforced their committed literals. The complete inventory
was frozen after cross-target comparison for native checks. See the C batch manifest for the exact selected roots and capture checks.

## Qualified native conversion

C baseline `8266e740` is pushed. Native integration moves fifteen live bodies and
two private owners to Go, retires the unused C minimap duplicate, and retains eleven
Go-backed C interfaces for remaining callers/registrations. The existing animation
frame/active methods are reused. Shared Blue/Violet caches remain in C for decoder
consumers. Working-tree C is 9,267 lines /39 files /zero reference C (−696), qualified; the reduction includes obsolete address headings and blank lines.

The first native compile identified two missing uint32-to-int conversions around
the existing effectType helper. Both were corrected; frozen expectations are
unchanged. The final focused run passes all drawable, alias and object-summon roots against
frozen expectations. Static-native-final passes. All three native sweeps pass 88 roots /69 exact capture hashes without skips. The ignored native
installer is consumed; its drafts no longer supersede reviewed source.

All sweep and production fingerprints match the final source. Fresh default,
highres and server binaries pass ABI/export audits. The full asset suite exactly
matches 1,553 known failure entries and 15 pass /3 fail /32 skip packages. Both
headless gameplay and explicit save/load pass against the spell-start references.
All jobs are joined. Detailed evidence is in
`client-drawable-state-native-qualification.json` and
`build/port-client-drawable/native-*`.

Completed verified cleanup reclaimed 509,147,276 bytes of obsolete legacy test
archives and 1,106,323,366 bytes of obsolete root test archives. It excluded current
drawable fixture archives and non-porttest artifacts. Identical completed captures
share storage (316,265,478 bytes across C targets; 474,398,217 native bytes shared
with C). Original inputs, qualified production binaries, logs and snapshots remain.
The cleanup applications are consumed.
