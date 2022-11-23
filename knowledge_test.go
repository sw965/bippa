package bippa

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {

}

func TestBuildMoveset(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	pokemonBuildKnowLedge := NewVenusaurBuildKnowledge()
	for i := 0; i < 100; i++ {
		pokemon, err := pokemonBuildKnowLedge.BuildMoveset(Pokemon{Moveset: Moveset{}}, Team{}, mtRandom)
		if err != nil {
			panic(err)
		}

		fmt.Println(pokemon.Moveset)
		if _, ok := pokemon.Moveset["ギガドレイン"]; !ok {
			break
		}
	}
}
