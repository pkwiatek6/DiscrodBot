package actions

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/pkwiatek6/DiscrodBot/data"
)

//RollD10 rolls a single d10 and returns the outcome
func RollD10() int {
	return rand.Intn(10) + 1
}

//Rolls fudge dice go from min to ten
func RollDF(minFudge int) int {
	if minFudge <= 10 && minFudge > 0 {
		return rand.Intn(10-minFudge+1) + minFudge
	} else {
		return rand.Intn(10-9+1) + 9
	}
}

//FlipCoin flips a coin and returns outcome
func FlipCoin(nick string) string {
	var Coin = "Tails"
	if rand.Intn(2) == 0 {
		Coin = "Heads"
	}
	return fmt.Sprintf("```%s flipped a coin and it came up %s```", nick, Coin)
}

//CountSuc counts the number of successes contained in diceReults
func CountSuc(diceResults []int, DC int) int {
	var successes = 0
	for i := 0; i < len(diceResults); i++ {
		if diceResults[i] == 10 {
			successes += 2
		} else if diceResults[i] >= DC {
			successes++
		} else if diceResults[i] == 1 && diceResults[i] < DC {
			successes--
		}
	}
	return successes
}

//RerollDice rerolls the 3 lowest dice that are not successes from a result
func RerollDice(character *data.Character) string {
	sort.Ints(character.LastRoll.Rolls)
	var failedRolls [3]int
	//new rolls is used to show what the 3 dice were rerolled into
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
	successes := CountSuc(character.LastRoll.Rolls, character.LastRoll.DC)
	var toPost string
	if successes >= 1 {
		toPost = fmt.Sprintf("```%s got %d Successes%s\nRerolls %v -> %v```", character.Name, successes, character.LastRoll.Reason, failedRolls, newRolls)
	} else if successes == 0 {
		toPost = fmt.Sprintf("```%s Failed%s\nRerolls %v -> %v```", character.Name, character.LastRoll.Reason, failedRolls, newRolls)
	} else {
		toPost = fmt.Sprintf("```%s got a Botch%s\nRerolls %v -> %v```", character.Name, character.LastRoll.Reason, failedRolls, newRolls)
	}
	return toPost
}

//RollDice rolls the dice for a check. DC is expected. Legacy function.
func RollDice(c string, channel string, session *discordgo.Session, character *data.Character, isSpecial bool) {
	toRoll := strings.Split(c, ",")
	//var reason string
	if len(toRoll) < 2 {
		log.Println(errors.New("roll dice: not enough inputs for command"))
		return
	} else if len(toRoll) == 2 {
		character.LastRoll.Reason = ""
	} else if len(toRoll) == 3 {
		reg := regexp.MustCompile(`^\s*(?:trying)*\s*(?:to)*\s*`)
		res := reg.ReplaceAllString(toRoll[2], "${1}")
		character.LastRoll.Reason = " trying to " + res
	}
	numDice, err := strconv.Atoi(toRoll[0])
	if err != nil {
		log.Println(errors.New("roll dice: numDice was not a number"))
		return
	}
	character.LastRoll.DC, err = strconv.Atoi(toRoll[1])
	if err != nil {
		log.Println(errors.New("roll dice: DC was not a number"))
		return
	}
	session.ChannelMessageSend(channel, RollDiceCommand(numDice, character.LastRoll.DC, character.LastRoll.Reason, character))

}

func RollDiceCommand(dicepool int, dc int, reason string, character *data.Character) string {
	numDice := dicepool
	character.LastRoll.DC = dc
	//makes an integer array the size of the number of dice rolled and populates it
	character.LastRoll.Rolls = make([]int, numDice)
	if character.FudgeRoll > 0 && character.FudgeRoll <= dicepool {
		for i := 0; i < character.FudgeRoll; i++ {
			character.LastRoll.Rolls[i] = RollDF(dc)
		}
		for i := character.FudgeRoll; i < numDice; i++ {
			character.LastRoll.Rolls[i] = RollD10()
		}
		character.FudgeRoll = 0
	} else {
		for i := 0; i < numDice; i++ {
			character.LastRoll.Rolls[i] = RollD10()
		}
	}

	successes := CountSuc(character.LastRoll.Rolls, character.LastRoll.DC)

	var toPost string
	if reason != "" {
		reg := regexp.MustCompile(`^\s*(?:trying)*\s*(?:to)*\s*`)
		res := reg.ReplaceAllString(reason, "${1}")
		character.LastRoll.Reason = " trying to " + res
		if successes >= 1 {
			toPost = fmt.Sprintf("```%s got %d Successes%s\nRolled %v```", character.Name, successes, character.LastRoll.Reason, character.LastRoll.Rolls)
		} else if successes == 0 {
			toPost = fmt.Sprintf("```%s Failed%s\nRolled %v```", character.Name, character.LastRoll.Reason, character.LastRoll.Rolls)
		} else {
			toPost = fmt.Sprintf("```%s got a Botch%s\nRolled %v```", character.Name, character.LastRoll.Reason, character.LastRoll.Rolls)
		}
	} else {
		if successes >= 1 {
			toPost = fmt.Sprintf("```%s got %d Successes\nRolled %v```", character.Name, successes, character.LastRoll.Rolls)
		} else if successes == 0 {
			toPost = fmt.Sprintf("```%s Failed\nRolled %v```", character.Name, character.LastRoll.Rolls)
		} else {
			toPost = fmt.Sprintf("```%s got a Botch\nRolled %v```", character.Name, character.LastRoll.Rolls)
		}
	}
	return toPost
}

//Sets the result for the next roll the invokee make
func WouldYouKindly(minResultCmd string, channel string, session *discordgo.Session, character *data.Character) {
	if character.DiscordUser == "Dublin07#9139" {
		reg := regexp.MustCompile(`^(?:wyk)\s*`)
		res := reg.ReplaceAllString(minResultCmd, "${1}")
		minResult, err := strconv.Atoi(res)
		if err != nil {
			log.Println(err)
			session.ChannelMessageSend(channel, "You did not enter a number after /wyk")
			return
		} else {
			character.FudgeRoll = minResult
			response := fmt.Sprintf("Fudge set to %d", minResult)
			session.ChannelMessageSend(channel, response)
		}
	} else {
		session.ChannelMessageSend(channel, "No, piss off")
	}
}

//isAuthUser will be used for the updated would you kindly func when permissions get added
func isAuthUser(authRole string, list []string) bool {
	for _, v := range list {
		if strings.ToLower(v) == strings.ToLower(authRole) {
			return true
		}
	}
	return false
}

//ScheduleSession saves reminders for next sessions, will ping everyone, should only be accessible to StoryTeller
func ScheduleSession() {
	//Not yet implemented
}
