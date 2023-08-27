package discord

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gkirito.com/autoUpgradeChainNode/pkg/handler"
	"gkirito.com/autoUpgradeChainNode/pkg/node_runner"
)

const Type = "discord"

type Discord[T handler.Message] struct {
	*discordgo.Session
	logger *log.Entry
}

func NewDiscord(token string) (*Discord[*discordgo.MessageCreate], error) {
	discord, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	return &Discord[*discordgo.MessageCreate]{
		Session: discord,
		logger:  log.WithField("listener", Discord[*discordgo.MessageCreate]{}),
	}, nil
}

func (d *Discord[T]) AddMsgHandler(runner node_runner.Runner, channel, KeyWordRegexp string, handler handler.Handler[T]) {
	d.AddHandler(handler.Discord(runner, channel, KeyWordRegexp, d.Send))
}

func (d *Discord[T]) Start() error {
	return d.Open()
}

func (d *Discord[T]) Send(relayMsg T, msg string) {
	func(relayMsg *discordgo.MessageCreate, msg string) {
		_, err := d.ChannelMessageSendReply(relayMsg.ChannelID, msg, &discordgo.MessageReference{
			MessageID: relayMsg.ID,
			ChannelID: relayMsg.ChannelID,
			GuildID:   relayMsg.GuildID,
		})
		if err != nil {
			d.logger.Error(err)
		}
	}(relayMsg, msg)
}
