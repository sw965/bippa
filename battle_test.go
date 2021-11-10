package bippa

import (
  "testing"
  "fmt"
  "math/rand"
  "github.com/seehuhn/mt19937"
  "time"
)

var mtRandom = rand.New(mt19937.New())

func TestNoEffectAttackMove(t *testing.T) {
  mtRandom.Seed(time.Now().UnixNano())
  p1Fighters := Fighters{
    NEW_TEST_POKEMON["カメックス"](), NEW_TEST_POKEMON["リザードン"](), NEW_TEST_POKEMON["フシギバナ"](),
  }

  p2Fighters := Fighters{
    NEW_TEST_POKEMON["フシギバナ"](), NEW_TEST_POKEMON["リザードン"](), NEW_TEST_POKEMON["カメックス"](),
  }

  initBattle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}

  p1NoCriticalDamageDetail := map[int]int{
    86:1, 90:3, 92:3, 96:3, 98:3, 102:2, 104:1,
  }
  p1CriticalDamageDetail := map[int]int{
    132:2, 134:2, 138:2, 140:2, 144:2, 146:2, 150:2, 152:1, 156:1,
  }
  p1DamageProbabilityDistribution := NewDamageProbabilityDistribution(p1NoCriticalDamageDetail, p1NoCriticalDamageDetail, 100)

  p1DamageTestCount := DamageTestCount{}
  p2DamageTestCount := DamageTestCount{}

  for i := 0; i < testNum; i++ {
    battle, err := initBattle.Run("ハイドロポンプ", mtRandom)
    if err != nil {
      panic(err)
    }

    battle, err = battle.Run("はなふぶき", mtRandom)
    if err != nil {
      panic(err)
    }

    p1Damage := battle.P1Fighters[0].CurrentDamage()
    p2Damage := battle.P2Fighters[0].CurrentDamage()

    p1DamageTestCount.Increment(p1Damage)
    p2DamageTestCount.Increment(p2Damage)
  }

  p1DamageCountPercentTest :=
  for p1Damage, percent := range p1DamageProbabilityDistribution {
    fmt.Println(p1Damage + " : ", p1Damage, float64(count) / float64(testNum))
  }

  for p2Damage, count := range p2DamageResults {
    fmt.Println("p2Damage = ", p2Damage, float64(count) / float64(testNum))
  }
}
