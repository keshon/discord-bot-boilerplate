package botsdef

import (
	"github.com/bwmarrin/discordgo"
	"github.com/gookit/slog"
	helloWorld "github.com/keshon/discord-bot-boilerplate/mod-helloworld/discord"
	hiGalaxy "github.com/keshon/discord-bot-boilerplate/mod-higalaxy/discord"
)

var Modules = []string{"hi", "hello"}

// CreateBotInstance creates a new bot instance based on the module name.
//
// Parameters:
// - session: a Discord session
// - module: the name of the module ("hi" or "hello")
// Returns a botsdef.Discord instance.
func CreateBotInstance(session *discordgo.Session, module string) Discord {
	switch module {
	case "hi":
		return hiGalaxy.NewDiscord(session)
	case "hello":
		return helloWorld.NewDiscord(session)

	// ..add more cases for other modules if needed

	default:
		slog.Printf("Unknown module: %s", module)
		return nil
	}
}
