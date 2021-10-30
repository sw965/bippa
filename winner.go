package bippa

type Winner struct {
  IsP1 bool
  IsP2 bool
}

var (
  WINNER_PLAYER1 = Winner{IsP1:true, IsP2:false}
  WINNER_PLAYER2 = Winner{IsP1:false, IsP2:true}
  DRAW = Winner{IsP1:false, IsP2:false}
)

func (winner *Winner) ToFloat64() float64 {
  winnerV := *winner
  if winnerV == WINNER_PLAYER1 {
    return 1.0
  }

  if winnerV == WINNER_PLAYER2 {
    return 0.0
  }
  return 0.5
}
