package bippa

// import (
// 	"math/rand"
// 	"github.com/sw965/crow/mcts/puct"
// 	"github.com/sw965/crow/game/sequential"
// )

// func NewTeamBuildMCTS(model *PokemonModel, r *rand.Rand) puct.MCTS[Team, TeamBuildActions, TeamBuildAction] {
// 	legalActions := func(team *Team) TeamBuildActions {
// 		return team.LegalBuildActions()
// 	}

// 	push := func(team Team, action TeamBuildAction) Team {
// 		return team.Push(&action)
// 	}

// 	isEnd := func(team *Team) bool {
// 		return len(team.LegalBuildActions()) == 0
// 	}

// 	equal := func(team1, team2 *Team) bool {
// 		return team1.Equal(*team2)
// 	}

// 	game := sequential.Game[Team, TeamBuildActions, TeamBuildAction]{
// 		LegalActions: legalActions,
// 		Push:         push,
// 		Equal:        equal,
// 		IsEnd:        isEnd,
// 	}

// 	game.SetRandomActionPlayer(r)

// 	leaf := func(team *Team) puct.LeafEvalY {
// 		y := 0.0
// 		for _, pokemon := range *team {
// 			y += model.Output(&pokemon)
// 		}
// 		return puct.LeafEvalY(y)
// 	}

// 	backward := func(y puct.LeafEvalY, team *Team) puct.BackwardEvalY {
// 		return puct.BackwardEvalY(y)
// 	}

// 	eval := puct.Eval[Team]{
// 		Leaf:     leaf,
// 		Backward: backward,
// 	}

// 	mcts := puct.MCTS[Team, TeamBuildActions, TeamBuildAction]{
// 		Game: game,
// 		Eval: eval,
// 	}
// 	mcts.SetNoPolicy()
// 	return mcts
// }