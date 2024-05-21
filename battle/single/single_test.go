package single_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
	sb "github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/dmgtools"
	orand "github.com/sw965/omw/rand"
)

func Test(t *testing.T) {
	r := orand.NewMt19937()
	p1Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
	}

	p2Fighters := sb.Fighters{
		bp.NewTemplateSquirtle(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateBulbasaur(),
	}
	battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
	fmt.Println(battle.P1Fighters[1].CurrentHP)
	push := sb.Push(dmgtools.RandBonuses{1.0}, r)
	nextBattle, err := push(battle, sb.Actions{ sb.Action{SwitchPokeName:bp.CHARMANDER, IsPlayer1:true}, sb.Action{CmdMoveName:bp.WATER_GUN, IsPlayer1:false} })
	if err != nil {
		panic(err)
	}
	fmt.Println(nextBattle.P1Fighters[0].CurrentHP)
}