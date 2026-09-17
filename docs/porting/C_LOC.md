# C source-size checkpoints

Update this table after every completed conversion chunk. Run
`python3 tools/porting/c_loc.py` from the repository root, or pass a Git revision
to reproduce a historical count.

Count physical lines, including blanks and comments, in tracked `src/**/*.c`.
Exclude headers, C embedded in Go preambles, external dependencies and generated
build outputs. Files beginning with `//go:build porttest` are test references and
are counted separately. This measures source size, not active code, semantic
reachability or remaining porting effort. See [the inventory](C_INVENTORY.md) for
target build selection and linker evidence.

| Completed chunk | Production .c files | Production C lines | Change | Test-reference C lines |
| --- | ---: | ---: | ---: | ---: |
| Before conversions (`0e9d2e1f`) | 153 | 142,665 | — | 0 |
| Protection checksum (`66fa7bd4`, 2026-09-10) | 153 | 142,637 | −28 | 33 |
| Retire checksum C test reference (2026-09-10) | 153 | 142,637 | 0 | 0 |
| Protection record lookup/index/swap (2026-09-10) | 153 | 142,570 | −67 | 0 |
| Protection spell/ability bitsets (2026-09-10) | 153 | 142,503 | −67 | 0 |
| Protection integer/float construction (2026-09-10) | 153 | 142,458 | −45 | 0 |
| Protection deletion/cleanup (2026-09-10) | 153 | 142,393 | −65 | 0 |
| Protection randomized insertion (2026-09-10) | 153 | 142,351 | −42 | 0 |
| Protection reserved records/handles (2026-09-10) | 153 | 142,327 | −24 | 0 |
| Protection rekey/shuffle (2026-09-10) | 153 | 142,265 | −62 | 0 |
| Retire unused protection C bridges (2026-09-10) | 153 | 142,265 | 0 | 0 |
| Protection integer/byte/word setters (2026-09-10) | 153 | 142,189 | −76 | 0 |
| Protection additive updates (2026-09-10) | 153 | 142,130 | −59 | 0 |
| Protection buffer validation (2026-09-10) | 153 | 142,115 | −15 | 0 |
| Protection object checksum/toggles (2026-09-10) | 153 | 141,984 | −131 | 0 |
| Protection float updates (2026-09-10) | 153 | 141,941 | −43 | 0 |
| Protection initialization (2026-09-10) | 153 | 141,914 | −27 | 0 |
| Protection floating RNG/state (2026-09-10) | 153 | 141,844 | −70 | 0 |
| Client unit-code/bit helpers (2026-09-10) | 153 | 141,821 | −23 | 0 |
| Dynamic unit-code/extent lookup (2026-09-10) | 153 | 141,786 | −35 | 0 |
| Waypoint allocation/link/predicates (2026-09-10) | 153 | 141,745 | −41 | 0 |
| Map-rule loading/parsing (2026-09-10) | 153 | 141,455 | −290 | 0 |
| Rule writing/online buffer fix (2026-09-10) | 153 | 141,351 | −104 | 0 |
| Rule-file deletion (2026-09-10) | 153 | 141,340 | −11 | 0 |
| Command-rule loading/dispatch (2026-09-10) | 153 | 141,215 | −125 | 0 |
| Spell-class eligibility / unused chat predicate (2026-09-10) | 153 | 141,180 | −35 | 0 |
| Player-ping minimum/average (2026-09-10) | 153 | 141,126 | −54 | 0 |
| Network alias table / exhaustion fix (2026-09-11) | 153 | 141,082 | −44 | 0 |
| Glyph/item eligibility and caches (2026-09-11) | 153 | 141,042 | −40 | 0 |
| Collision reflection / containment (2026-09-11) | 153 | 141,000 | −42 | 0 |
| Line projection / clamping (2026-09-11) | 153 | 140,903 | −97 | 0 |
| Durability classification (2026-09-11) | 153 | 140,879 | −24 | 0 |
| Waypoint link insertion (2026-09-11) | 153 | 140,845 | −34 | 0 |
| Tile selection state (2026-09-11) | 153 | 140,785 | −60 | 0 |
| Tile-fill worklist (2026-09-11) | 153 | 140,730 | −55 | 0 |
| Border selection / approved variation fix (2026-09-11) | 153 | 140,672 | −58 | 0 |
| Border edge mapping (2026-09-11) | 153 | 140,613 | −59 | 0 |
| Edge normalization (2026-09-11) | 153 | 140,578 | −35 | 0 |
| Subtile predicate / list lookup (2026-09-11) | 153 | 140,455 | −123 | 0 |
| Native Go grid route; C callers retained (2026-09-11) | 153 | 140,455 | 0 | 0 |
| Floor-rendering eligibility (2026-09-11) | 153 | 140,442 | −13 | 0 |
| Six AI movement actions and private helpers (2026-09-11) | 153 | 140,260 | −182 | 0 |
| Roaming history and successor selection (2026-09-11) | 153 | 140,082 | −178 | 0 |
| Main roaming update and unused history bridges (2026-09-11) | 153 | 139,951 | −131 | 0 |
| Guard, escort and sound investigation (2026-09-11) | 153 | 139,549 | −402 | 0 |
| Navigation and retreat actions/private policies (2026-09-11) | 153 | 139,165 | −384 | 0 |
| Movement-path execution and waypoint graph (2026-09-11) | 153 | 138,824 | −341 | 0 |
| Combat AI action owners and private helpers (2026-09-11) | 153 | 138,212 | −612 | 0 |
| AI lifecycle, revival and item searches (2026-09-11) | 153 | 137,882 | −330 | 0 |
| Monster commands, animation and state (2026-09-11) | 153 | 137,297 | −585 | 0 |
| Main monster AI and defensive reactions (2026-09-11) | 153 | 136,741 | −556 | 0 |
| Monster spell decisions and cast actions (2026-09-11) | 153 | 136,242 | −499 | 0 |
| Monster callback registration, strikes, death effects and loot (2026-09-11) | 153 | 135,294 | −948 | 0 |
| Object initialization and small death callbacks (2026-09-11) | 153 | 134,954 | −340 | 0 |
| Quest death-penalty policy and six private helpers (2026-09-11) | 153 | 134,569 | −385 | 0 |
| Generic object death callbacks across GAME5 and object-die C (2026-09-11) | 152 | 134,304 | −265 | 0 |
| Monster generator death, update, placement, spawn and copy (2026-09-11) | 152 | 133,836 | −468 | 0 |
| Spawn ownership, admission, visibility culling and periodic tick (2026-09-11) | 152 | 133,272 | −564 | 0 |
| Shop pricing, stock, sessions, offers, repair and sales (2026-09-11) | 152 | 131,935 | −1,337 | 0 |
| Trade opening, offer admission, purchases and sales (2026-09-11) | 151 | 131,120 | −815 | 0 |
| Health, poison, mana and gold; C-owned modifier-slot write fix (2026-09-11) | 150 | 130,474 | −646 | 0 |
| Inventory pickup/drop, placement/chest and specialized equipment pickup (2026-09-11) | 150 | 129,053 | −1,421 | 0 |
| Weapon/armor equipment, modifier dispatch and shield selection (2026-09-11) | 149 | 128,081 | −972 | 0 |
| Modifier effects, recharge and wand use (2026-09-11) | 149 | 127,104 | −977 | 0 |
| Temporary objects, projectile updates and Spark creation (2026-09-11) | 149 | 126,124 | −980 | 0 |
| World mechanisms, elevators/teleports and directional force (2026-09-11) | 149 | 125,303 | −821 | 0 |
| Objective objects, scoring and obelisk recharge | 149 | 124,395 | −908 | 0 |
| Player attack, melee/ranged helpers and reload (2026-09-11) | 149 | 123,290 | −1,105 | 0 |
| Projectile collisions, arrows/chakrams and traps (2026-09-11) | 149 | 122,284 | −1,006 | 0 |
| Damage dispatch, defense effects and durability (2026-09-11) | 149 | 120,952 | −1,332 | 0 |
| Object state, geometry and ownership (2026-09-11) | 149 | 119,818 | −1,134 | 0 |
| Reward generation (2026-09-11) | 149 | 117,956 | −1,862 | 0 |
| Locked-door padding fix (2026-09-12; no conversion) | 149 | 117,956 | 0 | 0 |
| Player controls, respawning, observers and bot transitions (2026-09-12) | 149 | 115,985 | −1,971 | 0 |
| Player-controls obsolete export cleanup (2026-09-12; no C algorithm change) | 149 | 115,985 | 0 | 0 |
| Spell casting, book queues and buff lifecycle (2026-09-12) | 149 | 114,659 | −1,326 | 0 |
| Instant spell effects, summons, charm and portals (2026-09-12) | 149 | 112,919 | −1,740 | 0 |
| Sustained spells and teleport callbacks (2026-09-12) | 149 | 110,283 | −2,636 | 0 |
| Room geometry, occupancy, exclusions and decoration selection (2026-09-13) | 149 | 108,877 | −1,406 | 0 |
| Map-painting C stack-record prerequisite (2026-09-13) | 149 | 108,858 | −19 | 0 |
| Map painting, walls, borders and door placement (2026-09-13) | 149 | 106,609 | −2,249 | 0 |
| Population stack-record and object-disposal prerequisites (2026-09-13; no conversion) | 149 | 106,611 | +2 | 0 |
| Population candidate arrays and complete item attributes prerequisite (2026-09-14) | 149 | 106,609 | −2 | 0 |
| Population, prefabs, inventories, exits and waypoints (2026-09-14) | 148 | 104,839 | −1,770 | 0 |
| Hallway second-corridor storage prerequisite (2026-09-14; no conversion) | 148 | 104,839 | 0 | 0 |
| Hallway routing, bends, obstruction admission and candidate connections (2026-09-14) | 148 | 104,307 | −532 | 0 |
| Theme value-buffer and modifier-counter prerequisites (2026-09-14; no conversion) | 148 | 104,300 | −7 | 0 |
| Inherited modifier removal prerequisite (2026-09-14; no conversion) | 148 | 104,300 | 0 | 0 |
| Complete theme parser C baseline (2026-09-14; 3,392 cases, no conversion) | 148 | 104,300 | 0 | 0 |
| Theme parser, conditional input, equipment and decoration definitions (2026-09-14) | 148 | 102,419 | −1,881 | 0 |
| Growth direction-array and merge-setting prerequisites (2026-09-14; no conversion) | 148 | 102,416 | −3 | 0 |
| Growth translated-door prerequisite and locked baseline (2026-09-14; no conversion) | 148 | 102,416 | 0 | 0 |
| Initial layouts, recursive growth, doors and waypoint connections (2026-09-14) | 148 | 101,335 | −1,081 | 0 |
| Wall-list correctness prerequisite and stable C generator integration (2026-09-14; no conversion) | 148 | 101,335 | 0 | 0 |
| Unflagged backdrop correction and expanded orchestration C baseline (2026-09-14; no conversion) | 148 | 101,335 | 0 | 0 |
| Map orchestration: generation, retries and save coordination (2026-09-14) | 148 | 101,128 | −207 | 0 |
| Journal padding prerequisite and first reporting baseline (2026-09-14) | 148 | 101,128 | 0 | 0 |
| Complete73-routine gameplay reporting C baseline (2026-09-14) | 148 | 101,128 | 0 | 0 |
| Gameplay reporting, inventory notifications and player-state aggregation (2026-09-14) | 148 | 99,791 | −1,337 | 0 |
| Gameplay text baseline (2026-09-14; 2,791 cases, no conversion) | 148 | 99,791 | 0 | 0 |
| Gameplay text, notifications and player iteration (2026-09-14) | 147 | 99,509 | −282 | 0 |
| Object lookup and 16-node net-code cache (2026-09-14) | 147 | 99,071 | −438 | 0 |
| Quest eligibility, modifiers and inventory limits (2026-09-14) | 147 | 98,506 | −565 | 0 |
| Client effects renderer/curve fixture checkpoint (2026-09-14; no conversion) | 147 | 98,506 | 0 | 0 |
| Missing plasma endpoint correction before effects baseline (2026-09-14; no conversion) | 147 | 98,504 | −2 | 0 |
| Client effects: initialize chain-lightning particle endpoint prerequisite (2026-09-14) | 147 | 98,505 | +1 | 0 |
| Client particles, rays, glow, lightning, plasma and curves (2026-09-14) | 142 | 96,216 | −2,289 | 0 |
| Drawable tail allocation guards before updates baseline (2026-09-14; no conversion) | 142 | 96,221 | +5 | 0 |
| Client drawable updates (2026-09-14) | 130 | 95,146 | −1,075 | 0 |
| Procedural particle drawing, lighting and color initialization (2026-09-14) | 124 | 94,415 | −731 | 0 |
| Sprite animation, frame parsing and boulder drawing (2026-09-15) | 118 | 93,775 | −640 | 0 |
| Object drawing prerequisites and C baseline (2026-09-15; no conversion) | 118 | 93,782 | +7 | 0 |
| Object drawing, materials, team colors and generator/summon lifecycle (2026-09-15) | 106 | 92,869 | −913 | 0 |
| Screen effects corrected C baseline (2026-09-15; no conversion) | 106 | 92,869 | 0 | 0 |
| Screen particles, remaining drawing and distance/raster helpers (2026-09-15) | 101 | 92,307 | −562 | 0 |
| Shared object renderer, ghost/shiny, beam and clipping helpers (2026-09-15) | 101 | 91,740 | −567 | 0 |
| Shared UI clipping/raster helpers and progress bars | 101 | 91,478 | −262 | 0 |
| Slider widgets, input and drawing | 101 | 90,879 | −599 | 0 |
| Radio-button selection, text and drawing | 101 | 90,583 | −296 | 0 |
| Text-entry ownership/bounds corrected C baseline | 101 | 90,594 | +11 | 0 |
| Text-entry input, composition, rendering and context | 101 | 89,969 | −625 | 0 |
| Listbox row/selection/scroll corrected C baseline | 101 | 89,979 | +10 | 0 |
| Listbox construction, rows, selection, scrolling and drawing | 100 | 88,572 | −1,407 | 0 |
| Inclusive window-ID range termination prerequisite | 100 | 88,574 | +2 | 0 |
| Window geometry, state, tree and draw-data helpers | 99 | 88,212 | −362 | 0 |
| Localized item-hover names and cursor tooltip storage | 98 | 87,884 | −328 | 0 |
| Meter baseline sound-observation ABI (tagged; no algorithm) | 98 | 87,888 | +4 | 0 |
| Health/mana, potion, weapon and charge meters | 97 | 86,818 | −1,070 | 0 |
| Client inventory queries, scalar state and item updates | 97 | 86,422 | −396 | 0 |
| Inventory transaction C prerequisites (pickup coordinates; compaction copy) | 97 | 86,424 | +2 | 0 |
| Client inventory stack/equipment transactions | 97 | 85,399 | −1,025 | 0 |
| Client inventory display, identification, stats and feedback | 97 | 84,237 | −1,162 | 0 |
| Client inventory windows, input, scrolling and lifecycle | 96 | 82,632 | −1,605 | 0 |
| Inventory cancellation ownership cleanup (Go only) | 96 | 82,632 | 0 | 0 |
| Quantity/trade UI qualified C prerequisites | 96 | 82,667 | +35 | 0 |
| Quantity dialog and player-to-player trade UI | 95 | 81,351 | −1,316 | 0 |
| Shop UI capacity, affordability and quantity ownership prerequisites | 95 | 81,373 | +22 | 0 |
| Shop UI storage, drawing, input and quantity transactions | 94 | 80,263 | −1,110 | 0 |
| Quest journal storage, reports, save integration and rendering | 93 | 79,924 | −339 | 0 |
| Quest briefing chapters, reports, sprite cache and presentation | 92 | 78,977 | −947 | 0 |
| Briefing window lifecycle, voice, input and transitions | 92 | 78,657 | −320 | 0 |
| Scoreboard collection, rank tables, rendering and mode state | 91 | 77,120 | −1,537 | 0 |
| Minimap rendering, zoom, objectives and AI debug traversal | 91 | 76,382 | −738 | 0 |
| World-wall rendering, projection and drawable visibility | 91 | 75,757 | −625 | 0 |
| Wall-edge RLE rasterizer and private helper interfaces | 91 | 75,473 | −284 | 0 |
| Tile texture/fill callbacks, packed raster and wrap setup | 91 | 74,029 | −1,444 | 0 |
| Tile composition, scrolling, overlays and private callback interfaces | 91 | 73,450 | −579 | 0 |
| Empty floor-definition binding prerequisite | 91 | 73,453 | +3 | 0 |
| Floor/edge asset readers, facade lookup and image-array cleanup | 91 | 72,940 | −513 | 0 |
| Things-section readers and aligned MemFile helper | 91 | 72,542 | −398 | 0 |
| Client sound-definition readers and sample lookup | 91 | 72,221 | −321 | 0 |
| Complete map decompressor, dictionary and adaptive tables | 90 | 71,252 | −969 | 0 |
| Complete map compressor, match search and adaptive writer | 88 | 69,342 | −1,910 | 0 |
| Quest selection zero-candidate fallback prerequisite | 88 | 69,345 | +3 | 0 |
| Complete map catalog, cycle parsing and quest rotation | 88 | 68,897 | −448 | 0 |
| Shared intrusive lists and player-group membership | 88 | 68,597 | −300 | 0 |
| Spellbook UI, page rendering, rewards and quickbar addition | 87 | 66,811 | −1,786 | 0 |
| Complete quickbar UI, activation, rows, traps, rendering and saved slots | 86 | 64,317 | −2,494 | 0 |
| Complete summon-creature panel, layout, commands and rendering | 85 | 63,323 | −994 | 0 |
| Binding-editor prompt refresh prerequisite | 85 | 63,326 | +3 | 0 |
| Complete in-game and main-menu binding editors | 83 | 62,571 | −755 | 0 |
| Options checkbox dispatcher prerequisite (C fix) | 83 | 62,475 | −96 | 0 |
| Complete main-menu and in-game options panels | 82 | 61,682 | −793 | 0 |
| Client map drawable readers | 82 | 61,065 | −617 | 0 |
| Colored-light degenerate-direction prerequisite (C fix) | 82 | 61,071 | +6 | 0 |
| Colored-light animation and unused map classification cleanup | 82 | 60,775 | −296 | 0 |
| Object-transfer rejection ownership prerequisite (C fix) | 82 | 60,779 | +4 | 0 |
| Server common and world-object serialization | 82 | 59,602 | −1,177 | 0 |
| Item serialization ownership and old wand defaults prerequisite (C fixes) | 82 | 59,604 | +2 | 0 |
| Item and reward serialization | 82 | 58,539 | −1,065 | 0 |
| Monster and NPC serialization (including trailing separator) | 82 | 56,968 | −1,571 | 0 |
| Visibility and effect reports, including orphan/EOF cleanup | 82 | 56,118 | −850 | 0 |
| Object and recipient reports, including private minimap count | 82 | 55,577 | −541 | 0 |
| Reliable queue pressure ownership prerequisite (C fix) | 82 | 55,584 | +7 | 0 |
| Reliable game-message queue | 82 | 54,962 | −622 | 0 |
| World collisions and interactions | 82 | 53,946 | −1,016 | 0 |
| Quest score constant width prerequisite (C fix) | 82 | 53,946 | 0 | 0 |
| Quest runtime, statistics and difficulty scaling | 82 | 53,049 | −897 | 0 |
| Roster padding, Flagball draw and player layout prerequisites | 82 | 53,046 | −3 | 0 |
| Match results and roster synchronization, private globals and EOF cleanup | 82 | 52,273 | −773 | 0 |

| Team message fields and clear/rebalance count prerequisites | 82 | 52,275 | +2 | 0 |

| Team runtime, membership and map objectives, including orphan/separator cleanup | 82 | 51,203 | −1,072 | 0 |

| Team UI row ownership, naming, selection and missing-resource prerequisites | 82 | 51,225 | +22 | 0 |
| Team HUD/player-list UI, private globals and translation-unit cleanup | 80 | 50,177 | −1,048 | 0 |
| Server-options missing-resource and panel ownership prerequisites | 80 | 50,181 | +4 | 0 |

The checksum removes two C function definitions from GAME5_2.c; C ABI entry
points remain as generated bridges into Go. Translation-unit counts do not fall
because the file still contains other functions. The 33 test-reference lines were subsequently removed after successful
differential validation; they remain recoverable from Git at `66fa7bd4`.
