package sim

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
)

// Hallow's End check for the Anniversary ilvl 141 re-issues from the Headless
// Horseman. Each item goes into the slot it would take, one slot at a time, and
// is compared against what was there: Koerry's and Koerrii's own gear first,
// then the Phase 3 presets of every spec that would want it. Rings are tried in
// both finger slots and the better swap is reported.
//
// The Witches Band (healing ring) is not simmed: no healer sim here reports a
// number it would move. The Helm's laugh and the Blade's Pumpkin Soldiers are
// flavour/unimplemented on-use effects, so both are pure stat sticks here.
//
//	HALLOWS_ITERATIONS=3000 go test -tags with_db ./sim/ -run TestHallowsEnd -v

const (
	itemRingOfGhoulishDelight = 281907 // 43 sta, 80 AP
	itemHorsemansSignetRing   = 281911 // 25 sta, 23 int, 47 spell damage
	itemHorsemansHelm         = 281919 // plate: 53 str, 62 sta, 41 crit, red/yellow/blue, no meta socket
	itemHorsemansBlade        = 281923 // main-hand sword 2.7: 18 agi, 20 sta, 60 AP
)

func hallowsIterationCount() int32 {
	if v := os.Getenv("HALLOWS_ITERATIONS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return int32(n)
		}
	}
	return 10000
}

// hallowsSwap puts one item into one slot, gemmed from its own sockets and
// keeping the slot's enchant.
func hallowsSwap(gear *proto.EquipmentSpec, slot int, id int32) *proto.EquipmentSpec {
	g := cloneGear(gear)
	g.Items[slot] = itemSpec(id, gear.Items[slot].Enchant, bestGemsFor(id)...)
	return g
}

func itemName(id int32) string {
	if item := core.GetItemByID(id); item != nil {
		return item.Name
	}
	return fmt.Sprintf("item %d", id)
}

type hallowsTry struct {
	label string
	item  int32
	slots []int
}

// hallowsReport sims every try against the base gear and prints the best slot
// for each, so a ring shows the finger it should replace.
func hallowsReport(who string, base *proto.EquipmentSpec, run func(*proto.EquipmentSpec) simResult, tries []hallowsTry) {
	b := run(base)
	fmt.Printf("HALLOWBASE | %-34s | DPS %8.1f (stdev %.1f)\n", who, b.dps, b.stdev)
	for _, tr := range tries {
		bestSlot, bestDps := -1, 0.0
		for _, slot := range tr.slots {
			if r := run(hallowsSwap(base, slot, tr.item)); bestSlot < 0 || r.dps > bestDps {
				bestSlot, bestDps = slot, r.dps
			}
		}
		replaced := itemName(base.Items[bestSlot].Id)
		fmt.Printf("HALLOW | %-34s | %-26s | replaces %-32s | %8.1f | %+7.1f (%+.2f%%)\n",
			who, tr.label, replaced, bestDps, bestDps-b.dps, 100*(bestDps-b.dps)/b.dps)
	}
}

var (
	hallowsRing  = hallowsTry{"Ring of Ghoulish Delight", itemRingOfGhoulishDelight, []int{slotFinger1, slotFinger2}}
	hallowsHelm  = hallowsTry{"The Horseman's Helm", itemHorsemansHelm, []int{slotHead}}
	hallowsBlade = hallowsTry{"The Horseman's Blade", itemHorsemansBlade, []int{slotMainHand}} // main hand only
	hallowsSig   = hallowsTry{"The Horseman's Signet Ring", itemHorsemansSignetRing, []int{slotFinger1, slotFinger2}}
)

func TestHallowsEndKoerry(t *testing.T) {
	RegisterAll()
	iterations := hallowsIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)
	fmt.Printf("ITERATIONS %d\n", iterations)

	warriorBaseGear = warriorCurrentGear
	run := func(gear *proto.EquipmentSpec) simResult {
		return runWarrior(t, gear, talentsFuryStandard, "fury", encounter, iterations)
	}
	tries := []hallowsTry{hallowsRing, hallowsHelm, hallowsBlade}
	hallowsReport("Koerry 4/5 T5 (goggles)", warriorCurrentGear(), run, tries)
	hallowsReport("Koerry 5/5 T5 (Destroyer helm)", warriorT5Gear(), run, tries)
}

func TestHallowsEndKoerrii(t *testing.T) {
	RegisterAll()
	iterations := hallowsIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)
	run := func(gear *proto.EquipmentSpec) simResult {
		return runShaman(t, gear, talentsKoerriiPvE, encounter, iterations)
	}
	// Shamans can't wear plate or use swords, so only the ring applies.
	hallowsReport("Koerrii (Enhancement)", currentGear(), run, []hallowsTry{hallowsRing})
}

// physicalTries maps the trinket benchmark's spec names to what each can use.
var physicalTries = map[string][]hallowsTry{
	"Warrior Fury":         {hallowsRing, hallowsHelm, hallowsBlade},
	"Warrior Arms (Kebab)": {hallowsRing, hallowsHelm, hallowsBlade},
	"Warrior Arms":         {hallowsRing, hallowsHelm},
	"Rogue":                {hallowsRing, hallowsBlade},
	"Hunter BM":            {hallowsRing},
	"Hunter SV":            {hallowsRing},
	"Enhancement Shaman":   {hallowsRing},
	"Feral Druid":          {hallowsRing},
	"Retribution Paladin":  {hallowsRing, hallowsHelm},
}

func TestHallowsEndPhysicalPresets(t *testing.T) {
	RegisterAll()
	iterations := hallowsIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)
	fmt.Printf("ITERATIONS %d\n", iterations)

	for _, spec := range benchSpecs() {
		tries, ok := physicalTries[spec.name]
		if !ok {
			continue
		}
		gear := core.GetGearSet(spec.gearDir, spec.gearFile).GearSet
		if spec.weaponMH != 0 {
			gear.Items[slotMainHand].Id = spec.weaponMH
			gear.Items[slotOffHand].Id = spec.weaponOH
		}
		var rotation *proto.APLRotation
		if spec.catRotation != nil {
			rotation = &proto.APLRotation{PriorityList: []*proto.APLListItem{{Action: &proto.APLAction{
				Action: &proto.APLAction_CatOptimalRotationAction{CatOptimalRotationAction: &proto.APLActionCatOptimalRotationAction{
					FinishingMove:      spec.catRotation.FinishingMove,
					Biteweave:          spec.catRotation.Biteweave,
					RipMinComboPoints:  spec.catRotation.RipMinComboPoints,
					BiteMinComboPoints: spec.catRotation.BiteMinComboPoints,
					MangleTrick:        spec.catRotation.MangleTrick,
					MaintainFaerieFire: spec.catRotation.MaintainFaerieFire,
				}},
			}}}}
		} else {
			rotation = core.GetAplRotation(spec.rotationDir, spec.rotationFile).Rotation
		}
		run := func(g *proto.EquipmentSpec) simResult {
			player := &proto.Player{
				Class:              spec.class,
				Race:               spec.race,
				Equipment:          g,
				Consumables:        spec.consumables,
				TalentsString:      spec.talents,
				Profession1:        spec.profession1,
				Profession2:        spec.profession2,
				DistanceFromTarget: spec.distance,
				Rotation:           rotation,
				Buffs:              core.FullIndividualBuffs,
			}
			core.WithSpec(player, spec.spec)
			return runHallowsPlayer(t, player, encounter, iterations)
		}
		hallowsReport(spec.name+" P3", gear, run, tries)
	}
}

type casterSpec struct {
	name         string
	class        proto.Class
	race         proto.Race
	spec         interface{}
	talents      string
	flask        int32
	gearDir      string
	gearFile     string
	rotationDir  string
	rotationFile string
}

func TestHallowsEndCasterPresets(t *testing.T) {
	RegisterAll()
	iterations := hallowsIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)
	fmt.Printf("ITERATIONS %d\n", iterations)

	const flaskBlindingLight, flaskPureDeath = 22861, 22866
	specs := []casterSpec{
		{
			name: "Arcane Mage", class: proto.Class_ClassMage, race: proto.Race_RaceTroll,
			spec: &proto.Player_Mage{Mage: &proto.Mage{Options: &proto.Mage_Options{
				ClassOptions: &proto.MageOptions{DefaultMageArmor: proto.MageArmor_MageArmorMageArmor},
			}}},
			talents: "2500052300030150330125--053500031003001", flask: flaskBlindingLight,
			gearDir: "../ui/specs/mage/dps/gear_sets", gearFile: "p3ArcaneStaff",
			rotationDir: "../ui/specs/mage/dps/apls", rotationFile: "arcane",
		},
		{
			name: "Destruction Warlock", class: proto.Class_ClassWarlock, race: proto.Race_RaceOrc,
			spec: &proto.Player_Warlock{Warlock: &proto.Warlock{Options: &proto.Warlock_Options{
				ClassOptions: &proto.WarlockOptions{
					Summon:          proto.WarlockOptions_Succubus,
					SacrificeSummon: true,
					Armor:           proto.WarlockOptions_FelArmor,
					CurseOptions:    proto.WarlockOptions_Agony,
				},
			}}},
			talents: "-20500301332101-50500051220051053105", flask: flaskPureDeath,
			gearDir: "../ui/specs/warlock/dps/gear_sets", gearFile: "t6",
			rotationDir: "../ui/specs/warlock/dps/apls", rotationFile: "destruction",
		},
		{
			name: "Shadow Priest", class: proto.Class_ClassPriest, race: proto.Race_RaceTroll,
			spec: &proto.Player_DpsPriest{DpsPriest: &proto.DpsPriest{Options: &proto.DpsPriest_Options{
				ClassOptions: &proto.PriestOptions{PreShadowform: true},
			}}},
			talents: "500230013--503250510240103051451", flask: flaskPureDeath,
			gearDir: "../ui/specs/priest/dps/gear_sets", gearFile: "p3",
			rotationDir: "../ui/specs/priest/dps/apls", rotationFile: "default",
		},
		{
			name: "Elemental Shaman", class: proto.Class_ClassShaman, race: proto.Race_RaceTroll,
			spec: &proto.Player_ElementalShaman{ElementalShaman: &proto.ElementalShaman{Options: &proto.ElementalShaman_Options{
				ClassOptions: &proto.ShamanOptions{},
			}}},
			talents: "55003105100213351051--05105301005", flask: flaskBlindingLight,
			gearDir: "../ui/specs/shaman/elemental/gear_sets", gearFile: "p3",
			rotationDir: "../ui/specs/shaman/elemental/apls", rotationFile: "default",
		},
		{
			name: "Balance Druid", class: proto.Class_ClassDruid, race: proto.Race_RaceNightElf,
			spec: &proto.Player_BalanceDruid{BalanceDruid: &proto.BalanceDruid{Options: &proto.BalanceDruid_Options{
				ClassOptions: &proto.DruidOptions{},
			}}},
			talents: "510022312503135231351--520033", flask: flaskBlindingLight,
			gearDir: "../ui/specs/druid/balance/gear_sets", gearFile: "p3",
			rotationDir: "../ui/specs/druid/balance/apls", rotationFile: "default",
		},
	}

	for _, spec := range specs {
		rotation := core.GetAplRotation(spec.rotationDir, spec.rotationFile).Rotation
		run := func(g *proto.EquipmentSpec) simResult {
			player := &proto.Player{
				Class:     spec.class,
				Race:      spec.race,
				Equipment: g,
				Consumables: &proto.ConsumesSpec{
					FlaskId: spec.flask, FoodId: 27657, PotId: 22839, ConjuredId: 12662,
				},
				TalentsString:      spec.talents,
				DistanceFromTarget: 20,
				Rotation:           rotation,
				Buffs:              core.FullIndividualBuffs,
			}
			core.WithSpec(player, spec.spec)
			return runHallowsPlayer(t, player, encounter, iterations)
		}
		hallowsReport(spec.name+" P3", core.GetGearSet(spec.gearDir, spec.gearFile).GearSet, run, []hallowsTry{hallowsSig})
	}
}

func runHallowsPlayer(t *testing.T, player *proto.Player, encounter *proto.Encounter, iterations int32) simResult {
	raid := core.SinglePlayerRaidProto(player, core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs)
	result := core.RunRaidSimConcurrent(&proto.RaidSimRequest{
		Raid:       raid,
		Encounter:  encounter,
		SimOptions: &proto.SimOptions{Iterations: iterations, RandomSeed: 101},
	})
	if result.Error != nil {
		t.Fatalf("sim failed: %s", result.Error.Message)
	}
	return simResult{dps: result.RaidMetrics.Dps.Avg, stdev: result.RaidMetrics.Dps.Stdev}
}

// Koerry's main hand: the Horseman's Blade against the Blacksmithing
// main-hand weapons he can craft, all with Mongoose. All are 2.7 speed and
// main-hand only; Dragonstrike and Dragonmaw are maces with a haste proc.
func TestHallowsEndKoerryBlacksmith(t *testing.T) {
	RegisterAll()
	iterations := hallowsIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)
	fmt.Printf("ITERATIONS %d\n", iterations)

	warriorBaseGear = warriorCurrentGear
	run := func(gear *proto.EquipmentSpec) simResult {
		return runWarrior(t, gear, talentsFuryStandard, "fury", encounter, iterations)
	}
	weapons := []hallowsTry{
		hallowsBlade,
		{"Dragonstrike", 28439, []int{slotMainHand}},
		{"Dragonmaw", 28438, []int{slotMainHand}},
		{"Wicked Edge of the Planes", 28433, []int{slotMainHand}},
	}
	hallowsReport("Koerry 4/5 T5 (goggles)", warriorCurrentGear(), run, weapons)
	hallowsReport("Koerry 5/5 T5 (Destroyer helm)", warriorT5Gear(), run, weapons)
}

// Koerrii can't use the Horseman's Blade (shamans can't wield swords), so this
// compares the Blacksmithing main hands instead, against the current
// Merciless Gladiator's Cleaver. They need Blacksmithing to equip, which is
// assumed here. Vengeful Gladiator's Cleaver is the known arena upgrade, as a
// reference. Koerrii is an Orc, so the axes get Axe Specialization.
func TestHallowsEndKoerriiBlacksmith(t *testing.T) {
	RegisterAll()
	iterations := hallowsIterationCount()
	encounter := core.MakeSingleTargetEncounter(5)
	fmt.Printf("ITERATIONS %d\n", iterations)

	run := func(gear *proto.EquipmentSpec) simResult {
		player := shamanPlayer(gear, talentsKoerriiPvE)
		player.Profession1 = proto.Profession_Blacksmithing
		return runHallowsPlayer(t, player, encounter, iterations)
	}
	hallowsReport("Koerrii (Enhancement)", currentGear(), run, []hallowsTry{
		{"Dragonstrike", 28439, []int{slotMainHand}},
		{"Dragonmaw", 28438, []int{slotMainHand}},
		{"Wicked Edge of the Planes", 28433, []int{slotMainHand}},
		{"Vengeful Gladiator's Cleaver", 33669, []int{slotMainHand}},
	})
}
