package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
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

	fmt.Printf("Part 1: %d\n", pt1(string(dat)))
	fmt.Printf("Part 2: %d\n", pt2(string(dat)))
}

//Part 1:
func pt1(input string) int {
	gifts := strings.Split(input, "\n")
	total := 0
	for i := range gifts {
		if gifts[i] != "" { //catch the empty line at end of input file
			total += paper(gifts[i])
		}
	}
	return total
}

//calculate wrapping paper needed
func paper(gift string) int {
	dim := strings.Split(gift, "x")
	l, errl := strconv.Atoi(dim[0])
	w, errw := strconv.Atoi(dim[1])
	h, errh := strconv.Atoi(dim[2])
	errchk(errl)
	errchk(errw)
	errchk(errh)
	sides := [3]int{l * w, w * h, h * l}
	return sides[0]*2 + sides[1]*2 + sides[2]*2 + min(sides)
}

//get minimum of 3 int array
func min(arr [3]int) int {
	elem := arr[0]
	for i := range arr {
		if arr[i] < elem {
			elem = arr[i]
		}
	}
	return elem
}

//Part 2:
func pt2(input string) int {
	gifts := strings.Split(input, "\n")
	total := 0
	for i := range gifts {
		if gifts[i] != "" { //catch the empty line at end of input file
			total += ribbon(gifts[i])
		}
	}
	return total
}

//Calculate required ribbon length
func ribbon(gift string) int {
	dim := strings.Split(gift, "x")
	l, errl := strconv.Atoi(dim[0])
	w, errw := strconv.Atoi(dim[1])
	h, errh := strconv.Atoi(dim[2])
	errchk(errl)
	errchk(errw)
	errchk(errh)
	return circ(l, w, h) + l*w*h
}

//Calculate shortest circumfrence
func circ(l int, w int, h int) int {
	if l > w {
		if l > h {
			return (w + h) * 2
		} else {
			return (l + w) * 2
		}
	} else {
		if w > h {
			return (l + h) * 2
		} else {
			return (l + w) * 2
		}
	}
}
