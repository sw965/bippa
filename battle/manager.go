package battle

import (
	"fmt"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"
	omwmath "github.com/sw965/omw/math"
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

type Manager struct {
	SelfLeadPokemons bp.Pokemons
	SelfBenchPokemons  bp.Pokemons
	OpponentLeadPokemons bp.Pokemons
	OpponentBenchPokemons bp.Pokemons

	Weather Weather
	RemainingTurnWeather int

	SelfFollowMePokemonPointers bp.PokemonPointers
	OpponentFollowMePokemonPointers bp.PokemonPointers

	RemainingTurnTrickRoom int

	IsSingle bool
	Turn int
	IsPlayer1View bool
}

func (m Manager) Clone() Manager {
	return Manager{
		SelfLeadPokemons:m.SelfLeadPokemons.Clone(),
		SelfBenchPokemons:m.SelfBenchPokemons.Clone(),
		OpponentLeadPokemons:m.OpponentLeadPokemons.Clone(),
		OpponentBenchPokemons:m.OpponentBenchPokemons.Clone(),
		Turn:m.Turn,
		Weather:m.Weather,
		RemainingTurnWeather:m.RemainingTurnWeather,
		IsPlayer1View:m.IsPlayer1View,
		IsSingle:m.IsSingle,
	}
}

func (m *Manager) IsTrickRoomState() bool {
	return m.RemainingTurnTrickRoom > 0
}

func (m *Manager) SwapView() {
	m.SelfLeadPokemons, m.SelfBenchPokemons, m.OpponentLeadPokemons, m.OpponentBenchPokemons =
		m.OpponentLeadPokemons, m.OpponentBenchPokemons, m.SelfLeadPokemons, m.SelfBenchPokemons
	m.IsPlayer1View = !m.IsPlayer1View
}

func (m *Manager) CalculateDamage(action *SoloAction, defender *bp.Pokemon, isSingleDmg bool, context *Context) (int, bool, error) {
	attacker := m.SelfLeadPokemons[action.SrcIndex]
	if _, ok := attacker.Moveset[action.MoveName]; !ok {
		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、覚えていない", attacker.Name.ToString(), action.MoveName.ToString())
 		return 0, false, fmt.Errorf(msg)
	}

	switch action.MoveName {
		//がむしゃら
		case bp.ENDEAVOR:
			if slices.Contains(defender.Types, bp.GHOST) {
				return 0, true, nil
			} else {
				dmg := defender.Stat.CurrentHP - attacker.Stat.CurrentHP
				isNoEffect := dmg <= 0
				return omwmath.Max(dmg, 0), isNoEffect, nil
			}
		//ふいうち
		case bp.SUCKER_PUNCH:
			if bp.MOVEDEX[defender.LastPlannedUseMoveName].Category == bp.STATUS {
				return 0, true, nil
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
func (m *Manager) TargetPokemonPointers(action *SoloAction, context *Context) bp.PokemonPointers {
	moveData := bp.MOVEDEX[action.MoveName]

	single := func() bp.PokemonPointers {
		switch moveData.Target {
			case bp.NORMAL_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
			case bp.OPPONENT_TWO_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
			case bp.SELF_TARGET:
				return m.SelfLeadPokemons.ToPointers().NotFainted()
			case bp.OTHERS_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
			case bp.ALL_TARGET:
				return bp.PokemonPointers{}
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
		}
		return bp.PokemonPointers{}
	}

	doubleNormalTarget := func() bp.PokemonPointers {
		followMe := m.OpponentFollowMePokemonPointers.NotFainted()
		if len(followMe) != 0 {
			return bp.PokemonPointers{followMe[0]}
		}

		//味方への攻撃
		if action.IsSelfLeadTarget {
			ally := m.SelfLeadPokemons[action.TargetIndex]
			//攻撃対象の味方が瀕死ならば、ランダムに相手を攻撃する。
			if ally.IsFainted() {
				ps:= m.OpponentLeadPokemons.ToPointers().NotFainted()
				if len(ps) == 0 {
					return bp.PokemonPointers{}
				}
				return omwrand.Sample(ps, 1, context.Rand)
			} else {
				return bp.PokemonPointers{&m.SelfLeadPokemons[action.TargetIndex]}
			}
		} else {
			target := m.OpponentLeadPokemons[action.TargetIndex]
			if target.IsFainted() {
				otherI := map[int]int{0:1, 1:0}[action.TargetIndex]
				other := m.OpponentLeadPokemons[otherI]
				if other.IsFainted() {
					return bp.PokemonPointers{}
				} else {
					return bp.PokemonPointers{&other}
				}						
			} else {
				return bp.PokemonPointers{&target}
			}
		}
	}

	doubleOpponentRandomOneTarget := func() bp.PokemonPointers {
		ps := m.OpponentLeadPokemons.ToPointers().NotFainted()
		if len(ps) == 0 {
			return bp.PokemonPointers{}
		} else {
			return omwrand.Sample(ps, 1, context.Rand)
		}
	}

	double := func() bp.PokemonPointers {
		switch moveData.Target {
			case bp.NORMAL_TARGET:
				return doubleNormalTarget()
			case bp.OPPONENT_TWO_TARGET:
				return m.OpponentLeadPokemons.ToPointers().NotFainted()
			case bp.SELF_TARGET:
				return bp.PokemonPointers{&m.SelfLeadPokemons[action.TargetIndex]}
			case bp.OTHERS_TARGET:
				allyI := map[int]int{0:1, 1:0}[action.SrcIndex]
				ally := m.SelfLeadPokemons[allyI]
				return bp.PokemonPointers{&ally, &m.OpponentLeadPokemons[0], &m.OpponentLeadPokemons[1]}.NotFainted()
			case bp.ALL_TARGET:
				return bp.PokemonPointers{}
			case bp.OPPONENT_RANDOM_ONE_TARGET:
				return doubleOpponentRandomOneTarget()
		}
		return bp.PokemonPointers{}
	}

	if m.IsSingle {
		return single()
	} else {
		return double()
	}
}

func (m *Manager) SelfLegalMoveSoloActions() SoloActions {
	selfLeadN := len(m.SelfLeadPokemons)
	opponentLeadN := len(m.OpponentLeadPokemons)
	actions := make(SoloActions, 0, (selfLeadN * opponentLeadN * bp.MAX_MOVESET_LENGTH) +  (selfLeadN * selfLeadN * bp.MAX_MOVESET_LENGTH))

	selfNotFaintedIdxs := m.SelfLeadPokemons.NotFaintedIndices()
	for _, srcI := range selfNotFaintedIdxs {
		src := m.SelfLeadPokemons[srcI]
		speed := src.Stat.Speed
		for _, usableMoveName := range src.UsableMoveNames() {
			//対象を指定して味方への攻撃や変化技を繰り出す場合
			for _, targetI := range selfNotFaintedIdxs {
				actions = append(actions, SoloAction{
				    MoveName:usableMoveName,
			        SrcIndex:srcI,
				    TargetIndex:targetI,
				    Speed:speed,
				    IsSelfLeadTarget:true,
					IsSelfView:true,
				})
			}

			opponentNotFaintedIdxs := m.OpponentLeadPokemons.NotFaintedIndices()
			//対象を指定して相手への攻撃や変化技を繰り出す場合
			for _, targetI := range opponentNotFaintedIdxs {
				actions = append(actions, SoloAction{
					MoveName:usableMoveName,
					SrcIndex:srcI,
					TargetIndex:targetI,
					Speed:speed,
					IsSelfLeadTarget:false,
					IsSelfView:true,
				})				
			}

			//対象指定なし
			actions = append(actions, SoloAction{
				MoveName:usableMoveName,
				SrcIndex:srcI,
				TargetIndex:-1,
				Speed:speed,
				IsSelfView:true,
			})
		}
	}

	return fn.Filter(actions, func(a SoloAction) bool {
		moveData := bp.MOVEDEX[a.MoveName]
		switch moveData.Target {
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

func (m *Manager) OpponentLegalMoveSoloActions() SoloActions {
	m.SwapView()
	as := m.SelfLegalMoveSoloActions()
	as.ToggleIsSelf()
	m.SwapView()
	return as
}

func (m *Manager) MoveUse(action *SoloAction, context *Context) error {
	pp, ok := m.SelfLeadPokemons[action.SrcIndex].Moveset[action.MoveName]
	if !ok {
		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、覚えていない", m.SelfLeadPokemons[action.SrcIndex].Name.ToString(), action.MoveName.ToString())
		return fmt.Errorf(msg)
	}

	if pp.Current <= 0 {
		msg := fmt.Sprintf("%sは %sを 繰り出そうとしたが、PPが0", m.SelfLeadPokemons[action.SrcIndex].Name.ToString(), action.MoveName.ToString())
		return fmt.Errorf(msg)
	}

	m.SelfLeadPokemons[action.SrcIndex].Moveset[action.MoveName].Current -= 1
	fmt.Println("moveuse", action.MoveName.ToString())
	context.Observer(m, MOVE_USE_EVENT, action.SrcIndex)
	move := GetMove(action.MoveName)
	move.Run(m, action, context)
	return nil
}

func (m *Manager) SelfLegalSwitchSoloActions() SoloActions {
	actions := make(SoloActions, 0, len(m.SelfLeadPokemons) * len(m.SelfBenchPokemons))
	for leadI, leadPokemon := range m.SelfLeadPokemons {
		for _, benchI := range m.SelfBenchPokemons.NotFaintedIndices() {
			actions = append(actions, SoloAction{
				MoveName:bp.EMPTY_MOVE_NAME,
				SrcIndex:leadI,
				TargetIndex:benchI,
				Speed:leadPokemon.Stat.Speed,
				IsSelfView:true,
			})
		}
	}
	return actions
}

func (m *Manager) OpponentLegalSwitchSoloActions() SoloActions {
	m.SwapView()
	as := m.SelfLegalSwitchSoloActions()
	as.ToggleIsSelf()
	m.SwapView()
	return as
}

func (m *Manager) Switch(leadIdx, benchIdx int, context *Context) error {
	if m.SelfBenchPokemons[benchIdx].IsFainted() {
		name := m.SelfBenchPokemons[benchIdx].Name
		msg := fmt.Sprintf("%d番目の %sに 交代しようとしたが、瀕死状態である為、交代出来ません。", benchIdx, name.ToString())
		return fmt.Errorf(msg)
	}

	m.SelfLeadPokemons[leadIdx], m.SelfBenchPokemons[benchIdx] = m.SelfBenchPokemons[benchIdx], m.SelfLeadPokemons[leadIdx]
	context.Observer(m, SWITCH_EVENT, leadIdx)
	return nil
}

func (m *Manager) SoloAction(action *SoloAction, context *Context) error {
	if action.IsMove() {
		return m.MoveUse(action, context)
	} else {
		return m.Switch(action.SrcIndex, action.TargetIndex, context)
	}
}

func (m *Manager) ToEasyRead() EasyReadManager {
	return EasyReadManager{
		SelfLeadPokemons:m.SelfLeadPokemons.ToEasyRead(),
		SelfBenchPokemons:m.SelfBenchPokemons.ToEasyRead(),

		OpponentLeadPokemons:m.OpponentLeadPokemons.ToEasyRead(),
		OpponentBenchPokemons:m.OpponentBenchPokemons.ToEasyRead(),

		Turn:m.Turn,
		IsPlayer1View:m.IsPlayer1View,
		Weather:m.Weather.ToString(),
		RemainingTurnWeather:m.RemainingTurnWeather,
	}
}

type EasyReadManager struct {
	SelfLeadPokemons bp.EasyReadPokemons
	SelfBenchPokemons bp.EasyReadPokemons

	OpponentLeadPokemons bp.EasyReadPokemons
	OpponentBenchPokemons bp.EasyReadPokemons

	Turn int
	IsPlayer1View bool

	Weather string
	RemainingTurnWeather int
}