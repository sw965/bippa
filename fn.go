package bippa

import (
  "fmt"
)

func Sigmoid(x float64) float64 {
  return 1.0 / (1 + math.Exp(-x))
}

func SigmoidDerivative(y, a float64) float64 {
  return a * (1 - y) * y
}

func MeanSquaredError(y, t float64) float64 {
  return (t - y) * (t - y)
}

func SigmoidMeanSquaredErrorDerivative(y, t, a float64) float64 {
  return ((2 * t) - (2 * y)) * -SigmoidDerivative(y, a)
}
