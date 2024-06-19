package cui

import (
	"fmt"
	"github.com/sw965/bippa/battle/single"
	battlemsg "github.com/sw965/bippa/battle/msg"
	omwcui "github.com/sw965/omw/cui"
	omwstrings "github.com/sw965/omw/strings"
	"time"
	"github.com/mattn/go-colorable"
	//"strings"
)

type Top string

func NewTop(selfPokemon, opponentPokemon *Pokemon, msg string) {
	fmt.Println(fmt.Sprintf("%s Lv.%d %s Lv.%d", selfPokemon.Name.ToString(), selfPokemon.Level, opponentPokemon.Name.ToString(), opponentPokemon.Level))
	fmt.Println(fmt.Sprintf("%d/%d %d/%d", selfPokemon.CurrentHP, selfPokemon.MaxHP, opponentPokemon.CurrentHP, opponentPokemon.MaxHP))
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("|                                                           |")
	fmt.Println("|                                                           |")
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("1.たたかう 2.ポケモン 3.にげる ")

}

func Cui(battle *single.Battle, opponentTrainerName string, second float64, messageboxHeight, messageboxWidth int) {
	omwcui.ClearConsole()
	//fmt.Println(len(strings.Fields(string(battlemsg.NewChallengeByTrainer(opponentTrainerName)))))
	selfPokemon := battle.SelfFighters[0]
	opponentPokemon := battle.OpponentFighters[0]
	fmt.Println(fmt.Sprintf("%s Lv.%d %s Lv.%d", selfPokemon.Name.ToString(), selfPokemon.Level, opponentPokemon.Name.ToString(), opponentPokemon.Level))
	fmt.Println(fmt.Sprintf("%d/%d %d/%d", selfPokemon.CurrentHP, selfPokemon.MaxHP, opponentPokemon.CurrentHP, opponentPokemon.MaxHP))
	fmt.Println("+-----------------------------------------------------------+")
	fmt.Println("|                                                           |")
	fmt.Println("|                                                           |")
	fmt.Println("+-----------------------------------------------------------+")
	out := colorable.NewColorableStdout()
	fmt.Fprint(out, "\033[4;3H")

	for _, msg := range battlemsg.NewChallengeByTrainer(opponentTrainerName, "\n\033[2C").ToSlice() {
		//omwcui.ClearConsole()
		fmt.Printf(msg)
		time.Sleep(time.Duration(100.0*second) * time.Millisecond)
	}

	fmt.Fprint(out, "\033[4;3H")
	fmt.Println("                          ")
	fmt.Fprint(out, "\033[5;3H")
	fmt.Println("                          ")
	fmt.Fprint(out, "\033[4;3H")

	for _, msg := range battlemsg.NewGo("", selfPokemon.Name, true, "\n\033[2C").ToSlice() {
		fmt.Printf(msg)
		time.Sleep(time.Duration(100.0*second) * time.Millisecond)
	}

	fmt.Fprint(out, "\033[4;3H")
	fmt.Println("                          ")
	fmt.Fprint(out, "\033[5;3H")
	fmt.Println("                          ")
	fmt.Fprint(out, "\033[4;3H")

	for _, msg := range battlemsg.NewGo(opponentTrainerName, opponentPokemon.Name, false, "\n\033[2C").ToSlice() {
		fmt.Printf(msg)
		time.Sleep(time.Duration(100.0*second) * time.Millisecond)
	}

	fmt.Fprint(out, "\033[?25l")
	for i := 1; i < 150; i++ {
		fmt.Fprint(out, "\033[2;1H")
		currentHP := fmt.Sprintf("%d", i)
		pad := ""
		for i := 0; i < 3 - omwstrings.Len(currentHP); i++ {
			pad += " "
		}
		fmt.Printf("%s", pad+currentHP)
		fmt.Fprint(out, "\033[2;1H")
		time.Sleep(time.Duration(100.0*second) * time.Millisecond)
	}
	fmt.Fprint(out, "\033[?25h")
	time.Sleep(time.Duration(100.0*second) * time.Millisecond)
}