# Briefing window lifecycle and transitions

Status: **C baseline qualified; expectations frozen** after qualified/pushed briefing presentation
**adfe6fa5**. Production unchanged: **78,977 C lines / 92 files / zero reference C**.

Scope: nine remaining GAME2.c routines, about 320 physical C lines: voice-start
callback, construction, input/background callbacks, draw dispatcher, scrolling
speed/completed-draw helpers, modal presentation and cleanup. Caller audit retains
the display bridge used by client__network__cdecode.c; constructor, lifecycle and
private drawing callbacks can use direct Go calls. Across briefing presentation
and window code, five C entries remain and twelve private interfaces can retire.

Reuse the briefing fixture's actual renderer, GUI parser/widgets, sprite, music,
string and player owners. Add a real dialogue owner with bounded private state:
voice queueing and clearing run without pumping physical audio. Preserve the
existing C baseline and all frozen presentation expectations.

Development-a passed ten selected roots (including the existing 100-cycle window
construction test) in **191.775s**, capturing **355 transition results / nine
groups**. Independent contracts cover all class/chapter begin/loss selections,
remembered-loss state, quest mode precedence, credits, input gates, float32
scrolling, strict elapsed-time boundaries and unsigned time subtraction, voice
pending, one-shot completion drawing, special sound identity/volume, resource
failure and cleanup. No production change or prerequisite was needed.

Added actual fade completion into the voice callback and full shipped
Briefing.wnd lifecycle through draw, fade and key handlers. Development-b is
running with `OPENNOX_BRIEFING_ASSETS`; expectations are not frozen yet.

The transition fixture uses the actual save-menu boundary in quest mode, where
it returns early. This does not establish ordinary save-selector pixels or
full-game credits playback. Ordinary chapter integration uses the tracked chapter
scenario; separately review whether additional owner setup is needed before
freezing. No replacement save, dialogue, fade or drawing algorithms are installed.

Next: join development, inspect/repeat captures, qualify affected targets and
freeze/commit C baseline. Reuse prior production builds/full-assets/gameplay only
after exact production fingerprint comparison. Finish the closed briefing family
with full accumulated native qualification at the subsystem boundary. Local scope,
caller audit, captures and unqualified drafts live under build/port-briefing-window.


Development-b passed all twelve selected roots in **27.441s** with shipped assets.
All previous 355 results repeated exactly; the two integration groups add 18
results, for **373 / eleven groups**. Three-target affected qualification is now
running. Production fingerprints must match the completed adfe6fa5 briefing
presentation before reusing its builds/full-assets/chapter gameplay. The save-menu
boundary remains outside this conversion; preserve the actual caller and document
the fixture's quest-mode limit rather than expanding into unrelated save UI.


## Frozen C qualification

Affected selection passed **251 default / 133.038s**, **249 server / 239.759s**,
and **251 highres / 155.705s**. All **373 results / eleven groups** repeated
exactly in each process; original presentation fixtures also passed unchanged.
All source fingerprints stayed identical during qualification. Production matches
adfe6fa5 exactly, so its three production builds, interface checks, full-asset
known-failure comparison and fresh chapter gameplay remain applicable.

The accumulated `^TestBriefing` selection already includes this family. Run
`^TestBriefingWindow` for its focused twelve roots, with `OPENNOX_BRIEFING_ASSETS`
set to the original Nox directory. Optional `OPENNOX_BRIEFING_WINDOW_CAPTURE`
writes the complete JSON. The final locked focused repeat passed all twelve roots in **26.047s** before
commit. All compiler/test jobs joined.
Local qualification: build/port-briefing-window/c-qualification.json. The native
window draft and apply/qualification scripts are prepared but unapplied.
