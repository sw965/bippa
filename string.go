package bippa

import (
	"fmt"
	"github.com/sw965/omw/fn"
)

func StringToPokeName(s string) (PokeName, error) {
	if pokeName, ok := STRING_TO_POKE_NAME[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_POKE_NAME に含まれていない為、PokeNameに変換出来ません。", s)
		return pokeName, fmt.Errorf(msg)
	} else {
		return pokeName, nil
	}
}

func StringsToPokeNames(ss []string) (PokeNames, error) {
	return fn.MapWithError[PokeNames](ss, StringToPokeName)
} 

func StringToMoveName(s string) (MoveName, error) {
	if moveName, ok := STRING_TO_MOVE_NAME[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_MOVE_NAME に含まれていない為、MoveNameに変換出来ません。", s)
		return moveName, fmt.Errorf(msg)
	} else {
		return moveName, nil
	}
}

func StringsToMoveNames(ss []string) (MoveNames, error) {
	return fn.MapWithError[MoveNames](ss, StringToMoveName)
}

func StringToMoveCategory(s string) (MoveCategory, error) {
	if category, ok := STRING_TO_MOVE_CATEGORY[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_MOVE_CATEGORY に含まれていない為、MoveCategoryに変換出来ません。", s)
		return category, fmt.Errorf(msg)
	} else {
		return category, nil
	}
}

func StringToType(s string) (Type, error) {
	if t, ok := STRING_TO_TYPE[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_TYPE に含まれていない為、Typeに変換出来ません。", s)
		return t, fmt.Errorf(msg)
	} else {
		return t, nil
	}
}

func StringsToTypes(ss []string) (Types, error) {
	return fn.MapWithError[Types](ss, StringToType)
}

func StringToTargetRange(s string) (TargetRange, error) {
	if target, ok := STRING_TO_TARGET_RANGE[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_TARGET_RANGE に含まれていない為、TargetRangeに変換出来ません。", s)
		return target, fmt.Errorf(msg)
	} else {
		return target, nil
	}
}

func StringToNature(s string) (Nature, error) {
	if nature, ok := STRING_TO_NATURE[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_NATURE に含まれていない為、Natureに変換出来ません。", s)
		return nature, fmt.Errorf(msg)
	} else {
		return nature, nil
	}
}

func StringsToNatures(ss []string) (Natures, error) {
	return fn.MapWithError[Natures](ss, StringToNature)
}

func StringToAbility(s string) (Ability, error) {
	if ability, ok := STRING_TO_ABILITY[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_ABILITY に含まれていない為、Natureに変換出来ません。", s)
		return ability, fmt.Errorf(msg)
	} else {
		return ability, nil
	}
}

func StringsToAbilities(ss []string) (Abilities, error) {
	return fn.MapWithError[Abilities](ss, StringToAbility)
}

func StringToItem(s string) (Item, error) {
	if item, ok := STRING_TO_ITEM[s]; !ok {
		msg := fmt.Sprintf("%s は STRING_TO_ITEM に含まれていない為、Itemに変換出来ません。", s)
		return item, fmt.Errorf(msg)
	} else {
		return item, nil
	}
}

func StringsToItems(ss []string) (Items, error) {
	return fn.MapWithError[Items](ss, StringToItem)
}