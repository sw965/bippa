package bippa

import (
  "fmt"
)

type FinalDefense int

func (spovb *SelfPointOfViewBattle) NewFinalDefenseCalc(moveName MoveName, isCritical bool) (FinalDefense, error) {
	moveData := MOVEDEX[moveName]

	var defValue int
	var rank int

	switch moveData.Category {
	case PHYSICS:
		defValue = spovb.SelfFighters[0].State.Def
		rank = spovb.SelfFighters[0].RankState.Def
	case SPECIAL:
		defValue = spovb.SelfFighters[0].State.SpDef
		rank = spovb.SelfFighters[0].RankState.SpDef
	default:
		return 0, fmt.Errorf("変化技以外でなければならない")
	}

	if rank > 0 && isCritical {
		rank = 0
	}

  rankBonus := RANK_TO_RANK_BONUS[rank]
	finalDefense := int(float64(defValue) * rankBonus)

	if finalDefense < 1 {
		return 1, nil
	}
	return finalDefense, nil
}
