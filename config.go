package main

import (
	"os"
	"log"
	"github.com/mitchellh/go-homedir"
	"path"
	"fmt"
	"io/ioutil"
	"encoding/json"
)

type ConfigData struct {
	Projects []ProjectConfig
	Templates []Template
}

type ProjectConfig struct {
	Name string
	Location string
	TemplateSettings []TemplateSetting
}

type TemplateSetting struct {
	Name string
	Settings map[string]interface{}
}

// Retrieve config returns a parsed config data from the tacklebox config directory in the home directory.
func RetrieveConfig() (ConfigData, error) {
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
		encodedConfigData, fileCreationOrReadError = InitializeConfigFile(configFilePath)
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

// InitializeConfigFile creates a default configuration file for use by tacklebox.
func InitializeConfigFile(configFilePath string) ([]byte, error) {
	file, createFileErr := os.Create(configFilePath)
	if createFileErr != nil {
		log.Fatalf("Failed to create config file at %s due to error: %s", configFilePath, createFileErr)
		return nil, createFileErr
	}
	defaultConfigFile := []byte(`{
		"Templates": []
	}`)
	_, writeFileErr := file.Write(defaultConfigFile)
	if writeFileErr != nil {
		log.Fatalf("Failed to create config file at %s due to error: %s", configFilePath, writeFileErr)
		return nil, writeFileErr
	}
	return defaultConfigFile, nil
}