package single

import (
	//"fmt"
	//"math"
	//"math/rand"
	bp "github.com/sw965/bippa"
	//omwrand "github.com/sw965/omw/math/rand"
	//omwmath "github.com/sw965/omw/math"
	"github.com/sw965/bippa/battle/dmgtools"
)

const (
	PLAYER_NUM = 2
	LEAD_NUM = 1
	BENCH_NUM = 2
	FIGHTERS_NUM = LEAD_NUM + BENCH_NUM
)

type Battle struct {
	SelfLeadPokemons bp.Pokemons
	SelfBenchPokemons  bp.Pokemons
	OpponentLeadPokemons bp.Pokemons
	OpponentBenchPokemons bp.Pokemons
	Turn int
	IsRealSelfView bool
}

func (b Battle) Clone() Battle {
	return Battle{
		SelfLeadPokemons:b.SelfLeadPokemons.Clone(),
		SelfBenchPokemons:b.SelfBenchPokemons.Clone(),
		OpponentLeadPokemons:b.OpponentLeadPokemons.Clone(),
		OpponentBenchPokemons:b.OpponentBenchPokemons.Clone(),
		Turn:b.Turn,
		IsRealSelfView:b.IsRealSelfView,
	}
}

func (b Battle) SwapView() Battle {
	b.SelfLeadPokemons, b.SelfBenchPokemons, b.OpponentLeadPokemons, b.OpponentBenchPokemons =
		b.OpponentLeadPokemons, b.OpponentBenchPokemons, b.SelfLeadPokemons, b.SelfBenchPokemons
	b.IsRealSelfView = !b.IsRealSelfView
	return b
}

func (b *Battle) CalcDamage(moveName bp.MoveName, selfLeadIdx, opponentLeadIdx int, context *Context) int {
	attacker := b.SelfLeadPokemons[selfLeadIdx]
	defender := b.OpponentLeadPokemons[opponentLeadIdx]
	calculator := dmgtools.Calculator{
		dmgtools.Attacker{
			PokeName:attacker.Name,
			Level:attacker.Level,
			Atk:attacker.Atk,
			SpAtk:attacker.SpAtk,
		},
		dmgtools.Defender{
			PokeName:defender.Name,
			Level:defender.Level,
			Def:defender.Def,
			SpDef:defender.SpDef,
		},
	}
	return calculator.Calculation(moveName, context.DamageRandBonus())
}

// func (b Battle) UseMove(moveName bp.MoveName, selfLeadIdx, opponentLeadIdx int, context *Context) (Battle, error) {
// 	pp, ok := b.SelfLeadPokemons[selfLeadIdx].Moveset[moveName]
// 	if !ok {
// 		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、覚えていない", b.SelfLeadPokemons[selfLeadIdx].Name.ToString(), moveName.ToString())
// 		return b, fmt.Errorf(msg)
// 	}

// 	if pp.Current <= 0 {
// 		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、PPが0", b.SelfLeadPokemons[selfLeadIdx].Name.ToString(), moveName.ToString())
// 		return b, fmt.Errorf(msg)
// 	}

// 	if b.SelfLeadPokemons[selfLeadIdx].IsFainted() {
// 		return b, nil
// 	}

// 	var opponentLeadIdxs []int
// 	if opponentLeadIdx == -1 {
// 		opponentLeadIdxs = omwslices.MakeInteger(0, len(b.OpponentBenchPokemons))
// 	} else {
// 		opponentLeadIdxs = []int{opponentLeadIdx}
// 	}

// 	b.SelfLeadPokemons = b.SelfLeadPokemons.Clone()
// 	b.OpponentLeadPokemons = b.OpponentLeadPokemons.Clone()
	
// 	b.SelfLeadPokemons[selfLeadIdx].Moveset[moveName].Current -= 1
// 	context.Observer(&b, MOVE_USE_EVENT)

// 	for _, idx := range opponentLeadIdxs {
// 		dmg := b.CalcDamage(moveName, selfLeadIdx, idx, context)
// 		dmg = omwmath.Min(dmg, b.OpponentLeadPokemons[idx].CurrentHP)
// 	}

// 	context.Observer(&b, ATTACK_MOVE_DAMAGE_EVENT)

// 	if moveName == bp.STRUGGLE {
// 		dmg := int(math.Round(float64(b.SelfLeadPokemons[0].CurrentHP) * 0.25))
// 		b.SelfLeadPokemons[0].CurrentHP -= dmg
// 		context.Observer(&b, RECOIL_EVENT)
// 	}
// 	return b, nil
// }

// func (b Battle) Switch(leadIdx, benchIdx int, context *Context) (Battle, error) {
// 	if b.SelfBenchPokemons[benchIdx].IsFainted() {
// 		name := b.SelfBenchPokemons[benchIdx].Name
// 		msg := fmt.Sprintf("%d番目の %sに 交代しようとしたが、瀕死状態である為、交代出来ません。", benchIdx, name)
// 		return b, fmt.Errorf(msg)
// 	}

// 	selfLeadPokemons := b.SelfLeadPokemons.Clone()
// 	selfBenchPokemons := b.SelfBenchPokemons.Clone()

// 	selfLeadPokemons[leadIdx], selfBenchPokemons[benchIdx] = selfBenchPokemons[benchIdx], selfLeadPokemons[leadIdx]
// 	b.SelfLeadPokemons = selfLeadPokemons
// 	b.SelfBenchPokemons = selfBenchPokemons
// 	context.Observer(&b, SWITCH_EVENT)
// 	return b, nil
// }

// func (b *Battle) Action(action Action, context *Context) (Battle, error) {
// 	if action.IsCommandMove() {
// 		return b.CommandMove(action.CmdMoveName, context)
// 	} else {
// 		return b.Switch(action.SwitchPokeName, context)
// 	}
// }

// func (b *Battle) SortActionsByOrder(actions Actions, r *rand.Rand) Actions {
// 	actions = actions.Clone()
// 	slices.SortFunc(actions, func(a1, a2 Action) {
// 		a1Priority := a1.Priority()
// 		a2Priority := a2.Priority()

// 		if a1Priority > a2Priority {
// 			return 1
// 		} else if a1Priority < a2Priority {
// 			return -1
// 		} else {
// 			a1Speed := a1.Speed
// 			a2Speed := a2.Speed
// 			if a1.Speed > a2.Speed {
// 				return 1
// 			} else if a1.Speed < a2.Speed {
// 				return 1
// 			} else {
// 				return omwrand.Choice([]int{-1, 1})
// 			}
// 		}
// 	})
// 	return actions
// }

// func (b *Battle) ToEasyRead() EasyReadBattle {
// 	return EasyReadBattle{
// 		SelfLeadPokemons:b.SelfLeadPokemons.ToEasyRead(),
// 		SelfBenchPokemons:b.SelfBenchPokemons.ToEasyRead(),

// 		OpponentLeadPokemons:b.OpponentFighters.ToEasyRead(),
// 		OpponentBenchPokemons:b.OpponentBenchPokemons.ToEasyRead(),

// 		Turn:b.Turn,
// 		IsRealSelfView:b.IsRealSelfView,
// 	}
// }