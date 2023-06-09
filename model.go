package bippa

import (
	"golang.org/x/exp/slices"
	"github.com/sw965/omw"
	"math/rand"
	"os"
	"fmt"
)

type TeamPokemonModelPart struct {
	Ability Ability
	Nature Nature
	Gender Gender
	Item Item
	MoveNames MoveNames

	HP State
	Atk State
	Def State
	SpAtk State
	SpDef State
	Speed State

	Value float64
	Number int
}

func (part *TeamPokemonModelPart) OK(pokemon *Pokemon) bool {
	if part.Ability != "" {
		if part.Ability != pokemon.Ability {
			return false
		}
	}

	if part.Nature != "" {
		if part.Nature != pokemon.Nature {
			return false
		}
	}

	if part.Gender != "" {
		if part.Gender != pokemon.Gender {
			return false
		}
	}

	if part.Item != "" {
		if part.Item != pokemon.Item {
			return false
		}
	}

	if len(part.MoveNames) != 0 {
		mclone := slices.Clone(part.MoveNames)
		mclone.Sort()
		pclone := slices.Clone(omw.Keys[MoveNames](pokemon.Moveset))
		pclone.Sort()

		if !slices.Equal(mclone, pclone) {
			return false
		}
	}

	if part.HP != EMPTY_STATE {
		if part.HP != pokemon.MaxHP {
			return false
		}
	}

	if part.Atk != EMPTY_STATE {
		if part.Atk != pokemon.Atk {
			return false
		}
	}

	if part.Def != EMPTY_STATE {
		if part.Def != pokemon.Def {
			return false
		}
	}

	if part.SpAtk != EMPTY_STATE {
		if part.SpAtk != pokemon.SpAtk {
			return false
		}
	}

	if part.SpDef != EMPTY_STATE {
		if part.SpDef != pokemon.SpDef {
			return false
		}
	}

	if part.Speed != EMPTY_STATE {
		if part.Speed != pokemon.Speed {
			return false
		}
	}

	return true
}

type TeamPokemonModel []*TeamPokemonModelPart

func NewMoveNamesAndAbilityTeamPokemonModel(moveNames MoveNames, r int, abilities Abilities) TeamPokemonModel {
	comb := omw.Combination[[]MoveNames, MoveNames](moveNames, r)
	result := make(TeamPokemonModel, 0, len(comb) * len(abilities))
	for _, mns := range comb {
		for _, ability := range abilities {
			part := TeamPokemonModelPart{Ability:ability, MoveNames:mns}
			result = append(result, &part)
		}
	}
	return result
}

func (model TeamPokemonModel) Init(r *rand.Rand) {
	for i, part := range model {
		part.Value = omw.RandFloat64(0, 16.0, r)
		part.Number = i
	}
}

func (model TeamPokemonModel) GetOK(pokemon *Pokemon) TeamPokemonModel {
	result := make(TeamPokemonModel, 0, len(model))
	for _, part := range model {
		if part.OK(pokemon) {
			result = append(result, part)
		}
	}
	return result
}

func (model TeamPokemonModel) Write(pokeName PokeName, modelName string, isOverwriteOk bool) error {
	folderPath := TEAM_POKEMON_MODEL_PATH + string(pokeName) + "/"
    _, err := os.Stat(folderPath)
	exist := !os.IsNotExist(err)

	if exist && !isOverwriteOk {
		return fmt.Errorf("上書き が 許可 されていないのに、上書きしようとした")
	}

	if !exist {
		err := os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	filePath := TEAM_POKEMON_MODEL_PATH + string(pokeName) + modelName + ".json"
	return omw.WriteJson(&model, filePath)
}