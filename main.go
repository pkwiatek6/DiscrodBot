package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"
	"unicode/utf8"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bwmarrin/discordgo"
	"github.com/pkwiatek6/DiscrodBot/actions"
	"github.com/pkwiatek6/DiscrodBot/data"
)

var (
	// Token for the bot
	Token string
	// LastRolls keeps track of the last player roll
	LastRolls map[string]data.RollHistory
	//Client is the cpnnection to the database
	Client *mongo.Client
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	LastRolls = make(map[string]data.RollHistory)
}

func main() {
	//opens connection the the database to load in relevant data, also closes it when program finishes running
	Client = actions.ConnectDB()
	//defer still exists to close connections when program returns, though only when it error's out
	defer Client.Disconnect(context.Background())

	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}
	//can add more handlers based on the discord api, the function passed must always accept a Session and a discord event
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
	go func() {
		<-sc
		//closes conentions upon reciviing an interupt
		fmt.Println("\r- Interrupt recived, Closing Bot")
		Client.Disconnect(context.Background())
		discord.Close()
	}()
	//if something breaks when closing it's the defer
	defer discord.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	var cmdGiven string

	//bot should never read it's own output
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.Compare(strings.ToLower(m.Content), "flip a coin") == 0 {
		go actions.FlipCoin(m.ChannelID, m.Member.Nick, s)
		return
	}
	//Deals with commands, consumes the prefix(default / or !)
	IsCommand, err := regexp.MatchString("[/!].*", m.Content)
	if err != nil {
		log.Printf("%s; offending Command %s\n", err, m.Content)
	} else {
		cmdGiven = trimPrefix(m.Content)
	}
	if IsCommand {
		//The Regex checks if you are rolling dice, I'm not using \s becuase it was giving me an error saying it's not a vaild escape sequence
		matched, err := regexp.MatchString("^[0-9]+,[0-9]+,?[a-zA-z\r\n\t\f\v]*", cmdGiven)
		if err != nil {
			log.Printf("%s; offending Command %s\n", err, m.Content)
			return
		}
		//checks what the other commands are, this should probably be made into a router
		if matched {
			go actions.RollDice(cmdGiven, m.Member.Nick, m.ChannelID, s, &LastRolls)
		} else if strings.Compare(strings.ToLower(cmdGiven), "reroll") == 0 || strings.Compare(strings.ToLower(cmdGiven), "r") == 0 {
			go actions.RerollDice(m.Member.Nick, m.ChannelID, s, &LastRolls)
		} else if strings.Compare(strings.ToLower(cmdGiven), "schedule") == 0 {
			//TODO make sceduling command for next session
		}
	}

}

//trims the prefix
func trimPrefix(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
