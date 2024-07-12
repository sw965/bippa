package bippa_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
)

func TestRomanStan2009Gyarados(t *testing.T) {
	pokemon := bp.NewRomanStan2009Gyarados()
	fmt.Println(pokemon)
}