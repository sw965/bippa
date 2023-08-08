package bippa_test

import (
	"testing"
	bp "github.com/sw965/bippa"
	"fmt"
)

func TestPOKEDEX(t *testing.T) {
	for pokeName, pokeData := range bp.POKEDEX {
		fmt.Println(pokeName, pokeData)
	}
}