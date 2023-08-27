package update

import (
	"context"
	"errors"
	"fmt"
	"github.com/compose-spec/compose-go/loader"
	"github.com/compose-spec/compose-go/types"
	"os"
	"strings"
)

func ChangeYamlFile(hash, composeFilePath string) error {
	config, err := loader.LoadWithContext(context.Background(), types.ConfigDetails{
		ConfigFiles: []types.ConfigFile{
			{
				Filename: composeFilePath,
			},
		},
	})
	if err != nil {
		return err
	}

	if len(config.Services) <= 0 {
		return errors.New("no services")
	}

	versionNames := strings.Split(config.Services[0].Image, ":")
	if len(versionNames) != 2 {
		return errors.New("not found image version in docker compose file")
	}

	if versionNames[1] == hash {
		return nil
	}

	config.Services[0].Image = fmt.Sprintf("%s:%s", versionNames[0], hash)

	newYamlData, err := config.MarshalYAML()
	if err != nil {
		return err
	}

	err = os.WriteFile(composeFilePath, newYamlData, 0644)
	if err != nil {
		return err
	}
	return nil
}
