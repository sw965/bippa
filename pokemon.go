package bippa

import (
	"fmt"
)

type Pokemon struct {
	Name            PokeName
	Level Level
	Nature          Nature
	Ability         Ability
	Gender          Gender
	Item            Item
	Moveset         Moveset

	Individual Individual
	Effort     Effort
	State State

	Types Types

	Rank  Rank
	StatusAilmentDetail StatusAilmentDetail
	ChoiceMoveName                 MoveName
	IsLeechSeed bool
}

func NewPokemon(pokeName PokeName, nature Nature, ability Ability, gender Gender, item Item,
	moveNames MoveNames, pointUps []PointUp, individual *Individual, effort *Effort) (Pokemon, error) {

	if !pokeName.IsValid() {
		errMsg := fmt.Sprintf("「%v」というポケモンは存在しない", pokeName)
		return Pokemon{}, fmt.Errorf(errMsg)
	}

	if !nature.IsValid() {
		errMsg := fmt.Sprintf("「%v」という性格は存在しない", nature)
		return Pokemon{}, fmt.Errorf(errMsg)
	}

	if !ability.IsValid(pokeName) {
		errMsg := fmt.Sprintf("特性「%v」の %v は不適", ability, pokeName)
		return Pokemon{}, fmt.Errorf(errMsg)
	}

	if !gender.IsValid(pokeName) {
		errMsg := fmt.Sprintf("性別 「%v」 の %v は不適", gender, pokeName)
		return Pokemon{}, fmt.Errorf(errMsg)
	}

	if !item.IsValid() {
		errMsg := fmt.Sprintf("「%v」というアイテムは存在しない", item)
		return Pokemon{}, fmt.Errorf(errMsg)
	}

	moveset, err := NewMoveset(pokeName, moveNames, pointUps)

	if err != nil {
		return Pokemon{}, err
	}

	pokeData := POKEDEX[pokeName]
	state, err := NewState(pokeName, nature, individual, effort, pokeData)

	if err != nil {
		return Pokemon{}, err
	}

	return Pokemon{Name: pokeName, Nature: nature, Ability: ability, Gender: gender, Item: item, Moveset: moveset,
		Individual: *individual, Effort: *effort,
		State: state, Types: pokeData.Types, Level:MAX_LEVEL}, nil
}

func (pokemon1 *Pokemon) Equal(pokemon2 *Pokemon) bool {
	if pokemon1.Name != pokemon2.Name {
		return false
	}

	if pokemon1.Nature != pokemon2.Nature {
		return false
	}

	if pokemon1.Ability != pokemon2.Ability {
		return false
	}

	if pokemon1.Gender != pokemon2.Gender {
		return false
	}

	if pokemon1.Item != pokemon2.Item {
		return false
	}

	if !pokemon1.Moveset.Equal(pokemon2.Moveset) {
		return false
	}

	if pokemon1.Individual != pokemon2.Individual {
		return false
	}

	if pokemon1.Effort != pokemon2.Effort {
		return false
	}

	if pokemon1.State != pokemon2.State {
		return false
	}

	for _, pokeType := range pokemon1.Types {
		if !pokemon2.Types.In(pokeType) {
			return false
		}
	}

	if pokemon1.Rank != pokemon2.Rank {
		return false
	}

	if pokemon1.StatusAilmentDetail != pokemon2.StatusAilmentDetail {
		return false
	}

	if pokemon1.ChoiceMoveName != pokemon2.ChoiceMoveName {
		return false
	}

	if pokemon1.IsLeechSeed != pokemon2.IsLeechSeed {
		return false
	}

	return true
}

func (pokemon *Pokemon) IsFullHP() bool {
	return pokemon.State.MaxHP == pokemon.State.CurrentHP
}

func (pokemon *Pokemon) IsFaint() bool {
	return pokemon.State.CurrentHP <= 0
}

func (pokemon *Pokemon) IsFaintDamage(damage int) bool {
	return damage >= int(pokemon.State.CurrentHP)
}

func (pokemon *Pokemon) CurrentDamage() int {
	return int(pokemon.State.MaxHP - pokemon.State.CurrentHP)
}

func (pokemon *Pokemon) NewSameTypeAttackBonus(moveName MoveName) SameTypeAttackBonus {
	moveType := MOVEDEX[moveName].Type
	inType := pokemon.Types.In(moveType)
	return BOOL_TO_SAME_TYPE_ATTACK_BONUS[inType]
}

func (pokemon *Pokemon) NewEffectivenessBonus(moveName MoveName) EffectivenessBonus {
	result := 1.0
	moveType := MOVEDEX[moveName].Type
	for _, pokeType := range pokemon.Types {
		result *= TYPEDEX[moveType][pokeType]
	}
	return EffectivenessBonus(result)
}

func (pokemon *Pokemon) BadPoisonDamage() int {
	return int(float64(pokemon.State.MaxHP) * float64(pokemon.StatusAilmentDetail.BadPoisonElapsedTurn) / 16.0)
}

func (pokemon *Pokemon) IsFocusSashOk(damage int) bool {
	return pokemon.IsFullHP() && pokemon.IsFaintDamage(damage) && pokemon.Item == "きあいのタスキ"
}

func (pokemon *Pokemon) StealthRockDamage() int {
	damagePercent := 1.0 / 8.0
	for _, pokeType := range pokemon.Types {
		damagePercent *= TYPEDEX[ROCK][pokeType]
	}
	return int(float64(pokemon.State.MaxHP) * damagePercent)
}
