package bippa

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"testing"
	"time"
)

// func Accuracy(moveNames MoveNames, random *rand.Rand) func(PokemonBuilder) float64 {
// 	result := func(pb PokemonBuilder) float64 {
// 		count := 0.0
// 		trialNum := 196

// 		for i := 0; i < trialNum; i++ {
// 			moveset, selectPokemonBuilder, err := pb.BuildMoveset(moveNames, Pokemon{}, Team{}, random)

// 			if err != nil {
// 				panic(err)
// 			}

// 			_, ok1 := moveset["ギガドレイン"]
// 			_, ok2 := moveset["ヘドロばくだん"]
// 			_, ok3 := moveset["やどりぎのタネ"]
// 			_, ok4 := moveset["まもる"]

// 			if ok1 {
// 				count += 1
// 			}

// 			if ok2 {
// 				count += 1
// 			}

// 			if ok3 {
// 				count += 1
// 			}

// 			if ok4 {
// 				count += 1
// 			}
// 		}
// 		result := float64(count) / float64(trialNum * 4)
// 		t := 1.0

// 		learningRate := (float64(result) - t) * (float64(result) - t)
// 	}
// 	return result
// }

func Accuracy(builder PokemonBuilder, moveNames MoveNames, testNum int, random *rand.Rand) float64 {
	count := 0
	for i := 0; i < testNum; i++ {
		getSelectPercents := func(pb PokemonBuilder) []float64 {
			return pb.SelectPercents()
		}

		moveset, _, err := builder.BuildMoveset(moveNames, Pokemon{}, Team{}, getSelectPercents, random)

		if err != nil {
			panic(err)
		}

		_, ok1 := moveset["ギガドレイン"]
		_, ok2 := moveset["ヘドロばくだん"]
		_, ok3 := moveset["やどりぎのタネ"]
		_, ok4 := moveset["まもる"]
		if ok1 && ok2 && ok3 && ok4 {
			count += 1
		}
	}
	return float64(count) / float64(testNum)
}

func TestBuildMoveset(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())
	moveNames := MoveNames{"ギガドレイン", "ヘドロばくだん", "こうごうせい", "やどりぎのタネ",
		"だいちのちから", "ソーラービーム", "まもる", "どくどく", "にほんばれ"}
	
	builder := NewInitPokemonBuilder(
		Abilities{"しんりょく", "ようりょくそ"},
		Items{"くろいヘドロ", "オボンのみ"},
		moveNames,
		Natures{"さみしがり", "ひかえめ", "いじっぱり"},
		mtRandom,
	).Init(mtRandom)

	//accuracy := Accuracy(moveNames, mtRandom)
	learningRate := 0.001
	getSelectPercents := func(pb PokemonBuilder) []float64 {
		return make([]float64, len(pb))
	}

	for i := 0; i < 1600000; i++ {
		//powerPoint := NewPowerPoint(10, 3)
		ms := Moveset{}
		moveset, actionHistory, err := builder.BuildMoveset(moveNames, Pokemon{Moveset:ms}, Team{}, getSelectPercents, mtRandom)

		if err != nil {
			panic(err)
		}

		ok := func(pokemon Pokemon, team Team) bool {
			_, ok1 := moveset["やどりぎのタネ"]
			_, ok2 := moveset["まもる"]
			_, ok3 := moveset["ギガドレイン"]
			_, ok4 := moveset["ヘドロばくだん"]
			return ok1 && ok2 && ok3 && ok4
		}

		teacher := PokemonBuilderTeacher{ActionHistory:actionHistory, Pokemon:Pokemon{}, Team:Team{}, OK:ok}
		err = builder.Optimizer(&teacher, 0.9, learningRate)
		if err != nil {
			panic(err)
		}

		if i%1600 == 0 {
			for _, pbck := range  builder {
				fmt.Println(pbck)
			}
			fmt.Println(Accuracy(builder, moveNames, 1280, mtRandom))
			fmt.Println("")
		}
	}
}
