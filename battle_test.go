package bippa

import (
  "testing"
  "fmt"
  // "math"
  // "math/rand"
  // "time"
  // "github.com/seehuhn/mt19937"
)

func NewHerakurosu() Pokemon {
  pokemon, err := NewPokemon("ヘラクロス", "いじっぱり", "こんじょう", "♂", "こだわりハチマキ",
    MoveNames{"インファイト", "かわらわり"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{HP:252, Atk:252, Speed:4},
  )
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewGenga() Pokemon {
  pokemon, err := NewPokemon("ゲンガー", "ひかえめ", "のろわれボディ", "♀", "くろいヘドロ",
    MoveNames{"シャドーボール", "きあいだま"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{HP:252, SpAtk:252, Speed:4},
  )
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewKureseria() Pokemon {
  pokemon, err := NewPokemon("クレセリア", "ずぶとい", "ふゆう", "♀", "ゴツゴツメット",
    MoveNames{"つきのひかり", "れいとうビーム"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{HP:252, Def:252, Speed:4},
  )
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewParusixen() Pokemon {
  pokemon, err := NewPokemon("パルシェン", "いじっぱり", "スキルリンク", "♂", "いのちのたま",
    MoveNames{"つららばり", "からをやぶる", "ロックブラスト"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{Atk:252, Speed:252, HP:4},
  )
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewEamudo() Pokemon {
  pokemon, err := NewPokemon("エアームド", "わんぱく", "がんじょう", "♀", "ゴツゴツメット",
    MoveNames{"つばさでうつ", "はねやすめ"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{HP:252, Def:252, Speed:4},
  )

  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewGaburiasu() Pokemon {
  pokemon, err := NewPokemon("ガブリアス", "ようき", "さめはだ", "♀", "こだわりハチマキ",
    MoveNames{"じしん", "げきりん"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{Atk:252, Speed:252, HP:4},
  )
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewHapinasu() Pokemon {
  pokemon, err := NewPokemon("ハピナス", "ひかえめ", "てんのめぐみ", "♀", "たべのこし",
    MoveNames{"シャドーボール", "タマゴうみ", "10まんボルト"}, PointUps{MAX_POINT_UP, MAX_POINT_UP, MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{SpAtk:252, SpDef:252, Speed:4},
  )
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewRateosu() Pokemon {
  pokemon, err := NewPokemon("ラティオス", "ひかえめ", "ふゆう", "♂", "こだわりメガネ",
    MoveNames{"りゅうのはどう", "サイコキネシス"}, PointUps{MAX_POINT_UP, MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{SpAtk:252, Speed:252, Def:4},
  )

  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewHusigibana() Pokemon {
  pokemon, err := NewPokemon("フシギバナ", "おだやか", "しんりょく", "♀", "なし",
    MoveNames{"エナジーボール"}, PointUps{MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{HP:252, SpDef:252, Speed:4})
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewRiza_don() Pokemon {
  pokemon, err := NewPokemon("リザードン", "ひかえめ", "もうか", "♂", "なし",
                        MoveNames{"ひのこ"}, PointUps{MAX_POINT_UP},
                        &ALL_MAX_INDIVIDUAL, &Effort{HP:252, SpAtk:252, Speed:4})
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewKamekkusu() Pokemon {
  pokemon, err := NewPokemon("カメックス", "ひかえめ", "げきりゅう", "♂", "なし",
    MoveNames{"みずでっぽう"}, PointUps{MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{HP:252, SpAtk:252, Speed:4})
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
  // p1Fighters := Fighters{NewHusigibana(), NewRiza_don(), NewHusigibana()}
  // p2Fighters := Fighters{NewKamekkusu(), NewRiza_don(), NewHusigibana()}
  // p2Fighters[0].CurrentHP = 0
  // battle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
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
  moveNames := MoveNames{"なみのり", "しんそく", "バレットパンチ", "れいとうビーム", "ねこだまし"}
  fmt.Println(moveNames.Sort(func(moveName MoveName) int {return MOVEDEX[moveName].PriorityRank}))
  // p1Fighters := Fighters{NewHusigibana(), NewKamekkusu(), NewRiza_don()}
  // p2Fighters := Fighters{NewKamekkusu(), NewRiza_don(), NewHusigibana()}
  //
  // //p1Fighters[0].CurrentHP = 0
  // p1Fighters[1].CurrentHP = 0
  // p1Fighters[2].CurrentHP = 0
  //
  // //p2Fighters[0].CurrentHP = 200
  // p2Fighters[1].CurrentHP = 0
  // p2Fighters[2].CurrentHP = 0
  //
  // battle := Battle{P1Fighters:p2Fighters, P2Fighters:p1Fighters}
  //
  // // battle.P1Fighters[1].CurrentHP = 0
  // // battle.P1Fighters[2].CurrentHP = 0
  // // battle.P2Fighters[1].CurrentHP = 0
  // // battle.P2Fighters[2].CurrentHP = 0
  // //
  // mtRandom := rand.New(mt19937.New())
  // mtRandom.Seed(time.Now().UnixNano())
  // // trainer := NewRandomInstructionTrainer(mtRandom)
  // // gameNum := 128000
  // // p1Count := 0
  // // p2Count := 0
  // //
  // // for i := 0; i < gameNum; i++ {
  // //   gameEndBattle, err := trainer.OneGame(trainer, battle, mtRandom)
  // //   if err != nil {
  // //     panic(err)
  // //   }
  // //
  // //   winner, err := gameEndBattle.Winner()
  // //   if err != nil {
  // //     panic(err)
  // //   }
  // //
  // //   if winner == WINNER_P1 {
  // //     p1Count += 1
  // //   }
  // //
  // //   if winner == WINNER_P2 {
  // //     p2Count += 1
  // //   }
  // // }
  // // fmt.Println(float64(p1Count) / float64(gameNum))
  // // fmt.Println(float64(p2Count) / float64(gameNum))
  //
  // randomPlayoutEval := NewRandomPlayoutEval(NewRandomInstructionTrainer(mtRandom), mtRandom)
  // allNodes, err := RunMCTS(battle, 19600, math.Sqrt(2), &randomPlayoutEval, mtRandom)
  //
  // if err != nil {
  //   panic(err)
  // }
  //
  // fmt.Println(allNodes[0].AverageReward())
  //
  // for actionCmd, ucb1 := range allNodes[0].ActionCmdUCB1s {
  //   fmt.Println(actionCmd)
  //   fmt.Println(ucb1)
  //   fmt.Println(ucb1.AverageReward())
  // }
}
