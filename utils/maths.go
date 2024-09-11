package utils

import (
	"math"
	"math/rand"
)

func Clamp(v int, v0 int, v1 int) int {
	if v < v0 {
		return v0
	}
	if v > v1 {
		return v1
	}
	return v
}

func AverageDiceRoll(count int, dType int) int {

	avg := (float64(dType+1) / 2.0)
	avg *= float64(count)
	avg = math.Floor(avg)
	return int(avg)

}

func RollDice(count int, dType int) int {

	out := 0

	for i := 0; i < count; i++ {
		roll := rand.Float64()*(float64(dType-1.0)) + 1.0
		roll = math.Round(roll)
		out += int(roll)

	}

	return out

}
