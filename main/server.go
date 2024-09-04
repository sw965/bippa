package main

import (
    "fmt"
	"net/http"
	"net/url"
    bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle"
    "encoding/json"
)

func dataHandler(w http.ResponseWriter, r *http.Request) {
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
    }
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(response)
}

var battleManager battle.Manager

func battleHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
	battleType, err := url.QueryUnescape(r.URL.Query().Get("battle_type"))
	//パニックじゃなくて、それように変える
	if err != nil {
		panic(err)
	}

	playerPokemonsStr, err := url.QueryUnescape(r.URL.Query().Get("player_pokemons"))
	if err != nil {
		panic(err)
	}

	caitlinPokemonsStr, err := url.QueryUnescape(r.URL.Query().Get("caitlin_pokemons"))
	if err != nil {
		panic(err)
	}

	var easyReadPlayerPokemons bp.EasyReadPokemons
	var easyReadCaitlinPokemons bp.EasyReadPokemons

	err = json.Unmarshal([]byte(playerPokemonsStr), &easyReadPlayerPokemons)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal([]byte(caitlinPokemonsStr), &easyReadCaitlinPokemons)
	if err != nil {
		panic(err)
	}

	playerPokemons, err := easyReadPlayerPokemons.From()
	if err != nil {
		panic(err)
	}

	caitlinPokemons, err := easyReadCaitlinPokemons.From()
	if err != nil {
		panic(err)
	}

	switch battleType {
		case "single":
			battleManager = battle.Manager{
				CurrentSelfLeadPokemons:playerPokemons[:1],
				CurrentSelfBenchPokemons:playerPokemons[1:],
				CurrentOpponentLeadPokemons:caitlinPokemons[:1],
				CurrentOpponentBenchPokemons:caitlinPokemons[1:],
			}
		case "double":
			battleManager = battle.Manager{
				CurrentSelfLeadPokemons:playerPokemons[:2],
				CurrentSelfBenchPokemons:playerPokemons[2:],
				CurrentOpponentLeadPokemons:caitlinPokemons[:2],
				CurrentOpponentBenchPokemons:caitlinPokemons[2:],
			}
	}

	fmt.Println(battleManager)

	response, err := json.Marshal("バトルタイプとポケモンの情報を受け取ったよ" + battleType)
	if err != nil {
		panic(err)
	}
    w.Write(response)
}

func main() {
	server := http.Server{
        Addr:":8080",
        Handler:nil,
    }

    http.HandleFunc("/data/", dataHandler)
	http.HandleFunc("/battle/", battleHandler)

    fmt.Println("サーバーが起動しました。")
	err := server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}