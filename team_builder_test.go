package bippa

import (
	"testing"
	"fmt"
	"math/rand"
	"time"
	"github.com/seehuhn/mt19937"
)

func Test(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	fmt.Println(ALL_VALID_EFFORTS)

	pbk, err := LoadJsonPokemonBuildCommonKnowledge("フシギバナ")
	if err != nil {
		panic(err)
	}
	fmt.Println(pbk)

	pscs := NewNatureAndIndividualAndEffortCombinations(&pbk, "HP", "Atk")
	if err != nil {
		panic(err)
	}

	pscms := NewPokemonStateCombinationModels(pscs, mtRandom)
	err = pscms.WriteJson("フシギバナ", "test.json")
	if err != nil {
		panic(err)
	}
}