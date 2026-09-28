# Native damage and monster callback dispatch

## Scope and producer review

Retire the foreign fallbacks for object damage and monster strike/death/dead
callbacks. Migrate five C observer fixtures, including their shared collision,
object-update and use routes. Remove the two now-unused generic dispatchers
`CallIntUPtr5` and `CallIntPtr`. Keep all native algorithms and frozen assertions.

All eleven production damage names register native exact-result and boolean
handlers. The default key, type parser and object constructor select/copy those
identities. Three foreign damage observer families remain in fixtures; aliases
through the shared target object reach combat, spells, effects, resources, object
state and projectile tests. Migrate those observers together. A raw offset audit
found no additional damage-slot writer. Monster tables/loaders select 26 native
handlers; the combat/lifecycle fixtures supply the remaining foreign keys.

The strike observer also occupies collision slot 696 in shield fixtures. The die
observer also occupies the object-update slot. Register both routes explicitly.
The lifecycle use observer has exact-result and boolean routes. Fixture-only
registration helpers populate dispatch maps without adding production parser
names. Initialize the monster handler map before init functions, then insert the
26 unchanged production entries without overwriting fixture registrations.

Preserve damage argument evaluation before callback lookup, typed nils, full
signed 32-bit results, nonzero boolean conversion, and all native KeepAlive calls.
Boolean-only test overrides remain independent from exact-result registration.
AI hit capture stops recording at capacity but still performs the requested
monster-definition mutation. Reset does not clear unused capture words. Lifecycle
capture keeps its original count/reset/word order; out-of-range writes now panic
instead of overrunning the C array, outside qualified captures. Private food-search
adapters have only two local callers; direct calls preserve signed32 pointer results.

Unsupported damage/monster keys panic after this producer audit. Other callback
families retain their existing fallback behavior. These are scoped, reversible
invalid-input decisions, not a claim to preserve calls through arbitrary addresses.

## Baseline and tests

Use the preceding qualified GUI source and broad owner evidence only after exact
source, supplemental source, environment, binary, log hash and discovered-name
checks. Run twelve focused original contracts freshly in each profile and four
asset-free damage forwarding contracts with safe allocation. The broad selection
is 2,482/2,471/2,482 roots, with the same eight explicitly excluded compiled roots;
it is not the complete root corpus or a claim that all exclusions are known failures.

Focused contracts cover damage forwarding and callback selection, exact values
under boolean override, combat/lifecycle shared observers, strike mutation,
poison/cloud behavior, collision dispatch and monster table loading. Frozen
assertions/captures remain unchanged. After conversion rerun the whole combined
selection, focused repeats, safe contracts/build/static checks, three production
builds/ABI, exact known-suite outcomes and fresh save/load.

## Review and delegation

Ten-file conversion, including one new porttest-only registration helper. Function
review: 17 added, 14 changed, four removed, 60 unchanged. All 26 monster handlers
and unchanged dispatcher bodies compare exactly; generated output matches the
edited generator. Primary implementation/review; Luna quota remains unavailable.
Ignored artifacts: `build/port-damage-monster-dispatch/`.

## Local recovery

Removed seven superseded section/particle binaries after committed-source,
qualified-replacement, hash and host-use checks, reclaiming 387,035,136 bytes.
Rebuild `6f03cc39` if those historical outputs are needed. Logs and metadata remain;
exact paths/hashes are in this batch's `section-cleanup-approved.json` and journal.
The replacements are qualified GUI binaries at `0cf5064c`.

Original baseline accepted at qualified GUI source `0cf5064c`: all exact-source
reuse checks passed, plus twelve fresh focused contracts in each profile and four
safe contracts. Conversion is qualified; accepted baseline commit `0e4c2fb2`. See [baseline](damage-monster-dispatch-baseline.json).

During native qualification, removed seven superseded duration/audio binaries
(387,031,040 allocated bytes) after committed-source, replacement, hash and host-use
checks. Active tests use separate new binaries. Rebuild `d3d759ba` for the old
outputs; original GUI baseline logs/metadata remain. Records:
`duration-cleanup-approved.json` and journal in this batch directory.

## Qualification and progress

Native broad selections passed **2,482/2,471/2,482** roots without skips. Twelve
focused contracts passed freshly in each profile before and after conversion;
four safe forwarding contracts passed on both versions. Native preflight,
safe build/static checks, three production builds/ABI, exact known-suite outcomes
and fresh save/load passed. All 1,654 original asset hashes remain unchanged.
Retired dispatchers/observers are absent; native handlers, other embedded bodies,
assertions and captures are unchanged. No conversion correction was required.

Embedded production C bodies: **10→8** (71/79 retired). Fixture C imports:
**88→83**. Production cgo remains **5 client/highres, 6 server**, across three
project packages. Selected legacy exports remain zero, headers 157 files /2,731
lines, standalone C **0 production /0 test**. Evidence:
[baseline](damage-monster-dispatch-baseline.json),
[qualification](damage-monster-dispatch-qualification.json),
[inventory](damage-monster-dispatch-inventory-after.json).

After qualification finished, 1,654 identical scenario asset copies were verified
against the original assets and checked for host process use, then removed,
recovering 560,435,200 allocated bytes. Original assets, outputs, binaries and
qualification records remain. Restore the scenario inputs with:

```sh
python3 build/port-artifact-cleanup/restore-recent-scenario.py build/baseline/runs/damage-monster-dispatch-save/deduplicated-assets.json
```
