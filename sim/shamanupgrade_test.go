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

// Benchmark: Phase 3 upgrade finder for the Enhancement Shaman
// "Koerrii" (Spineshatter-EU), whose current equipment is transcribed below
// from https://ironforge.pro/anniversary/player/Spineshatter/Koerrii/ .
//
// Every candidate is a single-slot swap against the unmodified armory gear, so
// each printed delta answers "what does replacing exactly this one item buy me?".
// A cumulative pass then stacks the winners in order to expose interactions
// (notably the 4pc Cataclysm Harness / T5 bonus, which survives one swap out of
// the current five tier pieces but not two).
//
// Kept in the tree deliberately (results written up in result_shaman.md) so the
// benchmark can be re-run under different assumptions: set SHAMAN_ITERATIONS to
// trade accuracy for speed, e.g.
//
//	SHAMAN_ITERATIONS=300 go test -tags with_db ./sim/ -run TestShamanP3Upgrades -v

// Equipment slot indices in proto.EquipmentSpec.Items.
const (
	slotHead = iota
	slotNeck
	slotShoulder
	slotBack
	slotChest
	slotWrist
	slotHands
	slotWaist
	slotLegs
	slotFeet
	slotFinger1
	slotFinger2
	slotTrinket1
	slotTrinket2
	slotMainHand
	slotOffHand
	slotRanged
)

var slotNames = [17]string{"Head", "Neck", "Shoulder", "Back", "Chest", "Wrist", "Hands",
	"Waist", "Legs", "Feet", "Finger1", "Finger2", "Trinket1", "Trinket2", "MainHand", "OffHand", "Ranged"}

// Enchant ids as used by ui/*/gear_sets/*.gear.json (enchant itemId, or effectId
// when the enchant has no item).
// NOTE: the engine resolves ItemSpec.Enchant through core.EnchantsByEffectID and
// silently ignores anything it cannot find (the panic in database.go is commented
// out). Several checked-in gear presets - including ui/shaman/enhancement/gear_sets
// p2/p3/p4 - store the enchant's *item* id instead (29192, 28888, 29535, 28279,
// 24003, 22559), which therefore applies no enchant at all. Everything here uses
// effect ids so the enchants actually land.
const (
	enchHeadGlyphOfFerocity   = 3003 // +34 AP, +16 hit        (item 29192)
	enchShoulderInscrVeng     = 2983 // +26 AP                 (item 28885)
	enchShoulderGrInscrVeng   = 2986 // +30 AP, +10 crit       (item 28888)
	enchCloakGreaterAgility   = 368  // +12 agi
	enchChestExceptionalStats = 2661 // +6 all stats           (item 24003)
	enchBracerBrawn           = 2647 // +12 str
	enchGlovesMajorStrength   = 684  // +15 str
	enchLegsCobrahide         = 3010 // +40 AP, +10 crit       (item 29533)
	enchLegsNethercobra       = 3012 // +50 AP, +12 crit       (item 29535)
	enchBootsCatsSwiftness    = 2939 // +6 agi, minor speed    (item 28279)
	enchBootsDexterity        = 2657 // +12 agi
	enchWeaponMongoose        = 2673 // proc                   (item 22559)
)

// enchantItemIDToEffectID translates the enchant item ids used by the checked-in
// enhancement gear presets into the effect ids the engine actually resolves.
var enchantItemIDToEffectID = map[int32]int32{
	22559: 2673, // Mongoose
	24003: 2661, // Exceptional Stats
	28279: 2939, // Cat's Swiftness
	28885: 2983, // Inscription of Vengeance
	28888: 2986, // Greater Inscription of Vengeance
	29192: 3003, // Glyph of Ferocity
	29533: 3010, // Cobrahide Leg Armor
	29535: 3012, // Nethercobra Leg Armor
	33307: 3225, // Executioner
}

func fixPresetEnchants(e *proto.EquipmentSpec) *proto.EquipmentSpec {
	for _, it := range e.Items {
		if fixed, ok := enchantItemIDToEffectID[it.Enchant]; ok {
			it.Enchant = fixed
		}
	}
	return e
}

// P3 epic gems, one per socket colour.
const (
	gemMetaRelentlessEarthstorm = 32409 // meta: +12 agi, 3% crit damage
	gemBoldCrimsonSpinel        = 32193 // red: +10 str
	gemSovereignShadowsong      = 32211 // purple (fits blue): +5 str, +7 sta
	gemInscribedPyrestone       = 32217 // orange (fits yellow): +5 str, +5 crit
)

func itemSpec(id, enchant int32, gems ...int32) *proto.ItemSpec {
	return &proto.ItemSpec{Id: id, Enchant: enchant, Gems: gems}
}

// bestGemsFor fills an item's actual sockets with the best P3 epic gem of the
// matching colour, so every socket bonus is satisfied and no phantom gems are
// added. (NewItem in database.go applies every gem it is handed, even beyond the
// item's socket count, so gem lists must be socket-accurate.)
func bestGemsFor(itemID int32) []int32 {
	item := core.GetItemByID(itemID)
	if item == nil || len(item.GemSockets) == 0 {
		return nil
	}
	gems := make([]int32, 0, len(item.GemSockets))
	for _, color := range item.GemSockets {
		switch color {
		case proto.GemColor_GemColorMeta:
			gems = append(gems, gemMetaRelentlessEarthstorm)
		case proto.GemColor_GemColorBlue:
			gems = append(gems, gemSovereignShadowsong)
		case proto.GemColor_GemColorYellow:
			gems = append(gems, gemInscribedPyrestone)
		default:
			gems = append(gems, gemBoldCrimsonSpinel)
		}
	}
	return gems
}

// allRedGemsFor ignores socket bonuses and stuffs every non-meta socket with the
// pure-Strength red gem.
func allRedGemsFor(itemID int32) []int32 {
	item := core.GetItemByID(itemID)
	if item == nil || len(item.GemSockets) == 0 {
		return nil
	}
	gems := make([]int32, 0, len(item.GemSockets))
	for _, color := range item.GemSockets {
		if color == proto.GemColor_GemColorMeta {
			gems = append(gems, gemMetaRelentlessEarthstorm)
		} else {
			gems = append(gems, gemBoldCrimsonSpinel)
		}
	}
	return gems
}

// currentGear is Koerrii's armory equipment, verbatim. The boots' "Minor Speed
// Increase" enchant carries no stats and has no entry in this repo's enchant db,
// so the feet slot is modelled as unenchanted.
func currentGear() *proto.EquipmentSpec {
	return &proto.EquipmentSpec{Items: []*proto.ItemSpec{
		itemSpec(30190, enchHeadGlyphOfFerocity, 32409, 24058), // Cataclysm Helm (T5)
		itemSpec(28509, 0), // Worgen Claw Necklace
		itemSpec(30194, enchShoulderInscrVeng, 24054, 24058),            // Cataclysm Shoulderplates (T5)
		itemSpec(30098, enchCloakGreaterAgility),                        // Razor-Scale Battlecloak
		itemSpec(30185, enchChestExceptionalStats, 24027, 24054, 24058), // Cataclysm Chestplate (T5)
		itemSpec(32816, enchBracerBrawn, 33782),                         // Veteran's Linked Bracers
		itemSpec(30189, enchGlovesMajorStrength),                        // Cataclysm Gauntlets (T5)
		itemSpec(30068, 0),                                              // Girdle of the Tidal Call
		itemSpec(30192, enchLegsCobrahide, 24027),                       // Cataclysm Legplates (T5)
		itemSpec(30039, 0),                                              // Boots of Utter Darkness (Minor Speed)
		itemSpec(30061, 0),                                              // Ancestral Ring of Conquest
		itemSpec(29282, 0),                                              // Violet Signet
		itemSpec(25937, 0),                                              // Terokkar Tablet of Precision
		itemSpec(28041, 0),                                              // Bladefist's Breadth
		itemSpec(31965, enchWeaponMongoose),                             // Merciless Gladiator's Cleaver
		itemSpec(31965, enchWeaponMongoose),                             // Merciless Gladiator's Cleaver
		itemSpec(24413, 0),                                              // Totem of the Thunderhead
	}}
}

func cloneGear(e *proto.EquipmentSpec) *proto.EquipmentSpec {
	out := &proto.EquipmentSpec{Items: make([]*proto.ItemSpec, len(e.Items))}
	for i, it := range e.Items {
		gems := make([]int32, len(it.Gems))
		copy(gems, it.Gems)
		out.Items[i] = &proto.ItemSpec{Id: it.Id, Enchant: it.Enchant, Gems: gems}
	}
	return out
}

// candidate is a single-slot replacement. Gems are the best P3 epic gems that
// match the replacement item's own sockets; the enchant is inherited from the
// slot's current enchant unless overridden.
type candidate struct {
	slot    int
	label   string
	id      int32
	gems    []int32 // unused: apply() always gems from the item's own sockets
	enchant int32   // 0 => inherit the baseline slot enchant
	note    string
}

func (c candidate) apply(e *proto.EquipmentSpec) {
	ench := c.enchant
	if ench == 0 {
		ench = currentGear().Items[c.slot].Enchant
	}
	// Gems always come from the replacement item's own sockets, never from the
	// hand-written list, so a socketless item can never inherit phantom gems.
	e.Items[c.slot] = itemSpec(c.id, ench, bestGemsFor(c.id)...)
}

func itemCandidates() []candidate {
	return []candidate{
		// ---- Head ----
		{slotHead, "Cursed Vision of Sargeras", 32235, nil, 0, "BT - Illidari Council"},
		{slotHead, "Forest Prowler's Helm", 32376, nil, 0, "Hyjal - Azgalor"},
		{slotHead, "Skyshatter Cover (T6)", 31015, nil, 0, "T6 head"},
		// ---- Neck ----
		{slotNeck, "Choker of Endless Nightmares", 32260, nil, 0, "BT - Reliquary of Souls"},
		{slotNeck, "Telonicus's Pendant of Mayhem", 30017, nil, 0, "TK - Void Reaver (P2)"},
		// ---- Shoulder ----
		{slotShoulder, "Mantle of Darkness", 32377, nil, 0, "BT - Teron Gorefiend"},
		{slotShoulder, "Shoulders of the Hidden Predator", 32264, nil, 0, "BT - Gurtogg"},
		{slotShoulder, "Shoulders of Lightning Reflexes", 32575, nil, 0, "BT trash"},
		{slotShoulder, "Swiftstrike Shoulders", 32581, nil, 0, "BT trash (crafted)"},
		{slotShoulder, "Skyshatter Pauldrons (T6)", 31024, nil, 0, "T6 shoulder"},
		// ---- Back ----
		{slotBack, "Shadowmoon Destroyer's Drape", 32323, nil, 0, "BT - Supremus"},
		// ---- Chest ----
		{slotChest, "Midnight Chestguard", 30905, nil, 0, "Hyjal - Archimonde"},
		{slotChest, "Mail of Fevered Pursuit", 30907, nil, 0, "Hyjal - Archimonde"},
		{slotChest, "Nether Shadow Tunic", 32252, nil, 0, "BT trash"},
		{slotChest, "Skyshatter Tunic (T6)", 31018, nil, 0, "T6 chest"},
		// ---- Wrist ----
		{slotWrist, "Deadly Cuffs", 30863, nil, 0, "Hyjal trash"},
		{slotWrist, "Insidious Bands", 32324, nil, 0, "BT - Supremus"},
		{slotWrist, "Bindings of Lightning Reflexes", 32574, nil, 0, "BT trash"},
		{slotWrist, "Swiftstrike Bracers", 32580, nil, 0, "BT trash (crafted)"},
		// ---- Hands ----
		{slotHands, "Grips of Damnation", 32347, nil, 0, "BT - Mother Shahraz"},
		{slotHands, "Fists of Mukoa", 32234, nil, 0, "BT trash"},
		{slotHands, "Skyshatter Grips (T6)", 31011, nil, 0, "T6 hands"},
		// ---- Waist ----
		{slotWaist, "Boneweave Girdle", 32346, nil, 0, "BT - Bloodboil"},
		{slotWaist, "Shadow-walker's Cord", 32265, nil, 0, "BT - Najentus"},
		{slotWaist, "Valestalker Girdle", 30919, nil, 0, "Hyjal - Anetheron"},
		{slotWaist, "Belt of One-Hundred Deaths", 30106, nil, 0, "SSC - Vashj (P2)"},
		// ---- Legs ----
		{slotLegs, "Bow-stitched Leggings", 30900, nil, 0, "Hyjal trash"},
		{slotLegs, "Shady Dealer's Pantaloons", 30898, nil, 0, "Hyjal trash"},
		{slotLegs, "Skyshatter Pants (T6)", 31021, nil, 0, "T6 legs"},
		{slotLegs, "Leggings of Murderous Intent", 29995, nil, 0, "TK - Kael'thas (P2)"},
		// ---- Feet ----
		{slotFeet, "Softstep Boots of Tracking", 32510, nil, 0, "BT - Gurtogg"},
		{slotFeet, "Black Featherlight Boots", 30891, nil, 0, "Hyjal trash"},
		{slotFeet, "Shadowmaster's Boots", 32366, nil, 0, "BT - Illidari Council"},
		// ---- Finger1 (replaces Ancestral Ring of Conquest) ----
		{slotFinger1, "Stormrage Signet Ring", 32497, nil, 0, "BT - Illidan"},
		{slotFinger1, "Ring of Deceitful Intent", 32266, nil, 0, "BT - Najentus"},
		{slotFinger1, "Band of the Ranger-General", 29997, nil, 0, "TK - Kael'thas (P2)"},
		// ---- Finger2 (replaces Violet Signet) ----
		{slotFinger2, "Stormrage Signet Ring", 32497, nil, 0, "BT - Illidan"},
		{slotFinger2, "Ring of Deceitful Intent", 32266, nil, 0, "BT - Najentus"},
		{slotFinger2, "Violet Signet of the Master Assassin", 29283, nil, 0, "Kara ring, exalted upgrade"},
		// ---- Trinket1 (replaces Terokkar Tablet of Precision) ----
		{slotTrinket1, "Madness of the Betrayer", 32505, nil, 0, "BT - Illidan"},
		{slotTrinket1, "Dragonspine Trophy", 28830, nil, 0, "Gruul (P1)"},
		{slotTrinket1, "Bloodlust Brooch", 29383, nil, 0, "badge / crafted (P1)"},
		{slotTrinket1, "Tsunami Talisman", 30627, nil, 0, "TK - Kael'thas (P2)"},
		{slotTrinket1, "Abacus of Violent Odds", 28288, nil, 0, "heroic Mech (P1)"},
		// ---- Trinket2 (replaces Bladefist's Breadth) ----
		{slotTrinket2, "Madness of the Betrayer", 32505, nil, 0, "BT - Illidan"},
		{slotTrinket2, "Dragonspine Trophy", 28830, nil, 0, "Gruul (P1)"},
		{slotTrinket2, "Bloodlust Brooch", 29383, nil, 0, "badge / crafted (P1)"},
		{slotTrinket2, "Tsunami Talisman", 30627, nil, 0, "TK - Kael'thas (P2)"},
		// ---- Main hand ----
		{slotMainHand, "Syphon of the Nathrezim", 32262, nil, 0, "BT - Gurtogg (2.8 mace)"},
		{slotMainHand, "Rising Tide", 32236, nil, 0, "BT - Najentus (2.6 axe)"},
		{slotMainHand, "Vengeful Gladiator's Cleaver", 33669, nil, 0, "S3 arena (2.6 axe)"},
		{slotMainHand, "The Brutalizer", 32254, nil, 0, "ZA - Zul'jin (1.6 axe)"},
		{slotMainHand, "Swiftsteel Bludgeon", 32943, nil, 0, "crafted (1.5 mace)"},
		// ---- Off hand ----
		{slotOffHand, "Syphon of the Nathrezim", 32262, nil, 0, "BT - Gurtogg (2.8 mace)"},
		{slotOffHand, "Rising Tide", 32236, nil, 0, "BT - Najentus (2.6 axe)"},
		{slotOffHand, "Vengeful Gladiator's Cleaver", 33669, nil, 0, "S3 arena (2.6 axe)"},
		{slotOffHand, "The Brutalizer", 32254, nil, 0, "ZA - Zul'jin (1.6 axe)"},
		// ---- Ranged (totem) ----
		{slotRanged, "Totem of the Astral Winds", 27815, nil, 0, "quest (P1) - +80 AP to Windfury"},
		{slotRanged, "Stonebreaker's Totem", 33507, nil, 0, "badge (P4)"},
	}
}

// enchantAndGemFixes are upgrades that need no drops at all.
type namedVariant struct {
	label string
	note  string
	apply func(e *proto.EquipmentSpec)
}

// regemAll re-gems every equipped item using the supplied socket-aware gemmer.
func regemAll(e *proto.EquipmentSpec, gemmer func(int32) []int32) {
	for _, it := range e.Items {
		it.Gems = gemmer(it.Id)
	}
}

// applyEnchantFixes upgrades the three enchants that are cheap wins regardless of
// which item currently occupies the slot.
func applyEnchantFixes(e *proto.EquipmentSpec) {
	e.Items[slotFeet].Enchant = enchBootsCatsSwiftness
	e.Items[slotShoulder].Enchant = enchShoulderGrInscrVeng
	e.Items[slotLegs].Enchant = enchLegsNethercobra
}

func enchantAndGemFixes() []namedVariant {
	return []namedVariant{
		{"Boots: none -> Cat's Swiftness", "+6 agi (currently Minor Speed, 0 stats)", func(e *proto.EquipmentSpec) {
			e.Items[slotFeet].Enchant = enchBootsCatsSwiftness
		}},
		{"Boots: none -> Dexterity", "+12 agi (no run speed)", func(e *proto.EquipmentSpec) {
			e.Items[slotFeet].Enchant = enchBootsDexterity
		}},
		{"Shoulder: Inscription -> Greater Inscription of Vengeance", "+4 AP, +10 crit (Aldor/Scryer exalted)", func(e *proto.EquipmentSpec) {
			e.Items[slotShoulder].Enchant = enchShoulderGrInscrVeng
		}},
		{"Legs: Cobrahide -> Nethercobra Leg Armor", "+10 AP, +2 crit", func(e *proto.EquipmentSpec) {
			e.Items[slotLegs].Enchant = enchLegsNethercobra
		}},
		{"Wrist: swap PvP gem for Bold Crimson Spinel", "Steady Talasite (sta/resil) -> +10 str", func(e *proto.EquipmentSpec) {
			e.Items[slotWrist].Gems = []int32{gemBoldCrimsonSpinel}
		}},
		{"Regem everything with P3 epic gems", "rare -> epic gems, colour-matched for socket bonuses", func(e *proto.EquipmentSpec) {
			regemAll(e, bestGemsFor)
		}},
		{"Regem everything, all red (ignore socket bonuses)", "every non-meta socket = Bold Crimson Spinel (+10 str)", func(e *proto.EquipmentSpec) {
			regemAll(e, allRedGemsFor)
		}},
		{"All enchant + gem fixes together", "Cat's Swiftness + Greater Inscription + Nethercobra + colour-matched regem", func(e *proto.EquipmentSpec) {
			applyEnchantFixes(e)
			regemAll(e, bestGemsFor)
		}},
	}
}

// ---------------------------------------------------------------------------
// Sim plumbing
// ---------------------------------------------------------------------------

// Koerrii's armory "primary" spec is 17/44/0; the second spec, 0/47/14, is the
// standard PvE Enhancement build and matches this repo's own default. All
// upgrade runs use the PvE build.
const (
	talentsKoerriiPvE     = "-502520210502133531151-50005301"
	talentsKoerriiPrimary = "250031501-500520210501133531151"
	talentsRepoDefault    = "03-500502210501133531151-50005301"
)

func shamanPlayer(gear *proto.EquipmentSpec, talents string) *proto.Player {
	player := &proto.Player{
		Class:     proto.Class_ClassShaman,
		Race:      proto.Race_RaceOrc,
		Equipment: gear,
		Consumables: &proto.ConsumesSpec{
			PotId: 22838, FlaskId: 22854, FoodId: 27658, ConjuredId: 22788, ExplosiveId: 30217,
			SuperSapper: true, GoblinSapper: true, ScrollAgi: true, ScrollStr: true,
		},
		TalentsString: talents,
		Rotation:      core.GetAplRotation("../ui/specs/shaman/enhancement/apls", "default").Rotation,
		Buffs:         core.FullIndividualBuffs,
	}
	core.WithSpec(player, &proto.Player_EnhancementShaman{
		EnhancementShaman: &proto.EnhancementShaman{
			Options: &proto.EnhancementShaman_Options{
				SyncType:    proto.ShamanSyncType_Auto,
				ImbueOh:     proto.ShamanImbue_WindfuryWeapon,
				ImbueOhSwap: proto.ShamanImbue_WindfuryWeapon,
				ClassOptions: &proto.ShamanOptions{
					ImbueMh:        proto.ShamanImbue_WindfuryWeapon,
					ImbueMhSwap:    proto.ShamanImbue_WindfuryWeapon,
					ShieldProcrate: 0.0,
				},
			},
		},
	})
	return player
}

type simResult struct {
	dps   float64
	stdev float64
	hit   float64
	crit  float64
	ap    float64
	haste float64
	exp   float64
}

func runShaman(t *testing.T, gear *proto.EquipmentSpec, talents string, encounter *proto.Encounter, iterations int32) simResult {
	raid := core.SinglePlayerRaidProto(shamanPlayer(gear, talents), core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs)

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

func iterationCount() int32 {
	if v := os.Getenv("SHAMAN_ITERATIONS"); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return int32(n)
		}
	}
	return 20000
}

func TestShamanP3Upgrades(t *testing.T) {
	RegisterAll()
	iterations := iterationCount()
	encounter := core.MakeSingleTargetEncounter(5)

	base := runShaman(t, currentGear(), talentsKoerriiPvE, encounter, iterations)
	fmt.Printf("ITERATIONS %d\n", iterations)
	fmt.Printf("BASELINE | current armory gear | DPS %.1f (stdev %.1f) | AP %.0f hit %.0f crit %.0f haste %.0f exp %.0f\n",
		base.dps, base.stdev, base.ap, base.hit, base.crit, base.haste, base.exp)

	// --- Talent sanity check --------------------------------------------------
	fmt.Println("\n=== TALENTS (same gear) ===")
	for _, tc := range []struct{ label, talents string }{
		{"Armory spec 2 (0/47/14, PvE) [baseline]", talentsKoerriiPvE},
		{"Armory spec 1 (17/44/0)", talentsKoerriiPrimary},
		{"Repo default enhancement (3/44/14)", talentsRepoDefault},
	} {
		r := runShaman(t, currentGear(), tc.talents, encounter, iterations)
		fmt.Printf("TALENT | %-42s | %8.1f | %+8.1f\n", tc.label, r.dps, r.dps-base.dps)
	}

	// --- Free upgrades: enchants and gems ------------------------------------
	fmt.Println("\n=== ENCHANT / GEM FIXES (no drops needed) ===")
	for _, v := range enchantAndGemFixes() {
		gear := cloneGear(currentGear())
		v.apply(gear)
		r := runShaman(t, gear, talentsKoerriiPvE, encounter, iterations)
		fmt.Printf("FIX | %-58s | %8.1f | %+7.1f | %s\n", v.label, r.dps, r.dps-base.dps, v.note)
	}

	// --- Single-slot item swaps ----------------------------------------------
	fmt.Println("\n=== SINGLE-SLOT ITEM SWAPS (vs current gear) ===")
	type scored struct {
		c candidate
		r simResult
	}
	all := []scored{}
	for _, c := range itemCandidates() {
		gear := cloneGear(currentGear())
		c.apply(gear)
		r := runShaman(t, gear, talentsKoerriiPvE, encounter, iterations)
		all = append(all, scored{c, r})
		fmt.Printf("SWAP | %-9s | %-38s | %6d | %8.1f | %+7.1f | AP %5.0f hit %4.0f crit %4.0f haste %4.0f exp %3.0f | %s\n",
			slotNames[c.slot], c.label, c.id, r.dps, r.dps-base.dps, r.ap, r.hit, r.crit, r.haste, r.exp, c.note)
	}

	// --- Cumulative stacking of the best option per slot ----------------------
	// Greedy, highest single-slot DPS first, refusing to hand the same item id to
	// two different slots (rings and trinkets are unique-equipped in game, and
	// this engine does not enforce that). Two copies of the same weapon are legal
	// and stay allowed.
	sortedAll := append([]scored(nil), all...)
	sort.Slice(sortedAll, func(i, j int) bool { return sortedAll[i].r.dps > sortedAll[j].r.dps })

	assigned := map[int]scored{}
	usedIDs := map[int32]bool{}
	ordered := []scored{}
	for _, s := range sortedAll {
		if s.r.dps <= base.dps {
			continue
		}
		if _, taken := assigned[s.c.slot]; taken {
			continue
		}
		isWeapon := s.c.slot == slotMainHand || s.c.slot == slotOffHand
		if !isWeapon && usedIDs[s.c.id] {
			continue
		}
		assigned[s.c.slot] = s
		usedIDs[s.c.id] = true
		ordered = append(ordered, s)
	}

	fmt.Println("\n=== CUMULATIVE (best item per slot, applied in single-slot rank order) ===")
	gear := cloneGear(currentGear())
	prev := base.dps
	for _, s := range ordered {
		s.c.apply(gear)
		r := runShaman(t, cloneGear(gear), talentsKoerriiPvE, encounter, iterations)
		fmt.Printf("CUMUL | +%-9s %-38s | %8.1f | step %+7.1f | total %+8.1f\n",
			slotNames[s.c.slot], s.c.label, r.dps, r.dps-prev, r.dps-base.dps)
		prev = r.dps
	}

	// --- Assembled builds -----------------------------------------------------
	// bestFor returns the highest-DPS candidate for a slot subject to a filter.
	bestFor := func(slot int, keep func(candidate) bool) (scored, bool) {
		var out scored
		found := false
		for _, s := range sortedAll {
			if s.c.slot != slot || (keep != nil && !keep(s.c)) {
				continue
			}
			if !found || s.r.dps > out.r.dps {
				out, found = s, true
			}
		}
		return out, found
	}

	tierSlots := []int{slotHead, slotShoulder, slotChest, slotLegs, slotHands}
	isTierSlot := func(slot int) bool {
		for _, ts := range tierSlots {
			if ts == slot {
				return true
			}
		}
		return false
	}

	// A: every slot upgraded to its best option (= end of the cumulative pass).
	buildAllBest := func() *proto.EquipmentSpec {
		g := cloneGear(currentGear())
		for _, s := range ordered {
			s.c.apply(g)
		}
		return g
	}

	// B: same, but only the single best tier slot is swapped, so 4pc Cataclysm
	// Harness (T5) survives.
	buildKeep4pcT5 := func() *proto.EquipmentSpec {
		g := cloneGear(currentGear())
		bestTier := scored{}
		haveTier := false
		for _, s := range ordered {
			if isTierSlot(s.c.slot) {
				if !haveTier || s.r.dps > bestTier.r.dps {
					bestTier, haveTier = s, true
				}
				continue
			}
			s.c.apply(g)
		}
		if haveTier {
			bestTier.c.apply(g)
		}
		return g
	}

	// C: 4pc Skyshatter Harness (T6) in head/shoulder/chest/legs, best non-tier
	// elsewhere, best non-tier-set option for hands.
	skyshatter := map[int]candidate{
		slotHead:     {slotHead, "Skyshatter Cover (T6)", 31015, nil, 0, ""},
		slotShoulder: {slotShoulder, "Skyshatter Pauldrons (T6)", 31024, nil, 0, ""},
		slotChest:    {slotChest, "Skyshatter Tunic (T6)", 31018, nil, 0, ""},
		slotLegs:     {slotLegs, "Skyshatter Pants (T6)", 31021, nil, 0, ""},
	}
	buildSkyshatter4pc := func() *proto.EquipmentSpec {
		g := cloneGear(currentGear())
		for _, s := range ordered {
			if _, isSky := skyshatter[s.c.slot]; isSky {
				continue
			}
			s.c.apply(g)
		}
		if s, ok := bestFor(slotHands, func(c candidate) bool { return c.id != 31011 }); ok {
			s.c.apply(g)
		}
		for _, c := range skyshatter {
			c.apply(g)
		}
		return g
	}

	applyFixes := func(g *proto.EquipmentSpec) *proto.EquipmentSpec {
		applyEnchantFixes(g)
		regemAll(g, bestGemsFor)
		return g
	}

	fmt.Println("\n=== ASSEMBLED BUILDS ===")
	for _, b := range []struct {
		label string
		gear  *proto.EquipmentSpec
	}{
		{"A. best item every slot", buildAllBest()},
		{"A+. best item every slot + enchant/gem fixes", applyFixes(buildAllBest())},
		{"B. best items, keep 4pc Cataclysm Harness (T5)", buildKeep4pcT5()},
		{"B+. B + enchant/gem fixes", applyFixes(buildKeep4pcT5())},
		{"C. best items with 4pc Skyshatter (T6)", buildSkyshatter4pc()},
		{"C+. C + enchant/gem fixes", applyFixes(buildSkyshatter4pc())},
	} {
		r := runShaman(t, cloneGear(b.gear), talentsKoerriiPvE, encounter, iterations)
		fmt.Printf("BUILD | %-46s | %8.1f | %+8.1f | AP %5.0f hit %4.0f crit %4.0f haste %4.0f exp %3.0f\n",
			b.label, r.dps, r.dps-base.dps, r.ap, r.hit, r.crit, r.haste, r.exp)
		for i, it := range b.gear.Items {
			fmt.Printf("BUILDITEM | %-46s | %-9s | %6d | ench %5d | gems %v\n", b.label, slotNames[i], it.Id, it.Enchant, it.Gems)
		}
	}

	// --- Reference points -----------------------------------------------------
	fmt.Println("\n=== REFERENCE BUILDS ===")
	for _, ref := range []struct {
		label string
		dir   string
		file  string
	}{
		{"repo P2 preset", "../ui/specs/shaman/enhancement/gear_sets", "p2"},
		{"repo P3 BiS preset", "../ui/specs/shaman/enhancement/gear_sets", "p3"},
		{"repo P4 BiS preset", "../ui/specs/shaman/enhancement/gear_sets", "p4"},
	} {
		raw := runShaman(t, core.GetGearSet(ref.dir, ref.file).GearSet, talentsKoerriiPvE, encounter, iterations)
		fmt.Printf("REF | %-22s (as checked in)     | %8.1f | %+8.1f | AP %5.0f hit %4.0f crit %4.0f haste %4.0f exp %3.0f\n",
			ref.label, raw.dps, raw.dps-base.dps, raw.ap, raw.hit, raw.crit, raw.haste, raw.exp)

		fixed := runShaman(t, fixPresetEnchants(core.GetGearSet(ref.dir, ref.file).GearSet), talentsKoerriiPvE, encounter, iterations)
		fmt.Printf("REF | %-22s (enchant ids fixed) | %8.1f | %+8.1f | AP %5.0f hit %4.0f crit %4.0f haste %4.0f exp %3.0f\n",
			ref.label, fixed.dps, fixed.dps-base.dps, fixed.ap, fixed.hit, fixed.crit, fixed.haste, fixed.exp)
	}
}
