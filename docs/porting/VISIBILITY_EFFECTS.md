# Visibility and effect reports

## Scope and status

Native conversion qualified: **850 physical C lines removed**, leaving **56,118
lines / 82 production files**, zero reference C. The original 849-line block plus
its trailing separator is replaced by five Go files (509 lines): 29 live behaviors
and retirement of one proven orphan. Six C exports remain; 24 are retired.

Default/server/highres each pass **383 affected roots / 10,772 leaf cases**, no
skips. All **8,962 records / 65 groups** match unchanged live expectations,
including supplemental scan-delay contracts. All three production builds and ABI
audits pass. Static checks pass. The full suite exactly retains 1,553 known failure
entries (15 packages pass / three fail / 32 skip). Gameplay, actual save/load and
flat rendering with exact map regeneration pass. Production qualification took
373.9s; source remained unchanged. Artifacts are under
build/port-visibility-effects/native-{default,server,highres,production}.
Native client SHA-256:
`58483f6c8802ebb8b80e465fdf302df1d612f2d7aca8d10774f48f4d86add3f7`.

## Baseline and independent contracts

Use real sparse player lists, observer camera objects, outbound message queues,
reliable-message pools and registered server/object owners. Compare packet bytes,
recipients, ordering, returns and state; do not replace lookup/delivery behavior.
Exercise strict viewport boundaries with adjacent float values, explicit float
stores, signed coordinates and low-word serialization. Perception needs full
seen-enemy collections, removal/replacement order, team/ray visibility, timing,
RNG consumption and callbacks. Freeze repeated C captures before translation.

Retain creature/item/common serialization checks as affected consumers. Qualify
all three targets and production builds, exports, known-failure suite comparison,
headless gameplay, save/load and flat rendering at the completed batch boundary.

## Decisions and issues for review

- Existing Go viewport delivery appears equivalent but must pass paired contracts
  before reuse. Existing vampire-effect Go and C entry points have different
  coordinate domains: C serializes and culls using unsigned destination low words;
  generator-spawn culling uses full signed destination coordinates. Preserve each
  entry point's behavior rather than merging by similar API names.
- Reuse the established real reliable-message pool owner for summon reports;
  these tests enqueue engine messages locally and need no network connection.

- The mana throttle compares an unsigned stored last value with a signed `short`
  current value. Equal values >= 32768 still enter the changed branch; the delta
  threshold or elapsed timer then decides whether to report. Preserve this observed
  promotion behavior. Independent full-word tests exposed it before baseline freeze.
- Baseline fixture setup initially omitted the reliable allocator's handle arena.
  Reused the isolated handle owner; no production correction was needed.

## Development evidence (before freeze)

Initial viewport/camera contracts passed 570 cases, including paired calls to the
existing server Go viewport method. Extended fixtures cover all thirty C entries:
point/destination/prediction/shield/chat packets, reliable reports, sparse and
observer players, full-word health/mana throttles, dirty flags and recipient gates,
seen collections through capacity sixteen, preferred/nearest target selection,
actual script callbacks, facing/invisibility gates, sight maintenance, real spatial
scans with seeded RNG and forced deadlines, global removal and destroy reports.

Development logs/captures are under build/port-visibility-effects/development*.
The expectations are now frozen. Final default/server/highres and repeat each
pass **84 roots / 11,688 leaf cases**, no skips. All **11,458 records / 65 groups**
match, including 4,908 existing serialization records. Static checks pass.
Production builds and all thirty C-symbol checks pass. The full suite exactly
retains 1,553 known failure entries (15 packages pass / three fail / 32 skip).
Gameplay matches 41 frames, actual save/load seven, and flat rendering fourteen
plus exact map regeneration. All final gates share **1,954 source-file
fingerprints**, unchanged. Artifacts: build/port-visibility-effects/
c-{default,server,highres,repeat}-final and c-production. Production qualification
took 223.6s; baseline client SHA-256 is
`b655cf4df17636626c104ef39a67d7299d034e0a8b64f44a2ae061586b9d82a3`.

Fixture corrections before freeze: prediction reads Float28 at byte 112, not
SpeedCur; player fixtures used in enemy tests need a valid type ID (zero aliases
missing fish-type IDs). An intermediate explanation confused ClassVisibleEnable
with ClassMonsterGenerator; the constants were checked and separate independent
cases now test visible-state reports and actual generator enemy classification.
These corrections changed no production behavior.

The first frozen default/server/highres sweeps passed all 84 selected roots and
all 65 capture groups (11,458 records: 6,550 visibility/effects plus 4,908
serialization consumers). Shared-state review then added explicit save/reset/restore
for the two zombie-type lookup caches used by sight/sound gates. The final sweeps
qualified this fixture-isolation correction with unchanged expectations.

## Reachability decision for native conversion

A second complete reference audit confirms `sub_528030` is orphaned: the only
production references are its own definition and header declaration. Its only
caller is the newly added primitive baseline fixture; no dynamic symbol lookup
path exists in this source tree. Retire its 47-line C block, declaration and
fixture rather than retaining an unused Go implementation. The committed C
baseline will preserve its 3,360 historical records for recovery; native checks
will explicitly omit that retired group's hash, keeping all live groups unchanged.
The signed-mana finding therefore describes an orphaned legacy helper, not an
observed player-facing problem. This reversible decision needs no user input.

Process improvement: perform this whole-repository reachability distinction
before writing new baseline fixtures. No-external-callers alone does not prove
unreachability when helpers are called inside the selected block.

## Frozen capture groups

| Group | Records |
| --- | ---: |
| camera | 40 |
| candidates | 40 |
| chat | 106 |
| destinations | 48 |
| destroy-reports | 225 |
| frame-copy | 16 |
| global-removal | 4 |
| health-throttle | 3360 |
| killable | 50 |
| point-packets | 240 |
| prediction | 64 |
| reliable | 24 |
| seen-cooldown | 84 |
| seen-insertion | 102 |
| seen-removal | 884 |
| shield | 265 |
| sight-maintenance | 8 |
| spatial-scan | 12 |
| special-updates | 108 |
| target-selection | 340 |
| viewport | 530 |

Expectations are in the fixture files and batch manifest. Native conversion will
retire only the explicitly orphaned health-throttle group; no live expectation
will be regenerated.

## Supplemental original-C coverage

Native review added the no-enemy scan-delay branch. Supplemental baseline
**e2616865** preserves the original C and adds 864 cases; both separate C runs
pass. Real circular bounds correct an origin-only fixture setup; the original
spatial-scan capture stays unchanged. See [VISIBILITY_SCAN_DELAY.md](VISIBILITY_SCAN_DELAY.md).
Native qualification now expects **8,962 records / 65 groups** and selects
**383 affected roots**. All live historical expectations remain unchanged.
