# Remaining fixture export bridges

Status: original baseline accepted; conversion draft reviewed but not installed.
Source baseline is `b2597f97`; no original test assertions or captures changed.

Seventy selected engine C exports have no production C consumer. Sixty existing
Go adapters are fixture-only and move under `porttest`; ten have real Go callers
and keep their bodies without C exports. The ten live production C callback
addresses remain for a later batch. Existing Go bodies/signatures, native owners,
layouts and external native-library bindings stay unchanged. Transitional C scalar
types in fixture adapters remain for subsequent cleanup.

Whole-source review checks actual Go calls/addresses, import-C preambles, headers
and macro aliases. The draft retains protection-set argument narrowing, temporary
update argument order/field offsets, real callback17 and independent C observers.
Five unused temporary callback addresses become nil while preserving numeric
capture IDs; the fixture skips nil identity registration to preserve zero.
The declaration review caught thirteen leftover pointer-returning declarations
before installation. Two C-int size assertions move with their fixture adapter.

The original full default corpus passed2,478 roots with only the established
`TestMapPopulationPrerequisiteProbe` skip. Each of305 newly selected roots passed
twice per profile (the full default supplies its first run). Prior same-source
book coverage of256 client/high-resolution and254 server roots is reused with
exact source/binary/environment verification. Combined focused selections are561
client/high-resolution and559 server roots. Broader indirect fixture sharing is
covered by the full default corpus, not claimed from the focused selection alone.

Converted acceptance requires full default plus focused server/high-resolution,
safe/static, three production/ABI builds, fresh preview and final save/load,
exact known-suite results and original asset hashes. The profile controller gains
an optional default-only pattern, allowing the full default sweep without a
redundant focused-default pass. Existing behavior without the option is unchanged.

Primary owns this batch locally: Luna remains quota-limited; no substitute model.
See [baseline](remaining-fixture-bridges-baseline.json),
[manifest](remaining-fixture-bridges-batch.json), and
[focused selection](remaining-fixture-bridges-tests.txt).
Local audit, draft and logs: `build/port-remaining-fixture-bridges/`.

Disk recovery removed four superseded binaries (187,367,424 allocated bytes;
rebuild from `1238c985`) and losslessly archived30 historical capture groups
(711,507,968 allocated bytes). Exact paths, hashes and restoration instructions
are in that folder's cleanup/archive journals; current binaries and original
assets remain. No qualifying source changes were made while original tests ran.
