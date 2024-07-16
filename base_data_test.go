package bippa_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
)

func TestPokedex(t *testing.T) {
	for pokeName, pokeData := range bp.POKEDEX {
		fmt.Println(pokeName, pokeData)
	}
}

func TestMovedex(t *testing.T) {
	for moveName, moveData := range bp.MOVEDEX {
		fmt.Println(moveName, moveData)
	}
}

func TestNaturedex(t *testing.T) {
	for nature, natureData := range bp.NATUREDEX {
		fmt.Println(nature.ToString(), natureData)
	}
}

func TestTypedex(t *testing.T) {
	for atkType, defTypeData := range bp.TYPEDEX {
		for defType, effect := range defTypeData {
			fmt.Println(atkType, defType, effect)
		}
	}
}