package main

import (
	"github.com/sw965/bippa/battle/single/cui"
	"github.com/sw965/bippa/battle/single"
	bp "github.com/sw965/bippa"
)

//testingを使ったテストだと、標準入力を受け付けないので、mainパッケージでテストをする

func main() {
	battle := single.Battle{
		SelfFighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
		OpponentFighters:single.Fighters{bp.NewTemplateBulbasaur(), bp.NewTemplateCharmander(), bp.NewTemplateSquirtle()},
	}

	cui.Cui(&battle, "してんのうの カトレア", 0.5)
}