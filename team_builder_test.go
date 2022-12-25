package bippa

import (
	"testing"
	//"fmt"
	"math/rand"
	"time"
	"github.com/seehuhn/mt19937"
)

func Test(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	pbCommonK1, err := LoadJsonPokemonBuildCommonKnowledge("フシギバナ")
	if err != nil {
		panic(err)
	}

	psce := NewPokemonStateCombinationEvaluator("フシギバナ", &pbCommonK1, mtRandom)
	psce.WriteJson("フシギバナ")

	tce, err := NewTeamCombinationEvaluator("フシギバナ", "リザードン", mtRandom)
	if err != nil {
		panic(err)
	}
	err = tce.WriteJson()
	if err != nil {
		panic(err)
	}

	//pbCommonK2 := ReadJsonPokemonBuildCommonKnowledge(path + "リザードン.json")

	// result := NewTeamCombinationFeatures("フシギバナ", "リザードン",
	// map[PokeName]*PokemonBuildCommonKnowledge{"フシギバナ":&pbCommonK1, "リザードン":&pbCommonK2})

	// for _, tcf := range result {
	// 	fmt.Println("---")
	// 	for k, v := range tcf {
	// 		fmt.Println(k, *v)
	// 	}
	// 	fmt.Println("---")
	// }
	// fmt.Println(result[0], result[1])
	// fmt.Println(len(result))
}