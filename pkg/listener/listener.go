package listener

import (
	"gkirito.com/autoUpgradeChainNode/pkg/handler"
	"gkirito.com/autoUpgradeChainNode/pkg/node_runner"
)

type Listener[T handler.Message] interface {
	Send(relayMsg T, msg string)
	Start() error
	AddMsgHandler(runner node_runner.Runner, channel, KeyWordRegexp string, handel handler.Handler[T])
}
