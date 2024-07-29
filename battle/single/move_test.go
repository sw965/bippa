package single_test

// import (
// 	"fmt"
// 	"testing"
// 	"github.com/sw965/bippa/battle/single"
// 	bp "github.com/sw965/bippa"
// 	"github.com/sw965/bippa/battle/dmgtools"
// 	omwrand "github.com/sw965/omw/math/rand"
// 	omwslices "github.com/sw965/omw/slices"
// )

// // 第4世代のダメージ計算は、このURLのツールと同じ結果になるかどうかをテストする。
// // https://pokemon-trainer.net/dp/dmjs4.html

// func Helper(battle single.Battle, action *single.SoloAction, testNum int, move func(*single.Battle, *single.SoloAction, *single.Context) error) ([]single.Battle, error) {
// 	r := omwrand.NewMt19937()
// 	context := single.Context{
// 		DamageRandBonuses:dmgtools.RandBonuses{1.0},
// 		Rand:r,
// 	}

// 	ret := make([]single.Battle, testNum)
// 	b := battle.Clone()
// 	for i := 0; i < testNum; i++ {
// 		err := move(&b, action, &context)
// 		if err != nil {
// 			panic(err)
// 		}
// 		ret[i] = b
// 		b = battle.Clone()
// 	}
// 	return ret, nil
// }

// //10まんボルト
// func TestThunderbolt(t *testing.T) {
// 	battle := single.Battle{
// 		SelfLeadPokemons:bp.Pokemons{
// 			bp.NewRomanStan2009Latios(),
// 			bp.NewRomanStan2009Metagross(),
// 		},

// 		OpponentLeadPokemons:bp.Pokemons{
// 			bp.NewKusanagi2009Toxicroak(),
// 			bp.NewKusanagi2009Empoleon(),
// 		},
// 	}

// 	action := single.SoloAction{
// 		MoveName:bp.THUNDERBOLT,
// 		SrcIndex:0,
// 		TargetIndex:0,
// 		IsOpponentLeadTarget:true,
// 		Target:bp.NORMAL_TARGET,
// 	}

// 	testNum := 1000
// 	results, err := Helper(battle, &action, testNum, single.Thunderbolt)
// 	if err != nil {
// 		panic(err)
// 	}

// 	noCritCount := omwslices.CountFunc(results, func(b single.Battle) bool {
// 		return b.OpponentLeadPokemons[0].CurrentHP == (158-91)
// 	})

// 	critCount := omwslices.CountFunc(results, func(b single.Battle) bool {
// 		return b.OpponentLeadPokemons[0].CurrentHP == 0
// 	})
	
// 	paralysisCount := omwslices.CountFunc(results, func(b single.Battle) bool {
// 		return b.OpponentLeadPokemons[0].StatusAilment == bp.PARALYSIS
// 	})

// 	fmt.Println("noCritProbability", float64(noCritCount) / float64(testNum))
// 	fmt.Println("critProbability", float64(critCount) / float64(testNum))
// 	fmt.Println("paralysisCountProbability", float64(paralysisCount) / float64(testNum))
// }

// //ストーンエッジ
// func TestStoneEdge(t *testing.T) {
// 	battle := single.Battle{
// 		SelfLeadPokemons:bp.Pokemons{
// 			bp.NewRomanStan2009Snorlax(),
// 			bp.NewRomanStan2009Gyarados(),
// 		},

// 		OpponentLeadPokemons:bp.Pokemons{
// 			bp.NewKusanagi2009Empoleon(),
// 			bp.NewRomanStan2009Gyarados(),
// 		},
// 	}
// 	battle.SelfLeadPokemons[1].Rank.Atk -= 1

// 	action := single.SoloAction{
// 		MoveName:bp.STONE_EDGE,
// 		SrcIndex:1,
// 		TargetIndex:1,
// 		IsOpponentLeadTarget:true,
// 		Target:bp.NORMAL_TARGET,
// 	}

// 	testNum := 10000
// 	results, err := Helper(battle, &action, testNum, single.StoneEdge)
// 	if err != nil {
// 		panic(err)
// 	}

// 	hitCount := omwslices.CountFunc(results, func(b single.Battle) bool {
// 		return !b.OpponentLeadPokemons[1].IsFullHP()
// 	})

// 	missCount := omwslices.CountFunc(results, func(b single.Battle) bool {
// 		return b.OpponentLeadPokemons[1].IsFullHP()
// 	})

// 	critCount := omwslices.CountFunc(results, func(b single.Battle) bool {
// 		return b.OpponentLeadPokemons[1].CurrentHP == 0
// 	})

// 	fmt.Println("hitCount", float64(hitCount) / float64(testNum))
// 	fmt.Println("missCount", float64(missCount) / float64(testNum))
// 	fmt.Println("critCount", float64(critCount) / float64(testNum))
// }

// func TestSurf(t *testing.T) {
// 	battle := single.Battle{
// 		SelfLeadPokemons:bp.Pokemons{
// 			bp.NewKusanagi2009Toxicroak(),
// 			bp.NewKusanagi2009Empoleon(),
// 		},

// 		OpponentLeadPokemons:bp.Pokemons{
// 			bp.NewRomanStan2009Metagross(),
// 			bp.NewRomanStan2009Latios(),
// 		},
// 	}

// 	action := single.SoloAction{
// 		MoveName:bp.SURF,
// 		SrcIndex:1,
// 		TargetIndex:-1,
// 		Target:bp.OTHERS_TARGET,
// 	}

// 	testNum := 10
// 	results, err := Helper(battle, &action, testNum, single.Surf)
// 	if err != nil {
// 		panic(err)
// 	}

// 	currentHPCount := omwslices.CountFunc(results, func(b single.Battle) bool {
// 		return b.SelfLeadPokemons[0].IsFullHP() &&
// 			//現状ダメージが2ずれる。原因不明
// 			b.OpponentLeadPokemons[0].CurrentHP == (185-66) &&
// 			//現状ダメージが1ずれる。原因不明
// 			b.OpponentLeadPokemons[1].CurrentHP == (155-33)
// 	})

// 	if currentHPCount != testNum {
// 		t.Errorf("テスト失敗")
// 	}
// }