package single_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle/single"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/battle/dmgtools"
)

func Test(t *testing.T) {
	r := omwrand.NewMt19937()
	p1Fighters := single.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
	}

	p2Fighters := single.Fighters{
		bp.NewTemplateSquirtle(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateBulbasaur(),
	}

	observer := func(battle *single.Battle, step single.Step) {
		switch step {
			case single.BEFORE_MOVE_USE_STEP:
				fmt.Println("ここ1")
			case single.AFTER_MOVE_USE_STEP:
				fmt.Println("ここ2")
			case single.BEFORE_MOVE_DAMAGE_STEP:
				fmt.Println("ここ3")
			case single.AFTER_MOVE_DAMAGE_STEP:
				fmt.Println("ここ4")
		}
	}

	battle := single.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters,  RandDmgBonuses:dmgtools.RandBonuses{1.0}, Observer:observer}
	push := single.NewPushFunc(r)
	nextBattle, err := push(battle, single.Actions{ single.Action{SwitchPokeName:bp.CHARMANDER, IsPlayer1:true}, single.Action{CmdMoveName:bp.WATER_GUN, IsPlayer1:false} })
	if err != nil {
		panic(err)
	}
	fmt.Println(nextBattle.P1Fighters[0].CurrentHP)
}