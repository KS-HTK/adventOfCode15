package main

import (
	"fmt"
	"strings"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

var input = []int{3, 1, 1, 3, 3, 2, 2, 1, 1, 3}

func main() {
	fmt.Printf("Part 1: %d\n", exec(40))
	fmt.Printf("Part 2: %d\n", exec(50))
}

func exec(a int) int {
	slc := input
	for i := 0; i < a; i++ {
		slc = lookAndSay(slc)
	}
	return len(slc)
}

func arrayToString(in []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(in), " ", delim, -1), "[]")
}

func lookAndSay(in []int) []int {
	field := []int{}
	acc := 0
	num := 0
	for _, n := range in {
		if num == 0 {
			num = n
			acc++
		} else if num == n {
			acc++
		} else {
			field = append(field, acc)
			field = append(field, num)
			acc = 1
			num = n
		}
	}
	field = append(field, acc)
	field = append(field, num)
	return field
}
