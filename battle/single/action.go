package single

import (
	bp "github.com/sw965/bippa"
)

type Action struct {
	CmdMoveName bp.MoveName
	SwitchPokeName bp.PokeName
	IsSelf bool
}

func StringToAction(s string, isPlayer1 bool) Action {
	moveName := bp.STRING_TO_MOVE_NAME[s]
	pokeName := bp.STRING_TO_POKE_NAME[s]
	return Action{CmdMoveName:moveName, SwitchPokeName:pokeName, IsSelf:isPlayer1}
}

func (a *Action) ToString() string {
	p := map[bool]string{true:"player1", false:"player2"}[a.IsSelf]
	return bp.MOVE_NAME_TO_STRING[a.CmdMoveName] + bp.POKE_NAME_TO_STRING[a.SwitchPokeName] + " " + p
}

func (a *Action) IsEmpty() bool {
	return a.CmdMoveName == bp.EMPTY_MOVE_NAME && a.SwitchPokeName == bp.EMPTY_POKE_NAME
}

func (a *Action) IsCommandMove() bool {
	return a.CmdMoveName != bp.EMPTY_MOVE_NAME
}

func (a *Action) IsSwitch() bool {
	return a.SwitchPokeName != bp.EMPTY_POKE_NAME
}

type Actions []Action

func (as Actions) IsAllEmpty() bool {
	for i := range as {
		if !as[i].IsEmpty() {
			return false
		}
	}
	return true
}

func (as Actions) ToStrings() []string {
	ret := make([]string, len(as))
	for i, a := range as {
		ret[i] = a.ToString()
	}
	return ret
}

type ActionSlices []Actions