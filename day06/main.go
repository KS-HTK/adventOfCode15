package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
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

	var gridpt1 [1000][1000]bool
	var gridpt2 [1000][1000]int
	command := strings.Split(string(dat), "\n")
	re := regexp.MustCompile(`[0-9]*,[0-9]*`)
	for _, cmd := range command {
		if cmd == "" {
			continue
		}
		res := re.FindAll([]byte(cmd), -1)
		from := strings.Split(string(res[0]), ",")
		to := strings.Split(string(res[1]), ",")

		fromx, err := strconv.Atoi(from[0])
		errchk(err)
		fromy, err := strconv.Atoi(from[1])
		errchk(err)
		tox, err := strconv.Atoi(to[0])
		errchk(err)
		toy, err := strconv.Atoi(to[1])
		errchk(err)

		if strings.Contains(cmd, "toggle") {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					gridpt1[i][j] = !gridpt1[i][j]
					gridpt2[i][j] += 2
				}
			}
		} else if strings.Contains(cmd, "off") {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					gridpt1[i][j] = false
					gridpt2[i][j]--
					if gridpt2[i][j] < 0 {
						gridpt2[i][j] = 0
					}
				}
			}
		} else {
			for i := fromx; i <= tox; i++ {
				for j := fromy; j <= toy; j++ {
					gridpt1[i][j] = true
					gridpt2[i][j]++
				}
			}
		}
	}
	acc1 := 0
	acc2 := 0
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			if gridpt1[i][j] {
				acc1++
			}
			acc2 += gridpt2[i][j]
		}
	}

	fmt.Printf("Part 1: %d\n", acc1)
	fmt.Printf("Part 2: %d\n", acc2)
}
