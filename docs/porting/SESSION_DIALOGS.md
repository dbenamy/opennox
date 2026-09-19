# Session dialogs and server filters — qualified native Go

Connected scope: the server filter, message of the day (MOTD), disconnect dialog
and quit menu. The selection and caller inventory identify 32 C functions, with
nine local owner candidates and the shared client state word. The 31 live bodies now use Go; one proven no-op was removed. Frozen C captures
and fresh native production qualification pass.

## Prerequisite correction to review

The MOTD file loader (`nox_motd_4463E0`) reads the full file. The transfer path in
`game.go` and `save.go` likewise preserves the message without limiting line length.
The display routine `nox_xxx_motd_4467F0` previously copied each line into a
256-byte local array, and `nox_xxx_motdAddSomeTextMB_446730` formatted it into a
256-unit wide array. Both operations lacked a length bound. Existing behavior
outside those buffers is undefined and must not become a compatibility golden.

Before freezing the C baseline, allocate the split buffer from message length
(with space for the blank-line substitution), and the wide buffer from line length.
Free both after their synchronous consumers finish. Preserve byte widening,
CR/LF/CRLF splitting, blank-line substitution and the final unterminated line.
The real listbox continues to copy at most 255 units per visible row. Allocation
failure skips the affected text; it does not introduce a new UI error dialog.
This reversible correction follows the user's standing authorization.

The corrected C baseline was **14,457 physical lines /58 files /zero reference C**,
+6 prerequisite lines relative to qualified parent `e6245c82`. This is not port
progress. Because production changed, the baseline received fresh three-target
production, ABI, known-suite and integration qualification.

## Baseline contracts

- Initial C line/filter contracts pass: 833 line cases and 8,192 predicates.
- Real resources pass 192 filter cases, three MOTD lifecycle cases and four
  disconnect viewport cases. File ownership and quit callback contracts pass.
- The first eight-root suite failed only because the quit fixture assumed no
  teams; its shared real server owner actually starts with two. The corrected
  contract controls counts 0, 2 and 256, including the C byte-return boundary.
- The display regression exercises lengths 0, 1, 254, 255, 256, 257, 1,024 and
  4,096 through the actual MOTD window, with byte widening, row-copy limits,
  unchanged input, display gates and pending-message flag checks.
- The tracked filter scenario has 13 screens. Original capture and independent
  repeat match exactly; reviewed checkbox/radio changes persist on reopening.
  Use `interact` for hover before pressing a checkbox. The earlier click-only
  capture did not exercise those changes and is not coverage evidence.

Rule comparison, filter configuration/actions and disconnect action/drawing
contracts now pass. They check all 28 compared rule bytes, slot-neighbor boundaries
and image coordinates. The MOTD gate uses runtime quest state at1556160, separately
from game-mode flags; both are explicitly owned and distinguished in the fixture.
The 33-root focused pass includes20 existing session-entry roots.

Final load/autosave checks pass. The recoverable C baseline was committed before replacing the selected bodies. Existing briefing/inventory/options
fixtures share selected owners and belong in the affected sweep.

## Frozen baseline captures

Default/highres:307 roots and227 artifacts; server:306 roots and226 artifacts.
No skips. The sole target difference is the existing `!server` world-selection
contract/capture. All shared artifacts match exactly across targets, all142 parent
artifacts are unchanged, and the nine new dialog captures independently repeat
between final-focused and default runs. All2,715 source fingerprints agree.
The manifest freezes these hashes; no expectations were regenerated to hide a
port difference. Static-final and all34 focused session roots pass. Fresh production builds/ABI, exact known-suite comparison and all three headless
scenarios pass. The suite retains exactly1,553 known failure entries and15 pass /
3 fail /32 skip packages. All13 filter screenshots match the unchanged C reference.
See [qualification](session-dialogs-c-qualification.json).

The historical filter cleanup1193360 has no writer or address escape in current
source/registrations. Its zero-initialized C global is read only by4896E0, called
from the Go entry cleanup wrapper. Retire that no-op call/helper/global with the
conversion, preserving the wrapper's independent Go context reset. This accounts
for one selected body; do not translate it solely to retain dead storage.

## Coverage boundaries

The focused action contracts observe external load/save/quit callbacks through
existing Go hooks; they do not themselves perform disk saves or leave a live
session. Real save/load and gameplay run separately in headless production.
Disconnect tests exercise wait and the ordinary timeout exit; the old download-
state cleanup branch is reviewed against its C calls during translation. Existing
options, inventory, briefing, rule and browser contracts cover their shared owners.
No fixture claims exhaustive GUI interaction coverage.

## Qualified native conversion

C baseline `29d96787` is committed/pushed. All 31 live bodies are translated or
inlined, the no-op cleanup removed, and shared fixture owner views updated.
Four selected C interfaces remain for unconverted callers:446360,445C40,446780
and446950. Retire them only when their last C callers move. The conversion removes
28 private interfaces, nine C owners and the empty quit-menu C translation unit.
Production C: **13,442 physical lines /57 files /zero reference**, **−1,015 lines**.

Final default/highres sweeps pass309 roots and227 artifacts; server passes308
roots and226 artifacts. Every frozen C artifact matches, with no skips. Two new
independent native contracts cover unrestricted-filter nil-record bypass and
nonzero initial disconnect resource positions; the baseline307/306 roots remain.
The existing world-selection test/capture is client-only. Source fingerprints
agree across all target and production runs. Static checking also passes.

Fresh client/highres/server builds and ABI checks pass, including all four retained
exports and removal of the28 interfaces/nine owners. The full asset suite retains
exactly1,553 known failure entries and15 pass /3 fail /32 skip package outcomes.
All13 filter screenshots, character-creation/gameplay and save/load scenarios
match their unchanged references. See
[native qualification](session-dialogs-native-qualification.json) and
[manifest](session-dialogs-native-batch.json). Local evidence lives under
`build/port-session-dialogs/native-qualified-*`.

Review corrections preserve unrestricted-mode return-before-record-read ordering,
conditional version/ping reads and raw MOTD visibility flags. Disconnect centering
uses EndPos at C offsets24/28, including initial resource position. The initial Go
draft used SizeVal; the shipped origin-zero resource hid the difference. The
corrected implementation and three-position regression preserve C behavior.
A cgo declaration spelling mismatch and a test fixture API typo were also fixed
before final qualification; no golden was changed to accommodate a mismatch.

The native MOTD file reader preserves byte contents, NUL termination, short-read
zero tails, counts and ownership for supported files. It rejects lengths that
cannot fit a signed32 allocation and handles allocation failure without the old
nil access. Those formerly undefined cases are recorded as a reversible choice
for review in DECISIONS.md, alongside the prerequisite long-line correction.

Completed C and initial native scenario copies were hash-deduplicated, reclaiming
3,338,302,701 bytes. Per-run manifests support restoration; original assets,
reference screenshots, changed saves and all qualification reports remain.
