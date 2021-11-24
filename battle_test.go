package bippa

import (
  "fmt"
  "testing"
  "math/rand"
  "github.com/seehuhn/mt19937"
  "time"
)

var mtRandom = rand.New(mt19937.New())

func init() {
  mtRandom.Seed(time.Now().UnixNano())
}

func TestGigaDrain(t *testing.T) {
  p1Fighters := Fighters{
    TEST_POKEMONS["フシギバナ"](),
    TEST_POKEMONS["リザードン"](),
    TEST_POKEMONS["カメックス"](),
  }

  p2Fighters := Fighters{
    TEST_POKEMONS["カメックス"](),
    TEST_POKEMONS["リザードン"](),
    TEST_POKEMONS["フシギバナ"](),
  }

  p1Fighters[0].State.CurrentHP = 111
  p1TestCurrentHPs := TestCurrentHPs{}
  p2TestCurrentHPs := TestCurrentHPs{}

  initBattle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
  testSimuNum := 12800

  for i := 0; i < testSimuNum; i++ {
    battle, err := initBattle.Run("ギガドレイン", mtRandom)
    if err != nil {
      panic(err)
    }

    battle, err = battle.Run("れいとうビーム", mtRandom)
    if err != nil {
      panic(err)
    }

    p1CurrentHP := battle.P1Fighters[0].State.CurrentHP
    p2CurrentHP := battle.P2Fighters[0].State.CurrentHP
    p1TestCurrentHPs.Increment(p1CurrentHP)
    p2TestCurrentHPs.Increment(p2CurrentHP)
  }

  //p1GigaDrainMinDamage := 74
  p1GigaDrainMaxDamage := 134
  p1BlackSludgeFloatHeal := 187.0 / 16.0
  p1BlackSludgeHeal := int(p1BlackSludgeFloatHeal)
  expectedP2MaxHP := 155

  expectedP1MinCurrentHP := State_(1 + p1BlackSludgeHeal)
  expectedP1MaxCurrentHP := State_((111 - 62) + (p1GigaDrainMaxDamage / 2) + p1BlackSludgeHeal)
  expectedP2MinCurrentHP := State_(expectedP2MaxHP - p1GigaDrainMaxDamage)
  expectedP2MaxCurrentHP := State_(expectedP2MaxHP)

  if expectedP1MinCurrentHP != p1TestCurrentHPs.Min() {
    t.Errorf("テスト失敗")
  }

  if expectedP1MaxCurrentHP != p1TestCurrentHPs.Max() {
    t.Errorf("テスト失敗")
  }

  if expectedP2MinCurrentHP != p2TestCurrentHPs.Min() {
    t.Errorf("テスト失敗")
  }

  if expectedP2MaxCurrentHP != p2TestCurrentHPs.Max() {
    t.Errorf("テスト失敗")
  }

  fmt.Println(0.1 * 0.8)
  fmt.Println(p2TestCurrentHPs.Percent(expectedP2MaxCurrentHP))
}
