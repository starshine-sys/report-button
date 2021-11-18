package main

import "github.com/diamondburned/arikawa/v3/discord"

var commands = []discord.Command{
	{
		Name:        "report",
		Type:        discord.ChatInputCommand,
		Description: "Report a situation in chat",
		Options: []discord.CommandOption{
			&discord.StringOption{
				OptionName:  "reason",
				Description: "The reason you're reporting this situation",
				Required:    true,
			},
		},
	},
	{
		Name: "Report message",
		Type: discord.MessageCommand,
	},
}
