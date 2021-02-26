package main

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
)

func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

//sums all numbers in a list of string numbers
func sum(slc []string) int {
	acc := 0
	for _, n := range slc {
		num, err := strconv.ParseInt(n, 10, 32)
		errchk(err)
		acc += int(num)
	}
	return acc
}

//returns all numbers in a string as slice of string
func getnums(in string) []string {
	re := regexp.MustCompile("-?[0-9]+")
	return re.FindAllString(in, -1)
}

func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	nums := getnums(string(dat))
	sumall := sum(nums)
	jobs := make(chan string, 100)
	results := make(chan int, 100)
	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)
	go worker(jobs, results)
	//define sum of anything to remove as 0
	redsum := 0
	//for all recived values add them to redsum
	for val := range results {
		redsum += val
	}
	fmt.Printf("Part 1: %d\n", sumall)
	fmt.Printf("Part 2: %d\n", sumall-redsum)
}

func worker(jobs chan string, results chan<- int) {
	for jobstr := range jobs {
		results <- findred(jobstr, jobs)
	}
}

var open = []rune("{")[0]
var close = []rune("}")[0]
var sqopen = []rune("[")[0]
var sqclose = []rune("]")[0]
var red = []rune("red")

func findred(dat string, jobs chan<- string) int {
	if len(dat) < 2 {
		return 0
	}
	runeDat := []rune(dat)
	//find first '{' symbol
	begin, end := -1, -1
	//queue for any partial string between '{' and '}'
	queue := []string{}

	for i := 0; i < len(runeDat); i++ {
		//count the bracket level to make sure that red and the correct closing bracket are only found if no other bracket has been opened.
		count := 0
		c := runeDat[i]
		if count == 0 {
			if c == open {
				begin = i
				count++
			}
			if c == sqopen {
				count++
			}
			if c == close {
				end = i
				queue = append(queue, string(runeDat[begin+1:end]))
			}
			if c == red[0] {
				if runeDat[i+1] == red[1] && runeDat[i+2] == red[2] {
					fmt.Println("Code RED!")
					//if red has been found the number sum is returned
					return sum(getnums(dat))
				}
			}
		} else if c == open || c == sqopen {
			count++
		} else if c == close || c == sqclose {
			count--
		}
	}
	fmt.Println("No red found. Adding neu jobs. ")
	//if no red has been found:
	for _, s := range queue {
		jobs <- s
	}
	return 0
}
