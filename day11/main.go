package main

import (
	"fmt"
)

var input = []byte("vzbxkghb")
var az = []byte("az")
var denied = []byte("iol")

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	nextPw(input)
	fmt.Printf("Part 1: %s\n", string(input))
	nextPw(input)
	fmt.Printf("Part 2: %s\n", string(input))
}

func nextPw(pw []byte) {
	last := len(pw) - 1
	//check for forbidden letters (i,o,l)
	for i, c := range pw {
		if c == denied[0] || c == denied[1] || c == denied[2] {
			for j := i + 1; j < len(pw); j++ {
				pw[j] = az[1]
			}
			break
		}
	}
	//increment the password
	for i := last; i > 0; i-- {
		pw[i]++
		if pw[i] > az[1] {
			pw[i] = az[0]
		} else {
			break
		}
	}
	//check if contains 3 consecutiv
	invalid := true
	for i := 0; i < len(pw)-3; i++ {
		if pw[i] == pw[i+1]-1 && pw[i] == pw[i+2]-2 {
			invalid = false
		}
	}
	if invalid {
		nextPw(pw)
		return
	}
	//check if contains two double chars
	invalid = true
	for i := 0; i < len(pw)-4; i++ {
		if pw[i] == pw[i+1] {
			for j := i + 1; j < len(pw)-1; j++ {
				if pw[j] == pw[j+1] {
					invalid = false
				}
			}
		}
	}
	if invalid {
		nextPw(pw)
	}
}
