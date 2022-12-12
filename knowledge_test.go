package bippa

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"testing"
	"time"
)

func TestBuildMoveset(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	fmt.Println(ALL_NATURES, len(POKEDEX["フシギバナ"].Learnset), POKEDEX["フシギバナ"].Learnset)

	pbCommonK := PokemonBuildCommonKnowledge{
		Items:     Items{"くろいヘドロ", "オボンのみ", "いのちのたま"},
		MoveNames: MoveNames{"ギガドレイン", "ヘドロばくだん", "やどりぎのタネ", "まもる", "どくどく", "こうごうせい", "ねむりごな", "だいちのちから"},
		//MoveNames:POKEDEX["フシギバナ"].Learnset,
		Natures: ALL_NATURES,
	}

	builder := NewInitPokemonBuilder(
		"フシギバナ", &pbCommonK,
	)

	builder.Init()

	learningRate := 0.0001
	finalWsFunc := func(averageWs []float64) []float64 {
		return make([]float64, len(averageWs))
	}

	actionHistories := make(PokemonBuilderActionHistories, 0, 10000)

	for i := 0; i < 1600000; i++ {
		powerPoint := NewPowerPoint(10, 3)
		moveset, actionHistory, err := builder.BuildMoveset(pbCommonK.MoveNames, Pokemon{}, Team{}, 1, finalWsFunc, mtRandom)

		if err != nil {
			panic(err)
		}

		gigaok := func() bool {
			_, ok := moveset["ギガドレイン"]
			return ok
		}

		if gigaok() {
			actionHistories = append(actionHistories, actionHistory)
		}

		moveset, actionHistory, err = builder.BuildMoveset(pbCommonK.MoveNames, Pokemon{Moveset: Moveset{"ギガドレイン": &powerPoint}}, Team{}, 1, finalWsFunc, mtRandom)

		hedook := func() bool {
			_, ok := moveset["ヘドロばくだん"]
			return ok
		}

		if hedook() {
			actionHistories = append(actionHistories, actionHistory)
		}

		moveset, actionHistory, err = builder.BuildMoveset(pbCommonK.MoveNames, Pokemon{Moveset: Moveset{"ギガドレイン": &powerPoint, "ヘドロばくだん": &powerPoint}}, Team{}, 1, finalWsFunc, mtRandom)

		mamook := func() bool {
			_, ok := moveset["まもる"]
			return ok
		}

		if mamook() {
			actionHistories = append(actionHistories, actionHistory)
		}

		moveset, actionHistory, err = builder.BuildMoveset(pbCommonK.MoveNames, Pokemon{Moveset: Moveset{"ギガドレイン": &powerPoint, "ヘドロばくだん": &powerPoint, "まもる": &powerPoint}}, Team{}, 1, finalWsFunc, mtRandom)

		dokuok := func() bool {
			_, ok := moveset["どくどく"]
			return ok
		}

		if dokuok() {
			actionHistories = append(actionHistories, actionHistory)
		}

		yadook := func() bool {
			_, ok := moveset["やどりぎのタネ"]
			return ok
		}

		if yadook() {
			actionHistories = append(actionHistories, actionHistory)
		}

		if i%1600 == 0 {
			err := actionHistories.Optimizer(128, 64, learningRate, mtRandom)
			if err != nil {
				panic(err)
			}
			actionHistories = make(PokemonBuilderActionHistories, 0, 1600)

			for _, pbck := range builder {
				fmt.Println(pbck)
			}
			fmt.Println(len(builder))
			fmt.Println(Accuracy(builder, pbCommonK.MoveNames, 10000, mtRandom))
			fmt.Println("")
		}
	}
}
