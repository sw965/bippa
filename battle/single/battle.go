package single

import (
	"fmt"
	//"math"
	"math/rand"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	//omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/bippa/battle/dmgtools"
	"github.com/sw965/omw/fn"
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

func (b *Battle) TargetPokemonPointers(action *SoloAction, r *rand.Rand) (bp.PokemonPointers, error) {
	isNotFainted := func(pokemon *bp.Pokemon) bool {
		return !pokemon.IsFainted()
	}

	selfLeadN := len(b.SelfLeadPokemons)
	opponentLeadN := len(b.OpponentLeadPokemons)
	selfLeadPokemons := b.SelfLeadPokemons.ToPointers()
	opponentLeadPokemons := b.OpponentLeadPokemons.ToPointers()
	var targetPokemons bp.PokemonPointers

	switch action.Target {
		case bp.NORMAL_TARGET:
			f := func(ps bp.PokemonPointers, n int) {
				if ps[action.TargetIndex].IsFainted() {
					if n == 1 {
						return
					} else {
						targetIdx := map[int]int{0:1, 1:0}[action.TargetIndex]
						targetPokemons = bp.PokemonPointers{ps[targetIdx]}
					}
				} else {
					targetPokemons = bp.PokemonPointers{ps[action.TargetIndex]}
				}
			}

			if action.IsSelfLeadTarget {
				f(selfLeadPokemons, selfLeadN)
			} else if action.IsOpponentLeadTarget {
				f(opponentLeadPokemons, opponentLeadN)
			} else {
				return targetPokemons, fmt.Errorf("SoloAction.Target == NORMAL_TARGET である場合、SoloAction.IsOpponentLeadTarget || SoloAction.IsSelfLeadTarget でなければならない。")
			}
			targetPokemons = fn.Filter(targetPokemons, isNotFainted)
		case bp.OPPONENT_TWO_TARGET:
			targetPokemons = opponentLeadPokemons
			targetPokemons = fn.Filter(targetPokemons, isNotFainted)
		case bp.SELF_TARGET:
			targetPokemons = bp.PokemonPointers{selfLeadPokemons[action.TargetIndex]}
			targetPokemons = fn.Filter(targetPokemons, isNotFainted)
		case bp.OTHERS_TARGET:
			var allyPokemons bp.PokemonPointers
			if len(b.SelfLeadPokemons) == 1 {
				allyPokemons = bp.PokemonPointers{}
			} else {
				targetIdx := map[int]int{0:1, 1:0}[action.SrcIndex]
				allyPokemons = bp.PokemonPointers{selfLeadPokemons[targetIdx]}
			}
			targetPokemons = omwslices.Concat(allyPokemons, opponentLeadPokemons)
			targetPokemons = fn.Filter(targetPokemons, isNotFainted)
		case bp.ALL_TARGET:
			targetPokemons = omwslices.Concat(selfLeadPokemons, opponentLeadPokemons)
			targetPokemons = fn.Filter(targetPokemons, isNotFainted)
		case bp.OPPONENT_RANDOM_ONE_TARGET:
			targetPokemons = fn.Filter(opponentLeadPokemons, isNotFainted)
			if len(targetPokemons) != 0 {
				targetPokemons = bp.PokemonPointers{omwrand.Choice(targetPokemons, r)}
			}
	}
	targetPokemons.SortBySpeed()
	return targetPokemons, nil
}

func (b *Battle) CalcDamage(action *SoloAction, defender *bp.Pokemon, isCrit, isSingleDmg bool, context *Context) (int, error) {
	attacker := b.SelfLeadPokemons[action.SrcIndex]
	if _, ok := attacker.Moveset[action.MoveName]; !ok {
		msg := fmt.Sprintf("%s は 覚えていないのに、%sを繰り出そうとした", attacker.Name.ToString(), action.MoveName.ToString())
		return 0, fmt.Errorf(msg)
	}

	calculator := dmgtools.Calculator{
		Attacker:dmgtools.Attacker{
			PokeName:attacker.Name,
			Level:attacker.Level,
			Atk:attacker.Atk,
			AtkRank:attacker.Rank.Atk,
			SpAtk:attacker.SpAtk,
			SpAtkRank:attacker.Rank.SpAtk,
		},

		Defender:dmgtools.Defender{
			PokeName:defender.Name,
			Level:defender.Level,
			Def:defender.Def,
			DefRank:defender.Rank.Def,
			SpDef:defender.SpDef,
			SpDefRank:defender.Rank.SpDef,
		},
		IsCritical:isCrit,
		IsSingleDamage:isSingleDmg,
	}
	return calculator.Calculation(action.MoveName, context.DamageRandBonus()), nil
}

func (b *Battle) LegalMoveSoloActions() SoloActions {
	selfLeadN := len(b.SelfLeadPokemons)
	opponentLeadN := len(b.OpponentLeadPokemons)
	actions := make(SoloActions, 0, (selfLeadN * opponentLeadN * bp.MAX_MOVESET) +  (selfLeadN * selfLeadN * bp.MAX_MOVESET))
	for i, selfPokemon := range b.SelfLeadPokemons {
		if selfPokemon.IsFainted() {
			continue
		}
		speed := selfPokemon.Speed
		moveset := selfPokemon.Moveset
		for moveName, pp := range moveset {
			if pp.Current <= 0 {
				continue
			}
			moveData := bp.MOVEDEX[moveName]

			for j, opponentPokemon := range b.OpponentLeadPokemons {
				if opponentPokemon.IsFainted() {
					continue
				}
				//場に出ている敵への対象指定
				actions = append(actions, SoloAction{
					MoveName:moveName,
					SrcIndex:i,
					TargetIndex:j,
					Speed:speed,
					IsOpponentLeadTarget:true,
					Target:moveData.Target,
				})
			}

			for j, selfTargetPokemon := range b.SelfLeadPokemons {
				if selfTargetPokemon.IsFainted() {
					continue
				}
				//場に出ている味方への対象指定
			    actions = append(actions, SoloAction{
				    MoveName:moveName,
			        SrcIndex:i,
				    TargetIndex:j,
				    Speed:speed,
				    IsSelfLeadTarget:true,
				    Target:moveData.Target,
				})
			}

			//対象指定なし
			actions = append(actions, SoloAction{
				MoveName:moveName,
				SrcIndex:i,
				TargetIndex:-1,
				Speed:speed,
				Target:moveData.Target,
			})
		}
	}
	
	return fn.Filter(actions, func(action SoloAction) bool {
		switch action.Target {
			case bp.NORMAL_TARGET:
				if action.IsOpponentLeadTarget {
					return true
				} else if action.IsSelfLeadTarget {
					return action.SrcIndex != action.TargetIndex
				} else {
					return false
				}
			case bp.OPPONENT_TWO_TARGET:
				return action.TargetIndex == -1
			case bp.SELF_TARGET:
				return action.IsSelfLeadTarget && action.SrcIndex == action.TargetIndex
			case bp.OTHERS_TARGET:
				return action.TargetIndex == -1
			case bp.ALL_TARGET:
				return action.TargetIndex == -1
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				return action.TargetIndex == -1
		}
		return true
	})
}

// func (b Battle) UseMove(action *SoloAction, context *Context) (Battle, error) {
// 	pp, ok := b.SelfLeadPokemons[action.SrcIndex].Moveset[moveName]
// 	if !ok {
// 		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、覚えていない", b.SelfLeadPokemons[action.SrcIndex].Name.ToString(), moveName.ToString())
// 		return b, fmt.Errorf(msg)
// 	}

// 	if pp.Current <= 0 {
// 		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、PPが0", b.SelfLeadPokemons[action.SrcIndex].Name.ToString(), moveName.ToString())
// 		return b, fmt.Errorf(msg)
// 	}

// 	if b.SelfLeadPokemons[action.SrcIndex].IsFainted() {
// 		return b, nil
// 	}

// 	b.SelfLeadPokemons = b.SelfLeadPokemons.Clone()
// 	b.OpponentLeadPokemons = b.OpponentLeadPokemons.Clone()
//  ここに挑発時の時の処理をする。
// 	b.SelfLeadPokemons[selfLeadIdx].Moveset[moveName].Current -= 1
// 	context.Observer(&b, MOVE_USE_EVENT)

// 	attack := func(attacker, defender *bp.Pokemon, isDoubleDamage bool) {
// 		dmg := calcDmg(attacker, defender isDoubleDamage, context)
// 		dmg = omwmath.Min(dmg, defender.CurrentHP)
// 		defender.CurrentHP = dmg
// 	}

// 	attacker := &b.SelfLeadPokemons[action.SrcIndex]
// 	normalTarget := func() {
// 		attack(attacker, b.OpponentLeadPokemons[action.TargetIndex])
// 	}

// 	opponentTwoTarget := func() {
// 		defenders := b.OpponentBenchPokemons.SortBySpeed()
// 		for _, defender := range defenders {
// 			attack(attacker, defender, isDoubleDamage)
// 		}
// 	}

// 	selfTarget := func() {
// 	}

// 	if moveName == bp.STRUGGLE {
// 		dmg := int(math.Round(float64(b.SelfLeadPokemons[0].CurrentHP) * 0.25))
// 		b.SelfLeadPokemons[0].CurrentHP -= dmg
// 		context.Observer(&b, RECOIL_EVENT)
// 	}
// 	return b, nil
// }

func (b *Battle) LegalSwitchSoloActions() SoloActions {
	actions := make(SoloActions, 0, len(b.SelfLeadPokemons) * len(b.SelfBenchPokemons))
	for i, leadPokemon := range b.SelfLeadPokemons {
		for j := range b.SelfBenchPokemons {
			actions = append(actions, SoloAction{
				MoveName:bp.EMPTY_MOVE_NAME,
				SrcIndex:i,
				TargetIndex:j,
				Speed:leadPokemon.Speed,
				IsSelfBenchTarget:true,
			})
		}
	}
	return fn.Filter(actions, func(action SoloAction) bool {
		return !b.SelfBenchPokemons[action.TargetIndex].IsFainted()
	})
}

// func (b Battle) Switch(leadIdx, benchIdx int, context *Context) (Battle, error) {
// ランクをリセットする。
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