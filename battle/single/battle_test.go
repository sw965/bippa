package single_test

import (
	"fmt"
	"testing"
	"github.com/sw965/bippa/battle/single/game"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/battle/dmgtools"
	bp "github.com/sw965/bippa"
	omwrand "github.com/sw965/omw/math/rand"

)

func Test(t *testing.T) {
	battle := single.Battle{
		SelfFighters:bp.Pokemons{
			bp.NewTemplateZapdos(),
			bp.NewTemplateSquirtle(),
			bp.NewTemplateBulbasaur(),
		},
		OpponentFighters:bp.Pokemons{
			bp.NewTemplateGarchomp(),
			bp.NewTemplateSuicune(),
			bp.NewTemplateBulbasaur(),
		},
		IsRealSelf:true,
	}

	rg := omwrand.NewMt19937()
	context := single.Context{
		Observer:single.EmptyObserver,
		Rand:rg,
		DamageRandBonuses:dmgtools.RandBonuses{1.0},
	}
	push := game.NewPushFunc(&context)
	battle, err := push(battle, single.Actions{
		single.Action{CmdMoveName:bp.THUNDERBOLT, IsSelf:true},
		single.Action{CmdMoveName:bp.STONE_EDGE, IsSelf:false},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(
		"selfCurrentHP", battle.SelfFighters[0].CurrentHP,
		"opponetCurrentHP", battle.OpponentFighters[0].CurrentHP,
	)
}