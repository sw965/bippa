package main

import (
	"fmt"
	omwrand "github.com/sw965/omw/math/rand"
	"github.com/sw965/bippa/feature"
	"github.com/sw965/crow/model/1d"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/battle/dmgtools"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle/single/mcts"
	"net/http"
	"encoding/json"
	"github.com/sw965/bippa/battle/single/game"
)

type ResponseDataElement struct {
	EasyReadBattle single.EasyReadBattle
	Step single.Step
}

type ResponseData []ResponseDataElement

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

	mctSearch := mcts.New(
		&single.Context{
			DamageRandBonuses:dmgtools.RandBonuses{1.0},
			Rand:rg,
			Observer:single.EmptyObserver,
		},
	)
	mctSearch.Game.SetRandActionPlayer(rg)

	featureFunc := feature.NewSingleBattleFunc(feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)
	mctSearch.LeafNodeJointEvalFunc = mcts.NewStandardLeafNodeJointEvalFunc(affine, featureFunc)

	battle := single.Battle{
		SelfFighters:bp.Pokemons{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		OpponentFighters:bp.Pokemons{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		IsRealSelf:true,
	}

	responseData := make(ResponseData, 0, 64)
	observer := func(battle *single.Battle, step single.Step) {
		responseData = append(responseData, ResponseDataElement{EasyReadBattle:battle.ToEasyRead(), Step:step})
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		actionStr := r.URL.Query().Get("action")
		if actionStr == "-1" {
			w.Header().Set("Content-Type", "application/json")
			responseData = append(responseData, ResponseDataElement{EasyReadBattle:battle.ToEasyRead(), Step:-1})
			jsonResponse, err := json.Marshal(responseData)
			if err != nil {
				fmt.Println("json errorを呼び出してよ")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(jsonResponse)
			responseData = make(ResponseData, 0, 64)
			return 
		}

		action := single.StringToAction(actionStr, true)
		rootNode := mctSearch.NewNode(&battle)
		err := mctSearch.Run(5120, rootNode, rg)
		if err != nil {
			panic(err)
		}

		jointAction := rootNode.SeparateUCBManager.JointActionByMaxTrial(rg)
		jointAvg := rootNode.SeparateUCBManager.JointAverageValue()
		fmt.Println("joint", jointAction[0].ToString(), jointAction[1].ToString(), jointAvg)
		fmt.Fprint(w, fmt.Sprintf("%s %s %v", jointAction[0].ToString(), jointAction[1].ToString(), jointAvg))

		battle, err = game.NewPushFunc(
			&single.Context{
				DamageRandBonuses:dmgtools.RandBonuses{1.0},
				Rand:rg,
				Observer:observer,
			},
		)(battle, single.Actions{action, jointAction[0]})
		if err != nil {
			panic(err)
		}
		if isEnd, _ := game.IsEnd(&battle); isEnd {
			fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
			fmt.Println("ゲームが終了した")
			return
		}
		fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
		responseData = make(ResponseData, 0, 64)
	}

	http.HandleFunc("/caitlin/", handler)
	fmt.Println("サーバ建て")
	fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
	err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}