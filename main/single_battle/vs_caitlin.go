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
	//"github.com/sw965/bippa/battle/single/game"
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

	battle := single.Battle{
		SelfFighters:bp.Pokemons{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		OpponentFighters:bp.Pokemons{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		IsRealSelf:true,
	}

	fmt.Println(battle)

	// selfViewLastBattle := battle
	// opponentViewLastBattle := battle.SwapPlayers()

	// observer := func(battle *single.Battle, step single.Step) {
	// 	var lastBattle single.Battle
	// 	if battle.IsRealSelf {
	// 		lastBattle = selfViewLastBattle
	// 	} else {
	// 		lastBattle = opponentViewLastBattle
	// 	}

	// 	switch step {
	// 		case single.AFTER_MOVE_USE_STEP:
	// 			var lastUsedMoveName bp.MoveName
	// 			lastMoveset := lastBattle.SelfFighters[0].Moveset
	// 			currentMoveset := battle.SelfFighters[0].Moveset
	// 			for moveName, pp := range currentMoveset {
	// 				if pp.Current != lastMoveset[moveName].Current {
	// 					lastUsedMoveName = moveName
	// 					break
	// 				}
	// 			}
	// 			operationData = append(operationData, []string{"battleMsg", "clear"})
	// 			for _, msg := range battlemsg.NewMoveUse(battle.SelfFighters[0].Name, lastUsedMoveName, battle.IsRealSelf).ToSlice() {
	// 				operationData = append(operationData, []string{"battleMsg", "add", msg})
	// 			}
	// 		case single.AFTER_OPPONENT_DAMAGE_STEP:
	// 			lastCurrentHP := lastBattle.OpponentFighters[0].CurrentHP
	// 			dmg := lastCurrentHP - battle.OpponentFighters[0].CurrentHP
	// 			lastMsg := omwslices.End[UIs](responseData).BattleMessage
	// 			for i := 1; i < dmg; i++ {
	// 				lastBattle = lastBattle.Clone()
	// 				lastBattle.OpponentFighters[0].CurrentHP -= dmg
	// 				responseData = append(responseData, NewUI(&lastBattle, lastMsg, lastBattle.IsRealSelf))
	// 			}
	// 		case single.AFTER_SELF_FAINT_STEP:
	// 			for _, msg := range battlemsg.NewFaint(battle.SelfFighters[0].Name, lastBattle.IsRealSelf).Accumulate() {
	// 				lastBattle = lastBattle.Clone()
	// 				responseData = append(responseData, NewUI(&lastBattle, msg, lastBattle.IsRealSelf))
	// 			}
	// 		case single.AFTER_OPPONENT_FAINT_STEP:
	// 			for _, msg := range battlemsg.NewFaint(battle.OpponentFighters[0].Name, !lastBattle.IsRealSelf).Accumulate() {
	// 				lastBattle = lastBattle.Clone()
	// 				responseData = append(responseData, NewUI(&lastBattle, msg, lastBattle.IsRealSelf))
	// 			}
	// 	}

	// 	if battle.IsRealSelf {
	// 		selfViewLastBattle = battle.Clone()
	// 	} else {
	// 		opponentViewLastBattle = battle.Clone()
	// 	}
	// }

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

	var selfTeam bp.EasyReadPokemons
	var opponentTeam bp.EasyReadPokemons

	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		selfTeamStr := r.URL.Query().Get("self_team")
		opponentTeamStr := r.URL.Query().Get("opponent_team")

		//どちらが片方だけの場合、エラーが起きるようにしておく。
		if selfTeamStr != "" && opponentTeamStr != "" {
			selfTeam = strToTeam(selfTeamStr)
			opponentTeam = strToTeam(opponentTeamStr)
			return
		}

		responseData, err := json.Marshal([][]string{[]string{"たいあたり", "神様"}, []string{"じしん", "ユウリちゃん"}})
		if err != nil {
			panic(err)
		}
		w.Write(responseData)
		// if r.URL.Query().Get("init") ==  "true" {
		// 	responseData = append(responseData, ResponseDataElement{EasyReadBattle:battle.ToEasyRead(), Step:-1})
		// 	jsonResponse, err := json.Marshal(responseData)
		// 	if err != nil {
		// 		http.Error(w, err.Error(), http.StatusInternalServerError)
		// 		return
		// 	}
		// 	w.Write(jsonResponse)
		// 	responseData = make(ResponseData, 0, 64)
		// 	return
		// }

		// action := single.StringToAction(actionStr, true)
		// fmt.Println("action", action)
		// rootNode := mctSearch.NewNode(&battle)
		// err := mctSearch.Run(5120, rootNode, rg)
		// if err != nil {
		// 	panic(err)
		// }

		// jointAction := rootNode.SeparateUCBManager.JointActionByMaxTrial(rg)
		// jointAvg := rootNode.SeparateUCBManager.JointAverageValue()
		// fmt.Println("joint", jointAction[0].ToString(), jointAction[1].ToString(), jointAvg)
		// battle, err = game.NewPushFunc(
		// 	&single.Context{
		// 		DamageRandBonuses:dmgtools.RandBonuses{1.0},
		// 		Rand:rg,
		// 		Observer:observer,
		// 	},
		// )(battle, single.Actions{action, jointAction[1]})
		// if err != nil {
		// 	panic(err)
		// }
		// if isEnd, _ := game.IsEnd(&battle); isEnd {
		// 	fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
		// 	fmt.Println("ゲームが終了した")
		// 	return
		// }
		// fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
		// jsonResponse, err := json.Marshal(responseData)
		// if err != nil {
		// 	fmt.Println("json errorを呼び出してよ")
		// 	http.Error(w, err.Error(), http.StatusInternalServerError)
		// 	return
		// }
		// w.Write(jsonResponse)
		// responseData = make(ResponseData, 0, 64)
	}

	http.HandleFunc("/caitlin/", handler)
	fmt.Println("サーバ建て")
	fmt.Println(battle.SelfFighters[0].Name.ToString(), battle.SelfFighters[0].CurrentHP, battle.OpponentFighters[0].Name.ToString(), battle.OpponentFighters[0].CurrentHP)
	err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}