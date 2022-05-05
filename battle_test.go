package bippa

import (
  "testing"
  "fmt"
  "math"
  "math/rand"
  "time"
  "github.com/seehuhn/mt19937"
)

func NewHusigibana() Pokemon {
  pokemon, err := NewPokemon("フシギバナ", "おだやか", "しんりょく", "♀", "くろいヘドロ",
    MoveNames{"つるのムチ"}, PointUps{MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{HP:252, Atk:252, Speed:4})
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewRiza_don() Pokemon {
  pokemon, err := NewPokemon("リザードン", "ひかえめ", "もうか", "♂", "いのちのたま",
                        MoveNames{"ひのこ"}, PointUps{MAX_POINT_UP},
                        &ALL_MAX_INDIVIDUAL, &Effort{HP:252, SpAtk:252, Speed:4})
  if err != nil {
    panic(err)
  }
  return pokemon
}

func NewKamekkusu() Pokemon {
  pokemon, err := NewPokemon("カメックス", "ひかえめ", "げきりゅう", "♂", "オボンのみ",
    MoveNames{"みずでっぽう"}, PointUps{MAX_POINT_UP},
    &ALL_MAX_INDIVIDUAL, &Effort{HP:252, SpAtk:252, Speed:4})
  if err != nil {
    panic(err)
  }
  return pokemon
}

func Test(t *testing.T) {
  p1Fighters := Fighters{NewHusigibana(), NewRiza_don(), NewKamekkusu()}
  p2Fighters := Fighters{NewKamekkusu(), NewRiza_don(), NewHusigibana()}
  battle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
  mtRandom := rand.New(mt19937.New())
  mtRandom.Seed(time.Now().UnixNano())
  randomPlayoutEval := NewRandomPlayoutEval(NewRandomInstructionTrainer(mtRandom), mtRandom)
  allNodes, err := RunDMCTS(battle, 5600, math.Sqrt(2), &randomPlayoutEval, mtRandom)

  if err != nil {
    panic(err)
  }

  for action, decoupledUCT := range allNodes[0].DecoupledUCTs {
    fmt.Println(action)
    fmt.Println(decoupledUCT)
  }
}
