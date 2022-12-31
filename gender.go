package bippa

import (
	"math/rand"
)

type Gender string

const (
	MALE         = Gender("♂")
	FEMALE       = Gender("♀")
	UNKNOWN      = Gender("不明")
)

func (gender Gender) IsValid(pokeName PokeName) bool {
	switch POKEDEX[pokeName].Gender {
	case "♂♀両方":
		return gender == MALE || gender == FEMALE
	case "♂のみ":
		return gender == MALE
	case "♀のみ":
		return gender == FEMALE
	default:
		return gender == UNKNOWN
	}
}

type Genders []Gender

var ALL_GENDERS = Genders{MALE, FEMALE, UNKNOWN}
var ALL_GENDERS_LENGTH = len(ALL_GENDERS)

func NewVaildGenders(pokeName PokeName) Genders {
	result := make(Genders, 0, len(ALL_GENDERS))
	for _, gender := range ALL_GENDERS {
		if gender.IsValid(pokeName) {
			result = append(result, gender)
		}
	}
	return result
}

func (genders Genders) Index(gender Gender) int {
	for i, iGender := range genders {
		if iGender == gender {
			return i
		}
	}
	return -1
}

func (genders Genders) RandomChoice(random *rand.Rand) Gender {
	index := random.Intn(len(genders))
	return genders[index]
}
