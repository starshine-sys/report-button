package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"

	_ "github.com/joho/godotenv/autoload"
)

var syncCommands bool

var (
	client    *api.Client
	publicKey ed25519.PublicKey

	reported   map[discord.MessageID]time.Time
	reportedMu sync.Mutex

	guildID   discord.GuildID
	channelID discord.ChannelID
)

func init() {
	reported = map[discord.MessageID]time.Time{}

	flag.BoolVar(&syncCommands, "sync", false, "Sync commands with Discord")
	flag.Parse()
}

func main() {
	key, err := hex.DecodeString(os.Getenv("PUBLIC_KEY"))
	if err != nil {
		log.Fatalf("Couldn't decode %q into a byte slice", os.Getenv("PUBLIC_KEY"))
	}
	publicKey = key

	sf, err := discord.ParseSnowflake(os.Getenv("GUILD_ID"))
	if err != nil {
		log.Fatalf("%q is not a valid snowflake", os.Getenv("GUILD_ID"))
	}
	guildID = discord.GuildID(sf)

	sf, err = discord.ParseSnowflake(os.Getenv("CHANNEL_ID"))
	if err != nil {
		log.Fatalf("%q is not a valid snowflake", os.Getenv("GUILD_ID"))
	}
	channelID = discord.ChannelID(sf)

	client = api.NewClient("Bot " + os.Getenv("TOKEN"))

	if syncCommands {
		log.Printf("Writing commands to Discord")

		me, err := client.CurrentApplication()
		if err != nil {
			log.Fatalf("Error getting application: %v", err)
		}

		_, err = client.BulkOverwriteGuildCommands(me.ID, guildID, commands)
		if err != nil {
			log.Fatalf("Error overwriting commands: %v", err)
		}
		log.Printf("Wrote commands!")
		return
	}

	http.HandleFunc("/", handle)

	log.Printf("Starting server")

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	port = strings.Trim(port, ":")

	log.Printf("Will listen on :%v", port)

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
