# Quest runtime, statistics and difficulty scaling

## Corrected C baseline

Following qualified world-collision conversion `2e19c78d`, this connected batch
contains **56 live functions / 894 physical C lines**: GAME3_2.c 004D6000 through
004D7B40 (656 lines / 44 functions) and GAME3_3.c 004E3CA0 through 004E4100
(238 lines / 12 functions). The unrelated player-list file parser is excluded.
This C baseline had **53,946 C lines / 82 files**, with zero reference C.

The final source review found that the first scoring fixture did not initialize
its exponent table. Its zero exponent masked an original-C width bug. The shipped
blob contains an eight-byte double near 1.9 at 0x581450+10088, but sub_4D66E0 read
it as long double. The compiled original uses fldt, reading an invalid extended
value. With shipped data, an independent stage-two/ten-point test returned
**2,147,483,648 instead of 37**. That result could then saturate the reported score.

Correct the C read to getMemDoublePtr and make every quest fixture install and
restore the shipped bytes. This is a reversible prerequisite correction made
before any Go translation. The corrected C passes all 24 roots; only the three
scoring/scoreboard groups need corrected expectations. The other 20 groups remain
identical. The initial zero-exponent gates are superseded and preserved locally;
no baseline containing them was committed.

The production-identity reuse audit is withdrawn because production C changed.
Fresh corrected-C default/server/highres and independent repeat each pass
**24 roots / 8,501 unique leaves**, no skips, **23 groups / 8,512 frozen records**.
Wall times are 36.48s / 122.66s / 44.42s / 8.14s. All four gates and fresh production
share **2,033 source fingerprints**. See c-fixed-coverage-audit.json.

All three production builds and ABI audits pass, including the 56 original quest
interfaces. The full asset suite exactly matches the known 1,553 failure entries
and package outcomes (15 pass / 3 fail / 32 no-test). Fresh gameplay (41), save/load (7), and flat rendering (14) frames pass against the preceding world-collision references; flat rendering
also verifies exact loaded-map regeneration. Production qualification takes 374.47s.
Client SHA-256: `cb7a26b46cfc040aa9fff33e1ddbf039a334620cc0c2e6c0c305743e9cebdcf8`.
Evidence is under build/port-quest-runtime/c-production and c-fixed-*.

The final caller audit also identifies nox_float2int16_abs (three physical lines
in GAME1_1.c) as private to the selected health scaler. Its original two-root
conversion suite passes default/server/highres in 6.12s / 5.66s / 7.05s: 1,397,827
float inputs plus 12 PC/RC combinations. Its C symbol is present in all production
binaries. Remove that helper with its final callers instead of retaining it solely
for tests. Total native scope is therefore **897 C lines / 57 functions**; keep the
existing independent Go absolute-conversion assertions when removing its C oracle.

## Contracts and findings

Fixtures reuse the real server/player roster, object factory/list, map index,
health and type metadata, gameplay balance lookup, observer/movement callbacks,
deferred audio and reliable message queue. Eight C-owned player objects exercise
the six-entry scoreboard and capacity boundary. Pointer returns are normalized
only to explicitly identified fixture objects, players or settings records.

- Statistics cover resets, dirty bits, wraparound, participation checks and
  reconnect/departure cleanup across all 32 slots. Non-target fields are checked.
- Scoring has exact rational stage-one weighting checks, broad original-C
  snapshots, zero/cap invariants, single-player composition and multiplayer
  aggregation, including inactive/nonparticipant distinctions.
- Scoreboard fields use distinct values to check ordering as well as narrowing;
  the actual 90-byte payload is checked through the real queue. Stage/map messages
  cover fixed buffer limits and byte/short narrowing.
- Participant and capacity checks cover 0..8 players, exact-one versus nonzero
  participation and headless-host exclusion. Gate fixtures cover movement,
  observer/camera state, buffs, messages/audio, closure and timeout boundaries.
- Difficulty caches, float-bit getters/setters, caps and healthy/damaged/disabled
  objects are checked through actual health updates and monster history refresh.

The original health scaler compares an unsigned current-health value with a
signed 16-bit maximum. Consequently healthy values at or above 32768 skip scaling.
Preserve that established behavior; boundary cases cover 32767/32768/65535.
The char*-declared sub_4D70B0 returns a settings record, not a string. Its fallback
also clears bit 0x10 in the real settings byte 100; the fixture checks both pointer
identity and the side effect.

Fixture development corrected a tick-rate argument type, restored a temporary
generator class before its original monster cleanup, and included the existing
buff-removal audio event before the gate-return sound. Expectations for the three scoring groups are regenerated only from the corrected
C baseline. This describes the committed C baseline; native implementation follows below.
Raw diagnostics remain ignored under build/port-quest-runtime/development*.

## Frozen groups

| Group | Records |
| --- | ---: |
| book-eligibility | 56 |
| departure | 1536 |
| difficulty-cache | 360 |
| float-scalars | 195 |
| gate-close | 72 |
| gate-return | 12 |
| health-scaling | 504 |
| observer-deadline | 900 |
| participant-capacity | 108 |
| participants | 128 |
| reconnect | 27 |
| reset-all | 32 |
| scalar-state | 69 |
| score-aggregation | 320 |
| score-primitive | 1512 |
| scoreboard | 48 |
| settings-fallback | 10 |
| slot-mask | 1280 |
| small-messages | 63 |
| soul-timeout | 240 |
| stage-messages | 216 |
| statistics | 504 |
| warp-admission | 320 |

## Qualification and reproduction

Source build/baseline/env.sh. Use tools/porting/run_batch.py with
docs/porting/quest-runtime-c-batch.json and c-default/c-server/c-highres/c-repeat
or production for corrected C, always with a fresh output directory. Use
quest-runtime-batch.json with native-default/native-server/native-highres and
production after conversion, comparing the corrected-C scenarios. Never edit source during gates.
The one-shot freeze script and installed fixture drafts are consumed.

The interface audit identified 24 entry points needed by remaining C and 33
retirements including the private absolute converter. No interface was retired
in the C baseline. Native ABI requirements and affected callers are qualified below.

## Native implementation

Five Go files implement statistics/scoring, gates/timers, state/messages, difficulty
and health scaling, and thin C exports. Production Go callers now invoke Go directly.
Twenty-four exports remain for C callers; 33 private interfaces are retired, including
the last-owner absolute converter. Working C size is **53,049 physical lines in
82 files**, zero reference C: **897 lines / 57 functions removed**.

The health-scaling export now returns void. Every remaining C and Go caller ignores
the old decompiler return, which was a mixture of temporary pointers/class/health
values. This removes an unused return contract; scaling mutations remain covered.
The signed-short maximum-health comparison, float32 intermediate stores, score
aggregation over all roster units and participant equality/nonzero distinctions
remain intentional compatibility behavior. Stage names retain the valid 31-byte
name plus terminator contract in each 32-byte field.

Static mapped-memory checks pass. Native discovery needed the flags
package's actual Go name (`noxflags`) and explicit player-index/C-object conversions.
The third focused run passes all 24 roots / 8,501 leaves and all 23 groups / 8,512
frozen records (127.98s), without behavioral fixes.
No frozen expectation has been changed for the Go conversion.

Native default/server/highres each pass **470 roots / 49,254 leaf cases**, no skips,
and **149 groups / 50,640 frozen records**, with identical 2,038 source fingerprints.
Times: 251.85s / 361.20s / 297.77s. Evidence is under
build/port-quest-runtime/native-default, native-server, native-highres and
native-coverage-audit.json.

Fresh production qualification passes (397.78s): all three 386/SSE2/CGO
builds and ABI audits; the exact known full-suite result (1,553 failure entries,
15 passing / 3 failing / 32 no-test packages); gameplay (41 frames), actual save/load
(7 frames), and flat rendering (14 frames) with exact compressed-map regeneration.
Production and affected tests share the same 2,038 source fingerprints.

Artifacts: build/port-quest-runtime/native-production and
build/baseline/runs/quest-runtime-native, quest-runtime-save-native,
quest-runtime-flat-native. Client SHA256: `5d157465f395b6c8644f2b23056265239edcdad1f4426915e59335bb55b57dc8`.

The native implementation is **558 Go lines in five files**, replacing 897 physical
C lines. No frozen expectation changed during conversion. No new user decision
or unresolved test failure was introduced.

Completed C and native scenario asset copies were deduplicated after qualification,
reclaiming 1,660,044,319 bytes per set. Ignored restoration manifests and all changed
saves, screenshots, result files and original assets remain available.
