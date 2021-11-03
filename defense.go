package bippa

import (
  "fmt"
)

type FinalDefense int

func NewFinalDefense(spovb *SelfPointOfViewBattle, moveName MoveName, isCritical bool) (FinalDefense, error) {
	moveData := MOVEDEX[moveName]

	var defenseState_ State_
	var rank_ Rank_

	switch moveData.Category {
	case PHYSICS:
		defenseState_ = spovb.SelfFighters[0].State.Def
		rank_ = spovb.SelfFighters[0].Rank.Def
	case SPECIAL:
		defenseState_ = spovb.SelfFighters[0].State.SpDef
		rank_ = spovb.SelfFighters[0].Rank.SpDef
	default:
		return 0, fmt.Errorf("変化技以外でなければならない")
	}

	if rank_ > 0 && isCritical {
		rank_ = 0
	}

  rankBonus := RANK__TO_RANK_BONUS[rank_]
	result := int(float64(defenseState_) * float64(rankBonus))

	if result < 1 {
		return 1, nil
	}
	return FinalDefense(result), nil
}
