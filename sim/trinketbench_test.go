package sim

import (
	"fmt"
	"testing"

	"github.com/wowsims/tbc/sim/core"
	"github.com/wowsims/tbc/sim/core/proto"
	"github.com/wowsims/tbc/sim/core/stats"
)

// Scratch benchmark: P3 trinket comparison across physical-DPS specs.
// Slot1 (Trinket1) is fixed to Bloodlust Brooch; Slot2 (Trinket2) cycles
// through the 4 candidate trinkets. Every run gets +30 bonus hit rating and
// the same core.FullRaidBuffs / FullPartyBuffs / FullIndividualBuffs / FullDebuffs
// so every spec is compared under identical raid conditions.
// Delete this file after extracting results.

const (
	itemBloodlustBrooch = 29383
	itemWarpSpringCoil  = 30450 // Rogue-only fixed trinket1 (replaces Bloodlust Brooch for Rogue)
	itemAbacus          = 28288
	itemMadness         = 32505
	itemDragonspineTr   = 28830
	itemTsunamiTalisman = 30627
	// Warrior-only Warglaive MH/OH substitute: S3 (Vengeful Gladiator) slow
	// one-hand axes. Warp Slicer was rejected since it's effectively a
	// legendary-tier item (Kael'thas quest chain reward). Rogues can't use
	// axes, so they get their own sword substitute below.
	itemVengefulGladiatorsCleaver = 33669 // MH, 1H axe, 2.6 speed
	itemVengefulGladiatorsChopper = 34015 // OH, 1H axe, 2.6 speed
	// Rogue-only Warglaive substitute: next-best non-legendary swords per
	// Wowhead P3 BiS guide.
	itemVengefulGladiatorsSlicer = 33762 // MH, 1H sword, 2.6 speed
	itemBladeOfSavagery          = 32369 // OH, 1H sword, 1.4 speed
	// Standard P3 BiS legendary weapons, used for the separate "(Glaives)" categories.
	itemWarglaiveMH = 32837
	itemWarglaiveOH = 32838
)

type trinketVariant struct {
	label    string
	trinket2 int32
}

var trinketVariants = []trinketVariant{
	{"Abacus of Violent Odds", itemAbacus},
	{"Madness of the Betrayer", itemMadness},
	{"Dragonspine Trophy", itemDragonspineTr},
	{"Tsunami Talisman", itemTsunamiTalisman},
}

type benchSpecConfig struct {
	name         string
	class        proto.Class
	race         proto.Race
	spec         interface{}
	talents      string
	consumables  *proto.ConsumesSpec
	profession1  proto.Profession
	profession2  proto.Profession
	distance     float64
	gearDir      string
	gearFile     string
	rotationDir  string
	rotationFile string
	catRotation  *proto.FeralCatDruid_Rotation // Feral Druid: builds a real APL priority list with CatOptimalRotationAction
	weaponMH     int32                         // if set, overrides Equipment.Items[14].Id
	weaponOH     int32                         // if set, overrides Equipment.Items[15].Id
	trinket1     int32                         // fixed Trinket1 override; defaults to Bloodlust Brooch if 0
	meleeWeave   string                        // Hunter: "true"/"false" override for "Melee weave" APL variable; "" = leave APL default (true)
}

func benchSpecs() []benchSpecConfig {
	warriorOptions := &proto.Player_DpsWarrior{
		DpsWarrior: &proto.DpsWarrior{
			Options: &proto.DpsWarrior_Options{
				ClassOptions: &proto.WarriorOptions{
					DefaultShout:  proto.WarriorShout_WarriorShoutBattle,
					DefaultStance: proto.WarriorStance_WarriorStanceBerserker,
				},
			},
		},
	}
	warriorConsumables := &proto.ConsumesSpec{
		PotId:       22838,
		FlaskId:     22854,
		FoodId:      27658,
		ConjuredId:  22788,
		ExplosiveId: 30217,
		SuperSapper: true,
		OhImbueId:   29453,
		ScrollAgi:   true,
		ScrollStr:   true,
	}

	hunterOptions := &proto.Player_Hunter{
		Hunter: &proto.Hunter{
			Options: &proto.Hunter_Options{
				ClassOptions: &proto.HunterOptions{
					Ammo:             proto.HunterOptions_AdamantiteStinger,
					PetSingleAbility: false,
					PetType:          proto.HunterOptions_Ravager,
					PetUptime:        100.0,
					QuiverBonus:      proto.HunterOptions_Speed15,
				},
			},
		},
	}
	hunterConsumables := &proto.ConsumesSpec{
		BattleElixirId:   22831,
		GuardianElixirId: 22840,
		FoodId:           27659,
		PotId:            22838,
		ConjuredId:       12662,
		ExplosiveId:      30217,
		PetFoodId:        33874,
		PetScrollAgi:     true,
		PetScrollStr:     true,
		SuperSapper:      true,
		GoblinSapper:     true,
		ScrollAgi:        true,
		ScrollStr:        true,
	}

	return []benchSpecConfig{
		{
			name: "Warrior Fury", class: proto.Class_ClassWarrior, race: proto.Race_RaceOrc,
			spec: warriorOptions, talents: "3500501130201-05050005505012050115",
			consumables: warriorConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 25,
			gearDir:  "../ui/warrior/dps/gear_sets", gearFile: "p3_fury",
			rotationDir: "../ui/warrior/dps/apls", rotationFile: "fury",
			weaponMH: itemVengefulGladiatorsCleaver, weaponOH: itemVengefulGladiatorsChopper,
		},
		{
			name: "Warrior Fury (Glaives)", class: proto.Class_ClassWarrior, race: proto.Race_RaceOrc,
			spec: warriorOptions, talents: "3500501130201-05050005505012050115",
			consumables: warriorConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 25,
			gearDir:  "../ui/warrior/dps/gear_sets", gearFile: "p3_fury",
			rotationDir: "../ui/warrior/dps/apls", rotationFile: "fury",
			// Standard P3 BiS gear already equips Warglaives - no weapon override.
		},
		{
			name: "Warrior Arms", class: proto.Class_ClassWarrior, race: proto.Race_RaceOrc,
			spec: warriorOptions, talents: "32005011352010500221-0550000500521203",
			consumables: warriorConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 25,
			gearDir:  "../ui/warrior/dps/gear_sets", gearFile: "p3_arms",
			rotationDir: "../ui/warrior/dps/apls", rotationFile: "arms",
		},
		{
			name: "Warrior Arms (Glaives)", class: proto.Class_ClassWarrior, race: proto.Race_RaceOrc,
			// Dual-wielding Warglaives needs the dual-wield "Kebab" talent build (has Dual Wield
			// Specialization), not the plain 2H Arms build - see ui/warrior/dps/presets.ts ArmsKebabTalents.
			spec: warriorOptions, talents: "34005021302010510321-0550000520501203",
			consumables: warriorConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 25,
			gearDir:  "../ui/warrior/dps/gear_sets", gearFile: "p3_arms",
			rotationDir: "../ui/warrior/dps/apls", rotationFile: "arms",
			weaponMH: itemWarglaiveMH, weaponOH: itemWarglaiveOH,
		},
		{
			name: "Warrior Arms (Kebab)", class: proto.Class_ClassWarrior, race: proto.Race_RaceOrc,
			// Affordable dual-wield Arms variant: same Kebab talents as "Warrior Arms (Glaives)"
			// above, but with cheaper S3 swords instead of legendary Warglaives.
			spec: warriorOptions, talents: "34005021302010510321-0550000520501203",
			consumables: warriorConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 25,
			gearDir:  "../ui/warrior/dps/gear_sets", gearFile: "p3_arms",
			rotationDir: "../ui/warrior/dps/apls", rotationFile: "arms",
			weaponMH: itemVengefulGladiatorsSlicer, weaponOH: itemBladeOfSavagery,
		},
		{
			name: "Rogue", class: proto.Class_ClassRogue, race: proto.Race_RaceHuman,
			spec:    &proto.Player_Rogue{Rogue: &proto.Rogue{Options: &proto.Rogue_Options{ClassOptions: &proto.RogueOptions{}}}},
			talents: "00532012502-023305200005015002321151",
			consumables: &proto.ConsumesSpec{
				FlaskId: 22854, FoodId: 33872, PotId: 22838, ConjuredId: 7676, OhImbueId: 27186,
			},
			gearDir: "../ui/rogue/dps/gear_sets", gearFile: "p3",
			rotationDir: "../ui/rogue/dps/apls", rotationFile: "swords",
			weaponMH: itemVengefulGladiatorsSlicer, weaponOH: itemBladeOfSavagery,
			trinket1: itemWarpSpringCoil,
		},
		{
			name: "Rogue (Glaives)", class: proto.Class_ClassRogue, race: proto.Race_RaceHuman,
			spec:    &proto.Player_Rogue{Rogue: &proto.Rogue{Options: &proto.Rogue_Options{ClassOptions: &proto.RogueOptions{}}}},
			talents: "00532012502-023305200005015002321151",
			consumables: &proto.ConsumesSpec{
				FlaskId: 22854, FoodId: 33872, PotId: 22838, ConjuredId: 7676, OhImbueId: 27186,
			},
			gearDir: "../ui/rogue/dps/gear_sets", gearFile: "p3",
			rotationDir: "../ui/rogue/dps/apls", rotationFile: "swords",
			trinket1: itemWarpSpringCoil,
			// Standard P3 BiS gear already equips Warglaives - no weapon override.
		},
		{
			name: "Hunter BM", class: proto.Class_ClassHunter, race: proto.Race_RaceOrc,
			spec: hunterOptions, talents: "512002005250122431051-0505201205",
			consumables: hunterConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 30,
			gearDir:  "../ui/hunter/dps/gear_sets/phase_3/bm", gearFile: "2h_6p",
			rotationDir: "../ui/hunter/dps/apls", rotationFile: "default",
			meleeWeave: "false",
		},
		{
			name: "Hunter BM (Weave)", class: proto.Class_ClassHunter, race: proto.Race_RaceOrc,
			spec: hunterOptions, talents: "512002005250122431051-0505201205",
			consumables: hunterConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 30,
			gearDir:  "../ui/hunter/dps/gear_sets/phase_3/bm", gearFile: "2h_6p",
			rotationDir: "../ui/hunter/dps/apls", rotationFile: "default",
			meleeWeave: "true",
		},
		{
			name: "Hunter SV", class: proto.Class_ClassHunter, race: proto.Race_RaceOrc,
			spec: hunterOptions, talents: "502-0550201205-333200022003223005103",
			consumables: hunterConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 30,
			gearDir:  "../ui/hunter/dps/gear_sets/phase_3/sv", gearFile: "2h_6p",
			rotationDir: "../ui/hunter/dps/apls", rotationFile: "default",
			meleeWeave: "false",
		},
		{
			name: "Hunter SV (Weave)", class: proto.Class_ClassHunter, race: proto.Race_RaceOrc,
			spec: hunterOptions, talents: "502-0550201205-333200022003223005103",
			consumables: hunterConsumables, profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			distance: 30,
			gearDir:  "../ui/hunter/dps/gear_sets/phase_3/sv", gearFile: "2h_6p",
			rotationDir: "../ui/hunter/dps/apls", rotationFile: "default",
			meleeWeave: "true",
		},
		{
			name: "Enhancement Shaman", class: proto.Class_ClassShaman, race: proto.Race_RaceTroll,
			spec: &proto.Player_EnhancementShaman{
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
			},
			talents: "03-500502210501133531151-50005301",
			consumables: &proto.ConsumesSpec{
				PotId: 22838, FlaskId: 22854, FoodId: 27658, ConjuredId: 22788, ExplosiveId: 30217,
				SuperSapper: true, GoblinSapper: true, ScrollAgi: true, ScrollStr: true,
			},
			gearDir: "../ui/shaman/enhancement/gear_sets", gearFile: "p3",
			rotationDir: "../ui/shaman/enhancement/apls", rotationFile: "default",
		},
		{
			name: "Feral Druid", class: proto.Class_ClassDruid, race: proto.Race_RaceNightElf,
			spec: &proto.Player_FeralCatDruid{
				FeralCatDruid: &proto.FeralCatDruid{
					Rotation: &proto.FeralCatDruid_Rotation{
						FinishingMove:      proto.FeralCatDruid_Rotation_Rip,
						Biteweave:          true,
						RipMinComboPoints:  5,
						BiteMinComboPoints: 5,
						MangleTrick:        true,
						MaintainFaerieFire: false,
					},
					Options: &proto.FeralCatDruid_Options{},
				},
			},
			talents: "-503032132322105301251-05503301",
			consumables: &proto.ConsumesSpec{
				PotId: 22838, BattleElixirId: 22831, GuardianElixirId: 32067, FoodId: 27664,
				MhImbueId: 34340, ConjuredId: 12662, SuperSapper: true, GoblinSapper: true,
				ScrollAgi: true, ScrollStr: true,
			},
			profession1: proto.Profession_Engineering, profession2: proto.Profession_Enchanting,
			gearDir: "../ui/druid/feralcat/gear_sets", gearFile: "p3_6p",
			catRotation: &proto.FeralCatDruid_Rotation{
				FinishingMove:      proto.FeralCatDruid_Rotation_Rip,
				Biteweave:          true,
				RipMinComboPoints:  5,
				BiteMinComboPoints: 5,
				MangleTrick:        true,
				MaintainFaerieFire: false,
			},
		},
		{
			name: "Retribution Paladin", class: proto.Class_ClassPaladin, race: proto.Race_RaceBloodElf,
			spec: &proto.Player_RetributionPaladin{
				RetributionPaladin: &proto.RetributionPaladin{
					Options: &proto.RetributionPaladin_Options{ClassOptions: &proto.PaladinOptions{}},
				},
			},
			talents: "5-053201-0523005120033125331051",
			consumables: &proto.ConsumesSpec{
				PotId: 22838, FlaskId: 22854, FoodId: 27658, ConjuredId: 12662, ExplosiveId: 30217,
				SuperSapper: true, GoblinSapper: true, ScrollAgi: true, ScrollStr: true,
			},
			profession1: proto.Profession_Engineering, profession2: proto.Profession_Blacksmithing,
			gearDir: "../ui/paladin/retribution/gear_sets", gearFile: "p3",
			rotationDir: "../ui/paladin/retribution/apls", rotationFile: "default",
		},
	}
}

func TestTrinketBenchmark(t *testing.T) {
	runTrinketBenchmark(t, core.MakeSingleTargetEncounter(5))
}

// Cleave scenario: same 13 specs x 4 trinkets grid, but against 2 targets
// instead of 1, to see how AOE/cleave-capable rotations shift the trinket
// ranking (e.g. Whirlwind, Fan of Knives, Swipe picking up extra targets).
func TestTrinketBenchmarkCleave(t *testing.T) {
	encounter := core.MakeSingleTargetEncounter(5)
	encounter.Targets = []*proto.Target{
		core.NewDefaultTarget(),
		core.NewDefaultTarget(),
	}
	runTrinketBenchmark(t, encounter)
}

// 4-target scenario: same 14 specs x 4 trinkets grid, but against 4 targets
// instead of 1 or 2, to see how cleave/AOE effects scale further with more
// enemies present (e.g. radius-based abilities benefit more at higher target
// counts than single-extra-target cleaves).
func TestTrinketBenchmarkCleave4(t *testing.T) {
	encounter := core.MakeSingleTargetEncounter(5)
	encounter.Targets = []*proto.Target{
		core.NewDefaultTarget(),
		core.NewDefaultTarget(),
		core.NewDefaultTarget(),
		core.NewDefaultTarget(),
	}
	runTrinketBenchmark(t, encounter)
}

func runTrinketBenchmark(t *testing.T, encounter *proto.Encounter) {
	RegisterAll()

	fmt.Printf("%-22s | %-26s | %10s | %10s\n", "Spec", "Trinket2", "DPS", "Hit Rating")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, spec := range benchSpecs() {
		for _, variant := range trinketVariants {
			gearSet := core.GetGearSet(spec.gearDir, spec.gearFile)
			equipment := gearSet.GearSet
			trinket1 := spec.trinket1
			if trinket1 == 0 {
				trinket1 = itemBloodlustBrooch
			}
			equipment.Items[12].Id = trinket1
			equipment.Items[13].Id = variant.trinket2

			if spec.weaponMH != 0 {
				equipment.Items[14].Id = spec.weaponMH
				equipment.Items[15].Id = spec.weaponOH
			}

			var rotation *proto.APLRotation
			if spec.catRotation != nil {
				// Build the real priority list (matches ui/druid/feralcat/sim.ts'
				// simpleRotation generator) instead of an empty TypeSimple stub, which
				// would leave the character with zero actions (auto-attack only).
				rotation = &proto.APLRotation{
					PriorityList: []*proto.APLListItem{
						{
							Action: &proto.APLAction{
								Action: &proto.APLAction_CatOptimalRotationAction{
									CatOptimalRotationAction: &proto.APLActionCatOptimalRotationAction{
										FinishingMove:      spec.catRotation.FinishingMove,
										Biteweave:          spec.catRotation.Biteweave,
										Ripweave:           spec.catRotation.Ripweave,
										RipMinComboPoints:  spec.catRotation.RipMinComboPoints,
										BiteMinComboPoints: spec.catRotation.BiteMinComboPoints,
										MangleTrick:        spec.catRotation.MangleTrick,
										RakeTrick:          spec.catRotation.RakeTrick,
										MaintainFaerieFire: spec.catRotation.MaintainFaerieFire,
									},
								},
							},
						},
					},
				}
			} else {
				rotation = core.GetAplRotation(spec.rotationDir, spec.rotationFile).Rotation
			}
			if spec.meleeWeave != "" {
				rotation.ValueVariables[2] = &proto.APLValueVariable{
					Name: "Melee weave",
					Value: &proto.APLValue{
						Value: &proto.APLValue_Const{
							Const: &proto.APLValueConst{Val: spec.meleeWeave},
						},
					},
				}
			}

			player := &proto.Player{
				Class:              spec.class,
				Race:               spec.race,
				Equipment:          equipment,
				Consumables:        spec.consumables,
				BonusStats:         &proto.UnitStats{Stats: stats.Stats{stats.MeleeHitRating: 30}.ToProtoArray()},
				TalentsString:      spec.talents,
				Profession1:        spec.profession1,
				Profession2:        spec.profession2,
				DistanceFromTarget: spec.distance,
				Rotation:           rotation,
				Buffs:              core.FullIndividualBuffs,
			}
			core.WithSpec(player, spec.spec)

			raid := core.SinglePlayerRaidProto(player, core.FullPartyBuffs, core.FullRaidBuffs, core.FullDebuffs)

			computeResult := core.ComputeStats(&proto.ComputeStatsRequest{Raid: raid, Encounter: encounter})
			hitRating := computeResult.RaidStats.Parties[0].Players[0].FinalStats.Stats[stats.MeleeHitRating]

			rsr := &proto.RaidSimRequest{
				Raid:      raid,
				Encounter: encounter,
				SimOptions: &proto.SimOptions{
					Iterations: 20000,
					RandomSeed: 101,
				},
			}
			result := core.RunRaidSimConcurrent(rsr)
			if result.Error != nil {
				t.Fatalf("%s / %s failed: %s", spec.name, variant.label, result.Error.Message)
			}

			fmt.Printf("%-22s | %-26s | %10.1f | %10.1f\n", spec.name, variant.label, result.RaidMetrics.Dps.Avg, hitRating)
		}
	}
}
