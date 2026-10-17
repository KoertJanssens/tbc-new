package sim

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

// Soft-reserve helper for Koerry (Fury Warrior) in Mount Hyjal and Black
// Temple: the warrior counterpart of sim/shamanhyjalsr_test.go and
// sim/shamanbtsr_test.go. Every drop from either instance that a Fury warrior
// could plausibly want is simmed as a single-slot swap.
//
// Baselines (pick with WARRIOR_SR_BASELINES, comma-separated; default all):
//   - "Oct13T5R": Oct13T5 plus Band of the Eternal Champion (exalted Scale of
//     the Sands) in place of Mithril Band of the Unscarred.
//   - "Oct13T5": the Oct13 gear with Destroyer Gauntlets back on the hands, for
//     5/5 Destroyer. One tier swap is free; Pillager's Gauntlets is a candidate.
//   - "Oct13": the live armory as of 2026-10-13 - 4/5 Destroyer (Pillager's
//     Gauntlets on hands), Empty Mug of Direbrew + Dragonspine Trophy,
//     Dragonstrike / Talon of Azshara. Any tier swap breaks the 4pc.
//   - "T5":  the 2026-08-31 PvE set from sim/warriorupgrade_test.go with
//     Destroyer Battle-Helm, 5/5 Destroyer. One tier swap is free.
//   - "4/5": the same set with Furious Gizmatic Goggles. Any tier swap breaks 4pc.
//
// Tier 6 drops as Protector tokens, so each token is simmed as the Onslaught
// Battlegear piece it buys. The Destroyer 2pc+4pc is also measured on its own
// by swapping in set-less copies of the tier pieces (same stats, no SetName).
//
// Besides the soft-reservable Hyjal/BT drops, the list carries items from other
// sources for the gear overview (softres_value.py ignores any zone that is not
// an instance of the raid): the exalted Scale of the Sands proc ring, and the
// Phase 4 (Zul'Aman era) items in the database. Those without a recorded source
// are tagged "P4".
//
// Loot tables from assets/database/db.json zones 3606 (Hyjal) and 3959 (BT).
// Left out on purpose: caster/healer/tank-only pieces, two-handers (Fury dual
// wields), and Beast-tamer's Shoulders (30892), which panics any non-hunter -
// see result_shaman_hyjal_sr.md.
//
//	WARRIOR_ITERATIONS=300 WARRIOR_SR_BASELINES=Oct13 go test -tags with_db ./sim/ -run TestWarriorHyjalBTSR -v -timeout 0

type warriorDrop struct {
	zone  string // "Hyjal", "BT", or "Hyjal+BT" for shared trash drops
	boss  string
	token string // soft-reserve name when the drop is a tier token, else ""
	c     candidate
}

func warriorHyjalBTDrops() []warriorDrop {
	return []warriorDrop{
		// ===================== Mount Hyjal =====================
		// --- Rage Winterchill ---
		{"Hyjal", "Rage Winterchill", "", candidate{slotWrist, "Deadly Cuffs", 30863, nil, 0, "leather"}},
		{"Hyjal", "Rage Winterchill", "", candidate{slotWrist, "Furious Shackles", 30861, nil, 0, "plate"}},
		{"Hyjal", "Rage Winterchill", "", candidate{slotWrist, "Bracers of the Pathfinder", 30864, nil, 0, "mail"}},
		{"Hyjal", "Rage Winterchill", "", candidate{slotShoulder, "Blood-stained Pauldrons", 30866, nil, 0, "plate"}},
		{"Hyjal", "Rage Winterchill", "", candidate{slotMainHand, "Tracker's Blade", 30865, nil, 0, "dagger 1.5"}},
		{"Hyjal", "Rage Winterchill", "", candidate{slotOffHand, "Tracker's Blade", 30865, nil, 0, "dagger 1.5"}},
		// --- Anetheron ---
		{"Hyjal", "Anetheron", "", candidate{slotWaist, "Don Alejandro's Money Belt", 30879, nil, 0, "leather"}},
		{"Hyjal", "Anetheron", "", candidate{slotFeet, "Quickstrider Moccasins", 30880, nil, 0, "mail"}},
		{"Hyjal", "Anetheron", "", candidate{slotMainHand, "Blade of Infamy", 30881, nil, 0, "sword 2.6"}},
		{"Hyjal", "Anetheron", "", candidate{slotOffHand, "Blade of Infamy", 30881, nil, 0, "sword 2.6"}},
		// --- Kaz'rogal ---
		{"Hyjal", "Kaz'rogal", "", candidate{slotWaist, "Belt of Seething Fury", 30915, nil, 0, "plate"}},
		{"Hyjal", "Kaz'rogal", "", candidate{slotWaist, "Valestalker Girdle", 30919, nil, 0, "mail"}},
		{"Hyjal", "Kaz'rogal", "", candidate{slotFeet, "Black Featherlight Boots", 30891, nil, 0, "leather"}},
		{"Hyjal", "Kaz'rogal", "", candidate{slotShoulder, "Razorfury Mantle", 30917, nil, 0, "leather"}},
		// --- Azgalor ---
		{"Hyjal", "Azgalor", "", candidate{slotLegs, "Shady Dealer's Pantaloons", 30898, nil, 0, "leather"}},
		{"Hyjal", "Azgalor", "", candidate{slotLegs, "Bow-stitched Leggings", 30900, nil, 0, "mail"}},
		{"Hyjal", "Azgalor", "", candidate{slotMainHand, "Boundless Agony", 30901, nil, 0, "dagger 1.8"}},
		{"Hyjal", "Azgalor", "", candidate{slotOffHand, "Boundless Agony", 30901, nil, 0, "dagger 1.8"}},
		{"Hyjal", "Azgalor", "Gloves of the Forgotten Protector", candidate{slotHands, "Onslaught Gauntlets (T6)", 30969, nil, 0, "token 31094"}},
		// --- Archimonde ---
		{"Hyjal", "Archimonde", "", candidate{slotChest, "Midnight Chestguard", 30905, nil, 0, "leather"}},
		{"Hyjal", "Archimonde", "", candidate{slotChest, "Mail of Fevered Pursuit", 30907, nil, 0, "mail"}},
		{"Hyjal", "Archimonde", "", candidate{slotLegs, "Legguards of Endless Rage", 30903, nil, 0, "plate"}},
		{"Hyjal", "Archimonde", "", candidate{slotRanged, "Bristleblitz Striker", 30906, nil, 0, "bow"}},
		{"Hyjal", "Archimonde", "Helm of the Forgotten Protector", candidate{slotHead, "Onslaught Battle-Helm (T6)", 30972, nil, 0, "token 31095"}},
		// --- Trash ---
		{"Hyjal+BT", "Trash", "", candidate{slotNeck, "Choker of Serrated Blades", 32591, nil, 0, ""}},
		{"Hyjal", "Trash", "", candidate{slotMainHand, "Claw of Molten Fury", 32946, nil, 0, "fist 2.7, MH only"}},
		{"Hyjal", "Trash", "", candidate{slotOffHand, "Fist of Molten Fury", 32945, nil, 0, "fist 1.5, OH only"}},

		// ===================== Black Temple =====================
		// --- High Warlord Naj'entus ---
		{"BT", "Naj'entus", "", candidate{slotShoulder, "Mantle of Darkness", 32377, nil, 0, "leather"}},
		{"BT", "Naj'entus", "", candidate{slotHands, "Fists of Mukoa", 32234, nil, 0, "mail"}},
		{"BT", "Naj'entus", "", candidate{slotMainHand, "Rising Tide", 32236, nil, 0, "axe 2.6"}},
		{"BT", "Naj'entus", "", candidate{slotOffHand, "Rising Tide", 32236, nil, 0, "axe 2.6"}},
		// --- Supremus ---
		{"BT", "Supremus", "", candidate{slotNeck, "Choker of Endless Nightmares", 32260, nil, 0, ""}},
		{"BT", "Supremus", "", candidate{slotChest, "Nether Shadow Tunic", 32252, nil, 0, "leather"}},
		{"BT", "Supremus", "", candidate{slotWrist, "Wraps of Precise Flight", 32251, nil, 0, "mail"}},
		{"BT", "Supremus", "", candidate{slotMainHand, "Syphon of the Nathrezim", 32262, nil, 0, "mace 2.8"}},
		{"BT", "Supremus", "", candidate{slotOffHand, "Syphon of the Nathrezim", 32262, nil, 0, "mace 2.8"}},
		{"BT", "Supremus", "", candidate{slotMainHand, "The Brutalizer", 32254, nil, 0, "axe 1.6, tank"}},
		{"BT", "Supremus", "", candidate{slotOffHand, "The Brutalizer", 32254, nil, 0, "axe 1.6, tank"}},
		{"BT", "Supremus", "", candidate{slotRanged, "Legionkiller", 32253, nil, 0, "crossbow"}},
		// --- Shade of Akama ---
		{"BT", "Shade of Akama", "", candidate{slotHands, "Grips of Silent Justice", 32278, nil, 0, "plate"}},
		{"BT", "Shade of Akama", "", candidate{slotShoulder, "Shoulders of the Hidden Predator", 32264, nil, 0, "mail"}},
		{"BT", "Shade of Akama", "", candidate{slotWaist, "Shadow-walker's Cord", 32265, nil, 0, "leather"}},
		{"BT", "Shade of Akama", "", candidate{slotFinger1, "Ring of Deceitful Intent", 32266, nil, 0, ""}},
		{"BT", "Shade of Akama", "", candidate{slotFinger2, "Ring of Deceitful Intent", 32266, nil, 0, ""}},
		// --- Teron Gorefiend ---
		{"BT", "Teron Gorefiend", "", candidate{slotBack, "Shadowmoon Destroyer's Drape", 32323, nil, 0, ""}},
		{"BT", "Teron Gorefiend", "", candidate{slotWrist, "Insidious Bands", 32324, nil, 0, "leather"}},
		{"BT", "Teron Gorefiend", "", candidate{slotFeet, "Softstep Boots of Tracking", 32510, nil, 0, "mail"}},
		{"BT", "Teron Gorefiend", "", candidate{slotRanged, "Twisted Blades of Zarak", 32326, nil, 0, "thrown"}},
		// --- Gurtogg Bloodboil ---
		{"BT", "Gurtogg Bloodboil", "", candidate{slotLegs, "Leggings of Divine Retribution", 32341, nil, 0, "plate"}},
		{"BT", "Gurtogg Bloodboil", "", candidate{slotChest, "Vest of Mounting Assault", 32334, nil, 0, "mail"}},
		{"BT", "Gurtogg Bloodboil", "", candidate{slotFinger1, "Unstoppable Aggressor's Ring", 32335, nil, 0, ""}},
		{"BT", "Gurtogg Bloodboil", "", candidate{slotFinger2, "Unstoppable Aggressor's Ring", 32335, nil, 0, ""}},
		{"BT", "Gurtogg Bloodboil", "", candidate{slotMainHand, "Messenger of Fate", 32269, nil, 0, "dagger 1.4"}},
		{"BT", "Gurtogg Bloodboil", "", candidate{slotOffHand, "Messenger of Fate", 32269, nil, 0, "dagger 1.4"}},
		// --- Reliquary of the Lost ---
		{"BT", "Reliquary of the Lost", "", candidate{slotFeet, "Dreadboots of the Legion", 32345, nil, 0, "plate"}},
		{"BT", "Reliquary of the Lost", "", candidate{slotHands, "Grips of Damnation", 32347, nil, 0, "leather"}},
		{"BT", "Reliquary of the Lost", "", candidate{slotWaist, "Boneweave Girdle", 32346, nil, 0, "mail"}},
		// --- Mother Shahraz ---
		{"BT", "Mother Shahraz", "", candidate{slotChest, "Heartshatter Breastplate", 32365, nil, 0, "plate"}},
		{"BT", "Mother Shahraz", "", candidate{slotFeet, "Shadowmaster's Boots", 32366, nil, 0, "leather"}},
		{"BT", "Mother Shahraz", "", candidate{slotMainHand, "Blade of Savagery", 32369, nil, 0, "sword 1.4"}},
		{"BT", "Mother Shahraz", "", candidate{slotOffHand, "Blade of Savagery", 32369, nil, 0, "sword 1.4"}},
		{"BT", "Mother Shahraz", "Pauldrons of the Forgotten Protector", candidate{slotShoulder, "Onslaught Shoulderblades (T6)", 30979, nil, 0, "token 31103"}},
		// --- The Illidari Council ---
		{"BT", "Illidari Council", "", candidate{slotHead, "Helm of the Illidari Shatterer", 32373, nil, 0, "plate"}},
		{"BT", "Illidari Council", "", candidate{slotHead, "Forest Prowler's Helm", 32376, nil, 0, "mail"}},
		{"BT", "Illidari Council", "", candidate{slotTrinket1, "Madness of the Betrayer", 32505, nil, 0, ""}},
		{"BT", "Illidari Council", "", candidate{slotTrinket2, "Madness of the Betrayer", 32505, nil, 0, ""}},
		{"BT", "Illidari Council", "Leggings of the Forgotten Protector", candidate{slotLegs, "Onslaught Greaves (T6)", 30977, nil, 0, "token 31100"}},
		// --- Illidan Stormrage ---
		{"BT", "Illidan Stormrage", "", candidate{slotHead, "Cursed Vision of Sargeras", 32235, nil, 0, "leather"}},
		{"BT", "Illidan Stormrage", "", candidate{slotFinger1, "Stormrage Signet Ring", 32497, nil, 0, ""}},
		{"BT", "Illidan Stormrage", "", candidate{slotFinger2, "Stormrage Signet Ring", 32497, nil, 0, ""}},
		{"BT", "Illidan Stormrage", "", candidate{slotMainHand, "Shard of Azzinoth", 32471, nil, 0, "dagger 1.9"}},
		{"BT", "Illidan Stormrage", "", candidate{slotOffHand, "Shard of Azzinoth", 32471, nil, 0, "dagger 1.9"}},
		{"BT", "Illidan Stormrage", "", candidate{slotMainHand, "Warglaive of Azzinoth (MH)", 32837, nil, 0, "legendary, MH only"}},
		{"BT", "Illidan Stormrage", "", candidate{slotOffHand, "Warglaive of Azzinoth (OH)", 32838, nil, 0, "legendary, OH only"}},
		{"BT", "Illidan Stormrage", "Chestguard of the Forgotten Protector", candidate{slotChest, "Onslaught Breastplate (T6)", 30975, nil, 0, "token 31091"}},
		// --- Trash ---
		{"BT", "Trash", "", candidate{slotFinger1, "Band of Devastation", 32526, nil, 0, ""}},
		{"BT", "Trash", "", candidate{slotFinger2, "Band of Devastation", 32526, nil, 0, ""}},
		{"BT", "Trash", "", candidate{slotHands, "Pillager's Gauntlets", 32608, nil, 0, "plate"}},
		{"BT", "Trash", "", candidate{slotWaist, "Girdle of the Lightbearer", 32606, nil, 0, "plate"}},
		{"BT", "Trash", "", candidate{slotFeet, "Treads of the Den Mother", 32593, nil, 0, "leather"}},
		{"BT", "Trash", "", candidate{slotMainHand, "Swiftsteel Bludgeon", 32943, nil, 0, "mace 1.5"}},
		{"BT", "Trash", "", candidate{slotOffHand, "Swiftsteel Bludgeon", 32943, nil, 0, "mace 1.5"}},

		// ===================== Not soft-reservable =====================
		// --- Scale of the Sands, exalted: 160 AP proc, 60s cooldown ---
		{"Rep", "Scale of the Sands", "", candidate{slotFinger1, "Band of the Eternal Champion", 29301, nil, 0, ""}},
		{"Rep", "Scale of the Sands", "", candidate{slotFinger2, "Band of the Eternal Champion", 29301, nil, 0, ""}},
		// --- Phase 4: Zul'Aman ---
		{"ZA", "Nalorakk", "", candidate{slotShoulder, "Pauldrons of Primal Fury", 33206, nil, 0, "mail"}},
		{"ZA", "Nalorakk", "", candidate{slotWaist, "Bladeangel's Money Belt", 33211, nil, 0, "leather"}},
		{"ZA", "Nalorakk", "", candidate{slotOffHand, "Fury", 33640, nil, 0, "fist 2.6, OH only"}},
		{"ZA", "Akil'zon", "", candidate{slotChest, "Bloodstained Elven Battlevest", 33215, nil, 0, "plate"}},
		{"ZA", "Akil'zon", "", candidate{slotMainHand, "Akil'zon's Talonblade", 33214, nil, 0, "sword 1.5"}},
		{"ZA", "Akil'zon", "", candidate{slotOffHand, "Akil'zon's Talonblade", 33214, nil, 0, "sword 1.5"}},
		{"ZA", "Jan'alai", "", candidate{slotChest, "Shadowtooth Trollskin Cuirass", 33329, nil, 0, "leather"}},
		{"ZA", "Halazzi", "", candidate{slotFeet, "Skullshatter Warboots", 33303, nil, 0, "plate"}},
		{"ZA", "Halazzi", "", candidate{slotNeck, "The Savage's Choker", 33297, nil, 0, ""}},
		{"ZA", "Halazzi", "", candidate{slotShoulder, "Shoulderpads of Dancing Blades", 33300, nil, 0, "leather"}},
		{"ZA", "Hex Lord Malacrass", "", candidate{slotMainHand, "Heartless", 33388, nil, 0, "sword 2.6"}},
		{"ZA", "Hex Lord Malacrass", "", candidate{slotOffHand, "Heartless", 33388, nil, 0, "sword 2.6"}},
		{"ZA", "Hex Lord Malacrass", "", candidate{slotMainHand, "Prowler's Strikeblade", 33298, nil, 0, "dagger 1.5"}},
		{"ZA", "Hex Lord Malacrass", "", candidate{slotOffHand, "Prowler's Strikeblade", 33298, nil, 0, "dagger 1.5"}},
		{"ZA", "Hex Lord Malacrass", "", candidate{slotOffHand, "Dagger of Bad Mojo", 33389, nil, 0, "dagger 1.8"}},
		{"ZA", "Zul'jin", "", candidate{slotHead, "Grimgrin Faceguard", 33479, nil, 0, "leather"}},
		{"ZA", "Zul'jin", "", candidate{slotTrinket1, "Berserker's Call", 33831, nil, 0, ""}},
		{"ZA", "Zul'jin", "", candidate{slotTrinket2, "Berserker's Call", 33831, nil, 0, ""}},
		{"ZA", "Zul'jin", "", candidate{slotRanged, "Ancient Amani Longbow", 33474, nil, 0, "bow"}},
		{"ZA", "Timed chest", "", candidate{slotFinger1, "Signet of Primal Wrath", 33496, nil, 0, ""}},
		{"ZA", "Timed chest", "", candidate{slotFinger2, "Signet of Primal Wrath", 33496, nil, 0, ""}},
		{"ZA", "Timed chest", "", candidate{slotBack, "Cloak of Fiends", 33590, nil, 0, ""}},
		{"ZA", "Timed chest", "", candidate{slotMainHand, "Rage", 33495, nil, 0, "fist 2.6, MH only"}},
		{"ZA", "Timed chest", "", candidate{slotRanged, "Tuskbreaker", 33491, nil, 0, "gun"}},
		// --- Phase 4, no source recorded in the database ---
		{"P4", "?", "", candidate{slotHead, "Amani Mask of Death", 33810, nil, 0, "plate"}},
		{"P4", "?", "", candidate{slotNeck, "Brooch of Deftness", 33296, nil, 0, ""}},
		{"P4", "?", "", candidate{slotShoulder, "Pauldrons of Gruesome Fate", 33514, nil, 0, "plate"}},
		{"P4", "?", "", candidate{slotWrist, "Eternium Rage-shackles", 33513, nil, 0, "plate"}},
		{"P4", "?", "", candidate{slotWrist, "Master Assassin Wristwraps", 33540, nil, 0, "leather"}},
		{"P4", "?", "", candidate{slotHands, "Furious Deathgrips", 33512, nil, 0, "plate"}},
		{"P4", "?", "", candidate{slotHands, "Trickster's Stickyfingers", 33539, nil, 0, "leather"}},
		{"P4", "?", "", candidate{slotWaist, "Chain of Unleashed Rage", 33331, nil, 0, "plate"}},
		{"P4", "?", "", candidate{slotWaist, "Waistguard of the Great Beast", 33583, nil, 0, "leather"}},
		{"P4", "?", "", candidate{slotLegs, "Bloodthirster's Wargreaves", 33501, nil, 0, "plate"}},
		{"P4", "?", "", candidate{slotLegs, "Shallow-grave Trousers", 33538, nil, 0, "leather"}},
		{"P4", "?", "", candidate{slotFeet, "Nyn'jah's Tabi Boots", 33222, nil, 0, "leather"}},
		// --- Season 3: Vengeful Gladiator's (arena) and Vindicator's (honor) ---
		{"S3", "Arena", "", candidate{slotHead, "Vengeful Gladiator's Plate Helm", 33730, nil, 0, "plate, PvP set"}},
		{"S3", "Arena", "", candidate{slotShoulder, "Vengeful Gladiator's Plate Shoulders", 33732, nil, 0, "plate, PvP set"}},
		{"S3", "Arena", "", candidate{slotChest, "Vengeful Gladiator's Plate Chestpiece", 33728, nil, 0, "plate, PvP set"}},
		{"S3", "Arena", "", candidate{slotHands, "Vengeful Gladiator's Plate Gauntlets", 33729, nil, 0, "plate, PvP set"}},
		{"S3", "Arena", "", candidate{slotLegs, "Vengeful Gladiator's Plate Legguards", 33731, nil, 0, "plate, PvP set"}},
		{"S3", "Honor", "", candidate{slotNeck, "Vindicator's Pendant of Triumph", 33923, nil, 0, ""}},
		{"S3", "Honor", "", candidate{slotWrist, "Vindicator's Plate Bracers", 33813, nil, 0, "plate"}},
		{"S3", "Honor", "", candidate{slotWaist, "Vindicator's Plate Belt", 33811, nil, 0, "plate"}},
		{"S3", "Honor", "", candidate{slotFeet, "Vindicator's Plate Greaves", 33812, nil, 0, "plate"}},
		{"S3", "Arena", "", candidate{slotMainHand, "Vengeful Gladiator's Cleaver", 33669, nil, 0, "axe 2.6"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Cleaver", 33669, nil, 0, "axe 2.6"}},
		{"S3", "Arena", "", candidate{slotMainHand, "Vengeful Gladiator's Slicer", 33762, nil, 0, "sword 2.6"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Slicer", 33762, nil, 0, "sword 2.6"}},
		{"S3", "Arena", "", candidate{slotMainHand, "Vengeful Gladiator's Pummeler", 33733, nil, 0, "mace 2.6"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Pummeler", 33733, nil, 0, "mace 2.6"}},
		{"S3", "Arena", "", candidate{slotMainHand, "Vengeful Gladiator's Right Ripper", 33737, nil, 0, "fist 2.6, MH only"}},
		{"S3", "Arena", "", candidate{slotMainHand, "Vengeful Gladiator's Shanker", 33754, nil, 0, "dagger 1.8"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Shanker", 33754, nil, 0, "dagger 1.8"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Hacker", 33689, nil, 0, "axe 1.5, OH only"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Quickblade", 33734, nil, 0, "sword 1.5, OH only"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Bonecracker", 33662, nil, 0, "mace 1.5, OH only"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Left Ripper", 33705, nil, 0, "fist 1.5, OH only"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Shiv", 33756, nil, 0, "dagger 1.4, OH only"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Mutilator", 33801, nil, 0, "dagger 1.8, OH only"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Chopper", 34015, nil, 0, "axe 2.6, OH only"}},
		{"S3", "Arena", "", candidate{slotOffHand, "Vengeful Gladiator's Left Render", 34016, nil, 0, "fist 2.6, OH only"}},
		{"S3", "Arena", "", candidate{slotRanged, "Vengeful Gladiator's War Edge", 33765, nil, 0, "thrown"}},
		{"S3", "Arena", "", candidate{slotRanged, "Vengeful Gladiator's Heavy Crossbow", 33006, nil, 0, "crossbow"}},
		{"S3", "Arena", "", candidate{slotRanged, "Vengeful Gladiator's Longbow", 34529, nil, 0, "bow"}},
		{"S3", "Arena", "", candidate{slotRanged, "Vengeful Gladiator's Rifle", 34530, nil, 0, "gun"}},
	}
}

// Koerry's raid gear from the armory as of 2026-10-13 (spec Fury), enchants
// translated to effect ids. The off hand's sharpening stone is the
// OhImbueId consumable in warriorPlayer, not a gem.
func warriorArmoryOct13Gear() *proto.EquipmentSpec {
	return &proto.EquipmentSpec{Items: []*proto.ItemSpec{
		itemSpec(30120, wEnchHeadGlyphOfFerocity, 32409, 24054), // Destroyer Battle-Helm (T5)
		itemSpec(30022, 0), // Pendant of the Perilous
		itemSpec(30122, wEnchShoulderGrInscrVeng, 24058, 24054),     // Destroyer Shoulderblades (T5)
		itemSpec(28777, wEnchCloakGreaterAgility),                   // Cloak of the Pit Stalker
		itemSpec(30118, wEnchChestExceptional, 24058, 31118, 24058), // Destroyer Breastplate (T5)
		itemSpec(28795, wEnchBracerBrawn, 24054, 30571),             // Bladespire Warbands
		itemSpec(32608, wEnchGlovesMajorStrength),                   // Pillager's Gauntlets
		itemSpec(30032, 0, 24054, 24058),                            // Red Belt of Battle
		itemSpec(30121, wEnchLegsNethercobra, 24058),                // Destroyer Greaves (T5)
		itemSpec(28608, wEnchBootsCatsSwiftness, 28362, 28363),      // Ironstriders of Urgency
		itemSpec(28730, 0),                   // Mithril Band of the Unscarred
		itemSpec(33919, 0),                   // Vindicator's Band of Triumph
		itemSpec(281739, 0),                  // Empty Mug of Direbrew
		itemSpec(28830, 0),                   // Dragonspine Trophy
		itemSpec(28439, wEnchWeaponMongoose), // Dragonstrike (mace, 2.7)
		itemSpec(30082, wEnchWeaponMongoose), // Talon of Azshara (sword, 2.7)
		itemSpec(32054, 0),                   // Merciless Gladiator's War Edge
	}}
}

// The Oct13 gear with Destroyer Gauntlets instead of Pillager's Gauntlets, for a
// 5/5 Destroyer baseline with one free tier swap.
func warriorArmoryOct13T5Gear() *proto.EquipmentSpec {
	g := warriorArmoryOct13Gear()
	g.Items[slotHands] = itemSpec(30119, wEnchGlovesMajorStrength) // Destroyer Gauntlets (T5)
	return g
}

// Oct13T5 wearing the exalted Hyjal reputation ring instead of Mithril Band of the
// Unscarred, the weaker of the two rings.
func warriorArmoryOct13T5RingGear() *proto.EquipmentSpec {
	g := warriorArmoryOct13T5Gear()
	g.Items[slotFinger1] = itemSpec(29301, 0) // Band of the Eternal Champion
	return g
}

// Two-item sets whose bonus a single-slot swap cannot see.
type warriorPair struct {
	zone, boss, label string
	mh, oh            int32
}

var warriorHyjalBTPairs = []warriorPair{
	{"BT", "Illidan Stormrage", "Warglaives of Azzinoth (pair)", 32837, 32838},
	{"Hyjal", "Trash", "Claw + Fist of Molten Fury (pair)", 32946, 32945},
}

// setless returns gear in which every Destroyer Battlegear piece is replaced by
// a stat-identical copy that belongs to no set, so the sim applies no 2pc/4pc.
// Set membership is read from Item.SetName at equip time (getSetBonuses in
// sim/core/item_sets.go), so a copy under a spare id is enough.
func setless(e *proto.EquipmentSpec) *proto.EquipmentSpec {
	const fakeIDBase = 900000
	g := cloneGear(e)
	for _, it := range g.Items {
		orig := core.GetItemByID(it.Id)
		if orig == nil || orig.SetName != "Destroyer Battlegear" {
			continue
		}
		fakeID := fakeIDBase + orig.ID
		core.AddToDatabase(&proto.SimDatabase{Items: []*proto.SimItem{{
			Id:               fakeID,
			Name:             orig.Name + " (set-less)",
			Type:             orig.Type,
			ArmorType:        orig.ArmorType,
			WeaponType:       orig.WeaponType,
			HandType:         orig.HandType,
			RangedWeaponType: orig.RangedWeaponType,
			GemSockets:       orig.GemSockets,
			SocketBonus:      orig.SocketBonus.ToProtoArray(),
			QualityModifier:  orig.QualityModifier,
			WeaponSpeed:      orig.SwingSpeed,
			ScalingOptions:   orig.ScalingOptions,
			ItemEffects:      orig.ItemEffects,
			LimitCategory:    orig.LimitCategory,
			Unique:           orig.Unique,
		}}})
		it.Id = fakeID
	}
	return g
}

type warriorSRBaseline struct {
	key   string
	gear  func() *proto.EquipmentSpec
	label string
}

// The baselines selected by WARRIOR_SR_BASELINES (comma-separated keys), or all.
func warriorSRBaselines() []warriorSRBaseline {
	all := []warriorSRBaseline{
		{"Oct13T5R", warriorArmoryOct13T5RingGear, "Oct13T5 + Band of the Eternal Champion"},
		{"Oct13T5", warriorArmoryOct13T5Gear, "armory 2026-10-13 + Destroyer Gauntlets, 5/5 Destroyer"},
		{"Oct13", warriorArmoryOct13Gear, "armory 2026-10-13, 4/5 Destroyer, Empty Mug"},
		{"T5", warriorT5Gear, "5/5 Destroyer (T5 helm)"},
		{"4/5", warriorCurrentGear, "4/5 Destroyer (Furious Gizmatic Goggles)"},
	}
	wanted := map[string]bool{}
	for _, k := range strings.Split(os.Getenv("WARRIOR_SR_BASELINES"), ",") {
		if k = strings.TrimSpace(k); k != "" {
			wanted[k] = true
		}
	}
	out := []warriorSRBaseline{}
	for _, b := range all {
		if len(wanted) == 0 || wanted[b.key] {
			out = append(out, b)
		}
	}
	return out
}

// tierSlotsWorn lists the slots holding a Destroyer piece.
func tierSlotsWorn(e *proto.EquipmentSpec) []string {
	slots := []string{}
	for i, it := range e.Items {
		if item := core.GetItemByID(it.Id); item != nil && item.SetName == "Destroyer Battlegear" {
			slots = append(slots, warriorSlotNames[i])
		}
	}
	return slots
}

func TestWarriorHyjalBTSR(t *testing.T) {
	RegisterAll()
	iterations := warriorIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)
	fmt.Printf("ITERATIONS %d\n", iterations)

	for _, b := range warriorSRBaselines() {
		// applyWarrior inherits enchants from warriorBaseGear.
		warriorBaseGear = b.gear
		run := func(gear *proto.EquipmentSpec) simResult {
			return runWarrior(t, gear, talentsFuryStandard, "fury", encounter, iterations)
		}

		base := run(b.gear())
		noSet := run(setless(b.gear()))
		fmt.Printf("WSRBASE | %s | %s | DPS %.1f (stdev %.1f) | AP %.0f hit %.0f crit %.0f exp %.0f\n",
			b.key, b.label, base.dps, base.stdev, base.ap, base.hit, base.crit, base.exp)
		fmt.Printf("WSRSET | %s | Destroyer set bonuses removed | DPS %.1f | %+.1f\n", b.key, noSet.dps, noSet.dps-base.dps)
		tierSlots := tierSlotsWorn(b.gear())
		fmt.Printf("WSRTIER | %s | %d | %s\n", b.key, len(tierSlots), strings.Join(tierSlots, ","))

		type scoredDrop struct {
			d warriorDrop
			r simResult
		}
		all := []scoredDrop{}
		for _, d := range warriorHyjalBTDrops() {
			item := core.GetItemByID(d.c.id)
			if d.c.slot == slotOffHand && item.HandType == proto.HandType_HandTypeMainHand ||
				d.c.slot == slotMainHand && item.HandType == proto.HandType_HandTypeOffHand {
				t.Fatalf("%s cannot go in %s", d.c.label, warriorSlotNames[d.c.slot])
			}
			gear := cloneGear(b.gear())
			d.c.applyWarrior(gear)
			r := run(gear)
			all = append(all, scoredDrop{d, r})
			fmt.Printf("WSR | %s | %-8s | %-18s | %-9s | %-34s | %6d | %8.1f | %+7.1f | %s | %s\n",
				b.key, d.zone, d.boss, warriorSlotNames[d.c.slot], d.c.label, d.c.id, r.dps, r.dps-base.dps, d.token, d.c.note)
		}
		for _, p := range warriorHyjalBTPairs {
			gear := cloneGear(b.gear())
			gear.Items[slotMainHand] = itemSpec(p.mh, gear.Items[slotMainHand].Enchant)
			gear.Items[slotOffHand] = itemSpec(p.oh, gear.Items[slotOffHand].Enchant)
			r := run(gear)
			fmt.Printf("WSRPAIR | %s | %-8s | %-18s | %-34s | %6d+%d | %8.1f | %+7.1f\n",
				b.key, p.zone, p.boss, p.label, p.mh, p.oh, r.dps, r.dps-base.dps)
		}

		sort.Slice(all, func(i, j int) bool { return all[i].r.dps > all[j].r.dps })
		fmt.Printf("\n=== RANKED (%s) ===\n", b.key)
		for i, s := range all {
			fmt.Printf("WSRRANK | %s | %2d | %-8s | %-18s | %-9s | %-34s | %+7.1f\n",
				b.key, i+1, s.d.zone, s.d.boss, warriorSlotNames[s.d.c.slot], s.d.c.label, s.r.dps-base.dps)
		}
		fmt.Println()
	}
}

// ---------------------------------------------------------------------------
// Adaptive variant: the baseline gets a fixed, large run, and every candidate is
// simmed in batches. A candidate stops early only once it is a proven downgrade
// (delta below zero at 95%); upgrades and undecided candidates run to the cap, so
// every upgrade's size is known, not just its sign.
//
// Seeds: the engine seeds iteration i of a run with RandomSeed+i (sim/core/sim.go,
// and sim_concurrent.go hands each worker the next block), so a run covers the
// seed range [RandomSeed, RandomSeed+iterations). Consecutive seeds would make
// consecutive batches near-copies of each other. Each batch therefore starts
// where the previous one ended, and candidates draw from a range disjoint from
// the baseline's, so batches pool into one independent sample and the delta's
// standard error is exact for independent runs.
//
// Output uses the same WSR line layout as TestWarriorHyjalBTSR, so
// softres_value.py reads it unchanged, with the iteration count, 95% interval
// and verdict appended.
//
//	WARRIOR_SR_BASELINES=Oct13T5R go test -tags with_db ./sim/ -run TestWarriorHyjalBTSRAdaptive -v -timeout 0
//
// Every batch also prints a WSRSTEP line (slot, item, iterations so far, delta,
// 95% half-width) for convergence plots. Knobs: WARRIOR_BASE_ITERATIONS (25000),
// WARRIOR_BATCH_ITERATIONS (5000),
// WARRIOR_MAX_ITERATIONS (25000 per candidate), WARRIOR_REFINE_IDS /
// WARRIOR_REFINE_ITERATIONS (fixed-length runs, default 100000), WARRIOR_TARGETS
// (cleave: number of boss targets, default 1), WARRIOR_CANDIDATE_SLOTS (e.g.
// MainHand,OffHand for a weapons-only run).
// ---------------------------------------------------------------------------

const adaptiveZ = 1.96 // two-sided 95%

func envInt(name string, def int) int {
	if v := os.Getenv(name); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return def
}

// dpsSample pools the per-iteration DPS of several runs from their means and
// standard deviations.
type dpsSample struct{ n, sum, sumSq float64 }

func (s *dpsSample) add(n int32, mean, stdev float64) {
	f := float64(n)
	s.n += f
	s.sum += f * mean
	s.sumSq += f * (stdev*stdev + mean*mean)
}

func (s dpsSample) mean() float64 { return s.sum / s.n }

func (s dpsSample) stdErr() float64 {
	m := s.mean()
	return math.Sqrt(math.Max(s.sumSq/s.n-m*m, 0) / s.n)
}

// runWarriorSeeded is runWarrior with a caller-chosen seed and no stat sheet.
func runWarriorSeeded(t *testing.T, gear *proto.EquipmentSpec, encounter *proto.Encounter, iterations int32, seed int64) (float64, float64) {
	raid := core.SinglePlayerRaidProto(warriorPlayer(gear, talentsFuryStandard, "fury"), core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs)
	result := core.RunRaidSimConcurrent(&proto.RaidSimRequest{
		Raid:       raid,
		Encounter:  encounter,
		SimOptions: &proto.SimOptions{Iterations: iterations, RandomSeed: seed},
	})
	if result.Error != nil {
		t.Fatalf("sim failed: %s", result.Error.Message)
	}
	return result.RaidMetrics.Dps.Avg, result.RaidMetrics.Dps.Stdev
}

func TestWarriorHyjalBTSRAdaptive(t *testing.T) {
	RegisterAll()
	encounter := core.MakeSingleTargetEncounter(5)
	// WARRIOR_TARGETS > 1 adds default boss targets; the fury APL switches to
	// Whirlwind/Cleave priorities when numberTargets > 1.
	targets := envInt("WARRIOR_TARGETS", 1)
	if targets > 1 {
		encounter.Targets = nil
		for i := 0; i < targets; i++ {
			encounter.Targets = append(encounter.Targets, core.NewDefaultTarget())
		}
	}
	// WARRIOR_CANDIDATE_SLOTS limits the candidates to some slots, e.g. MainHand,OffHand.
	slotFilter := map[string]bool{}
	for _, v := range strings.Split(os.Getenv("WARRIOR_CANDIDATE_SLOTS"), ",") {
		if v = strings.TrimSpace(v); v != "" {
			slotFilter[v] = true
		}
	}
	baseIters := int32(envInt("WARRIOR_BASE_ITERATIONS", 25000))
	batch := int32(envInt("WARRIOR_BATCH_ITERATIONS", 5000))
	maxIters := int32(envInt("WARRIOR_MAX_ITERATIONS", 25000))
	// Items listed in WARRIOR_REFINE_IDS skip the significance stop and run to
	// WARRIOR_REFINE_ITERATIONS, e.g. to size the gain of an item already owned.
	refineIters := int32(envInt("WARRIOR_REFINE_ITERATIONS", 100000))
	refine := map[int32]bool{}
	for _, v := range strings.Split(os.Getenv("WARRIOR_REFINE_IDS"), ",") {
		if id, err := strconv.Atoi(strings.TrimSpace(v)); err == nil {
			refine[int32(id)] = true
		}
	}
	fmt.Printf("ITERATIONS adaptive | baseline %d | batch %d | cap %d | z %.2f | targets %d\n", baseIters, batch, maxIters, adaptiveZ, targets)

	// Baseline and set-less runs share their seed range on purpose (common random
	// numbers make the set value tighter); candidates start far above both.
	const baseSeed, setSeed, firstBatchSeed = 101, 101, 1_000_000

	for _, b := range warriorSRBaselines() {
		if targets > 1 {
			b.key = fmt.Sprintf("%s@%dt", b.key, targets)
			b.label = fmt.Sprintf("%s, %d targets", b.label, targets)
		}
		warriorBaseGear = b.gear // applyWarrior inherits enchants from it
		stats := runWarrior(t, b.gear(), talentsFuryStandard, "fury", encounter, 1)

		var base, noSet dpsSample
		m, sd := runWarriorSeeded(t, b.gear(), encounter, baseIters, baseSeed)
		base.add(baseIters, m, sd)
		m, sd = runWarriorSeeded(t, setless(b.gear()), encounter, baseIters, setSeed)
		noSet.add(baseIters, m, sd)
		fmt.Printf("WSRBASE | %s | %s | DPS %.1f (se %.2f, %d iterations) | AP %.0f hit %.0f crit %.0f exp %.0f\n",
			b.key, b.label, base.mean(), base.stdErr(), baseIters, stats.ap, stats.hit, stats.crit, stats.exp)
		fmt.Printf("WSRSET | %s | Destroyer set bonuses removed | DPS %.1f | %+.1f\n", b.key, noSet.mean(), noSet.mean()-base.mean())
		tierSlots := tierSlotsWorn(b.gear())
		fmt.Printf("WSRTIER | %s | %d | %s\n", b.key, len(tierSlots), strings.Join(tierSlots, ","))

		// Round-robin: every round gives each undecided candidate one more batch, so
		// the whole list sharpens together and the log reads round by round. A
		// candidate's WSR line is printed the moment it is decided (or hits the
		// cap), so a partial log is already usable.
		type adaptiveCand struct {
			d         warriorDrop
			gear      *proto.EquipmentSpec
			s         dpsSample
			delta, se float64
			verdict   string // "" while undecided
		}
		cands := []*adaptiveCand{}
		worn := map[int32]bool{}
		for _, it := range b.gear().Items {
			worn[it.Id] = true
		}
		for _, d := range warriorHyjalBTDrops() {
			if worn[d.c.id] {
				continue // already in the baseline, e.g. the Hyjal ring on Oct13T5R
			}
			if len(slotFilter) > 0 && !slotFilter[warriorSlotNames[d.c.slot]] {
				continue
			}
			gear := cloneGear(b.gear())
			d.c.applyWarrior(gear)
			cands = append(cands, &adaptiveCand{d: d, gear: gear})
		}
		verdicts := map[string]int{}
		finish := func(c *adaptiveCand, verdict string) {
			c.verdict = verdict
			verdicts[verdict]++
			fmt.Printf("WSR | %s | %-8s | %-18s | %-9s | %-34s | %6d | %8.1f | %+7.1f | %s | %s | n=%d | ±%.1f | %s\n",
				b.key, c.d.zone, c.d.boss, warriorSlotNames[c.d.c.slot], c.d.c.label, c.d.c.id, c.s.mean(), c.delta,
				c.d.token, c.d.c.note, int(c.s.n), adaptiveZ*c.se, verdict)
		}

		totalIters := int64(0)
		for k := int64(0); ; k++ {
			open := 0
			for _, c := range cands {
				if c.verdict != "" {
					continue
				}
				m, sd := runWarriorSeeded(t, c.gear, encounter, batch, firstBatchSeed+k*int64(batch))
				c.s.add(batch, m, sd)
				totalIters += int64(batch)
				c.delta = c.s.mean() - base.mean()
				c.se = math.Hypot(c.s.stdErr(), base.stdErr())
				fmt.Printf("WSRSTEP | %s | %s | %s | %d | %d | %.2f | %.2f | %s | %s\n",
					b.key, warriorSlotNames[c.d.c.slot], c.d.c.label, c.d.c.id, int(c.s.n), c.delta, adaptiveZ*c.se,
					c.d.zone, c.d.boss)
				switch {
				case refine[c.d.c.id]:
					// Fixed-length run for items whose exact size matters.
					if int32(c.s.n) >= refineIters {
						finish(c, map[bool]string{true: "better", false: "worse"}[c.delta > 0]+" (fixed)")
					} else {
						open++
					}
				case c.delta < -adaptiveZ*c.se:
					// Proven downgrade: no need to know by how much.
					finish(c, "worse")
				case int32(c.s.n) >= maxIters:
					// Upgrades run to the cap so their size is known, not just their sign.
					if c.delta > adaptiveZ*c.se {
						finish(c, "better")
					} else {
						finish(c, "inconclusive")
					}
				default:
					open++
				}
			}
			fmt.Printf("WSRROUND | %s | round %d | %d iterations each | decided %d | open %d\n",
				b.key, k+1, int64(batch)*(k+1), len(cands)-open, open)
			if open == 0 {
				break
			}
		}
		fmt.Printf("WSRSUMMARY | %s | %d candidates | better %d | worse %d | inconclusive %d | %d iterations\n",
			b.key, len(cands), verdicts["better"], verdicts["worse"], verdicts["inconclusive"], totalIters)
	}
}

// Direct check of whole setups rather than summed single swaps. WARRIOR_SETUPS is
// a ';'-separated list of name=Slot:id,Slot:id applied to the first selected
// baseline (gems and enchants as in applyWarrior); every setup runs on the same
// seed range, so differences between them are tighter than the printed ± (which
// assumes independent runs).
//
//	WARRIOR_SR_BASELINES=Oct13T5R WARRIOR_SETUPS='today=Hands:32608;s3chest=Chest:33728' \
//	  go test -tags with_db ./sim/ -run TestWarriorSetups -v
func TestWarriorSetups(t *testing.T) {
	RegisterAll()
	encounter := core.MakeSingleTargetEncounter(5)
	iters := int32(envInt("WARRIOR_SETUP_ITERATIONS", 50000))
	b := warriorSRBaselines()[0]
	warriorBaseGear = b.gear
	slotIndex := map[string]int{}
	for i, n := range warriorSlotNames {
		slotIndex[n] = i
	}

	type setup struct {
		name string
		gear *proto.EquipmentSpec
	}
	setups := []setup{{"baseline (" + b.key + ")", b.gear()}}
	for _, spec := range strings.Split(os.Getenv("WARRIOR_SETUPS"), ";") {
		name, swaps, ok := strings.Cut(strings.TrimSpace(spec), "=")
		if !ok {
			continue
		}
		gear := cloneGear(b.gear())
		for _, sw := range strings.Split(swaps, ",") {
			slot, id, _ := strings.Cut(strings.TrimSpace(sw), ":")
			n, err := strconv.Atoi(id)
			if _, known := slotIndex[slot]; !known || err != nil {
				t.Fatalf("bad swap %q in setup %q", sw, name)
			}
			candidate{slotIndex[slot], "", int32(n), nil, 0, ""}.applyWarrior(gear)
		}
		setups = append(setups, setup{name, gear})
	}

	var first dpsSample
	for i, s := range setups {
		var smp dpsSample
		m, sd := runWarriorSeeded(t, s.gear, encounter, iters, 1_000_000)
		smp.add(iters, m, sd)
		if i == 0 {
			first = smp
		}
		tier := tierSlotsWorn(s.gear)
		fmt.Printf("SETUP | %-28s | DPS %8.1f | vs baseline %+6.1f ±%.1f | %d iterations | tier %d (%s)\n",
			s.name, smp.mean(), smp.mean()-first.mean(), adaptiveZ*math.Hypot(smp.stdErr(), first.stdErr()),
			iters, len(tier), strings.Join(tier, ","))
	}
}
