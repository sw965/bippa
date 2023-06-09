package bippa_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
	"golang.org/x/exp/slices"
)

func TestGetState(t *testing.T) {
	result := bp.GetState(102, 0, 0, bp.DOWN_NATURE_BONUS)
	expected := bp.State(96)
	if result != expected {
		t.Errorf("テスト失敗")
	}
}

func TestGetAllHPs(t *testing.T) {
	pokeData := bp.POKEDEX["フシギバナ"]
	result := bp.GetAllHPs(pokeData.BaseHP)
	fmt.Println("フシギバナ")
	fmt.Println(result)

	pokeData = bp.POKEDEX["ハピナス"]
	result = bp.GetAllHPs(pokeData.BaseHP)
	fmt.Println("ハピナス")
	fmt.Println(result)

	pokeData = bp.POKEDEX["ガブリアス"]
	result = bp.GetAllStates(pokeData.BaseSpeed)
	slices.Sort(result)
	fmt.Println("ガブリアス")
	fmt.Println(result)

	pokeData = bp.POKEDEX["サンダース"]
	result = bp.GetAllStates(pokeData.BaseSpeed)
	slices.Sort(result)
	fmt.Println("さんだす")
	fmt.Println(result)

	pokeData = bp.POKEDEX["リザードン"]
	result = bp.GetAllStates(pokeData.BaseSpeed)
	slices.Sort(result)
	fmt.Println("りざど")
	fmt.Println(result)
}