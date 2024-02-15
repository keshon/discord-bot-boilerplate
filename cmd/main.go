package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"

	"github.com/keshon/discord-bot-boilerplate/internal/botsdef"
	"github.com/keshon/discord-bot-boilerplate/internal/config"
	"github.com/keshon/discord-bot-boilerplate/internal/db"
	"github.com/keshon/discord-bot-boilerplate/internal/manager"
	"github.com/keshon/discord-bot-boilerplate/internal/rest"
	"github.com/keshon/discord-bot-boilerplate/internal/version"

	greetings "github.com/keshon/discord-bot-boilerplate/bot-greetings/discord"
	helloworld "github.com/keshon/discord-bot-boilerplate/bot-helloworld/discord"
)

// main is the entry point of the program.
//
// There are no parameters.
// There is no return type.
func main() {
	initLogger()
	config := loadConfig()
	initDatabase()
	discordSession := createDiscordSession(config.DiscordBotToken)
	bots := startBotHandlers(discordSession)
	handleDiscordSession(discordSession)
	startRestServer(config, bots)
	slog.Infof("%v is now running. Press Ctrl+C to exit", version.AppName)
	waitForExitSignal()
}

// initLogger initializes the logger with color theme and file handler.
//
// No parameters.
// No return type.
func initLogger() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		if textFormatter, ok := logger.Formatter.(*slog.TextFormatter); ok {
			textFormatter.EnableColor = true
			textFormatter.SetTemplate("[{{datetime}}] [{{level}}] [{{caller}}]\t{{message}} {{data}} {{extra}}\n")
			textFormatter.ColorTheme = slog.ColorTheme
		} else {
			slog.Error("Error: Text formatter is not a *slog.TextFormatter")
		}
	})
	if fileHandler, err := handler.NewFileHandler("./logs/all-levels.log", handler.WithLogLevels(slog.AllLevels)); err != nil {
		slog.Error("Error creating file handler:", err)
	} else {
		slog.PushHandler(fileHandler)
	}
}

// loadConfig loads the configuration and returns a pointer to config.Config.
//
// No parameters.
// Returns *config.Config.
func loadConfig() *config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Fatal("Error loading config", err)
		os.Exit(1)
	}
	slog.Info("Config loaded:\n" + cfg.String())
	return cfg
}

// initDatabase initializes the database.
//
// No parameters.
// No return type.
func initDatabase() {
	_, err := db.InitDB("./database.db")
	if err != nil {
		slog.Fatal("Error initializing the database", err)
		os.Exit(1)
	}
}

// createDiscordSession creates a new Discord session using the provided token.
//
// It takes a string token as a parameter and returns a *discordgo.Session.
func createDiscordSession(token string) *discordgo.Session {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("Error creating Discord session", err)
	}
	return session
}

// startBotHandlers initializes and starts Discord bot instances for each guild using the provided session.
//
// session: a pointer to a discordgo.Session
// map[string]*discord.BotInstance: a map of guild IDs to their corresponding BotInstance pointers
func startBotHandlers(session *discordgo.Session) []map[string]map[string]botsdef.Discord {
	bot1 := make(map[string]map[string]botsdef.Discord)
	bot2 := make(map[string]map[string]botsdef.Discord)

	guildIDs, err := db.GetAllGuildIDs()
	if err != nil {
		log.Fatal("Error retrieving or creating guilds", err)
	}

	for _, guildID := range guildIDs {
		bot1[guildID] = make(map[string]botsdef.Discord)
		bot1[guildID]["bot1"] = greetings.NewDiscord(session)
		bot1[guildID]["bot1"].Start(guildID)

		bot2[guildID] = make(map[string]botsdef.Discord)
		bot2[guildID]["bot2"] = helloworld.NewDiscord(session)
		bot2[guildID]["bot2"].Start(guildID)
	}

	var bots []map[string]map[string]botsdef.Discord
	for id, bot := range bot1 {
		bots = append(bots, map[string]map[string]botsdef.Discord{id: bot})
	}
	for id, bot := range bot2 {
		bots = append(bots, map[string]map[string]botsdef.Discord{id: bot})
	}

	guildManager := manager.NewGuildManager(session, bots)
	guildManager.Start()

	return bots
}

// handleDiscordSession is a Go function that opens a Discord session and handles any errors.
//
// It takes a parameter discordSession of type *discordgo.Session.
// It does not return any value.
func handleDiscordSession(discordSession *discordgo.Session) {
	if err := discordSession.Open(); err != nil {
		slog.Fatal("Error opening Discord session", err)
		os.Exit(1)
	}
	defer discordSession.Close()
}

// startRestServer starts the REST server based on the given configuration and bot instances.
//
// It takes a config.Config pointer and a map of string to *discord.BotInstance as parameters.
func startRestServer(config *config.Config, bots []map[string]map[string]botsdef.Discord) {
	if !config.RestEnabled {
		return
	}
	if config.RestGinRelease {
		gin.SetMode("release")
	}
	router := gin.Default()
	restAPI := rest.NewRest(bots)
	restAPI.Start(router)
	go func() {
		if len(config.RestHostname) == 0 {
			config.RestHostname = "localhost:8080"
			slog.Warn("Hostname is empty, setting to default:", config.RestHostname)
		}
		if err := router.Run(config.RestHostname); err != nil {
			slog.Fatal("Error starting REST API server:", err)
			return
		}
		slog.Info("REST API server started on", config.RestHostname)
	}()
}

// waitForExitSignal waits for an exit signal and handles it.
//
// No parameters.
// No return types.
func waitForExitSignal() {
	exitSignalChannel := make(chan os.Signal, 1)
	signal.Notify(exitSignalChannel, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-exitSignalChannel
}
