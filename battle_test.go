package bippa

import (
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

  testSimuNum := 12800
  p1Fighters[0].State.CurrentHP = 1
  initTestBattle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
  battleResults, err := initTestBattle.RepeatRun("ギガドレイン", "からをやぶる", testSimuNum, mtRandom)

  if err != nil {
    panic(err)
  }

  //条件付確率 p2のカメックスのCurrentHPが81(ギガドレインのダメージが74)であるという前提
  //初期HP + ギガドレインの回復量 + くろいヘドロの回復量
  p1GigaDrainHeal40Percent := battleResults.Filter(
    func(battle Battle) bool {
      return battle.P2Fighters[0].State.CurrentHP == 81
    },
  ).FilterPercent(
    func(battle Battle) bool {
      return battle.P1Fighters[0].State.CurrentHP == (1 + 37 + 11)
    },
  )

  if p1GigaDrainHeal40Percent != 1.0 {
    t.Errorf("テスト失敗")
  }

  p1GigaDrainHeal67Percent := battleResults.Filter(
    func(battle Battle) bool {
      return battle.P2Fighters[0].State.CurrentHP == 21
    },
  ).FilterPercent(
    func(battle Battle) bool {
      return battle.P1Fighters[0].State.CurrentHP == (1 + 67 + 11)
    },
  )

  if p1GigaDrainHeal67Percent != 1.0 {
    t.Errorf("テスト失敗")
  }
}
