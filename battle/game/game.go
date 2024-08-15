package game

import (
	"fmt"
	//"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/battle"
    bp "github.com/sw965/bippa"
	omwslices "github.com/sw965/omw/slices"
	//"golang.org/x/exp/slices"
)

// func Equal(m1, m2 *battle.Manager) bool {
// 	return m1.SelfLeadPokemons.Equal(m2.SelfLeadPokemons) &&
// 		m1.SelfBenchPokemons.Equal(m2.SelfBenchPokemons) &&
// 		m1.OpponentLeadPokemons.Equal(m2.OpponentLeadPokemons) &&
// 		m1.OpponentBenchPokemons.Equal(m2.OpponentBenchPokemons) &&
// 		m1.Weather == m2.Weather &&
// 		m1.RemainingTurnWeather == m2.RemainingTurnWeather &&
// 		slices.Equal(m1.SelfFollowMePokemonPointers, m2.SelfFollowMePokemonPointers) &&
// 		slices.Equal(m1.OpponentFollowMePokemonPointers, m2.OpponentFollowMePokemonPointers) &&
// 		m1.RemainingTurnTrickRoom == m2.RemainingTurnTrickRoom &&
// 		m1.IsSingle == m2.IsSingle &&
// 		m1.Turn == m2.Turn &&
// 		m1.IsPlayer1View == m2.IsPlayer1View
// }

func IsEnd(m *battle.Manager) (bool, []float64) {
	self := omwslices.Concat(m.CurrentSelfLeadPokemons, m.CurrentSelfBenchPokemons)
	opponent := omwslices.Concat(m.CurrentOpponentLeadPokemons, m.CurrentOpponentBenchPokemons)

	isSelfAllFaint := self.IsAllFainted()
	isOpponentAllFaint := opponent.IsAllFainted()

	if isSelfAllFaint && isOpponentAllFaint {
		return true, []float64{0.5, 0.5}
	} else if isSelfAllFaint {
		return true, []float64{0.0, 1.0}
	} else if isOpponentAllFaint {
		return true, []float64{1.0, 0.0}
	} else {
		return false, []float64{}
	}
}

// func LegalSeparateActions(m *battle.Manager) battle.ActionsSlice {
// 	selfMove := m.CurrentSelfLegalMoveSoloActions()
// 	selfSwitch := m.CurrentSelfLegalSwitchSoloActions()
// 	self := omwslices.Concat(selfMove, selfSwitch)

// 	opponentMove := m.CurrentOpponentLegalMoveSoloActions()
// 	opponentSwitch := m.CurrentOpponentLegalSwitchSoloActions()
// 	opponent := omwslices.Concat(opponentMove, opponentSwitch)

// 	return battle.ActionsSlice{self.ToActions(), opponent.ToActions()}
// }

// func NewPushFunc(context *battle.Context) func(battle.Manager, battle.Actions) (battle.Manager, error) {
// 	return func(m battle.Manager, actions battle.Actions) (battle.Manager, error) {
// 		m = m.Clone()
// 		soloActions := actions.ToSoloActions()
// 		soloActions.SortByOrder(context.Rand)

// 		m.ThisTurnSelfPlannedAction = actions[0]
// 		m.ThisTurnOpponentPlannedAction = actions[1]

// 		for _, soloAction := range soloActions {
// 			if soloAction.IsSelfView {
// 				m.SoloAction(&soloAction, context)
// 			} else {
// 				m.SwapView()
// 				m.SoloAction(&soloAction, context)
// 				m.SwapView()
// 			}
// 		}
// 		m.Turn += 1
// 		return m, nil
// 	}
// }

// func New(context *battle.Context) simultaneous.Game[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action] {
//     gm := simultaneous.Game[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action]{
//         Equal:                Equal,
//         IsEnd:                IsEnd,
//         LegalSeparateActions: LegalSeparateActions,
//         Push:                 NewPushFunc(context),
//     }
//     return gm
// }

//条件
//とりあえず、マネージャーもactionsも正しい前提。
//いずれかのプレイヤーが行動するときに呼び出される関数
func Push(manager battle.Manager, actions battle.Actions) (battle.Manager, error) {
	manager = manager.Clone()
	isSelfLeadAnyFainted := manager.CurrentSelfLeadPokemons.IsAnyFainted()
	isOpponentLeadAnyFainted := manager.CurrentOpponentLeadPokemons.IsAnyFainted()

	soloActions := actions.ToSoloActions()
	soloActions.SortByOrder(battle.GlobalContext.Rand)

	if isSelfLeadAnyFainted || isOpponentLeadAnyFainted {
		for _, soloAction := range soloActions {
			var err error
			if soloAction.IsSelf {
				err = manager.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
			} else {
				manager.SwapView()
				manager.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
				manager.SwapView()
			}
			if err != nil {
				return battle.Manager{}, err
			}
		}
		return battle.Manager{}, nil
	}

	for _, soloAction := range soloActions {
		if !soloAction.IsSelf {
			manager.SwapView()
		}

		srcPokemon := manager.GetHostLeadPokemons()[soloAction.SrcIndex]

		if soloAction.MoveName != bp.EMPTY_MOVE_NAME {
			if srcPokemon.IsFlinchState {
				if soloAction.IsSelf {
					manager.HostViewMessage = fmt.Sprintf("%sは ひるんで わざが だせない！", srcPokemon.Name.ToString())
				} else {
					manager.HostViewMessage = fmt.Sprintf("%sの %sは ひるんで わざが だせない！", manager.GuestHumanName, srcPokemon.Name.ToString())
				}

				battle.GlobalContext.Observer(&manager)
				continue
			}
			move := battle.GetMove(soloAction.MoveName)
			err := move.Run(&manager, &soloAction)
			if err != nil {
				return battle.Manager{}, err
			}
		} else {
			manager.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
		}

		if !soloAction.IsSelf {
			manager.SwapView()
		}
	}
	return manager, nil
}