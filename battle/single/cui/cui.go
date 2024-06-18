package cui

import (
	"fmt"
	"github.com/sw965/bippa/battle/single"
	battlemsg "github.com/sw965/bippa/battle/msg"
	omwcui "github.com/sw965/omw/cui"
	"time"
)

func Cui(battle *single.Battle, opponentTrainerName string, second float64) {
	omwcui.ClearConsole()
	for _, msg := range battlemsg.NewChallengeByTrainer(opponentTrainerName).Accumulate() {
		omwcui.ClearConsole()
		fmt.Println(msg)
		time.Sleep(time.Duration(100.0*second) * time.Millisecond)
	}

	selfPokemon := battle.SelfFighters[0]
	//opponentPokemon := battle.OpponentFighters[0]

	for _, msg := range battlemsg.NewGo(opponentTrainerName, selfPokemon.Name, true).Accumulate() {
		omwcui.ClearConsole()
		fmt.Println(msg)
		time.Sleep(time.Duration(100.0*second) * time.Millisecond)
	}
}