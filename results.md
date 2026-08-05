# Phase 3 Trinket Benchmark Results

20,000-iteration simulation comparing 4 candidate Trinket2 options across 13 physical-DPS
spec/gear configurations. Trinket1 is held fixed (Bloodlust Brooch, or Warp-Spring Coil for
Rogue). A flat +30 melee hit rating bonus is applied to every configuration.

## Assumptions

- **Iterations / seed**: 20,000 iterations, `RandomSeed: 101`, for every spec/trinket
  combination (`core.RunRaidSimConcurrent`).
- **Encounter**: single target dummy, `core.MakeSingleTargetEncounter` with duration
  variation of ±5s (180s base fight length).
- **Buffs/debuffs**: `FullRaidBuffs` / `FullPartyBuffs` / `FullIndividualBuffs` /
  `FullDebuffs` from `sim/core/test_utils.go` (includes Improved Faerie Fire).
- **Hit rating**: every spec gets a flat +30 `stats.MeleeHitRating` bonus via
  `Player.BonusStats`, on top of whatever hit rating its gear preset already provides. This
  was a simplification requested by the user in place of computing exact per-spec hit caps.
- **Gear**: each spec uses its existing (or newly hand-built, for Ret Paladin) checked-in
  Phase 3 BiS gear preset, with only Trinket1/Trinket2 swapped per run:
  - Warrior Fury: `ui/warrior/dps/gear_sets/p3_fury.gear.json`
  - Warrior Arms: `ui/warrior/dps/gear_sets/p3_arms.gear.json`
  - Rogue: `ui/rogue/dps/gear_sets/p3.gear.json`
  - Hunter BM / SV: `ui/hunter/dps/gear_sets/phase_3/{bm,sv}/2h_6p.gear.json`
  - Enhancement Shaman: `ui/shaman/enhancement/gear_sets/p3.gear.json`
  - Feral Druid: `ui/druid/feralcat/gear_sets/p3_6p.gear.json` ("BiS 6%" hit variant)
  - Retribution Paladin: `ui/paladin/retribution/gear_sets/p3.gear.json` (hand-built this
    session from a Wowhead P3 BiS guide, since no P3 preset previously existed in-repo;
    fixed during this project — see "Bugs found & fixed" below)
- **Weapons**: the legendary Warglaives of Azzinoth (32837/32838) are excluded from the
  primary Warrior Fury / Rogue rows per project constraint; a non-Warglaive substitute is
  used instead (S3 Vengeful Gladiator weapons). A second "(Glaives)" row is also included
  for Warrior Fury/Arms and Rogue showing DPS *with* actual Warglaives equipped, for
  reference:
  - Warrior Fury (non-Glaive): MH Vengeful Gladiator's Cleaver (33669), OH Vengeful
    Gladiator's Chopper (34015)
  - Rogue (non-Glaive): MH Vengeful Gladiator's Slicer (33762), OH Blade of Savagery (32369)
  - "(Glaives)" rows: MH/OH Warglaives of Azzinoth (32837/32838)
  - Warrior Arms uses a two-hander (Cataclysm's Edge, 30902) regardless — no Glaive
    variant applies.
- **Talents/consumables/rotations**: each spec uses its default/standard Phase 3
  talent build, consumables, and APL rotation as found in that spec's existing test suite
  (`sim/<class>/<spec>/*_test.go`) — see `/memories/session/trinket-benchmark-project.md`
  for the exact talent strings, consumable IDs, and rotation file names used per spec.
  Hunter's "Melee weave" APL rotation option is forced off for both BM and SV.
- **Trinket candidates (Trinket2)**:
  - Abacus of Violent Odds (28288)
  - Madness of the Betrayer (32505)
  - Dragonspine Trophy (28830)
  - Tsunami Talisman (30627)
- **Scratch benchmark code**: `sim/trinketbench_test.go` (temporary, not meant to remain
  in the repo — will be deleted once results are finalized).

## Bugs found & fixed during this project

1. **Feral Druid rotation was a no-op.** The scratch benchmark originally used a bare
   `&proto.APLRotation{Type: proto.APLRotation_TypeSimple}` with no `PriorityList`, which
   the Go engine silently treats as an empty rotation (only auto-attacks fire). Fixed by
   building an explicit `APLListItem`/`CatOptimalRotationAction` priority list (mirroring
   `ui/druid/feralcat/sim.ts`'s `simpleRotation` generator). DPS went from ~715-730 to a
   realistic ~2450-2470.
2. **Ret Paladin's hand-built P3 gear preset was missing enchants and had 11 empty gem
   sockets.** This caused its simulated DPS (~1721-1753) to score *below* the existing
   `TestRetribution.results` P1-gear fixture (1915 DPS) — a red flag, since P3 gear should
   never underperform P1. Fixed by adding the missing slot enchants (copied from
   `p1.gear.json`) and filling all empty gem sockets (hybrid gems to satisfy the Relentless
   Earthstorm Diamond meta gem's color requirement, pure Strength gems elsewhere). Final
   DPS: ~1853-1905, now sensibly above the P1 baseline.

## Full results table (20,000 iterations, seed 101)

| Spec | Trinket2 | DPS | Hit Rating |
|---|---|---:|---:|
| Warrior Fury | Abacus of Violent Odds | 2177.5 | 225.0 |
| Warrior Fury | Madness of the Betrayer | 2206.3 | 245.0 |
| Warrior Fury | Dragonspine Trophy | 2246.8 | 225.0 |
| Warrior Fury | Tsunami Talisman | 2201.3 | 235.0 |
| Warrior Fury (Glaives) | Abacus of Violent Odds | 2293.6 | 230.0 |
| Warrior Fury (Glaives) | Madness of the Betrayer | 2324.6 | 250.0 |
| Warrior Fury (Glaives) | Dragonspine Trophy | 2364.4 | 230.0 |
| Warrior Fury (Glaives) | Tsunami Talisman | 2322.6 | 240.0 |
| Warrior Arms | Abacus of Violent Odds | 2000.6 | 121.0 |
| Warrior Arms | Madness of the Betrayer | 2013.8 | 141.0 |
| Warrior Arms | Dragonspine Trophy | 2048.9 | 121.0 |
| Warrior Arms | Tsunami Talisman | 2011.9 | 131.0 |
| Warrior Arms (Glaives) | Abacus of Violent Odds | 1932.1 | 142.0 |
| Warrior Arms (Glaives) | Madness of the Betrayer | 1959.5 | 162.0 |
| Warrior Arms (Glaives) | Dragonspine Trophy | 1989.4 | 142.0 |
| Warrior Arms (Glaives) | Tsunami Talisman | 1955.4 | 152.0 |
| Rogue | Abacus of Violent Odds | 2176.4 | 348.0 |
| Rogue | Madness of the Betrayer | 2170.0 | 368.0 |
| Rogue | Dragonspine Trophy | 2235.4 | 348.0 |
| Rogue | Tsunami Talisman | 2168.5 | 358.0 |
| Rogue (Glaives) | Abacus of Violent Odds | 2345.7 | 346.0 |
| Rogue (Glaives) | Madness of the Betrayer | 2340.4 | 366.0 |
| Rogue (Glaives) | Dragonspine Trophy | 2406.2 | 346.0 |
| Rogue (Glaives) | Tsunami Talisman | 2337.8 | 356.0 |
| Hunter BM | Abacus of Violent Odds | 2901.2 | 109.0 |
| Hunter BM | Madness of the Betrayer | 2909.2 | 129.0 |
| Hunter BM | Dragonspine Trophy | 2931.8 | 109.0 |
| Hunter BM | Tsunami Talisman | 2909.7 | 119.0 |
| Hunter BM (Weave) | Abacus of Violent Odds | 3285.2 | 109.0 |
| Hunter BM (Weave) | Madness of the Betrayer | 3292.1 | 129.0 |
| Hunter BM (Weave) | Dragonspine Trophy | 3338.1 | 109.0 |
| Hunter BM (Weave) | Tsunami Talisman | 3289.8 | 119.0 |
| Hunter SV | Abacus of Violent Odds | 2564.4 | 58.0 |
| Hunter SV | Madness of the Betrayer | 2565.8 | 78.0 |
| Hunter SV | Dragonspine Trophy | 2587.7 | 58.0 |
| Hunter SV | Tsunami Talisman | 2568.2 | 68.0 |
| Hunter SV (Weave) | Abacus of Violent Odds | 2996.1 | 58.0 |
| Hunter SV (Weave) | Madness of the Betrayer | 2992.9 | 78.0 |
| Hunter SV (Weave) | Dragonspine Trophy | 3051.0 | 58.0 |
| Hunter SV (Weave) | Tsunami Talisman | 2994.6 | 68.0 |
| Enhancement Shaman | Abacus of Violent Odds | 2200.2 | 183.0 |
| Enhancement Shaman | Madness of the Betrayer | 2215.4 | 203.0 |
| Enhancement Shaman | Dragonspine Trophy | 2246.2 | 183.0 |
| Enhancement Shaman | Tsunami Talisman | 2205.4 | 193.0 |
| Feral Druid | Abacus of Violent Odds | 2453.2 | 116.0 |
| Feral Druid | Madness of the Betrayer | 2452.5 | 136.0 |
| Feral Druid | Dragonspine Trophy | 2458.8 | 116.0 |
| Feral Druid | Tsunami Talisman | 2465.8 | 126.0 |
| Retribution Paladin | Abacus of Violent Odds | 1869.2 | 88.0 |
| Retribution Paladin | Madness of the Betrayer | 1853.4 | 108.0 |
| Retribution Paladin | Dragonspine Trophy | 1904.9 | 88.0 |
| Retribution Paladin | Tsunami Talisman | 1858.5 | 98.0 |

## Trinket rankings per spec (best to worst, with DPS delta to next rank)

**Warrior Fury**
1. Dragonspine Trophy — 2246.8
2. Madness of the Betrayer — 2206.3 (-40.5)
3. Tsunami Talisman — 2201.3 (-5.0)
4. Abacus of Violent Odds — 2177.5 (-23.8)

**Warrior Fury (Glaives)**
1. Dragonspine Trophy — 2364.4
2. Madness of the Betrayer — 2324.6 (-39.8)
3. Tsunami Talisman — 2322.6 (-2.0)
4. Abacus of Violent Odds — 2293.6 (-29.0)

**Warrior Arms**
1. Dragonspine Trophy — 2048.9
2. Madness of the Betrayer — 2013.8 (-35.1)
3. Tsunami Talisman — 2011.9 (-1.9)
4. Abacus of Violent Odds — 2000.6 (-11.3)

**Warrior Arms (Glaives)**
1. Dragonspine Trophy — 1989.4
2. Madness of the Betrayer — 1959.5 (-29.9)
3. Tsunami Talisman — 1955.4 (-4.1)
4. Abacus of Violent Odds — 1932.1 (-23.3)

**Rogue**
1. Dragonspine Trophy — 2235.4
2. Abacus of Violent Odds — 2176.4 (-59.0)
3. Madness of the Betrayer — 2170.0 (-6.4)
4. Tsunami Talisman — 2168.5 (-1.5)

**Rogue (Glaives)**
1. Dragonspine Trophy — 2406.2
2. Abacus of Violent Odds — 2345.7 (-60.5)
3. Madness of the Betrayer — 2340.4 (-5.3)
4. Tsunami Talisman — 2337.8 (-2.6)

**Hunter BM**
1. Dragonspine Trophy — 2931.8
2. Tsunami Talisman — 2909.7 (-22.1)
3. Madness of the Betrayer — 2909.2 (-0.5)
4. Abacus of Violent Odds — 2901.2 (-8.0)

**Hunter BM (Weave)**
1. Dragonspine Trophy — 3338.1
2. Madness of the Betrayer — 3292.1 (-46.0)
3. Tsunami Talisman — 3289.8 (-2.3)
4. Abacus of Violent Odds — 3285.2 (-4.6)

**Hunter SV**
1. Dragonspine Trophy — 2587.7
2. Tsunami Talisman — 2568.2 (-19.5)
3. Madness of the Betrayer — 2565.8 (-2.4)
4. Abacus of Violent Odds — 2564.4 (-1.4)

**Hunter SV (Weave)**
1. Dragonspine Trophy — 3051.0
2. Abacus of Violent Odds — 2996.1 (-54.9)
3. Tsunami Talisman — 2994.6 (-1.5)
4. Madness of the Betrayer — 2992.9 (-1.7)

**Enhancement Shaman**
1. Dragonspine Trophy — 2246.2
2. Madness of the Betrayer — 2215.4 (-30.8)
3. Tsunami Talisman — 2205.4 (-10.0)
4. Abacus of Violent Odds — 2200.2 (-5.2)

**Feral Druid**
1. Tsunami Talisman — 2465.8
2. Dragonspine Trophy — 2458.8 (-7.0)
3. Abacus of Violent Odds — 2453.2 (-5.6)
4. Madness of the Betrayer — 2452.5 (-0.7)

**Retribution Paladin**
1. Dragonspine Trophy — 1904.9
2. Abacus of Violent Odds — 1869.2 (-35.7)
3. Tsunami Talisman — 1858.5 (-10.7)
4. Madness of the Betrayer — 1853.4 (-5.1)

## Summary

Dragonspine Trophy is the best Trinket2 choice in 12 of 13 spec categories, typically by a
clear 20-60 DPS margin over the runner-up. Feral Druid is the sole exception, where Tsunami
Talisman edges out Dragonspine Trophy by a small ~7 DPS — effectively a toss-up within
simulation noise.

## Cleave scenario (2 targets, 20,000 iterations, seed 101)

Same 13 spec/gear configurations and same 4 Trinket2 candidates as above, but the encounter
was changed to 2 targets instead of 1 (`Encounter.Targets` extended to
`[NewDefaultTarget(), NewDefaultTarget()]`, both full-health dummies, same 180s ±5s
duration/buffs/debuffs). No rotation/APL changes were made for this scenario — each spec
still runs its existing single-target-tuned APL as-is, so any DPS gains come purely from
incidental cleave/AOE effects already present in that rotation (e.g. Whirlwind, Cleave,
Swipe, multi-target-hitting procs), not from a dedicated AOE priority list. Treat this as a
"how much do I passively benefit from a second target" estimate rather than an optimized
cleave/AOE rotation benchmark.

### Full results table (cleave, 2 targets)

| Spec | Trinket2 | DPS | Hit Rating |
|---|---|---:|---:|
| Warrior Fury | Abacus of Violent Odds | 3318.4 | 225.0 |
| Warrior Fury | Madness of the Betrayer | 3378.5 | 245.0 |
| Warrior Fury | Dragonspine Trophy | 3431.7 | 225.0 |
| Warrior Fury | Tsunami Talisman | 3361.8 | 235.0 |
| Warrior Fury (Glaives) | Abacus of Violent Odds | 3402.3 | 230.0 |
| Warrior Fury (Glaives) | Madness of the Betrayer | 3467.7 | 250.0 |
| Warrior Fury (Glaives) | Dragonspine Trophy | 3515.7 | 230.0 |
| Warrior Fury (Glaives) | Tsunami Talisman | 3457.9 | 240.0 |
| Warrior Arms | Abacus of Violent Odds | 3115.6 | 121.0 |
| Warrior Arms | Madness of the Betrayer | 3154.8 | 141.0 |
| Warrior Arms | Dragonspine Trophy | 3218.0 | 121.0 |
| Warrior Arms | Tsunami Talisman | 3144.6 | 131.0 |
| Warrior Arms (Glaives) | Abacus of Violent Odds | 2883.2 | 142.0 |
| Warrior Arms (Glaives) | Madness of the Betrayer | 2940.6 | 162.0 |
| Warrior Arms (Glaives) | Dragonspine Trophy | 2982.0 | 142.0 |
| Warrior Arms (Glaives) | Tsunami Talisman | 2929.6 | 152.0 |
| Rogue | Abacus of Violent Odds | 2624.8 | 348.0 |
| Rogue | Madness of the Betrayer | 2589.5 | 368.0 |
| Rogue | Dragonspine Trophy | 2670.7 | 348.0 |
| Rogue | Tsunami Talisman | 2588.7 | 358.0 |
| Rogue (Glaives) | Abacus of Violent Odds | 2842.1 | 346.0 |
| Rogue (Glaives) | Madness of the Betrayer | 2807.1 | 366.0 |
| Rogue (Glaives) | Dragonspine Trophy | 2890.1 | 346.0 |
| Rogue (Glaives) | Tsunami Talisman | 2807.6 | 356.0 |
| Hunter BM | Abacus of Violent Odds | 3070.3 | 109.0 |
| Hunter BM | Madness of the Betrayer | 3075.9 | 129.0 |
| Hunter BM | Dragonspine Trophy | 3091.8 | 109.0 |
| Hunter BM | Tsunami Talisman | 3077.0 | 119.0 |
| Hunter BM (Weave) | Abacus of Violent Odds | 3445.1 | 109.0 |
| Hunter BM (Weave) | Madness of the Betrayer | 3454.2 | 129.0 |
| Hunter BM (Weave) | Dragonspine Trophy | 3506.8 | 109.0 |
| Hunter BM (Weave) | Tsunami Talisman | 3453.3 | 119.0 |
| Hunter SV | Abacus of Violent Odds | 2754.9 | 58.0 |
| Hunter SV | Madness of the Betrayer | 2765.2 | 78.0 |
| Hunter SV | Dragonspine Trophy | 2779.5 | 58.0 |
| Hunter SV | Tsunami Talisman | 2766.7 | 68.0 |
| Hunter SV (Weave) | Abacus of Violent Odds | 3199.8 | 58.0 |
| Hunter SV (Weave) | Madness of the Betrayer | 3200.6 | 78.0 |
| Hunter SV (Weave) | Dragonspine Trophy | 3252.1 | 58.0 |
| Hunter SV (Weave) | Tsunami Talisman | 3198.2 | 68.0 |
| Enhancement Shaman | Abacus of Violent Odds | 2310.5 | 183.0 |
| Enhancement Shaman | Madness of the Betrayer | 2326.9 | 203.0 |
| Enhancement Shaman | Dragonspine Trophy | 2358.1 | 183.0 |
| Enhancement Shaman | Tsunami Talisman | 2316.7 | 193.0 |
| Feral Druid | Abacus of Violent Odds | 2469.2 | 116.0 |
| Feral Druid | Madness of the Betrayer | 2467.7 | 136.0 |
| Feral Druid | Dragonspine Trophy | 2474.0 | 116.0 |
| Feral Druid | Tsunami Talisman | 2480.9 | 126.0 |
| Retribution Paladin | Abacus of Violent Odds | 1895.6 | 88.0 |
| Retribution Paladin | Madness of the Betrayer | 1881.0 | 108.0 |
| Retribution Paladin | Dragonspine Trophy | 1932.1 | 88.0 |
| Retribution Paladin | Tsunami Talisman | 1885.5 | 98.0 |

### Trinket rankings per spec, cleave scenario (best to worst, with DPS delta to next rank)

**Warrior Fury**
1. Dragonspine Trophy — 3431.7
2. Madness of the Betrayer — 3378.5 (-53.2)
3. Tsunami Talisman — 3361.8 (-16.7)
4. Abacus of Violent Odds — 3318.4 (-43.4)

**Warrior Fury (Glaives)**
1. Dragonspine Trophy — 3515.7
2. Madness of the Betrayer — 3467.7 (-48.0)
3. Tsunami Talisman — 3457.9 (-9.8)
4. Abacus of Violent Odds — 3402.3 (-55.6)

**Warrior Arms**
1. Dragonspine Trophy — 3218.0
2. Madness of the Betrayer — 3154.8 (-63.2)
3. Tsunami Talisman — 3144.6 (-10.2)
4. Abacus of Violent Odds — 3115.6 (-29.0)

**Warrior Arms (Glaives)**
1. Dragonspine Trophy — 2982.0
2. Madness of the Betrayer — 2940.6 (-41.4)
3. Tsunami Talisman — 2929.6 (-11.0)
4. Abacus of Violent Odds — 2883.2 (-46.4)

**Rogue**
1. Dragonspine Trophy — 2670.7
2. Abacus of Violent Odds — 2624.8 (-45.9)
3. Madness of the Betrayer — 2589.5 (-35.3)
4. Tsunami Talisman — 2588.7 (-0.8)

**Rogue (Glaives)**
1. Dragonspine Trophy — 2890.1
2. Abacus of Violent Odds — 2842.1 (-48.0)
3. Tsunami Talisman — 2807.6 (-34.5)
4. Madness of the Betrayer — 2807.1 (-0.5)

**Hunter BM**
1. Dragonspine Trophy — 3091.8
2. Tsunami Talisman — 3077.0 (-14.8)
3. Madness of the Betrayer — 3075.9 (-1.1)
4. Abacus of Violent Odds — 3070.3 (-5.6)

**Hunter BM (Weave)**
1. Dragonspine Trophy — 3506.8
2. Madness of the Betrayer — 3454.2 (-52.6)
3. Tsunami Talisman — 3453.3 (-0.9)
4. Abacus of Violent Odds — 3445.1 (-8.2)

**Hunter SV**
1. Dragonspine Trophy — 2779.5
2. Tsunami Talisman — 2766.7 (-12.8)
3. Madness of the Betrayer — 2765.2 (-1.5)
4. Abacus of Violent Odds — 2754.9 (-10.3)

**Hunter SV (Weave)**
1. Dragonspine Trophy — 3252.1
2. Madness of the Betrayer — 3200.6 (-51.5)
3. Abacus of Violent Odds — 3199.8 (-0.8)
4. Tsunami Talisman — 3198.2 (-1.6)

**Enhancement Shaman**
1. Dragonspine Trophy — 2358.1
2. Madness of the Betrayer — 2326.9 (-31.2)
3. Tsunami Talisman — 2316.7 (-10.2)
4. Abacus of Violent Odds — 2310.5 (-6.2)

**Feral Druid**
1. Tsunami Talisman — 2480.9
2. Dragonspine Trophy — 2474.0 (-6.9)
3. Abacus of Violent Odds — 2469.2 (-4.8)
4. Madness of the Betrayer — 2467.7 (-1.5)

**Retribution Paladin**
1. Dragonspine Trophy — 1932.1
2. Abacus of Violent Odds — 1895.6 (-36.5)
3. Tsunami Talisman — 1885.5 (-10.1)
4. Madness of the Betrayer — 1881.0 (-4.5)

### Cleave summary

The trinket ranking is unchanged from the single-target scenario: Dragonspine Trophy still
wins 12 of 13 categories, with Feral Druid again the sole exception favoring Tsunami
Talisman (by ~7 DPS, same as single-target — Feral melee abilities here don't cleave, so
its numbers barely move). Absolute DPS (Dragonspine Trophy row, single-target → cleave):

| Spec | Single-target DPS | Cleave DPS | Change |
|---|---:|---:|---:|
| Warrior Arms | 2048.9 | 3218.0 | +57.1% |
| Warrior Fury | 2246.8 | 3431.7 | +52.7% |
| Warrior Arms (Glaives) | 1989.4 | 2982.0 | +49.9% |
| Warrior Fury (Glaives) | 2364.4 | 3515.7 | +48.7% |
| Rogue (Glaives) | 2406.2 | 2890.1 | +20.1% |
| Rogue | 2235.4 | 2670.7 | +19.5% |
| Hunter SV | 2587.7 | 2779.5 | +7.4% |
| Hunter SV (Weave) | 3051.0 | 3252.1 | +6.6% |
| Hunter BM | 2931.8 | 3091.8 | +5.5% |
| Enhancement Shaman | 2246.2 | 2358.1 | +5.0% |
| Hunter BM (Weave) | 3338.1 | 3506.8 | +5.1% |
| Retribution Paladin | 1904.9 | 1932.1 | +1.4% |
| Feral Druid | 2458.8 | 2474.0 | +0.6% |

Warrior Fury/Arms see by far the largest gains (~48-57%), since Whirlwind is a core,
frequently-cast part of their rotation and now cleaves rage-free damage onto a full second
target. Rogue gains a more moderate ~20%, consistent with a talent/ability (e.g. Blade
Flurry-style cleave) splitting a portion of melee damage across both targets. Hunters and
Enhancement Shaman gain a modest ~5-7%, likely from partial multi-target effects (e.g.
Multi-Shot, cleave-adjacent procs) rather than a fully cleaving core rotation. Retribution
Paladin and Feral Druid are essentially flat (<2%), since neither kit used here has
meaningful passive multi-target damage. Since no spec's APL was re-tuned for an AOE
priority list, these numbers should be read as a floor, not a ceiling, on realistic 2-target
cleave performance.
