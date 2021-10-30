package bippa

import (
	"math/rand"
)

type PointUp int

const (
	MIN_POINT_UP = PointUp(0)
	MAX_POINT_UP = PointUp(3)
)

func NewRandomPointUp(random *rand.Rand) PointUp {
	return PointUp(random.Intn(int(MAX_POINT_UP) + 1))
}

func (pointUp PointUp) IsValid() bool {
	return pointUp >= MIN_POINT_UP && pointUp <= MAX_POINT_UP
}
