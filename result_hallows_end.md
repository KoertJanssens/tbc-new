# Hallow's End ilvl 141 items

10,000 iterations, seed 101, single target, 5 min, full raid buffs. Each item goes
into the slot it would take, gemmed from its own sockets with the slot's enchant
kept; rings try both fingers and the better one is shown. Run with
`HALLOWS_ITERATIONS=10000 go test -tags with_db ./sim/ -run TestHallowsEnd -v`
(raw log: `result_hallows_end.raw.txt`). Standard error is about 1 DPS, so
differences under ~3 DPS are noise.

## Own characters

| Character | Item | Replaces | DPS change |
|---|---|---|---|
| Koerry 4/5 T5 | **The Horseman's Blade** (MH) | Talon of Azshara | **+9.8 (+0.44%)** |
| Koerry 5/5 T5 | **The Horseman's Blade** (MH) | Talon of Azshara | **+10.0 (+0.45%)** |
| Koerry 4/5 T5 | Ring of Ghoulish Delight | Shapeshifter's Signet | -7.7 (-0.35%) |
| Koerry 5/5 T5 | Ring of Ghoulish Delight | Shapeshifter's Signet | -6.7 (-0.30%) |
| Koerry 4/5 T5 | The Horseman's Helm | Furious Gizmatic Goggles | -35.9 (-1.61%) |
| Koerry 5/5 T5 | The Horseman's Helm | Destroyer Battle-Helm | -37.2 (-1.67%) |
| Koerrii (Enh) | Ring of Ghoulish Delight | Violet Signet | -5.7 (-0.26%) |

The Helm has no meta socket, so taking it means losing Relentless Earthstorm.

## Phase 3 presets

| Spec | Ring of Ghoulish Delight | The Horseman's Helm | The Horseman's Blade (MH) |
|---|---|---|---|
| Warrior Fury | -1.31% | -3.06% | -1.66% |
| Warrior Arms | -1.21% | -3.24% | – |
| Warrior Arms (Kebab) | -1.28% | -3.18% | -0.71% |
| Rogue | -1.15% | – | +0.08% (noise) |
| Hunter BM | -1.03% | – | – |
| Hunter SV | -0.67% | – | – |
| Enhancement Shaman | -1.09% | – | – |
| Feral Druid | -0.93% | – | – |
| Retribution Paladin | -1.11% | -2.57% | – |

| Spec | The Horseman's Signet Ring |
|---|---|
| Arcane Mage | -0.33% |
| Destruction Warlock | -1.10% |
| Shadow Priest | -0.73% |
| Elemental Shaman | -0.84% |
| Balance Druid | -0.81% |

Witches Band (healing ring) is not simmed. The Helm's laugh and the Blade's
Pumpkin Soldiers on-use are not implemented, so both count as stats only.

## Koerry: Horseman's Blade vs Blacksmithing main hands

20,000 iterations, standard Fury talents, Mongoose on every weapon, replacing
Talon of Azshara (raw log: `result_hallows_end_blacksmith.raw.txt`).

| Main hand | 4/5 T5 | 5/5 T5 |
|---|---|---|
| **Dragonstrike** (mace, haste proc) | **+35.4 (+1.59%)** | **+35.4 (+1.59%)** |
| Wicked Edge of the Planes (axe) | +29.0 (+1.30%) | +28.1 (+1.26%) |
| Dragonmaw (mace, haste proc) | +16.9 (+0.76%) | +17.5 (+0.79%) |
| The Horseman's Blade (sword) | +10.4 (+0.47%) | +9.9 (+0.45%) |

## Koerrii: Blacksmithing main hands

Shamans can't use the Horseman's Blade (sword). 20,000 iterations, Mongoose on
every weapon, replacing the main-hand Merciless Gladiator's Cleaver, with
Blacksmithing assumed so the crafted weapons can be equipped (raw log:
`result_hallows_end_blacksmith_shaman.raw.txt`).

| Main hand | DPS change |
|---|---|
| Vengeful Gladiator's Cleaver (reference) | **+25.1 (+1.12%)** |
| Wicked Edge of the Planes (axe) | -25.9 (-1.16%) |
| Dragonstrike (mace, haste proc) | -28.4 (-1.27%) |
| Dragonmaw (mace, haste proc) | -48.3 (-2.16%) |
