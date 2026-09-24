# Spell-lifecycle and reward fixture owners

The proposed batch retires 13 spell-lifecycle and 12 reward-generation C exports,
routing fixtures to existing native owners. Nine reward initialization callbacks
remain live in production and keep their C addresses/signatures. Production
algorithms and frozen expected outputs stay unchanged.

All 51 focused roots pass in fresh default/server/highres processes against the
original `279e7867` source, using source-verified qualified binaries. The selection
includes spell effects and server reward orchestration as well as the two direct
fixture families. See [baseline](spell-reward-owners-baseline.json),
[selection](spell-reward-owners-tests.txt) and [manifest](spell-reward-owners-batch.json).

The fixtures register candidate function addresses into shared normalization maps.
Before retirement, an opt-in original-path probe will check raw-address lookup
hits, operation coverage and generated-ID behavior in two separate processes.
Use raw address keys: reward's normalized IDs overlap the generated-ID range,
so matching normalized numbers alone would misattribute hits. Keep raw map-size
diagnostics separate from the frozen outputs, which already account for known
address nondeterminism. Reserve removed registration slots only after this review.

GPT-6 Luna owns the isolated conversion draft. The primary owns probe design,
source review, baseline acceptance and qualification. Original source stays frozen
while any test/build runs; restore probe-only changes before accepting conversion.
No conversion is installed yet. Local artifacts: `build/port-spell-reward-owners/`.

Required converted gates: exact root-name sets in three profiles, safe/static,
three production builds/ABI checks, exact known-suite comparison, fresh headless
character creation/save/load/resume and unchanged original asset hashes. Record
actual cgo/export/header counts after qualification.


## Original address-use evidence

Two fresh default processes passed all 51 roots. Each recorded 202 outer results,
34 distinct registered addresses, every reward operation 1300–1320 and every
spell-lifecycle operation 1500–1528, and 189,482,055 normalization observations.
No registered function address was consumed by the snapshots. The nine retained
reward callbacks remain required by production regardless of that fixture result.

Address-use results, operation coverage and generated-ID counts agree between
processes. Raw map sizes differ in three records; these are diagnostics rather
than golden outputs. Preserve per-run cardinality accounting by reserving the
13 removed spell entries (existing reservation 1→14) and 12 removed reward entries
only when that fixture initializes. All frozen expectations remain unchanged.

Temporary instrumentation is restored. See [probe report](spell-reward-identity-probe.json)
and [reproducible patch](spell-reward-identity-probe.patch), applied to `279e7867`
with `git apply --unidiff-zero`. Probe generation initially rejected overlapping
edit context before installation; grouped-context replay resolved it, and both
runs used the same reviewed instrumentation. Conversion qualification is pending.
