package main

import (
    "fmt"
	"net/http"
    bp "github.com/sw965/bippa"
    "encoding/json"
)

func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")
    dataType := r.URL.Query().Get("data_type")
    fmt.Println("data_type", dataType)
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

func main() {
	server := http.Server{
        Addr:":8081",
        Handler:nil,
    }
    http.HandleFunc("/dawn/", handler)
	err := server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}