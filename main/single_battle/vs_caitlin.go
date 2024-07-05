package main

import (
	"fmt"
	omwrand "github.com/sw965/omw/math/rand"
	//omwslices "github.com/sw965/omw/slices"
	"github.com/sw965/bippa/feature"
	"github.com/sw965/crow/model/1d"
	"github.com/sw965/bippa/battle/single"
	"github.com/sw965/bippa/battle/dmgtools"
	bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle/single/mcts"
	"net/http"
	"encoding/json"
	"github.com/sw965/bippa/battle/single/game"
	//battlemsg "github.com/sw965/bippa/battle/msg"
	"net/url"
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
	param, err := model1d.LoadParamJSON("C:/Go/project/bippa/main/single_battle/nn_param.json")
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

	featureFunc := feature.NewSingleBattleFunc(feature.ExpectedDamageRatioToCurrentHP, feature.DPSRatioToCurrentHP)
	mctSearch.LeafNodeJointEvalFunc = mcts.NewLeafNodeJointEvalFunc(affine, featureFunc)

	strToTeam := func(teamStr string) bp.EasyReadPokemons {
		decodedTeamStr, err := url.QueryUnescape(teamStr)
		if err != nil {
			panic(err)
		}

		var team bp.EasyReadPokemons
		err = json.Unmarshal([]byte(decodedTeamStr), &team)
		if err != nil {
			panic(err)
		}
		return team
	}

	var ui single.ObserverUI
	var battle single.Battle
	var push func(single.Battle, single.Actions) (single.Battle, error)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		selfTeamStr := r.URL.Query().Get("self_team")
		opponentTeamStr := r.URL.Query().Get("opponent_team")

		//どちらが片方だけの場合、エラーが起きるようにしておく。
		if selfTeamStr != "null" && opponentTeamStr != "null" {
			selfTeam := strToTeam(selfTeamStr).From()
			opponentTeam := strToTeam(opponentTeamStr).From()
			fmt.Println("selfTeam", selfTeam, len(selfTeam))
			fmt.Println("opponentTeam", opponentTeam, len(opponentTeam))

			battle = single.Battle{
				SelfFighters:selfTeam[:3],
				OpponentFighters:opponentTeam[:3],
				IsRealSelf:true,
			}

			fmt.Println("battle", battle)
			initDisplayUIs, err := single.NewInitDisplayUIs("カトレア", &battle)
			if err != nil {
				panic(err)
			}

			response, err := json.Marshal(&initDisplayUIs)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(response)
			
			ui, err = single.NewObserverUI(&battle)
			if err != nil {
				panic(err)
			}
			ui.SelfTrainerName = "ユウリ"
			ui.OpponentTrainerName = "カトレア"
			push = game.NewPushFunc(
				&single.Context{
					DamageRandBonuses:dmgtools.RandBonuses{1.0},
					Rand:rg,
					Observer:ui.Observer,
				},
			)
			return
		}

		rootNode := mctSearch.NewNode(&battle)
		err := mctSearch.Run(5120, rootNode, rg)
		if err != nil {
			panic(err)
		}

		jointAction := rootNode.SeparateUCBManager.JointActionByMaxTrial(rg)
		jointAvg := rootNode.SeparateUCBManager.JointAverageValue()
		fmt.Println("joint", jointAction[0].ToString(), jointAction[1].ToString(), jointAvg)

		actionStr := r.URL.Query().Get("action")
		fmt.Println("actionStr", actionStr)
		action := single.StringToAction(actionStr, true)
		fmt.Println("jointAction", action, jointAction[1])

		battle, err = push(battle, single.Actions{action, jointAction[1]})
		if err != nil {
			panic(err)
		}

		if isEnd, _ := game.IsEnd(&battle); isEnd {
			fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
			fmt.Println("ゲームが終了した")
			return
		}

		response, err := json.Marshal(ui.Displays)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
		ui.Displays = make(single.DisplayUIs, 0, 128)
	}

	http.HandleFunc("/caitlin/", handler)
	fmt.Println("サーバ建て")
	err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}