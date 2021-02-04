package main

import (
	"crypto/md5"
	"encoding/hex"
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

var m map[string]int

func pt1(input string) int {
	for i := 1; true; i++ {
		hash := getMD5Hash(input, i)
		if strings.HasPrefix(hash, "00000") {
			return i
		}
	}
	return 0
}

func getMD5Hash(input string, num int) string {
	text := input + strconv.Itoa(num)
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func pt2(input string) int {
	for i := 1; true; i++ {
		hash := getMD5Hash(input, i)
		if strings.HasPrefix(hash, "000000") {
			return i
		}
	}
	return 0
}
