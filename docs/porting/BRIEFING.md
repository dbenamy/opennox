# Quest briefing presentation

Status: **Go conversion qualified**. C baseline **d1565ae5** is pushed.
The conversion removes **947 C lines / thirteen routines**, leaving **78,977
production C lines / 92 files / zero reference C**. All **467 frozen results /
eleven groups** match on default, server and high-resolution targets.

Scope: thirteen routines / 947 removable C lines: client__gui__guibrief.c and
GAME2.c helpers sub_44E110, sub_450960, sub_450AD0, sub_450AF0,
nox_gui_setQuestStage_450B00 and nox_gui_getQuestStage_450B10. Remaining briefing
window construction, presentation, fade and cleanup owners are integration callers.

Reuse inventoryWindowOwner for the actual renderer, GUI parser/widgets, player
records, font/string manager and sprite allocation/deletion owners. Give briefing
a full 640x480 or 800x480 pixel buffer, independent chapter strings/voice names,
a real music owner and deterministic clock. Save/restore all touched direct words
and mapped regions; preserve outside screen-particle/audio owners. Use actual
briefing construction/cleanup and inspect real captured-window state and fade
activation. No replacement decoder or test-only C algorithm is added.

Initial eleven roots cover all 66 class/chapter begin/loss combinations and
credits, 32-bit unsigned score comparison, caption/image fallback/stage helpers,
twelve sprite-cache allocation/reuse/cleanup cycles, 100 window lifecycles,
full-screen statistics/title/instructions rendering and all three frame/30 blink
state machines. Reports cover all 64 participating-player subsets packed as the
real producer does, independent exact field/order assertions, missing/duplicate
IDs and bounded sparse input, plus positive presentation and image/caption choices.
Development-e captured 461 results / ten groups; see qualification status below.

The server's sub_4D6770 advances its output slot only for participating players,
up to six entries. This establishes the packed-record invariant used by the
client qsort(count); sparse authored input is not evidence of a production bug.
Preserve it in C captures without changing supported behavior speculatively.

The development notes below describe fixture preparation before freezing.
No production prerequisite was needed. Local scope/caller notes are under
ignored build/port-briefing; committed tests preserve the frozen expectations.

Development-a stopped at discovery in 85.532s: an image handle was passed to a
pointer-only fixture helper. Use the handle’s numeric conversion. No roots ran;
all compiler jobs joined before this fixture-only correction. Production unchanged.

Development-b compiled and completed all eleven roots in 106.078s, but every
root stopped at the inherited fixture’s client/server type-index check. The
briefing extras duplicated QuestGoldChest and QuestGoldPile already registered
by inventoryWindowOwner. Filter those two extras; no production change.

Development-c passed title and statistics rendering, then stopped at the first
instruction draw with a null sprite draw callback (24.827s). Fixture-only extra
types lacked the static-image callback/data installed for the inventory’s base
items. Install that same actual production drawing callback and owned image data
for all briefing sprite types. Preserve the two completed captures compressed;
no baseline expectations are frozen and production remains unchanged.

Development-d stopped at discovery (1.318s): the new image-data word needs the
handle-to-uintptr conversion before uint32, as in the reused fixture. Corrected
after the driver joined; no roots ran.

Development-e passed all eleven roots in **25.876s**, capturing **461 results /
ten groups**. Added a twelfth root for shipped Briefing.wnd geometry, actual
modal presentation and all three drawing modes (six further results), using
`OPENNOX_BRIEFING_ASSETS` for original resources. Development-f runs with assets.

Development-f passed all twelve roots with shipped assets in 25.882s. Comparing
its repeated groups to development-e identified exactly two unnormalized pointer
classes: the interned empty UTF-16 caption (return plus caption field) and the
known mapped empty caption buffers. Other groups matched. Record those specific
owned identities and repeat; no expected hashes are frozen or behavior changed.

Development-g passed all twelve roots in **27.324s**, with **467 results / eleven
groups** after precise empty-caption identity normalization. Three-target affected
qualification is running with shipped assets. A fresh ordinary new-game chapter
scenario is being captured using unchanged qualified production; inspect its
chapter frames and repeat with updates disabled before claiming integration.

The first new-game capture completed, but both proposed chapter checkpoints
showed the captain scene on visual inspection. They do not establish briefing
pixels. The E2E harness intentionally reduces briefing duration to 10ms and fade
duration to ten frames; sample earlier frames without changing production timing.


## Fresh gameplay reference

The earlier-frame scenario completed and was visually inspected: frame 10 shows
chapter art/title text, frame 20 shows the chapter artwork, and frame 1 is the
loading transition. All eight checkpoints then matched in a fresh repeat with
NOX_E2E_OVERRIDE=false, Xvfb and null audio, using unchanged C briefing production
from 78ad58b2. The tracked briefing-chapter.yaml and briefing-chapter-pixels.json
preserve the route and decoded NRGBA hashes. The manifest hashes the reference
files; the separate fresh repeat supplies execution evidence. This covers the
ordinary warrior chapter transition and later gameplay, while quest statistics
and instructions are covered by actual-owner and shipped-window fixtures.


## Frozen C qualification

Affected-corpus qualification passed **240 default / 127.290s**, **238 server /
232.609s**, and **240 highres / 149.445s**. All **467 results / eleven groups**
matched development-g in each process. The selection includes briefing, journal,
client inventory/windows/entry/listbox/sliders/UI rendering/object rendering/meters,
gameplay reports/text, screen effects and particles. Production fingerprints are
identical to 78ad58b2, so that commit's three qualified production builds, exact
full-asset failure comparison and inventory gameplay remain applicable; the fresh
chapter capture/repeat supplies additional integration evidence. All source
fingerprints remained unchanged during these runs.

The accumulated test pattern now includes `^TestBriefing`. To repeat focused
checks, use the PORT.md 386/SSE2/CGO environment with the tagged test driver and
`^TestBriefing`, setting `OPENNOX_BRIEFING_ASSETS` to the original Nox data tree.
Run tags `porttest`, `porttest,server` and `porttest,highres`. The optional
`OPENNOX_BRIEFING_CAPTURE` prefix writes complete JSON. The remaining C boundary
will retain eight exports and retire five private helper exports/declarations;
see local remaining-callers.json for the audit. No reference C copy is introduced.

The final locked C repeat passed all twelve roots with assets in **25.714s**.
All jobs joined before baseline commit; production remains unchanged.


## Native translation

C baseline **d1565ae5** is committed and pushed. Replaced all thirteen scoped
routines with Go chapter/score structs, report decoding, sprite-cache management,
and title/statistics/instructions drawing. Eight C exports remain for actual
callers; five private helpers and their declarations are removed. No test-only C
algorithm remains. Applied size: **78,977 production C lines / 92 files**,
**947 fewer lines**, with zero reference C.

The first native run executed all twelve roots and matched **all 467 results /
eleven groups**. Its overall driver failed at Go vet: a zero-argument localized
`fmt.Sprintf` call used a dynamic format string. Joined every compiler job, then
introduced a variadic localization-format helper and kept the stage label as a
literal `%s` argument, matching the original concatenation. Frozen expectations
are unchanged. The corrected focused run subsequently passed; do not treat the initial
root passes as completed-batch qualification.

Preserved details include unsigned descending totals and observed stable ties,
UTF-16 code-unit name clipping, the existing `XX1` statistics label, the SoulGate
paragraph's height-based right edge, and instructions measuring the prompt with
the cached font while drawing with the window font. Existing C window creation,
modal presentation and cleanup remain real integration callers. All sprite draws
still use their actual callbacks and coordinate owner.


The corrected focused run passed all twelve roots in **184.274s**, with all
467 results / eleven groups identical to the frozen C baseline. The initial
run took 188.455s including its vet failure. Completed qualification is recorded below.
Old derived Go cache entries were pruned only after all readers joined: 166
entries / 3,013,755,798 bytes, recorded in the local cache-pruned manifest.
All capture evidence is preserved with verified lossless compression.


## Completed native qualification

Affected tests passed **240 default / 129.906s**, **238 server / 222.961s**, and
**240 highres / 143.610s**. Every variant reproduced all **467 results / eleven
groups** exactly. Production builds passed in **58.798s / 8.524s / 56.684s**
(default / highres / server); ELF32/i386/SSE2/CGO checks confirm eight retained
C bridges, five retired interfaces absent, and no tagged test helpers.

The full default asset suite completed in **50.660s** with exactly the established
**1,553 failure entries** and **15 pass / 3 fail / 32 skipped packages**. No failure
entry was added or removed. Fresh chapter gameplay completed in **37.462s**,
matching all eight checkpoints against the original C reference with updates
disabled, seeded randomness, fresh copied data/save, Xvfb and null audio.
All **1,636 Go/C/header fingerprints** stayed unchanged during qualification.
Every test, compiler and headless job joined before commit.

Local evidence: build/port-briefing/native-qualification.json, native-source-fingerprints.json,
native-qualified-capture-compression.json and native-binary-verification.json.
These artifacts are ignored; frozen tests, C baseline revision and the tracked
chapter scenario/pixel manifest preserve recovery. The next connected scope is
the remaining briefing window lifecycle and transition owners in GAME2.c.
