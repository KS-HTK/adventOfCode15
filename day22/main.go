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
	player       player
	boss         boss
	spell_timers [3]int
	nextSpell    int
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
func shield(p *player) (int, error) {
	const cost = 113
	const duration = 6
	if p.mana < cost {
		return -1, errors.New("Not enough Mana")
	}
	p.mana -= cost
	p.spent_mana += cost
	p.def += 7
	//update spell (remove shield if timer runs out)
	return duration, nil
}
func updateShield(time int, p *player) {
	if time == 0 {
		p.def -= 7
	}
}
func poison(p *player, b *boss) (int, error) {
	const cost = 173
	const duration = 6
	if p.mana < cost {
		return -1, errors.New("Not enough Mana")
	}
	p.mana -= cost
	p.spent_mana += cost
	//report true if spell runs out
	return duration, nil
}
func updatePoison(time int, b *boss) {
	if time > 0 {
		b.hp -= 3
	}
}
func recharge(p *player) (int, error) {
	const cost = 229
	const duration = 5
	if p.mana < cost {
		return -1, errors.New("Not enough Mana")
	}
	p.mana -= cost
	p.spent_mana += cost
	//report true if spell runs out
	return duration, nil
}
func updateRecharge(time int, p *player) {
	if time > 0 {
		p.mana += 101
	}
}

//update spells function:
func updateEffects(spells *[3]int, p *player, b *boss) {
	//spell timers may go negative. No match should last long enough for a int underflow
	//letting timers go negative prevents double trigger of shield removal, etc.
	updateShield(spells[0], p)
	spells[0] -= 1
	updatePoison(spells[1], b)
	spells[1] -= 1
	updateRecharge(spells[2], p)
	spells[2] -= 1
}

//addTrace to player for inspecting a game:
func addTracePoint(p *player, b boss, spl_time [3]int, nxtSpl int) {
	p.trace = append(p.trace, traceNode{player: *p, boss: b, spell_timers: spl_time, nextSpell: nxtSpl})
}

//game function
//returns player
func play(
	player player,
	boss boss,
	spell_timers [3]int,
	nextSpell int, res chan<- player) {
	addTracePoint(&player, boss, spell_timers, nextSpell)
	//player turn:
	//activate all effects that are timer based
	updateEffects(&spell_timers, &player, &boss)
	//check incase a spell killed the boss.
	if boss.hp <= 0 {
		res <- player
		return
	}

	//player casts spell
	//0 = Magic Missile
	//1 = Drain
	//2 = Shield
	//3 = Poison
	//4 = Recharge
	switch nextSpell {
	case 0:
		err := magicMissile(&player, &boss)
		errchk(err)
	case 1:
		err := drain(&player, &boss)
		errchk(err)
	case 2:
		duration, err := shield(&player)
		errchk(err)
		spell_timers[0] = duration
	case 3:
		duration, err := poison(&player, &boss)
		errchk(err)
		spell_timers[1] = duration
	case 4:
		duration, err := recharge(&player)
		errchk(err)
		spell_timers[2] = duration
	default:
		errchk(errors.New("No Spell Cast: " + fmt.Sprint(nextSpell)))
	}

	addTracePoint(&player, boss, spell_timers, nextSpell)
	//boss turn:
	//activate all effects that are timer based
	updateEffects(&spell_timers, &player, &boss)
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

	if player.mana >= 53 {
		play(player, boss, spell_timers, 0, res)
	} else {
		//if player cant cast cheapest spell it is GAME OVER for the player
		return
	}
	if player.mana >= 73 {
		play(player, boss, spell_timers, 1, res)
	}
	if player.mana >= 113 && spell_timers[0] <= 1 {
		play(player, boss, spell_timers, 2, res)
	}
	if player.mana >= 173 && spell_timers[1] <= 1 {
		play(player, boss, spell_timers, 3, res)
	}
	if player.mana >= 229 && spell_timers[2] <= 1 {
		play(player, boss, spell_timers, 4, res)
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
		printPlayer(&cheapest_winner)
		res1 <- cheapest_winner.spent_mana
	}()
	var wg sync.WaitGroup
	for spell := 0; spell < 5; spell++ {
		tmp := spell
		wg.Add(1)
		go func() {
			play(p, b, [3]int{-1, -1, -1}, tmp, res)
			wg.Done()
		}()
	}
	wg.Wait()
	close(res)
}

//print out a player Struct
func printPlayer(p *player) {
	readableSpell := [5]string{"Magic Missile", "Drain", "Shield", "Poison", "Recharge"}
	readableEffect := [3]string{"armor", "poison", "recharge"}
	for i, t := range p.trace {
		active := [3]string{"", "", ""}
		for i, v := range t.spell_timers {
			if v > 0 {
				active[i] = readableEffect[i]
			}
		}
		turn := ""
		action := ""
		if i%2 == 0 {
			turn = "Player"
			action = "\nPlayer casts " + readableSpell[t.nextSpell] + "."
		} else {
			turn = "Boss"
			action = "\nBoss attacks for " +
				strconv.Itoa(t.boss.dmg) + " - " +
				strconv.Itoa(t.player.def) + " = " +
				strconv.Itoa(t.boss.dmg-t.player.def) + " damage!"
		}
		fmt.Println("\n--", turn, "turn --",
			"\n- Player has", t.player.hp, "hit point,",
			t.player.def, "armor",
			t.player.mana, "mana",
			"\n- Boss has", t.boss.hp, "hit points",
			"\nactive spells:", active,
			action)
	}
	fmt.Println("The Player has Spent", p.spent_mana, "mana during this game.")
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
	//player starting values:
	p := player{hp: 50, def: 0, mana: 500} //for input (actual task)
	//p := player{hp: 10, def: 0, mana: 250} //for testinput (example)

	res1, res2 := make(chan int), make(chan int)
	go func() {
		playAll(p, b, res1)
		res2 <- -1
	}()
	//Tested results: 1362 (to high)
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
