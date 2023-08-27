package docker_compose

import (
	"bytes"
	"errors"
	"gkirito.com/autoUpgradeChainNode/pkg/update"
	"os/exec"
)

const DockerCompose = "docker-compose"

type DcRunner struct {
	composeFilePath string
}

func NewDcRunner(composeFilePath string) *DcRunner {
	return &DcRunner{
		composeFilePath: composeFilePath,
	}
}

func cmd(name string, args ...string) (string, error) {
	stop := exec.Command(name, args...)
	cmdLog := bytes.NewBuffer([]byte{})
	cmdErrLog := bytes.NewBuffer([]byte{})
	stop.Stdout = cmdLog
	stop.Stderr = cmdErrLog
	err := stop.Run()
	if err != nil {
		return "", err
	}
	if cmdErrLog.Len() > 0 {
		return "", errors.New(cmdErrLog.String())
	}
	return cmdLog.String(), nil
}

func (d DcRunner) Start() (string, error) {
	return cmd("docker compose", "-f", d.composeFilePath, "up", "-d")
}

func (d DcRunner) Stop() (string, error) {
	return cmd("docker compose", "-f", d.composeFilePath, "down")
}

func (d DcRunner) Upgrade(version string) error {
	return update.ChangeYamlFile(version, d.composeFilePath)
}