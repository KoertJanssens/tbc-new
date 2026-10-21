package sim

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Benchmark: Phase 3 upgrade finder for the warrior "Koerry" (Spineshatter-EU),
// whose current equipment is transcribed below from
// https://ironforge.pro/anniversary/player/Spineshatter/Koerry/ .
//
// Two deliberate departures from the armory, both at the owner's instruction:
//   - Trinket2 is Dragonspine Trophy, not the Solarian's Sapphire the armory shows.
//   - The character is Arms-specced on the armory (both saved specs). This is a
//     *Fury* report by request: every single row uses the repo's standard Fury
//     talents (3500501130201-05050005505012050115) and the fury APL
//     (ui/warrior/dps/apls/fury.apl.json). No Arms talents and no arms APL appear
//     anywhere in this benchmark.
//
// Same shape as sim/shamanupgrade_test.go: single-slot swaps against the real
// gear, then a cumulative pass that exposes tier-set interactions (here the 4pc
// Destroyer Battlegear, of which exactly 4 pieces are worn, so any tier swap
// breaks it).
//
// Re-run with:
//
//	SHAMAN_ITERATIONS unused here; set WARRIOR_ITERATIONS to trade accuracy for speed, e.g.
//	WARRIOR_ITERATIONS=300 go test -tags with_db ./sim/ -run TestWarriorP3Upgrades -v

var warriorSlotNames = [17]string{"Head", "Neck", "Shoulder", "Back", "Chest", "Wrist", "Hands",
	"Waist", "Legs", "Feet", "Finger1", "Finger2", "Trinket1", "Trinket2", "MainHand", "OffHand", "Ranged"}

// Enchant effect ids (see the note in shamanupgrade_test.go: the engine resolves
// ItemSpec.Enchant through core.EnchantsByEffectID and silently drops anything
// else). The warrior presets in ui/warrior/dps/gear_sets already use effect ids.
const (
	wEnchHeadGlyphOfFerocity = 3003 // +34 AP, +16 hit
	wEnchShoulderGrInscrVeng = 2986 // +30 AP, +10 crit
	wEnchCloakGreaterAgility = 368  // +12 agi
	wEnchChestExceptional    = 2661 // +6 all stats
	wEnchBracerBrawn         = 2647 // +12 str
	wEnchGlovesMajorStrength = 684  // +15 str
	wEnchLegsNethercobra     = 3012 // +50 AP, +12 crit
	wEnchBootsCatsSwiftness  = 2939 // +6 agi, minor speed
	wEnchBootsDexterity      = 2657 // +12 agi
	wEnchWeaponMongoose      = 2673 // proc: 120 agi / 2% haste
	wEnchWeaponExecutioner   = 3225 // proc: 840 armor penetration
)

// Koerry's equipment, verbatim from the armory except Trinket2 (Dragonspine
// Trophy replaces Solarian's Sapphire, per the owner). The offhand's "Sharpened"
// line in the armory is an Adamantite Sharpening Stone, which this engine models
// through ConsumesSpec.OhImbueId (29453) rather than as a gem, so it is not
// listed here.
func warriorCurrentGear() *proto.EquipmentSpec {
	return &proto.EquipmentSpec{Items: []*proto.ItemSpec{
		itemSpec(32461, wEnchHeadGlyphOfFerocity, 32409, 24054), // Furious Gizmatic Goggles (Engineering)
		itemSpec(30022, 0), // Pendant of the Perilous
		itemSpec(30122, wEnchShoulderGrInscrVeng, 24058, 24054),     // Destroyer Shoulderblades (T5)
		itemSpec(28777, wEnchCloakGreaterAgility),                   // Cloak of the Pit Stalker
		itemSpec(30118, wEnchChestExceptional, 24058, 31118, 24058), // Destroyer Breastplate (T5)
		itemSpec(28795, wEnchBracerBrawn, 23111, 30571),             // Bladespire Warbands
		itemSpec(30119, wEnchGlovesMajorStrength),                   // Destroyer Gauntlets (T5)
		itemSpec(30032, 0, 24054, 24058),                            // Red Belt of Battle
		itemSpec(30121, wEnchLegsNethercobra, 24058),                // Destroyer Greaves (T5)
		itemSpec(28608, wEnchBootsCatsSwiftness, 28362, 28363),      // Ironstriders of Urgency
		itemSpec(28730, 0),                   // Mithril Band of the Unscarred
		itemSpec(30834, 0),                   // Shapeshifter's Signet
		itemSpec(21670, 0),                   // Badge of the Swarmguard
		itemSpec(28830, 0),                   // Dragonspine Trophy (owner-supplied)
		itemSpec(30082, wEnchWeaponMongoose), // Talon of Azshara (sword, 2.7)
		itemSpec(32052, wEnchWeaponMongoose), // Merciless Gladiator's Slicer (sword, 2.6)
		itemSpec(32054, 0),                   // Merciless Gladiator's War Edge
	}}
}

// The gear every run starts from. Swapped by the two test entry points below:
// the armory gear (Engineering goggles in the head slot, so 4/5 Destroyer), or
// the same gear wearing Destroyer Battle-Helm instead (5/5 Destroyer).
var (
	warriorBaseGear  = warriorCurrentGear
	warriorBaseLabel = "current armory gear + Dragonspine Trophy (4/5 Destroyer)"
)

// Koerry's gear with the goggles replaced by the tier 5 helm, giving a full 5/5
// Destroyer Battlegear. With five pieces the first tier swap is free - the 4pc
// survives it - which changes the whole upgrade order, so it gets its own run.
func warriorT5Gear() *proto.EquipmentSpec {
	g := warriorCurrentGear()
	g.Items[slotHead] = itemSpec(30120, wEnchHeadGlyphOfFerocity, bestGemsFor(30120)...)
	return g
}

func (c candidate) applyWarrior(e *proto.EquipmentSpec) {
	ench := c.enchant
	if ench == 0 {
		ench = warriorBaseGear().Items[c.slot].Enchant
	}
	e.Items[c.slot] = itemSpec(c.id, ench, bestGemsFor(c.id)...)
}

func warriorCandidates() []candidate {
	return []candidate{
		// ---- Head ----
		{slotHead, "Cursed Vision of Sargeras", 32235, nil, 0, ""},
		{slotHead, "Helm of the Illidari Shatterer", 32373, nil, 0, ""},
		{slotHead, "Surestrike Goggles v2.0", 32474, nil, 0, ""},
		{slotHead, "Fel-Steel Warhelm", 29983, nil, 0, ""},
		{slotHead, "Onslaught Battle-Helm (T6)", 30972, nil, 0, ""},
		{slotHead, "Destroyer Battle-Helm (T5)", 30120, nil, 0, ""},
		{slotHead, "Furious Gizmatic Goggles", 32461, nil, 0, ""},
		// ---- Neck ----
		{slotNeck, "Choker of Endless Nightmares", 32260, nil, 0, ""},
		{slotNeck, "Choker of Serrated Blades", 32591, nil, 0, ""},
		// ---- Shoulder ----
		{slotShoulder, "Blood-stained Pauldrons", 30866, nil, 0, ""},
		{slotShoulder, "Mantle of Darkness", 32377, nil, 0, ""},
		{slotShoulder, "Swiftsteel Shoulders", 32570, nil, 0, ""},
		{slotShoulder, "Onslaught Shoulderblades (T6)", 30979, nil, 0, ""},
		// ---- Back ----
		{slotBack, "Shadowmoon Destroyer's Drape", 32323, nil, 0, ""},
		{slotBack, "Black-Iron Battlecloak", 30729, nil, 0, ""},
		// ---- Chest ----
		{slotChest, "Heartshatter Breastplate", 32365, nil, 0, ""},
		{slotChest, "Midnight Chestguard", 30905, nil, 0, ""},
		{slotChest, "Krakken-Heart Breastplate", 30102, nil, 0, ""},
		{slotChest, "Onslaught Breastplate (T6)", 30975, nil, 0, ""},
		// ---- Wrist ----
		{slotWrist, "Deadly Cuffs", 30863, nil, 0, ""},
		{slotWrist, "Bindings of Lightning Reflexes", 32574, nil, 0, ""},
		{slotWrist, "Swiftsteel Bracers", 32568, nil, 0, ""},
		{slotWrist, "Bracers of Eradication", 30057, nil, 0, ""},
		// ---- Hands ----
		{slotHands, "Pillager's Gauntlets", 32608, nil, 0, ""},
		{slotHands, "Grips of Silent Justice", 32278, nil, 0, ""},
		{slotHands, "Grips of Damnation", 32347, nil, 0, ""},
		{slotHands, "Onslaught Gauntlets (T6)", 30969, nil, 0, ""},
		// ---- Waist ----
		{slotWaist, "Boneweave Girdle", 32346, nil, 0, ""},
		{slotWaist, "Belt of One-Hundred Deaths", 30106, nil, 0, ""},
		{slotWaist, "Belt of Seething Fury", 30915, nil, 0, ""},
		{slotWaist, "Shadow-walker's Cord", 32265, nil, 0, ""},
		// ---- Legs ----
		{slotLegs, "Leggings of Divine Retribution", 32341, nil, 0, ""},
		{slotLegs, "Legguards of Endless Rage", 30903, nil, 0, ""},
		{slotLegs, "Shady Dealer's Pantaloons", 30898, nil, 0, ""},
		{slotLegs, "Onslaught Greaves (T6)", 30977, nil, 0, ""},
		// ---- Feet ----
		{slotFeet, "Softstep Boots of Tracking", 32510, nil, 0, ""},
		{slotFeet, "Black Featherlight Boots", 30891, nil, 0, ""},
		{slotFeet, "Dreadboots of the Legion", 32345, nil, 0, ""},
		{slotFeet, "Warboots of Obliteration", 30081, nil, 0, ""},
		// ---- Finger1 (replaces Mithril Band of the Unscarred) ----
		{slotFinger1, "Stormrage Signet Ring", 32497, nil, 0, ""},
		{slotFinger1, "Band of the Ranger-General", 29997, nil, 0, ""},
		{slotFinger1, "Ring of Deceitful Intent", 32266, nil, 0, ""},
		{slotFinger1, "Unstoppable Aggressor's Ring", 32335, nil, 0, ""},
		// ---- Finger2 (replaces Shapeshifter's Signet) ----
		{slotFinger2, "Stormrage Signet Ring", 32497, nil, 0, ""},
		{slotFinger2, "Band of the Ranger-General", 29997, nil, 0, ""},
		{slotFinger2, "Ring of Deceitful Intent", 32266, nil, 0, ""},
		{slotFinger2, "Unstoppable Aggressor's Ring", 32335, nil, 0, ""},
		// ---- Trinket1 (replaces Badge of the Swarmguard) ----
		{slotTrinket1, "Madness of the Betrayer", 32505, nil, 0, ""},
		{slotTrinket1, "Bloodlust Brooch", 29383, nil, 0, ""},
		{slotTrinket1, "Tsunami Talisman", 30627, nil, 0, ""},
		{slotTrinket1, "Abacus of Violent Odds", 28288, nil, 0, ""},
		{slotTrinket1, "Drake Fang Talisman", 19406, nil, 0, ""},
		// ---- Trinket2 (currently Dragonspine Trophy) ----
		{slotTrinket2, "Madness of the Betrayer", 32505, nil, 0, ""},
		{slotTrinket2, "Bloodlust Brooch", 29383, nil, 0, ""},
		{slotTrinket2, "Tsunami Talisman", 30627, nil, 0, ""},
		{slotTrinket2, "Abacus of Violent Odds", 28288, nil, 0, ""},
		// ---- Main hand ----
		{slotMainHand, "Vengeful Gladiator's Slicer", 33762, nil, 0, "S3 sword 2.6"},
		{slotMainHand, "Vengeful Gladiator's Cleaver", 33669, nil, 0, "S3 axe 2.6"},
		{slotMainHand, "Rising Tide", 32236, nil, 0, "axe 2.6"},
		{slotMainHand, "Blade of Infamy", 30881, nil, 0, "sword 2.6"},
		{slotMainHand, "Wicked Edge of the Planes", 28433, nil, 0, "axe 2.7"},
		{slotMainHand, "Syphon of the Nathrezim", 32262, nil, 0, "mace 2.8"},
		{slotMainHand, "Warglaive of Azzinoth (legendary)", 32837, nil, 0, "sword 2.8"},
		// ---- Off hand ----
		{slotOffHand, "Vengeful Gladiator's Slicer", 33762, nil, 0, "S3 sword 2.6"},
		{slotOffHand, "Vengeful Gladiator's Cleaver", 33669, nil, 0, "S3 axe 2.6"},
		{slotOffHand, "Blade of Savagery", 32369, nil, 0, "sword 1.4"},
		{slotOffHand, "Tracker's Blade", 30865, nil, 0, "dagger 1.5"},
		{slotOffHand, "Vengeful Gladiator's Shanker", 33754, nil, 0, "dagger 1.8"},
		{slotOffHand, "Rising Tide", 32236, nil, 0, "axe 2.6"},
		{slotOffHand, "Warglaive of Azzinoth (legendary)", 32838, nil, 0, "sword 2.8"},
		// ---- Ranged ----
		{slotRanged, "Serpent Spine Longbow", 30105, nil, 0, ""},
		{slotRanged, "Steelhawk Crossbow", 28504, nil, 0, ""},
	}
}

func warriorEnchantAndGemFixes() []namedVariant {
	return []namedVariant{
		{"Boots: Cat's Swiftness -> Dexterity", "+6 agi for the run speed", func(e *proto.EquipmentSpec) {
			e.Items[slotFeet].Enchant = wEnchBootsDexterity
		}},
		{"Weapons: Mongoose -> Executioner (both)", "840 armor pen proc instead of 120 agi / 2% haste", func(e *proto.EquipmentSpec) {
			e.Items[slotMainHand].Enchant = wEnchWeaponExecutioner
			e.Items[slotOffHand].Enchant = wEnchWeaponExecutioner
		}},
		{"Weapons: Executioner MH only", "Executioner main hand, Mongoose off hand", func(e *proto.EquipmentSpec) {
			e.Items[slotMainHand].Enchant = wEnchWeaponExecutioner
		}},
		{"Regem everything with P3 epic gems", "colour-matched for socket bonuses", func(e *proto.EquipmentSpec) {
			regemAll(e, bestGemsFor)
		}},
		{"Regem everything, all red (ignore socket bonuses)", "every non-meta socket = Bold Crimson Spinel (+10 str)", func(e *proto.EquipmentSpec) {
			regemAll(e, allRedGemsFor)
		}},
		{"All gem fixes together", "colour-matched regem", func(e *proto.EquipmentSpec) {
			regemAll(e, bestGemsFor)
		}},
	}
}

// ---------------------------------------------------------------------------
// Sim plumbing
// ---------------------------------------------------------------------------

const (
	// The repo's standard Fury build (sim/warrior/dps/dps_test.go).
	talentsFuryStandard = "3500501130201-05050005505012050115"
	// Standard Fury with 5 points moved from Improved Berserker Stance into
	// Sword Specialization - Koerry dual-wields two swords. Same string used by
	// the "Sword Spec Fury" variant in results_glaives.md.
	talentsFurySwordSpec = "350050113020105-0505000550501205011"
)

// warriorBonusHit is flat hit rating added to every player built by warriorPlayer,
// standing in for the regemming a hit-short gear set would get (see the T6R
// baseline in warriorhyjalbtsr_test.go). Zero for all the original reports.
var warriorBonusHit float64

func warriorPlayer(gear *proto.EquipmentSpec, talents string, rotationFile string) *proto.Player {
	player := &proto.Player{
		Class:     proto.Class_ClassWarrior,
		Race:      proto.Race_RaceOrc,
		Equipment: gear,
		Consumables: &proto.ConsumesSpec{
			PotId:       22838,
			FlaskId:     22854,
			FoodId:      27658,
			ConjuredId:  22788,
			ExplosiveId: 30217,
			SuperSapper: true,
			OhImbueId:   29453, // Adamantite Sharpening Stone, as worn
			ScrollAgi:   true,
			ScrollStr:   true,
		},
		TalentsString:      talents,
		Profession1:        proto.Profession_Engineering, // required to wear the goggles
		Profession2:        proto.Profession_Blacksmithing,
		DistanceFromTarget: 25,
		Rotation:           core.GetAplRotation("../ui/specs/warrior/dps/apls", rotationFile).Rotation,
		Buffs:              core.FullIndividualBuffs,
	}
	if warriorBonusHit != 0 {
		player.BonusStats = &proto.UnitStats{Stats: stats.Stats{stats.MeleeHitRating: warriorBonusHit}.ToProtoArray()}
	}
	core.WithSpec(player, &proto.Player_DpsWarrior{
		DpsWarrior: &proto.DpsWarrior{
			Options: &proto.DpsWarrior_Options{
				ClassOptions: &proto.WarriorOptions{
					DefaultShout:  proto.WarriorShout_WarriorShoutBattle,
					DefaultStance: proto.WarriorStance_WarriorStanceBerserker,
				},
			},
		},
	})
	return player
}

func runWarrior(t *testing.T, gear *proto.EquipmentSpec, talents, rotationFile string, encounter *proto.Encounter, iterations int32) simResult {
	raid := core.SinglePlayerRaidProto(warriorPlayer(gear, talents, rotationFile), core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs)

	cs := core.ComputeStats(&proto.ComputeStatsRequest{Raid: raid, Encounter: encounter})
	final := cs.RaidStats.Parties[0].Players[0].FinalStats.Stats

	result := core.RunRaidSimConcurrent(&proto.RaidSimRequest{
		Raid:       raid,
		Encounter:  encounter,
		SimOptions: &proto.SimOptions{Iterations: iterations, RandomSeed: 101},
	})
	if result.Error != nil {
		t.Fatalf("sim failed: %s", result.Error.Message)
	}
	return simResult{
		dps:   result.RaidMetrics.Dps.Avg,
		stdev: result.RaidMetrics.Dps.Stdev,
		hit:   final[stats.MeleeHitRating],
		crit:  final[stats.MeleeCritRating],
		ap:    final[stats.AttackPower],
		haste: final[stats.MeleeHasteRating],
		exp:   final[stats.ExpertiseRating],
	}
}

func warriorIterationCount() int32 {
	if v := os.Getenv("WARRIOR_ITERATIONS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return int32(n)
		}
	}
	return 20000
}

func TestWarriorP3Upgrades(t *testing.T) {
	RegisterAll()
	warriorBaseGear, warriorBaseLabel = warriorCurrentGear, "current armory gear + Dragonspine Trophy (4/5 Destroyer)"
	warriorUpgradeReport(t)
}

// Same benchmark from a 5/5 Destroyer Battlegear baseline: the Engineering
// goggles are replaced by Destroyer Battle-Helm, so the tier 4pc can survive one
// swap. Everything else is identical.
func TestWarriorP3UpgradesT5(t *testing.T) {
	RegisterAll()
	warriorBaseGear, warriorBaseLabel = warriorT5Gear, "armory gear + Dragonspine Trophy, 5/5 Destroyer (T5 helm)"
	warriorUpgradeReport(t)
}

func warriorUpgradeReport(t *testing.T) {
	iterations := warriorIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)

	run := func(gear *proto.EquipmentSpec) simResult {
		return runWarrior(t, gear, talentsFuryStandard, "fury", encounter, iterations)
	}

	base := run(warriorBaseGear())
	fmt.Printf("ITERATIONS %d\n", iterations)
	fmt.Printf("BASELINE | %s | DPS %.1f (stdev %.1f) | AP %.0f hit %.0f crit %.0f haste %.0f exp %.0f\n",
		warriorBaseLabel, base.dps, base.stdev, base.ap, base.hit, base.crit, base.haste, base.exp)

	// --- Talents ---------------------------------------------------------------
	fmt.Println("\n=== TALENTS / SPEC (same gear) ===")
	for _, tc := range []struct{ label, talents, rotation string }{
		{"Fury, repo standard [baseline]", talentsFuryStandard, "fury"},
		{"Fury + Sword Specialization", talentsFurySwordSpec, "fury"},
	} {
		r := runWarrior(t, warriorBaseGear(), tc.talents, tc.rotation, encounter, iterations)
		fmt.Printf("TALENT | %-40s | %8.1f | %+8.1f\n", tc.label, r.dps, r.dps-base.dps)
	}

	// --- Enchant / gem fixes ---------------------------------------------------
	fmt.Println("\n=== ENCHANT / GEM FIXES (no drops needed) ===")
	for _, v := range warriorEnchantAndGemFixes() {
		gear := cloneGear(warriorBaseGear())
		v.apply(gear)
		r := run(gear)
		fmt.Printf("FIX | %-52s | %8.1f | %+7.1f | %s\n", v.label, r.dps, r.dps-base.dps, v.note)
	}

	// --- Single-slot item swaps ------------------------------------------------
	fmt.Println("\n=== SINGLE-SLOT ITEM SWAPS (vs current gear) ===")
	type scored struct {
		c candidate
		r simResult
	}
	all := []scored{}
	for _, c := range warriorCandidates() {
		if c.id == warriorBaseGear().Items[c.slot].Id {
			continue // already worn in this baseline
		}
		gear := cloneGear(warriorBaseGear())
		c.applyWarrior(gear)
		r := run(gear)
		all = append(all, scored{c, r})
		fmt.Printf("SWAP | %-9s | %-34s | %6d | %8.1f | %+7.1f | AP %5.0f hit %4.0f crit %4.0f haste %4.0f exp %3.0f | %s\n",
			warriorSlotNames[c.slot], c.label, c.id, r.dps, r.dps-base.dps, r.ap, r.hit, r.crit, r.haste, r.exp, c.note)
	}

	// --- Cumulative, greedy, no duplicate item ids across slots -----------------
	sortedAll := append([]scored(nil), all...)
	sort.Slice(sortedAll, func(i, j int) bool { return sortedAll[i].r.dps > sortedAll[j].r.dps })

	pick := func(keep func(candidate) bool) []scored {
		assigned := map[int]bool{}
		usedIDs := map[int32]bool{}
		ordered := []scored{}
		for _, s := range sortedAll {
			if s.r.dps <= base.dps || assigned[s.c.slot] || (keep != nil && !keep(s.c)) {
				continue
			}
			isWeapon := s.c.slot == slotMainHand || s.c.slot == slotOffHand
			if !isWeapon && usedIDs[s.c.id] {
				continue
			}
			assigned[s.c.slot] = true
			usedIDs[s.c.id] = true
			ordered = append(ordered, s)
		}
		return ordered
	}

	notGlaive := func(c candidate) bool { return c.id != 32837 && c.id != 32838 }
	ordered := pick(notGlaive)

	fmt.Println("\n=== CUMULATIVE (best non-legendary item per slot, single-slot rank order) ===")
	gear := cloneGear(warriorBaseGear())
	prev := base.dps
	for _, s := range ordered {
		s.c.applyWarrior(gear)
		r := run(cloneGear(gear))
		fmt.Printf("CUMUL | +%-9s %-34s | %8.1f | step %+7.1f | total %+8.1f\n",
			warriorSlotNames[s.c.slot], s.c.label, r.dps, r.dps-prev, r.dps-base.dps)
		prev = r.dps
	}

	// --- Assembled builds ------------------------------------------------------
	// Which slots currently hold a Destroyer (T5) piece depends on the baseline:
	// 4 with the Engineering goggles in the head slot, 5 with the tier helm. That
	// count decides how many tier swaps are "free" - the 4pc survives down to four
	// pieces - so it is computed rather than hardcoded.
	destroyerIDs := map[int32]bool{30120: true, 30122: true, 30118: true, 30119: true, 30121: true}
	tierSlots := map[int]bool{}
	wornTier := 0
	for i, it := range warriorBaseGear().Items {
		if destroyerIDs[it.Id] {
			tierSlots[i] = true
			wornTier++
		}
	}
	freeTierSwaps := wornTier - 4
	if freeTierSwaps < 0 {
		freeTierSwaps = 0
	}
	fmt.Printf("TIERINFO | Destroyer pieces worn: %d | free tier swaps before 4pc breaks: %d\n", wornTier, freeTierSwaps)

	buildFrom := func(list []scored, skip func(int) bool) *proto.EquipmentSpec {
		g := cloneGear(warriorBaseGear())
		for _, s := range list {
			if skip != nil && skip(s.c.slot) {
				continue
			}
			s.c.applyWarrior(g)
		}
		return g
	}

	// Best candidate per slot ignoring whether the single-slot delta was positive -
	// used for the tierless build, where the tier slots individually look bad only
	// because of the 4pc they are breaking.
	bestPerSlot := func(keep func(candidate) bool) []scored {
		best := map[int]scored{}
		for _, s := range all {
			if keep != nil && !keep(s.c) {
				continue
			}
			// Never downgrade a slot. The exception is the tier slots, whose
			// single-slot deltas are negative only because they break the 4pc -
			// those are exactly what a tierless build wants to override.
			if s.r.dps <= base.dps && !tierSlots[s.c.slot] {
				continue
			}
			if cur, ok := best[s.c.slot]; !ok || s.r.dps > cur.r.dps {
				best[s.c.slot] = s
			}
		}
		usedIDs := map[int32]bool{}
		out := []scored{}
		for _, slot := range []int{slotTrinket1, slotTrinket2, slotFinger1, slotFinger2, slotHead, slotNeck,
			slotShoulder, slotBack, slotChest, slotWrist, slotHands, slotWaist, slotLegs, slotFeet,
			slotMainHand, slotOffHand, slotRanged} {
			s, ok := best[slot]
			if !ok {
				continue
			}
			isWeapon := slot == slotMainHand || slot == slotOffHand
			if !isWeapon && usedIDs[s.c.id] {
				continue
			}
			usedIDs[s.c.id] = true
			out = append(out, s)
		}
		return out
	}

	// Keep 4pc Destroyer: swap every non-tier slot, plus however many tier slots
	// the set can spare (one when 5/5 is worn, none when 4/5 is).
	tierBudget := map[int]bool{}
	{
		tierRanked := []scored{}
		for _, s := range sortedAll {
			if tierSlots[s.c.slot] && !tierBudget[s.c.slot] && len(tierRanked) < freeTierSwaps {
				tierBudget[s.c.slot] = true
				tierRanked = append(tierRanked, s)
			}
		}
	}
	keepT5 := buildFrom(ordered, func(slot int) bool { return tierSlots[slot] && !tierBudget[slot] })

	// Tierless: best item everywhere, tier bonuses be damned.
	tierless := buildFrom(bestPerSlot(notGlaive), nil)

	// 4pc Onslaught (T6): force the four T6 pieces in, best non-tier elsewhere.
	onslaught := []candidate{
		{slotShoulder, "Onslaught Shoulderblades (T6)", 30979, nil, 0, ""},
		{slotChest, "Onslaught Breastplate (T6)", 30975, nil, 0, ""},
		{slotLegs, "Onslaught Greaves (T6)", 30977, nil, 0, ""},
		{slotHands, "Onslaught Gauntlets (T6)", 30969, nil, 0, ""},
	}
	onslaughtSlots := map[int]bool{slotShoulder: true, slotChest: true, slotLegs: true, slotHands: true}
	t6 := buildFrom(bestPerSlot(notGlaive), func(slot int) bool { return onslaughtSlots[slot] })
	for _, c := range onslaught {
		c.applyWarrior(t6)
	}

	glaives := buildFrom(bestPerSlot(nil), nil)

	withExecutioner := func(g *proto.EquipmentSpec) *proto.EquipmentSpec {
		g.Items[slotMainHand].Enchant = wEnchWeaponExecutioner
		g.Items[slotOffHand].Enchant = wEnchWeaponExecutioner
		return g
	}

	fmt.Println("\n=== ASSEMBLED BUILDS ===")
	for _, b := range []struct {
		label string
		gear  *proto.EquipmentSpec
	}{
		{"A. stacked single-slot winners", buildFrom(ordered, nil)},
		{"B. best items, keeping 4pc Destroyer (T5)", keepT5},
		{"C. tierless: best item every slot", tierless},
		{"C+. tierless + Executioner on both weapons", withExecutioner(buildFrom(bestPerSlot(notGlaive), nil))},
		{"D. best items with 4pc Onslaught (T6)", t6},
		{"D+. 4pc Onslaught + Executioner", withExecutioner(func() *proto.EquipmentSpec {
			g := buildFrom(bestPerSlot(notGlaive), func(slot int) bool { return onslaughtSlots[slot] })
			for _, c := range onslaught {
				c.applyWarrior(g)
			}
			return g
		}())},
		{"E. tierless, Warglaives allowed (sim picks)", glaives},
		{"F. Warglaives allowed, keeping 4pc Destroyer (T5)", buildFrom(bestPerSlot(nil), func(slot int) bool { return tierSlots[slot] && !tierBudget[slot] })},
		{"G. tierless, both Warglaives forced", func() *proto.EquipmentSpec {
			g := buildFrom(bestPerSlot(notGlaive), nil)
			g.Items[slotMainHand] = itemSpec(32837, wEnchWeaponMongoose)
			g.Items[slotOffHand] = itemSpec(32838, wEnchWeaponMongoose)
			return g
		}()},
	} {
		r := runWarrior(t, cloneGear(b.gear), talentsFuryStandard, "fury", encounter, iterations)
		fmt.Printf("BUILD | %-48s | %8.1f | %+8.1f | AP %5.0f hit %4.0f crit %4.0f haste %4.0f exp %3.0f\n",
			b.label, r.dps, r.dps-base.dps, r.ap, r.hit, r.crit, r.haste, r.exp)
		for i, it := range b.gear.Items {
			fmt.Printf("BUILDITEM | %-48s | %-9s | %6d | ench %5d | gems %v\n", b.label, warriorSlotNames[i], it.Id, it.Enchant, it.Gems)
		}
	}

	// --- Reference builds ------------------------------------------------------
	fmt.Println("\n=== REFERENCE BUILDS ===")
	for _, ref := range []struct{ label, file string }{
		{"repo P2 fury preset", "p2_fury"},
		{"repo P3 fury preset (Warglaives)", "p3_fury"},
		{"repo P4 fury preset (Warglaives)", "p4_fury"},
	} {
		r := runWarrior(t, core.GetGearSet("../ui/specs/warrior/dps/gear_sets", ref.file).GearSet, talentsFuryStandard, "fury", encounter, iterations)
		fmt.Printf("REF | %-34s | %8.1f | %+8.1f | AP %5.0f hit %4.0f crit %4.0f haste %4.0f exp %3.0f\n",
			ref.label, r.dps, r.dps-base.dps, r.ap, r.hit, r.crit, r.haste, r.exp)
	}
}
