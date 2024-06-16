package game

import (
	"github.com/sw965/crow/game/simultaneous"
	"github.com/sw965/bippa/battle/single"
	"math/rand"
)

func New(r *rand.Rand) simultaneous.Game[single.Battle, single.ActionSlices, single.Actions, single.Action] {
    gm := simultaneous.Game[single.Battle, single.ActionSlices, single.Actions, single.Action]{
        Equal:                single.Equal,
        IsEnd:                single.IsEnd,
        LegalSeparateActions: single.LegalSeparateActions,
        Push:                 single.NewPushFunc(r),
    }
    gm.SetRandActionPlayer(r)
    return gm
}