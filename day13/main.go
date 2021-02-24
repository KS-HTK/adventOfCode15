package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), ".\n")

	fmt.Println(lines)
	fmt.Printf("Part 1: %d\n", pt1(string(dat)))
	fmt.Printf("Part 2: %d\n", pt2(string(dat)))
}

func pt1(input string) int {
	return strings.Count(input, "(") - strings.Count(input, ")")
}

func pt2(input string) int {
	lv := 0
	for i, c := range input {
		if lv == -1 {
			return i
		}
		if string(c) == "(" {
			lv++
		} else {
			lv--
		}
	}
	return -1
}
