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
	log.Println("Connection to Database established")
	discord, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}
	log.Println("Connection to Discord established")
	//can add more handlers based on the discord api, the function passed must always accept a Session and a discord event
	discord.AddHandler(messageCreate)

	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	err = discord.Open()
	if err != nil {
		log.Println("Error opening connection: ", err)
		return
	}
	log.Println("Connection to Discord opened")

	/*Characters, err = actions.LoadAllCharacters(Client)
	if err != nil {
		log.Println("Error loading all characters")
	}
	log.Println("All Characters loaded")
	*/
	fmt.Println("Bot is now running. Press CRTL-C or send SIGINT or SIGTERM to exit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	defer func() {
		<-sc
		//closes conentions upon reciviing an interupt
		log.Println("\r- Interrupt recived, Closing Bot")
		err = actions.SaveAllCharacters(Characters, Client)
		if err != nil {
			log.Println("Could not save characters: ", err)
		}
		err = Client.Disconnect(context.Background())
		if err != nil {
			log.Println("Could not close connection to database", err)
		}
		discord.Close()
	}()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	var cmdGiven string
	//bot should never read it's own output
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Author.Bot {
		return
	}
	//if there is no character it makes one or loads one in if it can
	if Characters[m.Author.ID] == nil {
		var err error
		Characters[m.Author.ID], err = actions.LoadCharacter(m.Member.Nick, m.Author.ID, Client)
		if err != nil {
			Characters[m.Author.ID] = new(data.Character)
			Characters[m.Author.ID].User = m.Author.ID
			Characters[m.Author.ID].Name = m.Member.Nick
			Characters[m.Author.ID].DiscordUser = m.Author.String()
			Characters[m.Author.ID].LastRoll = *new(data.RollHistory)
			err = actions.SaveCharacter(*Characters[m.Author.ID], Client)
			if err != nil {
				log.Println(err)
				return
			}
		}
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
		// m refferences the message
		//TODO change m to a better variable name
		if matched {
			go actions.RollDice(cmdGiven, m.ChannelID, s, Characters[m.Author.ID])
		} else if strings.Compare(strings.ToLower(cmdGiven), "reroll") == 0 || strings.Compare(strings.ToLower(cmdGiven), "r") == 0 {
			go actions.RerollDice(Characters[m.Author.ID], m.ChannelID, s)
		} else if strings.Compare(strings.ToLower(cmdGiven), "schedule") == 0 {
			//TODO make sceduling command for next session
		} else if strings.Compare(strings.ToLower(cmdGiven), "testsave") == 0 {
			//testing forcibly saves the character of the person who called it
			err := actions.SaveCharacter(*Characters[m.Author.ID], Client)
			if err != nil {
				log.Println(err)
			}
		}
		//TODO show and set commands for showing and setting data
	}

}

//trims the prefix
func trimPrefix(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
