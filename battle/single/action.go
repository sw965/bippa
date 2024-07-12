package single

import (
	bp "github.com/sw965/bippa"
)

type Action struct {
	MoveNames [2]bp.MoveName
	OpponentLeadIndices [2]int
	SelfBenchIndices [2]int
	Speeds [2]int
	IsSelf bool
}

// func (a *Action) IsMove() bool {
// 	return a.MoveName != bp.EMPTY_MOVE_NAME
// }

// func (a *Action) IsSwitch() bool {
// 	return a.PokeName != bp.EMPTY_POKE_NAME
// }

type Actions []Action

// func (as Actions) IsAllEmpty() bool {
// 	for i := range as {
// 		if !as[i].IsEmpty() {
// 			return false
// 		}
// 	}
// 	return true
// }

type ActionsSlice []Actions