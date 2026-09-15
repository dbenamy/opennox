# Quest journal storage and rendering

Status: **Go conversion qualified**, against committed/pushed C baseline
**2c3ea111**. Removes **339 C lines / eleven routines**, leaving **79,924 C lines /
93 files / zero reference C**. All 717 results / ten groups match unchanged.

Scope: 11 routines / 339 removable C lines: GAME1_1.c 00427490 through 004277B0,
plus client__gui__guijourn.c. This connects entry allocation/linking, first-match
updates/removals, per-player and all-player reports, mask cleanup, cached height
and rendering. The existing journal save/load serializer is an integration
caller, not part of this conversion. The NPC dialogue shell is outside this batch.

Reuse inventoryWindowOwner for actual players, GUI/font/rendering and localized
strings. Add actual player units/update records in slots 1, 7 and 31, plus an
active nil-unit slot to qualify sparse iteration. Own the real important-message
allocator/queue and inspect exact reports. Keep journal entries in the real
allocator; normalize known node identities and independently assert both links
and every unrelated player byte. No replacement serializer or test-only C
algorithm is introduced.

Initial tests cover 63-byte/NUL/non-ASCII names, full-width flags and padding,
first duplicate/head/middle/tail/missing operations, mask combinations, a
three-seed 384-step independent list model and 1,000 create/remove cycles.
Reports test local versus remote behavior and all-player iteration; rendering
covers labels, wrapping, multiline text, colors, oldest-first order and scroll
boundaries. Exact cached-height expectations use authored fixture strings and
the public renderer. The save/load integration invokes actual sub_41BEC0 with
a temporary real cryptfile, compares independent saved bytes, then checks the
reloaded list's links, order and fields. This does not establish whole-game save
byte equality.

At the frozen C baseline, all eleven roots passed; production C was unchanged. Drafts copied
from build/port-journal are now stale; do not reapply over fixes. The native draft
must wait for a qualified, frozen, committed C baseline.

The current shop-native binaries/full-asset results also qualify unchanged
journal C production. Reuse them only while production source is identical;
any justified prerequisite invalidates that reuse. Repeat new C captures in
all targets, qualify affected callers, freeze and commit before conversion.
Native qualification follows PORT.md's affected-corpus policy, with all three
production builds, interfaces, full-asset comparison and fresh inventory gameplay.
Keep the existing inventory journal-draw expectations unchanged.

Review observations: local flag updates currently do not recompute cached height
(the common update and decoder update both omit it). No bug correction is claimed.
The renderer's original 2,048-unit text scratch buffers are bounded in fixtures;
no overflowing C call is needed for a baseline. Native strings can preserve
valid layout without those fixed scratch buffers. Entries need C-compatible
storage/layout; the C load routine invokes retained remove/add owners rather
than freeing journal nodes outside those owners.

Evidence and draft qualification plan: build/port-journal/scope.json,
remaining-callers.json, qualification-plan.json, production-baseline.json and
develop-a logs/results/captures. Raw evidence remains local; this report and
frozen tests will be committed with the qualified baseline.


Development-a stopped at discovery in 82.784s: the new save/load test imported
common/flags using its path basename instead of its declared noxflags package
name. No roots ran. Add an explicit alias; production C is unchanged. The driver
and compiler jobs joined before this fixture-only correction. Development-b
retries the eleven roots with the corrected import.

Development-b passed all eleven roots in 101.554s: 717 captured results in ten
groups, with the additional lifecycle and independent contracts above. Capture
JSON was losslessly compressed and SHA-256 checked. The three-target C
qualification now repeats captures and exercises affected callers; production
fingerprints must match the preceding qualified shop conversion for evidence reuse.

## Frozen C baseline qualification

Three affected-corpus runs passed: default **245 roots / 149.627s**, server
**243 / 246.214s**, highres **245 / 171.243s**. All ten groups / 717 results
matched development-b byte for byte in each process. The selection covers
journal, client inventory/windows/entry/listbox/sliders/UI rendering/object
rendering/meters, gameplay reports/text, player controls and quest eligibility.
It includes existing journal inventory rendering expectations without changing them.
All Go/C/header fingerprints stayed unchanged during qualification. Production
fingerprints exactly match shop conversion 31ec1393, so its three production
builds/interfaces, exact full-asset failure comparison and fresh gameplay evidence
are reused for this unchanged-production C baseline. No new production build or
new gameplay run is claimed for this baseline.

Frozen expectations live in src/journal_*porttest_test.go. The adapter calls
actual C owners at this baseline; there is no reference algorithm copy. To repeat,
use the 386/SSE2/CGO environment in PORT.md, select `^TestJournal` through
`tools/porting/run_tests.py`, and run tags `porttest`, `porttest,server` and
`porttest,highres`. Optional `OPENNOX_JOURNAL_CAPTURE=/absolute/path/prefix` writes
full JSON for diagnosis. The full accumulated pattern now includes these tests.
The affected pattern is recorded above; fixture assets use the existing inventory
owner setup. Local detailed qualification: build/port-journal/c-qualification.json.

The locked focused C repeat passed all eleven roots in **22.837s**.

## Native translation

C baseline **2c3ea111** is committed and pushed. The draft Go implementation
replaces all eleven routines and removes client__gui__guijourn.c: **339 C lines**,
leaving **79,924 / 93 files / zero reference C**.
Ten C interfaces remain for the actual decoder, script and save/load callers.
Inventory drawing and the Go script caller invoke Go directly; the private draw
export and three now-unused report exports/declarations are removed. Tests invoke
the actual Go owners, while the production C save/load integration remains C.
No C algorithm is retained for testing. Entry allocation/free uses libc for the
shared layout and original nil-on-allocation-failure behavior; list ownership,
matching, mutations, reporting and presentation are Go. Direct Go string inputs
retain C's first-NUL semantics, 63-byte storage limit and byte truncation.

First native focused run: eleven roots pass in **183.705s**. All **717 results /
ten groups exactly match** the frozen C baseline, without a behavior correction
or changed expectation. The completed qualification is recorded below.

Native affected-corpus qualification passed **245 default / 155.828s**,
**243 server / 245.918s**, **245 highres / 168.299s**, each with all frozen captures
unchanged. Production default/highres/server builds passed in **58.476 / 8.548 /
55.310s**. All three are ELF32/i386/SSE2/CGO with ten retained journal interfaces,
four retired interfaces absent and no tagged test helpers. The full asset suite took **49.797s** and matched the exact known **1,553 failure
entries / 15 passing, three failing and 32 skipped packages**, exit 1. No failure
was added, removed or suppressed. Fresh inventory gameplay took **52.919s**,
exit 0, with all nine screenshots matching the committed inventory C reference
and updates disabled. All **1,627 Go/C/header fingerprints** remained unchanged
through qualification.


## Recovery and limits

The C baseline commit contains the fixtures and frozen expectations. The Go
conversion retains those expectations and the actual production serializer test.
Reconstruct detailed captures with `OPENNOX_JOURNAL_CAPTURE` and the tagged test
driver described above. Rebuild production targets without `porttest`; expected
journal C interfaces are the nine storage/report entry points in GAME1_1.h plus
nox_xxx_cliBuildJournalString_469BC0. The draw export and three journal report
exports must be absent. Use the tracked inventory-window.yaml and its recovery
instructions for fresh headless integration. Detailed local evidence is
build/port-journal/native-qualification.json; raw captures are losslessly gzip
compressed with checked SHA-256 manifests.

Gameplay qualifies the existing nine-screen inventory route, including journal
presentation. Authored owner fixtures cover storage/report boundaries and the
real journal save/load serializer; whole-game save byte equality is not claimed.
The existing local flag-update height behavior is preserved. No production
prerequisite or native behavior correction was needed. The one development
failure was the fixture import alias described above.
