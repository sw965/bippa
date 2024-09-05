package battle

import (
	"fmt"
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

type SoloActions []SoloAction

func NewLegalSoloActions(m *Manager) SoloActions {
	isSelfAnyFainted := m.CurrentSelfLeadPokemons.IsAnyFainted()
	isOpponentAnyFainted := m.CurrentOpponentLeadPokemons.IsAnyFainted()

	//相手だけ瀕死状態のポケモンがいるならば、自分は行動出来ない。
	if !isSelfAnyFainted && isOpponentAnyFainted {
		return SoloActions{}
	}

	if m.IsSingle() {
		as := make(SoloActions, 0, bp.MAX_MOVESET_LENGTH + len(m.CurrentSelfBenchPokemons))
		p := m.CurrentSelfLeadPokemons[0]
		speed := p.Stat.Speed

		//先頭のポケモンが瀕死状態ならば、技は使えない。
		if !p.IsFainted() {
			for _, moveName := range p.UsableMoveNames() {
				as = append(as, SoloAction{MoveName:moveName, Speed:speed})
			}
		}

		//先頭のポケモンが気絶しているしていないに関係なく、交代は出来る。
		for _, idx := range m.CurrentSelfBenchPokemons.NotFaintedIndices() {
			as = append(as, SoloAction{TargetIndex:idx, Speed:m.CurrentSelfBenchPokemons[idx].Stat.Speed})
		}
		return as
	}

	if isSelfAnyFainted {
		as := make(SoloActions, 0, DOUBLE_BATTLE_NUM * DOUBLE_BATTLE_NUM)
		for _, srcIdx := range m.CurrentSelfLeadPokemons.FaintedIndices() {
			src := m.CurrentSelfLeadPokemons[srcIdx]
			speed := src.Stat.Speed
			for _, targetIdx := range m.CurrentSelfBenchPokemons.NotFaintedIndices() {
				as = append(as, SoloAction{SrcIndex:srcIdx, TargetIndex:targetIdx, Speed:speed})
			}
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

type SoloActionsSlice []SoloActions

type Action [DOUBLE_BATTLE_NUM]SoloAction

func (a *Action) IsEmpty() bool {
	for _, solo := range a {
		if !solo.IsEmpty() {
			return false
		}
	}
	return true
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

	fmt.Println("len(soloActions) = ", len(soloActions))
	soloActionsSlice := omwslices.Combination[SoloActionsSlice, SoloActions](soloActions, DOUBLE_BATTLE_NUM)
	soloActionsSlice = fn.Filter(soloActionsSlice, func(soloActions SoloActions) bool {
		firstSoloAction := soloActions[0]
		secondSoloAction := soloActions[1]

		if secondSoloAction.IsEmpty() {
			return true
		}

		isFirstSwitch := !firstSoloAction.IsMove()
		isSecondSwitch := !secondSoloAction.IsMove()

		if isFirstSwitch && isSecondSwitch {
			/*
				自分の先頭のポケモン (ドクロッグ, エンペルト)
				自分の控えのポケモン (カビゴン, ボーマンダ)
				ドクロッグ → カビゴンと交代
				エンペルト → カビゴンと交代
				という風に、同じポケモンに交代する事は出来ない。
			*/
			if firstSoloAction.TargetIndex == secondSoloAction.TargetIndex {
				return false
			}
		}
		return firstSoloAction.SrcIndex != secondSoloAction.SrcIndex
	})

	fmt.Println("len(soloActionsSlice) = ", len(soloActionsSlice))
	fmt.Println("")
	as := make(Actions, 0, 128)
	for _, soloAs := range soloActionsSlice {
		as = append(as, Action{soloAs[0], soloAs[1]})
	}
	return as
}

func (as Actions) ToSoloActions() SoloActions {
	sas := make(SoloActions, 0, len(as) * DOUBLE_BATTLE_NUM)
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

type ActionsSlice []Actions