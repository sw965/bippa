package bippa

import (
	"math/rand"
)

const (
	PHYSICS = "物理"
	SPECIAL = "特殊"
	STATUS  = "変化"
)

type StatusMove func(Battle, *rand.Rand) Battle

//あさのひざし
func NewMorningSun(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//こうごうせい
func NewSynthesis(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//じこさいせい
func NewRecover(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//すなあつめ
func NewShoreUp(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//タマゴうみ
func NewSoftBoiled(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//つきのひかり
func NewMoonlight(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//なまける
func NewSlackOff(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//はねやすめ
func NewRoost(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//ミルクのみ
func NewMilkDrink(battle Battle, _ *rand.Rand) Battle {
	heal := int(float64(battle.P1Fighters[0].MaxHP) * (1.0 / 2.0))
	return battle.Heal(heal)
}

//どくどく
func NewToxic(battle Battle, _ *rand.Rand) Battle {
	if battle.P2Fighters[0].StatusAilment.Type != "" {
		return battle
	}

	if battle.P2Fighters[0].Types.In(POISON) {
		return battle
	}

	if battle.P2Fighters[0].Types.In(STEEL) {
		return battle
	}

	battle.P2Fighters[0].StatusAilment.Type = BAD_POISON
	return battle
}

//やどりぎのタネ
func NewLeechSeed(battle Battle, _ *rand.Rand) Battle {
	if battle.P2Fighters[0].Types.In(GRASS) {
		return battle
	}

	battle.P2Fighters[0].IsLeechSeed = true
	return battle
}

//つるぎのまい
func NewSwordsDance(battle Battle, _ *rand.Rand) Battle {
	return battle.RankFluctuation(&Rank{Atk: 2})
}

//りゅうのまい
func NewDragonDance(battle Battle, _ *rand.Rand) Battle {
	return battle.RankFluctuation(&Rank{Atk: 1, Speed: 1})
}

//からをやぶる
func NewShellSmash(battle Battle, _ *rand.Rand) Battle {
	return battle.RankFluctuation(&Rank{Atk: 2, Def: -1, SpAtk: 2, SpDef: -1, Speed: 2})
}

//てっぺき
func NewIronDefense(battle Battle, _ *rand.Rand) Battle {
	return battle.RankFluctuation(&Rank{Def: 2})
}

//めいそう
func NewCalmMind(battle Battle, _ *rand.Rand) Battle {
	return battle.RankFluctuation(&Rank{SpAtk: 1, SpDef: 1})
}

var STATUS_MOVES = map[MoveName]StatusMove{
	"あさのひざし":  NewMoonlight,
	"こうごうせい":  NewSynthesis,
	"じこさいせい":  NewRecover,
	"すなあつめ":   NewShoreUp,
	"タマゴうみ":   NewSoftBoiled,
	"つきのひかり":  NewMoonlight,
	"なまける":    NewSlackOff,
	"はねやすめ":   NewRoost,
	"ミルクのみ":   NewMorningSun,
	"どくどく":    NewToxic,
	"やどりぎのタネ": NewLeechSeed,
	"つるぎのまい":  NewSwordsDance,
	"りゅうのまい":  NewDragonDance,
	"からをやぶる":  NewShellSmash,
	"てっぺき":    NewIronDefense,
}
