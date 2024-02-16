package discord

import (
	"fmt"

	embed "github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

// handleRollCommand handles the roll command for Discord.

// handleExampleCommand handles the example command for the Discord bot.
//
// Parameters:
//
//	s - the Discord session
//	m - the message create event
//	param - the command parameter
//
// Return type: None
func (d *Discord) handleHiCommand(s *discordgo.Session, m *discordgo.MessageCreate, param string) {
	d.changeAvatar(s)

	message := "Hi Galaxy!"

	embed := embed.NewEmbed().
		SetTitle(fmt.Sprintf("%v", message)).
		SetColor(0x9f00d4).
		SetFooter("From mod-higalaxy")

	s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
}
