# Phase 3 Trinket Benchmark Results

20,000-iteration simulation comparing 4 candidate Trinket2 options across 14 physical-DPS
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
    fixed during this project, then later re-synced to a corrected BiS export — see "Bugs
    found & fixed" below for both)
- **Weapons**: the legendary Warglaives of Azzinoth (32837/32838) are excluded from the
  primary Warrior Fury / Rogue rows per project constraint; a non-Warglaive substitute is
  used instead (S3 Vengeful Gladiator weapons). A second "(Glaives)" row is also included
  for Warrior Fury/Arms and Rogue showing DPS *with* actual Warglaives equipped, for
  reference:
  - Warrior Fury (non-Glaive): MH Vengeful Gladiator's Cleaver (33669), OH Vengeful
    Gladiator's Chopper (34015)
  - Rogue (non-Glaive): MH Vengeful Gladiator's Slicer (33762), OH Blade of Savagery (32369)
  - "(Glaives)" rows: MH/OH Warglaives of Azzinoth (32837/32838)
  - Warrior Arms (base row) uses a two-hander (Cataclysm's Edge, 30902) with standard Arms
    talents. Warrior Arms (Glaives) and Warrior Arms (Kebab) are both dual-wield instead,
    and both use the dual-wield "Kebab" talent build
    (`34005021302010510321-0550000520501203`, invests in Dual Wield Specialization) rather
    than the 2H Arms talents — Glaives pairs those talents with actual Warglaives of
    Azzinoth, Kebab pairs them with the same affordable S3 swords as the Rogue non-Glaive
    row (Vengeful Gladiator's Slicer 33762 / Blade of Savagery 32369), for a cheap-dual-
    wield vs. legendary-dual-wield vs. 2H comparison. (The Glaives talent mismatch was a
    bug found & fixed during this project — see below.)
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
3. **`Warrior Arms (Glaives)` used the plain 2H Arms talent string despite dual-wielding
   Warglaives.** That meant 0 points in Dual Wield Specialization while wielding two
   one-handers — a genuine build/gear mismatch, most likely explaining why it (1932.1 DPS
   on Dragonspine Trophy) underperformed the plain 2H `Warrior Arms` row (2048.9 DPS)
   despite equipping top-tier legendary weapons. Fixed by switching it to the dual-wield
   "Kebab" talent build (`34005021302010510321-0550000520501203`, already defined in
   `ui/warrior/dps/presets.ts` as `ArmsKebabTalents`). DPS rose to 1984.2-2048.1
   (single-target) / 3024.6-3128.7 (cleave), now correctly at parity with the 2H row
   instead of below it. A new `Warrior Arms (Kebab)` row was also added using the same
   talent build with affordable S3 swords instead of Warglaives, for a 3-way comparison.
4. **Retribution Paladin gear further refined mid-project.** The gems used to fill the
   empty sockets in fix #2 were Phase 1 quality (a default pick at the time); a corrected
   BiS export supplied later in the project upgraded all 3 gem types to Phase 3 quality
   and swapped the chest piece (Lightbringer Breastplate → Midnight Chestguard). DPS rose
   again, by a consistent +25 to +26 across every trinket and scenario (1895-1931
   single-target, 1921-1958 cleave).
5. **`IndividualBuffs` were documented but never actually applied.** The Assumptions
   section always listed `FullIndividualBuffs` alongside the other three buff categories,
   but the benchmark's `player := &proto.Player{...}` construction never set a `Buffs`
   field — every spec was silently missing Blessing of Kings, Blessing of Wisdom/Might,
   and Unleashed Rage. Fixed by adding `Buffs: core.FullIndividualBuffs` to the player
   construction (applies uniformly to all 14 specs, so the comparison stays fair). DPS
   rose substantially across the board — roughly +11% to +17% for most specs, with
   Enhancement Shaman gaining noticeably less (~+8%), likely because it already
   self-provides an equivalent Unleashed-Rage-style AP buff via its own talents. All
   tables below reflect this fix.

## Full results table (20,000 iterations, seed 101)

| Spec | Trinket2 | DPS | Hit Rating |
|---|---|---:|---:|
| Warrior Fury | Abacus of Violent Odds | 2553.5 | 225.0 |
| Warrior Fury | Madness of the Betrayer | 2585.4 | 245.0 |
| Warrior Fury | Dragonspine Trophy | 2633.4 | 225.0 |
| Warrior Fury | Tsunami Talisman | 2579.7 | 235.0 |
| Warrior Fury (Glaives) | Abacus of Violent Odds | 2681.3 | 230.0 |
| Warrior Fury (Glaives) | Madness of the Betrayer | 2716.2 | 250.0 |
| Warrior Fury (Glaives) | Dragonspine Trophy | 2763.0 | 230.0 |
| Warrior Fury (Glaives) | Tsunami Talisman | 2712.9 | 240.0 |
| Warrior Arms | Abacus of Violent Odds | 2299.3 | 121.0 |
| Warrior Arms | Madness of the Betrayer | 2315.8 | 141.0 |
| Warrior Arms | Dragonspine Trophy | 2356.2 | 121.0 |
| Warrior Arms | Tsunami Talisman | 2313.7 | 131.0 |
| Warrior Arms (Glaives) | Abacus of Violent Odds | 2287.0 | 142.0 |
| Warrior Arms (Glaives) | Madness of the Betrayer | 2318.6 | 162.0 |
| Warrior Arms (Glaives) | Dragonspine Trophy | 2358.9 | 142.0 |
| Warrior Arms (Glaives) | Tsunami Talisman | 2314.7 | 152.0 |
| Warrior Arms (Kebab) | Abacus of Violent Odds | 2145.8 | 144.0 |
| Warrior Arms (Kebab) | Madness of the Betrayer | 2176.9 | 164.0 |
| Warrior Arms (Kebab) | Dragonspine Trophy | 2215.8 | 144.0 |
| Warrior Arms (Kebab) | Tsunami Talisman | 2173.6 | 154.0 |
| Rogue | Abacus of Violent Odds | 2495.3 | 348.0 |
| Rogue | Madness of the Betrayer | 2487.6 | 368.0 |
| Rogue | Dragonspine Trophy | 2563.8 | 348.0 |
| Rogue | Tsunami Talisman | 2485.5 | 358.0 |
| Rogue (Glaives) | Abacus of Violent Odds | 2687.2 | 346.0 |
| Rogue (Glaives) | Madness of the Betrayer | 2680.9 | 366.0 |
| Rogue (Glaives) | Dragonspine Trophy | 2757.6 | 346.0 |
| Rogue (Glaives) | Tsunami Talisman | 2679.3 | 356.0 |
| Hunter BM | Abacus of Violent Odds | 3249.3 | 109.0 |
| Hunter BM | Madness of the Betrayer | 3258.5 | 129.0 |
| Hunter BM | Dragonspine Trophy | 3283.6 | 109.0 |
| Hunter BM | Tsunami Talisman | 3258.6 | 119.0 |
| Hunter BM (Weave) | Abacus of Violent Odds | 3707.0 | 109.0 |
| Hunter BM (Weave) | Madness of the Betrayer | 3713.0 | 129.0 |
| Hunter BM (Weave) | Dragonspine Trophy | 3766.9 | 109.0 |
| Hunter BM (Weave) | Tsunami Talisman | 3711.7 | 119.0 |
| Hunter SV | Abacus of Violent Odds | 2844.6 | 58.0 |
| Hunter SV | Madness of the Betrayer | 2845.6 | 78.0 |
| Hunter SV | Dragonspine Trophy | 2871.5 | 58.0 |
| Hunter SV | Tsunami Talisman | 2849.1 | 68.0 |
| Hunter SV (Weave) | Abacus of Violent Odds | 3369.1 | 58.0 |
| Hunter SV (Weave) | Madness of the Betrayer | 3366.1 | 78.0 |
| Hunter SV (Weave) | Dragonspine Trophy | 3430.8 | 58.0 |
| Hunter SV (Weave) | Tsunami Talisman | 3366.2 | 68.0 |
| Enhancement Shaman | Abacus of Violent Odds | 2367.6 | 183.0 |
| Enhancement Shaman | Madness of the Betrayer | 2382.6 | 203.0 |
| Enhancement Shaman | Dragonspine Trophy | 2417.4 | 183.0 |
| Enhancement Shaman | Tsunami Talisman | 2372.8 | 193.0 |
| Feral Druid | Abacus of Violent Odds | 2822.5 | 116.0 |
| Feral Druid | Madness of the Betrayer | 2824.6 | 136.0 |
| Feral Druid | Dragonspine Trophy | 2828.0 | 116.0 |
| Feral Druid | Tsunami Talisman | 2839.7 | 126.0 |
| Retribution Paladin | Abacus of Violent Odds | 2174.2 | 80.0 |
| Retribution Paladin | Madness of the Betrayer | 2154.9 | 100.0 |
| Retribution Paladin | Dragonspine Trophy | 2216.5 | 80.0 |
| Retribution Paladin | Tsunami Talisman | 2160.6 | 90.0 |

## Trinket rankings per spec (best to worst, with DPS delta to next rank)

**Warrior Fury**
1. Dragonspine Trophy — 2633.4
2. Madness of the Betrayer — 2585.4 (-48.0)
3. Tsunami Talisman — 2579.7 (-5.7)
4. Abacus of Violent Odds — 2553.5 (-26.2)

**Warrior Fury (Glaives)**
1. Dragonspine Trophy — 2763.0
2. Madness of the Betrayer — 2716.2 (-46.8)
3. Tsunami Talisman — 2712.9 (-3.3)
4. Abacus of Violent Odds — 2681.3 (-31.6)

**Warrior Arms**
1. Dragonspine Trophy — 2356.2
2. Madness of the Betrayer — 2315.8 (-40.4)
3. Tsunami Talisman — 2313.7 (-2.1)
4. Abacus of Violent Odds — 2299.3 (-14.4)

**Warrior Arms (Glaives)**
1. Dragonspine Trophy — 2358.9
2. Madness of the Betrayer — 2318.6 (-40.3)
3. Tsunami Talisman — 2314.7 (-3.9)
4. Abacus of Violent Odds — 2287.0 (-27.7)

**Warrior Arms (Kebab)**
1. Dragonspine Trophy — 2215.8
2. Madness of the Betrayer — 2176.9 (-38.9)
3. Tsunami Talisman — 2173.6 (-3.3)
4. Abacus of Violent Odds — 2145.8 (-27.8)

**Rogue**
1. Dragonspine Trophy — 2563.8
2. Abacus of Violent Odds — 2495.3 (-68.5)
3. Madness of the Betrayer — 2487.6 (-7.7)
4. Tsunami Talisman — 2485.5 (-2.1)

**Rogue (Glaives)**
1. Dragonspine Trophy — 2757.6
2. Abacus of Violent Odds — 2687.2 (-70.4)
3. Madness of the Betrayer — 2680.9 (-6.3)
4. Tsunami Talisman — 2679.3 (-1.6)

**Hunter BM**
1. Dragonspine Trophy — 3283.6
2. Tsunami Talisman — 3258.6 (-25.0)
3. Madness of the Betrayer — 3258.5 (-0.1)
4. Abacus of Violent Odds — 3249.3 (-9.2)

**Hunter BM (Weave)**
1. Dragonspine Trophy — 3766.9
2. Madness of the Betrayer — 3713.0 (-53.9)
3. Tsunami Talisman — 3711.7 (-1.3)
4. Abacus of Violent Odds — 3707.0 (-4.7)

**Hunter SV**
1. Dragonspine Trophy — 2871.5
2. Tsunami Talisman — 2849.1 (-22.4)
3. Madness of the Betrayer — 2845.6 (-3.5)
4. Abacus of Violent Odds — 2844.6 (-1.0)

**Hunter SV (Weave)**
1. Dragonspine Trophy — 3430.8
2. Abacus of Violent Odds — 3369.1 (-61.7)
3. Tsunami Talisman — 3366.2 (-2.9)
4. Madness of the Betrayer — 3366.1 (-0.1)

**Enhancement Shaman**
1. Dragonspine Trophy — 2417.4
2. Madness of the Betrayer — 2382.6 (-34.8)
3. Tsunami Talisman — 2372.8 (-9.8)
4. Abacus of Violent Odds — 2367.6 (-5.2)

**Feral Druid**
1. Tsunami Talisman — 2839.7
2. Dragonspine Trophy — 2828.0 (-11.7)
3. Madness of the Betrayer — 2824.6 (-3.4)
4. Abacus of Violent Odds — 2822.5 (-2.1)

**Retribution Paladin**
1. Dragonspine Trophy — 2216.5
2. Abacus of Violent Odds — 2174.2 (-42.3)
3. Tsunami Talisman — 2160.6 (-13.6)
4. Madness of the Betrayer — 2154.9 (-5.7)

## Summary

Dragonspine Trophy is the best Trinket2 choice in 13 of 14 spec categories, typically by a
clear 20-70 DPS margin over the runner-up. Feral Druid is the sole exception, where Tsunami
Talisman edges out Dragonspine Trophy by a small ~12 DPS — effectively a toss-up within
simulation noise.

## Cleave scenario (2 targets, 20,000 iterations, seed 101)

Same 14 spec/gear configurations and same 4 Trinket2 candidates as above, but the encounter
was changed to 2 targets instead of 1 (`Encounter.Targets` extended to
`[NewDefaultTarget(), NewDefaultTarget()]`, both full-health dummies, same 180s ±5s
duration/buffs/debuffs). Each spec runs its existing single-target-tuned APL as-is, so any
DPS gains come purely from incidental cleave/AOE effects already present in that rotation
(e.g. Whirlwind, Cleave, Swipe, multi-target-hitting procs), not from a dedicated AOE
priority list. Treat this as a "how much do I passively benefit from a second target"
estimate rather than an optimized cleave/AOE rotation benchmark — with one exception:
Hunter received one small, verified cleave-specific APL tweak (see
[Cleave-tuned APL experiment](#cleave-tuned-apl-experiment) below), so the Hunter BM/SV rows
below already include it.

### Full results table (cleave, 2 targets)

| Spec | Trinket2 | DPS | Hit Rating |
|---|---|---:|---:|
| Warrior Fury | Abacus of Violent Odds | 3948.1 | 225.0 |
| Warrior Fury | Madness of the Betrayer | 4015.9 | 245.0 |
| Warrior Fury | Dragonspine Trophy | 4082.9 | 225.0 |
| Warrior Fury | Tsunami Talisman | 3998.5 | 235.0 |
| Warrior Fury (Glaives) | Abacus of Violent Odds | 4035.7 | 230.0 |
| Warrior Fury (Glaives) | Madness of the Betrayer | 4108.8 | 250.0 |
| Warrior Fury (Glaives) | Dragonspine Trophy | 4169.4 | 230.0 |
| Warrior Fury (Glaives) | Tsunami Talisman | 4097.0 | 240.0 |
| Warrior Arms | Abacus of Violent Odds | 3666.8 | 121.0 |
| Warrior Arms | Madness of the Betrayer | 3708.0 | 141.0 |
| Warrior Arms | Dragonspine Trophy | 3782.0 | 121.0 |
| Warrior Arms | Tsunami Talisman | 3697.9 | 131.0 |
| Warrior Arms (Glaives) | Abacus of Violent Odds | 3567.2 | 142.0 |
| Warrior Arms (Glaives) | Madness of the Betrayer | 3631.6 | 162.0 |
| Warrior Arms (Glaives) | Dragonspine Trophy | 3691.4 | 142.0 |
| Warrior Arms (Glaives) | Tsunami Talisman | 3616.5 | 152.0 |
| Warrior Arms (Kebab) | Abacus of Violent Odds | 3287.1 | 144.0 |
| Warrior Arms (Kebab) | Madness of the Betrayer | 3355.0 | 164.0 |
| Warrior Arms (Kebab) | Dragonspine Trophy | 3405.7 | 144.0 |
| Warrior Arms (Kebab) | Tsunami Talisman | 3341.8 | 154.0 |
| Rogue | Abacus of Violent Odds | 3012.9 | 348.0 |
| Rogue | Madness of the Betrayer | 2972.4 | 368.0 |
| Rogue | Dragonspine Trophy | 3066.7 | 348.0 |
| Rogue | Tsunami Talisman | 2970.9 | 358.0 |
| Rogue (Glaives) | Abacus of Violent Odds | 3259.8 | 346.0 |
| Rogue (Glaives) | Madness of the Betrayer | 3219.4 | 366.0 |
| Rogue (Glaives) | Dragonspine Trophy | 3315.6 | 346.0 |
| Rogue (Glaives) | Tsunami Talisman | 3219.3 | 356.0 |
| Hunter BM | Abacus of Violent Odds | 3440.3 | 109.0 |
| Hunter BM | Madness of the Betrayer | 3452.1 | 129.0 |
| Hunter BM | Dragonspine Trophy | 3478.5 | 109.0 |
| Hunter BM | Tsunami Talisman | 3453.2 | 119.0 |
| Hunter BM (Weave) | Abacus of Violent Odds | 3895.1 | 109.0 |
| Hunter BM (Weave) | Madness of the Betrayer | 3905.3 | 129.0 |
| Hunter BM (Weave) | Dragonspine Trophy | 3965.3 | 109.0 |
| Hunter BM (Weave) | Tsunami Talisman | 3905.4 | 119.0 |
| Hunter SV | Abacus of Violent Odds | 3043.7 | 58.0 |
| Hunter SV | Madness of the Betrayer | 3049.0 | 78.0 |
| Hunter SV | Dragonspine Trophy | 3075.1 | 58.0 |
| Hunter SV | Tsunami Talisman | 3050.7 | 68.0 |
| Hunter SV (Weave) | Abacus of Violent Odds | 3584.8 | 58.0 |
| Hunter SV (Weave) | Madness of the Betrayer | 3587.4 | 78.0 |
| Hunter SV (Weave) | Dragonspine Trophy | 3649.0 | 58.0 |
| Hunter SV (Weave) | Tsunami Talisman | 3587.4 | 68.0 |
| Enhancement Shaman | Abacus of Violent Odds | 2480.0 | 183.0 |
| Enhancement Shaman | Madness of the Betrayer | 2496.3 | 203.0 |
| Enhancement Shaman | Dragonspine Trophy | 2532.0 | 183.0 |
| Enhancement Shaman | Tsunami Talisman | 2485.1 | 193.0 |
| Feral Druid | Abacus of Violent Odds | 2839.2 | 116.0 |
| Feral Druid | Madness of the Betrayer | 2839.9 | 136.0 |
| Feral Druid | Dragonspine Trophy | 2843.1 | 116.0 |
| Feral Druid | Tsunami Talisman | 2855.5 | 126.0 |
| Retribution Paladin | Abacus of Violent Odds | 2200.7 | 80.0 |
| Retribution Paladin | Madness of the Betrayer | 2182.7 | 100.0 |
| Retribution Paladin | Dragonspine Trophy | 2243.6 | 80.0 |
| Retribution Paladin | Tsunami Talisman | 2187.8 | 90.0 |

### Trinket rankings per spec, cleave scenario (best to worst, with DPS delta to next rank)

**Warrior Fury**
1. Dragonspine Trophy — 4082.9
2. Madness of the Betrayer — 4015.9 (-67.0)
3. Tsunami Talisman — 3998.5 (-17.4)
4. Abacus of Violent Odds — 3948.1 (-50.4)

**Warrior Fury (Glaives)**
1. Dragonspine Trophy — 4169.4
2. Madness of the Betrayer — 4108.8 (-60.6)
3. Tsunami Talisman — 4097.0 (-11.8)
4. Abacus of Violent Odds — 4035.7 (-61.3)

**Warrior Arms**
1. Dragonspine Trophy — 3782.0
2. Madness of the Betrayer — 3708.0 (-74.0)
3. Tsunami Talisman — 3697.9 (-10.1)
4. Abacus of Violent Odds — 3666.8 (-31.1)

**Warrior Arms (Glaives)**
1. Dragonspine Trophy — 3691.4
2. Madness of the Betrayer — 3631.6 (-59.8)
3. Tsunami Talisman — 3616.5 (-15.1)
4. Abacus of Violent Odds — 3567.2 (-49.3)

**Warrior Arms (Kebab)**
1. Dragonspine Trophy — 3405.7
2. Madness of the Betrayer — 3355.0 (-50.7)
3. Tsunami Talisman — 3341.8 (-13.2)
4. Abacus of Violent Odds — 3287.1 (-54.7)

**Rogue**
1. Dragonspine Trophy — 3066.7
2. Abacus of Violent Odds — 3012.9 (-53.8)
3. Madness of the Betrayer — 2972.4 (-40.5)
4. Tsunami Talisman — 2970.9 (-1.5)

**Rogue (Glaives)**
1. Dragonspine Trophy — 3315.6
2. Abacus of Violent Odds — 3259.8 (-55.8)
3. Madness of the Betrayer — 3219.4 (-40.4)
4. Tsunami Talisman — 3219.3 (-0.1)

**Hunter BM**
1. Dragonspine Trophy — 3478.5
2. Tsunami Talisman — 3453.2 (-25.3)
3. Madness of the Betrayer — 3452.1 (-1.1)
4. Abacus of Violent Odds — 3440.3 (-11.8)

**Hunter BM (Weave)**
1. Dragonspine Trophy — 3965.3
2. Tsunami Talisman — 3905.4 (-59.9)
3. Madness of the Betrayer — 3905.3 (-0.1)
4. Abacus of Violent Odds — 3895.1 (-10.2)

**Hunter SV**
1. Dragonspine Trophy — 3075.1
2. Tsunami Talisman — 3050.7 (-24.4)
3. Madness of the Betrayer — 3049.0 (-1.7)
4. Abacus of Violent Odds — 3043.7 (-5.3)

**Hunter SV (Weave)**
1. Dragonspine Trophy — 3649.0
2. Madness of the Betrayer — 3587.4 (-61.6)
3. Tsunami Talisman — 3587.4 (-0.0, tied)
4. Abacus of Violent Odds — 3584.8 (-2.6)

**Enhancement Shaman**
1. Dragonspine Trophy — 2532.0
2. Madness of the Betrayer — 2496.3 (-35.7)
3. Tsunami Talisman — 2485.1 (-11.2)
4. Abacus of Violent Odds — 2480.0 (-5.1)

**Feral Druid**
1. Tsunami Talisman — 2855.5
2. Dragonspine Trophy — 2843.1 (-12.4)
3. Madness of the Betrayer — 2839.9 (-3.2)
4. Abacus of Violent Odds — 2839.2 (-0.7)

**Retribution Paladin**
1. Dragonspine Trophy — 2243.6
2. Abacus of Violent Odds — 2200.7 (-42.9)
3. Tsunami Talisman — 2187.8 (-12.9)
4. Madness of the Betrayer — 2182.7 (-5.1)

### Cleave summary

The trinket ranking is unchanged from the single-target scenario: Dragonspine Trophy still
wins 13 of 14 categories, with Feral Druid again the sole exception favoring Tsunami
Talisman (by ~12 DPS in both scenarios — Feral melee abilities here don't cleave, so its
numbers barely move). Absolute DPS (Dragonspine Trophy row, single-target → cleave):

| Spec | Single-target DPS | Cleave DPS | Change |
|---|---:|---:|---:|
| Warrior Arms | 2356.2 | 3782.0 | +60.5% |
| Warrior Arms (Glaives) | 2358.9 | 3691.4 | +56.5% |
| Warrior Fury | 2633.4 | 4082.9 | +55.0% |
| Warrior Arms (Kebab) | 2215.8 | 3405.7 | +53.7% |
| Warrior Fury (Glaives) | 2763.0 | 4169.4 | +50.9% |
| Rogue (Glaives) | 2757.6 | 3315.6 | +20.2% |
| Rogue | 2563.8 | 3066.7 | +19.6% |
| Hunter SV | 2871.5 | 3075.1 | +7.1% |
| Hunter SV (Weave) | 3430.8 | 3649.0 | +6.4% |
| Hunter BM | 3283.6 | 3478.5 | +5.9% |
| Hunter BM (Weave) | 3766.9 | 3965.3 | +5.3% |
| Enhancement Shaman | 2417.4 | 2532.0 | +4.7% |
| Retribution Paladin | 2216.5 | 2243.6 | +1.2% |
| Feral Druid | 2828.0 | 2843.1 | +0.5% |

Warrior Fury/Arms see by far the largest gains (~51-61%), since Whirlwind is a core,
frequently-cast part of their rotation and now cleaves rage-free damage onto a full second
target. Rogue gains a more moderate ~20%, consistent with a talent/ability (e.g. Blade
Flurry-style cleave) splitting a portion of melee damage across both targets. Hunters and
Enhancement Shaman gain a modest ~5-7%, likely from partial multi-target effects (e.g.
Multi-Shot, cleave-adjacent procs) rather than a fully cleaving core rotation (Hunter's
share already includes the small kept APL tweak described below). Retribution Paladin and
Feral Druid are essentially flat (<2%), since neither kit used here has meaningful passive
multi-target damage. Aside from Hunter's one small kept tweak, no other spec's APL was
re-tuned for an AOE priority list, so these numbers should be read as a floor, not a
ceiling, on realistic 2-target cleave performance.

## Cleave-tuned APL experiment

*Note: this experiment was originally run at 20,000 iterations before the `IndividualBuffs`
fix (bug #5) and, for Retribution Paladin, before the later gear re-sync (bug #4). It has
since been redone at 5,000 iterations (seed 101) on the current, fully-fixed benchmark
state, so the numbers below are now consistent with the rest of this document — no stale
pre-fix figures remain.*

Starting from the passive baseline above, three specs had genuine, already-implemented
cleave/AOE abilities that looked under-utilized in their default APL, so their rotation
files were edited to lean on those abilities specifically when `numberTargets > 1` (single-
target behavior provably unchanged, since the new logic is gated behind that condition or
additive to it):

- **Rogue** (`swords.apl.json`): widened the existing Blade Flurry condition with
  `OR numberTargets > 1`, so Blade Flurry (already talented, already cast under a buff-
  uptime condition) would also fire purely based on target count.
- **Retribution Paladin** (`default.apl.json`): flipped the already-wired but permanently-
  off `"Use Consecrate"` toggle to `numberTargets > 1`, activating Consecration (rank 6)
  only during cleave.
- **Hunter** (`default.apl.json`, shared by BM/SV): added a new top-priority Multi-Shot
  cast attempt gated on `numberTargets > 1` with a simple auto-attack-clipping safety
  check, supplementing the existing melee-weave-gated Multi-Shot logic.

Each change was re-benchmarked with the same 2-target cleave harness (5,000 iterations,
seed 101) and compared against the passive baseline above:

| Spec | Trinket2 | Baseline cleave DPS | Tuned cleave DPS | Change |
|---|---|---:|---:|---:|
| Rogue | Abacus of Violent Odds | 3011.6 | 2954.0 | -1.91% |
| Rogue | Madness of the Betrayer | 2974.0 | 2931.3 | -1.44% |
| Rogue | Dragonspine Trophy | 3067.0 | 3022.0 | -1.47% |
| Rogue | Tsunami Talisman | 2969.3 | 2936.1 | -1.12% |
| Retribution Paladin | Abacus of Violent Odds | 2200.6 | 2163.7 | -1.68% |
| Retribution Paladin | Madness of the Betrayer | 2183.6 | 2149.8 | -1.55% |
| Retribution Paladin | Dragonspine Trophy | 2243.3 | 2194.6 | -2.17% |
| Retribution Paladin | Tsunami Talisman | 2188.9 | 2150.9 | -1.74% |
| Hunter BM | Abacus of Violent Odds | 3427.8 | 3440.2 | +0.36% |
| Hunter BM | Madness of the Betrayer | 3436.0 | 3452.1 | +0.47% |
| Hunter BM | Dragonspine Trophy | 3454.5 | 3477.2 | +0.66% |
| Hunter BM | Tsunami Talisman | 3436.1 | 3453.6 | +0.51% |
| Hunter BM (Weave) | Abacus of Violent Odds | 3875.7 | 3896.2 | +0.53% |
| Hunter BM (Weave) | Madness of the Betrayer | 3888.2 | 3907.3 | +0.49% |
| Hunter BM (Weave) | Dragonspine Trophy | 3943.4 | 3965.5 | +0.56% |
| Hunter BM (Weave) | Tsunami Talisman | 3882.9 | 3906.6 | +0.61% |
| Hunter SV | Abacus of Violent Odds | 3049.6 | 3044.9 | -0.15% |
| Hunter SV | Madness of the Betrayer | 3057.3 | 3049.8 | -0.25% |
| Hunter SV | Dragonspine Trophy | 3075.1 | 3074.8 | -0.01% |
| Hunter SV | Tsunami Talisman | 3059.6 | 3050.0 | -0.31% |
| Hunter SV (Weave) | Abacus of Violent Odds | 3581.6 | 3584.2 | +0.07% |
| Hunter SV (Weave) | Madness of the Betrayer | 3584.0 | 3587.0 | +0.08% |
| Hunter SV (Weave) | Dragonspine Trophy | 3641.5 | 3648.6 | +0.19% |
| Hunter SV (Weave) | Tsunami Talisman | 3580.7 | 3587.4 | +0.19% |

**Findings:**

- **Rogue: reverted.** Widening Blade Flurry's condition made DPS consistently *worse*
  (~1.1-1.9%) instead of better. The original Slice-and-Dice + Expose-Armor-gated condition
  apparently times Blade Flurry into an already-convenient GCD window; forcing it the
  instant it comes off cooldown (regardless of buff state) instead displaces a combo-point
  finisher at a worse moment. `swords.apl.json` reverted to its original condition.
- **Retribution Paladin: reverted (at 2 targets only — see 4-target follow-up below).**
  Turning on Consecration during cleave made DPS consistently *worse* by ~1.5-2.2% at 2
  targets — smaller in magnitude than the original pre-fix run's ~5% regression (the
  current, stronger buff/gear baseline makes Consecrate's fixed GCD/mana tax proportionally
  cheaper, but it's still a clear, consistent loss at this target count, not a gain).
  Consecration's GCD-length cast plus mana cost crowds out higher-value melee
  actions/Judgement in this rotation's priority list at only 2 targets; the ability's low
  base damage doesn't make up for the lost tempo until more enemies are in its radius (see
  below). `default.apl.json` (Paladin) reverted to `"Use Consecrate": false`.
- **Hunter: kept.** The extra top-priority Multi-Shot attempt is a small, consistent, real
  gain for Beast Mastery (+0.36 to +0.66% across all trinkets and both weave variants) and
  essentially flat/noise-level for Survival (-0.31% to +0.19% across both variants). Since
  BM's gain outweighs SV's negligible change and neither regresses meaningfully, the change
  was kept in `default.apl.json` (Hunter). Single-target DPS is provably unaffected — the
  full single-target benchmark reproduces the current single-target Hunter numbers (e.g.
  Hunter BM/Dragonspine Trophy: 3283.6, Hunter BM (Weave): 3766.9, Hunter SV: 2871.5, Hunter
  SV (Weave): 3430.8), since the new action is gated behind `numberTargets > 1`.
- **Enhancement Shaman: no change made.** Its APL already unconditionally uses Magma Totem
  and twists Fire Nova Totem regardless of target count, so it already gets full cleave
  benefit with no APL edit needed.
- **Feral Druid: no change possible.** Swipe (the only multi-target melee ability in this
  codebase) is Bear Form-only, and the Cat Form rotation's Go implementation
  (`sim/druid/feralcat/rotation.go`) only ever targets `CurrentTarget` — there is no
  AOE/cleave lever available to a Cat Form APL in this sim.
- **Warrior: no change needed.** Its APL already branches on `numberTargets` (Whirlwind/
  Bloodthirst prioritized at 2+ targets, single-target-only logic like Execute/Sunder gated
  to 1 target).

### 4-target follow-up

To check whether these conclusions hold at a higher target count — and specifically
whether Consecrate's radius-based AOE finally becomes worthwhile with more nearby enemies —
the same three tweaks were re-benchmarked against a 4-target encounter instead of 2, using a
temporary, not-committed test scenario (still 5,000 iterations, seed 101):

| Spec | Trinket2 | Baseline cleave DPS | Tuned cleave DPS | Change |
|---|---|---:|---:|---:|
| Rogue | Abacus of Violent Odds | 3013.5 | 2954.2 | -1.97% |
| Rogue | Madness of the Betrayer | 2971.3 | 2932.5 | -1.31% |
| Rogue | Dragonspine Trophy | 3065.5 | 3020.3 | -1.47% |
| Rogue | Tsunami Talisman | 2970.7 | 2934.8 | -1.21% |
| Retribution Paladin | Abacus of Violent Odds | 2255.6 | 2330.2 | +3.31% |
| Retribution Paladin | Madness of the Betrayer | 2237.9 | 2318.7 | +3.61% |
| Retribution Paladin | Dragonspine Trophy | 2296.1 | 2357.5 | +2.67% |
| Retribution Paladin | Tsunami Talisman | 2233.8 | 2323.8 | +4.03% |
| Hunter BM | Abacus of Violent Odds | 3786.9 | 3870.3 | +2.20% |
| Hunter BM | Madness of the Betrayer | 3789.7 | 3889.5 | +2.63% |
| Hunter BM | Dragonspine Trophy | 3792.4 | 3913.4 | +3.19% |
| Hunter BM | Tsunami Talisman | 3789.5 | 3888.1 | +2.60% |
| Hunter BM (Weave) | Abacus of Violent Odds | 4217.1 | 4322.9 | +2.51% |
| Hunter BM (Weave) | Madness of the Betrayer | 4228.3 | 4341.9 | +2.69% |
| Hunter BM (Weave) | Dragonspine Trophy | 4303.1 | 4400.7 | +2.27% |
| Hunter BM (Weave) | Tsunami Talisman | 4226.4 | 4339.7 | +2.68% |
| Hunter SV | Abacus of Violent Odds | 3458.7 | 3514.4 | +1.61% |
| Hunter SV | Madness of the Betrayer | 3482.8 | 3523.9 | +1.18% |
| Hunter SV | Dragonspine Trophy | 3484.5 | 3547.0 | +1.79% |
| Hunter SV | Tsunami Talisman | 3483.6 | 3528.8 | +1.30% |
| Hunter SV (Weave) | Abacus of Violent Odds | 4005.9 | 4052.1 | +1.15% |
| Hunter SV (Weave) | Madness of the Betrayer | 4010.4 | 4065.6 | +1.38% |
| Hunter SV (Weave) | Dragonspine Trophy | 4062.9 | 4122.8 | +1.47% |
| Hunter SV (Weave) | Tsunami Talisman | 4006.8 | 4064.3 | +1.44% |

- **Rogue: conclusion unchanged (reverted).** Blade Flurry cleaves to only one extra nearby
  target in TBC regardless of total enemy count, so widening its condition is essentially
  the same tweak at 2 or 4 targets — DPS is still consistently *worse* (~1.2-2.0%, matching
  the 2-target result almost exactly).
- **Retribution Paladin: conclusion flips (favors keeping the tweak at 4 targets).** Unlike
  at 2 targets (a consistent ~1.5-2.2% loss), enabling Consecrate during cleave at 4 targets
  is a consistent ~2.7-4.0% *gain*. Consecrate damages all enemies in an 8-yard radius, so
  its value scales with nearby target count in a way Blade Flurry's fixed one-extra-target
  cleave doesn't; 2 targets isn't enough to offset its GCD/mana cost, but 4 clearly is. The
  right condition is likely target-count-dependent (e.g. `numberTargets > 2`) rather than a
  flat keep/revert verdict — worth tuning further rather than treating either 2-target or
  4-target result as final.
- **Hunter: conclusion unchanged (kept), and the gain grows with target count.** The extra
  Multi-Shot attempt is a clearly larger win at 4 targets than at 2: Beast Mastery
  +2.2-3.2% (vs. +0.36-0.66% at 2 targets) and Survival +1.2-1.8% (vs. essentially flat at 2
  targets). Multi-Shot hits all nearby enemies, so like Consecrate its value scales with
  target count.

This follow-up used a temporary, uncommitted 4-target test scenario and temporary APL edits
only for the duration of the experiment — no permanent test file or APL changes resulted
from it, and the main cleave benchmark above remains 2 targets.

The full results table, trinket rankings, and cleave summary above already reflect the kept
Hunter change — every other spec there is byte-for-byte the same APL as the single-target
scenario.
