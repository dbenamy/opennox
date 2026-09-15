# Briefing window lifecycle and transitions

Status: **Go conversion qualified**, after frozen C baseline **a46a62e5** and
briefing presentation **adfe6fa5**. Nine routines removed **320 physical C lines**,
leaving **78,657 / 92 production .c files / zero test-reference C**.

## Scope and compatibility

The remaining GAME2.c briefing routines are Go: voice-start callback,
construction, input/background callbacks, draw dispatch, scrolling speed,
completed drawing, modal presentation and cleanup. Go callers and window/fade
callbacks call Go directly. Across presentation and lifecycle, five C entries
remain for real callers and twelve private interfaces retire. The network decoder
still uses the display bridge.

Preserve nil GUI responses for zero callback results, shared pointer ownership,
original sprite deletion order, post-draw transition-state reads, float32 scroll
storage, strict elapsed-time comparison, unsigned timer subtraction and resource
failure ownership. No production prerequisite or behavior correction was needed.

## Baseline and independent checks

Fixtures reuse actual renderer, GUI parser/widgets, sprite, music, strings and
player owners. A real dialogue owner checks queued/cleared voice state without
pumping physical audio. Actual fades invoke the voice callback; actual key
handlers disable themselves after dismissal. Tests cover all three classes and
eleven chapters, begin/loss choices, remembered-loss state, quest mode precedence,
credits, input gates, timing boundaries, one-shot completion drawing, special
sound identity/volume, resource failure, cleanup and 100 construction cycles.
The shipped Briefing.wnd runs through draw, fade and key handlers.

The new **373 results / eleven groups** supplement **467 presentation results**:
**840 frozen results / 22 groups** remain exact. Development-a's 355 results
repeated unchanged when the 18 fade/shipped-lifecycle results were added.

C affected qualification passed **251 default / 133.038s**, **249 server /
239.759s**, **251 highres / 155.705s**, then a locked twelve-root repeat in
**26.047s**. Exact production fingerprints matched adfe6fa5, allowing its three
builds, full-assets and chapter gameplay evidence to be reused before freezing.

## Native qualification

Whole-family focused comparison passed **23 roots / 183.354s**, with all 840
results exact. At this briefing-family milestone, the **complete accumulated
corpus** ran in all three targets, including original briefing/shop/trade assets:

| Target | Selected roots completed | Wall time |
| --- | ---: | ---: |
| default | 974 | 551.131s |
| server | 971 | 628.189s |
| highres | 974 | 623.296s |

Each target includes the expected optional prerequisite skip; selected-root
completion is not a claim that skipped checks executed. No failures occurred.
The existing 600s per-package timeout was sufficient; the runner is unchanged.
However, highres root execution took 598.823s, leaving almost no margin. Before
the next complete-corpus milestone, add an explicit longer timeout for that
selection; affected-batch runs can retain 600s. Do not reinterpret a timeout as
a pass or discard the tests.
All jobs joined and all 1644 source fingerprints stayed identical.

Production builds passed: default **59.248s**,
highres **8.649s**, server
**55.531s**. All are ELF32/i386/SSE2/CGO,
with five retained entries, twelve retired entries and no test helpers linked.
The full-assets suite preserves exactly **1,553 failure entries**, with
**15 passing / three failing / 32 skipped packages**, exit 1. Existing failures
remain recorded; the full suite is not green.

A fresh headless chapter run completed in **38.596s** and
matched all eight tracked chapter/gameplay checkpoints against the original C
reference, with golden updates disabled and a fresh asset/save copy.

## Diagnostics and limits

Native-a stopped at discovery in 43.162s because a new cgo preamble used uint32_t
where the established map-draw flag declaration uses unsigned int. Corrected
after joining all compilers; no tests ran and no expectations changed. Native-b
and the final qualification passed. An apply-script export-removal assertion was
also corrected before compilation; it did not change an algorithm or oracle.

The transition fixture calls the actual save-menu boundary in quest mode, where
it returns early. This does not establish ordinary save-selector pixels or audible
full-game credits playback. Ordinary chapter integration is covered by the tracked
scenario; quest/credits drawing and transitions use actual owners. No replacement
save, dialogue, fade or rendering algorithms are installed.

## Reproduction

Use the documented 386/SSE2/CGO environment. Run `^TestBriefingWindow` for twelve
focused roots or `^TestBriefing` for the full 23-root family, with
`OPENNOX_BRIEFING_ASSETS` pointing to the original Nox data directory.
`OPENNOX_BRIEFING_WINDOW_CAPTURE` and `OPENNOX_BRIEFING_CAPTURE` optionally write
complete captures. Tracked tests contain frozen expectations; captures are not
needed from this VM to recover. The accumulated pattern already selects both.

Use [briefing-chapter.yaml](briefing-chapter.yaml),
[briefing-chapter-pixels.json](briefing-chapter-pixels.json) and
[RECOVERY.md](RECOVERY.md) for headless integration. Local evidence is under
build/port-briefing-window, especially native-qualification.json; the gameplay run
is build/baseline/runs/client-briefing-window-port. Captures were compressed with
round-trip/hash verification. Original assets and archive are unchanged.
