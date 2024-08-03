package game

import (
	"fmt"
	"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/battle"
    //bp "github.com/sw965/bippa"
	omwslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
)

func Equal(m1, m2 *battle.Manager) bool {
	return m1.SelfLeadPokemons.Equal(m2.SelfLeadPokemons) &&
		m1.SelfBenchPokemons.Equal(m2.SelfBenchPokemons) &&
		m1.OpponentLeadPokemons.Equal(m2.OpponentLeadPokemons) &&
		m1.OpponentBenchPokemons.Equal(m2.OpponentBenchPokemons) &&
		m1.Weather == m2.Weather &&
		m1.RemainingTurnWeather == m2.RemainingTurnWeather &&
		slices.Equal(m1.SelfFollowMePokemonPointers, m2.SelfFollowMePokemonPointers) &&
		slices.Equal(m1.OpponentFollowMePokemonPointers, m2.OpponentFollowMePokemonPointers) &&
		m1.RemainingTurnTrickRoom == m2.RemainingTurnTrickRoom &&
		m1.IsSingle == m2.IsSingle &&
		m1.Turn == m2.Turn &&
		m1.IsPlayer1View == m2.IsPlayer1View
}

func IsEnd(m *battle.Manager) (bool, []float64) {
	self := omwslices.Concat(m.SelfLeadPokemons, m.SelfBenchPokemons)
	opponent := omwslices.Concat(m.OpponentLeadPokemons, m.OpponentBenchPokemons)

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

func LegalSeparateActions(m *battle.Manager) battle.ActionsSlice {
	selfMove := m.SelfLegalMoveSoloActions()
	selfSwitch := m.SelfLegalSwitchSoloActions()
	self := omwslices.Concat(selfMove, selfSwitch)

	opponentMove := m.OpponentLegalMoveSoloActions()
	opponentSwitch := m.OpponentLegalSwitchSoloActions()
	opponent := omwslices.Concat(opponentMove, opponentSwitch)

	return battle.ActionsSlice{self.ToActions(), opponent.ToActions()}
}

func NewPushFunc(context *battle.Context) func(battle.Manager, battle.Actions) (battle.Manager, error) {
	return func(m battle.Manager, actions battle.Actions) (battle.Manager, error) {
		m = m.Clone()
		soloActions := actions.ToSoloActions()
		fmt.Println(soloActions)
		soloActions.SortByOrder(context.Rand)

		for _, soloAction := range soloActions {
			fmt.Println(soloAction.MoveName.ToString())
			if soloAction.IsSelfView {
				m.SoloAction(&soloAction, context)
			} else {
				m.SwapView()
				m.SoloAction(&soloAction, context)
				m.SwapView()
			}
		}
		m.Turn += 1
		return m, nil
	}
}

func New(context *battle.Context) simultaneous.Game[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action] {
    gm := simultaneous.Game[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action]{
        Equal:                Equal,
        IsEnd:                IsEnd,
        LegalSeparateActions: LegalSeparateActions,
        Push:                 NewPushFunc(context),
    }
    return gm
}