package manager

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/slog"

	"github.com/keshon/discord-bot-boilerplate/example-bot/discord"
	example_bot "github.com/keshon/discord-bot-boilerplate/example-bot/discord"
	"github.com/keshon/discord-bot-boilerplate/internal/config"
	"github.com/keshon/discord-bot-boilerplate/internal/db"
)

type GuildManager struct {
	Session       *discordgo.Session
	ExampleBots   map[string]*example_bot.ExampleBot
	commandPrefix string
}

// NewGuildManager creates a new GuildManager with the given discord session and bot instances.
//
// Parameters:
// - session: *discordgo.Session
// - botInstances: map[string]*discord.BotInstance
// Return type: *GuildManager
func NewGuildManager(session *discordgo.Session, examplebots map[string]*example_bot.ExampleBot) *GuildManager {
	config, err := config.NewConfig()
	if err != nil {
		slog.Fatalf("Error loading config:", err)
	}

	return &GuildManager{
		Session:       session,
		ExampleBots:   examplebots,
		commandPrefix: config.DiscordCommandPrefix,
	}
}

// Start starts the GuildManager.
func (gm *GuildManager) Start() {
	slog.Info("Guild manager started")
	gm.Session.AddHandler(gm.Commands)
}

// Commands handles the commands received in a Discord session message.
//
// Parameters:
//   - s: a pointer to the Discord session
//   - m: a pointer to the Discord message received
func (gm *GuildManager) Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	command, _, err := parseCommand(m.Message.Content, gm.commandPrefix)
	if err != nil {
		slog.Error(err)
		return
	}

	switch command {
	case "example", "help", "about":
		guildID := m.GuildID
		exists, err := db.DoesGuildExist(guildID)
		if err != nil {
			slog.Errorf("Error checking if guild is registered: %v", err)
			return
		}

		if !exists {
			gm.Session.ChannelMessageSend(m.Message.ChannelID, "Guild must be registered first.\nUse `"+gm.commandPrefix+"register` command.")
			return
		}
	}

	switch command {
	case "register":
		gm.handleRegisterCommand(s, m)
	case "unregister":
		gm.handleUnregisterCommand(s, m)
	}
}

// handleRegisterCommand handles the registration command for the GuildManager.
//
// Parameters:
// - s: The discord session.
// - m: The message create event.
func (gm *GuildManager) handleRegisterCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	channelID := m.Message.ChannelID
	guildID := m.GuildID

	exists, err := db.DoesGuildExist(guildID)
	if err != nil {
		slog.Errorf("Error checking if guild is registered: %v", err)
		return
	}

	if exists {
		gm.Session.ChannelMessageSend(channelID, "Guild is already registered")
		return
	}

	guild := db.Guild{ID: guildID, Name: ""}
	err = db.CreateGuild(guild)
	if err != nil {
		slog.Errorf("Error registering guild: %v", err)
		gm.Session.ChannelMessageSend(channelID, "Error registering guild")
		return
	}

	gm.setupBotInstance(gm.ExampleBots, s, guildID)
	gm.Session.ChannelMessageSend(channelID, "Guild registered successfully")
}

// handleUnregisterCommand handles the unregister command for the GuildManager.
//
// Parameters:
// - s: the discordgo Session
// - m: the discordgo MessageCreate
// Return type: none
func (gm *GuildManager) handleUnregisterCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	channelID := m.Message.ChannelID
	guildID := m.GuildID

	exists, err := db.DoesGuildExist(guildID)
	if err != nil {
		slog.Errorf("Error checking if guild is registered: %v", err)
		return
	}

	if !exists {
		gm.Session.ChannelMessageSend(channelID, "Guild is not registered")
		return
	}

	err = db.DeleteGuild(guildID)
	if err != nil {
		slog.Errorf("Error unregistering guild: %v", err)
		gm.Session.ChannelMessageSend(channelID, "Error unregistering guild")
		return
	}

	gm.removeBotInstance(guildID)
	gm.Session.ChannelMessageSend(channelID, "Guild unregistered successfully")
}

// setupBotInstance sets up a bot instance for the given guild.
//
// Parameters:
// - botInstances: a map of bot instances
// - session: a Discord session
// - guildID: the ID of the guild
// Return type: none
func (gm *GuildManager) setupBotInstance(exampleBots map[string]*example_bot.ExampleBot, session *discordgo.Session, guildID string) {
	exampleBots[guildID] = &example_bot.ExampleBot{
		Discord: discord.NewDiscord(session, guildID),
	}
	exampleBots[guildID].Discord.Start(guildID)
}

// removeBotInstance removes the bot instance for the given guild ID.
//
// guildID string
func (gm *GuildManager) removeBotInstance(guildID string) {
	instance, ok := gm.ExampleBots[guildID]
	if !ok {
		return // Guild instance not found, nothing to do
	}

	instance.Discord.IsInstanceActive = false

	delete(gm.ExampleBots, guildID)
}

// parseCommand parses the input based on the given pattern and returns the command and parameter.
//
// It takes two string parameters and returns two strings and an error.
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
