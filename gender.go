package bippa

type Gender string

const (
	MALE    = Gender("♂")
	FEMALE  = Gender("♀")
	UNKNOWN = Gender("不明")
)

func (gender Gender) IsValid(pokeName PokeName) bool {
	genderData := POKEDEX[pokeName].Gender

	if genderData == "♂♀両方" {
		return gender == MALE || gender == FEMALE
	}

	if genderData == "♂のみ" {
		return gender == MALE
	}

	if genderData == "♀のみ" {
		return gender == FEMALE
	}

	return gender == UNKNOWN
}
