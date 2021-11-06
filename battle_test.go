package bippa

import (
  "testing"
  "fmt"
  "math/rand"
  "github.com/seehuhn/mt19937"
  "time"
)

var mtRandom = rand.New(mt19937.New())

func NewVenusaur() Pokemon {
  result, err := NewPokemon(
    "フシギバナ", "しんちょう", "しんりょく",
    "♀", "なし", MoveNames{"つるのムチ"}, PointUps{0},
    &ALL_MAX_INDIVIDUAL, &CS252_H4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewCharizard() Pokemon {
  result, err := NewPokemon(
    "リザードン", "おくびょう", "もうか",
    "♂", "なし", MoveNames{"かえんほうしゃ"}, PointUps{1},
    &ALL_MAX_INDIVIDUAL, &CS252_H4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewBlastoise() Pokemon {
  result, err := NewPokemon(
    "カメックス", "ひかえめ", "げきりゅう",
    "♂", "なし", MoveNames{"ハイドロポンプ"}, PointUps{3},
    &ALL_MAX_INDIVIDUAL, &HC252_S4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func TestHydroPump(t *testing.T) {
  mtRandom.Seed(time.Now().UnixNano())
  p1Fighters := Fighters{NewBlastoise(), NewVenusaur(), NewCharizard()}
  p2Fighters := Fighters{NewVenusaur(), NewCharizard(), NewBlastoise()}
  initBattle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}

  testNum := 1280
  attackMissCount := 0

  for i := 0; i < testNum; i++ {
    battle, err := initBattle.Run("ハイドロポンプ", mtRandom)
    if err != nil {
      panic(err)
    }

    battle, err = battle.Run("つるのムチ", mtRandom)
    if err != nil {
      panic(err)
    }

    if battle.P2Fighters[0].State.CurrentHP == p2Fighters[0].State.MaxHP {
      attackMissCount += 1
    }
  }

  fmt.Println(float64(attackMissCount) / float64(testNum))
}
