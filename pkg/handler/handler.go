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

type Handler[T Message] func(logger *log.Entry, runner node_runner.Runner, msg T, KeyWordRegexp string, send func(T, string))

func (h Handler[T]) Discord(runner node_runner.Runner, channel, KeyWordRegexp string, send func(T, string)) func(session *discordgo.Session, msg *discordgo.MessageCreate) {
	logger := log.WithField("handler", "discord")
	return func(session *discordgo.Session, msg *discordgo.MessageCreate) {
		if msg.ChannelID != channel {
			return
		}
		logger.Infof("receive the message from disocrd channel[%s] msgId[%s] content: \n %s", msg.ChannelID, msg.ID, msg.Content)
		h(logger, runner, msg, KeyWordRegexp, send)
	}
}

func GetMsgInfo[T Message](msg T) (string, string, error) {
	if message, ok := any(msg).(*discordgo.MessageCreate); ok {
		return message.ID, message.Content, nil
	} else {
		return "", "", errors.New("unsupported message")
	}
}

func GetKeyword(msg, KeyWordRegexp string) (string, error) {
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
	findResult := rep.FindStringSubmatch(msg)
	if int64(len(findResult)) < resultIndex {
		return "", fmt.Errorf("the keyword index is range out of %d", resultIndex)
	}
	return findResult[resultIndex], nil
}

func MsgHandler[T Message](logger *log.Entry, runner node_runner.Runner, msg T, KeyWordRegexp string, send func(T, string)) {
	id, content, err := GetMsgInfo(msg)
	if err != nil {
		logger.Error("got msg info err: %v", err)
		return
	}
	msgLogger := logger.WithField("message", id)
	newInfo, err := GetKeyword(content, KeyWordRegexp)
	if err != nil {
		msgLogger.Infof("cannot find keyword")
		return
	}
	msgLogger.Infof("get new info: %s", newInfo)

	msgLogger.Infoln("start stop the node")
	stopLog, err := runner.Stop()
	if err != nil {
		msgLogger.Error(err)
		send(msg, err.Error())
		return
	}
	send(msg, stopLog)
	msgLogger.Infoln("node stop success")

	err = runner.Upgrade(newInfo)
	if err != nil {
		msgLogger.Error(err)
		send(msg, err.Error())
		return
	}
	msgLogger.Infoln("update success")

	startLog, err := runner.Start()
	if err != nil {
		msgLogger.Error(err)
		send(msg, err.Error())
		return
	}
	send(msg, startLog)
	msgLogger.Infoln("restart success")
}
