package main

import (
    "fmt"
	"net/http"
	"net/url"
    bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle"
	"github.com/sw965/bippa/battle/game"
    "encoding/json"
)

func dataQueryHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	requestType := r.URL.Query().Get("request_type")
    fmt.Println(requestType, "を送信しました。")

    var response []byte
    var err error
    switch requestType {
        case "all_poke_names":
            response, err = json.Marshal(bp.ALL_POKE_NAMES.ToStrings())
        case "pokedex":
            response, err = json.Marshal(bp.POKEDEX.ToEasyRead())
        case "all_move_names":
            response, err = json.Marshal(bp.ALL_MOVE_NAMES.ToStrings())
        case "movedex":
            response, err = json.Marshal(bp.MOVEDEX.ToEasyRead())
        case "all_natures":
            response, err = json.Marshal(bp.ALL_NATURES.ToStrings())
        case "naturedex":
            response, err = json.Marshal(bp.NATUREDEX.ToEasyRead())
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(response)
}

var aiTrainerTitle string
var aiTrainerName string
var playerLeadPokemons bp.Pokemons
var playerBenchPokemons bp.Pokemons
var aiLeadPokemons bp.Pokemons
var aiBenchPokemons bp.Pokemons
var battleManager battle.Manager

func battleInitHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
	
	aiTrainerTitle, err := url.QueryUnescape(r.URL.Query().Get("ai_trainer_title"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	aiTrainerName, err := url.QueryUnescape(r.URL.Query().Get("ai_trainer_name"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	fmt.Println("aiInfo", aiTrainerTitle, aiTrainerName)

	playerLeadPokemonsQuery:= r.URL.Query().Get("player_lead_pokemons")
	if playerLeadPokemonsQuery != "" {
		easyReadPlayerLeadPokemonsStr, err := url.QueryUnescape(playerLeadPokemonsQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println(easyReadPlayerLeadPokemonsStr)
		var easyReadPlayerLeadPokemons bp.EasyReadPokemons
		err = json.Unmarshal([]byte(easyReadPlayerLeadPokemonsStr), &easyReadPlayerLeadPokemons)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("通過")
		playerLeadPokemons, err = easyReadPlayerLeadPokemons.From()
		if err != nil {
			panic(err)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Println("通過2")
	}

	fmt.Println("playerLeadPokemons")
	fmt.Println(playerLeadPokemons.ToEasyRead())

	playerBenchPokemonsQuery:= r.URL.Query().Get("player_bench_pokemons")
	if playerBenchPokemonsQuery != "" {
		easyReadPlayerBenchPokemonsStr, err := url.QueryUnescape(playerBenchPokemonsQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var easyReadPlayerBenchPokemons bp.EasyReadPokemons
		err = json.Unmarshal([]byte(easyReadPlayerBenchPokemonsStr), &easyReadPlayerBenchPokemons)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		playerBenchPokemons, err = easyReadPlayerBenchPokemons.From()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	
	fmt.Println("playerBenchPokemons")
	fmt.Println(playerBenchPokemons.ToEasyRead())

	aiLeadPokemonsQuery:= r.URL.Query().Get("ai_lead_pokemons")
	if aiLeadPokemonsQuery != "" {
		easyReadAILeadPokemonsStr, err := url.QueryUnescape(aiLeadPokemonsQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var easyReadAILeadPokemons bp.EasyReadPokemons
		err = json.Unmarshal([]byte(easyReadAILeadPokemonsStr), &easyReadAILeadPokemons)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		aiLeadPokemons, err = easyReadAILeadPokemons.From()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("aiLeadPokemons")
	fmt.Println(aiLeadPokemons.ToEasyRead())

	aiBenchPokemonsQuery:= r.URL.Query().Get("ai_bench_pokemons")
	if aiBenchPokemonsQuery != "" {
		easyReadAIBenchPokemonsStr, err := url.QueryUnescape(aiBenchPokemonsQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var easyReadAIBenchPokemons bp.EasyReadPokemons
		err = json.Unmarshal([]byte(easyReadAIBenchPokemonsStr), &easyReadAIBenchPokemons)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		aiBenchPokemons, err = easyReadAIBenchPokemons.From()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	fmt.Println("aiBenchPokemons")
	fmt.Println(aiBenchPokemons.ToEasyRead())

	init, err := url.QueryUnescape(r.URL.Query().Get("init"))
	if init == "true" {
		if aiTrainerTitle == "" {
			response, err := json.Marshal("AIの肩書きが空文字なので、設定してください。(例：ポケモントレーナー、してんのう、チャンピオン、やまおとこ 等)")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(response)
			return
		}

		if aiTrainerName == "" {
			response, err := json.Marshal("AIのトレーナー名が空文字なので、設定してください。(例：リーフ、コトネ、ハルカ、ヒカリ、トウコ、メイ、セレナ 等)")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(response)
			return
		}

		n := len(playerLeadPokemons)
		if n != len(aiLeadPokemons) {
			response, err := json.Marshal("プレイヤーとAIの先頭のポケモンの数が違う為、シングルバトルなのかダブルバトルなのかを判断出来ません。")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(response)
			return
		}

		if n != battle.SINGLE && n != battle.DOUBLE {
			response, err := json.Marshal(fmt.Sprintf("プレイヤーの先頭のポケモンの数が %d匹である為、シングルバトルかダブルバトルを開始出来ません。先頭のポケモンの数は1匹か2匹にしてください。", n))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(response)
			return
		}

		battleManager = battle.Manager{
			CurrentSelfLeadPokemons:playerLeadPokemons,
			CurrentSelfBenchPokemons:playerBenchPokemons,
			CurrentOpponentLeadPokemons:aiLeadPokemons,
			CurrentOpponentBenchPokemons:aiBenchPokemons,
		}

		battleManager.Init(aiTrainerTitle, aiTrainerName)
		response, err := json.Marshal("バトルの初期化完了")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
	}
}

func battleQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	var response []byte
	requestType, err := url.QueryUnescape(r.URL.Query().Get("request_type"))
	switch requestType {
		case "battle":
			response, err = json.Marshal(battleManager.ToEasyRead())
		case "legal_separate_actions":
			response, err = json.Marshal(game.LegalSeparateActions(&battleManager).ToEasyRead())
	}

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	w.Write(response)
}

func battleCommandHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	requestType, err := url.QueryUnescape(r.URL.Query().Get("request_type"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	switch requestType {
		case "push":
			playerActionQuery := r.URL.Query().Get("player_action")
			if playerActionQuery == "" {
				response, err := json.Marshal("playerActionが空である為、pushを実行出来ません。")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(response)
				return 
			}

			easyReadPlayerActionStr, err := url.QueryUnescape(playerActionQuery)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var easyReadPlayerAction battle.EasyReadAction
			err = json.Unmarshal([]byte(easyReadPlayerActionStr), easyReadPlayerAction)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			playerAction, err := easyReadPlayerAction.From()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			aiActionQuery := r.URL.Query().Get("ai_action")
			if aiActionQuery == "" {
				response, err := json.Marshal("aiActionが空である為、pushを実行出来ません。")
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(response)
				return 
			}

			easyReadAIActionStr, err := url.QueryUnescape(aiActionQuery)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			var easyReadAIAction battle.EasyReadAction
			err = json.Unmarshal([]byte(easyReadAIActionStr), &easyReadAIAction)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			aiAction, err := easyReadAIAction.From()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			nextBattleManager, err := game.Push(battleManager, battle.Actions{playerAction, aiAction})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			battleManager = nextBattleManager
			response, err := json.Marshal(battleManager.ToEasyRead())
			w.Write(response)
	}
}

func main() {
	server := http.Server{
        Addr:":8080",
        Handler:nil,
    }

    http.HandleFunc("/data_query/", dataQueryHandler)
	http.HandleFunc("/battle_init/", battleInitHandle)
	http.HandleFunc("/battle_query/", battleQueryHandler)
	http.HandleFunc("/battle_command/", battleCommandHandler)

    fmt.Println("サーバーが起動しました。")
	err := server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}