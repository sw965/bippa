package dmgtools

import (
	"fmt"
	"math/rand"
	bp "github.com/sw965/bippa"
)

func CriticalN(rank bp.CriticalRank) (int, error) {
	if rank < 0 {
		return 0, fmt.Errorf("急所ランクは0以上でなければならない")
	}

	var n int
	switch rank {
		case 0:
			n = 16
		case 1:
			n = 8
		case 2:
			n = 4
		case 3:
			n = 3
		default:
			n = 2
	}
	return n, nil
}

func IsCritical(rank bp.CriticalRank, r *rand.Rand) (bool, error) {
	n, err := CriticalN(rank)
	if err != nil {
		return false, err
	}
	return r.Intn(n) == 0, nil
}