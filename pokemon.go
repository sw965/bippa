package bippa

import (
	"fmt"
)

type PokemonState struct {
	MaxHP     int
	CurrentHP int
	Atk       int
	Def       int
	SpAtk     int
	SpDef     int
	Speed     int
}

func NewPokemonState(pokeName PokeName, nature Nature, individualState *IndividualState, effortState *EffortState) (PokemonState, error) {
	if !individualState.IsAllValid() {
		errMsg := fmt.Sprintf("個体値は%v～%vでなければならない", MIN_INDIVIDUAL, MAX_INDIVIDUAL)
		return PokemonState{}, fmt.Errorf(errMsg)
	}

	if !effortState.IsAllValid() {
		errMsg := fmt.Sprintf("努力値は%v～%vでなければならない", MIN_EFFORT, MAX_EFFORT)
		return PokemonState{}, fmt.Errorf(errMsg)
	}

	if !effortState.IsValidSum() {
		errMsg := fmt.Sprintf("努力値の合計値は%vを超えてはならない", MAX_SUM_EFFORT)
		return PokemonState{}, fmt.Errorf(errMsg)
	}

	pokeData := POKEDEX[pokeName]
	natureData := NATUREDEX[nature]
	hp := HpStateCalc(pokeData.BaseHP, individualState.HP, effortState.HP)
	atk := StateCalc(pokeData.BaseAtk, individualState.Atk, effortState.Atk, natureData.AtkBonus)
	def := StateCalc(pokeData.BaseDef, individualState.Def, effortState.Def, natureData.DefBonus)
	spAtk := StateCalc(pokeData.BaseSpAtk, individualState.SpAtk, effortState.SpAtk, natureData.SpAtkBonus)
	spDef := StateCalc(pokeData.BaseSpDef, individualState.SpDef, effortState.SpDef, natureData.SpDefBonus)
	speed := StateCalc(pokeData.BaseSpeed, individualState.Speed, effortState.Speed, natureData.SpeedBonus)
	return PokemonState{MaxHP: hp, CurrentHP: hp, Atk: atk, Def: def, SpAtk: spAtk, SpDef: spDef, Speed: speed}, nil
}

type Pokemon struct {
	Name            PokeName
	Level int
	Nature          Nature
	Ability         Ability
	Gender          Gender
	Item            Item
	Moveset         Moveset
	IndividualState IndividualState
	EffortState     EffortState

	State PokemonState
	Types Types
	RankState  RankState

	StatusAilmentParam StatusAilmentParam

	ChoiceMoveName                 MoveName
}

func NewPokemon(pokeName PokeName, nature Nature, ability Ability, gender Gender, item Item,
	moveNames MoveNames, pointUps []PointUp, individualState *IndividualState, effortState *EffortState) (Pokemon, error) {

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

	state, err := NewPokemonState(pokeName, nature, individualState, effortState)

	if err != nil {
		return Pokemon{}, err
	}

	pokeData := POKEDEX[pokeName]
	return Pokemon{Name: pokeName, Nature: nature, Ability: ability, Gender: gender, Item: item, Moveset: moveset,
		IndividualState: *individualState, EffortState: *effortState,
		State: state, Types: pokeData.Types, Level:LEVEL}, nil
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

	if pokemon1.IndividualState != pokemon2.IndividualState {
		return false
	}

	if pokemon1.EffortState != pokemon2.EffortState {
		return false
	}

	if pokemon1.State != pokemon2.State {
		return false
	}

	for _, pokeType := range pokemon1.Types {
		if !pokemon2.InType(pokeType) {
			return false
		}
	}

	if pokemon1.Rank != pokemon2.Rank {
		return false
	}

	if pokemon1.StatusAilmentParam != pokemon2.StatusAilmentParam {
		return false
	}

	if pokemon1.ChoiceMoveName != pokemon2.ChoiceMoveName {
		return false
	}

	if pokemon1.IsProtectStats != pokemon2.IsProtectStats {
		return false
	}

	if pokemon1.ProtectConsecutiveSuccessCount != pokemon2.ProtectConsecutiveSuccessCount {
		return false
	}

	if pokemon1.IsLeechSeedState != pokemon2.IsLeechSeedState {
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
	return damage >= pokemon.State.CurrentHP
}

func (pokemon *Pokemon) CurrentDamage() int {
	return pokemon.State.MaxHP - pokemon.State.CurrentHP
}

func (pokemon *Pokemon) SameTypeAttackBonus(moveName MoveName) float64 {
	moveType := MOVEDEX[moveName].Type
	if pokemon.InType(moveType) {
		return 6144.0 / 4096.0
	}
	return 4096.0 / 4096.0
}

func (pokemon *Pokemon) EffectivenessBonus(moveName MoveName) float64 {
	result := 1.0
	moveType := MOVEDEX[moveName].Type
	for _, pokeType := range pokemon.Types {
		result *= TYPEDEX[moveType][pokeType]
	}
	return result
}

func (pokemon *Pokemon) BadPoisonDamage() int {
	return int(float64(pokemon.State.MaxHP) * float64(pokemon.StatusAilmentParam.BadPoisonElapsedTurn) / 16.0)
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
