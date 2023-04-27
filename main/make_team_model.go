package main

import (
	"fmt"
	"github.com/sw965/omw"
	bp "github.com/sw965/bippa"
	"github.com/sw965/crow"
	"math"
)

func main() {
	r := omw.NewMt19937()
	model := bp.NewInitPokemonModel("フシギバナ", r)
	fmt.Println(model.ParameterSize())
	//omw.NestMapPrint(model.MoveNameAndMoveName)

	fnCaller := bp.NewTeamBuildPUCTFunCaller(&model, r)
	puct := crow.PUCT[bp.Team, bp.TeamBuildCmd]{FunCaller:fnCaller}

	game := fnCaller.Game.Clone()
	game.SetPUCTPlayer(&puct, 1, math.Sqrt(2), r, 1.0)
	team1 := game.Playout(bp.Team{})
	team2 := game.Playout(bp.Team{})
}