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

// handleAboutCommand handles the about command for Discord.
func (d *Discord) handleAboutCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	d.changeAvatar(s)

	config, err := config.NewConfig()
	if err != nil {
		slog.Fatalf("Error loading config: %v", err)
	}

	var hostname string
	if os.Getenv("HOST") == "" {
		hostname = config.RestHostname
	} else {
		hostname = os.Getenv("HOST") // from docker environment
	}

	avatarUrl := utils.InferProtocolByPort(hostname, 443) + hostname + "/avatar/random?" + fmt.Sprint(time.Now().UnixNano())
	slog.Info(avatarUrl)

	title := "About Dice Roller"                                                        //getRandomAboutTitlePhrase()
	content := "example-bot roller is a simple implementation of rolling dices for D&D" //getRandomAboutDescriptionPhrase()

	embedStr := fmt.Sprintf("**%v**\n\n%v", title, content)

	embedMsg := embed.NewEmbed().
		SetDescription(embedStr).
		AddField("```"+version.BuildDate+"```", "Build date").
		AddField("```"+version.GoVersion+"```", "Go version").
		AddField("```Created by Innokentiy Sokolov```", "[Linkedin](https://www.linkedin.com/in/keshon), [GitHub](https://github.com/keshon), [Homepage](https://keshon.ru)").
		InlineAllFields().
		SetImage(avatarUrl).
		SetColor(0x9f00d4).SetFooter(version.AppFullName).MessageEmbed

	s.ChannelMessageSendEmbed(m.Message.ChannelID, embedMsg)
}
