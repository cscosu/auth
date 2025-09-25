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
	AlumniRoleId  string
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
			discordId := user.ID

			row := b.Db.QueryRow(`SELECT name_num, display_name, last_seen_timestamp, student, alum, employee, faculty FROM users WHERE discord_id = ?`, discordId)
			var (
				nameNum     string
				displayName string
				lastSeen    int
				student     bool
				alum        bool
				employee    bool
				faculty     bool
			)
			err := row.Scan(&nameNum, &displayName, &lastSeen, &student, &alum, &employee, &faculty)
			if err != nil {
				log.Println("/about command: discordId =", discordId, err)
				_ = b.Session.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "User has not linked their OSU account",
						Flags:   discordgo.MessageFlagsEphemeral,
					},
				})
				return
			}
			content := fmt.Sprintf("**[%s (%s)](<https://www.osu.edu/search/?query=%s>)**\nLast seen: <t:%d:f>\n",
				displayName,
				nameNum,
				url.QueryEscape(nameNum),
				lastSeen,
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
					Flags:   discordgo.MessageFlagsEphemeral,
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
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return false
	}

	return true
}

func (b *DiscordBot) Connect() {
	if b.Token == "" {
		log.Println("Failed to connect to discord: Missing token")
		return
	}
	if b.AdminRoleId == "" {
		log.Println("Failed to connect to discord: Missing admin role id")
		return
	}
	if b.GuildId == "" {
		log.Println("Failed to connect to discord: Missing guild id")
		return
	}
	if b.ClientId == "" {
		log.Println("Failed to connect to discord: Missing client id")
		return
	}
	if b.ClientSecret == "" {
		log.Println("Failed to connect to discord: Missing client secret")
		return
	}

	s, err := discordgo.New("Bot " + b.Token)
	if err != nil {
		log.Println("Failed to connect to discord:", err)
		return
	}

	b.Session = s

	s.Identify.Intents = discordgo.IntentGuildMembers

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Println("Logged in as", r.User.String())
	})

	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(b, i)
		}
	})

	s.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		row := b.Db.QueryRow("SELECT buck_id FROM users WHERE discord_id = ?", m.User.ID)
		if row != nil {
			b.GiveStudentRole(m.User.ID)
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
	if b.Session == nil {
		return fmt.Errorf("discord bot not connected")
	}
	return b.Session.GuildMemberRoleAdd(b.GuildId, discordId, b.StudentRoleId)
}

func (b *DiscordBot) AddStudentToGuild(discordId string, accessToken string) error {
	if b.Session == nil {
		return fmt.Errorf("discord bot not connected")
	}
	return b.Session.GuildMemberAdd(b.GuildId, discordId, &discordgo.GuildMemberAddParams{
		AccessToken: accessToken,
		Roles:       []string{b.StudentRoleId},
	})
}

func (b *DiscordBot) RemoveStudentRole(discordId string) error {
	if b.Session == nil {
		return fmt.Errorf("discord bot not connected")
	}
	return b.Session.GuildMemberRoleRemove(b.GuildId, discordId, b.StudentRoleId)
}

// Turns a student into an alum (removes student role and adds alumni role)
func (b *DiscordBot) Alumnify(discordId string) error {
	if b.Session == nil {
		return fmt.Errorf("discord bot not connected")
	}

	err := b.Session.GuildMemberRoleRemove(b.GuildId, discordId, b.StudentRoleId)
	if err != nil {
		return fmt.Errorf("failed to remove student role from alum: %s")
	}

	err = b.Session.GuildMemberRoleAdd(b.GuildId, discordId, b.AlumniRoleId)
	if err != nil {
		return fmt.Errorf("failed to add student role to alum: %s")
	}

	b.Db.Exec("UPDATE Users set alum=1, student=0 WHERE discord_id=?", discordId)

	return nil
}
