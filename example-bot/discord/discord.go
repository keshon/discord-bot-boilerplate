package discord

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/slog"
	"github.com/keshon/discord-bot-boilerplate/example-bot/utils"
	"github.com/keshon/discord-bot-boilerplate/internal/config"
)

type ExampleBot struct {
	Discord *Discord
}
type Discord struct {
	Session              *discordgo.Session
	GuildID              string
	IsInstanceActive     bool
	CommandPrefix        string
	LastChangeAvatarTime time.Time
	RateLimitDuration    time.Duration
}

// NewDiscord initializes a new Discord object with the given session and guild ID.
//
// Parameters:
// - session: a pointer to a discordgo.Session
// - guildID: a string representing the guild ID
// Returns a pointer to a Discord object.
func NewDiscord(session *discordgo.Session, guildID string) *Discord {
	config := loadConfig()

	return &Discord{
		Session:           session,
		IsInstanceActive:  true,
		CommandPrefix:     config.DiscordCommandPrefix,
		RateLimitDuration: time.Minute * 10,
	}
}

// loadConfig loads the configuration and returns a pointer to config.Config.
//
// No parameters.
// Returns a pointer to config.Config.
func loadConfig() *config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Fatal("Error loading config", err)
	}
	return cfg
}

// Start starts the Discord instance for the specified guild ID.
//
// guildID string
func (d *Discord) Start(guildID string) {
	slog.Info("Discord instance started for guild ID", guildID)
	d.Session.AddHandler(d.Commands)
	d.GuildID = guildID
}

// Commands handles the incoming Discord commands.
//
// Parameters:
//
//	s: the Discord session
//	m: the incoming Discord message
func (d *Discord) Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	slog.Warn(m.Message.Author.Username, ":", m.Message.Content)
	if m.GuildID != d.GuildID || !d.IsInstanceActive {
		return
	}

	command, parameter, err := parseCommand(m.Message.Content, d.CommandPrefix)
	if err != nil {
		return
	}

	switch getCanonicalCommand(command, [][]string{
		{"example", "e"},
		{"help", "h"},
		{"about", "a"},
	}) {
	case "example":
		slog.Error("in example")
		d.handleExampleCommand(s, m, parameter)
	case "help":
		d.handleHelpCommand(s, m)
	case "about":
		d.handleAboutCommand(s, m)
	}
}

// parseCommand parses the input based on the provided pattern
//
// input: the input string to be parsed
// pattern: the pattern to match at the beginning of the input
// string: the parsed command
// string: the parsed parameter
// error: an error if the pattern is not found or no command is found
func parseCommand(input, pattern string) (string, string, error) {
	input = strings.ToLower(input)
	pattern = strings.ToLower(pattern)

	if !strings.HasPrefix(input, pattern) {
		return "", "", fmt.Errorf("pattern not found")
	}

	input = input[len(pattern):]

	words := strings.Fields(input)
	if len(words) == 0 {
		return "", "", fmt.Errorf("no command found")
	}

	command := words[0]
	parameter := ""
	if len(words) > 1 {
		parameter = strings.Join(words[1:], " ")
		parameter = strings.TrimSpace(parameter)
	}
	return command, parameter, nil
}

// getCanonicalCommand finds the canonical command for the given alias.
//
// Parameters:
// - alias: a string representing the alias to be searched for.
// - commandAliases: a 2D slice containing the list of command aliases.
// Return type: string
func getCanonicalCommand(alias string, commandAliases [][]string) string {
	for _, aliases := range commandAliases {
		for _, command := range aliases {
			if command == alias {
				return aliases[0]
			}
		}
	}
	return ""
}

// changeAvatar changes the avatar of the Discord user.
//
// It takes a session as a parameter and does not return anything.
func (d *Discord) changeAvatar(s *discordgo.Session) {
	if time.Since(d.LastChangeAvatarTime) < d.RateLimitDuration {
		return
	}

	imgPath, err := utils.GetRandomImagePathFromPath("./assets/avatars")
	if err != nil {
		slog.Error("Error getting random image path:", err)
		return
	}

	avatar, err := utils.ReadFileToBase64(imgPath)
	if err != nil {
		slog.Error("Error reading file to base64:", err)
		return
	}

	_, err = s.UserUpdate("", avatar)
	if err != nil {
		slog.Error("Error updating user avatar:", err)
		return
	}

	d.LastChangeAvatarTime = time.Now()
}

// findUserVoiceState finds the voice state of the user in the given list of voice states.
//
// userID: The ID of the user to find the voice state for.
// voiceStates: The list of voice states to search through.
// *discordgo.VoiceState, bool: The found voice state and a boolean indicating if it was found.
func findUserVoiceState(userID string, voiceStates []*discordgo.VoiceState) (foundState *discordgo.VoiceState, found bool) {
	for _, state := range voiceStates {
		if state != nil && state.UserID == userID {
			foundState = state
			found = true
			return
		}
	}
	return
}
