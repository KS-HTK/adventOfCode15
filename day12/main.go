package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
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
	nums := getnums(string(dat))
	sumall := sum(nums)
	fmt.Printf("Part 1: %d\n", sumall)
	fmt.Printf("Part 2: %d\n", sumall-sumred(0, len(dat), dat))
}

func getnums(in string) []string {
	re := regexp.MustCompile("-?[0-9]+")
	return re.FindAllString(in, -1)
}

func sum(slc []string) int {
	acc := 0
	for _, n := range slc {
		num, err := strconv.ParseInt(n, 10, 32)
		errchk(err)
		acc += int(num)
	}
	return acc
}

var open = []byte("{")[0]
var close = []byte("}")[0]
var sqopen = []byte("[")[0]
var sqclose = []byte("]")[0]
var red = []byte("red")

func sumred(start int, end int, dat []byte) int {
	first := -1
	for i := start; i < end; i++ {
		if dat[i] == open {
			first = i
			break
		}
	}
	//if non found or start > end return 0
	if first == -1 {
		return 0
	}

	//search the matching closing bracket and look if red is contained
	last := -1
	containsRed := false
	for j := first; j < end; j++ {
		c := dat[j]
		count := 0
		if c == open || c == sqopen {
			count++
		} else if c == close || c == sqclose {
			//count can only be 0 if there is a "}" bracket.
			if count == 0 {
				last = j
				break
			} else {
				count--
			}
		} else if !containsRed && count == 0 && c == red[0] && dat[j+1] == red[1] && dat[j+2] == red[2] {
			containsRed = true
			fmt.Println("Code RED!")
		}
	}
	if containsRed {
		//found word red, return -sum of all numbers between first and last
		part := dat[first+1 : last]
		nums := getnums(string(part))
		//return sum + potential same level exclusions
		return sumred(last+1, end, dat) + sum(nums)
	}
	//if no red on this level
	return sumred(first+1, last, dat) + sumred(last+1, end, dat)
}