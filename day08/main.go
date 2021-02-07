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

	lstr := strings.Split(string(dat), "\n")
	//length of literal strings
	lstrlen := 0
	//length of interpreted strings
	istrlen := 0
	//length of escaped strings
	estrlen := 0
	escchr := regexp.MustCompile(`\\|\"`)
	for _, s := range lstr {
		if s == "" {
			continue
		}
		//Part 1 Counter
		lstrlen += len(s)
		str, err := strconv.Unquote(s)
		errchk(err)
		istrlen += len(str)
		estrlen += len(s)
		estrlen += 2
		estrlen += len(escchr.FindAll([]byte(s), -1))
	}
	test := lstr[0]
	fmt.Println(test, len(test), len(escchr.FindAll([]byte(test), -1)))

	fmt.Printf("Part 1: %d\n", lstrlen-istrlen)
	fmt.Printf("Part 2: %d\n", estrlen-lstrlen)
}
