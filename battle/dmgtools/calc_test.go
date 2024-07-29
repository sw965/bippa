package dmgtools_test

import (
	"testing"
	"fmt"
	"github.com/sw965/bippa/battle/dmgtools"
	bp "github.com/sw965/bippa"

)

func TestGarchompZapdos(t *testing.T) {
	attacker := bp.NewTemplateGarchomp()
	defender := bp.NewTemplateZapdos()

	attackerInfo := dmgtools.Attacker{
		PokeName:attacker.Name,
		Level:attacker.Level,
		Atk:attacker.Atk,
		SpAtk:attacker.SpAtk,
	}

	defenderInfo := dmgtools.Defender{
		PokeName:defender.Name,
		Level:defender.Level,
		Def:defender.Def,
		SpDef:defender.SpDef,
	}

	calculator := dmgtools.Calculator{Attacker:attackerInfo, Defender:defenderInfo}
	result := calculator.Calculation(bp.STONE_EDGE, 1.0)
	if result != 156 {
		fmt.Println("attacker", attacker)
		fmt.Println("defender", defender)
		t.Errorf("テスト失敗")
	}
}