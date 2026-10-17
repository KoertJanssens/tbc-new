# Phase 3 upgrade report — Koerrii (Enhancement Shaman, Spineshatter-EU)

Sim-backed answer to "what are my biggest Phase 3 upgrades?", built from the live armory
profile at <https://ironforge.pro/anniversary/player/Spineshatter/Koerrii/> and run through
this repo's engine at 20,000 iterations per configuration.

**Current gear sims at 2238 DPS** (stdev 79). Everything below is a delta against that number.

## TL;DR

1. **Your trinkets are the whole story.** Terokkar Tablet of Precision (ilvl 99) and Bladefist's
   Breadth (ilvl 88) are the two worst items you own by a wide margin. A **Dragonspine Trophy**
   — a *Phase 1* Gruul drop — is worth **+78 DPS** on its own, more than any single Black Temple
   armour piece. Pairing it with **Madness of the Betrayer** is the single biggest chunk of
   your remaining potential (78 + 39 ≈ 116 DPS between them).
2. **Do the free stuff first: +35 DPS for zero drops.** Your boots carry Minor Speed Increase
   (no stats at all), your shoulders have the non-exalted inscription, your legs have Cobrahide
   instead of Nethercobra, and your bracer is socketed with a PvP resilience gem.
3. **Do not take Syphon of the Nathrezim.** The repo's own P3 BiS preset equips it, but for you
   it is a **-33 DPS main hand**. You are an Orc: Axe Specialization gives +5 expertise with
   axes, and the sim models it. Swapping to that mace zeroes your expertise. Stay on axes.
4. **Never swap a second tier piece until you are replacing four of them.** You wear 5/5
   Cataclysm Harness; the 4pc (+5% extra Flurry haste) survives one swap but not two. Dropping
   below 4 pieces costs roughly 55 DPS, which is more than any single armour upgrade gains.
5. Ceiling: with the best item this sim found in every slot plus the enchant/gem fixes, you reach
   **2679 DPS (+442, +20%)**.

## Priority order — what to chase, in order

Ranked by single-slot DPS gain, with unique-equipped conflicts already resolved (only one
Dragonspine Trophy, one Stormrage Signet Ring, etc.).

| # | Slot | Upgrade | ilvl | DPS gain | Where it comes from |
|---|------|---------|------|----------|---------------------|
| 1 | Trinket1 | Dragonspine Trophy | 125 | **+78** | Gruul the Dragonkiller — Gruul's Lair |
| 2 | Waist | Belt of One-Hundred Deaths | 138 | **+38** | Lady Vashj — Serpentshrine Cavern |
| 3 | Head | Cursed Vision of Sargeras | 151 | **+34** | Illidan Stormrage — Black Temple |
| 4 | Trinket2 | Madness of the Betrayer | 141 | **+33** | The Illidari Council — Black Temple |
| 5 | Wrist | Bindings of Lightning Reflexes | 141 | **+32** | Crafted (Leatherworking) |
| 6 | Shoulder | Swiftstrike Shoulders | 141 | **+31** | Crafted (Leatherworking) **[requires Leatherworking]** |
| 7 | Chest | Midnight Chestguard | 151 | **+27** | Archimonde — Hyjal Summit |
| 8 | MainHand | Vengeful Gladiator's Cleaver | 146 | **+25** | Season 3 arena vendor |
| 9 | Legs | Bow-stitched Leggings | 141 | **+24** | Azgalor — Hyjal Summit |
| 10 | Finger2 | Stormrage Signet Ring | 151 | **+24** | Illidan Stormrage — Black Temple |
| 11 | Neck | Choker of Endless Nightmares | 141 | **+19** | Supremus — Black Temple |
| 12 | Hands | Grips of Damnation | 141 | **+19** | Reliquary of the Lost — Black Temple |
| 13 | Feet | Shadowmaster's Boots | 141 | **+16** | Mother Shahraz — Black Temple |
| 14 | OffHand | Vengeful Gladiator's Cleaver | 146 | **+15** | Season 3 arena vendor |
| 15 | Back | Shadowmoon Destroyer's Drape | 141 | **+12** | Teron Gorefiend — Black Temple |
| 16 | Finger1 | Band of the Ranger-General | 138 | **+11** | Kael'thas Sunstrider — Tempest Keep |
| 17 | Ranged | Totem of the Astral Winds | 115 | **+10** | Pandemonius — Mana-Tombs |

The DPS gains are *independent* single-slot measurements: each one replaces exactly that item
and nothing else. They do not simply add up — see [Stacking them](#stacking-them-the-4-piece-t5-cliff).

## Free upgrades — no drops required

| Change | DPS | Gain | Detail |
|--------|-----|------|--------|
| Boots: none -> Cat's Swiftness | 2241 | **+3** | +6 agi (currently Minor Speed, 0 stats) |
| Boots: none -> Dexterity | 2245 | **+7** | +12 agi (no run speed) |
| Shoulder: Inscription -> Greater Inscription of Vengeance | 2246 | **+8** | +4 AP, +10 crit (Aldor/Scryer exalted) |
| Legs: Cobrahide -> Nethercobra Leg Armor | 2243 | **+5** | +10 AP, +2 crit |
| Wrist: swap PvP gem for Bold Crimson Spinel | 2246 | **+9** | Steady Talasite (sta/resil) -> +10 str |
| Regem everything with P3 epic gems | 2255 | **+18** | rare -> epic gems, colour-matched for socket bonuses |
| Regem everything, all red (ignore socket bonuses) | 2261 | **+23** | every non-meta socket = Bold Crimson Spinel (+10 str) |
| All enchant + gem fixes together | 2272 | **+35** | Cat's Swiftness + Greater Inscription + Nethercobra + colour-matched regem |

Notes:

- **Boots.** Minor Speed Increase gives no stats. Cat's Swiftness keeps the run speed and adds
  +6 agi (+3); Dexterity drops the speed for +12 agi (+7). Take Cat's Swiftness unless you
  really want the extra ~4 DPS.
- **Gems.** Straight Bold Crimson Spinel everywhere (+23) beats colour-matching for socket
  bonuses (+18) on your current gear — your socket bonuses are small and mostly stamina. The
  bundled total below uses the colour-matched version, so it is the conservative number.
- Your bracer gem is a **Steady Talasite** (stamina + resilience) — a PvP gem doing nothing for
  raid DPS.

## Every candidate, by slot

Each row replaces only that one item on your current gear. Replacement items are gemmed with the
best P3 epic gem matching their own sockets and inherit your current enchant for that slot.

### Head — currently Cataclysm Helm (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Cursed Vision of Sargeras | 151 | 3 | 2272 | +34 | Illidan Stormrage — Black Temple |
| Forest Prowler's Helm | 141 | 3 | 2246 | +9 | The Illidari Council — Black Temple |
| Skyshatter Cover (T6) | 146 | 3 | 2245 | +8 | Tier 6 token — Hyjal / Black Temple |

### Neck — currently Worgen Claw Necklace (ilvl 115)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Choker of Endless Nightmares | 141 | 3 | 2257 | +19 | Supremus — Black Temple |
| Telonicus's Pendant of Mayhem | 138 | 2 | 2241 | +3 | Void Reaver — Tempest Keep |

### Shoulder — currently Cataclysm Shoulderplates (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Swiftstrike Shoulders | 141 | 3 | 2269 | +31 | Crafted (Leatherworking) **[requires Leatherworking]** |
| Shoulders of Lightning Reflexes | 141 | 3 | 2266 | +29 | Crafted (Leatherworking) **[requires Leatherworking]** |
| Mantle of Darkness | 141 | 3 | 2259 | +21 | High Warlord Naj'entus — Black Temple |
| Skyshatter Pauldrons (T6) | 146 | 3 | 2257 | +20 | Tier 6 token — Hyjal / Black Temple |
| Shoulders of the Hidden Predator | 141 | 3 | 2256 | +19 | Shade of Akama — Black Temple |

### Back — currently Razor-Scale Battlecloak (ilvl 128)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Shadowmoon Destroyer's Drape | 141 | 3 | 2249 | +12 | Teron Gorefiend — Black Temple |

### Chest — currently Cataclysm Chestplate (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Midnight Chestguard | 151 | 3 | 2265 | +27 | Archimonde — Hyjal Summit |
| Nether Shadow Tunic | 141 | 3 | 2253 | +16 | Supremus — Black Temple |
| Mail of Fevered Pursuit | 151 | 3 | 2242 | +4 | Archimonde — Hyjal Summit |
| Skyshatter Tunic (T6) | 146 | 3 | 2240 | +3 | Tier 6 token — Hyjal / Black Temple |

### Wrist — currently Veteran's Linked Bracers (ilvl 126)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Bindings of Lightning Reflexes | 141 | 3 | 2269 | +32 | Crafted (Leatherworking) |
| Insidious Bands | 141 | 3 | 2266 | +28 | Teron Gorefiend — Black Temple |
| Swiftstrike Bracers | 141 | 3 | 2266 | +28 | Crafted (Leatherworking) |
| Deadly Cuffs | 141 | 3 | 2265 | +28 | Rage Winterchill — Hyjal Summit |

### Hands — currently Cataclysm Gauntlets (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Grips of Damnation | 141 | 3 | 2257 | +19 | Reliquary of the Lost — Black Temple |
| Fists of Mukoa | 141 | 3 | 2256 | +18 | High Warlord Naj'entus — Black Temple |
| Skyshatter Grips (T6) | 146 | 3 | 2242 | +4 | Tier 6 token — Hyjal / Black Temple |

### Waist — currently Girdle of the Tidal Call (ilvl 128)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Belt of One-Hundred Deaths | 138 | 2 | 2276 | +38 | Lady Vashj — Serpentshrine Cavern |
| Shadow-walker's Cord | 141 | 3 | 2267 | +29 | Shade of Akama — Black Temple |
| Boneweave Girdle | 141 | 3 | 2267 | +29 | Reliquary of the Lost — Black Temple |
| Valestalker Girdle | 141 | 3 | 2267 | +29 | Kaz'rogal — Hyjal Summit |

### Legs — currently Cataclysm Legplates (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Bow-stitched Leggings | 141 | 3 | 2262 | +24 | Azgalor — Hyjal Summit |
| Shady Dealer's Pantaloons | 141 | 3 | 2251 | +13 | Azgalor — Hyjal Summit |
| Leggings of Murderous Intent | 138 | 2 | 2247 | +9 | Kael'thas Sunstrider — Tempest Keep |
| Skyshatter Pants (T6) | 146 | 3 | 2245 | +7 | Tier 6 token — Hyjal / Black Temple |

### Feet — currently Boots of Utter Darkness (ilvl 128)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Shadowmaster's Boots | 141 | 3 | 2254 | +16 | Mother Shahraz — Black Temple |
| Softstep Boots of Tracking | 141 | 3 | 2251 | +13 | Teron Gorefiend — Black Temple |
| Black Featherlight Boots | 141 | 3 | 2236 | -2 | Kaz'rogal — Hyjal Summit |

### Finger1 — currently Ancestral Ring of Conquest (ilvl 128)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Stormrage Signet Ring | 151 | 3 | 2257 | +20 | Illidan Stormrage — Black Temple |
| Band of the Ranger-General | 138 | 2 | 2249 | +11 | Kael'thas Sunstrider — Tempest Keep |
| Ring of Deceitful Intent | 141 | 3 | 2245 | +7 | Shade of Akama — Black Temple |

### Finger2 — currently Violet Signet (ilvl 120)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Stormrage Signet Ring | 151 | 3 | 2262 | +24 | Illidan Stormrage — Black Temple |
| Ring of Deceitful Intent | 141 | 3 | 2249 | +12 | Shade of Akama — Black Temple |
| Violet Signet of the Master Assassin | 130 | 1 | 2240 | +2 | Karazhan — Violet Eye reputation upgrade of your Violet Signet |

### Trinket1 — currently Terokkar Tablet of Precision (ilvl 99)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Dragonspine Trophy | 125 | 1 | 2315 | +78 | Gruul the Dragonkiller — Gruul's Lair |
| Madness of the Betrayer | 141 | 3 | 2277 | +39 | The Illidari Council — Black Temple |
| Tsunami Talisman | 128 | 2 | 2268 | +31 | Leotheras the Blind — Serpentshrine Cavern |
| Bloodlust Brooch | 110 | 1 | 2268 | +30 | Badge of Justice vendor |
| Abacus of Violent Odds | 115 | 1 | 2266 | +29 | Pathaleon the Calculator — The Mechanar |

### Trinket2 — currently Bladefist's Breadth (ilvl 88)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Dragonspine Trophy | 125 | 1 | 2309 | +72 | Gruul the Dragonkiller — Gruul's Lair |
| Madness of the Betrayer | 141 | 3 | 2271 | +33 | The Illidari Council — Black Temple |
| Tsunami Talisman | 128 | 2 | 2262 | +24 | Leotheras the Blind — Serpentshrine Cavern |
| Bloodlust Brooch | 110 | 1 | 2261 | +23 | Badge of Justice vendor |

### MainHand — currently Merciless Gladiator's Cleaver (ilvl 136)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Vengeful Gladiator's Cleaver | 146 | 3 | 2263 | +25 | Season 3 arena vendor |
| Rising Tide | 141 | 3 | 2246 | +8 | High Warlord Naj'entus — Black Temple |
| Syphon of the Nathrezim | 141 | 3 | 2204 | -33 | Supremus — Black Temple |
| Swiftsteel Bludgeon | 141 | 3 | 2098 | -139 | Black Temple |
| The Brutalizer | 141 | 3 | 2089 | -149 | Supremus — Black Temple |

### OffHand — currently Merciless Gladiator's Cleaver (ilvl 136)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Vengeful Gladiator's Cleaver | 146 | 3 | 2253 | +15 | Season 3 arena vendor |
| Rising Tide | 141 | 3 | 2241 | +4 | High Warlord Naj'entus — Black Temple |
| Syphon of the Nathrezim | 141 | 3 | 2206 | -32 | Supremus — Black Temple |
| The Brutalizer | 141 | 3 | 2112 | -126 | Supremus — Black Temple |

### Ranged — currently Totem of the Thunderhead (ilvl 94)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Totem of the Astral Winds | 115 | 1 | 2248 | +10 | Pandemonius — Mana-Tombs |
| Stonebreaker's Totem | 128 | 4 | 2238 | +0 | Badge of Justice vendor |

## Stacking them: the 4-piece T5 cliff

Applying the winners one at a time, best-first, re-simming after every step. The `step` column is
the marginal gain *on top of everything above it*.

| Step | Slot | Item | DPS | Step gain | Total gain |
|------|------|------|-----|-----------|------------|
| 1 | Trinket1 | Dragonspine Trophy | 2315 | +78 | +78 |
| 2 | Waist | Belt of One-Hundred Deaths | 2358 | +43 | +121 |
| 3 | Head | Cursed Vision of Sargeras | 2395 | +36 | +157 |
| 4 | Trinket2 | Madness of the Betrayer | 2434 | +39 | +196 |
| 5 | Wrist | Bindings of Lightning Reflexes | 2462 | +28 | +224 |
| 6 | Shoulder | Swiftstrike Shoulders | 2437 | -25 ⚠️ | +199 |
| 7 | Chest | Midnight Chestguard | 2467 | +31 | +230 |
| 8 | MainHand | Vengeful Gladiator's Cleaver | 2495 | +28 | +258 |
| 9 | Legs | Bow-stitched Leggings | 2521 | +26 | +283 |
| 10 | Finger2 | Stormrage Signet Ring | 2548 | +28 | +311 |
| 11 | Neck | Choker of Endless Nightmares | 2571 | +23 | +333 |
| 12 | Hands | Grips of Damnation | 2583 | +12 | +346 |
| 13 | Feet | Shadowmaster's Boots | 2600 | +17 | +363 |
| 14 | OffHand | Vengeful Gladiator's Cleaver | 2618 | +18 | +380 |
| 15 | Back | Shadowmoon Destroyer's Drape | 2633 | +15 | +396 |
| 16 | Finger1 | Band of the Ranger-General | 2647 | +14 | +410 |
| 17 | Ranged | Totem of the Astral Winds | 2660 | +12 | +422 |

Step 6 is the cliff. Swiftstrike Shoulders is worth **+31 DPS** on its own, but arriving
after Cursed Vision of Sargeras it is your *second* tier piece leaving, which drops Cataclysm
Harness from 4 pieces to 3 and kills the 4pc bonus (+5% additional Flurry haste). The step comes
out at **-25** — a swing of about 56 DPS. Everything after it is measured with the bonus
already gone, which is why the later tier swaps (chest, legs) look healthy again.

**Practical reading:** replace *one* tier piece and stop, or replace *four* in quick succession.
The middle ground is where DPS goes to die.

## Assembled builds

| Build | DPS | vs current | AP | Hit | Crit | Haste | Exp |
|-------|-----|-----------|----|-----|------|-------|-----|
| A. best item every slot | 2660 | +422 | 3363 | 188 | 327 | 102 | 45 |
| A+. best item every slot + enchant/gem fixes | 2679 | +442 | 3378 | 188 | 339 | 102 | 45 |
| B. best items, keep 4pc Cataclysm Harness (T5) | 2615 | +377 | 3299 | 223 | 284 | 27 | 45 |
| B+. B + enchant/gem fixes | 2644 | +407 | 3334 | 223 | 298 | 27 | 45 |
| C. best items with 4pc Skyshatter (T6) | 2606 | +368 | 3358 | 205 | 349 | 64 | 45 |
| C+. C + enchant/gem fixes | 2626 | +388 | 3374 | 205 | 361 | 64 | 45 |

Keeping 4pc Cataclysm Harness (build B+, 2644) costs about
35 DPS against going full non-tier (A+), so the T5 4pc is worth holding onto right up until
you can replace four slots at once — but it is not worth *building around* long term.

**4pc Skyshatter (T6) is not the answer** (2626, -54 vs A+). The 4pc bonus (+70 AP for 12s after
Stormstrike) is implemented in the sim and it still loses: the Skyshatter pieces themselves are
weaker than the best non-tier options in every one of those slots (see the per-slot tables — every
Skyshatter piece places last or near-last).

## Reference points

| Reference gear set | DPS | vs current |
|--------------------|-----|-----------|
| repo P2 preset (as checked in) | 2287 | +50 |
| repo P2 preset (enchant ids fixed) | 2465 | +227 |
| repo P3 BiS preset (as checked in) | 2435 | +197 |
| repo P3 BiS preset (enchant ids fixed) | 2623 | +385 |
| repo P4 BiS preset (as checked in) | 2448 | +210 |
| repo P4 BiS preset (enchant ids fixed) | 2634 | +397 |

Note that the repo's P3 BiS preset (2623) lands **57 DPS below** the build this report
recommends. That preset dual-wields Syphon of the Nathrezim; for an Orc, keeping axes is worth more
than the maces are.

The "as checked in" vs "enchant ids fixed" split is explained under [Bug found](#bug-found-gear-preset-enchants-silently-do-nothing).

## Talents

| Build | DPS | vs baseline |
|-------|-----|-------------|
| Armory spec 2 (0/47/14, PvE) [baseline] | 2238 | +0 |
| Armory spec 1 (17/44/0) | 2291 | +53 |
| Repo default enhancement (3/44/14) | 2245 | +8 |

Your armory second spec (0/47/14) was used as the baseline for every run in this report because it
is the conventional PvE enhancement build. But your **first spec — 17/44/0 — sims
+53 DPS higher** on this gear and this rotation. The 17 Elemental points are Convection 2,
Concussion 5, Call of Flame 3, Elemental Focus, Reverberation 5 and Improved Fire Totems 1, which
directly buff the shocks and the Magma/Fire Nova totem twisting that this repo's default
Enhancement APL leans on, while the Restoration points you would give up (Totemic Focus,
Nature's Guidance) buy mana and hit you largely do not need — see the hit note below.

Treat that as a sim result on this specific rotation, not a blanket recommendation: the shock-heavy
build is more mana-hungry, and its value depends on the fire-totem twisting actually happening.

## A note on hit rating

Your gear carries **207 melee hit rating ≈ 13.1%**. Special attacks cap at 9% against a raid boss
(142 rating at level 70), and your PvE spec adds another 3% from Nature's Guidance on top. So every
point of hit past that does nothing for Stormstrike or your shocks — it only helps white swings,
which sit far from their own (much higher) dual-wield cap.

This is why several high-hit Black Temple pieces score lower here than their item level suggests,
and why strength/AP/crit pieces win your slots. It is also part of why the 17/44/0 spec pulls ahead.

## The build this report ends at

Build A+ — 2679 DPS (+442).

| Slot | Item | ilvl | Enchant | Gems |
|------|------|------|---------|------|
| Head | Cursed Vision of Sargeras | 151 | Glyph of Ferocity | Relentless Earthstorm Diamond, Inscribed Pyrestone |
| Neck | Choker of Endless Nightmares | 141 | — | — |
| Shoulder | Swiftstrike Shoulders | 141 | Greater Inscription of Vengeance | — |
| Back | Shadowmoon Destroyer's Drape | 141 | Greater Agility | — |
| Chest | Midnight Chestguard | 151 | Exceptional Stats | Sovereign Shadowsong Amethyst, Bold Crimson Spinel, Inscribed Pyrestone |
| Wrist | Bindings of Lightning Reflexes | 141 | Brawn | — |
| Hands | Grips of Damnation | 141 | Major Strength | — |
| Waist | Belt of One-Hundred Deaths | 138 | — | Bold Crimson Spinel, Sovereign Shadowsong Amethyst |
| Legs | Bow-stitched Leggings | 141 | Nethercobra Leg Armor | Bold Crimson Spinel, Inscribed Pyrestone, Sovereign Shadowsong Amethyst |
| Feet | Shadowmaster's Boots | 141 | Cat's Swiftness | Bold Crimson Spinel, Inscribed Pyrestone |
| Finger1 | Band of the Ranger-General | 138 | — | — |
| Finger2 | Stormrage Signet Ring | 151 | — | — |
| Trinket1 | Dragonspine Trophy | 125 | — | — |
| Trinket2 | Madness of the Betrayer | 141 | — | — |
| MainHand | Vengeful Gladiator's Cleaver | 146 | Mongoose | — |
| OffHand | Vengeful Gladiator's Cleaver | 146 | Mongoose | — |
| Ranged | Totem of the Astral Winds | 115 | — | — |

## How this was measured

- **Harness**: `sim/shamanupgrade_test.go`, modelled on the Phase 3 trinket benchmark in
  [results.md](results.md). Re-run it with `go test -tags with_db ./sim/ -run TestShamanP3Upgrades -v`;
  set `SHAMAN_ITERATIONS` (default 20000) to trade accuracy for speed.
- **Iterations**: 20,000 per configuration, `RandomSeed: 101`, `core.RunRaidSimConcurrent`.
  ~85 configurations, ~24 minutes total.
- **Encounter**: single target, `core.MakeSingleTargetEncounter(5)` — 180s ±5s. Single target only;
  no cleave scenarios were run.
- **Buffs/debuffs**: `FullRaidBuffs` / `FullPartyBuffs` / `FullIndividualBuffs` / `FullDebuffs`.
  No bonus hit rating was added (unlike the trinket benchmark) — your real gear hit is used, so
  hit-rating upgrades are valued correctly.
- **Character**: Orc, Enhancement, Windfury on both weapons, `SyncType: Auto`, default
  `ui/shaman/enhancement/apls/default.apl.json` rotation, talents `-502520210502133531151-50005301`
  (your armory second spec). Consumables match the Enhancement config from the trinket benchmark:
  Haste Potion, Flask of Relentless Assault, food, Dark Rune, sappers/grenades, agi/str scrolls.
- **Gear source**: your equipped items, enchants and gems as reported by ironforge.pro, transcribed
  verbatim into `currentGear()`. Two details worth flagging:
  - Your boots' **Minor Speed Increase** has no entry in this repo's enchant database and grants no
    stats, so the feet slot is simmed as unenchanted — which is also what it is worth.
  - Your shirt (Stylish Red Shirt) is not a simmed slot.
- **Replacement gemming**: every candidate item is gemmed by reading its *own* sockets from
  `assets/database/db.json` and filling each with the best P3 epic gem for that colour (meta →
  Relentless Earthstorm Diamond, red → Bold Crimson Spinel, yellow → Inscribed Pyrestone, blue →
  Sovereign Shadowsong Amethyst). This matters: `NewItem` in `sim/core/database.go` applies every
  gem it is handed, even past the item's socket count, so a hand-written gem list on a socketless
  item would have silently granted free stats.
- **Unique-equipped**: the engine does not enforce it. The cumulative and assembled builds refuse
  to hand the same item id to two slots (weapons excepted, where two copies are legal), so no
  build here equips two Dragonspine Trophies.
- **Candidate pool**: for each slot, the Phase 3 items in this repo's database that a shaman can
  wear (cloth/leather/mail; axe/mace/dagger/fist one-handers; totems), plus the notable Phase 1-2
  items your current piece is losing to, plus the tier 6 (Skyshatter Harness) piece for that slot.
  Kael'thas' Infinity Blade (ilvl 175) was excluded as effectively legendary-tier, matching the
  constraint used in [results.md](results.md) for Warp Slicer.

## Caveats

- **Single target, patchwerk-style.** No movement, no adds, no target switching. Dragonspine
  Trophy and Madness of the Betrayer are both proc trinkets whose value would shift on a fight
  with downtime.
- **The rotation is the sim's, not yours.** The default Enhancement APL twists Windfury Totem,
  twists Flame Shock, and keeps Magma Totem up. If you do not do those things, items that feed them
  (and the 17/44/0 talent result) are worth less to you.
- **Meta gem colour requirements are not modelled.** The engine has no requirement field for meta
  gems, so the Relentless Earthstorm Diamond always applies. Check your own blue-gem count in game.
- **The T5 2pc bonus is not implemented** (it is a Lesser Healing Wave proc, deliberately skipped in
  `sim/shaman/item_sets.go`). Irrelevant to DPS, but it means "2pc T5" has no value in these numbers.
- **Profession requirements.** Swiftstrike Shoulders and Shoulders of Lightning Reflexes are
  Leatherworking-only to wear. If you are not a leatherworker, your best shoulder is
  **Mantle of Darkness** (+21), from Naj'entus in Black Temple. The bracers
  (Bindings of Lightning Reflexes, Swiftstrike Bracers) carry no such restriction in this database.
- **Phase tagging is this repo's.** Season 3 arena gear and The Brutalizer are tagged `phase: 3` in
  `assets/database/db.json`, which is why they appear as Phase 3 candidates here.
- **20,000 iterations** leaves roughly ±1 DPS of noise on each number. Differences under ~5 DPS
  between two candidates are not meaningful.

## Bug found: gear-preset enchants silently do nothing

`NewItem` resolves `ItemSpec.Enchant` through `core.EnchantsByEffectID` and silently ignores
anything it cannot find — the `panic` for an unknown enchant is commented out
(`sim/core/database.go:437-443`). The enchant database is keyed by **effect id**, but the
checked-in Enhancement Shaman gear presets (`ui/shaman/enhancement/gear_sets/p2|p3|p4.gear.json`)
store the enchant's **item id** instead:

| In the preset | Enchant | Real effect id |
|---------------|---------|----------------|
| 29192 | Glyph of Ferocity | 3003 |
| 28888 | Greater Inscription of Vengeance | 2986 |
| 24003 | Enchant Chest - Exceptional Stats | 2661 |
| 29535 | Nethercobra Leg Armor | 3012 |
| 28279 | Enchant Boots - Cat's Swiftness | 2939 |
| 22559 | Enchant Weapon - Mongoose | 2673 |
| 33307 | Enchant Weapon - Executioner | 3225 |

So those presets sim with **no head, shoulder, chest, leg or boot enchant and no Mongoose at all**.
The cost is large: the P3 preset goes from 2435 DPS as checked in to 2623 once the ids are
translated — a ~188 DPS understatement.

Auditing all 109 gear presets in `ui/*/*/gear_sets/`, exactly **11** contain enchant ids the engine
cannot resolve, and they are confined to two specs:

| Preset files | Unresolvable enchant ids |
|--------------|--------------------------|
| `ui/shaman/enhancement/gear_sets/` — preraid, p1-p5 | 22559, 24003, 28279, 28888, 29192, 29535, 33307 |
| `ui/druid/balance/gear_sets/` — p1_a, p2_a, p3, p4, p5 | 22534, 22536, 22560, 24003, 24274, 28272, 28886, 29191, 33150, 35297 |

Every other preset — including the elemental and restoration shaman ones — uses effect ids
correctly. Nothing was changed in the repo for this report: `sim/shamanupgrade_test.go` uses effect
ids throughout and translates preset ids via `fixPresetEnchants` for the reference rows, which is
why both an "as checked in" and an "enchant ids fixed" number are given above.

Knock-on effect on earlier work in this repo: of the 14 configurations in [results.md](results.md),
only the **Enhancement Shaman** row is affected (it uses `p3.gear.json` from the list above), and it
is understated by roughly the same ~190 DPS. The warrior, rogue, hunter, feral druid and
retribution paladin rows there are unaffected, as is all of [results_glaives.md](results_glaives.md)
(Fury Warrior and Rogue only). The cross-spec ranking in results.md therefore under-rates
Enhancement Shaman against the other 13 configurations.
