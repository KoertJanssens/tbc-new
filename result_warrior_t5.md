# Phase 3 upgrade report — Koerry (Fury Warrior, Spineshatter-EU) — from 5/5 T5

Sim-backed answer to "what are my biggest Phase 3 upgrades?", built from the live armory
profile at <https://ironforge.pro/anniversary/player/Spineshatter/Koerry/> and run through
this repo's engine at 20,000 iterations per configuration.

**Current gear sims at 2220 DPS** (stdev 103). Everything below is a delta against that number.

Departures from the armory:

- **Trinket 2 is Dragonspine Trophy**, not the Solarian's Sapphire the armory shows, as you told me.
- **The head slot is Destroyer Battle-Helm**, not your Furious Gizmatic Goggles, giving a full
  **5/5 Destroyer Battlegear** baseline as you asked. This is the whole point of this run: with
  five tier pieces the 4pc survives one swap, which changes the upgrade order. The 4/5 version
  is in [result_warrior.md](result_warrior.md).
- **Everything is simmed as Fury**, per your request — the repo's standard Fury build
  (`3500501130201-05050005505012050115`) with `ui/warrior/dps/apls/fury.apl.json`, on every single
  row. Your armory has two saved Arms specs; neither appears anywhere in this report.

## TL;DR

1. **Weapons first — they are your two biggest single upgrades.** Vengeful Gladiator's Cleaver in
   the main hand is **+48**, and it beats a *Warglaive of Azzinoth* in that slot (+44), because you
   are an Orc: Axe Specialization is +5 expertise and the sim models it (expertise goes 20 → 40).
   A fast off hand is the other half — Blade of Savagery (1.4 speed) is **+41**.
2. **Keep your Dragonspine Trophy.** Every trinket I tested against it lost, by 30-40 DPS. The
   weak trinket is Badge of the Swarmguard: Madness of the Betrayer replaces it for **+37**.
3. **Your first tier swap is free; the second one is the decision.** With 5/5 Destroyer the 4pc
   (-5 rage on Bloodthirst) survives one piece leaving. After that every further tier swap pays
   for the bonus — and it is still worth paying: playing it safe and defending the 4pc caps you at
   **2538**, against **2631** if you let it go (+93).
4. **4pc Onslaught (T6) is not worth chasing either** (2598), landing -7 against simply wearing the best
   available item in every slot. Take Onslaught pieces when they are individually better, not for the bonus.
5. Realistic ceiling: **2631 DPS (+411, +18%)** without legendaries. With *both* Warglaives
   of Azzinoth, **2722 (+502)** — see the note on why the pair is worth far more than either half.

For reference, going back to the goggles from this baseline is **-2 DPS** — the two helms are
a wash on raw stats. The entire value of 5/5 is the free tier swap, not the helm itself.

## Priority order — what to chase, in order

Ranked by single-slot DPS gain, with unique-equipped conflicts resolved and legendaries excluded.

| # | Slot | Upgrade | ilvl | DPS gain | Where it comes from |
|---|------|---------|------|----------|---------------------|
| 1 | MainHand | Vengeful Gladiator's Cleaver | 146 | **+48** | Season 3 arena vendor |
| 2 | OffHand | Blade of Savagery | 141 | **+41** | Mother Shahraz — Black Temple |
| 3 | Legs | Leggings of Divine Retribution | 141 | **+37** | Gurtogg Bloodboil — Black Temple |
| 4 | Trinket1 | Madness of the Betrayer | 141 | **+37** | The Illidari Council — Black Temple |
| 5 | Head | Cursed Vision of Sargeras | 151 | **+33** | Illidan Stormrage — Black Temple |
| 6 | Feet | Dreadboots of the Legion | 141 | **+31** | Reliquary of the Lost — Black Temple |
| 7 | Finger2 | Stormrage Signet Ring | 151 | **+23** | Illidan Stormrage — Black Temple |
| 8 | Hands | Grips of Silent Justice | 141 | **+22** | Shade of Akama — Black Temple |
| 9 | Chest | Midnight Chestguard | 151 | **+22** | Archimonde — Hyjal Summit |
| 10 | Shoulder | Onslaught Shoulderblades (T6) | 146 | **+20** | Tier 6 token — Hyjal / Black Temple |
| 11 | Finger1 | Unstoppable Aggressor's Ring | 141 | **+18** | Gurtogg Bloodboil — Black Temple |
| 12 | Waist | Belt of One-Hundred Deaths | 138 | **+17** | Lady Vashj — Serpentshrine Cavern |
| 13 | Back | Shadowmoon Destroyer's Drape | 141 | **+16** | Teron Gorefiend — Black Temple |
| 14 | Neck | Choker of Endless Nightmares | 141 | **+15** | Supremus — Black Temple |
| 15 | Wrist | Deadly Cuffs | 141 | **+13** | Rage Winterchill — Hyjal Summit |
| 16 | Ranged | Serpent Spine Longbow | 141 | **+4** | Lady Vashj — Serpentshrine Cavern |

These are *independent* single-slot measurements: each replaces exactly that one item. They do not
add up — see [Stacking them](#stacking-them).

## Free upgrades — no drops required

| Change | DPS | Gain | Detail |
|--------|-----|------|--------|
| Boots: Cat's Swiftness -> Dexterity | 2223 | **+4** | +6 agi for the run speed |
| Weapons: Mongoose -> Executioner (both) | 2213 | **-7** | 840 armor pen proc instead of 120 agi / 2% haste |
| Weapons: Executioner MH only | 2231 | **+11** | Executioner main hand, Mongoose off hand |
| Regem everything with P3 epic gems | 2238 | **+18** | colour-matched for socket bonuses |
| Regem everything, all red (ignore socket bonuses) | 2244 | **+24** | every non-meta socket = Bold Crimson Spinel (+10 str) |
| All gem fixes together | 2238 | **+18** | colour-matched regem |

Notes:

- **Gem all red.** Bold Crimson Spinel in every non-meta socket (+24) beats colour-matching for
  socket bonuses (+18). Your socket bonuses are small and stamina-heavy; +10 strength each is not.
- **Executioner belongs on the main hand only** (+11). On *both* weapons it is a net loss
  (-7) — Mongoose's agility and haste are worth more in the off hand than a second armor-pen proc.
- Your enchants are otherwise already correct: Glyph of Ferocity, Greater Inscription of Vengeance,
  Exceptional Stats, Nethercobra, Brawn, Major Strength, Greater Agility. Nothing to fix there.

## Every candidate, by slot

Each row replaces only that one item on your current gear. Replacements are gemmed with the best P3
epic gem matching their own sockets and inherit your current enchant for that slot.

### Head — currently Destroyer Battle-Helm (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Cursed Vision of Sargeras | 151 | 3 | 2253 | +33 | Illidan Stormrage — Black Temple |
| Helm of the Illidari Shatterer | 141 | 3 | 2238 | +18 | The Illidari Council — Black Temple |
| Onslaught Battle-Helm (T6) | 146 | 3 | 2232 | +12 | Tier 6 token — Hyjal / Black Temple |
| Furious Gizmatic Goggles | 127 | 2 | 2218 | -2 | Crafted (Engineering) **[requires Engineering]** |
| Surestrike Goggles v2.0 | 127 | 2 | 2212 | -8 | Crafted (Engineering) **[requires Engineering]** |
| Fel-Steel Warhelm | 128 | 2 | 2164 | -56 | Crafted (Blacksmithing) |

### Neck — currently Pendant of the Perilous (ilvl 128)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Choker of Endless Nightmares | 141 | 3 | 2235 | +15 | Supremus — Black Temple |
| Choker of Serrated Blades | 141 | 3 | 2225 | +5 | Hyjal Summit; Black Temple |

### Shoulder — currently Destroyer Shoulderblades (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Onslaught Shoulderblades (T6) | 146 | 3 | 2240 | +20 | Tier 6 token — Hyjal / Black Temple |
| Blood-stained Pauldrons | 141 | 3 | 2231 | +11 | Rage Winterchill — Hyjal Summit |
| Swiftsteel Shoulders | 141 | 3 | 2228 | +8 | Crafted (Blacksmithing) **[requires Blacksmithing]** |
| Mantle of Darkness | 141 | 3 | 2227 | +7 | High Warlord Naj'entus — Black Temple |

### Back — currently Cloak of the Pit Stalker (ilvl 125)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Shadowmoon Destroyer's Drape | 141 | 3 | 2236 | +16 | Teron Gorefiend — Black Temple |
| Black-Iron Battlecloak | 120 | 1 | 2227 | +7 | — |

### Chest — currently Destroyer Breastplate (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Midnight Chestguard | 151 | 3 | 2242 | +22 | Archimonde — Hyjal Summit |
| Onslaught Breastplate (T6) | 146 | 3 | 2235 | +15 | Tier 6 token — Hyjal / Black Temple |
| Heartshatter Breastplate | 141 | 3 | 2228 | +8 | Mother Shahraz — Black Temple |
| Krakken-Heart Breastplate | 138 | 2 | 2212 | -8 | Lady Vashj — Serpentshrine Cavern |

### Wrist — currently Bladespire Warbands (ilvl 125)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Deadly Cuffs | 141 | 3 | 2233 | +13 | Rage Winterchill — Hyjal Summit |
| Bindings of Lightning Reflexes | 141 | 3 | 2228 | +8 | Crafted (Leatherworking) |
| Bracers of Eradication | 128 | 2 | 2226 | +6 | The Lurker Below — Serpentshrine Cavern |
| Swiftsteel Bracers | 141 | 3 | 2222 | +2 | Crafted (Blacksmithing) |

### Hands — currently Destroyer Gauntlets (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Grips of Silent Justice | 141 | 3 | 2242 | +22 | Shade of Akama — Black Temple |
| Grips of Damnation | 141 | 3 | 2235 | +15 | Reliquary of the Lost — Black Temple |
| Onslaught Gauntlets (T6) | 146 | 3 | 2234 | +14 | Tier 6 token — Hyjal / Black Temple |
| Pillager's Gauntlets | 141 | 3 | 2229 | +9 | Black Temple |

### Waist — currently Red Belt of Battle (ilvl 128)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Belt of One-Hundred Deaths | 138 | 2 | 2237 | +17 | Lady Vashj — Serpentshrine Cavern |
| Boneweave Girdle | 141 | 3 | 2227 | +7 | Reliquary of the Lost — Black Temple |
| Shadow-walker's Cord | 141 | 3 | 2221 | +1 | Shade of Akama — Black Temple |
| Belt of Seething Fury | 141 | 3 | 2215 | -5 | Kaz'rogal — Hyjal Summit |

### Legs — currently Destroyer Greaves (ilvl 133)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Leggings of Divine Retribution | 141 | 3 | 2257 | +37 | Gurtogg Bloodboil — Black Temple |
| Legguards of Endless Rage | 151 | 3 | 2241 | +21 | Archimonde — Hyjal Summit |
| Onslaught Greaves (T6) | 146 | 3 | 2238 | +18 | Tier 6 token — Hyjal / Black Temple |
| Shady Dealer's Pantaloons | 141 | 3 | 2225 | +5 | Azgalor — Hyjal Summit |

### Feet — currently Ironstriders of Urgency (ilvl 115)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Dreadboots of the Legion | 141 | 3 | 2251 | +31 | Reliquary of the Lost — Black Temple |
| Softstep Boots of Tracking | 141 | 3 | 2239 | +19 | Teron Gorefiend — Black Temple |
| Warboots of Obliteration | 128 | 2 | 2225 | +5 | Morogrim Tidewalker — Serpentshrine Cavern |
| Black Featherlight Boots | 141 | 3 | 2213 | -7 | Kaz'rogal — Hyjal Summit |

### Finger1 — currently Mithril Band of the Unscarred (ilvl 115)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Stormrage Signet Ring | 151 | 3 | 2239 | +19 | Illidan Stormrage — Black Temple |
| Unstoppable Aggressor's Ring | 141 | 3 | 2238 | +18 | Gurtogg Bloodboil — Black Temple |
| Band of the Ranger-General | 138 | 2 | 2235 | +15 | Kael'thas Sunstrider — Tempest Keep |
| Ring of Deceitful Intent | 141 | 3 | 2224 | +4 | Shade of Akama — Black Temple |

### Finger2 — currently Shapeshifter's Signet (ilvl 100)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Stormrage Signet Ring | 151 | 3 | 2243 | +23 | Illidan Stormrage — Black Temple |
| Unstoppable Aggressor's Ring | 141 | 3 | 2241 | +21 | Gurtogg Bloodboil — Black Temple |
| Band of the Ranger-General | 138 | 2 | 2238 | +18 | Kael'thas Sunstrider — Tempest Keep |
| Ring of Deceitful Intent | 141 | 3 | 2227 | +7 | Shade of Akama — Black Temple |

### Trinket1 — currently Badge of the Swarmguard (ilvl 76)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Madness of the Betrayer | 141 | 3 | 2257 | +37 | The Illidari Council — Black Temple |
| Tsunami Talisman | 128 | 2 | 2253 | +33 | Leotheras the Blind — Serpentshrine Cavern |
| Bloodlust Brooch | 110 | 1 | 2246 | +26 | Badge of Justice vendor |
| Abacus of Violent Odds | 115 | 1 | 2242 | +22 | Heroic Mechanar — Pathaleon |
| Drake Fang Talisman | 75 | 1 | 2218 | -2 | Ahn'Qiraj (AQ40) |

### Trinket2 — currently Dragonspine Trophy (ilvl 125)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Madness of the Betrayer | 141 | 3 | 2181 | -39 | The Illidari Council — Black Temple |
| Tsunami Talisman | 128 | 2 | 2179 | -41 | Leotheras the Blind — Serpentshrine Cavern |
| Bloodlust Brooch | 110 | 1 | 2172 | -48 | Badge of Justice vendor |
| Abacus of Violent Odds | 115 | 1 | 2171 | -49 | Heroic Mechanar — Pathaleon |

### MainHand — currently Talon of Azshara (ilvl 134)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Vengeful Gladiator's Cleaver | 146 | 3 | 2268 | +48 | Season 3 arena vendor |
| Warglaive of Azzinoth (legendary) | 156 | 3 | 2264 | +44 | Illidan Stormrage — Black Temple (legendary) |
| Wicked Edge of the Planes | 136 | 2 | 2249 | +29 | Crafted (Blacksmithing) **[requires Blacksmithing]** |
| Vengeful Gladiator's Slicer | 146 | 3 | 2245 | +25 | Season 3 arena vendor |
| Rising Tide | 141 | 3 | 2242 | +22 | High Warlord Naj'entus — Black Temple |
| Blade of Infamy | 141 | 3 | 2235 | +15 | Anetheron — Hyjal Summit |
| Syphon of the Nathrezim | 141 | 3 | 2226 | +6 | Supremus — Black Temple |

### OffHand — currently Merciless Gladiator's Slicer (ilvl 136)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Warglaive of Azzinoth (legendary) | 156 | 3 | 2282 | +62 | Illidan Stormrage — Black Temple (legendary) |
| Blade of Savagery | 141 | 3 | 2261 | +41 | Mother Shahraz — Black Temple |
| Vengeful Gladiator's Cleaver | 146 | 3 | 2244 | +24 | Season 3 arena vendor |
| Tracker's Blade | 141 | 3 | 2243 | +23 | Rage Winterchill — Hyjal Summit |
| Vengeful Gladiator's Slicer | 146 | 3 | 2238 | +18 | Season 3 arena vendor |
| Vengeful Gladiator's Shanker | 146 | 3 | 2234 | +14 | Season 3 arena vendor |
| Rising Tide | 141 | 3 | 2222 | +2 | High Warlord Naj'entus — Black Temple |

### Ranged — currently Merciless Gladiator's War Edge (ilvl 136)

| Item | ilvl | Phase | DPS | vs current | Source |
|------|------|-------|-----|-----------|--------|
| Serpent Spine Longbow | 141 | 2 | 2224 | +4 | Lady Vashj — Serpentshrine Cavern |
| Steelhawk Crossbow | 115 | 1 | 2212 | -7 | Attumen the Huntsman — Karazhan |

## Stacking them

Applying the winners one at a time, best-first, re-simming after every step. The `step` column is the
marginal gain *on top of everything above it*.

| Step | Slot | Item | DPS | Step gain | Total gain |
|------|------|------|-----|-----------|------------|
| 1 | MainHand | Vengeful Gladiator's Cleaver | 2268 | +48 | +48 |
| 2 | OffHand | Blade of Savagery | 2310 | +42 | +90 |
| 3 | Legs | Leggings of Divine Retribution | 2352 | +42 | +132 |
| 4 | Trinket1 | Madness of the Betrayer | 2388 | +36 | +168 |
| 5 | Head | Cursed Vision of Sargeras | 2403 | +14 | +183 |
| 6 | Feet | Dreadboots of the Legion | 2438 | +35 | +218 |
| 7 | Finger2 | Stormrage Signet Ring | 2464 | +26 | +244 |
| 8 | Hands | Grips of Silent Justice | 2491 | +27 | +271 |
| 9 | Chest | Midnight Chestguard | 2511 | +20 | +291 |
| 10 | Shoulder | Onslaught Shoulderblades (T6) | 2536 | +25 | +316 |
| 11 | Finger1 | Unstoppable Aggressor's Ring | 2556 | +20 | +336 |
| 12 | Waist | Belt of One-Hundred Deaths | 2574 | +18 | +354 |
| 13 | Back | Shadowmoon Destroyer's Drape | 2593 | +19 | +373 |
| 14 | Neck | Choker of Endless Nightmares | 2611 | +18 | +391 |
| 15 | Wrist | Deadly Cuffs | 2626 | +15 | +406 |
| 16 | Ranged | Serpent Spine Longbow | 2631 | +5 | +411 |

Your tier pieces are the interesting part. You start at **5/5 Destroyer**, so the 4pc
(-5 rage on Bloodthirst) can absorb **1** swap before it dies.

- Tier swap 1: step 3, Leggings of Divine Retribution (+42) — free — still 4pc.
- Tier swap 2: step 5, Cursed Vision of Sargeras (+14) — **this is where the 4pc dies**.

Step 5 is the decision point, and it still pays (+14) even while giving up the bonus.
Everything after it is measured with the set already gone, which is why the remaining tier
swaps — Grips of Silent Justice (step 8, +27), Midnight Chestguard (step 9, +20), Onslaught Shoulderblades (step 10, +25) — look so much healthier than they do as single swaps.

So the ordering collapses to one question: take the free tier upgrade, then decide whether to
commit to losing the set. The build table below answers it.

## Assembled builds — the tier question

| Build | DPS | vs current | AP | Hit | Crit | Haste | Exp |
|-------|-----|-----------|----|-----|------|-------|-----|
| A. stacked single-slot winners | 2631 | +411 | 4091 | 222 | 434 | 0 | 45 |
| B. best items, keeping 4pc Destroyer (T5) | 2538 | +318 | 4035 | 211 | 463 | 0 | 45 |
| C. tierless: best item every slot | 2605 | +385 | 3997 | 222 | 404 | 0 | 65 |
| C+. tierless + Executioner on both weapons | 2607 | +387 | 3997 | 222 | 404 | 0 | 65 |
| D. best items with 4pc Onslaught (T6) | 2598 | +378 | 4023 | 208 | 324 | 0 | 65 |
| D+. 4pc Onslaught + Executioner | 2597 | +377 | 4023 | 208 | 324 | 0 | 65 |
| E. tierless, Warglaives allowed (sim picks) | 2630 | +410 | 3997 | 207 | 405 | 0 | 65 |
| F. Warglaives allowed, keeping 4pc Destroyer (T5) | 2535 | +315 | 3939 | 196 | 434 | 0 | 65 |
| G. tierless, both Warglaives forced | 2722 | +502 | 4009 | 220 | 384 | 0 | 45 |

Reading these:

- **Defending 4pc Destroyer costs you.** Build B spends its one free tier swap and then stops,
  landing at 2538 — -93 against build A. One free swap is worth taking; refusing the
  second, third and fourth is not.
- **4pc Onslaught (T6) does not rescue the idea** (2598, -7 vs tierless). Its 4pc is +5% Mortal
  Strike / Bloodthirst damage, which is real, but the Onslaught pieces themselves are armor-penetration
  and strength itemised and lose crit and hit against the best non-tier options — note the crit column
  dropping to 324 against 404.
- **Executioner adds a little to the tierless build** (+2) but subtracts from the T6 one; it is
  inside the noise either way at this iteration count.

### The Warglaives are a pair or they are nothing

Neither glaive is your best single-slot pick on its own — the main-hand Warglaive is +44 against
+48 for the S3 axe. But equipped **together** they are worth far more than the sum: build G
(both forced) reaches **2722**, while build E — where the optimiser is *allowed* glaives and picks
the axe main hand plus one glaive off hand — only reaches 2630.

That ~92 DPS gap is the two-piece set bonus, "The Twin Blades of Azzinoth"
(`sim/common/tbc/items_weapons.go:338`): a 1 PPM, 45s-ICD proc granting **450 haste rating for 10
seconds**, plus 200 attack power against demons. A greedy slot-by-slot search structurally cannot find
it, which is why build G forces the pair explicitly.

Practical consequence: **do not take a single Warglaive over the S3 axe.** One glaive is worse than an
axe you can buy with arena points. Only the pair is worth reorganising around.

## Talents

| Build | DPS | vs baseline |
|-------|-----|-------------|
| Fury, repo standard [baseline] | 2220 | +0 |
| Fury + Sword Specialization | 2208 | -12 |

You dual-wield two swords, so Sword Specialization looks tempting. It is not: moving 5 points out of
Improved Berserker Stance costs **10% attack power** (the sim shows ~374 AP) and the extra sword procs
recover only part of it — **-12 DPS**. Stay on the standard build.

## A note on hit and expertise

Your gear carries **136 melee hit rating ≈ 8.6%**, which is *below* the 9% special-attack cap
(142 rating at level 70). That makes hit genuinely valuable for you right now — unlike the shaman, who
is well over it — and it is why several hit-carrying pieces rank above higher-item-level alternatives.
Watch for upgrades that trade it away: Leggings of Divine Retribution wins its slot despite dropping
you to 106 hit, but only because 350 armor penetration and 35 crit outweigh it.

Expertise is the other half of the Orc story, and there is a trap in it. Your current 20 rating comes
entirely from **Shapeshifter's Signet** - your weakest-looking ring is quietly the only expertise you
own. Equipping an axe adds 20 more via the racial, and Belt of One-Hundred Deaths another 25.

That means several of the recommended ring upgrades *cost* you expertise: build A replaces
Shapeshifter's Signet and lands on 45 rating rather than 65, and still wins on raw stats. Worth
knowing if you find yourself choosing between rings — the sim has already priced the trade, but it is
not visible from the tooltips.

## The build this report ends at

Build A — 2631 DPS (+411), no legendaries.

| Slot | Item | ilvl | Enchant | Gems |
|------|------|------|---------|------|
| Head | Cursed Vision of Sargeras | 151 | Glyph of Ferocity | Relentless Earthstorm Diamond, Inscribed Pyrestone |
| Neck | Choker of Endless Nightmares | 141 | — | — |
| Shoulder | Onslaught Shoulderblades | 146 | Greater Inscription of Vengeance | Inscribed Pyrestone, Sovereign Shadowsong Amethyst |
| Back | Shadowmoon Destroyer's Drape | 141 | Greater Agility | — |
| Chest | Midnight Chestguard | 151 | Exceptional Stats | Sovereign Shadowsong Amethyst, Bold Crimson Spinel, Inscribed Pyrestone |
| Wrist | Deadly Cuffs | 141 | Brawn | Inscribed Pyrestone |
| Hands | Grips of Silent Justice | 141 | Major Strength | Bold Crimson Spinel, Bold Crimson Spinel |
| Waist | Belt of One-Hundred Deaths | 138 | — | Bold Crimson Spinel, Sovereign Shadowsong Amethyst |
| Legs | Leggings of Divine Retribution | 141 | Nethercobra Leg Armor | — |
| Feet | Dreadboots of the Legion | 141 | Cat's Swiftness | Inscribed Pyrestone, Inscribed Pyrestone |
| Finger1 | Unstoppable Aggressor's Ring | 141 | — | — |
| Finger2 | Stormrage Signet Ring | 151 | — | — |
| Trinket1 | Madness of the Betrayer | 141 | — | — |
| Trinket2 | Dragonspine Trophy | 125 | — | — |
| MainHand | Vengeful Gladiator's Cleaver | 146 | Mongoose | — |
| OffHand | Blade of Savagery | 141 | Mongoose | — |
| Ranged | Serpent Spine Longbow | 141 | — | — |

Gems shown are what the harness equipped (colour-matched). Per the free-upgrades table, all-red is
worth a few more DPS on the pieces you keep.

## Reference points

| Reference gear set | DPS | vs current |
|--------------------|-----|-----------|
| repo P2 fury preset | 2338 | +118 |
| repo P3 fury preset (Warglaives) | 2761 | +541 |
| repo P4 fury preset (Warglaives) | 2795 | +575 |

The repo's P3 Fury preset (2761) sits +39 above this report's Warglaive build. It equips the
glaive pair too, and beats build G because it is a full hand-picked BiS list rather than something
assembled from a fixed candidate pool — it also uses Madness of the Betrayer + Dragonspine Trophy and
better legs, shoulders and bracers than my pool contained for those slots. Treat it as the true
P3 ceiling, and this report as the route there from where you actually are.

Unlike the enhancement shaman presets, the warrior presets store enchants as *effect* ids, so they
resolve correctly and need no translation — see the bug section of [result_shaman.md](result_shaman.md).

## How this was measured

- **Harness**: `sim/warriorupgrade_test.go`. Re-run with
  `go test -tags with_db ./sim/ -run TestWarriorP3Upgrades -v`; set `WARRIOR_ITERATIONS` (default
  20000) to trade accuracy for speed.
- **Iterations**: 20,000 per configuration, `RandomSeed: 101`, `core.RunRaidSimConcurrent`.
  ~90 configurations, ~11 minutes.
- **Encounter**: single target, `core.MakeSingleTargetEncounter(5)` — 180s ±5s. No cleave scenarios.
- **Buffs/debuffs**: `FullRaidBuffs` / `FullPartyBuffs` / `FullIndividualBuffs` / `FullDebuffs`. No
  bonus hit rating added — your real gear hit is used, so hit upgrades are valued correctly.
- **Character**: Orc, Fury, Battle Shout, Berserker stance, 25 yd starting distance, Engineering +
  Blacksmithing, consumables matching `sim/warrior/dps/dps_test.go` (Haste Potion, Flask of
  Relentless Assault, Roasted Clefthoof, Flame Cap, Adamantite Grenade, Super Sapper, agi/str scrolls,
  and **Adamantite Sharpening Stone** as `OhImbueId: 29453` — this is the "Sharpened" line the armory
  shows on your off hand; the engine models weapon stones as a consumable, not a gem).
- **Gear source**: your equipped items, enchants and gems as reported by ironforge.pro, transcribed
  into `warriorCurrentGear()`, with Dragonspine Trophy in trinket 2 as instructed. Shirt and tabard
  are not simmed slots.
- **Replacement gemming**: each candidate is gemmed by reading its *own* sockets from
  `assets/database/db.json` and filling each with the best P3 epic gem of that colour. This matters —
  `NewItem` in `sim/core/database.go` applies every gem it is handed, even past the item's socket
  count, so a hand-written list on a socketless item would silently grant free stats.
- **Unique-equipped** is not enforced by the engine, so the assembled builds refuse to hand the same
  item id to two slots (two copies of a weapon excepted, which is legal).
- **Candidate pool**: for each slot, the Phase 3 items in this database a warrior can wear, plus the
  notable Phase 1-2 items your current piece is losing to, plus the tier 6 (Onslaught Battlegear)
  piece for that slot, plus both Warglaives. Kael'thas' Warp Slicer and Infinity Blade (ilvl 175) were
  excluded as effectively legendary-tier, matching the constraint used in [results.md](results.md).

## Caveats

- **Single target, patchwerk-style.** No movement, no adds, no execute-phase weighting beyond what the
  APL does naturally. Dragonspine Trophy, Madness of the Betrayer and the Warglaive set proc are all
  proc effects whose value shifts on fights with downtime.
- **Badge of the Swarmguard is an armor-penetration proc trinket**; how favourably it sims depends on
  this engine's armor model and the target's armor value. It loses here, but it is the one result in
  this report most sensitive to that modelling choice.
- **Meta gem colour requirements are not modelled** — the engine has no requirement field, so your
  Relentless Earthstorm Diamond always applies. Check your own blue-gem count in game.
- **The T5 2pc (Overpower grants 100 AP)** is implemented but near-worthless to Fury, since Overpower
  is a Battle-stance ability and you fight in Berserker stance. Do not read the 2pc as a reason to
  keep tier.
- **Profession assumptions.** Your goggles require Engineering, so Engineering is set. Blacksmithing
  is also set (matching the repo default), which is what Swiftsteel Shoulders would require to wear;
  it loses its slot anyway, so nothing hinges on it.
- **Phase tagging is this repo's.** Season 3 arena gear is tagged `phase: 3` in
  `assets/database/db.json`, which is why Vengeful Gladiator weapons appear as Phase 3 candidates.
- **20,000 iterations** leaves roughly ±1 DPS of noise per number. Differences under ~5 DPS between
  two candidates are not meaningful — several of the assembled builds are inside that band.
- **The build assembler is greedy, not exhaustive.** It takes the best candidate per slot and refuses
  to reuse an item id, but it does not fall back to a second choice when its first pick is claimed by
  another slot: in builds C-F the Finger2 slot therefore keeps your Shapeshifter's Signet, because
  Stormrage Signet Ring went to Finger1. Those builds are lower bounds, not optima. Build G also shows
  the deeper limitation — a slot-by-slot search cannot find a *pair* bonus like the Warglaives.
