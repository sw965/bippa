package main

import (
	"fmt"
	bp "github.com/sw965/bippa"
	orand "github.com/sw965/omw/rand"
	"github.com/sw965/bippa/dmgtools"
	oslices "github.com/sw965/omw/slices"
	"github.com/sw965/crow/model"
	sb "github.com/sw965/bippa/battle/single"
	"golang.org/x/exp/slices"
	"github.com/sw965/crow/tensor"
)

func main() {
	r := orand.NewMt19937()
	testNum := 12800
	pokemons := bp.Team{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
		bp.NewTemplateSuicune(),
		bp.NewTemplateGarchomp(),
	}

	nn, _ := model.NewThreeLayerAffineParamReLUInput1DOutputSigmoid1D(
		(len(bp.ALL_MOVE_NAMES) + (2* (len(bp.ALL_TYPES) + len(oslices.Combination[[]bp.Types, bp.Types](bp.ALL_TYPES, 2))))) * 6,
		64, 16, 1, 0.001, 64.0, r,
	)

	game := sb.NewGame(dmgtools.RAND_BONUSES, r)
	game.SetRandActionPlayer(r)

	for i := 0; i < testNum; i++ {
		team1 := slices.Clone(pokemons)
		orand.Shuffle(team1, r)
		team1 = team1[:3]
		f1 := sb.Fighters{}
		for i, pokemon := range team1 {
			f1[i] = pokemon
		}

		team2 := slices.Clone(pokemons)
		orand.Shuffle(team2, r)
		team2 = team2[:3]
		f2 := sb.Fighters{}
		for i, pokemon := range team2 {
			f2[i] = pokemon
		}

		fmt.Println("i=", oslices.Concat(team1, team2).PokeNames().ToStrings())

		x := bp.TeamEvalFeature(oslices.Concat(team1, team2))
		y, err := nn.Predict(x)
		if err != nil {
			panic(err)
		}
		fmt.Println("y=", y)

		initBattle := sb.Battle{P1Fighters:f1, P2Fighters:f2}
		endBattle, err := game.Playout(initBattle)
		if err != nil {
			panic(err)
		}
		ys, err := endBattle.EndLeafNodeEvalYs()
		if err != nil {
			panic(err)
		}
		nn.SGD(x, tensor.D1{float64(ys[0])}, 0.01)
	}
}