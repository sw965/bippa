package main

import (
	"fmt"
	"github.com/sw965/omw"
	bp "github.com/sw965/bippa"
)

func main() {
	r := omw.NewMt19937()
	learnset := bp.POKEDEX["フシギバナ"].Learnset
	modelPart := bp.NewCombination2FeatureTeamModelPart[bp.MoveNames, bp.Natures](learnset, bp.ALL_NATURES, r)
	for k1, m := range modelPart {
		for k2, v := range m {
			fmt.Println(k1, k2, v)
		}
	}
}