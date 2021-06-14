package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strconv"
	"strings"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func interpret(lines []string) []int {
	var cont []int
	for _, line := range lines {
		if line != "" {
			num, err := strconv.Atoi(line)
			errchk(err)
			cont = append(cont, num)
		}
	}
	return cont
}

func findComb(start int, acc int, target int, local []int, cont []int, res chan<- []int) {
	if acc == target {
		res <- local
		return
	}
	for i := start; i < len(cont); i++ {
		if acc+cont[i] > target {
			continue
		}
		loc2 := append(make([]int, len(local)), local...)
		loc2 = append(loc2, cont[i])
		findComb(i+1, acc+cont[i], target, loc2, cont, res)
	}
}

func logRes(in <-chan []int, res1 chan<- int, res2 chan<- int) {
	min := 100
	count := 0
	allCount := 0
	for comb := range in {
		allCount += 1
		l := len(comb)
		if l < min {
			min = l
			count = 1
		} else if l == min {
			count += 1
		}
	}
	res1 <- allCount
	res2 <- count
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), "\n")
	var containers []int
	containers = interpret(lines)
	sort.Ints(containers)

	eggnog := 150
	res1, res2 := make(chan int), make(chan int)
	go func() {
		res := make(chan []int)
		go logRes(res, res1, res2)
		findComb(0, 0, eggnog, make([]int, 0), containers, res)
		close(res)
	}()

	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
