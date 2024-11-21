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

// RollD10 rolls a single d10 and returns the outcome
func RollD10() int {
	return rand.Intn(10) + 1
}

// Rolls fudge dice go from min to ten
func RollDF(minFudge int) int {
	if minFudge >= 1 && minFudge <= 10 {
		return rand.Intn(11-minFudge) + minFudge
	}
	return 10
}

// FlipCoin flips a coin and returns outcome
func FlipCoin(nick string) string {
	var Coin = "Tails"
	if rand.Intn(2) == 0 {
		Coin = "Heads"
	}
	return fmt.Sprintf("```%s flipped a coin and it came up %s```", nick, Coin)
}

// CountSuc counts the number of successes contained in diceReults
func CountSuc(diceResults []int, DC int) int {
	var successes = 0
	for _, roll := range diceResults {
		if roll == 10 {
			//DC 10s only give a success if you you roll a ten but, not 2 successes
			if DC < 10 {
				successes += 2
			} else {
				successes += 1
			}
			continue // Skip further checks for this roll
		}

		if roll >= DC {
			successes += 1
		}
		if roll == 1 {
			successes--
		}
	}
	return successes
}

// RerollDice rerolls the 3 lowest dice that are not successes from a result
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

// RollDice rolls the dice for a check. DC is expected. Legacy function.
func RollDice(c string, channel string, session *discordgo.Session, character *data.Character) {
	toRoll := strings.Split(c, ",")
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
		session.ChannelMessageSend(channel, "Number of dice is not a number (did you add a space?)")
		return
	}
	character.LastRoll.DC, err = strconv.Atoi(toRoll[1])
	if err != nil {
		log.Println(errors.New("roll dice: DC was not a number"))
		session.ChannelMessageSend(channel, "DC is not a number (did you add a space?)")
		return
	}
	session.ChannelMessageSend(channel, RollDiceCommand(numDice, character.LastRoll.DC, character.LastRoll.Reason, character))

}

func RollDiceCommand(dicepool int, dc int, reason string, character *data.Character) string {
	numDice := dicepool
	character.LastRoll.DC = dc
	character.LastRoll.Rolls = make([]int, numDice)

	if character.FudgeRoll > 0 {
		if character.FudgeRoll > numDice {
			character.FudgeRoll = numDice //Can't fudge more than you have dice
		}

		log.Println("Fudging the roll")
		// Fudges the first N rolls
		for i := 0; i < character.FudgeRoll; i++ {
			character.LastRoll.Rolls[i] = RollDF(dc)
		}
		// Roll the remaining dice normally
		for i := character.FudgeRoll; i < numDice; i++ {
			character.LastRoll.Rolls[i] = RollD10()
		}
		log.Printf("Fudged roll to: [%v]", character.LastRoll.Rolls)
		shuffle(character.LastRoll.Rolls)
		log.Println("Clearing fudge")
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
		character.LastRoll.Reason = ""
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

// Sets the minimum results for the next roll the invokee makes
func WouldYouKindly(minResults int, character *data.Character) string {
	character.FudgeRoll = minResults
	return fmt.Sprintf("Fudge set to minimum %d successes", minResults)
}

// go passes slices by reference so it should shuffle in place
func shuffle(nums []int) []int {
	for i := len(nums) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		nums[i], nums[j] = nums[j], nums[i]
	}
	return nums
}
