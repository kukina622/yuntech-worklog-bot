package bot

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

type discordBot struct {
	session *discordgo.Session
}

var instance *discordBot

func GetDiscordBotInstance() *discordBot {
	if instance == nil {
		instance = &discordBot{}
	}
	return instance
}

func (bot *discordBot) InitSession(token string) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	bot.session = session
	err = bot.session.Open()
	if err != nil {
		panic(err)
	}
}

func (bot *discordBot) SendMessage(message string, channel string) {
	bot.session.ChannelMessageSend(channel, message)
}
