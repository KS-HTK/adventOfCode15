package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//check error, exit program if present.
func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	availableTeaspoons := 100
	//read input file and check error
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	//split at newline into lines array
	lines := strings.Split(string(dat), "\n")
	//create map of ingredientes and their properties
	ingredients := make(map[string]map[string]int)
	ingKeys := make([]string, 0)
	for _, line := range lines {
		if line == "" {
			continue
		}
		tmp := strings.Split(line, ": ")
		name := tmp[0]
		ingKeys = append(ingKeys, name)
		properties := strings.Split(tmp[1], ", ")
		ingredients[name] = make(map[string]int)
		for _, prop := range properties {
			propertie := strings.Split(prop, " ")
			propName := propertie[0]
			propValue, err := strconv.Atoi(propertie[1])
			errchk(err)
			ingredients[name][propName] = propValue
		}
	}

	//get all distributions of availableTeaspoons among 4
	dists := dist4(availableTeaspoons)
	//find max result of all distributions
	maxRes := 0
	maxResPt2 := 0
	for i := 0; i<len(dists); i++ {
		cIng := distIngred(ingKeys, ingredients, dists[i])
		cookieScore := rate(cIng)
		if cookieScore > maxRes {
			maxRes = cookieScore
		}
		if cIng["calories"] == 500 && cookieScore > maxResPt2 {
			maxResPt2 = cookieScore
		}
	}
	fmt.Printf("Part 1: %d\n", maxRes)
	fmt.Printf("Part 2: %d\n", maxResPt2)
}

//get all distributions
func dist4(max int) [][]int{
	res := make([][]int, 0, 4)
	for a := 0; a <= max; a++ {
		for b := 0; b <= max-a; b++ {
			for c := 0; c <= max-a-b; c++{
				res = append(res, []int{a, b, c, max-a-b-c})
			}
		}
	}
	return res
}

func rate(a map[string]int) int{
	res := 1
	for k, v := range a {
		if k != "calories" {
			if v < 1 {
				return 0
			}
			res *= v
		}
	}
	return res
}

func distIngred(ingKeys []string, ing map[string]map[string]int, dist []int) map[string]int{
	res := make([]map[string]int, 4)
	for i := 0; i < 4; i++ {
		res[i] = mulMapConst(ing[ingKeys[i]], dist[i])
	}
	return sumMapArray(res)
}

//multiply field by const
func mulMapConst(f map[string]int, i int) map[string]int{
	res := make(map[string]int)
	for k, v := range f {
		res[k] = v*i
	}
	return res
}

func sumMapArray(f []map[string]int) map[string]int{
	res := make(map[string]int)
	for i := 0; i<len(f); i++ {
		for k, v := range f[i] {
			res[k] += v
		}
	}
	return res
}