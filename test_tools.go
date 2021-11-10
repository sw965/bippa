package bippa

type DamageCountTest map[int]int

func (damageCountTest DamageCountTest) Increment(key int) {
  if _, ok := damageCountTest[key]; ok {
    damageCountTest[key] += 1
  } else {
    damageCountTest[key] = 1
  }

}

func (damageCountTest DamageCountTest) TotalCount() int {
  result := 0
  for _, count := range damageCountTest {
    result += count
  }
  return result
}

func (damageCountTest DamageCountTest) NewDamageProbabilityDistribution() DamageProbabilityDistribution {
  result := DamageProbabilityDistribution{}
  totalCount := damageCountTest.TotalCount()
  for damage, count := range damageCountTest {
    result[damage] = float64(count) / float64(totalCount)
  }
  return result
}

type DamageProbabilityDistribution map[int]float64

func NewDamageProbabilityDistribution(noCriticalDamageDetail criticalDamageDetail map[int]int, accuracy int) map[int]float64 {
	result := map[int]float64{}
	for damage, count := range noCriticalDamageDetail {
		result[damage] = float64(count) * float64(DAMAGE_RS_LENGTH) * float64(accuracy) * float64(NO_CRITICAL_PERCENT)
	}

	for damage, count := range criticalDamageDetail {
		percent := float64(count) * float64(DAMAGE_RS_LENGTH) * float64(accuracy) * float64(CRITICAL_PERCENT)
		if _, ok := result[damage]; ok {
			//確率の加法定理
			result[damage] += percent
		} else {
			result[damage] = percent
		}
	}
	return result
}

func (damageProbabilityDistribution DamageProbabilityDistribution) PrintErrorValue() {

}
