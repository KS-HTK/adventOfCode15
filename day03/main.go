package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)

	fmt.Printf("Part 1: %d\n", pt1(string(dat)))
	fmt.Printf("Part 2: %d\n", pt2(string(dat)))
}

var m map[string]int

func pt1(input string) int {
	x := 0
	y := 0
	m = make(map[string]int)
	m[strconv.Itoa(x)+":"+strconv.Itoa(y)] = 0
	for i, c := range input {
		switch string(c) {
		case "^":
			y++
		case "v":
			y--
		case "<":
			x--
		case ">":
			x++
		}
		m[strconv.Itoa(x)+":"+strconv.Itoa(y)] = i
	}
	return len(m)
}

func pt2(input string) int {
	x := 0
	y := 0
	x2 := 0
	y2 := 0
	robo := false
	m = make(map[string]int)
	m[strconv.Itoa(x)+":"+strconv.Itoa(y)] = 0
	for i, c := range input {
		if robo {
			switch string(c) {
			case "^":
				y2++
			case "v":
				y2--
			case "<":
				x2--
			case ">":
				x2++
			}
			m[strconv.Itoa(x2)+":"+strconv.Itoa(y2)] = i
		} else {
			switch string(c) {
			case "^":
				y++
			case "v":
				y--
			case "<":
				x--
			case ">":
				x++
			}
			m[strconv.Itoa(x)+":"+strconv.Itoa(y)] = i
		}
		robo = !robo
	}
	return len(m)
}
