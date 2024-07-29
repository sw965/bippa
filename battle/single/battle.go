package single

import (
	"fmt"
	"math"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	omwmath "github.com/sw965/omw/math"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/bippa/battle/dmgtools"
	"github.com/sw965/omw/fn"
	"golang.org/x/exp/slices"
)

const (
	PLAYER_NUM = 2
	LEAD_NUM = 1
	BENCH_NUM = 2
	FIGHTERS_NUM = LEAD_NUM + BENCH_NUM
	DOUBLE = 2
)

type Battle struct {
	SelfLeadPokemons bp.Pokemons
	SelfBenchPokemons  bp.Pokemons
	OpponentLeadPokemons bp.Pokemons
	OpponentBenchPokemons bp.Pokemons

	Weather Weather
	RemainingTurnWeather int

	SelfFollowMePokemons PokemonPointers
	OpponentFollowMePokemons PokemonPointers

	IsSingle bool
	Turn int
	IsPlayer1 bool
}

func (b Battle) Clone() Battle {
	return Battle{
		SelfLeadPokemons:b.SelfLeadPokemons.Clone(),
		SelfBenchPokemons:b.SelfBenchPokemons.Clone(),
		OpponentLeadPokemons:b.OpponentLeadPokemons.Clone(),
		OpponentBenchPokemons:b.OpponentBenchPokemons.Clone(),
		Turn:b.Turn,
		Weather:b.Weather,
		RemainingTurnWeather:b.RemainingTurnWeather,
		IsPlayer1:b.IsPlayer1,
		IsSingle:b.IsSingle,
	}
}

func (b *Battle) SwapView() {
	b.SelfLeadPokemons, b.SelfBenchPokemons, b.OpponentLeadPokemons, b.OpponentBenchPokemons =
		b.OpponentLeadPokemons, b.OpponentBenchPokemons, b.SelfLeadPokemons, b.SelfBenchPokemons
	b.IsPlayer1 = !b.IsPlayer1
}

func (b *Battle) CalculateDamage(action *SoloAction, defender *bp.Pokemon, isSingleDmg bool, context *Context) (int, bool, error) {
	attacker := b.SelfLeadPokemons[action.SrcIndex]
	if _, ok := attacker.Moveset[action.MoveName]; !ok {
		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、覚えていない", attacker.Name.ToString(), action.MoveName.ToString())
 		return 0, false, fmt.Errorf(msg)
	}

	switch action.MoveName {
		case bp.ENDEAVOR:
			if slices.Contains(defender.Types, bp.GHOST) {
				return 0, true, nil
			} else {
				dmg := int(math.Abs(float64(defender.CurrentHP) - float64(attacker.CurrentHP)))
				isNoEffect := dmg <= 0
				return omwmath.Max(dmg, 0), isNoEffect, nil
			}
	}

	moveData := bp.MOVEDEX[action.MoveName]
	critRank := moveData.CriticalRank
	isCrit, err := dmgtools.IsCritical(critRank, context.Rand)
	if err != nil {
		return 0, false, err
	}

	calculator := dmgtools.Calculator{
		Attacker:dmgtools.NewAttacker(&attacker),
		Defender:dmgtools.NewDefender(defender),
		IsSingleDamage:isSingleDmg,
		IsCritical:isCrit,
		RandBonus:context.DamageRandBonus(),
	}
	dmg, isNoEffect := calculator.Calculation(action.MoveName)
	return dmg, isNoEffect, nil
}

//攻撃する側のポケモンは瀕死ではない事が前提で呼び出す関数
func (b *Battle) TargetPokemonsIndices(action *SoloAction, context *Context) ([]int, []int) {
	isNotFainted := func(p bp.Pokemon) bool {
		return !p.IsFainted()
	}

	switch action.MoveTarget {
		case bp.NORMAL_TARGET:
			opponentIdxs := b.FollowMeStateOpponentLeadPokemonsIndices()
			if len(opponentIdxs) != 0 {
				return []int{}, opponentIdxs
			}

			if b.IsSingle {
				target := b.OpponentLeadPokemons[action.TargetIndex]
				if target.IsFainted() {
					return []int{}, []int{}
				} else {
					return []int{}, []int{action.TargetIndex}
				}
			} else {
				//味方への攻撃の場合
				if action.IsSelfLeadTarget {
					ally := b.SelfLeadPokemons[action.TargetIndex]
					//攻撃対象の味方のポケモンが瀕死ならば
					if ally.IsFainted() {
						//攻撃対象が相手のポケモンになる。
						notFaintedIdxs := omwslices.IndicesFunc[bp.Pokemons](b.OpponentLeadPokemons, isNotFainted)
						if len(notFaintedIdxs) == 0 {
							return []int{}, []int{}
						} else {
							//攻撃対象になる相手のポケモンはランダム
							return []int{}, omwrand.Sample(notFaintedIdxs, 1, context.Rand)
						}
					} else {
						return []int{action.TargetIndex}, []int{}
					}
				} else {
					target := b.OpponentLeadPokemons[action.TargetIndex]
					if target.IsFainted() {
						otherIdx := map[int]int{0:1, 1:0}[action.TargetIndex]
						other := b.OpponentLeadPokemons[otherIdx]
						if other.IsFainted() {
							return []int{}, []int{}
						} else {
							return []int{}, []int{otherIdx}
						}
					} else {
						return []int{}, []int{action.TargetIndex}
					}
				}
			}
		case bp.OPPONENT_TWO_TARGET:
			if b.IsSingle {
				target := b.OpponentLeadPokemons[0]
				if target.IsFainted() {
					return []int{}, []int{}
				} else {
					return []int{0}, []int{}
				}
			} else {
				idxs := omwslices.IndicesFunc[bp.Pokemons](b.OpponentLeadPokemons, isNotFainted)
				return []int{}, idxs
			}
		case bp.SELF_TARGET:
			return []int{action.TargetIndex}, []int{}
		case bp.OTHERS_TARGET:
			if b.IsSingle {
				target := b.OpponentLeadPokemons[0]
				if target.IsFainted() {
					return []int{}, []int{}
				} else {
					return []int{0}, []int{}
				}
			} else {
				allyIdx := map[int]int{0:1, 1:0}[action.SrcIndex]
				ally := b.SelfLeadPokemons[allyIdx]
				var allyIdxs []int
				if ally.IsFainted() {
					allyIdxs = []int{}					
				} else {
					allyIdxs = []int{allyIdx}
				}
				opponentIdxs := omwslices.IndicesFunc[bp.Pokemons](b.OpponentLeadPokemons, isNotFainted)
				return allyIdxs, opponentIdxs
			}
		case bp.ALL_TARGET:
			return []int{}, []int{}
		case bp.OPPONENT_RANDOM_ONE_TARGET:
			opponentIdxs := b.FollowMeStateOpponentLeadPokemonsIndices()
			if len(opponentIdxs) != 0 {
				return []int{}, opponentIdxs
			}

			if b.IsSingle {
				target := b.OpponentLeadPokemons[0]
				if target.IsFainted() {
					return []int{}, []int{}
				} else {
					return []int{0}, []int{}
				}
			} else {
				idxs := omwslices.IndicesFunc[bp.Pokemons](b.OpponentBenchPokemons, isNotFainted)
				if len(idxs) == 0 {
					return []int{}, []int{}
				} else {
					return []int{}, omwrand.Sample(idxs, 1, context.Rand)
				}
			}
		default:
			return []int{}, []int{}
	}
}

func (b *Battle) TargetPokemonPointers(action *SoloAction, context *Context) bp.PokemonPointers {
	selfIdxs, opponentIdxs := b.TargetPokemonsIndices(action, context)

	self := make(bp.PokemonPointers, len(selfIdxs))
	for i, idx := range selfIdxs {
		self[i] = &b.SelfLeadPokemons[idx]
	}

	opponent := make(bp.PokemonPointers, len(opponentIdxs))
	for i, idx := range opponentIdxs {
		opponent[i] = &b.OpponentLeadPokemons[idx]
	}
	
	return omwslices.Concat(self, opponent)
}

func (b *Battle) SelfLegalMoveSoloActions() SoloActions {
	selfLeadN := len(b.SelfLeadPokemons)
	opponentLeadN := len(b.OpponentLeadPokemons)
	actions := make(SoloActions, 0, (selfLeadN * opponentLeadN * bp.MAX_MOVESET) +  (selfLeadN * selfLeadN * bp.MAX_MOVESET))

	selfNotFaintedIdxs := b.SelfLeadPokemons.NotFaintedIndices()
	for _, srcI := range selfBotFaintedIdxs {
		src := b.SelfLeadPokemons[idx]
		speed := src.Speed
		moveset := src.Moveset

		for _, usableMoveName := range src.UsableMoveNames() {
			moveData := bp.MOVEDEX[usableMoveName]

			//対象を指定して味方への攻撃や変化技を繰り出す場合
			for _, targetI := range notFaintedIdxs {
				target := b.SelfLeadPokemons[targetI]
				actions = append(actions, SoloAction{
				    MoveName:moveName,
			        SrcIndex:srcI,
				    TargetIndex:targetI,
				    Speed:speed,
				    IsSelfLeadTarget:true,
				    MoveTarget:moveData.Target,
					IsSelfView:true,
				})
			}

			opponentNotFaintedIdxs := b.OpponentLeadPokemons.NotFaintedIndices()
			//対象を指定して相手への攻撃や変化技を繰り出す場合
			for _, targetI := range opponentNotFaintedIdxs {
				actions = append(actions, SoloAction{
					MoveName:moveName,
					SrcIndex:srcI,
					TargetIndex:targetI,
					Speed:speed,
					IsSelfLeadTarget:false,
					MoveTarget:moveData.Target,
					IsSelfView:true,
				})				
			}

			//対象指定なし
			actions = append(actions, SoloAction{
				MoveName:moveName,
				SrcIndex:i,
				TargetIndex:-1,
				Speed:speed,
				MoveTarget:moveData.Target,
				IsSelfView:true,
			})
	}

	return fn.Filter(actions, func(a SoloAction) bool {
		switch a.MoveTarget {
			case bp.NORMAL_TARGET:
				if a.TargetIndex == -1 {
					return false
				} else if a.IsSelfLeadTarget {
					return a.SrcIndex != a.TargetIndex
				} else {
					return true
				}
			case bp.OPPONENT_TWO_TARGET:
				return a.TargetIndex == -1
			case bp.SELF_TARGET:
				return a.IsSelfLeadTarget && a.SrcIndex == a.TargetIndex
			case bp.OTHERS_TARGET:
				return a.TargetIndex == -1
			case bp.ALL_TARGET:
				return a.TargetIndex == -1
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				return a.TargetIndex == -1
			default:
				return true
		}
	})
}

func (b *Battle) OpponentLegalMoveSoloActions() SoloActions {
	b.SwapView()
	ret := b.SelfLegalMoveSoloActions()
	ret.ToggleIsSelf()
	b.SwapView()
	return ret
}

func (b *Battle) LegalSeparateMoveSoloActionsSlice() SoloActionsSlice {
	self := b.SelfLegalMoveSoloActions()
	opponent := b.OpponentLegalMoveSoloActions()
	return SoloActionsSlice{self, opponent}
}

func (b *Battle) MoveUse(action *SoloAction, context *Context) error {
	pp, ok := b.SelfLeadPokemons[action.SrcIndex].Moveset[action.MoveName]
	if !ok {
		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、覚えていない", b.SelfLeadPokemons[action.SrcIndex].Name.ToString(), action.MoveName.ToString())
		return fmt.Errorf(msg)
	}

	if pp.Current <= 0 {
		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、PPが0", b.SelfLeadPokemons[action.SrcIndex].Name.ToString(), action.MoveName.ToString())
		return fmt.Errorf(msg)
	}

	b.SelfLeadPokemons[action.SrcIndex].Moveset[action.MoveName].Current -= 1
	return GetMoveFunc(action.MoveName)(b, action, context)
}

func (b *Battle) LegalSwitchSoloActions() SoloActions {
	actions := make(SoloActions, 0, len(b.SelfLeadPokemons) * len(b.SelfBenchPokemons))
	for i, leadPokemon := range b.SelfLeadPokemons {
		for j := range b.SelfBenchPokemons {
			actions = append(actions, SoloAction{
				MoveName:bp.EMPTY_MOVE_NAME,
				SrcIndex:i,
				TargetIndex:j,
				Speed:leadPokemon.Speed,
				IsSelfView:true,
			})
		}
	}
	return fn.Filter(actions, func(action SoloAction) bool {
		return !b.SelfBenchPokemons[action.TargetIndex].IsFainted()
	})
}

func (b *Battle) Switch(leadIdx, benchIdx int, context *Context) error {
	if b.SelfBenchPokemons[benchIdx].IsFainted() {
		name := b.SelfBenchPokemons[benchIdx].Name
		msg := fmt.Sprintf("%d番目の %sに 交代しようとしたが、瀕死状態である為、交代出来ません。", benchIdx, name.ToString())
		return fmt.Errorf(msg)
	}

	b.SelfLeadPokemons[leadIdx], b.SelfBenchPokemons[benchIdx] = b.SelfBenchPokemons[benchIdx], b.SelfLeadPokemons[leadIdx]
	context.Observer(b, SWITCH_EVENT)
	return nil
}

func (b *Battle) SoloAction(action *SoloAction, context *Context) error {
	if action.IsMove() {
		return b.MoveUse(action, context)
	} else {
		return b.Switch(action.SrcIndex, action.TargetIndex, context)
	}
}

func (b *Battle) ToEasyRead() EasyReadBattle {
	return EasyReadBattle{
		SelfLeadPokemons:b.SelfLeadPokemons.ToEasyRead(),
		SelfBenchPokemons:b.SelfBenchPokemons.ToEasyRead(),

		OpponentLeadPokemons:b.OpponentLeadPokemons.ToEasyRead(),
		OpponentBenchPokemons:b.OpponentBenchPokemons.ToEasyRead(),

		Turn:b.Turn,
		IsPlayer1:b.IsPlayer1,
		Weather:b.Weather.ToString(),
		RemainingTurnWeather:b.RemainingTurnWeather,
	}
}