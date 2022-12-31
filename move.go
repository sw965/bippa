package bippa

import (
	"fmt"
	"math/rand"
)

const (
	PHYSICS = "物理"
	SPECIAL = "特殊"
	STATUS  = "変化"
)

type StatusMove func(Battle, *rand.Rand) Battle

// あさのひざし
func MorningSun(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// こうごうせい
func Synthesis(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// じこさいせい
func Recover(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// すなあつめ
func ShoreUp(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// タマゴうみ
func SoftBoiled(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// つきのひかり
func Moonlight(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// なまける
func SlackOff(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// はねやすめ
func Roost(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// ミルクのみ
func MilkDrink(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

// どくどく
func Toxic(battle Battle, _ *rand.Rand) Battle {
	if battle.P2Fighters[0].StatusAilment != "" {
		return battle
	}

	if battle.P2Fighters[0].Types.In(POISON) {
		return battle
	}

	if battle.P2Fighters[0].Types.In(STEEL) {
		return battle
	}

	battle.P2Fighters[0].StatusAilment = BAD_POISON
	return battle
}

// やどりぎのタネ
func LeechSeed(battle Battle, _ *rand.Rand) Battle {
	if battle.P2Fighters[0].Types.In(GRASS) {
		return battle
	}

	battle.P2Fighters[0].IsLeechSeed = true
	return battle
}

// つるぎのまい
func SwordsDance(battle Battle, _ *rand.Rand) Battle {
	return battle.RankStateFluctuation(&RankState{Atk: 2})
}

// りゅうのまい
func DragonDance(battle Battle, _ *rand.Rand) Battle {
	return battle.RankStateFluctuation(&RankState{Atk: 1, Speed: 1})
}

// からをやぶる
func ShellSmash(battle Battle, _ *rand.Rand) Battle {
	return battle.RankStateFluctuation(&RankState{Atk: 2, Def: -1, SpAtk: 2, SpDef: -1, Speed: 2})
}

// てっぺき
func IronDefense(battle Battle, _ *rand.Rand) Battle {
	return battle.RankStateFluctuation(&RankState{Def: 2})
}

// めいそう
func CalmMind(battle Battle, _ *rand.Rand) Battle {
	return battle.RankStateFluctuation(&RankState{SpAtk: 1, SpDef: 1})
}

var STATUS_MOVES = map[MoveName]StatusMove{
	"あさのひざし":  Moonlight,
	"こうごうせい":  Synthesis,
	"じこさいせい":  Recover,
	"すなあつめ":   ShoreUp,
	"タマゴうみ":   SoftBoiled,
	"つきのひかり":  Moonlight,
	"なまける":    SlackOff,
	"はねやすめ":   Roost,
	"ミルクのみ":   MorningSun,
	"どくどく":    Toxic,
	"やどりぎのタネ": LeechSeed,
	"つるぎのまい":  SwordsDance,
	"りゅうのまい":  DragonDance,
	"からをやぶる":  ShellSmash,
	"てっぺき":    IronDefense,
}

func init() {
	for moveName, _ := range STATUS_MOVES {
		if _, ok := MOVEDEX[moveName]; !ok {
			errMsg := fmt.Sprintf("STATUS_MOVES の Key に 存在する %v は たぶんタイピングミスってるんで、作者に連絡よろろん twitter @chuusotunamapai", moveName)

			fmt.Println(errMsg)
		}
	}
}
