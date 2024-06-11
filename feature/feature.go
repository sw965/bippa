package feature

import (
	bp "github.com/sw965/bippa"
	"github.com/sw965/crow/tensor"
	"github.com/sw965/bippa/battle/dmgtools"
	"github.com/sw965/bippa/battle/single"
	omwslices "github.com/sw965/omw/slices"
	omwmath "github.com/sw965/omw/math"
	"github.com/sw965/bippa/team"
	"math"
)

func FirePowerIndex(pokemon *bp.Pokemon) tensor.D1 {
	ret := make(tensor.D1, len(bp.ALL_MOVE_NAMES))
	for i, moveName := range bp.ALL_MOVE_NAMES {
		if _, ok := pokemon.Moveset[moveName]; ok {
			moveData := bp.MOVEDEX[moveName]
			if moveData.Category == bp.PHYSICS {
				ret[i] = float64(moveData.Power * pokemon.Atk) / 100.0
			} else if moveData.Category == bp.SPECIAL {
				ret[i] = float64(moveData.Power * pokemon.SpAtk) / 100.0
			} else {
				ret[i] = 1.0
			}
		}
	}
	return ret
}

func DefenseIndex(pokemon *bp.Pokemon) tensor.D1 {
	pokeTypes := bp.POKEDEX[pokemon.Name].Types
	defIndexFeature := make(tensor.D1, len(bp.ALL_TYPES))
	spDefIndexFeature := make(tensor.D1, len(bp.ALL_TYPES))
	for i, t := range bp.ALL_TYPES {
		effect := dmgtools.Effectiveness(t, pokeTypes)
		defIndexFeature[i] = float64(effect * float64(pokemon.Def)) / 100.0
		spDefIndexFeature[i] = float64(effect * float64(pokemon.SpDef)) / 100.0
	}
	return omwslices.Concat(defIndexFeature, spDefIndexFeature)
}


type TeamFunc func(team.Team) tensor.D1

func NewTeamFunc(capacity int, fs ...func(*bp.Pokemon) tensor.D1) TeamFunc {
	return func(party team.Team) tensor.D1 {
		ret := make(tensor.D1, 0, capacity)
		for _, pokemon := range party {
			feature := make(tensor.D1, 0, n)
			for _, f := range fs {
				feature = append(feature, f(&pokemon)...)
			}
			ret = append(ret, feature...)
		}
		return ret
	}
}

func ExpectedDamageRatioToCurrentHP(p1Pokemon, p2Pokemon *bp.Pokemon) tensor.D1 {
	dmgCalc := dmgtools.Calculator{
		Attacker:dmgtools.Attacker{
			PokeName:p1Pokemon.Name,
			Level:bp.DEFAULT_LEVEL,
			Atk:p1Pokemon.Atk,
			SpAtk:p1Pokemon.SpAtk,
		},

		Defender:dmgtools.Defender{
			PokeName:p2Pokemon.Name,
			Level:bp.DEFAULT_LEVEL,
			Def:p2Pokemon.Def,
			SpDef:p2Pokemon.SpDef,
		},
	}

	ret := make([]float64, 0, bp.MAX_MOVESET)
	for moveName, _ := range p1Pokemon.Moveset {
		dmg := dmgCalc.Expected(moveName)
		accuracy := bp.MOVEDEX[moveName].Accuracy
		expected := dmg / float64(p2Pokemon.CurrentHP) * float64(accuracy) / 100.0
		ret = append(ret, omwmath.Min(expected, 1.0))
	}
	return tensor.D1{omwmath.Max(ret...)}
}

func DPSRatioToCurrentHP(p1Pokemon, p2Pokemon *bp.Pokemon) tensor.D1 {
	dmgCalc := dmgtools.Calculator{
		Attacker:dmgtools.Attacker{
			PokeName:p1Pokemon.Name,
			Level:bp.DEFAULT_LEVEL,
			Atk:p1Pokemon.Atk,
			SpAtk:p1Pokemon.SpAtk,
		},

		Defender:dmgtools.Defender{
			PokeName:p2Pokemon.Name,
			Level:bp.DEFAULT_LEVEL,
			Def:p2Pokemon.Def,
			SpDef:p2Pokemon.SpDef,
		},
	}

	ret := make([]float64, 0, bp.MAX_MOVESET)
	for moveName, _ := range p1Pokemon.Moveset {
		accuracy := bp.MOVEDEX[moveName].Accuracy
		dmg := omwmath.Max(dmgCalc.Expected(moveName) , float64(p2Pokemon.CurrentHP))
		koHit := math.Ceil(float64(p2Pokemon.CurrentHP) / float64(dmg))
		dpsRatio := float64(dmg) / float64(p2Pokemon.CurrentHP) * math.Pow(float64(accuracy) / 100.0, koHit)
		ret = append(ret, dpsRatio)
	}
	return tensor.D1{omwmath.Max(ret...)}
}

type SingleBattleFunc func(*single.Battle) tensor.D1

func NewSingleBattle(capacity int, fs ...func(*bp.Pokemon, *bp.Pokemon) tensor.D1) SingleBattleFunc {
	SPEED_WIN_IDX := 0
	SPEED_LOSS_IDX := 1
	P1_FAINT_IDX := 2
	P2_FAINT_IDX := 3

	return func(battle *single.Battle) tensor.D1 {
		ret := make(tensor.D1, 0, capacity)
		for _, p1Pokemon := range battle.P1Fighters {
			for _, p2Pokemon := range battle.P2Fighters {
				splited := make(tensor.D2, P2_FAINT_IDX+1)
				splited[SPEED_WIN_IDX] = tensor.NewD1Zeros(n*2)
				splited[SPEED_LOSS_IDX] = tensor.NewD1Zeros(n*2)
				splited[P1_FAINT_IDX] = tensor.NewD1Zeros(1)
				splited[P2_FAINT_IDX] = tensor.NewD1Zeros(1)
	
				both := make(tensor.D1 , 0, n*n*4+2)
				for _, f := range fs {
					both = append(both, omwslices.Concat(f(&p1Pokemon, &p2Pokemon), f(&p2Pokemon, &p1Pokemon))...)
				}

				isP1Faint := p1Pokemon.IsFaint()
				isP2Faint := p2Pokemon.IsFaint()
				isNotFaint := !isP1Faint && !isP2Faint
	
				if isNotFaint && p1Pokemon.Speed >= p2Pokemon.Speed {
					splited[SPEED_WIN_IDX] = both
				}
	
				if isNotFaint && p1Pokemon.Speed <= p2Pokemon.Speed {
					splited[SPEED_LOSS_IDX] = both
				}
	
				if isP1Faint {
					splited[P1_FAINT_IDX] = tensor.D1{1}
				}
	
				if isP2Faint {
					splited[P2_FAINT_IDX] = tensor.D1{1}
				}
				for _, v := range splited {
					ret = append(ret, v...)
				}
			}
	
		}
		return ret
	}
}