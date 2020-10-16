package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
	"github.com/pkwiatek6/DiscrodBot/data"
)

/*
type rollHistory struct {
	rolls  []int
	dc     int
	reason string
}
*/
var (
	// Token for the bot
	Token string
	// LastRolls keeps track of the last player roll
	LastRolls map[string]data.RollHistory
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	LastRolls = make(map[string]data.RollHistory)
}

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CRTL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	var cmdGiven string

	//bot should never read it's own output
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Compare(strings.ToLower(m.Content), "flip a coin") == 0 {
		go flipCoin(m.ChannelID, m.Member.Nick, s)
		return
	}
	//Deals with commands, consumes the prefix(default / or !)
	IsCommand, err := regexp.MatchString("[/!].*", m.Content)
	if err != nil {
		log.Printf("%s; offending Command %s\n", err, m.Content)
	} else {
		cmdGiven = trimSlash(m.Content)
	}
	if IsCommand {
		//The Regex checks if you are rolling dice, I'm not using \s becuase it was giving me an error for some reason
		matched, err := regexp.MatchString("^[0-9]+,[0-9]+,?[a-zA-z\r\n\t\f\v]*", cmdGiven)
		if err != nil {
			log.Printf("%s; offending Command %s\n", err, m.Content)
			return
		}
		if matched {
			go rollDice(cmdGiven, m.Member.Nick, m.ChannelID, s)
		} else if strings.Compare(strings.ToLower(cmdGiven), "reroll") == 0 || strings.Compare(strings.ToLower(cmdGiven), "r") == 0 {
			go rerollDice(m.Member.Nick, m.ChannelID, s)
		} else if strings.Compare(strings.ToLower(cmdGiven), "schedule") == 0 {

		}
	}

}
func rollDice(c string, nick string, channel string, session *discordgo.Session) {
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
		diceResults[i] = rollD10()
	}
	successes := countSuc(diceResults, DC)
	//TODO: recode to user unique identifer instead of nickname so rolls aren;t lost
	LastRolls[nick] = data.RollHistory{Rolls: diceResults, DC: DC, Reason: reason}
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

func rerollDice(name string, channel string, session *discordgo.Session) {
	var oldResults = LastRolls[name].Rolls
	sort.Ints(oldResults)
	var tempDC = LastRolls[name].DC
	var failedRolls [3]int
	var newRolls [3]int
	var max = 3
	if len(oldResults) < 3 {
		max = len(oldResults)
	}
	for i := 0; i < max; i++ {
		if oldResults[i] < tempDC && oldResults[i] != 10 {
			failedRolls[i] = oldResults[i]
			newRolls[i] = rollD10()
			oldResults[i] = newRolls[i]
		}
	}
	successes := countSuc(oldResults, tempDC)
	LastRolls[name] = data.RollHistory{Rolls: oldResults, DC: tempDC, Reason: LastRolls[name].Reason}
	if successes >= 1 {
		toPost := fmt.Sprintf("```%s got %d Successes%s\nRerolls %v -> %v```", name, successes, LastRolls[name].Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)

	} else if successes == 0 {
		toPost := fmt.Sprintf("```%s Failed%s\nRerolls %v -> %v```", name, LastRolls[name].Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)
	} else {
		toPost := fmt.Sprintf("```%s got a Botch%s\nRerolls %v -> %v```", name, LastRolls[name].Reason, failedRolls, newRolls)
		session.ChannelMessageSend(channel, toPost)
	}
}

func trimSlash(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func rollD10() int {
	return 1 + rand.Intn(9)
}

func flipCoin(channel string, nick string, session *discordgo.Session) {
	if rand.Intn(2) == 0 {
		session.ChannelMessageSend(channel, fmt.Sprintf("```%s flipped a coin and it came up %s```", nick, "Heads"))
		return
	}
	session.ChannelMessageSend(channel, fmt.Sprintf("```%s flipped a coin and it came up %s```", nick, "Tails"))
	return
}
func countSuc(diceResults []int, DC int) int {
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
