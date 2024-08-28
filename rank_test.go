package bippa_test

import (
	"testing"
	"fmt"
	bp "github.com/sw965/bippa"
)

func TestRankStatAdjustFluctuation(t *testing.T) {
	rs := bp.RankStat{
		Atk:5,
		Def:-5,
		SpAtk:6,
		SpDef:-6,
		Speed:0,
	}

	fluctuation := bp.RankStat{
		Atk:2,
		Def:-2,
		SpAtk:1,
		SpDef:-1,
		Speed:10, 
	}

	result := fluctuation.AdjustFluctuation(&rs)
	expected := bp.RankStat{
		Atk:1,
		Def:-1,
		SpAtk:0,
		SpDef:0,
		Speed:6,
	}
	fmt.Println(result)
	fmt.Println(expected)
}