package botsdef

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/keshon/discord-bot-boilerplate/example-bot/discord"
)

type Discord struct {
	Session              *discordgo.Session
	GuildID              string
	IsInstanceActive     bool
	CommandPrefix        string
	LastChangeAvatarTime time.Time
	RateLimitDuration    time.Duration
}

type BotInstance struct {
	ExampleBot *Discord
}

func ConvertToBotsDefDiscord(discordInstance *discord.Discord) *Discord {
	return &Discord{
		// Copy the fields from discordInstance to the corresponding fields in botsdef.Discord
		Session:              discordInstance.Session,
		GuildID:              discordInstance.GuildID,
		IsInstanceActive:     discordInstance.IsInstanceActive,
		CommandPrefix:        discordInstance.CommandPrefix,
		LastChangeAvatarTime: discordInstance.LastChangeAvatarTime,
		RateLimitDuration:    discordInstance.RateLimitDuration,
	}
}
