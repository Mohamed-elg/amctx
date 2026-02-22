package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

var amtoolConfigPath string
var amctxConfigPath string

func init() {
	var err error

	amtoolConfigPath, err = getAbsolutePath(".config/amtool/config.yml")
	if err != nil {
		panic(err)
	}

	amctxConfigPath, err = getAbsolutePath(".config/amctx/config.yml")
	if err != nil {
		panic(err)
	}
}

func getAbsolutePath(path string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, path), nil
}

func getOrCreateConfigFilePath(path string) (string, bool, error) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			return "", false, err
		}
		file, err := os.Create(path)
		if err != nil {
			return "", false, err
		}
		err = file.Close()
		if err != nil {
			return "", false, err
		}
		return path, true, nil
	}
	return path, false, nil
}

func getConfigFileContent(filePath string) (map[string]string, bool, error) {
	path, created, err := getOrCreateConfigFilePath(filePath)

	if err != nil {
		return nil, created, err
	}

	ymlBytes, err := os.ReadFile(path)

	if err != nil {
		return nil, created, err
	}

	data := make(map[string]string)
	err = yaml.Unmarshal(ymlBytes, &data)

	if err != nil {
		return nil, created, err
	}
	return data, created, nil
}

func updateKeyValue(filePath string, key string, value string) (bool, error) {
	data, created, err := getConfigFileContent(filePath)

	if err != nil {
		return created, err
	}

	data[key] = value

	ymlBytes, err := yaml.Marshal(&data)
	if err != nil {
		return created, err
	}

	err = os.WriteFile(filePath, ymlBytes, 0644)

	if err != nil {
		return created, err
	}
	return created, nil
}

func setAlertmanagerUrl(url string) (bool, error) {
	created, err := updateKeyValue(amtoolConfigPath, "alertmanager.url", url)

	if err != nil {
		return created, err
	}

	return created, nil
}

func CreateOrUpdateAlertmanagerAlias(alias string, url string) (bool, error) {
	created, err := updateKeyValue(amctxConfigPath, alias, url)

	if err != nil {
		return created, err
	}

	return created, nil
}

func ListAliases() ([]string, bool, error) {
	data, created, err := getConfigFileContent(amctxConfigPath)

	if err != nil {
		return nil, created, err
	}

	aliases := []string{}
	for key := range data {
		aliases = append(aliases, key)
	}
	sort.Strings(aliases)

	return aliases, created, nil
}

func SwitchContext(alias string) (bool, error) {
	data, created, err := getConfigFileContent(amctxConfigPath)

	if err != nil {
		return created, err
	}

	url, ok := data[alias]
	if !ok {
		return created, fmt.Errorf("no alias exists with the name: %s", alias)
	}

	created, err = setAlertmanagerUrl(url)
	if err != nil {
		return created, err
	}

	return created, nil
}
