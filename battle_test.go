package bippa

import (
	"fmt"
	//"github.com/seehuhn/mt19937"
	//"math"
	//"math/rand"
	"testing"
	//"time"
)

func NewGyaradosu() Pokemon {
	pokemon, err := NewPokemon("ギャラドス", "いじっぱり", "いかく", "♀", "たべのこし",
		MoveNames{"たきのぼり", "こおりのキバ", "りゅうのまい"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{Atk:252, Speed:252},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewHerakurosu() Pokemon {
	pokemon, err := NewPokemon("ヘラクロス", "いじっぱり", "こんじょう", "♂", "こだわりハチマキ",
		MoveNames{"インファイト", "かわらわり"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, Atk: 252, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewGenga() Pokemon {
	pokemon, err := NewPokemon("ゲンガー", "ひかえめ", "のろわれボディ", "♀", "くろいヘドロ",
		MoveNames{"シャドーボール", "きあいだま"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, SpAtk: 252, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewKureseria() Pokemon {
	pokemon, err := NewPokemon("クレセリア", "ずぶとい", "ふゆう", "♀", "ゴツゴツメット",
		MoveNames{"つきのひかり", "れいとうビーム"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, Def: 252, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewParusixen() Pokemon {
	pokemon, err := NewPokemon("パルシェン", "いじっぱり", "スキルリンク", "♂", "いのちのたま",
		MoveNames{"つららおとし", "からをやぶる", "ロックブラスト"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{Atk: 252, Speed: 252, HP: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewEamudo() Pokemon {
	pokemon, err := NewPokemon("エアームド", "わんぱく", "がんじょう", "♀", "ゴツゴツメット",
		MoveNames{"つばさでうつ", "はねやすめ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, Def: 252, Speed: 4},
	)

	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewBangirasu() Pokemon {
	pokemon, err := NewPokemon("バンギラス", "いじっぱり", "すなおこし", "♀", "なし",
		MoveNames{"かみくだく", "いわなだれ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, Atk: 252, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewHassamu() Pokemon {
	pokemon, err := NewPokemon("ハッサム", "いじっぱり", "テクニシャン", "♂", "こだわりハチマキ",
		MoveNames{"シザークロス", "バレットパンチ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, Atk: 252, Speed: 4},
	)

	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewGaburiasu() Pokemon {
	pokemon, err := NewPokemon("ガブリアス", "ようき", "さめはだ", "♀", "こだわりハチマキ",
		MoveNames{"じしん", "げきりん"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{Atk: 252, Speed: 252, HP: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewHapinasu() Pokemon {
	pokemon, err := NewPokemon("ハピナス", "ひかえめ", "てんのめぐみ", "♀", "たべのこし",
		MoveNames{"シャドーボール", "タマゴうみ", "10まんボルト"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{SpAtk: 252, SpDef: 252, Speed: 4},
	)
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewRateosu() Pokemon {
	pokemon, err := NewPokemon("ラティオス", "ひかえめ", "ふゆう", "♂", "こだわりメガネ",
		MoveNames{"りゅうのはどう", "サイコキネシス", "かみなり"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{SpAtk: 252, Speed: 252, Def: 4},
	)

	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewHusigibana() Pokemon {
	pokemon, err := NewPokemon("フシギバナ", "おだやか", "しんりょく", "♀", "なし",
		MoveNames{"エナジーボール", "やどりぎのタネ", "こうごうせい"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, SpDef: 252, Speed: 4})
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewRiza_don() Pokemon {
	pokemon, err := NewPokemon("リザードン", "ひかえめ", "もうか", "♂", "たべのこし",
		MoveNames{"かえんほうしゃ"}, PointUps{MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, SpAtk: 252, Speed: 4})
	if err != nil {
		panic(err)
	}
	return pokemon
}

func NewKamekkusu() Pokemon {
	pokemon, err := NewPokemon("カメックス", "ひかえめ", "げきりゅう", "♂", "なし",
		MoveNames{"なみのり", "れいとうビーム"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
		&ALL_MAX_INDIVIDUAL, &Effort{HP: 252, SpAtk: 252, Speed: 4})
	if err != nil {
		panic(err)
	}
	return pokemon
}

// func TestBattle(t *testing.T) {
//   p1Fighters := Fighters{NewHusigibana(), NewRiza_don(), NewKamekkusu()}
//   p2Fighters := Fighters{NewKamekkusu(), NewRiza_don(), NewHusigibana()}
//   battle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
// }

func TestPush(t *testing.T) {
	// mtRandom := rand.New(mt19937.New())
	// mtRandom.Seed(time.Now().UnixNano())
	//
	//p1Fighters := Fighters{NewHusigibana(), NewRiza_don(), NewHusigibana()}
	//p2Fighters := Fighters{NewKamekkusu(), NewRiza_don(), NewHusigibana()}
	// p2Fighters[0].CurrentHP = 0
	//battle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
	// var err error
	//
	// battle, err = battle.Push("リザードン", mtRandom)
	// if err != nil {
	//   panic(err)
	// }
	//
	// battle, err = battle.Push("リーフストーム", mtRandom)
	// if err != nil {
	//   panic(err)
	// }
	//
	// battle, err = battle.Push("ひのこ", mtRandom)
	// if err != nil {
	//   panic(err)
	// }
	//
	// fmt.Println(battle.P1Fighters[0].CurrentHP)
	// fmt.Println(battle.P1Fighters[1].CurrentHP)
	// fmt.Println(battle.P1Fighters[2].CurrentHP)
	//
	// fmt.Println(battle.P2Fighters[0].CurrentHP)
	// fmt.Println(battle.P2Fighters[1].CurrentHP)
	// fmt.Println(battle.P2Fighters[2].CurrentHP)
}

func Test(t *testing.T) {
	p1Fighters := Fighters{NewGyaradosu(), NewKureseria(), NewRiza_don()}
	p2Fighters := Fighters{NewRateosu(), NewGenga(), NewRiza_don()}

	p2Fighters[0].CurrentHP /= 2
	p1Fighters[2].CurrentHP = 0
	p2Fighters[2].CurrentHP = 0
	initBattle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}

	// mtRandom := rand.New(mt19937.New())
	// mtRandom.Seed(time.Now().UnixNano())
	//
	// randomInstructionTrainer := NewRandomInstructionTrainer(mtRandom)
	// randomPlayoutEval := NewPlayoutEval(randomInstructionTrainer, mtRandom)
	//
	// allNodes, err := RunMCTS(initBattle, 1, math.Sqrt(2), NoPolicy, &randomPlayoutEval, mtRandom)
	// if err != nil {
	// 	panic(err)
	// }
	//
	// for actionCmd, pucb := range allNodes[0].ActionCmdPUCBs {
	// 	fmt.Println(actionCmd, pucb.AverageReward(), pucb.Trial)
	// }
	fmt.Println(initBattle.ExpectedAttackDamage("だいもんじ"))
	//fmt.Println(initBattle.AttackDamageProbabilityDistribution("シャドーボール"))
	fmt.Println(NewNotBad(&initBattle))
	//mctsTrainer1960 := NewMCTSTrainer(196, math.Sqrt(2), &randomPlayoutEval, mtRandom)
	//mctsTrainer196 := NewMCTSTrainer(16, math.Sqrt(2), &randomPlayoutEval, mtRandom)

	// gameNum := 196
	// winCount := 0
	// for i := 0; i < gameNum; i++ {
	// 	winner, err := mctsTrainer1960.Playout(mctsTrainer196, initBattle, mtRandom)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	if winner == WINNER_P1 {
	// 		winCount += 1
	// 	}
	// }
	// fmt.Println(float64(winCount) / float64(gameNum))
}
