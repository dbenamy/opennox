# Equipment modifier callback identities

Scope: retire 49 C export bridges around equipment effects. Forty addresses are
registered production callbacks (17 damage, seven defend, four update and twelve
engage/disengage); nine additional exports are used only by fixtures. Preserve
all eight callback fields in the 144-byte ModifierEff record, parser behavior,
callback order, scalar widths, mutation and all existing frozen captures.

Use distinct process-lifetime Go byte keys and typed dispatch for the three-,
five- and six-argument families. Keep separate result/discard paths and exact C
fallback conventions for unknown addresses. Migrate nine production dispatch
sites, including the raw offset-88 inversion route, all getters/fixture routes,
and six blob aliases. Retain the nine fixture-only distinct tokens where captures
or comparisons observe identity. External native bindings are unchanged.

## Original contracts and testing

Seven new porttest files add five roots. Independent C observers exercise 840
argument/result cases, including nil combinations, signed 32-bit words, scalar
writes and surrounding canaries. The actual item-update owner covers 1,296 class,
equipped-state and callback-slot combinations. Fire hooks are replaced repeatedly
across 96 calls to verify late binding through the real five-argument production
route. Registry contracts check all 40 keys in foreign records across GC, and
parser contracts exercise 93 compatible field routes, valid/invalid values and
unchanged surrounding fields. Existing effects, equipment, damage, controls,
collisions, inventory, transfer, update and resource contracts complete the
331-root selection.

Production source is identical to qualified duration commit d5d80c42; reuse its
production evidence. Require exact focused root sets in default/server/highres.
After conversion add an independent native-dispatch contract and run all focused
profiles, the complete default corpus, safe/static checks, three production/ABI
builds, the exact known asset-suite comparison and a fresh headless character
creation/save/load/resume scenario. Preserve original assets and frozen captures.
Independent test-only C observers qualify foreign callback fallback; they do not
retain engine algorithms.

## Review decisions

- equipmentEffects returns the last pointer or callback word, but all eight
  production callers discard it. Preserve configured raw observer results;
  native formerly-void owners use unused zero results.
- Fire callbacks historically declare four arguments but production dispatch
  supplies five on the qualified 386 target. Test that actual route and preserve
  the ignored fifth buffer; changing the test to four arguments loses coverage.
- Preserve blob alias initialization and both existing tables' numeric offsets.
  Do not reset global blob data inside new isolated registry fixtures.
- Remove the unexported, unreachable nox_xxx_useByNetCode_53F8E0 shim and its unused
  effectsMod helper when retiring effects_exports.go. Neither counts among the
  49 exports. Whole-source search found only the orphan's definition.

Luna audited reachability and drafted registry/parser contracts; primary checked
all 359 unique reference lines and nine production dispatch sites. Primary caught
an omitted raw offset-88 route and strengthened the lifetime test to actually
store keys in foreign records and read them again after GC, plus parser canaries.
Luna reviewed the primary observer/update/hook fixtures. Implementation is drafted
separately while original tests run; primary owns integration and acceptance.


## Frozen original result

All 331 roots pass in each default/server/highres profile without skips. Exact
run/pass name sets match the independent selection and source fingerprints match
throughout. Only the seven new porttest files differ from qualified d5d80c42;
production qualification is reused from that revision. No runtime baseline
fixture correction was required.

[Original baseline](modifier-identities-baseline.json),
[test selection](modifier-identities-tests.txt),
[qualification manifest](modifier-identities-batch.json).
Local evidence: `build/port-modifier-identities/baseline-original/`.


Disk housekeeping after the baseline reclaimed 559,857,664 allocated bytes from
1,654 verified original-asset copies in the completed duration scenario, and
764,968,960 bytes from 14 obsolete root/legacy Go cache archives predating qualified
d5d80c42. Original assets, saves, binaries and current modifier caches remain.
Host-use/hash checks and recovery manifests are recorded in PORTING_STATE.md.


## Conversion review before qualification

The installed conversion changes 32 source paths, retiring exactly 49 exports.
All retained export signatures/bodies and existing assertions/captures are
unchanged. Primary independently checked all 40 registration mappings, nine
production dispatch sites, six blob aliases, 39 sparse C fixture identities,
three existing native item-use identities and 37 C fixture operations.

Luna supplied 28 changed paths and two deletions. Primary caught and corrected
a missing signed 32-bit fixture argument conversion and three fixture routes
that read scalar data before checking early-return guards. Production adapters
already had those guards. Primary added the native API/784-case dispatch contract,
used semantic effect names for the keys and formatted the reviewed source.
Final qualification results are recorded below.

An existing no-op benchmark was captured from the frozen original binary with
GOMAXPROCS=2, three 100ms samples per callback and allocation reporting. Compare
against the converted binary in this VM only; this measures callback overhead,
not game-frame performance. Local evidence: baseline-benchmark.json/txt.


The converted no-op spot check reduced the nine modifier callbacks' median costs
from about 230–350 ns/op to 18–41 ns/op, with zero allocations in both versions.
The unchanged duration control measured 26 versus 36 ns/op, illustrating the noise
in these short VM samples. This confirms a reduction in the measured callback
boundary cost, without establishing an overall gameplay speedup. Exact samples,
commands and log/binary hashes are in `build/port-modifier-identities/` under
`benchmark-comparison.json`, `baseline-benchmark.json` and the two benchmark logs.


## Qualified conversion

All 49 selected export bridges are retired. Forty registered modifier callback
identities now use stable Go keys with typed dispatch; nine fixture-only identities
remain distinct test-only keys. All nine production dispatch sites and six blob
aliases migrate with their consumers. Original names, parsers, record layout,
sparse fixture IDs, callback guards, signed widths and late-bound fire hooks remain.
Unknown C callbacks retain their exact three-/five-/six-argument fallback paths.

All 332 focused roots pass in default/server/highres, with no skips and exact
expected name sets. The full default corpus has 2,457 passes plus the established
TestMapPopulationPrerequisiteProbe skip, matching all 2,458 independent expected
names. Safe/static checks, three production/ABI builds, the exact known asset-suite
outcomes and fresh headless character creation/save/load/resume pass. All 1,654
original asset hashes remain unchanged; accepted phases have identical source
fingerprints. Existing assertions and captured expectations are unchanged.

Selected production cgo files: 195→193; legacy export bridges: 686→637.
The 157 tracked headers contain 3,394 physical lines. Standalone production
and test-reference C remain zero. Immediate-phase net progress is 270/463
cgo files and 1,253/1,890 legacy export bridges eliminated. External native-library
bindings remain unchanged.

[Qualification](modifier-identities-qualification.json),
[updated inventory](modifier-identities-inventory-after.json).
Original baseline: 09f15464. Local evidence: `build/port-modifier-identities/`.
