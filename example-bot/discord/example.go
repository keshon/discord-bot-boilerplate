package discord

import (
	"fmt"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

// handleRollCommand handles the roll command for Discord.
func (d *Discord) handleExampleCommand(s *discordgo.Session, m *discordgo.MessageCreate, param string) {
	d.changeAvatar(s)

	helloWorld := "Hello World"

	embedMsg := embed.NewEmbed().
		SetTitle(fmt.Sprintf("%v", helloWorld)).
		SetColor(0x9f00d4)

	s.ChannelMessageSendEmbed(m.ChannelID, embedMsg.MessageEmbed)
}
