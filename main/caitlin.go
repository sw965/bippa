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
)

type ResponseData struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

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

	mctSearch := mcts.New(rg)
	mctSearch.Game.SetRandActionPlayer(rg)
	f := feature.NewSingleBattleFunc(feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)
	mctSearch.LeafNodeJointEvalFunc = mcts.NewLeafNodeJointEvalFunc(affine, f)

	battle := single.Battle{
		SelfFighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		OpponentFighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		IsRealSelf:true,
		RandDmgBonuses:dmgtools.RandBonuses{1.0},
		Observer:func(_ *single.Battle, _ single.Step) {},
	}
	
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		procType := r.URL.Query().Get("procType")
		switch procType {
			case "0":
				fmt.Fprint(w, "")
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

		battle, err = single.NewPushFunc(rg)(battle, single.Actions{action, jointAction[1]})
		if err != nil {
			panic(err)
		}
		if isEnd, _ := single.IsEnd(&battle); isEnd {
			fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
			fmt.Println("ゲームが終了した")
			return
		}
		fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
	}

	http.HandleFunc("/caitlin/", handler)
	fmt.Println("サーバ建て")
	fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
	err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}