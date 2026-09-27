# Native book/UI callback identities

Replaced five C tooltip addresses with existing GUI-native dispatch registrations,
and two book-page image-completion addresses with stable native identities. The
image-completion call keeps its original foreign fallback and ImageRef argument.
Book callbacks still ignore that argument and update the same extracted globals.
No C algorithms, renderer backend, image layout or page-selection math changed.

Removed the unused `sub_4C3260` export and five fixture-only C-typed reward/removal
wrappers. Fixture operations call the same Go owners and keep explicit signed 32-bit
conversions. Live identity entries use native keys; the unused entry becomes nil
without changing the name/index order. Existing state and pixel goldens are intact.

Independent original contracts cover last-frame entry, repetition, nonzero start,
frame period, exact reference pointer, callback mutation order, loop/nil/unsupported
paths, actual spell/creature owners and record layout. The minimal foreign observer
is Go code behind a test-only C ABI entrypoint, not a C reference algorithm.
The original preflight found a fixture assertion reading inert backing-blob words;
corrected it to extracted owner globals before freezing the baseline.

Qualification passed:

- 256 default/high-resolution roots and 254 server roots, with exact original
  names, no failures/skips and unchanged frozen expectations. The two server
  exclusions are existing `!server` rendering tests.
- Original captures passed twice per profile. Baseline `b5831dc9` adds only two
  porttest files, so qualified production `1238c985` supplies the production baseline.
- Safe/static and all three production/ABI builds passed; eight retired exports
  are absent. Preview/final character creation with save/load/resume passed on
  the same final binary.
- The full asset suite matches the known result: 304 failure events, 17 passing,
  two failing and 32 skipped packages. All 1,654 original asset hashes are unchanged.
- Measured: 95 client/96 server production cgo files (down seven), 80 selected
  legacy exports (down eight), 157 headers/2,809 physical lines. Embedded C bodies
  remain 77; standalone production/test-reference C remains 0.

[Baseline](book-callbacks-baseline.json),
[qualification](book-callbacks-qualification.json),
[inventory](book-callbacks-inventory-after.json),
[test selection](book-callbacks-tests.txt),
[batch manifest](book-callbacks-batch.json).

The callback scope is closed and audited: one image-completion call site, two
engine keys, five tooltips using an unchanged dispatcher. A full accumulated
porttest corpus is not rerun for this bounded family. The last full default corpus
remains earlier-source `f8332e3d`; it is not claimed as converted-source evidence.

Luna was unavailable due to quota. Primary performed the draft, caller/identity
review, baseline and acceptance locally; no alternative model or agent fleet.

Cleanup retained originals, source, saves and logs. Thirteen superseded binaries
were removed after source/replacement/hash/host-use checks: rebuild `f8332e3d`,
`1238c985`, `0121a5f3` or `b5831dc9` as recorded in the batch journals. Fifteen
obsolete project cache archives were removed. Three historical failed-setup
binaries were losslessly archived; restore with `gzip -dk FILE.test.gz`,
`chmod 755 FILE.test` and the recorded SHA256. Preview/final each deduplicated
1,654 original-asset copies after passed qualification. See the checkpoint for
restoration paths; completed cleanup/finalizer scripts must not be rerun.
