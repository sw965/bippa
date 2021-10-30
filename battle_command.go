package bippa

import (
	"math/rand"
)

type BattleCommand string

func (battleCommand BattleCommand) IsPokeName() bool {
	_, ok := POKEDEX[PokeName(battleCommand)]
	return ok
}

func (battleCommand BattleCommand) IsMoveName() bool {
	_, ok := MOVEDEX[MoveName(battleCommand)]
	return ok
}

func (battleCommand BattleCommand) PriorityRank() int {
	if battleCommand == BattleCommand(STRUGGLE) {
		return 0
	}

	if battleCommand.IsMoveName() {
		return MOVEDEX[MoveName(battleCommand)].PriorityRank
	}
	return 999
}

type BattleCommands []BattleCommand

func (battleCommands BattleCommands) RandomChoice(random *rand.Rand) BattleCommand {
	index := random.Intn(len(battleCommands))
	return battleCommands[index]
}
