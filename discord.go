package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

func getenv(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalln(name, "environment variable not defined")
	}
	return value
}

var (
	token       string
	guildId     string
	adminRoleId string

	commands = []*discordgo.ApplicationCommand{
		// All commands and options must have a description
		// Commands/options without description will fail the registration
		// of the command.
		{
			Name: "about",

			Description: "admin only - lookup all data about a user",
		},
		{
			Name: "archive-channel",

			Description: "admin only - archive a channel",
		},
		{
			Name: "help",

			Description: "list commands a user can use",
		},
	}

	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"about": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if !requireAdmin(s, i) {
				return
			}
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This command has not been finished",
				},
			})
		},
		"archive-channel": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This command has not been finished",
				},
			})
		},
		"help": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			// s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			// 	Type: discordgo.InteractionResponseChannelMessageWithSource,
			// 	Data: &discordgo.InteractionResponseData{
			// 		if isAdmin(s, i.user) {
			// 			Content: "This command has not been finished",
			// 		} else {
			// 			Content:
			// 		}
			// 	},
			// })
		},
	}
)

func isAdmin(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == adminRoleId {
			return true
		}
	}
	return false
}

func requireAdmin(s *discordgo.Session, i *discordgo.InteractionCreate) bool {
	if !isAdmin(i.Member) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command requires admin",
			},
		})
		return false
	}

	return true
}

func discordBot() {
	token = getenv("DISCORD_BOT_TOKEN")
	guildId = getenv("DISCORD_GUILD_ID")
	adminRoleId = getenv("DISCORD_ADMIN_ROLE_ID")

	s, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Panicln("Failed to connect", err)
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Logged in as", r.User.String())
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	err = s.Open()
	if err != nil {
		log.Panicln("Failed to open session", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}
