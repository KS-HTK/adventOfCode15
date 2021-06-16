package main

import (
	"fmt"
)

func pt1(target int) int {
	houses := make([]int, target/10+1)
	for elf := 1; elf < len(houses); elf++ {
		for house := elf; house < len(houses); house += elf {
			houses[house] += elf * 10
		}
	}
	for house, presents := range houses {
		if presents > target {
			return house
		}
	}
	return -1
}

func pt2(target int) int {
	houses := make([]int, target/11+1)
	for elf := 1; elf < len(houses); elf++ {
		visited := 0
		for house := elf; house < len(houses) && visited < 50; house += elf {
			visited++
			houses[house] += elf * 11
		}
	}
	for house, presents := range houses {
		if presents > target {
			return house
		}
	}
	return -1
}

//main function.
func main() {
	input := 34000000
	res1, res2 := make(chan int), make(chan int)
	go func() {
		res1 <- pt1(input)
		res2 <- pt2(input)
	}()
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
