package single_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
	sb "github.com/sw965/bippa/battle/single"
	orand "github.com/sw965/omw/rand"
	"github.com/sw965/bippa/dmgtools"
	"github.com/sw965/crow/tensor"
	"github.com/sw965/crow/model"
	"github.com/sw965/crow/mcts/duct"
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

	r := orand.NewMt19937()
	battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
	push := sb.Push(dmgtools.RandBonuses{1.0}, r)
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
	battle.P1Fighters[0].CurrentHP = 0
	fmt.Println(battle.P1Fighters[0].CurrentHP, battle.P2Fighters[0].CurrentHP)
}

func Test2(t *testing.T) {
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

	p1Fighters[0].CurrentHP = 0
	battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
	legalActionss := sb.LegalActionss(&battle)
	fmt.Println(legalActionss)
}

func TestActionData(t *testing.T) {
	r := orand.NewMt19937()
	p1Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
	}

	p2Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateSquirtle(),
		bp.NewTemplateCharmander(),
	}

	// p1Fighters[0].CurrentHP = 1
	// p1Fighters[1].CurrentHP = 0
	// p1Fighters[2].CurrentHP = 0

	// p2Fighters[0].CurrentHP = 1
	// p2Fighters[1].CurrentHP = 0
	// p2Fighters[2].CurrentHP = 0

	battle := sb.Battle{
		P1Fighters:p1Fighters,
		P2Fighters:p2Fighters,
	}
	battle.P1Fighters[2].CurrentHP = 0
	battle.P2Fighters[2].CurrentHP = 0

	xn := len(bp.ALL_TYPES) * 4 * sb.FIGHTER_NUM * 2
	h1 := 64
	h2 := 16
	yn := 1
	affine, _ := model.NewThreeLayerAffineLeakyReLUInput1DOutputSigmoid1D(xn, h1, h2, yn, 0.0001, 64.0, r)

	mcts := sb.NewMCTS(dmgtools.RAND_BONUSES, r)
	mcts.LeafNodeEvalsFunc = func(battle *sb.Battle) duct.LeafNodeEvalYs {
		b1 := battle.P1Fighters.IsAllFaint()
		b2 := battle.P2Fighters.IsAllFaint()
		if b1 && b2 {
			return duct.LeafNodeEvalYs{0.5, 0.5}
		}

		if b1 {
			return duct.LeafNodeEvalYs{0.0, 1.0}
		}

		if b2 {
			return duct.LeafNodeEvalYs{1.0, 0.0}
		}

		x := battle.IndexFeature(15000)
		v, err := affine.Predict(x)
		if err != nil {
			panic(err)
		}
		return duct.LeafNodeEvalYs{duct.LeafNodeEvalY(v[0]), 1.0 - duct.LeafNodeEvalY(v[0])}
	}

	trainXs := make(tensor.D2, 0, 1280)
	trainYs := make(tensor.D2, 0, 1280)
	battles := make([]sb.Battle, 0, 1280)
	jointActions := make(sb.Actionss, 0, 1280)

	game := mcts.Game.Clone()
	game.Player = func(b *sb.Battle) (sb.Actions, error) {
		rootNode := mcts.NewNode(b)
		err := mcts.Run(128, rootNode, r)

		trainXs = append(trainXs, b.IndexFeature(15000))
		trainYs = append(trainYs, tensor.D1{rootNode.UCBManagers[0].AverageValue()})
		battles = append(battles, *b)

		swaped := b.SwapPlayers()
		trainXs = append(trainXs, swaped.IndexFeature(15000))
		trainYs = append(trainYs, tensor.D1{rootNode.UCBManagers[1].AverageValue()})
		battles = append(battles, swaped)

		jointAction := rootNode.UCBManagers.JointActionByMaxTrial(r)
		jointActions = append(jointActions, jointAction)
		jointActions = append(jointActions, sb.Actions{jointAction[1], jointAction[0]})
		return jointAction, err
	}

	for i := 0; i < 12800; i++ {
		_, err := game.Playout(battle)
		if err != nil {
			panic(err)
		}
		if i%10 == 0 {
			fmt.Println("ok")
			n := len(trainXs)
			for j := 0; j < n; j++ {
				idx := r.Intn(n)
				affine.SGD(trainXs[idx], trainYs[idx], 0.01)
			}
			for j, b := range battles {
				fmt.Println(
					"j =", j,
					bp.POKE_NAME_TO_STRING[b.P1Fighters[0].Name],
					b.P1Fighters[0].CurrentHP,
					bp.MOVE_NAME_TO_STRING[jointActions[j][0].CmdMoveName],
					bp.POKE_NAME_TO_STRING[jointActions[j][0].SwitchPokeName],

					bp.POKE_NAME_TO_STRING[b.P2Fighters[0].Name],
					b.P2Fighters[0].CurrentHP,
					bp.MOVE_NAME_TO_STRING[jointActions[j][1].CmdMoveName],
					bp.POKE_NAME_TO_STRING[jointActions[j][1].SwitchPokeName],
					"v =", trainYs[j],
				)
			}
			trainXs = make(tensor.D2, 0, 1280)
			trainYs = make(tensor.D2, 0, 1280)
			battles = make([]sb.Battle, 0, 1280)
			jointActions = make(sb.Actionss, 0, 1280)
		}
	}
}

func Test3(t *testing.T) {
	r := orand.NewMt19937()
	p1Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
	}

	p2Fighters := sb.Fighters{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateSquirtle(),
		bp.NewTemplateCharmander(),
	}
	p1Fighters[0].CurrentHP = 1
	battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}

	xn := len(bp.ALL_TYPES) * 4 * sb.FIGHTER_NUM * 2
	h1 := 64
	h2 := 16
	yn := 1

	affine, _ := model.NewThreeLayerAffineLeakyReLUInput1DOutputSigmoid1D(xn, h1, h2, yn, 0.0001, 64.0, r)
	fmt.Println(len(battle.IndexFeature(15000)), xn)
	swaped := battle.SwapPlayers()
	for i := 0; i < 5120; i++ {
		affine.SGD(battle.IndexFeature(15000), tensor.D1{0.1}, 0.01)
		affine.SGD(swaped.IndexFeature(1500), tensor.D1{0.9}, 0.01)
	}
	fmt.Println(affine.Predict(battle.IndexFeature(15000)))
	fmt.Println(affine.Predict(swaped.IndexFeature(15000)))
}