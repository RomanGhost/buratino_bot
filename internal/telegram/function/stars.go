package function

import "math"

const STAR_PRICE = 1790

func CountStar(integerPart, fractionalPart uint64) int64 {
	amount := integerPart*1000 + fractionalPart*10
	return int64(math.Ceil(float64(amount) / float64(STAR_PRICE)))
}

func GetMoneyFromStar(starCount int) (uint64, uint64) {
	amount := starCount * STAR_PRICE
	return uint64(amount / 1000), uint64(amount % 1000)
}
