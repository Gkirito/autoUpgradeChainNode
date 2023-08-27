package node_runner

import (
	"gkirito.com/autoUpgradeChainNode/pkg/node_runner/docker_compose"
	"gkirito.com/autoUpgradeChainNode/pkg/node_runner/systemd"
)

type Runner interface {
	Start() (string, error)
	Upgrade(name string) error
	Stop() (string, error)
}

type Config interface {
	~string
}

func GetRunner[C Config](runnerType string, config C) Runner {
	switch runnerType {
	case docker_compose.DockerCompose:
		return docker_compose.NewDcRunner(string(config))
	case systemd.SystemdRunner:
		return nil
	default:
		panic("unKnown runner type")
	}
}
