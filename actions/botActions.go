package actions

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkwiatek6/DiscrodBot/data"
)

//RollD10 rolls a single d10 and returns the outcome
func RollD10() int {
	return 1 + rand.Intn(9)
}

//FlipCoin flips a coin and returns outcome
func FlipCoin(channel string, nick string, session *discordgo.Session) {
	if rand.Intn(2) == 0 {
		session.ChannelMessageSend(channel, fmt.Sprintf("```%s flipped a coin and it came up %s```", nick, "Heads"))
		return
	}
	session.ChannelMessageSend(channel, fmt.Sprintf("```%s flipped a coin and it came up %s```", nick, "Tails"))
	return
}

//CountSuc counts the number of successes contained in diceReults
func CountSuc(diceResults []int, DC int) int {
	var successes = 0
	for i := 0; i < len(diceResults); i++ {
		if diceResults[i] == 10 {
			successes += 2
		} else if diceResults[i] >= DC && diceResults[i] != 10 {
			successes++
		} else if diceResults[i] == 1 && diceResults[i] < DC {
			successes--
		}
	}
	return successes
}

//RerollDice rerolls the 3 lowest dice that are not successes from a result
func RerollDice(name string, channel string, session *discordgo.Session, lastRolls *map[string]data.RollHistory) {
	var oldResults = (*lastRolls)[name].Rolls
	sort.Ints(oldResults)
	var tempDC = (*lastRolls)[name].DC
	var failedRolls [3]int
	var newRolls [3]int
	var max = 3
	if len(oldResults) < 3 {
		max = len(oldResults)
	}
	for i := 0; i < max; i++ {
		if oldResults[i] < tempDC && oldResults[i] != 10 {
			failedRolls[i] = oldResults[i]
			newRolls[i] = RollD10()
			oldResults[i] = newRolls[i]
		}
	}
	successes := CountSuc(oldResults, tempDC)
	(*lastRolls)[name] = data.RollHistory{Rolls: oldResults, DC: tempDC, Reason: (*lastRolls)[name].Reason}
	if successes >= 1 {
		toPost := fmt.Sprintf("```%s got %d Successes%s\nRerolls %v -> %v```", name, successes, (*lastRolls)[name].Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)

	} else if successes == 0 {
		toPost := fmt.Sprintf("```%s Failed%s\nRerolls %v -> %v```", name, (*lastRolls)[name].Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)
	} else {
		toPost := fmt.Sprintf("```%s got a Botch%s\nRerolls %v -> %v```", name, (*lastRolls)[name].Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)
	}
}

//RollDice rolls the dice for a check. DC is expected
func RollDice(c string, nick string, channel string, session *discordgo.Session, lastRolls *map[string]data.RollHistory) {
	var reason string
	toRoll := strings.Split(c, ",")
	if len(toRoll) < 2 {
		log.Println(errors.New("Roll Dice: Not enough inputs for command"))
		return
	}
	if len(toRoll) == 3 {
		reason = " trying to " + toRoll[2]
	}
	numDice, err := strconv.Atoi(toRoll[0])
	if err != nil {
		log.Println(errors.New("Roll Dice: numDice was not a number"))
		return
	}
	DC, err := strconv.Atoi(toRoll[1])
	if err != nil {
		log.Println(errors.New("Roll Dice: DC was not a number"))
		return
	}
	//makes an integer array the size of the number of dice rolled and populates it
	diceResults := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		diceResults[i] = RollD10()
	}
	successes := CountSuc(diceResults, DC)
	//TODO: recode to user unique identifer instead of nickname so rolls aren;t lost
	(*lastRolls)[nick] = data.RollHistory{Rolls: diceResults, DC: DC, Reason: reason}
	if successes >= 1 {
		toPost := fmt.Sprintf("```%s got %d Successes%s\nRolled %v```", nick, successes, reason, diceResults)
		session.ChannelMessageSend(channel, toPost)
	} else if successes == 0 {
		toPost := fmt.Sprintf("```%s Failed%s\nRolled %v```", nick, reason, diceResults)
		session.ChannelMessageSend(channel, toPost)
	} else {
		toPost := fmt.Sprintf("```%s got a Botch%s\nRolled %v```", nick, reason, diceResults)
		session.ChannelMessageSend(channel, toPost)
	}
}
