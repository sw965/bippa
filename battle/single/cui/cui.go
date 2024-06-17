package cui

import (
	"time"
)

type Face struct {
	P1PokeName bp.PokeName
	P1Level bp.Level
	P1Gender bp.Gender
	P1MaxHP int
	P1CurrentHP int

	P1PokeName bp.PokeName
	P1Level bp.Level
	P1Gender bp.Gender
	P2MaxHP int
	P2CurrentHP int

	Message string
}

func (f *Face) Println() {
	data := []string{
		fmt.Sprintf("Lv.%d %s %s\n%d/%d", f.P1Level, f.P1Gender.ToString(), P1PokeName.ToString(), f.P1CurrentHP, f.P1MaxHP),
		fmt.Sprintf("Lv.%d %s %s\n%d/%d", f.P2Level, f.P2Gender.ToString(), P2PokeName.ToString(), f.P2CurrentHP, f.P2MaxHP)),
	}

	for _, ele := range data {
		fmt.Println(ele)
	}
}

type Faces []Face

type Bench struct {
}