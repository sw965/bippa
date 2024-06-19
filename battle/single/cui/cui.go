package cui

import (
	"fmt"
	"github.com/sw965/bippa/battle/single"
	//battlemsg "github.com/sw965/bippa/battle/msg"
	omwcui "github.com/sw965/omw/cui"
	omwstrings "github.com/sw965/omw/strings"
	"time"
	"github.com/mattn/go-colorable"
	//"strings"
	bp "github.com/sw965/bippa"
)

type Top string

func NewTop(selfPokemon, opponentPokemon *bp.Pokemon) {
	ret += fmt.Sprintf("%s Lv.%d         %s Lv.%d", selfPokemon.Name.ToString(), selfPokemon.Level, opponentPokemon.Name.ToString(), opponentPokemon.Level)
	ret += fmt.Sprintf("%d/%d            %d/%d\n", selfPokemon.CurrentHP, selfPokemon.MaxHP, opponentPokemon.CurrentHP, opponentPokemon.MaxHP)
	ret += "+-----------------------------------------------------------+\n"
	ret += "|                                                           |\n"
	ret += "|                                                           |\n"
	ret += "+-----------------------------------------------------------+\n"
	ret += "1.たたかう 2.ポケモン 3.にげる\n"
	return ret
}

func Cui(battle *single.Battle, opponentTrainerName string, second float64, messageboxHeight, messageboxWidth int) {
	omwcui.ClearConsole()
	top := NewTop(&battle.SelfFighters[0], &battle.OpponentFighters[0])
	fmt.Println(top)
}