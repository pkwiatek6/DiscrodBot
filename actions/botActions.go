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
// Extra logic is needed because 10's count as 2 success only when it is your specialty otherwise they count as 1
func CountSuc(diceResults []int, DC int, isSpecial bool) int {
	var successes = 0
	for i := 0; i < len(diceResults); i++ {
		if diceResults[i] == 10 && isSpecial {
			successes += 2
		} else if diceResults[i] >= DC && !isSpecial {
			successes++
		} else if diceResults[i] == 1 && diceResults[i] < DC {
			successes--
		}
	}
	return successes
}

//RerollDice rerolls the 3 lowest dice that are not successes from a result
func RerollDice(character *data.Character, channel string, session *discordgo.Session) {
	sort.Ints(character.LastRoll.Rolls)
	var failedRolls [3]int
	var newRolls [3]int
	var max = 3
	if len(character.LastRoll.Rolls) < 3 {
		max = len(character.LastRoll.Rolls)
	}
	for i := 0; i < max; i++ {
		if character.LastRoll.Rolls[i] < character.LastRoll.DC && character.LastRoll.Rolls[i] != 10 {
			failedRolls[i] = character.LastRoll.Rolls[i]
			newRolls[i] = RollD10()
			character.LastRoll.Rolls[i] = newRolls[i]
		}
	}
	//TODO: Find a way to keep track if previous roll was special.
	successes := CountSuc(character.LastRoll.Rolls, character.LastRoll.DC, false)
	//new rolls is used to show what the 3 dice were rerolled into
	if successes >= 1 {
		toPost := fmt.Sprintf("```%s got %d Successes%s\nRerolls %v -> %v```", character.Name, successes, character.LastRoll.Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)

	} else if successes == 0 {
		toPost := fmt.Sprintf("```%s Failed%s\nRerolls %v -> %v```", character.Name, character.LastRoll.Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)
	} else {
		toPost := fmt.Sprintf("```%s got a Botch%s\nRerolls %v -> %v```", character.Name, character.LastRoll.Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)
	}
}

//RollDice rolls the dice for a check. DC is expected
func RollDice(c string, channel string, session *discordgo.Session, character *data.Character, isSpecial bool) {
	toRoll := strings.Split(c, ",")
	//var reason string
	if len(toRoll) < 2 {
		log.Println(errors.New("Roll Dice: Not enough inputs for command"))
		return
	} else if len(toRoll) == 2 {
		character.LastRoll.Reason = ""
	} else if len(toRoll) == 3 {
		character.LastRoll.Reason = " trying to " + toRoll[2]
	}
	numDice, err := strconv.Atoi(toRoll[0])
	if err != nil {
		log.Println(errors.New("Roll Dice: numDice was not a number"))
		return
	}
	character.LastRoll.DC, err = strconv.Atoi(toRoll[1])
	if err != nil {
		log.Println(errors.New("Roll Dice: DC was not a number"))
		return
	}
	//makes an integer array the size of the number of dice rolled and populates it
	character.LastRoll.Rolls = make([]int, numDice)
	for i := 0; i < numDice; i++ {
		character.LastRoll.Rolls[i] = RollD10()
	}
	successes := CountSuc(character.LastRoll.Rolls, character.LastRoll.DC, isSpecial)
	//character.LastRoll = data.RollHistory{Rolls: diceResults, DC: character.LastRoll.DC, Reason: reason}
	if successes >= 1 {
		toPost := fmt.Sprintf("```%s got %d Successes%s\nRolled %v```", character.Name, successes, character.LastRoll.Reason, character.LastRoll.Rolls)
		session.ChannelMessageSend(channel, toPost)
	} else if successes == 0 {
		toPost := fmt.Sprintf("```%s Failed%s\nRolled %v```", character.Name, character.LastRoll.Reason, character.LastRoll.Rolls)
		session.ChannelMessageSend(channel, toPost)
	} else {
		toPost := fmt.Sprintf("```%s got a Botch%s\nRolled %v```", character.Name, character.LastRoll.Reason, character.LastRoll.Rolls)
		session.ChannelMessageSend(channel, toPost)
	}
}

//ScheduleSession saves reminders for next sessions, will ping everyone, should only be accessible to StoryTeller
func ScheduleSession() {

}
