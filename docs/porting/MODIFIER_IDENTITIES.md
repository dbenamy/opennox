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
