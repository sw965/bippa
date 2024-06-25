package main

import (
	"net/http"
    bp "github.com/sw965/bippa"
    "encoding/json"
)

func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    dataType := r.URL.Query().Get("data_type")
    var response []byte
    var err error
    switch dataType {
        case "pokedex":
            response, err = json.Marshal(bp.POKEDEX.ToEasyRead())
        case "movedex":
            response, err = json.Marshal(bp.MOVEDEX.ToEasyRead())
        case "all_poke_names":
            response, err = json.Marshal(bp.ALL_POKE_NAMES.ToStrings())
        case "all_move_names":
            response, err = json.Marshal(bp.ALL_MOVE_NAMES.ToStrings())
    }

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Write(response)
}

func main() {
	server := http.Server{
        Addr:":8080",
        Handler:nil,
    }
    http.HandleFunc("/dawn/", handler)
	err := server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}