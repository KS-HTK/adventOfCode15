package main

import (
	"fmt"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func getCodeCount(col int, row int) int {
	//sum all num to col+row-1
	sumTo := (col + row - 2)
	sum := (sumTo + 1) * (sumTo / 2)
	if sumTo%2 == 1 {
		sum += sumTo/2 + 1
	}
	return sum + col
}
func nextCode(code int) int {
	return (code * 252533) % 33554393
}

func main() {
	code := 20151125
	row, col := 2981, 3075
	codeCount := getCodeCount(col, row)
	for i := 1; i < codeCount; i++ {
		code = nextCode(code)
	}

	//pt1 guess: 5 449 385< x <31 562 160, x != 16 776 489, 4 600 451, 1 350 806
	fmt.Printf("Part 1: %d\n", code)
}
