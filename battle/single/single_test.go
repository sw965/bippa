package single_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
	sb "github.com/sw965/bippa/battle/single"
	"github.com/sw965/omw"
)

func Test(t *testing.T) {
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

	r := omw.NewMt19937()
	battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
	push := sb.Push(r)
	battle, err := push(
		battle,
		sb.Actions{
			sb.Action{CmdMoveName:bp.TACKLE, IsPlayer1:true},
			sb.Action{CmdMoveName:bp.WATER_GUN, IsPlayer1:false},
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(battle.P1Fighters[0].CurrentHP, battle.P2Fighters[0].CurrentHP)
}

func TestMCTS(t *testing.T) {
	r := omw.NewMt19937()
	mcts := sb.NewMCTS(r)
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

	allNodes, err := mcts.Run(1960, battle, 1.41, r)
	if err != nil {
		panic(err)
	}
	for _, as := range allNodes[0].MaxTrialActionsPath(r, 8) {
		fmt.Println(as)
	}
}