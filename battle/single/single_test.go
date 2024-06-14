package single_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/battle/dmgtools"
	omwrand "github.com/sw965/omw/math/rand"
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
	battle := single.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
	fmt.Println(battle.P1Fighters[1].CurrentHP)
	push := single.NewPushFunc(dmgtools.RandBonuses{1.0}, r)
	nextBattle, err := push(battle, single.Actions{ single.Action{SwitchPokeName:bp.CHARMANDER, IsPlayer1:true}, single.Action{CmdMoveName:bp.WATER_GUN, IsPlayer1:false} })
	if err != nil {
		panic(err)
	}
	fmt.Println(nextBattle.P1Fighters[0].CurrentHP)
}