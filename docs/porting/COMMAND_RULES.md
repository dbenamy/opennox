# Command-rule loading — 2026-09-10

Scope: map-name adapter 57A950, user/map selection 4D0550, command-file reader 4D0670
and exact section lookup 57AE30. This pipeline dispatches console commands; it
is separate from the settings-rule parser and must not mutate ruleLoaderContext.

## Original-C baseline

The fixture records ExecConsoleCmd calls instead of executing game commands.
It saves/restores the working directory, header table, filename blob ranges,
handle arena, callback and game flags. It checks file contents, table/constants,
registered handles and the separate settings-parser context. Callback return
values alternate; optional callback flag changes must affect subsequent commands
and remain visible until the fixture restores the caller's state.

Header tests cover all seven section names, case/whitespace/NUL variations and
per-character substitutions. 68 direct-file cases cover eleven flag patterns,
LF/CRLF/final-no-LF, static and changing flags, empty/missing files, blank and
whitespace lines, comments, exact headers, NUL, Latin-1/UTF-8 byte widening, bare
CR and whole-physical-line truncation to 254 bytes. 73 selection cases cover
user-file presence/emptiness, map fallback, the Go wrapper, case-insensitive paths,
backslash-only directory recognition, NUL arguments, map/extension variants,
defined short map names and the nullable path entry. The dead internet branch
must never dispatch an existing internet.rul.

Compatibility details established against C:

- A section is active when ANY game flag intersects its mask, including COMMON.
- Headers match the entire case-sensitive line; spaces/CR prevent recognition.
- The filesystem bridge normalizes CRLF and reads the whole physical line before
  truncating its copy. This reader keeps 254 bytes, unlike the settings reader's 255.
- Bytes are widened individually to U+00xx before Go command dispatch.
- Failed command execution does not stop reading. Callback flag changes propagate.
- With no backslash, Arena.map first probes Arenauser.rul, then Arena.rul.
- 57A950 strips four characters from `maps\\` plus the name, not the name alone:
  short names a/ab/abc first probe user.rul/mauser.rul/mapsuser.rul respectively.
- A successful empty user file blocks fallback; missing both files still returns 1
  from the path/map helpers, while a nil path returns 0.

All original-C cases pass on 386 before replacement. Inputs that underflow the C
path buffer or cause a non-EOF read-error loop are excluded from that baseline.
Production C before conversion: **141,340 physical lines**, 153 files, zero
reference C lines. Local artifacts: `build/port-rule-command/`.
