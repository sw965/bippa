package bippa_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
)

func TestTeamPokemonModelPart(t *testing.T) {
	moveNames := bp.MoveNames{"ギガドレイン", "ヘドロばくだん"}
	pokemon, err := bp.NewPokemon(
		"フシギバナ", bp.MALE, "ひかえめ", "しんりょく",  "くろいヘドロ",
		moveNames, bp.NewMaxPowerPointUps(len(moveNames)),
		&bp.ALL_MAX_INDIVIDUAL_STATE, &bp.EffortState{HP:bp.MAX_EFFORT, SpAtk:bp.MAX_EFFORT},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(pokemon)
}