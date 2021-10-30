package bippa

import (
	"encoding/json"
	"io/ioutil"
)

const (
	FIGHTERS_LENGTH = 3
)

type Fighters [FIGHTERS_LENGTH]Pokemon

func (fighters1 *Fighters) Equal(fighters2 *Fighters) bool {
	for i, pokemon := range fighters1 {
		if !pokemon.Equal(&fighters2[i]) {
			return false
		}
	}
	return true
}

func (fighters *Fighters) Index(pokeName PokeName) int {
	for i, pokemon := range fighters {
		if pokemon.Name == pokeName {
			return i
		}
	}
	return -1
}

func (fighters *Fighters) PokeNames() PokeNames {
	result := make(PokeNames, FIGHTERS_LENGTH)
	for i, pokemon := range fighters {
		result[i] = pokemon.Name
	}
	return result
}

func (fighters *Fighters) IsUnique() bool {
	return fighters.PokeNames().IsUnique()
}

func (fighters *Fighters) IsAllFaint() bool {
	for _, pokemon := range fighters {
		if !pokemon.IsFaint() {
			return false
		}
	}
	return true
}

func (fighters *Fighters) AvailableMoveNames() MoveNames {
	if fighters[0].IsFaint() {
		return MoveNames{}
	}

	moveset := fighters[0].Moveset
	var u MoveNames

	if fighters[0].ChoiceMoveName != "" {
		u = MoveNames{fighters[0].ChoiceMoveName}
	} else {
		u = moveset.MoveNames()
	}

	result := make(MoveNames, 0)
	for _, moveName := range u {
		powerPoint := moveset[moveName]
		if powerPoint > 0 {
			result = append(result, moveName)
		}
	}

	if len(result) == 0 {
		return MoveNames{STRUGGLE}
	}
	return result
}

func (fighters *Fighters) SwitchablePokeNames() []PokeName {
	result := make([]PokeName, 0)
	for _, pokemon := range fighters[1:] {
		if !pokemon.IsFaint() {
			result = append(result, pokemon.Name)
		}
	}
	return result
}

func (fighters *Fighters) AvailableBattleCommands() BattleCommands {
	availableMoveNames := fighters.AvailableMoveNames()
	switchablePokeNames := fighters.SwitchablePokeNames()
	result := make([]BattleCommand, 0, len(availableMoveNames)+len(switchablePokeNames))

	for _, moveName := range availableMoveNames {
		result = append(result, BattleCommand(moveName))
	}

	for _, pokeName := range switchablePokeNames {
		result = append(result, BattleCommand(pokeName))
	}
	return result
}

func (fighters *Fighters) IsAvailableBattleCommand(battleCommand BattleCommand) bool {
	availableBattleCommands := fighters.AvailableBattleCommands()
	for _, iBattleCommand := range availableBattleCommands {
		if iBattleCommand == battleCommand {
			return true
		}
	}
	return false
}

func (fighters *Fighters) Save(filePath string) error {
	file, err := json.MarshalIndent(fighters, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, file, 0644)
}

func ReadFighters(filePath string) (Fighters, error) {
	fighters := Fighters{}
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return Fighters{}, err
	}
	err = json.Unmarshal(file, &fighters)
	return fighters, err
}
