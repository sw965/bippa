package bippa

import (
  "testing"
  "fmt"
  "math/rand"
  "github.com/seehuhn/mt19937"
  "time"
)

var mtRandom = rand.New(mt19937.New())

func init() {
  mtRandom.Seed(time.Now().UnixNano())
}

func NoEffectAttackMoveHelper(initBattle Battle, p1BattleCommand, p2BattleCommand BattleCommand,
  p1PtpDamage, p2PtpDamage *PtpDamage, p1Accuracy, p2Accuracy, testNum int, permissionErrorValue float64) string {
  p1TrueDamageProbabilityDistribution := p1PtpDamage.NewDamageProbabilityDistribution(p2Accuracy)
  p2TrueDamageProbabilityDistribution := p2PtpDamage.NewDamageProbabilityDistribution(p1Accuracy)
  p1TestDamageData := TestDamageData{0:0}
  p2TestDamageData := TestDamageData{0:0}

  for i := 0; i < testNum; i++ {
    battle, err := initBattle.Run(p1BattleCommand, mtRandom)
    if err != nil {
      panic(err)
    }

    battle, err = battle.Run(p2BattleCommand, mtRandom)
    if err != nil {
      panic(err)
    }

    p1Damage := battle.P1Fighters[0].CurrentDamage()
    p2Damage := battle.P2Fighters[0].CurrentDamage()
    p1TestDamageData.Increment(p1Damage)
    p2TestDamageData.Increment(p2Damage)
  }

  p1TestDamageProbabilityDistribution := p1TestDamageData.NewDamageProbabilityDistribution()
  p1ErrorValue := p1TrueDamageProbabilityDistribution.ErrorValue(p1TestDamageProbabilityDistribution)
  for damage, errorValue := range p1ErrorValue {
    if permissionErrorValue < errorValue {
      return fmt.Sprintf("p1Damage %v errorValue = %v", damage, errorValue)
    }
  }

  p2TestDamageProbabilityDistribution := p2TestDamageData.NewDamageProbabilityDistribution()
  p2ErrorValue := p2TrueDamageProbabilityDistribution.ErrorValue(p2TestDamageProbabilityDistribution)
  for damage, errorValue := range p2ErrorValue {
    if permissionErrorValue < errorValue {
      return fmt.Sprintf("p2Damage %v errorValue = %v", damage, errorValue)
    }
  }
  return ""
}

func TestNoEffectAttackMove(t *testing.T) {
  p1Fighters := Fighters{
    NEW_TEST_POKEMON["カメックス"](), NEW_TEST_POKEMON["リザードン"](), NEW_TEST_POKEMON["フシギバナ"](),
  }

  p2Fighters := Fighters{
    NEW_TEST_POKEMON["フシギバナ"](), NEW_TEST_POKEMON["リザードン"](), NEW_TEST_POKEMON["カメックス"](),
  }

  battle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
  p1PtpDamage := PtpDamage{}
  p1PtpDamage.NoCritical = map[int]int{86:1, 90:3, 92:3, 96:3, 98:3, 102:2, 104:1}
  p1PtpDamage.Critical = map[int]int{132:2, 134:2, 138:2, 140:2, 144:2, 146:2, 150:2, 152:1, 156:1}

  p2PtpDamage := PtpDamage{}
  p2PtpDamage.NoCritical = map[int]int{28:2, 29:2, 30:5, 31:2, 32:2, 33:3}
  p2PtpDamage.Critical = map[int]int{42:3, 43:2, 44:1, 45:3, 46:2, 47:1, 48:3, 49:1, 50:1}

  errMsg := NoEffectAttackMoveHelper(battle, "ハイドロポンプ", "はなふぶき", &p1PtpDamage, &p2PtpDamage,
    80, 100, 12800, 0.01)

  if errMsg != "" {
    t.Errorf(errMsg)
  }
}
