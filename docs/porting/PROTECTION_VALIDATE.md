# Protection buffer validation — 2026-09-10

`sub_56FB00` now uses Go for signed-handle eligibility, first-match record lookup
and checksum validation. The C export remains for the player wide-name validation
caller. Ineligible or missing handles return zero before evaluating the buffer.
A match compares its encrypted value against key XOR the existing length-aware
checksum helper. No manager, RNG or input state changes.

The same 515 cases passed against original C (`d948eaa2`) and are retained for
the conversion. They cover aligned/unaligned data, complete/trailing words, zero
and high-bit keys/handles, missing IDs and first-match duplicates. Readable input
bytes and full list/manager/RNG snapshots must remain unchanged. Nil buffers
with sizes through UINT_MAX retain the already-ported checksum's zero result.

Existing memguard PROT_NONE pages prove that ineligible and missing handles
never read the buffer, even with a huge declared size, and that a matched record
with 0–3 bytes never reads an incomplete word. These are short-circuit tests,
not permission to dereference invalid buffers for positive complete-word sizes.
The checksum oracle uses independent byte lanes.

Final original-C baseline: build/port-validate/c-before-final.log. The corpus
ensures zero keys also occur with live matching records. Production C is
**142,115 physical lines**, down **15**, in 153 files; C references **0**.
See [C_LOC.md](C_LOC.md). Final target/integration results follow.

All accumulated protection tests pass for default/server/highres on 386. All
three production targets build, with the required validation C entry confirmed
as Go-backed. Fresh `validate-port` warrior gameplay exits 0 against both
preserved screenshots with overrides disabled. Evidence: build/port-validate.
The latest full-suite comparison remains the rekey milestone (no new failures).
