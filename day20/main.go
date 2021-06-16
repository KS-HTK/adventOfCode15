package main

import (
	"fmt"
	"sort"
)

//error checking function
//causes panic if error occurs
func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

//returns prime factorisation of number num
func primefactors(num int) []int {
	p := 2
	for num >= p*p {
		if num%p == 0 {
			return append(primefactors(num/p), p)
		}
		p++
	}
	return []int{num}
}

//returns all divisors of number num
func divisors(num int) []int {
	pf := primefactors(num)
	prodArr := getAllProducts(pf)
	//writing all to a temporary map and back to slice to remove duplicates
	tmpMap := make(map[int]int)
	for _, prod := range prodArr {
		tmpMap[prod] = 0
	}
	for _, prime := range pf {
		tmpMap[prime] = 0
	}
	div := make([]int, 0)
	for key := range tmpMap {
		div = append(div, key)
	}
	div = append(div, 1)
	sort.Ints(div)
	return div
}

func getAllProducts(in []int) []int {
	if len(in) < 2 {
		return in
	}
	products := getAllProducts(in[1:])
	tmp := make([]int, 0)
	for _, a := range products {
		tmp = append(tmp, a*in[0])
	}
	tmp = append(tmp, in[0])
	return append(products, tmp...)
}

func sum(in []int) int {
	result := 0
	for _, v := range in {
		result += v
	}
	return result
}

func pt1(target int) int {
	target = target / 10
	for i := 13; i < target; i++ {
		if sum(divisors(i)) >= target {
			return i
		}
	}
	return -1
}

func modSum(in []int, i int) int {
	result := 0
	for _, v := range in {
		if v*50 <= i {
			result += v
		}
	}
	return result
}

//4989600, 720720
func pt2(target int) int {
	for i := 13; i <= target; i++ {
		if modSum(divisors(i), i)*11 >= target {
			return i
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
