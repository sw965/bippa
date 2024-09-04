package game

import (
	"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/battle"
    bp "github.com/sw965/bippa"
	omwslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
)

func Equal(m1, m2 *battle.Manager) bool {
	return m1.CurrentSelfLeadPokemons.Equal(m2.CurrentSelfLeadPokemons) &&
		m1.CurrentSelfBenchPokemons.Equal(m2.CurrentSelfBenchPokemons) &&
		m1.CurrentOpponentLeadPokemons.Equal(m2.CurrentOpponentLeadPokemons) &&
		m1.CurrentOpponentBenchPokemons.Equal(m2.CurrentOpponentBenchPokemons) &&

		slices.Equal(m1.CurrentSelfFollowMePokemonPointers, m2.CurrentSelfFollowMePokemonPointers) &&
		slices.Equal(m1.CurrentOpponentFollowMePokemonPointers, m2.CurrentOpponentFollowMePokemonPointers) &&

		m1.Weather == m2.Weather &&
		m1.RemainingTurn == m2.RemainingTurn &&
		m1.Turn == m2.Turn &&

		m1.CurrentSelfIsHost == m2.CurrentSelfIsHost
}

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

func LegalSeparateActions(m *battle.Manager) battle.ActionsSlice {
	self := battle.NewLegalActions(m)
	for i, a := range self {
		for j := range a {
			self[i][j].IsCurrentSelf = true
		} 
	}

	m.SwapView()
	opponent := battle.NewLegalActions(m)
	m.SwapView()
	return battle.ActionsSlice{self, opponent}
}

func Push(m battle.Manager, actions battle.Actions) (battle.Manager, error) {
	m = m.Clone()
	isSelfLeadAnyFainted := m.CurrentSelfLeadPokemons.IsAnyFainted()
	isOpponentLeadAnyFainted := m.CurrentOpponentLeadPokemons.IsAnyFainted()

	soloActions := actions.ToSoloActions()
	soloActions.SortByOrder(&m)

	if isSelfLeadAnyFainted || isOpponentLeadAnyFainted {
		for _, soloAction := range soloActions {
			var err error
			if soloAction.IsCurrentSelf {
				err = m.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
			} else {
				m.SwapView()
				m.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
				m.SwapView()
			}
			if err != nil {
				return battle.Manager{}, err
			}
		}
		return m, nil
	}

	//ふいうちの為の処理
	for _, soloAction := range soloActions {
		if !soloAction.IsMove() {
			continue
		}
		if soloAction.IsCurrentSelf {
			m.CurrentSelfLeadPokemons[soloAction.SrcIndex].ThisTurnPlannedUseMoveName = soloAction.MoveName
		} else {
			m.CurrentOpponentLeadPokemons[soloAction.SrcIndex].ThisTurnPlannedUseMoveName = soloAction.MoveName
		}
	}

	for _, soloAction := range soloActions {
		if !soloAction.IsCurrentSelf {
			m.SwapView()
		}

		if soloAction.MoveName != bp.EMPTY_MOVE_NAME {
			move := battle.GetMove(soloAction.MoveName)
			err := move.Run(&m, &soloAction)
			if err != nil {
				return battle.Manager{}, err
			}
		} else {
			m.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
		}

		if !soloAction.IsCurrentSelf {
			m.SwapView()
		}
	}
	err := m.TurnEnd()
	return m, err
}

func New(context *battle.Context) simultaneous.Game[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action] {
    gm := simultaneous.Game[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action]{
        Equal:                Equal,
        IsEnd:                IsEnd,
        LegalSeparateActions: LegalSeparateActions,
        Push:                 Push,
    }
    return gm
}