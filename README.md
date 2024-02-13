# Discord Bot Boilerplate

Basic boilerplate for my Discord bots.

## Getting Started

### Adding the Bot to a Discord Server

To add Dicer Roller to your Discord server:

1. Create a bot at the [Discord Developer Portal](https://discord.com/developers/applications) and acquire the Bot's CLIENT_ID.
2. Use the following link: `discord.com/oauth2/authorize?client_id=YOUR_CLIENT_ID_HERE&scope=bot&permissions=36727824`
   - Replace `YOUR_CLIENT_ID_HERE` with your Bot's Client ID from step 1.
3. The Discord authorization page will open in your browser, allowing you to select a server.
4. Choose the server where you want to add Dicer Roller and click "Authorize".
5. If prompted, complete the reCAPTCHA verification.
6. Grant Dicer Roller the necessary permissions for it to function correctly.
7. Click "Authorize" to add Dicer Roller to your server.

### Building Dicer Roller

Dicer Roller is written in Go language, allowing it to run on a *server* or as a *local* program.

**Local Usage**
Follow the provided scripts to build Dicer Roller locally:
  - `bash-and-run.bat` (or `.sh` for Linux): Build the debug version and execute.
  - `build-release.bat` (or `.sh` for Linux): Build the release version.

For local usage, run these scripts for your operating system and rename `.env.example` to `.env`, storing your Discord Bot Token in the `DISCORD_BOT_TOKEN` variable.

**Server Usage**
For Docker deployment, refer to the `deploy/README.md` for specific instructions.

## Where to get support

If you have any questions you can ask me in my [Discord server](https://discord.gg/NVtdTka8ZT) to get support. Bear in mind there is no community whatsoever — just me.

## License

Discord Bot Boilerplate is licensed under the [MIT License](https://opensource.org/licenses/MIT).
