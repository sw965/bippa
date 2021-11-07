package bippa

import (
  "testing"
  "fmt"
  "math/rand"
  "github.com/seehuhn/mt19937"
  "time"
)

var mtRandom = rand.New(mt19937.New())

func TestHydroPump(t *testing.T) {
  mtRandom.Seed(time.Now().UnixNano())
  p1Fighters := Fighters{
    NEW_TEST_POKEMON["カメックス"](), NEW_TEST_POKEMON["リザードン"], NEW_TEST_POKEMON["フシギバナ"],
  }

  p2Fighters := Fighters{
    NEW_TEST_POKEMON["フシギバナ"](), NEW_TEST_POKEMON["リザードン"], NEW_TEST_POKEMON["カメックス"],
  }

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
