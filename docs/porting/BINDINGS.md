# Binding editors

The next connected batch covers the in-game and main-menu key-binding editors:
21 C functions / 757 physical lines. The C implementation is still installed;
the corrected C baseline is qualified: **7,425 frozen records / nine groups**,
matching default, repeat, server and highres, plus all 32 affected tests without
skips. Gameplay qualification and translation remain.

## Scope and contracts

- GAME3_1.c: 004C3500 through the boundary at 004CA540, and 004CB880
  through the end of the file.
- client__gui__guiinput.c and client__shell__inputcfg__inputcfg.c.
- Real listbox rows, localized key/event titles and control-event bindings.
- Assignment removes all matching titles from both columns, then updates the
  selected row and clears selection. Invalid keyboard codes leave state alone;
  the mouse-title adapter instead represents invalid codes as an empty string.
- Escape cancels capture. Mouse buttons and wheel directions have distinct codes.
- Applying walks each row's secondary key before its primary key and always
  restores Escape to ToggleQuitMenu. Test the resulting binding owner, including
  order and duplicates, rather than incidental pointer-valued C returns.
- Constructors, shared scrolling, capture focus, apply/cancel, cleanup and actual
  configuration output require separate lifecycle coverage. The menu version
  also owns a transition animation and a completion callback.

## Review notes

No production behavior changes have been made. Sharing Go helpers is appropriate
only where the two original owners agree. Keep animation ownership separate.
The assignment fixture uses real listboxes and independently checks every registered
key plus invalid codes, both columns, absent selection and first/last rows.
It explicitly owns live C globals separately from backing blob bytes, and restores
localized string constants on cleanup.

C remaining at batch start: **63,323 lines / 85 files**, zero reference C.
Ignored working audit and diagnostics: build/port-bindings.

### Findings during baseline development

- The 256 apply combinations pass against C in both editors (512 records),
  independently checking insertion order, clearing old bindings and mandatory
  Escape. The final hashes are frozen in bindings-captures.json.
- Escape clears the selected row but retains the selected-list global. This is
  distinct from successful assignment, which clears both.
- The key-entry callbacks push/pop the GUI modal stack. They do not acquire or
  release mouse capture. Tests exercise and inspect the actual modal stack.
- At the synthetic 480-pixel screen width, the in-game 640-pixel root uses an
  unsigned subtraction for centering. Its wrapped endpoint is reordered by
  Window.SetPos. Preserve this existing unsupported-size behavior in translation;
  reconsider it only with a separate decision about supported window sizes.
- Listboxes reject selecting empty text rows. Invalid-key cases must account for
  that widget behavior before invoking the editor handler.


## Initial qualified C baseline (d6879118)

| Group | Records |
| --- | ---: |
| Assignments | 5,004 |
| Modal input | 1,800 |
| Apply | 512 |
| Construction | 12 |
| Shared routing | 60 |
| Wheel filtering | 24 |

Eight additional temporary configuration-file round trips exercise the production
hotkey serializer and parser. Duplicate keys serialize at their first position
with their last action; this includes the final mandatory Escape override. This
checks hotkey persistence, not the unrelated video/audio/server configuration
writer or full apply/return-to-options lifecycle; gameplay coverage remains.

Evidence: build/port-bindings/c-default (29.118s), c-repeat (6.560s), c-server
(112.963s), c-highres (35.330s), c-affected (16.427s). Each records unchanged source
identity. c-index-proof.json matches the staged source to all five phases.
Default/server use GOMAXPROCS=2; highres/repeat/affected use 1; heap limit 768MiB.
The static memory-access preflight also passes. Earlier numbered runs are fixture
development diagnostics only, not accepted baseline evidence.

The baseline commit does not change production algorithms or C LOC. A preliminary
Go state-owner draft is ignored at build/port-bindings/gui_bindings_state.go;
it has not been compiled or installed and must not be treated as completed work.

## Prompt prerequisite discovered in gameplay

The first real menu run exposed a blank key-entry prompt. The static-text widget
clears its pointer when initialized with an empty string. Construction supplies
the empty dynamic buffer; later C routing formats that buffer but never updates
the widget. The new TestBindingEditorPromptText reproduces the failure in both
editors: the buffer contains the expected action, while visible text is empty.

A scoped correction refreshes child 981 after formatting the prompt in each C
router. It avoids changing static-text behavior for every other UI owner. This is
an intentional, reversible prerequisite fix to review, not part of translation.
It adds three physical C lines, including one needed declaration header include.
The regression checks consecutive prompts with different action titles, including
Unicode. Its after-fix qualification passes; the earlier committed six-group baseline
remains the unmodified-C evidence.

The development gameplay run bindings-menu-develop reached all five screenshots
but lacked a quit step and timed out after completing the script. It is diagnostic
only. No child process remains. Add an explicit quit step to subsequent scenarios.


### Qualified prompt correction

The nine-group corrected baseline adds four visible-prompt records, six auxiliary
window cases and three cancel/visibility cases. All earlier frozen hashes remain
unchanged. Default/repeat/server/highres and all 32 affected tests pass without
skips. The static memory check passes. Source/index proof: c-prompt-index-proof.json.
- c-prompt-default: 33.636s, GOMAXPROCS=1.
- c-prompt-repeat: 7.064s, GOMAXPROCS=1.
- c-prompt-server: 177.865s, GOMAXPROCS=1.
- c-prompt-highres: 33.631s, GOMAXPROCS=1.
- c-prompt-affected: 16.847s, GOMAXPROCS=1.

Current C remaining: **63,326 lines / 85 files**, zero reference C. The original
757-line candidate now contains three prerequisite lines to translate as well.

The corrected menu development replay visibly names the action, assigns A to
Move Forward, removes its former duplicate, and returns to options. Explicit
post-transition waits are required before interacting: an animation completion
can otherwise change focus after an early prompt opens. Final combined menu and
in-game gameplay qualification is still pending.

For disk space, verified duplicates in completed summon-native, summon-flat-native
and bindings-menu-develop2 runs were removed using the existing hash-and-restore
manifest procedure (1,660,074,083 bytes). Original assets, screenshots, logs,
modified files and saves remain. Restore with the recorded per-run manifest.
