# Console command handling

Selected: all 45 bodies /900 body lines in client__system__parsecmd.c,
plus map-column formatter sub_4417E0 and shared score-limit setter
sub_409FB0_settings. The latter is the only outside caller of variadic
sub_440A20; moving it can retire that C-only formatting path entirely.
Do not keep a varargs bridge solely for converted callers.

Ownership: actual console registry is already Go in parsecmd.go. Current
WrapCommandC sets server/client context, allocates UTF16 tokens, and calls C.
Use that real registered path for baseline handlers and compare printed output,
return/arity behavior, settings, player statuses and message queues. Localized
strings must provide meaningful formats and populated command-token tables.
The remote-command dispatcher borrows a player owner, records current sender,
and calls existing observer/script/console/camera/broadcast/audio services.
Cover dispatch and actual owners where practical, observing existing service
interfaces only when those services are outside conversion scope.

Compatibility concerns: the remote mode-gated early return required the
borrowed-sender correction described below. Numeric parsing is wcstol with
32-bit results followed by byte/word narrowing, not unrestricted Go integers.
Player limit clamps signed values to0..999, while lessons narrow before clamping.
Name concatenation and map rows use C byte/UTF16 formatting and fixed buffers;
control separator and formatter table data. Case-insensitive tokens and player
names must use compatible comparison. LAN and T1 share rate/type but differ in
online flag. mute/unmute differ for local/server routes and reserved slot31.

Test groups: arity/registration, toggle/settings/name/password/quality, mode
limits/notifications, monster/spell selection, roster/mute/admission, map lists,
GUI action routes, gameplay cheats, remote sender/observer/script/broadcast/camera.
Reuse ServerConfig, ServerOptions, MatchRoster, PlayerControls, TeamRuntime,
MapCatalog, ReliableReport and existing command-rules contracts. Qualify scope
based on actual shared owners; baseline production can reuse the voting native
qualification only if production source stays identical.

The baseline uses the actual Go command registry and existing C handlers.
Local discovery references are in build/port-console-commands/references.json;
original and repaired selected-body hashes are tracked in
console-commands-selection.json.

## Reproduced prerequisite correction

The real registered monster-toggle handler printed “respawn OFF” while leaving
bit8 enabled. Independent contracts failed for both initial bit states, mixed
case input and client/server contexts (eight failures in second.json). Change the
`off` branch from the flag-add helper to the existing flag-remove helper before
freezing. This is an authorized reversible correctness fix; it requires fresh
production qualification for the C baseline.

Fixture corrections used the existing UTF16 allocator, pointed at the owned
player instead of a copy, initialized all three sequence-list sentinel words,
and inspected broadcast Kind1 queues rather than the reliable-report queue.
These corrections did not change production behavior or regenerate a frozen oracle.

The remote-dispatch contract also reproduced a borrowed-sender lifetime bug:
15 combinations of gated mode/action returned with the current sender still set.
Clear it on that early return, consistent with every other completed dispatch
branch. Nil-player and missing-unit behavior is unchanged. The correction is
covered before freezing and will receive fresh C production qualification.

Map listing exposed a fixture-model mismatch rather than a production defect:
Nox's custom wide formatter ignores width and precision for `%S`. Its actual
output appends the byte name and `.map` suffix, followed by two tabs, without
padding to20columns. The contract now asserts that behavior after inspecting
noxstring.c; native formatting must preserve it. User listing includes localized
mute labels and Unicode names. Spell tests cover title/name lookup, unchanged
states, chat-mode restriction and the special flagball spell exemption.

## Frozen contracts and qualification scope

Two independent processes (`final-a` and `final-b`) passed 23 focused roots /
527 test events, with no skips. Their 17 byte-identical captures contain 571
records; SHA256 expectations are frozen in the tests and listed in
[console-commands-captures.json](console-commands-captures.json). Contract loops
also assert behavior independently of these captures.

Coverage includes client/server command context, arity, localized switches,
nonzero connection-rate tables, numeric parsing and byte/word narrowing,
UTF16-to-byte server names, password storage, spell lookup and restrictions,
case-insensitive player matching, reserved slot31, admission lists, quest paths,
map rows, menu dispatch, team/mode settings, level changes and viewport requests.
Remote dispatch uses actual players, script VM, camera, packet queues and audio
queues. Existing service hooks observe console execution, observer admission and
ability cancellation, whose implementations are outside this conversion.
Permission variants and every completed dispatch assert sender cleanup.

Numeric boundaries include ASCII signs/whitespace, suffixes, 32-bit saturation,
non-ASCII whitespace and fullwidth digits. Server-name contracts preserve the
legacy low byte of each UTF16 code unit and embedded-NUL behavior, including
raw stored bytes to avoid JSON replacement of invalid UTF8 masking a mismatch.
The native implementation should remove unsafe fixed intermediate buffers while
preserving defined string results and the final 15-byte stored-name limit.

The selected 298-root regression gate includes rules commands, server settings,
options, teams, player controls, roster/gameplay reports, reliable reports,
map catalog, scripts, voting and client audio. The shared player-control fixture
only gains an isolated operation57 for the actual registered level command;
existing operations and allocation/capture infrastructure remain unchanged.
Static mapped-memory validation passes. All three target sweeps and fresh production
qualification pass. The first driver invocation rejected a multiline
pattern before running tests; the corrected pattern is one anchored expression.

Limitations: the headless MOTD contract verifies dispatch/reset, while existing
GUI owners cover options; physical display and audible output remain release
checks. Existing unimplemented commands retain their current return behavior.
Full production qualification is required because the two C fixes change code.


## Qualified C baseline

All three targets executed 298 roots /20,209 passing tests, with no skips;
168 captures /32,009 records match across targets. All four gates saw identical
source. Three fresh ELF32/SSE2/CGO binaries pass ABI checks. The full asset suite
matches all 1,553 known failure entries and package outcomes (15 pass /3 fail /
32 skip). Options/gameplay, save/load and flat-render scenarios all match their
references. See [console-commands-c-qualification.json](console-commands-c-qualification.json).
Working production C is 34,344 lines /71files; reference C remains zero.

Artifacts: `build/port-console-commands/c-{default,server,highres}-qualified`,
`c-production`, `c-static.log`, `final-a.*`, `final-b.*`. Completed voting-native
asset copies were deduplicated against SHA256-verified originals, reclaiming
1,660,044,319 bytes; per-run restoration manifests preserve reproducibility.
