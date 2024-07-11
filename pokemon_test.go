package bippa_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
)

func TestRomanStanGyarados(t *testing.T) {
	pokemon := bp.NewRomanStanGyarados()
	fmt.Println(pokemon)
}