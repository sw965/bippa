package bippa

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"github.com/sw965/omw"
	"math/rand"
	"testing"
	"time"
)

func Accuracy(moveNames MoveNames, random *rand.Rand) func(PokemonBuilder) float64 {
	result := func(pb PokemonBuilder) float64 {
		// result := 0.0
		// for _, pbk := range pb {
		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"ギガドレイン"}}) {
		// 		if pbk.Tier == TIER1 {
		// 			result += 1.0
		// 		}
		// 	}

		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"ヘドロばくだん"}}) {
		// 		if pbk.Tier == TIER1 {
		// 			result += 1.0
		// 		}
		// 	}

		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"まもる"}}) {
		// 		if pbk.Tier == TIER1 {
		// 			result += 1.0
		// 		}
		// 	}

		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"やどりぎのタネ"}}) {
		// 		if pbk.Tier == TIER1 {
		// 			result += 1.0
		// 		}
		// 	}

		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"だいちのちから"}}) {
		// 		if pbk.Tier != TIER1 {
		// 			result += 1.0
		// 		}
		// 	}

		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"どくどく"}}) {
		// 		if pbk.Tier != TIER1 {
		// 			result += 1.0
		// 		}
		// 	}

		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"こうごうせい"}}) {
		// 		if pbk.Tier != TIER1 {
		// 			result += 1.0
		// 		}
		// 	}

		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"にほんばれ"}}) {
		// 		if pbk.Tier != TIER1 {
		// 			result += 1.0
		// 		}
		// 	}

		// 	if pbk.NearlyEqual(&PokemonBuildKnowledge{SelfMoveNames:MoveNames{"ソーラービーム"}}) {
		// 		if pbk.Tier != TIER1 {
		// 			result += 1.0
		// 		}
		// 	}
		// }

		trialNum := 196
		count := 0.0

		for i := 0; i < trialNum; i++ {
			moveset, err := pb.BuildMoveset(moveNames, Pokemon{}, Team{}, random)
			if err != nil {
				panic(err)
			}
			_, ok1 := moveset["ギガドレイン"]
			_, ok2 := moveset["ヘドロばくだん"]
			_, ok3 := moveset["やどりぎのタネ"]
			_, ok4 := moveset["まもる"]

			if ok1 {
				count += 1.0
			}

			if ok2 {
				count += 1.0
			}

			if ok3 {
				count += 1.0
			}

			if ok4 {
				count += 1.0
			}
		}
		return count
		//return result
	}
	return result
}

func TestBuildMoveset(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())
	moveNames := MoveNames{"ギガドレイン", "ヘドロばくだん", "こうごうせい", "やどりぎのタネ",
		"だいちのちから", "ソーラービーム", "まもる", "どくどく", "にほんばれ"}

	size := 512
	builders := make(PokemonBuilders, size)
	for i := 0; i < size; i++ {
		builders[i] = NewInitPokemonBuilder(
			Abilities{"しんりょく", "ようりょくそ"},
			Items{"くろいヘドロ", "オボンのみ"},
			moveNames,
			Natures{"さみしがり", "ひかえめ", "いじっぱり"},
			mtRandom,
		)
	}
	accuracy := Accuracy(moveNames, mtRandom)

	for i := 0; i < 12800; i++ {
		builders, err := builders.NextGeneration(accuracy, 96, 0.5, 0.01, 1, mtRandom)
		if err != nil {
			panic(err)
		}

		if i%1 == 0 {
			accuracyYs := builders.AccuracyYs(accuracy)
			elite := builders.Elite(accuracyYs).RandomChoice(mtRandom)
			fmt.Println("len = ", len(builders), "accyracy = ", omw.MaxFloat64(accuracyYs...), "i = ", i)

			fmt.Println("tier1", elite.Filter(func(pbk *PokemonBuildKnowledge) bool { return pbk.Tier == TIER1 }))
			fmt.Println("tier2", elite.Filter(func(pbk *PokemonBuildKnowledge) bool { return pbk.Tier == TIER2 }))
			fmt.Println("tier3", elite.Filter(func(pbk *PokemonBuildKnowledge) bool { return pbk.Tier == TIER3 }))
			fmt.Println("tier4", elite.Filter(func(pbk *PokemonBuildKnowledge) bool { return pbk.Tier == TIER4 }))
			fmt.Println("tier5", elite.Filter(func(pbk *PokemonBuildKnowledge) bool { return pbk.Tier == TIER5 }))
		}
	}
}
