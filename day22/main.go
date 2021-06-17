package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"
)

//error checking function
//causes panic if error occurs
func errchk(e error) {
	if e != nil {
		panic(e)
	}
}

type player struct {
	hp         int
	def        int
	mana       int
	spent_mana int
	trace      []traceNode
}
type boss struct {
	hp  int
	dmg int
}

type traceNode struct {
	player        player
	boss          boss
	active_spells [3]func() bool
	nextSpell     int
}

var cheapest_winner player

func magicMissile(player *player, boss *boss) error {
	const cost = 53
	if player.mana < cost {
		return errors.New("Not enough Mana")
	}
	player.mana -= cost
	player.spent_mana += cost
	boss.hp -= 4
	return nil
}
func drain(player *player, boss *boss) error {
	const cost = 73
	if player.mana < cost {
		return errors.New("Not enough Mana")
	}
	player.mana -= cost
	player.spent_mana += cost
	boss.hp -= 2
	player.hp += 2
	return nil
}
func shield(player *player) (func() bool, error) {
	const cost = 113
	if player.mana < cost {
		return nil, errors.New("Not enough Mana")
	}
	player.mana -= cost
	player.spent_mana += cost
	timer := 6
	player.def += 7
	//report true if spell runs out
	return func() bool {
		if timer > 0 {
			timer--
			if timer == 0 {
				player.def -= 7
				return true
			}
			return false
		}
		return true
	}, nil
}
func poison(player *player, boss *boss) (func() bool, error) {
	const cost = 173
	if player.mana < cost {
		return nil, errors.New("Not enough Mana")
	}
	player.mana -= cost
	player.spent_mana += cost
	timer := 6
	//report true if spell runs out
	return func() bool {
		if timer > 0 {
			boss.hp -= 3
			timer--
		}
		if timer > 0 {
			return false
		}
		return true
	}, nil
}
func recharge(player *player) (func() bool, error) {
	const cost = 229
	if player.mana < cost {
		return nil, errors.New("Not enough Mana")
	}
	player.mana -= cost
	player.spent_mana += cost
	timer := 5
	//report true if spell runs out
	return func() bool {
		if timer > 0 {
			player.mana += 101
			timer--
		}
		if timer > 0 {
			return false
		}
		return true
	}, nil
}

//update spells function:
func updateEffects(spells *[3]func() bool) {
	for i, f := range spells {
		if f != nil {
			b := f()
			if b {
				spells[i] = nil
			}
		}
	}
}

//addTrace to player for inspecting a game:
func addTracePoint(p *player, b boss, act_spl [3]func() bool, nxtSpl int) {
	p.trace = append(p.trace, traceNode{player: *p, boss: b, active_spells: act_spl, nextSpell: nxtSpl})
}

//game function
//returns player
func play(
	player player,
	boss boss,
	active_spells [3]func() bool,
	nextSpell int, res chan<- player) {
	addTracePoint(&player, boss, active_spells, nextSpell)
	//player turn:
	//activate all effects that are timer based
	updateEffects(&active_spells)
	//check incase a spell killed the boss.
	if boss.hp <= 0 {
		res <- player
		return
	}

	//player casts spell
	//0 = Magic Missile
	//1 = Drain
	//2 = Shield   = ac…_sp…0
	//3 = Poison   = ac…_sp…1
	//4 = Recharge = ac…_sp…2
	switch nextSpell {
	case 0:
		err := magicMissile(&player, &boss)
		errchk(err)
	case 1:
		err := drain(&player, &boss)
		errchk(err)
	case 2:
		effect, err := shield(&player)
		errchk(err)
		active_spells[0] = effect
	case 3:
		effect, err := poison(&player, &boss)
		errchk(err)
		active_spells[1] = effect
	case 4:
		effect, err := recharge(&player)
		errchk(err)
		active_spells[2] = effect
	default:
		errchk(errors.New("No Spell Cast: " + fmt.Sprint(nextSpell)))
	}

	addTracePoint(&player, boss, active_spells, nextSpell)
	//boss turn:
	//activate all effects that are timer based
	updateEffects(&active_spells)
	//check incase a spell killed the boss.
	if boss.hp <= 0 {
		res <- player
		return
	}
	dmg := boss.dmg - player.def
	if dmg < 1 {
		dmg = 1
	}
	player.hp -= dmg
	if player.hp <= 0 {
		//res <- player
		return
	}
	//if no one wins play next round:
	//optimization: if spent mana runs over current cheapest_win exit
	if player.spent_mana > cheapest_winner.spent_mana {
		return
	}

	if player.mana > 53 {
		play(player, boss, active_spells, 0, res)
	} else {
		//if player cant cast cheapest spell it is GAME OVER for the player
		return
	}
	if player.mana > 73 {
		play(player, boss, active_spells, 1, res)
	}
	if player.mana > 113 && active_spells[0] == nil {
		play(player, boss, active_spells, 2, res)
	}
	if player.mana > 173 && active_spells[1] == nil {
		play(player, boss, active_spells, 3, res)
	}
	if player.mana > 229 && active_spells[2] == nil {
		play(player, boss, active_spells, 4, res)
	}
}

func playAll(p player, b boss, res1 chan<- int) {
	cheapest_winner = player{spent_mana: 6000}
	res := make(chan player, 2000)
	go func() {
		for pl := range res {
			if pl.hp > 0 && pl.spent_mana < cheapest_winner.spent_mana {
				cheapest_winner = pl
			}
		}
		//printPlayer(&cheapest_winner)
		res1 <- cheapest_winner.spent_mana
	}()
	var wg sync.WaitGroup
	for spell := 0; spell < 5; spell++ {
		tmp := spell
		wg.Add(1)
		go func() {
			play(p, b, [3]func() bool{nil, nil, nil}, tmp, res)
			wg.Done()
		}()
	}
	wg.Wait()
	close(res)
}

//print out a player Struct
func printPlayer(p *player) {
	readableSpell := [3]string{"armor", "poison", "recharge"}
	for i, t := range p.trace {
		active := [3]string{"", "", ""}
		for i, v := range t.active_spells {
			if v != nil {
				active[i] = readableSpell[i]
			}
		}
		fmt.Println("Round:", i/2,
			"Player: ",
			"hp:", t.player.hp,
			"mana:", t.player.mana,
			"Boss hp:", t.boss.hp,
			"next spell:", t.nextSpell,
			"active spells:", active)
	}
	fmt.Println("hp: ", p.hp)
	fmt.Println("mana: ", p.mana)
	fmt.Println("spent mana: ", p.spent_mana)
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
	b := boss{hp: hp, dmg: dmg}
	p := player{hp: 50, def: 0, mana: 500}

	res1, res2 := make(chan int), make(chan int)
	go func() {
		playAll(p, b, res1)
		res2 <- -1
	}()
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
