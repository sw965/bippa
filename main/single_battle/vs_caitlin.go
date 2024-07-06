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

func stringToEasyReadPokemons(pokemonsStr string) bp.EasyReadPokemons {
	decodedPokemonsStr, err := url.QueryUnescape(pokemonsStr)
	if err != nil {
		panic(err)
	}

	var pokemons bp.EasyReadPokemons
	err = json.Unmarshal([]byte(decodedPokemonsStr), &pokemons)
	if err != nil {
		panic(err)
	}
	return pokemons
}

func stringToAction(actionStr string, isSelf bool) single.Action {
	decodedActionStr, err := url.QueryUnescape(actionStr)
	if err != nil {
		panic(err)
	}
	return single.StringToAction(decodedActionStr, true)
}

type ResponseBattle struct {
	SelfFighters bp.EasyReadPokemons
	OpponentFighters bp.EasyReadPokemons
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

	var ui single.ObserverUI
	var battle single.Battle
	var push func(single.Battle, single.Actions) (single.Battle, error)

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		selfTeamStr := r.URL.Query().Get("self_team")
		opponentTeamStr := r.URL.Query().Get("opponent_team")

		if selfTeamStr != "null" && opponentTeamStr != "null" {
			easyReadSelfPokemons := stringToEasyReadPokemons(selfTeamStr)
			easyReadOpponentPokemons := stringToEasyReadPokemons(opponentTeamStr)

			response, err := json.Marshal(&ResponseBattle{SelfFighters:easyReadSelfPokemons[:3], OpponentFighters:easyReadOpponentPokemons[:3]})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(response)

			selfFighters := easyReadSelfPokemons.From()[:3]
			opponentFighters := easyReadOpponentPokemons.From()[:3]

			battle = single.Battle{
				SelfFighters:selfFighters,
				OpponentFighters:opponentFighters,
				IsRealSelf:true,
			}
			
			ui, err = single.NewObserverUI(&battle, 128)
			if err != nil {
				panic(err)
			}
			ui.SelfTrainerName = "ユウリ"
			ui.OpponentTrainerName = "カトレア"

			context := &single.Context{
				DamageRandBonuses:dmgtools.RandBonuses{1.0},
				Rand:rg,
				Observer:ui.Observer,
			}
			push = game.NewPushFunc(context)
			ui.Displays = append(ui.Displays, single.NewDisplayUI(&battle, ""))
			return
		} else if selfTeamStr != "null" || opponentTeamStr != "null" {
			err := fmt.Errorf("どちらか一方のチームの情報しか渡されていません。両チームの情報は、同時に渡してください。")
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		action := stringToAction(actionStr, true)
		fmt.Println("jointAction", action, jointAction[1])

		battle, err = push(battle, single.Actions{action, jointAction[1]})
		if err != nil {
			panic(err)
		}

		for i, display := range ui.Displays {
			fmt.Println(i, display)
		}

		if isEnd, _ := game.IsEnd(&battle); isEnd {
			return
		}

		response, err := json.Marshal(ui.Displays)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
		ui.Displays = ui.Displays.LastElementSlice()
	}

	http.HandleFunc("/caitlin/", handler)
	fmt.Println("サーバ建て")
	err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}