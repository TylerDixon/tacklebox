package main

import (
	"github.com/mitchellh/go-homedir"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"log"
	"encoding/json"
)

type ConfigData struct {
	TemplateDir string

}

func main() {
	fmt.Printf("Cool\n")
	configData, err := initialize()
	if err != nil {
		return;
	}
	fmt.Printf("%s", configData.TemplateDir)
}

func initialize() (ConfigData, error) {
	var unmarshaledConfigData ConfigData
	homeDirectory, _ := homedir.Dir()
	mkDirErr := os.Mkdir(path.Join(homeDirectory, ".tacklebox"), 0777)
	if !os.IsExist(mkDirErr) {
		fmt.Printf("Create tacklebox config directory error: %s", mkDirErr)
		return unmarshaledConfigData, mkDirErr
	}
	configFilePath := path.Join(homeDirectory, ".tacklebox", "config.json")

	// Create or retrieve config data
	var encodedConfigData []byte
	var fileCreationOrReadError error
	if _, err := os.Stat(configFilePath); err == nil {
		encodedConfigData, fileCreationOrReadError = ioutil.ReadFile(configFilePath)
	} else {
		encodedConfigData, fileCreationOrReadError = initializeConfigFile(configFilePath)
	}
	if fileCreationOrReadError != nil {
		return unmarshaledConfigData, fileCreationOrReadError;
	}

	unmarshalConfigError := json.Unmarshal(encodedConfigData, &unmarshaledConfigData)

	if unmarshalConfigError != nil {
		fmt.Printf("Failed to unmarshal config data with error %s", unmarshalConfigError)
		return unmarshaledConfigData, unmarshalConfigError
	}
	return unmarshaledConfigData, nil
}

func initializeConfigFile(configFilePath string) ([]byte, error) {
	file, createFileErr := os.Create(configFilePath)
	if createFileErr != nil {
		log.Fatalf("Failed to create config file at %s due to error: %s", configFilePath, createFileErr)
		return nil, createFileErr
	}
	defaultConfigFile := []byte(`{
		"TemplateDir": "./templates"
	}`)
	_, writeFileErr := file.Write(defaultConfigFile)
	if writeFileErr != nil {
		log.Fatalf("Failed to create config file at %s due to error: %s", configFilePath, writeFileErr)
		return nil, writeFileErr
	}
	return defaultConfigFile, nil
}
