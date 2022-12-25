package bippa

import (
	"fmt"
)

type Pokemon struct {
	Name    PokeName
	Types   Types
	Level   Level
	Gender  Gender
	Nature  Nature
	Ability Ability
	Item    Item
	Moveset Moveset

	Individual Individual
	Effort     Effort

	MaxHP     int
	CurrentHP int
	Atk       int
	Def       int
	SpAtk     int
	SpDef     int
	Speed     int

	StatusAilment        StatusAilment
	BadPoisonElapsedTurn int
	Rank                 Rank
	ChoiceMoveName       MoveName

	IsLeechSeed bool
}

func NewPokemon(pokeName PokeName, nature Nature, ability Ability, gender Gender, item Item,
	moveNames MoveNames, pointUps PointUps, individual *Individual, effort *Effort) (Pokemon, error) {

	if pokeName == "" {
		return Pokemon{}, fmt.Errorf("ポケモン名が、ゼロ値になっている")
	}

	if nature == "" {
		return Pokemon{}, fmt.Errorf("性格が、ゼロ値になっている")
	}

	if ability == "" {
		return Pokemon{}, fmt.Errorf("特性が、ゼロ値になっている")
	}

	if gender == "" {
		return Pokemon{}, fmt.Errorf("性別が、ゼロ値になっている")
	}

	if item == "" {
		return Pokemon{}, fmt.Errorf("アイテムが、ゼロ値になっている。何も持たせない場合は、EMPTY_ITEMを使って。")
	}

	pokeData, ok := POKEDEX[pokeName]

	if !ok {
		return Pokemon{}, fmt.Errorf("ポケモン名 %v は 不適", pokeName)
	}

	natureData, ok := NATUREDEX[nature]

	if !ok {
		return Pokemon{}, fmt.Errorf("性格 %v は 不適", nature)
	}

	if !ability.IsValid(pokeName) {
		return Pokemon{}, fmt.Errorf("特性 %v の %v は不適", ability, pokeName)
	}

	if !gender.IsValid(pokeName) {
		return Pokemon{}, fmt.Errorf("性別 %v の %v は不適", gender, pokeName)
	}

	if !item.IsValid() {
		return Pokemon{}, fmt.Errorf("アイテム %v は 不適", item)
	}

	moveset, err := NewMoveset(pokeName, moveNames, pointUps)

	if err != nil {
		return Pokemon{}, err
	}

	hp := CalcHp(pokeData.BaseHP, individual.HP, effort.HP)
	atk := CalcState(pokeData.BaseAtk, individual.Atk, effort.Atk, natureData.AtkBonus)
	def := CalcState(pokeData.BaseDef, individual.Def, effort.Def, natureData.DefBonus)
	spAtk := CalcState(pokeData.BaseSpAtk, individual.SpAtk, effort.SpAtk, natureData.SpAtkBonus)
	spDef := CalcState(pokeData.BaseSpDef, individual.SpDef, effort.SpDef, natureData.SpDefBonus)
	speed := CalcState(pokeData.BaseSpeed, individual.Speed, effort.Speed, natureData.SpeedBonus)

	return Pokemon{Name: pokeName, Nature: nature, Ability: ability, Gender: gender, Item: item, Moveset: moveset,
		Individual: *individual, Effort: *effort,
		MaxHP: hp, CurrentHP: hp, Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed,
		Types: pokeData.Types, Level: DEFAULT_LEVEL}, nil
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

	if pokemon1.MaxHP != pokemon2.MaxHP {
		return false
	}

	if pokemon1.CurrentHP != pokemon2.CurrentHP {
		return false
	}

	if pokemon1.Atk != pokemon2.Atk {
		return false
	}

	if pokemon1.Def != pokemon2.Def {
		return false
	}

	if pokemon1.SpAtk != pokemon2.SpAtk {
		return false
	}

	if pokemon1.SpDef != pokemon2.SpDef {
		return false
	}

	if pokemon1.Speed != pokemon2.Speed {
		return false
	}

	if pokemon1.Individual != pokemon2.Individual {
		return false
	}

	if pokemon1.Effort != pokemon2.Effort {
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

	if pokemon1.StatusAilment != pokemon2.StatusAilment {
		return false
	}

	if pokemon1.BadPoisonElapsedTurn != pokemon2.BadPoisonElapsedTurn {
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
	return pokemon.MaxHP == pokemon.CurrentHP
}

func (pokemon *Pokemon) IsFaint() bool {
	return pokemon.CurrentHP <= 0
}

func (pokemon *Pokemon) IsFaintDamage(damage int) bool {
	return damage >= int(pokemon.CurrentHP)
}

func (pokemon *Pokemon) CurrentDamage() int {
	return int(pokemon.MaxHP - pokemon.CurrentHP)
}

func (pokemon *Pokemon) SameTypeAttackBonus(moveName MoveName) SameTypeAttackBonus {
	moveType := MOVEDEX[moveName].Type
	inType := pokemon.Types.In(moveType)
	return BOOL_TO_SAME_TYPE_ATTACK_BONUS[inType]
}

func (pokemon *Pokemon) EffectivenessBonus(moveName MoveName) EffectivenessBonus {
	result := 1.0
	moveType := MOVEDEX[moveName].Type
	for _, pokeType := range pokemon.Types {
		result *= TYPEDEX[moveType][pokeType]
	}
	return EffectivenessBonus(result)
}

func (pokemon *Pokemon) BadPoisonDamage() int {
	damage := int(float64(pokemon.MaxHP) * float64(pokemon.BadPoisonElapsedTurn) / 16.0)
	if damage < 1 {
		return 1
	} else {
		return damage
	}
}

type Pokemons []Pokemon

func (pokemons Pokemons) Filter(f func(Pokemon) bool) Pokemons {
	result := make(Pokemons, 0, len(pokemons))
	for _, pokemon := range pokemons {
		if f(pokemon) {
			result = append(result, pokemon)
		}
	}
	return result
}

func (pokemons Pokemons) Count(pokemon Pokemon) int {
	result := 0
	for _, iPokemon := range pokemons {
		iPokemonPointer := &iPokemon
		if iPokemonPointer.Equal(&pokemon) {
			result += 1
		}
	}
	return result
}

func (pokemons Pokemons) Counts() []int {
	result := make([]int, len(pokemons))
	for i, pokemon := range pokemons {
		result[i] = pokemons.Count(pokemon)
	}
	return result
}
