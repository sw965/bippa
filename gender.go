package bippa

type Gender int

const (
	UNKNOWN Gender = iota
	MALE
	FEMALE
)

func (g Gender) ToString() string {
	return GENDER_TO_STRING[g]
}

type Genders []Gender

func (gs Genders) ToStrings() []string {
	ss := make([]string, len(gs))
	for i, g := range gs {
		ss[i] = g.ToString()
	}
	return ss
}
