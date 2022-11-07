package bippa

import (
	"testing"
	"fmt"
	"math/rand"
	"github.com/seehuhn/mt19937"
	"time"
	"github.com/sw965/omw"
)

func TestPush(t *testing.T) {
	mtRandom := rand.New(mt19937.New())
	mtRandom.Seed(time.Now().UnixNano())

	p1Fighters := NewRentalFighters("フシギバナ", "リザードン", "カメックス")
	p2Fighters := NewRentalFighters("ギャラドス", "サンダー", "ゲンガー")

	expectedP1AttackDamages := []int{43, 45, 46, 48, 49, 51, 52, 66, 67, 69, 70, 72, 73, 75, 76, 78}

	expectedP2NoCriticalAttackDamages, err := omw.MakeSliceIntRange(92, 111, 2)
	if err != nil {
		panic(err)
	}

	expectedP2CriticalAttackDamages, err := omw.MakeSliceIntRange(138, 165, 2)
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
		battle := Battle{P1Fighters:p1Fighters, P2Fighters:p2Fighters}
		
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

	if !omw.IsSliceIntContains(expectedP1AttackDamages, p1AttackDamageResults...) {
		t.Errorf("テスト失敗")
	}

	if !omw.IsSliceIntContains(expectedP2AttackDamages, p2AttackDamageResults...) {
		fmt.Println(p2AttackDamageResults)
		t.Errorf("テスト失敗")
	}
}