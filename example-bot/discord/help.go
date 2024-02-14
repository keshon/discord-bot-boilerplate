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

// handleHelpCommand handles the help command for the Discord bot.
//
// Takes in a session and a message create, and does not return any value.
func (d *Discord) handleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
	d.changeAvatar(s)

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Fatal("Error loading config", err)
		return
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = cfg.RestHostname
	}

	avatarURL := utils.InferProtocolByPort(host, 443) + host + "/avatar/random?" + fmt.Sprint(time.Now().UnixNano())
	prefix := d.commandPrefix
	example := fmt.Sprintf("`%vexample` - prints Hello World\n", prefix)
	help := fmt.Sprintf("**Show help**: `%vhelp` \nAliases: `%vh`\n", prefix, prefix)
	about := fmt.Sprintf("**Show version**: `%vabout` \nAliases: `%va`\n", prefix, prefix)
	register := fmt.Sprintf("**Enable commands listening**: `%vregister`\n", prefix)
	unregister := fmt.Sprintf("**Disable commands listening**: `%vunregister`", prefix)

	embedMsg := embed.NewEmbed().
		SetTitle(fmt.Sprintf("ℹ️ %v — Command Usage", version.AppName)).
		SetDescription("Some commands are aliased for shortness.\n").
		AddField("", "*Example*\n"+example).
		AddField("", "").
		AddField("", "*General*\n"+help+about).
		AddField("", "").
		AddField("", "*Administration*\n"+register+unregister).
		SetThumbnail(avatarURL).
		SetColor(0x9f00d4).
		SetFooter(version.AppFullName).
		MessageEmbed

	_, err = s.ChannelMessageSendEmbed(m.Message.ChannelID, embedMsg)
	if err != nil {
		slog.Fatal("Error sending embed message", err)
	}
}
