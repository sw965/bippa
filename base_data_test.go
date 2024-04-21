package bippa_test

import (
	"testing"
	"fmt"
	bippa "github.com/sw965/bippa"
)

func TestPokedex(t *testing.T) {
	for pokeName, pokeData := range bippa.POKEDEX {
		fmt.Println(pokeName, pokeData)
	}
}

func TestMovedex(t *testing.T) {
	for moveName, moveData := range bippa.MOVEDEX {
		fmt.Println(moveName, moveData)
	}
}

func TestTypedex(t *testing.T) {
	for atkType, defTypeData := range bippa.TYPEDEX {
		for defType, effect := range defTypeData {
			fmt.Println(atkType, defType, effect)
		}
	}
}