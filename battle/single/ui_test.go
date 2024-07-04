package single_test

import (
	"testing"
	"fmt"
	"github.com/sw965/bippa/battle/single"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle/single/game"
	omwrand "github.com/sw965/omw/math/rand"
)

func TestUI(t *testing.T) {
	battle := single.Battle{
		SelfFighters:bp.Pokemons{
			bp.NewTemplateGarchomp(),
			bp.NewTemplateCharmander(),
			bp.NewTemplateBulbasaur(),
		},
		OpponentFighters:bp.Pokemons{
			bp.NewTemplateSuicune(),
			bp.NewTemplateBulbasaur(),
			bp.NewTemplateGarchomp(),
		},
		IsRealSelf:true,
	}
	ui := single.UI{}
	rg := omwrand.NewMt19937()
	context := single.NewContext(rg)
	context.Observer = ui.Observer

	push := game.NewPushFunc(&context)
	actions := single.Actions{
		single.Action{CmdMoveName:bp.STONE_EDGE, IsSelf:true},
		single.Action{CmdMoveName:bp.SURF, IsSelf:false},
	}
	battle, err := push(battle, actions)
	if err != nil {
		t.Errorf(fmt.Sprintf("%v", err))
	}
	fmt.Println(battle)
}