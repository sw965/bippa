package main

import (
	"fmt"
	"net/http"
	bpServer "github.com/sw965/bippa/server"
	"os"
	"strconv"
)

func main() {
	argsN := len(os.Args)
	ipAddr := 8080
	var err error

	if argsN == 2 {
		ipAddr, err = strconv.Atoi(os.Args[1])
	} else if argsN != 1 {
		err = fmt.Errorf("コマンドラインに指定できる引数は1つまでです。")
	}

	if err != nil {
		panic(err)
	}

	server := http.Server{
        Addr:fmt.Sprintf(":%d", ipAddr),
        Handler:nil,
    }

	http.HandleFunc("/data_query/", bpServer.HandleDataQuery)
	http.HandleFunc("/battle_init/", bpServer.HandleBattleInit)
	http.HandleFunc("/battle_query/", bpServer.HandleBattleQuery)
	http.HandleFunc("/battle_command/", bpServer.HandleBattleCommand)

	fmt.Println("Bippa Main Server が起動しました。")
	fmt.Println("ipAddr =", ipAddr)
	err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}