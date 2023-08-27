package handler

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gkirito.com/autoUpgradeChainNode/pkg/node_runner"
	"regexp"
	"strconv"
	"strings"
)

type Message interface {
	~*discordgo.MessageCreate
}

type Handler[T Message] func(runner node_runner.Runner, msg T, KeyWordRegexp string, send func(T, string))

func (m Handler[T]) Discord(runner node_runner.Runner, channel, KeyWordRegexp string, send func(T, string)) func(session *discordgo.Session, msg *discordgo.MessageCreate) {
	return func(session *discordgo.Session, msg *discordgo.MessageCreate) {
		if msg.ChannelID != channel {
			return
		}
		log.Infof("receive the message from disocrd channel[%s] msgId[%s] content: \n %s", msg.ChannelID, msg.ID, msg.Content)
		m(runner, msg, KeyWordRegexp, send)
	}
}

func GetKeyword[T Message](msg T, KeyWordRegexp string) (string, error) {
	rp := strings.Split(KeyWordRegexp, "|")
	if len(rp) != 2 {
		return "", errors.New("the keyword regexp must use '|' to split the regexp and result index")
	}
	resultIndex, err := strconv.ParseInt(rp[1], 10, 64)
	if err != nil {
		return "", err
	}
	rep, err := regexp.Compile(rp[0])
	if err != nil {
		return "", err
	}
	var content string
	if message, ok := any(msg).(*discordgo.MessageCreate); ok {
		content = message.Content
	} else {
		return "", errors.New("unsupported message")
	}
	findResult := rep.FindStringSubmatch(content)
	if int64(len(findResult)) < resultIndex {
		return "", fmt.Errorf("the keyword index is range out of %d", resultIndex)
	}
	return findResult[resultIndex], nil
}

func MsgHandler[T Message](runner node_runner.Runner, msg T, KeyWordRegexp string, send func(T, string)) {
	newInfo, err := GetKeyword(msg, KeyWordRegexp)
	if err != nil {
		log.Infof("cannot find keyword")
		return
	}
	log.Infof("get new info: %s", newInfo)

	log.Infoln("start stop the node")
	stopLog, err := runner.Stop()
	if err != nil {
		log.Error(err)
		send(msg, err.Error())
		return
	}
	send(msg, stopLog)
	log.Infoln("node stop success")

	err = runner.Upgrade(newInfo)
	if err != nil {
		log.Error(err)
		send(msg, err.Error())
		return
	}
	log.Infoln("update success")

	startLog, err := runner.Start()
	if err != nil {
		log.Error(err)
		send(msg, err.Error())
		return
	}
	send(msg, startLog)
	log.Infoln("restart success")
}
