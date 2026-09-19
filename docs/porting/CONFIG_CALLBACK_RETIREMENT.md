# Orphaned configuration callbacks

Remove twelve C callbacks whose command-file parser and caller were deleted in
2022. This is unreachable-code cleanup, not a replacement parser. The qualified
browser parent is `c857b71f`. The cleanup removes 216 function-body lines and 24
separator lines: qualified C is 14,451 physical lines in 58 files, zero reference C.
The cleanup is qualified; C is down another 240 physical lines.

## Reachability

The [audit](config-callback-retirement-audit.json) records the former table
references, callback slots, and historical evidence. The old table occupies
`0x587000 + [171912,172172)`, with twelve 20-byte records plus a terminator.
Its only contemporary uses were the callback registrations. No raw pointer into
this range occurs in any embedded blob, scanning every byte alignment. Dynamic
sound, modifier, AI, reward, GUI, player-file and map-theme table helpers use
separate ranges; no current reader reaches this command table.

Commit `391c59fd` removed the file-reading caller of `sub_4A7DF0`; `d15e0aec`
removed that parser itself. The historical parser selected callback slots at
`171928 + 20*index`, establishing which table the remaining registrations served.
No current direct calls, callback registrations outside this table, or parser
replacement references these functions.

Remove the twelve definitions, header declarations, Go C-function bindings,
blob pointer struct fields and relocation assignments. Leave embedded blob bytes
and offsets unchanged. The browser still reads the start-coordinate backing words;
those words and their live Go readers remain. Do not translate unreachable string
setters just to retain tests for them.

## Validation

The existing 197-root browser/startup suite passes on default, server and highres,
checking all 142 frozen artifacts against original C. Static mapped-access checks,
fresh three-target binaries and ABI inventories pass with all twelve symbols
absent. The full suite matches the exact known failures: 1,553 entries; 15 pass /3 fail /32
skip packages. Browser, character-creation/gameplay and save/load scenarios pass. No new algorithm fixtures or changed goldens.

Local evidence is under `build/port-config-callbacks/final-*`; the tracked
[manifest](config-callback-retirement-batch.json) specifies all commands and retired
symbols. The source removal script is consumed and must never be replayed.

See the [qualification report](config-callback-retirement-qualification.json).
The final qualified browser asset copies were deduplicated after hash/inactivity
checks, reclaiming 1,669,136,451 bytes; captures, changed saves and restoration
manifests remain.
