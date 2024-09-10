package battle

import (
	//"fmt"
	bp "github.com/sw965/bippa"
	"golang.org/x/exp/slices"
	omwrand "github.com/sw965/omw/math/rand"
	omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/omw/fn"
)

type SoloAction struct {
	MoveName bp.MoveName
	SrcIndex int
	TargetIndex int
	IsSelfLeadTarget bool
	Speed int
	IsCurrentSelf bool
}

func (a *SoloAction) IsEmpty() bool {
	return a.Speed == 0
}

func (a *SoloAction) IsMove() bool {
	return a.MoveName != bp.EMPTY_MOVE_NAME
}

func (a *SoloAction) Priority() int {
	if a.IsMove() {
		moveData := bp.MOVEDEX[a.MoveName]
		return moveData.PriorityRank
	} else {
		return 999
	}
}

func (a *SoloAction) ToEasyRead() EasyReadSoloAction {
	return EasyReadSoloAction{
		MoveName:a.MoveName.ToString(),
		SrcIndex:a.SrcIndex,
		TargetIndex:a.TargetIndex,
		IsSelfLeadTarget:a.IsSelfLeadTarget,
		Speed:a.Speed,
		IsCurrentSelf:a.IsCurrentSelf,
	}
}

type EasyReadSoloAction struct {
	MoveName string
	SrcIndex int
	TargetIndex int
	IsSelfLeadTarget bool
	Speed int
	IsCurrentSelf bool
}

func (e *EasyReadSoloAction) From() (SoloAction, error) {
	moveName, err := bp.StringToMoveName(e.MoveName)
	return SoloAction{
		MoveName:moveName,
		SrcIndex:e.SrcIndex,
		TargetIndex:e.TargetIndex,
		IsSelfLeadTarget:e.IsSelfLeadTarget,
		Speed:e.Speed,
		IsCurrentSelf:e.IsCurrentSelf,
	}, err
} 

type SoloActions []SoloAction

func NewLegalSoloActions(m *Manager) SoloActions {
	selfMustSwitch, opponentMustSwitch := m.MustSwitch()

	//相手だけ交代しなければならない状態であれば、自分は行動出来ない。
	if !selfMustSwitch && opponentMustSwitch {
		return SoloActions{}
	}

	if m.IsSingle() {
		as := make(SoloActions, 0, bp.MAX_MOVESET_LENGTH + len(m.CurrentSelfBenchPokemons))
		src := m.CurrentSelfLeadPokemons[0]
		speed := src.Stat.Speed

		//先頭のポケモンが瀕死状態ではないならば、技を使える。
		if !src.IsFainted() {
			for _, moveName := range src.UsableMoveNames() {
				as = append(as, SoloAction{MoveName:moveName, Speed:speed})
			}
		}

		//先頭のポケモンの瀕死状態の有無に関わらず、交代は出来る。
		for _, idx := range m.CurrentSelfBenchPokemons.NotFaintedIndices() {
			as = append(as, SoloAction{TargetIndex:idx, Speed:speed})
		}
		return as
	}

	if selfMustSwitch {
		as := make(SoloActions, 0, DOUBLE * DOUBLE)
		//複数のポケモンが瀕死状態であっても、1匹ずつ交代する為、0番目のインデックスにアクセスする。
		srcIdx := m.CurrentSelfLeadPokemons.FaintedIndices()[0]
		src := m.CurrentSelfLeadPokemons[srcIdx]
		speed := src.Stat.Speed
		for _, targetIdx := range m.CurrentSelfBenchPokemons.NotFaintedIndices() {
			as = append(as, SoloAction{SrcIndex:srcIdx, TargetIndex:targetIdx, Speed:speed})
		}
		return as
	}

	as := make(SoloActions, 0, 128)
	for _, srcIdx := range m.CurrentSelfLeadPokemons.NotFaintedIndices() {
		src := m.CurrentSelfLeadPokemons[srcIdx]
		speed := src.Stat.Speed
		for _, moveName := range src.UsableMoveNames() {
			//相手を対象指定して、技を繰り出す。
			for _, targetIdx := range m.CurrentOpponentLeadPokemons.NotFaintedIndices() {
				as = append(as, SoloAction{MoveName:moveName, SrcIndex:srcIdx, TargetIndex:targetIdx, Speed:speed})
			}

			//味方を対象指定して、技を繰り出す。
			for _, targetIdx := range m.CurrentSelfLeadPokemons.NotFaintedIndices() {
				as = append(as, SoloAction{MoveName:moveName, SrcIndex:srcIdx, TargetIndex:targetIdx, IsSelfLeadTarget:true, Speed:speed})
			}

			//対象指定せずに、技を繰り出す。
			as = append(as, SoloAction{MoveName:moveName, SrcIndex:srcIdx, TargetIndex:-1, Speed:speed})
		}

		//交代
		for _, targetIdx := range m.CurrentSelfBenchPokemons.NotFaintedIndices() {
			as = append(as, SoloAction{SrcIndex:srcIdx, TargetIndex:targetIdx, Speed:speed})
		}
	}

	return fn.Filter(as, func(a SoloAction) bool {
		if a.IsMove() {
			moveData := bp.MOVEDEX[a.MoveName]
			switch moveData.Target {
				case bp.NORMAL_TARGET:
					//自分自身への攻撃は出来ない
					if a.IsSelfLeadTarget {
						return a.SrcIndex != a.TargetIndex
					}
					return a.TargetIndex != - 1
				case bp.SELF_TARGET:
					return a.SrcIndex == a.TargetIndex && a.IsSelfLeadTarget
				/*
					defaultは 下記のTargetを想定している。
					OPPONENT_TWO_TARGET (いわなだれ 等)
					OTHERS_TARGET (じばく/だいばくはつ 等)
					ALL_TARGET (あまごい 等)
					OPPONENT_RANDOM_ONE_TARGET (わるあがぎ 等)
				*/
				default:
					return a.TargetIndex == -1
			}
		}
		return true
	})
}

func (as SoloActions) SortByOrder(m *Manager) {
	slices.SortFunc(as, func(a1, a2 SoloAction) bool {
		a1Priority := a1.Priority()
		a2Priority := a2.Priority()
		if a1Priority > a2Priority {
			return true
		} else if a1Priority < a2Priority {
			return false
		} else {
			a1Speed := a1.Speed
			a2Speed := a2.Speed
			if a1Speed > a2Speed {
				return !m.IsTrickRoomState()
			} else if a1Speed < a2Speed {
				return m.IsTrickRoomState()
			} else {
				return omwrand.Bool(GlobalContext.Rand)
			}
		}
	})
}

func (as SoloActions) FilterByNotEmpty() SoloActions {
	s := make(SoloActions, 0, len(as))
	for _, a := range as {
		if !a.IsEmpty() {
			s = append(s, a)
		}
	}
	return s
}

func (as SoloActions) FilterByEqualSrcIndex(idx int) SoloActions {
	s := make(SoloActions, 0, len(as))
	for _, a := range as {
		if a.SrcIndex == idx {
			s = append(s, a)
		}
	}
	return s
}

type SoloActionsSlice []SoloActions

type Action [DOUBLE]SoloAction

func (a *Action) IsEmpty() bool {
	for _, solo := range a {
		if !solo.IsEmpty() {
			return false
		}
	}
	return true
}

func (a *Action) ToEasyRead() EasyReadAction {
	e := EasyReadAction{}
	for i, soloAction := range a {
		e[i] = soloAction.ToEasyRead()
	}
	return e
}

type EasyReadAction [DOUBLE]EasyReadSoloAction

func (e *EasyReadAction) From() (Action, error) {
	var err error
	a := Action{}
	for i, solo := range e {
		a[i], err = solo.From()
		if err != nil {
			return Action{}, err
		}
	}
	return a, nil
}

type Actions []Action

func NewLegalActions(m *Manager) Actions {
	soloActions := NewLegalSoloActions(m)
	if m.IsSingle() {
		as := make(Actions, len(soloActions))
		for i, soloAction := range soloActions {
			as[i] = Action{soloAction}
		}
		return as
	}

	groupedSoloActionsSliceBySrcIdx := make(SoloActionsSlice, DOUBLE)
	for i := 0; i < DOUBLE; i++ {
		groupedSoloActionsSliceBySrcIdx[i] = soloActions.FilterByEqualSrcIndex(i)
	}

	for srcIdx, soloActions := range groupedSoloActionsSliceBySrcIdx {
		if len(soloActions) == 0 {
			groupedSoloActionsSliceBySrcIdx[srcIdx] = SoloActions{SoloAction{SrcIndex:srcIdx}}
		}
	}

	soloActionsSlice := omwslices.CartesianProduct[SoloActionsSlice](groupedSoloActionsSliceBySrcIdx...)
	soloActionsSlice = fn.Filter(soloActionsSlice, func(soloActions SoloActions) bool {
		first := soloActions[0]
		second := soloActions[1]

		isFirstSwitch := !first.IsMove()
		isSecondSwitch := !second.IsMove()

		if isFirstSwitch && !first.IsEmpty() && isSecondSwitch && !second.IsEmpty() {
			/*
				自分の先頭のポケモン (ドクロッグ, エンペルト)
				自分の控えのポケモン (カビゴン, ボーマンダ)
				ドクロッグ → カビゴンと交代
				エンペルト → カビゴンと交代
				みたいに、同じポケモンに交代する事は出来ない。
			*/
			if first.TargetIndex == second.TargetIndex {
				return false
			}
		}
		return true
	})

	as := make(Actions, len(soloActionsSlice))
	for i, soloActions := range soloActionsSlice {
		as[i] = Action{soloActions[0], soloActions[1]}
	}
	return as
}

func (as Actions) ToSoloActions() SoloActions {
	sas := make(SoloActions, 0, len(as) * DOUBLE)
	for _, a := range as {
		for _, soloAction := range a {
			sas = append(sas, soloAction)
		}
	}
	return sas
}

func (as Actions) FilterByNotEmpty() Actions {
	s := make(Actions, 0, len(as))
	for _, a := range as {
		if !a.IsEmpty() {
			s = append(s, a)
		}
	}
	return s
}

func (as Actions) ToEasyRead() EasyReadActions {
	es := make(EasyReadActions, len(as))
	for i, a := range as {
		es[i] = a.ToEasyRead()
	}
	return es
}

type EasyReadActions []EasyReadAction

type ActionsSlice []Actions

func (ass ActionsSlice) ToEasyRead() EasyReadActionsSlice {
	ess := make(EasyReadActionsSlice, len(ass))
	for i, as := range ass {
		ess[i] = as.ToEasyRead()
	}
	return ess
}

type EasyReadActionsSlice []EasyReadActions