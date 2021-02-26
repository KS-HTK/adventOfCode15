package main

import (
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
	lines := strings.Split(string(dat), "\n")
	//race end time
	end := 2503
	//list of channel from deers
	pointsChannels := make(map[string]chan int)
	//scores of deers
	scores := make(map[string]int)
	//channel for returning winning Distance/Score
	pt1res := make(chan int)
	pt2res := make(chan int)

	for _, line := range lines {
		if line == "" {
			continue
		}
		//get values for deer
		words := strings.Split(line, " ")
		name := words[0]
		speed, speederr := strconv.Atoi(words[3])
		time, timeerr := strconv.Atoi(words[6])
		pause, pauseerr := strconv.Atoi(words[13])
		errchk(speederr)
		errchk(timeerr)
		errchk(pauseerr)

		//store channels for sending points back
		pChannel := make(chan int)
		pointsChannels[name] = pChannel
		scores[name] = 0

		//start race for this deer
		go func() {
			calcDist(speed, time, pause, end, pChannel)
		}()
	}

	go func() {
		maxDist := 0
		for t := 0; t < end; t++ {
			maxDist = 0
			distance := make(map[string]int)
			for k, channel := range pointsChannels {
				distance[k] = <-channel
			}
			for _, dist := range distance {
				if dist > maxDist {
					maxDist = dist
				}
			}
			for name, dist := range distance {
				if dist == maxDist {
					scores[name] = scores[name] + 1
				}
			}
		}
		pt1res <- maxDist
		close(pt1res)
		maxScore := 0
		for _, score := range scores {
			if score > maxScore {
				maxScore = score
			}
		}
		pt2res <- maxScore
		close(pt2res)
	}()

	fmt.Printf("Part 1: %d\n", <-pt1res)
	fmt.Printf("Part 2: %d\n", <-pt2res)
}

func calcDist(speed int, time int, pause int, end int, dist chan<- int) {
	distance := 0
	travel := time
	tpause := pause
	for t := 0; t < end; t++ {
		if travel > 0 {
			travel--
			distance += speed
		} else if tpause > 0 {
			tpause--
		}
		if travel == 0 && tpause == 0 {
			travel = time
			tpause = pause
		}
		dist <- distance
	}
	close(dist)
}
