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
var cheapest_hm_winner player

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
	p.def = 7
	//update spell (remove shield if timer runs out)
	return duration, nil
}
func updateShield(time int, p *player) {
	if time == 0 {
		p.def = 0
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
	nextSpell int,
	hardMode bool,
	res chan<- player) {
	addTracePoint(&player, boss, spell_timers, nextSpell)
	//player turn:
	//if hard mode deal 1 dmg to player (check if dead)
	if hardMode {
		player.hp -= 1
		if player.hp <= 0 {
			return
		}
	}
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
		if err != nil {
			return
		}
	case 1:
		err := drain(&player, &boss)
		if err != nil {
			return
		}
	case 2:
		duration, err := shield(&player)
		if err != nil || spell_timers[0] > 0 {
			return
		}
		spell_timers[0] = duration
	case 3:
		duration, err := poison(&player, &boss)
		if err != nil || spell_timers[1] > 0 {
			return
		}
		spell_timers[1] = duration
	case 4:
		duration, err := recharge(&player)
		if err != nil || spell_timers[2] > 0 {
			return
		}
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
	var wg sync.WaitGroup
	//determining possible spells here is not as trivial as was implemented before
	//recharge can provide mana for use.
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(s int) {
			play(player, boss, spell_timers, s, hardMode, res)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

func playAll(p player, b boss, res1 chan<- int, res2 chan<- int) {
	cheapest_winner = player{spent_mana: 12000}
	res := make(chan player, 2000)
	cheapest_hm_winner = player{spent_mana: 12000}
	res_hm := make(chan player, 2000)
	go func() {
		for pl := range res {
			if pl.hp > 0 && pl.spent_mana < cheapest_winner.spent_mana {
				cheapest_winner = pl
			}
		}
		res1 <- cheapest_winner.spent_mana
	}()
	go func() {
		for pl := range res_hm {
			if pl.hp > 0 && pl.spent_mana < cheapest_winner.spent_mana {
				cheapest_hm_winner = pl
			}
		}
		printGameTrace(&cheapest_hm_winner)
		res2 <- cheapest_hm_winner.spent_mana
	}()
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	for spell := 0; spell < 5; spell++ {
		wg.Add(1)
		go func(s int) {
			play(player{hp: p.hp, def: p.def, mana: p.mana}, boss{hp: b.hp, dmg: b.dmg}, [3]int{-1, -1, -1}, s, false, res)
			wg.Done()
		}(spell)
		wg2.Add(1)
		go func(s int) {
			play(p, b, [3]int{-1, -1, -1}, s, true, res_hm)
			wg2.Done()
		}(spell)
	}
	go func() {
		wg2.Wait()
		close(res_hm)
	}()
	wg.Wait()
	close(res)
}

//print out a player Struct
func printGameTrace(p *player) {
	readableSpell := [5]string{"Magic Missile", "Drain", "Shield", "Poison", "Recharge"}
	readableEffect := [3]string{"Shield", "Poison", "Recharge"}
	for i, t := range p.trace {
		active := ""
		for i, v := range t.spell_timers {
			if v > 0 {
				active += readableEffect[i] + ": " + strconv.Itoa(v) + ", "
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
		playAll(p, b, res1, res2)
	}()
	//Tested results: 1362 (to high)
	fmt.Printf("Part 1: %d\n", <-res1)
	fmt.Printf("Part 2: %d\n", <-res2)
}
