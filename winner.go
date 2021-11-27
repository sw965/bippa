package bippa

type Winner struct {
  IsP1 bool
  IsP2 bool
}

var (
  WINNER_P1 = Winner{IsP1:true, IsP2:false}
  WINNER_P2 = Winner{IsP1:false, IsP2:true}
  DRAW = Winner{IsP1:false, IsP2:false}
)

var WINNER_TO_FLOAT64 = map[Winner]float64{WINNER_P1:1.0, WINNER_P2:0.0, DRAW:0.5}
