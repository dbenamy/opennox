# Navigation and retreat actions — 2026-09-11

Connected scope: move-to, far-move, dodge, flee, return-home, retreat and
retreat-to-master, their lifecycle/private policies, and preceding-action lookup
50A040 (its last owner is move-to). Keep shared C movement/path engines, spell
policy, direction calculation, audio and food lookup outside this batch.

Original C passes 38,783 shared guarded-fixture scenarios. Coverage includes all
seven registered actions, lifecycle status masks, partial/full stacks, preceding
conditions, speed and dodge threshold neighbors/nonfinite inputs, health max zero,
resume threshold neighbors/NaN, enemy/owner presence, food hit/miss through the real
spatial index and vision query, path statuses/results and unsigned frame wrap.
Timer/path combinations are crossed independently. The retreat generator boundary
records input points, capacity and supplies zero/one/multiple point results; actual
C path setup/movement still run. Deeper spell-policy behavior is gated off.

Independent assertions check all three health/cast predicates and food lookup
selection. Full normalized changed-word states, action/RNG state and boundary
traces match SHA256 7996e036ac64ebdf8c2786f48052d3f3dadfd0f5fab36c2e6113907a17134f3d.
Object/monster/health/definition/target storage is C-owned and guarded; only spatial
visitation tokens may change in the food target, and those changes enter the hash.
The shared one-shot movement flag and food-search scratch globals are captured
and restored; game flags are restored too. No copied C test implementation.

Artifacts: build/port-ai-navigation (ignored). Baseline source count remains
139,549 physical C lines, 153 files, zero reference C. Native conversion and whole
batch qualification follow; no production implementation has changed yet.
