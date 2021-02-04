package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
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

func pt1(input string) int {
	str := strings.Split(input, "\n")
	dour := `aa|bb|cc|dd|ee|ff|gg|hh|ii|jj|kk|ll|mm|nn|oo|pp|qq|rr|ss|tt|uu|vv|ww|xx|yy|zz`
	nmatch := `ab|cd|pq|xy`
	nice := 0
	for _, s := range str {
		vowels, err := regexp.Match(`([aeiou].*){3,}`, []byte(s))
		errchk(err)
		doubles, err := regexp.Match(dour, []byte(s))
		errchk(err)
		noughtys, err := regexp.Match(nmatch, []byte(s))
		errchk(err)
		if !vowels || !doubles || noughtys {
			continue
		}
		nice++
	}
	return nice
}

func pt2(input string) int {
	str := strings.Split(input, "\n")
	drep := `a.a|b.b|c.c|d.d|e.e|f.f|g.g|h.h|i.i|j.j|k.k|l.l|m.m|n.n|o.o|p.p|q.q|r.r|s.s|t.t|u.u|v.v|w.w|x.x|y.y|z.z`
	doub := false
	nice := 0
	for _, s := range str {
		for i, _ := range s {
			for j := i + 2; j < len(s)-1; j++ {
				if s[i] == s[j] && s[i+1] == s[j+1] {
					doub = true
				}
			}
		}
		drepe, err := regexp.Match(drep, []byte(s))
		errchk(err)
		if !doub || !drepe {
			continue
		}
		nice++
	}
	return nice
}
