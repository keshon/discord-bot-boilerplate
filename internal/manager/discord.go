package manager

import (
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/gookit/slog"

	"github.com/keshon/discord-bot-boilerplate/example-bot/discord"
	"github.com/keshon/discord-bot-boilerplate/internal/config"
	"github.com/keshon/discord-bot-boilerplate/internal/db"
)

var (
	commandPrefix = os.Getenv("COMMAND_PREFIX")
)

type GuildManager struct {
	Session       *discordgo.Session
	BotInstances  map[string]*discord.BotInstance
	commandPrefix string
}

// NewGuildManager creates a new GuildManager with the given discord session and bot instances.
//
// Parameters:
// - session: *discordgo.Session
// - botInstances: map[string]*discord.BotInstance
// Return type: *GuildManager
func NewGuildManager(session *discordgo.Session, botInstances map[string]*discord.BotInstance) *GuildManager {
	config, err := config.NewConfig()
	if err != nil {
		slog.Fatalf("Error loading config:", err)
	}

	return &GuildManager{
		Session:       session,
		BotInstances:  botInstances,
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
			gm.Session.ChannelMessageSend(m.Message.ChannelID, "Guild must be registered first. Use `"+gm.commandPrefix+"register` command.")
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

	gm.setupBotInstance(gm.BotInstances, s, guildID)
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
func (gm *GuildManager) setupBotInstance(botInstances map[string]*discord.BotInstance, session *discordgo.Session, guildID string) {
	botInstances[guildID] = &discord.BotInstance{
		ExampleBot: discord.NewDiscord(session, guildID),
	}
	botInstances[guildID].ExampleBot.Start(guildID)
}

// removeBotInstance removes the bot instance for the given guild ID.
//
// guildID string
func (gm *GuildManager) removeBotInstance(guildID string) {
	instance, ok := gm.BotInstances[guildID]
	if !ok {
		return // Guild instance not found, nothing to do
	}

	instance.ExampleBot.InstanceActive = false

	delete(gm.BotInstances, guildID)
}
