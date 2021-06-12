package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func interpret(lines []string) []map[string]int {
	var data []map[string]int
	for _, line := range lines {
		if line == "" {
			continue
		}
		tmp := make(map[string]int)
		sub1 := strings.SplitN(line, ": ", 2)
		mappings := strings.Split(sub1[1], ", ")
		for _, v := range mappings {
			str := strings.Split(v, ": ")
			val, err := strconv.Atoi(str[1])
			errchk(err)
			tmp[str[0]] = val
		}
		data = append(data, tmp)
	}
	return data
}

func verify(res chan<- int, index int, data map[string]int) {

	fmt.Println(index, data)
	res <- index + 1
	close(res)
}

func findSue(result chan<- int, data []map[string]int) {
	origin := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}
	for i, m := range data {
		b := true
		for k, v := range m {
			if origin[k] != v {
				b = false
			}
		}
		if b {
			result <- i + 1
		}
	}
	close(result)
}

func findSue2(result chan<- int, data []map[string]int) {
	origin := map[string]int{
		"children":    3,
		"cats":        7,
		"samoyeds":    2,
		"pomeranians": 3,
		"akitas":      0,
		"vizslas":     0,
		"goldfish":    5,
		"trees":       3,
		"cars":        2,
		"perfumes":    1,
	}
	for i, m := range data {
		b := true
		for k, v := range m {
			switch k {
			case "cats":
				if !(origin[k] < v) {
					b = false
				}
			case "trees":
				if !(origin[k] < v) {
					b = false
				}
			case "pomeranians":
				if !(origin[k] > v) {
					b = false
				}
			case "goldfish":
				if !(origin[k] > v) {
					b = false
				}
			default:
				if !(origin[k] == v) {
					b = false
				}
			}
		}
		if b {
			result <- i + 1
		}
	}
	close(result)
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), "\n")
	var dataMap, dataMap2 []map[string]int
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		dataMap = interpret(lines)
		wg.Done()
	}()
	go func() {
		dataMap2 = interpret(lines)
		wg.Done()
	}()
	wg.Wait()

	res1, res2 := make(chan int), make(chan int)
	go findSue(res1, dataMap)
	go findSue2(res2, dataMap2)

	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
