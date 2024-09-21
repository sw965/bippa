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

func parseAndConvQuery[T, U any](query string, from func(U) (T, error)) (T, error) {
	decoded, err := url.QueryUnescape(query)
	if err != nil {
		var t T
		return t, err
	}
	var u U
	err = json.Unmarshal([]byte(decoded), &u)
	if err != nil {
		var t T
		return t, err
	}

	t, err := from(u)
	return t, err
}

func dataQueryHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	dataType := r.URL.Query().Get("data_type")
    fmt.Println(dataType, "を送信しました。")

    var response []byte
    var err error
    switch dataType {
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
		case "all_items":
			response, err = json.Marshal(bp.ALL_ITEMS.ToStrings())
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

	if aiTrainerTitle == "" {
		err := fmt.Errorf(fmt.Sprintf("ai_trainer_titleが空である為、バトルを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	aiTrainerName, err := url.QueryUnescape(r.URL.Query().Get("ai_trainer_name"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	if aiTrainerName == "" {
		err := fmt.Errorf(fmt.Sprintf("ai_trainer_nameが空である為、バトルを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	from := func(es bp.EasyReadPokemons) (bp.Pokemons, error) { return es.From() }

	playerLeadPokemonsQuery := r.URL.Query().Get("player_lead_pokemons")
	if playerLeadPokemonsQuery == "" {
		err := fmt.Errorf(fmt.Sprintf("player_lead_pokemonsが空である為、バトルを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	playerLeadPokemons, err := parseAndConvQuery(playerLeadPokemonsQuery, from)

	playerBenchPokemonsQuery:= r.URL.Query().Get("player_bench_pokemons")
	if playerBenchPokemonsQuery == "" {
		err := fmt.Errorf(fmt.Sprintf("player_bench_pokemonsが空である為、バトルを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	playerBenchPokemons, err := parseAndConvQuery(playerBenchPokemonsQuery, from)

	aiLeadPokemonsQuery:= r.URL.Query().Get("ai_lead_pokemons")
	if aiLeadPokemonsQuery == "" {
		err := fmt.Errorf(fmt.Sprintf("ai_lead_pokemonsが空である為、バトルを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	aiLeadPokemons, err := parseAndConvQuery(aiLeadPokemonsQuery, from)

	aiBenchPokemonsQuery:= r.URL.Query().Get("ai_bench_pokemons")
	if aiBenchPokemonsQuery == "" {
		err := fmt.Errorf(fmt.Sprintf("ai_bench_pokemonsが空である為、バトルを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	aiBenchPokemons, err := parseAndConvQuery(aiBenchPokemonsQuery, from)

	battleManager = battle.Manager{
		CurrentSelfLeadPokemons:playerLeadPokemons,
		CurrentSelfBenchPokemons:playerBenchPokemons,
		CurrentOpponentLeadPokemons:aiLeadPokemons,
		CurrentOpponentBenchPokemons:aiBenchPokemons,
	}

	ms := make(battle.Managers, 0, 128)
	battle.GlobalContext.Observer = func(m *battle.Manager) {
		ms = append(ms, m.Clone())
	}
	battleManager.Init(aiTrainerTitle, aiTrainerName)

	response, err := json.Marshal(ms.ToEasyRead())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func battleQueryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	var response []byte
	queryType, err := url.QueryUnescape(r.URL.Query().Get("query_type"))
	switch queryType {
		case "battle":
			response, err = json.Marshal(battleManager.ToEasyRead())
		case "separate_legal_actions":
			response, err = json.Marshal(game.SeparateLegalActions(&battleManager).ToEasyRead())
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

	commandType, err := url.QueryUnescape(r.URL.Query().Get("command_type"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	switch commandType {
		case "push":
			from := func(e *battle.EasyReadAction) (battle.Action, error) { return e.From() }

			playerActionQuery := r.URL.Query().Get("player_action")
			if playerActionQuery == "" {
				err := fmt.Errorf("player_actionが空である為、pushを実行出来ません。")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			playerAction, err := parseAndConvQuery(playerActionQuery, from)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return			
			}
			for i := range playerAction {
				playerAction[i].IsCurrentSelf = true
			}

			aiActionQuery := r.URL.Query().Get("ai_action")
			if aiActionQuery == "" {
				err := fmt.Errorf("ai_actionが空である為、pushを実行出来ません。")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			aiAction, err := parseAndConvQuery(aiActionQuery, from)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return			
			}

			fmt.Println("actionQuery", playerAction.ToEasyRead(), aiAction.ToEasyRead())

			ms := make(battle.Managers, 0, 256)
			ms = append(ms, battleManager.Clone())
			battle.GlobalContext.Observer = func(m *battle.Manager) {
				c := m.Clone()
				if !c.CurrentSelfIsHost {
					c.SwapView()
				}
				ms = append(ms, c)
			}
			
			nextBattleManager, err := game.Push(battleManager, battle.Actions{playerAction, aiAction})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			battle.GlobalContext.Observer = battle.EmptyObserver
			battleManager = nextBattleManager
			response, err := json.Marshal(ms.ToEasyRead())
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