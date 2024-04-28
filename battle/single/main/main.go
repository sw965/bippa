package main

import (
	"fmt"
	"github.com/sw965/crow/model"
	bp "github.com/sw965/bippa"
	"github.com/sw965/crow/tensor"
)

func main() {
	n := len(bp.ALL_MOVE_NAMES)
	w := model.D2Var{
		Param:tensor.NewD2Zeros(n, n),
	}
	w.Init(0.9)

	b := model.D1Var{
		Param:tensor.NewD1Zeros(n),
	}
	b.Init(0.9)

	linear := model.NewD2LinearSumTanhMSE(w, b, 0.0001)
	fmt.Println(w.Param, linear)
	y, err := linear.Predict()
}