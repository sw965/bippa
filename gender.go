package bippa

type Gender string

const (
	MALE    = Gender("♂")
	FEMALE  = Gender("♀")
	UNKNOWN = Gender("不明")
	EMPTY_GENDER = Gender("なし")
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

func (genders Genders) Index(gender Gender) int {
	for i, iGender := range genders {
		if iGender == gender {
			return i
		}
	}
	return -1
}