# report button

(not actually a button, it's an application command)

## Setup

Set up an `.env` file with the following keys:

- `TOKEN` and `PUBLIC_KEY`: your bot's token and public key from Discord
- `GUILD_ID`: the guild ID to respond to commands from, and to write commands to.
- `CHANNEL_ID`: the channel ID to send reports to.

Also set `PORT` if you don't want to use the default (3100)

You should also change some of the messages in `consts.go` (especially `reportError`)

Then build the application (use `build.sh` to build without cgo) and run it with the `-sync` flag. This will sync the commands to Discord and immediately exit.  
Then rerun the executable with no arguments to run the report button!

## License

I don't care what you do with this, so it's licensed under the Unlicense :]  
Check LICENSE in the repository's root for the full text!