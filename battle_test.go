package bippa

import (
	"fmt"
	"github.com/seehuhn/mt19937"
	"github.com/sw965/omw"
	"math/rand"
	"testing"
	"time"
)

func TestPush(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	p1Fighters := NewRentalFighters("フシギバナ", "リザードン", "カメックス")
	p2Fighters := NewRentalFighters("ギャラドス", "サンダー", "ゲンガー")

	expectedP1AttackDamages := []int{43, 45, 46, 48, 49, 51, 52, 66, 67, 69, 70, 72, 73, 75, 76, 78}

	expectedP2NoCriticalAttackDamages, err := omw.MakeSliceInt(92, 111, 2)
	if err != nil {
		panic(err)
	}

	expectedP2CriticalAttackDamages, err := omw.MakeSliceInt(138, 165, 2)
	if err != nil {
		panic(err)
	}

	expectedP2AttackDamages := append(expectedP2NoCriticalAttackDamages, expectedP2CriticalAttackDamages...)
	expectedP2AttackDamages = append(expectedP2AttackDamages, 0)

	blackSludgeHeal := 11
	p1MaxHP := 187
	p2MaxHP := 171

	testNum := 1960
	p1AttackDamageResults := make([]int, 1960)
	p2AttackDamageResults := make([]int, 1960)

	for i := 0; i < testNum; i++ {
		battle := Battle{P1Fighters: p1Fighters, P2Fighters: p2Fighters}

		battle, err := battle.Push("ギガドレイン", mtRandom)
		if err != nil {
			panic(err)
		}

		battle, err = battle.Push("こおりのキバ", mtRandom)
		if err != nil {
			panic(err)
		}

		p1AttackDamageResults[i] = p2MaxHP - battle.P2Fighters[0].CurrentHP

		p2AttackDamageResult := p1MaxHP - battle.P1Fighters[0].CurrentHP

		if p2AttackDamageResult > 0 {
			p2AttackDamageResult += blackSludgeHeal
		}

		p2AttackDamageResults[i] = p2AttackDamageResult
	}

	if !omw.SliceIntContains(expectedP1AttackDamages, p1AttackDamageResults...) {
		t.Errorf("テスト失敗")
	}

	if !omw.SliceIntContains(expectedP2AttackDamages, p2AttackDamageResults...) {
		fmt.Println(p2AttackDamageResults)
		t.Errorf("テスト失敗")
	}
}

func TestAttackDamageProbabilityDistribution(t *testing.T) {
	attackPokemon := NEW_RENTAL_POKEMONS["フシギバナ"]()
	defensePokemon := NEW_RENTAL_POKEMONS["カメックス"]()
	defensePokemon.Types = Types{"みず", "じめん"}
	adpd, err := NewAttackDamageProbabilityDistribution(&attackPokemon, &defensePokemon, "ギガドレイン", 100, 16)
	if err != nil {
		panic(err)
	}
	expected := adpd.RatioExpected(float64(defensePokemon.CurrentHP))
	fmt.Println(expected)

	adpd, _ = NewAttackDamageProbabilityDistribution(&defensePokemon, &attackPokemon, "れいとうビーム", 100, 16)
	expected = adpd.RatioExpected(float64(attackPokemon.CurrentHP))
	fmt.Println(expected)
}
