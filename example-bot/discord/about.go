package discord

import (
	"fmt"
	"os"
	"time"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"github.com/gookit/slog"

	"github.com/keshon/discord-bot-boilerplate/example-bot/utils"
	"github.com/keshon/discord-bot-boilerplate/internal/config"
	"github.com/keshon/discord-bot-boilerplate/internal/version"
)

// handleAboutCommand is a function to handle the about command in Discord.
//
// It takes a Discord session and a Discord message as parameters and does not return anything.
func (d *Discord) handleAboutCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	d.changeAvatar(s)

	config, err := config.NewConfig()
	if err != nil {
		slog.Fatal("Error loading config", err)
		return
	}

	hostname := os.Getenv("HOST")
	if hostname == "" {
		hostname = config.RestHostname
	}

	avatarUrl := fmt.Sprintf("%s://%s/avatar/random?%d", utils.InferProtocolByPort(hostname, 443), hostname, time.Now().UnixNano())
	slog.Info(avatarUrl)

	title := fmt.Sprintf("ℹ️ %v — About", version.AppName)
	content := fmt.Sprintf("**%v**\n\n%v", version.AppName, version.AppDescription)

	embedMsg := embed.NewEmbed().
		SetDescription(fmt.Sprintf("**%v**\n\n%v", title, content)).
		AddField("```"+version.BuildDate+"```", "Build date").
		AddField("```"+version.GoVersion+"```", "Go version").
		AddField("```Created by Innokentiy Sokolov```", "[Linkedin](https://www.linkedin.com/in/keshon), [GitHub](https://github.com/keshon), [Homepage](https://keshon.ru)").
		InlineAllFields().
		SetImage(avatarUrl).
		SetColor(0x9f00d4).SetFooter(version.AppFullName).MessageEmbed

	_, err = s.ChannelMessageSendEmbed(m.Message.ChannelID, embedMsg)
	if err != nil {
		slog.Fatal("Error sending embed message", err)
	}
}
