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
	// Characters keeps track of players
	Characters map[string]*data.Character
	//Client is the cpnnection to the database
	Client *mongo.Client
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
	rand.Seed(time.Now().UnixNano())
	Characters = make(map[string]*data.Character)
}

func main() {
	//opens connection the the database to load in relevant data, also closes it when program finishes running
	var err error
	Client, err = actions.ConnectDB()
	if err != nil {
		log.Println(err)
		return
	}
	loadCharacters()
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
	defer func() {
		<-sc
		//closes conentions upon reciviing an interupt
		fmt.Println("\r- Interrupt recived, Closing Bot")
		Client.Disconnect(context.Background())
		discord.Close()
	}()
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
		Characters[m.Author.ID].User = m.Author.ID
		log.Printf("WTF:%v\n", Characters[m.Author.ID].User)
		//checks what the other commands are, this should probably be made into a router
		// m refferences the message
		//TODO change m to a better variable name
		if matched {
			go actions.RollDice(cmdGiven, m.ChannelID, s, Characters[m.Author.ID])
		} else if strings.Compare(strings.ToLower(cmdGiven), "reroll") == 0 || strings.Compare(strings.ToLower(cmdGiven), "r") == 0 {
			go actions.RerollDice(Characters[m.Author.ID], m.ChannelID, s)
		} else if strings.Compare(strings.ToLower(cmdGiven), "schedule") == 0 {
			//TODO make sceduling command for next session
		} else if strings.Compare(strings.ToLower(cmdGiven), "testing") == 0 {
			err := actions.SaveCharacter(*Characters[m.Author.ID], Client)
			if err != nil {
				log.Println(err)
			}
		}
	}

}

//trims the prefix
func trimPrefix(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

//loadCharacters caches all the characters from the DB
func loadCharacters() {
	Client.Database(actions.Database).Collection(actions.Collection)
}
