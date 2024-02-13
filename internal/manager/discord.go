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
	prefix = os.Getenv("COMMAND_PREFIX")
)

// GuildManager manages guild-related functionality.
type GuildManager struct {
	Session      *discordgo.Session
	BotInstances map[string]*discord.BotInstance
	prefix       string
}

// NewGuildManager creates a new instance of GuildManager.
func NewGuildManager(session *discordgo.Session, botInstances map[string]*discord.BotInstance) *GuildManager {
	config, err := config.NewConfig()
	if err != nil {
		slog.Fatalf("Error loading config:", err)
	}

	return &GuildManager{
		Session:      session,
		BotInstances: botInstances,
		prefix:       config.DiscordCommandPrefix,
	}
}

// Start starts the GuildManager instance.
func (gm *GuildManager) Start() {
	slog.Info("Guild manager started")
	gm.Session.AddHandler(gm.Commands)
}

// Commands handles incoming Discord commands.
func (gm *GuildManager) Commands(s *discordgo.Session, m *discordgo.MessageCreate) {
	command, _, err := parseCommand(m.Message.Content, gm.prefix)
	if err != nil {
		slog.Error(err)
		return
	}

	// Check if the command is "roll" or "dice"
	if command == "example" || command == "help" || command == "about" {
		guildID := m.GuildID
		exists, err := db.DoesGuildExist(guildID)
		if err != nil {
			slog.Errorf("Error checking if guild is registered: %v", err)
			return
		}

		if !exists {
			gm.Session.ChannelMessageSend(m.Message.ChannelID, "Guild must be registered first. Use `"+gm.prefix+"register` command.")
			return
		}
	}

	switch command {
	case "register":
		gm.handleRegisterCommand(s, m)
	case "unregister":
		gm.handleUnregisterCommand(s, m)
	default:
		// log.Println("Unknown command")
	}
}

// handleRegisterCommand handles the registration of a guild.
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

// handleUnregisterCommand handles the unregistration of a guild.
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

// setupBotInstance sets up a new BotInstance for a guild.
func (gm *GuildManager) setupBotInstance(botInstances map[string]*discord.BotInstance, session *discordgo.Session, guildID string) {
	botInstances[guildID] = &discord.BotInstance{
		ExampleBot: discord.NewDiscord(session, guildID),
	}
	botInstances[guildID].ExampleBot.Start(guildID)
}

// removeBotInstance removes a BotInstance for a guild.
func (gm *GuildManager) removeBotInstance(guildID string) {
	instance, ok := gm.BotInstances[guildID]
	if !ok {
		return // Guild instance not found, nothing to do
	}

	instance.ExampleBot.InstanceActive = false

	delete(gm.BotInstances, guildID)
}
