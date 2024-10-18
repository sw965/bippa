package main

import (
	"fmt"
	"os"
	"net/http"
	"net/url"
	"encoding/json"
	"strconv"
	omwHTTP "github.com/sw965/omw/http"
	"github.com/sw965/crow/mcts/duct"
	"github.com/sw965/bippa/battle"
)

var mctSearcher duct.MCTS[battle.Manager, battle.ActionsSlice, battle.Actions, battle.Action] 
var c float64 = 5
var simulationNum int = 1960

func handleQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	queryType, err := url.QueryUnescape(r.URL.Query().Get("query_type"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	var response []byte
	switch queryType {
		case "joint_action":
			battleQuery := r.URL.Query().Get("battle")
			if battleQuery == "" {
				err := fmt.Errorf(fmt.Sprintf("battleが空である為、Battle MCTS Server から joint_action を取得出来ません。"))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			easyReadManager, err := omwHTTP.QueryToType[battle.EasyReadManager](battleQuery)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			manager, err := easyReadManager.From()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rootNode := mctSearcher.NewNode(&manager)
			fmt.Println(fmt.Sprintf("バトルのモンテカルロ木探索を実行します。試行回数 = %d", simulationNum))
			err = mctSearcher.Run(simulationNum, rootNode, battle.GlobalContext.Rand)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			fmt.Println("バトルのモンテカルロ木探索の実行が完了しました。")

			jointAction := rootNode.SeparateUCBManager.JointActionByMaxTrial(battle.GlobalContext.Rand)
			response, err = json.Marshal(jointAction.ToEasyRead())
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	}
	w.Write(response)
}

const (
	HOST = "host"
	GUEST = "guest"
)

func main() {
	argsN := len(os.Args)
	if argsN > 3 {
		err := fmt.Errorf("コマンドラインに指定できる引数は2つまでです。")
		panic(err)
	}

	ipAddr := 8082
	userRole := GUEST
	var err error

	if argsN >= 2 {
		userRole = os.Args[1]
		if userRole != HOST && userRole != GUEST {
			msg := fmt.Sprintf("コマンドラインの第1引数は、「%s」もしくは「%s」以外は指定できません。", HOST, GUEST)
			err := fmt.Errorf(msg)
			if err != nil {
				panic(err)
			}
		}

		if userRole == HOST {
			ipAddr = 8081
		}
	}

	if argsN == 3 {
		ipAddr, err = strconv.Atoi(os.Args[2])
		if err != nil {
			panic(err)
		}
	}

	server := http.Server{
        Addr:fmt.Sprintf(":%d", ipAddr),
        Handler:nil,
    }

	pattern := "/" + userRole + "_player/"
	http.HandleFunc(pattern, handleQuery)

	fmt.Println("Battle MCTS Server が起動しました。")
	fmt.Println("userRole =", userRole)
	fmt.Println("ipAddr =", ipAddr)
	fmt.Println("pattern =", pattern)

	err = server.ListenAndServe()
    if err != nil {
        panic(err)
    }
}