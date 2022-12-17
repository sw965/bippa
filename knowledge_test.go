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
		Items:     Items{"くろいヘドロ", "オボンのみ", "いのちのたま", "しろいハーブ"},
		MoveNames: MoveNames{"ギガドレイン", "ヘドロばくだん", "やどりぎのタネ", "まもる", "どくどく", "こうごうせい", "ねむりごな", "だいちのちから", "リーフストーム"},
		Natures: Natures{"ずぶとい", "しんちょう", "ひかえめ"},
	}

	builder := NewInitPokemonBuildCombinationKnowledgeList(
		"フシギバナ", &pbCommonK,
	)

	builder.Init()

	learningRate := 0.001
	selectKnowledgeListsU := make(PokemonBuildCombinationKnowledgeLists, 0, 10000)
	usableKnowledgeListsU := make(PokemonBuildCombinationKnowledgeLists, 0, 10000)

	for i := 0; i < 1600000; i++ {
		pokemon, selectKnowledgeLists, usableKnowledgeLists, err := builder.Run(
			Pokemon{Name:"フシギバナ"}, Team{}, &pbCommonK,
			func(meanPolicies []float64) []float64 { return meanPolicies}, mtRandom,
		)

		if err != nil {
			panic(err)
		}

		ok := func() bool {
			moveset := pokemon.Moveset
			_, ok1 := moveset["ギガドレイン"]
			_, ok2 := moveset["ヘドロばくだん"]
			_, ok3 := moveset["やどりぎのタネ"]
			_, ok4 := moveset["まもる"]
			return ok1 && ok2 && ok3 && ok4
		}

		if ok() {
			selectKnowledgeListsU = append(selectKnowledgeListsU, selectKnowledgeLists...)
			usableKnowledgeListsU = append(usableKnowledgeListsU, usableKnowledgeLists...)
		}

		if i%12800== 0 {
			selectKnowledgeListsU.PolicyOptimizer(usableKnowledgeListsU, 1.0, learningRate, 0.001)
			selectKnowledgeListsU = make(PokemonBuildCombinationKnowledgeLists, 0, 10000)
			usableKnowledgeListsU = make(PokemonBuildCombinationKnowledgeLists, 0, 10000)

			testPoke, _, _, err := builder.Run(
				Pokemon{Name:"フシギバナ"}, Team{}, &pbCommonK,
				func(meanPolicies []float64) []float64 {return meanPolicies }, mtRandom,
			)
			if err != nil {
				panic(err)
			}

			for _, pbCombK := range builder {
				fmt.Println(pbCombK)
			}

			fmt.Println(testPoke.Moveset)
		}
	}
}
