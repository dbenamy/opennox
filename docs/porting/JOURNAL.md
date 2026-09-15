# Quest journal storage and rendering

Status: **C baseline qualified; expectations frozen** after qualified/pushed shop conversion
31ec1393. Production remains 80,263 C lines / 94 files / zero reference C.

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

All eleven roots pass with frozen expectations; no production correction or
journal translation has been applied. Drafts copied
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
