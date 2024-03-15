# Discord Bot Template

This is a basic boilerplate that I use as a starting point for my Discord bot projects. It was derived from the [Melodix Player](https://github.com/keshon/discord-bot-template) project.

## Getting Started

### Adding the Bot to a Discord Server

To add the bot to your Discord server:

1. Create an application at the [Discord Developer Portal](https://discord.com/developers/applications) and acquire the CLIENT_ID from OAuth2 section.
2. Use the following link: `discord.com/oauth2/authorize?client_id=YOUR_CLIENT_ID_HERE&scope=bot&permissions=36727824`
   - Replace `YOUR_CLIENT_ID_HERE` with your Bot's Client ID from step 1.
3. The Discord authorization page will open in your browser, allowing you to select a server.
4. Choose the server where you want to add the bot and click "Authorize".
5. If prompted, complete the reCAPTCHA verification.

### Building

### Locally ###

Follow the provided scripts :
  - `bash-and-run.bat` (or `.sh` for Linux): Build the debug version and execute.
  - `build-release.bat` (or `.sh` for Linux): Build the release version.

For local usage, run these scripts for your operating system and rename `.env.example` to `.env`, storing your Discord Bot Token in the `DISCORD_BOT_TOKEN` variable.

### Docker ###

For Docker deployment, refer to the `docker/README.md` for specific instructions.

## Where to get support

If you have any questions you can ask me in my [Discord server](https://discord.gg/NVtdTka8ZT) to get support.

## License

Discord Bot Template is licensed under the [MIT License](https://opensource.org/licenses/MIT).
