package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/utils/json/option"
	"github.com/go-chi/render"
	"github.com/starshine-sys/report-button/interaction"
)

func handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r := recover()
		if r != nil {
			log.Printf("Recovered from error: %v", r)
		}
	}()

	if !interaction.Verify(r, publicKey) {
		http.Error(w, "signature mismatch", http.StatusUnauthorized)
	}

	defer r.Body.Close()

	var i interaction.Interaction
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		log.Printf("Error unmarshaling body: %v", err)
		return
	}

	if i.Type == discord.PingInteractionType {
		render.JSON(w, r, api.InteractionResponse{
			Type: api.PongInteraction,
		})
		return
	}

	if i.Data == nil {
		log.Printf("Data for interaction %v was nil!", i.ID)
		return
	}

	if i.GuildID != guildID {
		if i.GuildID.IsValid() {
			respond(w, r, "This command doesn't work in DMs.")
		} else {
			respond(w, r, "This command doesn't work in this server.")
		}
		return
	}

	switch i.Data.Name {
	case "report":
		err = report(w, r, i)
	case "Report message":
		err = reportMessage(w, r, i)
	default:
		err = fmt.Errorf("unknown command %q", i.Data.Name)
	}
}

func report(w http.ResponseWriter, r *http.Request, i interaction.Interaction) (err error) {
	msg := "No report reason (this is a bug!)"

	for _, o := range i.Data.Options {
		if o.Name == "reason" {
			msg = o.String()
			break
		}
	}

	u := i.GetUser()
	if u == nil {
		respond(w, r, "No user associated with this command (this is a bug)!")
		return
	}

	content := fmt.Sprintf(
		"New report by %v/%v from %v\nLink to context: <https://discord.com/channels/%v/%v/%v>",
		u.Mention(), u.Tag(), i.ChannelID.Mention(),
		i.GuildID, i.ChannelID, discord.NewSnowflake(time.Now().UTC()),
	)

	_, err = client.SendMessageComplex(channelID, api.SendMessageData{
		Content: content,
		Embeds: []discord.Embed{{
			Title:       "Reason",
			Description: msg,
			Timestamp:   discord.NowTimestamp(),
			Footer: &discord.EmbedFooter{
				Text: "User ID: " + u.ID.String(),
			},
			Color: 0xc1140b,
		}},
		AllowedMentions: &api.AllowedMentions{},
	})
	if err != nil {
		respond(w, r, reportError)
		return err
	}

	respond(w, r, reportOK)
	return nil
}

func reportMessage(w http.ResponseWriter, r *http.Request, i interaction.Interaction) (err error) {
	msg := i.Data.Message()
	if msg == nil {
		respond(w, r, "Your command didn't have a message attached (this is a bug)!")
		return
	}

	reportedMu.Lock()
	v, ok := reported[msg.ID]
	reportedMu.Unlock()
	if ok {
		if v.Add(timeout).After(time.Now()) {
			respond(w, r, timeoutMsg)
			return
		}
	}

	u := i.GetUser()
	if u == nil {
		respond(w, r, "No user associated with this command (this is a bug)!")
		return
	}

	content := fmt.Sprintf(
		"Message in %v reported by %v/%v\nLink to message: <https://discord.com/channels/%v/%v/%v>",
		i.ChannelID.Mention(), u.Mention(), u.Tag(),
		i.GuildID, msg.ChannelID, msg.ID,
	)

	_, err = client.SendMessageComplex(channelID, api.SendMessageData{
		Content: content,
		Embeds: []discord.Embed{{
			Author: &discord.EmbedAuthor{
				Icon: msg.Author.AvatarURLWithType(discord.PNGImage),
				Name: msg.Author.Tag() + " / " + msg.Author.ID.String(),
			},
			Description: msg.Content,
			Timestamp:   msg.Timestamp,
			Footer: &discord.EmbedFooter{
				Text: "Message ID: " + msg.ID.String(),
			},
			Color: 0xc1140b,
		}},
		AllowedMentions: &api.AllowedMentions{},
	})
	if err != nil {
		respond(w, r, reportError)
		return err
	}

	reportedMu.Lock()
	reported[msg.ID] = time.Now()
	reportedMu.Unlock()

	respond(w, r, reportOK)
	return nil
}

func respond(w http.ResponseWriter, r *http.Request, tmpl string, v ...interface{}) {
	render.JSON(w, r, api.InteractionResponse{
		Type: api.MessageInteractionWithSource,
		Data: &api.InteractionResponseData{
			Content: option.NewNullableString(
				fmt.Sprintf(tmpl, v...),
			),
			Flags: api.EphemeralResponse,
		},
	})
}
