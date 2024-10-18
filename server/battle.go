package server

import (
    "fmt"
	"net/http"
	"net/url"
    "encoding/json"
    bp "github.com/sw965/bippa"
	"github.com/sw965/bippa/battle"
	"github.com/sw965/bippa/battle/game"
	"github.com/sw965/crow/game/simultaneous"
	omwHTTP "github.com/sw965/omw/http"
)

func queryToPokemons(query string) (bp.Pokemons, error) {
	easyReadPokemons, err := omwHTTP.QueryToType[bp.EasyReadPokemons](query)
	if err != nil {
		return bp.Pokemons{}, err
	}
	return easyReadPokemons.From()
}

func queryToAction(query string) (battle.Action, error) {
	easyReadAction, err := omwHTTP.QueryToType[battle.EasyReadAction](query)
	if err != nil {
		return battle.Action{}, err
	}
	return easyReadAction.From()
}

var battleManager battle.Manager

func GetBattleManager() battle.Manager {
	return battleManager
}

func HandleBattleInit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	guestTrainerTitle, err := url.QueryUnescape(r.URL.Query().Get("guest_trainer_title"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	if guestTrainerTitle == "" {
		err := fmt.Errorf(fmt.Sprintf("guest_trainer_titleが空である為、Battleを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	guestTrainerName, err := url.QueryUnescape(r.URL.Query().Get("guest_trainer_name"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	if guestTrainerName == "" {
		err := fmt.Errorf(fmt.Sprintf("guest_trainer_nameが空である為、Battleを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hostLeadPokemonsQuery := r.URL.Query().Get("host_lead_pokemons")
	if hostLeadPokemonsQuery == "" {
		err := fmt.Errorf(fmt.Sprintf("host_lead_pokemonsが空である為、Battleを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hostLeadPokemons, err := queryToPokemons(hostLeadPokemonsQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return		
	}

	hostBenchPokemonsQuery:= r.URL.Query().Get("host_bench_pokemons")
	if hostBenchPokemonsQuery == "" {
		err := fmt.Errorf(fmt.Sprintf("host_bench_pokemonsが空である為、Battleを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	hostBenchPokemons, err := queryToPokemons(hostBenchPokemonsQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return		
	}

	guestLeadPokemonsQuery:= r.URL.Query().Get("guest_lead_pokemons")
	if guestLeadPokemonsQuery == "" {
		err := fmt.Errorf(fmt.Sprintf("guest_lead_pokemonsが空である為、Battleを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	guestLeadPokemons, err := queryToPokemons(guestLeadPokemonsQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return	
	}

	guestBenchPokemonsQuery:= r.URL.Query().Get("guest_bench_pokemons")
	if guestBenchPokemonsQuery == "" {
		err := fmt.Errorf(fmt.Sprintf("guest_bench_pokemonsが空である為、Battleを初期化出来ません。"))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	guestBenchPokemons, err := queryToPokemons(guestBenchPokemonsQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return	
	}

	battleManager = battle.Manager{
		CurrentSelfLeadPokemons:hostLeadPokemons,
		CurrentSelfBenchPokemons:hostBenchPokemons,
		CurrentOpponentLeadPokemons:guestLeadPokemons,
		CurrentOpponentBenchPokemons:guestBenchPokemons,
	}

	ms := make(battle.Managers, 0, 128)
	battle.GlobalContext.Observer = func(m *battle.Manager) {
		ms = append(ms, m.Clone())
	}
	battleManager.Init(guestTrainerTitle, guestTrainerName)

	response, err := json.Marshal(ms.ToEasyRead())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
	fmt.Println("Battleを初期化しました。")
}

// GoからJavascriptに送られる構造体のフィールドは、全て大文字で始まるように一貫を持たせる為に、main package でも大文字から始める。
type UsableMoveNamesInfo struct {
    CurrentSelfLeadPokemons [][]string
    CurrentSelfBenchPokemons [][]string
    CurrentOpponentLeadPokemons [][]string
    CurrentOpponentBenchPokemons [][]string
}

func NewUsableMoveNamesInfo(m *battle.Manager) UsableMoveNamesInfo {
	get := func(ps bp.Pokemons) [][]string {
		sss := make([][]string, len(ps))
		for i, p := range ps {
			sss[i] = p.UsableMoveNames().ToStrings()
		}
		return sss
	}

	return UsableMoveNamesInfo{
		CurrentSelfLeadPokemons:get(m.CurrentSelfLeadPokemons),
		CurrentSelfBenchPokemons:get(m.CurrentSelfBenchPokemons),
		CurrentOpponentLeadPokemons:get(m.CurrentOpponentLeadPokemons),
		CurrentOpponentBenchPokemons:get(m.CurrentOpponentBenchPokemons),
	}
}

// GoからJavascriptに送られる構造体のフィールドは、全て大文字で始まるように一貫を持たせる為に、main package でも大文字から始める。
type BattleEndInfo struct {
	IsEnd bool
	JointEval simultaneous.JointEval
}

func HandleBattleQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	var response []byte
	queryType, err := url.QueryUnescape(r.URL.Query().Get("query_type"))

	switch queryType {
		case "battle":
			response, err = json.Marshal(battleManager.ToEasyRead())
		case "separate_legal_actions":
			response, err = json.Marshal(game.SeparateLegalActions(&battleManager).ToEasyRead())
		case "battle_end_info":
			isEnd, jointEval := game.IsEnd(&battleManager)
			response, err = json.Marshal(BattleEndInfo{IsEnd:isEnd, JointEval:jointEval})
		case "usable_move_names_info":
			response, err = json.Marshal(NewUsableMoveNamesInfo(&battleManager))
	}
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
	w.Write(response)
}

func HandleBattleCommand(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Type", "application/json")

	commandType, err := url.QueryUnescape(r.URL.Query().Get("command_type"))
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	switch commandType {
		case "push":
			hostActionQuery := r.URL.Query().Get("host_action")
			if hostActionQuery == "" {
				err := fmt.Errorf("host_actionが空である為、pushを実行出来ません。")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			hostAction, err := queryToAction(hostActionQuery)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return			
			}
			for i := range hostAction {
				hostAction[i].IsCurrentSelf = true
			}

			guestActionQuery := r.URL.Query().Get("guest_action")
			if guestActionQuery == "" {
				err := fmt.Errorf("guest_actionが空である為、pushを実行出来ません。")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			guestAction, err := queryToAction(guestActionQuery)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return			
			}

			ms := make(battle.Managers, 0, 256)
			ms = append(ms, battleManager.Clone())

			battle.GlobalContext.Observer = func(m *battle.Manager) {
				mv := m.Clone()
				if !mv.CurrentSelfIsHost {
					mv.SwapView()
				}
				ms = append(ms, mv)
			}
			
			nextBattleManager, err := game.Push(battleManager, battle.Actions{hostAction, guestAction})
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