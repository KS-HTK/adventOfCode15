package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

//error checking function
//causes panic if error occurs
func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

type human struct {
	hp         int
	dmg        int
	def        int
	weapon     item
	armor      item
	ring1      item
	ring2      item
	spent_gold int
}

type item struct {
	name string
	cost int
	dmg  int
	def  int
}

var weapons []item
var armors []item
var rings []item

//game function
//returns true if p1 wins
func play(p1 human, p2 human) bool {
	//p1 attack:
	dmg := p1.dmg - p2.def
	if dmg < 1 {
		dmg = 1
	}
	p2.hp -= dmg
	if p2.hp <= 0 {
		return true
	}
	//p2 attack:
	dmg = p2.dmg - p1.def
	if dmg < 1 {
		dmg = 1
	}
	p1.hp -= dmg
	if p1.hp <= 0 {
		return false
	}
	//if no one wins play next round:
	return play(p1, p2)
}

func updatePlayerStats(p human) human {
	p.dmg = p.weapon.dmg + p.armor.dmg + p.ring1.dmg + p.ring2.dmg
	p.def = p.weapon.def + p.armor.def + p.ring1.def + p.ring2.def
	p.spent_gold = p.weapon.cost + p.armor.cost + p.ring1.cost + p.ring2.cost
	return p
}

func playAll(boss human) []int {
	cheapest_win := 10000000
	expensiv_loss := 0
	for _, weapon := range weapons {
		for _, armor := range armors {
			for i, ring1 := range rings {
				for j, ring2 := range rings {
					if i == j {
						continue
					}
					player := human{hp: 100, dmg: 0, def: 0, weapon: weapon, armor: armor, ring1: ring1, ring2: ring2, spent_gold: 0}
					player = updatePlayerStats(player)

					cloned_boss := human{hp: boss.hp, dmg: boss.dmg, def: boss.def}
					if play(player, cloned_boss) {
						if player.spent_gold < cheapest_win {
							cheapest_win = player.spent_gold
						}
					} else {
						if player.spent_gold > expensiv_loss {
							expensiv_loss = player.spent_gold
						}
					}
				}
			}
		}
	}
	return []int{cheapest_win, expensiv_loss}
}

//main function.
func main() {
	dat, err := ioutil.ReadFile("input")
	errchk(err)
	lines := strings.Split(string(dat), "\n")
	hp, err := strconv.Atoi(strings.Split(lines[0], ": ")[1])
	errchk(err)
	dmg, err := strconv.Atoi(strings.Split(lines[1], ": ")[1])
	errchk(err)
	def, err := strconv.Atoi(strings.Split(lines[2], ": ")[1])
	errchk(err)
	boss := human{hp: hp, dmg: dmg, def: def}
	weapons = []item{
		{name: "Dagger", cost: 8, dmg: 4, def: 0},
		{name: "Shortsword", cost: 10, dmg: 5, def: 0},
		{name: "Warhammer", cost: 25, dmg: 6, def: 0},
		{name: "Longsword", cost: 40, dmg: 7, def: 0},
		{name: "Greataxe", cost: 74, dmg: 8, def: 0},
	}
	armors = []item{
		{name: "No Armor", cost: 0, dmg: 0, def: 0},
		{name: "Leather", cost: 13, dmg: 0, def: 1},
		{name: "Chainmail", cost: 31, dmg: 0, def: 2},
		{name: "Splintmail", cost: 53, dmg: 0, def: 3},
		{name: "Bandedmail", cost: 75, dmg: 0, def: 4},
		{name: "Platemail", cost: 102, dmg: 0, def: 5},
	}
	rings = []item{
		{name: "No Ring", cost: 0, dmg: 0, def: 0},
		{name: "No Ring", cost: 0, dmg: 0, def: 0},
		{name: "Damage +1", cost: 25, dmg: 1, def: 0},
		{name: "Damage +2", cost: 50, dmg: 2, def: 0},
		{name: "Damage +3", cost: 100, dmg: 3, def: 0},
		{name: "Defense +1", cost: 20, dmg: 0, def: 1},
		{name: "Defense +2", cost: 40, dmg: 0, def: 2},
		{name: "Defense +3", cost: 80, dmg: 0, def: 3},
	}
	res1, res2 := make(chan int), make(chan int)
	go func() {
		res := playAll(boss)
		res1 <- res[0]
		res2 <- res[1]
	}()
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
