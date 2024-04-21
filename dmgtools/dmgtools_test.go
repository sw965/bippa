package dmgtools_test

import (
	"testing"
	"fmt"
	"golang.org/x/exp/maps"
	"github.com/sw965/bippa/dmgtools"
	bp "github.com/sw965/bippa"

)

func TestCalculator(t *testing.T) {
	attacker, err := bp.NewPokemon(bp.CHARMANDER, bp.MoveNames{bp.EMBER})
	if err != nil {
		panic(err)
	}

	defender, err := bp.NewPokemon(bp.SQUIRTLE, bp.MoveNames{bp.WATER_GUN})
	if err != nil {
		panic(err)
	}
	fmt.Println(attacker)
	fmt.Println(defender)

	attackerInfo := dmgtools.Attacker{
		PokeName:attacker.Name,
		Level:attacker.Level,
		Atk:attacker.Atk,
		SpAtk:attacker.SpAtk,
		MoveName:maps.Keys(attacker.Moveset)[0],
	}

	defenderInfo := dmgtools.Defender{
		PokeName:defender.Name,
		Level:defender.Level,
		Def:defender.Def,
		SpDef:defender.SpDef,
	}

	calc := dmgtools.Calculator{Attacker:attackerInfo, Defender:defenderInfo}
	result := calc.Execute(1.0)
	fmt.Println(result)
}