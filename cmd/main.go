package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"

	"github.com/keshon/discord-bot-boilerplate/example-bot/discord"
	"github.com/keshon/discord-bot-boilerplate/internal/config"
	"github.com/keshon/discord-bot-boilerplate/internal/db"
	"github.com/keshon/discord-bot-boilerplate/internal/manager"
	"github.com/keshon/discord-bot-boilerplate/internal/rest"
	"github.com/keshon/discord-bot-boilerplate/internal/version"
)

func main() {
	initLogger()

	config := initConfig()

	initDatabase()

	dg := initDiscordSession(config.DiscordBotToken)

	bots := initBots(dg)

	handleDiscordSession(dg)

	startRestServer(config, bots)

	slog.Infof("%v is now running. Press Ctrl+C to exit", version.AppName)

	waitForExitSignal()
}

func initLogger() {
	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
		f.SetTemplate("[{{datetime}}] [{{level}}] [{{caller}}]\t{{message}} {{data}} {{extra}}\n")
		f.ColorTheme = slog.ColorTheme
	})

	h1 := handler.MustFileHandler("./logs/all-levels.log", handler.WithLogLevels(slog.AllLevels))
	slog.PushHandler(h1)
}

func initConfig() *config.Config {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Fatal("Error loading config", err)
		os.Exit(1)
	}
	slog.Info("Config loaded:\n" + cfg.String())
	return cfg
}

func initDatabase() {
	_, err := db.InitDB("./dbase.db")
	if err != nil {
		slog.Fatal("Error initializing the database", err)
		os.Exit(1)
	}
}

func initDiscordSession(token string) *discordgo.Session {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		slog.Fatal("Error creating Discord session", err)
		os.Exit(1)
	}
	return dg
}

func initBots(session *discordgo.Session) map[string]*discord.BotInstance {
	bots := make(map[string]*discord.BotInstance)

	guildManager := manager.NewGuildManager(session, bots)
	guildManager.Start()

	guildIDs, err := db.GetAllGuildIDs()
	if err != nil {
		slog.Fatal("Error retrieving or creating guilds", err)
		os.Exit(1)
	}

	for _, guildID := range guildIDs {
		bots = startBotInstances(session, guildID, bots)
	}

	return bots
}

func handleDiscordSession(dg *discordgo.Session) {
	err := dg.Open()
	if err != nil {
		slog.Fatal("Error opening Discord session", err)
		os.Exit(1)
	}
	defer dg.Close()
}

func startBotInstances(session *discordgo.Session, guildID string, bots map[string]*discord.BotInstance) map[string]*discord.BotInstance {
	bots[guildID] = &discord.BotInstance{
		ExampleBot: discord.NewDiscord(session, guildID),
	}
	bots[guildID].ExampleBot.Start(guildID)

	return bots
}

func startRestServer(config *config.Config, bots map[string]*discord.BotInstance) {
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
		hostname := config.RestHostname

		if len(hostname) == 0 {
			hostname = "localhost:8080"
			slog.Warn("Hostname is empty, setting to default", hostname)
		}

		_, _, err := net.SplitHostPort(hostname)
		if err != nil {
			slog.Error("Could not detect host and port from sring", hostname)
			return
		}

		if err := router.Run(hostname); err != nil {
			slog.Fatal("Error starting REST API server", err)
			return
		}
		slog.Info("REST API server started on", hostname)
	}()
}

func waitForExitSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
