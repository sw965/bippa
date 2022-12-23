package bippa

import (
	"testing"
	"fmt"
	"os"
)

func Test(t *testing.T) {
	path := os.Getenv("GOPATH") + "ratta/pokemon_build_common_knowledge/" 
	pbCommonK1 := ReadJsonPokemonBuildCommonKnowledge(path + "フシギバナ.json")
	pbCommonK2 := ReadJsonPokemonBuildCommonKnowledge(path + "リザードン.json")

	result := NewTeamCombinationFeatures("フシギバナ", "リザードン",
	map[PokeName]*PokemonBuildCommonKnowledge{"フシギバナ":&pbCommonK1, "リザードン":&pbCommonK2})

	for _, tcf := range result {
		fmt.Println("---")
		for k, v := range tcf {
			fmt.Println(k, *v)
		}
		fmt.Println("---")
	}

	fmt.Println(len(result))
}