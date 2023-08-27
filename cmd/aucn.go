package main

import (
	"flag"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"gkirito.com/autoUpgradeChainNode/pkg/handler"
	"gkirito.com/autoUpgradeChainNode/pkg/listener/discord"
	"gkirito.com/autoUpgradeChainNode/pkg/node_runner"
	"os"
	"os/signal"
)

var (
	ListenType      string
	DiscordBotToken string
	KeyWordRegexp   string
	RunnerType      string
	DCPath          string
	SystemdName     string
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	flag.StringVar(&ListenType, "lt", "discord", "use -lt to set the listen type, default: discord")
	flag.StringVar(&DiscordBotToken, "dt", "", "set discord bot token")
	flag.StringVar(&KeyWordRegexp, "kwr", "", "set keyword regexp, ex. \\bCommit Sha: ([0-9a-fA-F]+)|1")
	flag.StringVar(&RunnerType, "rt", "docker-compose", "set node runner type: docker-compose | systemd")
	flag.StringVar(&DCPath, "dcp", "./docker-compose.yaml", "set docker compose path, defual: ./docker-compose.yaml")
	flag.StringVar(&SystemdName, "sn", "", "set systemd name")
}

func main() {
	flag.Parse()
	runner := node_runner.GetRunner(RunnerType, DCPath)
	switch ListenType {
	case discord.Type:
		l, err := discord.NewDiscord(DiscordBotToken)
		if err != nil {
			log.Error(err)
			return
		}
		l.AddMsgHandler(runner, KeyWordRegexp, handler.MsgHandler[*discordgo.MessageCreate])
		err = l.Start()
		if err != nil {
			log.Error(err)
			return
		}
	}
	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, os.Kill, os.Interrupt)
	select {
	case <-osSignal:
		log.Infoln("exit")
		return
	}
}
