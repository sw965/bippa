package main

import (
	"fmt"
	"math"
	"math/rand"
	"github.com/sw965/omw"
	"github.com/sw965/crow/model"
	bp "github.com/sw965/bippa"
	"github.com/sw965/crow/tensor"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/maps"
	sb "github.com/sw965/bippa/battle/single"
	"github.com/sw965/crow/mlfuncs"
)

func GetFeature(bt *sb.Battle) tensor.D2 {
	names := bt.P1Fighters.Names()
	ks := maps.Keys(bt.P2Fighters[0].Moveset)
	feature := make(tensor.D1, 0, len(bp.ALL_MOVE_NAMES) + len(bp.ALL_POKE_NAMES))
	for _, name := range bp.ALL_MOVE_NAMES {
		if slices.Contains(ks, name) {
			feature = append(feature, 1.0)
		} else {
			feature = append(feature, 0.0)
		}
	}

	for _, name := range bp.ALL_POKE_NAMES {
		if slices.Contains(names[1:], name) {
			feature = append(feature, 1.0)		
		} else {
			feature = append(feature, 0.0)
		}
	}
	ret := make(tensor.D2, len(bp.ALL_MOVE_NAMES) + len(bp.ALL_POKE_NAMES))
	for i := 0; i < len(ret); i++ {
		ret[i] = feature.Clone()
	}
	return ret
}

func GetLabel(ps map[sb.Action]float64) tensor.D1 {
	label := make(tensor.D1, len(bp.ALL_MOVE_NAMES) + len(bp.ALL_POKE_NAMES))
	for i, name := range bp.ALL_MOVE_NAMES {
		action := sb.Action{CmdMoveName:name, SwitchPokeName:bp.EMPTY_POKE_NAME, IsPlayer1:true}
		p, ok := ps[action]
		if ok {
			label[i] = p
		}
	}

	n := len(bp.ALL_MOVE_NAMES)
	for i, name := range bp.ALL_POKE_NAMES {
		action := sb.Action{CmdMoveName:bp.EMPTY_MOVE_NAME, SwitchPokeName:name, IsPlayer1:true}
		p, ok := ps[action]
		if ok {
			label[i+n] = p 
		}
	}

	sum := 0.0
	for i := range label {
		sum += label[i]
	}
	for i := range label {
		label[i] /= sum
	}

	return mlfuncs.D1SigmoidToTanh(label)
}

func MakeFighters(r *rand.Rand) sb.Fighters {
	pokemons := []bp.Pokemon{
		bp.NewTemplateBulbasaur(),
		bp.NewTemplateCharmander(),
		bp.NewTemplateSquirtle(),
		//bp.NewTemplateGarchomp(),
	}
	omw.ShuffleSlice(pokemons, r)
	ret := sb.Fighters{}
	for i := 0; i < sb.FIGHTER_NUM; i++ {
		ret[i] = pokemons[i]
	}
	return ret
}

func main() {
	r := omw.NewMt19937()
	n := len(bp.ALL_MOVE_NAMES) + len(bp.ALL_POKE_NAMES)
	w := map[bp.PokeName]tensor.D2{}
	for _, pokeName := range bp.ALL_POKE_NAMES {
		w[pokeName] = tensor.NewD2Zeros(n, n)
	}

	b := map[bp.PokeName]tensor.D1{}
	for _, pokeName := range bp.ALL_POKE_NAMES {
		b[pokeName] = tensor.NewD1Zeros(n)
	}

	linear := model.NewD2LinearSumTanhMSE(0.0001)
	mcts := sb.NewMCTS(r)

	num := 196
	simulation := 2560
	for i := 0; i < num; i++ {
		p1Fighters := MakeFighters(r)
		p2Fighters := MakeFighters(r)
		battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
		allNodes, err := mcts.Run(simulation, battle, math.Sqrt(2), r)
		if err != nil {
			panic(err)
		}
		ps := allNodes[0].PUCBManagers[0].TrialPercents()
		x := GetFeature(&battle)
		t := GetLabel(ps)
		linear.SetParam(w[battle.P1Fighters[0].Name], b[battle.P1Fighters[0].Name])
		err = linear.Train(x, t, 0.1)
		if err != nil {
			panic(err)
		}
	}

	testNum := 128
	for i := 0; i < testNum; i++ {
		p1Fighters := MakeFighters(r)
		p2Fighters := MakeFighters(r)
		battle := sb.Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
		fmt.Println(
			bp.POKE_NAME_TO_STRING[p1Fighters[0].Name],
			bp.POKE_NAME_TO_STRING[p1Fighters[1].Name],
			bp.POKE_NAME_TO_STRING[p1Fighters[2].Name],

			bp.POKE_NAME_TO_STRING[p2Fighters[0].Name],
			bp.POKE_NAME_TO_STRING[p2Fighters[1].Name],
			bp.POKE_NAME_TO_STRING[p2Fighters[2].Name],
		)
		x := GetFeature(&battle)
		linear.SetParam(w[p1Fighters[0].Name], b[p1Fighters[0].Name])
		y, err := linear.Predict(x)
		if err != nil {
			panic(err)
		}
		fmt.Println(y)
	}
}