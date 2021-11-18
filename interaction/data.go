// Package interaction provides types for handling Discord webhook interactions.
package interaction

import "github.com/diamondburned/arikawa/v3/discord"

// Interaction is an interaction received from Discord.
type Interaction struct {
	ID      discord.InteractionID `json:"id"`
	Token   string                `json:"token"`
	Version uint                  `json:"version"`

	ApplicationID discord.AppID               `json:"application_id"`
	Type          discord.InteractionDataType `json:"type"`

	Data *Data `json:"data"`

	GuildID   discord.GuildID   `json:"guild_id"`
	ChannelID discord.ChannelID `json:"channel_id"`

	Member *discord.Member `json:"member"`
	User   *discord.User   `json:"user"`

	Message *discord.Message `json:"message"`
}

// UserID returns the user ID of this interaction, if any.
func (i Interaction) UserID() discord.UserID {
	if u := i.GetUser(); u != nil {
		return u.ID
	}
	return 0
}

// GetUser gets this interaction's user, if any.
func (i Interaction) GetUser() *discord.User {
	if i.Member != nil {
		return &i.Member.User
	}

	if i.User != nil {
		return i.User
	}

	return nil
}

// Data is the interaction data.
type Data struct {
	// Application commands only
	ID       discord.CommandID                  `json:"id"`
	Name     string                             `json:"name"`
	Type     discord.CommandType                `json:"type"`
	Resolved ResolvedData                       `json:"resolved"`
	Options  []discord.CommandInteractionOption `json:"options"`

	// Components only
	CustomID      string                `json:"custom_id"`
	ComponentType discord.ComponentType `json:"component_type"`
	// Select only
	Values []string `json:"values"`

	// User/message commands only
	TargetID discord.Snowflake `json:"target_id"`
}

// User gets the targeted user, if any.
func (d Data) User() *discord.User {
	if !d.TargetID.IsValid() {
		return nil
	}

	v, ok := d.Resolved.Users[discord.UserID(d.TargetID)]
	if !ok {
		return nil
	}
	return &v
}

// Message gets the targeted message, if any.
func (d Data) Message() *discord.Message {
	if !d.TargetID.IsValid() {
		return nil
	}

	v, ok := d.Resolved.Messages[discord.MessageID(d.TargetID)]
	if !ok {
		return nil
	}
	return &v
}

// ResolvedData is the resolved data for this interaction.
type ResolvedData struct {
	Users    map[discord.UserID]discord.User       `json:"users"`
	Members  map[discord.UserID]discord.Member     `json:"members"`
	Roles    map[discord.RoleID]discord.Role       `json:"roles"`
	Channels map[discord.ChannelID]discord.Channel `json:"channels"`
	Messages map[discord.MessageID]discord.Message `json:"messages"`
}
