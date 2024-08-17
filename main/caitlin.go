package main

import (
    "fmt"
	"net/http"
	bp "github.com/sw965/bippa"
    "github.com/sw965/bippa/battle"
    "encoding/json"
	"net/url"
)

func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    selfTeam := r.URL.Query().Get("selfTeam")
	decodedStr, err := url.QueryUnescape(selfTeam)
	if err != nil {
		panic(err)
	}
	fmt.Println(decodedStr)

	opponentTeam := r.URL.Query().Get("opponentTeam")
	fmt.Println("opponentTeam", opponentTeam)

	battle := battle.Manager{
		CurrentSelfLeadPokemons:bp.Pokemons{
			bp.NewKusanagi2009Empoleon(),
			bp.NewKusanagi2009Toxicroak(),
		},

		CurrentOpponentLeadPokemons:bp.Pokemons{
			bp.NewMoruhu2007Bronzong(),
			bp.NewMoruhu2007Metagross(),
		},
	}
	response, err := json.Marshal(battle)
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
    http.HandleFunc("/caitlin/", handler)
    fmt.Println("サーバーを建てたよ")
	err := server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}
