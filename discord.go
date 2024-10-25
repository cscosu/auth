package main

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Token         string
	GuildId       string
	AdminRoleId   string
	StudentRoleId string
	ClientId      string
	ClientSecret  string

	Session *discordgo.Session
}

var (
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
	}

	commandHandlers = map[string]func(b *DiscordBot, i *discordgo.InteractionCreate){
		"about": func(b *DiscordBot, i *discordgo.InteractionCreate) {
			if !requireAdmin(b, i) {
				return
			}
			b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This command has not been finished",
				},
			})
		},
		"archive-channel": func(b *DiscordBot, i *discordgo.InteractionCreate) {
			b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "This command has not been finished",
				},
			})
		},
	}
)

func isAdmin(b *DiscordBot, m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == b.AdminRoleId {
			return true
		}
	}
	return false
}

func requireAdmin(b *DiscordBot, i *discordgo.InteractionCreate) bool {
	if !isAdmin(b, i.Member) {
		b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "This command requires admin",
			},
		})
		return false
	}

	return true
}

func (b *DiscordBot) Connect() {
	if b.Token == "" {
		log.Fatalln("Missing token")
	}
	if b.AdminRoleId == "" {
		log.Fatalln("Missing admin role id")
	}
	if b.GuildId == "" {
		log.Fatalln("Missing guild id")
	}
	if b.ClientId == "" {
		log.Fatalln("Missing client id")
	}
	if b.ClientSecret == "" {
		log.Fatalln("Missing client secret")
	}

	s, err := discordgo.New("Bot " + b.Token)
	if err != nil {
		log.Fatalln("Failed to connect", err)
	}

	b.Session = s

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Logged in as", r.User.String())
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(b, i)
		}
	})

	err = s.Open()
	if err != nil {
		log.Fatalln("Failed to open session", err)
	}

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, b.GuildId, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}
}

func (b *DiscordBot) GiveStudentRole(discordId string) error {
	return b.Session.GuildMemberRoleAdd(b.GuildId, discordId, b.StudentRoleId)
}

func (b *DiscordBot) AddStudentToGuild(discordId string, accessToken string) error {
	return b.Session.GuildMemberAdd(b.GuildId, discordId, &discordgo.GuildMemberAddParams{
		AccessToken: accessToken,
		Roles:       []string{b.StudentRoleId},
	})
}

func (b *DiscordBot) RemoveStudentRole(discordId string) error {
	return b.Session.GuildMemberRoleRemove(b.GuildId, discordId, b.StudentRoleId)
}
