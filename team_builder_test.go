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

	fmt.Println(EFFECTIVE_EFFORTS)

	pbk1, err := LoadJsonPokemonBuildCommonKnowledge("フシギバナ")
	if err != nil {
		panic(err)
	}
	fmt.Println(pbk1)

	pbk2, err := LoadJsonPokemonBuildCommonKnowledge("リザードン")
	if err != nil {
		panic(err)
	}
	fmt.Println(pbk2)

	mpscs := NewPokemon1MoveNameAndNatureAndPokemon2EffortCombinations(&pbk1, &pbk2, "HP")

	mpscms := NewMultiplePokemonStateCombinationModels(mpscs, mtRandom)
	err = mpscms.WriteJson("フシギバナ", "リザードン")
	if err != nil {
		panic(err)
	}
}