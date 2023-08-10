package bippa

import (
	"fmt"
	"github.com/sw965/omw/fn"
)

type Gender int

const (
	MALE Gender = iota
	FEMALE
	UNKNOWN
)

func NewGender(s string) (Gender, error) {
	switch s {
		case "♂":
			return MALE, nil
		case "♀":
			return FEMALE, nil
		case "不明":
			return UNKNOWN, nil
		default:
			return -1, fmt.Errorf("不適なgender")
	}
}

type Genders []Gender

var ALL_GENDERS = Genders{MALE, FEMALE, UNKNOWN}

func NewGenders(ss []string) (Genders, error) {
	return fn.MapError[Genders](ss, NewGender)
}