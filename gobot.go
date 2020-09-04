package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/discordgo"
)

var (
	// Token for the bot
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
}

func main() {
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running. Press CRTL-C to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	var diceResult string
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.HasPrefix(m.Content, "/") {
		diceResult = rollDice(trimSlash(m.Content))
		s.ChannelMessageSend(m.ChannelID, diceResult)
	}

}

func rollDice(c string) string {
	toRoll := strings.Split(c, ",")
	numDice, _ := strconv.Atoi(toRoll[0])
	DC, _ := strconv.Atoi(toRoll[1])
	var successes int
	diceResults := make([]int, numDice)
	for i := 0; i < numDice; i++ {
		diceResults[i] = rollD10()
		if diceResults[i] == 10 {
			successes += 2
		} else if diceResults[i] >= DC && diceResults[i] != 10 {
			successes++
		} else if diceResults[i] == 1 && diceResults[i] < DC {
			successes--
		}
	}
	if successes >= 1 {
		return fmt.Sprintf("```%d Successes\nRolled %v```", successes, diceResults)
	} else if successes == 0 {
		return fmt.Sprintf("```\nBotched\nRolled %v```", diceResults)
	}
	return fmt.Sprintf("```\nFailed\nRolled %v```", diceResults)
}

func trimSlash(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func rollD10() int {
	return 1 + rand.Intn(9)
}
