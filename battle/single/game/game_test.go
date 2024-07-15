package game_test

import (
	"testing"
	"fmt"
	"github.com/sw965/bippa/battle/single"
	// "github.com/sw965/bippa/battle/single/game"
	bp "github.com/sw965/bippa"
	"github.com/sw965/omw/fn"
)

func TestSingleBattleLegalSeparateActions(t *testing.T) {
	battle := single.Battle{
		SelfLeadPokemons:bp.Pokemons{bp.NewRomanStan2009Gyarados(), bp.NewRomanStan2009Metagross()},
		SelfBenchPokemons:bp.Pokemons{bp.NewRomanStan2009Latios()},
		OpponentLeadPokemons:bp.Pokemons{bp.NewKusanagi2009Empoleon(), bp.NewKusanagi2009Snorlax()},
		OpponentBenchPokemons:bp.Pokemons{bp.NewKusanagi2009Toxicroak()},
	}
	fmt.Println(LegalFunc(&battle).ToEasyRead())
	fmt.Println(LegalFunc2(&battle).ToEasyRead())
}

func TestDoubleBattleLegalSeparateActions(t *testing.T) {

}