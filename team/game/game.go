package game

import (
	"github.com/sw965/crow/game/sequential"
	"github.com/sw965/bippa/team"
)

func New() sequential.Game[team.Team, team.Actions, team.Action] {
	return sequential.Game[team.Team, team.Actions, team.Action]{
		LegalActions:team.LegalActions,
		Equal:team.Equal,
		Push:team.Push,
		IsEnd:team.IsEnd,
	}
}