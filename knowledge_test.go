package bippa

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"testing"
	"time"
)

func Accuracy(builder PokemonBuilder, moveNames MoveNames, testNum int, random *rand.Rand) (float64, float64) {
	count1 := 0
	count2 := 0
	for i := 0; i < testNum; i++ {
		getWs := func(pb PokemonBuilder) []float64 {
			return pb.Ws()
		}

		//powerPoint := NewPowerPoint(10, 3)
		moveset, _, err := builder.BuildMoveset(moveNames, Pokemon{}, Team{}, getWs, 4, random)

		if err != nil {
			panic(err)
		}

		_, ok1 := moveset["ギガドレイン"]
		_, ok2 := moveset["ヘドロばくだん"]

		_, ok3 := moveset["まもる"]
		_, ok4 := moveset["やどりぎのタネ"]
		
		_, ok5 := moveset["こうごうせい"]
		_, ok6 := moveset["ねむりごな"]

		if ok1 && ok2 && ok3 && ok4 {
			count1 += 1
		}

		if ok1 && ok2 && ok5 && ok6 {
			count2 += 1
		}
	}
	return float64(count1) / float64(testNum), float64(count2) / float64(testNum)
}

func TestBuildMoveset(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())
	moveNames := MoveNames{"ギガドレイン", "ヘドロばくだん", "こうごうせい", "やどりぎのタネ",
		"だいちのちから", "ソーラービーム", "まもる", "どくどく", "にほんばれ", "ねむりごな"}
	
	builder := NewInitPokemonBuilder(
		Abilities{"しんりょく", "ようりょくそ"},
		Items{"くろいヘドロ", "オボンのみ"},
		moveNames,
		Natures{"さみしがり", "ひかえめ", "いじっぱり"},
		mtRandom,
	).Init(mtRandom)

	//accuracy := Accuracy(moveNames, mtRandom)
	learningRate := 0.0001
	getWs := func(pb PokemonBuilder) []float64 {
		return make([]float64, len(pb))
	}

	actionHistories := make(PokemonBuilderActionHistories, 0, 12800)

	for i := 0; i < 1600000; i++ {
		//powerPoint := NewPowerPoint(10, 3)
		moveset, actionHistory, err := builder.BuildMoveset(moveNames, Pokemon{}, Team{}, getWs, 4, mtRandom)

		if err != nil {
			panic(err)
		}

		ok1 := func() bool {
			_, ok1 := moveset["ギガドレイン"]
			_, ok2 := moveset["ヘドロばくだん"]
			_, ok3 := moveset["やどりぎのタネ"]
			_, ok4 := moveset["まもる"]
			return ok1 && ok2 && ok3 && ok4
		}

		ok2 := func() bool {
			_, ok1 := moveset["ギガドレイン"]
			_, ok2 := moveset["ヘドロばくだん"]
			_, ok3 := moveset["こうごうせい"]
			_, ok4 := moveset["ねむりごな"]
			return ok1 && ok2 && ok3 && ok4
		}

		if ok1() {
			err = builder.Optimizer(&actionHistory, learningRate)
			if err != nil {
				panic(err)
			}
			actionHistories = append(actionHistories, actionHistory)
		}

		if ok2() {
			err = builder.Optimizer(&actionHistory, learningRate)
			if err != nil {
				panic(err)
			}
			actionHistories = append(actionHistories, actionHistory)		
		}

		if i%1600 == 0 {
			for i := 0; i < 64; i++ {
				for _, actionHistory := range actionHistories.RandomChoices(128, mtRandom) {
					err := builder.Optimizer(&actionHistory, learningRate)
					if err != nil {
						panic(err)
					}
				}
			}

			actionHistories = make(PokemonBuilderActionHistories, 0, 12800)

			for _, pbck := range  builder {
				fmt.Println(pbck)
			}

			fmt.Println(Accuracy(builder, moveNames, 1280, mtRandom))
			fmt.Println("")
		}
	}
}
