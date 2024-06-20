package main

import (
	"fmt"
	"github.com/sw965/bippa/battle/single/game"
	"github.com/sw965/crow/mcts/duct"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/feature"
	"github.com/sw965/crow/model/1d"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/battle/dmgtools"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle/single/mcts"
	"net/http"
)

func main() {
	server := http.Server{
        Addr:":8080",
        Handler:nil,
    }

	rg := omwrand.NewMt19937()
	xn := 90
	u1n := 64
	u2n := 16
	yn := 1
	affine, variable := model1d.NewStandardAffine(xn, u1n, u2n, yn, 0.0001, 64.0, rg)
	param, err := model1d.LoadParamJSON("C:/Users/kuroko/Desktop/test.json")
	if err != nil {
		panic(err)
	}
	variable.SetParam(param)

	gm := game.New(rg)

	mctSearch := mcts.New(rg)
	mctSearch.Game.SetRandActionPlayer(rg)

	leafNodeJointEvalFunc := func(battle *single.Battle) (duct.LeafNodeJointEvalY, error) {
		if isEnd, gameRetJointVal := single.IsEnd(battle); isEnd {
			y := make(duct.LeafNodeJointEvalY, len(gameRetJointVal))
			for i, v := range gameRetJointVal {
				y[i] = v
			}
			return y, nil
		} else {
			x := feature.NewSingleBattleFunc(2, feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)(battle)
			y, err := affine.Predict(x)
			v := y[0]
			return duct.LeafNodeJointEvalY{v, 1.0-v}, err
		}
	}

	mctSearch.LeafNodeJointEvalFunc = leafNodeJointEvalFunc

	initBattle := single.Battle{
		SelfFighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		OpponentFighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		IsRealSelf:true,
		RandDmgBonuses:dmgtools.RandBonuses{1.0},
		Observer:func(_ *single.Battle, _ single.Step) {},
	}

	initBattle.OpponentFighters[0].CurrentHP = 0

	w.Header().Set("Access-Control-Allow-Origin", "*")

	hoge := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprintf(w, "神の宣告")
	}

	gm.Player = func(battle *single.Battle) (single.Actions, []float64, error) {
		action := single.StringToAction(r.URL.Query().Get("action"), false)
		fmt.Println("action", action.ToString())
		battle.Observer = func(_ *single.Battle, _ single.Step) {}
		jointAction, jointQ, err := mctSearch.NewPlayer(5120, rg)(battle)
		if err != nil {
			return single.Actions{}, []float64{}, err
		}
		fmt.Println("jointQ", jointQ)
		return single.Actions{jointAction[0], action}, []float64{}, nil
	}

	endBattle, _, _, _, err := gm.PlayoutWithHistory(initBattle, 128)
	if err != nil {
		panic(err)
	}
	isEnd, gameRetJointVal := single.IsEnd(&endBattle)
	fmt.Println("isEnd", isEnd, gameRetJointVal)

	http.HandleFunc("/hoge/", hoge)
	fmt.Println("サーバ建て")
	err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}