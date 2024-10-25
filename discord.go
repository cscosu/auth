package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/bwmarrin/discordgo"
)

type DiscordBot struct {
	Token         string
	GuildId       string
	AdminRoleId   string
	StudentRoleId string
	ClientId      string
	ClientSecret  string
	Db            *sql.DB

	Session *discordgo.Session
}

var (
	commands = []*discordgo.ApplicationCommand{
		// All commands and options must have a description
		// Commands/options without description will fail the registration
		// of the command.
		{
			Name:        "about",
			Description: "Lookup information about a discord user who has linked their OSU account",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user",
					Description: "User",
					Required:    true,
				},
			},
		},
	}

	commandHandlers = map[string]func(b *DiscordBot, i *discordgo.InteractionCreate){
		"about": func(b *DiscordBot, i *discordgo.InteractionCreate) {
			if !b.requireAdmin(i) {
				return
			}
			options := i.ApplicationCommandData().Options
			user := options[0].UserValue(nil)

			row := b.Db.QueryRow(`SELECT name_num, display_name, last_login, student, alum, employee, faculty FROM users WHERE discord_id = ?`, user.ID)
			var (
				nameNum     string
				displayName string
				lastLogin   int
				student     bool
				alum        bool
				employee    bool
				faculty     bool
			)
			err := row.Scan(&nameNum, &displayName, &lastLogin, &student, &alum, &employee, &faculty)
			if err != nil {
				_ = b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "User has not linked their OSU account",
					},
				})
				return
			}
			content := fmt.Sprintf("**[%s (%s)](<https://www.osu.edu/search/?query=%s>)**\nLast login: <t:%d:f>\n",
				displayName,
				nameNum,
				url.QueryEscape(nameNum),
				lastLogin,
			)

			sep := ""
			if student {
				content += sep + "Student"
				sep = ", "
			}
			if alum {
				content += sep + "Alum"
				sep = ", "
			}
			if employee {
				content += sep + "Employee"
				sep = ", "
			}
			if faculty {
				content += sep + "Faculty"
			}

			_ = b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
	}
)

func (b *DiscordBot) isAdmin(m *discordgo.Member) bool {
	for _, role := range m.Roles {
		if role == b.AdminRoleId {
			return true
		}
	}
	return false
}

func (b *DiscordBot) requireAdmin(i *discordgo.InteractionCreate) bool {
	if !b.isAdmin(i.Member) {
		_ = b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
