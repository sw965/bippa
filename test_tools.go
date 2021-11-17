package bippa

import (
  "math"
)

//PtpDamageの説明 (Pokemon Trainer Paradiseの略)
//https://pokemon-trainer.net/swsh/damage/ (ポケモントレーナー天国のダメージ計算ツール)
//↑のサイトでダメージ計算を行うと↓ような計算結果が出力される。

//ダメージ詳細
//はなびらのまい
//フシギバナ → カメックス
//ダメージ
//86, 90, 90, 90, 92, 92, 92, 96, 96, 96, 98, 98, 98, 102, 102, 104
//急所ダメージ
//132, 132, 134, 134, 138, 138, 140, 140, 144, 144, 146, 146, 150, 150, 152, 156
//以下略

//上記のダメージを元に、PtpDamageDataを作る。例は次のようになる。
//ptpDamage.NoCritical = map[int]int{86:1, 90:3, 92:3, 96:3, 98:3, 102:2, 104:1}
//ptpDamage.Critical = map[int]int{132:2, 134:2, 138:2, 140:2, 144:2, 146:2, 150:2, 152:1, 156:1}

//PtpDamageDataは、真の確率分布を作る為のデータであり、
//テストで出力された確率分布が、真の値とどの程度の誤差があるのかを確認する為にある。
type PtpDamage struct {
  NoCritical map[int]int
  Critical map[int]int
}

func (ptpDamage *PtpDamage) NewDamageProbabilityDistribution(accuracy int) DamageProbabilityDistribution {
	result := DamageProbabilityDistribution{}
  accuracyPercent := float64(accuracy) / float64(100.0)

	for damage, count := range ptpDamage.NoCritical {
    rPercent := float64(count) / float64(DAMAGE_RS_LENGTH)
		result[damage] = rPercent * accuracyPercent * float64(NO_CRITICAL_PERCENT)
	}

	for damage, count := range ptpDamage.Critical {
    rPercent := float64(count) / float64(DAMAGE_RS_LENGTH)
		percent := rPercent * accuracyPercent * float64(CRITICAL_PERCENT)
		if _, ok := result[damage]; ok {
			//確率の加法定理
			result[damage] += percent
		} else {
			result[damage] = percent
		}
	}
  result[0] = 1.0 - accuracyPercent
	return result
}

//key valueに入力する値はPtpDamageDataと同じである。
type TestDamageData map[int]int

func (testDamageData TestDamageData) Increment(key int) {
  if _, ok := testDamageData[key]; ok {
    testDamageData[key] += 1
  } else {
    testDamageData[key] = 1
  }
}

func (testDamageData TestDamageData) TotalCount() int {
  result := 0
  for _, count := range testDamageData {
    result += count
  }
  return result
}

func (testDamageData TestDamageData) NewDamageProbabilityDistribution() DamageProbabilityDistribution {
  result := DamageProbabilityDistribution{}
  totalCount := testDamageData.TotalCount()
  for damage, count := range testDamageData {
    result[damage] = float64(count) / float64(totalCount)
  }
  return result
}

type DamageProbabilityDistribution map[int]float64

func (d1 DamageProbabilityDistribution) ErrorValue(d2 DamageProbabilityDistribution) map[int]float64 {
  result := map[int]float64{}
  for damage, percent1 := range d1 {
    if _, ok := d2[damage]; !ok {
      result[damage] = 999.0
    } else {
      percent2 := d2[damage]
      result[damage] = math.Abs(percent1 - percent2)
    }
  }
  return result
}

func NewTestVenusaur() Pokemon {
  result, err := NewPokemon(
    "フシギバナ", "しんちょう", "しんりょく",
    "♀", "なし", MoveNames{"はなふぶき"}, PointUps{0},
    &ALL_MAX_INDIVIDUAL, &HD252_S4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

func NewTestCharizard() Pokemon {
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

func NewTestBlastoise() Pokemon {
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

func NewTestAerodactyl() Pokemon {
  result, err := NewPokemon(
    "プテラ", "ようき", "プレッシャー", "♂", "きあいのタスキ",
    MoveNames{"ストーンエッジ"}, PointUps{3},
    &ALL_MAX_INDIVIDUAL, &AS252_H4,
  )

  if err != nil {
    panic(err)
  }
  return result
}

var TEST_POKEMONS = map[PokeName]func()Pokemon{
	"フシギバナ":NewTestVenusaur,
	"リザードン":NewTestCharizard,
	"カメックス":NewTestBlastoise,
	"プテラ":NewTestAerodactyl,
}
