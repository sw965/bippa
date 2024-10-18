package game

import (
	"fmt"
	"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/battle"
    bp "github.com/sw965/bippa"
	omwslices "github.com/sw965/omw/slices"
	"golang.org/x/exp/slices"
)

func SeparateLegalActionsProvider(m *battle.Manager) battle.ActionsSlice {
	self := battle.NewLegalActions(m)
	if len(self) == 0 {
		self = battle.Actions{battle.Action{}}
	}

	for i, a := range self {
		for j := range a {
			self[i][j].IsCurrentSelf = true
		} 
	}

	m.SwapView()
	opponent := battle.NewLegalActions(m)
	m.SwapView()
	if len(opponent) == 0 {
		opponent = battle.Actions{battle.Action{}}
	}
	return battle.ActionsSlice{self, opponent}
}

func Transitioner(m battle.Manager, actions battle.Actions) (battle.Manager, error) {
	m = m.Clone()
	actions = actions.FilterByNotEmpty()
	if len(actions) == 0 {
		return battle.Manager{}, fmt.Errorf("両プレイヤーのActionが、Action.IsEmpty() == true である為、処理を続行出来ません。")
	}

	soloActions := actions.ToSoloActions()
	soloActions = soloActions.FilterByNotEmpty()
	soloActions.SortByOrder(&m)

	selfMustSwitch, opponentMustSwitch := m.MustSwitch()

	if selfMustSwitch || opponentMustSwitch {
		for _, soloAction := range soloActions {
			var err error
			if soloAction.IsCurrentSelf {
				err = m.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
			} else {
				m.SwapView()
				err = m.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
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

		var err error
		if soloAction.MoveName != bp.EMPTY_MOVE_NAME {
			move := battle.GetMove(soloAction.MoveName)
			err = move.Run(&m, &soloAction)
		} else {
			err = m.Switch(soloAction.SrcIndex, soloAction.TargetIndex)
		}

		if err != nil {
			return battle.Manager{}, err
		}

		if !soloAction.IsCurrentSelf {
			m.SwapView()
		}
	}
	err := m.TurnEnd()
	return m, err
}

func Comparator(m1, m2 *battle.Manager) bool {
	return m1.CurrentSelfLeadPokemons.Equal(m2.CurrentSelfLeadPokemons) &&
		m1.CurrentSelfBenchPokemons.Equal(m2.CurrentSelfBenchPokemons) &&
		m1.CurrentOpponentLeadPokemons.Equal(m2.CurrentOpponentLeadPokemons) &&
		m1.CurrentOpponentBenchPokemons.Equal(m2.CurrentOpponentBenchPokemons) &&

		slices.Equal(m1.CurrentSelfAttentionPokemonPointers, m2.CurrentSelfAttentionPokemonPointers) &&
		slices.Equal(m1.CurrentOpponentAttentionPokemonPointers, m2.CurrentOpponentAttentionPokemonPointers) &&

		m1.Weather == m2.Weather &&
		m1.RemainingTurn == m2.RemainingTurn &&
		m1.Turn == m2.Turn &&

		m1.CurrentSelfIsHost == m2.CurrentSelfIsHost
}

func EndChecker(m *battle.Manager) bool {
	self := omwslices.Concat(m.CurrentSelfLeadPokemons, m.CurrentSelfBenchPokemons)
	opponent := omwslices.Concat(m.CurrentOpponentLeadPokemons, m.CurrentOpponentBenchPokemons)
	isSelfAllFainted := self.IsAllFainted()
	isOpponentAllFainted := opponent.IsAllFainted()
	return isSelfAllFainted || isOpponentAllFainted
}

func NewLogic() simultaneous.Logic[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action] {
    g := simultaneous.Logic[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action]{
		SeparateLegalActionsProvider: SeparateLegalActionsProvider,
        Transitioner:                 Transitioner,
        Comparator:                Comparator,
        EndChecker:                EndChecker,
    }
    return g
}

func ResultJointEvaluator(m *battle.Manager) (simultaneous.ResultJointEval, error) {
	if !m.CurrentSelfIsHost {
		fmt.Println("ここ")
		return simultaneous.ResultJointEval{}, fmt.Errorf("battle.Manager.CurrentSelfIsHost == false のとき、ResultJointEvaluatorを呼び出してはならない。")
	}

	self := omwslices.Concat(m.CurrentSelfLeadPokemons, m.CurrentSelfBenchPokemons)
	opponent := omwslices.Concat(m.CurrentOpponentLeadPokemons, m.CurrentOpponentBenchPokemons)
	isSelfAllFainted := self.IsAllFainted()
	isOpponentAllFainted := opponent.IsAllFainted()
	if isSelfAllFainted && isOpponentAllFainted {
		fmt.Println("ここ2")
		return simultaneous.ResultJointEval{0.5, 0.5}, nil
	} else if isSelfAllFainted {
		fmt.Println("ここ3")
		return simultaneous.ResultJointEval{0.0, 1.0}, nil
	} else if isOpponentAllFainted {
		fmt.Println("ここ4")
		return simultaneous.ResultJointEval{1.0, 0.0}, nil
	} else {
		fmt.Println("ここ5")
		return simultaneous.ResultJointEval{}, fmt.Errorf("バトルが終了していないとき、ResultJointEvaluatorを呼び出してはならない。")
	}
}