package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/bwmarrin/discordgo"
	"github.com/pkwiatek6/DiscrodBot/actions"
	"github.com/pkwiatek6/DiscrodBot/data"
)

var (
	discord *discordgo.Session
	// Token for the bot
	Token          = flag.String("t", "", "Bot acess token")
	GuildID        = flag.String("GID", "", "Test Guild ID. IF not passed - bot registers commands globally")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

func init() { flag.Parse() }

func init() {
	var err error
	discord, err = discordgo.New("Bot " + *Token)
	if err != nil {
		log.Fatalln("Error creating Discord session: ", err)
	}
	log.Println("Connection to Discord established")
}

var (
	adminMemeberPermissions int64 = discordgo.PermissionAdministrator

	Characters map[string]*data.Character
	//Client is the connection to the database
	Client *mongo.Client

	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "coin",
			Description: "Flips a coin",
		},
		{
			Name:        "reroll",
			Description: "Re-rolls lowest 3 dice that are lower than the DC by using willpower.",
		},
		{
			Name:                     "wyk",
			Description:              "Sets the minimum number of success you will get on your next roll",
			DefaultMemberPermissions: &adminMemeberPermissions,
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "min-results",
					Description: "Minimum number of results wanted",
					Required:    true,
				},
			},
		},
		{
			Name:        "saveall",
			Description: "Saves all users to database, debug tool do not use unless Peter told you too",
		},
		{
			Name:        "roll",
			Description: "Rolls a dice pool against a dc with and optional reason",
			Options: []*discordgo.ApplicationCommandOption{

				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "dice-pool",
					Description: "Number of dice to roll",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "dc",
					Description: "DC of the check",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "action",
					Description: "What action you are trying to do.",
					Required:    false,
				},
			},
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"coin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: actions.FlipCoin(i.Member.Nick),
				},
			})
		},
		"reroll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: actions.RerollDice(Characters[i.Member.User.ID]),
				},
			})
		},
		// To be implemented when permissions are added in discordgo
		"wyk": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var minResults = int(i.ApplicationCommandData().Options[0].IntValue())
			message := actions.WouldYouKindly(minResults, Characters[i.Member.User.ID])
			discord.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: message,
				},
			})
		},
		"saveall": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			for key := range Characters {
				err := actions.SaveCharacter(*Characters[key], Client)
				if err != nil {
					log.Println(err)
				}
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "All data saved :)",
				},
			})
		},
		"roll": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if Characters[i.Member.User.ID] == nil {
				Characters[i.Member.User.ID] = new(data.Character)
				Characters[i.Member.User.ID].User = i.Member.User.ID
				Characters[i.Member.User.ID].Name = i.Member.Nick
				Characters[i.Member.User.ID].DiscordUser = i.Member.User.String()
				Characters[i.Member.User.ID].LastRoll = *new(data.RollHistory)
				actions.SaveCharacter(*Characters[i.Member.User.ID], Client)

			}
			var dicepool = int(i.ApplicationCommandData().Options[0].IntValue())
			var dc = int(i.ApplicationCommandData().Options[1].IntValue())
			var msg string
			if len(i.ApplicationCommandData().Options) == 3 {
				var reason = i.ApplicationCommandData().Options[2].StringValue()
				msg = actions.RollDiceCommand(dicepool, dc, reason, Characters[i.Member.User.ID])
			} else {
				msg = actions.RollDiceCommand(dicepool, dc, "", Characters[i.Member.User.ID])
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: msg,
				},
			})
		},
	}
	globalCommands = []*discordgo.ApplicationCommand{}
)

func init() {
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			handler(s, i)
		}

	})
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func init() {
	var err error
	Characters = make(map[string]*data.Character)
	Client, err = actions.ConnectDB()
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Connection to Database established")
	Characters, err = actions.LoadAllCharacters(Client)
	if err != nil {
		log.Fatalln("Error loading all characters")
	}
	log.Println("All Characters loaded")

}

func main() {
	err := discord.Open()
	if err != nil {
		log.Fatalln("Error opening connection: ", err)
	}
	log.Println("Connection to Discord opened")

	fmt.Println("Bot is now running. Press CRTL-C or send SIGINT or SIGTERM to exit")

	for _, v := range commands {
		gCMD, err := discord.ApplicationCommandCreate(discord.State.User.ID, *GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		globalCommands = append(globalCommands, gCMD)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
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
		if *RemoveCommands {
			for _, v := range globalCommands {
				err := discord.ApplicationCommandDelete(discord.State.User.ID, *GuildID, v.ID)
				if err != nil {
					log.Printf("Cannot delete '%v' command: %v", v.ID, err)
				}
			}
		}
		discord.Close()
	}()

}
