package game

import (
	"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/battle/dmgtools"
	"math/rand"
)

func New(randDmgBonuses dmgtools.RandBonuses, r *rand.Rand) simultaneous.Game[single.Battle, single.Actionss, single.Actions, single.Action] {
	ret := simultaneous.Game[single.Battle, single.Actionss, single.Actions, single.Action]{
		Equal:single.Equal,
		IsEnd:single.IsEnd,
		LegalActionss:single.LegalActionss,
		Push:single.NewPushFunc(randDmgBonuses, r),
	}
	ret.SetRandActionPlayer(r)
	return ret
}
