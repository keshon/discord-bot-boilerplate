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

// handleHelpCommand handles the help command for Discord.
func (d *Discord) handleHelpCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	example := fmt.Sprintf("`%vexample` - prints Hello World\n", d.prefix)
	help := fmt.Sprintf("**Show help**: `%vhelp` \nAliases: `%vh`\n", d.prefix, d.prefix)
	about := fmt.Sprintf("**Show version**: `%vabout` \nAliases: `%va`\n", d.prefix, d.prefix)
	register := fmt.Sprintf("**Enable commands listening**: `%vregister`\n", d.prefix)
	unregister := fmt.Sprintf("**Disable commands listening**: `%vunregister`", d.prefix)

	embedMsg := embed.NewEmbed().
		SetTitle(fmt.Sprintf("ℹ️ %v — Command Usage", version.AppName)).
		SetDescription("Some commands are aliased for shortness.\n").
		AddField("", "*Example*\n"+example).
		AddField("", "").
		AddField("", "*General*\n"+help+about).
		AddField("", "").
		AddField("", "*Administration*\n"+register+unregister).
		SetThumbnail(avatarUrl). // TODO: move out to config .env file
		SetColor(0x9f00d4).SetFooter(version.AppFullName).MessageEmbed

	s.ChannelMessageSendEmbed(m.Message.ChannelID, embedMsg)
}
